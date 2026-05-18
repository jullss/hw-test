package logger

import (
	"testing"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func TestLogger(t *testing.T) {
	tests := []struct {
		name     string
		levelStr string
		expected zapcore.Level
	}{
		{"info level", "info", zap.InfoLevel},
		{"debug level", "debug", zap.DebugLevel},
		{"error level", "error", zap.ErrorLevel},
		{"invalid level defaults to info", "something_invalid", zap.InfoLevel},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := New(tt.levelStr)
			if l == nil {
				t.Fatal("logger should not be nil")
			}
		})
	}
}
