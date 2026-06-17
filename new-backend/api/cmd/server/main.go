package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/borkanie/parliament-api/internal/config"
	"github.com/borkanie/parliament-api/internal/db"
	"github.com/borkanie/parliament-api/internal/handlers"
	"github.com/borkanie/parliament-api/internal/repository"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func main() {
	cfg := config.Load()
	setupLogging(cfg.LogLevel)

	database, err := db.Connect(cfg.MongoURI, cfg.DBName)
	if err != nil {
		slog.Error("failed to connect to MongoDB", "err", err)
		os.Exit(1)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := db.EnsureIndexes(ctx, database); err != nil {
		slog.Error("failed to create indexes", "err", err)
		os.Exit(1)
	}

	// Repositories
	partyRepo := repository.NewPartyRepo(database)
	politiciansRepo := repository.NewPoliticiansRepo(database)
	votingRepo := repository.NewVotingRepo(database)

	// Handlers
	partyHandler := handlers.NewPartyHandler(partyRepo)
	politiciansHandler := handlers.NewPoliticiansHandler(politiciansRepo)
	votingHandler := handlers.NewVotingHandler(votingRepo)

	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "OPTIONS"},
		AllowedHeaders: []string{"Accept", "Content-Type"},
	}))

	router.Get("/health", handlers.Health)

	// Match original C# controller casing exactly
	router.Route("/api/Party", func(subrouter chi.Router) {
		subrouter.Get("/all", partyHandler.All)
		subrouter.Get("/GetById/", partyHandler.GetByID)
		subrouter.Get("/query", partyHandler.Query)
	})

	router.Route("/api/Politicians", func(subrouter chi.Router) {
		subrouter.Get("/getAllPoliticians", politiciansHandler.GetAll)
		subrouter.Get("/GetById/", politiciansHandler.GetByID)
		subrouter.Get("/GetByName/", politiciansHandler.GetByName)
	})

	router.Route("/api/Voting", func(subrouter chi.Router) {
		subrouter.Get("/getAllRounds", votingHandler.GetAllRounds)
		subrouter.Get("/getRoundById/", votingHandler.GetRoundById)
		subrouter.Get("/GetResultForVote/", votingHandler.GetResultForVote)
		subrouter.Get("/GetAllVotesForARoundById/", votingHandler.GetAllVotesForARoundById)
	})

	server := &http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.Port),
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		slog.Info("server starting", "port", cfg.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("server error", "err", err)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	slog.Info("shutting down server")
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()
	if err := server.Shutdown(shutdownCtx); err != nil {
		slog.Error("server forced to shutdown", "err", err)
	}
}

func setupLogging(level string) {
	var logLevel slog.Level
	switch level {
	case "debug":
		logLevel = slog.LevelDebug
	case "info":
		logLevel = slog.LevelInfo
	case "warning", "warn":
		logLevel = slog.LevelWarn
	default:
		logLevel = slog.LevelWarn
	}
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: logLevel})))
}
