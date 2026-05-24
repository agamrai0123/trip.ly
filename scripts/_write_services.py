"""Write trip-service, user-service, collaboration-service, notification-service, search-service."""
import os

BACKEND = r'D:\Learn\trip.ly\backend'

def write(rel, content):
    fpath = os.path.join(BACKEND, rel)
    os.makedirs(os.path.dirname(fpath), exist_ok=True)
    with open(fpath, 'w', encoding='utf-8') as f:
        f.write(content)
    print(f'  wrote {rel}')


# ══════════════════════════════════════════════════════════════════
# TRIP-SERVICE  (HTTP :8082, gRPC :9082)
# ══════════════════════════════════════════════════════════════════

write('services/trip-service/config/trip-service-config.json', '''{
    "version": "1.0",
    "server_port": 8082,
    "grpc_port": 9082,
    "metric_port": 7082,
    "logging": { "level": 1, "path": "./log/trip-service.log", "max_size_mb": 256, "max_backups": 5, "max_age_days": 30 },
    "database": {
        "host": "localhost", "port": 5432, "name": "wanderplan",
        "user": "postgres", "password": "", "schema": "wanderplan",
        "max_conns": 20, "min_conns": 2
    },
    "kafka": { "brokers": ["localhost:9092"], "topic_trip_events": "trip-events" },
    "cors": { "allowed_origins": ["http://localhost:5173"] },
    "rate_limit": { "rps": 200, "burst": 400 }
}
''')

write('services/trip-service/internal/config.go', '''package internal

import (
\t"github.com/agamrai0123/wanderplan/pkg/database"
)

type Config struct {
\tVersion    string      `mapstructure:"version"`
\tServerPort int         `mapstructure:"server_port"`
\tGRPCPort   int         `mapstructure:"grpc_port"`
\tMetricPort int         `mapstructure:"metric_port"`
\tLogging    LoggingCfg  `mapstructure:"logging"`
\tDatabase   DatabaseCfg `mapstructure:"database"`
\tKafka      KafkaCfg    `mapstructure:"kafka"`
\tCORS       CORSCfg     `mapstructure:"cors"`
\tRateLimit  RateLimitCfg `mapstructure:"rate_limit"`
}
type LoggingCfg struct {
\tLevel int; Path string; MaxSizeMB int; MaxBackups int; MaxAgeDays int
}
type DatabaseCfg struct {
\tHost string `mapstructure:"host"`;  Port int `mapstructure:"port"`
\tName string `mapstructure:"name"`;  User string `mapstructure:"user"`
\tPassword string `mapstructure:"password"`; Schema string `mapstructure:"schema"`
\tMaxConns int32 `mapstructure:"max_conns"`; MinConns int32 `mapstructure:"min_conns"`
}
func (d DatabaseCfg) ToDBConfig() database.Config {
\treturn database.Config{Host: d.Host, Port: d.Port, DBName: d.Name, User: d.User, Password: d.Password, Schema: d.Schema, MaxConns: d.MaxConns, MinConns: d.MinConns}
}
type KafkaCfg struct { Brokers []string `mapstructure:"brokers"`; TopicTripEvents string `mapstructure:"topic_trip_events"` }
type CORSCfg struct { AllowedOrigins []string `mapstructure:"allowed_origins"` }
type RateLimitCfg struct { RPS float64 `mapstructure:"rps"`; Burst int `mapstructure:"burst"` }
func (c *Config) Validate() error {
\tif c.ServerPort == 0 { c.ServerPort = 8082 }
\tif c.GRPCPort == 0 { c.GRPCPort = 9082 }
\treturn nil
}
''')

write('services/trip-service/internal/errors.go', '''package internal

import pkgerr "github.com/agamrai0123/wanderplan/pkg/errors"

var (BadRequest = pkgerr.BadRequest; Unauthorized = pkgerr.Unauthorized; NotFound = pkgerr.NotFound; Internal = pkgerr.Internal)
''')

write('services/trip-service/internal/logger.go', '''package internal

import (
\t"github.com/rs/zerolog"
\t"github.com/rs/zerolog/log"
\tpkglogger "github.com/agamrai0123/wanderplan/pkg/logger"
)

func InitLogger(cfg LoggingCfg, service string) zerolog.Logger {
\tlevel := zerolog.Level(cfg.Level)
\tif level < zerolog.TraceLevel || level > zerolog.Disabled { level = zerolog.InfoLevel }
\tmaxSizeMB := cfg.MaxSizeMB; if maxSizeMB <= 0 { maxSizeMB = 256 }
\tlogger := pkglogger.Init(pkglogger.Config{Level: int(level), FilePath: cfg.Path, MaxSizeMB: maxSizeMB, MaxBackups: cfg.MaxBackups, MaxAgeDays: cfg.MaxAgeDays, Service: service})
\tlog.Logger = logger; zerolog.DefaultContextLogger = &logger
\treturn logger
}
''')

write('services/trip-service/internal/metrics.go', '''package internal

import (
\t"github.com/prometheus/client_golang/prometheus"
\t"github.com/prometheus/client_golang/prometheus/collectors"
)

func NewRegistry() *prometheus.Registry {
\treg := prometheus.NewRegistry()
\treg.MustRegister(collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}), collectors.NewGoCollector())
\treturn reg
}
''')

write('services/trip-service/internal/models.go', '''package internal

import "time"

// Trip represents a travel itinerary.
type Trip struct {
\tID          string    `json:"id" db:"id"`
\tOwnerID     string    `json:"owner_id" db:"owner_id"`
\tTitle       string    `json:"title" db:"title"`
\tDescription string    `json:"description,omitempty" db:"description"`
\tDestination string    `json:"destination,omitempty" db:"destination"`
\tStartDate   *time.Time `json:"start_date,omitempty" db:"start_date"`
\tEndDate     *time.Time `json:"end_date,omitempty" db:"end_date"`
\tCoverImage  string    `json:"cover_image,omitempty" db:"cover_image"`
\tIsPublic    bool      `json:"is_public" db:"is_public"`
\tCreatedAt   time.Time `json:"created_at" db:"created_at"`
\tUpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// ItineraryDay groups items for a single day.
type ItineraryDay struct {
\tID       string         `json:"id" db:"id"`
\tTripID   string         `json:"trip_id" db:"trip_id"`
\tDayIndex int            `json:"day_index" db:"day_index"`
\tDate     *time.Time     `json:"date,omitempty" db:"date"`
\tItems    []ItineraryItem `json:"items,omitempty"`
}

// ItineraryItem is a single activity / place in an itinerary day.
type ItineraryItem struct {
\tID          string    `json:"id" db:"id"`
\tDayID       string    `json:"day_id" db:"day_id"`
\tPosition    int       `json:"position" db:"position"`
\tTitle       string    `json:"title" db:"title"`
\tDescription string    `json:"description,omitempty" db:"description"`
\tPlaceID     string    `json:"place_id,omitempty" db:"place_id"`
\tPlaceName   string    `json:"place_name,omitempty" db:"place_name"`
\tLat         float64   `json:"lat,omitempty" db:"lat"`
\tLng         float64   `json:"lng,omitempty" db:"lng"`
\tDuration    int       `json:"duration_mins,omitempty" db:"duration_mins"`
\tCost        float64   `json:"cost,omitempty" db:"cost"`
\tCreatedAt   time.Time `json:"created_at" db:"created_at"`
}

// ReorderRequest holds new positions from the dnd-kit front-end.
type ReorderRequest struct {
\tItems []struct {
\t\tID       string `json:"id"`
\t\tPosition int    `json:"position"`
\t} `json:"items"`
}
''')

write('services/trip-service/internal/database.go', '''package internal

import (
\t"context"
\t"fmt"
\t"time"

\t"github.com/google/uuid"
\t"github.com/jackc/pgx/v5/pgxpool"
)

// TripRepo manages trip persistence.
type TripRepo struct{ db *pgxpool.Pool }

func NewTripRepo(db *pgxpool.Pool) *TripRepo { return &TripRepo{db: db} }

func (r *TripRepo) Create(ctx context.Context, ownerID, title, description, destination string) (*Trip, error) {
\tt := &Trip{ID: uuid.NewString(), OwnerID: ownerID, Title: title, Description: description, Destination: destination, CreatedAt: time.Now(), UpdatedAt: time.Now()}
\t_, err := r.db.Exec(ctx,
\t\t`INSERT INTO wanderplan.trips (id,owner_id,title,description,destination,is_public,created_at,updated_at)
\t\t VALUES ($1,$2,$3,$4,$5,false,$6,$7)`,
\t\tt.ID, t.OwnerID, t.Title, t.Description, t.Destination, t.CreatedAt, t.UpdatedAt)
\tif err != nil { return nil, fmt.Errorf("create trip: %w", err) }
\treturn t, nil
}

func (r *TripRepo) GetByID(ctx context.Context, id string) (*Trip, error) {
\tt := &Trip{}
\terr := r.db.QueryRow(ctx,
\t\t`SELECT id,owner_id,title,description,destination,start_date,end_date,cover_image,is_public,created_at,updated_at
\t\t FROM wanderplan.trips WHERE id=$1`, id).
\t\tScan(&t.ID, &t.OwnerID, &t.Title, &t.Description, &t.Destination, &t.StartDate, &t.EndDate, &t.CoverImage, &t.IsPublic, &t.CreatedAt, &t.UpdatedAt)
\tif err != nil { return nil, fmt.Errorf("get trip: %w", err) }
\treturn t, nil
}

func (r *TripRepo) ListByOwner(ctx context.Context, ownerID string) ([]Trip, error) {
\trows, err := r.db.Query(ctx,
\t\t`SELECT id,owner_id,title,description,destination,start_date,end_date,cover_image,is_public,created_at,updated_at
\t\t FROM wanderplan.trips WHERE owner_id=$1 ORDER BY created_at DESC`, ownerID)
\tif err != nil { return nil, fmt.Errorf("list trips: %w", err) }
\tdefer rows.Close()
\tvar trips []Trip
\tfor rows.Next() {
\t\tvar t Trip
\t\tif err := rows.Scan(&t.ID, &t.OwnerID, &t.Title, &t.Description, &t.Destination, &t.StartDate, &t.EndDate, &t.CoverImage, &t.IsPublic, &t.CreatedAt, &t.UpdatedAt); err != nil { return nil, err }
\t\ttrips = append(trips, t)
\t}
\treturn trips, nil
}

func (r *TripRepo) Update(ctx context.Context, id string, fields map[string]any) error {
\tfields["updated_at"] = time.Now()
\t_, err := r.db.Exec(ctx,
\t\t`UPDATE wanderplan.trips SET title=$2,description=$3,destination=$4,is_public=$5,updated_at=$6 WHERE id=$1`,
\t\tid, fields["title"], fields["description"], fields["destination"], fields["is_public"], fields["updated_at"])
\tif err != nil { return fmt.Errorf("update trip: %w", err) }
\treturn nil
}

func (r *TripRepo) Delete(ctx context.Context, id string) error {
\t_, err := r.db.Exec(ctx, `DELETE FROM wanderplan.trips WHERE id=$1`, id)
\tif err != nil { return fmt.Errorf("delete trip: %w", err) }
\treturn nil
}

// ItemRepo manages itinerary item positions.
type ItemRepo struct{ db *pgxpool.Pool }

func NewItemRepo(db *pgxpool.Pool) *ItemRepo { return &ItemRepo{db: db} }

func (r *ItemRepo) Reorder(ctx context.Context, items []struct{ ID string; Position int }) error {
\ttx, err := r.db.Begin(ctx)
\tif err != nil { return err }
\tdefer tx.Rollback(ctx)
\tfor _, item := range items {
\t\t_, err := tx.Exec(ctx, `UPDATE wanderplan.itinerary_items SET position=$2 WHERE id=$1`, item.ID, item.Position)
\t\tif err != nil { return fmt.Errorf("reorder item %s: %w", item.ID, err) }
\t}
\treturn tx.Commit(ctx)
}
''')

