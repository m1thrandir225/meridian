package domain

import (
	"time"

	"github.com/google/uuid"
)

// Analytics is the aggregate root for the analytics domain
// It represents the analytics state and coordinates all analytics operations
type Analytics struct {
	ID      AnalyticsID
	Metrics []*AnalyticsMetric
	Version int64
}

func NewAnalytics() (*Analytics, error) {
	id, err := NewAnalyticsID()
	if err != nil {
		return nil, err
	}

	return &Analytics{
		ID:      *id,
		Metrics: make([]*AnalyticsMetric, 0),
		Version: 1,
	}, nil
}

// Domain methods for the Analytics aggregate
func (a *Analytics) TrackMetric(name string, value float64, channelID *uuid.UUID, userID *uuid.UUID, metadata map[string]interface{}) error {
	metricID, err := NewMetricID()
	if err != nil {
		return err
	}

	metric := &AnalyticsMetric{
		ID:        *metricID,
		Name:      name,
		Value:     value,
		ChannelID: channelID,
		UserID:    userID,
		Timestamp: time.Now(),
		Metadata:  metadata,
		Version:   1,
	}

	a.Metrics = append(a.Metrics, metric)
	a.Version++

	return nil
}

func (a *Analytics) TrackUserRegistration(userID string, timestamp time.Time) error {
	parsedUserID, err := uuid.Parse(userID)
	if err != nil {
		return err
	}

	return a.TrackMetric("user_registration", 1, nil, &parsedUserID, map[string]interface{}{
		"registration_time": timestamp,
	})
}

func (a *Analytics) TrackMessageSent(messageID, channelID, senderID string, contentLength int, timestamp time.Time) error {
	parsedChannelID, err := uuid.Parse(channelID)
	if err != nil {
		return err
	}

	parsedSenderID, err := uuid.Parse(senderID)
	if err != nil {
		return err
	}

	return a.TrackMetric("message_sent", 1, &parsedChannelID, &parsedSenderID, map[string]interface{}{
		"message_id":     messageID,
		"content_length": contentLength,
		"timestamp":      timestamp,
	})
}

func (a *Analytics) TrackChannelCreated(channelID, creatorID string, timestamp time.Time) error {
	parsedChannelID, err := uuid.Parse(channelID)
	if err != nil {
		return err
	}

	parsedCreatorID, err := uuid.Parse(creatorID)
	if err != nil {
		return err
	}

	return a.TrackMetric("channel_created", 1, &parsedChannelID, &parsedCreatorID, map[string]interface{}{
		"creation_time": timestamp,
	})
}
