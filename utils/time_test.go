package utils

import (
	"testing"
	"time"
)

func TestFormatDuration(t *testing.T) {
	tests := []struct {
		name     string
		duration uint
		expected time.Duration
	}{
		{
			name:     "Test FormatDuration with zero",
			duration: 0,
			expected: 0,
		},
		{
			name:     "Test FormatDuration with positive number",
			duration: 5,
			expected: 5 * time.Second,
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				result := FormatDuration(tt.duration)
				if result != tt.expected {
					t.Errorf("Expected: %v, but got: %v", tt.expected, result)
				}
			},
		)
	}
}
