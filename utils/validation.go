package utils

import (
	"regexp"
)

func ValidateEmail(email string) bool {
	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	match, _ := regexp.MatchString(pattern, email)
	return match
}

func ValidateUsername(username string) bool {
	return len(username) >= 3 && len(username) <= 30
}

func ValidatePassword(password string) bool {
	return len(password) >= 8
}
