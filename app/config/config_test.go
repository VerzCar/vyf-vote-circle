package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewConfig(t *testing.T) {
	tests := []struct {
		name        string
		configPath  string
		environment string
		expected    *Config
	}{
		{
			name:        "NewConfig with development environment",
			configPath:  ".",
			environment: EnvironmentDev,
			expected: &Config{
				Environment: EnvironmentDev,
				// Fill in the rest of the expected Config fields here
			},
		},
		{
			name:        "NewConfig with production environment",
			configPath:  ".",
			environment: EnvironmentProd,
			expected: &Config{
				Environment: EnvironmentProd,
				// Fill in the rest of the expected Config fields here
			},
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				// Set the environment variable
				os.Setenv("ENVIRONMENT", tt.environment)

				// Call the function
				result := NewConfig(tt.configPath)

				// Assert that the result matches the expected value
				assert.Equal(t, tt.expected, result)

				// Unset the environment variable
				os.Unsetenv("ENVIRONMENT")
			},
		)
	}
}
