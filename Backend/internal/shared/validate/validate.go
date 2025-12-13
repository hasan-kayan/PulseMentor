package validate

import (
	"regexp"
	"strings"
)

var (
	emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
)

func Email(email string) bool {
	if email == "" {
		return false
	}
	return emailRegex.MatchString(email)
}

func Password(password string) bool {
	if len(password) < 8 {
		return false
	}
	return true
}

func NonEmpty(s string) bool {
	return strings.TrimSpace(s) != ""
}

