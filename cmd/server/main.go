package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rs/zerolog"

	"github.com/aftaab/trelay/internal/api"
	"github.com/aftaab/trelay/internal/config"
	"github.com/aftaab/trelay/internal/core/analytics"
	"github.com/aftaab/trelay/internal/core/auth"
	"github.com/aftaab/trelay/internal/core/folder"
	"github.com/aftaab/trelay/internal/core/link"
	"github.com/aftaab/trelay/internal/storage/sqlite"
)

func main() {
	// Initialize logger
	logger := zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}).
		With().
		Timestamp().
		Logger()

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to load configuration")
	}

	// Initialize database
	db, err := sqlite.Open(cfg.DSN())
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to open database")
	}
	defer db.Close()

	// Run migrations
	if err := db.Migrate(); err != nil {
		logger.Fatal().Err(err).Msg("failed to run migrations")
	}

	// Initialize repositories
	linkRepo := sqlite.NewLinkRepository(db)
	clickRepo := sqlite.NewClickRepository(db)
	folderRepo := sqlite.NewFolderRepository(db)

	// Initialize services
	linkService := link.NewService(
		linkRepo,
		cfg.App.SlugLength,
		cfg.App.CustomDomains,
	)

	analyticsService := analytics.NewService(
		clickRepo,
		cfg.App.IPAnonymization,
		cfg.App.AnalyticsEnabled,
	)

	folderService := folder.NewService(folderRepo)

	// Hash API key for comparison
	apiKeyHash := auth.HashAPIKey(cfg.Auth.APIKey)

	// Initialize router
	router := api.NewRouter(api.RouterConfig{
		APIKeyHash:      apiKeyHash,
		JWTSecret:       cfg.Auth.JWTSecret,
		TokenExpiry:     cfg.Auth.TokenExpiry,
		RateLimitPerMin: cfg.App.RateLimitPerMin,
		Logger:          logger,
		StaticDir:       cfg.App.StaticDir,
	}, linkService, analyticsService, folderService)

	// Initialize server
	server := api.NewServer(api.ServerConfig{
		Address:         cfg.Address(),
		ReadTimeout:     cfg.Server.ReadTimeout,
		WriteTimeout:    cfg.Server.WriteTimeout,
		ShutdownTimeout: cfg.Server.ShutdownTimeout,
	}, router, logger)

	// Start server in goroutine
	go func() {
		if err := server.Start(); err != nil {
			logger.Fatal().Err(err).Msg("server error")
		}
	}()

	logger.Info().
		Str("address", cfg.Address()).
		Str("base_url", cfg.App.BaseURL).
		Msg("trelay server started")

	// Wait for shutdown signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), cfg.Server.ShutdownTimeout)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Error().Err(err).Msg("server shutdown error")
	}

	logger.Info().Msg("server stopped")
}
