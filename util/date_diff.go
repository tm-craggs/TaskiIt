package util

import (
	"fmt"
	"time"
)

// dateDiff returns a string like "3 days" or "1 day" from the difference between two dates
func dateDiff(a, b time.Time) string {
	a = a.Truncate(24 * time.Hour)
	b = b.Truncate(24 * time.Hour)
	diff := int(b.Sub(a).Hours() / 24)
	if diff < 0 {
		diff = -diff
	}
	if diff == 1 {
		return "1 day"
	}
	return fmt.Sprintf("%d days", diff)
}
