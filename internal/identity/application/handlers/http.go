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

	claims, claimsExist := auth.TokenClaimsFromContext(ctx.Request.Context())
	if !claimsExist {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: TokenClaims not found  in token."})
		return
	}

	//TODO: fetch user from repository and return details
	ctx.JSON(http.StatusOK, gin.H{
		"userId": userId,
		"email":  claims.Custom.Email})
}

// PUT /api/v1/update-profile

func (h *HTTPHandler) UpdateCurrentUser(ctx *gin.Context) {
	//TODO: implement
	ctx.JSON(http.StatusOK, gin.H{"message": "To be implemented"})
}

// DELETE /api/v1/me
func (h *HTTPHandler) DeleteCurrentUser(ctx *gin.Context) {
	//TODO: implement
	ctx.JSON(http.StatusOK, gin.H{"message": "To be implemented"})
}
