package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/AbubakarMahmood1/Serverless-URL-Analytics-Platform/config"
	"github.com/AbubakarMahmood1/Serverless-URL-Analytics-Platform/internal/models"
	"github.com/AbubakarMahmood1/Serverless-URL-Analytics-Platform/internal/repository/dynamodb"
	"github.com/AbubakarMahmood1/Serverless-URL-Analytics-Platform/internal/repository/redis"
	"github.com/AbubakarMahmood1/Serverless-URL-Analytics-Platform/internal/utils"
)

type URLService struct {
	urlRepo *dynamodb.URLRepository
	cache   *redis.Cache
	config  *config.Config
}

func NewURLService(urlRepo *dynamodb.URLRepository, cache *redis.Cache, cfg *config.Config) *URLService {
	return &URLService{
		urlRepo: urlRepo,
		cache:   cache,
		config:  cfg,
	}
}

// ShortenURL creates a shortened URL
func (s *URLService) ShortenURL(ctx context.Context, req *models.ShortenRequest, createdBy string) (*models.ShortenResponse, error) {
	// Validate URL
	if !utils.ValidateURL(req.URL) {
		return nil, errors.New("invalid URL")
	}

	// Normalize URL
	normalizedURL := utils.NormalizeURL(req.URL)

	var shortCode string
	var customAlias bool

	// Handle custom alias or generate random code
	if req.CustomAlias != "" {
		if !utils.IsValidCustomAlias(req.CustomAlias) {
			return nil, errors.New("invalid custom alias")
		}
		shortCode = req.CustomAlias
		customAlias = true
	} else {
		// Generate random short code with retry logic
		maxRetries := 5
		for i := 0; i < maxRetries; i++ {
			code, err := utils.GenerateShortCode(s.config.ShortCodeLength)
			if err != nil {
				return nil, fmt.Errorf("failed to generate short code: %w", err)
			}

			// Check if code already exists
			existing, err := s.urlRepo.GetByShortCode(ctx, code)
			if err != nil && !errors.Is(err, dynamodb.ErrURLNotFound) {
				return nil, err
			}

			if existing == nil {
				shortCode = code
				break
			}
		}

		if shortCode == "" {
			return nil, errors.New("failed to generate unique short code")
		}
	}

	// Create URL model
	url := &models.URL{
		ShortCode:   shortCode,
		OriginalURL: normalizedURL,
		CustomAlias: customAlias,
		CreatedAt:   time.Now().Unix(),
		CreatedBy:   createdBy,
		IsActive:    true,
	}

	// Save to database
	if err := s.urlRepo.Create(ctx, url); err != nil {
		if errors.Is(err, dynamodb.ErrURLAlreadyExists) {
			return nil, errors.New("short code already exists, please choose a different alias")
		}
		return nil, fmt.Errorf("failed to create shortened URL: %w", err)
	}

	// Cache the URL
	if err := s.cache.SetURL(ctx, shortCode, normalizedURL, 1*time.Hour); err != nil {
		// Log error but don't fail the request
		fmt.Printf("Warning: failed to cache URL: %v\n", err)
	}

	// Build response
	shortURL := fmt.Sprintf("%s/%s", s.config.BaseURL, shortCode)
	qrCodeURL := fmt.Sprintf("%s/api/qr/%s", s.config.BaseURL, shortCode)

	return &models.ShortenResponse{
		ShortURL:    shortURL,
		ShortCode:   shortCode,
		QRCodeURL:   qrCodeURL,
		OriginalURL: normalizedURL,
		CreatedAt:   time.Unix(url.CreatedAt, 0),
	}, nil
}

// GetOriginalURL retrieves the original URL from a short code
func (s *URLService) GetOriginalURL(ctx context.Context, shortCode string) (string, error) {
	// Check cache first
	cachedURL, err := s.cache.GetURL(ctx, shortCode)
	if err != nil {
		fmt.Printf("Warning: cache error: %v\n", err)
	}
	if cachedURL != "" {
		return cachedURL, nil
	}

	// Fetch from database
	url, err := s.urlRepo.GetByShortCode(ctx, shortCode)
	if err != nil {
		if errors.Is(err, dynamodb.ErrURLNotFound) {
			return "", errors.New("URL not found")
		}
		return "", err
	}

	// Check if URL is active
	if !url.IsActive {
		return "", errors.New("URL has been deactivated")
	}

	// Check expiration
	if url.ExpiresAt > 0 && time.Now().Unix() > url.ExpiresAt {
		return "", errors.New("URL has expired")
	}

	// Cache the URL
	if err := s.cache.SetURL(ctx, shortCode, url.OriginalURL, 1*time.Hour); err != nil {
		fmt.Printf("Warning: failed to cache URL: %v\n", err)
	}

	return url.OriginalURL, nil
}

// GetURLInfo retrieves detailed information about a shortened URL
func (s *URLService) GetURLInfo(ctx context.Context, shortCode string) (*models.URLInfo, error) {
	url, err := s.urlRepo.GetByShortCode(ctx, shortCode)
	if err != nil {
		if errors.Is(err, dynamodb.ErrURLNotFound) {
			return nil, errors.New("URL not found")
		}
		return nil, err
	}

	// Get click count from cache
	clickCount, err := s.cache.GetClickCount(ctx, shortCode)
	if err != nil {
		fmt.Printf("Warning: failed to get click count: %v\n", err)
		clickCount = 0
	}

	return &models.URLInfo{
		ShortCode:   url.ShortCode,
		OriginalURL: url.OriginalURL,
		CustomAlias: url.CustomAlias,
		CreatedAt:   time.Unix(url.CreatedAt, 0),
		TotalClicks: int(clickCount),
		IsActive:    url.IsActive,
	}, nil
}

// DeleteURL deactivates a shortened URL
func (s *URLService) DeleteURL(ctx context.Context, shortCode string) error {
	if err := s.urlRepo.Delete(ctx, shortCode); err != nil {
		return err
	}

	// Remove from cache
	if err := s.cache.DeleteURL(ctx, shortCode); err != nil {
		fmt.Printf("Warning: failed to delete from cache: %v\n", err)
	}

	return nil
}

// GetUserURLs retrieves all URLs created by a user
func (s *URLService) GetUserURLs(ctx context.Context, userID string) ([]*models.URLInfo, error) {
	urls, err := s.urlRepo.GetByCreatedBy(ctx, userID)
	if err != nil {
		return nil, err
	}

	urlInfos := make([]*models.URLInfo, 0, len(urls))
	for _, url := range urls {
		// Get click count
		clickCount, _ := s.cache.GetClickCount(ctx, url.ShortCode)

		urlInfos = append(urlInfos, &models.URLInfo{
			ShortCode:   url.ShortCode,
			OriginalURL: url.OriginalURL,
			CustomAlias: url.CustomAlias,
			CreatedAt:   time.Unix(url.CreatedAt, 0),
			TotalClicks: int(clickCount),
			IsActive:    url.IsActive,
		})
	}

	return urlInfos, nil
}
