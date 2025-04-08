package aggregates

import (
	"time"
)

type Order struct {
	ID             string    `json:"id"`
	CustomerID     string    `json:"customer_id"`
	RestaurantID   string    `json:"restaurant_id"`
	Total          float64   `json:"total"`
	TotalDiscount  float64   `json:"total_discount"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	DeletedAt      *time.Time `json:"deleted_at,omitempty"`
} 