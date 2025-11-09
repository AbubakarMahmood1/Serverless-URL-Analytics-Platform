package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"net"
	"net/url"
	"strings"
)

// ValidateURL checks if a URL is valid and safe
func ValidateURL(urlStr string) bool {
	// Parse the URL
	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		return false
	}

	// Check scheme
	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		return false
	}

	// Check if host is present
	if parsedURL.Host == "" {
		return false
	}

	// Basic blacklist check for localhost and private IPs
	if isLocalOrPrivate(parsedURL.Host) {
		return false
	}

	return true
}

// isLocalOrPrivate checks if a host is localhost or a private IP
func isLocalOrPrivate(host string) bool {
	// Remove port if present
	hostWithoutPort := host
	if strings.Contains(host, ":") {
		hostWithoutPort, _, _ = net.SplitHostPort(host)
	}

	// Check for localhost
	if hostWithoutPort == "localhost" || hostWithoutPort == "127.0.0.1" || hostWithoutPort == "::1" {
		return true
	}

	// Parse as IP
	ip := net.ParseIP(hostWithoutPort)
	if ip == nil {
		return false
	}

	// Check if private IP
	return ip.IsPrivate() || ip.IsLoopback()
}

// HashIP creates a SHA-256 hash of an IP address for privacy
func HashIP(ip string) string {
	hash := sha256.Sum256([]byte(ip))
	return hex.EncodeToString(hash[:])
}

// NormalizeURL normalizes a URL by removing trailing slashes and fragments
func NormalizeURL(urlStr string) string {
	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		return urlStr
	}

	// Remove fragment
	parsedURL.Fragment = ""

	// Remove trailing slash from path
	parsedURL.Path = strings.TrimSuffix(parsedURL.Path, "/")

	return parsedURL.String()
}
