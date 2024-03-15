package database

import (
	"testing"

	"github.com/VerzCar/vyf-vote-circle/app/config"
	"github.com/stretchr/testify/assert"
)

func TestDsn(t *testing.T) {
	tests := []struct {
		name     string
		conf     *config.Config
		expected string
	}{
		{
			name: "dsn with production environment",
			conf: &config.Config{
				Environment: config.EnvironmentProd,
				Db: struct {
					Host      string
					Port      uint16
					Name      string
					User      string
					Password  string
					Migration bool
					Test      struct {
						Host     string
						Port     uint16
						Name     string
						User     string
						Password string
					}
				}{
					Host:      "localhost",
					Port:      uint16(5432),
					Name:      "test",
					User:      "test",
					Password:  "test",
					Migration: false,
					Test: struct {
						Host     string
						Port     uint16
						Name     string
						User     string
						Password string
					}{
						Host:     "localhost",
						Port:     uint16(5432),
						Name:     "test",
						User:     "test",
						Password: "test",
					},
				},
			},
			expected: "host=localhost port=5432 user=test dbname=test password=test sslmode=require",
		},
		{
			name: "dsn with non-production environment",
			conf: &config.Config{
				Environment: config.EnvironmentDev,
				Db: struct {
					Host      string
					Port      uint16
					Name      string
					User      string
					Password  string
					Migration bool
					Test      struct {
						Host     string
						Port     uint16
						Name     string
						User     string
						Password string
					}
				}{
					Host:      "localhost",
					Port:      uint16(5432),
					Name:      "test",
					User:      "test",
					Password:  "test",
					Migration: false,
					Test: struct {
						Host     string
						Port     uint16
						Name     string
						User     string
						Password string
					}{
						Host:     "localhost",
						Port:     uint16(5432),
						Name:     "test",
						User:     "test",
						Password: "test",
					},
				},
			},
			expected: "host=localhost port=5432 user=test dbname=test password=test sslmode=disable",
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				result := dsn(tt.conf)
				assert.Equal(t, tt.expected, result)
			},
		)
	}
}
