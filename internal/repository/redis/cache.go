package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/AbubakarMahmood1/Serverless-URL-Analytics-Platform/config"
	"github.com/redis/go-redis/v9"
)

type Cache struct {
	client *redis.Client
}

// NewCache creates a new Redis cache client
func NewCache(cfg *config.Config) (*Cache, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisURL,
		Password: cfg.RedisPassword,
		DB:       cfg.RedisDB,
	})

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	return &Cache{client: client}, nil
}

// GetURL retrieves a cached URL
func (c *Cache) GetURL(ctx context.Context, shortCode string) (string, error) {
	key := fmt.Sprintf("url:%s", shortCode)
	result, err := c.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", nil // Cache miss
	}
	if err != nil {
		return "", err
	}
	return result, nil
}

// SetURL caches a URL with TTL
func (c *Cache) SetURL(ctx context.Context, shortCode, originalURL string, ttl time.Duration) error {
	key := fmt.Sprintf("url:%s", shortCode)
	return c.client.Set(ctx, key, originalURL, ttl).Err()
}

// DeleteURL removes a URL from cache
func (c *Cache) DeleteURL(ctx context.Context, shortCode string) error {
	key := fmt.Sprintf("url:%s", shortCode)
	return c.client.Del(ctx, key).Err()
}

// IncrementClickCount increments the click counter for analytics
func (c *Cache) IncrementClickCount(ctx context.Context, shortCode string) (int64, error) {
	key := fmt.Sprintf("analytics:%s:count", shortCode)
	return c.client.Incr(ctx, key).Result()
}

// GetClickCount retrieves the cached click count
func (c *Cache) GetClickCount(ctx context.Context, shortCode string) (int64, error) {
	key := fmt.Sprintf("analytics:%s:count", shortCode)
	result, err := c.client.Get(ctx, key).Int64()
	if err == redis.Nil {
		return 0, nil
	}
	return result, err
}

// CheckRateLimit checks if a request should be rate limited
func (c *Cache) CheckRateLimit(ctx context.Context, identifier string, maxRequests int, window time.Duration) (bool, error) {
	key := fmt.Sprintf("ratelimit:%s", identifier)

	// Increment counter
	count, err := c.client.Incr(ctx, key).Result()
	if err != nil {
		return false, err
	}

	// Set expiry on first request
	if count == 1 {
		if err := c.client.Expire(ctx, key, window).Err(); err != nil {
			return false, err
		}
	}

	// Check if limit exceeded
	return count > int64(maxRequests), nil
}

// GetRateLimitInfo returns current rate limit status
func (c *Cache) GetRateLimitInfo(ctx context.Context, identifier string) (int64, time.Duration, error) {
	key := fmt.Sprintf("ratelimit:%s", identifier)

	count, err := c.client.Get(ctx, key).Int64()
	if err == redis.Nil {
		return 0, 0, nil
	}
	if err != nil {
		return 0, 0, err
	}

	ttl, err := c.client.TTL(ctx, key).Result()
	if err != nil {
		return 0, 0, err
	}

	return count, ttl, nil
}

// Close closes the Redis connection
func (c *Cache) Close() error {
	return c.client.Close()
}
