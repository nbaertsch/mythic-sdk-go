//go:build integration

package integration

import (
	"context"
	"testing"
	"time"

	"github.com/nbaertsch/mythic-sdk-go/pkg/mythic/types"
)

func TestCredentials_GetCredentials(t *testing.T) {

	client := AuthenticateTestClient(t)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	credentials, err := client.GetCredentials(ctx)
	if err != nil {
		t.Fatalf("GetCredentials failed: %v", err)
	}

	if credentials == nil {
		t.Fatal("GetCredentials returned nil")
	}

	t.Logf("Found %d credential(s)", len(credentials))

	// If there are credentials, verify structure
	if len(credentials) > 0 {
		cred := credentials[0]
		if cred.ID == 0 {
			t.Error("Credential ID should not be 0")
		}
		if cred.Type == "" {
			t.Error("Credential type should not be empty")
		}
		t.Logf("First credential: %s", cred.String())
	}
}

func TestCredentials_CreateAndManageCredential(t *testing.T) {

	client := AuthenticateTestClient(t)

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	// Create new credential
	createReq := &types.CreateCredentialRequest{
		Type:       "plaintext",
		Account:    "testuser_" + time.Now().Format("150405"),
		Realm:      "TESTDOMAIN",
		Credential: "TestPassword123!",
		Comment:    "Integration test credential - " + time.Now().Format("2006-01-02 15:04:05"),
		Metadata:   `{"test":true,"source":"sdk-integration-test"}`,
	}

	cred, err := client.CreateCredential(ctx, createReq)
	if err != nil {
		t.Fatalf("CreateCredential failed: %v", err)
	}

	if cred == nil {
		t.Fatal("CreateCredential returned nil")
	}

	if cred.ID == 0 {
		t.Error("Created credential should have an ID")
	}

	if cred.Type != createReq.Type {
		t.Errorf("Expected type %q, got %q", createReq.Type, cred.Type)
	}

	if cred.Account != createReq.Account {
		t.Errorf("Expected account %q, got %q", createReq.Account, cred.Account)
	}

	if cred.Realm != createReq.Realm {
		t.Errorf("Expected realm %q, got %q", createReq.Realm, cred.Realm)
	}

	t.Logf("Created credential: %s (ID: %d)", cred.String(), cred.ID)

	// Update the credential
	newComment := "Updated test credential - " + time.Now().Format("15:04:05")
	newType := "hash"

	updateReq := &types.UpdateCredentialRequest{
		ID:      cred.ID,
		Type:    &newType,
		Comment: &newComment,
	}

	updated, err := client.UpdateCredential(ctx, updateReq)
	if err != nil {
		t.Fatalf("UpdateCredential failed: %v", err)
	}

	if updated.Type != newType {
		t.Errorf("Expected type %q, got %q", newType, updated.Type)
	}

	if updated.Comment != newComment {
		t.Errorf("Expected comment %q, got %q", newComment, updated.Comment)
	}

	t.Logf("Updated credential: %s", updated.String())

	// Get credentials to verify it exists
	credentials, err := client.GetCredentials(ctx)
	if err != nil {
		t.Fatalf("GetCredentials failed: %v", err)
	}

	found := false
	for _, c := range credentials {
		if c.ID == cred.ID {
			found = true
			if c.Comment != newComment {
				t.Errorf("Retrieved credential comment %q doesn't match updated %q", c.Comment, newComment)
			}
			break
		}
	}

	if !found {
		t.Error("Could not find created credential in list")
	}

	// Cleanup: Delete the credential
	err = client.DeleteCredential(ctx, cred.ID)
	if err != nil {
		t.Logf("Warning: Failed to delete test credential: %v", err)
	} else {
		t.Log("Cleaned up test credential")
	}

	// Verify it's deleted
	credentials, err = client.GetCredentials(ctx)
	if err != nil {
		t.Fatalf("GetCredentials failed after delete: %v", err)
	}

	for _, c := range credentials {
		if c.ID == cred.ID {
			t.Error("Deleted credential still appears in list")
		}
	}
}

func TestCredentials_CreateCredential_InvalidRequest(t *testing.T) {

	client := AuthenticateTestClient(t)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Test with nil request
	_, err := client.CreateCredential(ctx, nil)
	if err == nil {
		t.Fatal("Expected error for nil request, got nil")
	}
	t.Logf("Expected error for nil: %v", err)

	// Test with empty type
	req := &types.CreateCredentialRequest{
		Account:    "user",
		Credential: "password",
	}
	_, err = client.CreateCredential(ctx, req)
	if err == nil {
		t.Fatal("Expected error for empty type, got nil")
	}
	t.Logf("Expected error for empty type: %v", err)

	// Test with empty account
	req2 := &types.CreateCredentialRequest{
		Type:       "plaintext",
		Credential: "password",
	}
	_, err = client.CreateCredential(ctx, req2)
	if err == nil {
		t.Fatal("Expected error for empty account, got nil")
	}
	t.Logf("Expected error for empty account: %v", err)

	// Test with empty credential
	req3 := &types.CreateCredentialRequest{
		Type:    "plaintext",
		Account: "user",
	}
	_, err = client.CreateCredential(ctx, req3)
	if err == nil {
		t.Fatal("Expected error for empty credential, got nil")
	}
	t.Logf("Expected error for empty credential: %v", err)
}

