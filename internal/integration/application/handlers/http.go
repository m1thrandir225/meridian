package handlers

import (
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/m1thrandir225/meridian/internal/integration/application/services"
	"github.com/m1thrandir225/meridian/internal/integration/domain"
)

type HTTPHandler struct {
	integrationService *services.IntegrationService
}

func NewHttpHandler(service *services.IntegrationService) *HTTPHandler {
	return &HTTPHandler{
		integrationService: service,
	}
}

type RegisterIntegrationRequest struct {
	ServiceName      string   `json:"service_name" binding:"required"`
	TargetChannelIDs []string `json:"target_channel_ids" binding:"required"`
}

type RegisterIntegrationResponse struct {
	ServiceName    string   `json:"service_name"`
	TargetChannels []string `json:"target_channels"`
	Token          string   `json:"token"`
}

type RevokeIntegrationRequest struct {
	IntegrationID string `json:"integration_id" binding:"required"`
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
	ctx.Status(http.StatusNoContent)
}
