package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/m1thrandir225/meridian/internal/integration/application/services"
	"github.com/m1thrandir225/meridian/internal/integration/domain"
	messagingpb "github.com/m1thrandir225/meridian/internal/messaging/infrastructure/api"
	"github.com/m1thrandir225/meridian/pkg/cache"
)

type HTTPHandler struct {
	integrationService *services.IntegrationService
	messageClient      *services.MessagingClient
	cache              *cache.RedisCache
}

func NewHttpHandler(
	service *services.IntegrationService,
	cache *cache.RedisCache,
	messageClient *services.MessagingClient,
) *HTTPHandler {
	return &HTTPHandler{
		integrationService: service,
		cache:              cache,
		messageClient:      messageClient,
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

func (h *HTTPHandler) handleWebhookMessage(ctx *gin.Context) {
	integrationID := ctx.GetHeader("X-Integration-ID")
	targetChannels := ctx.GetHeader("X-Integration-Target-Channels")

	if integrationID == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Missing integration ID"})
		return
	}

	var req WebhookMessageRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var channelIDs []string
	if err := json.Unmarshal([]byte(targetChannels), &channelIDs); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid target channels format"})
		return
	}

	targetChannelID := req.TargetChannelID
	if targetChannelID == "" {
		if len(channelIDs) == 0 {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "No target channels available"})
			return
		}
	}
	targetChannelID = channelIDs[0]

	channelAllowed := false
	for _, allowedChannel := range channelIDs {
		if allowedChannel == targetChannelID {
			channelAllowed = true
			break
		}
	}

	if !channelAllowed {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "Channel not allowed for this integration"})
		return
	}

	messageReq := &messagingpb.SendMessageRequest{
		Content:          req.ContentText,
		SenderId:         integrationID,
		SenderType:       "integration",
		SenderName:       "Integration Bot",
		TargetChannelIds: []string{targetChannelID},
		MessageType:      "integration_webhook",
		Metadata:         req.Metadata,
	}

	resp, err := h.messageClient.SendMessage(ctx, messageReq)
	if err != nil {
		log.Printf("ERROR: Failed to send message via gRPC: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send message"})
		return
	}

	if !resp.Success {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Message delivery failed"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success":    true,
		"message_id": resp.Responses[0].MessageId,
		"channel_id": targetChannelID,
	})
}

func (h *HTTPHandler) handleCallbackMessage(ctx *gin.Context) {
	integrationID := ctx.GetHeader("X-Integration-ID")
	targetChannels := ctx.GetHeader("X-Integration-Target-Channels")

	if integrationID == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Missing integration ID"})
		return
	}

	var req CallbackMessageRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var channelIDs []string
	if err := json.Unmarshal([]byte(targetChannels), &channelIDs); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid target channels format"})
		return
	}

	targetChannelID := req.TargetChannelID
	if targetChannelID == "" {
		if len(channelIDs) == 0 {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "No target channels available"})
			return
		}
		targetChannelID = channelIDs[0]
	}

	channelAllowed := false
	for _, allowedChannel := range channelIDs {
		if allowedChannel == targetChannelID {
			channelAllowed = true
			break
		}
	}

	if !channelAllowed {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "Channel not allowed for this integration"})
		return
	}

	messageReq := &messagingpb.SendMessageRequest{
		Content:          req.ContentText,
		SenderId:         integrationID,
		SenderType:       "integration",
		SenderName:       "Integration Bot",
		TargetChannelIds: []string{targetChannelID},
		MessageType:      "integration_callback",
		Metadata:         req.Metadata,
	}

	resp, err := h.messageClient.SendMessage(ctx, messageReq)
	if err != nil {
		log.Printf("ERROR: Failed to send message via gRPC: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send message"})
		return
	}

	if !resp.Success {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Message delivery failed"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success":    true,
		"message_id": resp.Responses[0].MessageId,
		"channel_id": targetChannelID,
	})
}