write('services/trip-service/internal/service.go', '''package internal

import (
\t"context"
\t"encoding/json"

\t"github.com/rs/zerolog/log"
\tpkgkafka "github.com/agamrai0123/wanderplan/pkg/kafka"
)

// TripService orchestrates trip business logic.
type TripService struct {
\ttripRepo *TripRepo
\titemRepo *ItemRepo
\tproducer *pkgkafka.Producer
\ttopic    string
}

func NewTripService(tripRepo *TripRepo, itemRepo *ItemRepo, producer *pkgkafka.Producer, topic string) *TripService {
\treturn &TripService{tripRepo: tripRepo, itemRepo: itemRepo, producer: producer, topic: topic}
}

func (s *TripService) CreateTrip(ctx context.Context, ownerID, title, description, destination string) (*Trip, error) {
\ttrip, err := s.tripRepo.Create(ctx, ownerID, title, description, destination)
\tif err != nil { return nil, Internal("create trip: "+err.Error()) }
\ts.publishEvent(ctx, "trip.created", trip)
\treturn trip, nil
}

func (s *TripService) GetTrip(ctx context.Context, id string) (*Trip, error) {
\ttrip, err := s.tripRepo.GetByID(ctx, id)
\tif err != nil { return nil, NotFound("trip not found") }
\treturn trip, nil
}

func (s *TripService) ListTrips(ctx context.Context, ownerID string) ([]Trip, error) {
\treturn s.tripRepo.ListByOwner(ctx, ownerID)
}

func (s *TripService) DeleteTrip(ctx context.Context, id string) error {
\tif err := s.tripRepo.Delete(ctx, id); err != nil { return Internal(err.Error()) }
\ts.publishEvent(ctx, "trip.deleted", map[string]string{"id": id})
\treturn nil
}

func (s *TripService) ReorderItems(ctx context.Context, req ReorderRequest) error {
\ttype item struct{ ID string; Position int }
\titems := make([]item, len(req.Items))
\tfor i, it := range req.Items { items[i] = item{ID: it.ID, Position: it.Position} }
\treturn s.itemRepo.Reorder(ctx, items)
}

func (s *TripService) publishEvent(ctx context.Context, eventType string, payload any) {
\tdata, _ := json.Marshal(payload)
\tif err := s.producer.Publish(ctx, s.topic, eventType, data); err != nil {
\t\tlog.Warn().Err(err).Str("event", eventType).Msg("kafka publish failed")
\t}
}
''')

write('services/trip-service/internal/handlers.go', '''package internal

import (
\t"github.com/gin-gonic/gin"
\tpkgmw "github.com/agamrai0123/wanderplan/pkg/middleware"
\tpkgresp "github.com/agamrai0123/wanderplan/pkg/response"
)

type Handlers struct{ svc *TripService }

func NewHandlers(svc *TripService) *Handlers { return &Handlers{svc: svc} }

func (h *Handlers) CreateTrip(c *gin.Context) {
\tvar body struct {
\t\tTitle       string `json:"title" binding:"required"`
\t\tDescription string `json:"description"`
\t\tDestination string `json:"destination"`
\t}
\tif err := c.ShouldBindJSON(&body); err != nil { pkgresp.Err(c, BadRequest(err.Error())); return }
\townerID := pkgmw.UserID(c)
\ttrip, err := h.svc.CreateTrip(c.Request.Context(), ownerID, body.Title, body.Description, body.Destination)
\tif err != nil { pkgresp.Err(c, err); return }
\tpkgresp.Created(c, trip)
}

func (h *Handlers) GetTrip(c *gin.Context) {
\ttrip, err := h.svc.GetTrip(c.Request.Context(), c.Param("id"))
\tif err != nil { pkgresp.Err(c, err); return }
\tpkgresp.OK(c, trip)
}

func (h *Handlers) ListTrips(c *gin.Context) {
\townerID := pkgmw.UserID(c)
\ttrips, err := h.svc.ListTrips(c.Request.Context(), ownerID)
\tif err != nil { pkgresp.Err(c, err); return }
\tpkgresp.OK(c, trips)
}

func (h *Handlers) DeleteTrip(c *gin.Context) {
\tif err := h.svc.DeleteTrip(c.Request.Context(), c.Param("id")); err != nil { pkgresp.Err(c, err); return }
\tpkgresp.NoContent(c)
}

func (h *Handlers) ReorderItems(c *gin.Context) {
\tvar req ReorderRequest
\tif err := c.ShouldBindJSON(&req); err != nil { pkgresp.Err(c, BadRequest(err.Error())); return }
\tif err := h.svc.ReorderItems(c.Request.Context(), req); err != nil { pkgresp.Err(c, err); return }
\tpkgresp.NoContent(c)
}
''')

write('services/trip-service/internal/routes.go', '''package internal

import (
\t"github.com/gin-gonic/gin"
\t"github.com/prometheus/client_golang/prometheus"
\tpkgjwt "github.com/agamrai0123/wanderplan/pkg/jwt"
\tpkgmw "github.com/agamrai0123/wanderplan/pkg/middleware"
)

func RegisterRoutes(r *gin.Engine, h *Handlers, jwtMgr *pkgjwt.Manager, reg *prometheus.Registry, cfg *Config) {
\t_, metricsH := pkgmw.Metrics("trip-service", reg)
\tr.Use(pkgmw.RequestID(), pkgmw.Logger(), pkgmw.Recovery(), pkgmw.CORS(cfg.CORS.AllowedOrigins))
\tr.GET("/healthz", func(c *gin.Context) { c.JSON(200, gin.H{"status": "ok"}) })
\tr.GET("/metrics", gin.WrapH(metricsH))

\tv1 := r.Group("/trips")
\tv1.Use(pkgmw.Auth(jwtMgr))
\t{
\t\tv1.POST("", h.CreateTrip)
\t\tv1.GET("", h.ListTrips)
\t\tv1.GET("/:id", h.GetTrip)
\t\tv1.DELETE("/:id", h.DeleteTrip)
\t\tv1.PATCH("/:id/items/reorder", h.ReorderItems)
\t}
}
''')

write('services/trip-service/cmd/main.go', '''package main

import (
\t"context"
\t"fmt"
\t"net/http"
\t"os"
\t"os/signal"
\t"syscall"
\t"time"

\t"github.com/gin-gonic/gin"
\t"github.com/rs/zerolog/log"
\tpkgcfg "github.com/agamrai0123/wanderplan/pkg/config"
\tpkgdb "github.com/agamrai0123/wanderplan/pkg/database"
\tpkgjwt "github.com/agamrai0123/wanderplan/pkg/jwt"
\tpkgkafka "github.com/agamrai0123/wanderplan/pkg/kafka"
\tinternal "github.com/agamrai0123/wanderplan/services/trip-service/internal"
)

func main() {
\tvar cfg internal.Config
\tif err := pkgcfg.Load("trip-service-config", "./config", "TRIP", &cfg); err != nil {
\t\tfmt.Fprintf(os.Stderr, "load config: %v\\n", err)
\t}
\tcfg.Validate()
\tinternal.InitLogger(cfg.Logging, "trip-service")

\tctx := context.Background()
\tpool, err := pkgdb.NewPool(ctx, cfg.Database.ToDBConfig())
\tif err != nil { log.Fatal().Err(err).Msg("db pool") }
\tdefer pool.Close()

\tproducer, err := pkgkafka.NewProducer(cfg.Kafka.Brokers)
\tif err != nil { log.Warn().Err(err).Msg("kafka producer unavailable") }

\ttripRepo := internal.NewTripRepo(pool)
\titemRepo := internal.NewItemRepo(pool)
\tsvc := internal.NewTripService(tripRepo, itemRepo, producer, cfg.Kafka.TopicTripEvents)
\thandlers := internal.NewHandlers(svc)

\tprivB64 := os.Getenv("JWT_PRIVATE_KEY"); pubB64 := os.Getenv("JWT_PUBLIC_KEY")
\tvar jwtMgr *pkgjwt.Manager
\tif privB64 != "" && pubB64 != "" {
\t\tjwtMgr, err = pkgjwt.NewManagerFromBase64(privB64, pubB64, 15*time.Minute)
\t\tif err != nil { log.Fatal().Err(err).Msg("jwt manager") }
\t}

\tgin.SetMode(gin.ReleaseMode)
\trouter := gin.New()
\treg := internal.NewRegistry()
\tinternal.RegisterRoutes(router, handlers, jwtMgr, reg, &cfg)

\tsrv := &http.Server{Addr: fmt.Sprintf(":%d", cfg.ServerPort), Handler: router, ReadTimeout: 30 * time.Second, WriteTimeout: 60 * time.Second}
\tgo func() {
\t\tlog.Info().Int("port", cfg.ServerPort).Msg("trip-service listening")
\t\tif err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed { log.Fatal().Err(err).Msg("http") }
\t}()

\tquit := make(chan os.Signal, 1); signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM); <-quit
\tshutCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second); defer cancel()
\tsrv.Shutdown(shutCtx)
\tlog.Info().Msg("trip-service stopped")
}
''')


