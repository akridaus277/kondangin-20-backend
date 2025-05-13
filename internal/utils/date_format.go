package utils

import (
	"time"
)

// ParseDate parses string in "YYYY-MM-DD" format
func ParseDate(dateStr string, dateLayout string) (time.Time, error) {
	return time.Parse(dateLayout, dateStr)
}

func IsValidDateFormat(dateStr string, dateLayout string) bool {
	_, err := ParseDate(dateStr, dateLayout)
	return err == nil
}
