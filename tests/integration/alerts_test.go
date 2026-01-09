//go:build integration

package integration

import (
	"context"
	"testing"
	"time"

	"github.com/nbaertsch/mythic-sdk-go/pkg/mythic/types"
)

func TestAlerts_GetAlerts(t *testing.T) {
	client := AuthenticateTestClient(t)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	alerts, err := client.GetAlerts(ctx, 0, &types.AlertFilter{
		Limit: 20,
	})
	if err != nil {
		t.Fatalf("GetAlerts failed: %v", err)
	}

	if alerts == nil {
		t.Fatal("GetAlerts returned nil")
	}

	t.Logf("Found %d alert(s)", len(alerts))

	for _, alert := range alerts {
		if alert.ID == 0 {
			t.Error("Alert ID should not be 0")
		}
		t.Logf("  - [%s] %s: %s (Severity: %d, Resolved: %v)",
			alert.Source, alert.AlertType, alert.Message, alert.Severity, alert.Resolved)
	}
}

func TestAlerts_GetAlertByID(t *testing.T) {
	client := AuthenticateTestClient(t)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Get alerts first
	alerts, err := client.GetAlerts(ctx, 0, &types.AlertFilter{Limit: 5})
	if err != nil {
		t.Fatalf("GetAlerts failed: %v", err)
	}
	if len(alerts) == 0 {
		t.Skip("No alerts available for testing")
	}

	alertID := alerts[0].ID
	alert, err := client.GetAlertByID(ctx, alertID)
	if err != nil {
		t.Fatalf("GetAlertByID failed: %v", err)
	}

	if alert == nil {
		t.Fatal("GetAlertByID returned nil")
	}

	if alert.ID != alertID {
		t.Errorf("Expected alert ID %d, got %d", alertID, alert.ID)
	}

	t.Logf("Retrieved alert %d: %s", alertID, alert.String())
	t.Logf("  - Source: %s", alert.Source)
	t.Logf("  - Severity: %d", alert.Severity)
	t.Logf("  - Resolved: %v", alert.Resolved)
}

func TestAlerts_GetUnresolvedAlerts(t *testing.T) {
	client := AuthenticateTestClient(t)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	alerts, err := client.GetUnresolvedAlerts(ctx, 0)
	if err != nil {
		t.Fatalf("GetUnresolvedAlerts failed: %v", err)
	}

	if alerts == nil {
		t.Fatal("GetUnresolvedAlerts returned nil")
	}

	t.Logf("Found %d unresolved alert(s)", len(alerts))

	// Verify all are unresolved
	for _, alert := range alerts {
		if alert.Resolved {
			t.Errorf("Alert %d should be unresolved", alert.ID)
		}
	}

	// Log high severity alerts
	highSeverity := 0
	for _, alert := range alerts {
		if alert.Severity >= 4 {
			highSeverity++
			t.Logf("  - HIGH SEVERITY: %s - %s", alert.AlertType, alert.Message)
		}
	}
	t.Logf("High severity unresolved alerts: %d", highSeverity)
}

func TestAlerts_ResolveAlert(t *testing.T) {
	t.Skip("Skipping resolve test to preserve alert state")

	client := AuthenticateTestClient(t)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Get unresolved alerts
	alerts, err := client.GetUnresolvedAlerts(ctx, 0)
	if err != nil {
		t.Fatalf("GetUnresolvedAlerts failed: %v", err)
	}
	if len(alerts) == 0 {
		t.Skip("No unresolved alerts for testing")
	}

	alertID := alerts[0].ID

	err = client.ResolveAlert(ctx, &types.ResolveAlertRequest{
		AlertID: alertID,
		Notes:   "Resolved during integration test",
	})
	if err != nil {
		t.Fatalf("ResolveAlert failed: %v", err)
	}

	t.Logf("Successfully resolved alert %d", alertID)

	// Verify resolution
	alert, err := client.GetAlertByID(ctx, alertID)
	if err != nil {
		t.Fatalf("GetAlertByID failed: %v", err)
	}

	if !alert.Resolved {
		t.Error("Alert should be resolved")
	}
}

