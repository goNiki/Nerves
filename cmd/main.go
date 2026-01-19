package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/goNiki/Nerves/internal/config"
	"github.com/goNiki/Nerves/internal/di"
	"github.com/goNiki/Nerves/internal/infrastructure/httpserver"
)

func main() {
	cfg, err := config.Load("./internal/config/.env")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	container, err := di.NewContainer(cfg)
	if err != nil {
		log.Fatalf("Failed to create container: %v", err)
	}
	defer container.Close()

	logger := container.Logger.Log

	if err := container.RunMigrations(); err != nil {
		logger.Error("migrations failed", "error", err)
		os.Exit(1)
	}

	router := container.SetupRouter()

	workerCtx, cancelWorker := context.WithCancel(context.Background())

	go container.WebhookWorker.Run(workerCtx)
	logger.Info("webhook worker started")

	go container.CacheWorker.Run(workerCtx)
	logger.Info("cache refresh worker started")

	server := httpserver.New(router, cfg.Server.Port(), logger)

	go func() {
		if err := server.Start(); err != nil {
			logger.Error("server error", "error", err)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	cancelWorker()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		logger.Error("shutdown error", "error", err)
	}
}
