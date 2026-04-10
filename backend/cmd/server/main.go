package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"kbmanage/backend/internal/api/router"
	"kbmanage/backend/internal/repository"
)

func main() {
	cfg, err := repository.LoadConfig()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	db, err := repository.NewGormDB(cfg)
	if err != nil {
		log.Fatalf("failed to initialize database: %v", err)
	}

	rdb, err := repository.NewRedisClient(cfg)
	if err != nil {
		log.Fatalf("failed to initialize redis: %v", err)
	}
	defer func() {
		_ = rdb.Close()
	}()

	engine := router.NewRouter(db, rdb, cfg)

	httpSrv := &http.Server{
		Addr:              cfg.HTTPAddr,
		Handler:           engine,
		ReadHeaderTimeout: 5 * time.Second,
	}

	go func() {
		log.Printf("server listening on %s", cfg.HTTPAddr)
		if err := httpSrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen failed: %v", err)
		}
	}()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := httpSrv.Shutdown(ctx); err != nil {
		log.Printf("graceful shutdown failed: %v", err)
	}
}
