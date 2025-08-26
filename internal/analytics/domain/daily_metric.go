package domain

import "time"

// DailyMetric represents aggregated metrics for a specific day
type DailyMetric struct {
	ID       MetricID
	Name     string
	Date     time.Time
	Value    float64
	Count    int64
	Metadata map[string]interface{}
	Version  int64
}
