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

		TradeName    string             `json:"trade_name"`
		LegalName    string             `json:"legal_name"`
		Cnpj         string             `json:"cnpj"`
		ContactPhone string             `json:"contact_phone"`
		Slug         string             `json:"slug"`
		Address      types.Address      `json:"address"`
		LogoUrl      string             `json:"logo_url"`
		BannerUrl    string             `json:"banner_url"`
		Settings     RestaurantSettings `json:"settings"`
	}
)
