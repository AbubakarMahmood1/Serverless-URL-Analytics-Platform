package models

import "time"

// URL represents a shortened URL entry
type URL struct {
	ShortCode   string    `json:"short_code" dynamodbav:"short_code"`
	OriginalURL string    `json:"original_url" dynamodbav:"original_url"`
	CustomAlias bool      `json:"custom_alias" dynamodbav:"custom_alias"`
	CreatedAt   int64     `json:"created_at" dynamodbav:"created_at"`
	CreatedBy   string    `json:"created_by,omitempty" dynamodbav:"created_by,omitempty"`
	ExpiresAt   int64     `json:"expires_at,omitempty" dynamodbav:"expires_at,omitempty"`
	IsActive    bool      `json:"is_active" dynamodbav:"is_active"`
}

// ShortenRequest represents the request to shorten a URL
type ShortenRequest struct {
	URL         string `json:"url" validate:"required,url"`
	CustomAlias string `json:"custom_alias,omitempty" validate:"omitempty,alphanum,min=3,max=50"`
}

// ShortenResponse represents the response after shortening a URL
type ShortenResponse struct {
	ShortURL   string `json:"short_url"`
	ShortCode  string `json:"short_code"`
	QRCodeURL  string `json:"qr_code_url"`
	OriginalURL string `json:"original_url"`
	CreatedAt  time.Time `json:"created_at"`
}

// URLInfo represents detailed information about a shortened URL
type URLInfo struct {
	ShortCode   string    `json:"short_code"`
	OriginalURL string    `json:"original_url"`
	CustomAlias bool      `json:"custom_alias"`
	CreatedAt   time.Time `json:"created_at"`
	TotalClicks int       `json:"total_clicks"`
	IsActive    bool      `json:"is_active"`
}
