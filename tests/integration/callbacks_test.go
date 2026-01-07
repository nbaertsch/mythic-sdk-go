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
