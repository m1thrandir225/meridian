package domain

import (
	"time"
)

type Command interface {
	CommandName() string
}

type TrackUserRegistrationCommand struct {
	UserID    string
	Timestamp time.Time
}

func (c TrackUserRegistrationCommand) CommandName() string {
	return "TrackUserRegistration"
}

type TrackMessageSentCommand struct {
	MessageID     string
	ChannelID     string
	SenderID      string
	Timestamp     time.Time
	ContentLength int
}

func (c TrackMessageSentCommand) CommandName() string {
	return "TrackMessageSent"
}

type TrackChannelCreatedCommand struct {
	ChannelID string
	CreatorID string
	Timestamp time.Time
}

func (c TrackChannelCreatedCommand) CommandName() string {
	return "TrackChannelCreated"
}

type TrackUserJoinedChannelCommand struct {
	UserID    string
	ChannelID string
	Timestamp time.Time
}

func (c TrackUserJoinedChannelCommand) CommandName() string {
	return "TrackUserJoinedChannel"
}

type TrackReactionAddedCommand struct {
	ReactionID   string
	MessageID    string
	UserID       string
	ReactionType string
	Timestamp    time.Time
}

func (c TrackReactionAddedCommand) CommandName() string {
	return "TrackReactionAdded"
}

type GetAnalyticsCommand struct {
	StartDate time.Time
	EndDate   time.Time
	Metrics   []string
}

func (c GetAnalyticsCommand) CommandName() string {
	return "GetAnalytics"
}

type GetUserActivityCommand struct {
	UserID string
}

func (c GetUserActivityCommand) CommandName() string {
	return "GetUserActivity"
}

type GetChannelActivityCommand struct {
	ChannelID string
}

func (c GetChannelActivityCommand) CommandName() string {
	return "GetChannelActivity"
}
