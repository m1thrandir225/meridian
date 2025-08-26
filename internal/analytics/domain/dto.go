package domain

import (
	"time"
)

// DashboardData represents the main dashboard metrics
type DashboardData struct {
	TotalUsers             int64     `json:"total_users"`
	ActiveUsers            int64     `json:"active_users"`
	NewUsersToday          int64     `json:"new_users_today"`
	MessagesToday          int64     `json:"messages_today"`
	TotalChannels          int64     `json:"total_channels"`
	ActiveChannels         int64     `json:"active_channels"`
	AverageMessagesPerUser float64   `json:"average_messages_per_user"`
	PeakHour               int       `json:"peak_hour"`
	LastUpdated            time.Time `json:"last_updated"`
}

// UserGrowthData represents user growth over time
type UserGrowthData struct {
	Period     string  `json:"period"`
	NewUsers   int64   `json:"new_users"`
	TotalUsers int64   `json:"total_users"`
	GrowthRate float64 `json:"growth_rate"`
}

// MessageVolumeData represents message volume over time
type MessageVolumeData struct {
	Period    string  `json:"period"`
	Messages  int64   `json:"messages"`
	Channels  int64   `json:"channels"`
	AvgLength float64 `json:"avg_length"`
}

// ChannelActivityData represents channel activity metrics
type ChannelActivityData struct {
	ChannelID     string    `json:"channel_id"`
	ChannelName   string    `json:"channel_name"`
	MessagesCount int64     `json:"messages_count"`
	MembersCount  int64     `json:"members_count"`
	LastMessageAt time.Time `json:"last_message_at"`
	ActivityScore float64   `json:"activity_score"`
}

// TopUserData represents top active users
type TopUserData struct {
	UserID         string    `json:"user_id"`
	Username       string    `json:"username"`
	MessagesSent   int64     `json:"messages_sent"`
	ChannelsJoined int64     `json:"channels_joined"`
	ReactionsGiven int64     `json:"reactions_given"`
	LastActiveAt   time.Time `json:"last_active_at"`
}

// ReactionUsageData represents reaction usage statistics
type ReactionUsageData struct {
	ReactionType string  `json:"reaction_type"`
	Count        int64   `json:"count"`
	Percentage   float64 `json:"percentage"`
}

// TimeSeriesData represents time-series data for charts
type TimeSeriesData struct {
	Timestamp time.Time `json:"timestamp"`
	Value     float64   `json:"value"`
	Count     int64     `json:"count"`
}
