package domain

import "time"

// HourlyMetric represents aggregated metrics for a specific hour
type HourlyMetric struct {
	ID       MetricID
	Name     string
	DateTime time.Time
	Value    float64
	Count    int64
	Metadata map[string]interface{}
	Version  int64
}
