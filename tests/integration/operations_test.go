//go:build integration

package integration

import (
	"context"
	"testing"
	"time"

	"github.com/nbaertsch/mythic-sdk-go/pkg/mythic/types"
)

func TestOperations_GetOperations(t *testing.T) {
	SkipIfNoMythic(t)

	client := AuthenticateTestClient(t)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	operations, err := client.GetOperations(ctx)
	if err != nil {
		t.Fatalf("GetOperations failed: %v", err)
	}

	if operations == nil {
		t.Fatal("GetOperations returned nil")
	}

	if len(operations) == 0 {
		t.Fatal("Expected at least one operation (default operation)")
	}

	// Verify operation structure
	op := operations[0]
	if op.ID == 0 {
		t.Error("Operation ID should not be 0")
	}
	if op.Name == "" {
		t.Error("Operation name should not be empty")
	}
	if op.AdminID == 0 {
		t.Error("Operation should have an admin")
	}

	t.Logf("Found %d operation(s)", len(operations))
	t.Logf("First operation: %s (ID: %d, Admin: %d)", op.Name, op.ID, op.AdminID)
}

func TestOperations_GetOperationByID(t *testing.T) {
	SkipIfNoMythic(t)

	client := AuthenticateTestClient(t)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Get all operations to find a valid ID
	operations, err := client.GetOperations(ctx)
	if err != nil {
		t.Fatalf("GetOperations failed: %v", err)
	}

	if len(operations) == 0 {
		t.Skip("No operations available for testing")
	}

	// Get specific operation
	op, err := client.GetOperationByID(ctx, operations[0].ID)
	if err != nil {
		t.Fatalf("GetOperationByID failed: %v", err)
	}

	if op == nil {
		t.Fatal("GetOperationByID returned nil")
	}

	if op.ID != operations[0].ID {
		t.Errorf("Expected operation ID %d, got %d", operations[0].ID, op.ID)
	}

	if op.Name == "" {
		t.Error("Operation name should not be empty")
	}

	t.Logf("Retrieved operation: %s", op.String())
}

func TestOperations_GetOperationByID_NotFound(t *testing.T) {
	SkipIfNoMythic(t)

	client := AuthenticateTestClient(t)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Try to get non-existent operation
	_, err := client.GetOperationByID(ctx, 999999)
	if err == nil {
		t.Fatal("Expected error for non-existent operation, got nil")
	}

	t.Logf("Expected error: %v", err)
}

func TestOperations_CreateAndUpdateOperation(t *testing.T) {
	SkipIfNoMythic(t)

	client := AuthenticateTestClient(t)

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	// Get current user to use as admin
	me, err := client.GetMe(ctx)
	if err != nil {
		t.Fatalf("GetMe failed: %v", err)
	}

	// Create new operation
	createReq := &types.CreateOperationRequest{
		Name:    "Test Operation " + time.Now().Format("20060102-150405"),
		AdminID: &me.ID,
		Channel: "#test-channel",
		Webhook: "https://hooks.example.com/test",
	}

	op, err := client.CreateOperation(ctx, createReq)
	if err != nil {
		t.Fatalf("CreateOperation failed: %v", err)
	}

	if op == nil {
		t.Fatal("CreateOperation returned nil")
	}

	if op.ID == 0 {
		t.Error("Created operation should have an ID")
	}

	if op.Name != createReq.Name {
		t.Errorf("Expected operation name %q, got %q", createReq.Name, op.Name)
	}

	t.Logf("Created operation: %s (ID: %d)", op.Name, op.ID)

	// Update the operation
	newName := "Updated Test Operation " + time.Now().Format("150405")
	complete := false
	bannerText := "Test Banner"
	bannerColor := "#00ff00"

	updateReq := &types.UpdateOperationRequest{
		OperationID: op.ID,
		Name:        &newName,
		Complete:    &complete,
		BannerText:  &bannerText,
		BannerColor: &bannerColor,
	}

	updatedOp, err := client.UpdateOperation(ctx, updateReq)
	if err != nil {
		t.Fatalf("UpdateOperation failed: %v", err)
	}

	if updatedOp.Name != newName {
		t.Errorf("Expected updated name %q, got %q", newName, updatedOp.Name)
	}

	if updatedOp.BannerText != bannerText {
		t.Errorf("Expected banner text %q, got %q", bannerText, updatedOp.BannerText)
	}

	if updatedOp.BannerColor != bannerColor {
		t.Errorf("Expected banner color %q, got %q", bannerColor, updatedOp.BannerColor)
	}

	t.Logf("Updated operation: %s", updatedOp.String())

	// Cleanup: Mark as complete
	completeTrue := true
	cleanupReq := &types.UpdateOperationRequest{
		OperationID: op.ID,
		Complete:    &completeTrue,
	}
	_, err = client.UpdateOperation(ctx, cleanupReq)
	if err != nil {
		t.Logf("Warning: Failed to mark test operation as complete: %v", err)
	}
}

