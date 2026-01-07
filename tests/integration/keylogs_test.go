//go:build integration

package integration

import (
	"context"
	"testing"
	"time"
)

func TestKeylogs_GetKeylogs(t *testing.T) {
	SkipIfNoMythic(t)

	client := AuthenticateTestClient(t)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	keylogs, err := client.GetKeylogs(ctx)
	if err != nil {
		t.Fatalf("GetKeylogs failed: %v", err)
	}

	if keylogs == nil {
		t.Fatal("GetKeylogs returned nil")
	}

	t.Logf("Found %d keylog(s)", len(keylogs))

	// If there are keylogs, verify structure
	if len(keylogs) > 0 {
		kl := keylogs[0]
		if kl.ID == 0 {
			t.Error("Keylog ID should not be 0")
		}
		if kl.CallbackID == 0 {
			t.Error("Keylog should have a callback ID")
		}
		t.Logf("First keylog: %s", kl.String())
		t.Logf("  - Window: %s", kl.Window)
		t.Logf("  - User: %s", kl.User)
		t.Logf("  - Has keystrokes: %v", kl.HasKeystrokes())
		if kl.HasKeystrokes() {
			t.Logf("  - Keystrokes length: %d characters", len(kl.Keystrokes))
		}
	}
}

func TestKeylogs_GetKeylogsByOperation(t *testing.T) {
	SkipIfNoMythic(t)

	client := AuthenticateTestClient(t)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Get current operation
	currentOpID := client.GetCurrentOperation()
	if currentOpID == nil {
		t.Skip("No current operation set")
	}

	keylogs, err := client.GetKeylogsByOperation(ctx, *currentOpID)
	if err != nil {
		t.Fatalf("GetKeylogsByOperation failed: %v", err)
	}

	if keylogs == nil {
		t.Fatal("GetKeylogsByOperation returned nil")
	}

	t.Logf("Found %d keylog(s) for operation %d", len(keylogs), *currentOpID)

	// Verify all keylogs belong to the operation
	for _, kl := range keylogs {
		if kl.OperationID != *currentOpID {
			t.Errorf("Expected operation ID %d, got %d", *currentOpID, kl.OperationID)
		}
	}

	// Verify keylogs are sorted by timestamp (descending)
	if len(keylogs) > 1 {
		for i := 1; i < len(keylogs); i++ {
			if keylogs[i].Timestamp.After(keylogs[i-1].Timestamp) {
				t.Error("Keylogs should be sorted by timestamp (descending/newest first)")
				break
			}
		}
	}
}

func TestKeylogs_GetKeylogsByOperation_InvalidID(t *testing.T) {
	SkipIfNoMythic(t)

	client := AuthenticateTestClient(t)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Test with zero ID
	_, err := client.GetKeylogsByOperation(ctx, 0)
	if err == nil {
		t.Fatal("Expected error for zero operation ID, got nil")
	}

	t.Logf("Expected error: %v", err)
}

func TestKeylogs_GetKeylogsByCallback(t *testing.T) {
	SkipIfNoMythic(t)

	client := AuthenticateTestClient(t)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Get active callbacks first
	callbacks, err := client.GetAllActiveCallbacks(ctx)
	if err != nil {
		t.Fatalf("GetAllActiveCallbacks failed: %v", err)
	}

	if len(callbacks) == 0 {
		t.Skip("No active callbacks available for testing")
	}

	// Get keylogs for first callback
	callbackID := callbacks[0].ID
	keylogs, err := client.GetKeylogsByCallback(ctx, callbackID)
	if err != nil {
		t.Fatalf("GetKeylogsByCallback failed: %v", err)
	}

	if keylogs == nil {
		t.Fatal("GetKeylogsByCallback returned nil")
	}

	t.Logf("Found %d keylog(s) for callback %d", len(keylogs), callbackID)

	// Verify all keylogs belong to the callback
	for _, kl := range keylogs {
		if kl.CallbackID != callbackID {
			t.Errorf("Expected callback ID %d, got %d", callbackID, kl.CallbackID)
		}
	}

	// Verify keylogs are sorted by timestamp (descending)
	if len(keylogs) > 1 {
		for i := 1; i < len(keylogs); i++ {
			if keylogs[i].Timestamp.After(keylogs[i-1].Timestamp) {
				t.Error("Keylogs should be sorted by timestamp (descending/newest first)")
				break
			}
		}
	}

	// Log details of keylogs if any exist
	if len(keylogs) > 0 {
		kl := keylogs[0]
		t.Logf("Most recent keylog:")
		t.Logf("  - ID: %d", kl.ID)
		t.Logf("  - Window: %s", kl.Window)
		t.Logf("  - User: %s", kl.User)
		t.Logf("  - Task ID: %d", kl.TaskID)
		t.Logf("  - Timestamp: %s", kl.Timestamp.Format("2006-01-02 15:04:05"))
		if kl.HasKeystrokes() {
			t.Logf("  - Keystrokes: %d characters", len(kl.Keystrokes))
		}
	}
}

