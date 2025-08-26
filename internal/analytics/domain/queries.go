package domain

import (
	"time"
)

type Query interface {
	QueryName() string
}

type GetDashboardDataQuery struct {
	TimeRange time.Duration
}

func (q GetDashboardDataQuery) QueryName() string {
	return "GetDashboardData"
}

type GetUserGrowthQuery struct {
	StartDate time.Time
	EndDate   time.Time
	Interval  string // "daily", "weekly", "monthly"
}

func (q GetUserGrowthQuery) QueryName() string {
	return "GetUserGrowth"
}

type GetMessageVolumeQuery struct {
	StartDate time.Time
	EndDate   time.Time
	ChannelID *string
}

func (q GetMessageVolumeQuery) QueryName() string {
	return "GetMessageVolume"
}

type GetChannelActivityQuery struct {
	StartDate time.Time
	EndDate   time.Time
	Limit     int
}

func (q GetChannelActivityQuery) QueryName() string {
	return "GetChannelActivity"
}

type GetTopUsersQuery struct {
	StartDate time.Time
	EndDate   time.Time
	Limit     int
}

func (q GetTopUsersQuery) QueryName() string {
	return "GetTopUsers"
}

type GetReactionUsageQuery struct {
	StartDate time.Time
	EndDate   time.Time
}

func (q GetReactionUsageQuery) QueryName() string {
	return "GetReactionUsage"
}
