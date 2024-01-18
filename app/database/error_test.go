package database

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestRecordNotFound(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected bool
	}{
		{
			name:     "RecordNotFound with gorm.ErrRecordNotFound",
			err:      gorm.ErrRecordNotFound,
			expected: true,
		},
		{
			name:     "RecordNotFound with other error",
			err:      errors.New("other error"),
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				result := RecordNotFound(tt.err)
				assert.Equal(t, tt.expected, result)
			},
		)
	}
}
