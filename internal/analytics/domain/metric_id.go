package domain

import "github.com/google/uuid"

type MetricID struct {
	value uuid.UUID
}

func NewMetricID() (*MetricID, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}
	return &MetricID{value: id}, nil
}

func NewMetricIDFromString(id string) (*MetricID, error) {
	parsedID, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}
	return &MetricID{value: parsedID}, nil
}

func (m MetricID) String() string {
	return m.value.String()
}

func (m MetricID) Value() uuid.UUID {
	return m.value
}
