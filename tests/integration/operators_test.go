//go:build integration
// +build integration

package integration

import (
	"context"
	"testing"
	"time"

	"github.com/nbaertsch/mythic-sdk-go/pkg/mythic/types"
)

// TestOperators_GetOperators tests retrieving all operators
func TestOperators_GetOperators(t *testing.T) {
	ctx := context.Background()

	operators, err := client.GetOperators(ctx)
	if err != nil {
		t.Fatalf("Failed to get operators: %v", err)
	}

	t.Logf("Retrieved %d operators", len(operators))

	// There should be at least the current user
	if len(operators) == 0 {
		t.Error("Expected at least one operator")
	}

	// Verify structure if operators exist
	if len(operators) > 0 {
		op := operators[0]
		if op.Username == "" {
			t.Error("Operator Username should not be empty")
		}

		// Test String method
		str := op.String()
		if str == "" {
			t.Error("Operator.String() should not return empty string")
		}
		t.Logf("Operator: %s", str)
		t.Logf("  ID: %d", op.ID)
		t.Logf("  Username: %s", op.Username)
		t.Logf("  Admin: %v", op.Admin)
		t.Logf("  Active: %v", op.Active)
		t.Logf("  Deleted: %v", op.Deleted)
		t.Logf("  Account Type: %s", op.AccountType)
		t.Logf("  Failed Login Count: %d", op.FailedLoginCount)

		// Test helper methods
		if op.IsAdmin() {
			t.Logf("  Operator has admin privileges")
		}
		if op.IsActive() {
			t.Logf("  Operator is active")
		}
		if op.IsBotAccount() {
			t.Logf("  Operator is a bot account")
		}
		if op.IsLocked() {
			t.Logf("  Operator account is locked")
		}
	}
}

// TestOperators_GetOperatorByID tests retrieving a specific operator
func TestOperators_GetOperatorByID(t *testing.T) {
	ctx := context.Background()

	// Get all operators first
	operators, err := client.GetOperators(ctx)
	if err != nil || len(operators) == 0 {
		t.Skip("No operators available for testing")
	}

	// Get first operator by ID
	operatorID := operators[0].ID
	operator, err := client.GetOperatorByID(ctx, operatorID)
	if err != nil {
		t.Fatalf("Failed to get operator by ID: %v", err)
	}

	if operator.ID != operatorID {
		t.Errorf("Expected operator ID %d, got %d", operatorID, operator.ID)
	}
	if operator.Username == "" {
		t.Error("Operator Username should not be empty")
	}

	t.Logf("Retrieved operator: %s", operator.String())
}

// TestOperators_GetOperatorByID_InvalidInput tests with invalid input
func TestOperators_GetOperatorByID_InvalidInput(t *testing.T) {
	ctx := context.Background()

	// Try to get operator with ID 0
	_, err := client.GetOperatorByID(ctx, 0)
	if err == nil {
		t.Error("Expected error for operator ID 0")
	}
}

// TestOperators_CreateOperator tests creating a new operator
func TestOperators_CreateOperator(t *testing.T) {
	ctx := context.Background()

	// Generate unique username
	username := "testuser_" + time.Now().Format("20060102150405")
	password := "SecurePassword123!"

	request := &types.CreateOperatorRequest{
		Username: username,
		Password: password,
	}

	operator, err := client.CreateOperator(ctx, request)
	if err != nil {
		// This might fail if the current user doesn't have admin privileges
		t.Logf("Failed to create operator (may require admin): %v", err)
		t.Skip("Skipping operator creation test - may require admin privileges")
	}

	if operator.Username != username {
		t.Errorf("Expected username %q, got %q", username, operator.Username)
	}

	t.Logf("Created operator: %s", operator.String())
	t.Logf("  ID: %d", operator.ID)

	// Clean up: mark operator as deleted
	deleteReq := &types.UpdateOperatorStatusRequest{
		OperatorID: operator.ID,
		Deleted:    boolPtr(true),
	}
	if err := client.UpdateOperatorStatus(ctx, deleteReq); err != nil {
		t.Logf("Failed to clean up test operator: %v", err)
	}
}

