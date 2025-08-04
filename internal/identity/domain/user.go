package domain

import (
	"fmt"
	"github.com/m1thrandir225/meridian/pkg/common"
	"time"
)

type User struct {
	ID               UserID
	Username         Username
	FirstName        string
	LastName         string
	Email            UserEmail
	PasswordHash     PasswordHash
	Version          int64
	RegistrationTime time.Time

	events []common.DomainEvent
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
		events:           make([]common.DomainEvent, 0),
	}

	event := CreateUserRegisteredEvent(&user)

	user.addEvent(event)

	return &user, nil
}

func (u *User) addEvent(event common.DomainEvent) {
	u.events = append(u.events, event)
}

func (u *User) Events() []common.DomainEvent {
	return u.events
}

func (u *User) ClearEvents() {
	u.events = nil
}

func (u *User) Authenticate(rawPassword string) error {
	if !u.PasswordHash.Match(rawPassword) {
		return ErrPasswordPolicy
	}
	return nil
}

func (u *User) UpdatePassword(newPassword string) error {
	if u.PasswordHash.Match(newPassword) {
		return fmt.Errorf("new password is the same as the old one")
	}
	newPasswordHash, err := NewPasswordHash(newPassword)
	if err != nil {
		return err
	}
	u.PasswordHash = newPasswordHash
	u.Version++

	event := CreateUserProfileUpdated(u, map[string]any{"password": "REDACTED"})
	u.addEvent(event)

	return nil
}

func (u *User) UpdateProfile(newUsernameStr, newEmailStr, newFirstNameStr, newLastNameStr *string) error {
	emailChanged := false
	firstNameChanged := false
	lastNameChanged := false
	usernameChanged := false

	if newUsernameStr != nil {
		oldUsername := u.Username
		newUsername, err := NewUsername(*newUsernameStr)
		if err != nil {
			return err
		}
		if newUsername != oldUsername {
			u.Username = newUsername
			usernameChanged = true
		}
	}

	if newEmailStr != nil {
		oldEmail := u.Email
		newEmail, err := NewEmail(*newEmailStr)
		if err != nil {
			return err
		}
		if newEmail != oldEmail {
			u.Email = newEmail
			emailChanged = true
		}
	}
	if newFirstNameStr != nil {
		oldFirstName := u.FirstName
		if *newFirstNameStr != oldFirstName {
			u.FirstName = *newFirstNameStr
			firstNameChanged = true
		}
	}

	if newLastNameStr != nil {
		oldLastName := u.LastName
		if *newLastNameStr != oldLastName {
			u.LastName = *newLastNameStr
			lastNameChanged = true
		}
	}

	changedFields := make(map[string]any)
	if emailChanged {
		changedFields["email"] = u.Email
	}
	if firstNameChanged {
		changedFields["first_name"] = u.FirstName
	}
	if lastNameChanged {
		changedFields["last_name"] = u.LastName
	}
	if usernameChanged {
		changedFields["username"] = u.Username
	}

	u.Version++

	event := CreateUserProfileUpdated(u, changedFields)
	u.addEvent(event)

	return nil
}
