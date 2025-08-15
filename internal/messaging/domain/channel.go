package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/m1thrandir225/meridian/pkg/common"
)

type Channel struct {
	ID              uuid.UUID
	Name            string
	Topic           string
	CreationTime    time.Time
	CreatorUserID   uuid.UUID
	Members         []Member
	Messages        []Message
	LastMessageTime time.Time
	IsArchived      bool
	Version         int64
	pendingEvents   []common.DomainEvent
}

func (c *Channel) addEvent(event common.DomainEvent) {
	c.pendingEvents = append(c.pendingEvents, event)
}

func (c *Channel) GetPendingEvents() []common.DomainEvent {
	return c.pendingEvents
}

func (c *Channel) ClearPendingEvents() {
	c.pendingEvents = []common.DomainEvent{}
}

func NewChannel(name, topic string, creatorUserID uuid.UUID) (*Channel, error) {
	if name == "" {
		return nil, errors.New("channel name cannot be empty")
	}

	now := time.Now().UTC()
	channelID := uuid.New()
	creator := newMember(creatorUserID, "owner", now, now)

	channel := &Channel{
		ID:              channelID,
		Name:            name,
		Topic:           topic,
		CreatorUserID:   creatorUserID,
		CreationTime:    now,
		Members:         []Member{creator},
		Messages:        []Message{},
		LastMessageTime: now,
		IsArchived:      false,
		Version:         1,
	}

	channel.addEvent(CreateChannelCreatedEvent(channel))
	return channel, nil
}

func (c *Channel) AddMember(userID uuid.UUID) error {
	for _, member := range c.Members {
		if member.GetId() == userID {
			return errors.New("user is already a member of the channel")
		}
	}
	now := time.Now().UTC()
	member := newMember(userID, "member", now, now)
	c.Members = append(c.Members, member)

	c.addEvent(CreateUserJoinedChannelEvent(c, member))
	c.Version++
	return nil
}

func (c *Channel) RemoveMember(userID uuid.UUID) error {
	found := false
	var searchMember Member
	for _, member := range c.Members {
		if member.GetId() == userID {
			lastIdx := len(c.Members) - 1
			searchMember = member
			c.Members = c.Members[:lastIdx]
			found = true
			break
		}
	}
	if !found {
		return errors.New("user not apart of the channel")
	}

	c.addEvent(CreateUserLeftChannelEvent(c, searchMember.id))

	return nil
}

func (c *Channel) ArchiveChannel(userID uuid.UUID) error {
	if c.CreatorUserID != userID {
		return errors.New("only the channel owner can archive it")
	}

	c.IsArchived = true
	c.Version++

	c.addEvent(CreateChannelArchivedEvent(c, userID))
	return nil
}

func (c *Channel) UnarchiveChannel(userId uuid.UUID) error {
	if c.CreatorUserID != userId {
		return errors.New("only the channel owner can archive it")
	}

	c.IsArchived = false
	c.Version++

	c.addEvent(CreateChannelUnarchivedEvent(c, userId))
	return nil
}

func (c *Channel) SetTopic(userID uuid.UUID, topic string) error {
	if c.CreatorUserID != userID {
		return errors.New("user does not have permission to do this action")
	}

	c.Topic = topic
	c.Version++

	c.addEvent(CreateChannelTopicChangedEvent(c, userID))
	return nil
}

func (c *Channel) CanUserPostMessage(userID uuid.UUID) bool {
	for _, member := range c.Members {
		if member.GetId() == userID {
			return true
		}
	}
	return false
}

// Normal chat message sent by a user
func (c *Channel) PostMessage(senderUserID uuid.UUID, content MessageContent, parentMessageID *uuid.UUID) (*Message, error) {
	if !c.CanUserPostMessage(senderUserID) {
		return nil, errors.New("user is not allowed to post in this channel")
	}

	if parentMessageID != nil {
		parentFound := false
		for _, msg := range c.Messages {
			if msg.GetId() == *parentMessageID {
				parentFound = true
				break
			}
		}
		if !parentFound {
			return nil, errors.New("parent message not found")
		}
	}

	now := time.Now().UTC()
	messageID, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}

	message := newMessage(
		messageID,
		c.ID,
		&senderUserID,
		nil,
		parentMessageID,
		content,
		[]Reaction{},
		now,
	)

	c.Messages = append(c.Messages, message)
	c.LastMessageTime = now
	c.Version++
	c.addEvent(CreateMessageSentEvent(c, &message))
	return &message, nil
}

// Also a message but one sent by an integration service bot
func (c *Channel) PostNotification(integrationID uuid.UUID, content MessageContent) (*Message, error) {
	now := time.Now().UTC()
	messageID, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}

	message := newMessage(
		messageID,
		c.ID,
		nil,
		&integrationID,
		nil,
		content,
		[]Reaction{},
		now,
	)

	c.Messages = append(c.Messages, message)
	c.LastMessageTime = now
	c.Version++
	c.addEvent(CreateMessageSentEvent(c, &message))
	return &message, nil
}

func (c *Channel) AddReaction(messageID, userID uuid.UUID, reactionType string) (*Reaction, error) {
	if !c.CanUserPostMessage(userID) {
		return nil, errors.New("user is not allowed to react in this channel")
	}
	var targetMessage *Message
	for i := range c.Messages {
		if c.Messages[i].GetId() == messageID {
			targetMessage = &c.Messages[i]
			break
		}
	}

	if targetMessage == nil {
		return nil, errors.New("message not found")
	}

	for _, reaction := range targetMessage.GetReactions() {
		if reaction.GetUserId() == userID && reaction.GetReactionType() == reactionType {
			return nil, errors.New("user already added this reaction")
		}
	}

	now := time.Now().UTC()
	reactionID, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}
	reaction := newReaction(
		reactionID,
		messageID,
		userID,
		reactionType,
		now,
	)

	targetMessage.setReactions(append(targetMessage.GetReactions(), reaction))
	c.Version++

	c.addEvent(CreateReactionAddedEvent(c, &reaction))

	return &reaction, nil
}

func (c *Channel) RemoveReaction(messageID, userID uuid.UUID, reactionType string) (*Reaction, error) {
	var targetMessage *Message
	for i := range c.Messages {
		if c.Messages[i].GetId() == messageID {
			targetMessage = &c.Messages[i]
			break
		}
	}

	if targetMessage == nil {
		return nil, errors.New("message not found")
	}

	found := false
	var reaction Reaction
	for i, reaction := range targetMessage.GetReactions() {
		if reaction.GetUserId() == userID && reaction.GetReactionType() == reactionType {
			lastIdx := len(targetMessage.GetReactions()) - 1
			targetMessage.GetReactions()[i] = targetMessage.GetReactions()[lastIdx]
			reaction = targetMessage.GetReactions()[lastIdx]
			targetMessage.setReactions(targetMessage.GetReactions()[:lastIdx])

			found = true

			break
		}
	}
	if !found {
		return nil, errors.New("reaction not found")
	}
	c.Version++
	c.addEvent(CreateReactionRemovedEvent(c, messageID, userID, reactionType))
	return &reaction, nil
}

func (c *Channel) AddBotMember(integrationID uuid.UUID) error {
	for _, member := range c.Members {
		if member.GetId() == integrationID {
			return errors.New("bot is already a member of the channel")
		}
	}

	now := time.Now().UTC()
	member := newMember(integrationID, "bot", now, now)
	c.Members = append(c.Members, member)

	c.addEvent(CreateBotJoinedChannelEvent(c, member))
	c.Version++
	return nil
}