// TestOperators_CreateOperator_InvalidInput tests operator creation with invalid input
func TestOperators_CreateOperator_InvalidInput(t *testing.T) {
	ctx := context.Background()

	// Test with short password (less than 12 characters)
	request := &types.CreateOperatorRequest{
		Username: "testuser",
		Password: "short",
	}

	_, err := client.CreateOperator(ctx, request)
	if err == nil {
		t.Error("Expected error for password less than 12 characters")
	}

	// Test with empty username
	request = &types.CreateOperatorRequest{
		Username: "",
		Password: "SecurePassword123!",
	}

	_, err = client.CreateOperator(ctx, request)
	if err == nil {
		t.Error("Expected error for empty username")
	}

	// Test with nil request
	_, err = client.CreateOperator(ctx, nil)
	if err == nil {
		t.Error("Expected error for nil request")
	}
}

// TestOperators_UpdateOperatorStatus tests updating operator status
func TestOperators_UpdateOperatorStatus(t *testing.T) {
	ctx := context.Background()

	// Get operators to find one to update
	operators, err := client.GetOperators(ctx)
	if err != nil || len(operators) == 0 {
		t.Skip("No operators available for testing")
	}

	// Find a non-admin operator to test with
	var testOperatorID int
	for _, op := range operators {
		if !op.Admin && !op.Deleted && op.AccountType != types.AccountTypeBot {
			testOperatorID = op.ID
			break
		}
	}

	if testOperatorID == 0 {
		t.Skip("No suitable operator found for status update test")
	}

	// Test updating active status
	active := false
	request := &types.UpdateOperatorStatusRequest{
		OperatorID: testOperatorID,
		Active:     &active,
	}

	err = client.UpdateOperatorStatus(ctx, request)
	if err != nil {
		t.Logf("Failed to update operator status (may require admin): %v", err)
		t.Skip("Skipping operator status update test - may require admin privileges")
	}

	t.Logf("Updated operator %d status", testOperatorID)

	// Revert the change
	active = true
	request.Active = &active
	if err := client.UpdateOperatorStatus(ctx, request); err != nil {
		t.Logf("Failed to revert operator status: %v", err)
	}
}

// TestOperators_UpdateOperatorStatus_InvalidInput tests with invalid input
func TestOperators_UpdateOperatorStatus_InvalidInput(t *testing.T) {
	ctx := context.Background()

	// Test with operator ID 0
	request := &types.UpdateOperatorStatusRequest{
		OperatorID: 0,
		Active:     boolPtr(true),
	}

	err := client.UpdateOperatorStatus(ctx, request)
	if err == nil {
		t.Error("Expected error for operator ID 0")
	}

	// Test with no fields to update
	request = &types.UpdateOperatorStatusRequest{
		OperatorID: 1,
	}

	err = client.UpdateOperatorStatus(ctx, request)
	if err == nil {
		t.Error("Expected error when no fields are provided for update")
	}
}

// TestOperators_GetOperatorPreferences tests retrieving operator preferences
func TestOperators_GetOperatorPreferences(t *testing.T) {
	ctx := context.Background()

	// Get operators to find one to test with
	operators, err := client.GetOperators(ctx)
	if err != nil || len(operators) == 0 {
		t.Skip("No operators available for testing")
	}

	operatorID := operators[0].ID
	prefs, err := client.GetOperatorPreferences(ctx, operatorID)
	if err != nil {
		t.Logf("Failed to get operator preferences: %v", err)
		t.Skip("Skipping preferences test - may not be supported")
	}

	if prefs.OperatorID != operatorID {
		t.Errorf("Expected OperatorID %d, got %d", operatorID, prefs.OperatorID)
	}

	t.Logf("Retrieved preferences for operator %d", operatorID)
	t.Logf("  Preferences JSON length: %d bytes", len(prefs.PreferencesJSON))
}

