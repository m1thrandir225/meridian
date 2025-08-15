package domain

import (
	"time"

	"github.com/google/uuid"
)

type IntegrationBot struct {
	id          uuid.UUID
	serviceName string
	createdAt   time.Time
	isRevoked   bool
}

func NewIntegrationBot(id uuid.UUID, serviceName string, createdAt time.Time, isRevoked bool) *IntegrationBot {
	return &IntegrationBot{
		id:          id,
		serviceName: serviceName,
		createdAt:   createdAt,
		isRevoked:   isRevoked,
	}
}

func (i *IntegrationBot) GetId() uuid.UUID {
	return i.id
}

func (i *IntegrationBot) GetServiceName() string {
	return i.serviceName
}

func (i *IntegrationBot) GetCreatedAt() time.Time {
	return i.createdAt
}

func (i *IntegrationBot) GetIsRevoked() bool {
	return i.isRevoked
}
