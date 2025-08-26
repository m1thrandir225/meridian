package persistence

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/m1thrandir225/meridian/internal/analytics/domain"
)

type AnalyticsRepository interface {
	// Metrics
	SaveMetric(ctx context.Context, metric *domain.AnalyticsMetric) error
	GetMetrics(ctx context.Context, name string, startDate, endDate time.Time) ([]*domain.AnalyticsMetric, error)

	// Daily Metrics
	SaveDailyMetric(ctx context.Context, metric *domain.DailyMetric) error
	GetDailyMetrics(ctx context.Context, name string, startDate, endDate time.Time) ([]*domain.DailyMetric, error)

	// Hourly Metrics
	SaveHourlyMetric(ctx context.Context, metric *domain.HourlyMetric) error
	GetHourlyMetrics(ctx context.Context, name string, startDate, endDate time.Time) ([]*domain.HourlyMetric, error)

	// User Activity
	SaveUserActivity(ctx context.Context, activity *domain.UserActivity) error
	GetUserActivity(ctx context.Context, userID uuid.UUID) (*domain.UserActivity, error)
	GetActiveUsers(ctx context.Context, timeRange time.Duration) (int64, error)
	GetTotalUsers(ctx context.Context) (int64, error)
	GetNewUsersCount(ctx context.Context, date time.Time) (int64, error)

	// Channel Activity
	SaveChannelActivity(ctx context.Context, activity *domain.ChannelActivity) error
	GetChannelActivity(ctx context.Context, channelID uuid.UUID) (*domain.ChannelActivity, error)
	GetActiveChannels(ctx context.Context, timeRange time.Duration) (int64, error)
	GetTotalChannels(ctx context.Context) (int64, error)

	// Message Analytics
	GetMessagesCount(ctx context.Context, date time.Time) (int64, error)
	GetTotalMessages(ctx context.Context) (int64, error)
	GetPeakHour(ctx context.Context, startDate, endDate time.Time) (int, error)

	// Analytics Queries
	GetUserGrowth(ctx context.Context, startDate, endDate time.Time, interval string) ([]domain.UserGrowthData, error)
	GetMessageVolume(ctx context.Context, startDate, endDate time.Time, channelID *string) ([]domain.MessageVolumeData, error)
	GetChannelActivityList(ctx context.Context, startDate, endDate time.Time, limit int) ([]domain.ChannelActivityData, error)
	GetTopUsers(ctx context.Context, startDate, endDate time.Time, limit int) ([]domain.TopUserData, error)
	GetReactionUsage(ctx context.Context, startDate, endDate time.Time) ([]domain.ReactionUsageData, error)
}
