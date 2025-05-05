package domain

import "github.com/google/uuid"

// NOTE: might be redundant
type MessageContent struct {
	Text      string
	Mentions  []uuid.UUID
	Links     []string
	Formatted bool
}
