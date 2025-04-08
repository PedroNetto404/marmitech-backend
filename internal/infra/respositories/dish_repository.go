package respositories

import (
	"github.com/PedroNetto404/marmitech-backend/internal/domain/aggregates"
	"github.com/PedroNetto404/marmitech-backend/pkg/database"
	"github.com/PedroNetto404/marmitech-backend/pkg/types"
)

type dishRepository struct {
	db *database.Db
}

func NewDishRepository(db *database.Db) *dishRepository {
	return &dishRepository{
		db: db,
	}
}

func (r *dishRepository) Find(args types.FindArgs) (*types.PagedSlice[aggregates.Dish], error) {
	// Implementation of the Find method
	return nil, nil
}
func (r *dishRepository) FindById(id string) (*aggregates.Dish, error) {
	// Implementation of the FindById method
	return nil, nil
}
func (r *dishRepository) Create(dish *aggregates.Dish) error {
	// Implementation of the Create method
	return nil
}
func (r *dishRepository) Update(dish *aggregates.Dish) error {
	// Implementation of the Update method
	return nil
}
func (r *dishRepository) Delete(id string) error {
	// Implementation of the Delete method
	return nil
}

func (r *dishRepository) Exists(name string, restaurantId string) (bool, error) {
	// Implementation of the Exists method
	return false, nil
}
