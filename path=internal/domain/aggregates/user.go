package aggregates

import (
	"time"

	"github.com/PedroNetto404/marmitech-backend/internal/domain/abstractions"
)

type User struct {
	abstractions.AggregateRoot
	ID          string    `json:"id"`
	Email       string    `json:"email"`
	PwdHash     string    `json:"pwd_hash"`
	Active      bool      `json:"active"`
	Role        string    `json:"role"`
	Permissions []string  `json:"permissions"` // Assuming permissions are stored as a JSON array of strings
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty"`
}

var PermissionsByRole = map[string][]string{
	// Roles and permissions mapping
	"admin": {
		// permissions for admin
	},
	"customer": {
		// permissions for customer
	},
} 