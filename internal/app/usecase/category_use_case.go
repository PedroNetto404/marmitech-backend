package usecase

import (
	"errors"
	"fmt"

	"github.com/PedroNetto404/marmitech-backend/internal/domain/aggregates"
	"github.com/PedroNetto404/marmitech-backend/internal/domain/ports"
	"github.com/PedroNetto404/marmitech-backend/pkg/types"
)

const (
	categoryBucket = "product_categories"
)

var (
	ErrcategoryNotFound      = errors.New("product category not found")
	ErrcategoryAlreadyExists = errors.New("product category already exists")
)

type (
	CategoryPayload struct {
		Name       string                       `json:"name"`
		Restaurant aggregates.PartialRestaurant `json:"restaurant"`
		Priority   int                          `json:"priority"`
	}

	IcategoryUseCase interface {
		Find(args types.FindArgs) (*types.PagedSlice[aggregates.Category], error)
		FindById(id string) (*aggregates.Category, error)
		Create(category *CategoryPayload) (*aggregates.Category, error)
		Update(id string, category *CategoryPayload) (*aggregates.Category, error)
		Delete(id string) error
		SetPicture(id string, picture *types.FilePayload) (*aggregates.Category, error)
		DeletePicture(id string) (*aggregates.Category, error)
		Activate(id string) (*aggregates.Category, error)
		Deactivate(id string) (*aggregates.Category, error)
		Reorder(id string, priority int, destinationId string) (*aggregates.Category, error)
	}

	categoryUseCase struct {
		categoryRepository   ports.ICategoryRepository
		restaurantRepository ports.IRestaurantRepository
		blockStorage         ports.IBlockStorage
	}
)

func NewCategoryUseCase(
	categoryRepository ports.ICategoryRepository,
	restaurantRepository ports.IRestaurantRepository,
	blockStorage ports.IBlockStorage,
) IcategoryUseCase {
	return &categoryUseCase{
		categoryRepository:   categoryRepository,
		restaurantRepository: restaurantRepository,
		blockStorage:         blockStorage,
	}
}

func (p *categoryUseCase) Find(args types.FindArgs) (*types.PagedSlice[aggregates.Category], error) {
	productCategories, err := p.categoryRepository.Find(args)
	if err != nil {
		return nil, err
	}

	return productCategories, nil
}

func (p *categoryUseCase) FindById(id string) (*aggregates.Category, error) {
	category, err := p.categoryRepository.FindById(id)
	if err != nil {
		return nil, err
	}

	if category == nil {
		return nil, ErrcategoryNotFound
	}

	return category, nil
}

func (p *categoryUseCase) Create(payload *CategoryPayload) (*aggregates.Category, error) {
	exists, err := p.categoryRepository.Exists(payload.Restaurant.Id, payload.Name)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, ErrcategoryAlreadyExists
	}

	category := aggregates.Newcategory(
		payload.Restaurant.Id,
		payload.Name,
		payload.Priority,
	)

	err = p.categoryRepository.Create(category)
	if err != nil {
		return nil, err
	}

	return category, nil
}

func (p *categoryUseCase) Update(id string, payload *CategoryPayload) (*aggregates.Category, error) {
	exists, err := p.categoryRepository.Exists(payload.Restaurant.Id, payload.Name)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, ErrcategoryAlreadyExists
	}

	category, err := p.categoryRepository.FindById(id)
	if err != nil {
		return nil, err
	}

	if category == nil {
		return nil, ErrcategoryNotFound
	}

	category.Name = payload.Name
	category.Priority = payload.Priority

	err = p.categoryRepository.Update(category)
	if err != nil {
		return nil, err
	}

	return category, nil
}

func (p *categoryUseCase) Delete(id string) error {
	category, err := p.categoryRepository.FindById(id)
	if err != nil {
		return err
	}

	if category == nil {
		return ErrcategoryNotFound
	}
	err = p.categoryRepository.Delete(id)
	if err != nil {
		return err
	}
	return nil
}

func (p *categoryUseCase) SetPicture(id string, picture *types.FilePayload) (*aggregates.Category, error) {
	category, err := p.categoryRepository.FindById(id)
	if err != nil {
		return nil, err
	}
	if category == nil {
		return nil, ErrcategoryNotFound
	}

	key := fmt.Sprintf("product_category_picture_%s", category.Id)
	url, err := p.blockStorage.Save(key, categoryBucket, picture.Content)
	if err != nil {
		return nil, err
	}

	category.PictureUrl = url
	err = p.categoryRepository.Update(category)
	if err != nil {
		return nil, err
	}
	return category, nil
}

func (p *categoryUseCase) DeletePicture(id string) (*aggregates.Category, error) {
	category, err := p.categoryRepository.FindById(id)
	if err != nil {
		return nil, err
	}

	if category == nil {
		return nil, ErrcategoryNotFound
	}

	if category.PictureUrl == "" {
		return nil, nil
	}

	key := fmt.Sprintf("product_category_picture_%s", category.Id)
	err = p.blockStorage.Delete(key, categoryBucket)
	if err != nil {
		return nil, err
	}
	category.PictureUrl = ""
	err = p.categoryRepository.Update(category)
	if err != nil {
		return nil, err
	}
	return category, nil
}

func (p *categoryUseCase) Activate(id string) (*aggregates.Category, error) {
	return p.setActiveStatus(id, true)
}

func (p *categoryUseCase) Deactivate(id string) (*aggregates.Category, error) {
	return p.setActiveStatus(id, false)
}

func (p *categoryUseCase) Reorder(id string, priority int, swapId string) (*aggregates.Category, error) {
	category, err := p.categoryRepository.FindById(id)
	if err != nil {
		return nil, err
	}
	if category == nil {
		return nil, ErrcategoryNotFound
	}

	if swapId == "" {
		return p.swapPriorities(category, swapId)
	} else {
		category.Priority = priority
	}

	err = p.categoryRepository.Update(category)
	if err != nil {
		return nil, err
	}

	return category, nil
}

func (p *categoryUseCase) setActiveStatus(id string, status bool) (*aggregates.Category, error) {
	category, err := p.categoryRepository.FindById(id)
	if err != nil {
		return nil, err
	}
	if category == nil {
		return nil, ErrcategoryNotFound
	}

	category.Active = status

	err = p.categoryRepository.Update(category)
	if err != nil {
		return nil, err
	}
	return category, nil
}

func (p *categoryUseCase) swapPriorities(category *aggregates.Category, swapId string) (*aggregates.Category, error) {
	swapcategory, err := p.categoryRepository.FindById(fmt.Sprintf("%d", swapId))
	if err != nil {
		return nil, err
	}
	if swapcategory == nil {
		return nil, ErrcategoryNotFound
	}

	category.Priority, swapcategory.Priority = swapcategory.Priority, category.Priority

	err = p.categoryRepository.Update(category)
	if err != nil {
		return nil, err
	}

	err = p.categoryRepository.Update(swapcategory)
	if err != nil {
		return nil, err
	}

	return category, nil
}
