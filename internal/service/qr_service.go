package service

import (
	"fmt"

	"github.com/skip2/go-qrcode"
)

type QRService struct {
	baseURL string
}

func NewQRService(baseURL string) *QRService {
	return &QRService{
		baseURL: baseURL,
	}
}

// GenerateQRCode generates a QR code for a short URL
func (s *QRService) GenerateQRCode(shortCode string, size int) ([]byte, error) {
	if size <= 0 {
		size = 256 // Default size
	}

	// Construct the full URL
	url := fmt.Sprintf("%s/%s", s.baseURL, shortCode)

	// Generate QR code
	qrCode, err := qrcode.Encode(url, qrcode.Medium, size)
	if err != nil {
		return nil, fmt.Errorf("failed to generate QR code: %w", err)
	}

	return qrCode, nil
}

// GenerateQRCodeWithLevel generates a QR code with a specific error correction level
func (s *QRService) GenerateQRCodeWithLevel(shortCode string, size int, level qrcode.RecoveryLevel) ([]byte, error) {
	if size <= 0 {
		size = 256
	}

	url := fmt.Sprintf("%s/%s", s.baseURL, shortCode)

	qrCode, err := qrcode.Encode(url, level, size)
	if err != nil {
		return nil, fmt.Errorf("failed to generate QR code: %w", err)
	}

	return qrCode, nil
}
