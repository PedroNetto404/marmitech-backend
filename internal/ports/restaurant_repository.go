package ports

import "github.com/PedroNetto404/marmitech-backend/internal/domain/aggregates"


type IRestaurantRepository interface {
	IRepository[*aggregates.Restaurant]
	FindByDocument(cnpj string) (*aggregates.Restaurant, error)
}
