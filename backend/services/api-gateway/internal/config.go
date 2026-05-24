package internal

// Config is the top-level api-gateway configuration.
type Config struct {
	Version    string     `mapstructure:"version"`
	ServerPort int        `mapstructure:"server_port"`
	GRPCPort   int        `mapstructure:"grpc_port"`
	MetricPort int        `mapstructure:"metric_port"`
	Logging    LoggingCfg `mapstructure:"logging"`
	Services   ServicesCfg `mapstructure:"services"`
	CORS       CORSCfg    `mapstructure:"cors"`
	RateLimit  RateLimitCfg `mapstructure:"rate_limit"`
}

type LoggingCfg struct {
	Level      int    `mapstructure:"level"`
	Path       string `mapstructure:"path"`
	MaxSizeMB  int    `mapstructure:"max_size_mb"`
	MaxBackups int    `mapstructure:"max_backups"`
	MaxAgeDays int    `mapstructure:"max_age_days"`
}

type ServicesCfg struct {
	AuthAddr          string `mapstructure:"auth_addr"`
	TripAddr          string `mapstructure:"trip_addr"`
	UserAddr          string `mapstructure:"user_addr"`
	CollaborationAddr string `mapstructure:"collaboration_addr"`
	NotificationAddr  string `mapstructure:"notification_addr"`
	SearchAddr        string `mapstructure:"search_addr"`
}

type CORSCfg struct {
	AllowedOrigins []string `mapstructure:"allowed_origins"`
}

type RateLimitCfg struct {
	RPS   float64 `mapstructure:"rps"`
	Burst int     `mapstructure:"burst"`
}

// Validate returns an error if required fields are missing.
func (c *Config) Validate() error {
	if c.Services.AuthAddr == "" {
		c.Services.AuthAddr = "localhost:9081"
	}
	if c.ServerPort == 0 {
		c.ServerPort = 8080
	}
	if c.MetricPort == 0 {
		c.MetricPort = 7080
	}
	if c.RateLimit.RPS == 0 {
		c.RateLimit.RPS = 100
		c.RateLimit.Burst = 200
	}
	if len(c.CORS.AllowedOrigins) == 0 {
		c.CORS.AllowedOrigins = []string{"http://localhost:5173"}
	}
	return nil
}
