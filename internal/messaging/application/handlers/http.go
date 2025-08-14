package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/m1thrandir225/meridian/internal/messaging/application/services"
	"github.com/m1thrandir225/meridian/internal/messaging/domain"
)

type HTTPHandler struct {
	channelService *services.ChannelService
}

func NewHttpHandler(service *services.ChannelService) *HTTPHandler {
	return &HTTPHandler{
		channelService: service,
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

	cmd := domain.GetUserChannelsCommand{
		UserID: userID,
	}

	channels, err := h.channelService.HandleGetUserChannels(ctx, cmd)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "something went wrong: " + err.Error()})
		return
	}

	channelsResponse := make([]ChannelResponse, len(channels))
	for i, channel := range channels {
		channelsResponse[i] = ToChannelResponse(channel)
	}

	ctx.JSON(http.StatusOK, channelsResponse)
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

	channelDTO := ToChannelResponse(channel)

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

	channel, err := h.channelService.HandleGetChannel(ctx, domain.GetChannelCommand{
		ChannelID: channelId,
	})
	if err != nil {
		ctx.JSON(http.StatusNotFound, errorResponse(err))
		return
	}

	channelDTO := ToChannelResponse(channel)

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

	channelDTO := ToChannelResponse(channel)

	ctx.JSON(http.StatusOK, channelDTO)
}

// POST /api/v1/channels/:channelId/messages
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

	message, err := h.channelService.HandleMessageSent(ctx, domain.SendMessageCommand{
		ChannelID:       channelId,
		SenderUserID:    senderID,
		ParentMessageID: parentMessageID,
		Content:         domain.MessageContent{},
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	messageDTO := ToMessageResponse(message)

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

	channel, err := h.channelService.HandleListMessages(ctx, domain.ListMessagesForChannelCommand{
		ChannelID: channelId,
		Limit:     50,
		Offset:    0,
	})
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	messages := channel.Messages

	messagesDTO := make([]MessageResponse, len(messages))
	for i, message := range messages {
		messagesDTO[i] = ToMessageResponse(&message)
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

	reaction, err := h.channelService.HandleAddReaction(ctx, domain.AddReactionCommand{
		ChannelID:    channelId,
		MessageID:    messageId,
		UserID:       userId,
		ReactionType: req.ReactionType,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	reactionDTO := ToReactionResponse(reaction)

	ctx.JSON(http.StatusCreated, reactionDTO)
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

	_, err = h.channelService.HandleRemoveReaction(ctx, domain.RemoveReactionCommand{
		ChannelID:    channelId,
		MessageID:    messageId,
		UserID:       userId,
		ReactionType: req.ReactionType,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.Status(http.StatusOK)
}
