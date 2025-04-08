package aggregates

import (
	"github.com/PedroNetto404/marmitech-backend/internal/domain/aggregates"
)

type Product struct {
	ID           string             `json:"id"`
	Name         string             `json:"name"`
	Description  string             `json:"description,omitempty"`
	SalesPrice   float64            `json:"sales_price"`
	CostPrice    float64            `json:"cost_price"`
	PictureUrl   string             `json:"picture_url,omitempty"`
	DishTypeMap  map[string]int     `json:"dish_type_map"` // Assuming dish_type_map is a JSON object with string keys and integer values
	Active       bool               `json:"active"`
	CategoryID   string             `json:"category_id"`
	RestaurantID string             `json:"restaurant_id"`
} 