package domain

import (
	"fmt"
	"strings"
	"unicode"
)

type username string

func NewUsername(input string) (username, error) {
	name := strings.TrimSpace(input)

	if len(name) < 3 || len(name) > 30 {
		return "", ErrUsernameInvalid
	}
	for _, r := range name {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) {
			return "", fmt.Errorf("%w: contains invalid characters", ErrUsernameInvalid)
		}
	}
	return username(input), nil
}

func (u username) String() string {
	return string(u)
}
