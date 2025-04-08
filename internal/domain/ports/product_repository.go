package ports

import "github.com/PedroNetto404/marmitech-backend/internal/domain/aggregates"

type IProductRepository interface {
	IRepository[aggregates.Product]
	Exists(name, restaurantId string) (bool, error)
}
