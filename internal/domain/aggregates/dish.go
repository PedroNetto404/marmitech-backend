package aggregates

import (
	"github.com/PedroNetto404/marmitech-backend/internal/domain/abstractions"
	dishtype "github.com/PedroNetto404/marmitech-backend/pkg/enums/dish_type"
)

type Dish struct {
	abstractions.AggregateRoot
	Restaurant           PartialRestaurant `json:"restaurant"`
	Name                 string            `json:"name"`
	Type                 dishtype.DishType `json:"type"`
	PictureUrl           string            `json:"picture_url"`
	PriceWhenUsedAsAddOn float64           `json:"price_when_used_as_add_on"`
}

func NewDish(
	restaurantId string,
	name string,
	dishType dishtype.DishType,
	priceWhenUsedAsAddOn float64,
) *Dish {
	return &Dish{
		AggregateRoot: abstractions.NewAggregateRoot(),
		Restaurant: PartialRestaurant{
			Id:        restaurantId,
			TradeName: "",
		},
		Name:                 name,
		Type:                 dishType,
		PriceWhenUsedAsAddOn: priceWhenUsedAsAddOn,
	}
}
