package handlers

import (
	"context"
	"encoding/json"
	"time"

	"sync"

	"github.com/m1thrandir225/meridian/internal/analytics/application/services"
	"github.com/m1thrandir225/meridian/internal/analytics/domain"
	"github.com/m1thrandir225/meridian/pkg/kafka"
	"github.com/m1thrandir225/meridian/pkg/logging"
	"go.uber.org/zap"
)

type AnalyticsEventHandler struct {
	analyticsService *services.AnalyticsService
	logger           *logging.Logger
	processedMsgs    map[string]time.Time
	mu               sync.RWMutex
}

func NewAnalyticsEventHandler(analyticsService *services.AnalyticsService, logger *logging.Logger) *AnalyticsEventHandler {
	return &AnalyticsEventHandler{
		analyticsService: analyticsService,
		logger:           logger,
		processedMsgs:    make(map[string]time.Time),
		mu:               sync.RWMutex{},
	}
}

func (h *AnalyticsEventHandler) HandleEvent(ctx context.Context, event kafka.Event) error {
	logger := h.logger.WithMethod("HandleEvent")
	logger.Info("Handling analytics event",
		zap.String("topic", event.Topic),
		zap.String("key", event.Key),
		zap.Int32("partition", event.Partition),
		zap.Int64("offset", event.Offset))

	var rawEvent map[string]interface{}
	if err := json.Unmarshal(event.Data, &rawEvent); err != nil {
		logger.Error("Failed to parse raw event", zap.Error(err))
		return err
	}
	logger.Info("Raw event structure", zap.Any("event", rawEvent))

	var baseEvent struct {
		Name string `json:"name"`
	}

	if err := json.Unmarshal(event.Data, &baseEvent); err != nil {
		logger.Error("Failed to parse event name", zap.Error(err))
		return err
	}

	logger.Info("Parsed event name", zap.String("event_name", baseEvent.Name))

	switch baseEvent.Name {
	case "UserRegistered":
		return h.handleUserRegistered(ctx, event)
	case "UserProfileUpdated":
		return h.handleUserProfileUpdated(ctx, event)
	case "MessageSent":
		return h.handleMessageSent(ctx, event.Data)
	case "ChannelCreated":
		return h.handleChannelCreated(ctx, event)
	case "UserJoinedChannel":
		return h.handleUserJoinedChannel(ctx, event)
	case "UserLeftChannel":
		return h.handleUserLeftChannel(ctx, event)
	case "ReactionAdded":
		return h.handleReactionAdded(ctx, event)
	case "ReactionRemoved":
		return h.handleReactionRemoved(ctx, event)
	case "IntegrationRegistered":
		return h.handleIntegrationRegistered(ctx, event)
	default:
		logger.Info("Unhandled event type", zap.String("event_name", baseEvent.Name))
		return nil
	}
}

// handleUserRegistered is a handler for the UserRegistered event
func (h *AnalyticsEventHandler) handleUserRegistered(ctx context.Context, event kafka.Event) error {
	logger := h.logger.WithMethod("handleUserRegistered")
	logger.Info("Processing user registration event")

	var userEvent UserRegisteredEvent
	if err := json.Unmarshal(event.Data, &userEvent); err != nil {
		logger.Error("Failed to unmarshal user registration event", zap.Error(err))
		return err
	}

	cmd := domain.TrackUserRegistrationCommand{
		UserID:    userEvent.UserID,
		Timestamp: userEvent.Timestamp,
	}

	return h.analyticsService.TrackUserRegistration(ctx, cmd)
}

