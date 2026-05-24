package main

import (
	"context"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	pkgcfg "github.com/agamrai0123/wanderplan/pkg/config"
	"github.com/agamrai0123/wanderplan/pkg/database"
	"github.com/agamrai0123/wanderplan/pkg/jwt"
	"github.com/agamrai0123/wanderplan/pkg/kafka"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"

	internal "github.com/agamrai0123/wanderplan/services/auth-service/internal"
	proto "github.com/agamrai0123/wanderplan/proto/gen/wanderplan/v1"

	"google.golang.org/grpc"
)

func main() {
	var cfg internal.Config
	if err := pkgcfg.Load("auth-server-config", "./config", "WANDERPLAN", &cfg); err != nil {
		log.Fatal().Err(err).Msg("failed to load config")
	}
	if err := cfg.Validate(); err != nil {
		log.Fatal().Err(err).Msg("invalid config")
	}

	// Override config from flat env vars set by Render (or any deployment env).
	if v := os.Getenv("DB_HOST"); v != "" {
		cfg.Database.Host = v
	}
	if v := os.Getenv("DB_PORT"); v != "" {
		if n := atoi(v); n > 0 {
			cfg.Database.Port = n
		}
	}
	if v := os.Getenv("DB_NAME"); v != "" {
		cfg.Database.Name = v
	}
	if v := os.Getenv("DB_USER"); v != "" {
		cfg.Database.User = v
	}
	if v := os.Getenv("DB_PASSWORD"); v != "" {
		cfg.Database.Password = v
	}
	if v := os.Getenv("KAFKA_BROKERS"); v != "" {
		cfg.Kafka.Brokers = strings.Split(v, ",")
	}
	if v := os.Getenv("GOOGLE_CLIENT_ID"); v != "" {
		cfg.OAuth.Google.ClientID = v
	}
	if v := os.Getenv("GOOGLE_CLIENT_SECRET"); v != "" {
		cfg.OAuth.Google.ClientSecret = v
	}
	if v := os.Getenv("GITHUB_CLIENT_ID"); v != "" {
		cfg.OAuth.GitHub.ClientID = v
	}
	if v := os.Getenv("GITHUB_CLIENT_SECRET"); v != "" {
		cfg.OAuth.GitHub.ClientSecret = v
	}

	internal.InitLogger(cfg.Logging, "auth-service")

	ctx := context.Background()
	pool, err := database.NewPool(ctx, cfg.Database.ToDBConfig())
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect to database")
	}
	defer pool.Close()

	privKey := os.Getenv("JWT_PRIVATE_KEY")
	pubKey := os.Getenv("JWT_PUBLIC_KEY")
	if privKey == "" || pubKey == "" {
		log.Fatal().Msg("JWT_PRIVATE_KEY and JWT_PUBLIC_KEY env vars are required")
	}
	jwtMgr, err := jwt.NewManagerFromBase64(privKey, pubKey, 15*time.Minute)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create JWT manager")
	}

	producer, err := kafka.NewProducer(cfg.Kafka.Brokers)
	if err != nil {
		log.Warn().Err(err).Msg("kafka unavailable — events will be dropped")
	}

	userRepo := internal.NewUserRepo(pool)
	rtRepo := internal.NewRefreshTokenRepo(pool)
	svc := internal.NewAuthService(&cfg, userRepo, rtRepo, jwtMgr, producer)
	handlers := internal.NewHandlers(svc, &cfg)
	reg := internal.NewRegistry()

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
		log.Info().Str("addr", httpSrv.Addr).Msg("auth-service HTTP listening")
		if err := httpSrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal().Err(err).Msg("HTTP server error")
		}
	}()

	grpcAddr := ":" + cfg.GRPCPort
	lis, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		log.Fatal().Err(err).Str("addr", grpcAddr).Msg("failed to bind gRPC port")
	}
	grpcSrv := grpc.NewServer()
	proto.RegisterAuthServiceServer(grpcSrv, &authGRPCServer{svc: svc})
	go func() {
		log.Info().Str("addr", grpcAddr).Msg("auth-service gRPC listening")
		if err := grpcSrv.Serve(lis); err != nil {
			log.Fatal().Err(err).Msg("gRPC server error")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Info().Msg("shutting down auth-service...")

	shutCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	grpcSrv.GracefulStop()
	if err := httpSrv.Shutdown(shutCtx); err != nil {
		log.Error().Err(err).Msg("HTTP shutdown error")
	}
	if producer != nil {
		_ = producer.Close()
	}
	log.Info().Msg("auth-service stopped")
}

type authGRPCServer struct {
	proto.UnimplementedAuthServiceServer
	svc *internal.AuthService
}

func (s *authGRPCServer) ValidateToken(ctx context.Context, req *proto.ValidateTokenRequest) (*proto.ValidateTokenResponse, error) {
	claims, err := s.svc.ValidateToken(req.Token)
	if err != nil {
		return nil, err
	}
	return &proto.ValidateTokenResponse{
		UserId:    claims.UserID,
		Email:     claims.Email,
		Name:      claims.Name,
		AvatarUrl: claims.AvatarURL,
	}, nil
}

// atoi converts a string to int, returning 0 on error.
func atoi(s string) int {
	n := 0
	for _, c := range s {
		if c < '0' || c > '9' {
			return 0
		}
		n = n*10 + int(c-'0')
	}
	return n
