package domain

import (
	"errors"
	"time"
)

var ErrInvalidDateRange = errors.New("end date must be on or after start date")

func dateOnly(value time.Time) time.Time {
	return time.Date(value.Year(), value.Month(), value.Day(), 0, 0, 0, 0, time.UTC)
}

func NextOccurrence(original, fromDate time.Time, isYearly bool) time.Time {
	original = dateOnly(original)
	fromDate = dateOnly(fromDate)
	if !isYearly {
		return original
	}
	day := original.Day()
	maxDay := daysInMonth(fromDate.Year(), original.Month())
	if day > maxDay {
		day = maxDay
	}
	candidate := time.Date(fromDate.Year(), original.Month(), day, 0, 0, 0, 0, time.UTC)
	if candidate.Before(fromDate) {
		nextYear := fromDate.Year() + 1
		day = original.Day()
		maxDay = daysInMonth(nextYear, original.Month())
		if day > maxDay {
			day = maxDay
		}
		candidate = time.Date(nextYear, original.Month(), day, 0, 0, 0, 0, time.UTC)
	}
	return candidate
}

func daysInMonth(year int, month time.Month) int {
	return time.Date(year, month+1, 0, 0, 0, 0, 0, time.UTC).Day()
}
