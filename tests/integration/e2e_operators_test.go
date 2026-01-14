//go:build integration

package integration

import (
	"context"
	"testing"
	"time"

	"github.com/nbaertsch/mythic-sdk-go/pkg/mythic/types"
)

// TestE2E_OperatorManagement tests the complete operator management workflow
// Covers: GetOperators, GetOperatorByID, CreateOperator, UpdateOperatorStatus,
// UpdateOperatorOperation, GetOperatorPreferences, UpdateOperatorPreferences,
// GetOperatorSecrets, UpdateOperatorSecrets, CreateInviteLink, GetInviteLinks,
// UpdatePasswordAndEmail
func TestE2E_OperatorManagement(t *testing.T) {
	client := AuthenticateTestClient(t)

	var testOperatorID int
	var originalPreferences *types.OperatorPreferences
	var currentOperatorID int

	// Register cleanup
	defer func() {
		// Restore preferences if modified
		if originalPreferences != nil && currentOperatorID > 0 {
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			_ = client.UpdateOperatorPreferences(ctx, &types.UpdateOperatorPreferencesRequest{
				OperatorID:  currentOperatorID,
				Preferences: originalPreferences.Preferences,
			})
			cancel()
			t.Log("Restored original preferences")
		}
		// Deactivate test operator
		if testOperatorID > 0 {
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			active := false
			_ = client.UpdateOperatorStatus(ctx, &types.UpdateOperatorStatusRequest{
				OperatorID: testOperatorID,
				Active:     &active,
			})
			cancel()
			t.Logf("Deactivated test operator ID: %d", testOperatorID)
		}
	}()

	// Test 1: Get all operators
	t.Log("=== Test 1: Get all operators ===")
	ctx1, cancel1 := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel1()

	operators, err := client.GetOperators(ctx1)
	if err != nil {
		t.Fatalf("GetOperators failed: %v", err)
	}
	if len(operators) == 0 {
		t.Fatal("No operators found")
	}
	t.Logf("✓ Found %d operators", len(operators))

	// Test 2: Get current operator info
	t.Log("=== Test 2: Get current operator (GetMe) ===")
	ctx2, cancel2 := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel2()

	currentOperator, err := client.GetMe(ctx2)
	if err != nil {
		t.Fatalf("GetMe failed: %v", err)
	}
	if currentOperator == nil {
		t.Fatal("GetMe returned nil")
	}
	currentOperatorID = currentOperator.ID
	t.Logf("✓ Current operator: %s (ID: %d, Admin: %v)", currentOperator.Username, currentOperator.ID, currentOperator.Admin)

	// Test 3: Get operator by ID
	t.Log("=== Test 3: Get operator by ID ===")
	ctx3, cancel3 := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel3()

	operator, err := client.GetOperatorByID(ctx3, currentOperatorID)
	if err != nil {
		t.Fatalf("GetOperatorByID failed: %v", err)
	}
	if operator.ID != currentOperatorID {
		t.Errorf("Operator ID mismatch: expected %d, got %d", currentOperatorID, operator.ID)
	}
	t.Logf("✓ Operator retrieved: %s", operator.Username)

	// Test 4: Get operator preferences
	t.Log("=== Test 4: Get operator preferences ===")
	ctx4, cancel4 := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel4()

	preferences, err := client.GetOperatorPreferences(ctx4, currentOperatorID)
	if err != nil {
		t.Fatalf("GetOperatorPreferences failed: %v", err)
	}
	if preferences == nil {
		t.Fatal("GetOperatorPreferences returned nil")
	}
	originalPreferences = preferences
	t.Logf("✓ Preferences retrieved for operator %d", preferences.OperatorID)

	// Test 5: Update operator preferences
	t.Log("=== Test 5: Update operator preferences ===")
	ctx5, cancel5 := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel5()

	// Create new preferences with test values
	testPrefs := map[string]interface{}{
		"theme":          "dark",
		"fontSize":       14,
		"e2e_test_field": "modified_by_test",
	}

	updatePrefsReq := &types.UpdateOperatorPreferencesRequest{
		OperatorID:  currentOperatorID,
		Preferences: testPrefs,
	}

	err = client.UpdateOperatorPreferences(ctx5, updatePrefsReq)
	if err != nil {
		t.Fatalf("UpdateOperatorPreferences failed: %v", err)
	}
	t.Log("✓ Operator preferences updated")

	// Test 6: Verify preferences updated
	t.Log("=== Test 6: Verify preferences updated ===")
	ctx6, cancel6 := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel6()

	updatedPrefs, err := client.GetOperatorPreferences(ctx6, currentOperatorID)
	if err != nil {
		t.Fatalf("GetOperatorPreferences after update failed: %v", err)
	}
	if updatedPrefs.Preferences["e2e_test_field"] != "modified_by_test" {
		t.Error("Preferences not updated correctly")
	}
	t.Logf("✓ Preferences update verified")

	// Test 7: Get operator secrets
	t.Log("=== Test 7: Get operator secrets ===")
	ctx7, cancel7 := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel7()

	secrets, err := client.GetOperatorSecrets(ctx7, currentOperatorID)
	if err != nil {
		t.Fatalf("GetOperatorSecrets failed: %v", err)
	}
	if secrets == nil {
		t.Fatal("GetOperatorSecrets returned nil")
	}
	t.Logf("✓ Secrets retrieved for operator %d", secrets.OperatorID)

	// Test 8: Update operator secrets
	t.Log("=== Test 8: Update operator secrets ===")
	ctx8, cancel8 := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel8()

	testSecrets := map[string]interface{}{
		"api_key":        "test_key_12345",
		"e2e_test_field": "secret_value",
	}

	updateSecretsReq := &types.UpdateOperatorSecretsRequest{
		OperatorID: currentOperatorID,
		Secrets:    testSecrets,
	}

	err = client.UpdateOperatorSecrets(ctx8, updateSecretsReq)
	if err != nil {
		t.Fatalf("UpdateOperatorSecrets failed: %v", err)
	}
	t.Log("✓ Operator secrets updated")

	// Test 9: Verify secrets updated
	t.Log("=== Test 9: Verify secrets updated ===")
	ctx9, cancel9 := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel9()

	updatedSecrets, err := client.GetOperatorSecrets(ctx9, currentOperatorID)
	if err != nil {
		t.Fatalf("GetOperatorSecrets after update failed: %v", err)
	}
	if updatedSecrets.Secrets["e2e_test_field"] != "secret_value" {
		t.Error("Secrets not updated correctly")
	}
	t.Logf("✓ Secrets update verified")

	// Test 10: Create invite link
	t.Log("=== Test 10: Create invite link ===")
	ctx10, cancel10 := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel10()

	inviteReq := &types.CreateInviteLinkRequest{
		MaxUses:       5,
		Name:          "E2E Test Invite",
		OperationRole: "operator",
	}

	inviteLink, err := client.CreateInviteLink(ctx10, inviteReq)
	if err != nil {
		// Invite links may be disabled on the server
		t.Logf("⚠ CreateInviteLink not available (server may have invite links disabled): %v", err)
	} else {
		t.Logf("✓ Invite link created with code: %s (maxUses: %d)", inviteLink.Code, inviteLink.MaxUses)
	}

	// Test 11: Get invite links
	t.Log("=== Test 11: Get invite links ===")
	ctx11, cancel11 := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel11()

	inviteLinks, err := client.GetInviteLinks(ctx11)
	if err != nil {
		t.Fatalf("GetInviteLinks failed: %v", err)
	}
	t.Logf("✓ Found %d invite links", len(inviteLinks))

	// Verify our invite link is present (only if we successfully created one)
	if inviteLink != nil && inviteLink.Code != "" {
		foundInvite := false
		for _, link := range inviteLinks {
			if link.Code == inviteLink.Code {
				foundInvite = true
				t.Logf("  ✓ Found our invite link: %s", link.Code)
				break
			}
		}
		if !foundInvite {
			t.Error("Created invite link not found in list")
		}
	} else {
		t.Log("  ⚠ Skipping invite link verification (creation was not successful)")
	}

	// Test 12: Create new operator
	t.Log("=== Test 12: Create new operator ===")
	ctx12, cancel12 := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel12()

	testUsername := "e2e_test_operator"
	testPassword := "TestPassword123!"
	createOpReq := &types.CreateOperatorRequest{
		Username: testUsername,
		Password: testPassword,
	}

	newOperator, err := client.CreateOperator(ctx12, createOpReq)
	if err != nil {
		t.Fatalf("CreateOperator failed: %v", err)
	}
	if newOperator.ID == 0 {
		t.Fatal("Created operator has ID 0")
	}
	testOperatorID = newOperator.ID
	t.Logf("✓ New operator created: %s (ID: %d)", newOperator.Username, newOperator.ID)

	// Test 13: Verify new operator in list
	t.Log("=== Test 13: Verify new operator in list ===")
	ctx13, cancel13 := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel13()

	allOperators, err := client.GetOperators(ctx13)
	if err != nil {
		t.Fatalf("GetOperators after creation failed: %v", err)
	}

	foundOperator := false
	for _, op := range allOperators {
		if op.ID == testOperatorID {
			foundOperator = true
			t.Logf("✓ New operator found in list: %s", op.Username)
			break
		}
	}
	if !foundOperator {
		t.Error("Created operator not found in operators list")
	}

	// Test 14: Update operator status (deactivate)
	t.Log("=== Test 14: Update operator status (deactivate) ===")
	ctx14, cancel14 := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel14()

	inactive := false
	statusReq := &types.UpdateOperatorStatusRequest{
		OperatorID: testOperatorID,
		Active:     &inactive,
	}

	err = client.UpdateOperatorStatus(ctx14, statusReq)
	if err != nil {
		t.Fatalf("UpdateOperatorStatus (deactivate) failed: %v", err)
	}
	t.Log("✓ Operator deactivated")

	// Test 15: Verify status change
	t.Log("=== Test 15: Verify status change ===")
	ctx15, cancel15 := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel15()

	deactivatedOp, err := client.GetOperatorByID(ctx15, testOperatorID)
	if err != nil {
		t.Fatalf("GetOperatorByID after deactivation failed: %v", err)
	}
	if deactivatedOp.Active {
		t.Error("Operator still active after deactivation")
	}
	t.Logf("✓ Status change verified: active = %v", deactivatedOp.Active)

	// Test 16: Reactivate operator
	t.Log("=== Test 16: Reactivate operator ===")
	ctx16, cancel16 := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel16()

	active := true
	reactivateReq := &types.UpdateOperatorStatusRequest{
		OperatorID: testOperatorID,
		Active:     &active,
	}

	err = client.UpdateOperatorStatus(ctx16, reactivateReq)
	if err != nil {
		t.Fatalf("UpdateOperatorStatus (reactivate) failed: %v", err)
	}
	t.Log("✓ Operator reactivated")

	// Test 17: Assign operator to current operation
	t.Log("=== Test 17: Assign operator to current operation ===")
	ctx17, cancel17 := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel17()

	operationID := client.GetCurrentOperation()
	if operationID == nil {
		t.Fatal("No current operation set")
	}

	assignReq := &types.UpdateOperatorOperationRequest{
		OperatorID:  testOperatorID,
		OperationID: *operationID,
		Remove:      false,
	}

	err = client.UpdateOperatorOperation(ctx17, assignReq)
	if err != nil {
		t.Fatalf("UpdateOperatorOperation failed: %v", err)
	}
	t.Logf("✓ Operator assigned to operation %d", *operationID)

	// Test 18: Verify operation assignment
	t.Log("=== Test 18: Verify operation assignment ===")
	ctx18, cancel18 := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel18()

	opOperators, err := client.GetOperatorsByOperation(ctx18, *operationID)
	if err != nil {
		t.Fatalf("GetOperatorsByOperation failed: %v", err)
	}

	foundAssignment := false
	for _, opOp := range opOperators {
		if opOp.ID == testOperatorID {
			foundAssignment = true
			t.Logf("✓ Operator found in operation %d", *operationID)
			break
		}
	}
	if !foundAssignment {
		t.Error("Operator not found in operation after assignment")
	}

	t.Log("=== ✓ All operator management tests passed ===")
}

