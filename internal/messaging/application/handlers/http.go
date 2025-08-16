package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/m1thrandir225/meridian/internal/messaging/application/services"
	"github.com/m1thrandir225/meridian/internal/messaging/domain"
	"github.com/m1thrandir225/meridian/pkg/cache"
)

type HTTPHandler struct {
	channelService *services.ChannelService
	messageService *services.MessageService
	cache          *cache.RedisCache
}

func NewHttpHandler(
	channelService *services.ChannelService,
	messageService *services.MessageService,
	cache *cache.RedisCache,
) *HTTPHandler {
	return &HTTPHandler{
		channelService: channelService,
		messageService: messageService,
		cache:          cache,
	}
}

func (h *HTTPHandler) GetUserChannels(ctx *gin.Context) {
	userIDStr := ctx.GetHeader("X-User-ID")
	if userIDStr == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "user id not valid"})
		return
	}

	cacheKey := fmt.Sprintf("user_channels:%s", userID)
	var cachedChannels []domain.Channel
	if hit, _ := h.cache.GetWithMetrics(ctx, cacheKey, &cachedChannels); hit {
		ctx.JSON(http.StatusOK, cachedChannels)
		return
	}

	cmd := domain.GetUserChannelsCommand{
		UserID: userID,
	}

	channels, err := h.channelService.HandleGetUserChannels(ctx, cmd)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	channelsDTO, err := h.channelService.ReturnChannelDTOs(ctx, channels)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	h.cache.Set(ctx, cacheKey, channelsDTO, 15*time.Minute)

	ctx.JSON(http.StatusOK, channelsDTO)
}

