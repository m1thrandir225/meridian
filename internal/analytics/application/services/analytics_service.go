package services

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/m1thrandir225/meridian/internal/analytics/domain"
	"github.com/m1thrandir225/meridian/internal/analytics/infrastructure/persistence"
	"github.com/m1thrandir225/meridian/pkg/logging"
	"go.uber.org/zap"
)

type AnalyticsService struct {
	repo   persistence.AnalyticsRepository
	logger *logging.Logger
}

func NewAnalyticsService(repo persistence.AnalyticsRepository, logger *logging.Logger) *AnalyticsService {
	return &AnalyticsService{
		repo:   repo,
		logger: logger,
	}
}

// TrackUserRegistration tracks when a new user registers
func (s *AnalyticsService) TrackUserRegistration(ctx context.Context, cmd domain.TrackUserRegistrationCommand) error {
	logger := s.logger.WithMethod("TrackUserRegistration")
	logger.Info("Tracking user registration", zap.String("user_id", cmd.UserID))

	userID, err := uuid.Parse(cmd.UserID)
	if err != nil {
		logger.Error("Invalid user ID", zap.Error(err))
		return err
	}

	metricID, err := domain.NewMetricID()
	if err != nil {
		logger.Error("Failed to create metric ID", zap.Error(err))
		return err
	}

	metric := &domain.AnalyticsMetric{
		ID:        *metricID,
		Name:      "user_registration",
		Value:     1,
		UserID:    &userID,
		Timestamp: cmd.Timestamp,
		Version:   1,
	}

	if err := s.repo.SaveMetric(ctx, metric); err != nil {
		logger.Error("Failed to save user registration metric", zap.Error(err))
		return err
	}

	// Update daily aggregation
	dailyMetric := &domain.DailyMetric{
		ID:      *metricID,
		Name:    "user_registrations_daily",
		Date:    cmd.Timestamp.Truncate(24 * time.Hour),
		Value:   1,
		Count:   1,
		Version: 1,
	}

	if err := s.repo.SaveDailyMetric(ctx, dailyMetric); err != nil {
		logger.Error("Failed to save daily user registration metric", zap.Error(err))
		return err
	}

	logger.Info("User registration tracked successfully", zap.String("user_id", cmd.UserID))
	return nil
}

// TrackMessageSent tracks when a message is sent
func (s *AnalyticsService) TrackMessageSent(ctx context.Context, cmd domain.TrackMessageSentCommand) error {
	logger := s.logger.WithMethod("TrackMessageSent")
	logger.Info("Tracking message sent", zap.String("message_id", cmd.MessageID))

	channelID, err := uuid.Parse(cmd.ChannelID)
	if err != nil {
		logger.Error("Invalid channel ID", zap.Error(err))
		return err
	}

	senderID, err := uuid.Parse(cmd.SenderID)
	if err != nil {
		logger.Error("Invalid sender ID", zap.Error(err))
		return err
	}

	metricID, err := domain.NewMetricID()
	if err != nil {
		logger.Error("Failed to create metric ID", zap.Error(err))
		return err
	}

	// Track message metric
	metric := &domain.AnalyticsMetric{
		ID:        *metricID,
		Name:      "message_sent",
		Value:     1,
		ChannelID: &channelID,
		UserID:    &senderID,
		Timestamp: cmd.Timestamp,
		Metadata: map[string]interface{}{
			"content_length": cmd.ContentLength,
		},
		Version: 1,
	}

	if err := s.repo.SaveMetric(ctx, metric); err != nil {
		logger.Error("Failed to save message metric", zap.Error(err))
		return err
	}

	// Update daily aggregation
	dailyMetric := &domain.DailyMetric{
		ID:      *metricID,
		Name:    "messages_daily",
		Date:    cmd.Timestamp.Truncate(24 * time.Hour),
		Value:   1,
		Count:   1,
		Version: 1,
	}

	if err := s.repo.SaveDailyMetric(ctx, dailyMetric); err != nil {
		logger.Error("Failed to save daily message metric", zap.Error(err))
		return err
	}

	// Update hourly aggregation
	hourlyMetric := &domain.HourlyMetric{
		ID:       *metricID,
		Name:     "messages_hourly",
		DateTime: cmd.Timestamp.Truncate(time.Hour),
		Value:    1,
		Count:    1,
		Version:  1,
	}

	if err := s.repo.SaveHourlyMetric(ctx, hourlyMetric); err != nil {
		logger.Error("Failed to save hourly message metric", zap.Error(err))
		return err
	}

	// Update user activity for message
	if err := s.updateUserActivityForMessage(ctx, senderID, cmd.Timestamp); err != nil {
		logger.Error("Failed to update user activity", zap.Error(err))
		return err
	}

	// Update channel activity for message
	if err := s.updateChannelActivityForMessage(ctx, channelID, cmd.Timestamp); err != nil {
		logger.Error("Failed to update channel activity", zap.Error(err))
		return err
	}

	logger.Info("Message sent tracked successfully", zap.String("message_id", cmd.MessageID))
	return nil
}

