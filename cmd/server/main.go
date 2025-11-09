package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/AbubakarMahmood1/Serverless-URL-Analytics-Platform/config"
	"github.com/AbubakarMahmood1/Serverless-URL-Analytics-Platform/internal/api"
	"github.com/AbubakarMahmood1/Serverless-URL-Analytics-Platform/internal/api/handlers"
	"github.com/AbubakarMahmood1/Serverless-URL-Analytics-Platform/internal/repository/dynamodb"
	"github.com/AbubakarMahmood1/Serverless-URL-Analytics-Platform/internal/repository/redis"
	"github.com/AbubakarMahmood1/Serverless-URL-Analytics-Platform/internal/service"
	"github.com/gofiber/fiber/v2"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Initialize DynamoDB client
	dynamoClient, err := dynamodb.NewClient(cfg)
	if err != nil {
		log.Fatalf("Failed to create DynamoDB client: %v", err)
	}

	// Initialize Redis cache
	cache, err := redis.NewCache(cfg)
	if err != nil {
		log.Fatalf("Failed to create Redis cache: %v", err)
	}
	defer cache.Close()

	// Initialize repositories
	urlRepo := dynamodb.NewURLRepository(dynamoClient, cfg.DynamoDBURLsTable)
	analyticsRepo := dynamodb.NewAnalyticsRepository(dynamoClient, cfg.DynamoDBAnalyticsTable)

	// Initialize services
	urlService := service.NewURLService(urlRepo, cache, cfg)
	analyticsService := service.NewAnalyticsService(analyticsRepo, cache)
	qrService := service.NewQRService(cfg.BaseURL)

	// Initialize handlers
	urlHandler := handlers.NewURLHandler(urlService, analyticsService)
	analyticsHandler := handlers.NewAnalyticsHandler(analyticsService)
	qrHandler := handlers.NewQRHandler(qrService)
	healthHandler := handlers.NewHealthHandler()

	// Create Fiber app
	app := fiber.New(fiber.Config{
		AppName:      "URL Shortener API",
		ServerHeader: "URLShortener",
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			return c.Status(code).JSON(fiber.Map{
				"error": err.Error(),
			})
		},
	})

	// Setup routes
	api.SetupRoutes(app, urlHandler, analyticsHandler, qrHandler, healthHandler, analyticsService, cache, cfg)

	// Start server in a goroutine
	go func() {
		addr := fmt.Sprintf(":%s", cfg.Port)
		log.Printf("Server starting on %s", addr)
		log.Printf("Environment: %s", cfg.Env)
		log.Printf("Base URL: %s", cfg.BaseURL)

		if err := app.Listen(addr); err != nil {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")
	if err := app.Shutdown(); err != nil {
		log.Fatalf("Server shutdown failed: %v", err)
	}

	log.Println("Server stopped successfully")
}
