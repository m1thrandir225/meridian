package handlers

import (
	"context"
	"encoding/json"
	"time"

	"github.com/m1thrandir225/meridian/internal/analytics/application/services"
	"github.com/m1thrandir225/meridian/internal/analytics/domain"
	"github.com/m1thrandir225/meridian/pkg/kafka"
	"github.com/m1thrandir225/meridian/pkg/logging"
	"go.uber.org/zap"
)

type AnalyticsEventHandler struct {
	analyticsService *services.AnalyticsService
	logger           *logging.Logger
}

func NewAnalyticsEventHandler(analyticsService *services.AnalyticsService, logger *logging.Logger) *AnalyticsEventHandler {
	return &AnalyticsEventHandler{
		analyticsService: analyticsService,
		logger:           logger,
	}
}

func (h *AnalyticsEventHandler) HandleEvent(ctx context.Context, event kafka.Event) error {
	logger := h.logger.WithMethod("HandleEvent")
	logger.Info("Handling analytics event",
		zap.String("topic", event.Topic),
		zap.String("key", event.Key),
		zap.Int32("partition", event.Partition),
		zap.Int64("offset", event.Offset))

	// Debug: Log the raw JSON structure to see what we're actually receiving
	var rawEvent map[string]interface{}
	if err := json.Unmarshal(event.Data, &rawEvent); err != nil {
		logger.Error("Failed to parse raw event", zap.Error(err))
		return err
	}
	logger.Info("Raw event structure", zap.Any("event", rawEvent))

	// Parse the event type from the BaseDomainEvent structure
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

func (h *AnalyticsEventHandler) handleUserProfileUpdated(ctx context.Context, event kafka.Event) error {
	logger := h.logger.WithMethod("handleUserProfileUpdated")
	logger.Info("Processing user profile updated event")

	var userEvent UserProfileUpdatedEvent
	if err := json.Unmarshal(event.Data, &userEvent); err != nil {
		logger.Error("Failed to unmarshal user profile updated event", zap.Error(err))
		return err
	}

	// Track profile update as a metric
	cmd := domain.TrackMessageSentCommand{
		MessageID:     "profile_update_" + userEvent.UserID,
		ChannelID:     "system",
		SenderID:      userEvent.UserID,
		Timestamp:     userEvent.Timestamp,
		ContentLength: len(userEvent.UpdatedFields),
	}

	return h.analyticsService.TrackMessageSent(ctx, cmd)
}

func (h *AnalyticsEventHandler) handleMessageSent(ctx context.Context, eventData []byte) error {
	logger := h.logger.WithMethod("handleMessageSent")

	var messageEvent MessageSentEvent
	if err := json.Unmarshal(eventData, &messageEvent); err != nil {
		logger.Error("Failed to unmarshal message event", zap.Error(err))
		// Log the raw data to see what we're trying to parse
		logger.Error("Raw message event data", zap.String("data", string(eventData)))
		return err
	}

	// Log the parsed event to debug
	logger.Info("Parsed message event",
		zap.String("message_id", messageEvent.MessageID),
		zap.String("aggr_id", messageEvent.AggrID),
		zap.String("sender_user_id", func() string {
			if messageEvent.SenderUserID != nil {
				return *messageEvent.SenderUserID
			}
			return "nil"
		}()),
		zap.String("integration_id", func() string {
			if messageEvent.IntegrationID != nil {
				return *messageEvent.IntegrationID
			}
			return "nil"
		}()))

	// Check if we have a sender (either user or integration)
	if messageEvent.SenderUserID == nil && messageEvent.IntegrationID == nil {
		logger.Warn("No sender ID found in message event",
			zap.String("message_id", messageEvent.MessageID))
		return nil // Skip this message, don't treat as error
	}

	// Parse timestamp
	timestamp, err := time.Parse(time.RFC3339, messageEvent.Timestamp)
	if err != nil {
		logger.Error("Failed to parse timestamp", zap.Error(err))
		return err
	}

	// Use the sender ID (either user or integration)
	senderID := ""
	if messageEvent.SenderUserID != nil {
		senderID = *messageEvent.SenderUserID
	} else if messageEvent.IntegrationID != nil {
		senderID = *messageEvent.IntegrationID
	}

	// Create command for tracking message sent
	cmd := domain.TrackMessageSentCommand{
		MessageID:     messageEvent.MessageID,
		ChannelID:     messageEvent.AggrID, // Channel ID is the aggregate ID
		SenderID:      senderID,
		Timestamp:     timestamp,
		ContentLength: len(messageEvent.Content.Text), // Use the actual text content
	}

	// Track the message using the service method
	if err := h.analyticsService.TrackMessageSent(ctx, cmd); err != nil {
		logger.Error("Failed to track message sent", zap.Error(err))
		return err
	}

	logger.Info("Tracked message sent",
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

	// Track user leaving channel (could be used for engagement metrics)
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

	// Track reaction removal (could be used for engagement metrics)
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

	// Track integration registration
	cmd := domain.TrackMessageSentCommand{
		MessageID:     "integration_registered_" + integrationEvent.IntegrationID,
		ChannelID:     "system",
		SenderID:      integrationEvent.CreatorUserID,
		Timestamp:     integrationEvent.RegisteredAt,
		ContentLength: len(integrationEvent.TargetChannelIDs),
	}

	return h.analyticsService.TrackMessageSent(ctx, cmd)
}
