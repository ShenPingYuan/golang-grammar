package model

import "time"

const (
	OrderStatusPending = "pending"
	OrderStatusPaid    = "paid"
)

type Order struct {
	ID        string
	UserID    string
	Product   string
	Amount    float64
	Status    string
	CreatedAt time.Time
	UpdatedAt time.Time
}