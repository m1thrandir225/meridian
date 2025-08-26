package kafka

import (
	"log"

	"github.com/IBM/sarama"
)

type SaramaConsumerGroupHandler struct {
	handler EventHandler
}

func (h *SaramaConsumerGroupHandler) Setup(sarama.ConsumerGroupSession) error {
	return nil
}
func (h *SaramaConsumerGroupHandler) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}
func (h *SaramaConsumerGroupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for message := range claim.Messages() {
		event := Event{
			Topic:     message.Topic,
			Partition: message.Partition,
			Offset:    message.Offset,
			Key:       string(message.Key),
			Data:      message.Value,
		}

		if err := h.handler.HandleEvent(session.Context(), event); err != nil {
			log.Printf("Error handling event: %v", err)
		}

		session.MarkMessage(message, "")
	}

	return nil

}
