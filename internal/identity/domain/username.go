package domain

import (
	"fmt"
	"strings"
	"unicode"
)

type Username string

func NewUsername(input string) (Username, error) {
	name := strings.TrimSpace(input)

	if len(name) < 3 || len(name) > 30 {
		return "", ErrUsernameInvalid
	}
	for _, r := range name {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) {
			return "", fmt.Errorf("%w: contains invalid characters", ErrUsernameInvalid)
		}
	}
	return Username(input), nil
}

func (u Username) String() string {
	return string(u)
}
