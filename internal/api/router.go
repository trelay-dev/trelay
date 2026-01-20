package api

import (
	"time"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/rs/zerolog"

	"github.com/aftaab/trelay/internal/api/handler"
	"github.com/aftaab/trelay/internal/api/middleware"
	"github.com/aftaab/trelay/internal/core/analytics"
	"github.com/aftaab/trelay/internal/core/auth"
	"github.com/aftaab/trelay/internal/core/link"
)

// RouterConfig holds configuration for the router.
type RouterConfig struct {
	APIKeyHash      string
	JWTSecret       string
	TokenExpiry     time.Duration
	RateLimitPerMin int
	Logger          zerolog.Logger
}

// NewRouter creates a new chi router with all routes configured.
func NewRouter(
	cfg RouterConfig,
	linkService *link.Service,
	analyticsService *analytics.Service,
) *chi.Mux {
	r := chi.NewRouter()

	// JWT manager
	jwtManager := auth.NewJWTManager(cfg.JWTSecret, cfg.TokenExpiry, cfg.TokenExpiry*7)

	// Rate limiter
	rateLimiter := middleware.NewRateLimiter(cfg.RateLimitPerMin, time.Minute)

	// Global middleware
	r.Use(chimiddleware.RequestID)
	r.Use(chimiddleware.RealIP)
	r.Use(middleware.Logging(cfg.Logger))
	r.Use(chimiddleware.Recoverer)
	r.Use(middleware.RateLimit(rateLimiter))

	// CORS
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-API-Key"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// Initialize handlers
	healthHandler := handler.NewHealthHandler()
	authHandler := handler.NewAuthHandler(jwtManager, cfg.APIKeyHash)
	linkHandler := handler.NewLinkHandler(linkService)
	statsHandler := handler.NewStatsHandler(linkService, analyticsService)
	redirectHandler := handler.NewRedirectHandler(linkService, analyticsService)

	// Health check routes (no auth)
	r.Get("/healthz", healthHandler.Health)
	r.Get("/health", healthHandler.Health)
	r.Get("/readyz", healthHandler.Ready)

	// API routes
	r.Route("/api/v1", func(r chi.Router) {
		// Auth routes (no auth required)
		r.Post("/auth/login", authHandler.Login)
		r.Post("/auth/refresh", authHandler.Refresh)

		// Protected routes
		r.Group(func(r chi.Router) {
			r.Use(middleware.Auth(cfg.APIKeyHash, jwtManager))

			// Links
			r.Post("/links", linkHandler.Create)
			r.Get("/links", linkHandler.List)
			r.Get("/links/{slug}", linkHandler.Get)
			r.Patch("/links/{slug}", linkHandler.Update)
			r.Delete("/links/{slug}", linkHandler.Delete)
			r.Post("/links/{slug}/restore", linkHandler.Restore)

			// Stats
			r.Get("/stats/{slug}", statsHandler.GetStats)
			r.Get("/stats/{slug}/daily", statsHandler.GetDailyStats)
			r.Get("/stats/{slug}/monthly", statsHandler.GetMonthlyStats)
			r.Get("/stats/{slug}/referrers", statsHandler.GetReferrers)
		})
	})

	// Redirect route (public, at root level)
	r.Get("/{slug}", redirectHandler.Redirect)

	return r
}
