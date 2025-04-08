package aggregates

import (
	"github.com/PedroNetto404/marmitech-backend/internal/domain/abstractions"
	dishtype "github.com/PedroNetto404/marmitech-backend/pkg/enums/dish_type"
)

type PartialDish struct {
	Id string `json:"id"`
	Name string `json:"name"`
	Type dishtype.DishType `json:"type"`
}

type Dish struct {
	abstractions.AggregateRoot
	Restaurant           PartialRestaurant `json:"restaurant"`
	Name                 string            `json:"name"`
	Type                 dishtype.DishType `json:"type"`
	PictureUrl           string            `json:"picture_url"`
}

func NewDish(
	restaurantId string,
	name string,
	dishType dishtype.DishType,
) *Dish {
	return &Dish{
		AggregateRoot: abstractions.NewAggregateRoot(),
		Restaurant: PartialRestaurant{
			Id:        restaurantId,
			TradeName: "",
		},
		Name:                 name,
		Type:                 dishType,
	}
}
