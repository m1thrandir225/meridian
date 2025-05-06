package domain

import (
	"time"

	"github.com/google/uuid"
)

type Reaction struct {
	id           uuid.UUID
	messageId    uuid.UUID
	userId       uuid.UUID
	reactionType string
	createdAt    time.Time
}

func newReaction(id, messageId, userId uuid.UUID, reactionType string, createdAt time.Time) Reaction {
	return Reaction{
		id:           id,
		messageId:    messageId,
		userId:       userId,
		reactionType: reactionType,
		createdAt:    createdAt,
	}
}

func RehydrateReaction(id, messageId, userId uuid.UUID, reactionType string, createdAt time.Time) Reaction {
	return Reaction{
		id:           id,
		messageId:    messageId,
		userId:       userId,
		reactionType: reactionType,
		createdAt:    createdAt,
	}
}

func (r *Reaction) GetId() uuid.UUID {
	return r.id
}

func (r *Reaction) setId(id uuid.UUID) {
	r.id = id
}

func (r *Reaction) GetMessageId() uuid.UUID {
	return r.messageId
}

func (r *Reaction) setMessageId(id uuid.UUID) {
	r.messageId = id
}

func (r *Reaction) GetUserId() uuid.UUID {
	return r.userId
}

func (r *Reaction) setUserId(id uuid.UUID) {
	r.userId = id
}

func (r *Reaction) GetReactionType() string {
	return r.reactionType
}

func (r *Reaction) setReactionType(reaction string) {
	r.reactionType = reaction
}

func (r *Reaction) GetCreatedAt() time.Time {
	return r.createdAt
}

func (r *Reaction) setCreatedat(createdAt time.Time) {
	r.createdAt = createdAt
}
