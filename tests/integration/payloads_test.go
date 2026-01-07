//go:build integration

package integration

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/nbaertsch/mythic-sdk-go/pkg/mythic/types"
)

func TestPayloads_GetPayloadTypes(t *testing.T) {
	SkipIfNoMythic(t)

	client := AuthenticateTestClient(t)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	payloadTypes, err := client.GetPayloadTypes(ctx)
	if err != nil {
		t.Fatalf("GetPayloadTypes failed: %v", err)
	}

	if payloadTypes == nil {
		t.Fatal("GetPayloadTypes returned nil")
	}

	if len(payloadTypes) == 0 {
		t.Skip("No payload types available in test Mythic instance")
	}

	// Verify payload type structure
	pt := payloadTypes[0]
	if pt.ID == 0 {
		t.Error("Payload type ID should not be 0")
	}
	if pt.Name == "" {
		t.Error("Payload type name should not be empty")
	}

	t.Logf("Found %d payload type(s)", len(payloadTypes))
	for _, payloadType := range payloadTypes {
		t.Logf("  - %s", payloadType.String())
	}
}

func TestPayloads_GetPayloads(t *testing.T) {
	SkipIfNoMythic(t)

	client := AuthenticateTestClient(t)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	payloads, err := client.GetPayloads(ctx)
	if err != nil {
		t.Fatalf("GetPayloads failed: %v", err)
	}

	if payloads == nil {
		t.Fatal("GetPayloads returned nil")
	}

	t.Logf("Found %d payload(s)", len(payloads))

	// If there are payloads, verify structure
	if len(payloads) > 0 {
		p := payloads[0]
		if p.ID == 0 {
			t.Error("Payload ID should not be 0")
		}
		if p.UUID == "" {
			t.Error("Payload UUID should not be empty")
		}
		t.Logf("First payload: %s", p.String())
	}
}

func TestPayloads_CreateAndManagePayload(t *testing.T) {
	SkipIfNoMythic(t)

	client := AuthenticateTestClient(t)

	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Second) // Longer timeout for build
	defer cancel()

	// Get available payload types
	payloadTypes, err := client.GetPayloadTypes(ctx)
	if err != nil {
		t.Fatalf("GetPayloadTypes failed: %v", err)
	}

	if len(payloadTypes) == 0 {
		t.Skip("No payload types available for testing")
	}

	// Find a supported payload type
	var selectedType *types.PayloadType
	for _, pt := range payloadTypes {
		if pt.Supported && pt.ContainerRunning {
			selectedType = pt
			break
		}
	}

	if selectedType == nil {
		t.Skip("No supported and running payload types available")
	}

	t.Logf("Using payload type: %s", selectedType.Name)

	// For now, skip getting commands since GetPayloadCommands requires a payload ID
	// We'll just create a payload without specifying commands
	var selectedCommands []string

	// Create payload
	createReq := &types.CreatePayloadRequest{
		PayloadType: selectedType.Name,
		Description: "Integration test payload - " + time.Now().Format("2006-01-02 15:04:05"),
		Commands:    selectedCommands,
	}

	payload, err := client.CreatePayload(ctx, createReq)
	if err != nil {
		t.Fatalf("CreatePayload failed: %v", err)
	}

	if payload == nil {
		t.Fatal("CreatePayload returned nil")
	}

	if payload.UUID == "" {
		t.Error("Created payload should have a UUID")
	}

	t.Logf("Created payload: %s (UUID: %s)", payload.String(), payload.UUID)

	// Wait for payload to complete building (with timeout)
	t.Log("Waiting for payload to build...")
	err = client.WaitForPayloadComplete(ctx, payload.UUID, 120) // 120 seconds timeout
	if err != nil {
		t.Logf("WaitForPayloadComplete error (may still be building): %v", err)
	}

	// Get final payload status
	finalPayload, err := client.GetPayloadByUUID(ctx, payload.UUID)
	if err != nil {
		t.Fatalf("GetPayloadByUUID failed: %v", err)
	}

	t.Logf("Payload build status: %s (Phase: %s)", finalPayload.String(), finalPayload.BuildPhase)

	if finalPayload.IsFailed() {
		t.Logf("Build failed with message: %s", finalPayload.BuildMessage)
		t.Logf("Build stderr: %s", finalPayload.BuildStderr)
	}

	// Get payload by UUID
	retrieved, err := client.GetPayloadByUUID(ctx, payload.UUID)
	if err != nil {
		t.Fatalf("GetPayloadByUUID failed: %v", err)
	}

	if retrieved.UUID != payload.UUID {
		t.Errorf("Expected payload UUID %s, got %s", payload.UUID, retrieved.UUID)
	}

	// Update payload
	newDescription := "Updated test payload - " + time.Now().Format("15:04:05")
	callbackAlert := true

	updateReq := &types.UpdatePayloadRequest{
		UUID:          payload.UUID,
		Description:   &newDescription,
		CallbackAlert: &callbackAlert,
	}

	updated, err := client.UpdatePayload(ctx, updateReq)
	if err != nil {
		t.Fatalf("UpdatePayload failed: %v", err)
	}

	if updated.Description != newDescription {
		t.Errorf("Expected description %q, got %q", newDescription, updated.Description)
	}

	if updated.CallbackAlert != callbackAlert {
		t.Errorf("Expected callback alert %v, got %v", callbackAlert, updated.CallbackAlert)
	}

	t.Logf("Updated payload description to: %s", newDescription)

	// Export payload config
	config, err := client.ExportPayloadConfig(ctx, payload.UUID)
	if err != nil {
		t.Fatalf("ExportPayloadConfig failed: %v", err)
	}

	if config == "" {
		t.Error("Exported config should not be empty")
	}

	t.Logf("Exported config length: %d bytes", len(config))

	// Download payload if it's ready
	if finalPayload.IsReady() {
		t.Log("Attempting to download payload...")
		data, err := client.DownloadPayload(ctx, payload.UUID)
		if err != nil {
			t.Logf("Warning: DownloadPayload failed: %v", err)
		} else {
			if len(data) == 0 {
				t.Error("Downloaded payload should not be empty")
			} else {
				t.Logf("Downloaded payload size: %d bytes", len(data))
			}
		}
	}

	// Cleanup: Delete the payload
	deleted := true
	cleanupReq := &types.UpdatePayloadRequest{
		UUID:    payload.UUID,
		Deleted: &deleted,
	}

	_, err = client.UpdatePayload(ctx, cleanupReq)
	if err != nil {
		t.Logf("Warning: Failed to delete test payload: %v", err)
	} else {
		t.Log("Cleaned up test payload")
	}
}

