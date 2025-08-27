package app

import (
	"Ex-L0/internal/cache"
	"Ex-L0/internal/config"
	"Ex-L0/internal/http"

	"Ex-L0/internal/kafka"
	"Ex-L0/internal/logger"
	"Ex-L0/internal/repository/pg"
	"Ex-L0/internal/service"
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func Run(ctx context.Context) {
	cfg := config.Load()

	log := logger.NewLogger(cfg.Log)
	log.Info("starting order service", "env", cfg.Env)

	//DB
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		cfg.DB.User, cfg.DB.Password, cfg.DB.Host, cfg.DB.Port, cfg.DB.DBName, cfg.DB.SSLMode,
	)
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		panic(fmt.Errorf("pgxpool new: %w", err))
	}
	defer pool.Close()

	if err := pool.Ping(ctx); err != nil {
		panic(fmt.Errorf("db ping: %w", err))
	}
	log.Info("connected to postgres")

	//Repo/Service
	repo := pg.NewRepo(pool)
	ordersService := &service.OrdersService{Repository: repo}

	//Cache
	c := cache.New(cfg.Cache.MaxEntries)
	preloadN := cfg.Cache.MaxEntries
	if preloadN > 0 {
		_ = http.PrewarmCache(ctx, repo, c, preloadN)
	}
	log.Info("cache ready")

	//Kafka consumer
	cons := kafka.NewConsumer(cfg.Kafka, log, ordersService, c)
	go cons.Run(ctx)

	//HTTP server
	server := http.NewHTTPServer(cfg.HTTPServer, log, ordersService, c)
	if err := server.StartServer(ctx); err != nil {
		log.Error("http server stopped", "err", err)
	}

	time.Sleep(100 * time.Millisecond)
	log.Info("shutdown complete")
}
