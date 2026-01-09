//go:build integration

package integration

import (
	"context"
	"testing"
)

// TestBlockList_DeleteInvalidInput tests input validation for DeleteBlockList.
func TestBlockList_DeleteInvalidInput(t *testing.T) {
	client := getTestClient(t)
	ctx := context.Background()

	// Test with zero block list ID
	_, err := client.DeleteBlockList(ctx, 0)
	if err == nil {
		t.Fatal("DeleteBlockList with zero ID should return error")
	}
	t.Logf("Zero ID error: %v", err)

	// Test with negative block list ID
	_, err = client.DeleteBlockList(ctx, -1)
	if err == nil {
		t.Fatal("DeleteBlockList with negative ID should return error")
	}
	t.Logf("Negative ID error: %v", err)
}

// TestBlockList_DeleteNonexistent tests deleting a block list that doesn't exist.
func TestBlockList_DeleteNonexistent(t *testing.T) {
	client := getTestClient(t)
	ctx := context.Background()

	// Try to delete a block list that likely doesn't exist
	result, err := client.DeleteBlockList(ctx, 999999)
	if err == nil {
		t.Fatal("DeleteBlockList with nonexistent ID should return error")
	}
	t.Logf("Nonexistent block list error (expected): %v", err)

	// The result might be nil or contain error details
	if result != nil {
		t.Logf("Response: %s", result.String())
		if result.IsSuccessful() {
			t.Error("Should not be successful when block list doesn't exist")
		}
	}
}

// TestBlockListEntry_DeleteInvalidInput tests input validation for DeleteBlockListEntry.
func TestBlockListEntry_DeleteInvalidInput(t *testing.T) {
	client := getTestClient(t)
	ctx := context.Background()

	// Test with empty entry IDs list
	_, err := client.DeleteBlockListEntry(ctx, []int{})
	if err == nil {
		t.Fatal("DeleteBlockListEntry with empty list should return error")
	}
	t.Logf("Empty list error: %v", err)

	// Test with zero entry ID
	_, err = client.DeleteBlockListEntry(ctx, []int{0})
	if err == nil {
		t.Fatal("DeleteBlockListEntry with zero ID should return error")
	}
	t.Logf("Zero ID error: %v", err)

	// Test with negative entry ID
	_, err = client.DeleteBlockListEntry(ctx, []int{-1})
	if err == nil {
		t.Fatal("DeleteBlockListEntry with negative ID should return error")
	}
	t.Logf("Negative ID error: %v", err)

	// Test with mixed valid and invalid IDs
	_, err = client.DeleteBlockListEntry(ctx, []int{123, 0, 456})
	if err == nil {
		t.Fatal("DeleteBlockListEntry with zero in list should return error")
	}
	t.Logf("Mixed IDs error: %v", err)
}

// TestBlockListEntry_DeleteDuplicates tests handling of duplicate entry IDs.
func TestBlockListEntry_DeleteDuplicates(t *testing.T) {
	client := getTestClient(t)
	ctx := context.Background()

	// Test with duplicate entry IDs
	_, err := client.DeleteBlockListEntry(ctx, []int{123, 456, 123})
	if err == nil {
		t.Fatal("DeleteBlockListEntry with duplicate IDs should return error")
	}
	t.Logf("Duplicate IDs error: %v", err)

	// Test with all same IDs
	_, err = client.DeleteBlockListEntry(ctx, []int{789, 789, 789})
	if err == nil {
		t.Fatal("DeleteBlockListEntry with all duplicate IDs should return error")
	}
	t.Logf("All duplicate IDs error: %v", err)
}

// TestBlockListEntry_DeleteNonexistent tests deleting entries that don't exist.
func TestBlockListEntry_DeleteNonexistent(t *testing.T) {
	client := getTestClient(t)
	ctx := context.Background()

	// Try to delete entries that likely don't exist
	entryIDs := []int{999999, 888888, 777777}
	result, err := client.DeleteBlockListEntry(ctx, entryIDs)
	if err == nil {
		t.Fatal("DeleteBlockListEntry with nonexistent IDs should return error")
	}
	t.Logf("Nonexistent entries error (expected): %v", err)

	// The result might contain error details
	if result != nil {
		t.Logf("Response: %s", result.String())
		if result.IsSuccessful() {
			t.Error("Should not be successful when entries don't exist")
		}
		t.Logf("Deleted count: %d", result.DeletedCount)
	}
}

