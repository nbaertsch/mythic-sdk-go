//go:build integration

package integration

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/nbaertsch/mythic-sdk-go/pkg/mythic/types"
)

// TestE2E_PayloadLifecycle tests the complete payload build and deployment workflow
// Covers: GetPayloadTypes, GetPayloads, GetPayloadByUUID, CreatePayload, UpdatePayload,
// DeletePayload, DownloadPayload, ExportPayloadConfig, RebuildPayload, GetPayloadCommands,
// GetPayloadOnHost, GetBuildParametersByPayloadType, GetBuildParameterInstancesByPayload
func TestE2E_PayloadLifecycle(t *testing.T) {
	client := AuthenticateTestClient(t)

	var createdPayloadUUIDs []string
	var testPayloadUUID string
	var poseidonTypeID int

	// Register cleanup
	defer func() {
		// Delete all created payloads
		for _, uuid := range createdPayloadUUIDs {
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			_ = client.DeletePayload(ctx, uuid)
			cancel()
			t.Logf("Cleaned up payload: %s", uuid)
		}
	}()

	// Test 1: Get all payload types
	t.Log("=== Test 1: Get all payload types ===")
	ctx1, cancel1 := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel1()

	payloadTypes, err := client.GetPayloadTypes(ctx1)
	if err != nil {
		t.Fatalf("GetPayloadTypes failed: %v", err)
	}
	if len(payloadTypes) == 0 {
		t.Fatal("No payload types found")
	}
	t.Logf("✓ Found %d payload types", len(payloadTypes))

	// Test 2: Find Poseidon payload type
	t.Log("=== Test 2: Find Poseidon payload type ===")
	var poseidonType *types.PayloadType
	for _, pt := range payloadTypes {
		t.Logf("  - %s (ID: %d, Running: %v)", pt.Name, pt.ID, pt.ContainerRunning)
		if pt.Name == "poseidon" {
			poseidonType = pt
			break
		}
	}

	if poseidonType == nil {
		t.Skip("Poseidon payload type not found - skipping payload build tests")
	}
	poseidonTypeID = poseidonType.ID

	if !poseidonType.ContainerRunning {
		t.Skipf("Poseidon container not running (ID: %d) - skipping payload build tests", poseidonType.ID)
	}

	t.Logf("✓ Poseidon payload type found: ID %d (Container Running: %v)", poseidonType.ID, poseidonType.ContainerRunning)

	// Test 3: Get build parameters for Poseidon
	t.Log("=== Test 3: Get build parameters for Poseidon ===")
	ctx2, cancel2 := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel2()

	buildParams, err := client.GetBuildParametersByPayloadType(ctx2, poseidonTypeID)
	if err != nil {
		t.Logf("⚠ GetBuildParametersByPayloadType failed: %v", err)
	} else {
		t.Logf("✓ Found %d build parameters for Poseidon", len(buildParams))
		for _, param := range buildParams {
			t.Logf("  - %s: %s", param.Name, param.Description)
		}
	}

	// Test 4: Get baseline payload count
	t.Log("=== Test 4: Get baseline payload count ===")
	ctx3, cancel3 := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel3()

	baselinePayloads, err := client.GetPayloads(ctx3)
	if err != nil {
		t.Fatalf("GetPayloads failed: %v", err)
	}
	baselineCount := len(baselinePayloads)
	t.Logf("✓ Baseline payload count: %d", baselineCount)

	// Test 5: Get C2 profiles for payload creation
	t.Log("=== Test 5: Get C2 profiles for payload ===")
	ctx4, cancel4 := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel4()

	c2Profiles, err := client.GetC2Profiles(ctx4)
	if err != nil {
		t.Fatalf("GetC2Profiles failed: %v", err)
	}
	if len(c2Profiles) == 0 {
		t.Fatal("No C2 profiles available")
	}

	// Find HTTP profile or use first available
	var c2Profile *types.C2Profile
	for _, profile := range c2Profiles {
		if profile.Name == "http" {
			c2Profile = profile
			break
		}
	}
	if c2Profile == nil {
		c2Profile = c2Profiles[0]
	}
	t.Logf("✓ Using C2 profile: %s (ID: %d)", c2Profile.Name, c2Profile.ID)

	// Test 6: Create payload build request
	t.Log("=== Test 6: Create Poseidon payload ===")
	ctx5, cancel5 := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel5()

	payloadReq := &types.CreatePayloadRequest{
		PayloadType: "poseidon",
		OS:          "linux",
		Description: "E2E Test Payload - Poseidon Linux",
		Filename:    "e2e_test_poseidon",
		Commands: []string{
			"shell", "download", "upload", "ps", "ls",
		},
		C2Profiles: []types.C2ProfileConfig{
			{
				Name: c2Profile.Name,
				Parameters: map[string]interface{}{
					"callback_host":     "127.0.0.1",
					"callback_port":     "7443",
					"callback_interval": "10",
				},
			},
		},
		BuildParameters: map[string]interface{}{},
	}

	newPayload, err := client.CreatePayload(ctx5, payloadReq)
	if err != nil {
		t.Fatalf("CreatePayload failed: %v", err)
	}
	if newPayload.UUID == "" {
		t.Fatal("Created payload has empty UUID")
	}
	testPayloadUUID = newPayload.UUID
	createdPayloadUUIDs = append(createdPayloadUUIDs, testPayloadUUID)
	t.Logf("✓ Payload created: UUID %s (Build Phase: %s)", newPayload.UUID, newPayload.BuildPhase)

	// Test 7: Wait for payload build to complete
	t.Log("=== Test 7: Wait for payload build ===")
	ctx6, cancel6 := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel6()

	err = client.WaitForPayloadComplete(ctx6, testPayloadUUID, 90)
	if err != nil {
		t.Fatalf("WaitForPayloadComplete failed: %v", err)
	}
	t.Logf("✓ Payload build completed")

	// Test 8: Get payload by UUID and verify status
	t.Log("=== Test 8: Get payload by UUID ===")
	ctx7, cancel7 := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel7()

	builtPayload, err := client.GetPayloadByUUID(ctx7, testPayloadUUID)
	if err != nil {
		t.Fatalf("GetPayloadByUUID failed: %v", err)
	}
	if builtPayload.UUID != testPayloadUUID {
		t.Errorf("Payload UUID mismatch: expected %s, got %s", testPayloadUUID, builtPayload.UUID)
	}
	if !builtPayload.IsReady() {
		t.Errorf("Payload not ready: BuildPhase=%s, BuildMessage=%s", builtPayload.BuildPhase, builtPayload.BuildMessage)
		if builtPayload.BuildStderr != "" {
			t.Logf("Build stderr: %s", builtPayload.BuildStderr)
		}
	}
	t.Logf("✓ Payload retrieved: %s (Status: %s)", builtPayload.UUID, builtPayload.BuildPhase)

	// Test 9: Verify payload appears in payload list
	t.Log("=== Test 9: Verify payload in list ===")
	ctx8, cancel8 := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel8()

	allPayloads, err := client.GetPayloads(ctx8)
	if err != nil {
		t.Fatalf("GetPayloads after creation failed: %v", err)
	}
	found := false
	for _, p := range allPayloads {
		if p.UUID == testPayloadUUID {
			found = true
			t.Logf("✓ Payload found in list: %s", p.UUID)
			break
		}
	}
	if !found {
		t.Error("Created payload not found in payload list")
	}
	newCount := len(allPayloads)
	t.Logf("✓ Total payloads now: %d (baseline: %d)", newCount, baselineCount)

	// Test 10: Download payload binary
	t.Log("=== Test 10: Download payload binary ===")
	ctx9, cancel9 := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel9()

	payloadBytes, err := client.DownloadPayload(ctx9, testPayloadUUID)
	if err != nil {
		t.Fatalf("DownloadPayload failed: %v", err)
	}
	if len(payloadBytes) == 0 {
		t.Fatal("Downloaded payload is empty")
	}
	t.Logf("✓ Payload downloaded: %d bytes", len(payloadBytes))

	// Verify it's a valid binary (check for ELF header on Linux)
	if len(payloadBytes) >= 4 {
		if payloadBytes[0] == 0x7f && payloadBytes[1] == 'E' && payloadBytes[2] == 'L' && payloadBytes[3] == 'F' {
			t.Log("✓ Downloaded file is a valid ELF binary")
		} else {
			t.Logf("⚠ Downloaded file may not be ELF (first 4 bytes: %x %x %x %x)", payloadBytes[0], payloadBytes[1], payloadBytes[2], payloadBytes[3])
		}
	}

	// Test 11: Export payload config
	t.Log("=== Test 11: Export payload config ===")
	ctx10, cancel10 := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel10()

	configJSON, err := client.ExportPayloadConfig(ctx10, testPayloadUUID)
	if err != nil {
		t.Logf("⚠ ExportPayloadConfig failed (may not be supported): %v", err)
	} else {
		if configJSON == "" {
			t.Error("Exported config is empty")
		} else {
			// Verify it's valid JSON
			var configData interface{}
			if err := json.Unmarshal([]byte(configJSON), &configData); err != nil {
				t.Errorf("Exported config is not valid JSON: %v", err)
			} else {
				t.Logf("✓ Payload config exported: %d bytes", len(configJSON))
			}
		}
	}

	// Test 12: Update payload description
	t.Log("=== Test 12: Update payload description ===")
	ctx11, cancel11 := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel11()

	newDesc := "Updated: E2E Test Payload - Modified Description"
	updateReq := &types.UpdatePayloadRequest{
		UUID:        testPayloadUUID,
		Description: &newDesc,
	}

	updatedPayload, err := client.UpdatePayload(ctx11, updateReq)
	if err != nil {
		t.Fatalf("UpdatePayload failed: %v", err)
	}
	if updatedPayload.Description != newDesc {
		t.Errorf("Description not updated: expected %s, got %s", newDesc, updatedPayload.Description)
	}
	t.Logf("✓ Payload description updated")

	// Test 13: Get payload commands
	t.Log("=== Test 13: Get payload commands ===")
	ctx12, cancel12 := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel12()

	commands, err := client.GetPayloadCommands(ctx12, builtPayload.ID)
	if err != nil {
		t.Logf("⚠ GetPayloadCommands failed: %v", err)
	} else {
		t.Logf("✓ Payload has %d commands", len(commands))
		for _, cmd := range commands {
			t.Logf("  - %s", cmd)
		}
	}

	// Test 14: Get build parameter instances for payload
	t.Log("=== Test 14: Get build parameter instances ===")
	ctx13, cancel13 := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel13()

	buildParamInstances, err := client.GetBuildParameterInstancesByPayload(ctx13, builtPayload.ID)
	if err != nil {
		t.Logf("⚠ GetBuildParameterInstancesByPayload failed: %v", err)
	} else {
		t.Logf("✓ Found %d build parameter instances", len(buildParamInstances))
	}

	// Test 15: Get payloads on host (should be empty - no deployments yet)
	t.Log("=== Test 15: Get payloads on host ===")
	ctx14, cancel14 := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel14()

	operationID := client.GetCurrentOperation()
	if operationID == nil {
		t.Fatal("No current operation set")
	}

	payloadsOnHost, err := client.GetPayloadOnHost(ctx14, *operationID)
	if err != nil {
		t.Logf("⚠ GetPayloadOnHost failed: %v", err)
	} else {
		t.Logf("✓ Found %d payloads on host for operation %d", len(payloadsOnHost), *operationID)
	}

	// Test 16: Create second payload for rebuild testing
	t.Log("=== Test 16: Create second payload for rebuild test ===")
	ctx15, cancel15 := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel15()

	rebuildReq := &types.CreatePayloadRequest{
		PayloadType: "poseidon",
		OS:          "linux",
		Description: "E2E Test Payload - For Rebuild Testing",
		Filename:    "e2e_test_rebuild",
		Commands:    []string{"shell", "ps"},
		C2Profiles: []types.C2ProfileConfig{
			{
				Name: c2Profile.Name,
				Parameters: map[string]interface{}{
					"callback_host": "127.0.0.1",
					"callback_port": "7443",
				},
			},
		},
	}

	rebuildPayload, err := client.CreatePayload(ctx15, rebuildReq)
	if err != nil {
		t.Fatalf("CreatePayload (for rebuild) failed: %v", err)
	}
	createdPayloadUUIDs = append(createdPayloadUUIDs, rebuildPayload.UUID)
	t.Logf("✓ Second payload created: UUID %s", rebuildPayload.UUID)

	// Test 17: Wait for second payload build
	t.Log("=== Test 17: Wait for second payload build ===")
	ctx16, cancel16 := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel16()

	err = client.WaitForPayloadComplete(ctx16, rebuildPayload.UUID, 90)
	if err != nil {
		t.Logf("⚠ Second payload build failed: %v", err)
	} else {
		t.Log("✓ Second payload build completed")
	}

	// Test 18: Rebuild the second payload
	t.Log("=== Test 18: Rebuild payload ===")
	ctx17, cancel17 := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel17()

	rebuiltPayload, err := client.RebuildPayload(ctx17, rebuildPayload.UUID)
	if err != nil {
		t.Logf("⚠ RebuildPayload failed (may require fully built payload): %v", err)
	} else {
		t.Logf("✓ Payload rebuild initiated: UUID %s (Build Phase: %s)", rebuiltPayload.UUID, rebuiltPayload.BuildPhase)

		// Test 19: Wait for rebuild to complete
		t.Log("=== Test 19: Wait for rebuild completion ===")
		ctx18, cancel18 := context.WithTimeout(context.Background(), 120*time.Second)
		defer cancel18()

		err = client.WaitForPayloadComplete(ctx18, rebuiltPayload.UUID, 90)
		if err != nil {
			t.Logf("⚠ Rebuild wait failed: %v", err)
		} else {
			t.Log("✓ Rebuild completed")
		}
	}

	// Test 20: Delete second payload
	t.Log("=== Test 20: Delete second payload ===")
	ctx19, cancel19 := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel19()

	err = client.DeletePayload(ctx19, rebuildPayload.UUID)
	if err != nil {
		t.Errorf("DeletePayload failed: %v", err)
	} else {
		t.Logf("✓ Payload %s deleted", rebuildPayload.UUID)
		// Remove from cleanup list since we deleted it
		for i, uuid := range createdPayloadUUIDs {
			if uuid == rebuildPayload.UUID {
				createdPayloadUUIDs = append(createdPayloadUUIDs[:i], createdPayloadUUIDs[i+1:]...)
				break
			}
		}
	}

	// Test 21: Verify deletion
	t.Log("=== Test 21: Verify payload deletion ===")
	ctx20, cancel20 := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel20()

	finalPayloads, err := client.GetPayloads(ctx20)
	if err != nil {
		t.Fatalf("GetPayloads after deletion failed: %v", err)
	}

	for _, p := range finalPayloads {
		if p.UUID == rebuildPayload.UUID {
			if p.Deleted {
				t.Logf("✓ Deleted payload marked as deleted in list")
			} else {
				t.Error("Deleted payload still active in list")
			}
			break
		}
	}
	t.Logf("✓ Final payload count: %d", len(finalPayloads))

	t.Log("=== ✓ All payload lifecycle tests passed ===")
}