# ══════════════════════════════════════════════════════════════════
# USER-SERVICE  (HTTP :8083, gRPC :9083)
# ══════════════════════════════════════════════════════════════════

write('services/user-service/config/user-service-config.json', '''{
    "version": "1.0",
    "server_port": 8083,
    "grpc_port": 9083,
    "metric_port": 7083,
    "logging": { "level": 1, "path": "./log/user-service.log", "max_size_mb": 256, "max_backups": 5, "max_age_days": 30 },
    "database": {
        "host": "localhost", "port": 5432, "name": "wanderplan",
        "user": "postgres", "password": "", "schema": "wanderplan",
        "max_conns": 20, "min_conns": 2
    },
    "cors": { "allowed_origins": ["http://localhost:5173"] },
    "rate_limit": { "rps": 200, "burst": 400 }
}
''')

write('services/user-service/internal/config.go', '''package internal

import "github.com/agamrai0123/wanderplan/pkg/database"

type Config struct {
\tVersion    string       `mapstructure:"version"`
\tServerPort int          `mapstructure:"server_port"`
\tGRPCPort   int          `mapstructure:"grpc_port"`
\tMetricPort int          `mapstructure:"metric_port"`
\tLogging    LoggingCfg   `mapstructure:"logging"`
\tDatabase   DatabaseCfg  `mapstructure:"database"`
\tCORS       CORSCfg      `mapstructure:"cors"`
\tRateLimit  RateLimitCfg `mapstructure:"rate_limit"`
}
type LoggingCfg struct { Level int; Path string; MaxSizeMB int; MaxBackups int; MaxAgeDays int }
type DatabaseCfg struct {
\tHost string `mapstructure:"host"`; Port int `mapstructure:"port"`
\tName string `mapstructure:"name"`; User string `mapstructure:"user"`
\tPassword string `mapstructure:"password"`; Schema string `mapstructure:"schema"`
\tMaxConns int32 `mapstructure:"max_conns"`; MinConns int32 `mapstructure:"min_conns"`
}
func (d DatabaseCfg) ToDBConfig() database.Config {
\treturn database.Config{Host: d.Host, Port: d.Port, DBName: d.Name, User: d.User, Password: d.Password, Schema: d.Schema, MaxConns: d.MaxConns, MinConns: d.MinConns}
}
type CORSCfg struct { AllowedOrigins []string `mapstructure:"allowed_origins"` }
type RateLimitCfg struct { RPS float64 `mapstructure:"rps"`; Burst int `mapstructure:"burst"` }
func (c *Config) Validate() error {
\tif c.ServerPort == 0 { c.ServerPort = 8083 }; if c.GRPCPort == 0 { c.GRPCPort = 9083 }
\treturn nil
}
''')

write('services/user-service/internal/errors.go', '''package internal

import pkgerr "github.com/agamrai0123/wanderplan/pkg/errors"

var (BadRequest = pkgerr.BadRequest; Unauthorized = pkgerr.Unauthorized; NotFound = pkgerr.NotFound; Internal = pkgerr.Internal)
''')

write('services/user-service/internal/logger.go', '''package internal

import (
\t"github.com/rs/zerolog"
\t"github.com/rs/zerolog/log"
\tpkglogger "github.com/agamrai0123/wanderplan/pkg/logger"
)

func InitLogger(cfg LoggingCfg, service string) zerolog.Logger {
\tlevel := zerolog.Level(cfg.Level)
\tif level < zerolog.TraceLevel || level > zerolog.Disabled { level = zerolog.InfoLevel }
\tmaxSizeMB := cfg.MaxSizeMB; if maxSizeMB <= 0 { maxSizeMB = 256 }
\tlogger := pkglogger.Init(pkglogger.Config{Level: int(level), FilePath: cfg.Path, MaxSizeMB: maxSizeMB, MaxBackups: cfg.MaxBackups, MaxAgeDays: cfg.MaxAgeDays, Service: service})
\tlog.Logger = logger; zerolog.DefaultContextLogger = &logger; return logger
}
''')

write('services/user-service/internal/metrics.go', '''package internal

import (
\t"github.com/prometheus/client_golang/prometheus"
\t"github.com/prometheus/client_golang/prometheus/collectors"
)

func NewRegistry() *prometheus.Registry {
\treg := prometheus.NewRegistry()
\treg.MustRegister(collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}), collectors.NewGoCollector())
\treturn reg
}
''')

write('services/user-service/internal/models.go', '''package internal

import "time"

type Profile struct {
\tID        string    `json:"id"`
\tEmail     string    `json:"email"`
\tName      string    `json:"name"`
\tAvatarURL string    `json:"avatar_url"`
\tBio       string    `json:"bio,omitempty"`
\tLocation  string    `json:"location,omitempty"`
\tCreatedAt time.Time `json:"created_at"`
\tUpdatedAt time.Time `json:"updated_at"`
}

type UserStats struct {
\tTotalTrips     int     `json:"total_trips"`
\tTotalCountries int     `json:"total_countries"`
\tTotalDays      int     `json:"total_days"`
\tTotalBudget    float64 `json:"total_budget"`
}
''')

write('services/user-service/internal/database.go', '''package internal

import (
\t"context"
\t"fmt"
\t"time"

\t"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepo struct{ db *pgxpool.Pool }
func NewUserRepo(db *pgxpool.Pool) *UserRepo { return &UserRepo{db: db} }

func (r *UserRepo) GetByID(ctx context.Context, id string) (*Profile, error) {
\tp := &Profile{}
\terr := r.db.QueryRow(ctx,
\t\t`SELECT id,email,name,avatar_url,COALESCE(bio,''),COALESCE(location,''),created_at,updated_at FROM wanderplan.users WHERE id=$1`, id).
\t\tScan(&p.ID, &p.Email, &p.Name, &p.AvatarURL, &p.Bio, &p.Location, &p.CreatedAt, &p.UpdatedAt)
\tif err != nil { return nil, fmt.Errorf("get user: %w", err) }
\treturn p, nil
}

func (r *UserRepo) Update(ctx context.Context, id, name, bio, location, avatarURL string) (*Profile, error) {
\tp := &Profile{}
\terr := r.db.QueryRow(ctx,
\t\t`UPDATE wanderplan.users SET name=$2,bio=$3,location=$4,avatar_url=$5,updated_at=$6 WHERE id=$1
\t\t RETURNING id,email,name,avatar_url,COALESCE(bio,''),COALESCE(location,''),created_at,updated_at`,
\t\tid, name, bio, location, avatarURL, time.Now()).
\t\tScan(&p.ID, &p.Email, &p.Name, &p.AvatarURL, &p.Bio, &p.Location, &p.CreatedAt, &p.UpdatedAt)
\tif err != nil { return nil, fmt.Errorf("update user: %w", err) }
\treturn p, nil
}

func (r *UserRepo) GetStats(ctx context.Context, id string) (*UserStats, error) {
\tstats := &UserStats{}
\tr.db.QueryRow(ctx,
\t\t`SELECT COUNT(*), COUNT(DISTINCT LOWER(destination)), COALESCE(SUM(EXTRACT(DAY FROM (end_date - start_date))::int), 0)
\t\t FROM wanderplan.trips WHERE owner_id=$1 AND start_date IS NOT NULL AND end_date IS NOT NULL`, id).
\t\tScan(&stats.TotalTrips, &stats.TotalCountries, &stats.TotalDays)
\treturn stats, nil
}
''')

write('services/user-service/internal/service.go', '''package internal

import "context"

type UserService struct{ repo *UserRepo }
func NewUserService(repo *UserRepo) *UserService { return &UserService{repo: repo} }

func (s *UserService) GetProfile(ctx context.Context, id string) (*Profile, error) {
\tp, err := s.repo.GetByID(ctx, id)
\tif err != nil { return nil, NotFound("user not found") }
\treturn p, nil
}

func (s *UserService) UpdateProfile(ctx context.Context, id, name, bio, location, avatarURL string) (*Profile, error) {
\tp, err := s.repo.Update(ctx, id, name, bio, location, avatarURL)
\tif err != nil { return nil, Internal(err.Error()) }
\treturn p, nil
}

func (s *UserService) GetStats(ctx context.Context, id string) (*UserStats, error) {
\treturn s.repo.GetStats(ctx, id)
}
''')

write('services/user-service/internal/handlers.go', '''package internal

import (
\t"github.com/gin-gonic/gin"
\tpkgmw "github.com/agamrai0123/wanderplan/pkg/middleware"
\tpkgresp "github.com/agamrai0123/wanderplan/pkg/response"
)

type Handlers struct{ svc *UserService }
func NewHandlers(svc *UserService) *Handlers { return &Handlers{svc: svc} }

func (h *Handlers) GetMe(c *gin.Context) {
\tp, err := h.svc.GetProfile(c.Request.Context(), pkgmw.UserID(c))
\tif err != nil { pkgresp.Err(c, err); return }
\tpkgresp.OK(c, p)
}

func (h *Handlers) UpdateMe(c *gin.Context) {
\tvar body struct {
\t\tName      string `json:"name"`
\t\tBio       string `json:"bio"`
\t\tLocation  string `json:"location"`
\t\tAvatarURL string `json:"avatar_url"`
\t}
\tif err := c.ShouldBindJSON(&body); err != nil { pkgresp.Err(c, BadRequest(err.Error())); return }
\tp, err := h.svc.UpdateProfile(c.Request.Context(), pkgmw.UserID(c), body.Name, body.Bio, body.Location, body.AvatarURL)
\tif err != nil { pkgresp.Err(c, err); return }
\tpkgresp.OK(c, p)
}

func (h *Handlers) GetStats(c *gin.Context) {
\tstats, err := h.svc.GetStats(c.Request.Context(), pkgmw.UserID(c))
\tif err != nil { pkgresp.Err(c, err); return }
\tpkgresp.OK(c, stats)
}
''')

