package util

import (
	"fmt"
	"time"
)

// dateDiff calculates the difference between two dates
// it returns a human-readable string showing the number of days between them
func dateDiff(a, b time.Time) string {

	// ignore time of day by truncating both to midnight
	a = a.Truncate(24 * time.Hour)
	b = b.Truncate(24 * time.Hour)

	// calculate the difference in hours, then convert to days
	diff := int(b.Sub(a).Hours() / 24)

	// take the absolute value to ensure positive day count
	if diff < 0 {
		diff = -diff
	}

	// handle single/plural form
	if diff == 1 {
		return "1 day"
	}

	return fmt.Sprintf("%d days", diff)
}
