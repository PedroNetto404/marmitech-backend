package respositories

import (
	"github.com/PedroNetto404/marmitech-backend/pkg/database"
	"github.com/PedroNetto404/marmitech-backend/internal/domain/aggregates"
	"github.com/PedroNetto404/marmitech-backend/internal/domain/ports"
	"github.com/PedroNetto404/marmitech-backend/pkg/types"
)

type productRepository struct {
	database *database.Database
}

func NewProductRepository(db *database.Database) ports.IProductRepository {
	return &productRepository{
		database: db,
	}
}

func (r *productRepository) Create(product *aggregates.Product) error {
	panic("implement me")
}

func (r *productRepository) Find(args types.FindArgs) (*types.PagedSlice[aggregates.Product], error) {
	panic("implement me")
}

func (r *productRepository) FindById(id string) (*aggregates.Product, error) {
	panic("implement me")
}

func (r *productRepository) Update(product *aggregates.Product) error {
	panic("implement me")
}

func (r *productRepository) Delete(id string) error {
	panic("implement me")
}

func (r *productRepository) Exists(name, restaurantId string) (bool, error) {
	panic("implement me")
}


