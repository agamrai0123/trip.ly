package internal

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	pkglogger "github.com/agamrai0123/wanderplan/pkg/logger"
)

// InitLogger initialises the global zerolog logger for the api-gateway.
func InitLogger(cfg LoggingCfg, service string) zerolog.Logger {
	level := zerolog.Level(cfg.Level)
	if level < zerolog.TraceLevel || level > zerolog.Disabled {
		level = zerolog.InfoLevel
	}
	maxSizeMB := cfg.MaxSizeMB
	if maxSizeMB <= 0 {
		maxSizeMB = 256
	}
	logger := pkglogger.Init(pkglogger.Config{
		Level:      int(level),
		FilePath:   cfg.Path,
		MaxSizeMB:  maxSizeMB,
		MaxBackups: cfg.MaxBackups,
		MaxAgeDays: cfg.MaxAgeDays,
		Service:    service,
	})
	log.Logger = logger
	zerolog.DefaultContextLogger = &logger
	return logger
}
