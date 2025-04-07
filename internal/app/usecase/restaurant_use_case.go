package usecase

import (
	"fmt"
	"time"

	"github.com/PedroNetto404/marmitech-backend/internal/app/dtos"
	"github.com/PedroNetto404/marmitech-backend/internal/domain/aggregates"
	"github.com/PedroNetto404/marmitech-backend/internal/domain/ports"
	"github.com/PedroNetto404/marmitech-backend/pkg/types"
)

const bucket = "restaurants"

type (
	IRestaurantUseCase interface {
		Create(input *dtos.RestaurantPayload) (*dtos.RestaurantDto, error)
		Update(id string, input *dtos.RestaurantPayload) (*dtos.RestaurantDto, error)
		SetImages(id string, payload *dtos.SetRestaurantImagesPayload) (*dtos.RestaurantDto, error)
		GetById(id string) (*dtos.RestaurantDto, error)
		GetAll(request types.FindArgs) (*types.PagedSlice[dtos.RestaurantDto], error)
		Delete(id string) error
	}

	restaurantUseCase struct {
		restaurantRepository ports.IRestaurantRepository
		fileStorage          ports.IBlockStorage
	}
)

func NewRestaurantUseCase(
	restaurantRepository ports.IRestaurantRepository,
	fileStorage ports.IBlockStorage,
) IRestaurantUseCase {
	return &restaurantUseCase{
		restaurantRepository: restaurantRepository,
		fileStorage:          fileStorage,
	}
}

func (r *restaurantUseCase) Create(input *dtos.RestaurantPayload) (*dtos.RestaurantDto, error) {
	exists, err := r.RestaurantExists(input.Slug, input.Cnpj)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, fmt.Errorf("restaurant with slug %s or cnpj %s already exists", input.Slug, input.Cnpj)
	}

	restaurant := aggregates.NewRestaurant(
		input.TradeName,
		input.LegalName,
		input.Cnpj,
		input.ContactPhone,
		input.WhatsAppPhone,
		input.Email,
		input.Slug,
		input.Address,
		input.Settings,
	)

	err = r.restaurantRepository.Create(restaurant)
	if err != nil {
		return nil, err
	}

	restaurantDto := dtos.MapRestauntToDto(restaurant)
	return restaurantDto, nil
}

func (r *restaurantUseCase) Update(id string, input *dtos.RestaurantPayload) (*dtos.RestaurantDto, error) {
	exists, err := r.RestaurantExists(input.Slug, input.Cnpj)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, fmt.Errorf("restaurant with slug %s or cnpj %s already exists", input.Slug, input.Cnpj)
	}

	restaurant, err := r.restaurantRepository.FindById(id)
	if err != nil {
		return nil, err
	}
	if restaurant == nil {
		return nil, fmt.Errorf("restaurant with id %s not found", id)
	}

	restaurant.TradeName = input.TradeName
	restaurant.LegalName = input.LegalName
	restaurant.CNPJ = input.Cnpj
	restaurant.ContactPhone = input.ContactPhone
	restaurant.WhatsAppPhone = input.WhatsAppPhone
	restaurant.Email = input.Email
	restaurant.Slug = input.Slug
	restaurant.Address = input.Address
	restaurant.Settings = input.Settings
	restaurant.UpdatedAt = time.Now()
	err = r.restaurantRepository.Update(restaurant)
	if err != nil {
		return nil, err
	}

	restaurantDto := dtos.MapRestauntToDto(restaurant)
	return restaurantDto, nil
}

func (r *restaurantUseCase) SetImages(id string, payload *dtos.SetRestaurantImagesPayload) (*dtos.RestaurantDto, error) {
	restaurant, err := r.restaurantRepository.FindById(id)
	if err != nil {
		return nil, err
	}

	if restaurant == nil {
		return nil, fmt.Errorf("restaurant with id %s not found", id)
	}

	if payload.Logo != nil {
		key := fmt.Sprintf("restaurants_%s_logo", restaurant.Id)

		url, err := r.fileStorage.Save(key, bucket, payload.Logo.Content)
		if err != nil {
			return nil, err
		}
		restaurant.LogoUrl = url
	}

	if payload.Banner != nil {
		key := fmt.Sprintf("restaurants_%s_banner", restaurant.Id)

		url, err := r.fileStorage.Save(key, bucket, payload.Banner.Content)
		if err != nil {
			return nil, err
		}
		restaurant.BannerUrl = url
	}

	err = r.restaurantRepository.Update(restaurant)
	if err != nil {
		return nil, err
	}

	restaurantDto := dtos.MapRestauntToDto(restaurant)
	return restaurantDto, nil
}

func (r *restaurantUseCase) GetById(id string) (*dtos.RestaurantDto, error) {
	restaurant, err := r.restaurantRepository.FindById(id)
	if err != nil {
		return nil, err
	}

	if restaurant == nil {
		return nil, fmt.Errorf("restaurant with id %s not found", id)
	}

	restaurantDto := dtos.MapRestauntToDto(restaurant)
	return restaurantDto, nil
}

func (r *restaurantUseCase) GetAll(args types.FindArgs) (*types.PagedSlice[dtos.RestaurantDto], error) {
	restaurants, err := r.restaurantRepository.Find(args)
	if err != nil {
		return nil, err
	}

	if restaurants == nil {
		return nil, fmt.Errorf("restaurants not found")
	}

	mapped := types.MapPagedSlice(*restaurants, func(restaurant aggregates.Restaurant) dtos.RestaurantDto {
		return *dtos.MapRestauntToDto(&restaurant)
	})
	return &mapped, nil
}

func (r *restaurantUseCase) Delete(id string) error {
	restaurant, err := r.restaurantRepository.FindById(id)
	if err != nil {
		return err
	}

	if restaurant == nil {
		return fmt.Errorf("restaurant with id %s not found", id)
	}

	r.restaurantRepository.Delete(id)
	return nil
}

func (r *restaurantUseCase) RestaurantExists(slug, cnpj string) (bool, error) {
	restaurants, err := r.restaurantRepository.Find(types.FindArgs{
		Filter: types.Filter{
			"slug": slug,
			"cnpj": cnpj,
		},
	})

	if err != nil {
		return false, err
	}

	if restaurants == nil {
		return false, nil
	}

	return true, nil
}
