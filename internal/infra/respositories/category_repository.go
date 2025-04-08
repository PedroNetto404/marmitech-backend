package respositories

import (
	"github.com/PedroNetto404/marmitech-backend/internal/domain/aggregates"
	"github.com/PedroNetto404/marmitech-backend/pkg/database"
	"github.com/PedroNetto404/marmitech-backend/pkg/types"
)

type categoryRepository struct {
	db *database.Db
}

func NewCategoryRepository(db *database.Db) *categoryRepository {
	return &categoryRepository{
		db: db,
	}
}

func (r *categoryRepository) Create(category *aggregates.Category) error {
	panic("implement me")
}

func (r *categoryRepository) Update(category *aggregates.Category) error {
	panic("implement me")
}
func (r *categoryRepository) Delete(id string) error {
	panic("implement me")
}

func (r *categoryRepository) FindById(id string) (*aggregates.Category, error) {
	panic("implement me")
}
func (r *categoryRepository) Find(args types.FindArgs) (*types.PagedSlice[aggregates.Category], error) {
	panic("implement me")
}
func (r *categoryRepository) Exists(name string, restaurantId string) (bool, error) {
	panic("implement me")
}
func (r *categoryRepository) FindByRestaurantId(restaurantId string) ([]aggregates.Category, error) {
	panic("implement me")
}
