package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/IBM/sarama"
	"github.com/m1thrandir225/meridian/pkg/common"
)

type SaramaEventPublisher struct {
	producer sarama.SyncProducer
	topic    string
}

func NewSaramaEventPublisher(producer sarama.SyncProducer, defaultTopic string) *SaramaEventPublisher {
	return &SaramaEventPublisher{producer: producer, topic: defaultTopic}
}

func (p *SaramaEventPublisher) PublishEvents(ctx context.Context, domainEvents []common.DomainEvent) error {
	for _, event := range domainEvents {
		topic := p.determineTopic(event)

		payload, err := json.Marshal(event)
		if err != nil {
			log.Printf("ERROR: Failed to marshal event %s (%s): %v", event.EventName(), event.EventID(), err)
			continue
		}

		msg := &sarama.ProducerMessage{
			Topic: topic,
			Key:   sarama.StringEncoder(event.AggregateID()),
			Value: sarama.ByteEncoder(payload),
		}

		partition, offset, err := p.producer.SendMessage(msg)
		if err != nil {
			log.Printf("ERROR: Failed to publish event %s (%s) to topic %s: %v", event.EventName(), event.EventID(), topic, err)

			return fmt.Errorf("failed to push event %s: %w", event.EventName(), err)
		}
		log.Printf("Published event %s (%s) to topic %s, partition %d, offset %d", event.EventName(), event.EventID(), topic, partition, offset)
	}
	return nil
}

func (p *SaramaEventPublisher) PublishEvent(ctx context.Context, event common.DomainEvent) error {
	topic := p.determineTopic(event)

	payload, err := json.Marshal(event)
	if err != nil {
		log.Printf("ERROR: Failed to marshal event %s (%s): %v", event.EventName(), event.EventID(), err)
	}

	msg := &sarama.ProducerMessage{
		Topic: topic,
		Key:   sarama.StringEncoder(event.AggregateID()),
		Value: sarama.ByteEncoder(payload),
	}

	partition, offset, err := p.producer.SendMessage(msg)
	if err != nil {
		log.Printf("ERROR: Failed to publish event %s (%s) to topic %s: %v", event.EventName(), event.EventID(), topic, err)

		return fmt.Errorf("failed to push event %s: %w", event.EventName(), err)
	}
	log.Printf("Published event %s (%s) to topic %s, partition %d, offset %d", event.EventName(), event.EventID(), topic, partition, offset)

	return nil
}

func (p *SaramaEventPublisher) determineTopic(event common.DomainEvent) string {
	return p.topic
}

func (p *SaramaEventPublisher) Close() error {
	if p.producer != nil {
		return p.producer.Close()
	}
	return nil
}

var _ EventPublisher = (*SaramaEventPublisher)(nil)
