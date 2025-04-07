package ports

import (
	"github.com/PedroNetto404/marmitech-backend/internal/domain/abstractions"
	"github.com/PedroNetto404/marmitech-backend/pkg/pagination"
)

type (
	Filter map[string]any

	FindArgs struct {
		Limit   int    `json:"limit"`
		Offset  int    `json:"offset"`
		SortBy  string `json:"sort_by"`
		SortAsc bool   `json:"sort_asc"`
		Filter  Filter `json:"filter"`
	}

	IRepository[T abstractions.IAggreagateRoot] interface {
		Find(args FindArgs) (pagination.PagedSlice[T], error)
		FindById(id string) (T, error)
		Create(record T) error
		Update(record T) error
		Delete(id string) error
	}
)