// TrackChannelCreated tracks when a new channel is created
func (s *AnalyticsService) TrackChannelCreated(ctx context.Context, cmd domain.TrackChannelCreatedCommand) error {
	logger := s.logger.WithMethod("TrackChannelCreated")
	logger.Info("Tracking channel created", zap.String("channel_id", cmd.ChannelID))

	channelID, err := uuid.Parse(cmd.ChannelID)
	if err != nil {
		logger.Error("Invalid channel ID", zap.Error(err))
		return err
	}

	creatorID, err := uuid.Parse(cmd.CreatorID)
	if err != nil {
		logger.Error("Invalid creator ID", zap.Error(err))
		return err
	}

	metricID, err := domain.NewMetricID()
	if err != nil {
		logger.Error("Failed to create metric ID", zap.Error(err))
		return err
	}

	metric := &domain.AnalyticsMetric{
		ID:        *metricID,
		Name:      "channel_created",
		Value:     1,
		ChannelID: &channelID,
		UserID:    &creatorID,
		Timestamp: cmd.Timestamp,
		Version:   1,
	}

	if err := s.repo.SaveMetric(ctx, metric); err != nil {
		logger.Error("Failed to save channel created metric", zap.Error(err))
		return err
	}

	// Initialize channel activity
	channelActivity := &domain.ChannelActivity{
		ID:            *metricID,
		ChannelID:     channelID,
		MessagesCount: 0,
		MembersCount:  1,
		LastMessageAt: cmd.Timestamp,
		ActivityScore: 0,
		Version:       1,
	}

	if err := s.repo.SaveChannelActivity(ctx, channelActivity); err != nil {
		logger.Error("Failed to save channel activity", zap.Error(err))
		return err
	}

	logger.Info("Channel created tracked successfully", zap.String("channel_id", cmd.ChannelID))
	return nil
}

// TrackUserJoinedChannel tracks when a user joins a channel
func (s *AnalyticsService) TrackUserJoinedChannel(ctx context.Context, cmd domain.TrackUserJoinedChannelCommand) error {
	logger := s.logger.WithMethod("TrackUserJoinedChannel")
	logger.Info("Tracking user joined channel", zap.String("user_id", cmd.UserID), zap.String("channel_id", cmd.ChannelID))

	userID, err := uuid.Parse(cmd.UserID)
	if err != nil {
		logger.Error("Invalid user ID", zap.Error(err))
		return err
	}

	channelID, err := uuid.Parse(cmd.ChannelID)
	if err != nil {
		logger.Error("Invalid channel ID", zap.Error(err))
		return err
	}

	metricID, err := domain.NewMetricID()
	if err != nil {
		logger.Error("Failed to create metric ID", zap.Error(err))
		return err
	}

	metric := &domain.AnalyticsMetric{
		ID:        *metricID,
		Name:      "user_joined_channel",
		Value:     1,
		ChannelID: &channelID,
		UserID:    &userID,
		Timestamp: cmd.Timestamp,
		Version:   1,
	}

	if err := s.repo.SaveMetric(ctx, metric); err != nil {
		logger.Error("Failed to save user joined channel metric", zap.Error(err))
		return err
	}

	// Update user activity for channel join
	if err := s.updateUserActivityForChannelJoin(ctx, userID, cmd.Timestamp); err != nil {
		logger.Error("Failed to update user activity", zap.Error(err))
		return err
	}

	// Update channel activity for member join
	if err := s.updateChannelActivityForMemberJoin(ctx, channelID, cmd.Timestamp); err != nil {
		logger.Error("Failed to update channel activity", zap.Error(err))
		return err
	}

	logger.Info("User joined channel tracked successfully", zap.String("user_id", cmd.UserID), zap.String("channel_id", cmd.ChannelID))
	return nil
}

