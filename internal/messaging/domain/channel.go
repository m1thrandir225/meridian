package domain

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/m1thrandir225/meridian/pkg/common"
)

// Channel represents a chat channel in the system
// It is the aggregate root for the channel domain
// It contains all the information about a channel, including its members, messages, and invites
type Channel struct {
	ID              uuid.UUID
	Name            string
	Topic           string
	CreationTime    time.Time
	CreatorUserID   uuid.UUID
	Members         []Member
	Messages        []Message
	Invites         []ChannelInvite
	LastMessageTime time.Time
	IsArchived      bool
	Version         int64
	pendingEvents   []common.DomainEvent
}

// NewChannel creates a new channel
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
		Invites:         []ChannelInvite{},
		LastMessageTime: now,
		IsArchived:      false,
		Version:         1,
	}

	channel.addEvent(CreateChannelCreatedEvent(channel))
	return channel, nil
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

// AddMember adds a member to a channel
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

// RemoveMember removes a member from a channel
func (c *Channel) RemoveMember(memberID uuid.UUID) error {
	found := false
	var searchMember Member
	for i, member := range c.Members {
		if member.GetId() == memberID {
			lastIdx := len(c.Members) - 1
			searchMember = member
			c.Members[i] = c.Members[lastIdx]
			c.Members = c.Members[:lastIdx]
			found = true
			break
		}
	}
	if !found {
		return errors.New("member not apart of the channel")
	}
	c.Version++
	c.addEvent(CreateUserLeftChannelEvent(c, searchMember.GetId()))

	return nil
}

// ArchiveChannel archives a channel
func (c *Channel) ArchiveChannel(userID uuid.UUID) error {
	if c.CreatorUserID != userID {
		return errors.New("only the channel owner can archive it")
	}

	c.IsArchived = true
	c.Version++

	c.addEvent(CreateChannelArchivedEvent(c, userID))
	return nil
}

// UnarchiveChannel unarchives a channel
func (c *Channel) UnarchiveChannel(userId uuid.UUID) error {
	if c.CreatorUserID != userId {
		return errors.New("only the channel owner can archive it")
	}

	c.IsArchived = false
	c.Version++

	c.addEvent(CreateChannelUnarchivedEvent(c, userId))
	return nil
}

// SetTopic sets the topic of a channel
func (c *Channel) SetTopic(userID uuid.UUID, topic string) error {
	if c.CreatorUserID != userID {
		return errors.New("user does not have permission to do this action")
	}

	c.Topic = topic
	c.Version++

	c.addEvent(CreateChannelTopicChangedEvent(c, userID))
	return nil
}

// canUserPostMessage checks if a user is allowed to post a message
func (c *Channel) canUserPostMessage(userID uuid.UUID) bool {
	for _, member := range c.Members {
		if member.GetId() == userID {
			return true
		}
	}
	return false
}

// PostMessage posts a message to a channel
func (c *Channel) PostMessage(senderUserID uuid.UUID, content MessageContent, parentMessageID *uuid.UUID) (*Message, error) {
	if !c.canUserPostMessage(senderUserID) {
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

// PostNotification posts a notification to a channel
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

// AddReaction adds a reaction to a message
func (c *Channel) AddReaction(messageID, userID uuid.UUID, reactionType string) (*Reaction, error) {
	if !c.canUserPostMessage(userID) {
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

// RemoveReaction removes a reaction from a message
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
	var removedReaction Reaction
	for i, r := range targetMessage.GetReactions() {
		if r.GetUserId() == userID && r.GetReactionType() == reactionType {
			lastIdx := len(targetMessage.GetReactions()) - 1
			targetMessage.GetReactions()[i] = targetMessage.GetReactions()[lastIdx]
			removedReaction = targetMessage.GetReactions()[lastIdx]
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
	return &removedReaction, nil
}

// AddBotMember adds a bot to a channel
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

// CreateInvite creates a new invite for a channelj
func (c *Channel) CreateInvite(createdByUserID uuid.UUID, expiresAt time.Time, maxUses *int) (*ChannelInvite, error) {
	isMember := false
	for _, member := range c.Members {
		if member.GetId() == createdByUserID {
			isMember = true
			break
		}
	}
	if !isMember {
		return nil, errors.New("user is not a member of the channel")
	}
	inviteCode, err := c.generateInviteCode()
	if err != nil {
		return nil, fmt.Errorf("failed to generate invite code: %w", err)
	}

	invite := NewChannelInvite(
		c.ID,
		createdByUserID,
		inviteCode,
		expiresAt,
		maxUses,
	)

	c.Invites = append(c.Invites, *invite)
	c.Version++

	c.addEvent(CreateChannelInviteCreatedEvent(invite, c))

	return invite, nil
}

// AcceptInvite accepts an invite to a channel and adds the user to the channel
func (c *Channel) AcceptInvite(inviteCode string, userID uuid.UUID) error {
	var targetInvite *ChannelInvite
	for i := range c.Invites {
		if c.Invites[i].GetInviteCode() == inviteCode {
			targetInvite = &c.Invites[i]
			break
		}
	}

	if targetInvite == nil {
		return errors.New("invite not found")
	}

	if !targetInvite.CanBeUsed() {
		return errors.New("invite has expired or reached max uses")
	}

	for _, member := range c.Members {
		if member.GetId() == userID {
			return nil
		}
	}

	err := targetInvite.Use()
	if err != nil {
		return nil
	}

	err = c.AddMember(userID)
	if err != nil {
		return err
	}

	c.addEvent(CreateChannelInviteUsedEvent(targetInvite, c))
	return nil
}

// DeactivateInvite deactivates an invite to a channel
func (c *Channel) DeactivateInvite(inviteID uuid.UUID, userID uuid.UUID) error {
	if c.CreatorUserID != userID {
		var inviteCreatorID uuid.UUID
		for _, invite := range c.Invites {
			if invite.GetID() == inviteID {
				inviteCreatorID = invite.GetCreatedByUserID()
				break
			}
		}

		if inviteCreatorID != userID {
			return errors.New("user does not have permission to deactivate this invite")
		}
	}

	for i := range c.Invites {
		if c.Invites[i].GetID() == inviteID {
			c.Invites[i].Deactivate()
			c.Version++
			c.addEvent(CreateChannelInviteDeactivatedEvent(&c.Invites[i], c))
			return nil
		}
	}

	return errors.New("invite not found")
}

// GetActiveInvites gets all active invites for a channel
func (c *Channel) GetActiveInvites() []ChannelInvite {
	var activeInvites []ChannelInvite
	for _, invite := range c.Invites {
		if invite.GetIsActive() {
			activeInvites = append(activeInvites, invite)
		}
	}
	return activeInvites
}

// generateInviteCode generates a random 4 byte hexcode
func (c *Channel) generateInviteCode() (string, error) {
	bytes := make([]byte, 4)

	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
