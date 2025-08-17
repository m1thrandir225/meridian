package domain

import "errors"

type UserIDRef string

func (u UserIDRef) String() string {
	return string(u)
}

func NewUserIDRef(id string) (UserIDRef, error) {
	if id == "" {
		return "", errors.New("user ID cannot be empty")
	}
	return UserIDRef(id), nil
}