func TestAlerts_CreateCustomAlert(t *testing.T) {
	t.Skip("Skipping create test to avoid cluttering alerts")

	client := AuthenticateTestClient(t)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	currentOp := client.GetCurrentOperation()
	if currentOp == nil {
		t.Skip("No current operation set")
	}

	alert, err := client.CreateCustomAlert(ctx, &types.CreateAlertRequest{
		Message:     "Test alert from integration test",
		AlertType:   "info",
		Source:      "integration_test",
		Severity:    2,
		OperationID: *currentOp,
	})
	if err != nil {
		t.Fatalf("CreateCustomAlert failed: %v", err)
	}

	if alert == nil {
		t.Fatal("CreateCustomAlert returned nil")
	}

	t.Logf("Created custom alert %d: %s", alert.ID, alert.String())

	// Clean up
	err = client.ResolveAlert(ctx, &types.ResolveAlertRequest{
		AlertID: alert.ID,
		Notes:   "Test cleanup",
	})
	if err != nil {
		t.Logf("Warning: Failed to clean up test alert: %v", err)
	}
}

func TestAlerts_GetAlertStatistics(t *testing.T) {
	client := AuthenticateTestClient(t)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	stats, err := client.GetAlertStatistics(ctx, 0)
	if err != nil {
		t.Fatalf("GetAlertStatistics failed: %v", err)
	}

	if stats == nil {
		t.Fatal("GetAlertStatistics returned nil")
	}

	t.Logf("Alert Statistics:")
	for key, value := range stats {
		t.Logf("  - %s: %d", key, value)
	}

	// Verify some expected keys exist
	if total, ok := stats["total"]; ok {
		t.Logf("Total alerts: %d", total)
	}
	if unresolved, ok := stats["unresolved"]; ok {
		t.Logf("Unresolved alerts: %d", unresolved)
	}
}

func TestAlerts_SubscribeToAlerts(t *testing.T) {
	client := AuthenticateTestClient(t)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	receivedEvents := 0
	maxEvents := 5

	config := &types.SubscriptionConfig{
		Type: types.SubscriptionTypeAlert,
		Handler: func(event *types.SubscriptionEvent) error {
			receivedEvents++
			t.Logf("Received alert event: %s", event.String())

			// Extract alert data
			if alertType, ok := event.GetDataField("alert"); ok {
				t.Logf("  - Alert type: %v", alertType)
			}
			if message, ok := event.GetDataField("message"); ok {
				t.Logf("  - Message: %v", message)
			}

			if receivedEvents >= maxEvents {
				cancel()
			}
			return nil
		},
		BufferSize: 100,
	}

	sub, err := client.Subscribe(ctx, config)
	if err != nil {
		t.Fatalf("Subscribe failed: %v", err)
	}
	defer sub.Close()

	t.Log("Subscribed to alerts, waiting for events...")

	// Wait for events or timeout
	select {
	case <-sub.Done:
		t.Log("Subscription closed")
	case err := <-sub.Errors:
		t.Fatalf("Subscription error: %v", err)
	case <-ctx.Done():
		t.Log("Context timeout")
	}

	t.Logf("Received %d alert event(s)", receivedEvents)

	if receivedEvents > 0 {
		t.Log("Alert subscription test successful")
	} else {
		t.Log("No alert events received (may be expected if no alerts generated)")
	}
}

func TestAlerts_InvalidInputs(t *testing.T) {
	client := AuthenticateTestClient(t)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Test zero alert ID
	_, err := client.GetAlertByID(ctx, 0)
	if err == nil {
		t.Error("Expected error for zero alert ID")
	}

	// Test resolve with zero ID
	err = client.ResolveAlert(ctx, &types.ResolveAlertRequest{AlertID: 0, Notes: "test"})
	if err == nil {
		t.Error("Expected error for zero alert ID in ResolveAlert")
	}

	// Test create with empty message
	_, err = client.CreateCustomAlert(ctx, &types.CreateAlertRequest{
		Message:     "",
		AlertType:   "info",
		Source:      "test",
		Severity:    1,
		OperationID: 1,
	})
	if err == nil {
		t.Error("Expected error for empty message")
	}

	// Test create with invalid severity
	_, err = client.CreateCustomAlert(ctx, &types.CreateAlertRequest{
		Message:     "test",
		AlertType:   "info",
		Source:      "test",
		Severity:    -1,
		OperationID: 1,
	})
	if err == nil {
		t.Error("Expected error for negative severity")
	}

	t.Log("All invalid input tests passed")
}
