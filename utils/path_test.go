package utils

import (
	"path/filepath"
	"testing"
)

func TestBase(t *testing.T) {
	tests := []struct {
		name     string
		expected string
	}{
		{
			name:     "Test Base",
			expected: filepath.Dir(basePath) + "/",
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				result := Base()
				if result != tt.expected {
					t.Errorf("Expected: %v, but got: %v", tt.expected, result)
				}
			},
		)
	}
}

func TestFromBase(t *testing.T) {
	tests := []struct {
		name     string
		addPath  string
		expected string
	}{
		{
			name:     "Test FromBase with empty string",
			addPath:  "",
			expected: Base(),
		},
		{
			name:     "Test FromBase with non-empty string",
			addPath:  "test",
			expected: Base() + "test",
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				result := FromBase(tt.addPath)
				if result != tt.expected {
					t.Errorf("Expected: %v, but got: %v", tt.expected, result)
				}
			},
		)
	}
}
