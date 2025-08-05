package domain

import "github.com/google/uuid"

type IntegrationID struct {
	value uuid.UUID
}

func NewIntegrationID() (*IntegrationID, error) {
	id, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}
	return &IntegrationID{
		value: id,
	}, nil
}

func (id *IntegrationID) String() string {
	return id.value.String()
}

func (id *IntegrationID) Value() uuid.UUID {
	return id.value
}

func NewIntegrationIDFromString(input string) (*IntegrationID, error) {
	id, err := uuid.Parse(input)
	if err != nil {
		return nil, err
	}
	return &IntegrationID{
		value: id,
	}, nil
}
