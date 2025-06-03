package services

import (
	"context"

	"github.com/m1thrandir225/meridian/pkg/common"
)

type EventPublisher interface {
	PublishEvents(ctx context.Context, events []common.DomainEvent) error
}
