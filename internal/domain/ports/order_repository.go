package ports

import "github.com/PedroNetto404/marmitech-backend/internal/domain/aggregates"

type IOrderRepository interface {
	IRepository[aggregates.Order]
}
