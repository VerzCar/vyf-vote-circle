package utils

import "time"

// FormatDuration to seconds to use it as a
// standard duration format.
func FormatDuration(duration uint) time.Duration {
	return time.Duration(duration) * time.Second
}

func IsTimeBetween(t, min, max time.Time) bool {
	if min.After(max) {
		min, max = max, min
	}
	return (t.Equal(min) || t.After(min)) && (t.Equal(max) || t.Before(max))
}
