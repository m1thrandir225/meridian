package handlers

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
