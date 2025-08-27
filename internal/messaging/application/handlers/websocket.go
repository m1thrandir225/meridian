package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/m1thrandir225/meridian/internal/messaging/application/services"
	"github.com/m1thrandir225/meridian/internal/messaging/domain"
	"github.com/m1thrandir225/meridian/pkg/logging"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type WebSocketHandler struct {
	upgrader       websocket.Upgrader
	clients        map[string]*websocket.Conn
	mu             sync.RWMutex
	channelService *services.ChannelService
	messageService *services.MessageService
	redisClient    *redis.Client
	identityClient *services.IdentityClient
	logger         *logging.Logger
}

func NewWebSocketHandler(
	channelService *services.ChannelService,
	messageService *services.MessageService,
	redisClient *redis.Client,
	identityClient *services.IdentityClient,
	logger *logging.Logger,
) *WebSocketHandler {
	handler := &WebSocketHandler{
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true //TODO fix for production
			},
		},
		clients:        make(map[string]*websocket.Conn),
		channelService: channelService,
		messageService: messageService,
		redisClient:    redisClient,
		identityClient: identityClient,
		logger:         logger,
	}

	if redisClient != nil {
		go handler.subscribeToRedisMessages()
	}

	return handler
}

func (h *WebSocketHandler) HandleWebSocket(c *gin.Context) {
	logger := h.logger.WithMethod("HandleWebSocket")
	logger.Info("Handling WebSocket")

	token := c.Query("token")
	if token == "" {
		h.logger.Error("No token provided")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	userID, err := h.validateToken(c.Request.Context(), token)
	if err != nil {
		logger.Error("Failed to validate token", zap.Error(err))
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	conn, err := h.upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		logger.Error("Failed to upgrade connection", zap.Error(err))
		return
	}
	defer conn.Close()

	h.addClient(userID, conn)
	defer h.removeClient(userID)

	logger.Info("WebSocket connection established", zap.String("user_id", userID))

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
				logger.Error("WebSocket error", zap.Error(err))
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
				logger.Error("Failed to handle message from user", zap.String("user_id", userID), zap.Error(err))
				h.sendToClient(userID, WebSocketMessage{
					Type:    "error",
					Payload: map[string]string{"message": "Failed to send message", "error": err.Error()},
				})
			}
		case "add_reaction":
			err := h.handleIncomingReaction(userID, msg.Payload)
			if err != nil {
				logger.Error("Failed to handle reaction from user", zap.String("user_id", userID), zap.Error(err))
				h.sendToClient(userID, WebSocketMessage{
					Type:    "error",
					Payload: map[string]string{"message": "Failed to add reaction", "error": err.Error()},
				})
			}
		case "remove_reaction":
			err := h.handleRemoveReaction(userID, msg.Payload)
			if err != nil {
				logger.Error("Failed to handle reaction from user", zap.String("user_id", userID), zap.Error(err))
				h.sendToClient(userID, WebSocketMessage{
					Type:    "error",
					Payload: map[string]string{"message": "Failed to remove reaction", "error": err.Error()},
				})
			}
		case "typing_start":
			h.handleTypingIndicator(userID, msg.Payload, "typing_start")
		case "typing_stop":
			h.handleTypingIndicator(userID, msg.Payload, "typing_stop")
		default:
			logger.Error("Unknown message type", zap.String("message_type", msg.Type), zap.String("user_id", userID))
		}
	}
}

