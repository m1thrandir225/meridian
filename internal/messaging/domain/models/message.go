package models

import (
	"time"

	"github.com/google/uuid"
)

type Message struct {
	ID              uuid.UUID
	ChannelID       uuid.UUID
	SenderUserID    *uuid.UUID // may be nil
	IntegrationID   *uuid.UUID // may be nil
	Content         MessageContent
	Timestamp       time.Time
	ParentMessageID *uuid.UUID
	Reactions       []Reaction
}
