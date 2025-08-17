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

// NewHttpHandler creates a new HTTPHandler
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
		ctx.JSON(http.StatusUnauthorized, errorResponse(fmt.Errorf("missing creator ID")))
		return
	}
	var req RegisterIntegrationRequest
	if err := ctx.BindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	cmd := domain.RegisterIntegrationCommand{
		ServiceName:    req.ServiceName,
		CreatorUserID:  creatorID,
		TargetChannels: req.TargetChannelIDs,
	}

	integration, token, err := h.integrationService.RegisterIntegration(ctx, cmd)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	registerBotReq := &messagingpb.RegisterBotRequest{
		IntegrationId: integration.ID.String(),
		ChannelIds:    integration.TargetChannelIDsAsStringSlice(),
	}

	grpcResp, err := h.messageClient.RegisterBot(ctx, registerBotReq)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	if !grpcResp.Success {
		ctx.JSON(http.StatusInternalServerError, errorResponse(errors.New(grpcResp.Error)))
		return
	}

	// Invalidate channel caches for all affected channels
	for _, channelID := range req.TargetChannelIDs {
		channelCacheKey := fmt.Sprintf("channel:%s", channelID)
		h.cache.Delete(ctx, channelCacheKey)
		log.Printf("Invalidated cache for channel: %s", channelID)

		userChannelsCacheKey := fmt.Sprintf("user_channels:%s", creatorID)
		h.cache.Delete(ctx, userChannelsCacheKey)
		log.Printf("Invalidated user channels cache for user: %s", creatorID)
	}

	resp := ToIntegrationDTO(integration, token)

	cacheKey := fmt.Sprintf("integration:%s", integration.ID.String())
	h.cache.Set(ctx, cacheKey, resp, 15*time.Minute)

	ctx.JSON(http.StatusOK, resp)
}

// DELETE /api/v1/integrations
func (h *HTTPHandler) handleRevokeIntegration(ctx *gin.Context) {
	requestorID := ctx.GetHeader("X-User-ID")
	if requestorID == "" {
		ctx.JSON(http.StatusUnauthorized, errorResponse(fmt.Errorf("missing requestorID")))
		return
	}

	var req RevokeIntegrationRequest
	if err := ctx.BindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	cmd := domain.RevokeTokenCommand{
		IntegrationID: req.IntegrationID,
		RequestorID:   requestorID,
	}

	err := h.integrationService.RevokeToken(ctx, cmd)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			ctx.JSON(http.StatusNotFound, errorResponse(fmt.Errorf("integration not found")))
		} else if errors.Is(err, domain.ErrForbidden) {
			ctx.JSON(http.StatusForbidden, errorResponse(fmt.Errorf("forbidden")))
		} else {
			log.Printf("ERROR: Failed to revoke integration: %v", err)
			ctx.JSON(http.StatusInternalServerError, errorResponse(fmt.Errorf("internal server error")))
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
		ctx.JSON(http.StatusUnauthorized, errorResponse(fmt.Errorf("missing integration ID")))
		return
	}

	var req WebhookMessageRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var channelIDs []string
	if err := json.Unmarshal([]byte(targetChannels), &channelIDs); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("invalid target channels format")))
		return
	}

	targetChannelID := req.TargetChannelID
	if targetChannelID == "" {
		if len(channelIDs) == 0 {
			ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("no target channels available")))
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
		ctx.JSON(http.StatusForbidden, errorResponse(fmt.Errorf("channel not allowed for this integration")))
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
		ctx.JSON(http.StatusInternalServerError, errorResponse(fmt.Errorf("failed to send message")))
		return
	}

	if !resp.Success {
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("message delivery failed")))
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
		ctx.JSON(http.StatusUnauthorized, errorResponse(fmt.Errorf("missing integration ID")))
		return
	}

	var req CallbackMessageRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var channelIDs []string
	if err := json.Unmarshal([]byte(targetChannels), &channelIDs); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("invalid target channels format")))
		return
	}

	targetChannelID := req.TargetChannelID
	if targetChannelID == "" {
		if len(channelIDs) == 0 {
			ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("no target channels available")))
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
		ctx.JSON(http.StatusForbidden, errorResponse(fmt.Errorf("channel not allowed for this integration")))
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
		ctx.JSON(http.StatusInternalServerError, errorResponse(fmt.Errorf("failed to send message")))
		return
	}

	if !resp.Success {
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("message delivery failed")))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success":    true,
		"message_id": resp.Responses[0].MessageId,
		"channel_id": targetChannelID,
	})
}
