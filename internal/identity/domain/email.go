package domain

import "strings"

type email string

func NewEmail(emailAddr string) (email, error) {
	emailAddr = strings.ToLower(strings.Trim(emailAddr, " "))

	if !strings.Contains(emailAddr, "@") || len(emailAddr) < 5 {
		return "", ErrEmailInvalid
	}
	return email(emailAddr), nil
}

func (e *email) String() string {
	return string(*e)
}
