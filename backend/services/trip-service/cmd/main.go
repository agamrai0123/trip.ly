package main

import (
	"context"
	"fmt"
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
	proto "github.com/agamrai0123/wanderplan/proto/gen/wanderplan/v1"
	internal "github.com/agamrai0123/wanderplan/services/trip-service/internal"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
)

func main() {
	var cfg internal.Config
	if err := pkgcfg.Load("trip-service-config", "./config", "TRIP", &cfg); err != nil {
		fmt.Fprintf(os.Stderr, "load config: %v\n", err)
	}
	if err := cfg.Validate(); err != nil {
		fmt.Fprintf(os.Stderr, "config validation: %v\n", err)
		os.Exit(1)
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

	internal.InitLogger(cfg.Logging, "trip-service")

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

	producer, err := kafka.NewProducer(cfg.Kafka.Brokers)
	if err != nil {
		log.Warn().Err(err).Msg("kafka unavailable — events will be dropped")
	}

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

	tripRepo := internal.NewTripRepo(pool)
	dayRepo := internal.NewDayRepo(pool)
	itemRepo := internal.NewItemRepo(pool)
	svc := internal.NewTripService(tripRepo, dayRepo, itemRepo, producer)
	handlers := internal.NewHandlers(svc)
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
		log.Info().Str("addr", httpSrv.Addr).Msg("trip-service HTTP listening")
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
	proto.RegisterTripServiceServer(grpcSrv, &tripGRPCServer{svc: svc})
	go func() {
		log.Info().Str("addr", grpcAddr).Msg("trip-service gRPC listening")
		if err := grpcSrv.Serve(lis); err != nil {
			log.Fatal().Err(err).Msg("gRPC server error")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Info().Msg("shutting down trip-service...")

	shutCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	grpcSrv.GracefulStop()
	if err := httpSrv.Shutdown(shutCtx); err != nil {
		log.Error().Err(err).Msg("HTTP shutdown error")
	}
	if producer != nil {
		_ = producer.Close()
	}
	log.Info().Msg("trip-service stopped")
}

// tripGRPCServer implements proto.TripServiceServer.
type tripGRPCServer struct {
	proto.UnimplementedTripServiceServer
	svc *internal.TripService
}

// ListTripsByUser returns trip summaries for the given user.
func (g *tripGRPCServer) ListTripsByUser(ctx context.Context, req *proto.ListTripsByUserRequest) (*proto.ListTripsByUserResponse, error) {
	trips, err := g.svc.ListTrips(ctx, req.UserId)
	if err != nil {
		return nil, err
	}
	var summaries []*proto.TripSummary
	for _, t := range trips {
		s := &proto.TripSummary{
			Id:          t.ID,
			Title:       t.Title,
			Destination: t.Destination,
			Status:      t.Status,
		}
		if t.CoverImageURL != "" {
			s.CoverImageUrl = t.CoverImageURL
		}
		if t.StartDate != nil {
			s.StartDate = t.StartDate.Format("2006-01-02")
		}
		if t.EndDate != nil {
			s.EndDate = t.EndDate.Format("2006-01-02")
		}
		summaries = append(summaries, s)
	}
	return &proto.ListTripsByUserResponse{Trips: summaries}, nil
}

// GetTripStats returns aggregate statistics for the user's trips.
func (g *tripGRPCServer) GetTripStats(ctx context.Context, req *proto.GetTripStatsRequest) (*proto.GetTripStatsResponse, error) {
	stats, err := g.svc.GetStats(ctx, req.UserId)
	if err != nil {
		return nil, err
	}
	return &proto.GetTripStatsResponse{
		TotalTrips:     int32(stats.TotalTrips),
		TotalCountries: int32(stats.TotalCountries),
		TotalDays:      int32(stats.TotalDays),
		TotalBudget:    stats.TotalBudget,
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
}
