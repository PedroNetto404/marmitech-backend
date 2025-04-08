package usecase

import (
	"errors"
	"fmt"

	"github.com/PedroNetto404/marmitech-backend/internal/domain/aggregates"
	"github.com/PedroNetto404/marmitech-backend/internal/domain/ports"
	dishtype "github.com/PedroNetto404/marmitech-backend/pkg/enums/dish_type"
	"github.com/PedroNetto404/marmitech-backend/pkg/types"
)

const (
	dishBucket = "dishes"
)

var (
	ErrDishNotFound      = errors.New("dish not found")
	ErrDishAlreadyExists = errors.New("dish already exists")
)

type (
	DishPayload struct {
		Name                 string                   `json:"name"`
		Type                 dishtype.DishType        `json:"type"`
		PriceWhenUsedAsAddOn float64                  `json:"price_when_used_as_add_on"`
		Restaurant           aggregates.PartialRestaurant `json:"restaurant"`
	}

	IDishUseCase interface {
		Find(args types.FindArgs) (*types.PagedSlice[aggregates.Dish], error)
		FindById(id string) (*aggregates.Dish, error)
		Create(dish *DishPayload) (*aggregates.Dish, error)
		Update(id string, dish *DishPayload) (*aggregates.Dish, error)
		Delete(id string) error
		SetPicture(id string, picture *types.FilePayload) (*aggregates.Dish, error)
		DeletePicture(id string) (*aggregates.Dish, error)
	}

	dishUseCase struct {
		dishRepository       ports.IDishRepository
		restaurantRepository ports.IRestaurantRepository
		blockStorage         ports.IBlockStorage
	}
)

func NewDishUseCase(
	dishRepository ports.IDishRepository,
	restaurantRepository ports.IRestaurantRepository,
	blockStorage ports.IBlockStorage,
) IDishUseCase {
	return &dishUseCase{
		dishRepository:       dishRepository,
		restaurantRepository: restaurantRepository,
		blockStorage:         blockStorage,
	}
}

func (d *dishUseCase) Find(args types.FindArgs) (*types.PagedSlice[aggregates.Dish], error) {
	dishes, err := d.dishRepository.Find(args)
	if err != nil {
		return nil, err
	}

	return dishes, nil
}

func (d *dishUseCase) FindById(id string) (*aggregates.Dish, error) {
	dish, err := d.dishRepository.FindById(id)
	if err != nil {
		return nil, err
	}

	if dish == nil {
		return nil, ErrDishNotFound
	}

	return dish, nil
}

func (d *dishUseCase) Create(dishPayload *DishPayload) (*aggregates.Dish, error) {
	exists, err := d.Exists(dishPayload.Name, dishPayload.Restaurant.Id)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, ErrDishAlreadyExists
	}

	dish := aggregates.NewDish(
		dishPayload.Restaurant.Id,
		dishPayload.Name,
		dishPayload.Type,
	)

	err = d.dishRepository.Create(dish)
	if err != nil {
		return nil, err
		}

	return dish, nil
}

func (d *dishUseCase) Update(id string, dishPayload *DishPayload) (*aggregates.Dish, error) {
	exists, err := d.Exists(dishPayload.Name, dishPayload.Restaurant.Id)
	if err != nil {
		return nil, err
	}

	if exists {
		return nil, ErrDishAlreadyExists
	}

	dish, err := d.dishRepository.FindById(id)
	if err != nil {
		return nil, err
	}

	if dish == nil {
		return nil, ErrDishNotFound
	}

	dish.Name = dishPayload.Name
	dish.Type = dishPayload.Type
	dish.Restaurant.Id = dishPayload.Restaurant.Id

	err = d.dishRepository.Update(dish)
	if err != nil {
		return nil, err
	}

	return dish, nil
}

func (d *dishUseCase) Delete(id string) error {
	dish, err := d.dishRepository.FindById(id)
	if err != nil {
		return err
	}

	if dish == nil {
		return ErrDishNotFound
	}

	err = d.dishRepository.Delete(id)
	if err != nil {
		return err
	}

	return nil
}

func (d *dishUseCase) SetPicture(id string, picture *types.FilePayload) (*aggregates.Dish, error) {
	dish, err := d.dishRepository.FindById(id)
	if err != nil {
		return nil, err
	}

	if dish == nil {
		return nil, ErrDishNotFound
	}

	key := fmt.Sprintf("dish_picture_%s", dish.Id)
	url, err := d.blockStorage.Save(key, dishBucket, picture.Content)
	if err != nil {
		return nil, err
	}

	dish.PictureUrl = url

	err = d.dishRepository.Update(dish)
	if err != nil {
		return nil, err
	}

	return dish, nil
}

func (d *dishUseCase) Exists(name, restaurantId string) (bool, error) {
	dish, err := d.dishRepository.Find(types.FindArgs{
		Filter: types.Filter{
			"name":          name,
			"restaurant_id": restaurantId,
		},
	})
	if err != nil {
		return false, err
	}

	if dish != nil && len(dish.Records) > 0 {
		return true, nil
	}

	return false, nil
}

func (d *dishUseCase) DeletePicture(id string) (*aggregates.Dish, error) {
	dish, err := d.dishRepository.FindById(id)
	if err != nil {
		return nil, err
	}

	if dish == nil {
		return nil, ErrDishNotFound
	}

	key := fmt.Sprintf("dish_picture_%s", dish.Id)
	err = d.blockStorage.Delete(key, dishBucket)
	if err != nil {
		return nil, err
	}

	dish.PictureUrl = ""

	err = d.dishRepository.Update(dish)
	if err != nil {
		return nil, err
	}

	return dish, nil
}