func (h *WebSocketHandler) handleIncomingMessage(senderID string, payload interface{}) error {
	logger := h.logger.WithMethod("handleIncomingMessage")
	logger.Info("Handling incoming message")

	// Parse the message payload
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		logger.Error("Failed to marshal payload", zap.Error(err))
		return err
	}

	var incomingMsg IncomingMessagePayload
	if err := json.Unmarshal(payloadBytes, &incomingMsg); err != nil {
		logger.Error("Failed to unmarshal payload", zap.Error(err))
		return err
	}

	// Validate required fields
	if incomingMsg.ChannelID == "" {
		logger.Error("Channel ID is required")
		return fmt.Errorf("channel_id is required")
	}
	if incomingMsg.Content == "" {
		logger.Error("Content is required")
		return fmt.Errorf("content is required")
	}

	// Parse UUIDs
	senderUUID, err := uuid.Parse(senderID)
	if err != nil {
		logger.Error("Invalid sender ID", zap.Error(err))
		return fmt.Errorf("invalid sender ID: %w", err)
	}

	channelUUID, err := uuid.Parse(incomingMsg.ChannelID)
	if err != nil {
		logger.Error("Invalid channel ID", zap.Error(err))
		return fmt.Errorf("invalid channel ID: %w", err)
	}

	var parentMessageUUID *uuid.UUID
	if incomingMsg.ParentMessageID != "" {
		parentUUID, err := uuid.Parse(incomingMsg.ParentMessageID)
		if err != nil {
			logger.Error("Invalid parent message ID", zap.Error(err))
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

	message, err := h.messageService.HandleMessageSent(ctx, cmd)
	if err != nil {
		logger.Error("Failed to send message", zap.Error(err))
		return fmt.Errorf("failed to send message: %w", err)
	}
	messageDTO, err := h.messageService.ToMessageDTO(ctx, message)
	if err != nil {
		logger.Error("Failed to convert message to DTO", zap.Error(err))
		return fmt.Errorf("failed to convert message to DTO: %w", err)
	}

	outgoingMsg := OutgoingMessagePayload{
		ID:        messageDTO.ID,
		Content:   messageDTO.ContentText,
		ChannelID: messageDTO.ChannelID,
		Timestamp: messageDTO.CreatedAt,
	}

	// Handle sender ID safely
	if messageDTO.SenderUserID != nil {
		outgoingMsg.SenderUserID = *messageDTO.SenderUserID
	} else if message.GetSenderUserId() != nil {
		outgoingMsg.SenderUserID = message.GetSenderUserId().String()
	}

	// Handle integration ID safely
	if messageDTO.IntegrationID != nil {
		outgoingMsg.IntegrationID = *messageDTO.IntegrationID
	} else if message.GetIntegrationId() != nil {
		outgoingMsg.IntegrationID = message.GetIntegrationId().String()
	}

	// Handle parent message ID safely
	if messageDTO.ParentMessageID != nil {
		outgoingMsg.ParentMessageID = *messageDTO.ParentMessageID
	} else if message.GetParentMessageId() != nil {
		outgoingMsg.ParentMessageID = message.GetParentMessageId().String()
	}

	// Handle sender user information safely
	if messageDTO.SenderUser != nil {
		outgoingMsg.SenderUser = &UserDTO{
			ID:        messageDTO.SenderUser.ID,
			Username:  messageDTO.SenderUser.Username,
			Email:     messageDTO.SenderUser.Email,
			FirstName: messageDTO.SenderUser.FirstName,
			LastName:  messageDTO.SenderUser.LastName,
		}
	}

	// Handle integration bot information safely
	if messageDTO.IntegrationBot != nil {
		outgoingMsg.IntegrationBot = &IntegrationBotDTO{
			ID:          messageDTO.IntegrationBot.ID,
			ServiceName: messageDTO.IntegrationBot.ServiceName,
			CreatedAt:   messageDTO.IntegrationBot.CreatedAt,
			IsRevoked:   messageDTO.IntegrationBot.IsRevoked,
		}
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
func (h *WebSocketHandler) handleIncomingReaction(userID string, payload interface{}) error {
	logger := h.logger.WithMethod("handleIncomingReaction")
	logger.Info("Handling incoming reaction")

	// Parse the reaction payload
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		logger.Error("Failed to marshal payload", zap.Error(err))
		return err
	}

	var incomingReaction IncomingReactionPayload
	if err := json.Unmarshal(payloadBytes, &incomingReaction); err != nil {
		logger.Error("Failed to unmarshal payload", zap.Error(err))
		return err
	}

	// Validate required fields
	if incomingReaction.MessageID == "" {
		logger.Error("Message ID is required")
		return fmt.Errorf("message_id is required")
	}
	if incomingReaction.ChannelID == "" {
		logger.Error("Channel ID is required")
		return fmt.Errorf("channel_id is required")
	}
	if incomingReaction.ReactionType == "" {
		logger.Error("Reaction type is required")
		return fmt.Errorf("reaction_type is required")
	}

	// Parse UUIDs
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		logger.Error("Invalid user ID", zap.Error(err))
		return fmt.Errorf("invalid user ID: %w", err)
	}

	channelUUID, err := uuid.Parse(incomingReaction.ChannelID)
	if err != nil {
		logger.Error("Invalid channel ID", zap.Error(err))
		return fmt.Errorf("invalid channel ID: %w", err)
	}

	messageUUID, err := uuid.Parse(incomingReaction.MessageID)
	if err != nil {
		logger.Error("Invalid message ID", zap.Error(err))
		return fmt.Errorf("invalid message ID: %w", err)
	}

	cmd := domain.AddReactionCommand{
		ChannelID:    channelUUID,
		MessageID:    messageUUID,
		UserID:       userUUID,
		ReactionType: incomingReaction.ReactionType,
	}

	// Handle through domain service
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	reaction, err := h.messageService.HandleAddReaction(ctx, cmd)
	if err != nil {
		logger.Error("Failed to add reaction", zap.Error(err))
		return fmt.Errorf("failed to add reaction: %w", err)
	}

	outgoingReaction := OutgoingReactionPayload{
		ID:           reaction.GetId().String(),
		MessageID:    reaction.GetMessageId().String(),
		ChannelID:    incomingReaction.ChannelID,
		UserID:       reaction.GetUserId().String(),
		ReactionType: reaction.GetReactionType(),
		Timestamp:    reaction.GetCreatedAt(),
	}

	if h.redisClient != nil {
		go h.publishReactionToRedis(outgoingReaction)
	} else {
		go h.broadcastToChannel(incomingReaction.ChannelID, WebSocketMessage{
			Type:    "reaction_added",
			Payload: outgoingReaction,
		})
	}

	return nil
}

func (h *WebSocketHandler) handleRemoveReaction(userID string, payload interface{}) error {
	logger := h.logger.WithMethod("handleRemoveReaction")
	logger.Info("Handling remove reaction")

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		logger.Error("Failed to marshal payload", zap.Error(err))
		return err
	}

	var incomingReaction IncomingReactionPayload
	if err := json.Unmarshal(payloadBytes, &incomingReaction); err != nil {
		logger.Error("Failed to unmarshal payload", zap.Error(err))
		return err
	}

	if incomingReaction.MessageID == "" {
		logger.Error("Message ID is required")
		return fmt.Errorf("message_id is required")
	}
	if incomingReaction.ChannelID == "" {
		logger.Error("Channel ID is required")
		return fmt.Errorf("channel_id is required")
	}
	if incomingReaction.ReactionType == "" {
		logger.Error("Reaction type is required")
		return fmt.Errorf("reaction_type is required")
	}

	userUUID, err := uuid.Parse(userID)
	if err != nil {
		logger.Error("Invalid user ID", zap.Error(err))
		return fmt.Errorf("invalid user ID: %w", err)
	}

	channelUUID, err := uuid.Parse(incomingReaction.ChannelID)
	if err != nil {
		logger.Error("Invalid channel ID", zap.Error(err))
		return fmt.Errorf("invalid channel ID: %w", err)
	}

	messageUUID, err := uuid.Parse(incomingReaction.MessageID)
	if err != nil {
		logger.Error("Invalid message ID", zap.Error(err))
		return fmt.Errorf("invalid message ID: %w", err)
	}

	cmd := domain.RemoveReactionCommand{
		ChannelID:    channelUUID,
		MessageID:    messageUUID,
		UserID:       userUUID,
		ReactionType: incomingReaction.ReactionType,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	reaction, err := h.messageService.HandleRemoveReaction(ctx, cmd)
	if err != nil {
		logger.Error("Failed to remove reaction", zap.Error(err))
		return fmt.Errorf("failed to remove reaction: %w", err)
	}

	outgoingReaction := OutgoingReactionPayload{
		ID:           reaction.GetId().String(),
		MessageID:    reaction.GetMessageId().String(),
		ChannelID:    incomingReaction.ChannelID,
		UserID:       reaction.GetUserId().String(),
		ReactionType: reaction.GetReactionType(),
		Timestamp:    reaction.GetCreatedAt(),
	}

	if h.redisClient != nil {
		go h.publishReactionRemovedToRedis(outgoingReaction)
	} else {
		go h.broadcastToChannel(incomingReaction.ChannelID, WebSocketMessage{
			Type:    "reaction_removed",
			Payload: outgoingReaction,
		})
	}

	return nil
}

func (h *WebSocketHandler) handleTypingIndicator(userID string, payload interface{}, typingType string) {
	logger := h.logger.WithMethod("handleTypingIndicator")
	logger.Info("Handling typing indicator")

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		logger.Error("Failed to marshal typing payload", zap.Error(err))
		return
	}

	var typingPayload TypingPayload
	if err := json.Unmarshal(payloadBytes, &typingPayload); err != nil {
		logger.Error("Failed to unmarshal typing payload", zap.Error(err))
		return
	}

	if typingPayload.ChannelID == "" {
		logger.Error("Channel ID is required")
		return
	}

	// Set user ID from authenticated user
	typingPayload.UserID = userID

	// Get user information from identity service
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	userInfo, err := h.identityClient.GetUserByID(ctx, userID)
	if err != nil {
		logger.Error("Failed to get user info for typing indicator", zap.Error(err))
		// Continue without user info
	} else {
		typingPayload.Username = userInfo.User.Username
		typingPayload.User = &UserDTO{
			ID:        userInfo.User.Id,
			Username:  userInfo.User.Username,
			Email:     userInfo.User.Email,
			FirstName: userInfo.User.FirstName,
			LastName:  userInfo.User.LastName,
		}
	}

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
	logger := h.logger.WithMethod("addClient")
	logger.Info("Adding client")

	h.mu.Lock()
	defer h.mu.Unlock()
	h.clients[userID] = conn
}

