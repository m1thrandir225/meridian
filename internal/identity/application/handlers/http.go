package handlers

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/m1thrandir225/meridian/internal/identity/application/services"
	"github.com/m1thrandir225/meridian/internal/identity/domain"
	"github.com/m1thrandir225/meridian/pkg/auth"
	"log"
	"net/http"
	"time"
)

type HTTPHandler struct {
	userService *services.IdentityService
}

func NewHTTPHandler(userService *services.IdentityService) *HTTPHandler {
	return &HTTPHandler{
		userService: userService,
	}
}

type RegisterUserRequest struct {
	Username  string `json:"username" binding:"required"`
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required"`
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
}

type RegisterUserResponse struct {
	UserID    string `json:"user_id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type AuthenticateUserRequest struct {
	Login    string `json:"login" binding:"required"` // Username or Email
	Password string `json:"password" binding:"required"`
}

type AuthenticateUserResponse struct {
	UserID         string    `json:"user_id"`
	Username       string    `json:"username"`
	AccessToken    string    `json:"access_token"`
	TokenType      string    `json:"token_type"`
	ExpirationDate time.Time `json:"expiration_date"`
	//RefreshToken string `json:"refresh_token,omitempty"`
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

type UpdateProfileResponse struct {
	UserID    string `json:"user_id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

// POST /api/v1/register
func (h *HTTPHandler) handleRegisterRequest(ctx *gin.Context) {
	var req RegisterUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	cmd := domain.RegisterUserCommand{
		Username:  req.Username,
		Email:     req.Email,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Password:  req.Password,
	}

	user, err := h.userService.RegisterUser(ctx, cmd)
	if err != nil {
		if errors.Is(err, domain.ErrUsernameTaken) || errors.Is(err, domain.ErrEmailTaken) {
			ctx.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		} else if errors.Is(err, domain.ErrPasswordPolicy) || errors.Is(err, domain.ErrUsernameInvalid) || errors.Is(err, domain.ErrEmailInvalid) {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		} else {
			log.Printf("ERROR registering user: %v", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register user"})
		}
		return
	}

	response := RegisterUserResponse{
		UserID:    user.ID.String(),
		Username:  user.Username.String(),
		Email:     user.Email.String(),
		FirstName: user.FirstName,
		LastName:  user.LastName,
	}
	ctx.JSON(http.StatusOK, response)
}

// POST /api/v1/login
func (h *HTTPHandler) handleLoginRequest(ctx *gin.Context) {
	var req AuthenticateUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	cmd := domain.AuthenticateUserCommand{
		LoginIdentifier: req.Login,
		Password:        req.Password,
	}

	tokenString, claims, err := h.userService.AuthenticateUser(ctx, cmd)
	if err != nil {
		if errors.Is(err, domain.ErrAuthFailed) {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		} else {
			log.Printf("ERROR authenticating user: %v", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Authentication process failed"})
		}
		return
	}

	resp := AuthenticateUserResponse{
		UserID:         claims.Custom.UserID,
		Username:       claims.Custom.Email,
		AccessToken:    tokenString,
		TokenType:      "bearer",
		ExpirationDate: claims.ExpirationDate,
	}

	ctx.JSON(http.StatusOK, resp)
}

// GET /api/v1/me
func (h *HTTPHandler) handleGetCurrentUser(ctx *gin.Context) {
	userId, exists := auth.UserIDFromContext(ctx.Request.Context())
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: User ID not found in token"})
		return
	}

	cmd := domain.GetUserCommand{
		UserID: userId,
	}

	user, err := h.userService.GetUser(ctx, cmd)
	if err != nil {
		log.Printf("ERROR fetching user: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user, reason: " + err.Error()})
		return
	}

	response := RegisterUserResponse{
		UserID:    user.ID.String(),
		Username:  user.Username.String(),
		Email:     user.Email.String(),
		FirstName: user.FirstName,
		LastName:  user.LastName,
	}

	ctx.JSON(http.StatusOK, response)
}

func (req *UpdateProfileRequest) HasAnyField() bool {
	return req.Username != nil || req.FirstName != nil || req.LastName != nil || req.Email != nil
}

// PUT /api/v1/update-profile
func (h *HTTPHandler) UpdateCurrentUser(ctx *gin.Context) {
	userId, exists := auth.UserIDFromContext(ctx.Request.Context())
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: User ID not found in token"})
		return
	}

	var req UpdateProfileRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if !req.HasAnyField() {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "At least one field must be provided to update profile"})
	}

	var cmd domain.UpdateUserProfileCommand
	cmd.UserID = userId
	if req.Username != nil {
		cmd.NewUsername = req.Username
	}
	if req.FirstName != nil {
		cmd.NewFirstName = req.FirstName
	}
	if req.LastName != nil {
		cmd.NewLastName = req.LastName
	}
	if req.Email != nil {
		cmd.NewEmail = req.Email
	}

	updatedUser, err := h.userService.UpdateUserProfile(ctx, cmd)
	if err != nil {
		log.Printf("ERROR updating user profile: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user profile, reason: " + err.Error()})
		return
	}

	response := UpdateProfileResponse{
		UserID:    updatedUser.ID.String(),
		Username:  updatedUser.Username.String(),
		Email:     updatedUser.Email.String(),
		FirstName: updatedUser.FirstName,
		LastName:  updatedUser.LastName,
	}
	ctx.JSON(http.StatusOK, response)
}

// PATCH /api/v1/me/password
func (h *HTTPHandler) UpdateUserPassword(ctx *gin.Context) {
	userId, exists := auth.UserIDFromContext(ctx.Request.Context())
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: User ID not found in token"})
		return
	}
	var req UpdatePasswordRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	cmd := domain.UpdateUserPasswordCommand{
		UserID:      userId,
		NewPassword: req.NewPassword,
	}

	err := h.userService.UpdateUserPassword(ctx, cmd)
	if err != nil {
		log.Printf("ERROR updating user password: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user password, reason: " + err.Error()})
		return
	}
	ctx.Status(http.StatusAccepted)
}

// DELETE /api/v1/me
func (h *HTTPHandler) DeleteCurrentUser(ctx *gin.Context) {
	userId, exists := auth.UserIDFromContext(ctx.Request.Context())
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: User ID not found in token"})
		return
	}

	cmd := domain.DeleteUserCommand{
		UserID: userId,
	}
	err := h.userService.DeleteUser(ctx, cmd)
	if err != nil {
		log.Printf("ERROR deleting user: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user, reason: " + err.Error()})
		return
	}
	ctx.Status(http.StatusAccepted)
}
