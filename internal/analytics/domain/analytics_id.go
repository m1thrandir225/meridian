package domain

import "github.com/google/uuid"

type AnalyticsID struct {
	value uuid.UUID
}

func NewAnalyticsID() (*AnalyticsID, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}
	return &AnalyticsID{value: id}, nil
}

func NewAnalyticsIDFromString(id string) (*AnalyticsID, error) {
	parsedID, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}
	return &AnalyticsID{value: parsedID}, nil
}

func (a AnalyticsID) String() string {
	return a.value.String()
}

func (a AnalyticsID) Value() uuid.UUID {
	return a.value
}
