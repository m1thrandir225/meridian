package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/m1thrandir225/meridian/internal/messaging/application/services"
	"github.com/m1thrandir225/meridian/internal/messaging/domain"
	"github.com/redis/go-redis/v9"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"
)

type WebSocketHandler struct {
	upgrader       websocket.Upgrader
	clients        map[string]*websocket.Conn
	mu             sync.RWMutex
	channelService *services.ChannelService
	redisClient    *redis.Client
}

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

func NewWebSocketHandler(channelService *services.ChannelService, redisClient *redis.Client) *WebSocketHandler {
	handler := &WebSocketHandler{
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true //TODO fix for production
			},
		},
		clients:        make(map[string]*websocket.Conn),
		channelService: channelService,
		redisClient:    redisClient,
	}

	if redisClient != nil {
		go handler.subscribeToRedisMessages()
	}

	return handler
}

func (h *WebSocketHandler) HandleWebSocket(c *gin.Context) {
	userID := c.GetHeader("X-User-Id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	conn, err := h.upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("Failed to upgrade connection: %v", err)
		return
	}
	defer conn.Close()

	h.addClient(userID, conn)
	defer h.removeClient(userID)

	log.Printf("WebSocket connection established for user: %s", userID)

	// Send connection confirmation
	h.sendToClient(userID, WebSocketMessage{
		Type:    "connected",
		Payload: map[string]string{"user_id": userID, "timestamp": time.Now().UTC().Format(time.RFC3339)},
	})

	// Handle incoming messages
	for {
		var msg WebSocketMessage
		err := conn.ReadJSON(&msg)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket error: %v", err)
			}
			break
		}

		// Process different message types
		switch msg.Type {
		case "ping":
			h.sendToClient(userID, WebSocketMessage{
				Type:    "pong",
				Payload: map[string]string{"timestamp": time.Now().UTC().Format(time.RFC3339)},
			})
		case "message":
			err := h.handleIncomingMessage(userID, msg.Payload)
			if err != nil {
				log.Printf("Failed to handle message from user %s: %v", userID, err)
				h.sendToClient(userID, WebSocketMessage{
					Type:    "error",
					Payload: map[string]string{"message": "Failed to send message", "error": err.Error()},
				})
			}
		case "typing_start":
			h.handleTypingIndicator(userID, msg.Payload, "typing_start")
		case "typing_stop":
			h.handleTypingIndicator(userID, msg.Payload, "typing_stop")
		default:
			log.Printf("Unknown message type: %s from user: %s", msg.Type, userID)
		}
	}
}

func (h *WebSocketHandler) handleIncomingMessage(senderID string, payload interface{}) error {
	// Parse the message payload
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	var incomingMsg IncomingMessagePayload
	if err := json.Unmarshal(payloadBytes, &incomingMsg); err != nil {
		return err
	}

	// Validate required fields
	if incomingMsg.ChannelID == "" {
		return fmt.Errorf("channel_id is required")
	}
	if incomingMsg.Content == "" {
		return fmt.Errorf("content is required")
	}

	// Parse UUIDs
	senderUUID, err := uuid.Parse(senderID)
	if err != nil {
		return fmt.Errorf("invalid sender ID: %w", err)
	}

	channelUUID, err := uuid.Parse(incomingMsg.ChannelID)
	if err != nil {
		return fmt.Errorf("invalid channel ID: %w", err)
	}

	var parentMessageUUID *uuid.UUID
	if incomingMsg.ParentMessageID != "" {
		parentUUID, err := uuid.Parse(incomingMsg.ParentMessageID)
		if err != nil {
			return fmt.Errorf("invalid parent message ID: %w", err)
		}
		parentMessageUUID = &parentUUID
	}

	messageContent := domain.NewMessageContent(incomingMsg.Content)

	cmd := domain.SendMessageCommand{
		ChannelID:       channelUUID,
		SenderUserID:    senderUUID,
		Content:         messageContent,
		ParentMessageID: parentMessageUUID,
	}

	// Handle through domain service
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	message, err := h.channelService.HandleMessageSent(ctx, cmd)
	if err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}

	outgoingMsg := OutgoingMessagePayload{
		ID:        message.GetId().String(),
		Content:   message.GetContent().GetText(),
		SenderID:  message.GetSenderUserId().String(),
		ChannelID: message.GetChannelId().String(),
		Timestamp: message.GetCreatedAt(),
	}

	if message.GetParentMessageId() != nil {
		outgoingMsg.ParentMessageID = message.GetParentMessageId().String()
	}

	if h.redisClient != nil {
		go h.publishMessageToRedis(outgoingMsg)
	} else {
		go h.broadcastToChannel(incomingMsg.ChannelID, WebSocketMessage{
			Type:    "new_message",
			Payload: outgoingMsg,
		})
	}

	return nil
}