// handleUserProfileUpdated is a handler for the UserProfileUpdated event
func (h *AnalyticsEventHandler) handleUserProfileUpdated(ctx context.Context, event kafka.Event) error {
	logger := h.logger.WithMethod("handleUserProfileUpdated")
	logger.Info("Processing user profile updated event")

	var userEvent UserProfileUpdatedEvent
	if err := json.Unmarshal(event.Data, &userEvent); err != nil {
		logger.Error("Failed to unmarshal user profile updated event", zap.Error(err))
		return err
	}

	cmd := domain.TrackMessageSentCommand{
		MessageID:     "profile_update_" + userEvent.UserID,
		ChannelID:     "system",
		SenderID:      userEvent.UserID,
		Timestamp:     userEvent.Timestamp,
		ContentLength: len(userEvent.UpdatedFields),
	}

	return h.analyticsService.TrackMessageSent(ctx, cmd)
}

// handleMessageSent is a handler for the MessageSent event
func (h *AnalyticsEventHandler) handleMessageSent(ctx context.Context, eventData []byte) error {
	logger := h.logger.WithMethod("handleMessageSent")

	var messageEvent MessageSentEvent
	if err := json.Unmarshal(eventData, &messageEvent); err != nil {
		logger.Error("Failed to unmarshal message event", zap.Error(err))
		logger.Error("Raw message event data", zap.String("data", string(eventData)))
		return err
	}

	h.mu.Lock()
	if lastProcessed, exists := h.processedMsgs[messageEvent.MessageID]; exists {
		if time.Since(lastProcessed) < time.Minute {
			logger.Info("Message already processed recently, skipping",
				zap.String("message_id", messageEvent.MessageID))
			h.mu.Unlock()
			return nil
		}
	}
	h.processedMsgs[messageEvent.MessageID] = time.Now()
	h.mu.Unlock()

	logger.Info("Processing message event",
		zap.String("message_id", messageEvent.MessageID),
		zap.String("aggr_id", messageEvent.AggrID),
		zap.String("timestamp", messageEvent.Timestamp))

	if messageEvent.SenderUserID == nil && messageEvent.IntegrationID == nil {
		logger.Warn("No sender ID found in message event",
			zap.String("message_id", messageEvent.MessageID))
		return nil
	}

	timestamp, err := time.Parse(time.RFC3339, messageEvent.Timestamp)
	if err != nil {
		logger.Error("Failed to parse timestamp", zap.Error(err))
		return err
	}

	senderID := ""
	if messageEvent.SenderUserID != nil {
		senderID = *messageEvent.SenderUserID
	} else if messageEvent.IntegrationID != nil {
		senderID = *messageEvent.IntegrationID
	}

	cmd := domain.TrackMessageSentCommand{
		MessageID:     messageEvent.MessageID,
		ChannelID:     messageEvent.AggrID, // Channel ID is the aggregate ID
		SenderID:      senderID,
		Timestamp:     timestamp,
		ContentLength: len(messageEvent.Content.Text),
	}

	if err := h.analyticsService.TrackMessageSent(ctx, cmd); err != nil {
		logger.Error("Failed to track message sent", zap.Error(err))
		return err
	}

	logger.Info("Successfully tracked message sent",
		zap.String("message_id", messageEvent.MessageID),
		zap.String("channel_id", messageEvent.AggrID))

	return nil
}

func (h *AnalyticsEventHandler) handleChannelCreated(ctx context.Context, event kafka.Event) error {
	logger := h.logger.WithMethod("handleChannelCreated")
	logger.Info("Processing channel created event")

	var channelEvent ChannelCreatedEvent
	if err := json.Unmarshal(event.Data, &channelEvent); err != nil {
		logger.Error("Failed to unmarshal channel created event", zap.Error(err))
		return err
	}

	cmd := domain.TrackChannelCreatedCommand{
		ChannelID: channelEvent.AggrID, // Channel ID is the aggregate ID
		CreatorID: channelEvent.CreatorUserID,
		Timestamp: channelEvent.Time,
	}

	return h.analyticsService.TrackChannelCreated(ctx, cmd)
}

