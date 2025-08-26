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

// GET /api/v1/integrations
func (h *HTTPHandler) handleListIntegrations(ctx *gin.Context) {
	requestorID := ctx.GetHeader("X-User-ID")
	if requestorID == "" {
		ctx.JSON(http.StatusUnauthorized, errorResponse(fmt.Errorf("missing requestor ID")))
		return
	}

	cmd := domain.ListIntegrationsCommand{
		CreatorUserID: requestorID,
	}

	integrations, err := h.integrationService.ListIntegrations(ctx, cmd)
	if err != nil {
		log.Printf("ERROR: Failed to list integrations: %v", err)
		ctx.JSON(http.StatusInternalServerError, errorResponse(fmt.Errorf("internal server error")))
		return
	}

	// Convert to DTOs
	var dtos []IntegrationDTO
	for _, integration := range integrations {
		dto := ToIntegrationDTO(integration, "")
		dtos = append(dtos, dto)
	}

	ctx.JSON(http.StatusOK, gin.H{
		"integrations": dtos,
	})
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
		RequestorId:   creatorID,
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

// POST /api/v1/integrations/upvoke
func (h *HTTPHandler) handleUpvokeIntegration(ctx *gin.Context) {
	requestorID := ctx.GetHeader("X-User-ID")
	if requestorID == "" {
		ctx.JSON(http.StatusUnauthorized, errorResponse(fmt.Errorf("missing requestorID")))
		return
	}

	var req UpvokeIntegrationRequest
	if err := ctx.BindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	cmd := domain.UpvokeIntegrationCommand{
		IntegrationID: req.IntegrationID,
		RequestorID:   requestorID,
	}

	integration, newToken, err := h.integrationService.UpvokeIntegration(ctx, cmd)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	integrationDTO := ToIntegrationDTO(integration, newToken)

	cacheKey := fmt.Sprintf("integration:%s", req.IntegrationID)
	grpcCacheKey := fmt.Sprintf("grpc_integration:%s", req.IntegrationID)
	h.cache.Delete(ctx, cacheKey)
	h.cache.Delete(ctx, grpcCacheKey)

	ctx.JSON(http.StatusOK, integrationDTO)
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
	grpcCacheKey := fmt.Sprintf("grpc_integration:%s", req.IntegrationID)
	h.cache.Delete(ctx, cacheKey)
	h.cache.Delete(ctx, grpcCacheKey)

	ctx.Status(http.StatusNoContent)
}

// PUT /api/v1/integrations/:id
func (h *HTTPHandler) handleUpdateIntegration(ctx *gin.Context) {
	requestorID := ctx.GetHeader("X-User-ID")
	if requestorID == "" {
		ctx.JSON(http.StatusUnauthorized, errorResponse(fmt.Errorf("missing requestor ID")))
		return
	}
	var uri IntegrationURI
	if err := ctx.ShouldBindUri(&uri); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var req UpdateIntegrationRequest
	if err := ctx.BindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	getCMD := domain.GetIntegrationCommand{
		IntegrationID: uri.IntegrationID,
	}

	current, err := h.integrationService.GetIntegration(ctx, getCMD)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(fmt.Errorf("failed to load integration")))
		return
	}
	oldSet := map[string]struct{}{}
	for _, id := range current.TargetChannelIDsAsStringSlice() {
		oldSet[id] = struct{}{}
	}
	newSet := map[string]struct{}{}
	for _, id := range req.TargetChannelIDs {
		newSet[id] = struct{}{}
	}
	var toAdd, toRemove []string
	for id := range newSet {
		if _, ok := oldSet[id]; !ok {
			toAdd = append(toAdd, id)
		}
	}
	for id := range oldSet {
		if _, ok := newSet[id]; !ok {
			toRemove = append(toRemove, id)
		}
	}

	// 1) Apply membership changes first
	if len(toAdd) > 0 {
		if _, err := h.messageClient.RegisterBot(ctx, &messagingpb.RegisterBotRequest{
			IntegrationId: current.ID.String(),
			ChannelIds:    toAdd,
			RequestorId:   requestorID,
		}); err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(fmt.Errorf("failed adding bot to channels: %w", err)))
			return
		}
	}
	if len(toRemove) > 0 {
		if _, err := h.messageClient.RemoveBot(ctx, &messagingpb.RemoveBotRequest{
			IntegrationId: current.ID.String(),
			ChannelIds:    toRemove,
			RequestorId:   requestorID,
		}); err != nil {
			// rollback adds if any
			if len(toAdd) > 0 {
				_, _ = h.messageClient.RemoveBot(ctx, &messagingpb.RemoveBotRequest{
					IntegrationId: current.ID.String(),
					ChannelIds:    toAdd,
					RequestorId:   requestorID,
				})
			}
			ctx.JSON(http.StatusInternalServerError, errorResponse(fmt.Errorf("failed removing bot from channels: %w", err)))
			return
		}
	}

	// 2) Persist new integration config
	integration, err := h.integrationService.UpdateIntegration(ctx, domain.UpdateIntegrationCommand{
		IntegrationID:    uri.IntegrationID,
		RequestorID:      requestorID,
		TargetChannelIDs: req.TargetChannelIDs,
	})
	if err != nil {
		// best-effort rollback messaging to previous state
		if len(toAdd) > 0 {
			_, _ = h.messageClient.RemoveBot(ctx, &messagingpb.RemoveBotRequest{
				IntegrationId: current.ID.String(),
				ChannelIds:    toAdd,
				RequestorId:   requestorID,
			})
		}
		if len(toRemove) > 0 {
			_, _ = h.messageClient.RegisterBot(ctx, &messagingpb.RegisterBotRequest{
				IntegrationId: current.ID.String(),
				ChannelIds:    toRemove,
				RequestorId:   requestorID,
			})
		}

		if errors.Is(err, domain.ErrNotFound) {
			ctx.JSON(http.StatusNotFound, errorResponse(fmt.Errorf("integration not found")))
		} else if errors.Is(err, domain.ErrForbidden) {
			ctx.JSON(http.StatusForbidden, errorResponse(fmt.Errorf("forbidden")))
		} else if errors.Is(err, domain.ErrIntegrationRevoked) {
			ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("cannot update revoked integration")))
		} else {
			ctx.JSON(http.StatusInternalServerError, errorResponse(fmt.Errorf("internal server error")))
		}
		return
	}

	// Invalidate caches
	for _, id := range append(append([]string{}, toAdd...), toRemove...) {
		h.cache.Delete(ctx, fmt.Sprintf("channel:%s", id))
		h.cache.Delete(ctx, fmt.Sprintf("user_channels:%s", requestorID))
	}

	ctx.JSON(http.StatusOK, ToIntegrationDTO(integration, ""))
}

