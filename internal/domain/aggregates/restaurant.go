package aggregates

import (
	"time"

	"github.com/PedroNetto404/marmitech-backend/internal/domain/abstractions"
	"github.com/PedroNetto404/marmitech-backend/pkg/types"
)

const (
	RestaurantCreatedEvent abstractions.EventName = "restaurant.created"
	RestaurantUpdatedEvent abstractions.EventName = "restaurant.updated"
	RestaurantDeletedEvent abstractions.EventName = "restaurant.deleted"
)

type (
	DeliveryConfig struct {
		Enabled            bool `json:"enabled"`
		FeePerKm           int  `json:"fee_per_km"`
		MinimumOrderValue  int  `json:"minimum_order_value"`
		MaxRadiusDelivery  int  `json:"max_radius_delivery"`
		AverageTimeMinutes int  `json:"average_time_minutes"`
	}

	EcommerceFeatures struct {
		MinimumOrderValue int        `json:"minimum_order_value"`
		Acquired          bool       `json:"acquired"`
		AcquiredAt        *time.Time `json:"acquired_at,omitempty"`
		Online            bool       `json:"online"`
	}

	PostPaidFeatures struct {
		Acquired   bool       `json:"acquired"`
		AcquiredAt *time.Time `json:"acquired_at,omitempty"`
		Enabled    bool       `json:"enabled"`
	}

	RestaurantSettings struct {
		ShowCnpjInReceipt bool                 `json:"show_cnpj_in_receipt"`
		WeeklySchedule    types.WeeklySchedule `json:"weekly_schedule"`
		Delivery          DeliveryConfig       `json:"delivery"`
		Ecommerce         EcommerceFeatures    `json:"ecommerce"`
		PostPaid          PostPaidFeatures     `json:"post_paid"`
	}

	Restaurant struct {
		abstractions.AggregateRoot

		TradeName     string             `json:"trade_name"`
		LegalName     string             `json:"legal_name"`
		CNPJ          string             `json:"cnpj"`
		ContactPhone  string             `json:"contact_phone"`
		WhatsAppPhone string             `json:"whatsapp_phone"`
		Email         string             `json:"email"`
		Slug          string             `json:"slug"`
		Address       types.Address      `json:"address"`
		LogoUrl       string             `json:"logo_url"`
		BannerUrl     string             `json:"banner_url"`
		Settings      RestaurantSettings `json:"settings"`
		CreatedAt     time.Time          `json:"created_at"`
		UpdatedAt     time.Time          `json:"updated_at"`
		Active        bool               `json:"active"`
	}
)

func NewRestaurant(
	tradeName,
	legalName,
	cnpj,
	contactPhone,
	whatsappPhone,
	email,
	slug string,
	address types.Address,
	settings RestaurantSettings,
) *Restaurant {
	return &Restaurant{
		AggregateRoot: abstractions.NewAggregateRoot(),
		TradeName:     tradeName,
		LegalName:     legalName,
		CNPJ:          cnpj,
		ContactPhone:  contactPhone,
		WhatsAppPhone: whatsappPhone,
		Email:         email,
		Slug:          slug,
		Address:       address,
		Settings:      settings,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
		Active:        true,
	}
}
