package aggregates

import (
	"time"

	"github.com/PedroNetto404/marmitech-backend/internal/domain/abstractions"
)

type User struct {
	abstractions.AggregateRoot
	Email       string
	Active      bool
	Role        string
	Permissions []string
	PwdHash     string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   time.Time
}

var PermissionsByRole = map[string][]string{
	// Dono do negócio
	"admin": {
		"create:user",
		"read:user",
		"update:user",
		"delete:user",
		"create:restaurant",
		"read:restaurant",
		"update:restaurant",
		"delete:restaurant",
		"create:dish",
		"read:dish",
		"update:dish",
		"delete:dish",
		"create:category",
		"read:category",
		"update:category",
		"delete:category",
		"create:product",
		"read:product",
		"update:product",
		"delete:product",
	},
	// cliente do negócio, pode ser um cliente do restaurante, etc.
	"customer": {
		"read:restaurant",
		"read:dish",
		"read:category",
		"read:product",
		"read:user",
		"read:order",
		"read:payment",
		"create:order",
		"update:order",
		"create:payment",
		"update:customer",
	},
}
