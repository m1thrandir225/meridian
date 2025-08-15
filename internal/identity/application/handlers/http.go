package handlers

import (
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/m1thrandir225/meridian/internal/identity/application/services"
	"github.com/m1thrandir225/meridian/internal/identity/domain"
	"github.com/m1thrandir225/meridian/pkg/auth"
)

type HTTPHandler struct {
	userService *services.IdentityService
}

func NewHTTPHandler(userService *services.IdentityService) *HTTPHandler {
	return &HTTPHandler{
		userService: userService,
	}
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

	response := UserResponse{
		ID:        user.ID.String(),
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

	accessToken, refreshToken, claims, err := h.userService.AuthenticateUser(ctx, cmd)
	if err != nil {
		if errors.Is(err, domain.ErrAuthFailed) {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		} else {
			log.Printf("ERROR authenticating user: %v", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Authentication process failed"})
		}
		return
	}

	getUserCMD := domain.GetUserCommand{
		UserID: claims.Custom.UserID,
	}

	user, err := h.userService.GetUser(ctx, getUserCMD)
	if err != nil {
		log.Printf("ERROR fetching user: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user, reason: " + err.Error()})
		return
	}

	resp := AuthenticateUserResponse{
		User: UserResponse{
			ID:        user.ID.String(),
			Username:  user.Username.String(),
			Email:     user.Email.String(),
			FirstName: user.FirstName,
			LastName:  user.LastName,
		},
		Tokens: AuthenticateTokensResponse{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
			TokenType:    "bearer",
			ExpiresIn:    int64(h.userService.AuthTokenValidity.Seconds()),
		},
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

	response := UserResponse{
		ID:        user.ID.String(),
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
func (h *HTTPHandler) handleUpdateCurrentUserRequest(ctx *gin.Context) {
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

	response := UserResponse{
		ID:        updatedUser.ID.String(),
		Username:  updatedUser.Username.String(),
		Email:     updatedUser.Email.String(),
		FirstName: updatedUser.FirstName,
		LastName:  updatedUser.LastName,
	}
	ctx.JSON(http.StatusOK, response)
}

// PATCH /api/v1/me/password
func (h *HTTPHandler) handleUpdateUserPasswordRequest(ctx *gin.Context) {
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
func (h *HTTPHandler) handleDeleteUserRequest(ctx *gin.Context) {
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

// POST /api/v1/me/refresh-token
func (h *HTTPHandler) handleRefreshTokenRequest(ctx *gin.Context) {
	var req RefreshTokenRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	cmd := domain.RefreshTokenCommand{
		Device:          ctx.Request.UserAgent(),
		IPAddress:       ctx.ClientIP(),
		RawRefreshToken: req.RefreshToken,
	}

	newAccessToken, newRefreshToken, err := h.userService.RefreshAuthentication(ctx, cmd)
	if err != nil {
		log.Printf("Refresh token failed: %v", err)
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid refresh token"})
		return
	}

	resp := AuthenticateTokensResponse{
		AccessToken:  newAccessToken,
		TokenType:    "bearer",
		ExpiresIn:    int64(h.userService.AuthTokenValidity.Seconds()),
		RefreshToken: newRefreshToken,
	}
	ctx.JSON(http.StatusOK, resp)
}