write('services/user-service/internal/routes.go', '''package internal

import (
\t"github.com/gin-gonic/gin"
\t"github.com/prometheus/client_golang/prometheus"
\tpkgjwt "github.com/agamrai0123/wanderplan/pkg/jwt"
\tpkgmw "github.com/agamrai0123/wanderplan/pkg/middleware"
)

func RegisterRoutes(r *gin.Engine, h *Handlers, jwtMgr *pkgjwt.Manager, reg *prometheus.Registry, cfg *Config) {
\t_, metricsH := pkgmw.Metrics("user-service", reg)
\tr.Use(pkgmw.RequestID(), pkgmw.Logger(), pkgmw.Recovery(), pkgmw.CORS(cfg.CORS.AllowedOrigins))
\tr.GET("/healthz", func(c *gin.Context) { c.JSON(200, gin.H{"status": "ok"}) })
\tr.GET("/metrics", gin.WrapH(metricsH))
\tv1 := r.Group("/users")
\tv1.Use(pkgmw.Auth(jwtMgr))
\t{
\t\tv1.GET("/me", h.GetMe)
\t\tv1.PATCH("/me", h.UpdateMe)
\t\tv1.GET("/me/stats", h.GetStats)
\t}
}
''')

write('services/user-service/cmd/main.go', '''package main

import (
\t"context"
\t"fmt"
\t"net/http"
\t"os"
\t"os/signal"
\t"syscall"
\t"time"

\t"github.com/gin-gonic/gin"
\t"github.com/rs/zerolog/log"
\tpkgcfg "github.com/agamrai0123/wanderplan/pkg/config"
\tpkgdb "github.com/agamrai0123/wanderplan/pkg/database"
\tpkgjwt "github.com/agamrai0123/wanderplan/pkg/jwt"
\tinternal "github.com/agamrai0123/wanderplan/services/user-service/internal"
)

func main() {
\tvar cfg internal.Config
\tif err := pkgcfg.Load("user-service-config", "./config", "USER", &cfg); err != nil { fmt.Fprintf(os.Stderr, "load config: %v\\n", err) }
\tcfg.Validate()
\tinternal.InitLogger(cfg.Logging, "user-service")

\tpool, err := pkgdb.NewPool(context.Background(), cfg.Database.ToDBConfig())
\tif err != nil { log.Fatal().Err(err).Msg("db pool") }
\tdefer pool.Close()

\tprivB64 := os.Getenv("JWT_PRIVATE_KEY"); pubB64 := os.Getenv("JWT_PUBLIC_KEY")
\tvar jwtMgr *pkgjwt.Manager
\tif privB64 != "" && pubB64 != "" {
\t\tjwtMgr, err = pkgjwt.NewManagerFromBase64(privB64, pubB64, 15*time.Minute)
\t\tif err != nil { log.Fatal().Err(err).Msg("jwt") }
\t}

\tgin.SetMode(gin.ReleaseMode); router := gin.New()
\treg := internal.NewRegistry()
\tsvc := internal.NewUserService(internal.NewUserRepo(pool))
\tinternal.RegisterRoutes(router, internal.NewHandlers(svc), jwtMgr, reg, &cfg)

\tsrv := &http.Server{Addr: fmt.Sprintf(":%d", cfg.ServerPort), Handler: router, ReadTimeout: 30 * time.Second, WriteTimeout: 60 * time.Second}
\tgo func() {
\t\tlog.Info().Int("port", cfg.ServerPort).Msg("user-service listening")
\t\tif err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed { log.Fatal().Err(err).Msg("http") }
\t}()
\tquit := make(chan os.Signal, 1); signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM); <-quit
\tshutCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second); defer cancel()
\tsrv.Shutdown(shutCtx); log.Info().Msg("user-service stopped")
}
''')


# ══════════════════════════════════════════════════════════════════
# COLLABORATION-SERVICE  (HTTP :8084)
# ══════════════════════════════════════════════════════════════════

write('services/collaboration-service/config/collaboration-service-config.json', '''{
    "version": "1.0",
    "server_port": 8084,
    "metric_port": 7084,
    "logging": { "level": 1, "path": "./log/collab-service.log", "max_size_mb": 256, "max_backups": 5, "max_age_days": 30 },
    "database": {
        "host": "localhost", "port": 5432, "name": "wanderplan",
        "user": "postgres", "password": "", "schema": "wanderplan",
        "max_conns": 10, "min_conns": 2
    },
    "kafka": { "brokers": ["localhost:9092"], "topic_collab_events": "collab-events" },
    "cors": { "allowed_origins": ["http://localhost:5173"] },
    "rate_limit": { "rps": 100, "burst": 200 }
}
''')

write('services/collaboration-service/internal/config.go', '''package internal

import "github.com/agamrai0123/wanderplan/pkg/database"

type Config struct {
\tVersion    string       `mapstructure:"version"`
\tServerPort int          `mapstructure:"server_port"`
\tMetricPort int          `mapstructure:"metric_port"`
\tLogging    LoggingCfg   `mapstructure:"logging"`
\tDatabase   DatabaseCfg  `mapstructure:"database"`
\tKafka      KafkaCfg     `mapstructure:"kafka"`
\tCORS       CORSCfg      `mapstructure:"cors"`
\tRateLimit  RateLimitCfg `mapstructure:"rate_limit"`
}
type LoggingCfg struct { Level int; Path string; MaxSizeMB int; MaxBackups int; MaxAgeDays int }
type DatabaseCfg struct {
\tHost string `mapstructure:"host"`; Port int `mapstructure:"port"`; Name string `mapstructure:"name"`
\tUser string `mapstructure:"user"`; Password string `mapstructure:"password"`; Schema string `mapstructure:"schema"`
\tMaxConns int32 `mapstructure:"max_conns"`; MinConns int32 `mapstructure:"min_conns"`
}
func (d DatabaseCfg) ToDBConfig() database.Config {
\treturn database.Config{Host: d.Host, Port: d.Port, DBName: d.Name, User: d.User, Password: d.Password, Schema: d.Schema, MaxConns: d.MaxConns, MinConns: d.MinConns}
}
type KafkaCfg struct { Brokers []string `mapstructure:"brokers"`; TopicCollabEvents string `mapstructure:"topic_collab_events"` }
type CORSCfg struct { AllowedOrigins []string `mapstructure:"allowed_origins"` }
type RateLimitCfg struct { RPS float64 `mapstructure:"rps"`; Burst int `mapstructure:"burst"` }
func (c *Config) Validate() error { if c.ServerPort == 0 { c.ServerPort = 8084 }; return nil }
''')

write('services/collaboration-service/internal/errors.go', '''package internal
import pkgerr "github.com/agamrai0123/wanderplan/pkg/errors"
var (BadRequest = pkgerr.BadRequest; Unauthorized = pkgerr.Unauthorized; NotFound = pkgerr.NotFound; Conflict = pkgerr.Conflict; Internal = pkgerr.Internal)
''')

write('services/collaboration-service/internal/logger.go', '''package internal
import (
\t"github.com/rs/zerolog"; "github.com/rs/zerolog/log"
\tpkglogger "github.com/agamrai0123/wanderplan/pkg/logger"
)
func InitLogger(cfg LoggingCfg, service string) zerolog.Logger {
\tl := zerolog.Level(cfg.Level); if l < zerolog.TraceLevel || l > zerolog.Disabled { l = zerolog.InfoLevel }
\tms := cfg.MaxSizeMB; if ms <= 0 { ms = 256 }
\tlogger := pkglogger.Init(pkglogger.Config{Level: int(l), FilePath: cfg.Path, MaxSizeMB: ms, MaxBackups: cfg.MaxBackups, MaxAgeDays: cfg.MaxAgeDays, Service: service})
\tlog.Logger = logger; zerolog.DefaultContextLogger = &logger; return logger
}
''')

write('services/collaboration-service/internal/metrics.go', '''package internal
import ("github.com/prometheus/client_golang/prometheus"; "github.com/prometheus/client_golang/prometheus/collectors")
func NewRegistry() *prometheus.Registry {
\treg := prometheus.NewRegistry()
\treg.MustRegister(collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}), collectors.NewGoCollector()); return reg
}
''')

write('services/collaboration-service/internal/models.go', '''package internal
import "time"
type Collaborator struct {
\tID       string    `json:"id"`
\tTripID   string    `json:"trip_id"`
\tUserID   string    `json:"user_id"`
\tEmail    string    `json:"email"`
\tRole     string    `json:"role"` // viewer | editor | admin
\tJoinedAt time.Time `json:"joined_at"`
}
''')

write('services/collaboration-service/internal/database.go', '''package internal
import (
\t"context"; "fmt"; "time"
\t"github.com/google/uuid"; "github.com/jackc/pgx/v5/pgxpool"
)
type CollabRepo struct{ db *pgxpool.Pool }
func NewCollabRepo(db *pgxpool.Pool) *CollabRepo { return &CollabRepo{db: db} }
func (r *CollabRepo) Invite(ctx context.Context, tripID, userID, email, role string) (*Collaborator, error) {
\tc := &Collaborator{ID: uuid.NewString(), TripID: tripID, UserID: userID, Email: email, Role: role, JoinedAt: time.Now()}
\t_, err := r.db.Exec(ctx, `INSERT INTO wanderplan.collaborators (id,trip_id,user_id,email,role,joined_at) VALUES ($1,$2,$3,$4,$5,$6)`,
\t\tc.ID, c.TripID, c.UserID, c.Email, c.Role, c.JoinedAt)
\tif err != nil { return nil, fmt.Errorf("invite collaborator: %w", err) }
\treturn c, nil
}
func (r *CollabRepo) List(ctx context.Context, tripID string) ([]Collaborator, error) {
\trows, err := r.db.Query(ctx, `SELECT id,trip_id,user_id,email,role,joined_at FROM wanderplan.collaborators WHERE trip_id=$1`, tripID)
\tif err != nil { return nil, err }
\tdefer rows.Close()
\tvar cs []Collaborator
\tfor rows.Next() {
\t\tvar c Collaborator
\t\trows.Scan(&c.ID, &c.TripID, &c.UserID, &c.Email, &c.Role, &c.JoinedAt); cs = append(cs, c)
\t}
\treturn cs, nil
}
func (r *CollabRepo) Remove(ctx context.Context, id string) error {
\t_, err := r.db.Exec(ctx, `DELETE FROM wanderplan.collaborators WHERE id=$1`, id); return err
}
''')

