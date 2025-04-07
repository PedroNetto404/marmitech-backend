package ports

import (
	"github.com/PedroNetto404/marmitech-backend/pkg/types"
)

type IRepository[T any] interface {
	Find(args types.FindArgs) (*types.PagedSlice[T], error)
	FindById(id string) (*T, error)
	Create(record *T) error
	Update(record *T) error
	Delete(id string) error
}
