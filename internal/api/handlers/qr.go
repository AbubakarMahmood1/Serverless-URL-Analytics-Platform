package handlers

import (
	"strconv"

	"github.com/AbubakarMahmood1/Serverless-URL-Analytics-Platform/internal/service"
	"github.com/gofiber/fiber/v2"
)

type QRHandler struct {
	qrService *service.QRService
}

func NewQRHandler(qrService *service.QRService) *QRHandler {
	return &QRHandler{
		qrService: qrService,
	}
}

// GenerateQR handles GET /api/qr/:shortCode
func (h *QRHandler) GenerateQR(c *fiber.Ctx) error {
	shortCode := c.Params("shortCode")
	if shortCode == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Short code is required",
		})
	}

	// Get size from query parameter (default: 256)
	size := 256
	if sizeStr := c.Query("size"); sizeStr != "" {
		if parsedSize, err := strconv.Atoi(sizeStr); err == nil && parsedSize > 0 {
			size = parsedSize
		}
	}

	// Generate QR code
	qrCode, err := h.qrService.GenerateQRCode(shortCode, size)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to generate QR code",
		})
	}

	// Set content type and return image
	c.Set("Content-Type", "image/png")
	c.Set("Cache-Control", "public, max-age=86400") // Cache for 24 hours

	return c.Send(qrCode)
}
