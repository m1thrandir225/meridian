package handlers

import "time"

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
	ID              string    `json:"id"`
	Content         string    `json:"content"`
	SenderID        string    `json:"sender_id"`
	ChannelID       string    `json:"channel_id"`
	ParentMessageID string    `json:"parent_message_id,omitempty"`
	Timestamp       time.Time `json:"timestamp"`
}

type TypingPayload struct {
	ChannelID string `json:"channel_id"`
	UserID    string `json:"user_id"`
}
