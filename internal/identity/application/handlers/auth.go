package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/m1thrandir225/meridian/internal/identity/application/services"
	"github.com/m1thrandir225/meridian/pkg/auth"
	"github.com/m1thrandir225/meridian/pkg/cache"
)

var (
	ErrInvalidToken                    = errors.New("invalid token")
	ErrInvalidAPIKey                   = errors.New("invalid api key")
	ErrUnauthorized                    = errors.New("unauthorized")
	ErrInvalidAuthorizationFormat      = errors.New("invalid authorization format. Use 'Bearer <token>' or 'ApiKey <key>'")
	ErrIntegrationEndpoint             = errors.New("unauthorized: Integration endpoint")
	ErrFailedToCreateIntegrationClient = errors.New("failed to create integration client")
)

type AuthHandler struct {
	userService        *services.IdentityService
	cache              *cache.RedisCache
	tokenVerifier      auth.TokenVerifier
	integrationGrpcURL string
}

func NewAuthHandler(
	userService *services.IdentityService,
	cache *cache.RedisCache,
	tokenVerifier auth.TokenVerifier,
	integrationGRPCURL string,
) *AuthHandler {
	return &AuthHandler{
		userService:        userService,
		cache:              cache,
		tokenVerifier:      tokenVerifier,
		integrationGrpcURL: integrationGRPCURL,
	}
}

func (h *AuthHandler) ValidateToken(ctx *gin.Context) {
	authHeader := ctx.GetHeader("Authorization")
	if authHeader == "" {
		ctx.JSON(http.StatusUnauthorized, errorResponse(ErrUnauthorized))
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
		ctx.JSON(http.StatusUnauthorized, errorResponse(ErrInvalidAuthorizationFormat))
	}

}

func (h *AuthHandler) handlePasetoAuth(ctx *gin.Context, token string, isIntegrationEndpoint bool, path string) {
	cacheKey := fmt.Sprintf("token_validation:%s", token)
	var cachedResponse PasetoValidateResponse
	if hit, _ := h.cache.GetWithMetrics(ctx.Request.Context(), cacheKey, &cachedResponse); hit {
		ctx.Header("X-User-ID", cachedResponse.UserID)
		ctx.Header("X-User-Email", cachedResponse.Email)
		ctx.Header("X-User-Type", "user")
		ctx.Header("X-Auth-Method", "paseto")
		ctx.JSON(http.StatusOK, cachedResponse)
		return
	}

	claims, err := h.tokenVerifier.Verify(token)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(ErrInvalidToken))
		return
	}

	if isIntegrationEndpoint {
		ctx.JSON(http.StatusUnauthorized, errorResponse(ErrIntegrationEndpoint))
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

	h.cache.Set(ctx.Request.Context(), cacheKey, response, 15*time.Minute)

	ctx.JSON(http.StatusOK, response)
}

func (h *AuthHandler) handleAPITokenAuth(ctx *gin.Context, apiKey string, isIntegrationEndpoint bool, path string) {
	if !isIntegrationEndpoint {
		ctx.JSON(http.StatusUnauthorized, errorResponse(ErrIntegrationEndpoint))
		return
	}

	integrationClient, err := services.NewIntegrationClient(h.integrationGrpcURL)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(ErrFailedToCreateIntegrationClient))
		return
	}

	defer integrationClient.Close()

	resp, err := integrationClient.ValidateAPIToken(apiKey)
	if err != nil || !resp.Valid {
		ctx.JSON(http.StatusUnauthorized, errorResponse(ErrInvalidAPIKey))
		return
	}

	targetChannelsJSON := "[]"
	if len(resp.TargetChannelIds) > 0 {
		targetChannelsJSON = fmt.Sprintf(`["%s"]`, strings.Join(resp.TargetChannelIds, `","`))
	}

	ctx.Header("X-Integration-ID", resp.IntegrationId)
	ctx.Header("X-Integration-Name", resp.IntegrationName)
	ctx.Header("X-User-Type", "integration")
	ctx.Header("X-Auth-Method", "api-token")
	ctx.Header("X-Integration-Target-Channels", targetChannelsJSON)

	response := APIKeyValidateResponse{
		Valid:                     true,
		IntegrationID:             resp.IntegrationId,
		IntegrationName:           resp.IntegrationName,
		IntegrationTargetChannels: targetChannelsJSON,
	}

	ctx.JSON(http.StatusOK, response)
}
