//go:build integration

package integration

import (
	"context"
	"testing"
	"time"

	"github.com/nbaertsch/mythic-sdk-go/pkg/mythic/types"
)

func TestCallbacks_GetAllCallbacks(t *testing.T) {
	SkipIfNoMythic(t)

	client := AuthenticateTestClient(t)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	callbacks, err := client.GetAllCallbacks(ctx)
	if err != nil {
		t.Fatalf("GetAllCallbacks failed: %v", err)
	}

	// Fresh Mythic instance may have no callbacks
	if callbacks == nil {
		t.Fatal("GetAllCallbacks returned nil")
	}

	t.Logf("Found %d total callbacks", len(callbacks))

	// If we have callbacks, verify their structure
	for i, cb := range callbacks {
		if cb == nil {
			t.Errorf("Callback at index %d is nil", i)
			continue
		}

		if cb.ID == 0 {
			t.Errorf("Callback %d has invalid ID", i)
		}

		if cb.DisplayID == 0 {
			t.Errorf("Callback %d has invalid DisplayID", i)
		}

		// Basic field validation
		if cb.Host == "" {
			t.Logf("Warning: Callback %d has empty Host", cb.DisplayID)
		}

		t.Logf("Callback %d: %s@%s (Active: %v, OS: %s)",
			cb.DisplayID, cb.User, cb.Host, cb.Active, cb.OS)
	}
}

func TestCallbacks_GetAllActiveCallbacks(t *testing.T) {
	SkipIfNoMythic(t)

	client := AuthenticateTestClient(t)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Get all callbacks first
	allCallbacks, err := client.GetAllCallbacks(ctx)
	if err != nil {
		t.Fatalf("GetAllCallbacks failed: %v", err)
	}

	// Get active callbacks
	activeCallbacks, err := client.GetAllActiveCallbacks(ctx)
	if err != nil {
		t.Fatalf("GetAllActiveCallbacks failed: %v", err)
	}

	if activeCallbacks == nil {
		t.Fatal("GetAllActiveCallbacks returned nil")
	}

	t.Logf("Found %d active callbacks out of %d total", len(activeCallbacks), len(allCallbacks))

	// All returned callbacks should be active
	for _, cb := range activeCallbacks {
		if !cb.Active {
			t.Errorf("GetAllActiveCallbacks returned inactive callback %d", cb.DisplayID)
		}
	}

	// Active count should be <= total count
	if len(activeCallbacks) > len(allCallbacks) {
		t.Errorf("Active callback count (%d) exceeds total count (%d)",
			len(activeCallbacks), len(allCallbacks))
	}
}

func TestCallbacks_GetCallbackByID(t *testing.T) {
	SkipIfNoMythic(t)

	client := AuthenticateTestClient(t)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Get all callbacks to find one to test with
	callbacks, err := client.GetAllCallbacks(ctx)
	if err != nil {
		t.Fatalf("GetAllCallbacks failed: %v", err)
	}

	if len(callbacks) == 0 {
		t.Skip("No callbacks available for testing")
	}

	// Test getting the first callback by ID
	targetCallback := callbacks[0]
	t.Logf("Testing with callback %d", targetCallback.DisplayID)

	callback, err := client.GetCallbackByID(ctx, targetCallback.DisplayID)
	if err != nil {
		t.Fatalf("GetCallbackByID failed: %v", err)
	}

	if callback == nil {
		t.Fatal("GetCallbackByID returned nil")
	}

	// Verify we got the right callback
	if callback.ID != targetCallback.ID {
		t.Errorf("Expected callback ID %d, got %d", targetCallback.ID, callback.ID)
	}

	if callback.DisplayID != targetCallback.DisplayID {
		t.Errorf("Expected display ID %d, got %d", targetCallback.DisplayID, callback.DisplayID)
	}

	// Verify key fields match
	if callback.Host != targetCallback.Host {
		t.Errorf("Expected host %q, got %q", targetCallback.Host, callback.Host)
	}

	if callback.User != targetCallback.User {
		t.Errorf("Expected user %q, got %q", targetCallback.User, callback.User)
	}

	t.Logf("Successfully retrieved callback: %s@%s", callback.User, callback.Host)
}