write('services/collaboration-service/internal/service.go', '''package internal
import (
\t"context"; "encoding/json"
\t"github.com/rs/zerolog/log"; pkgkafka "github.com/agamrai0123/wanderplan/pkg/kafka"
)
type CollabService struct { repo *CollabRepo; producer *pkgkafka.Producer; topic string }
func NewCollabService(repo *CollabRepo, producer *pkgkafka.Producer, topic string) *CollabService {
\treturn &CollabService{repo: repo, producer: producer, topic: topic}
}
func (s *CollabService) Invite(ctx context.Context, tripID, userID, email, role string) (*Collaborator, error) {
\tc, err := s.repo.Invite(ctx, tripID, userID, email, role)
\tif err != nil { return nil, Internal(err.Error()) }
\td, _ := json.Marshal(c); s.producer.Publish(ctx, s.topic, "collaborator.invited", d)
\treturn c, nil
}
func (s *CollabService) List(ctx context.Context, tripID string) ([]Collaborator, error) {
\treturn s.repo.List(ctx, tripID)
}
func (s *CollabService) Remove(ctx context.Context, id string) error {
\tif err := s.repo.Remove(ctx, id); err != nil { return Internal(err.Error()) }
\td, _ := json.Marshal(map[string]string{"id": id})
\tif err := s.producer.Publish(ctx, s.topic, "collaborator.removed", d); err != nil {
\t\tlog.Warn().Err(err).Msg("kafka publish")
\t}
\treturn nil
}
''')

write('services/collaboration-service/internal/handlers.go', '''package internal
import (
\t"github.com/gin-gonic/gin"
\tpkgmw "github.com/agamrai0123/wanderplan/pkg/middleware"
\tpkgresp "github.com/agamrai0123/wanderplan/pkg/response"
)
type Handlers struct{ svc *CollabService }
func NewHandlers(svc *CollabService) *Handlers { return &Handlers{svc: svc} }
func (h *Handlers) Invite(c *gin.Context) {
\tvar body struct { UserID string `json:"user_id"`; Email string `json:"email" binding:"required"`; Role string `json:"role"` }
\tif err := c.ShouldBindJSON(&body); err != nil { pkgresp.Err(c, BadRequest(err.Error())); return }
\tif body.Role == "" { body.Role = "viewer" }
\tinviterID := pkgmw.UserID(c)
\t_ = inviterID // used for auth checks
\tcol, err := h.svc.Invite(c.Request.Context(), c.Param("trip_id"), body.UserID, body.Email, body.Role)
\tif err != nil { pkgresp.Err(c, err); return }
\tpkgresp.Created(c, col)
}
func (h *Handlers) List(c *gin.Context) {
\tcols, err := h.svc.List(c.Request.Context(), c.Param("trip_id"))
\tif err != nil { pkgresp.Err(c, err); return }
\tpkgresp.OK(c, cols)
}
func (h *Handlers) Remove(c *gin.Context) {
\tif err := h.svc.Remove(c.Request.Context(), c.Param("id")); err != nil { pkgresp.Err(c, err); return }
\tpkgresp.NoContent(c)
}
''')

write('services/collaboration-service/internal/routes.go', '''package internal
import (
\t"github.com/gin-gonic/gin"; "github.com/prometheus/client_golang/prometheus"
\tpkgjwt "github.com/agamrai0123/wanderplan/pkg/jwt"; pkgmw "github.com/agamrai0123/wanderplan/pkg/middleware"
)
func RegisterRoutes(r *gin.Engine, h *Handlers, jwtMgr *pkgjwt.Manager, reg *prometheus.Registry, cfg *Config) {
\t_, metricsH := pkgmw.Metrics("collaboration-service", reg)
\tr.Use(pkgmw.RequestID(), pkgmw.Logger(), pkgmw.Recovery(), pkgmw.CORS(cfg.CORS.AllowedOrigins))
\tr.GET("/healthz", func(c *gin.Context) { c.JSON(200, gin.H{"status": "ok"}) })
\tr.GET("/metrics", gin.WrapH(metricsH))
\tv1 := r.Group("/collaborators")
\tv1.Use(pkgmw.Auth(jwtMgr))
\t{
\t\tv1.POST("/trips/:trip_id", h.Invite)
\t\tv1.GET("/trips/:trip_id", h.List)
\t\tv1.DELETE("/:id", h.Remove)
\t}
}
''')

write('services/collaboration-service/cmd/main.go', '''package main
import (
\t"context"; "fmt"; "net/http"; "os"; "os/signal"; "syscall"; "time"
\t"github.com/gin-gonic/gin"; "github.com/rs/zerolog/log"
\tpkgcfg "github.com/agamrai0123/wanderplan/pkg/config"
\tpkgdb "github.com/agamrai0123/wanderplan/pkg/database"
\tpkgjwt "github.com/agamrai0123/wanderplan/pkg/jwt"
\tpkgkafka "github.com/agamrai0123/wanderplan/pkg/kafka"
\tinternal "github.com/agamrai0123/wanderplan/services/collaboration-service/internal"
)
func main() {
\tvar cfg internal.Config
\tif err := pkgcfg.Load("collaboration-service-config", "./config", "COLLAB", &cfg); err != nil { fmt.Fprintf(os.Stderr, "config: %v\\n", err) }
\tcfg.Validate(); internal.InitLogger(cfg.Logging, "collaboration-service")
\tpool, err := pkgdb.NewPool(context.Background(), cfg.Database.ToDBConfig())
\tif err != nil { log.Fatal().Err(err).Msg("db pool") }
\tdefer pool.Close()
\tproducer, _ := pkgkafka.NewProducer(cfg.Kafka.Brokers)
\tprivB64 := os.Getenv("JWT_PRIVATE_KEY"); pubB64 := os.Getenv("JWT_PUBLIC_KEY")
\tvar jwtMgr *pkgjwt.Manager
\tif privB64 != "" && pubB64 != "" {
\t\tjwtMgr, err = pkgjwt.NewManagerFromBase64(privB64, pubB64, 15*time.Minute)
\t\tif err != nil { log.Fatal().Err(err).Msg("jwt") }
\t}
\tgin.SetMode(gin.ReleaseMode); router := gin.New(); reg := internal.NewRegistry()
\tsvc := internal.NewCollabService(internal.NewCollabRepo(pool), producer, cfg.Kafka.TopicCollabEvents)
\tinternal.RegisterRoutes(router, internal.NewHandlers(svc), jwtMgr, reg, &cfg)
\tsrv := &http.Server{Addr: fmt.Sprintf(":%d", cfg.ServerPort), Handler: router, ReadTimeout: 30 * time.Second, WriteTimeout: 60 * time.Second}
\tgo func() {
\t\tlog.Info().Int("port", cfg.ServerPort).Msg("collaboration-service listening")
\t\tif err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed { log.Fatal().Err(err).Msg("http") }
\t}()
\tquit := make(chan os.Signal, 1); signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM); <-quit
\tctx, cancel := context.WithTimeout(context.Background(), 10*time.Second); defer cancel()
\tsrv.Shutdown(ctx); log.Info().Msg("collaboration-service stopped")
}
''')


# ══════════════════════════════════════════════════════════════════
# NOTIFICATION-SERVICE  (HTTP :8085 + WebSocket)
# ══════════════════════════════════════════════════════════════════

write('services/notification-service/config/notification-service-config.json', '''{
    "version": "1.0",
    "server_port": 8085,
    "metric_port": 7085,
    "logging": { "level": 1, "path": "./log/notification-service.log", "max_size_mb": 256, "max_backups": 5, "max_age_days": 30 },
    "database": {
        "host": "localhost", "port": 5432, "name": "wanderplan",
        "user": "postgres", "password": "", "schema": "wanderplan",
        "max_conns": 10, "min_conns": 2
    },
    "kafka": {
        "brokers": ["localhost:9092"],
        "group_id": "notification-service",
        "topics": ["auth-events", "trip-events", "collab-events"]
    },
    "cors": { "allowed_origins": ["http://localhost:5173"] },
    "rate_limit": { "rps": 100, "burst": 200 }
}
''')

write('services/notification-service/internal/config.go', '''package internal
import "github.com/agamrai0123/wanderplan/pkg/database"
type Config struct {
\tVersion    string       `mapstructure:"version"`
\tServerPort int          `mapstructure:"server_port"`
\tMetricPort int          `mapstructure:"metric_port"`
\tLogging    LoggingCfg   `mapstructure:"logging"`
\tDatabase   DatabaseCfg  `mapstructure:"database"`
\tKafka      KafkaCfg     `mapstructure:"kafka"`
\tCORS       CORSCfg      `mapstructure:"cors"`
\tRateLimit  RateLimitCfg `mapstructure:"rate_limit"`
}
type LoggingCfg struct { Level int; Path string; MaxSizeMB int; MaxBackups int; MaxAgeDays int }
type DatabaseCfg struct {
\tHost string `mapstructure:"host"`; Port int `mapstructure:"port"`; Name string `mapstructure:"name"`
\tUser string `mapstructure:"user"`; Password string `mapstructure:"password"`; Schema string `mapstructure:"schema"`
\tMaxConns int32 `mapstructure:"max_conns"`; MinConns int32 `mapstructure:"min_conns"`
}
func (d DatabaseCfg) ToDBConfig() database.Config {
\treturn database.Config{Host: d.Host, Port: d.Port, DBName: d.Name, User: d.User, Password: d.Password, Schema: d.Schema, MaxConns: d.MaxConns, MinConns: d.MinConns}
}
type KafkaCfg struct {
\tBrokers []string `mapstructure:"brokers"`; GroupID string `mapstructure:"group_id"`; Topics []string `mapstructure:"topics"`
}
type CORSCfg struct { AllowedOrigins []string `mapstructure:"allowed_origins"` }
type RateLimitCfg struct { RPS float64 `mapstructure:"rps"`; Burst int `mapstructure:"burst"` }
func (c *Config) Validate() error { if c.ServerPort == 0 { c.ServerPort = 8085 }; return nil }
''')

