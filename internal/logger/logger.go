package logger

import (
	"Ex-L0/internal/config"
	"log/slog"
	"os"
)

type Logger struct{ *slog.Logger }

func NewLogger(cfg config.Log) *Logger {
	var h slog.Handler
	switch cfg.Format {
	case "json":
		h = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: level(cfg.Level)})
	default:
		h = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: level(cfg.Level)})
	}
	return &Logger{slog.New(h)}
}

func level(s string) slog.Level {
	switch s {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
