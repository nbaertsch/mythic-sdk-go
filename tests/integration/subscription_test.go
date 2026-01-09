//go:build integration

package integration

import (
	"context"
	"strings"
	"testing"

	"github.com/nbaertsch/mythic-sdk-go/pkg/mythic/types"
)

// TestSubscribe_NotImplemented tests that Subscribe returns not implemented error.
// WebSocket subscriptions require additional infrastructure and are marked as
// not yet implemented in the current SDK version.
func TestSubscribe_NotImplemented(t *testing.T) {
	client := getTestClient(t)
	ctx := context.Background()

	handler := func(event *types.SubscriptionEvent) error {
		t.Logf("Received event: %s", event.String())
		return nil
	}

	config := &types.SubscriptionConfig{
		Type:       types.SubscriptionTypeTaskOutput,
		Handler:    handler,
		BufferSize: 100,
	}

	// Subscribe should return "not implemented" error
	sub, err := client.Subscribe(ctx, config)
	if err == nil {
		t.Fatal("Subscribe should return error for unimplemented WebSocket support")
	}

	if !strings.Contains(err.Error(), "not yet implemented") {
		t.Errorf("Expected 'not yet implemented' error, got: %v", err)
	}

	t.Logf("Subscribe correctly returns not implemented error: %v", err)

	if sub != nil {
		t.Error("Subscribe should return nil subscription when error occurs")
	}
}

// TestSubscribe_InvalidConfig tests validation of subscription configuration.
func TestSubscribe_InvalidConfig(t *testing.T) {
	client := getTestClient(t)
	ctx := context.Background()

	tests := []struct {
		name      string
		config    *types.SubscriptionConfig
		shouldErr bool
		errMsg    string
	}{
		{
			name: "empty subscription type",
			config: &types.SubscriptionConfig{
				Type:    "",
				Handler: func(e *types.SubscriptionEvent) error { return nil },
			},
			shouldErr: true,
			errMsg:    "type cannot be empty",
		},
		{
			name: "nil handler",
			config: &types.SubscriptionConfig{
				Type:    types.SubscriptionTypeCallback,
				Handler: nil,
			},
			shouldErr: true,
			errMsg:    "handler cannot be nil",
		},
		{
			name: "negative buffer size",
			config: &types.SubscriptionConfig{
				Type:       types.SubscriptionTypeFile,
				Handler:    func(e *types.SubscriptionEvent) error { return nil },
				BufferSize: -10,
			},
			shouldErr: true,
			errMsg:    "buffer size cannot be negative",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sub, err := client.Subscribe(ctx, tt.config)

			if tt.shouldErr {
				if err == nil {
					t.Error("Subscribe should return validation error")
				} else if !strings.Contains(err.Error(), tt.errMsg) {
					t.Errorf("Error should contain %q, got: %v", tt.errMsg, err)
				} else {
					t.Logf("Validation error (expected): %v", err)
				}
			}

			if sub != nil {
				t.Error("Subscribe should return nil when error occurs")
			}
		})
	}
}

// TestSubscribe_ValidConfigNotImplemented tests valid config still returns not implemented.
func TestSubscribe_ValidConfigNotImplemented(t *testing.T) {
	client := getTestClient(t)
	ctx := context.Background()

	// Test all subscription types
	types := []types.SubscriptionType{
		types.SubscriptionTypeTaskOutput,
		types.SubscriptionTypeCallback,
		types.SubscriptionTypeFile,
		types.SubscriptionTypeAll,
	}

	for _, subType := range types {
		t.Run(string(subType), func(t *testing.T) {
			config := &types.SubscriptionConfig{
				Type:       subType,
				Handler:    func(e *types.SubscriptionEvent) error { return nil },
				BufferSize: 50,
			}

			sub, err := client.Subscribe(ctx, config)

			// Should get not implemented error even with valid config
			if err == nil {
				t.Error("Subscribe should return not implemented error")
			} else if !strings.Contains(err.Error(), "not yet implemented") {
				t.Errorf("Expected not implemented error, got: %v", err)
			} else {
				t.Logf("Type %s: Not implemented (expected)", subType)
			}

			if sub != nil {
				t.Error("Subscribe should return nil subscription")
			}
		})
	}
}

// TestUnsubscribe_NilSubscription tests unsubscribe with nil subscription.
func TestUnsubscribe_NilSubscription(t *testing.T) {
	client := getTestClient(t)
	ctx := context.Background()

	err := client.Unsubscribe(ctx, nil)
	if err == nil {
		t.Fatal("Unsubscribe with nil should return error")
	}

	if !strings.Contains(err.Error(), "cannot be nil") {
		t.Errorf("Error should mention nil, got: %v", err)
	}

	t.Logf("Nil subscription error (expected): %v", err)
}

