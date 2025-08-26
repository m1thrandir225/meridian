package persistence

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/m1thrandir225/meridian/internal/analytics/domain"
)

type PostgresAnalyticsRepository struct {
	db *pgxpool.Pool
}

func NewPostgresAnalyticsRepository(db *pgxpool.Pool) *PostgresAnalyticsRepository {
	return &PostgresAnalyticsRepository{
		db: db,
	}
}

// Metrics
func (r *PostgresAnalyticsRepository) SaveMetric(ctx context.Context, metric *domain.AnalyticsMetric) error {
	query := `
		INSERT INTO analytics_metrics (id, name, value, channel_id, user_id, timestamp, metadata, version)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		ON CONFLICT (id) DO UPDATE SET
			value = EXCLUDED.value,
			metadata = EXCLUDED.metadata,
			version = EXCLUDED.version
	`

	metadataJSON, err := json.Marshal(metric.Metadata)
	if err != nil {
		return fmt.Errorf("failed to marshal metadata: %w", err)
	}

	_, err = r.db.Exec(ctx, query,
		metric.ID.Value(),
		metric.Name,
		metric.Value,
		metric.ChannelID,
		metric.UserID,
		metric.Timestamp,
		metadataJSON,
		metric.Version,
	)

	return err
}

func (r *PostgresAnalyticsRepository) GetMetrics(ctx context.Context, name string, startDate, endDate time.Time) ([]*domain.AnalyticsMetric, error) {
	query := `
		SELECT id, name, value, channel_id, user_id, timestamp, metadata, version
		FROM analytics_metrics
		WHERE name = $1 AND timestamp >= $2 AND timestamp <= $3
		ORDER BY timestamp ASC
	`

	rows, err := r.db.Query(ctx, query, name, startDate, endDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var metrics []*domain.AnalyticsMetric
	for rows.Next() {
		var metric domain.AnalyticsMetric
		var idStr string
		var metadataJSON []byte
		var channelID, userID sql.NullString

		err := rows.Scan(
			&idStr,
			&metric.Name,
			&metric.Value,
			&channelID,
			&userID,
			&metric.Timestamp,
			&metadataJSON,
			&metric.Version,
		)
		if err != nil {
			return nil, err
		}

		metricID, err := domain.NewMetricIDFromString(idStr)
		if err != nil {
			return nil, err
		}
		metric.ID = *metricID

		if channelID.Valid {
			parsedChannelID, err := uuid.Parse(channelID.String)
			if err != nil {
				return nil, err
			}
			metric.ChannelID = &parsedChannelID
		}

		if userID.Valid {
			parsedUserID, err := uuid.Parse(userID.String)
			if err != nil {
				return nil, err
			}
			metric.UserID = &parsedUserID
		}

		if len(metadataJSON) > 0 {
			if err := json.Unmarshal(metadataJSON, &metric.Metadata); err != nil {
				return nil, err
			}
		}

		metrics = append(metrics, &metric)
	}

	return metrics, nil
}

// Daily Metrics
func (r *PostgresAnalyticsRepository) SaveDailyMetric(ctx context.Context, metric *domain.DailyMetric) error {
	query := `
		INSERT INTO daily_metrics (id, name, date, value, count, metadata, version)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT (name, date) DO UPDATE SET
			value = daily_metrics.value + EXCLUDED.value,
			count = daily_metrics.count + EXCLUDED.count,
			version = EXCLUDED.version
	`

	metadataJSON, err := json.Marshal(metric.Metadata)
	if err != nil {
		return fmt.Errorf("failed to marshal metadata: %w", err)
	}

	_, err = r.db.Exec(ctx, query,
		metric.ID.Value(),
		metric.Name,
		metric.Date,
		metric.Value,
		metric.Count,
		metadataJSON,
		metric.Version,
	)

	return err
}

func (r *PostgresAnalyticsRepository) GetDailyMetrics(ctx context.Context, name string, startDate, endDate time.Time) ([]*domain.DailyMetric, error) {
	query := `
		SELECT id, name, date, value, count, metadata, version
		FROM daily_metrics
		WHERE name = $1 AND date >= $2 AND date <= $3
		ORDER BY date ASC
	`

	rows, err := r.db.Query(ctx, query, name, startDate, endDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var metrics []*domain.DailyMetric
	for rows.Next() {
		var metric domain.DailyMetric
		var idStr string
		var metadataJSON []byte

		err := rows.Scan(
			&idStr,
			&metric.Name,
			&metric.Date,
			&metric.Value,
			&metric.Count,
			&metadataJSON,
			&metric.Version,
		)
		if err != nil {
			return nil, err
		}

		metricID, err := domain.NewMetricIDFromString(idStr)
		if err != nil {
			return nil, err
		}
		metric.ID = *metricID

		if len(metadataJSON) > 0 {
			if err := json.Unmarshal(metadataJSON, &metric.Metadata); err != nil {
				return nil, err
			}
		}

		metrics = append(metrics, &metric)
	}

	return metrics, nil
}

// Hourly Metrics
func (r *PostgresAnalyticsRepository) SaveHourlyMetric(ctx context.Context, metric *domain.HourlyMetric) error {
	query := `
		INSERT INTO hourly_metrics (id, name, date_time, value, count, metadata, version)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT (name, date_time) DO UPDATE SET
			value = hourly_metrics.value + EXCLUDED.value,
			count = hourly_metrics.count + EXCLUDED.count,
			version = EXCLUDED.version
	`

	metadataJSON, err := json.Marshal(metric.Metadata)
	if err != nil {
		return fmt.Errorf("failed to marshal metadata: %w", err)
	}

	_, err = r.db.Exec(ctx, query,
		metric.ID.Value(),
		metric.Name,
		metric.DateTime,
		metric.Value,
		metric.Count,
		metadataJSON,
		metric.Version,
	)

	return err
}

func (r *PostgresAnalyticsRepository) GetHourlyMetrics(ctx context.Context, name string, startDate, endDate time.Time) ([]*domain.HourlyMetric, error) {
	query := `
		SELECT id, name, date_time, value, count, metadata, version
		FROM hourly_metrics
		WHERE name = $1 AND date_time >= $2 AND date_time <= $3
		ORDER BY date_time ASC
	`

	rows, err := r.db.Query(ctx, query, name, startDate, endDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var metrics []*domain.HourlyMetric
	for rows.Next() {
		var metric domain.HourlyMetric
		var idStr string
		var metadataJSON []byte

		err := rows.Scan(
			&idStr,
			&metric.Name,
			&metric.DateTime,
			&metric.Value,
			&metric.Count,
			&metadataJSON,
			&metric.Version,
		)
		if err != nil {
			return nil, err
		}

		metricID, err := domain.NewMetricIDFromString(idStr)
		if err != nil {
			return nil, err
		}
		metric.ID = *metricID

		if len(metadataJSON) > 0 {
			if err := json.Unmarshal(metadataJSON, &metric.Metadata); err != nil {
				return nil, err
			}
		}

		metrics = append(metrics, &metric)
	}

	return metrics, nil
}

// User Activity
func (r *PostgresAnalyticsRepository) SaveUserActivity(ctx context.Context, activity *domain.UserActivity) error {
	query := `
		INSERT INTO user_activities (id, user_id, last_active_at, messages_sent, channels_joined, reactions_given, session_duration, version)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		ON CONFLICT (user_id) DO UPDATE SET
			last_active_at = EXCLUDED.last_active_at,
			messages_sent = user_activities.messages_sent + EXCLUDED.messages_sent,
			channels_joined = user_activities.channels_joined + EXCLUDED.channels_joined,
			reactions_given = user_activities.reactions_given + EXCLUDED.reactions_given,
			session_duration = EXCLUDED.session_duration,
			version = EXCLUDED.version
	`

	_, err := r.db.Exec(ctx, query,
		activity.ID.Value(),
		activity.UserID,
		activity.LastActiveAt,
		activity.MessagesSent,
		activity.ChannelsJoined,
		activity.ReactionsGiven,
		activity.SessionDuration,
		activity.Version,
	)

	return err
}

func (r *PostgresAnalyticsRepository) GetUserActivity(ctx context.Context, userID uuid.UUID) (*domain.UserActivity, error) {
	query := `
		SELECT id, user_id, last_active_at, messages_sent, channels_joined, reactions_given, session_duration, version
		FROM user_activities
		WHERE user_id = $1
	`

	var activity domain.UserActivity
	var idStr string

	err := r.db.QueryRow(ctx, query, userID).Scan(
		&idStr,
		&activity.UserID,
		&activity.LastActiveAt,
		&activity.MessagesSent,
		&activity.ChannelsJoined,
		&activity.ReactionsGiven,
		&activity.SessionDuration,
		&activity.Version,
	)

	if err != nil {
		if err.Error() == "no rows in result set" {
			return nil, fmt.Errorf("user activity not found")
		}
		return nil, err
	}

	metricID, err := domain.NewMetricIDFromString(idStr)
	if err != nil {
		return nil, err
	}
	activity.ID = *metricID

	return &activity, nil
}

func (r *PostgresAnalyticsRepository) GetActiveUsers(ctx context.Context, timeRange time.Duration) (int64, error) {
	query := `
		SELECT COUNT(DISTINCT user_id)
		FROM user_activities
		WHERE last_active_at >= $1
	`

	var count int64
	err := r.db.QueryRow(ctx, query, time.Now().Add(-timeRange)).Scan(&count)
	return count, err
}

func (r *PostgresAnalyticsRepository) GetTotalUsers(ctx context.Context) (int64, error) {
	query := `SELECT COUNT(*) FROM user_activities`

	var count int64
	err := r.db.QueryRow(ctx, query).Scan(&count)
	return count, err
}

func (r *PostgresAnalyticsRepository) GetNewUsersCount(ctx context.Context, date time.Time) (int64, error) {
	query := `
		SELECT COALESCE(SUM(value), 0)
		FROM daily_metrics
		WHERE name = 'user_registrations_daily' AND date = $1
	`

	var count int64
	err := r.db.QueryRow(ctx, query, date.Truncate(24*time.Hour)).Scan(&count)
	return count, err
}

// Channel Activity
func (r *PostgresAnalyticsRepository) SaveChannelActivity(ctx context.Context, activity *domain.ChannelActivity) error {
	query := `
		INSERT INTO channel_activities (id, channel_id, messages_count, members_count, last_message_at, activity_score, version)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT (channel_id) DO UPDATE SET
			messages_count = channel_activities.messages_count + EXCLUDED.messages_count,
			members_count = EXCLUDED.members_count,
			last_message_at = EXCLUDED.last_message_at,
			activity_score = EXCLUDED.activity_score,
			version = EXCLUDED.version
	`

	_, err := r.db.Exec(ctx, query,
		activity.ID.Value(),
		activity.ChannelID,
		activity.MessagesCount,
		activity.MembersCount,
		activity.LastMessageAt,
		activity.ActivityScore,
		activity.Version,
	)

	return err
}

func (r *PostgresAnalyticsRepository) GetChannelActivity(ctx context.Context, channelID uuid.UUID) (*domain.ChannelActivity, error) {
	query := `
		SELECT id, channel_id, messages_count, members_count, last_message_at, activity_score, version
		FROM channel_activities
		WHERE channel_id = $1
	`

	var activity domain.ChannelActivity
	var idStr string

	err := r.db.QueryRow(ctx, query, channelID).Scan(
		&idStr,
		&activity.ChannelID,
		&activity.MessagesCount,
		&activity.MembersCount,
		&activity.LastMessageAt,
		&activity.ActivityScore,
		&activity.Version,
	)

	if err != nil {
		if err.Error() == "no rows in result set" {
			return nil, fmt.Errorf("channel activity not found")
		}
		return nil, err
	}

	metricID, err := domain.NewMetricIDFromString(idStr)
	if err != nil {
		return nil, err
	}
	activity.ID = *metricID

	return &activity, nil
}

func (r *PostgresAnalyticsRepository) GetActiveChannels(ctx context.Context, timeRange time.Duration) (int64, error) {
	query := `
		SELECT COUNT(*)
		FROM channel_activities
		WHERE last_message_at >= $1
	`

	var count int64
	err := r.db.QueryRow(ctx, query, time.Now().Add(-timeRange)).Scan(&count)
	return count, err
}

func (r *PostgresAnalyticsRepository) GetTotalChannels(ctx context.Context) (int64, error) {
	query := `SELECT COUNT(*) FROM channel_activities`

	var count int64
	err := r.db.QueryRow(ctx, query).Scan(&count)
	return count, err
}

// Message Analytics
func (r *PostgresAnalyticsRepository) GetMessagesCount(ctx context.Context, date time.Time) (int64, error) {
	query := `
		SELECT COALESCE(SUM(value), 0)
		FROM daily_metrics
		WHERE name = 'messages_daily' AND date = $1
	`

	var count int64
	err := r.db.QueryRow(ctx, query, date.Truncate(24*time.Hour)).Scan(&count)
	return count, err
}

func (r *PostgresAnalyticsRepository) GetTotalMessages(ctx context.Context) (int64, error) {
	query := `SELECT COALESCE(SUM(value), 0) FROM daily_metrics WHERE name = 'messages_daily'`

	var count int64
	err := r.db.QueryRow(ctx, query).Scan(&count)
	return count, err
}

func (r *PostgresAnalyticsRepository) GetPeakHour(ctx context.Context, startDate, endDate time.Time) (int, error) {
	query := `
		SELECT EXTRACT(hour FROM date_time) as hour
		FROM hourly_metrics
		WHERE name = 'messages_hourly' AND date_time >= $1 AND date_time <= $2
		GROUP BY hour
		ORDER BY SUM(value) DESC
		LIMIT 1
	`

	var hour float64
	err := r.db.QueryRow(ctx, query, startDate, endDate).Scan(&hour)
	if err != nil {
		return 0, err
	}

	return int(hour), nil
}

// Analytics Queries
func (r *PostgresAnalyticsRepository) GetUserGrowth(ctx context.Context, startDate, endDate time.Time, interval string) ([]domain.UserGrowthData, error) {
	var query string
	switch interval {
	case "daily":
		query = `
			SELECT
				date as period,
				SUM(value) as new_users,
				SUM(SUM(value)) OVER (ORDER BY date) as total_users,
				LAG(SUM(value)) OVER (ORDER BY date) as prev_new_users
			FROM daily_metrics
			WHERE name = 'user_registrations_daily' AND date >= $1 AND date <= $2
			GROUP BY date
			ORDER BY date
		`
	case "weekly":
		query = `
			SELECT
				DATE_TRUNC('week', date) as period,
				SUM(value) as new_users,
				SUM(SUM(value)) OVER (ORDER BY DATE_TRUNC('week', date)) as total_users,
				LAG(SUM(value)) OVER (ORDER BY DATE_TRUNC('week', date)) as prev_new_users
			FROM daily_metrics
			WHERE name = 'user_registrations_daily' AND date >= $1 AND date <= $2
			GROUP BY DATE_TRUNC('week', date)
			ORDER BY period
		`
	case "monthly":
		query = `
			SELECT
				DATE_TRUNC('month', date) as period,
				SUM(value) as new_users,
				SUM(SUM(value)) OVER (ORDER BY DATE_TRUNC('month', date)) as total_users,
				LAG(SUM(value)) OVER (ORDER BY DATE_TRUNC('month', date)) as prev_new_users
			FROM daily_metrics
			WHERE name = 'user_registrations_daily' AND date >= $1 AND date <= $2
			GROUP BY DATE_TRUNC('month', date)
			ORDER BY period
		`
	default:
		return nil, fmt.Errorf("unsupported interval: %s", interval)
	}

	rows, err := r.db.Query(ctx, query, startDate, endDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var growthData []domain.UserGrowthData
	for rows.Next() {
		var data domain.UserGrowthData
		var period time.Time
		var prevNewUsers sql.NullInt64

		err := rows.Scan(&period, &data.NewUsers, &data.TotalUsers, &prevNewUsers)
		if err != nil {
			return nil, err
		}

		data.Period = period.Format("2006-01-02")

		if prevNewUsers.Valid && prevNewUsers.Int64 > 0 {
			data.GrowthRate = float64(data.NewUsers-prevNewUsers.Int64) / float64(prevNewUsers.Int64) * 100
		}

		growthData = append(growthData, data)
	}

	return growthData, nil
}

func (r *PostgresAnalyticsRepository) GetMessageVolume(ctx context.Context, startDate, endDate time.Time, channelID *string) ([]domain.MessageVolumeData, error) {
	query := `
		SELECT
			date as period,
			SUM(value) as messages,
			COUNT(DISTINCT channel_id) as channels,
			AVG(metadata->>'content_length') as avg_length
		FROM analytics_metrics
		WHERE name = 'message_sent' AND timestamp >= $1 AND timestamp <= $2
	`

	args := []interface{}{startDate, endDate}
	if channelID != nil {
		query += " AND channel_id = $3"
		args = append(args, *channelID)
	}

	query += " GROUP BY date ORDER BY date"

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var volumeData []domain.MessageVolumeData
	for rows.Next() {
		var data domain.MessageVolumeData
		var period time.Time
		var avgLength sql.NullFloat64

		err := rows.Scan(&period, &data.Messages, &data.Channels, &avgLength)
		if err != nil {
			return nil, err
		}

		data.Period = period.Format("2006-01-02")
		if avgLength.Valid {
			data.AvgLength = avgLength.Float64
		}

		volumeData = append(volumeData, data)
	}

	return volumeData, nil
}

func (r *PostgresAnalyticsRepository) GetChannelActivityList(ctx context.Context, startDate, endDate time.Time, limit int) ([]domain.ChannelActivityData, error) {
	query := `
		SELECT
			channel_id,
			messages_count,
			members_count,
			last_message_at,
			activity_score
		FROM channel_activities
		WHERE last_message_at >= $1 AND last_message_at <= $2
		ORDER BY activity_score DESC, messages_count DESC
		LIMIT $3
	`

	rows, err := r.db.Query(ctx, query, startDate, endDate, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var activityData []domain.ChannelActivityData
	for rows.Next() {
		var data domain.ChannelActivityData

		err := rows.Scan(
			&data.ChannelID,
			&data.MessagesCount,
			&data.MembersCount,
			&data.LastMessageAt,
			&data.ActivityScore,
		)
		if err != nil {
			return nil, err
		}

		data.ChannelID = data.ChannelID
		activityData = append(activityData, data)
	}

	return activityData, nil
}

func (r *PostgresAnalyticsRepository) GetTopUsers(ctx context.Context, startDate, endDate time.Time, limit int) ([]domain.TopUserData, error) {
	query := `
		SELECT
			user_id,
			messages_sent,
			channels_joined,
			reactions_given,
			last_active_at
		FROM user_activities
		WHERE last_active_at >= $1 AND last_active_at <= $2
		ORDER BY messages_sent DESC, last_active_at DESC
		LIMIT $3
	`

	rows, err := r.db.Query(ctx, query, startDate, endDate, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var topUsers []domain.TopUserData
	for rows.Next() {
		var data domain.TopUserData

		err := rows.Scan(
			&data.UserID,
			&data.MessagesSent,
			&data.ChannelsJoined,
			&data.ReactionsGiven,
			&data.LastActiveAt,
		)
		if err != nil {
			return nil, err
		}

		data.UserID = data.UserID // Convert to string if needed
		topUsers = append(topUsers, data)
	}

	return topUsers, nil
}

func (r *PostgresAnalyticsRepository) GetReactionUsage(ctx context.Context, startDate, endDate time.Time) ([]domain.ReactionUsageData, error) {
	query := `
		SELECT
			metadata->>'reaction_type' as reaction_type,
			COUNT(*) as count
		FROM analytics_metrics
		WHERE name = 'reaction_added' AND timestamp >= $1 AND timestamp <= $2
		GROUP BY metadata->>'reaction_type'
		ORDER BY count DESC
	`

	rows, err := r.db.Query(ctx, query, startDate, endDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reactionData []domain.ReactionUsageData
	var totalCount int64

	// First pass to get total count
	totalQuery := `
		SELECT COUNT(*)
		FROM analytics_metrics
		WHERE name = 'reaction_added' AND timestamp >= $1 AND timestamp <= $2
	`
	err = r.db.QueryRow(ctx, totalQuery, startDate, endDate).Scan(&totalCount)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var data domain.ReactionUsageData

		err := rows.Scan(&data.ReactionType, &data.Count)
		if err != nil {
			return nil, err
		}

		if totalCount > 0 {
			data.Percentage = float64(data.Count) / float64(totalCount) * 100
		}

		reactionData = append(reactionData, data)
	}

	return reactionData, nil
}
