//go:build integration
// +build integration

package integration

import (
	"context"
	"testing"
	"time"

	"github.com/nbaertsch/mythic-sdk-go/pkg/mythic/types"
)

// TestTokens_GetTokens tests retrieving all tokens
func TestTokens_GetTokens(t *testing.T) {
	ctx := context.Background()

	// Get tokens for current operation
	tokens, err := client.GetTokens(ctx)
	if err != nil {
		t.Fatalf("Failed to get tokens: %v", err)
	}

	t.Logf("Retrieved %d tokens", len(tokens))

	// Verify structure if any tokens exist
	if len(tokens) > 0 {
		token := tokens[0]
		if token.ID == 0 {
			t.Error("Token ID should not be 0")
		}
		if token.OperationID == 0 {
			t.Error("Token OperationID should not be 0")
		}
		if token.Timestamp.IsZero() {
			t.Error("Token timestamp should not be zero")
		}

		// Test String method
		str := token.String()
		if str == "" {
			t.Error("Token.String() should not return empty string")
		}
		t.Logf("Token: %s", str)

		// Test helper methods
		t.Logf("Token deleted: %v", token.IsDeleted())
		t.Logf("Token has task: %v", token.HasTask())
		t.Logf("Token integrity level: %s", token.GetIntegrityLevelString())
	}
}

// TestTokens_GetTokensByOperation tests retrieving tokens for a specific operation
func TestTokens_GetTokensByOperation(t *testing.T) {
	ctx := context.Background()

	// Get current operation
	operationID := client.GetCurrentOperation()
	if operationID == nil {
		t.Skip("No current operation set")
	}

	// Get tokens for operation
	tokens, err := client.GetTokensByOperation(ctx, *operationID)
	if err != nil {
		t.Fatalf("Failed to get tokens by operation: %v", err)
	}

	t.Logf("Retrieved %d tokens for operation %d", len(tokens), *operationID)

	// Verify all tokens belong to the operation
	for _, token := range tokens {
		if token.OperationID != *operationID {
			t.Errorf("Expected operation ID %d, got %d", *operationID, token.OperationID)
		}
		if token.IsDeleted() {
			t.Error("Should not return deleted tokens")
		}
	}
}

// TestTokens_GetTokenByID_InvalidID tests getting a token with invalid ID
func TestTokens_GetTokenByID_InvalidID(t *testing.T) {
	ctx := context.Background()

	// Try to get token with ID 0
	_, err := client.GetTokenByID(ctx, 0)
	if err == nil {
		t.Error("Expected error for token ID 0")
	}

	// Try to get token with very high ID (likely doesn't exist)
	_, err = client.GetTokenByID(ctx, 999999999)
	if err == nil {
		t.Error("Expected error for non-existent token")
	}
}

// TestTokens_GetCallbackTokens tests retrieving callback tokens
func TestTokens_GetCallbackTokens(t *testing.T) {
	ctx := context.Background()

	// Get callback tokens for current operation
	callbackTokens, err := client.GetCallbackTokens(ctx)
	if err != nil {
		t.Fatalf("Failed to get callback tokens: %v", err)
	}

	t.Logf("Retrieved %d callback tokens", len(callbackTokens))

	// Verify structure if any callback tokens exist
	if len(callbackTokens) > 0 {
		ct := callbackTokens[0]
		if ct.ID == 0 {
			t.Error("CallbackToken ID should not be 0")
		}
		if ct.CallbackID == 0 {
			t.Error("CallbackToken CallbackID should not be 0")
		}
		if ct.TokenID == 0 {
			t.Error("CallbackToken TokenID should not be 0")
		}
		if ct.Timestamp.IsZero() {
			t.Error("CallbackToken timestamp should not be zero")
		}

		// Test String method
		str := ct.String()
		if str == "" {
			t.Error("CallbackToken.String() should not return empty string")
		}
		t.Logf("CallbackToken: %s", str)
	}
}