func (h *HTTPHandler) handleDeleteIntegration(ctx *gin.Context) {
	requestorID := ctx.GetHeader("X-User-ID")
	if requestorID == "" {
		ctx.JSON(http.StatusUnauthorized, errorResponse(fmt.Errorf("missing requestor ID")))
		return
	}
	var uri IntegrationURI
	if err := ctx.ShouldBindUri(&uri); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	getCMD := domain.GetIntegrationCommand{
		IntegrationID: uri.IntegrationID,
	}

	integration, err := h.integrationService.GetIntegration(ctx, getCMD)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			ctx.JSON(http.StatusNotFound, errorResponse(fmt.Errorf("integration not found")))
			return
		} else if errors.Is(err, domain.ErrForbidden) {
			ctx.JSON(http.StatusForbidden, errorResponse(fmt.Errorf("forbidden")))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(fmt.Errorf("failed to get integration: %w", err)))
		return
	}

	channels := integration.TargetChannelIDsAsStringSlice()
	if len(channels) > 0 {
		_, err := h.messageClient.RemoveBot(ctx, &messagingpb.RemoveBotRequest{
			IntegrationId: integration.ID.String(),
			ChannelIds:    channels,
			RequestorId:   requestorID,
		})
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(fmt.Errorf("failed to remove bot from channels: %w", err)))
			return
		}
	}

	deleteCMD := domain.DeleteIntegrationCommand{
		IntegrationID: uri.IntegrationID,
		RequestorID:   requestorID,
	}

	err = h.integrationService.DeleteIntegration(ctx, deleteCMD)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			ctx.JSON(http.StatusNotFound, errorResponse(fmt.Errorf("integration not found")))
		} else if errors.Is(err, domain.ErrForbidden) {
			ctx.JSON(http.StatusForbidden, errorResponse(fmt.Errorf("forbidden")))
		} else {
			ctx.JSON(http.StatusInternalServerError, errorResponse(fmt.Errorf("internal server error")))
		}
		return
	}
	cacheKey := fmt.Sprintf("integration:%s", uri.IntegrationID)
	grpcCacheKey := fmt.Sprintf("grpc_integration:%s", uri.IntegrationID)
	h.cache.Delete(ctx, grpcCacheKey)
	h.cache.Delete(ctx, cacheKey)

	for _, channelID := range channels {
		h.cache.Delete(ctx, fmt.Sprintf("channel:%s", channelID))
	}
	h.cache.Delete(ctx, fmt.Sprintf("user_channels:%s", requestorID))

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

	if len(channelIDs) == 0 {
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("no target channels available")))
		return
	}

	messageReq := &messagingpb.SendMessageRequest{
		Content:          req.ContentText,
		SenderId:         integrationID,
		SenderType:       "integration",
		SenderName:       "Integration Bot",
		TargetChannelIds: channelIDs,
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
		"success":     true,
		"message_id":  resp.Responses[0].MessageId,
		"channel_ids": channelIDs,
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