func TestCredentials_UpdateCredential_InvalidRequest(t *testing.T) {

	client := AuthenticateTestClient(t)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Test with nil request
	_, err := client.UpdateCredential(ctx, nil)
	if err == nil {
		t.Fatal("Expected error for nil request, got nil")
	}
	t.Logf("Expected error for nil: %v", err)

	// Test with zero ID
	req := &types.UpdateCredentialRequest{
		ID: 0,
	}
	_, err = client.UpdateCredential(ctx, req)
	if err == nil {
		t.Fatal("Expected error for zero ID, got nil")
	}
	t.Logf("Expected error for zero ID: %v", err)

	// Test with no fields to update
	req2 := &types.UpdateCredentialRequest{
		ID: 1,
	}
	_, err = client.UpdateCredential(ctx, req2)
	if err == nil {
		t.Fatal("Expected error for no fields to update, got nil")
	}
	t.Logf("Expected error for no fields: %v", err)
}

func TestCredentials_UpdateCredential_NotFound(t *testing.T) {

	client := AuthenticateTestClient(t)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Try to update non-existent credential
	comment := "test"
	req := &types.UpdateCredentialRequest{
		ID:      999999,
		Comment: &comment,
	}

	_, err := client.UpdateCredential(ctx, req)
	if err == nil {
		t.Fatal("Expected error for non-existent credential, got nil")
	}

	t.Logf("Expected error: %v", err)
}

func TestCredentials_DeleteCredential_InvalidID(t *testing.T) {

	client := AuthenticateTestClient(t)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Test with zero ID
	err := client.DeleteCredential(ctx, 0)
	if err == nil {
		t.Fatal("Expected error for zero ID, got nil")
	}

	t.Logf("Expected error: %v", err)
}

func TestCredentials_GetCredentialsByOperation(t *testing.T) {

	client := AuthenticateTestClient(t)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Get current operation
	currentOpID := client.GetCurrentOperation()
	if currentOpID == nil {
		t.Skip("No current operation set")
	}

	credentials, err := client.GetCredentialsByOperation(ctx, *currentOpID)
	if err != nil {
		t.Fatalf("GetCredentialsByOperation failed: %v", err)
	}

	if credentials == nil {
		t.Fatal("GetCredentialsByOperation returned nil")
	}

	t.Logf("Found %d credential(s) for operation %d", len(credentials), *currentOpID)

	// Verify all credentials belong to the operation
	for _, cred := range credentials {
		if cred.OperationID != *currentOpID {
			t.Errorf("Expected operation ID %d, got %d", *currentOpID, cred.OperationID)
		}
	}
}

func TestCredentials_CreateDifferentTypes(t *testing.T) {

	client := AuthenticateTestClient(t)

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	testCases := []struct {
		name       string
		credType   string
		account    string
		credential string
		realm      string
	}{
		{
			name:       "plaintext password",
			credType:   "plaintext",
			account:    "admin",
			credential: "P@ssw0rd!",
			realm:      "WORKSTATION",
		},
		{
			name:       "NTLM hash",
			credType:   "hash",
			account:    "sqlserver",
			credential: "aad3b435b51404eeaad3b435b51404ee:8846f7eaee8fb117ad06bdd830b7586c",
			realm:      "DOMAIN.LOCAL",
		},
		{
			name:       "SSH key",
			credType:   "key",
			account:    "root",
			credential: "-----BEGIN RSA PRIVATE KEY-----\nMIIEpAIBAAKCAQEA...\n-----END RSA PRIVATE KEY-----",
			realm:      "server.example.com",
		},
	}

	createdIDs := []int{}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := &types.CreateCredentialRequest{
				Type:       tc.credType,
				Account:    tc.account + "_" + time.Now().Format("150405"),
				Credential: tc.credential,
				Realm:      tc.realm,
				Comment:    "Test credential: " + tc.name,
			}

			cred, err := client.CreateCredential(ctx, req)
			if err != nil {
				t.Fatalf("CreateCredential failed: %v", err)
			}

			if cred.Type != tc.credType {
				t.Errorf("Expected type %q, got %q", tc.credType, cred.Type)
			}

			t.Logf("Created %s: %s", tc.name, cred.String())
			createdIDs = append(createdIDs, cred.ID)
		})
	}

	// Cleanup all created credentials
	for _, id := range createdIDs {
		err := client.DeleteCredential(ctx, id)
		if err != nil {
			t.Logf("Warning: Failed to delete credential %d: %v", id, err)
		}
	}
}

func TestCredentials_PartialUpdate(t *testing.T) {

	client := AuthenticateTestClient(t)

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	// Create a credential
	createReq := &types.CreateCredentialRequest{
		Type:       "plaintext",
		Account:    "partial_test_" + time.Now().Format("150405"),
		Realm:      "TESTDOMAIN",
		Credential: "original_password",
		Comment:    "Original comment",
	}

	cred, err := client.CreateCredential(ctx, createReq)
	if err != nil {
		t.Fatalf("CreateCredential failed: %v", err)
	}

	// Update only the comment
	newComment := "Only comment updated"
	updateReq := &types.UpdateCredentialRequest{
		ID:      cred.ID,
		Comment: &newComment,
	}

	updated, err := client.UpdateCredential(ctx, updateReq)
	if err != nil {
		t.Fatalf("UpdateCredential failed: %v", err)
	}

	// Verify only comment changed
	if updated.Comment != newComment {
		t.Errorf("Expected comment %q, got %q", newComment, updated.Comment)
	}

	if updated.Type != cred.Type {
		t.Error("Type should not have changed")
	}

	if updated.Account != cred.Account {
		t.Error("Account should not have changed")
	}

	if updated.Credential != cred.Credential {
		t.Error("Credential should not have changed")
	}

	t.Logf("Partial update successful: %s", updated.String())

	// Cleanup
	err = client.DeleteCredential(ctx, cred.ID)
	if err != nil {
		t.Logf("Warning: Failed to delete test credential: %v", err)
	}
}
