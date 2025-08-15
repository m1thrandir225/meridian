package handlers

type UserResponse struct {
	ID        string `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type AuthenticateTokensResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int64  `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
}

type RegisterUserRequest struct {
	Username  string `json:"username" binding:"required"`
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required"`
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
}

type AuthenticateUserRequest struct {
	Login    string `json:"login" binding:"required"` // Username or Email
	Password string `json:"password" binding:"required"`
}

type AuthenticateUserResponse struct {
	User   UserResponse               `json:"user"`
	Tokens AuthenticateTokensResponse `json:"tokens"`
}

type UpdateProfileRequest struct {
	Username  *string `json:"username" binding:"omitempty,max=255"`
	FirstName *string `json:"first_name" binding:"omitempty,max=255"`
	LastName  *string `json:"last_name" binding:"omitempty,max=255"`
	Email     *string `json:"email" binding:"omitempty,email"`
}

type UpdatePasswordRequest struct {
	NewPassword string `json:"new_password" binding:"required"`
}
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}
