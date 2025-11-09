package api

import (
	"github.com/AbubakarMahmood1/Serverless-URL-Analytics-Platform/config"
	"github.com/AbubakarMahmood1/Serverless-URL-Analytics-Platform/internal/api/handlers"
	"github.com/AbubakarMahmood1/Serverless-URL-Analytics-Platform/internal/api/middleware"
	"github.com/AbubakarMahmood1/Serverless-URL-Analytics-Platform/internal/repository/redis"
	"github.com/AbubakarMahmood1/Serverless-URL-Analytics-Platform/internal/service"
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(
	app *fiber.App,
	urlHandler *handlers.URLHandler,
	analyticsHandler *handlers.AnalyticsHandler,
	qrHandler *handlers.QRHandler,
	healthHandler *handlers.HealthHandler,
	analyticsService *service.AnalyticsService,
	cache *redis.Cache,
	cfg *config.Config,
) {
	// Global middleware
	app.Use(middleware.Logger())
	app.Use(middleware.CORS())

	// Health check
	app.Get("/health", healthHandler.HealthCheck)

	// API routes with rate limiting
	api := app.Group("/api", middleware.RateLimiter(cache, cfg))
	{
		// URL shortening
		api.Post("/shorten", urlHandler.ShortenURL)

		// Link management
		api.Get("/links", urlHandler.ListUserURLs)
		api.Get("/links/:shortCode", urlHandler.GetURLInfo)
		api.Delete("/links/:shortCode", urlHandler.DeleteURL)

		// Analytics
		api.Get("/analytics/:shortCode", analyticsHandler.GetAnalytics)

		// QR code generation
		api.Get("/qr/:shortCode", qrHandler.GenerateQR)
	}

	// URL redirect (no rate limiting for better UX)
	// Apply analytics tracking middleware
	app.Get("/:shortCode", middleware.AnalyticsTracker(analyticsService), urlHandler.RedirectURL)
}