func TestCallbacks_GetCallbackByID_NotFound(t *testing.T) {
	SkipIfNoMythic(t)

	client := AuthenticateTestClient(t)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Use a very high display ID that shouldn't exist
	invalidID := 999999

	_, err := client.GetCallbackByID(ctx, invalidID)
	if err == nil {
		t.Fatal("GetCallbackByID should fail for non-existent callback")
	}

	t.Logf("Expected error for invalid callback ID: %v", err)
}

func TestCallbacks_UpdateCallback(t *testing.T) {
	SkipIfNoMythic(t)

	client := AuthenticateTestClient(t)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Get all callbacks to find one to test with
	callbacks, err := client.GetAllCallbacks(ctx)
	if err != nil {
		t.Fatalf("GetAllCallbacks failed: %v", err)
	}

	if len(callbacks) == 0 {
		t.Skip("No callbacks available for testing")
	}

	targetCallback := callbacks[0]
	originalDescription := targetCallback.Description

	t.Logf("Testing update with callback %d (current description: %q)",
		targetCallback.DisplayID, originalDescription)

	// Update the description
	newDescription := "Integration test updated description"
	err = client.UpdateCallback(ctx, &types.CallbackUpdateRequest{
		CallbackDisplayID: targetCallback.DisplayID,
		Description:       &newDescription,
	})
	if err != nil {
		t.Fatalf("UpdateCallback failed: %v", err)
	}

	// Verify the update
	updated, err := client.GetCallbackByID(ctx, targetCallback.DisplayID)
	if err != nil {
		t.Fatalf("GetCallbackByID after update failed: %v", err)
	}

	if updated.Description != newDescription {
		t.Errorf("Expected description %q, got %q", newDescription, updated.Description)
	}

	t.Logf("Successfully updated callback description to: %q", updated.Description)

	// Restore original description
	if originalDescription != "" {
		err = client.UpdateCallback(ctx, &types.CallbackUpdateRequest{
			CallbackDisplayID: targetCallback.DisplayID,
			Description:       &originalDescription,
		})
		if err != nil {
			t.Logf("Warning: Failed to restore original description: %v", err)
		} else {
			t.Log("Restored original description")
		}
	}
}

func TestCallbacks_UpdateCallback_Locked(t *testing.T) {
	SkipIfNoMythic(t)

	client := AuthenticateTestClient(t)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Get all callbacks to find one to test with
	callbacks, err := client.GetAllCallbacks(ctx)
	if err != nil {
		t.Fatalf("GetAllCallbacks failed: %v", err)
	}

	if len(callbacks) == 0 {
		t.Skip("No callbacks available for testing")
	}

	targetCallback := callbacks[0]
	originalLocked := targetCallback.Locked

	t.Logf("Testing locked field with callback %d (current: %v)",
		targetCallback.DisplayID, originalLocked)

	// Toggle locked status
	newLocked := !originalLocked
	err = client.UpdateCallback(ctx, &types.CallbackUpdateRequest{
		CallbackDisplayID: targetCallback.DisplayID,
		Locked:            &newLocked,
	})
	if err != nil {
		t.Fatalf("UpdateCallback failed: %v", err)
	}

	// Verify the update
	updated, err := client.GetCallbackByID(ctx, targetCallback.DisplayID)
	if err != nil {
		t.Fatalf("GetCallbackByID after update failed: %v", err)
	}

	if updated.Locked != newLocked {
		t.Errorf("Expected locked=%v, got %v", newLocked, updated.Locked)
	}

	t.Logf("Successfully toggled locked status to: %v", updated.Locked)

	// Restore original locked status
	err = client.UpdateCallback(ctx, &types.CallbackUpdateRequest{
		CallbackDisplayID: targetCallback.DisplayID,
		Locked:            &originalLocked,
	})
	if err != nil {
		t.Logf("Warning: Failed to restore original locked status: %v", err)
	} else {
		t.Log("Restored original locked status")
	}
}

