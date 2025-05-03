package models

import (
	"errors"
	"time"

	"github.com/google/uuid"
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
}

func NewChannel(name string, creatorUserID uuid.UUID) (*Channel, error) {
	if name == "" {
		return nil, errors.New("channel name cannot be empty")
	}

	now := time.Now().UTC()
	channelID := uuid.New()

	channel := &Channel{
		ID:            channelID,
		Name:          name,
		CreatorUserID: creatorUserID,
		CreationTime:  now,
		Members: []Member{ // TODO: add to member
			{
				ID:       creatorUserID,
				Role:     "owner",
				JoinedAt: now,
				LastRead: now,
			},
		},
		Messages:        []Message{},
		LastMessageTime: now,
		IsArchived:      false,
		Version:         1,
	}
	return channel, nil
}

func (c *Channel) AddMember(userID uuid.UUID) error {
	for _, member := range c.Members {
		if member.ID == userID {
			return errors.New("user is already a member of the channel")
		}
	}
	now := time.Now().UTC()
	c.Members = append(c.Members, Member{
		ID:       userID,
		Role:     "member",
		JoinedAt: now,
		LastRead: now,
	})
	return nil
}

func (c *Channel) ArchiveChannel(userID uuid.UUID) error {
	if c.CreatorUserID != userID {
		return errors.New("only the channel owner can archive it")
	}

	c.IsArchived = true
	c.Version++

	return nil
}

func (c *Channel) UnarchiveChannel(userId uuid.UUID) error {
	if c.CreatorUserID != userId {
		return errors.New("only the channel owner can archive it")
	}

	c.IsArchived = false
	c.Version++
	return nil
}

func (c *Channel) SetTopic(userID uuid.UUID, topic string) error {
	if c.CreatorUserID != userID {
		return errors.New("user does not have permission to do this action")
	}

	c.Topic = topic
	c.Version++

	return nil
}

func (c *Channel) CanUserPostMessage(userID uuid.UUID) bool {
	for _, member := range c.Members {
		if member.ID == userID {
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
			if msg.ID == *parentMessageID {
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

	message := Message{
		ID:              messageID,
		ChannelID:       c.ID,
		SenderUserID:    &senderUserID,
		IntegrationID:   nil,
		Content:         content,
		Timestamp:       now,
		ParentMessageID: parentMessageID,
		Reactions:       []Reaction{},
	}

	c.Messages = append(c.Messages, message)
	c.LastMessageTime = now
	c.Version++
	return &message, nil
}

// Also a message but one sent by an integration service bot
func (c *Channel) PostNotification(integrationID uuid.UUID, content MessageContent) (*Message, error) {
	now := time.Now().UTC()
	messageID, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}

	message := Message{
		ID:              messageID,
		ChannelID:       c.ID,
		SenderUserID:    nil,
		IntegrationID:   &integrationID,
		Content:         content,
		Timestamp:       now,
		ParentMessageID: nil,
		Reactions:       []Reaction{},
	}

	c.Messages = append(c.Messages, message)
	c.LastMessageTime = now
	c.Version++

	return &message, nil
}

func (c *Channel) AddReaction(messageID, userID uuid.UUID, reactionType string) (*Reaction, error) {
	if !c.CanUserPostMessage(userID) {
		return nil, errors.New("user is not allowed to react in this channel")
	}
	var targetMessage *Message
	for i := range c.Messages {
		if c.Messages[i].ID == messageID {
			targetMessage = &c.Messages[i]
			break
		}
	}

	if targetMessage == nil {
		return nil, errors.New("message not found")
	}

	for _, reaction := range targetMessage.Reactions {
		if reaction.UserID == userID && reaction.ReactionType == reactionType {
			return nil, errors.New("user already added this reaction")
		}
	}

	now := time.Now().UTC()
	reactionID, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}
	reaction := Reaction{
		ID:           reactionID,
		UserID:       userID,
		MessageID:    messageID,
		ReactionType: reactionType,
		Timestamp:    now,
	}

	targetMessage.Reactions = append(targetMessage.Reactions, reaction)
	c.Version++

	return &reaction, nil
}

func (c *Channel) RemoveReaction(messageID, userID uuid.UUID, reactionType string) error {
	var targetMessage *Message
	for i := range c.Messages {
		if c.Messages[i].ID == messageID {
			targetMessage = &c.Messages[i]
			break
		}
	}

	if targetMessage == nil {
		return errors.New("message not found")
	}

	found := false
	for i, reaction := range targetMessage.Reactions {
		if reaction.UserID == userID && reaction.ReactionType == reactionType {
			lastIdx := len(targetMessage.Reactions) - 1
			targetMessage.Reactions[i] = targetMessage.Reactions[lastIdx]
			targetMessage.Reactions = targetMessage.Reactions[:lastIdx]
			found = true
			break
		}
	}
	if !found {
		return errors.New("reaction not found")
	}
	c.Version++
	return nil
}
