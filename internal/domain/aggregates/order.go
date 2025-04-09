package aggregates

import (
	"time"

	"github.com/PedroNetto404/marmitech-backend/internal/domain/abstractions"
)

type (
	PartialProduct struct {
		Id   string
		Name string
	}

	OrderItem struct {
		ID          string
		OrderID     string
		ProductID   string
		Quantity    int
		UnitPrice   float64
		TotalPrice  float64
		Observation string
		CreatedAt   time.Time
		UpdatedAt   time.Time
		DeletedAt   *time.Time
	}

	OrderDelivery struct {
		ID                 string
		OrderID            string
		Address            string
		Fee                float64
		Distance           float64
		AverageTimeMinutes int
		Status             string
		CreatedAt          time.Time
		UpdatedAt          time.Time
		DeletedAt          *time.Time
	}

	OrderPayment struct {
		ID        string
		OrderID   string
		Amount    float64
		Method    string
		Status    string
		CreatedAt time.Time
		UpdatedAt time.Time
		DeletedAt *time.Time
	}

	Order struct {
		abstractions.AggregateRoot
		CustomerID    string
		RestaurantID  string
		Status        string
		Total         float64
		TotalDiscount float64
		Discount      float64
		Observation   string
		CreatedAt     time.Time
		UpdatedAt     time.Time
		DeletedAt     *time.Time
		Delivery      *OrderDelivery
		Items         []OrderItem
		Payments      []OrderPayment
	}
)

func (o *Order) HasBeenFullyPaidVirtual() bool {
	totalPaid := 0.0
	for _, payment := range o.Payments {
		if payment.Status == "PAID" {
			totalPaid += payment.Amount
		}
	}
	return totalPaid >= o.Total
}
