package handlers

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/m1thrandir225/meridian/internal/messaging/domain"
)

type ChannelIDUri struct {
	ChannelID string `uri:"channelId" binding:"required,uuid"`
}

type MessageIDUri struct {
	MessageID string `uri:"messageId" binding:"required,uuid"`
}

type CreateChannelRequest struct {
	Name  string `json:"name"  binding:"required"`
	Topic string `json:"topic" `
}

type ChannelResponse struct {
	ID              string    `json:"id"`
	Name            string    `json:"name"`
	Topic           string    `json:"topic"`
	CreatorUserID   string    `json:"creator_user_id"`
	CreationTime    time.Time `json:"creation_time"`
	LastMessageTime time.Time `json:"last_message_time"`
	IsArchived      bool      `json:"is_archived"`
	MembersCount    int       `json:"members_count"`
}

func ToChannelResponse(channel *domain.Channel) ChannelResponse {
	return ChannelResponse{
		ID:              channel.ID.String(),
		Name:            channel.Name,
		Topic:           channel.Topic,
		CreatorUserID:   channel.CreatorUserID.String(),
		CreationTime:    channel.CreationTime,
		LastMessageTime: channel.LastMessageTime,
		IsArchived:      channel.IsArchived,
		MembersCount:    len(channel.Members),
	}
}

type SendMessageRequest struct {
	ContentText          string  `json:"content_text" binding:"required"`
	IsIntegrationMessage *bool   `json:"is_integration_message" binding:"required"`
	ParentMessageID      *string `json:"parent_message_id,omitempty" binding:"omitempty"`
}

type MessageResponse struct {
	ID              string    `json:"id"`
	ChannelID       string    `json:"channel_id"`
	SenderUserID    *string   `json:"sender_user_id,omitempty"`
	IntegrationID   *string   `json:"integration_id,omitempty"`
	ContentText     string    `json:"content_text"`
	CreatedAt       time.Time `json:"created_at"`
	ParentMessageID *string
	SenderUser      *UserResponse      `json:"sender_user,omitempty"`
	Reactions       []ReactionResponse `json:"reactions,omitempty"`
}
type UserResponse struct {
	ID        string `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

func ToUserResponse(user *domain.User) UserResponse {
	return UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
	}
}

func ToMessageResponse(message *domain.Message) MessageResponse {
	var senderId, integrationId, parentId *string
	if message.GetSenderUserId() != nil {
		sId := message.GetSenderUserId().String()
		senderId = &sId
	}

	if message.GetIntegrationId() != nil {
		iId := message.GetIntegrationId().String()
		integrationId = &iId
	}

	if message.GetParentMessageId() != nil {
		pId := message.GetParentMessageId().String()
		parentId = &pId
	}

	reactionsDTO := make([]ReactionResponse, len(message.GetReactions()))
	for i, reaction := range message.GetReactions() {
		reactionsDTO[i] = ToReactionResponse(&reaction)
	}

	var senderUser *UserResponse
	if message.GetSenderUser() != nil {
		sender := message.GetSenderUser()
		senderUser = &UserResponse{
			ID:        sender.ID,
			Username:  sender.Username,
			Email:     sender.Email,
			FirstName: sender.FirstName,
			LastName:  sender.LastName,
		}
	}

	return MessageResponse{
		ID:              message.GetId().String(),
		ChannelID:       message.GetChannelId().String(),
		SenderUserID:    senderId,
		IntegrationID:   integrationId,
		ContentText:     message.GetContent().GetText(),
		CreatedAt:       message.GetCreatedAt(),
		ParentMessageID: parentId,
		Reactions:       reactionsDTO,
		SenderUser:      senderUser,
	}
}

type JoinChannelRequest struct {
	UserID string `json:"user_id" binding:"required,uuid"`
}

type AddReactionRequest struct {
	UserID       string `json:"user_id" binding:"required,uuid"`
	ReactionType string `json:"reaction_type" binding:"required"`
}

type ReactionResponse struct {
	ID           string    `json:"id"`
	MessageID    string    `json:"message_id"`
	UserID       string    `json:"user_id"`
	ReactionType string    `json:"reaction_type"`
	Timestamp    time.Time `json:"timestamp"`
}

func ToReactionResponse(reaction *domain.Reaction) ReactionResponse {
	return ReactionResponse{
		ID:           reaction.GetId().String(),
		MessageID:    reaction.GetMessageId().String(),
		UserID:       reaction.GetUserId().String(),
		ReactionType: reaction.GetReactionType(),
		Timestamp:    reaction.GetCreatedAt(),
	}
}

type RemoveReactionRequest struct {
	UserID       string `json:"user_id" binding:"required,uuid"`
	ReactionType string `json:"reaction_type" binding:"required"`
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
