package handlers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/m1thrandir225/meridian/internal/identity/application/services"
	"github.com/m1thrandir225/meridian/pkg/auth"
)

type AuthHandler struct {
	userService        *services.IdentityService
	tokenVerifier      auth.TokenVerifier
	integrationGrpcURL string
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

func NewAuthHandler(userService *services.IdentityService, tokenVerifier auth.TokenVerifier, integrationGRPCURL string) *AuthHandler {
	return &AuthHandler{
		userService:        userService,
		tokenVerifier:      tokenVerifier,
		integrationGrpcURL: integrationGRPCURL,
	}
}

func (h *AuthHandler) ValidateToken(ctx *gin.Context) {
	authHeader := ctx.GetHeader("Authorization")
	if authHeader == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "No authorization header"})
		return
	}
	originalPath := ctx.GetHeader("X-Forwarded-Uri")
	if originalPath == "" {
		originalPath = ctx.Request.URL.Path
	}

	originalMethod := ctx.GetHeader("X-Forwarded-Method")
	if originalMethod == "" {
		originalMethod = ctx.Request.Method
	}

	isIntegrationEndpoint := strings.Contains(originalPath, "/integrations/webhook") ||
		strings.Contains(originalPath, "/integrations/callback")

	if strings.HasPrefix(authHeader, "Bearer ") {
		token := strings.TrimPrefix(authHeader, "Bearer ")
		h.handlePasetoAuth(ctx, token, isIntegrationEndpoint, originalPath)
	} else if strings.HasPrefix(authHeader, "ApiKey ") {
		apiKey := strings.TrimPrefix(authHeader, "ApiKey ")
		h.handleAPITokenAuth(ctx, apiKey, isIntegrationEndpoint, originalPath)
	} else {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization format. Use 'Bearer <token>' or 'ApiKey <key>'"})
	}

}

func (h *AuthHandler) handlePasetoAuth(ctx *gin.Context, token string, isIntegrationEndpoint bool, path string) {
	claims, err := h.tokenVerifier.Verify(token)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	if isIntegrationEndpoint {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: Integration endpoint"})
		return
	}

	ctx.Header("X-User-ID", claims.Custom.UserID)
	ctx.Header("X-User-Email", claims.Custom.Email)
	ctx.Header("X-User-Type", "user")
	ctx.Header("X-Auth-Method", "paseto")

	response := PasetoValidateResponse{
		Valid:  true,
		UserID: claims.Custom.UserID,
		Email:  claims.Custom.Email,
	}

	ctx.JSON(http.StatusOK, response)
}

func (h *AuthHandler) handleAPITokenAuth(ctx *gin.Context, apiKey string, isIntegrationEndpoint bool, path string) {
	if !isIntegrationEndpoint {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: Integration endpoint"})
		return
	}

	integrationClient, err := services.NewIntegrationClient(h.integrationGrpcURL)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create integration client"})
		return
	}

	defer integrationClient.Close()

	resp, err := integrationClient.ValidateAPIToken(apiKey)
	if err != nil || !resp.Valid {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid API token"})
		return
	}

	ctx.Header("X-Integration-ID", resp.IntegrationId)
	ctx.Header("X-Integration-Name", resp.IntegrationName)
	ctx.Header("X-User-Type", "integration")
	ctx.Header("X-Auth-Method", "api-token")
	ctx.Header("X-Integration-Target-Channels", "") //TODO: fix this

	response := APIKeyValidateResponse{
		Valid:                     true,
		IntegrationID:             resp.IntegrationId,
		IntegrationName:           resp.IntegrationId,
		IntegrationTargetChannels: "",
	}
	ctx.JSON(http.StatusOK, response)

}
