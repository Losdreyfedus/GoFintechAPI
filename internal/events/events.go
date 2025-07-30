package events

import (
	"time"

	"github.com/google/uuid"
)

// Event represents a domain event
type Event struct {
	ID          string                 `json:"id"`
	Type        string                 `json:"type"`
	AggregateID string                 `json:"aggregate_id"`
	Version     int                    `json:"version"`
	Data        map[string]interface{} `json:"data"`
	Metadata    map[string]interface{} `json:"metadata"`
	Timestamp   time.Time              `json:"timestamp"`
}

// NewEvent creates a new event
func NewEvent(eventType, aggregateID string, data map[string]interface{}) *Event {
	return &Event{
		ID:          uuid.New().String(),
		Type:        eventType,
		AggregateID: aggregateID,
		Version:     1,
		Data:        data,
		Metadata:    make(map[string]interface{}),
		Timestamp:   time.Now(),
	}
}

// Event types
const (
	EventUserCreated          = "user.created"
	EventUserUpdated          = "user.updated"
	EventTransactionCreated   = "transaction.created"
	EventTransactionCompleted = "transaction.completed"
	EventTransactionFailed    = "transaction.failed"
	EventBalanceUpdated       = "balance.updated"
)

// EventStore interface for storing and retrieving events
type EventStore interface {
	Append(events []*Event) error
	GetEvents(aggregateID string) ([]*Event, error)
	GetEventsByType(eventType string) ([]*Event, error)
	GetEventsSince(timestamp time.Time) ([]*Event, error)
}

// EventPublisher interface for publishing events
type EventPublisher interface {
	Publish(event *Event) error
	Subscribe(eventType string, handler func(*Event)) error
}

// EventHandler interface for handling events
type EventHandler interface {
	Handle(event *Event) error
}

// EventBus combines EventStore and EventPublisher
type EventBus struct {
	store     EventStore
	publisher EventPublisher
	handlers  map[string][]EventHandler
}

// NewEventBus creates a new event bus
func NewEventBus(store EventStore, publisher EventPublisher) *EventBus {
	return &EventBus{
		store:     store,
		publisher: publisher,
		handlers:  make(map[string][]EventHandler),
	}
}

// Publish publishes an event
func (eb *EventBus) Publish(event *Event) error {
	// Store the event
	if err := eb.store.Append([]*Event{event}); err != nil {
		return err
	}

	// Publish the event
	if err := eb.publisher.Publish(event); err != nil {
		return err
	}

	// Handle the event
	return eb.handleEvent(event)
}

// Subscribe subscribes to an event type
func (eb *EventBus) Subscribe(eventType string, handler EventHandler) {
	eb.handlers[eventType] = append(eb.handlers[eventType], handler)
}

// handleEvent handles an event by calling all registered handlers
func (eb *EventBus) handleEvent(event *Event) error {
	handlers, exists := eb.handlers[event.Type]
	if !exists {
		return nil
	}

	for _, handler := range handlers {
		if err := handler.Handle(event); err != nil {
			return err
		}
	}

	return nil
}

// RebuildState rebuilds the state from event stream
func (eb *EventBus) RebuildState(aggregateID string, initialState interface{}, applyFunc func(interface{}, *Event) interface{}) (interface{}, error) {
	events, err := eb.store.GetEvents(aggregateID)
	if err != nil {
		return nil, err
	}

	state := initialState
	for _, event := range events {
		state = applyFunc(state, event)
	}

	return state, nil
}
