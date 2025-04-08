package aggregates

import (
	"github.com/PedroNetto404/marmitech-backend/internal/domain/abstractions"
	dishtype "github.com/PedroNetto404/marmitech-backend/pkg/enums/dish_type"
)

type (
	DishTypeMap map[dishtype.DishType]int

	Product struct {
		abstractions.AggregateRoot

		Name        string `json:"name"`
		Description string `json:"description"`
		SalesPrice  string `json:"sales_price"`
		// se o produto for uma marmita, o preço de custo é da embalagem
		CostPrice   string `json:"cost_price"`
		PictureUrl string `json:"picture_url"`
		// marmita p: 2 carnes, 1 guarnição e 2 acompanhamentos
		DishTypeMap DishTypeMap `json:"dish_type_map"`
		Active      bool   `json:"active"`
		Category PartialCategory `json:"category"`
		Restaurant PartialRestaurant `json:"restaurant"`
	}
)

func NewProduct(
	name string,
	description string,
	salesPrice string,
	costPrice string,
	dishTypeMap DishTypeMap,
	categoryId string,
	restaurantId string,
) *Product {
	return &Product{
		AggregateRoot: abstractions.NewAggregateRoot(),
		Name:        name,
		Description: description,
		SalesPrice:  salesPrice,
		CostPrice:   costPrice,
		DishTypeMap: dishTypeMap,
		Category: PartialCategory{
			Id: categoryId,
		},
		Restaurant: PartialRestaurant{
			Id: restaurantId,	
		},
	}
}