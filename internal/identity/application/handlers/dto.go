package handlers

import "github.com/gin-gonic/gin"

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

type AuthenticateUserResponse struct {
	User   UserResponse               `json:"user"`
	Tokens AuthenticateTokensResponse `json:"tokens"`
}

type PasetoValidateResponse struct {
	Valid  bool   `json:"valid"`
	UserID string `json:"user_id"`
	Email  string `json:"email"`
}

type APIKeyValidateResponse struct {
	Valid                     bool   `json:"valid"`
	IntegrationID             string `json:"integration_id"`
	IntegrationName           string `json:"integration_name"`
	IntegrationTargetChannels string `json:"integration_target_channels"`
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
