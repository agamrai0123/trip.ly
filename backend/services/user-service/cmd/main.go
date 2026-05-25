package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	pkgcfg "github.com/agamrai0123/wanderplan/pkg/config"
	"github.com/agamrai0123/wanderplan/pkg/database"
	"github.com/agamrai0123/wanderplan/pkg/jwt"
	proto "github.com/agamrai0123/wanderplan/proto/gen/wanderplan/v1"
	internal "github.com/agamrai0123/wanderplan/services/user-service/internal"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
)

func main() {
	var cfg internal.Config
	if err := pkgcfg.Load("user-service-config", "./config", "USER", &cfg); err != nil {
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
	internal.InitLogger(cfg.Logging, "user-service")

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

	tripSvcAddr := cfg.TripServiceAddr
	if addr := os.Getenv("TRIP_SERVICE_ADDR"); addr != "" {
		tripSvcAddr = addr
	}
	if tripSvcAddr == "" {
		tripSvcAddr = "localhost:9082"
	}

	userRepo := internal.NewUserRepo(pool)
	svc := internal.NewUserService(userRepo, tripSvcAddr)
	handlers := internal.NewHandlers(svc)
	reg := internal.NewRegistry()
	reg.MustRegister(database.NewPoolCollector(pool, "user-service"))

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
		log.Info().Str("addr", httpSrv.Addr).Msg("user-service HTTP listening")
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
	proto.RegisterUserServiceServer(grpcSrv, &userGRPCServer{svc: svc})
	go func() {
		log.Info().Str("addr", grpcAddr).Msg("user-service gRPC listening")
		if err := grpcSrv.Serve(lis); err != nil {
			log.Fatal().Err(err).Msg("gRPC server error")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Info().Msg("shutting down user-service...")

	shutCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	grpcSrv.GracefulStop()
	if err := httpSrv.Shutdown(shutCtx); err != nil {
		log.Error().Err(err).Msg("HTTP shutdown error")
	}
	log.Info().Msg("user-service stopped")
}

// userGRPCServer implements proto.UserServiceServer.
type userGRPCServer struct {
	proto.UnimplementedUserServiceServer
	svc *internal.UserService
}

// GetUser returns a user profile for internal gRPC callers.
func (g *userGRPCServer) GetUser(ctx context.Context, req *proto.GetUserRequest) (*proto.GetUserResponse, error) {
	u, err := g.svc.GetMe(ctx, req.UserId)
	if err != nil {
		return nil, err
	}
	return &proto.GetUserResponse{
		Id:        u.ID,
		Email:     u.Email,
		Name:      u.Name,
		AvatarUrl: u.AvatarURL,
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
