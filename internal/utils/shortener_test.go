package utils

import (
	"testing"
)

func TestGenerateShortCode(t *testing.T) {
	tests := []struct {
		name   string
		length int
		want   int
	}{
		{"Generate 6 character code", 6, 6},
		{"Generate 8 character code", 8, 8},
		{"Generate 10 character code", 10, 10},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GenerateShortCode(tt.length)
			if err != nil {
				t.Errorf("GenerateShortCode() error = %v", err)
				return
			}
			if len(got) != tt.want {
				t.Errorf("GenerateShortCode() length = %v, want %v", len(got), tt.want)
			}

			// Check if all characters are valid
			for _, char := range got {
				if !isValidCharacter(char) {
					t.Errorf("GenerateShortCode() contains invalid character: %c", char)
				}
			}
		})
	}
}

func TestGenerateShortCodeUniqueness(t *testing.T) {
	codes := make(map[string]bool)
	iterations := 1000

	for i := 0; i < iterations; i++ {
		code, err := GenerateShortCode(6)
		if err != nil {
			t.Errorf("GenerateShortCode() error = %v", err)
			return
		}
		if codes[code] {
			t.Errorf("GenerateShortCode() generated duplicate: %s", code)
		}
		codes[code] = true
	}
}

func TestIsValidCustomAlias(t *testing.T) {
	tests := []struct {
		name  string
		alias string
		want  bool
	}{
		{"Valid alias", "mylink", true},
		{"Valid alias with numbers", "test123", true},
		{"Valid alias with hyphen", "my-link", true},
		{"Valid alias with underscore", "my_link", true},
		{"Too short", "ab", false},
		{"Too long", "this-is-a-very-long-alias-that-exceeds-the-maximum-allowed-length", false},
		{"Invalid characters", "my link", false},
		{"Invalid characters special", "my@link", false},
		{"Empty string", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsValidCustomAlias(tt.alias); got != tt.want {
				t.Errorf("IsValidCustomAlias() = %v, want %v", got, tt.want)
			}
		})
	}
}

func isValidCharacter(char rune) bool {
	return (char >= 'a' && char <= 'z') ||
		(char >= 'A' && char <= 'Z') ||
		(char >= '0' && char <= '9')
}

func BenchmarkGenerateShortCode(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GenerateShortCode(6)
	}
}
