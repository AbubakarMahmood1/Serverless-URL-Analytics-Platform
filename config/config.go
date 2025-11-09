package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	// Server
	Port string
	Env  string

	// DynamoDB
	AWSRegion            string
	AWSAccessKeyID       string
	AWSSecretAccessKey   string
	DynamoDBURLsTable    string
	DynamoDBAnalyticsTable string

	// Redis
	RedisURL      string
	RedisPassword string
	RedisDB       int

	// Application
	BaseURL          string
	ShortCodeLength  int
	RateLimitRequests int
	RateLimitWindow  int // seconds

	// GeoIP
	GeoIPDBPath string
}

var AppConfig *Config

// Load loads configuration from environment variables
func Load() *Config {
	// Load .env file if it exists (development)
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	redisDB, _ := strconv.Atoi(getEnv("REDIS_DB", "0"))
	shortCodeLength, _ := strconv.Atoi(getEnv("SHORT_CODE_LENGTH", "6"))
	rateLimitRequests, _ := strconv.Atoi(getEnv("RATE_LIMIT_REQUESTS", "100"))
	rateLimitWindow, _ := strconv.Atoi(getEnv("RATE_LIMIT_WINDOW", "60"))

	AppConfig = &Config{
		// Server
		Port: getEnv("PORT", "8080"),
		Env:  getEnv("ENV", "development"),

		// DynamoDB
		AWSRegion:              getEnv("AWS_REGION", "us-east-1"),
		AWSAccessKeyID:         getEnv("AWS_ACCESS_KEY_ID", ""),
		AWSSecretAccessKey:     getEnv("AWS_SECRET_ACCESS_KEY", ""),
		DynamoDBURLsTable:      getEnv("DYNAMODB_URLS_TABLE", "urls"),
		DynamoDBAnalyticsTable: getEnv("DYNAMODB_ANALYTICS_TABLE", "analytics"),

		// Redis
		RedisURL:      getEnv("REDIS_URL", "localhost:6379"),
		RedisPassword: getEnv("REDIS_PASSWORD", ""),
		RedisDB:       redisDB,

		// Application
		BaseURL:           getEnv("BASE_URL", "http://localhost:8080"),
		ShortCodeLength:   shortCodeLength,
		RateLimitRequests: rateLimitRequests,
		RateLimitWindow:   rateLimitWindow,

		// GeoIP
		GeoIPDBPath: getEnv("GEOIP_DB_PATH", ""),
	}

	return AppConfig
}

// getEnv retrieves an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
