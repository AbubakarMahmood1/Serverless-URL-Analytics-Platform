package service

import (
	"context"
	"time"

	"github.com/AbubakarMahmood1/Serverless-URL-Analytics-Platform/internal/models"
	"github.com/AbubakarMahmood1/Serverless-URL-Analytics-Platform/internal/repository/dynamodb"
	"github.com/AbubakarMahmood1/Serverless-URL-Analytics-Platform/internal/repository/redis"
	"github.com/google/uuid"
)

type AnalyticsService struct {
	analyticsRepo *dynamodb.AnalyticsRepository
	cache         *redis.Cache
}

func NewAnalyticsService(analyticsRepo *dynamodb.AnalyticsRepository, cache *redis.Cache) *AnalyticsService {
	return &AnalyticsService{
		analyticsRepo: analyticsRepo,
		cache:         cache,
	}
}

// RecordClick records a click event
func (s *AnalyticsService) RecordClick(ctx context.Context, event *models.ClickEvent) error {
	// Generate unique click ID
	event.ClickID = uuid.New().String()
	event.Timestamp = time.Now().Unix()

	// Save to database
	if err := s.analyticsRepo.RecordClick(ctx, event); err != nil {
		return err
	}

	// Increment cached counter
	if _, err := s.cache.IncrementClickCount(ctx, event.ShortCode); err != nil {
		// Log but don't fail
		// In production, you might want to use a proper logger
	}

	return nil
}

// GetAnalyticsSummary retrieves aggregated analytics for a URL
func (s *AnalyticsService) GetAnalyticsSummary(ctx context.Context, shortCode string) (*models.AnalyticsSummary, error) {
	// Get all clicks
	clicks, err := s.analyticsRepo.GetClicksByShortCode(ctx, shortCode, 0)
	if err != nil {
		return nil, err
	}

	summary := &models.AnalyticsSummary{
		ShortCode:      shortCode,
		TotalClicks:    len(clicks),
		UniqueVisitors: s.countUniqueVisitors(clicks),
		GeographicData: s.aggregateGeographicData(clicks),
		DeviceBreakdown: s.aggregateDeviceBreakdown(clicks),
		TrafficOverTime: s.aggregateTrafficOverTime(clicks),
		TopReferrers:    s.aggregateTopReferrers(clicks),
	}

	return summary, nil
}

// GetAnalyticsByTimeRange retrieves analytics for a specific time range
func (s *AnalyticsService) GetAnalyticsByTimeRange(ctx context.Context, shortCode string, startTime, endTime time.Time) (*models.AnalyticsSummary, error) {
	clicks, err := s.analyticsRepo.GetClicksByTimeRange(ctx, shortCode, startTime.Unix(), endTime.Unix())
	if err != nil {
		return nil, err
	}

	summary := &models.AnalyticsSummary{
		ShortCode:       shortCode,
		TotalClicks:     len(clicks),
		UniqueVisitors:  s.countUniqueVisitors(clicks),
		GeographicData:  s.aggregateGeographicData(clicks),
		DeviceBreakdown: s.aggregateDeviceBreakdown(clicks),
		TrafficOverTime: s.aggregateTrafficOverTime(clicks),
		TopReferrers:    s.aggregateTopReferrers(clicks),
	}

	return summary, nil
}

// countUniqueVisitors counts unique visitors based on IP hash
func (s *AnalyticsService) countUniqueVisitors(clicks []*models.ClickEvent) int {
	uniqueIPs := make(map[string]bool)
	for _, click := range clicks {
		uniqueIPs[click.IPHash] = true
	}
	return len(uniqueIPs)
}

// aggregateGeographicData aggregates clicks by country and city
func (s *AnalyticsService) aggregateGeographicData(clicks []*models.ClickEvent) []models.GeographicBreakdown {
	geoMap := make(map[string]map[string]int) // country -> city -> count

	for _, click := range clicks {
		if click.Country == "" {
			continue
		}

		if geoMap[click.Country] == nil {
			geoMap[click.Country] = make(map[string]int)
		}

		city := click.City
		if city == "" {
			city = "Unknown"
		}
		geoMap[click.Country][city]++
	}

	var result []models.GeographicBreakdown
	for country, cities := range geoMap {
		for city, count := range cities {
			result = append(result, models.GeographicBreakdown{
				Country: country,
				City:    city,
				Clicks:  count,
			})
		}
	}

	return result
}

// aggregateDeviceBreakdown aggregates clicks by device type, OS, and browser
func (s *AnalyticsService) aggregateDeviceBreakdown(clicks []*models.ClickEvent) models.DeviceBreakdown {
	breakdown := models.DeviceBreakdown{
		OS:      make(map[string]int),
		Browser: make(map[string]int),
	}

	for _, click := range clicks {
		switch click.DeviceType {
		case "desktop":
			breakdown.Desktop++
		case "mobile":
			breakdown.Mobile++
		case "tablet":
			breakdown.Tablet++
		}

		if click.OS != "" {
			breakdown.OS[click.OS]++
		}

		if click.Browser != "" {
			breakdown.Browser[click.Browser]++
		}
	}

	return breakdown
}

// aggregateTrafficOverTime aggregates clicks over time (by hour)
func (s *AnalyticsService) aggregateTrafficOverTime(clicks []*models.ClickEvent) []models.TrafficDataPoint {
	timeMap := make(map[int64]int) // hour timestamp -> count

	for _, click := range clicks {
		// Round to hour
		hourTimestamp := (click.Timestamp / 3600) * 3600
		timeMap[hourTimestamp]++
	}

	var result []models.TrafficDataPoint
	for timestamp, count := range timeMap {
		result = append(result, models.TrafficDataPoint{
			Timestamp: time.Unix(timestamp, 0),
			Clicks:    count,
		})
	}

	return result
}

// aggregateTopReferrers aggregates clicks by referrer
func (s *AnalyticsService) aggregateTopReferrers(clicks []*models.ClickEvent) []models.ReferrerData {
	referrerMap := make(map[string]int)

	for _, click := range clicks {
		referrer := click.Referrer
		if referrer == "" {
			referrer = "Direct"
		}
		referrerMap[referrer]++
	}

	var result []models.ReferrerData
	for referrer, count := range referrerMap {
		result = append(result, models.ReferrerData{
			Referrer: referrer,
			Clicks:   count,
		})
	}

	return result
}