func (h *WebSocketHandler) removeClient(userID string) {
	logger := h.logger.WithMethod("removeClient")
	logger.Info("Removing client")

	h.mu.Lock()
	defer h.mu.Unlock()
	delete(h.clients, userID)
}

func (h *WebSocketHandler) sendToClient(userID string, message WebSocketMessage) error {
	logger := h.logger.WithMethod("sendToClient")
	logger.Info("Sending to client")

	h.mu.RLock()
	conn, exists := h.clients[userID]
	h.mu.RUnlock()

	if !exists {
		logger.Error("Client not connected")
		return nil // Client not connected
	}

	return conn.WriteJSON(message)
}

func (h *WebSocketHandler) publishMessageToRedis(message OutgoingMessagePayload) {
	logger := h.logger.WithMethod("publishMessageToRedis")
	logger.Info("Publishing message to Redis")

	if h.redisClient == nil {
		return
	}

	ctx := context.Background()

	messageJSON, err := json.Marshal(WebSocketMessage{
		Type:    "new_message",
		Payload: message,
	})

	if err != nil {
		logger.Error("Failed to marshal message", zap.Error(err))
		return
	}

	channelKey := fmt.Sprintf("channel:%s", message.ChannelID)
	err = h.redisClient.Publish(ctx, channelKey, messageJSON).Err()
	if err != nil {
		logger.Error("Failed to publish message to Redis", zap.Error(err))
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
	logger := h.logger.WithMethod("publishTypingToRedis")
	logger.Info("Publishing typing to Redis")

	if h.redisClient == nil {
		return
	}

	ctx := context.Background()

	typingJSON, err := json.Marshal(typingMsg)
	if err != nil {
		logger.Error("Failed to marshal typing message", zap.Error(err))
		return
	}

	channelKey := fmt.Sprintf("channel:%s", channelID)
	err = h.redisClient.Publish(ctx, channelKey, typingJSON).Err()
	if err != nil {
		logger.Error("Failed to publish typing message to Redis", zap.Error(err))
	}
}

func (h *WebSocketHandler) subscribeToRedisMessages() {
	logger := h.logger.WithMethod("subscribeToRedisMessages")
	logger.Info("Subscribing to Redis messages")

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
			logger.Error("Failed to unmarshal message", zap.Error(err))
			continue
		}
		if strings.HasPrefix(msg.Channel, "channel:") {
			channelID := strings.TrimPrefix(msg.Channel, "channel:")
			h.broadcastToChannel(channelID, wsMessage)
		}
	}
}