func TestPayloads_GetPayloadByUUID_NotFound(t *testing.T) {
	SkipIfNoMythic(t)

	client := AuthenticateTestClient(t)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Try to get non-existent payload
	_, err := client.GetPayloadByUUID(ctx, "00000000-0000-0000-0000-000000000000")
	if err == nil {
		t.Fatal("Expected error for non-existent payload, got nil")
	}

	t.Logf("Expected error: %v", err)
}

func TestPayloads_GetPayloadCommands(t *testing.T) {
	SkipIfNoMythic(t)

	client := AuthenticateTestClient(t)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Get existing payloads
	payloads, err := client.GetPayloads(ctx)
	if err != nil {
		t.Fatalf("GetPayloads failed: %v", err)
	}

	if len(payloads) == 0 {
		t.Skip("No payloads available for testing")
	}

	// Get commands for first payload
	commands, err := client.GetPayloadCommands(ctx, payloads[0].ID)
	if err != nil {
		t.Fatalf("GetPayloadCommands failed: %v", err)
	}

	if commands == nil {
		t.Fatal("GetPayloadCommands returned nil")
	}

	t.Logf("Payload %d has %d command(s)", payloads[0].ID, len(commands))

	// Verify commands if any exist
	if len(commands) > 0 {
		t.Logf("First command: %s", commands[0])
	}
}

func TestPayloads_GetPayloadOnHost(t *testing.T) {
	SkipIfNoMythic(t)

	client := AuthenticateTestClient(t)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Get current operation
	currentOpID := client.GetCurrentOperation()
	if currentOpID == nil {
		t.Skip("No current operation set")
	}

	// This might return empty results if no payloads have been tracked on hosts
	payloadsOnHost, err := client.GetPayloadOnHost(ctx, *currentOpID)
	if err != nil {
		t.Fatalf("GetPayloadOnHost failed: %v", err)
	}

	if payloadsOnHost == nil {
		t.Fatal("GetPayloadOnHost returned nil")
	}

	t.Logf("Found %d payload(s) on host(s) for operation %d", len(payloadsOnHost), *currentOpID)

	// If there are results, verify structure
	if len(payloadsOnHost) > 0 {
		poh := payloadsOnHost[0]
		if poh.ID == 0 {
			t.Error("PayloadOnHost ID should not be 0")
		}
		if poh.PayloadID == 0 {
			t.Error("PayloadID should not be 0")
		}
		t.Logf("First entry: %s", poh.String())
	}
}

