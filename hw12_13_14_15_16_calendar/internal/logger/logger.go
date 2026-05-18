package logger

import (
	"log/slog"
	"os"
	"strings"
)

type Logger struct {
	*slog.Logger
}

func New(level string) *Logger {
	var result slog.Level

	switch strings.ToLower(level) {
	case "info":
		result = slog.LevelInfo
	case "error":
		result = slog.LevelError
	case "warn":
		result = slog.LevelWarn
	case "debug":
		result = slog.LevelDebug
	default:
		result = slog.LevelInfo
	}

	opts := &slog.HandlerOptions{Level: result}
	handler := slog.NewTextHandler(os.Stdout, opts)

	return &Logger{
		Logger: slog.New(handler),
	}
}
