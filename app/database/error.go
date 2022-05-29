package database

import (
	"errors"
	"gorm.io/gorm"
)

// RecordNotFound checks whether the error is from
// gorm.ErrRecordNotFound and if match it returns true, otherwise
// false.
func RecordNotFound(err error) bool {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return true
	}
	return false
}