func TestPayloads_RebuildPayload(t *testing.T) {
	SkipIfNoMythic(t)

	client := AuthenticateTestClient(t)

	ctx, cancel := context.WithTimeout(context.Background(), 180*time.Second)
	defer cancel()

	// First, get existing payloads
	payloads, err := client.GetPayloads(ctx)
	if err != nil {
		t.Fatalf("GetPayloads failed: %v", err)
	}

	// Find a payload that's already built
	var targetPayload *types.Payload
	for _, p := range payloads {
		if p.IsReady() && !p.Deleted {
			targetPayload = p
			break
		}
	}

	if targetPayload == nil {
		t.Skip("No suitable payload available for rebuild test")
	}

	t.Logf("Rebuilding payload: %s", targetPayload.UUID)

	// Rebuild the payload
	rebuiltPayload, err := client.RebuildPayload(ctx, targetPayload.UUID)
	if err != nil {
		t.Fatalf("RebuildPayload failed: %v", err)
	}

	if rebuiltPayload == nil {
		t.Fatal("RebuildPayload returned nil")
	}

	if rebuiltPayload.UUID == "" {
		t.Error("Rebuilt payload should have a UUID")
	}

	t.Logf("Rebuilt payload: %s (UUID: %s)", rebuiltPayload.String(), rebuiltPayload.UUID)

	// Wait a bit for rebuild to start
	time.Sleep(2 * time.Second)

	// Check status
	status, err := client.GetPayloadByUUID(ctx, rebuiltPayload.UUID)
	if err != nil {
		t.Fatalf("GetPayloadByUUID failed after rebuild: %v", err)
	}

	t.Logf("Rebuild status: %s", status.BuildPhase)
}

func TestPayloads_DeletePayload(t *testing.T) {
	SkipIfNoMythic(t)

	client := AuthenticateTestClient(t)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Get payload types
	payloadTypes, err := client.GetPayloadTypes(ctx)
	if err != nil || len(payloadTypes) == 0 {
		t.Skip("Cannot create test payload for deletion")
	}

	// Find a supported type
	var selectedType *types.PayloadType
	for _, pt := range payloadTypes {
		if pt.Supported && pt.ContainerRunning {
			selectedType = pt
			break
		}
	}

	if selectedType == nil {
		t.Skip("No supported payload types available")
	}

	// Create a simple payload to delete
	createReq := &types.CreatePayloadRequest{
		PayloadType: selectedType.Name,
		Description: "Test payload for deletion",
	}

	payload, err := client.CreatePayload(ctx, createReq)
	if err != nil {
		t.Fatalf("CreatePayload failed: %v", err)
	}

	t.Logf("Created test payload for deletion: %s", payload.UUID)

	// Delete it using UpdatePayload
	deleted := true
	deleteReq := &types.UpdatePayloadRequest{
		UUID:    payload.UUID,
		Deleted: &deleted,
	}

	result, err := client.UpdatePayload(ctx, deleteReq)
	if err != nil {
		t.Fatalf("Delete (UpdatePayload) failed: %v", err)
	}

	if !result.Deleted {
		t.Error("Payload should be marked as deleted")
	}

	t.Logf("Successfully deleted payload: %s", payload.UUID)

	// Verify using DeletePayload method
	err = client.DeletePayload(ctx, payload.UUID)
	if err != nil {
		t.Logf("DeletePayload call: %v", err)
	}
}

func TestPayloads_CreatePayload_InvalidRequest(t *testing.T) {
	SkipIfNoMythic(t)

	client := AuthenticateTestClient(t)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Test with empty payload type
	req := &types.CreatePayloadRequest{
		PayloadType: "",
	}

	_, err := client.CreatePayload(ctx, req)
	if err == nil {
		t.Fatal("Expected error for empty payload type, got nil")
	}

	t.Logf("Expected error: %v", err)

	// Test with nil request
	_, err = client.CreatePayload(ctx, nil)
	if err == nil {
		t.Fatal("Expected error for nil request, got nil")
	}

	t.Logf("Expected error for nil: %v", err)
}