// TestTokens_GetCallbackTokensByCallback tests retrieving tokens for a specific callback
func TestTokens_GetCallbackTokensByCallback(t *testing.T) {
	ctx := context.Background()

	// First get all callbacks
	callbacks, err := client.GetAllCallbacks(ctx)
	if err != nil {
		t.Fatalf("Failed to get callbacks: %v", err)
	}

	if len(callbacks) == 0 {
		t.Skip("No callbacks available for testing")
	}

	// Get tokens for the first callback
	callbackID := callbacks[0].ID
	callbackTokens, err := client.GetCallbackTokensByCallback(ctx, callbackID)
	if err != nil {
		t.Fatalf("Failed to get callback tokens: %v", err)
	}

	t.Logf("Retrieved %d tokens for callback %d", len(callbackTokens), callbackID)

	// Verify all tokens belong to the callback
	for _, ct := range callbackTokens {
		if ct.CallbackID != callbackID {
			t.Errorf("Expected callback ID %d, got %d", callbackID, ct.CallbackID)
		}
	}
}

// TestTokens_GetCallbackTokensByCallback_InvalidID tests invalid callback ID
func TestTokens_GetCallbackTokensByCallback_InvalidID(t *testing.T) {
	ctx := context.Background()

	// Try to get callback tokens with ID 0
	_, err := client.GetCallbackTokensByCallback(ctx, 0)
	if err == nil {
		t.Error("Expected error for callback ID 0")
	}
}

// TestTokens_GetAPITokens tests retrieving API tokens
func TestTokens_GetAPITokens(t *testing.T) {
	ctx := context.Background()

	// Get API tokens
	apiTokens, err := client.GetAPITokens(ctx)
	if err != nil {
		t.Fatalf("Failed to get API tokens: %v", err)
	}

	t.Logf("Retrieved %d API tokens", len(apiTokens))

	// Verify structure if any API tokens exist
	if len(apiTokens) > 0 {
		token := apiTokens[0]
		if token.ID == 0 {
			t.Error("APIToken ID should not be 0")
		}
		if token.TokenType == "" {
			t.Error("APIToken TokenType should not be empty")
		}
		if token.CreationTime.IsZero() {
			t.Error("APIToken CreationTime should not be zero")
		}

		// Test String method
		str := token.String()
		if str == "" {
			t.Error("APIToken.String() should not return empty string")
		}
		t.Logf("APIToken: %s", str)

		// Test helper methods
		t.Logf("APIToken active: %v", token.IsActive())
		t.Logf("APIToken deleted: %v", token.IsDeleted())
	}
}

// TestTokens_CreateAndDeleteAPIToken tests creating and deleting an API token
func TestTokens_CreateAndDeleteAPIToken(t *testing.T) {
	ctx := context.Background()

	// Create a new API token
	tokenValue, err := client.CreateAPIToken(ctx)
	if err != nil {
		t.Fatalf("Failed to create API token: %v", err)
	}

	t.Logf("Created API token: %s", tokenValue[:20]+"...")

	// Wait a moment for token to be fully created
	time.Sleep(2 * time.Second)

	// Get all API tokens to find the one we just created
	apiTokens, err := client.GetAPITokens(ctx)
	if err != nil {
		t.Fatalf("Failed to get API tokens: %v", err)
	}

	// Find our token
	var createdTokenID int
	for _, token := range apiTokens {
		if token.TokenValue == tokenValue {
			createdTokenID = token.ID
			break
		}
	}

	if createdTokenID == 0 {
		t.Fatal("Could not find created API token")
	}

	t.Logf("Found created API token with ID: %d", createdTokenID)

	// Delete the token
	err = client.DeleteAPIToken(ctx, createdTokenID)
	if err != nil {
		t.Fatalf("Failed to delete API token: %v", err)
	}

	t.Logf("Successfully deleted API token %d", createdTokenID)

	// Verify token is deleted
	time.Sleep(2 * time.Second)
	apiTokens, err = client.GetAPITokens(ctx)
	if err != nil {
		t.Fatalf("Failed to get API tokens after deletion: %v", err)
	}

	// Check that the token is not in the list (deleted tokens are filtered out)
	for _, token := range apiTokens {
		if token.ID == createdTokenID {
			t.Error("Deleted token should not appear in API tokens list")
		}
	}
}

// TestTokens_DeleteAPIToken_InvalidID tests deleting with invalid ID
func TestTokens_DeleteAPIToken_InvalidID(t *testing.T) {
	ctx := context.Background()

	// Try to delete token with ID 0
	err := client.DeleteAPIToken(ctx, 0)
	if err == nil {
		t.Error("Expected error for token ID 0")
	}

	// Try to delete token with non-existent ID
	err = client.DeleteAPIToken(ctx, 999999999)
	if err == nil {
		t.Error("Expected error for non-existent token")
	}
}

