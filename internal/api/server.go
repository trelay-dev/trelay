package api

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/rs/zerolog"
)

// Server wraps the HTTP server with graceful shutdown support.
type Server struct {
	server *http.Server
	logger zerolog.Logger
}

// ServerConfig holds server configuration.
type ServerConfig struct {
	Address         string
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	ShutdownTimeout time.Duration
}

// NewServer creates a new HTTP server.
func NewServer(cfg ServerConfig, handler http.Handler, logger zerolog.Logger) *Server {
	return &Server{
		server: &http.Server{
			Addr:         cfg.Address,
			Handler:      handler,
			ReadTimeout:  cfg.ReadTimeout,
			WriteTimeout: cfg.WriteTimeout,
		},
		logger: logger,
	}
}

// Start begins listening for HTTP requests.
func (s *Server) Start() error {
	s.logger.Info().Str("address", s.server.Addr).Msg("starting HTTP server")
	if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("server error: %w", err)
	}
	return nil
}

// Shutdown gracefully stops the server.
func (s *Server) Shutdown(ctx context.Context) error {
	s.logger.Info().Msg("shutting down HTTP server")
	return s.server.Shutdown(ctx)
}
