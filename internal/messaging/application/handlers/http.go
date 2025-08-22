package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/m1thrandir225/meridian/internal/messaging/application/services"
	"github.com/m1thrandir225/meridian/internal/messaging/domain"
	"github.com/m1thrandir225/meridian/pkg/cache"
	"github.com/m1thrandir225/meridian/pkg/logging"
	"go.uber.org/zap"
)

var (
	ErrUnauthorized = errors.New("unauthorized")
)

type HTTPHandler struct {
	channelService *services.ChannelService
	messageService *services.MessageService
	cache          *cache.RedisCache
	logger         *logging.Logger
}

func NewHttpHandler(
	channelService *services.ChannelService,
	messageService *services.MessageService,
	cache *cache.RedisCache,
	logger *logging.Logger,
) *HTTPHandler {
	return &HTTPHandler{
		channelService: channelService,
		messageService: messageService,
		cache:          cache,
		logger:         logger,
	}
}

func (h *HTTPHandler) handleGetUserChannels(ctx *gin.Context) {
	logger := h.logger.WithMethod("handleGetUserChannels")
	logger.Info("Getting user channels")

	userIDStr := ctx.GetHeader("X-User-ID")
	if userIDStr == "" {
		ctx.JSON(http.StatusUnauthorized, errorResponse(ErrUnauthorized))
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		logger.Error("Failed to parse user ID", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// cacheKey := fmt.Sprintf("user_channels:%s", userID.String())
	// var cachedChannels []domain.ChannelDTO
	// if hit, _ := h.cache.GetWithMetrics(ctx.Request.Context(), cacheKey, &cachedChannels); hit {
	// 	ctx.JSON(http.StatusOK, cachedChannels)
	// 	return
	// }

	cmd := domain.GetUserChannelsCommand{
		UserID: userID,
	}

	channels, err := h.channelService.HandleGetUserChannels(ctx, cmd)
	if err != nil {
		logger.Error("Failed to get user channels", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	channelsDTO, err := h.channelService.ReturnChannelDTOs(ctx, channels)
	if err != nil {
		logger.Error("Failed to return channel DTOs", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	//h.cache.Set(ctx.Request.Context(), cacheKey, channelsDTO, 5*time.Minute)

	ctx.JSON(http.StatusOK, channelsDTO)
}

// POST /api/v1/channels/
func (h *HTTPHandler) handleCreateChannel(ctx *gin.Context) {
	logger := h.logger.WithMethod("handleCreateChannel")
	logger.Info("Creating channel")

	creatorID := ctx.GetHeader("X-User-ID")
	if creatorID == "" {
		ctx.JSON(http.StatusUnauthorized, errorResponse(ErrUnauthorized))
		return
	}
	var req CreateChannelRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	creatorUserID, err := uuid.Parse(creatorID)
	if err != nil {
		logger.Error("Failed to parse creator user ID", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	channel, err := h.channelService.HandleCreateChannel(ctx, domain.CreateChannelCommand{
		CreatorUserID: creatorUserID,
		Name:          req.Name,
		Topic:         req.Topic,
	})
	if err != nil {
		logger.Error("Failed to create channel", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	cacheKey := fmt.Sprintf("user_channels:%s", creatorUserID.String())
	h.cache.Delete(ctx.Request.Context(), cacheKey)

	channelDTO, err := h.channelService.ReturnChannelDTO(ctx, channel)
	if err != nil {
		logger.Error("Failed to return channel DTO", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusCreated, channelDTO)
}

// GET /api/v1/channels/:channelId
func (h *HTTPHandler) handleGetChannel(ctx *gin.Context) {
	logger := h.logger.WithMethod("handleGetChannel")
	logger.Info("Getting channel")

	var req ChannelIDUri
	if err := ctx.ShouldBindUri(&req); err != nil {
		logger.Error("Failed to bind URI", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	channelId, err := uuid.Parse(req.ChannelID)
	if err != nil {
		logger.Error("Failed to parse channel ID", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	cacheKey := fmt.Sprintf("channel:%s", channelId.String())
	var cachedChannel domain.ChannelDTO
	if hit, _ := h.cache.GetWithMetrics(ctx.Request.Context(), cacheKey, &cachedChannel); hit {
		logger.Info("Channel retrieved from cache", zap.String("channel_id", channelId.String()))
		ctx.JSON(http.StatusOK, cachedChannel)
		return
	}

	channel, err := h.channelService.HandleGetChannel(ctx, domain.GetChannelCommand{
		ChannelID: channelId,
	})
	if err != nil {
		logger.Error("Failed to get channel", zap.Error(err))
		ctx.JSON(http.StatusNotFound, errorResponse(err))
		return
	}

	channelDTO, err := h.channelService.ReturnChannelDTO(ctx, channel)
	if err != nil {
		logger.Error("Failed to return channel DTO", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	h.cache.Set(ctx.Request.Context(), cacheKey, *channelDTO, 10*time.Minute)
	logger.Info("Channel retrieved", zap.String("channel_id", channel.ID.String()))
	ctx.JSON(http.StatusOK, channelDTO)
}

// POST /api/v1/channels/:channelId/bots
func (h *HTTPHandler) handleAddBotToChannel(ctx *gin.Context) {
	logger := h.logger.WithMethod("handleAddBotToChannel")
	logger.Info("Adding bot to channel")

	var channelIdUri ChannelIDUri
	if err := ctx.ShouldBindUri(&channelIdUri); err != nil {
		logger.Error("Failed to bind URI", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var req AddBotToChannelRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		logger.Error("Failed to bind JSON", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	channelId, err := uuid.Parse(channelIdUri.ChannelID)
	if err != nil {
		logger.Error("Failed to parse channel ID", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	integrationId, err := uuid.Parse(req.IntegrationID)
	if err != nil {
		logger.Error("Failed to parse integration ID", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	channel, err := h.channelService.HandleAddBotToChannel(ctx, domain.AddBotToChannelCommand{
		ChannelID:     channelId,
		IntegrationID: integrationId,
	})
	if err != nil {
		logger.Error("Failed to add bot to channel", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	channelDTO, err := h.channelService.ReturnChannelDTO(ctx, channel)
	if err != nil {
		logger.Error("Failed to return channel DTO", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	logger.Info("Bot added to channel", zap.String("channel_id", channel.ID.String()))
	ctx.JSON(http.StatusOK, channelDTO)
}

// PUT /api/v1/channels/:channelId/archive
func (h *HTTPHandler) handleArchiveChannel(ctx *gin.Context) {
	logger := h.logger.WithMethod("handleArchiveChannel")
	logger.Info("Archiving channel")

	var req ChannelIDUri
	if err := ctx.ShouldBindUri(&req); err != nil {
		logger.Error("Failed to bind URI", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	// TODO: implement authentication
	ctx.Status(http.StatusOK)
}

// PUT /api/v1/channels/:channelId/unarchive
func (h *HTTPHandler) handleUnarchiveChannel(ctx *gin.Context) {
	logger := h.logger.WithMethod("handleUnarchiveChannel")
	logger.Info("Unarchiving channel")

	var req ChannelIDUri
	if err := ctx.ShouldBindUri(&req); err != nil {
		logger.Error("Failed to bind URI", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// TODO: implement authentication

	ctx.Status(http.StatusOK)
}

// POST /api/v1/channels/:channelId/join
func (h *HTTPHandler) handleJoinChannel(ctx *gin.Context) {
	logger := h.logger.WithMethod("handleJoinChannel")
	logger.Info("Joining channel")

	var req JoinChannelRequest
	var uriReq ChannelIDUri

	if err := ctx.ShouldBindJSON(&req); err != nil {
		logger.Error("Failed to bind JSON", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	if err := ctx.ShouldBindUri(&uriReq); err != nil {
		logger.Error("Failed to bind URI", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	userId, err := uuid.Parse(req.UserID)
	if err != nil {
		logger.Error("Failed to parse user ID", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	channelId, err := uuid.Parse(uriReq.ChannelID)
	if err != nil {
		logger.Error("Failed to parse channel ID", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	channel, err := h.channelService.HandleJoinChannel(ctx, domain.JoinChannelCommand{
		ChannelID: channelId,
		UserID:    userId,
	})
	if err != nil {
		logger.Error("Failed to join channel", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	channelDTO, err := h.channelService.ReturnChannelDTO(ctx, channel)
	if err != nil {
		logger.Error("Failed to return channel DTO", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	cacheKey := fmt.Sprintf("user_channels:%s", userId.String())
	h.cache.Delete(ctx.Request.Context(), cacheKey)

	logger.Info("Channel joined", zap.String("channel_id", channel.ID.String()))
	ctx.JSON(http.StatusOK, channelDTO)
}

// POST /api/v1/channels/:channelId/messages
// FIXME: redundant, should be removed
func (h *HTTPHandler) handleSendMessage(ctx *gin.Context) {
	logger := h.logger.WithMethod("handleSendMessage")
	logger.Info("Sending message")

	userID := ctx.GetHeader("X-User-ID")
	if userID == "" {
		ctx.JSON(http.StatusUnauthorized, errorResponse(ErrUnauthorized))
		return
	}
	var req SendMessageRequest
	var uriReq ChannelIDUri

	if err := ctx.ShouldBindJSON(&req); err != nil {
		logger.Error("Failed to bind JSON", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	if err := ctx.ShouldBindUri(&uriReq); err != nil {
		logger.Error("Failed to bind URI", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	channelId, err := uuid.Parse(uriReq.ChannelID)
	if err != nil {
		logger.Error("Failed to parse channel ID", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	senderID, err := uuid.Parse(userID)
	if err != nil {
		logger.Error("Failed to parse sender user ID", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var parentMessageID *uuid.UUID
	if req.ParentMessageID != nil {
		parsed, err := uuid.Parse(*req.ParentMessageID)
		if err != nil {
			logger.Error("Failed to parse parent message ID", zap.Error(err))
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
		logger.Error("Failed to send message", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	logger.Info("Message sent", zap.String("message_id", message.GetId().String()))
	ctx.JSON(http.StatusCreated, messageDTO)
}

// GET /api/v1/channels/:channelId/messages
func (h *HTTPHandler) handleGetMessages(ctx *gin.Context) {
	logger := h.logger.WithMethod("handleGetMessages")
	logger.Info("Getting messages")

	var uriReq ChannelIDUri

	if err := ctx.ShouldBindUri(&uriReq); err != nil {
		logger.Error("Failed to bind URI", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	channelId, err := uuid.Parse(uriReq.ChannelID)
	if err != nil {
		logger.Error("Failed to parse channel ID", zap.Error(err))
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
		logger.Error("Failed to list messages", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	messagesDTO, err := h.messageService.ToMessageDTOs(ctx, messages)
	if err != nil {
		logger.Error("Failed to convert messages to DTOs", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	logger.Info("Messages retrieved", zap.String("channel_id", channelId.String()))
	ctx.JSON(http.StatusOK, messagesDTO)
}

// POST /api/v1/channels/:channelId/messages/:messageId/reactions
func (h *HTTPHandler) handleAddReaction(ctx *gin.Context) {
	logger := h.logger.WithMethod("handleAddReaction")
	logger.Info("Adding reaction")

	var req AddReactionRequest
	var channelIdUri ChannelIDUri
	var messageIdUri MessageIDUri

	if err := ctx.ShouldBindJSON(&req); err != nil {
		logger.Error("Failed to bind JSON", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	if err := ctx.ShouldBindUri(&channelIdUri); err != nil {
		logger.Error("Failed to bind URI", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	if err := ctx.ShouldBindUri(&messageIdUri); err != nil {
		logger.Error("Failed to bind URI", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	userId, err := uuid.Parse(req.UserID)
	if err != nil {
		logger.Error("Failed to parse user ID", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	channelId, err := uuid.Parse(channelIdUri.ChannelID)
	if err != nil {
		logger.Error("Failed to parse channel ID", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	messageId, err := uuid.Parse(messageIdUri.MessageID)
	if err != nil {
		logger.Error("Failed to parse message ID", zap.Error(err))
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
		logger.Error("Failed to add reaction", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	logger.Info("Reaction added", zap.String("reaction_id", reaction.GetId().String()))
	//TODO: return the reaction as a DTO? or return the message with the reaction?
	ctx.JSON(http.StatusCreated, reaction)
}

// DELETE /api/v1/channels/:channelId/messages/:messageId/reactions
func (h *HTTPHandler) handleRemoveReaction(ctx *gin.Context) {
	logger := h.logger.WithMethod("handleRemoveReaction")
	logger.Info("Removing reaction")

	var req RemoveReactionRequest
	var channelIdUri ChannelIDUri
	var messageIdUri MessageIDUri

	if err := ctx.ShouldBindJSON(&req); err != nil {
		logger.Error("Failed to bind JSON", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	if err := ctx.ShouldBindUri(&channelIdUri); err != nil {
		logger.Error("Failed to bind URI", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	if err := ctx.ShouldBindUri(&messageIdUri); err != nil {
		logger.Error("Failed to bind URI", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	userId, err := uuid.Parse(req.UserID)
	if err != nil {
		logger.Error("Failed to parse user ID", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	channelId, err := uuid.Parse(channelIdUri.ChannelID)
	if err != nil {
		logger.Error("Failed to parse channel ID", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	messageId, err := uuid.Parse(messageIdUri.MessageID)
	if err != nil {
		logger.Error("Failed to parse message ID", zap.Error(err))
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
		logger.Error("Failed to remove reaction", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	logger.Info("Reaction removed", zap.String("reaction_id", cmd.MessageID.String()))
	ctx.Status(http.StatusOK)
}

// POST /api/v1/channels/:channelId/invites
func (h *HTTPHandler) handleCreateChannelInvite(ctx *gin.Context) {
	logger := h.logger.WithMethod("handleCreateChannelInvite")
	logger.Info("Creating channel invite")

	var channelIdUri ChannelIDUri
	if err := ctx.ShouldBindUri(&channelIdUri); err != nil {
		logger.Error("Failed to bind URI", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	channelId, err := uuid.Parse(channelIdUri.ChannelID)
	if err != nil {
		logger.Error("Failed to parse channel ID", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	var req CreateChannelInviteRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		logger.Error("Failed to bind JSON", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	creatorID := ctx.GetHeader("X-User-ID")
	if creatorID == "" {
		logger.Error("Unauthorized")
		ctx.JSON(http.StatusUnauthorized, errorResponse(ErrUnauthorized))
		return
	}

	creatorUserID, err := uuid.Parse(creatorID)
	if err != nil {
		logger.Error("Failed to parse creator user ID", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var maxUses *int
	if req.MaxUses != 0 {
		maxUses = &req.MaxUses
	}

	cmd := domain.CreateChannelInviteCommand{
		ChannelID:       channelId,
		CreatedByUserID: creatorUserID,
		ExpiresAt:       req.ExpiresAt,
		MaxUses:         maxUses,
	}

	_, invite, err := h.channelService.HandleCreateChannelInvite(ctx, cmd)
	if err != nil {
		logger.Error("Failed to create channel invite", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	inviteDTO := domain.ToChannelInviteDTO(invite)

	logger.Info("Channel invite created", zap.String("invite_id", invite.ID.String()))
	ctx.JSON(http.StatusCreated, inviteDTO)
}

// GET /api/v1/channels/:channelId/invites
func (h *HTTPHandler) handleGetChannelInvites(ctx *gin.Context) {
	logger := h.logger.WithMethod("handleGetChannelInvites")
	logger.Info("Getting channel invites")

	var channelIdUri ChannelIDUri
	if err := ctx.ShouldBindUri(&channelIdUri); err != nil {
		logger.Error("Failed to bind URI", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	channelId, err := uuid.Parse(channelIdUri.ChannelID)
	if err != nil {
		logger.Error("Failed to parse channel ID", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	cmd := domain.GetChannelInvitesCommand{
		ChannelID: channelId,
	}

	channel, err := h.channelService.HandleGetChannelInvites(ctx, cmd)
	if err != nil {
		logger.Error("Failed to get channel invites", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	invites := channel.GetActiveInvites()
	invitesDTO := make([]domain.ChannelInviteDTO, len(invites))
	for i, invite := range invites {
		invitesDTO[i] = domain.ToChannelInviteDTO(&invite)
	}

	logger.Info("Channel invites retrieved", zap.String("channel_id", channel.ID.String()))
	ctx.JSON(http.StatusOK, invitesDTO)
}

// POST /api/v1/invites/accept
func (h *HTTPHandler) handleAcceptChannelInvite(ctx *gin.Context) {
	logger := h.logger.WithMethod("handleAcceptChannelInvite")
	logger.Info("Accepting channel invite")

	var req AcceptChannelInviteRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		logger.Error("Failed to bind JSON", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	userIdStr := ctx.GetHeader("X-User-ID")
	if userIdStr == "" {
		logger.Error("Unauthorized")
		ctx.JSON(http.StatusUnauthorized, errorResponse(ErrUnauthorized))
		return
	}

	userId, err := uuid.Parse(userIdStr)
	if err != nil {
		logger.Error("Failed to parse user ID", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	cmd := domain.AcceptChannelInviteCommand{
		InviteCode: req.InviteCode,
		UserID:     userId,
	}

	channel, err := h.channelService.HandleAcceptChannelInvite(ctx, cmd)
	if err != nil {
		logger.Error("Failed to accept channel invite", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	channelDTO, err := h.channelService.ReturnChannelDTO(ctx, channel)
	if err != nil {
		logger.Error("Failed to return channel DTO", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	logger.Info("Channel invite accepted", zap.String("channel_id", channel.ID.String()))
	cacheKey := fmt.Sprintf("user_channels:%s", userId.String())
	h.cache.Delete(ctx.Request.Context(), cacheKey)

	ctx.JSON(http.StatusOK, channelDTO)
}

// DELETE /api/v1/invites/:inviteId
func (h *HTTPHandler) handleDeactivateChannelInvite(ctx *gin.Context) {
	logger := h.logger.WithMethod("handleDeactivateChannelInvite")
	logger.Info("Deactivating channel invite")

	var inviteIdUri InvideIDUri
	if err := ctx.ShouldBindUri(&inviteIdUri); err != nil {
		logger.Error("Failed to bind URI", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	inviteId, err := uuid.Parse(inviteIdUri.InvideID)
	if err != nil {
		logger.Error("Failed to parse invite ID", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	userIDStr := ctx.GetHeader("X-User-ID")
	if userIDStr == "" {
		logger.Error("Unauthorized")
		ctx.JSON(http.StatusUnauthorized, errorResponse(ErrUnauthorized))
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		logger.Error("Failed to parse user ID", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	cmd := domain.DeactivateChannelInviteCommand{
		InviteID: inviteId,
		UserID:   userID,
	}

	_, err = h.channelService.HandleDeactivateChannelInvite(ctx, cmd)
	if err != nil {
		logger.Error("Failed to deactivate channel invite", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	logger.Info("Channel invite deactivated", zap.String("invite_id", inviteId.String()))
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