func TestCallbacks_IntegrityLevelHelpers(t *testing.T) {
	SkipIfNoMythic(t)

	client := AuthenticateTestClient(t)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	callbacks, err := client.GetAllCallbacks(ctx)
	if err != nil {
		t.Fatalf("GetAllCallbacks failed: %v", err)
	}

	if len(callbacks) == 0 {
		t.Skip("No callbacks available for testing")
	}

	// Test integrity level helper methods
	for _, cb := range callbacks {
		isHigh := cb.IsHigh()
		isSystem := cb.IsSystem()

		// Integrity level validation
		expectedHigh := cb.IntegrityLevel >= types.IntegrityLevelHigh
		if isHigh != expectedHigh {
			t.Errorf("Callback %d: IsHigh()=%v, expected %v (level=%d)",
				cb.DisplayID, isHigh, expectedHigh, cb.IntegrityLevel)
		}

		expectedSystem := cb.IntegrityLevel == types.IntegrityLevelSystem
		if isSystem != expectedSystem {
			t.Errorf("Callback %d: IsSystem()=%v, expected %v (level=%d)",
				cb.DisplayID, isSystem, expectedSystem, cb.IntegrityLevel)
		}

		// System should always be high
		if isSystem && !isHigh {
			t.Errorf("Callback %d: IsSystem=true but IsHigh=false", cb.DisplayID)
		}

		t.Logf("Callback %d: Integrity=%d, IsHigh=%v, IsSystem=%v",
			cb.DisplayID, cb.IntegrityLevel, isHigh, isSystem)
	}
}

func TestCallbacks_CallbackFields(t *testing.T) {
	SkipIfNoMythic(t)

	client := AuthenticateTestClient(t)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	callbacks, err := client.GetAllCallbacks(ctx)
	if err != nil {
		t.Fatalf("GetAllCallbacks failed: %v", err)
	}

	if len(callbacks) == 0 {
		t.Skip("No callbacks available for testing")
	}

	// Test field presence and validity
	callback := callbacks[0]

	tests := []struct {
		name  string
		value interface{}
		check func(interface{}) bool
	}{
		{"ID", callback.ID, func(v interface{}) bool { return v.(int) > 0 }},
		{"DisplayID", callback.DisplayID, func(v interface{}) bool { return v.(int) > 0 }},
		{"AgentCallbackID", callback.AgentCallbackID, func(v interface{}) bool { return v.(string) != "" }},
		{"InitCallback", callback.InitCallback, func(v interface{}) bool { return !v.(time.Time).IsZero() }},
		{"LastCheckin", callback.LastCheckin, func(v interface{}) bool { return !v.(time.Time).IsZero() }},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !tt.check(tt.value) {
				t.Errorf("Field %s has invalid value: %v", tt.name, tt.value)
			} else {
				t.Logf("Field %s: %v", tt.name, tt.value)
			}
		})
	}

	// Test IP parsing
	if len(callback.IP) > 0 {
		t.Logf("Callback has %d IP addresses: %v", len(callback.IP), callback.IP)
	}

	// Test time calculations
	timeSinceInit := time.Since(callback.InitCallback)
	timeSinceCheckin := time.Since(callback.LastCheckin)

	t.Logf("Time since init: %v", timeSinceInit)
	t.Logf("Time since last checkin: %v", timeSinceCheckin)

	if timeSinceInit < 0 {
		t.Error("InitCallback is in the future")
	}

	if timeSinceCheckin < 0 {
		t.Error("LastCheckin is in the future")
	}
}

