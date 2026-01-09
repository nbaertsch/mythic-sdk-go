package types

import "fmt"

// SubscriptionType represents the type of subscription.
type SubscriptionType string

const (
	// SubscriptionTypeTaskOutput subscribes to task output updates
	SubscriptionTypeTaskOutput SubscriptionType = "task_output"
	// SubscriptionTypeCallback subscribes to callback status changes
	SubscriptionTypeCallback SubscriptionType = "callback"
	// SubscriptionTypeFile subscribes to new file uploads
	SubscriptionTypeFile SubscriptionType = "file"
	// SubscriptionTypeAlert subscribes to operational alerts
	SubscriptionTypeAlert SubscriptionType = "alert"
	// SubscriptionTypeAll subscribes to all events
	SubscriptionTypeAll SubscriptionType = "all"
)

// SubscriptionEvent represents a real-time event received via subscription.
type SubscriptionEvent struct {
	Type      SubscriptionType       `json:"type"`
	Data      map[string]interface{} `json:"data"`
	Timestamp string                 `json:"timestamp"`
}

// String returns a human-readable representation of the subscription event.
func (s *SubscriptionEvent) String() string {
	return fmt.Sprintf("Event: %s at %s", s.Type, s.Timestamp)
}

// GetDataField retrieves a specific field from the event data.
func (s *SubscriptionEvent) GetDataField(key string) (interface{}, bool) {
	if s.Data == nil {
		return nil, false
	}
	val, ok := s.Data[key]
	return val, ok
}

// SubscriptionHandler is a callback function for handling subscription events.
type SubscriptionHandler func(*SubscriptionEvent) error

// SubscriptionConfig configures a GraphQL subscription.
type SubscriptionConfig struct {
	// Type of subscription to create
	Type SubscriptionType

	// Handler function called for each event
	Handler SubscriptionHandler

	// Filter criteria for the subscription (optional)
	Filter map[string]interface{}

	// OperationID to filter events (optional, uses current if not set)
	OperationID int

	// BufferSize for the event channel (default: 100)
	BufferSize int
}

// String returns a human-readable representation of the subscription config.
func (s *SubscriptionConfig) String() string {
	if s.OperationID > 0 {
		return fmt.Sprintf("Subscription: %s (operation %d)", s.Type, s.OperationID)
	}
	return fmt.Sprintf("Subscription: %s", s.Type)
}

// Validate checks if the subscription config is valid.
func (s *SubscriptionConfig) Validate() error {
	if s.Type == "" {
		return fmt.Errorf("subscription type cannot be empty")
	}
	if s.Handler == nil {
		return fmt.Errorf("subscription handler cannot be nil")
	}
	if s.BufferSize < 0 {
		return fmt.Errorf("buffer size cannot be negative")
	}
	return nil
}

// Subscription represents an active GraphQL subscription.
type Subscription struct {
	ID     string
	Type   SubscriptionType
	Active bool
	Events chan *SubscriptionEvent
	Errors chan error
	Done   chan struct{}
}

// String returns a human-readable representation of the subscription.
func (s *Subscription) String() string {
	status := "active"
	if !s.Active {
		status = "inactive"
	}
	return fmt.Sprintf("Subscription %s: %s (%s)", s.ID, s.Type, status)
}

// Close closes the subscription and all its channels.
func (s *Subscription) Close() {
	if s.Active {
		s.Active = false
		close(s.Done)
		close(s.Events)
		close(s.Errors)
	}
}