func TestOperations_UpdateCurrentOperation(t *testing.T) {
	SkipIfNoMythic(t)

	client := AuthenticateTestClient(t)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Get current user
	me, err := client.GetMe(ctx)
	if err != nil {
		t.Fatalf("GetMe failed: %v", err)
	}

	// Store original operation ID
	originalOpID := me.CurrentOperationID

	// Get available operations
	operations, err := client.GetOperations(ctx)
	if err != nil {
		t.Fatalf("GetOperations failed: %v", err)
	}

	if len(operations) == 0 {
		t.Skip("No operations available for testing")
	}

	// Switch to first operation
	targetOpID := operations[0].ID
	err = client.UpdateCurrentOperationForUser(ctx, targetOpID)
	if err != nil {
		t.Fatalf("UpdateCurrentOperationForUser failed: %v", err)
	}

	// Verify client updated its current operation
	currentOpID := client.GetCurrentOperation()
	if currentOpID == nil || *currentOpID != targetOpID {
		t.Errorf("Expected current operation ID %d, got %v", targetOpID, currentOpID)
	}

	t.Logf("Switched to operation ID: %d", targetOpID)

	// Restore original operation if it was set
	if originalOpID != nil {
		err = client.UpdateCurrentOperationForUser(ctx, *originalOpID)
		if err != nil {
			t.Logf("Warning: Failed to restore original operation: %v", err)
		}
	}
}

func TestOperations_GetOperatorsByOperation(t *testing.T) {
	SkipIfNoMythic(t)

	client := AuthenticateTestClient(t)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Get available operations
	operations, err := client.GetOperations(ctx)
	if err != nil {
		t.Fatalf("GetOperations failed: %v", err)
	}

	if len(operations) == 0 {
		t.Skip("No operations available for testing")
	}

	// Get operators for first operation
	operators, err := client.GetOperatorsByOperation(ctx, operations[0].ID)
	if err != nil {
		t.Fatalf("GetOperatorsByOperation failed: %v", err)
	}

	if operators == nil {
		t.Fatal("GetOperatorsByOperation returned nil")
	}

	// Should have at least the admin operator
	if len(operators) == 0 {
		t.Fatal("Expected at least one operator in operation")
	}

	// Verify operator structure
	opOp := operators[0]
	if opOp.OperatorID == 0 {
		t.Error("Operator ID should not be 0")
	}
	if opOp.OperationID != operations[0].ID {
		t.Errorf("Expected operation ID %d, got %d", operations[0].ID, opOp.OperationID)
	}
	if opOp.ViewMode == "" {
		t.Error("View mode should not be empty")
	}

	t.Logf("Found %d operator(s) in operation %s", len(operators), operations[0].Name)
	for _, op := range operators {
		if op.Operator != nil {
			t.Logf("  - Operator %s (ID: %d) with view mode: %s", op.Operator.Username, op.OperatorID, op.ViewMode)
		}
	}
}

