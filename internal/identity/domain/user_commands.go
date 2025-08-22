package domain

type Command interface {
	CommandName() string
}

type RegisterUserCommand struct {
	Username  string
	Email     string
	FirstName string
	LastName  string
	Password  string
}

func (c RegisterUserCommand) CommandName() string {
	return "RegisterUser"
}

type AuthenticateUserCommand struct {
	LoginIdentifier string
	Password        string
	IPAddress       string
	Device          string
}

func (c AuthenticateUserCommand) CommandName() string {
	return "AuthenticateUser"
}

type GetUserCommand struct {
	UserID string
}

func (c GetUserCommand) CommandName() string {
	return "GetUser"
}

type GetUsersCommand struct {
	UserIds []string
}

func (c GetUsersCommand) CommandName() string {
	return "GetUsers"
}

type UpdateUserProfileCommand struct {
	UserID       string
	NewEmail     *string
	NewUsername  *string
	NewFirstName *string
	NewLastName  *string
}

func (c UpdateUserProfileCommand) CommandName() string {
	return "UpdateUserProfile"
}

type UpdateUserPasswordCommand struct {
	UserID      string
	NewPassword string
}

func (c UpdateUserPasswordCommand) CommandName() string {
	return "UpdateUserPassword"
}

type DeleteUserCommand struct {
	UserID string
}

func (c DeleteUserCommand) CommandName() string {
	return "DeleteUser"
}

type RefreshTokenCommand struct {
	RawRefreshToken string
	Device          string
	IPAddress       string
}

func (c RefreshTokenCommand) CommandName() string {
	return "RefreshToken"
}

type RevokeTokenCommand struct {
	Token string
}

func (c RevokeTokenCommand) CommandName() string {
	return "RevokeToken"
}

type RevokeAllTokensCommand struct {
	UserID string
}

func (c RevokeAllTokensCommand) CommandName() string {
	return "RevokeAllTokens"
}
