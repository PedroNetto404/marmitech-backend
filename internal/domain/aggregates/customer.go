package aggregates

import (
	"github.com/PedroNetto404/marmitech-backend/internal/domain/abstractions"
	"github.com/PedroNetto404/marmitech-backend/pkg/types"
)

type (
	PartialCustomer struct {
		Id        string
		FirstName string
		LastName  string
		Email     string
	}

	Customer struct {
		abstractions.AggregateRoot
		FirstName    string
		LastName     string
		ContactEmail string
		ContactPhone string
		Address      types.Address
	}
)