// TrackReactionAdded tracks when a reaction is added to a message
func (s *AnalyticsService) TrackReactionAdded(ctx context.Context, cmd domain.TrackReactionAddedCommand) error {
	logger := s.logger.WithMethod("TrackReactionAdded")
	logger.Info("Tracking reaction added", zap.String("reaction_id", cmd.ReactionID))

	userID, err := uuid.Parse(cmd.UserID)
	if err != nil {
		logger.Error("Invalid user ID", zap.Error(err))
		return err
	}

	metricID, err := domain.NewMetricID()
	if err != nil {
		logger.Error("Failed to create metric ID", zap.Error(err))
		return err
	}

	metric := &domain.AnalyticsMetric{
		ID:        *metricID,
		Name:      "reaction_added",
		Value:     1,
		UserID:    &userID,
		Timestamp: cmd.Timestamp,
		Metadata: map[string]interface{}{
			"reaction_type": cmd.ReactionType,
			"message_id":    cmd.MessageID,
		},
		Version: 1,
	}

	if err := s.repo.SaveMetric(ctx, metric); err != nil {
		logger.Error("Failed to save reaction metric", zap.Error(err))
		return err
	}

	// Update user activity for reaction
	if err := s.updateUserActivityForReaction(ctx, userID, cmd.Timestamp); err != nil {
		logger.Error("Failed to update user activity", zap.Error(err))
		return err
	}

	logger.Info("Reaction added tracked successfully", zap.String("reaction_id", cmd.ReactionID))
	return nil
}

// GetDashboardData retrieves dashboard metrics
func (s *AnalyticsService) GetDashboardData(ctx context.Context, query domain.GetDashboardDataQuery) (*domain.DashboardData, error) {
	logger := s.logger.WithMethod("GetDashboardData")
	logger.Info("Getting dashboard data")

	now := time.Now()
	startDate := now.Add(-query.TimeRange)

	totalUsers, err := s.repo.GetTotalUsers(ctx)
	if err != nil {
		logger.Error("Failed to get total users", zap.Error(err))
		return nil, err
	}

	activeUsers, err := s.repo.GetActiveUsers(ctx, query.TimeRange)
	if err != nil {
		logger.Error("Failed to get active users", zap.Error(err))
		return nil, err
	}

	newUsersToday, err := s.repo.GetNewUsersCount(ctx, now.Truncate(24*time.Hour))
	if err != nil {
		logger.Error("Failed to get new users today", zap.Error(err))
		return nil, err
	}

	messagesToday, err := s.repo.GetMessagesCount(ctx, now.Truncate(24*time.Hour))
	if err != nil {
		logger.Error("Failed to get messages today", zap.Error(err))
		return nil, err
	}

	totalChannels, err := s.repo.GetTotalChannels(ctx)
	if err != nil {
		logger.Error("Failed to get total channels", zap.Error(err))
		return nil, err
	}

	activeChannels, err := s.repo.GetActiveChannels(ctx, query.TimeRange)
	if err != nil {
		logger.Error("Failed to get active channels", zap.Error(err))
		return nil, err
	}

	averageMessagesPerUser := float64(0)
	if totalUsers > 0 {
		totalMessages, err := s.repo.GetTotalMessages(ctx)
		if err != nil {
			logger.Error("Failed to get total messages", zap.Error(err))
			return nil, err
		}
		averageMessagesPerUser = float64(totalMessages) / float64(totalUsers)
	}

	peakHour, err := s.repo.GetPeakHour(ctx, startDate, now)
	if err != nil {
		logger.Error("Failed to get peak hour", zap.Error(err))
		return nil, err
	}

	dashboardData := &domain.DashboardData{
		TotalUsers:             totalUsers,
		ActiveUsers:            activeUsers,
		NewUsersToday:          newUsersToday,
		MessagesToday:          messagesToday,
		TotalChannels:          totalChannels,
		ActiveChannels:         activeChannels,
		AverageMessagesPerUser: averageMessagesPerUser,
		PeakHour:               peakHour,
		LastUpdated:            now,
	}

	logger.Info("Dashboard data retrieved successfully")
	return dashboardData, nil
}

