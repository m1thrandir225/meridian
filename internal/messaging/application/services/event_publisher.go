package services

import (
	"context"

	"github.com/m1thrandir225/meridian/internal/messaging/domain"
)

type EventPublisher interface {
	PublishEvents(ctx context.Context, events []domain.DomainEvent) error
}