func TestKeylogs_GetKeylogsByCallback_InvalidID(t *testing.T) {
	SkipIfNoMythic(t)

	client := AuthenticateTestClient(t)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Test with zero ID
	_, err := client.GetKeylogsByCallback(ctx, 0)
	if err == nil {
		t.Fatal("Expected error for zero callback ID, got nil")
	}

	t.Logf("Expected error: %v", err)
}

func TestKeylogs_KeylogHelperMethods(t *testing.T) {
	SkipIfNoMythic(t)

	client := AuthenticateTestClient(t)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	keylogs, err := client.GetKeylogs(ctx)
	if err != nil {
		t.Fatalf("GetKeylogs failed: %v", err)
	}

	if len(keylogs) == 0 {
		t.Skip("No keylogs available for testing helper methods")
	}

	kl := keylogs[0]

	// Test String() method
	str := kl.String()
	if str == "" {
		t.Error("String() should not return empty string")
	}
	t.Logf("Keylog string: %s", str)

	// Test HasKeystrokes()
	if kl.HasKeystrokes() {
		t.Logf("Keylog has keystrokes (%d characters)", len(kl.Keystrokes))
	} else {
		t.Log("Keylog has no keystrokes captured")
	}
}

func TestKeylogs_MultipleCallbacks(t *testing.T) {
	SkipIfNoMythic(t)

	client := AuthenticateTestClient(t)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Get all keylogs
	allKeylogs, err := client.GetKeylogs(ctx)
	if err != nil {
		t.Fatalf("GetKeylogs failed: %v", err)
	}

	if len(allKeylogs) == 0 {
		t.Skip("No keylogs available for testing")
	}

	// Count keylogs per callback
	callbackCounts := make(map[int]int)
	for _, kl := range allKeylogs {
		callbackCounts[kl.CallbackID]++
	}

	t.Logf("Keylogs distributed across %d callback(s)", len(callbackCounts))
	for callbackID, count := range callbackCounts {
		t.Logf("  - Callback %d: %d keylog(s)", callbackID, count)

		// Verify by fetching directly
		keylogs, err := client.GetKeylogsByCallback(ctx, callbackID)
		if err != nil {
			t.Errorf("Failed to get keylogs for callback %d: %v", callbackID, err)
			continue
		}

		if len(keylogs) != count {
			t.Errorf("Expected %d keylogs for callback %d, got %d", count, callbackID, len(keylogs))
		}
	}
}

func TestKeylogs_TimestampOrdering(t *testing.T) {
	SkipIfNoMythic(t)

	client := AuthenticateTestClient(t)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	keylogs, err := client.GetKeylogs(ctx)
	if err != nil {
		t.Fatalf("GetKeylogs failed: %v", err)
	}

	if len(keylogs) < 2 {
		t.Skip("Need at least 2 keylogs to test timestamp ordering")
	}

	// Verify descending order (newest first)
	for i := 1; i < len(keylogs); i++ {
		if keylogs[i].Timestamp.After(keylogs[i-1].Timestamp) {
			t.Errorf("Timestamp ordering broken at index %d: %s > %s",
				i,
				keylogs[i].Timestamp.Format("2006-01-02 15:04:05"),
				keylogs[i-1].Timestamp.Format("2006-01-02 15:04:05"))
		}
	}

	t.Log("Timestamp ordering verified (newest first)")
	t.Logf("  - Newest: %s", keylogs[0].Timestamp.Format("2006-01-02 15:04:05"))
	t.Logf("  - Oldest: %s", keylogs[len(keylogs)-1].Timestamp.Format("2006-01-02 15:04:05"))
}
