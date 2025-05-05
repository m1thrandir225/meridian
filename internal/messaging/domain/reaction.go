package domain

import (
	"time"

	"github.com/google/uuid"
)

type Reaction struct {
	ID           uuid.UUID
	MessageID    uuid.UUID
	UserID       uuid.UUID
	ReactionType string
	Timestamp    time.Time
}
