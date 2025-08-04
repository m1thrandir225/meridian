package kafka

import (
	"context"
	"github.com/IBM/sarama"
	"log"
)

type EventConsumer struct {
	brokers       []string
	consumerGroup string
	config        *sarama.Config
}

func NewEventConsumer(brokers []string, consumerGroup string, config *sarama.Config) *EventConsumer {
	return &EventConsumer{
		brokers:       brokers,
		consumerGroup: consumerGroup,
		config:        config,
	}
}

func (c *EventConsumer) ConsumeEvents(ctx context.Context, topics []string, handler EventHandler) error {
	consumerGroup, err := sarama.NewConsumerGroup(c.brokers, c.consumerGroup, c.config)
	if err != nil {
		return err
	}
	defer consumerGroup.Close()

	consumer := &SaramaConsumerGroupHandler{handler: handler}
	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			if err := consumerGroup.Consume(ctx, topics, consumer); err != nil {
				log.Printf("Error consuming: %v", err)
				return err
			}
		}
	}

}
