package kafka

import "context"

type EventHandler interface {
	HandleEvent(ctx context.Context, event Event) error
}