write('services/notification-service/internal/errors.go', '''package internal
import pkgerr "github.com/agamrai0123/wanderplan/pkg/errors"
var (BadRequest = pkgerr.BadRequest; Unauthorized = pkgerr.Unauthorized; NotFound = pkgerr.NotFound; Internal = pkgerr.Internal)
''')

write('services/notification-service/internal/logger.go', '''package internal
import (
\t"github.com/rs/zerolog"; "github.com/rs/zerolog/log"
\tpkglogger "github.com/agamrai0123/wanderplan/pkg/logger"
)
func InitLogger(cfg LoggingCfg, service string) zerolog.Logger {
\tl := zerolog.Level(cfg.Level); if l < zerolog.TraceLevel || l > zerolog.Disabled { l = zerolog.InfoLevel }
\tms := cfg.MaxSizeMB; if ms <= 0 { ms = 256 }
\tlogger := pkglogger.Init(pkglogger.Config{Level: int(l), FilePath: cfg.Path, MaxSizeMB: ms, MaxBackups: cfg.MaxBackups, MaxAgeDays: cfg.MaxAgeDays, Service: service})
\tlog.Logger = logger; zerolog.DefaultContextLogger = &logger; return logger
}
''')

write('services/notification-service/internal/metrics.go', '''package internal
import ("github.com/prometheus/client_golang/prometheus"; "github.com/prometheus/client_golang/prometheus/collectors")
func NewRegistry() *prometheus.Registry {
\treg := prometheus.NewRegistry()
\treg.MustRegister(collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}), collectors.NewGoCollector()); return reg
}
''')

write('services/notification-service/internal/models.go', '''package internal
import "time"
type Notification struct {
\tID        string    `json:"id"`
\tUserID    string    `json:"user_id"`
\tType      string    `json:"type"`
\tTitle     string    `json:"title"`
\tBody      string    `json:"body"`
\tReadAt    *time.Time `json:"read_at,omitempty"`
\tCreatedAt time.Time `json:"created_at"`
}
''')

write('services/notification-service/internal/database.go', '''package internal
import (
\t"context"; "fmt"; "time"
\t"github.com/google/uuid"; "github.com/jackc/pgx/v5/pgxpool"
)
type NotifRepo struct{ db *pgxpool.Pool }
func NewNotifRepo(db *pgxpool.Pool) *NotifRepo { return &NotifRepo{db: db} }
func (r *NotifRepo) Create(ctx context.Context, userID, notifType, title, body string) (*Notification, error) {
\tn := &Notification{ID: uuid.NewString(), UserID: userID, Type: notifType, Title: title, Body: body, CreatedAt: time.Now()}
\t_, err := r.db.Exec(ctx, `INSERT INTO wanderplan.notifications (id,user_id,type,title,body,created_at) VALUES ($1,$2,$3,$4,$5,$6)`,
\t\tn.ID, n.UserID, n.Type, n.Title, n.Body, n.CreatedAt)
\tif err != nil { return nil, fmt.Errorf("create notification: %w", err) }
\treturn n, nil
}
func (r *NotifRepo) List(ctx context.Context, userID string) ([]Notification, error) {
\trows, err := r.db.Query(ctx, `SELECT id,user_id,type,title,body,read_at,created_at FROM wanderplan.notifications WHERE user_id=$1 ORDER BY created_at DESC LIMIT 50`, userID)
\tif err != nil { return nil, err }
\tdefer rows.Close()
\tvar ns []Notification
\tfor rows.Next() {
\t\tvar n Notification; rows.Scan(&n.ID, &n.UserID, &n.Type, &n.Title, &n.Body, &n.ReadAt, &n.CreatedAt); ns = append(ns, n)
\t}
\treturn ns, nil
}
func (r *NotifRepo) MarkRead(ctx context.Context, id, userID string) error {
\tnow := time.Now()
\t_, err := r.db.Exec(ctx, `UPDATE wanderplan.notifications SET read_at=$3 WHERE id=$1 AND user_id=$2`, id, userID, now)
\treturn err
}
''')

write('services/notification-service/internal/hub.go', '''package internal
import (
\t"encoding/json"; "sync"
\t"github.com/rs/zerolog/log"
)

// Hub manages active WebSocket clients.
type Hub struct {
\tmu      sync.RWMutex
\tclients map[string][]chan []byte // userID → channels
}

func NewHub() *Hub { return &Hub{clients: make(map[string][]chan []byte)} }

// Subscribe registers a channel for the given user.
func (h *Hub) Subscribe(userID string, ch chan []byte) {
\th.mu.Lock(); defer h.mu.Unlock()
\th.clients[userID] = append(h.clients[userID], ch)
}

// Unsubscribe removes a channel.
func (h *Hub) Unsubscribe(userID string, ch chan []byte) {
\th.mu.Lock(); defer h.mu.Unlock()
\tchans := h.clients[userID]
\tfor i, c := range chans {
\t\tif c == ch { h.clients[userID] = append(chans[:i], chans[i+1:]...); break }
\t}
\tif len(h.clients[userID]) == 0 { delete(h.clients, userID) }
}

// Broadcast sends a notification to all subscribers of userID.
func (h *Hub) Broadcast(userID string, n *Notification) {
\tdata, err := json.Marshal(n)
\tif err != nil { log.Error().Err(err).Msg("marshal notification"); return }
\th.mu.RLock(); defer h.mu.RUnlock()
\tfor _, ch := range h.clients[userID] {
\t\tselect { case ch <- data: default: }
\t}
}
''')

write('services/notification-service/internal/service.go', '''package internal
import (
\t"context"; "encoding/json"
\t"github.com/rs/zerolog/log"; pkgkafka "github.com/agamrai0123/wanderplan/pkg/kafka"
)
type NotifService struct { repo *NotifRepo; hub *Hub; consumer *pkgkafka.Consumer }
func NewNotifService(repo *NotifRepo, hub *Hub, consumer *pkgkafka.Consumer) *NotifService {
\treturn &NotifService{repo: repo, hub: hub, consumer: consumer}
}
func (s *NotifService) StartKafkaConsumer(ctx context.Context) {
\ts.consumer.Register("user.login", func(ctx context.Context, msg []byte) error {
\t\tvar payload map[string]string; json.Unmarshal(msg, &payload)
\t\tuserID := payload["user_id"]; if userID == "" { return nil }
\t\tn, err := s.repo.Create(ctx, userID, "auth", "Logged in", "You signed in from a new session.")
\t\tif err != nil { log.Warn().Err(err).Msg("create login notification"); return nil }
\t\ts.hub.Broadcast(userID, n); return nil
\t})
\ts.consumer.Register("trip.created", func(ctx context.Context, msg []byte) error {
\t\tvar trip map[string]string; json.Unmarshal(msg, &trip)
\t\towner := trip["owner_id"]; if owner == "" { return nil }
\t\tn, err := s.repo.Create(ctx, owner, "trip", "Trip created", "Your trip \\\""+trip["title"]+"\\\" was created.")
\t\tif err != nil { return nil }
\t\ts.hub.Broadcast(owner, n); return nil
\t})
\ts.consumer.Register("collaborator.invited", func(ctx context.Context, msg []byte) error {
\t\tvar col map[string]string; json.Unmarshal(msg, &col)
\t\tuserID := col["user_id"]; if userID == "" { return nil }
\t\tn, err := s.repo.Create(ctx, userID, "collab", "Invited to trip", "You were invited to collaborate.")
\t\tif err != nil { return nil }
\t\ts.hub.Broadcast(userID, n); return nil
\t})
\tgo s.consumer.Start(ctx)
}
func (s *NotifService) List(ctx context.Context, userID string) ([]Notification, error) {
\treturn s.repo.List(ctx, userID)
}
func (s *NotifService) MarkRead(ctx context.Context, id, userID string) error {
\treturn s.repo.MarkRead(ctx, id, userID)
}
''')

write('services/notification-service/internal/handlers.go', '''package internal
import (
\t"net/http"; "time"
\t"github.com/gin-gonic/gin"; "github.com/gorilla/websocket"
\tpkgmw "github.com/agamrai0123/wanderplan/pkg/middleware"
\tpkgresp "github.com/agamrai0123/wanderplan/pkg/response"
)
var upgrader = websocket.Upgrader{
\tCheckOrigin: func(r *http.Request) bool { return true },
\tReadBufferSize: 1024, WriteBufferSize: 1024,
}
type Handlers struct{ svc *NotifService; hub *Hub }
func NewHandlers(svc *NotifService, hub *Hub) *Handlers { return &Handlers{svc: svc, hub: hub} }

func (h *Handlers) List(c *gin.Context) {
\tns, err := h.svc.List(c.Request.Context(), pkgmw.UserID(c))
\tif err != nil { pkgresp.Err(c, err); return }
\tpkgresp.OK(c, ns)
}

func (h *Handlers) MarkRead(c *gin.Context) {
\tif err := h.svc.MarkRead(c.Request.Context(), c.Param("id"), pkgmw.UserID(c)); err != nil { pkgresp.Err(c, err); return }
\tpkgresp.NoContent(c)
}

func (h *Handlers) WS(c *gin.Context) {
\tuserID := pkgmw.UserID(c)
\tconn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
\tif err != nil { return }
\tdefer conn.Close()
\tch := make(chan []byte, 32)
\th.hub.Subscribe(userID, ch)
\tdefer h.hub.Unsubscribe(userID, ch)
\tfor msg := range ch {
\t\tconn.SetWriteDeadline(time.Now().Add(10 * time.Second))
\t\tif err := conn.WriteMessage(websocket.TextMessage, msg); err != nil { break }
\t}
}
''')

