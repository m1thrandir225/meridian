package handlers

import (
	"time"
)

type WebSocketMessage struct {
	Type    string      `json:"type"`
	Payload interface{} `json:"payload"`
}

type IncomingMessagePayload struct {
	Content         string `json:"content"`
	ChannelID       string `json:"channel_id"`
	ParentMessageID string `json:"parent_message_id,omitempty"`
}

type OutgoingMessagePayload struct {
	ID              string             `json:"id"`
	Content         string             `json:"content"`
	SenderUserID    string             `json:"sender_user_id"`
	IntegrationID   string             `json:"integration_id"`
	ChannelID       string             `json:"channel_id"`
	ParentMessageID string             `json:"parent_message_id,omitempty"`
	Timestamp       time.Time          `json:"timestamp"`
	SenderUser      *UserDTO           `json:"sender_user,omitempty"`
	IntegrationBot  *IntegrationBotDTO `json:"integration_bot,omitempty"`
}

type IncomingReactionPayload struct {
	MessageID    string `json:"message_id"`
	ChannelID    string `json:"channel_id"`
	ReactionType string `json:"reaction_type"`
}

type OutgoingReactionPayload struct {
	ID           string    `json:"id"`
	MessageID    string    `json:"message_id"`
	ChannelID    string    `json:"channel_id"`
	UserID       string    `json:"user_id"`
	ReactionType string    `json:"reaction_type"`
	Timestamp    time.Time `json:"timestamp"`
}

type UserDTO struct {
	ID        string `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type IntegrationBotDTO struct {
	ID          string    `json:"id"`
	ServiceName string    `json:"service_name"`
	CreatedAt   time.Time `json:"created_at"`
	IsRevoked   bool      `json:"is_revoked"`
}

type TypingPayload struct {
	ChannelID string `json:"channel_id"`
	UserID    string `json:"user_id"`
}
