package events

import (
	"time"

	"github.com/google/uuid"
	"github.com/m1thrandir225/meridian/internal/messaging/domain/models"
)

type ChannelCreatedEvent struct {
	BaseDomainEvent
	Name          string
	CreatorUserID string
	Topic         string
}

type UserJoinedChannelEvent struct {
	BaseDomainEvent
	UserID   string
	Role     string
	JoinedAt time.Time
}

type UserLeftChannelEvent struct {
	BaseDomainEvent
	UserID string
}

type MessageSentEvent struct {
	BaseDomainEvent
	MessageID       string
	SenderUserID    *string
	IntegrationID   *string
	Content         models.MessageContent
	Timestamp       time.Time
	ParentMessageID *string
}

type NotificationSentEvent struct {
	BaseDomainEvent
	MessageID     string
	IntegrationID string
	Content       models.MessageContent
	Timestamp     time.Time
}

type ReactionAddedEvent struct {
	BaseDomainEvent
	ReactionID   string
	MessageID    string
	UserID       string
	ReactionType string
	Timestamp    time.Time
}

type ReactionRemovedEvent struct {
	BaseDomainEvent
	MessageID    string
	UserID       string
	ReactionType string
}

type ChannelArchivedEvent struct {
	BaseDomainEvent
	ArchivedBy string
}

type ChannelUnarchivedEvent struct {
	BaseDomainEvent
	UnarchivedBy string
}

type ChannelTopicChangedEvent struct {
	BaseDomainEvent
	Topic     string
	ChangedBy string
}

func CreateChannelCreatedEvent(channel *models.Channel) ChannelCreatedEvent {
	base := NewBaseDomainEvent("ChannelCreated", channel.ID, channel.Version)
	return ChannelCreatedEvent{
		BaseDomainEvent: base,
		Name:            channel.Name,
		CreatorUserID:   channel.CreatorUserID.String(),
		Topic:           channel.Topic,
	}
}

func CreateUserJoinedChannelEvent(channel *models.Channel, member models.Member) UserJoinedChannelEvent {
	base := NewBaseDomainEvent("UserJoinedChannel", channel.ID, channel.Version)
	return UserJoinedChannelEvent{
		BaseDomainEvent: base,
		UserID:          member.ID.String(),
		Role:            member.Role,
		JoinedAt:        member.JoinedAt,
	}
}

func CreateUserLeftChannelEvent(channel *models.Channel, userID uuid.UUID) UserLeftChannelEvent {
	base := NewBaseDomainEvent("UserLeftChannel", channel.ID, channel.Version)
	return UserLeftChannelEvent{
		BaseDomainEvent: base,
		UserID:          userID.String(),
	}
}

func CreateMessageSentEvent(channel *models.Channel, message *models.Message) MessageSentEvent {
	base := NewBaseDomainEvent("MessageSent", channel.ID, channel.Version)

	var senderUserIDStr *string
	if message.SenderUserID != nil {
		id := message.SenderUserID.String()
		senderUserIDStr = &id
	}

	var integrationIDStr *string
	if message.IntegrationID != nil {
		id := message.IntegrationID.String()
		integrationIDStr = &id
	}

	var parentMessageIDStr *string
	if message.ParentMessageID != nil {
		id := message.ParentMessageID.String()
		parentMessageIDStr = &id
	}

	return MessageSentEvent{
		BaseDomainEvent: base,
		MessageID:       message.ID.String(),
		SenderUserID:    senderUserIDStr,
		IntegrationID:   integrationIDStr,
		Content:         message.Content,
		Timestamp:       message.Timestamp,
		ParentMessageID: parentMessageIDStr,
	}
}

func CreateNotificationSentEvent(channel *models.Channel, message *models.Message) NotificationSentEvent {
	base := NewBaseDomainEvent("NotificationSent", channel.ID, channel.Version)

	return NotificationSentEvent{
		BaseDomainEvent: base,
		MessageID:       message.ID.String(),
		IntegrationID:   message.IntegrationID.String(),
		Content:         message.Content,
		Timestamp:       message.Timestamp,
	}
}

func CreateReactionAddedEvent(channel *models.Channel, reaction *models.Reaction) ReactionAddedEvent {
	base := NewBaseDomainEvent("ReactionAdded", channel.ID, channel.Version)

	return ReactionAddedEvent{
		BaseDomainEvent: base,
		ReactionID:      reaction.ID.String(),
		MessageID:       reaction.MessageID.String(),
		UserID:          reaction.UserID.String(),
		ReactionType:    reaction.ReactionType,
		Timestamp:       reaction.Timestamp,
	}
}

func CreateReactionRemovedEvent(channel *models.Channel, messageID uuid.UUID, userID uuid.UUID, reactionType string) ReactionRemovedEvent {
	base := NewBaseDomainEvent("ReactionRemoved", channel.ID, channel.Version)

	return ReactionRemovedEvent{
		BaseDomainEvent: base,
		MessageID:       messageID.String(),
		UserID:          userID.String(),
		ReactionType:    reactionType,
	}
}

func CreateChannelTopicChangedEvent(channel *models.Channel, changedBy uuid.UUID) ChannelTopicChangedEvent {
	base := NewBaseDomainEvent("ChannelTopicChanged", channel.ID, channel.Version)

	return ChannelTopicChangedEvent{
		BaseDomainEvent: base,
		Topic:           channel.Topic,
		ChangedBy:       changedBy.String(),
	}
}

func CreateChannelArchivedEvent(channel *models.Channel, archivedBy uuid.UUID) ChannelArchivedEvent {
	base := NewBaseDomainEvent("ChannelArchived", channel.ID, channel.Version)

	return ChannelArchivedEvent{
		BaseDomainEvent: base,
		ArchivedBy:      archivedBy.String(),
	}
}

func CreateChannelUnarchivedEvent(channel *models.Channel, unarchivedBy uuid.UUID) ChannelUnarchivedEvent {
	base := NewBaseDomainEvent("ChannelUnarchived", channel.ID, channel.Version)

	return ChannelUnarchivedEvent{
		BaseDomainEvent: base,
		UnarchivedBy:    unarchivedBy.String(),
	}
}