// TestBlockListEntry_DeleteSingleEntry tests deleting a single entry.
func TestBlockListEntry_DeleteSingleEntry(t *testing.T) {
	client := getTestClient(t)
	ctx := context.Background()

	// Try to delete a single entry (will fail if it doesn't exist)
	result, err := client.DeleteBlockListEntry(ctx, []int{999999})
	if err != nil {
		t.Logf("Single entry deletion error (expected if not exists): %v", err)
	}

	if result != nil {
		t.Logf("Response: %s", result.String())
	}
}

// TestBlockListEntry_DeleteMultipleEntries tests deleting multiple entries.
func TestBlockListEntry_DeleteMultipleEntries(t *testing.T) {
	client := getTestClient(t)
	ctx := context.Background()

	// Try to delete multiple entries (will fail if they don't exist)
	entryIDs := []int{111111, 222222, 333333, 444444, 555555}
	result, err := client.DeleteBlockListEntry(ctx, entryIDs)
	if err != nil {
		t.Logf("Multiple entries deletion error (expected if not exists): %v", err)
	}

	if result != nil {
		t.Logf("Response: %s", result.String())
		t.Logf("Requested deletion of %d entries", len(entryIDs))
		if result.IsSuccessful() {
			t.Logf("Successfully deleted %d entries", result.DeletedCount)
		}
	}
}

// TestBlockListEntry_DeleteValidation tests ID validation logic.
func TestBlockListEntry_DeleteValidation(t *testing.T) {
	client := getTestClient(t)
	ctx := context.Background()

	testCases := []struct {
		name      string
		entryIDs  []int
		shouldErr bool
		reason    string
	}{
		{
			name:      "valid single ID",
			entryIDs:  []int{123},
			shouldErr: false, // Should pass validation
			reason:    "single positive ID",
		},
		{
			name:      "valid multiple IDs",
			entryIDs:  []int{123, 456, 789},
			shouldErr: false,
			reason:    "multiple positive unique IDs",
		},
		{
			name:      "invalid zero ID",
			entryIDs:  []int{123, 0},
			shouldErr: true,
			reason:    "contains zero",
		},
		{
			name:      "invalid negative ID",
			entryIDs:  []int{123, -5},
			shouldErr: true,
			reason:    "contains negative",
		},
		{
			name:      "invalid duplicate IDs",
			entryIDs:  []int{123, 456, 123},
			shouldErr: true,
			reason:    "contains duplicates",
		},
		{
			name:      "invalid empty list",
			entryIDs:  []int{},
			shouldErr: true,
			reason:    "empty list",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := client.DeleteBlockListEntry(ctx, tc.entryIDs)

			if tc.shouldErr {
				// Should get validation error
				if err == nil {
					t.Errorf("Expected validation error for %s, but got none", tc.reason)
				} else {
					t.Logf("Got expected validation error for %s: %v", tc.reason, err)
				}
			} else {
				// Should not get validation error (might get operation error)
				if err != nil {
					errStr := err.Error()
					// Check if it's a validation error
					if contains(errStr, "must be positive") || contains(errStr, "cannot be empty") || contains(errStr, "duplicate") {
						t.Errorf("Got unexpected validation error for %s: %v", tc.reason, err)
					} else {
						t.Logf("Got operation error (not validation) for %s: %v", tc.reason, err)
					}
				}
			}
		})
	}
}

// TestBlockList_DeleteBoundaries tests boundary conditions for block list IDs.
func TestBlockList_DeleteBoundaries(t *testing.T) {
	client := getTestClient(t)
	ctx := context.Background()

	testCases := []struct {
		id          int
		shouldPass  bool
		description string
	}{
		{0, false, "zero"},
		{1, true, "minimum valid"},
		{100, true, "normal ID"},
		{999999, true, "large ID"},
		{-1, false, "negative"},
		{-100, false, "large negative"},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			_, err := client.DeleteBlockList(ctx, tc.id)

			if tc.shouldPass {
				// Should pass validation (might fail on operation)
				if err != nil {
					errStr := err.Error()
					if contains(errStr, "must be positive") {
						t.Errorf("ID %d should pass validation but got: %v", tc.id, err)
					} else {
						t.Logf("ID %d passed validation, got operation error: %v", tc.id, err)
					}
				}
			} else {
				// Should fail validation
				if err == nil {
					t.Errorf("ID %d should fail validation but passed", tc.id)
				} else {
					t.Logf("ID %d correctly failed validation: %v", tc.id, err)
				}
			}
		})
	}
}
