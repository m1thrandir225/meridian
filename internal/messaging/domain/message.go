package domain

import (
	"time"

	"github.com/google/uuid"
)

type Message struct {
	id              uuid.UUID
	channelId       uuid.UUID
	senderUserId    *uuid.UUID // may be nil
	integrationId   *uuid.UUID // may be nil
	content         MessageContent
	createdAt       time.Time
	parentMessageId *uuid.UUID
	reactions       []Reaction
}

func newMessage(id uuid.UUID, channelId uuid.UUID, senderUserId, integrationId, parentMessageId *uuid.UUID, content MessageContent, reactions []Reaction, timestamp time.Time) Message {
	return Message{
		id:              id,
		channelId:       channelId,
		senderUserId:    senderUserId,
		integrationId:   integrationId,
		parentMessageId: parentMessageId,
		content:         content,
		reactions:       reactions,
		createdAt:       timestamp,
	}
}

// For external usage
func RehydrateMessage(id uuid.UUID, channelId uuid.UUID, senderUserId, integrationId, parentMessageId *uuid.UUID, content MessageContent, reactions []Reaction, timestamp time.Time) Message {
	return Message{
		id:              id,
		channelId:       channelId,
		senderUserId:    senderUserId,
		integrationId:   integrationId,
		parentMessageId: parentMessageId,
		content:         content,
		reactions:       reactions,
		createdAt:       timestamp,
	}
}

func (m *Message) GetId() uuid.UUID {
	return m.id
}

func (m *Message) setId(id uuid.UUID) {
	m.id = id
}

func (m *Message) GetChannelId() uuid.UUID {
	return m.channelId
}

func (m *Message) setChannelId(id uuid.UUID) {
	m.channelId = id
}

func (m *Message) GetSenderUserId() *uuid.UUID {
	return m.senderUserId
}

func (m *Message) setSenderUserId(id *uuid.UUID) {
	m.senderUserId = id
}

func (m *Message) GetIntegrationId() *uuid.UUID {
	return m.integrationId
}

func (m *Message) setIntegrationId(id *uuid.UUID) {
	m.integrationId = id
}

func (m *Message) GetContent() *MessageContent {
	return &m.content
}

func (m *Message) setContent(content MessageContent) {
	m.content = content
}

func (m *Message) GetCreatedAt() time.Time {
	return m.createdAt
}

func (m *Message) setCreatedAt(timestamp time.Time) {
	m.createdAt = timestamp
}

func (m *Message) GetParentMessageId() *uuid.UUID {
	return m.parentMessageId
}

func (m *Message) setParentMessageId(id *uuid.UUID) {
	m.parentMessageId = id
}

func (m *Message) GetReactions() []Reaction {
	return m.reactions
}

func (m *Message) setReactions(reactions []Reaction) {
	m.reactions = reactions
}
