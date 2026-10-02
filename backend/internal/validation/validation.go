package validation

import (
	"errors"
	"strings"
	"time"
)

var (
	ErrRequired     = errors.New("required field is missing")
	ErrInvalidDate  = errors.New("invalid date")
	ErrInvalidTime  = errors.New("invalid time")
	ErrInvalidEmail = errors.New("invalid email")
	ErrInvalidSlot  = errors.New("invalid appointment slot")
	ErrInvalidID    = errors.New("invalid ID")
)

func Required(value string) bool {
	return strings.TrimSpace(value) != ""
}

func Email(value string) bool {
	value = strings.TrimSpace(value)

	return strings.Contains(value, "@") &&
		strings.Contains(value, ".")
}

func Date(value string) (time.Time, error) {
	return time.Parse("2006-01-02", value)
}

func Time(value string) (time.Time, error) {
	return time.Parse("15:04", value)
}

func Is30MinuteSlot(value string) bool {
	t, err := Time(value)
	if err != nil {
		return false
	}

	return t.Minute() == 0 || t.Minute() == 30
}