func TestPayloads_UpdatePayload_InvalidRequest(t *testing.T) {
	SkipIfNoMythic(t)

	client := AuthenticateTestClient(t)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Test with empty UUID
	req := &types.UpdatePayloadRequest{
		UUID: "",
	}

	_, err := client.UpdatePayload(ctx, req)
	if err == nil {
		t.Fatal("Expected error for empty UUID, got nil")
	}

	t.Logf("Expected error: %v", err)

	// Test with nil request
	_, err = client.UpdatePayload(ctx, nil)
	if err == nil {
		t.Fatal("Expected error for nil request, got nil")
	}

	t.Logf("Expected error for nil: %v", err)
}

func TestPayloads_DownloadPayload_NotReady(t *testing.T) {
	SkipIfNoMythic(t)

	client := AuthenticateTestClient(t)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Try to download a non-existent payload
	_, err := client.DownloadPayload(ctx, "00000000-0000-0000-0000-000000000000")
	if err == nil {
		t.Fatal("Expected error downloading non-existent payload, got nil")
	}

	t.Logf("Expected error: %v", err)
}

func TestPayloads_ExportPayloadConfig(t *testing.T) {
	SkipIfNoMythic(t)

	client := AuthenticateTestClient(t)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Get existing payloads
	payloads, err := client.GetPayloads(ctx)
	if err != nil {
		t.Fatalf("GetPayloads failed: %v", err)
	}

	if len(payloads) == 0 {
		t.Skip("No payloads available for export test")
	}

	// Export config for first payload
	config, err := client.ExportPayloadConfig(ctx, payloads[0].UUID)
	if err != nil {
		t.Fatalf("ExportPayloadConfig failed: %v", err)
	}

	if config == "" {
		t.Error("Exported config should not be empty")
	}

	t.Logf("Exported config for payload %s: %d bytes", payloads[0].UUID, len(config))
}

func TestPayloads_WaitForPayloadComplete_Timeout(t *testing.T) {
	SkipIfNoMythic(t)

	client := AuthenticateTestClient(t)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Try to wait for a non-existent payload (should error immediately)
	err := client.WaitForPayloadComplete(ctx, "00000000-0000-0000-0000-000000000000", 5) // 5 seconds timeout
	if err == nil {
		t.Fatal("Expected error waiting for non-existent payload, got nil")
	}

	t.Logf("Expected error: %v", err)
}

func TestPayloads_DownloadToFile(t *testing.T) {
	SkipIfNoMythic(t)

	client := AuthenticateTestClient(t)

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	// Get existing ready payloads
	payloads, err := client.GetPayloads(ctx)
	if err != nil {
		t.Fatalf("GetPayloads failed: %v", err)
	}

	// Find a ready payload
	var readyPayload *types.Payload
	for _, p := range payloads {
		if p.IsReady() && !p.Deleted {
			readyPayload = p
			break
		}
	}

	if readyPayload == nil {
		t.Skip("No ready payloads available for download test")
	}

	// Download to memory first
	data, err := client.DownloadPayload(ctx, readyPayload.UUID)
	if err != nil {
		t.Fatalf("DownloadPayload failed: %v", err)
	}

	if len(data) == 0 {
		t.Error("Downloaded payload should not be empty")
	}

	t.Logf("Downloaded payload %s: %d bytes", readyPayload.UUID, len(data))

	// Test writing to a temporary file
	tmpfile := "/tmp/test-payload-" + time.Now().Format("20060102-150405")
	defer os.Remove(tmpfile)

	err = os.WriteFile(tmpfile, data, 0644)
	if err != nil {
		t.Fatalf("Failed to write payload to file: %v", err)
	}

	// Verify file size
	stat, err := os.Stat(tmpfile)
	if err != nil {
		t.Fatalf("Failed to stat downloaded file: %v", err)
	}

	if stat.Size() != int64(len(data)) {
		t.Errorf("File size mismatch: expected %d, got %d", len(data), stat.Size())
	}

	t.Logf("Successfully wrote payload to file: %s (%d bytes)", tmpfile, stat.Size())
}
