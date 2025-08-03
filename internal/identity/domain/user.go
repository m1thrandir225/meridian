package domain

import (
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

	events []any
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

	return &user, nil
}

func (u *User) addEvent(event interface{}) {
	u.events = append(u.events, event)
}

func (u *User) Events() []any {
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

func (u *User) UpdateProfile(
	newEmailStr,
	newFirstNameStr,
	newLastNameStr string) error {

	emailChanged := false
	firstNameChanged := false
	lastNameChanged := false

	if newEmailStr != "" {
		oldEmail := u.Email
		newEmail, err := NewEmail(newEmailStr)
		if err != nil {
			return err
		}
		if newEmail != oldEmail {
			u.Email = newEmail
			emailChanged = true
		}
	}
	if newFirstNameStr != "" {
		oldFirstName := u.FirstName
		if newFirstNameStr != oldFirstName {
			u.FirstName = newFirstNameStr
			firstNameChanged = true
		}
	}

	if newLastNameStr != "" {
		oldLastName := u.LastName
		if newLastNameStr != oldLastName {
			u.LastName = newLastNameStr
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

	u.addEvent(UserProfileUpdatedEvent{
		BaseDomainEvent: common.NewBaseDomainEvent("updated_user_profile", u.ID.value, u.Version),
		UserID:          u.ID.String(),
		UpdatedFields:   changedFields,
		Timestamp:       time.Now(),
	})
	return nil
}