// TestE2E_PayloadErrorHandling tests error scenarios for payload operations
func TestE2E_PayloadErrorHandling(t *testing.T) {
	client := AuthenticateTestClient(t)

	// Test 1: Get non-existent payload by UUID
	t.Log("=== Test 1: Get non-existent payload ===")
	ctx1, cancel1 := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel1()

	_, err := client.GetPayloadByUUID(ctx1, "00000000-0000-0000-0000-000000000000")
	if err == nil {
		t.Error("Expected error for non-existent payload UUID")
	}
	t.Logf("✓ Non-existent payload rejected: %v", err)

	// Test 2: Download non-existent payload
	t.Log("=== Test 2: Download non-existent payload ===")
	ctx2, cancel2 := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel2()

	_, err = client.DownloadPayload(ctx2, "00000000-0000-0000-0000-000000000000")
	if err == nil {
		t.Error("Expected error for non-existent payload download")
	}
	t.Logf("✓ Non-existent payload download rejected: %v", err)

	// Test 3: Delete non-existent payload
	t.Log("=== Test 3: Delete non-existent payload ===")
	ctx3, cancel3 := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel3()

	err = client.DeletePayload(ctx3, "00000000-0000-0000-0000-000000000000")
	if err == nil {
		t.Error("Expected error for non-existent payload deletion")
	}
	t.Logf("✓ Non-existent payload deletion rejected: %v", err)

	// Test 4: Create payload with invalid payload type
	t.Log("=== Test 4: Create payload with invalid type ===")
	ctx4, cancel4 := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel4()

	invalidReq := &types.CreatePayloadRequest{
		PayloadType: "invalid_payload_type_xyz",
		OS:          "linux",
	}

	_, err = client.CreatePayload(ctx4, invalidReq)
	if err == nil {
		t.Error("Expected error for invalid payload type")
	}
	t.Logf("✓ Invalid payload type rejected: %v", err)

	t.Log("=== ✓ All error handling tests passed ===")
}
