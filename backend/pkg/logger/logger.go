// Package logger provides a zerolog-based structured logger with lumberjack log rotation.
package logger

import (
	"io"
	"os"
	"sync"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gopkg.in/natefinch/lumberjack.v2"
)

var once sync.Once

// Config holds logger initialisation parameters.
type Config struct {
	// Level is the zerolog integer level (-1=trace, 0=debug, 1=info, 2=warn, 3=error).
	Level int
	// FilePath is the path to the rotating log file. Empty means file logging is disabled.
	FilePath string
	// MaxSizeMB is the maximum log file size before rotation.
	MaxSizeMB int
	// MaxBackups is the number of rotated files to keep.
	MaxBackups int
	// MaxAgeDays is the maximum age of a rotated file.
	MaxAgeDays int
	// Service is the service name embedded in every log line.
	Service string
}

// Init initialises the global zerolog logger exactly once.
// Subsequent calls are no-ops and return the existing logger.
func Init(cfg Config) zerolog.Logger {
	once.Do(func() {
		zerolog.TimeFieldFormat = time.RFC3339Nano

		var writers []io.Writer

		// Console writer for local development
		writers = append(writers, zerolog.ConsoleWriter{
			Out:        os.Stderr,
			TimeFormat: time.RFC3339,
		})

		// File writer with rotation (when configured)
		if cfg.FilePath != "" {
			maxSize := cfg.MaxSizeMB
			if maxSize <= 0 {
				maxSize = 100
			}
			maxBackups := cfg.MaxBackups
			if maxBackups <= 0 {
				maxBackups = 10
			}
			maxAge := cfg.MaxAgeDays
			if maxAge <= 0 {
				maxAge = 14
			}
			writers = append(writers, &lumberjack.Logger{
				Filename:   cfg.FilePath,
				MaxSize:    maxSize,
				MaxBackups: maxBackups,
				MaxAge:     maxAge,
				Compress:   true,
			})
		}

		var w io.Writer
		if len(writers) == 1 {
			w = writers[0]
		} else {
			w = zerolog.MultiLevelWriter(writers...)
		}

		l := zerolog.New(w).
			Level(zerolog.Level(cfg.Level)).
			With().
			Timestamp().
			Str("service", cfg.Service).
			Logger()

		log.Logger = l
	})
	return log.Logger
}

// Get returns the current global zerolog logger.
func Get() zerolog.Logger { return log.Logger }
