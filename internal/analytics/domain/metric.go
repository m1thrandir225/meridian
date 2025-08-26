package domain

import (
	"time"

	"github.com/google/uuid"
)

// AnalyticsMetric represents a single metric measurement
type AnalyticsMetric struct {
	ID        MetricID
	Name      string
	Value     float64
	ChannelID *uuid.UUID
	UserID    *uuid.UUID
	Timestamp time.Time
	Metadata  map[string]interface{}
	Version   int64
}
