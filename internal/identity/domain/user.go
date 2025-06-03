package domain

import (
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
)

var (
	ErrUserIDInvalid  = errors.New("user ID is invalid")
	ErrEmailInvalid   = errors.New("email format is invalid")
	ErrEmailTaken     = errors.New("email is already taken")
	ErrPasswordPolicy = errors.New("password does not meet policy requirements")
	ErrAuthentication = errors.New("authentication failed: invalid credentials")
)

type UserID struct {
	value uuid.UUID
}

func NewUserID() UserID {
	return UserID{
		value: uuid.New(),
	}
}

func UserIDFromString(input string) (*UserID, error) {
	id, err := uuid.Parse(input)
	if err != nil {
		return nil, err
	}
	return &UserID{
		value: id,
	}, nil
}

type Email string

func NewEmail(emailAddr string) (Email, error) {
	emailAddr = strings.ToLower(strings.Trim(emailAddr, " "))

	if !strings.Contains(emailAddr, "@") || len(emailAddr) < 5 {
		return "", ErrEmailInvalid
	}
	return Email(emailAddr), nil
}

func (e *Email) String() string {
	return string(*e)
}

type User struct {
	id               UserID
	fullName         string
	email            Email
	passwordHash     PasswordHash
	RegistrationTime time.Time
}

func NewUser(emailAddr, rawPasssword, name string) (*User, error) {
	email, err := NewEmail(emailAddr)
	if err != nil {
		return nil, err
	}

	passwordHash, err := NewPasswordHash(rawPasssword)
	if err != nil {
		return nil, err
	}

	id := NewUserID()

	return &User{
		id:               id,
		fullName:         name,
		email:            email,
		passwordHash:     passwordHash,
		RegistrationTime: time.Now(),
	}, nil
}

func (u *User) UpdateProfile(updated any) {
}
