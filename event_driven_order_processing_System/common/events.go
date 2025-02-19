package common

type OrderCreatedEvent struct {
	OrderID    string  `json:"order_id"`
	UserID     int     `json:"user_id"`
	TotalPrice float64 `json:"total_price"`
}

type PaymentProcessedEvent struct {
	OrderID       string `json:"order_id"`
	PaymentStatus string `json:"payment_status"`
}
