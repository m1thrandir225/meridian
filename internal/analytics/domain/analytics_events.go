package domain

import (
	"time"

	"github.com/m1thrandir225/meridian/pkg/common"
)

type MetricTrackedEvent struct {
	common.BaseDomainEvent
	MetricName string                 `json:"metric_name"`
	Value      float64                `json:"value"`
	ChannelID  *string                `json:"channel_id,omitempty"`
	UserID     *string                `json:"user_id,omitempty"`
	Metadata   map[string]interface{} `json:"metadata"`
	Timestamp  time.Time              `json:"timestamp"`
}

type AnalyticsAggregatedEvent struct {
	common.BaseDomainEvent
	AggregationType string    `json:"aggregation_type"`
	Period          string    `json:"period"`
	Value           float64   `json:"value"`
	Count           int64     `json:"count"`
	Timestamp       time.Time `json:"timestamp"`
}

func CreateMetricTrackedEvent(analytics *Analytics, metric *AnalyticsMetric) MetricTrackedEvent {
	base := common.NewBaseDomainEvent("MetricTracked", analytics.ID.value, analytics.Version, "Analytics")

	var channelID, userID *string
	if metric.ChannelID != nil {
		channelIDStr := metric.ChannelID.String()
		channelID = &channelIDStr
	}
	if metric.UserID != nil {
		userIDStr := metric.UserID.String()
		userID = &userIDStr
	}

	return MetricTrackedEvent{
		BaseDomainEvent: base,
		MetricName:      metric.Name,
		Value:           metric.Value,
		ChannelID:       channelID,
		UserID:          userID,
		Metadata:        metric.Metadata,
		Timestamp:       metric.Timestamp,
	}
}

func CreateAnalyticsAggregatedEvent(analytics *Analytics, aggregationType, period string, value float64, count int64) AnalyticsAggregatedEvent {
	base := common.NewBaseDomainEvent("AnalyticsAggregated", analytics.ID.value, analytics.Version, "Analytics")

	return AnalyticsAggregatedEvent{
		BaseDomainEvent: base,
		AggregationType: aggregationType,
		Period:          period,
		Value:           value,
		Count:           count,
		Timestamp:       time.Now(),
	}
}