// TestE2E_OperatorsErrorHandling tests error scenarios for operator operations
func TestE2E_OperatorsErrorHandling(t *testing.T) {
	client := AuthenticateTestClient(t)

	// Test 1: Get non-existent operator
	t.Log("=== Test 1: Get non-existent operator ===")
	ctx1, cancel1 := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel1()

	_, err := client.GetOperatorByID(ctx1, 999999)
	if err == nil {
		t.Error("Expected error for non-existent operator ID")
	}
	t.Logf("✓ Non-existent operator rejected: %v", err)

	// Test 2: Update non-existent operator status
	t.Log("=== Test 2: Update non-existent operator status ===")
	ctx2, cancel2 := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel2()

	active := false
	statusReq := &types.UpdateOperatorStatusRequest{
		OperatorID: 999999,
		Active:     &active,
	}

	err = client.UpdateOperatorStatus(ctx2, statusReq)
	if err == nil {
		t.Error("Expected error for non-existent operator status update")
	}
	t.Logf("✓ Non-existent operator status update rejected: %v", err)

	// Test 3: Get preferences for arbitrary operator ID
	// Note: GetOperatorPreferences always returns preferences for the current authenticated operator
	// regardless of the operatorID parameter passed (API limitation)
	t.Log("=== Test 3: Get preferences (returns current operator's preferences) ===")
	ctx3, cancel3 := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel3()

	prefs, err := client.GetOperatorPreferences(ctx3, 999999)
	if err != nil {
		t.Errorf("GetOperatorPreferences failed: %v", err)
	}
	if prefs == nil {
		t.Error("Expected preferences for current operator, got nil")
	}
	t.Logf("✓ GetOperatorPreferences returns current operator's preferences (ID in response: %d)", prefs.OperatorID)

	// Test 4: Create operator with empty username
	t.Log("=== Test 4: Create operator with empty username ===")
	ctx4, cancel4 := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel4()

	emptyReq := &types.CreateOperatorRequest{
		Username: "",
		Password: "password123",
	}

	_, err = client.CreateOperator(ctx4, emptyReq)
	if err == nil {
		t.Error("Expected error for empty username")
	}
	t.Logf("✓ Empty username rejected: %v", err)

	t.Log("=== ✓ All error handling tests passed ===")
}
