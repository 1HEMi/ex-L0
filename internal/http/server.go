package http

import (
	"Ex-L0/internal/cache"
	"Ex-L0/internal/config"
	"Ex-L0/internal/logger"
	"Ex-L0/internal/repository"
	"Ex-L0/internal/service"
	"context"
	"embed"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

//go:embed ui/*
var uiFS embed.FS

type HTTPServer struct {
	cfg    config.HTTPServer
	log    *logger.Logger
	orders *service.OrdersService
	cache  *cache.Cache
}

func NewHTTPServer(cfg config.HTTPServer, log *logger.Logger,
	orders *service.OrdersService, c *cache.Cache) *HTTPServer {
	return &HTTPServer{
		cfg:    cfg,
		log:    log,
		orders: orders,
		cache:  c,
	}
}

func PrewarmCache(ctx context.Context, repo repository.OrdersReader, c *cache.Cache, n int) error {
	uids, err := repo.SelectRecentIDS(ctx, n)
	if err != nil {
		return err
	}

	for _, uid := range uids {
		if o, err := repo.GetByID(ctx, uid); err == nil {
			c.Set(uid, o)
		}
	}
	return nil
}

func (s *HTTPServer) StartServer(ctx context.Context) error {
	router := mux.NewRouter()

	router.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	}).Methods(http.MethodGet)

	h := &Handlers{Log: s.log, OrdersService: s.orders, Cache: s.cache}
	router.HandleFunc("/order/{order_uid}", h.HandlerGetOrderByID).
		Methods(http.MethodGet)
	router.HandleFunc("/order", func(w http.ResponseWriter, r *http.Request) {
		b, err := uiFS.ReadFile("ui/order.html")
		if err != nil {
			http.Error(w, "order page not found", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(b)
	}).Methods(http.MethodGet)
	router.PathPrefix("/").Handler(http.FileServer(http.FS(uiFS)))
	server := &http.Server{
		Addr:              s.cfg.Address,
		Handler:           loggingMiddleware(s.log, router),
		ReadTimeout:       parseDuration(string(s.cfg.Timeout)),
		WriteTimeout:      parseDuration(string(s.cfg.Timeout)),
		IdleTimeout:       parseDuration(string(s.cfg.Timeout)),
		ReadHeaderTimeout: 3 * time.Second,
	}
	return server.ListenAndServe()
}

func parseDuration(s string) time.Duration {
	d, err := time.ParseDuration(s)
	if err != nil {
		return 4 * time.Second
	}

	return d
}