// GetUserGrowth retrieves user growth data
func (s *AnalyticsService) GetUserGrowth(ctx context.Context, query domain.GetUserGrowthQuery) ([]domain.UserGrowthData, error) {
	logger := s.logger.WithMethod("GetUserGrowth")
	logger.Info("Getting user growth data")

	growthData, err := s.repo.GetUserGrowth(ctx, query.StartDate, query.EndDate, query.Interval)
	if err != nil {
		logger.Error("Failed to get user growth data", zap.Error(err))
		return nil, err
	}

	logger.Info("User growth data retrieved successfully", zap.Int("data_points", len(growthData)))
	return growthData, nil
}

// GetMessageVolume retrieves message volume data
func (s *AnalyticsService) GetMessageVolume(ctx context.Context, query domain.GetMessageVolumeQuery) ([]domain.MessageVolumeData, error) {
	logger := s.logger.WithMethod("GetMessageVolume")
	logger.Info("Getting message volume data")

	volumeData, err := s.repo.GetMessageVolume(ctx, query.StartDate, query.EndDate, query.ChannelID)
	if err != nil {
		logger.Error("Failed to get message volume data", zap.Error(err))
		return nil, err
	}

	logger.Info("Message volume data retrieved successfully", zap.Int("data_points", len(volumeData)))
	return volumeData, nil
}

// GetChannelActivity retrieves channel activity data
func (s *AnalyticsService) GetChannelActivity(ctx context.Context, query domain.GetChannelActivityQuery) ([]domain.ChannelActivityData, error) {
	logger := s.logger.WithMethod("GetChannelActivity")
	logger.Info("Getting channel activity data")

	activityData, err := s.repo.GetChannelActivityList(ctx, query.StartDate, query.EndDate, query.Limit)
	if err != nil {
		logger.Error("Failed to get channel activity data", zap.Error(err))
		return nil, err
	}

	logger.Info("Channel activity data retrieved successfully", zap.Int("data_points", len(activityData)))
	return activityData, nil
}

// GetTopUsers retrieves top active users
func (s *AnalyticsService) GetTopUsers(ctx context.Context, query domain.GetTopUsersQuery) ([]domain.TopUserData, error) {
	logger := s.logger.WithMethod("GetTopUsers")
	logger.Info("Getting top users data")

	topUsers, err := s.repo.GetTopUsers(ctx, query.StartDate, query.EndDate, query.Limit)
	if err != nil {
		logger.Error("Failed to get top users data", zap.Error(err))
		return nil, err
	}

	logger.Info("Top users data retrieved successfully", zap.Int("users_count", len(topUsers)))
	return topUsers, nil
}

// GetReactionUsage retrieves reaction usage statistics
func (s *AnalyticsService) GetReactionUsage(ctx context.Context, query domain.GetReactionUsageQuery) ([]domain.ReactionUsageData, error) {
	logger := s.logger.WithMethod("GetReactionUsage")
	logger.Info("Getting reaction usage data")

	reactionData, err := s.repo.GetReactionUsage(ctx, query.StartDate, query.EndDate)
	if err != nil {
		logger.Error("Failed to get reaction usage data", zap.Error(err))
		return nil, err
	}

	logger.Info("Reaction usage data retrieved successfully", zap.Int("reaction_types", len(reactionData)))
	return reactionData, nil
}

// Helper methods
func (s *AnalyticsService) updateUserActivity(ctx context.Context, userID uuid.UUID, timestamp time.Time) error {
	activity, err := s.repo.GetUserActivity(ctx, userID)
	if err != nil {
		// Create new user activity if not exists
		metricID, err := domain.NewMetricID()
		if err != nil {
			return err
		}

		activity = &domain.UserActivity{
			ID:              *metricID,
			UserID:          userID,
			LastActiveAt:    timestamp,
			MessagesSent:    0,
			ChannelsJoined:  0,
			ReactionsGiven:  0,
			SessionDuration: 0,
			Version:         1,
		}
	}

	activity.LastActiveAt = timestamp
	activity.MessagesSent++
	activity.Version++

	return s.repo.SaveUserActivity(ctx, activity)
}

