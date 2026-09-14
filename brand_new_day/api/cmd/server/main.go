// Command server is the brand_new_day API entrypoint. It wires configuration
// into a MongoDB connection, builds the repository → service → controller
// stack, and serves the contract-first routes generated from openapi.yaml
// alongside a Swagger UI.
package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"

	"github.com/borkanie/brand-new-day-api/internal/config"
	"github.com/borkanie/brand-new-day-api/internal/controller"
	"github.com/borkanie/brand-new-day-api/internal/db"
	"github.com/borkanie/brand-new-day-api/internal/generated"
	"github.com/borkanie/brand-new-day-api/internal/repository"
	"github.com/borkanie/brand-new-day-api/internal/service"
	"github.com/borkanie/brand-new-day-api/internal/swagger"
)

const (
	serverReadHeaderTimeout = 10 * time.Second
	serverShutdownTimeout   = 15 * time.Second
	ensureIndexesTimeout    = 30 * time.Second
)

func main() {
	if err := run(); err != nil {
		slog.Error("server terminated with an error", "error", err)
		os.Exit(1)
	}
}

func run() error {
	configuration := config.Load()

	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: configuration.SlogLevel(),
	})))

	mongoClient, database, err := db.ConnectClient(configuration.MongoURI, configuration.DatabaseName)
	if err != nil {
		return err
	}
	defer func() {
		disconnectContext, cancelDisconnect := context.WithTimeout(context.Background(), serverShutdownTimeout)
		defer cancelDisconnect()
		if disconnectErr := mongoClient.Disconnect(disconnectContext); disconnectErr != nil {
			slog.Error("failed to disconnect from MongoDB", "error", disconnectErr)
		}
	}()

	ensureIndexesContext, cancelEnsureIndexes := context.WithTimeout(context.Background(), ensureIndexesTimeout)
	defer cancelEnsureIndexes()
	if err := db.EnsureIndexes(ensureIndexesContext, database); err != nil {
		return err
	}

	// Repositories — one per collection; votes live on voting_rounds, so there
	// is deliberately no vote repository.
	partyRepository := repository.NewPartyRepository(database)
	politicianRepository := repository.NewPoliticianRepository(database)
	lawBucketRepository := repository.NewLawBucketRepository(database)
	normativeRepository := repository.NewNormativeRepository(database)
	votingRoundRepository := repository.NewVotingRoundRepository(database)

	// Services.
	politicianService := service.NewPoliticianService(politicianRepository, partyRepository, votingRoundRepository)
	votingService := service.NewVotingService(
		votingRoundRepository,
		normativeRepository,
		lawBucketRepository,
		politicianRepository,
		partyRepository,
	)
	lawService := service.NewLawService(lawBucketRepository, normativeRepository, votingRoundRepository)

	// The single Controller implementing the generated ServerInterface.
	apiController := controller.NewController(politicianService, votingService, lawService)

	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Recoverer)
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{http.MethodGet, http.MethodOptions},
		AllowedHeaders: []string{"Accept", "Content-Type"},
		MaxAge:         300,
	}))

	// Swagger UI and the raw spec are mounted before the generated routes so
	// they sit outside the spec-validation middleware below — /swagger/ and
	// /openapi.yaml are not themselves described by the contract.
	swagger.MountRoutes(router)

	specValidatorMiddleware, err := swagger.NewSpecValidatorMiddleware()
	if err != nil {
		return err
	}

	// HandlerFromMux registers every generated route onto a sub-router that
	// validates incoming requests against the embedded spec.
	router.Group(func(validatedRouter chi.Router) {
		validatedRouter.Use(specValidatorMiddleware)
		generated.HandlerFromMux(apiController, validatedRouter)
	})

	httpServer := &http.Server{
		Addr:              ":" + configuration.Port,
		Handler:           router,
		ReadHeaderTimeout: serverReadHeaderTimeout,
	}

	serverErrors := make(chan error, 1)
	go func() {
		slog.Info("brand_new_day API listening",
			"port", configuration.Port,
			"database", configuration.DatabaseName,
			"swaggerUI", "http://localhost:"+configuration.Port+"/swagger/",
		)
		if listenErr := httpServer.ListenAndServe(); listenErr != nil && !errors.Is(listenErr, http.ErrServerClosed) {
			serverErrors <- listenErr
			return
		}
		serverErrors <- nil
	}()

	shutdownSignals := make(chan os.Signal, 1)
	signal.Notify(shutdownSignals, os.Interrupt, syscall.SIGTERM)

	select {
	case listenErr := <-serverErrors:
		return listenErr
	case receivedSignal := <-shutdownSignals:
		slog.Info("shutdown signal received, draining connections", "signal", receivedSignal.String())

		shutdownContext, cancelShutdown := context.WithTimeout(context.Background(), serverShutdownTimeout)
		defer cancelShutdown()
		if shutdownErr := httpServer.Shutdown(shutdownContext); shutdownErr != nil {
			return shutdownErr
		}
		slog.Info("server stopped cleanly")
		return nil
	}
}
