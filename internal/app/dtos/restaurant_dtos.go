package dtos

import (
	"github.com/PedroNetto404/marmitech-backend/internal/domain/aggregates"
	"github.com/PedroNetto404/marmitech-backend/pkg/types"
)

type (
	RestaurantPayload struct {
		TradeName     string                        `json:"trade_name"`
		LegalName     string                        `json:"legal_name"`
		Cnpj          string                        `json:"cnpj"`
		ContactPhone  string                        `json:"contact_phone"`
		WhatsAppPhone string                        `json:"whatsapp_phone"`
		Email         string                        `json:"email"`
		Slug          string                        `json:"slug"`
		Address       types.Address                 `json:"address"`
		Settings      aggregates.RestaurantSettings `json:"settings"`
	}

	RestaurantDto struct {
		Id            string                        `json:"id"`
		TradeName     string                        `json:"trade_name"`
		LegalName     string                        `json:"legal_name"`
		ContactPhone  string                        `json:"contact_phone"`
		WhatsAppPhone string                        `json:"whatsapp_phone"`
		Email         string                        `json:"email"`
		Slug          string                        `json:"slug"`
		Address       string                        `json:"address"`
		LogoUrl       string                        `json:"logo_url"`
		BannerUrl     string                        `json:"banner_url"`
		Settings      aggregates.RestaurantSettings `json:"settings"`
		CreatedAt     string                        `json:"created_at"`
		UpdatedAt     string                        `json:"updated_at"`
	}

	SetRestaurantImagesPayload struct {
		Logo   *types.FilePayload
		Banner *types.FilePayload
	}
)

func MapRestauntToDto(restaurant *aggregates.Restaurant) *RestaurantDto {
	return &RestaurantDto{
		Id:            restaurant.Id,
		TradeName:     restaurant.TradeName,
		LegalName:     restaurant.LegalName,
		ContactPhone:  restaurant.ContactPhone,
		WhatsAppPhone: restaurant.WhatsAppPhone,
		Email:         restaurant.Email,
		Slug:          restaurant.Slug,
		Address:       restaurant.Address.String(),
		LogoUrl:       restaurant.LogoUrl,
		BannerUrl:     restaurant.BannerUrl,
		Settings:      restaurant.Settings,
	}
}