func TestCallbacks_CreateCallback(t *testing.T) {
	t.Skip("Skipping CreateCallback test - requires valid payload UUID and manual verification")

	SkipIfNoMythic(t)
	client := AuthenticateTestClient(t)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Get a payload to use for callback creation
	payloads, err := client.GetPayloads(ctx)
	if err != nil {
		t.Fatalf("GetPayloads failed: %v", err)
	}
	if len(payloads) == 0 {
		t.Skip("No payloads available for callback creation")
	}

	payloadUUID := payloads[0].UUID
	user := "test_user"
	host := "test_host"
	ip := "192.168.1.100"

	input := &CreateCallbackInput{
		PayloadUUID: payloadUUID,
		User:        &user,
		Host:        &host,
		IP:          &ip,
	}

	err = client.CreateCallback(ctx, input)
	if err != nil {
		t.Fatalf("CreateCallback failed: %v", err)
	}

	t.Log("Successfully created callback (manual cleanup required)")
}

func TestCallbacks_CreateCallback_InvalidInput(t *testing.T) {
	SkipIfNoMythic(t)
	client := AuthenticateTestClient(t)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Test nil input
	err := client.CreateCallback(ctx, nil)
	if err == nil {
		t.Error("Expected error for nil input")
	}

	// Test empty payload UUID
	err = client.CreateCallback(ctx, &CreateCallbackInput{})
	if err == nil {
		t.Error("Expected error for empty payload UUID")
	}

	t.Log("All invalid input tests passed")
}

func TestCallbacks_DeleteCallback(t *testing.T) {
	t.Skip("Skipping DeleteCallback test - destructive operation")

	SkipIfNoMythic(t)
	client := AuthenticateTestClient(t)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Note: This test is skipped by default as it's destructive
	// To test manually, create a test callback first, then delete it

	t.Log("DeleteCallback is a destructive operation and should be tested manually")
}

func TestCallbacks_DeleteCallback_InvalidInput(t *testing.T) {
	SkipIfNoMythic(t)
	client := AuthenticateTestClient(t)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Test empty callback list
	err := client.DeleteCallback(ctx, []int{})
	if err == nil {
		t.Error("Expected error for empty callback list")
	}

	// Test nonexistent callback
	err = client.DeleteCallback(ctx, []int{999999})
	if err == nil {
		t.Log("Note: DeleteCallback might succeed even for non-existent callbacks")
	}

	t.Log("All invalid input tests completed")
}

func TestCallbacks_ExportCallbackConfig(t *testing.T) {
	SkipIfNoMythic(t)
	client := AuthenticateTestClient(t)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Get a callback to export
	callbacks, err := client.GetAllActiveCallbacks(ctx)
	if err != nil {
		t.Fatalf("GetAllActiveCallbacks failed: %v", err)
	}
	if len(callbacks) == 0 {
		t.Skip("No active callbacks for config export")
	}

	agentCallbackID := callbacks[0].AgentCallbackID
	config, err := client.ExportCallbackConfig(ctx, agentCallbackID)
	if err != nil {
		t.Fatalf("ExportCallbackConfig failed: %v", err)
	}

	if config == "" {
		t.Error("ExportCallbackConfig returned empty config")
	}

	t.Logf("Successfully exported config (length: %d bytes)", len(config))

	// Verify it's valid JSON
	var configMap map[string]interface{}
	if err := parseJSON([]byte(config), &configMap); err != nil {
		t.Errorf("Exported config is not valid JSON: %v", err)
	} else {
		t.Logf("Config contains %d top-level keys", len(configMap))
	}
}

func TestCallbacks_ExportCallbackConfig_InvalidInput(t *testing.T) {
	SkipIfNoMythic(t)
	client := AuthenticateTestClient(t)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Test empty agent callback ID
	_, err := client.ExportCallbackConfig(ctx, "")
	if err == nil {
		t.Error("Expected error for empty agent_callback_id")
	}

	// Test nonexistent callback
	_, err = client.ExportCallbackConfig(ctx, "nonexistent-callback-id")
	if err == nil {
		t.Error("Expected error for nonexistent callback")
	}

	t.Log("All invalid input tests passed")
}

