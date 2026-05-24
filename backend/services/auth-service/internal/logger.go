package internal

import (
	"os"

	pkglogger "github.com/agamrai0123/wanderplan/pkg/logger"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// InitLogger initialises the global zerolog logger for the auth-service.
func InitLogger(cfg LoggingCfg, service string) zerolog.Logger {
	level := zerolog.Level(cfg.Level)
	if level < zerolog.TraceLevel || level > zerolog.Disabled {
		level = zerolog.InfoLevel
	}
	maxSizeMB := cfg.MaxSizeMB
	if maxSizeMB <= 0 {
		maxSizeMB = 100
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
	if os.Getenv("ENV") != "production" {
		consoleWriter := zerolog.ConsoleWriter{Out: os.Stderr}
		log.Logger = zerolog.New(zerolog.MultiLevelWriter(consoleWriter, logger)).
			With().Timestamp().Str("service", service).Logger()
	}
	return logger
}
