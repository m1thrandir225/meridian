package domain

import "github.com/google/uuid"

type User struct {
	id        uuid.UUID
	username  string
	firstName string
	lastName  string
	email     string
}

func NewUser(id uuid.UUID, username, firstName, lastName, email string) *User {
	return &User{
		id:        id,
		username:  username,
		firstName: firstName,
		lastName:  lastName,
		email:     email,
	}
}

func (u *User) GetId() uuid.UUID {
	return u.id
}

func (u *User) GetUsername() string {
	return u.username
}

func (u *User) GetFirstName() string {
	return u.firstName
}

func (u *User) GetLastName() string {
	return u.lastName
}

func (u *User) GetEmail() string {
	return u.email
}
