package logging

import (
	"os"
	"time"

	"github.com/rs/zerolog"
)

// Ensure returns the provided logger or a safe console logger for local development.
func Ensure(log *zerolog.Logger) *zerolog.Logger {
	if log != nil {
		return log
	}
	writer := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
	logger := zerolog.New(writer).With().Timestamp().Str("module", "hrms").Logger()
	return &logger
}

// First returns the first logger from an optional constructor argument list.
func First(logs ...*zerolog.Logger) *zerolog.Logger {
	if len(logs) == 0 {
		return nil
	}
	return logs[0]
}

// Component adds a stable component field while preserving the caller's logger settings.
func Component(log *zerolog.Logger, component string) *zerolog.Logger {
	logger := Ensure(log).With().Str("component", component).Logger()
	return &logger
}
