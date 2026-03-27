package util

import (
	"regexp"
	"strings"
)

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

func IsValidEmail(email string) bool {
	return emailRegex.MatchString(email)
}

func IsNotBlank(s string) bool {
	return strings.TrimSpace(s) != ""
}

func IsMinLength(s string, min int) bool {
	return len(strings.TrimSpace(s)) >= min
}