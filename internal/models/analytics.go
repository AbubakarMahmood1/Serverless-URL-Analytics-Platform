package models

import "time"

// ClickEvent represents a single click event on a shortened URL
type ClickEvent struct {
	ShortCode  string `json:"short_code" dynamodbav:"short_code"`
	Timestamp  int64  `json:"timestamp" dynamodbav:"timestamp"`
	ClickID    string `json:"click_id" dynamodbav:"click_id"`
	IPHash     string `json:"ip_hash" dynamodbav:"ip_hash"`
	Country    string `json:"country,omitempty" dynamodbav:"country,omitempty"`
	City       string `json:"city,omitempty" dynamodbav:"city,omitempty"`
	DeviceType string `json:"device_type,omitempty" dynamodbav:"device_type,omitempty"` // mobile, desktop, tablet
	OS         string `json:"os,omitempty" dynamodbav:"os,omitempty"`
	Browser    string `json:"browser,omitempty" dynamodbav:"browser,omitempty"`
	Referrer   string `json:"referrer,omitempty" dynamodbav:"referrer,omitempty"`
	UserAgent  string `json:"user_agent,omitempty" dynamodbav:"user_agent,omitempty"`
}

// AnalyticsSummary represents aggregated analytics for a URL
type AnalyticsSummary struct {
	ShortCode       string                 `json:"short_code"`
	TotalClicks     int                    `json:"total_clicks"`
	UniqueVisitors  int                    `json:"unique_visitors"`
	GeographicData  []GeographicBreakdown  `json:"geographic_data"`
	DeviceBreakdown DeviceBreakdown        `json:"device_breakdown"`
	TrafficOverTime []TrafficDataPoint     `json:"traffic_over_time"`
	TopReferrers    []ReferrerData         `json:"top_referrers"`
}

// GeographicBreakdown represents clicks by geographic location
type GeographicBreakdown struct {
	Country string `json:"country"`
	City    string `json:"city,omitempty"`
	Clicks  int    `json:"clicks"`
}

// DeviceBreakdown represents clicks by device type
type DeviceBreakdown struct {
	Desktop int            `json:"desktop"`
	Mobile  int            `json:"mobile"`
	Tablet  int            `json:"tablet"`
	OS      map[string]int `json:"os"`
	Browser map[string]int `json:"browser"`
}

// TrafficDataPoint represents clicks at a specific time
type TrafficDataPoint struct {
	Timestamp time.Time `json:"timestamp"`
	Clicks    int       `json:"clicks"`
}

// ReferrerData represents clicks from a specific referrer
type ReferrerData struct {
	Referrer string `json:"referrer"`
	Clicks   int    `json:"clicks"`
}
