package middleware

import (
	"github.com/AbubakarMahmood1/Serverless-URL-Analytics-Platform/internal/models"
	"github.com/AbubakarMahmood1/Serverless-URL-Analytics-Platform/internal/service"
	"github.com/AbubakarMahmood1/Serverless-URL-Analytics-Platform/internal/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/mileusna/useragent"
)

// AnalyticsTracker is middleware that tracks analytics for URL redirects
func AnalyticsTracker(analyticsService *service.AnalyticsService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get short code from context (should be set by handler)
		shortCode := c.Locals("short_code")
		if shortCode == nil {
			return c.Next()
		}

		// Parse User-Agent
		ua := useragent.Parse(c.Get("User-Agent"))

		// Determine device type
		deviceType := "desktop"
		if ua.Mobile {
			deviceType = "mobile"
		} else if ua.Tablet {
			deviceType = "tablet"
		}

		// Create click event
		event := &models.ClickEvent{
			ShortCode:  shortCode.(string),
			IPHash:     utils.HashIP(c.IP()),
			DeviceType: deviceType,
			OS:         ua.OS,
			Browser:    ua.Name,
			Referrer:   c.Get("Referer"),
			UserAgent:  c.Get("User-Agent"),
			// Geographic data would be filled by GeoIP service
		}

		// Record click asynchronously to not slow down redirect
		go func() {
			if err := analyticsService.RecordClick(c.Context(), event); err != nil {
				// Log error (in production, use proper logging)
			}
		}()

		return c.Next()
	}
}