func TestOperations_OperationEventLog(t *testing.T) {
	SkipIfNoMythic(t)

	client := AuthenticateTestClient(t)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Get available operations
	operations, err := client.GetOperations(ctx)
	if err != nil {
		t.Fatalf("GetOperations failed: %v", err)
	}

	if len(operations) == 0 {
		t.Skip("No operations available for testing")
	}

	operationID := operations[0].ID

	// Create an event log entry
	createReq := &types.CreateOperationEventLogRequest{
		OperationID: operationID,
		Message:     "Integration test event - " + time.Now().Format("2006-01-02 15:04:05"),
		Level:       "info",
		Source:      "sdk-integration-test",
	}

	log, err := client.CreateOperationEventLog(ctx, createReq)
	if err != nil {
		t.Fatalf("CreateOperationEventLog failed: %v", err)
	}

	if log == nil {
		t.Fatal("CreateOperationEventLog returned nil")
	}

	if log.ID == 0 {
		t.Error("Event log ID should not be 0")
	}

	if log.Message != createReq.Message {
		t.Errorf("Expected message %q, got %q", createReq.Message, log.Message)
	}

	if log.Level != createReq.Level {
		t.Errorf("Expected level %q, got %q", createReq.Level, log.Level)
	}

	t.Logf("Created event log: %s", log.String())

	// Retrieve event logs
	logs, err := client.GetOperationEventLog(ctx, operationID, 10)
	if err != nil {
		t.Fatalf("GetOperationEventLog failed: %v", err)
	}

	if logs == nil {
		t.Fatal("GetOperationEventLog returned nil")
	}

	// Should have at least the one we just created
	if len(logs) == 0 {
		t.Error("Expected at least one event log entry")
	}

	// Verify we can find our created log
	found := false
	for _, l := range logs {
		if l.ID == log.ID {
			found = true
			if l.Message != createReq.Message {
				t.Errorf("Retrieved log message %q doesn't match created %q", l.Message, createReq.Message)
			}
			break
		}
	}

	if !found {
		t.Error("Could not find created event log in retrieved logs")
	}

	t.Logf("Retrieved %d event log(s)", len(logs))
}

func TestOperations_CreateOperation_InvalidRequest(t *testing.T) {
	SkipIfNoMythic(t)

	client := AuthenticateTestClient(t)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Test with empty name
	req := &types.CreateOperationRequest{
		Name: "",
	}

	_, err := client.CreateOperation(ctx, req)
	if err == nil {
		t.Fatal("Expected error for empty operation name, got nil")
	}

	t.Logf("Expected error: %v", err)

	// Test with nil request
	_, err = client.CreateOperation(ctx, nil)
	if err == nil {
		t.Fatal("Expected error for nil request, got nil")
	}

	t.Logf("Expected error for nil: %v", err)
}

func TestOperations_UpdateOperation_InvalidRequest(t *testing.T) {
	SkipIfNoMythic(t)

	client := AuthenticateTestClient(t)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Test with zero operation ID
	req := &types.UpdateOperationRequest{
		OperationID: 0,
	}

	_, err := client.UpdateOperation(ctx, req)
	if err == nil {
		t.Fatal("Expected error for zero operation ID, got nil")
	}

	t.Logf("Expected error: %v", err)

	// Test with no fields to update
	req2 := &types.UpdateOperationRequest{
		OperationID: 1,
	}

	_, err = client.UpdateOperation(ctx, req2)
	if err == nil {
		t.Fatal("Expected error for no fields to update, got nil")
	}

	t.Logf("Expected error: %v", err)
}

func TestOperations_CreateEventLog_InvalidRequest(t *testing.T) {
	SkipIfNoMythic(t)

	client := AuthenticateTestClient(t)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Test with zero operation ID
	req := &types.CreateOperationEventLogRequest{
		OperationID: 0,
		Message:     "Test message",
	}

	_, err := client.CreateOperationEventLog(ctx, req)
	if err == nil {
		t.Fatal("Expected error for zero operation ID, got nil")
	}

	t.Logf("Expected error: %v", err)

	// Test with empty message
	req2 := &types.CreateOperationEventLogRequest{
		OperationID: 1,
		Message:     "",
	}

	_, err = client.CreateOperationEventLog(ctx, req2)
	if err == nil {
		t.Fatal("Expected error for empty message, got nil")
	}

	t.Logf("Expected error: %v", err)
}

func TestOperations_ViewModes(t *testing.T) {
	// Test view mode constants
	viewModes := []types.OperatorViewMode{
		types.ViewModeOperator,
		types.ViewModeSpectator,
		types.ViewModeLead,
	}

	for _, mode := range viewModes {
		if mode == "" {
			t.Errorf("View mode should not be empty")
		}
		t.Logf("View mode: %s", mode)
	}
}
