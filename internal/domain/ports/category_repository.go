package ports

import "github.com/PedroNetto404/marmitech-backend/internal/domain/aggregates"

type ICategoryRepository interface {
	IRepository[aggregates.Category]
	Exists(
		restaurantId string,
		name string,
	) (bool, error)
	FindByRestaurantId(restaurantId string) ([]aggregates.Category, error)
}
