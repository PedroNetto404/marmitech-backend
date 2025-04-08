package ports

import "github.com/PedroNetto404/marmitech-backend/internal/domain/aggregates"

type IUserRepository interface {
	IRepository[aggregates.User]
	FindByEmail(email string) (*aggregates.User, error)
}