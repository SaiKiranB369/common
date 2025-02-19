package main

import (
	"net/http"

	"github.com/IBM/sarama"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type OrderRequest struct {
	UserID     int     `json:"user_id"`
	ProductID  int     `json:"product_id"`
	Quantity   int     `json:"quantity"`
	TotalPrice float64 `json:"total_price"`
}

func createOrderHandler(db *gorm.DB, producer sarama.SyncProducer) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req OrderRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Save to database
		order := Order{
			UserID:     req.UserID,
			ProductID:  req.ProductID,
			Quantity:   req.Quantity,
			TotalPrice: req.TotalPrice,
			Status:     "PENDING",
		}
		db.Create(&order)

		// Publish Kafka event
		event := OrderCreatedEvent{
			OrderID:    order.ID.String(),
			UserID:     order.UserID,
			TotalPrice: order.TotalPrice,
		}
		publishOrderCreatedEvent(producer, event)

		c.JSON(http.StatusCreated, gin.H{
			"order_id": order.ID,
			"status":   order.Status,
		})
	}
}
