package common

import (
	"time"

	"gorm.io/gorm"
)

type Order struct {
	ID         string `gorm:"primaryKey"`
	UserID     int    `gorm:"index"`
	ProductID  int
	Quantity   int
	TotalPrice float64
	Status     string
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  gorm.DeletedAt `gorm:"index"`
}
