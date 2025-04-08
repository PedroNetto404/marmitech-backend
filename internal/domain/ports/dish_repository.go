package ports

import "github.com/PedroNetto404/marmitech-backend/internal/domain/aggregates"

type IDishRepository interface {
	IRepository[aggregates.Dish]
}