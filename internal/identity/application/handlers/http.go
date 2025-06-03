package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/m1thrandir225/meridian/internal/identity/application/services"
)

type HTTPHandler struct {
	userService *services.UserService
}

// POST /api/v1/register
func (h *HTTPHandler) Register(ctx *gin.Context) {}

// POST /api/v1/login
func (h *HTTPHandler) Login(ctx *gin.Context) {}

// GET /api/v1/me
func (h *HTTPHandler) GetCurrentUser(ctx *gin.Context) {}

// PUT /api/v1/update-profile

func (h *HTTPHandler) UpdateCurrentUser(ctx *gin.Context) {}

// DELETE /api/v1/me
func (h *HTTPHandler) DeleteCurrentUser(ctx *gin.Context) {}
