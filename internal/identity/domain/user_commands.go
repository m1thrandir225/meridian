package domain

type RegisterUserCommand struct {
	Username  string
	Email     string
	FirstName string
	LastName  string
	Password  string
}

type AuthenticateUserCommand struct {
	LoginIdentifier string
	Password        string
}

type GetUserCommand struct {
	UserID string
}

type UpdateUserProfileCommand struct {
	UserID       string
	NewEmail     *string
	NewUsername  *string
	NewFirstName *string
	NewLastName  *string
}

type UpdateUserPasswordCommand struct {
	UserID      string
	NewPassword string
}

type DeleteUserCommand struct {
	UserID string
}
