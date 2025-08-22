package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// ChannelInvite represents an invite to a channel
type ChannelInvite struct {
	ID              uuid.UUID
	ChannelID       uuid.UUID
	CreatedByUserID uuid.UUID
	InviteCode      string
	ExpiresAt       time.Time
	MaxUse          *int
	CurrentUses     int
	CreatedAt       time.Time
	IsActive        bool
}

func NewChannelInvite(channelID, createdByUserID uuid.UUID, inviteCode string, expiresAt time.Time, maxUse *int) *ChannelInvite {
	return &ChannelInvite{
		ID:              uuid.New(),
		ChannelID:       channelID,
		CreatedByUserID: createdByUserID,
		InviteCode:      inviteCode,
		ExpiresAt:       expiresAt,
		MaxUse:          maxUse,
		CurrentUses:     0,
		CreatedAt:       time.Now().UTC(),
		IsActive:        true,
	}
}

func (ci *ChannelInvite) GetID() uuid.UUID {
	return ci.ID
}

func (ci *ChannelInvite) GetChannelID() uuid.UUID {
	return ci.ChannelID
}

func (ci *ChannelInvite) GetCreatedByUserID() uuid.UUID {
	return ci.CreatedByUserID
}

func (ci *ChannelInvite) GetInviteCode() string {
	return ci.InviteCode
}

func (ci *ChannelInvite) GetExpiresAt() time.Time {
	return ci.ExpiresAt
}

func (ci *ChannelInvite) GetMaxUse() *int {
	return ci.MaxUse
}

func (ci *ChannelInvite) GetCurrentUses() int {
	return ci.CurrentUses
}

func (ci *ChannelInvite) GetCreatedAt() time.Time {
	return ci.CreatedAt
}

func (ci *ChannelInvite) GetIsActive() bool {
	return ci.IsActive
}

func (ci *ChannelInvite) IsExpired() bool {
	return ci.ExpiresAt.Before(time.Now().UTC())
}

func (ci *ChannelInvite) HasReachedMaxUse() bool {
	return ci.MaxUse != nil && ci.CurrentUses >= *ci.MaxUse
}

func (ci *ChannelInvite) CanBeUsed() bool {
	return ci.IsActive && !ci.IsExpired() && !ci.HasReachedMaxUse()
}

func (ci *ChannelInvite) Use() error {
	if !ci.CanBeUsed() {
		return errors.New("invite cannot be used")
	}
	ci.CurrentUses++

	if ci.HasReachedMaxUse() {
		ci.Deactivate()
	}

	//TODO: add event for channel invite used
	return nil
}

func (ci *ChannelInvite) Deactivate() {
	ci.IsActive = false

}
