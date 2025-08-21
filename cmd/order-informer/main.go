package main

import (
	"Ex-L0/internal/config"
	"Ex-L0/internal/lib/logger/sl"
	"Ex-L0/internal/storage/postgresql"

	"log/slog"
	"os"
)

func main() {
	cfg := config.MustLoad()
	log := setupLogger(cfg.Env)
	log.Info("Starting Order Informer Service")
	log.Debug("debug messages are enabled")

	storage, err := postgresql.New(cfg.DB.URL)
	if err != nil {
		log.Error("Failed to initialize storage", sl.Err(err))
		os.Exit(1)
	}

	_ = storage
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	log = slog.New(
		slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
	)

	return log
}
