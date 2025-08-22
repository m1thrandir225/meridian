package common

import (
	"time"

	"github.com/google/uuid"
)

// Event that has occured in the domain
type DomainEvent interface {
	EventID() string
	EventName() string
	EventTime() time.Time
	AggregateID() string
	AggregateType() string
	AggregateVersion() int64
}

// Common fields for all domain events
type BaseDomainEvent struct {
	ID          string
	Name        string
	Time        time.Time
	AggrID      string
	AggrType    string
	AggrVersion int64
}

func (e BaseDomainEvent) EventID() string {
	return e.ID
}

func (e BaseDomainEvent) EventName() string {
	return e.Name
}

func (e BaseDomainEvent) EventTime() time.Time {
	return e.Time
}

func (e BaseDomainEvent) AggregateID() string {
	return e.AggrID
}

func (e BaseDomainEvent) AggregateType() string {
	return e.AggrType
}

func (e BaseDomainEvent) AggregateVersion() int64 {
	return e.AggrVersion
}

func NewBaseDomainEvent(eventName string, aggregateID uuid.UUID, version int64, aggregateType string) BaseDomainEvent {
	event, err := uuid.NewV7()
	if err != nil {
		panic("error while creating an id")
	}
	return BaseDomainEvent{
		ID:          event.String(),
		Name:        eventName,
		Time:        time.Now().UTC(),
		AggrID:      aggregateID.String(),
		AggrType:    aggregateType,
		AggrVersion: version,
	}
}
