package main

import (
	"log"

	"github.com/IBM/sarama"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// Database connection
	dsn := "host=localhost user=user password=password dbname=orders port=5433 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Kafka producer
	producer, err := sarama.NewSyncProducer([]string{"localhost:9093"}, nil)
	if err != nil {
		log.Fatal("Failed to create Kafka producer:", err)
	}
	defer producer.Close()

	// Setup routes
	router := gin.Default()
	router.POST("/orders", createOrderHandler(db, producer))

	log.Fatal(router.Run(":8081"))
}
