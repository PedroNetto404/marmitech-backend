package usecase

import (
	"errors"
	"fmt"

	"github.com/PedroNetto404/marmitech-backend/internal/domain/aggregates"
	"github.com/PedroNetto404/marmitech-backend/internal/domain/ports"
	"github.com/PedroNetto404/marmitech-backend/pkg/types"
)

const (
	productBucket = "products"
)

var (
	ErrProductNotFound      = errors.New("product not found")
	ErrCategoryNotFound     = errors.New("category not found")
	ErrRestaurantNotFound   = errors.New("restaurant not found")
	ErrProductAlreadyExists = errors.New("product already exists")
)

type (
	ProductPayload struct {
		Name        string                       `json:"name"`
		Description string                       `json:"description"`
		SalesPrice  string                       `json:"sales_price"`
		CostPrice   string                       `json:"cost_price"`
		DishTypeMap aggregates.DishTypeMap       `json:"dish_type_map"`
		Category    aggregates.PartialCategory   `json:"category"`
		Restaurant  aggregates.PartialRestaurant `json:"restaurant"`
	}

	IProductUseCase interface {
		Create(payload *ProductPayload) (*aggregates.Product, error)
		Find(args types.FindArgs) (*types.PagedSlice[aggregates.Product], error)
		FindById(id string) (*aggregates.Product, error)
		Update(id string, payload *ProductPayload) (*aggregates.Product, error)
		Delete(id string) error
		SetPicture(id string, payload *types.FilePayload) (*aggregates.Product, error)
		DeletePicture(id string) (*aggregates.Product, error)
	}

	productUseCase struct {
		productRepository    ports.IProductRepository
		categoryRepository   ports.ICategoryRepository
		restaurantRepository ports.IRestaurantRepository
		blockStorage         ports.IBlockStorage
	}
)

func NewProductUseCase(
	productRepository ports.IProductRepository,
	categoryRepository ports.ICategoryRepository,
	restaurantRepository ports.IRestaurantRepository,
	blockStorage ports.IBlockStorage,
) IProductUseCase {
	return &productUseCase{
		productRepository:  productRepository,
		categoryRepository: categoryRepository,
		blockStorage:       blockStorage,
		restaurantRepository: restaurantRepository,
	}
}

func (u *productUseCase) Create(payload *ProductPayload) (*aggregates.Product, error) {
	restaurant, err := u.restaurantRepository.FindById(payload.Restaurant.Id)
	if err != nil {
		return nil, err
	}
	if restaurant == nil {
		return nil, ErrRestaurantNotFound
	}

	category, err := u.categoryRepository.FindById(payload.Category.Id)
	if err != nil {
		return nil, err
	}
	if category == nil {
		return nil, ErrCategoryNotFound
	}

	exists, err := u.productRepository.Exists(payload.Name, restaurant.Id)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, ErrProductAlreadyExists
	}

	product := aggregates.NewProduct(
		payload.Name,
		payload.Description,
		payload.SalesPrice,
		payload.CostPrice,
		payload.DishTypeMap,
		payload.Category.Id,
		payload.Restaurant.Id,
	)
	err = u.productRepository.Create(product)
	if err != nil {
		return nil, err
	}

	return product, nil
}

func (u *productUseCase) Find(args types.FindArgs) (*types.PagedSlice[aggregates.Product], error) {
	products, err := u.productRepository.Find(args)
	if err != nil {
		return nil, err
	}

	return products, nil
}

func (u *productUseCase) FindById(id string) (*aggregates.Product, error) {
	product, err := u.productRepository.FindById(id)
	if err != nil {
		return nil, err
	}
	if product == nil {
		return nil, ErrProductNotFound
	}
	return product, nil
}

func (u *productUseCase) Update(id string, payload *ProductPayload) (*aggregates.Product, error) {
	restaurant, err := u.restaurantRepository.FindById(payload.Restaurant.Id)
	if err != nil {
		return nil, err
	}
	if restaurant == nil {
		return nil, ErrRestaurantNotFound
	}

	category, err := u.categoryRepository.FindById(payload.Category.Id)
	if err != nil {
		return nil, err
	}
	if category == nil {
		return nil, ErrCategoryNotFound
	}

	product, err := u.productRepository.FindById(id)
	if err != nil {
		return nil, err
	}
	if product == nil {
		return nil, ErrProductNotFound
	}

	product.Name = payload.Name
	product.Description = payload.Description
	product.SalesPrice = payload.SalesPrice
	product.CostPrice = payload.CostPrice
	product.DishTypeMap = payload.DishTypeMap
	product.Category.Id = payload.Category.Id
	product.Restaurant.Id = payload.Restaurant.Id

	err = u.productRepository.Update(product)
	if err != nil {
		return nil, err
	}

	return product, nil
}

func (u *productUseCase) Delete(id string) error {
	product, err := u.productRepository.FindById(id)
	if err != nil {
		return err
	}

	if product == nil {
		return ErrProductNotFound
	}

	return u.productRepository.Delete(id)
}

func (u *productUseCase) SetPicture(id string, payload *types.FilePayload) (*aggregates.Product, error) {
	product, err := u.productRepository.FindById(id)
	if err != nil {
		return nil, err
	}
	if product == nil {
		return nil, ErrProductNotFound
	}

	key := fmt.Sprintf("product_%s_picture", id)
	url, err := u.blockStorage.Save(key, productBucket, payload.Content)
	if err != nil {
		return nil, err
	}

	product.PictureUrl = url
	err = u.productRepository.Update(product)

	if err != nil {
		return nil, err
	}

	return product, nil
}

func (u *productUseCase) DeletePicture(id string) (*aggregates.Product, error) {
	product, err := u.productRepository.FindById(id)
	if err != nil {
		return nil, err
	}
	if product == nil {
		return nil, ErrProductNotFound
	}

	if product.PictureUrl == "" {
		return product, nil
	}

	key := fmt.Sprintf("product_%s_picture", id)
	err = u.blockStorage.Delete(key, productBucket)
	if err != nil {	
		return nil, err
	}
	product.PictureUrl = ""

	err = u.productRepository.Update(product)
	if err != nil {
		return nil, err
	}

	return product, nil
}
