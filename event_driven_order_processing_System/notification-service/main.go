// main.go
package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/IBM/sarama"
)

func main() {
	// Initialize Redis
	rdb := NewRedisClient()
	defer rdb.Close()

	// Start HTTP server in goroutine
	go StartHTTPServer(rdb)

	// Configure Kafka consumer
	config := sarama.NewConfig()
	config.Version = sarama.V2_8_0_0
	config.Consumer.Group.Rebalance.Strategy =
		sarama.NewBalanceStrategyRoundRobin()
	config.Consumer.Offsets.Initial = sarama.OffsetOldest

	// Create consumer group
	consumer, err := sarama.NewConsumerGroup(
		[]string{"localhost:9094"},
		"notification-service",
		config,
	)
	if err != nil {
		log.Fatalf("Consumer group creation failed: %v", err)
	}
	defer consumer.Close()

	// Handle shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start consuming
	handler := &ConsumerHandler{redisClient: rdb}

	go func() {
		for {
			if err := consumer.Consume(ctx, []string{"payment_processed"}, handler); err != nil {
				log.Printf("Consume error: %v", err)
				time.Sleep(5 * time.Second)
			}
		}
	}()

	// Wait for termination
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig
	log.Println("Shutting down...")
}
