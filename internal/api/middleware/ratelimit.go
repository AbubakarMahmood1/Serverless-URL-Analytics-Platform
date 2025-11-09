package middleware

import (
	"time"

	"github.com/AbubakarMahmood1/Serverless-URL-Analytics-Platform/config"
	"github.com/AbubakarMahmood1/Serverless-URL-Analytics-Platform/internal/repository/redis"
	"github.com/gofiber/fiber/v2"
)

func RateLimiter(cache *redis.Cache, cfg *config.Config) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get IP address
		ip := c.IP()

		// Check rate limit
		exceeded, err := cache.CheckRateLimit(
			c.Context(),
			ip,
			cfg.RateLimitRequests,
			time.Duration(cfg.RateLimitWindow)*time.Second,
		)

		if err != nil {
			// Log error but don't block request
			return c.Next()
		}

		if exceeded {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"error": "Rate limit exceeded. Please try again later.",
			})
		}

		// Get rate limit info for headers
		count, ttl, err := cache.GetRateLimitInfo(c.Context(), ip)
		if err == nil {
			c.Set("X-RateLimit-Limit", string(rune(cfg.RateLimitRequests)))
			c.Set("X-RateLimit-Remaining", string(rune(cfg.RateLimitRequests-int(count))))
			c.Set("X-RateLimit-Reset", string(rune(time.Now().Add(ttl).Unix())))
		}

		return c.Next()
	}
}
