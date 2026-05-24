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
	"github.com/agamrai0123/wanderplan/pkg/database"
	"github.com/agamrai0123/wanderplan/pkg/jwt"
	"github.com/agamrai0123/wanderplan/pkg/kafka"
	internal "github.com/agamrai0123/wanderplan/services/notification-service/internal"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func main() {
	var cfg internal.Config
	if err := pkgcfg.Load("notification-service-config", "./config", "NOTIFY", &cfg); err != nil {
		fmt.Fprintf(os.Stderr, "load config: %v\n", err)
	}
	if err := cfg.Validate(); err != nil {
		fmt.Fprintf(os.Stderr, "config validation: %v\n", err)
		os.Exit(1)
	}

	internal.InitLogger(cfg.Logging, "notification-service")

	dbCfg := database.Config{
		Host:     cfg.Database.Host,
		Port:     cfg.Database.Port,
		DBName:   cfg.Database.Name,
		User:     cfg.Database.User,
		Password: cfg.Database.Password,
		Schema:   cfg.Database.Schema,
		MaxConns: cfg.Database.MaxConns,
		MinConns: cfg.Database.MinConns,
	}
	pool, err := database.NewPool(context.Background(), dbCfg)
	if err != nil {
		log.Fatal().Err(err).Msg("connect to database")
	}
	defer pool.Close()

	privB64 := os.Getenv("JWT_PRIVATE_KEY")
	pubB64 := os.Getenv("JWT_PUBLIC_KEY")
	var jwtMgr *jwt.Manager
	if privB64 != "" && pubB64 != "" {
		jwtMgr, err = jwt.NewManagerFromBase64(privB64, pubB64, 15*time.Minute)
		if err != nil {
			log.Fatal().Err(err).Msg("init jwt manager")
		}
	} else {
		log.Warn().Msg("JWT keys not set — Auth middleware disabled")
	}

	repo := internal.NewNotificationRepo(pool)
	hub := internal.NewHub()
	svc := internal.NewNotificationService(repo)
	handlers := internal.NewHandlers(svc, hub)
	reg := internal.NewRegistry()

	// Start Kafka consumer if brokers are configured.
	var consumer *internal.EventConsumer
	if len(cfg.Kafka.Brokers) > 0 {
		topics := []string{
			kafka.TopicAuthEvents,
			kafka.TopicTripEvents,
			kafka.TopicCollabEvents,
		}
		consumer, err = internal.NewEventConsumer(cfg.Kafka.Brokers, cfg.Kafka.ConsumerGroupID, topics, repo, hub)
		if err != nil {
			log.Warn().Err(err).Msg("kafka consumer init failed — notifications from events disabled")
		} else {
			consumer.Start()
			log.Info().Strs("topics", topics).Msg("kafka consumer started")
		}
	}

	if os.Getenv("ENV") == "production" {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.New()
	internal.RegisterRoutes(r, handlers, jwtMgr, reg, &cfg)

	httpSrv := &http.Server{
		Addr:         ":" + cfg.ServerPort,
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}
	go func() {
		log.Info().Str("addr", httpSrv.Addr).Msg("notification-service HTTP listening")
		if err := httpSrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal().Err(err).Msg("HTTP server error")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Info().Msg("shutting down notification-service...")

	shutCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := httpSrv.Shutdown(shutCtx); err != nil {
		log.Error().Err(err).Msg("HTTP shutdown error")
	}
	if consumer != nil {
		_ = consumer.Close()
	}
	log.Info().Msg("notification-service stopped")
}
