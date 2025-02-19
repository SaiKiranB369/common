// consumer.go
package main

import (
	"encoding/json"
	"log"

	"github.com/IBM/sarama"
	"github.com/redis/go-redis/v9"
)

type PaymentProcessedEvent struct {
	OrderID       string `json:"order_id"`
	PaymentStatus string `json:"payment_status"`
}

type ConsumerHandler struct {
	redisClient *redis.Client
}

func (h *ConsumerHandler) Setup(sarama.ConsumerGroupSession) error {
	log.Println("Initializing Kafka consumer")
	return nil
}

func (h *ConsumerHandler) Cleanup(sarama.ConsumerGroupSession) error {
	log.Println("Closing Kafka consumer")
	return nil
}

func (h *ConsumerHandler) ConsumeClaim(session sarama.ConsumerGroupSession,
	claim sarama.ConsumerGroupClaim) error {
	for message := range claim.Messages() {
		var event PaymentProcessedEvent
		if err := json.Unmarshal(message.Value, &event); err != nil {
			log.Printf("Failed to parse message: %v", err)
			continue
		}

		// Process notification (explained next section)
		h.processNotification(event)
		session.MarkMessage(message, "")
	}
	return nil
}