// TestOperators_GetOperatorSecrets tests retrieving operator secrets
func TestOperators_GetOperatorSecrets(t *testing.T) {
	ctx := context.Background()

	// Get operators to find one to test with
	operators, err := client.GetOperators(ctx)
	if err != nil || len(operators) == 0 {
		t.Skip("No operators available for testing")
	}

	operatorID := operators[0].ID
	secrets, err := client.GetOperatorSecrets(ctx, operatorID)
	if err != nil {
		t.Logf("Failed to get operator secrets: %v", err)
		t.Skip("Skipping secrets test - may not be supported")
	}

	if secrets.OperatorID != operatorID {
		t.Errorf("Expected OperatorID %d, got %d", operatorID, secrets.OperatorID)
	}

	t.Logf("Retrieved secrets for operator %d", operatorID)
	t.Logf("  Secrets JSON length: %d bytes", len(secrets.SecretsJSON))
}

// TestOperators_GetInviteLinks tests retrieving invite links
func TestOperators_GetInviteLinks(t *testing.T) {
	ctx := context.Background()

	links, err := client.GetInviteLinks(ctx)
	if err != nil {
		t.Logf("Failed to get invite links: %v", err)
		t.Skip("Skipping invite links test - may require admin or feature not available")
	}

	t.Logf("Retrieved %d invite links", len(links))

	// Verify structure if links exist
	if len(links) > 0 {
		link := links[0]
		if link.Code == "" {
			t.Error("InviteLink Code should not be empty")
		}

		// Test String method
		str := link.String()
		if str == "" {
			t.Error("InviteLink.String() should not return empty string")
		}
		t.Logf("Invite link: %s", str)
		t.Logf("  ID: %d", link.ID)
		t.Logf("  Code: %s", link.Code)
		t.Logf("  Max Uses: %d", link.MaxUses)
		t.Logf("  Current Uses: %d", link.CurrentUses)
		t.Logf("  Active: %v", link.Active)
		t.Logf("  Expires: %s", link.ExpiresAt.Format(time.RFC3339))

		// Test helper methods
		if link.IsExpired() {
			t.Logf("  Link is expired")
		}
		if link.IsActive() {
			t.Logf("  Link is active and not expired")
		}
		if link.HasUsesRemaining() {
			t.Logf("  Link has uses remaining")
		}
	}
}

// TestOperators_CreateInviteLink tests creating an invite link
func TestOperators_CreateInviteLink(t *testing.T) {
	ctx := context.Background()

	expiresAt := time.Now().Add(7 * 24 * time.Hour) // 7 days from now
	request := &types.CreateInviteLinkRequest{
		MaxUses:   10,
		ExpiresAt: expiresAt,
	}

	link, err := client.CreateInviteLink(ctx, request)
	if err != nil {
		t.Logf("Failed to create invite link (may require admin): %v", err)
		t.Skip("Skipping invite link creation test - may require admin privileges")
	}

	if link.Code == "" {
		t.Error("Invite link Code should not be empty")
	}
	if link.MaxUses != 10 {
		t.Errorf("Expected MaxUses 10, got %d", link.MaxUses)
	}
	if link.CurrentUses != 0 {
		t.Errorf("Expected CurrentUses 0, got %d", link.CurrentUses)
	}

	t.Logf("Created invite link: %s", link.String())
	t.Logf("  Code: %s", link.Code)
}

// TestOperators_CreateInviteLink_InvalidInput tests with invalid input
func TestOperators_CreateInviteLink_InvalidInput(t *testing.T) {
	ctx := context.Background()

	// Test with max uses 0
	request := &types.CreateInviteLinkRequest{
		MaxUses:   0,
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}

	_, err := client.CreateInviteLink(ctx, request)
	if err == nil {
		t.Error("Expected error for max uses 0")
	}

	// Test with expiration in the past
	request = &types.CreateInviteLinkRequest{
		MaxUses:   10,
		ExpiresAt: time.Now().Add(-24 * time.Hour),
	}

	_, err = client.CreateInviteLink(ctx, request)
	if err == nil {
		t.Error("Expected error for expiration date in the past")
	}
}

