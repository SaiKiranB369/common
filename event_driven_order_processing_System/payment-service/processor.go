// payment-service/processor.go
package main

import (
	"encoding/json"
	"log"

	"event_driven_order_processing_system/common/"

	"github.com/IBM/sarama"
	"gorm.io/gorm"
)

type PaymentProcessor struct {
	db       *gorm.DB
	producer sarama.SyncProducer
}

func (p *PaymentProcessor) Setup(sarama.ConsumerGroupSession) error {
	log.Println("Payment processor initialized")
	return nil
}

func (p *PaymentProcessor) Cleanup(sarama.ConsumerGroupSession) error {
	log.Println("Payment processor shutting down")
	return nil
}

func (p *PaymentProcessor) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		var event common.OrderCreatedEvent
		if err := json.Unmarshal(msg.Value, &event); err != nil {
			log.Printf("Failed to parse event: %v", err)
			continue
		}

		p.processPayment(event)
		session.MarkMessage(msg, "")
	}
	return nil
}

func (p *PaymentProcessor) processPayment(event common.OrderCreatedEvent) {
	status := "REJECTED"
	if event.TotalPrice < 1000 {
		status = "APPROVED"
	}

	paymentEvent := common.PaymentProcessedEvent{
		OrderID:       event.OrderID,
		PaymentStatus: status,
	}

	if err := p.publishPaymentEvent(paymentEvent); err != nil {
		log.Printf("Failed to publish payment event: %v", err)
		return
	}

	log.Printf("Processed payment for order %s: %s", event.OrderID, status)
}

func (p *PaymentProcessor) publishPaymentEvent(event common.PaymentProcessedEvent) error {
	jsonData, _ := json.Marshal(event)
	msg := &sarama.ProducerMessage{
		Topic: "payment_processed",
		Value: sarama.ByteEncoder(jsonData),
	}

	_, _, err := p.producer.SendMessage(msg)
	return err
}
