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

func (h *HTTPHandler) CreateChannel(ctx *gin.Context) {
	ctx.Status(http.StatusCreated)
}

// TODO:: implement
func (h *HTTPHandler) GetChannel(ctx *gin.Context) {
	ctx.Status(http.StatusOK)
}

// TODO: implement
func (h *HTTPHandler) JoinChannel(ctx *gin.Context) {
	ctx.Status(http.StatusOK)
}

// TODO: implement
func (h *HTTPHandler) SendMessage(ctx *gin.Context) {
	ctx.Status(http.StatusOK)
}