func (h *AnalyticsEventHandler) handleUserJoinedChannel(ctx context.Context, event kafka.Event) error {
	logger := h.logger.WithMethod("handleUserJoinedChannel")
	logger.Info("Processing user joined channel event")

	var joinEvent UserJoinedChannelEvent
	if err := json.Unmarshal(event.Data, &joinEvent); err != nil {
		logger.Error("Failed to unmarshal user joined channel event", zap.Error(err))
		return err
	}

	cmd := domain.TrackUserJoinedChannelCommand{
		UserID:    joinEvent.UserID,
		ChannelID: joinEvent.AggrID, // Channel ID is the aggregate ID
		Timestamp: joinEvent.JoinedAt,
	}

	return h.analyticsService.TrackUserJoinedChannel(ctx, cmd)
}

func (h *AnalyticsEventHandler) handleUserLeftChannel(ctx context.Context, event kafka.Event) error {
	logger := h.logger.WithMethod("handleUserLeftChannel")
	logger.Info("Processing user left channel event")

	var leftEvent UserLeftChannelEvent
	if err := json.Unmarshal(event.Data, &leftEvent); err != nil {
		logger.Error("Failed to unmarshal user left channel event", zap.Error(err))
		return err
	}

	cmd := domain.TrackMessageSentCommand{
		MessageID:     "user_left_" + leftEvent.UserID + "_" + leftEvent.AggrID,
		ChannelID:     leftEvent.AggrID,
		SenderID:      leftEvent.UserID,
		Timestamp:     leftEvent.Time,
		ContentLength: 0,
	}

	return h.analyticsService.TrackMessageSent(ctx, cmd)
}

func (h *AnalyticsEventHandler) handleReactionAdded(ctx context.Context, event kafka.Event) error {
	logger := h.logger.WithMethod("handleReactionAdded")
	logger.Info("Processing reaction added event")

	var reactionEvent ReactionAddedEvent
	if err := json.Unmarshal(event.Data, &reactionEvent); err != nil {
		logger.Error("Failed to unmarshal reaction added event", zap.Error(err))
		return err
	}

	cmd := domain.TrackReactionAddedCommand{
		ReactionID:   reactionEvent.ReactionID,
		MessageID:    reactionEvent.MessageID,
		UserID:       reactionEvent.UserID,
		ReactionType: reactionEvent.ReactionType,
		Timestamp:    reactionEvent.Timestamp,
	}

	return h.analyticsService.TrackReactionAdded(ctx, cmd)
}

func (h *AnalyticsEventHandler) handleReactionRemoved(ctx context.Context, event kafka.Event) error {
	logger := h.logger.WithMethod("handleReactionRemoved")
	logger.Info("Processing reaction removed event")

	var reactionEvent ReactionRemovedEvent
	if err := json.Unmarshal(event.Data, &reactionEvent); err != nil {
		logger.Error("Failed to unmarshal reaction removed event", zap.Error(err))
		return err
	}

	cmd := domain.TrackReactionAddedCommand{
		ReactionID:   "removed_" + reactionEvent.MessageID + "_" + reactionEvent.UserID,
		MessageID:    reactionEvent.MessageID,
		UserID:       reactionEvent.UserID,
		ReactionType: reactionEvent.ReactionType,
		Timestamp:    reactionEvent.Time,
	}

	return h.analyticsService.TrackReactionAdded(ctx, cmd)
}

func (h *AnalyticsEventHandler) handleIntegrationRegistered(ctx context.Context, event kafka.Event) error {
	logger := h.logger.WithMethod("handleIntegrationRegistered")
	logger.Info("Processing integration registered event")

	var integrationEvent IntegrationRegisteredEvent
	if err := json.Unmarshal(event.Data, &integrationEvent); err != nil {
		logger.Error("Failed to unmarshal integration registered event", zap.Error(err))
		return err
	}

	cmd := domain.TrackMessageSentCommand{
		MessageID:     "integration_registered_" + integrationEvent.IntegrationID,
		ChannelID:     "system",
		SenderID:      integrationEvent.CreatorUserID,
		Timestamp:     integrationEvent.RegisteredAt,
		ContentLength: len(integrationEvent.TargetChannelIDs),
	}

	return h.analyticsService.TrackMessageSent(ctx, cmd)
}
