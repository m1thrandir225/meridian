package services

import (
	"context"
	"database/sql"
	"errors"
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

	// Check if user already exists
	existingUser, err := s.repo.GetUserActivity(ctx, userID)
	if err == nil && existingUser != nil {
		logger.Info("User already exists, skipping registration", zap.String("user_id", cmd.UserID))
		return nil
	}

	// Create new user activity
	metricID, err := domain.NewMetricID()
	if err != nil {
		logger.Error("Failed to create metric ID", zap.Error(err))
		return err
	}

	userActivity := &domain.UserActivity{
		ID:              *metricID,
		UserID:          userID,
		LastActiveAt:    cmd.Timestamp,
		MessagesSent:    0,
		ChannelsJoined:  0,
		ReactionsGiven:  0,
		SessionDuration: 0,
		Version:         1,
	}

	if err := s.repo.SaveUserActivity(ctx, userActivity); err != nil {
		logger.Error("Failed to save user activity", zap.Error(err))
		return err
	}

	// Save metric for tracking
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

	logger.Info("User registration tracked successfully", zap.String("user_id", cmd.UserID))
	return nil
}

// TrackMessageSent tracks when a message is sent
func (s *AnalyticsService) TrackMessageSent(ctx context.Context, cmd domain.TrackMessageSentCommand) error {
	logger := s.logger.WithMethod("TrackMessageSent")
	logger.Info("Tracking message sent", zap.String("message_id", cmd.MessageID))

	// Check if message already processed
	existingMetric, err := s.repo.GetMetricByMessageID(ctx, cmd.MessageID)
	if err == nil && existingMetric != nil {
		logger.Info("Message already processed, skipping", zap.String("message_id", cmd.MessageID))
		return nil
	}

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

	// Save the message metric
	metric := &domain.AnalyticsMetric{
		ID:        *metricID,
		Name:      "message_sent",
		Value:     1,
		ChannelID: &channelID,
		UserID:    &senderID,
		Timestamp: cmd.Timestamp,
		Metadata: map[string]interface{}{
			"content_length": cmd.ContentLength,
			"message_id":     cmd.MessageID,
		},
		Version: 1,
	}

	if err := s.repo.SaveMetric(ctx, metric); err != nil {
		logger.Error("Failed to save message metric", zap.Error(err))
		return err
	}

	// Update user activity
	if err := s.incrementUserMessages(ctx, senderID, cmd.Timestamp); err != nil {
		logger.Error("Failed to update user activity", zap.Error(err))
		return err
	}

	// Update channel activity
	if err := s.incrementChannelMessages(ctx, channelID, cmd.Timestamp); err != nil {
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

	// Check if channel already exists
	existingChannel, err := s.repo.GetChannelActivity(ctx, channelID)
	if err == nil && existingChannel != nil {
		logger.Info("Channel already exists, skipping creation", zap.String("channel_id", cmd.ChannelID))
		return nil
	}

	metricID, err := domain.NewMetricID()
	if err != nil {
		logger.Error("Failed to create metric ID", zap.Error(err))
		return err
	}

	// Create channel activity
	channelActivity := &domain.ChannelActivity{
		ID:            *metricID,
		ChannelID:     channelID,
		MessagesCount: 0,
		MembersCount:  1, // Creator is the first member
		LastMessageAt: cmd.Timestamp,
		ActivityScore: 2.0, // Base score for creation
		Version:       1,
	}

	if err := s.repo.SaveChannelActivity(ctx, channelActivity); err != nil {
		logger.Error("Failed to save channel activity", zap.Error(err))
		return err
	}

	// Save metric
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

	// Save metric
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

	// Update user activity
	if err := s.incrementUserChannelsJoined(ctx, userID, cmd.Timestamp); err != nil {
		logger.Error("Failed to update user activity", zap.Error(err))
		return err
	}

	// Update channel activity
	if err := s.incrementChannelMembers(ctx, channelID, cmd.Timestamp); err != nil {
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

	// Save metric
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

	// Update user activity
	if err := s.incrementUserReactions(ctx, userID, cmd.Timestamp); err != nil {
		logger.Error("Failed to update user activity", zap.Error(err))
		return err
	}

	logger.Info("Reaction added tracked successfully", zap.String("reaction_id", cmd.ReactionID))
	return nil
}

// Helper methods for atomic increments
func (s *AnalyticsService) incrementUserMessages(ctx context.Context, userID uuid.UUID, timestamp time.Time) error {
	return s.repo.IncrementUserMessages(ctx, userID, timestamp)
}

func (s *AnalyticsService) incrementUserChannelsJoined(ctx context.Context, userID uuid.UUID, timestamp time.Time) error {
	return s.repo.IncrementUserChannelsJoined(ctx, userID, timestamp)
}

func (s *AnalyticsService) incrementUserReactions(ctx context.Context, userID uuid.UUID, timestamp time.Time) error {
	return s.repo.IncrementUserReactions(ctx, userID, timestamp)
}

func (s *AnalyticsService) incrementChannelMessages(ctx context.Context, channelID uuid.UUID, timestamp time.Time) error {
	return s.repo.IncrementChannelMessages(ctx, channelID, timestamp)
}

func (s *AnalyticsService) incrementChannelMembers(ctx context.Context, channelID uuid.UUID, timestamp time.Time) error {
	return s.repo.IncrementChannelMembers(ctx, channelID, timestamp)
}

// GetDashboardData returns basic dashboard data for the system
func (s *AnalyticsService) GetDashboardData(ctx context.Context, query domain.GetDashboardDataQuery) (*domain.DashboardData, error) {
	logger := s.logger.WithMethod("GetDashboardData")
	logger.Info("Getting dashboard data")

	now := time.Now()
	startDate := now.Add(-query.TimeRange)

	totalUsers, err := s.repo.GetTotalUsers(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			logger.Error("No users found", zap.Error(err))
		} else {
			logger.Error("Failed to get total users", zap.Error(err))
			return nil, err
		}
	}

	activeUsers, err := s.repo.GetActiveUsers(ctx, query.TimeRange)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			logger.Error("No users found", zap.Error(err))
		} else {
			logger.Error("Failed to get active users", zap.Error(err))
			return nil, err
		}
	}

	newUsersToday, err := s.repo.GetNewUsersCount(ctx, now.Truncate(24*time.Hour))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			logger.Error("No users found", zap.Error(err))
		} else {
			logger.Error("Failed to get new users today", zap.Error(err))
			return nil, err
		}
	}

	messagesToday, err := s.repo.GetMessagesCount(ctx, now.Truncate(24*time.Hour))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			logger.Error("No messages found", zap.Error(err))
		} else {
			logger.Error("Failed to get messages today", zap.Error(err))
			return nil, err
		}
	}

	totalChannels, err := s.repo.GetTotalChannels(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			logger.Error("No channels found", zap.Error(err))
		} else {
			logger.Error("Failed to get total channels", zap.Error(err))
			return nil, err
		}
	}

	activeChannels, err := s.repo.GetActiveChannels(ctx, query.TimeRange)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			logger.Error("No channels found", zap.Error(err))
		} else {
			logger.Error("Failed to get active channels", zap.Error(err))
			return nil, err
		}
	}

	averageMessagesPerUser := float64(0)
	if totalUsers > 0 {
		totalMessages, err := s.repo.GetTotalMessages(ctx)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				logger.Error("No messages in the system.", zap.Error(err))
			} else {
				logger.Error("Failed to get total messages", zap.Error(err))
				return nil, err
			}
		}
		averageMessagesPerUser = float64(totalMessages) / float64(totalUsers)
	}

	peakHour, err := s.repo.GetPeakHour(ctx, startDate, now)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			logger.Error("No peak our yet", zap.Error(err))
		} else {
			logger.Error("Failed to get peak hour", zap.Error(err))
			return nil, err
		}
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

