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
