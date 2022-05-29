package utils

import (
	"github.com/gofrs/uuid"
	"time"
)

// UniqueKey creates a unique key with current time-date stamp
// and uuid4.
// Returns the time with appended uuid as string
func UniqueKey() string {
	currentTime := time.Now()
	formattedTime := currentTime.Format(time.RFC3339)

	uID, _ := uuid.NewV4()

	return formattedTime + uID.String()
}