// TestTokens_TimestampOrdering tests that tokens are ordered by timestamp
func TestTokens_TimestampOrdering(t *testing.T) {
	ctx := context.Background()

	// Get tokens
	tokens, err := client.GetTokens(ctx)
	if err != nil {
		t.Fatalf("Failed to get tokens: %v", err)
	}

	if len(tokens) < 2 {
		t.Skip("Need at least 2 tokens to test ordering")
	}

	// Verify tokens are in descending order by timestamp (newest first)
	for i := 0; i < len(tokens)-1; i++ {
		if tokens[i].Timestamp.Before(tokens[i+1].Timestamp) {
			t.Errorf("Tokens not in descending timestamp order: %v before %v",
				tokens[i].Timestamp, tokens[i+1].Timestamp)
		}
	}

	t.Logf("Verified %d tokens are properly ordered by timestamp", len(tokens))
}

// TestTokens_IntegrityLevels tests token integrity level handling
func TestTokens_IntegrityLevels(t *testing.T) {
	// Test all integrity levels
	levels := map[int]string{
		0:  "Untrusted",
		1:  "Low",
		2:  "Medium",
		3:  "High",
		4:  "System",
		99: "Unknown",
	}

	for level, expected := range levels {
		token := types.Token{IntegrityLevel: level}
		result := token.GetIntegrityLevelString()
		if result != expected {
			t.Errorf("Expected integrity level %q for %d, got %q", expected, level, result)
		}
	}
}

// TestTokens_NoCurrentOperation tests error when no current operation is set
func TestTokens_NoCurrentOperation(t *testing.T) {
	ctx := context.Background()

	// Save current operation
	originalOp := client.GetCurrentOperation()

	// Clear current operation
	client.SetCurrentOperation(nil)

	// Try to get tokens without operation set
	_, err := client.GetTokens(ctx)
	if err == nil {
		t.Error("Expected error when no current operation is set")
	}

	// Try to get callback tokens without operation set
	_, err = client.GetCallbackTokens(ctx)
	if err == nil {
		t.Error("Expected error when no current operation is set")
	}

	// Restore original operation
	client.SetCurrentOperation(originalOp)
}

// TestTokens_CallbackTokenOrdering tests callback token timestamp ordering
func TestTokens_CallbackTokenOrdering(t *testing.T) {
	ctx := context.Background()

	// Get callback tokens
	callbackTokens, err := client.GetCallbackTokens(ctx)
	if err != nil {
		t.Fatalf("Failed to get callback tokens: %v", err)
	}

	if len(callbackTokens) < 2 {
		t.Skip("Need at least 2 callback tokens to test ordering")
	}

	// Verify callback tokens are in descending order by timestamp (newest first)
	for i := 0; i < len(callbackTokens)-1; i++ {
		if callbackTokens[i].Timestamp.Before(callbackTokens[i+1].Timestamp) {
			t.Errorf("Callback tokens not in descending timestamp order: %v before %v",
				callbackTokens[i].Timestamp, callbackTokens[i+1].Timestamp)
		}
	}

	t.Logf("Verified %d callback tokens are properly ordered by timestamp", len(callbackTokens))
}

// TestTokens_APITokenOrdering tests API token creation time ordering
func TestTokens_APITokenOrdering(t *testing.T) {
	ctx := context.Background()

	// Get API tokens
	apiTokens, err := client.GetAPITokens(ctx)
	if err != nil {
		t.Fatalf("Failed to get API tokens: %v", err)
	}

	if len(apiTokens) < 2 {
		t.Skip("Need at least 2 API tokens to test ordering")
	}

	// Verify API tokens are in descending order by creation time (newest first)
	for i := 0; i < len(apiTokens)-1; i++ {
		if apiTokens[i].CreationTime.Before(apiTokens[i+1].CreationTime) {
			t.Errorf("API tokens not in descending creation time order: %v before %v",
				apiTokens[i].CreationTime, apiTokens[i+1].CreationTime)
		}
	}

	t.Logf("Verified %d API tokens are properly ordered by creation time", len(apiTokens))
}
