package internal

import "fmt"

// Config is the top-level notification-service configuration.
type Config struct {
	Version    string       `mapstructure:"version"`
	ServerPort string       `mapstructure:"server_port"`
	Logging    LoggingCfg   `mapstructure:"logging"`
	Database   DatabaseCfg  `mapstructure:"database"`
	Kafka      KafkaCfg     `mapstructure:"kafka"`
	CORS       CORSCfg      `mapstructure:"cors"`
	RateLimit  RateLimitCfg `mapstructure:"rate_limiting"`
}

// LoggingCfg holds log output settings.
type LoggingCfg struct {
	Level      int    `mapstructure:"level"`
	Path       string `mapstructure:"path"`
	MaxSizeMB  int    `mapstructure:"max_size_mb"`
	MaxBackups int    `mapstructure:"max_backups"`
	MaxAgeDays int    `mapstructure:"max_age_days"`
}

// DatabaseCfg holds PostgreSQL connection settings.
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

// KafkaCfg holds Kafka settings for consuming events.
type KafkaCfg struct {
	Brokers         []string `mapstructure:"brokers"`
	ConsumerGroupID string   `mapstructure:"consumer_group_id"`
}

// CORSCfg holds allowed origin settings.
type CORSCfg struct {
	AllowedOrigins []string `mapstructure:"allowed_origins"`
}

// RateLimitCfg holds rate-limit parameters.
type RateLimitCfg struct {
	RPS   float64 `mapstructure:"rps"`
	Burst int     `mapstructure:"burst"`
}

// Validate checks required config fields.
func (c *Config) Validate() error {
	if c.ServerPort == "" {
		return fmt.Errorf("server_port is required")
	}
	return nil
}
