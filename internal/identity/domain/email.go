package domain

import "strings"

type UserEmail string

func NewEmail(emailAddr string) (UserEmail, error) {
	emailAddr = strings.ToLower(strings.Trim(emailAddr, " "))

	if !strings.Contains(emailAddr, "@") || len(emailAddr) < 5 {
		return "", ErrEmailInvalid
	}
	return UserEmail(emailAddr), nil
}

func (e *UserEmail) String() string {
	return string(*e)
}