func TestCallbacks_ImportCallbackConfig(t *testing.T) {
	t.Skip("Skipping ImportCallbackConfig test - requires valid exported config")

	SkipIfNoMythic(t)
	client := AuthenticateTestClient(t)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Note: This test requires a valid exported config
	// In practice, you would export a config first, then import it

	t.Log("ImportCallbackConfig should be tested with a valid exported config")
}

func TestCallbacks_ImportCallbackConfig_InvalidInput(t *testing.T) {
	SkipIfNoMythic(t)
	client := AuthenticateTestClient(t)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Test empty config
	err := client.ImportCallbackConfig(ctx, "")
	if err == nil {
		t.Error("Expected error for empty config")
	}

	// Test invalid JSON
	err = client.ImportCallbackConfig(ctx, "not valid json")
	if err == nil {
		t.Error("Expected error for invalid JSON")
	}

	t.Log("All invalid input tests passed")
}

func TestCallbacks_AddCallbackGraphEdge(t *testing.T) {
	t.Skip("Skipping AddCallbackGraphEdge test - requires multiple callbacks")

	SkipIfNoMythic(t)
	client := AuthenticateTestClient(t)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Get active callbacks
	callbacks, err := client.GetAllActiveCallbacks(ctx)
	if err != nil {
		t.Fatalf("GetAllActiveCallbacks failed: %v", err)
	}
	if len(callbacks) < 2 {
		t.Skip("Need at least 2 callbacks for graph edge test")
	}

	// Get C2 profiles
	profiles, err := client.GetC2Profiles(ctx, 0)
	if err != nil || len(profiles) == 0 {
		t.Skip("No C2 profiles available")
	}

	sourceID := callbacks[0].DisplayID
	destID := callbacks[1].DisplayID
	c2ProfileName := profiles[0].Name

	err = client.AddCallbackGraphEdge(ctx, sourceID, destID, c2ProfileName)
	if err != nil {
		t.Fatalf("AddCallbackGraphEdge failed: %v", err)
	}

	t.Logf("Successfully added graph edge: %d -> %d (profile: %s)", sourceID, destID, c2ProfileName)
}

func TestCallbacks_AddCallbackGraphEdge_InvalidInput(t *testing.T) {
	SkipIfNoMythic(t)
	client := AuthenticateTestClient(t)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Test zero source ID
	err := client.AddCallbackGraphEdge(ctx, 0, 1, "http")
	if err == nil {
		t.Error("Expected error for zero source ID")
	}

	// Test zero destination ID
	err = client.AddCallbackGraphEdge(ctx, 1, 0, "http")
	if err == nil {
		t.Error("Expected error for zero destination ID")
	}

	// Test empty C2 profile name
	err = client.AddCallbackGraphEdge(ctx, 1, 2, "")
	if err == nil {
		t.Error("Expected error for empty C2 profile name")
	}

	t.Log("All invalid input tests passed")
}

func TestCallbacks_RemoveCallbackGraphEdge(t *testing.T) {
	t.Skip("Skipping RemoveCallbackGraphEdge test - requires existing graph edge")

	SkipIfNoMythic(t)
	client := AuthenticateTestClient(t)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Note: This test requires an existing edge ID
	// In practice, you would add an edge first, then remove it

	t.Log("RemoveCallbackGraphEdge should be tested with an existing edge ID")
}

func TestCallbacks_RemoveCallbackGraphEdge_InvalidInput(t *testing.T) {
	SkipIfNoMythic(t)
	client := AuthenticateTestClient(t)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Test zero edge ID
	err := client.RemoveCallbackGraphEdge(ctx, 0)
	if err == nil {
		t.Error("Expected error for zero edge ID")
	}

	// Test negative edge ID
	err = client.RemoveCallbackGraphEdge(ctx, -1)
	if err == nil {
		t.Error("Expected error for negative edge ID")
	}

	t.Log("All invalid input tests passed")
}
