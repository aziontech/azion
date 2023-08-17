package create

import (
	"fmt"
	"time"
)

// ParseExpirationDate parses a string representation of an expiration date and returns a time.Time value representing the expiration date.
// If the string is in interval format (e.g. "1d", "2w", "2m", "1y"), it calculates the expiration date from the current date and the specified interval.
// If the string is in full date format (e.g. "18/08/2023", "2023-02-12"), it tries to parse it as a full date and returns the parsed date if it is in the future.
// If the string doesn't match any of the supported formats, it returns an error.
func ParseExpirationDate(currentDate time.Time, expirationString string) (time.Time, error) {
	// Map the suffixes to the corresponding time intervals
	suffixMapping := map[byte]time.Duration{
		'd': 24 * time.Hour,
		'w': 7 * 24 * time.Hour,
		'm': 30 * 24 * time.Hour,  // Using an approximate value for months
		'y': 365 * 24 * time.Hour, // Using an approximate value for years
	}

	// If the string contains a suffix, it is a range format
	lastIndex := len(expirationString) - 1
	lastChar := expirationString[lastIndex]
	if interval, ok := suffixMapping[lastChar]; ok {
		value := expirationString[:lastIndex]
		intervalValue := 0
		fmt.Sscanf(value, "%d", &intervalValue)
		expirationDate := currentDate.Add(time.Duration(intervalValue) * interval)
		return expirationDate, nil
	}

	// Try to analyze as full date (yyyy-mm-dd) or (dd/mm/yyyy)
	formats := []string{"2006-01-02", "02/01/2006"}
	for _, format := range formats {
		if expirationDate, err := time.Parse(format, expirationString); err == nil {
			if expirationDate.After(currentDate) {
				return expirationDate, nil
			}
		}
	}

	return time.Time{}, fmt.Errorf("invalid date format, what do we expect: \"1d\", \"2w\", \"2m\", \"1y\", \"18/08/2023\", \"2023-02-12\"")
}