// TestOperators_OperatorAccountTypes tests different account types
func TestOperators_OperatorAccountTypes(t *testing.T) {
	ctx := context.Background()

	operators, err := client.GetOperators(ctx)
	if err != nil {
		t.Fatalf("Failed to get operators: %v", err)
	}

	userCount := 0
	botCount := 0
	adminCount := 0

	for _, op := range operators {
		if op.IsBotAccount() {
			botCount++
		} else {
			userCount++
		}
		if op.IsAdmin() {
			adminCount++
		}
	}

	t.Logf("Total operators: %d", len(operators))
	t.Logf("  User accounts: %d", userCount)
	t.Logf("  Bot accounts: %d", botCount)
	t.Logf("  Administrators: %d", adminCount)
}

// TestOperators_OperatorStatusChecks tests operator status checks
func TestOperators_OperatorStatusChecks(t *testing.T) {
	ctx := context.Background()

	operators, err := client.GetOperators(ctx)
	if err != nil {
		t.Fatalf("Failed to get operators: %v", err)
	}

	activeCount := 0
	inactiveCount := 0
	deletedCount := 0
	lockedCount := 0

	for _, op := range operators {
		if op.IsActive() {
			activeCount++
		} else {
			inactiveCount++
		}
		if op.IsDeleted() {
			deletedCount++
		}
		if op.IsLocked() {
			lockedCount++
		}
	}

	t.Logf("Operator status breakdown:")
	t.Logf("  Active: %d", activeCount)
	t.Logf("  Inactive: %d", inactiveCount)
	t.Logf("  Deleted: %d", deletedCount)
	t.Logf("  Locked: %d", lockedCount)
}

// boolPtr is a helper function to get a pointer to a bool value
func boolPtr(b bool) *bool {
	return &b
}

// TestOperators_UpdateOperatorOperation tests updating operator's operation settings
func TestOperators_UpdateOperatorOperation(t *testing.T) {
	t.Skip("Skipping UpdateOperatorOperation to avoid modifying operator permissions")

	client := AuthenticateTestClient(t)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Get operators
	operators, err := client.GetOperators(ctx)
	if err != nil || len(operators) == 0 {
		t.Skip("No operators available for testing")
	}

	// Get current operation
	currentOpID := client.GetCurrentOperation()
	if currentOpID == nil {
		t.Skip("No current operation set")
	}

	// Find a non-admin operator
	var testOperatorID int
	for _, op := range operators {
		if !op.Admin && !op.Deleted {
			testOperatorID = op.ID
			break
		}
	}

	if testOperatorID == 0 {
		t.Skip("No suitable operator found for testing")
	}

	// Update operator's view mode in the operation
	viewMode := types.ViewModeOperator
	req := &types.UpdateOperatorOperationRequest{
		OperatorID:  testOperatorID,
		OperationID: *currentOpID,
		ViewMode:    &viewMode,
	}

	err = client.UpdateOperatorOperation(ctx, req)
	if err != nil {
		t.Fatalf("UpdateOperatorOperation failed: %v", err)
	}

	t.Logf("Successfully updated operator %d operation settings", testOperatorID)
}

// TestOperators_UpdateOperatorOperation_InvalidInput tests with invalid input
func TestOperators_UpdateOperatorOperation_InvalidInput(t *testing.T) {

	client := AuthenticateTestClient(t)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Test with nil request
	err := client.UpdateOperatorOperation(ctx, nil)
	if err == nil {
		t.Fatal("Expected error for nil request, got nil")
	}
	t.Logf("Nil request error: %v", err)

	// Test with zero operator ID
	err = client.UpdateOperatorOperation(ctx, &types.UpdateOperatorOperationRequest{
		OperatorID:  0,
		OperationID: 1,
	})
	if err == nil {
		t.Fatal("Expected error for zero operator ID, got nil")
	}
	t.Logf("Zero operator ID error: %v", err)

	// Test with zero operation ID
	err = client.UpdateOperatorOperation(ctx, &types.UpdateOperatorOperationRequest{
		OperatorID:  1,
		OperationID: 0,
	})
	if err == nil {
		t.Fatal("Expected error for zero operation ID, got nil")
	}
	t.Logf("Zero operation ID error: %v", err)
}

