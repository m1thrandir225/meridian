package handlers

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/m1thrandir225/meridian/internal/identity/application/services"
	"github.com/m1thrandir225/meridian/internal/identity/domain"
	"github.com/m1thrandir225/meridian/pkg/auth"
	"github.com/m1thrandir225/meridian/pkg/cache"
)

type HTTPHandler struct {
	userService *services.IdentityService
	cache       *cache.RedisCache
}

func NewHTTPHandler(userService *services.IdentityService, cache *cache.RedisCache) *HTTPHandler {
	return &HTTPHandler{
		userService: userService,
		cache:       cache,
	}
}

// POST /api/v1/register
func (h *HTTPHandler) handleRegisterRequest(ctx *gin.Context) {
	var req RegisterUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
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
			ctx.JSON(http.StatusConflict, errorResponse(err))
		} else if errors.Is(err, domain.ErrPasswordPolicy) || errors.Is(err, domain.ErrUsernameInvalid) || errors.Is(err, domain.ErrEmailInvalid) {
			ctx.JSON(http.StatusBadRequest, errorResponse(err))
		} else {
			log.Printf("ERROR registering user: %v", err)
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
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
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
	}

	cmd := domain.AuthenticateUserCommand{
		LoginIdentifier: req.Login,
		Password:        req.Password,
	}

	accessToken, refreshToken, claims, err := h.userService.AuthenticateUser(ctx, cmd)
	if err != nil {
		if errors.Is(err, domain.ErrAuthFailed) {
			ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		} else {
			log.Printf("ERROR authenticating user: %v", err)
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		}
		return
	}

	getUserCMD := domain.GetUserCommand{
		UserID: claims.Custom.UserID,
	}

	user, err := h.userService.GetUser(ctx, getUserCMD)
	if err != nil {
		log.Printf("ERROR fetching user: %v", err)
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
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
		ctx.JSON(http.StatusUnauthorized, errorResponse(ErrUnauthorized))
		return
	}

	cacheKey := fmt.Sprintf("user_profile:%s", userId)
	var cachedUser UserResponse
	if hit, _ := h.cache.GetWithMetrics(ctx.Request.Context(), cacheKey, &cachedUser); hit {
		ctx.JSON(http.StatusOK, cachedUser)
		return
	}

	cmd := domain.GetUserCommand{
		UserID: userId,
	}

	user, err := h.userService.GetUser(ctx, cmd)
	if err != nil {
		log.Printf("ERROR fetching user: %v", err)
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	response := UserResponse{
		ID:        user.ID.String(),
		Username:  user.Username.String(),
		Email:     user.Email.String(),
		FirstName: user.FirstName,
		LastName:  user.LastName,
	}

	h.cache.Set(ctx.Request.Context(), cacheKey, response, 15*time.Minute)

	ctx.JSON(http.StatusOK, response)
}

func (req *UpdateProfileRequest) HasAnyField() bool {
	return req.Username != nil || req.FirstName != nil || req.LastName != nil || req.Email != nil
}

// PUT /api/v1/update-profile
func (h *HTTPHandler) handleUpdateCurrentUserRequest(ctx *gin.Context) {
	userId, exists := auth.UserIDFromContext(ctx.Request.Context())
	if !exists {
		ctx.JSON(http.StatusUnauthorized, errorResponse(ErrUnauthorized))
		return
	}

	var req UpdateProfileRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	if !req.HasAnyField() {
		ctx.JSON(http.StatusBadRequest, errorResponse(errors.New("at least one field must be provided to update profile")))
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
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	cacheKey := fmt.Sprintf("user_profile:%s", userId)
	h.cache.Delete(ctx.Request.Context(), cacheKey)
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
		ctx.JSON(http.StatusUnauthorized, errorResponse(ErrUnauthorized))
		return
	}
	var req UpdatePasswordRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	cmd := domain.UpdateUserPasswordCommand{
		UserID:      userId,
		NewPassword: req.NewPassword,
	}

	err := h.userService.UpdateUserPassword(ctx, cmd)
	if err != nil {
		log.Printf("ERROR updating user password: %v", err)
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.Status(http.StatusAccepted)
}

// DELETE /api/v1/me
func (h *HTTPHandler) handleDeleteUserRequest(ctx *gin.Context) {
	userId, exists := auth.UserIDFromContext(ctx.Request.Context())
	if !exists {
		ctx.JSON(http.StatusUnauthorized, errorResponse(ErrUnauthorized))
		return
	}

	cmd := domain.DeleteUserCommand{
		UserID: userId,
	}
	err := h.userService.DeleteUser(ctx, cmd)
	if err != nil {
		log.Printf("ERROR deleting user: %v", err)
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.Status(http.StatusAccepted)
}

// POST /api/v1/me/refresh-token
func (h *HTTPHandler) handleRefreshTokenRequest(ctx *gin.Context) {
	var req RefreshTokenRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
	}

	cmd := domain.RefreshTokenCommand{
		Device:          ctx.Request.UserAgent(),
		IPAddress:       ctx.ClientIP(),
		RawRefreshToken: req.RefreshToken,
	}

	newAccessToken, newRefreshToken, err := h.userService.RefreshAuthentication(ctx, cmd)
	if err != nil {
		log.Printf("Refresh token failed: %v", err)
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
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

// GET /api/v1/me/metrics
func (h *HTTPHandler) handleGetMetrics(ctx *gin.Context) {
	metrics := h.cache.GetMetrics()

	ctx.JSON(http.StatusOK, gin.H{
		"hits":     metrics.GetHits(),
		"misses":   metrics.GetMisses(),
		"hit_rate": metrics.GetHitRate(),
	})
}

// GET /api/v1/health
func (h *HTTPHandler) handleGetHealth(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"status":  "healthy",
		"service": "identity",
	})
}
