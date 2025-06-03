package util

import (
	"fmt"
	"os"
	"time"
)

func VerifyDate(date string) {

	_, err := time.Parse("2006-01-02", date)
	if err != nil {
		fmt.Println("Invalid date format. Use YYYY-MM-DD")
		os.Exit(1)
	}

}