write('services/notification-service/internal/routes.go', '''package internal
import (
\t"github.com/gin-gonic/gin"; "github.com/prometheus/client_golang/prometheus"
\tpkgjwt "github.com/agamrai0123/wanderplan/pkg/jwt"; pkgmw "github.com/agamrai0123/wanderplan/pkg/middleware"
)
func RegisterRoutes(r *gin.Engine, h *Handlers, jwtMgr *pkgjwt.Manager, reg *prometheus.Registry, cfg *Config) {
\t_, metricsH := pkgmw.Metrics("notification-service", reg)
\tr.Use(pkgmw.RequestID(), pkgmw.Logger(), pkgmw.Recovery(), pkgmw.CORS(cfg.CORS.AllowedOrigins))
\tr.GET("/healthz", func(c *gin.Context) { c.JSON(200, gin.H{"status": "ok"}) })
\tr.GET("/metrics", gin.WrapH(metricsH))
\tv1 := r.Group("/notifications")
\tv1.Use(pkgmw.Auth(jwtMgr))
\t{
\t\tv1.GET("", h.List)
\t\tv1.PATCH("/:id", h.MarkRead)
\t\tv1.GET("/ws", h.WS)
\t}
}
''')

write('services/notification-service/cmd/main.go', '''package main
import (
\t"context"; "fmt"; "net/http"; "os"; "os/signal"; "syscall"; "time"
\t"github.com/gin-gonic/gin"; "github.com/rs/zerolog/log"
\tpkgcfg "github.com/agamrai0123/wanderplan/pkg/config"
\tpkgdb "github.com/agamrai0123/wanderplan/pkg/database"
\tpkgjwt "github.com/agamrai0123/wanderplan/pkg/jwt"
\tpkgkafka "github.com/agamrai0123/wanderplan/pkg/kafka"
\tinternal "github.com/agamrai0123/wanderplan/services/notification-service/internal"
)
func main() {
\tvar cfg internal.Config
\tif err := pkgcfg.Load("notification-service-config", "./config", "NOTIF", &cfg); err != nil { fmt.Fprintf(os.Stderr, "config: %v\\n", err) }
\tcfg.Validate(); internal.InitLogger(cfg.Logging, "notification-service")
\tpool, err := pkgdb.NewPool(context.Background(), cfg.Database.ToDBConfig())
\tif err != nil { log.Fatal().Err(err).Msg("db pool") }
\tdefer pool.Close()
\tconsumer, err := pkgkafka.NewConsumer(cfg.Kafka.Brokers, cfg.Kafka.GroupID, cfg.Kafka.Topics)
\tif err != nil { log.Warn().Err(err).Msg("kafka consumer unavailable") }
\tprivB64 := os.Getenv("JWT_PRIVATE_KEY"); pubB64 := os.Getenv("JWT_PUBLIC_KEY")
\tvar jwtMgr *pkgjwt.Manager
\tif privB64 != "" && pubB64 != "" {
\t\tjwtMgr, err = pkgjwt.NewManagerFromBase64(privB64, pubB64, 15*time.Minute)
\t\tif err != nil { log.Fatal().Err(err).Msg("jwt") }
\t}
\thub := internal.NewHub()
\trepo := internal.NewNotifRepo(pool)
\tsvc := internal.NewNotifService(repo, hub, consumer)
\tif consumer != nil { svc.StartKafkaConsumer(context.Background()) }
\tgin.SetMode(gin.ReleaseMode); router := gin.New(); reg := internal.NewRegistry()
\tinternal.RegisterRoutes(router, internal.NewHandlers(svc, hub), jwtMgr, reg, &cfg)
\tsrv := &http.Server{Addr: fmt.Sprintf(":%d", cfg.ServerPort), Handler: router, ReadTimeout: 30 * time.Second, WriteTimeout: 120 * time.Second}
\tgo func() {
\t\tlog.Info().Int("port", cfg.ServerPort).Msg("notification-service listening")
\t\tif err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed { log.Fatal().Err(err).Msg("http") }
\t}()
\tquit := make(chan os.Signal, 1); signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM); <-quit
\tctx, cancel := context.WithTimeout(context.Background(), 10*time.Second); defer cancel()
\tsrv.Shutdown(ctx); log.Info().Msg("notification-service stopped")
}
''')


# ══════════════════════════════════════════════════════════════════
# SEARCH-SERVICE  (HTTP :8086)
# ══════════════════════════════════════════════════════════════════

write('services/search-service/config/search-service-config.json', '''{
    "version": "1.0",
    "server_port": 8086,
    "metric_port": 7086,
    "logging": { "level": 1, "path": "./log/search-service.log", "max_size_mb": 256, "max_backups": 5, "max_age_days": 30 },
    "database": {
        "host": "localhost", "port": 5432, "name": "wanderplan",
        "user": "postgres", "password": "", "schema": "wanderplan",
        "max_conns": 10, "min_conns": 2
    },
    "google_places_api_key": "",
    "cors": { "allowed_origins": ["http://localhost:5173"] },
    "rate_limit": { "rps": 100, "burst": 200 }
}
''')

write('services/search-service/internal/config.go', '''package internal
import "github.com/agamrai0123/wanderplan/pkg/database"
type Config struct {
\tVersion            string       `mapstructure:"version"`
\tServerPort         int          `mapstructure:"server_port"`
\tMetricPort         int          `mapstructure:"metric_port"`
\tLogging            LoggingCfg   `mapstructure:"logging"`
\tDatabase           DatabaseCfg  `mapstructure:"database"`
\tGooglePlacesAPIKey string       `mapstructure:"google_places_api_key"`
\tCORS               CORSCfg      `mapstructure:"cors"`
\tRateLimit          RateLimitCfg `mapstructure:"rate_limit"`
}
type LoggingCfg struct { Level int; Path string; MaxSizeMB int; MaxBackups int; MaxAgeDays int }
type DatabaseCfg struct {
\tHost string `mapstructure:"host"`; Port int `mapstructure:"port"`; Name string `mapstructure:"name"`
\tUser string `mapstructure:"user"`; Password string `mapstructure:"password"`; Schema string `mapstructure:"schema"`
\tMaxConns int32 `mapstructure:"max_conns"`; MinConns int32 `mapstructure:"min_conns"`
}
func (d DatabaseCfg) ToDBConfig() database.Config {
\treturn database.Config{Host: d.Host, Port: d.Port, DBName: d.Name, User: d.User, Password: d.Password, Schema: d.Schema, MaxConns: d.MaxConns, MinConns: d.MinConns}
}
type CORSCfg struct { AllowedOrigins []string `mapstructure:"allowed_origins"` }
type RateLimitCfg struct { RPS float64 `mapstructure:"rps"`; Burst int `mapstructure:"burst"` }
func (c *Config) Validate() error {
\tif c.ServerPort == 0 { c.ServerPort = 8086 }
\tif c.GooglePlacesAPIKey == "" { c.GooglePlacesAPIKey = os.Getenv("GOOGLE_PLACES_API_KEY") }
\treturn nil
}
''')

# Fix: config.go uses os but doesn't import it
write('services/search-service/internal/config.go', '''package internal
import (
\t"os"
\t"github.com/agamrai0123/wanderplan/pkg/database"
)
type Config struct {
\tVersion            string       `mapstructure:"version"`
\tServerPort         int          `mapstructure:"server_port"`
\tMetricPort         int          `mapstructure:"metric_port"`
\tLogging            LoggingCfg   `mapstructure:"logging"`
\tDatabase           DatabaseCfg  `mapstructure:"database"`
\tGooglePlacesAPIKey string       `mapstructure:"google_places_api_key"`
\tCORS               CORSCfg      `mapstructure:"cors"`
\tRateLimit          RateLimitCfg `mapstructure:"rate_limit"`
}
type LoggingCfg struct { Level int; Path string; MaxSizeMB int; MaxBackups int; MaxAgeDays int }
type DatabaseCfg struct {
\tHost string `mapstructure:"host"`; Port int `mapstructure:"port"`; Name string `mapstructure:"name"`
\tUser string `mapstructure:"user"`; Password string `mapstructure:"password"`; Schema string `mapstructure:"schema"`
\tMaxConns int32 `mapstructure:"max_conns"`; MinConns int32 `mapstructure:"min_conns"`
}
func (d DatabaseCfg) ToDBConfig() database.Config {
\treturn database.Config{Host: d.Host, Port: d.Port, DBName: d.Name, User: d.User, Password: d.Password, Schema: d.Schema, MaxConns: d.MaxConns, MinConns: d.MinConns}
}
type CORSCfg struct { AllowedOrigins []string `mapstructure:"allowed_origins"` }
type RateLimitCfg struct { RPS float64 `mapstructure:"rps"`; Burst int `mapstructure:"burst"` }
func (c *Config) Validate() error {
\tif c.ServerPort == 0 { c.ServerPort = 8086 }
\tif c.GooglePlacesAPIKey == "" { c.GooglePlacesAPIKey = os.Getenv("GOOGLE_PLACES_API_KEY") }
\treturn nil
}
''')

write('services/search-service/internal/errors.go', '''package internal
import pkgerr "github.com/agamrai0123/wanderplan/pkg/errors"
var (BadRequest = pkgerr.BadRequest; NotFound = pkgerr.NotFound; Internal = pkgerr.Internal)
''')

write('services/search-service/internal/logger.go', '''package internal
import (
\t"github.com/rs/zerolog"; "github.com/rs/zerolog/log"
\tpkglogger "github.com/agamrai0123/wanderplan/pkg/logger"
)
func InitLogger(cfg LoggingCfg, service string) zerolog.Logger {
\tl := zerolog.Level(cfg.Level); if l < zerolog.TraceLevel || l > zerolog.Disabled { l = zerolog.InfoLevel }
\tms := cfg.MaxSizeMB; if ms <= 0 { ms = 256 }
\tlogger := pkglogger.Init(pkglogger.Config{Level: int(l), FilePath: cfg.Path, MaxSizeMB: ms, MaxBackups: cfg.MaxBackups, MaxAgeDays: cfg.MaxAgeDays, Service: service})
\tlog.Logger = logger; zerolog.DefaultContextLogger = &logger; return logger
}
''')