func (h *WebSocketHandler) publishReactionToRedis(reaction OutgoingReactionPayload) {
	logger := h.logger.WithMethod("publishReactionToRedis")
	logger.Info("Publishing reaction to Redis")

	if h.redisClient == nil {
		return
	}

	ctx := context.Background()

	reactionJSON, err := json.Marshal(WebSocketMessage{
		Type:    "reaction_added",
		Payload: reaction,
	})

	if err != nil {
		logger.Error("Failed to marshal reaction", zap.Error(err))
		return
	}

	channelKey := fmt.Sprintf("channel:%s", reaction.ChannelID)
	err = h.redisClient.Publish(ctx, channelKey, reactionJSON).Err()
	if err != nil {
		logger.Error("Failed to publish reaction to Redis", zap.Error(err))
	}
}

func (h *WebSocketHandler) publishReactionRemovedToRedis(reaction OutgoingReactionPayload) {
	logger := h.logger.WithMethod("publishReactionRemovedToRedis")
	logger.Info("Publishing reaction removed to Redis")

	if h.redisClient == nil {
		return
	}

	ctx := context.Background()

	reactionJSON, err := json.Marshal(WebSocketMessage{
		Type:    "reaction_removed",
		Payload: reaction,
	})

	if err != nil {
		logger.Error("Failed to marshal reaction removal", zap.Error(err))
		return
	}

	channelKey := fmt.Sprintf("channel:%s", reaction.ChannelID)
	err = h.redisClient.Publish(ctx, channelKey, reactionJSON).Err()
	if err != nil {
		logger.Error("Failed to publish reaction removal to Redis", zap.Error(err))
	}

}

