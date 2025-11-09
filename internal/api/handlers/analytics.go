package handlers

import (
	"errors"
	"time"

	"github.com/AbubakarMahmood1/Serverless-URL-Analytics-Platform/internal/service"
	"github.com/gofiber/fiber/v2"
)

type AnalyticsHandler struct {
	analyticsService *service.AnalyticsService
}

func NewAnalyticsHandler(analyticsService *service.AnalyticsService) *AnalyticsHandler {
	return &AnalyticsHandler{
		analyticsService: analyticsService,
	}
}

// GetAnalytics handles GET /api/analytics/:shortCode
func (h *AnalyticsHandler) GetAnalytics(c *fiber.Ctx) error {
	shortCode := c.Params("shortCode")
	if shortCode == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Short code is required",
		})
	}

	// Check for time range query parameters
	startTimeStr := c.Query("start_time")
	endTimeStr := c.Query("end_time")

	var summary interface{}
	var err error

	if startTimeStr != "" && endTimeStr != "" {
		// Parse time range
		startTime, err1 := time.Parse(time.RFC3339, startTimeStr)
		endTime, err2 := time.Parse(time.RFC3339, endTimeStr)

		if err1 != nil || err2 != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid time format. Use RFC3339 format (e.g., 2024-01-01T00:00:00Z)",
			})
		}

		summary, err = h.analyticsService.GetAnalyticsByTimeRange(c.Context(), shortCode, startTime, endTime)
	} else {
		// Get all analytics
		summary, err = h.analyticsService.GetAnalyticsSummary(c.Context(), shortCode)
	}

	if err != nil {
		if errors.Is(err, errors.New("URL not found")) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Analytics not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve analytics",
		})
	}

	return c.JSON(summary)
}
