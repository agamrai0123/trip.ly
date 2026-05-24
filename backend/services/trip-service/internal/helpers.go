package internal

import (
	"fmt"
	"time"
)

// coalesce returns s if non-empty, else fallback.
func coalesce(s, fallback string) string {
	if s == "" {
		return fallback
	}
	return s
}

// parseDate parses a date string in RFC3339 or date-only format.
func parseDate(s string) (time.Time, error) {
	if t, err := time.Parse(time.RFC3339, s); err == nil {
		return t, nil
	}
	if t, err := time.Parse("2006-01-02", s); err == nil {
		return t, nil
	}
	return time.Time{}, fmt.Errorf("unrecognised date format: %s", s)
}
