package domain

import (
	"time"

	"github.com/google/uuid"
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
	Content         MessageContent
	Timestamp       time.Time
	ParentMessageID *string
}

type NotificationSentEvent struct {
	BaseDomainEvent
	MessageID     string
	IntegrationID string
	Content       MessageContent
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

func CreateChannelCreatedEvent(channel *Channel) ChannelCreatedEvent {
	base := NewBaseDomainEvent("ChannelCreated", channel.ID, channel.Version)
	return ChannelCreatedEvent{
		BaseDomainEvent: base,
		Name:            channel.Name,
		CreatorUserID:   channel.CreatorUserID.String(),
		Topic:           channel.Topic,
	}
}

func CreateUserJoinedChannelEvent(channel *Channel, member Member) UserJoinedChannelEvent {
	base := NewBaseDomainEvent("UserJoinedChannel", channel.ID, channel.Version)
	return UserJoinedChannelEvent{
		BaseDomainEvent: base,
		UserID:          member.GetId().String(),
		Role:            member.GetRole(),
		JoinedAt:        member.GetJoinedAt(),
	}
}

func CreateUserLeftChannelEvent(channel *Channel, userID uuid.UUID) UserLeftChannelEvent {
	base := NewBaseDomainEvent("UserLeftChannel", channel.ID, channel.Version)
	return UserLeftChannelEvent{
		BaseDomainEvent: base,
		UserID:          userID.String(),
	}
}

func CreateMessageSentEvent(channel *Channel, message *Message) MessageSentEvent {
	base := NewBaseDomainEvent("MessageSent", channel.ID, channel.Version)

	var senderUserIDStr *string
	if message.GetSenderUserId() != nil {
		id := message.GetSenderUserId().String()
		senderUserIDStr = &id
	}

	var integrationIDStr *string
	if message.GetIntegrationId() != nil {
		id := message.GetIntegrationId().String()
		integrationIDStr = &id
	}

	var parentMessageIDStr *string
	if message.GetParentMessageId() != nil {
		id := message.GetParentMessageId().String()
		parentMessageIDStr = &id
	}

	return MessageSentEvent{
		BaseDomainEvent: base,
		MessageID:       message.GetId().String(),
		SenderUserID:    senderUserIDStr,
		IntegrationID:   integrationIDStr,
		Content:         message.GetContent(),
		Timestamp:       message.GetTimestamp(),
		ParentMessageID: parentMessageIDStr,
	}
}

func CreateNotificationSentEvent(channel *Channel, message *Message) NotificationSentEvent {
	base := NewBaseDomainEvent("NotificationSent", channel.ID, channel.Version)

	return NotificationSentEvent{
		BaseDomainEvent: base,
		MessageID:       message.GetId().String(),
		IntegrationID:   message.GetIntegrationId().String(),
		Content:         message.GetContent(),
		Timestamp:       message.GetTimestamp(),
	}
}

func CreateReactionAddedEvent(channel *Channel, reaction *Reaction) ReactionAddedEvent {
	base := NewBaseDomainEvent("ReactionAdded", channel.ID, channel.Version)

	return ReactionAddedEvent{
		BaseDomainEvent: base,
		ReactionID:      reaction.GetId().String(),
		MessageID:       reaction.GetMessageId().String(),
		UserID:          reaction.GetUserId().String(),
		ReactionType:    reaction.GetReactionType(),
		Timestamp:       reaction.GetCreatedAt(),
	}
}

func CreateReactionRemovedEvent(channel *Channel, messageID uuid.UUID, userID uuid.UUID, reactionType string) ReactionRemovedEvent {
	base := NewBaseDomainEvent("ReactionRemoved", channel.ID, channel.Version)

	return ReactionRemovedEvent{
		BaseDomainEvent: base,
		MessageID:       messageID.String(),
		UserID:          userID.String(),
		ReactionType:    reactionType,
	}
}

func CreateChannelTopicChangedEvent(channel *Channel, changedBy uuid.UUID) ChannelTopicChangedEvent {
	base := NewBaseDomainEvent("ChannelTopicChanged", channel.ID, channel.Version)

	return ChannelTopicChangedEvent{
		BaseDomainEvent: base,
		Topic:           channel.Topic,
		ChangedBy:       changedBy.String(),
	}
}

func CreateChannelArchivedEvent(channel *Channel, archivedBy uuid.UUID) ChannelArchivedEvent {
	base := NewBaseDomainEvent("ChannelArchived", channel.ID, channel.Version)

	return ChannelArchivedEvent{
		BaseDomainEvent: base,
		ArchivedBy:      archivedBy.String(),
	}
}

func CreateChannelUnarchivedEvent(channel *Channel, unarchivedBy uuid.UUID) ChannelUnarchivedEvent {
	base := NewBaseDomainEvent("ChannelUnarchived", channel.ID, channel.Version)

	return ChannelUnarchivedEvent{
		BaseDomainEvent: base,
		UnarchivedBy:    unarchivedBy.String(),
	}
}
