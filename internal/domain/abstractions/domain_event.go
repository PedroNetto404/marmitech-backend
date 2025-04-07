package abstractions

import "time"

type EventName string

type DomainEvent struct {
	Name        EventName `json:"event_type"`   
	AggregateId string    `json:"aggregate_id"` 
	OccuredAt   time.Time `json:"occurred_at"`  
	ProcessedAt time.Time `json:"processed_at"` 
}

func NewDomainEvent(name EventName, aggregateId string) DomainEvent {
	return DomainEvent{
		Name:        name,
		AggregateId: aggregateId,
		OccuredAt:   time.Now(),
	}
}

func (e *DomainEvent) SetProcessedAt() {
	e.ProcessedAt = time.Now()
}