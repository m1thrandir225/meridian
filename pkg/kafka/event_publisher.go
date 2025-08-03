package kafka

import (
	"context"
	"github.com/m1thrandir225/meridian/pkg/common"
)

type EventPublisher interface {
	PublishEvents(ctx context.Context, events []common.DomainEvent) error
	PublishEvent(ctx context.Context, event common.DomainEvent) error
}
