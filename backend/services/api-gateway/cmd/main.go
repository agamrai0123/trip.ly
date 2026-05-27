package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	pkgcfg "github.com/agamrai0123/wanderplan/pkg/config"
	pkgjwt "github.com/agamrai0123/wanderplan/pkg/jwt"
	internal "github.com/agamrai0123/wanderplan/services/api-gateway/internal"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func main() {
	// ── Config ────────────────────────────────────────────────
	var cfg internal.Config
	if err := pkgcfg.Load("api-gateway-config", "./config", "GW", &cfg); err != nil {
		fmt.Fprintf(os.Stderr, "load config: %v\n", err)
	}
	if err := cfg.Validate(); err != nil {
		fmt.Fprintf(os.Stderr, "config validation: %v\n", err)
		os.Exit(1)
	}

	// Override service addresses from Render/deployment env vars.
	if v := os.Getenv("AUTH_SERVICE_ADDR"); v != "" {
		cfg.Services.AuthAddr = v
	}
	if v := os.Getenv("AUTH_SERVICE_GRPC_ADDR"); v != "" {
		cfg.Services.AuthGRPCAddr = v
	}
	if v := os.Getenv("TRIP_SERVICE_ADDR"); v != "" {
		cfg.Services.TripAddr = v
	}
	if v := os.Getenv("USER_SERVICE_ADDR"); v != "" {
		cfg.Services.UserAddr = v
	}
	if v := os.Getenv("COLLAB_SERVICE_ADDR"); v != "" {
		cfg.Services.CollaborationAddr = v
	}
	if v := os.Getenv("NOTIFY_SERVICE_ADDR"); v != "" {
		cfg.Services.NotificationAddr = v
	}
	if v := os.Getenv("SEARCH_SERVICE_ADDR"); v != "" {
		cfg.Services.SearchAddr = v
	}

	// ── Logger ────────────────────────────────────────────────
	internal.InitLogger(cfg.Logging, "api-gateway")

	// ── JWT Manager (for Auth middleware on protected routes) ─
	privB64 := os.Getenv("JWT_PRIVATE_KEY")
	pubB64  := os.Getenv("JWT_PUBLIC_KEY")
	var jwtMgr *pkgjwt.Manager
	if privB64 != "" && pubB64 != "" {
		var err error
		jwtMgr, err = pkgjwt.NewManagerFromBase64(privB64, pubB64, 15*time.Minute)
		if err != nil {
			log.Fatal().Err(err).Msg("init jwt manager")
		}
	} else {
		log.Warn().Msg("JWT_PRIVATE_KEY / JWT_PUBLIC_KEY not set — Auth middleware disabled")
	}

	// ── Auth-service gRPC validator ───────────────────────────
	authValidator, err := internal.NewAuthValidator(cfg.Services.AuthGRPCAddr)
	if err != nil {
		log.Fatal().Err(err).Msg("dial auth-service")
	}
	defer authValidator.Close()

	// ── HTTP server ───────────────────────────────────────────
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()

	reg := internal.NewRegistry()
	internal.RegisterRoutes(router, internal.NewHandlers(authValidator, cfg.Services), jwtMgr, reg, &cfg)

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.ServerPort),
		Handler:      router,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 60 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	go func() {
		log.Info().Int("port", cfg.ServerPort).Msg("api-gateway HTTP listening")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal().Err(err).Msg("http server")
		}
	}()

	// ── Graceful shutdown ─────────────────────────────────────
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info().Msg("shutting down api-gateway…")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Error().Err(err).Msg("graceful shutdown failed")
	}
	log.Info().Msg("api-gateway stopped")
}