// POST /api/v1/channels/
func (h *HTTPHandler) CreateChannel(ctx *gin.Context) {
	creatorID := ctx.GetHeader("X-User-ID")
	if creatorID == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	var req CreateChannelRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	creatorUserID, err := uuid.Parse(creatorID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	channel, err := h.channelService.HandleCreateChannel(ctx, domain.CreateChannelCommand{
		CreatorUserID: creatorUserID,
		Name:          req.Name,
		Topic:         req.Topic,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	cacheKey := fmt.Sprintf("user_channels:%s", creatorUserID.String())
	h.cache.Delete(ctx.Request.Context(), cacheKey)

	channelDTO, err := h.channelService.ReturnChannelDTO(ctx, channel)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusCreated, channelDTO)
}

// GET /api/v1/channels/:channelId
func (h *HTTPHandler) GetChannel(ctx *gin.Context) {
	var req ChannelIDUri
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	channelId, err := uuid.Parse(req.ChannelID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	cacheKey := fmt.Sprintf("channel:%s", channelId.String())
	var cachedChannel interface{}
	if hit, _ := h.cache.GetWithMetrics(ctx.Request.Context(), cacheKey, &cachedChannel); hit {
		ctx.JSON(http.StatusOK, cachedChannel)
		return
	}

	channel, err := h.channelService.HandleGetChannel(ctx, domain.GetChannelCommand{
		ChannelID: channelId,
	})
	if err != nil {
		ctx.JSON(http.StatusNotFound, errorResponse(err))
		return
	}

	channelDTO, err := h.channelService.ReturnChannelDTO(ctx, channel)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	h.cache.Set(ctx.Request.Context(), cacheKey, channelDTO, 15*time.Minute)

	ctx.JSON(http.StatusOK, channelDTO)
}

// POST /api/v1/channels/:channelId/bots
func (h *HTTPHandler) AddBotToChannel(ctx *gin.Context) {
	var channelIdUri ChannelIDUri

	if err := ctx.ShouldBindUri(&channelIdUri); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var req AddBotToChannelRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	channelId, err := uuid.Parse(channelIdUri.ChannelID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	integrationId, err := uuid.Parse(req.IntegrationID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	channel, err := h.channelService.HandleAddBotToChannel(ctx, domain.AddBotToChannelCommand{
		ChannelID:     channelId,
		IntegrationID: integrationId,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	channelDTO, err := h.channelService.ReturnChannelDTO(ctx, channel)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, channelDTO)
}

// PUT /api/v1/channels/:channelId/archive
func (h *HTTPHandler) ArchiveChannel(ctx *gin.Context) {
	var req ChannelIDUri
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	// TODO: implement authentication
	ctx.Status(http.StatusOK)
}

// PUT /api/v1/channels/:channelId/unarchive
func (h *HTTPHandler) UnarchiveChannel(ctx *gin.Context) {
	var req ChannelIDUri
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// TODO: implement authentication

	ctx.Status(http.StatusOK)
}

// POST /api/v1/channels/:channelId/join
func (h *HTTPHandler) JoinChannel(ctx *gin.Context) {
	var req JoinChannelRequest
	var uriReq ChannelIDUri

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	if err := ctx.ShouldBindUri(&uriReq); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	userId, err := uuid.Parse(req.UserID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	channelId, err := uuid.Parse(uriReq.ChannelID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	channel, err := h.channelService.HandleJoinChannel(ctx, domain.JoinChannelCommand{
		ChannelID: channelId,
		UserID:    userId,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	channelDTO, err := h.channelService.ReturnChannelDTO(ctx, channel)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	cacheKey := fmt.Sprintf("user_channels:%s", userId.String())
	h.cache.Delete(ctx.Request.Context(), cacheKey)

	ctx.JSON(http.StatusOK, channelDTO)
}

// POST /api/v1/channels/:channelId/messages
// FIXME: redundant, should be removed
func (h *HTTPHandler) SendMessage(ctx *gin.Context) {
	userID := ctx.GetHeader("X-User-ID")
	if userID == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	var req SendMessageRequest
	var uriReq ChannelIDUri

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	if err := ctx.ShouldBindUri(&uriReq); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	channelId, err := uuid.Parse(uriReq.ChannelID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	senderID, err := uuid.Parse(userID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var parentMessageID *uuid.UUID
	if req.ParentMessageID != nil {
		parsed, err := uuid.Parse(*req.ParentMessageID)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, errorResponse(err))
			return
		}
		parentMessageID = &parsed
	}

	content := domain.NewMessageContent(req.ContentText) //TODO: add content type

	message, err := h.messageService.HandleMessageSent(ctx, domain.SendMessageCommand{
		ChannelID:       channelId,
		SenderUserID:    senderID,
		ParentMessageID: parentMessageID,
		Content:         content,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	messageDTO, err := h.messageService.ToMessageDTO(ctx, message)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusCreated, messageDTO)
}

// GET /api/v1/channels/:channelId/messages
func (h *HTTPHandler) GetMessages(ctx *gin.Context) {
	var uriReq ChannelIDUri

	if err := ctx.ShouldBindUri(&uriReq); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	channelId, err := uuid.Parse(uriReq.ChannelID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	cmd := domain.ListMessagesForChannelCommand{
		ChannelID: channelId,
		Limit:     50,
		Offset:    0,
	}

	messages, err := h.messageService.HandleListMessages(ctx, cmd)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	messagesDTO, err := h.messageService.ToMessageDTOs(ctx, messages)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, messagesDTO)
}

// POST /api/v1/channels/:channelId/messages/:messageId/reactions
func (h *HTTPHandler) AddReaction(ctx *gin.Context) {
	var req AddReactionRequest
	var channelIdUri ChannelIDUri
	var messageIdUri MessageIDUri

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	if err := ctx.ShouldBindUri(&channelIdUri); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	if err := ctx.ShouldBindUri(&messageIdUri); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	userId, err := uuid.Parse(req.UserID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	channelId, err := uuid.Parse(channelIdUri.ChannelID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	messageId, err := uuid.Parse(messageIdUri.MessageID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	cmd := domain.AddReactionCommand{
		ChannelID:    channelId,
		MessageID:    messageId,
		UserID:       userId,
		ReactionType: req.ReactionType,
	}

	reaction, err := h.messageService.HandleAddReaction(ctx, cmd)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	//TODO: return the reaction as a DTO? or return the message with the reaction?
	ctx.JSON(http.StatusCreated, reaction)
}

// DELETE /api/v1/channels/:channelId/messages/:messageId/reactions
func (h *HTTPHandler) RemoveReaction(ctx *gin.Context) {
	var req RemoveReactionRequest
	var channelIdUri ChannelIDUri
	var messageIdUri MessageIDUri

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	if err := ctx.ShouldBindUri(&channelIdUri); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	if err := ctx.ShouldBindUri(&messageIdUri); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	userId, err := uuid.Parse(req.UserID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	channelId, err := uuid.Parse(channelIdUri.ChannelID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	messageId, err := uuid.Parse(messageIdUri.MessageID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	cmd := domain.RemoveReactionCommand{
		ChannelID:    channelId,
		MessageID:    messageId,
		UserID:       userId,
		ReactionType: req.ReactionType,
	}

	_, err = h.messageService.HandleRemoveReaction(ctx, cmd)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.Status(http.StatusOK)
}

// GET /api/v1/metrics
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
		"service": "messaging",
	})
}
