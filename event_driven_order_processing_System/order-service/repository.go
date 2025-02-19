// order-service/repository.go
package main

import (
	"encoding/json"
	"fmt"

	"event_driven_order_processing_system/common/"

	"github.com/IBM/sarama"
	"gorm.io/gorm"
)

type OrderRepository struct {
	db       *gorm.DB
	producer sarama.SyncProducer
}

func NewOrderRepository(db *gorm.DB, producer sarama.SyncProducer) *OrderRepository {
	return &OrderRepository{db: db, producer: producer}
}

func (r *OrderRepository) CreateOrder(order *common.Order) error {
	if err := r.db.Create(order).Error; err != nil {
		return fmt.Errorf("database create failed: %w", err)
	}

	event := common.OrderCreatedEvent{
		OrderID:    order.ID.String(),
		UserID:     order.UserID,
		TotalPrice: order.TotalPrice,
	}

	return r.publishOrderCreated(event)
}

func (r *OrderRepository) publishOrderCreated(event common.OrderCreatedEvent) error {
	jsonData, _ := json.Marshal(event)
	msg := &sarama.ProducerMessage{
		Topic: "order_created",
		Value: sarama.ByteEncoder(jsonData),
	}

	if _, _, err := r.producer.SendMessage(msg); err != nil {
		return fmt.Errorf("failed to publish event: %w", err)
	}
	return nil
}
