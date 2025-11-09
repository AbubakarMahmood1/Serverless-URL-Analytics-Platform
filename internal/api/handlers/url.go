package handlers

import (
	"errors"

	"github.com/AbubakarMahmood1/Serverless-URL-Analytics-Platform/internal/models"
	"github.com/AbubakarMahmood1/Serverless-URL-Analytics-Platform/internal/service"
	"github.com/gofiber/fiber/v2"
)

type URLHandler struct {
	urlService       *service.URLService
	analyticsService *service.AnalyticsService
}

func NewURLHandler(urlService *service.URLService, analyticsService *service.AnalyticsService) *URLHandler {
	return &URLHandler{
		urlService:       urlService,
		analyticsService: analyticsService,
	}
}

// ShortenURL handles POST /api/shorten
func (h *URLHandler) ShortenURL(c *fiber.Ctx) error {
	var req models.ShortenRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// In production, extract user ID from JWT or API key
	createdBy := "anonymous"

	response, err := h.urlService.ShortenURL(c.Context(), &req, createdBy)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(response)
}

// RedirectURL handles GET /:shortCode
func (h *URLHandler) RedirectURL(c *fiber.Ctx) error {
	shortCode := c.Params("shortCode")
	if shortCode == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Short code is required",
		})
	}

	// Get original URL
	originalURL, err := h.urlService.GetOriginalURL(c.Context(), shortCode)
	if err != nil {
		if errors.Is(err, errors.New("URL not found")) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "URL not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve URL",
		})
	}

	// Set short code in context for analytics middleware
	c.Locals("short_code", shortCode)

	// Redirect
	return c.Redirect(originalURL, fiber.StatusFound)
}

// GetURLInfo handles GET /api/links/:shortCode
func (h *URLHandler) GetURLInfo(c *fiber.Ctx) error {
	shortCode := c.Params("shortCode")
	if shortCode == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Short code is required",
		})
	}

	info, err := h.urlService.GetURLInfo(c.Context(), shortCode)
	if err != nil {
		if errors.Is(err, errors.New("URL not found")) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "URL not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve URL info",
		})
	}

	return c.JSON(info)
}

// DeleteURL handles DELETE /api/links/:shortCode
func (h *URLHandler) DeleteURL(c *fiber.Ctx) error {
	shortCode := c.Params("shortCode")
	if shortCode == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Short code is required",
		})
	}

	if err := h.urlService.DeleteURL(c.Context(), shortCode); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete URL",
		})
	}

	return c.JSON(fiber.Map{
		"message": "URL deleted successfully",
	})
}

// ListUserURLs handles GET /api/links
func (h *URLHandler) ListUserURLs(c *fiber.Ctx) error {
	// In production, extract user ID from JWT or API key
	userID := "anonymous"

	urls, err := h.urlService.GetUserURLs(c.Context(), userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve URLs",
		})
	}

	return c.JSON(fiber.Map{
		"urls": urls,
	})
}
