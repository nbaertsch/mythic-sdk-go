//go:build integration

package integration

import (
	"context"
	"testing"
)

// TestGetStagingInfo_Success tests retrieving staging information.
func TestGetStagingInfo_Success(t *testing.T) {
	client := getTestClient(t)
	ctx := context.Background()

	// Get staging info for current operation
	stagingList, err := client.GetStagingInfo(ctx)
	if err != nil {
		t.Fatalf("GetStagingInfo failed: %v", err)
	}

	t.Logf("Found %d staging entries", len(stagingList))

	// Log details of each staging entry
	for i, staging := range stagingList {
		t.Logf("Staging %d: %s", i+1, staging.String())
		t.Logf("  ID: %d", staging.ID)
		t.Logf("  Payload ID: %d", staging.PayloadID)
		t.Logf("  Staging UUID: %s", staging.StagingUUID)
		t.Logf("  C2 Profile ID: %d", staging.C2ProfileID)
		t.Logf("  Active: %v", staging.IsActive())
		t.Logf("  Has Encryption: %v", staging.HasEncryption())
		t.Logf("  Has Expiration: %v", staging.HasExpiration())

		if staging.HasExpiration() {
			t.Logf("  Expiration Time: %s", staging.ExpirationTime)
		}

		// Validate required fields
		if staging.ID <= 0 {
			t.Errorf("Staging %d has invalid ID: %d", i+1, staging.ID)
		}
		if staging.StagingUUID == "" {
			t.Errorf("Staging %d has empty UUID", i+1)
		}
		if staging.PayloadID <= 0 {
			t.Errorf("Staging %d has invalid payload ID: %d", i+1, staging.PayloadID)
		}
		if staging.C2ProfileID <= 0 {
			t.Errorf("Staging %d has invalid C2 profile ID: %d", i+1, staging.C2ProfileID)
		}
	}
}

// TestGetStagingInfo_EmptyOperation tests staging info when no staging exists.
func TestGetStagingInfo_EmptyOperation(t *testing.T) {
	client := getTestClient(t)
	ctx := context.Background()

	// This might return empty list if no staging entries exist
	stagingList, err := client.GetStagingInfo(ctx)
	if err != nil {
		t.Fatalf("GetStagingInfo failed: %v", err)
	}

	t.Logf("Found %d staging entries (might be empty if no staging configured)", len(stagingList))

	// Empty list is valid - not an error
	if len(stagingList) == 0 {
		t.Log("No staging entries found (expected if operation has no staged payloads)")
	}
}

// TestGetStagingInfo_ValidateStructure tests the structure of returned staging info.
func TestGetStagingInfo_ValidateStructure(t *testing.T) {
	client := getTestClient(t)
	ctx := context.Background()

	stagingList, err := client.GetStagingInfo(ctx)
	if err != nil {
		t.Fatalf("GetStagingInfo failed: %v", err)
	}

	if len(stagingList) == 0 {
		t.Skip("No staging entries to validate (operation has no staged payloads)")
	}

	// Take first entry for validation
	staging := stagingList[0]

	// Test String() method
	str := staging.String()
	if str == "" {
		t.Error("String() returned empty string")
	}
	t.Logf("String representation: %s", str)

	// Test IsActive() method
	isActive := staging.IsActive()
	t.Logf("IsActive: %v (Active=%v, Deleted=%v)", isActive, staging.Active, staging.Deleted)

	// Test IsDeleted() method
	isDeleted := staging.IsDeleted()
	t.Logf("IsDeleted: %v", isDeleted)

	// Test HasEncryption() method
	hasEncryption := staging.HasEncryption()
	t.Logf("HasEncryption: %v", hasEncryption)
	if hasEncryption {
		if staging.EncryptionKey != "" {
			t.Logf("  Has encryption key (length: %d)", len(staging.EncryptionKey))
		}
		if staging.DecryptionKey != "" {
			t.Logf("  Has decryption key (length: %d)", len(staging.DecryptionKey))
		}
	}

	// Test HasExpiration() method
	hasExpiration := staging.HasExpiration()
	t.Logf("HasExpiration: %v", hasExpiration)
	if hasExpiration {
		t.Logf("  Expiration time: %s", staging.ExpirationTime)
	}

	// Validate timestamp format
	if staging.CreationTime == "" {
		t.Error("CreationTime is empty")
	}
	t.Logf("Creation time: %s", staging.CreationTime)
}

// TestGetStagingInfo_FilterActive tests that only active staging is returned.
func TestGetStagingInfo_FilterActive(t *testing.T) {
	client := getTestClient(t)
	ctx := context.Background()

	stagingList, err := client.GetStagingInfo(ctx)
	if err != nil {
		t.Fatalf("GetStagingInfo failed: %v", err)
	}

	if len(stagingList) == 0 {
		t.Skip("No staging entries to check")
	}

	// Check that all returned staging entries are not deleted
	for i, staging := range stagingList {
		if staging.IsDeleted() {
			t.Errorf("Staging %d (ID: %d) is marked as deleted but was returned", i+1, staging.ID)
		}
	}

	t.Logf("Verified %d staging entries are not deleted", len(stagingList))
}

// TestGetStagingInfo_UUIDFormat tests that staging UUIDs are present and valid.
func TestGetStagingInfo_UUIDFormat(t *testing.T) {
	client := getTestClient(t)
	ctx := context.Background()

	stagingList, err := client.GetStagingInfo(ctx)
	if err != nil {
		t.Fatalf("GetStagingInfo failed: %v", err)
	}

	if len(stagingList) == 0 {
		t.Skip("No staging entries to validate UUIDs")
	}

	// Check that all UUIDs are non-empty and have reasonable length
	for i, staging := range stagingList {
		if staging.StagingUUID == "" {
			t.Errorf("Staging %d (ID: %d) has empty UUID", i+1, staging.ID)
		}

		// UUIDs should be at least a few characters
		if len(staging.StagingUUID) < 10 {
			t.Errorf("Staging %d has suspiciously short UUID: %s", i+1, staging.StagingUUID)
		}

		t.Logf("Staging %d UUID: %s (length: %d)", i+1, staging.StagingUUID, len(staging.StagingUUID))
	}
}

// TestGetStagingInfo_Relationships tests that relationship IDs are valid.
func TestGetStagingInfo_Relationships(t *testing.T) {
	client := getTestClient(t)
	ctx := context.Background()

	stagingList, err := client.GetStagingInfo(ctx)
	if err != nil {
		t.Fatalf("GetStagingInfo failed: %v", err)
	}

	if len(stagingList) == 0 {
		t.Skip("No staging entries to validate relationships")
	}

	for i, staging := range stagingList {
		// Validate payload ID
		if staging.PayloadID <= 0 {
			t.Errorf("Staging %d has invalid payload ID: %d", i+1, staging.PayloadID)
		}

		// Validate C2 profile ID
		if staging.C2ProfileID <= 0 {
			t.Errorf("Staging %d has invalid C2 profile ID: %d", i+1, staging.C2ProfileID)
		}

		// Validate operation ID
		if staging.OperationID <= 0 {
			t.Errorf("Staging %d has invalid operation ID: %d", i+1, staging.OperationID)
		}

		// Validate operator ID
		if staging.OperatorID <= 0 {
			t.Errorf("Staging %d has invalid operator ID: %d", i+1, staging.OperatorID)
		}

		t.Logf("Staging %d relationships: Payload=%d, C2Profile=%d, Operation=%d, Operator=%d",
			i+1, staging.PayloadID, staging.C2ProfileID, staging.OperationID, staging.OperatorID)
	}
}
