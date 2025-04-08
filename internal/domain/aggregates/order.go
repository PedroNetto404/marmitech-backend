package aggregates

import (
	"github.com/PedroNetto404/marmitech-backend/internal/domain/abstractions"
	"github.com/PedroNetto404/marmitech-backend/pkg/types"
)

type (
	OrderItem struct {
		abstractions.Entity
		Quantity int
		Observation string
		ProductUnitPrice float64
		Discount float64
		ItemTotal float64
	}		

	OrderDelivery struct {
		abstractions.Entity

	}

	Order struct {
		abstractions.Entity
		Customer  PartialCustomer
		Restaurant PartialRestaurant
		Items     []OrderItem
		Total     float64
		TotalDiscount  float64
		DeliveryFee float64
		DeliveryAddress types.Address
	}

)