package abstractions

import "github.com/google/uuid"

type (
	IEntity interface {
		GetId() string
	}

	Entity struct {
		Id string
	}
)

func NewEntity() Entity {
	return Entity{
		Id: uuid.NewString(),
	}
}

func (e *Entity) GetId() string {
	return e.Id
}
