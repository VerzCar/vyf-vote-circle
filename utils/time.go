package utils

import "time"

// FormatDuration to seconds to use it as a
// standard duration format.
func FormatDuration(duration uint) time.Duration {
	return time.Duration(duration) * time.Second
}
