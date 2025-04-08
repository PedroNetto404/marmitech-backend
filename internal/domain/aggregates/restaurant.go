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
	PartialRestaurant struct {
		Id        string `json:"id"`
		TradeName string `json:"trade_name,omitempty"`
	}

	DeliveryConfig struct {
		Enabled            bool `json:"enabled"`
		FeePerKm           int  `json:"fee_per_km"`
		MinimumOrderValue  int  `json:"minimum_order_value"`
		MaxRadiusKm        int  `json:"max_radius_km"`
		AverageTimeMinutes int  `json:"average_time_minutes"`
	}

	EcommerceFeatures struct {
		MinimumOrderValue int  `json:"minimum_order_value"`
		Enabled           bool `json:"enabled"`
		Acquired          bool `json:"acquired"`
		AcquiredAt        *time.Time `json:"acquired_at,omitempty"`
	}

	PostPaidFeatures struct {
		Enabled                  bool `json:"enabled"`
		MinimumOrderValue        int  `json:"minimum_order_value"`
		AverageTimeMinutes       int  `json:"average_time_minutes"`
		DeliveryFeePerKm         int  `json:"delivery_fee_per_km"`
		DeliveryMaxRadiusKm      int  `json:"delivery_max_radius_km"`
		DeliveryAverageTimeMinutes int `json:"delivery_average_time_minutes"`
	}

	RestaurantSettings struct {
		ShowCnpjInReceipt      bool              `json:"show_cnpj_in_receipt"`
		Delivery               DeliveryConfig    `json:"delivery"`
		Ecommerce              EcommerceFeatures `json:"ecommerce"`
		CustomerPostPaidOrders PostPaidFeatures  `json:"customer_post_paid_orders"`
	}

	PixKey struct {
		Key string `json:"key"`
		Alias string `json:"alias"`
	}

	Payments struct {
		AcceptedPaymentMethods []string `json:"accepted_payment_methods"`
		PixKeys                []PixKey `json:"pix_keys"`
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
		DeletedAt     *time.Time         `json:"deleted_at,omitempty"`
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
