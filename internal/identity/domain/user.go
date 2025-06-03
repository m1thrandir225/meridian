package domain

import (
	"errors"
	"time"
)

var (
	ErrUserIDInvalid   = errors.New("user ID is invalid")
	ErrUsernameInvalid = errors.New("username is invalid (length, characters)")
	ErrUsernameTaken   = errors.New("username is already taken")
	ErrEmailInvalid    = errors.New("email format is invalid")
	ErrEmailTaken      = errors.New("email is already taken")
	ErrPasswordPolicy  = errors.New("password does not meet policy requirements")
	ErrAuthentication  = errors.New("authentication failed: invalid credentials")
)

type User struct {
	ID               userID
	Username         username
	FirstName        string
	LastName         string
	Email            email
	PasswordHash     passwordHash
	Version          int64
	RegistrationTime time.Time

	events []interface{}
}

func (u *User) Events() []interface{} {
	return u.events
}

func (u *User) ClearEvents() {
	u.events = nil
}

func NewUser(usernameStr, emailStr, firstName, lastName, rawPassword string) (*User, error) {
	email, err := NewEmail(emailStr)
	if err != nil {
		return nil, err
	}

	passwordHash, err := NewPasswordHash(rawPassword)
	if err != nil {
		return nil, err
	}

	username, err := NewUsername(usernameStr)
	if err != nil {
		return nil, err
	}

	id, err := NewUserID()
	if err != nil {
		return nil, err
	}

	user := User{
		ID:               *id,
		Username:         username,
		Email:            email,
		FirstName:        firstName,
		LastName:         lastName,
		PasswordHash:     passwordHash,
		Version:          1,
		RegistrationTime: time.Now(),
	}

	// user.addEvent()
	return &user, nil
}

func (u *User) addEvent(event interface{}) {
	u.events = append(u.events, event)
}

func (u *User) UpdateProfile(updated any) {
}
