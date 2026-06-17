package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/borkanie/parliament-scraper/internal/config"
	"github.com/borkanie/parliament-scraper/internal/db"
	"github.com/borkanie/parliament-scraper/internal/orchestrator"
)

func main() {
	cfg := config.Load()
	setupLogging(cfg.LogLevel)

	database, err := db.Connect(cfg.MongoURI, cfg.DBName)
	if err != nil {
		slog.Error("failed to connect to MongoDB", "err", err)
		os.Exit(1)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	orc := orchestrator.New(database, cfg.OpenAIKey, cfg.ProxyURL)

	if cfg.RunOnce {
		// Kubernetes CronJob mode: run once and exit
		orc.RunOnce(ctx)
		return
	}

	// Daemon mode: run on 12h ticker until SIGINT/SIGTERM
	go orc.Run(ctx)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	cancel()
}

func setupLogging(level string) {
	var logLevel slog.Level
	switch level {
	case "debug":
		logLevel = slog.LevelDebug
	case "info":
		logLevel = slog.LevelInfo
	default:
		logLevel = slog.LevelWarn
	}
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: logLevel})))
}