// TestOperators_UpdateOperatorPreferences tests updating operator preferences
func TestOperators_UpdateOperatorPreferences(t *testing.T) {
	t.Skip("Skipping UpdateOperatorPreferences to avoid modifying operator settings")

	client := AuthenticateTestClient(t)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Get operators
	operators, err := client.GetOperators(ctx)
	if err != nil || len(operators) == 0 {
		t.Skip("No operators available for testing")
	}

	operatorID := operators[0].ID

	// Update operator preferences
	req := &types.UpdateOperatorPreferencesRequest{
		OperatorID: operatorID,
		Preferences: map[string]interface{}{
			"theme":    "dark",
			"fontSize": 14,
		},
	}

	err = client.UpdateOperatorPreferences(ctx, req)
	if err != nil {
		t.Fatalf("UpdateOperatorPreferences failed: %v", err)
	}

	t.Logf("Successfully updated preferences for operator %d", operatorID)

	// Verify preferences were updated
	prefs, err := client.GetOperatorPreferences(ctx, operatorID)
	if err != nil {
		t.Logf("Could not verify preferences: %v", err)
		return
	}

	t.Logf("Preferences updated: %d bytes", len(prefs.PreferencesJSON))
}

// TestOperators_UpdateOperatorPreferences_InvalidInput tests with invalid input
func TestOperators_UpdateOperatorPreferences_InvalidInput(t *testing.T) {

	client := AuthenticateTestClient(t)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Test with nil request
	err := client.UpdateOperatorPreferences(ctx, nil)
	if err == nil {
		t.Fatal("Expected error for nil request, got nil")
	}
	t.Logf("Nil request error: %v", err)

	// Test with zero operator ID
	err = client.UpdateOperatorPreferences(ctx, &types.UpdateOperatorPreferencesRequest{
		OperatorID:  0,
		Preferences: map[string]interface{}{},
	})
	if err == nil {
		t.Fatal("Expected error for zero operator ID, got nil")
	}
	t.Logf("Zero operator ID error: %v", err)

	// Test with nil preferences
	err = client.UpdateOperatorPreferences(ctx, &types.UpdateOperatorPreferencesRequest{
		OperatorID:  1,
		Preferences: nil,
	})
	if err == nil {
		t.Fatal("Expected error for nil preferences, got nil")
	}
	t.Logf("Nil preferences error: %v", err)
}

// TestOperators_UpdateOperatorSecrets tests updating operator secrets
func TestOperators_UpdateOperatorSecrets(t *testing.T) {
	t.Skip("Skipping UpdateOperatorSecrets to avoid modifying operator secrets")

	client := AuthenticateTestClient(t)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Get operators
	operators, err := client.GetOperators(ctx)
	if err != nil || len(operators) == 0 {
		t.Skip("No operators available for testing")
	}

	operatorID := operators[0].ID

	// Update operator secrets
	req := &types.UpdateOperatorSecretsRequest{
		OperatorID: operatorID,
		Secrets: map[string]interface{}{
			"api_key": "test-key-12345",
		},
	}

	err = client.UpdateOperatorSecrets(ctx, req)
	if err != nil {
		t.Fatalf("UpdateOperatorSecrets failed: %v", err)
	}

	t.Logf("Successfully updated secrets for operator %d", operatorID)

	// Verify secrets were updated
	secrets, err := client.GetOperatorSecrets(ctx, operatorID)
	if err != nil {
		t.Logf("Could not verify secrets: %v", err)
		return
	}

	t.Logf("Secrets updated: %d bytes", len(secrets.SecretsJSON))
}