// TestUnsubscribe_InactiveSubscription tests unsubscribe with inactive subscription.
func TestUnsubscribe_InactiveSubscription(t *testing.T) {
	client := getTestClient(t)
	ctx := context.Background()

	// Create an inactive subscription
	sub := &types.Subscription{
		ID:     "test-sub",
		Type:   types.SubscriptionTypeTaskOutput,
		Active: false,
		Events: make(chan *types.SubscriptionEvent),
		Errors: make(chan error),
		Done:   make(chan struct{}),
	}

	err := client.Unsubscribe(ctx, sub)
	if err == nil {
		t.Fatal("Unsubscribe inactive subscription should return error")
	}

	if !strings.Contains(err.Error(), "not active") {
		t.Errorf("Error should mention not active, got: %v", err)
	}

	t.Logf("Inactive subscription error (expected): %v", err)
}

// TestUnsubscribe_ActiveSubscription tests unsubscribe with active subscription.
func TestUnsubscribe_ActiveSubscription(t *testing.T) {
	client := getTestClient(t)
	ctx := context.Background()

	// Create an active subscription
	sub := &types.Subscription{
		ID:     "test-sub",
		Type:   types.SubscriptionTypeCallback,
		Active: true,
		Events: make(chan *types.SubscriptionEvent, 10),
		Errors: make(chan error, 10),
		Done:   make(chan struct{}),
	}

	// Unsubscribe should succeed
	err := client.Unsubscribe(ctx, sub)
	if err != nil {
		t.Fatalf("Unsubscribe active subscription should succeed: %v", err)
	}

	// Verify subscription is now inactive
	if sub.Active {
		t.Error("Subscription should be inactive after Unsubscribe")
	}

	t.Log("Successfully unsubscribed active subscription")
}

// TestSubscriptionConfig_WithFilters tests subscription with filter criteria.
func TestSubscriptionConfig_WithFilters(t *testing.T) {
	client := getTestClient(t)
	ctx := context.Background()

	config := &types.SubscriptionConfig{
		Type:    types.SubscriptionTypeTaskOutput,
		Handler: func(e *types.SubscriptionEvent) error { return nil },
		Filter: map[string]interface{}{
			"task_id": 123,
			"status":  "running",
		},
		BufferSize: 100,
	}

	// Should still return not implemented even with filters
	sub, err := client.Subscribe(ctx, config)
	if err == nil {
		t.Fatal("Subscribe with filters should return not implemented error")
	}

	if !strings.Contains(err.Error(), "not yet implemented") {
		t.Errorf("Expected not implemented error, got: %v", err)
	}

	if sub != nil {
		t.Error("Subscribe should return nil subscription")
	}

	t.Log("Subscribe with filters correctly returns not implemented")
}

// TestSubscriptionConfig_WithOperationID tests subscription with specific operation.
func TestSubscriptionConfig_WithOperationID(t *testing.T) {
	client := getTestClient(t)
	ctx := context.Background()

	config := &types.SubscriptionConfig{
		Type:        types.SubscriptionTypeCallback,
		Handler:     func(e *types.SubscriptionEvent) error { return nil },
		OperationID: 1,
		BufferSize:  50,
	}

	// Should still return not implemented
	sub, err := client.Subscribe(ctx, config)
	if err == nil {
		t.Fatal("Subscribe with operation ID should return not implemented error")
	}

	if !strings.Contains(err.Error(), "not yet implemented") {
		t.Errorf("Expected not implemented error, got: %v", err)
	}

	if sub != nil {
		t.Error("Subscribe should return nil subscription")
	}

	t.Log("Subscribe with operation ID correctly returns not implemented")
}

// TestSubscription_Close tests subscription Close method.
func TestSubscription_Close(t *testing.T) {
	// Create a test subscription
	sub := &types.Subscription{
		ID:     "close-test",
		Type:   types.SubscriptionTypeAll,
		Active: true,
		Events: make(chan *types.SubscriptionEvent, 10),
		Errors: make(chan error, 10),
		Done:   make(chan struct{}),
	}

	if !sub.Active {
		t.Fatal("Subscription should start as active")
	}

	// Close the subscription
	sub.Close()

	if sub.Active {
		t.Error("Subscription should be inactive after Close()")
	}

	// Verify Done channel is closed
	select {
	case _, ok := <-sub.Done:
		if ok {
			t.Error("Done channel should be closed")
		}
	default:
		t.Error("Done channel should be closed and readable")
	}

	t.Log("Subscription closed successfully")
}

// TestSubscriptionEvent_DataAccess tests accessing event data fields.
func TestSubscriptionEvent_DataAccess(t *testing.T) {
	event := &types.SubscriptionEvent{
		Type:      types.SubscriptionTypeTaskOutput,
		Timestamp: "2026-01-08T12:00:00Z",
		Data: map[string]interface{}{
			"task_id": 123,
			"output":  "command executed successfully",
			"status":  "completed",
		},
	}

	// Test GetDataField with existing field
	taskID, ok := event.GetDataField("task_id")
	if !ok {
		t.Error("GetDataField should find task_id")
	}
	if taskID != 123 {
		t.Errorf("task_id = %v, want 123", taskID)
	}

	// Test GetDataField with missing field
	_, ok = event.GetDataField("missing_field")
	if ok {
		t.Error("GetDataField should return false for missing field")
	}

	// Test String method
	str := event.String()
	if !strings.Contains(str, "task_output") {
		t.Errorf("String() should contain subscription type")
	}

	t.Log("Event data access works correctly")
}
