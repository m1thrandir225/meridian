package handlers

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/m1thrandir225/meridian/internal/integration/application/services"
	"github.com/m1thrandir225/meridian/internal/integration/domain"
	"github.com/m1thrandir225/meridian/pkg/cache"
)

type HTTPHandler struct {
	integrationService *services.IntegrationService
	cache              *cache.RedisCache
}

func NewHttpHandler(
	service *services.IntegrationService,
	cache *cache.RedisCache,
) *HTTPHandler {
	return &HTTPHandler{
		integrationService: service,
		cache:              cache,
	}
}

// POST /api/v1/integrations
func (h *HTTPHandler) handleRegisterIntegration(ctx *gin.Context) {
	creatorID := ctx.GetHeader("X-User-ID")
	if creatorID == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	var req RegisterIntegrationRequest
	if err := ctx.BindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cmd := domain.RegisterIntegrationCommand{
		ServiceName:    req.ServiceName,
		CreatorUserID:  creatorID,
		TargetChannels: req.TargetChannelIDs,
	}

	integration, token, err := h.integrationService.RegisterIntegration(ctx, cmd)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp := RegisterIntegrationResponse{
		ServiceName:    integration.ServiceName,
		TargetChannels: integration.TargetChannelIDsAsStringSlice(),
		Token:          token,
	}

	cacheKey := fmt.Sprintf("integration:%s", integration.ID)
	h.cache.Set(ctx, cacheKey, resp, 15*time.Minute)

	ctx.JSON(http.StatusOK, resp)
}

// DELETE /api/v1/integrations
func (h *HTTPHandler) handleRevokeIntegration(ctx *gin.Context) {
	requestorID := ctx.GetHeader("X-User-ID")
	if requestorID == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var req RevokeIntegrationRequest
	if err := ctx.BindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cmd := domain.RevokeTokenCommand{
		IntegrationID: req.IntegrationID,
		RequestorID:   requestorID,
	}

	err := h.integrationService.RevokeToken(ctx, cmd)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Integration not found"})
		} else if errors.Is(err, domain.ErrForbidden) {
			ctx.JSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
		} else {
			log.Printf("ERROR: Failed to revoke integration: %v", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		}
		return
	}

	cacheKey := fmt.Sprintf("integration:%s", req.IntegrationID)
	h.cache.Delete(ctx, cacheKey)

	ctx.Status(http.StatusNoContent)
}

func (h *HTTPHandler) handleGetMetrics(ctx *gin.Context) {
	metrics := h.cache.GetMetrics()
	ctx.JSON(http.StatusOK, gin.H{
		"hits":     metrics.GetHits(),
		"misses":   metrics.GetMisses(),
		"hit_rate": metrics.GetHitRate(),
	})
}

func (h *HTTPHandler) handleGetHealth(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"status":  "healthy",
		"service": "integration",
	})
}
