package domain

type User struct {
	ID        string
	Username  string
	Email     string
	FirstName string
	LastName  string
}

func NewUser(id, username, email, firstName, lastName string) *User {
	return &User{
		ID:        id,
		Username:  username,
		Email:     email,
		FirstName: firstName,
		LastName:  lastName,
	}
}
