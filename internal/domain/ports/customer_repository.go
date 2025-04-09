package ports

import "github.com/PedroNetto404/marmitech-backend/internal/domain/aggregates"

type ICustomerRepository interface {
	Create(customer *aggregates.Customer) error
}