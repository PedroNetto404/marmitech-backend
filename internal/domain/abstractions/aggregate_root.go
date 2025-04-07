package abstractions

type (
	IAggreagateRoot interface {
		IEntity

		DomainEvents() []DomainEvent
		ClearDomainEvents()
		RaiseDomainEvent(event DomainEvent)
	}

	AggregateRoot struct {
		Entity

		Events []DomainEvent
	}
)

func NewAggregateRoot() AggregateRoot {
	return AggregateRoot{
		Entity: NewEntity(),
	}
}

func (a *AggregateRoot) DomainEvents() []DomainEvent {
	return a.Events
}

func (a *AggregateRoot) ClearDomainEvents() {
	a.Events = make([]DomainEvent, 0)
}

func (a *AggregateRoot) RaiseDomainEvent(event DomainEvent) {
	a.Events = append(a.Events, event)
}
