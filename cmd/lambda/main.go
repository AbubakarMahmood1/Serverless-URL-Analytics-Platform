package main

import (
	"context"
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
	"github.com/aws/aws-lambda-go/lambda"
	dynamodb_sdk "github.com/aws/aws-sdk-go-v2/service/dynamodb"
	fiberadapter "github.com/awslabs/aws-lambda-go-api-proxy/fiber"
	"github.com/gofiber/fiber/v2"
)

var fiberLambda *fiberadapter.FiberLambda

// Initialize the Fiber app for Lambda (runs once on cold start)
func init() {
	if isLambda() {
		app := setupApp()
		fiberLambda = fiberadapter.New(app)
	}
}

// Handler is the Lambda function handler
func Handler(ctx context.Context, req interface{}) (interface{}, error) {
	return fiberLambda.ProxyWithContext(ctx, req)
}

// setupApp creates and configures the Fiber application
func setupApp() *fiber.App {
	// Load configuration
	cfg := config.Load()

	// Initialize DynamoDB client
	var dynamoClient *dynamodb_sdk.Client
	var err error

	if cfg.UseLocalDynamoDB {
		log.Printf("Using Local DynamoDB at %s", cfg.DynamoDBLocalEndpoint)
		dynamoClient, err = dynamodb.NewLocalClient(cfg.DynamoDBLocalEndpoint)
	} else {
		log.Printf("Using AWS DynamoDB in region %s", cfg.AWSRegion)
		dynamoClient, err = dynamodb.NewClient(cfg)
	}

	if err != nil {
		log.Fatalf("Failed to create DynamoDB client: %v", err)
	}

	// Initialize Redis cache
	cache, err := redis.NewCache(cfg)
	if err != nil {
		log.Fatalf("Failed to create Redis cache: %v", err)
	}

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

	return app
}

// isLambda checks if running in AWS Lambda environment
func isLambda() bool {
	return os.Getenv("AWS_LAMBDA_FUNCTION_NAME") != ""
}

func main() {
	if isLambda() {
		// Running in Lambda - use Lambda handler
		log.Println("Starting in AWS Lambda mode")
		lambda.Start(Handler)
	} else {
		// Running as traditional server
		log.Println("Starting in server mode")
		runServer()
	}
}

// runServer runs the traditional HTTP server (non-Lambda)
func runServer() {
	app := setupApp()
	cfg := config.Load()

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
