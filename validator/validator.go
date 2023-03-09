package validator

import (
	"regexp"
	"strings"
)

func IsValidInput(city, state string) bool {
	// Ensure that city and state are not empty
	if strings.TrimSpace(city) == "" || strings.TrimSpace(state) == "" {
		return false
	}

	// Ensure that city and state contain only alphabetic characters and spaces
	if !regexp.MustCompile(`^[A-Za-z ]+$`).MatchString(city) || !regexp.MustCompile(`^[A-Za-z ]+$`).MatchString(state) {
		return false
	}

	return true
}