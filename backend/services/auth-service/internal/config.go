package internal

import (
	"fmt"

	"github.com/agamrai0123/wanderplan/pkg/database"
)

// Config is the top-level auth-service configuration.
type Config struct {
	Version     string       `mapstructure:"version"`
	ServerPort  string       `mapstructure:"server_port"`
	GRPCPort    string       `mapstructure:"grpc_port"`
	MetricPort  int          `mapstructure:"metric_port"`
	Logging     LoggingCfg   `mapstructure:"logging"`
	Database    DatabaseCfg  `mapstructure:"database"`
	Kafka       KafkaCfg     `mapstructure:"kafka"`
	OAuth       OAuthCfg     `mapstructure:"oauth"`
	CORS        CORSCfg      `mapstructure:"cors"`
	RateLimit   RateLimitCfg `mapstructure:"rate_limiting"`
	Cookie      CookieCfg    `mapstructure:"cookie"`
	JWTPrivKey  string       `mapstructure:"jwt_private_key"`
	JWTPubKey   string       `mapstructure:"jwt_public_key"`
	FrontendURL string       `mapstructure:"frontend_url"`
}

type LoggingCfg struct {
	Level      int    `mapstructure:"level"`
	Path       string `mapstructure:"path"`
	MaxSizeMB  int    `mapstructure:"max_size_mb"`
	MaxBackups int    `mapstructure:"max_backups"`
	MaxAgeDays int    `mapstructure:"max_age_days"`
}

type DatabaseCfg struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Name     string `mapstructure:"name"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	SSLMode  string `mapstructure:"ssl_mode"`
	Schema   string `mapstructure:"schema"`
	MaxConns int32  `mapstructure:"max_conns"`
	MinConns int32  `mapstructure:"min_conns"`
}

func (d DatabaseCfg) ToDBConfig() database.Config {
	return database.Config{
		Host:     d.Host,
		Port:     d.Port,
		DBName:   d.Name,
		User:     d.User,
		Password: d.Password,
		Schema:   d.Schema,
		MaxConns: d.MaxConns,
		MinConns: d.MinConns,
	}
}

type KafkaCfg struct {
	Brokers         []string `mapstructure:"brokers"`
	TopicAuthEvents string   `mapstructure:"topic_auth_events"`
}

type OAuthProviderCfg struct {
	ClientID     string `mapstructure:"client_id"`
	ClientSecret string `mapstructure:"client_secret"`
	RedirectURL  string `mapstructure:"redirect_url"`
}

type OAuthCfg struct {
	Google OAuthProviderCfg `mapstructure:"google"`
	GitHub OAuthProviderCfg `mapstructure:"github"`
}

type CORSCfg struct {
	AllowedOrigins []string `mapstructure:"allowed_origins"`
}

type RateLimitCfg struct {
	RPS   float64 `mapstructure:"rps"`
	Burst int     `mapstructure:"burst"`
}

type CookieCfg struct {
	Domain          string `mapstructure:"domain"`
	Secure          bool   `mapstructure:"secure"`
	StateMaxAgeSecs int    `mapstructure:"state_max_age_secs"`
}

func (c *Config) Validate() error {
	if c.ServerPort == "" {
		return fmt.Errorf("server_port is required")
	}
	if c.GRPCPort == "" {
		c.GRPCPort = "9081"
	}
	if c.Database.Host == "" {
		return fmt.Errorf("database.host is required")
	}
	return nil
}
