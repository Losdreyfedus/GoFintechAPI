package logger

import (
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// InitLogger initializes the global logger with proper configuration
func InitLogger(level string, pretty bool) {
	// Set log level
	var logLevel zerolog.Level
	switch level {
	case "debug":
		logLevel = zerolog.DebugLevel
	case "info":
		logLevel = zerolog.InfoLevel
	case "warn":
		logLevel = zerolog.WarnLevel
	case "error":
		logLevel = zerolog.ErrorLevel
	default:
		logLevel = zerolog.InfoLevel
	}

	// Configure logger
	zerolog.SetGlobalLevel(logLevel)
	zerolog.TimeFieldFormat = time.RFC3339

	// Use pretty console logging in development
	if pretty {
		log.Logger = log.Output(zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: time.RFC3339,
		})
	} else {
		// Use JSON format in production
		log.Logger = zerolog.New(os.Stdout).With().Timestamp().Logger()
	}
}

// GetLogger returns the global logger instance
func GetLogger() zerolog.Logger {
	return log.Logger
}

// Debug logs a debug message
func Debug(msg string, fields map[string]interface{}) {
	log.Debug().Fields(fields).Msg(msg)
}

// Info logs an info message
func Info(msg string, fields map[string]interface{}) {
	log.Info().Fields(fields).Msg(msg)
}

// Warn logs a warning message
func Warn(msg string, fields map[string]interface{}) {
	log.Warn().Fields(fields).Msg(msg)
}

// Error logs an error message
func Error(msg string, err error, fields map[string]interface{}) {
	if fields == nil {
		fields = make(map[string]interface{})
	}
	fields["error"] = err.Error()
	log.Error().Fields(fields).Msg(msg)
}

// Fatal logs a fatal message and exits
func Fatal(msg string, err error, fields map[string]interface{}) {
	if fields == nil {
		fields = make(map[string]interface{})
	}
	fields["error"] = err.Error()
	log.Fatal().Fields(fields).Msg(msg)
}
