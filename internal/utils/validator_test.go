package utils

import (
	"testing"
)

func TestValidateURL(t *testing.T) {
	tests := []struct {
		name string
		url  string
		want bool
	}{
		{"Valid HTTP URL", "http://example.com", true},
		{"Valid HTTPS URL", "https://example.com", true},
		{"Valid URL with path", "https://example.com/path/to/page", true},
		{"Valid URL with query", "https://example.com?query=value", true},
		{"Valid URL with fragment", "https://example.com#section", true},
		{"Invalid scheme", "ftp://example.com", false},
		{"No scheme", "example.com", false},
		{"Localhost", "http://localhost", false},
		{"127.0.0.1", "http://127.0.0.1", false},
		{"Private IP", "http://192.168.1.1", false},
		{"Empty string", "", false},
		{"Just scheme", "http://", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ValidateURL(tt.url); got != tt.want {
				t.Errorf("ValidateURL(%s) = %v, want %v", tt.url, got, tt.want)
			}
		})
	}
}

func TestHashIP(t *testing.T) {
	tests := []struct {
		name string
		ip   string
	}{
		{"IPv4 address", "192.168.1.1"},
		{"IPv6 address", "2001:0db8:85a3:0000:0000:8a2e:0370:7334"},
		{"Localhost", "127.0.0.1"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hash1 := HashIP(tt.ip)
			hash2 := HashIP(tt.ip)

			// Same IP should produce same hash
			if hash1 != hash2 {
				t.Errorf("HashIP() not consistent for %s", tt.ip)
			}

			// Hash should be non-empty
			if hash1 == "" {
				t.Errorf("HashIP() returned empty string")
			}

			// Hash should be different from original
			if hash1 == tt.ip {
				t.Errorf("HashIP() returned unhashed IP")
			}

			// Hash should be hex string
			if len(hash1) != 64 { // SHA-256 produces 64 character hex string
				t.Errorf("HashIP() length = %d, want 64", len(hash1))
			}
		})
	}
}

func TestHashIPUniqueness(t *testing.T) {
	ip1 := "192.168.1.1"
	ip2 := "192.168.1.2"

	hash1 := HashIP(ip1)
	hash2 := HashIP(ip2)

	if hash1 == hash2 {
		t.Errorf("Different IPs produced same hash")
	}
}

func TestNormalizeURL(t *testing.T) {
	tests := []struct {
		name string
		url  string
		want string
	}{
		{
			"Remove trailing slash",
			"https://example.com/",
			"https://example.com",
		},
		{
			"Remove fragment",
			"https://example.com/page#section",
			"https://example.com/page",
		},
		{
			"Remove both",
			"https://example.com/page/#section",
			"https://example.com/page",
		},
		{
			"Already normalized",
			"https://example.com/page",
			"https://example.com/page",
		},
		{
			"Keep query parameters",
			"https://example.com/page?query=value",
			"https://example.com/page?query=value",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NormalizeURL(tt.url); got != tt.want {
				t.Errorf("NormalizeURL() = %v, want %v", got, tt.want)
			}
		})
	}
}

func BenchmarkValidateURL(b *testing.B) {
	url := "https://example.com/path/to/page"
	for i := 0; i < b.N; i++ {
		ValidateURL(url)
	}
}

func BenchmarkHashIP(b *testing.B) {
	ip := "192.168.1.1"
	for i := 0; i < b.N; i++ {
		HashIP(ip)
	}
}
