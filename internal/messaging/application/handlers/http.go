package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/m1thrandir225/meridian/internal/messaging/application/services"
)

type HTTPHandler struct {
	channelService *services.ChannelService
}

func NewHttpHandler(service *services.ChannelService) *HTTPHandler {
	return &HTTPHandler{
		channelService: service,
	}
}

// POST /api/v1/channels/
func (h *HTTPHandler) CreateChannel(ctx *gin.Context) {
	var req CreateChannelRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	// TODO: implement
	ctx.Status(http.StatusCreated)
}

// GET /api/v1/channels/:channelId
func (h *HTTPHandler) GetChannel(ctx *gin.Context) {
	var req ChannelIDUri
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// TODO:: implement
	ctx.Status(http.StatusOK)
}

// PUT /api/v1/channels/:channelId/archive
func (h *HTTPHandler) ArchiveChannel(ctx *gin.Context) {
	var req ChannelIDUri
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
}

// PUT /api/v1/channels/:channelId/unarchive
func (h *HTTPHandler) UnarchiveChannel(ctx *gin.Context) {
	var req ChannelIDUri
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	ctx.Status(http.StatusOK)
}

// POST /api/v1/channels/:channelId/join
func (h *HTTPHandler) JoinChannel(ctx *gin.Context) {
	var req JoinChannelRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// TODO: implement
	ctx.Status(http.StatusOK)
}

// POST /api/v1/channels/:channelId/messages
func (h *HTTPHandler) SendMessage(ctx *gin.Context) {
	var req SendMessageRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// TODO: implement
	ctx.Status(http.StatusOK)
}

// POST /api/v1/channels/:channelId/messages/:messageId/reactions
func (h *HTTPHandler) AddReaction(ctx *gin.Context) {
	var req AddReactionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	// TODO: implement
	ctx.Status(http.StatusOK)
}

// DELETE /api/v1/channels/:channelId/messages/:messageId/reactions
func (h *HTTPHandler) RemoveReaction(ctx *gin.Context) {
	var req RemoveReactionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	// TODO: implement
	ctx.Status(http.StatusOK)
}
