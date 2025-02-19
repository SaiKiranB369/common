// http.go
package main

import (
	"context"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func StartHTTPServer(rdb *redis.Client) {
	router := gin.Default()

	// Add health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// Notification endpoint with proper Redis connection
	router.GET("/notifications/:order_id", func(c *gin.Context) {
		orderID := c.Param("order_id")

		// Get notification from Redis
		val, err := rdb.Get(context.Background(), orderID).Result()

		if err == redis.Nil {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Notification not found for order " + orderID,
			})
			return
		}
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to retrieve notification",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"order_id": orderID,
			"message":  val,
		})
	})

	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Failed to start HTTP server: %v", err)
	}
}
