// Package database provides a pgx/v5 connection pool factory and golang-migrate runner.
package database

import (
	"context"
	"fmt"
	"net/url"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
)

// Config holds all parameters needed to connect to PostgreSQL.
type Config struct {
	Host            string
	Port            int
	DBName          string
	User            string
	Password        string
	Schema          string
	MaxConns        int32
	MinConns        int32
	MaxConnLifetime time.Duration
	MaxConnIdleTime time.Duration
}

// DSN returns a libpq-compatible connection string.
// All values are properly encoded to prevent injection.
func (c Config) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%d dbname=%s user=%s password=%s sslmode=disable search_path=%s,public",
		c.Host, c.Port, c.DBName, c.User, c.Password, c.Schema,
	)
}

// MigrationURL builds a URL suitable for golang-migrate.
func (c Config) MigrationURL() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=disable&search_path=%s&x-migrations-table=schema_migrations",
		url.QueryEscape(c.User),
		url.QueryEscape(c.Password),
		c.Host,
		c.Port,
		c.DBName,
		c.Schema,
	)
}

// NewPool creates a validated pgxpool connection pool.
func NewPool(ctx context.Context, cfg Config) (*pgxpool.Pool, error) {
	poolCfg, err := pgxpool.ParseConfig(cfg.DSN())
	if err != nil {
		return nil, fmt.Errorf("parse pool config: %w", err)
	}

	if cfg.MaxConns > 0 {
		poolCfg.MaxConns = cfg.MaxConns
	}
	if cfg.MinConns > 0 {
		poolCfg.MinConns = cfg.MinConns
	}
	if cfg.MaxConnLifetime > 0 {
		poolCfg.MaxConnLifetime = cfg.MaxConnLifetime
	}
	if cfg.MaxConnIdleTime > 0 {
		poolCfg.MaxConnIdleTime = cfg.MaxConnIdleTime
	}

	pool, err := pgxpool.NewWithConfig(ctx, poolCfg)
	if err != nil {
		return nil, fmt.Errorf("create pool: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("ping database: %w", err)
	}

	log.Info().
		Str("host", cfg.Host).
		Int("port", cfg.Port).
		Str("db", cfg.DBName).
		Str("schema", cfg.Schema).
		Int32("max_conns", poolCfg.MaxConns).
		Msg("database pool ready")

	return pool, nil
}

// RunMigrations applies all pending up migrations from sourcePath.
// It is safe to call on every startup — if no changes exist it is a no-op.
func RunMigrations(dbURL, sourcePath string) error {
	m, err := migrate.New("file://"+sourcePath, dbURL)
	if err != nil {
		return fmt.Errorf("init migrator: %w", err)
	}
	defer m.Close()

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("run migrations: %w", err)
	}

	v, dirty, _ := m.Version()
	log.Info().
		Uint("version", v).
		Bool("dirty", dirty).
		Str("source", sourcePath).
		Msg("migrations applied")

	return nil
}
