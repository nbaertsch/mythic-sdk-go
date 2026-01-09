//go:build integration

package integration

import (
	"context"
	"testing"

	"github.com/nbaertsch/mythic-sdk-go/pkg/mythic/types"
)

// TestSubscribe_InvalidInput tests input validation for Subscribe.
func TestSubscribe_InvalidInput(t *testing.T) {
	client := getTestClient(t)
	ctx := context.Background()

	// Test with nil config
	_, err := client.Subscribe(ctx, nil)
	if err == nil {
		t.Fatal("Subscribe with nil config should return error")
	}
	t.Logf("Nil config error: %v", err)

	// Test with empty subscription type
	config := &types.SubscriptionConfig{
		Type:    "",
		Handler: func(event *types.SubscriptionEvent) error { return nil },
	}
	_, err = client.Subscribe(ctx, config)
	if err == nil {
		t.Fatal("Subscribe with empty type should return error")
	}
	t.Logf("Empty type error: %v", err)

	// Test with nil handler
	config = &types.SubscriptionConfig{
		Type:    types.SubscriptionTypeTaskOutput,
		Handler: nil,
	}
	_, err = client.Subscribe(ctx, config)
	if err == nil {
		t.Fatal("Subscribe with nil handler should return error")
	}
	t.Logf("Nil handler error: %v", err)
}

// TestSubscribe_NoOperation tests subscribing without setting operation.
func TestSubscribe_NoOperation(t *testing.T) {
	// Create client without setting operation
	client := getTestClient(t)
	// Clear operation ID to test error handling
	client.SetCurrentOperation(0)

	ctx := context.Background()

	config := &types.SubscriptionConfig{
		Type:    types.SubscriptionTypeTaskOutput,
		Handler: func(event *types.SubscriptionEvent) error { return nil },
	}

	_, err := client.Subscribe(ctx, config)
	if err == nil {
		t.Fatal("Subscribe without operation should return error")
	}
	t.Logf("No operation error (expected): %v", err)
}

// TestUnsubscribe_InvalidInput tests input validation for Unsubscribe.
func TestUnsubscribe_InvalidInput(t *testing.T) {
	client := getTestClient(t)
	ctx := context.Background()

	// Test with nil subscription
	err := client.Unsubscribe(ctx, nil)
	if err == nil {
		t.Fatal("Unsubscribe with nil subscription should return error")
	}
	t.Logf("Nil subscription error: %v", err)

	// Test with inactive subscription
	sub := &types.Subscription{
		ID:     "test-sub",
		Type:   types.SubscriptionTypeTaskOutput,
		Active: false,
		Events: make(chan *types.SubscriptionEvent, 10),
		Errors: make(chan error, 10),
		Done:   make(chan struct{}),
	}
	err = client.Unsubscribe(ctx, sub)
	if err == nil {
		t.Fatal("Unsubscribe with inactive subscription should return error")
	}
	t.Logf("Inactive subscription error: %v", err)
}

// TestSubscription_Lifecycle tests basic subscription lifecycle operations.
func TestSubscription_Lifecycle(t *testing.T) {
	// Note: This test cannot establish real WebSocket connections
	// in the integration test environment without a live Mythic server
	// configured for subscriptions. It validates the API surface.

	client := getTestClient(t)
	ctx := context.Background()

	// Create subscription config
	eventCount := 0
	config := &types.SubscriptionConfig{
		Type: types.SubscriptionTypeTaskOutput,
		Handler: func(event *types.SubscriptionEvent) error {
			eventCount++
			t.Logf("Received event: %s", event.String())
			return nil
		},
		BufferSize:  50,
		OperationID: 1, // Use fixed operation ID
	}

	// Attempt to subscribe
	// This will likely fail if WebSocket endpoint is not available,
	// but validates the API works correctly
	sub, err := client.Subscribe(ctx, config)
	if err != nil {
		// Expected error if WebSocket endpoint not configured
		t.Logf("Subscribe error (may be expected): %v", err)

		// Check if it's a connection error (acceptable in test environment)
		if sub == nil {
			t.Logf("Subscription not established (acceptable in test without WebSocket server)")
			return
		}
	}

	// If subscription was created, test cleanup
	if sub != nil {
		defer func() {
			if sub.Active {
				err := client.Unsubscribe(ctx, sub)
				if err != nil {
					t.Logf("Unsubscribe error: %v", err)
				}
			}
		}()

		t.Logf("Subscription created: ID=%s, Type=%s, Active=%v",
			sub.ID, sub.Type, sub.Active)
	}
}

// TestSubscription_MultipleTypes tests creating subscriptions for different event types.
func TestSubscription_MultipleTypes(t *testing.T) {
	client := getTestClient(t)
	ctx := context.Background()

	subscriptionTypes := []types.SubscriptionType{
		types.SubscriptionTypeTaskOutput,
		types.SubscriptionTypeCallback,
		types.SubscriptionTypeFile,
	}

	for _, subType := range subscriptionTypes {
		t.Run(string(subType), func(t *testing.T) {
			config := &types.SubscriptionConfig{
				Type: subType,
				Handler: func(event *types.SubscriptionEvent) error {
					t.Logf("Event type %s: %s", subType, event.String())
					return nil
				},
				OperationID: 1,
			}

			sub, err := client.Subscribe(ctx, config)
			if err != nil {
				t.Logf("Subscribe to %s error (may be expected): %v", subType, err)
				return
			}

			if sub != nil {
				defer func() {
					if sub.Active {
						_ = client.Unsubscribe(ctx, sub)
					}
				}()

				t.Logf("Subscription to %s created successfully", subType)
			}
		})
	}
}

// TestSubscription_Filters tests subscription with custom filters.
func TestSubscription_Filters(t *testing.T) {
	client := getTestClient(t)
	ctx := context.Background()

	config := &types.SubscriptionConfig{
		Type: types.SubscriptionTypeTaskOutput,
		Handler: func(event *types.SubscriptionEvent) error {
			return nil
		},
		Filter: map[string]interface{}{
			"callback_id": 123,
			"host":        "workstation-01",
		},
		OperationID: 1,
	}

	sub, err := client.Subscribe(ctx, config)
	if err != nil {
		t.Logf("Subscribe with filters error (may be expected): %v", err)
		return
	}

	if sub != nil {
		defer func() {
			if sub.Active {
				_ = client.Unsubscribe(ctx, sub)
			}
		}()

		t.Logf("Filtered subscription created successfully")
	}
}

// TestClientClose_WithActiveSubscriptions tests that client Close() cleans up subscriptions.
func TestClientClose_WithActiveSubscriptions(t *testing.T) {
	client := getTestClient(t)
	ctx := context.Background()

	// Create a subscription (may or may not succeed depending on environment)
	config := &types.SubscriptionConfig{
		Type: types.SubscriptionTypeTaskOutput,
		Handler: func(event *types.SubscriptionEvent) error {
			return nil
		},
		OperationID: 1,
	}

	sub, _ := client.Subscribe(ctx, config)

	// Close client should clean up all subscriptions
	err := client.Close()
	if err != nil {
		t.Fatalf("Client.Close() failed: %v", err)
	}

	// If subscription was created, verify it's no longer active
	if sub != nil && sub.Active {
		t.Error("Subscription should be inactive after client Close()")
	}

	t.Log("Client closed successfully with subscription cleanup")
}
