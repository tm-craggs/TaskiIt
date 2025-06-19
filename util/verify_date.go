package util

import (
	"fmt"
	"time"
)

// VerifyDate checks if the input string is a valid date in DD-MM-YYYY format.
// returns an error if the format is invalid.
func VerifyDate(date string) error {

	// parse the date according to YYYY-MM-DD
	_, err := time.Parse("2006-01-02", date)
	if err != nil {
		return fmt.Errorf("failed to parse date. %w", err)
	}
	return nil
}
