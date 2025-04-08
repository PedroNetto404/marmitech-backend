package aggregates

import "github.com/PedroNetto404/marmitech-backend/internal/domain/abstractions"

type PartialCategory struct {
	Id string `json:"id"`
	Name string `json:"name"`
	PictureUrl string `json:"picture_url"`
}

type Category struct {
	abstractions.AggregateRoot
	Name       string            `json:"name"`
	PictureUrl string            `json:"picture_url"`
	Restaurant PartialRestaurant `json:"restaurant"`
	Priority   int               `json:"priority"`
	Active     bool              `json:"active"`
}

func Newcategory(
	name string,
	restaurantId string,
	priority int,
) *Category {
	return &Category{
		AggregateRoot: abstractions.NewAggregateRoot(),
		Name:          name,
		Restaurant: PartialRestaurant{
			Id: restaurantId,
		},
		Priority: priority,
		Active:   true,
	}
}