// TestOperators_UpdateOperatorSecrets_InvalidInput tests with invalid input
func TestOperators_UpdateOperatorSecrets_InvalidInput(t *testing.T) {

	client := AuthenticateTestClient(t)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Test with nil request
	err := client.UpdateOperatorSecrets(ctx, nil)
	if err == nil {
		t.Fatal("Expected error for nil request, got nil")
	}
	t.Logf("Nil request error: %v", err)

	// Test with zero operator ID
	err = client.UpdateOperatorSecrets(ctx, &types.UpdateOperatorSecretsRequest{
		OperatorID: 0,
		Secrets:    map[string]interface{}{},
	})
	if err == nil {
		t.Fatal("Expected error for zero operator ID, got nil")
	}
	t.Logf("Zero operator ID error: %v", err)

	// Test with nil secrets
	err = client.UpdateOperatorSecrets(ctx, &types.UpdateOperatorSecretsRequest{
		OperatorID: 1,
		Secrets:    nil,
	})
	if err == nil {
		t.Fatal("Expected error for nil secrets, got nil")
	}
	t.Logf("Nil secrets error: %v", err)
}

// TestOperators_UpdatePasswordAndEmail tests updating password and email
func TestOperators_UpdatePasswordAndEmail(t *testing.T) {
	t.Skip("Skipping UpdatePasswordAndEmail to avoid modifying operator credentials")

	client := AuthenticateTestClient(t)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Get operators
	operators, err := client.GetOperators(ctx)
	if err != nil || len(operators) == 0 {
		t.Skip("No operators available for testing")
	}

	// Find current operator
	var currentOperatorID int
	for _, op := range operators {
		// Assuming first non-admin is test operator
		if !op.Admin {
			currentOperatorID = op.ID
			break
		}
	}

	if currentOperatorID == 0 {
		t.Skip("No suitable operator found for testing")
	}

	// Update email only (password change requires old password)
	newEmail := "test-updated@example.com"
	req := &types.UpdatePasswordAndEmailRequest{
		OperatorID:  currentOperatorID,
		OldPassword: "test-password",
		Email:       &newEmail,
	}

	err = client.UpdatePasswordAndEmail(ctx, req)
	if err != nil {
		// Expected to fail without correct password
		t.Logf("UpdatePasswordAndEmail failed (expected): %v", err)
		return
	}

	t.Logf("Successfully updated password/email for operator %d", currentOperatorID)
}

// TestOperators_UpdatePasswordAndEmail_InvalidInput tests with invalid input
func TestOperators_UpdatePasswordAndEmail_InvalidInput(t *testing.T) {

	client := AuthenticateTestClient(t)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Test with nil request
	err := client.UpdatePasswordAndEmail(ctx, nil)
	if err == nil {
		t.Fatal("Expected error for nil request, got nil")
	}
	t.Logf("Nil request error: %v", err)

	// Test with zero operator ID
	err = client.UpdatePasswordAndEmail(ctx, &types.UpdatePasswordAndEmailRequest{
		OperatorID:  0,
		OldPassword: "test",
	})
	if err == nil {
		t.Fatal("Expected error for zero operator ID, got nil")
	}
	t.Logf("Zero operator ID error: %v", err)

	// Test with empty old password
	err = client.UpdatePasswordAndEmail(ctx, &types.UpdatePasswordAndEmailRequest{
		OperatorID:  1,
		OldPassword: "",
	})
	if err == nil {
		t.Fatal("Expected error for empty old password, got nil")
	}
	t.Logf("Empty old password error: %v", err)

	// Test with no new password or email
	err = client.UpdatePasswordAndEmail(ctx, &types.UpdatePasswordAndEmailRequest{
		OperatorID:  1,
		OldPassword: "test",
	})
	if err == nil {
		t.Fatal("Expected error when no new password or email provided, got nil")
	}
	t.Logf("No update fields error: %v", err)
}
