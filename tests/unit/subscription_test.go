package unit

import (
	"strings"
	"testing"

	"github.com/nbaertsch/mythic-sdk-go/pkg/mythic/types"
)

func TestSubscriptionEvent_String(t *testing.T) {
	event := types.SubscriptionEvent{
		Type:      types.SubscriptionTypeTaskOutput,
		Timestamp: "2026-01-08T10:00:00Z",
		Data:      map[string]interface{}{"output": "test data"},
	}

	result := event.String()

	if !strings.Contains(result, "task_output") {
		t.Errorf("String() = %q, should contain 'task_output'", result)
	}
	if !strings.Contains(result, "2026-01-08") {
		t.Errorf("String() = %q, should contain timestamp", result)
	}
}

func TestSubscriptionEvent_GetDataField(t *testing.T) {
	tests := []struct {
		name      string
		event     types.SubscriptionEvent
		field     string
		expectVal interface{}
		expectOk  bool
	}{
		{
			name: "existing field",
			event: types.SubscriptionEvent{
				Data: map[string]interface{}{
					"output":  "test data",
					"task_id": 123,
				},
			},
			field:     "output",
			expectVal: "test data",
			expectOk:  true,
		},
		{
			name: "missing field",
			event: types.SubscriptionEvent{
				Data: map[string]interface{}{
					"output": "test data",
				},
			},
			field:     "missing",
			expectVal: nil,
			expectOk:  false,
		},
		{
			name: "nil data",
			event: types.SubscriptionEvent{
				Data: nil,
			},
			field:     "any",
			expectVal: nil,
			expectOk:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			val, ok := tt.event.GetDataField(tt.field)
			if ok != tt.expectOk {
				t.Errorf("GetDataField() ok = %v, want %v", ok, tt.expectOk)
			}
			if ok && val != tt.expectVal {
				t.Errorf("GetDataField() value = %v, want %v", val, tt.expectVal)
			}
		})
	}
}

func TestSubscriptionConfig_String(t *testing.T) {
	tests := []struct {
		name     string
		config   types.SubscriptionConfig
		contains []string
	}{
		{
			name: "with operation ID",
			config: types.SubscriptionConfig{
				Type:        types.SubscriptionTypeCallback,
				OperationID: 42,
			},
			contains: []string{"callback", "42"},
		},
		{
			name: "without operation ID",
			config: types.SubscriptionConfig{
				Type: types.SubscriptionTypeTaskOutput,
			},
			contains: []string{"task_output"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.config.String()

			for _, substr := range tt.contains {
				if !strings.Contains(result, substr) {
					t.Errorf("String() = %q, should contain %q", result, substr)
				}
			}
		})
	}
}

func TestSubscriptionConfig_Validate(t *testing.T) {
	tests := []struct {
		name      string
		config    types.SubscriptionConfig
		shouldErr bool
		errMsg    string
	}{
		{
			name: "valid config",
			config: types.SubscriptionConfig{
				Type:       types.SubscriptionTypeTaskOutput,
				Handler:    func(e *types.SubscriptionEvent) error { return nil },
				BufferSize: 100,
			},
			shouldErr: false,
		},
		{
			name: "empty type",
			config: types.SubscriptionConfig{
				Handler:    func(e *types.SubscriptionEvent) error { return nil },
				BufferSize: 100,
			},
			shouldErr: true,
			errMsg:    "type cannot be empty",
		},
		{
			name: "nil handler",
			config: types.SubscriptionConfig{
				Type:       types.SubscriptionTypeCallback,
				BufferSize: 100,
			},
			shouldErr: true,
			errMsg:    "handler cannot be nil",
		},
		{
			name: "negative buffer size",
			config: types.SubscriptionConfig{
				Type:       types.SubscriptionTypeFile,
				Handler:    func(e *types.SubscriptionEvent) error { return nil },
				BufferSize: -10,
			},
			shouldErr: true,
			errMsg:    "buffer size cannot be negative",
		},
		{
			name: "zero buffer size (valid)",
			config: types.SubscriptionConfig{
				Type:       types.SubscriptionTypeAll,
				Handler:    func(e *types.SubscriptionEvent) error { return nil },
				BufferSize: 0,
			},
			shouldErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()

			if tt.shouldErr {
				if err == nil {
					t.Error("Validate() should return error but got nil")
				} else if !strings.Contains(err.Error(), tt.errMsg) {
					t.Errorf("Validate() error = %q, should contain %q", err.Error(), tt.errMsg)
				}
			} else {
				if err != nil {
					t.Errorf("Validate() should not error but got: %v", err)
				}
			}
		})
	}
}

func TestSubscription_String(t *testing.T) {
	tests := []struct {
		name     string
		sub      types.Subscription
		contains []string
	}{
		{
			name: "active subscription",
			sub: types.Subscription{
				ID:     "sub-123",
				Type:   types.SubscriptionTypeTaskOutput,
				Active: true,
			},
			contains: []string{"sub-123", "task_output", "active"},
		},
		{
			name: "inactive subscription",
			sub: types.Subscription{
				ID:     "sub-456",
				Type:   types.SubscriptionTypeCallback,
				Active: false,
			},
			contains: []string{"sub-456", "callback", "inactive"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.sub.String()

			for _, substr := range tt.contains {
				if !strings.Contains(result, substr) {
					t.Errorf("String() = %q, should contain %q", result, substr)
				}
			}
		})
	}
}

func TestSubscription_Close(t *testing.T) {
	sub := &types.Subscription{
		ID:     "test-sub",
		Type:   types.SubscriptionTypeAll,
		Active: true,
		Events: make(chan *types.SubscriptionEvent, 10),
		Errors: make(chan error, 10),
		Done:   make(chan struct{}),
	}

	if !sub.Active {
		t.Error("Subscription should start as active")
	}

	// Close the subscription
	sub.Close()

	if sub.Active {
		t.Error("Subscription should be inactive after Close()")
	}

	// Verify channels are closed by checking if receive returns immediately
	select {
	case _, ok := <-sub.Done:
		if ok {
			t.Error("Done channel should be closed")
		}
	default:
		t.Error("Done channel should be readable (closed)")
	}
}

func TestSubscriptionType_Constants(t *testing.T) {
	tests := []struct {
		name     string
		subType  types.SubscriptionType
		expected string
	}{
		{"task output", types.SubscriptionTypeTaskOutput, "task_output"},
		{"callback", types.SubscriptionTypeCallback, "callback"},
		{"file", types.SubscriptionTypeFile, "file"},
		{"all", types.SubscriptionTypeAll, "all"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if string(tt.subType) != tt.expected {
				t.Errorf("SubscriptionType = %q, want %q", tt.subType, tt.expected)
			}
		})
	}
}
