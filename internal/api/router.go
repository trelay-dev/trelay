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
	"github.com/aftaab/trelay/internal/core/folder"
	"github.com/aftaab/trelay/internal/core/link"
	"github.com/aftaab/trelay/internal/core/preview"
)

type RouterConfig struct {
	APIKeyHash      string
	JWTSecret       string
	TokenExpiry     time.Duration
	RateLimitPerMin int
	Logger          zerolog.Logger
}

func NewRouter(
	cfg RouterConfig,
	linkService *link.Service,
	analyticsService *analytics.Service,
	folderService *folder.Service,
) *chi.Mux {
	r := chi.NewRouter()

	jwtManager := auth.NewJWTManager(cfg.JWTSecret, cfg.TokenExpiry, cfg.TokenExpiry*7)
	rateLimiter := middleware.NewRateLimiter(cfg.RateLimitPerMin, time.Minute)

	r.Use(chimiddleware.RequestID)
	r.Use(chimiddleware.RealIP)
	r.Use(middleware.SecureHeaders)
	r.Use(middleware.Logging(cfg.Logger))
	r.Use(chimiddleware.Recoverer)
	r.Use(middleware.RateLimit(rateLimiter))
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-API-Key"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	previewService := preview.NewService()
	healthHandler := handler.NewHealthHandler()
	authHandler := handler.NewAuthHandler(jwtManager, cfg.APIKeyHash)
	linkHandler := handler.NewLinkHandler(linkService)
	statsHandler := handler.NewStatsHandler(linkService, analyticsService)
	previewHandler := handler.NewPreviewHandler(previewService)
	folderHandler := handler.NewFolderHandler(folderService)
	redirectHandler := handler.NewRedirectHandler(linkService, analyticsService)

	r.Get("/healthz", healthHandler.Health)
	r.Get("/health", healthHandler.Health)
	r.Get("/readyz", healthHandler.Ready)

	r.Route("/api/v1", func(r chi.Router) {
		r.Post("/auth/login", authHandler.Login)
		r.Post("/auth/refresh", authHandler.Refresh)

		r.Group(func(r chi.Router) {
			r.Use(middleware.Auth(cfg.APIKeyHash, jwtManager))

			r.Post("/links", linkHandler.Create)
			r.Get("/links", linkHandler.List)
			r.Delete("/links", linkHandler.BulkDelete)
			r.Get("/links/{slug}", linkHandler.Get)
			r.Patch("/links/{slug}", linkHandler.Update)
			r.Delete("/links/{slug}", linkHandler.Delete)
			r.Post("/links/{slug}/restore", linkHandler.Restore)

			r.Get("/preview", previewHandler.Fetch)

			r.Get("/stats/{slug}", statsHandler.GetStats)
			r.Get("/stats/{slug}/daily", statsHandler.GetDailyStats)
			r.Get("/stats/{slug}/monthly", statsHandler.GetMonthlyStats)
			r.Get("/stats/{slug}/referrers", statsHandler.GetReferrers)

			r.Post("/folders", folderHandler.Create)
			r.Get("/folders", folderHandler.List)
			r.Get("/folders/{id}", folderHandler.Get)
			r.Delete("/folders/{id}", folderHandler.Delete)
		})
	})

	r.Get("/{slug}", redirectHandler.Redirect)

	return r
}
