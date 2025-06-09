package util

import (
	"fmt"
	"time"
)

// VerifyDate checks if the input string is in YYYY-MM-DD format.
// Returns an error if the format is invalid.
func VerifyDate(date string) error {
	_, err := time.Parse("2006-01-02", date)
	if err != nil {
		return fmt.Errorf("failed to parse date. %w", err)
	}
	return nil
}
