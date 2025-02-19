// notification.go
package main

import (
	"context"
	"fmt"
	"log"
	"time"
)

// In notification.go
func (h *ConsumerHandler) processNotification(event PaymentProcessedEvent) {
	log.Printf("▶️ Received event: %+v", event) // Debug log

	msg := generateMessage(event)
	log.Printf("ℹ️ Generated message: %s", msg)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := h.redisClient.Set(ctx, event.OrderID, msg, 24*time.Hour).Err(); err != nil {
		log.Printf("‼️ REDIS ERROR: %v", err) // Explicit error prefix
	} else {
		log.Printf("✅ Stored notification for %s", event.OrderID) // Success indicator
	}
}

func generateMessage(event PaymentProcessedEvent) string {
	if event.PaymentStatus == "SUCCESS" {
		return fmt.Sprintf("Payment successful for order %s", event.OrderID)
	}
	return fmt.Sprintf("Payment failed for order %s", event.OrderID)
}