// Other query methods remain the same...
func (s *AnalyticsService) GetUserGrowth(ctx context.Context, query domain.GetUserGrowthQuery) ([]domain.UserGrowthData, error) {
	return s.repo.GetUserGrowth(ctx, query.StartDate, query.EndDate, query.Interval)
}

func (s *AnalyticsService) GetMessageVolume(ctx context.Context, query domain.GetMessageVolumeQuery) ([]domain.MessageVolumeData, error) {
	return s.repo.GetMessageVolume(ctx, query.StartDate, query.EndDate, query.ChannelID)
}

func (s *AnalyticsService) GetChannelActivity(ctx context.Context, query domain.GetChannelActivityQuery) ([]domain.ChannelActivityData, error) {
	return s.repo.GetChannelActivityList(ctx, query.StartDate, query.EndDate, query.Limit)
}

func (s *AnalyticsService) GetTopUsers(ctx context.Context, query domain.GetTopUsersQuery) ([]domain.TopUserData, error) {
	return s.repo.GetTopUsers(ctx, query.StartDate, query.EndDate, query.Limit)
}

func (s *AnalyticsService) GetReactionUsage(ctx context.Context, query domain.GetReactionUsageQuery) ([]domain.ReactionUsageData, error) {
	return s.repo.GetReactionUsage(ctx, query.StartDate, query.EndDate)
}