func (s *AnalyticsService) updateChannelActivity(ctx context.Context, channelID uuid.UUID, timestamp time.Time) error {
	activity, err := s.repo.GetChannelActivity(ctx, channelID)
	if err != nil {
		// Create new channel activity if not exists
		metricID, err := domain.NewMetricID()
		if err != nil {
			return err
		}

		activity = &domain.ChannelActivity{
			ID:            *metricID,
			ChannelID:     channelID,
			MessagesCount: 0,
			MembersCount:  0,
			LastMessageAt: timestamp,
			ActivityScore: 0,
			Version:       1,
		}
	}

	activity.LastMessageAt = timestamp
	activity.MessagesCount++
	activity.Version++

	return s.repo.SaveChannelActivity(ctx, activity)
}

func (s *AnalyticsService) updateUserActivityForMessage(ctx context.Context, userID uuid.UUID, timestamp time.Time) error {
	activity, err := s.repo.GetUserActivity(ctx, userID)
	if err != nil {
		// Create new user activity if not exists
		metricID, err := domain.NewMetricID()
		if err != nil {
			return err
		}

		activity = &domain.UserActivity{
			ID:              *metricID,
			UserID:          userID,
			LastActiveAt:    timestamp,
			MessagesSent:    0,
			ChannelsJoined:  0,
			ReactionsGiven:  0,
			SessionDuration: 0,
			Version:         1,
		}
	}

	activity.LastActiveAt = timestamp
	activity.MessagesSent++
	activity.Version++

	return s.repo.SaveUserActivity(ctx, activity)
}

func (s *AnalyticsService) updateUserActivityForChannelJoin(ctx context.Context, userID uuid.UUID, timestamp time.Time) error {
	activity, err := s.repo.GetUserActivity(ctx, userID)
	if err != nil {
		// Create new user activity if not exists
		metricID, err := domain.NewMetricID()
		if err != nil {
			return err
		}

		activity = &domain.UserActivity{
			ID:              *metricID,
			UserID:          userID,
			LastActiveAt:    timestamp,
			MessagesSent:    0,
			ChannelsJoined:  0,
			ReactionsGiven:  0,
			SessionDuration: 0,
			Version:         1,
		}
	}

	activity.LastActiveAt = timestamp
	activity.ChannelsJoined++
	activity.Version++

	return s.repo.SaveUserActivity(ctx, activity)
}

func (s *AnalyticsService) updateUserActivityForReaction(ctx context.Context, userID uuid.UUID, timestamp time.Time) error {
	activity, err := s.repo.GetUserActivity(ctx, userID)
	if err != nil {
		// Create new user activity if not exists
		metricID, err := domain.NewMetricID()
		if err != nil {
			return err
		}

		activity = &domain.UserActivity{
			ID:              *metricID,
			UserID:          userID,
			LastActiveAt:    timestamp,
			MessagesSent:    0,
			ChannelsJoined:  0,
			ReactionsGiven:  0,
			SessionDuration: 0,
			Version:         1,
		}
	}

	activity.LastActiveAt = timestamp
	activity.ReactionsGiven++
	activity.Version++

	return s.repo.SaveUserActivity(ctx, activity)
}

func (s *AnalyticsService) updateChannelActivityForMessage(ctx context.Context, channelID uuid.UUID, timestamp time.Time) error {
	activity, err := s.repo.GetChannelActivity(ctx, channelID)
	if err != nil {
		// Create new channel activity if not exists
		metricID, err := domain.NewMetricID()
		if err != nil {
			return err
		}

		activity = &domain.ChannelActivity{
			ID:            *metricID,
			ChannelID:     channelID,
			MessagesCount: 0,
			MembersCount:  0,
			LastMessageAt: timestamp,
			ActivityScore: 0,
			Version:       1,
		}
	}

	activity.LastMessageAt = timestamp
	activity.MessagesCount++
	activity.Version++

	return s.repo.SaveChannelActivity(ctx, activity)
}

func (s *AnalyticsService) updateChannelActivityForMemberJoin(ctx context.Context, channelID uuid.UUID, timestamp time.Time) error {
	activity, err := s.repo.GetChannelActivity(ctx, channelID)
	if err != nil {
		// Create new channel activity if not exists
		metricID, err := domain.NewMetricID()
		if err != nil {
			return err
		}

		activity = &domain.ChannelActivity{
			ID:            *metricID,
			ChannelID:     channelID,
			MessagesCount: 0,
			MembersCount:  0,
			LastMessageAt: timestamp,
			ActivityScore: 0,
			Version:       1,
		}
	}

	activity.LastMessageAt = timestamp
	activity.MembersCount++
	activity.Version++

	return s.repo.SaveChannelActivity(ctx, activity)
}
