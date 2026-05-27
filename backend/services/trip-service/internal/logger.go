package internal

import (
	pkglogger "github.com/agamrai0123/wanderplan/pkg/logger"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// InitLogger initialises the global zerolog logger for trip-service.
func InitLogger(cfg LoggingCfg, service string) zerolog.Logger {
	maxSizeMB := cfg.MaxSizeMB
	if maxSizeMB <= 0 {
		maxSizeMB = 100
	}
	logger := pkglogger.Init(pkglogger.Config{
		Level:      cfg.Level,
		FilePath:   cfg.Path,
		MaxSizeMB:  maxSizeMB,
		MaxBackups: cfg.MaxBackups,
		MaxAgeDays: cfg.MaxAgeDays,
		Service:    service,
	})
	log.Logger = logger
	return logger
}