func (h *WebSocketHandler) handleTypingIndicator(userID string, payload interface{}, typingType string) {
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Failed to marshal typing payload: %v", err)
		return
	}

	var typingPayload TypingPayload
	if err := json.Unmarshal(payloadBytes, &typingPayload); err != nil {
		log.Printf("Failed to unmarshal typing payload: %v", err)
		return
	}

	if typingPayload.ChannelID == "" {
		return
	}

	// Set user ID from authenticated user
	typingPayload.UserID = userID

	typingMsg := WebSocketMessage{
		Type:    typingType,
		Payload: typingPayload,
	}

	// Publish typing indicator via Redis (ephemeral)
	if h.redisClient != nil {
		go h.publishTypingToRedis(typingPayload.ChannelID, typingMsg)
	} else {
		// Fallback: broadcast directly to connected clients
		go h.broadcastToChannel(typingPayload.ChannelID, typingMsg)
	}
}

func (h *WebSocketHandler) addClient(userID string, conn *websocket.Conn) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.clients[userID] = conn
}

func (h *WebSocketHandler) removeClient(userID string) {
	h.mu.Lock()
	defer h.mu.Unlock()
	delete(h.clients, userID)
}

func (h *WebSocketHandler) sendToClient(userID string, message WebSocketMessage) error {
	h.mu.RLock()
	conn, exists := h.clients[userID]
	h.mu.RUnlock()

	if !exists {
		return nil // Client not connected
	}

	return conn.WriteJSON(message)
}

func (h *WebSocketHandler) publishMessageToRedis(message OutgoingMessagePayload) {
	if h.redisClient == nil {
		return
	}

	ctx := context.Background()

	messageJSON, err := json.Marshal(WebSocketMessage{
		Type:    "new_message",
		Payload: message,
	})

	if err != nil {
		log.Printf("Failed to marshal message: %v", err)
		return
	}

	channelKey := fmt.Sprintf("channel:%s", message.ChannelID)
	err = h.redisClient.Publish(ctx, channelKey, messageJSON).Err()
	if err != nil {
		log.Printf("Failed to publish message to Redis: %v", err)
	}

	messageKey := fmt.Sprintf("message:%s", message.ID)
	messageCacheJSON, _ := json.Marshal(message)
	h.redisClient.Set(ctx, messageKey, messageCacheJSON, 24*time.Hour)

	recentKey := fmt.Sprintf("channel:%s:recent", message.ChannelID)
	h.redisClient.LPush(ctx, recentKey, message.ID)
	h.redisClient.LTrim(ctx, recentKey, 0, 100)
	h.redisClient.Expire(ctx, recentKey, 24*time.Hour)

}

func (h *WebSocketHandler) publishTypingToRedis(channelID string, typingMsg WebSocketMessage) {
	if h.redisClient == nil {
		return
	}

	ctx := context.Background()

	typingJSON, err := json.Marshal(typingMsg)
	if err != nil {
		log.Printf("Failed to marshal typing message: %v", err)
		return
	}

	channelKey := fmt.Sprintf("channel:%s", channelID)
	err = h.redisClient.Publish(ctx, channelKey, typingJSON).Err()
	if err != nil {
		log.Printf("Failed to publish typing message to Redis: %v", err)
	}
}

func (h *WebSocketHandler) subscribeToRedisMessages() {
	if h.redisClient == nil {
		return
	}

	ctx := context.Background()

	pubsub := h.redisClient.PSubscribe(ctx, "channel:*")
	defer pubsub.Close()

	ch := pubsub.Channel()
	for msg := range ch {
		var wsMessage WebSocketMessage
		if err := json.Unmarshal([]byte(msg.Payload), &wsMessage); err != nil {
			log.Printf("Failed to unmarshal message: %v", err)
			continue
		}
		if strings.HasPrefix(msg.Channel, "channel:") {
			channelID := strings.TrimPrefix(msg.Channel, "channel:")
			h.broadcastToChannel(channelID, wsMessage)
		}
	}
}

func (h *WebSocketHandler) broadcastToChannel(channelID string, message WebSocketMessage) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	for userID, conn := range h.clients {
		err := conn.WriteJSON(message)
		//TODO: check if the current user is a member of the channel
		if err != nil {
			log.Printf("Failed to send message to user %s: %v", userID, err)
			conn.Close()
			delete(h.clients, userID)
		}
	}
}

func (h *WebSocketHandler) BroadcastMessage(message *domain.Message) {
	outgoingMsg := OutgoingMessagePayload{
		ID:        message.GetId().String(),
		Content:   message.GetContent().GetText(),
		ChannelID: message.GetChannelId().String(),
		Timestamp: message.GetCreatedAt(),
	}

	if message.GetSenderUserId() != nil {
		outgoingMsg.SenderID = message.GetSenderUserId().String()
	}

	if message.GetParentMessageId() != nil {
		outgoingMsg.ParentMessageID = message.GetParentMessageId().String()
	}

	wsMessage := WebSocketMessage{
		Type:    "new_message",
		Payload: outgoingMsg,
	}

	if h.redisClient != nil {
		h.publishMessageToRedis(outgoingMsg)
	} else {
		h.broadcastToChannel(message.GetChannelId().String(), wsMessage)
	}
}

func (h *WebSocketHandler) SendToUser(userID string, message WebSocketMessage) error {
	return h.sendToClient(userID, message)
}
