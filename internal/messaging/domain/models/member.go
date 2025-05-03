package models

import (
	"time"

	"github.com/google/uuid"
)

type Member struct {
	ID       uuid.UUID
	Role     string
	JoinedAt time.Time
	LastRead time.Time
}
