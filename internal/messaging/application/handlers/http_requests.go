package handlers

import (
	"time"

	"github.com/gin-gonic/gin"
)

type ChannelIDUri struct {
	ChannelID string `uri:"channelId" binding:"required,uuid"`
}

type MessageIDUri struct {
	MessageID string `uri:"messageId" binding:"required,uuid"`
}
type InvideIDUri struct {
	InvideID string `uri:"inviteId" binding:"required,uuid"`
}

type CreateChannelRequest struct {
	Name  string `json:"name"  binding:"required"`
	Topic string `json:"topic" `
}

type SendMessageRequest struct {
	ContentText          string  `json:"content_text" binding:"required"`
	IsIntegrationMessage *bool   `json:"is_integration_message" binding:"required"`
	ParentMessageID      *string `json:"parent_message_id,omitempty" binding:"omitempty"`
}

type JoinChannelRequest struct {
	UserID string `json:"user_id" binding:"required,uuid"`
}

type AddReactionRequest struct {
	UserID       string `json:"user_id" binding:"required,uuid"`
	ReactionType string `json:"reaction_type" binding:"required"`
}

type RemoveReactionRequest struct {
	UserID       string `json:"user_id" binding:"required,uuid"`
	ReactionType string `json:"reaction_type" binding:"required"`
}

type AddBotToChannelRequest struct {
	IntegrationID string `json:"integration_id" binding:"required,uuid"`
}

type CreateChannelInviteRequest struct {
	ExpiresAt time.Time `json:"expires_at" binding:"required"`
	MaxUses   int       `json:"max_uses,omitempty"`
}

type AcceptChannelInviteRequest struct {
	InviteCode string `json:"invite_code" binding:"required"`
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