func (h *WebSocketHandler) broadcastToChannel(channelID string, message WebSocketMessage) {
	logger := h.logger.WithMethod("broadcastToChannel")
	logger.Info("Broadcasting to channel")

	h.mu.RLock()
	defer h.mu.RUnlock()

	for userID, conn := range h.clients {
		err := conn.WriteJSON(message)
		//TODO: check if the current user is a member of the channel
		if err != nil {
			logger.Error("Failed to send message to user", zap.String("user_id", userID), zap.Error(err))
			conn.Close()
			delete(h.clients, userID)
		}
	}
}

func (h *WebSocketHandler) BroadcastMessage(message *domain.Message) {
	logger := h.logger.WithMethod("BroadcastMessage")
	logger.Info("Broadcasting message")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	messageDTO, err := h.messageService.ToMessageDTO(ctx, message)
	if err != nil {
		logger.Error("Failed to convert message to DTO", zap.Error(err))
		return
	}

	outgoingMsg := OutgoingMessagePayload{
		ID:        messageDTO.ID,
		Content:   messageDTO.ContentText,
		ChannelID: messageDTO.ChannelID,
		Timestamp: messageDTO.CreatedAt,
	}

	// Handle sender ID safely
	if messageDTO.SenderUserID != nil {
		outgoingMsg.SenderUserID = *messageDTO.SenderUserID
	} else if message.GetSenderUserId() != nil {
		outgoingMsg.SenderUserID = message.GetSenderUserId().String()
	}

	// Handle integration ID safely
	if messageDTO.IntegrationID != nil {
		outgoingMsg.IntegrationID = *messageDTO.IntegrationID
	} else if message.GetIntegrationId() != nil {
		outgoingMsg.IntegrationID = message.GetIntegrationId().String()
	}

	// Handle parent message ID safely
	if messageDTO.ParentMessageID != nil {
		outgoingMsg.ParentMessageID = *messageDTO.ParentMessageID
	} else if message.GetParentMessageId() != nil {
		outgoingMsg.ParentMessageID = message.GetParentMessageId().String()
	}

	// Handle sender user information safely
	if messageDTO.SenderUser != nil {
		outgoingMsg.SenderUser = &UserDTO{
			ID:        messageDTO.SenderUser.ID,
			Username:  messageDTO.SenderUser.Username,
			Email:     messageDTO.SenderUser.Email,
			FirstName: messageDTO.SenderUser.FirstName,
			LastName:  messageDTO.SenderUser.LastName,
		}
	}

	// Handle integration bot information safely
	if messageDTO.IntegrationBot != nil {
		outgoingMsg.IntegrationBot = &IntegrationBotDTO{
			ID:          messageDTO.IntegrationBot.ID,
			ServiceName: messageDTO.IntegrationBot.ServiceName,
			CreatedAt:   messageDTO.IntegrationBot.CreatedAt,
			IsRevoked:   messageDTO.IntegrationBot.IsRevoked,
		}
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
	logger := h.logger.WithMethod("SendToUser")
	logger.Info("Sending to user")

	return h.sendToClient(userID, message)
}

func (h *WebSocketHandler) validateToken(ctx context.Context, token string) (string, error) {
	logger := h.logger.WithMethod("validateToken")
	logger.Info("Validating token")

	resp, err := h.identityClient.ValidateToken(ctx, token)
	if err != nil {
		logger.Error("Failed to validate token", zap.Error(err))
		return "", fmt.Errorf("failed to validate token: %w", err)
	}

	return resp.UserId, nil
}
