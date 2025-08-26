package domain

import (
	"time"

	"github.com/google/uuid"
)

// UserActivity represents user engagement metrics
type UserActivity struct {
	ID              MetricID
	UserID          uuid.UUID
	LastActiveAt    time.Time
	MessagesSent    int64
	ChannelsJoined  int64
	ReactionsGiven  int64
	SessionDuration time.Duration
	Version         int64
}