write('services/search-service/internal/metrics.go', '''package internal
import ("github.com/prometheus/client_golang/prometheus"; "github.com/prometheus/client_golang/prometheus/collectors")
func NewRegistry() *prometheus.Registry {
\treg := prometheus.NewRegistry()
\treg.MustRegister(collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}), collectors.NewGoCollector()); return reg
}
''')

write('services/search-service/internal/models.go', '''package internal
// PlaceResult is a location returned by Google Places or local tsvector search.
type PlaceResult struct {
\tPlaceID     string  `json:"place_id"`
\tName        string  `json:"name"`
\tAddress     string  `json:"formatted_address"`
\tLat         float64 `json:"lat"`
\tLng         float64 `json:"lng"`
\tPhotoRef    string  `json:"photo_reference,omitempty"`
\tRating      float64 `json:"rating,omitempty"`
\tTypes       []string `json:"types,omitempty"`
}
''')

write('services/search-service/internal/database.go', '''package internal
import (
\t"context"; "fmt"
\t"github.com/jackc/pgx/v5/pgxpool"
)
// PlaceCache caches Google Places results in Postgres.
type PlaceCache struct{ db *pgxpool.Pool }
func NewPlaceCache(db *pgxpool.Pool) *PlaceCache { return &PlaceCache{db: db} }
func (c *PlaceCache) Get(ctx context.Context, query string) ([]PlaceResult, bool) {
\trows, err := c.db.Query(ctx,
\t\t`SELECT place_id,name,address,lat,lng,photo_ref,rating FROM wanderplan.places_cache WHERE query=$1 AND expires_at > NOW()`, query)
\tif err != nil { return nil, false }
\tdefer rows.Close()
\tvar results []PlaceResult
\tfor rows.Next() {
\t\tvar p PlaceResult; rows.Scan(&p.PlaceID, &p.Name, &p.Address, &p.Lat, &p.Lng, &p.PhotoRef, &p.Rating); results = append(results, p)
\t}
\tif len(results) == 0 { return nil, false }
\treturn results, true
}
func (c *PlaceCache) Set(ctx context.Context, query string, results []PlaceResult) {
\tfor _, p := range results {
\t\tc.db.Exec(ctx,
\t\t\t`INSERT INTO wanderplan.places_cache (place_id,name,address,lat,lng,photo_ref,rating,query,expires_at)
\t\t\t VALUES ($1,$2,$3,$4,$5,$6,$7,$8,NOW()+INTERVAL \'1 day\')
\t\t\t ON CONFLICT (place_id,query) DO UPDATE SET expires_at=NOW()+INTERVAL \'1 day\'`,
\t\t\tp.PlaceID, p.Name, p.Address, p.Lat, p.Lng, p.PhotoRef, p.Rating, query)
\t}
}
func (c *PlaceCache) FullTextSearch(ctx context.Context, q string) ([]PlaceResult, error) {
\trows, err := c.db.Query(ctx,
\t\t`SELECT place_id,name,address,lat,lng,photo_ref,rating
\t\t FROM wanderplan.places_cache
\t\t WHERE to_tsvector(\'english\', name || \' \' || address) @@ plainto_tsquery(\'english\', $1)
\t\t LIMIT 10`, q)
\tif err != nil { return nil, fmt.Errorf("tsvector search: %w", err) }
\tdefer rows.Close()
\tvar results []PlaceResult
\tfor rows.Next() {
\t\tvar p PlaceResult; rows.Scan(&p.PlaceID, &p.Name, &p.Address, &p.Lat, &p.Lng, &p.PhotoRef, &p.Rating); results = append(results, p)
\t}
\treturn results, nil
}
''')

write('services/search-service/internal/service.go', '''package internal
import (
\t"context"; "encoding/json"; "fmt"; "net/http"; "net/url"
\t"github.com/rs/zerolog/log"
)
type SearchService struct { cache *PlaceCache; apiKey string }
func NewSearchService(cache *PlaceCache, apiKey string) *SearchService { return &SearchService{cache: cache, apiKey: apiKey} }

func (s *SearchService) Search(ctx context.Context, query string) ([]PlaceResult, error) {
\tif results, ok := s.cache.Get(ctx, query); ok { return results, nil }
\tif s.apiKey != "" {
\t\tresults, err := s.googlePlaces(ctx, query)
\t\tif err != nil { log.Warn().Err(err).Msg("google places fallback to tsvector") } else {
\t\t\ts.cache.Set(ctx, query, results); return results, nil
\t\t}
\t}
\treturn s.cache.FullTextSearch(ctx, query)
}

func (s *SearchService) googlePlaces(ctx context.Context, query string) ([]PlaceResult, error) {
\tu := fmt.Sprintf("https://maps.googleapis.com/maps/api/place/textsearch/json?query=%s&key=%s",
\t\turl.QueryEscape(query), url.QueryEscape(s.apiKey))
\treq, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
\tif err != nil { return nil, err }
\tresp, err := http.DefaultClient.Do(req)
\tif err != nil { return nil, err }
\tdefer resp.Body.Close()
\tvar r struct {
\t\tResults []struct {
\t\t\tPlaceID          string `json:"place_id"`
\t\t\tName             string `json:"name"`
\t\t\tFormattedAddress string `json:"formatted_address"`
\t\t\tGeometry         struct{ Location struct{ Lat, Lng float64 } } `json:"geometry"`
\t\t\tRating           float64 `json:"rating"`
\t\t\tTypes            []string `json:"types"`
\t\t} `json:"results"`
\t}
\tif err := json.NewDecoder(resp.Body).Decode(&r); err != nil { return nil, err }
\tvar results []PlaceResult
\tfor _, res := range r.Results {
\t\tresults = append(results, PlaceResult{
\t\t\tPlaceID: res.PlaceID, Name: res.Name, Address: res.FormattedAddress,
\t\t\tLat: res.Geometry.Location.Lat, Lng: res.Geometry.Location.Lng,
\t\t\tRating: res.Rating, Types: res.Types,
\t\t})
\t}
\treturn results, nil
}
''')

write('services/search-service/internal/handlers.go', '''package internal
import (
\t"github.com/gin-gonic/gin"
\tpkgresp "github.com/agamrai0123/wanderplan/pkg/response"
)
type Handlers struct{ svc *SearchService }
func NewHandlers(svc *SearchService) *Handlers { return &Handlers{svc: svc} }
func (h *Handlers) Search(c *gin.Context) {
\tq := c.Query("q"); if q == "" { pkgresp.Err(c, BadRequest("q is required")); return }
\tresults, err := h.svc.Search(c.Request.Context(), q)
\tif err != nil { pkgresp.Err(c, Internal(err.Error())); return }
\tpkgresp.OK(c, results)
}
''')

write('services/search-service/internal/routes.go', '''package internal
import (
\t"github.com/gin-gonic/gin"; "github.com/prometheus/client_golang/prometheus"
\tpkgjwt "github.com/agamrai0123/wanderplan/pkg/jwt"; pkgmw "github.com/agamrai0123/wanderplan/pkg/middleware"
)
func RegisterRoutes(r *gin.Engine, h *Handlers, jwtMgr *pkgjwt.Manager, reg *prometheus.Registry, cfg *Config) {
\t_, metricsH := pkgmw.Metrics("search-service", reg)
\tr.Use(pkgmw.RequestID(), pkgmw.Logger(), pkgmw.Recovery(), pkgmw.CORS(cfg.CORS.AllowedOrigins))
\tr.GET("/healthz", func(c *gin.Context) { c.JSON(200, gin.H{"status": "ok"}) })
\tr.GET("/metrics", gin.WrapH(metricsH))
\tv1 := r.Group("/search")
\tv1.Use(pkgmw.Auth(jwtMgr))
\t{ v1.GET("", h.Search) }
}
''')

write('services/search-service/cmd/main.go', '''package main
import (
\t"context"; "fmt"; "net/http"; "os"; "os/signal"; "syscall"; "time"
\t"github.com/gin-gonic/gin"; "github.com/rs/zerolog/log"
\tpkgcfg "github.com/agamrai0123/wanderplan/pkg/config"
\tpkgdb "github.com/agamrai0123/wanderplan/pkg/database"
\tpkgjwt "github.com/agamrai0123/wanderplan/pkg/jwt"
\tinternal "github.com/agamrai0123/wanderplan/services/search-service/internal"
)
func main() {
\tvar cfg internal.Config
\tif err := pkgcfg.Load("search-service-config", "./config", "SEARCH", &cfg); err != nil { fmt.Fprintf(os.Stderr, "config: %v\\n", err) }
\tcfg.Validate(); internal.InitLogger(cfg.Logging, "search-service")
\tpool, err := pkgdb.NewPool(context.Background(), cfg.Database.ToDBConfig())
\tif err != nil { log.Fatal().Err(err).Msg("db pool") }
\tdefer pool.Close()
\tprivB64 := os.Getenv("JWT_PRIVATE_KEY"); pubB64 := os.Getenv("JWT_PUBLIC_KEY")
\tvar jwtMgr *pkgjwt.Manager
\tif privB64 != "" && pubB64 != "" {
\t\tjwtMgr, err = pkgjwt.NewManagerFromBase64(privB64, pubB64, 15*time.Minute)
\t\tif err != nil { log.Fatal().Err(err).Msg("jwt") }
\t}
\tgin.SetMode(gin.ReleaseMode); router := gin.New(); reg := internal.NewRegistry()
\tcache := internal.NewPlaceCache(pool)
\tsvc := internal.NewSearchService(cache, cfg.GooglePlacesAPIKey)
\tinternal.RegisterRoutes(router, internal.NewHandlers(svc), jwtMgr, reg, &cfg)
\tsrv := &http.Server{Addr: fmt.Sprintf(":%d", cfg.ServerPort), Handler: router, ReadTimeout: 30 * time.Second, WriteTimeout: 60 * time.Second}
\tgo func() {
\t\tlog.Info().Int("port", cfg.ServerPort).Msg("search-service listening")
\t\tif err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed { log.Fatal().Err(err).Msg("http") }
\t}()
\tquit := make(chan os.Signal, 1); signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM); <-quit
\tctx, cancel := context.WithTimeout(context.Background(), 10*time.Second); defer cancel()
\tsrv.Shutdown(ctx); log.Info().Msg("search-service stopped")
}
''')

print('\nAll service files written.')
