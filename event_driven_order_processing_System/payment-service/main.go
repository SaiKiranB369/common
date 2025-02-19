package main

import (
	"context"
	"log"

	"github.com/IBM/sarama"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// Database connection
	dsn := "host=localhost user=user password=password dbname=orders port=5434 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Kafka consumer
	config := sarama.NewConfig()
	config.Consumer.Offsets.Initial = sarama.OffsetOldest
	consumer, err := sarama.NewConsumerGroup(
		[]string{"localhost:9094"},
		"payment-service-group",
		config,
	)
	if err != nil {
		log.Fatal("Failed to create consumer group:", err)
	}

	processor := &PaymentProcessor{db: db}
	go func() {
		for {
			err := consumer.Consume(context.Background(), []string{"order_created"}, processor)
			if err != nil {
				log.Println("Consume error:", err)
			}
		}
	}()

	log.Println("Payment service started...")
	select {}
}
