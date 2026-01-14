//go:build integration

package integration

import (
	"context"
	"testing"
	"time"

	"github.com/nbaertsch/mythic-sdk-go/pkg/mythic/types"
)

// TestE2E_TokenRetrieval tests all token retrieval operations.
// Covers: GetTokens, GetTokensByOperation, GetTokenByID
func TestE2E_TokenRetrieval(t *testing.T) {
	client := AuthenticateTestClient(t)

	// Test 1: Get all tokens for current operation
	t.Log("=== Test 1: Get all tokens ===")
	ctx1, cancel1 := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel1()

	tokens, err := client.GetTokens(ctx1)
	if err != nil {
		t.Fatalf("GetTokens failed: %v", err)
	}
	t.Logf("✓ Retrieved %d tokens for current operation", len(tokens))

	if len(tokens) == 0 {
		t.Log("⚠ No tokens found, this is expected if no credential collection has occurred")
		// Continue with other tests that don't depend on existing tokens
	} else {
		// Validate token structure
		for _, token := range tokens {
			if token.ID == 0 {
				t.Error("Token has ID 0")
			}
			if token.OperationID == 0 {
				t.Error("Token has OperationID 0")
			}
		}

		// Show sample tokens
		sampleCount := 5
		if len(tokens) < sampleCount {
			sampleCount = len(tokens)
		}
		t.Logf("  Sample tokens:")
		for i := 0; i < sampleCount; i++ {
			tok := tokens[i]
			t.Logf("    [%d] %s (Host: %s)", i+1, tok.User, tok.Host)
		}

		// Test 2: Get tokens by operation
		testOperationID := tokens[0].OperationID
		t.Log("=== Test 2: Get tokens by operation ===")
		ctx2, cancel2 := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel2()

		opTokens, err := client.GetTokensByOperation(ctx2, testOperationID)
		if err != nil {
			t.Fatalf("GetTokensByOperation failed: %v", err)
		}
		t.Logf("✓ Retrieved %d tokens for operation %d", len(opTokens), testOperationID)

		// Verify all tokens belong to operation
		for _, token := range opTokens {
			if token.OperationID != testOperationID {
				t.Errorf("Token %d has wrong OperationID: expected %d, got %d",
					token.ID, testOperationID, token.OperationID)
			}
		}

		// Test 3: Get token by ID
		t.Log("=== Test 3: Get token by ID ===")
		ctx3, cancel3 := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel3()

		testTokenID := tokens[0].ID
		token, err := client.GetTokenByID(ctx3, testTokenID)
		if err != nil {
			t.Fatalf("GetTokenByID failed: %v", err)
		}
		if token.ID != testTokenID {
			t.Errorf("Token ID mismatch: expected %d, got %d", testTokenID, token.ID)
		}
		t.Logf("✓ Token %d retrieved: %s (Host: %s)", token.ID, token.User, token.Host)
	}

	t.Log("=== ✓ Token retrieval tests passed ===")
}

// TestE2E_CallbackTokens tests callback token operations.
// Covers: GetCallbackTokens, GetCallbackTokensByCallback
func TestE2E_CallbackTokens(t *testing.T) {
	client := AuthenticateTestClient(t)

	// Test 1: Get all callback tokens
	t.Log("=== Test 1: Get all callback tokens ===")
	ctx1, cancel1 := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel1()

	callbackTokens, err := client.GetCallbackTokens(ctx1)
	if err != nil {
		t.Fatalf("GetCallbackTokens failed: %v", err)
	}
	t.Logf("✓ Retrieved %d callback tokens", len(callbackTokens))

	if len(callbackTokens) == 0 {
		t.Log("⚠ No callback tokens found")
		return
	}

	// Validate structure
	for _, ct := range callbackTokens {
		if ct.ID == 0 {
			t.Error("CallbackToken has ID 0")
		}
		if ct.CallbackID == 0 {
			t.Error("CallbackToken has CallbackID 0")
		}
		if ct.TokenID == 0 {
			t.Error("CallbackToken has TokenID 0")
		}
	}

	// Test 2: Get callback tokens for specific callback
	if len(callbackTokens) > 0 {
		t.Log("=== Test 2: Get callback tokens for specific callback ===")
		ctx2, cancel2 := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel2()

		testCallbackID := callbackTokens[0].CallbackID
		cbTokens, err := client.GetCallbackTokensByCallback(ctx2, testCallbackID)
		if err != nil {
			t.Fatalf("GetCallbackTokensByCallback failed: %v", err)
		}
		t.Logf("✓ Retrieved %d tokens for callback %d", len(cbTokens), testCallbackID)

		// Verify all belong to callback
		for _, ct := range cbTokens {
			if ct.CallbackID != testCallbackID {
				t.Error("CallbackToken has wrong CallbackID")
			}
		}
	}

	t.Log("=== ✓ Callback token tests passed ===")
}

// TestE2E_APITokens tests API token operations.
// Covers: GetAPITokens, DeleteAPIToken
func TestE2E_APITokens(t *testing.T) {
	client := AuthenticateTestClient(t)

	// Test: Get all API tokens
	t.Log("=== Test: Get all API tokens ===")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	apiTokens, err := client.GetAPITokens(ctx)
	if err != nil {
		t.Fatalf("GetAPITokens failed: %v", err)
	}
	t.Logf("✓ Retrieved %d API tokens", len(apiTokens))

	if len(apiTokens) == 0 {
		t.Log("⚠ No API tokens found")
		return
	}

	// Validate structure
	for _, token := range apiTokens {
		if token.ID == 0 {
			t.Error("APIToken has ID 0")
		}
		if token.TokenValue == "" {
			t.Error("APIToken has empty TokenValue")
		}
	}

	// Show sample (mask token values for security)
	sampleCount := 3
	if len(apiTokens) < sampleCount {
		sampleCount = len(apiTokens)
	}
	t.Logf("  Sample API tokens:")
	for i := 0; i < sampleCount; i++ {
		token := apiTokens[i]
		maskedValue := token.TokenValue
		if len(maskedValue) > 10 {
			maskedValue = maskedValue[:10] + "..." + maskedValue[len(maskedValue)-4:]
		}
		t.Logf("    [%d] %s (Active: %v)", i+1, maskedValue, token.Active)
	}

	// Note: DeleteAPIToken not tested to avoid deleting production tokens
	t.Log("  ⚠ DeleteAPIToken not tested to avoid data loss")

	t.Log("=== ✓ API token tests passed ===")
}

// TestE2E_TokenAttributes tests token attribute analysis.
func TestE2E_TokenAttributes(t *testing.T) {
	client := AuthenticateTestClient(t)

	t.Log("=== Test: Analyze token attributes ===")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	tokens, err := client.GetTokens(ctx)
	if err != nil {
		t.Fatalf("GetTokens failed: %v", err)
	}

	if len(tokens) == 0 {
		t.Skip("No tokens found for attribute analysis")
	}

	t.Logf("✓ Analyzing %d tokens", len(tokens))

	// Analyze attributes
	users := make(map[string]int)
	hosts := make(map[string]int)
	integrityLevels := make(map[string]int)

	for _, token := range tokens {
		if token.User != "" {
			users[token.User]++
		}
		if token.Host != "" {
			hosts[token.Host]++
		}
		integrityLevels[token.GetIntegrityLevelString()]++
	}

	t.Logf("  Top users:")
	userCount := 0
	for user, count := range users {
		if userCount < 5 {
			t.Logf("    %s: %d", user, count)
			userCount++
		}
	}

	t.Logf("  Top hosts:")
	hostCount := 0
	for host, count := range hosts {
		if hostCount < 5 {
			t.Logf("    %s: %d", host, count)
			hostCount++
		}
	}

	t.Logf("  Integrity level distribution:")
	for level, count := range integrityLevels {
		t.Logf("    %s: %d", level, count)
	}

	// Show privileged tokens (high or system integrity)
	var privilegedTokens []*types.Token
	for _, token := range tokens {
		if token.IntegrityLevel >= 3 { // High or System
			privilegedTokens = append(privilegedTokens, token)
		}
	}

	if len(privilegedTokens) > 0 {
		t.Logf("  Found %d privileged tokens (High/System integrity):", len(privilegedTokens))
		showCount := 3
		if len(privilegedTokens) < showCount {
			showCount = len(privilegedTokens)
		}
		for i := 0; i < showCount; i++ {
			token := privilegedTokens[i]
			t.Logf("    %s (Host: %s, Integrity: %s)", token.User, token.Host, token.GetIntegrityLevelString())
		}
	}

	t.Log("=== ✓ Token attribute analysis complete ===")
}

// TestE2E_TokenErrorHandling tests error scenarios for token operations.
func TestE2E_TokenErrorHandling(t *testing.T) {
	client := AuthenticateTestClient(t)

	// Test 1: Get tokens by invalid operation
	t.Log("=== Test 1: Get tokens by invalid operation ===")
	ctx1, cancel1 := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel1()

	_, err := client.GetTokensByOperation(ctx1, 0)
	if err == nil {
		t.Error("Expected error for invalid operation ID")
	}
	t.Logf("✓ Invalid operation ID rejected: %v", err)

	// Test 2: Get token by invalid ID
	t.Log("=== Test 2: Get token by invalid ID ===")
	ctx2, cancel2 := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel2()

	_, err = client.GetTokenByID(ctx2, 0)
	if err == nil {
		t.Error("Expected error for invalid token ID")
	}
	t.Logf("✓ Invalid token ID rejected: %v", err)

	// Test 3: Get token by non-existent ID
	t.Log("=== Test 3: Get token by non-existent ID ===")
	ctx3, cancel3 := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel3()

	_, err = client.GetTokenByID(ctx3, 999999)
	if err == nil {
		t.Error("Expected error for non-existent token")
	}
	t.Logf("✓ Non-existent token rejected: %v", err)

	// Test 4: Get callback tokens by invalid callback
	t.Log("=== Test 4: Get callback tokens by invalid callback ===")
	ctx4, cancel4 := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel4()

	_, err = client.GetCallbackTokensByCallback(ctx4, 0)
	if err == nil {
		t.Error("Expected error for invalid callback ID")
	}
	t.Logf("✓ Invalid callback ID rejected: %v", err)

	// Test 5: Delete API token with invalid ID
	t.Log("=== Test 5: Delete API token with invalid ID ===")
	ctx5, cancel5 := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel5()

	err = client.DeleteAPIToken(ctx5, 0)
	if err == nil {
		t.Error("Expected error for invalid token ID")
	}
	t.Logf("✓ Invalid API token ID rejected: %v", err)

	t.Log("=== ✓ All error handling tests passed ===")
}

// TestE2E_TokenTimestamps tests token timestamp ordering.
func TestE2E_TokenTimestamps(t *testing.T) {
	client := AuthenticateTestClient(t)

	t.Log("=== Test: Token timestamp analysis ===")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	tokens, err := client.GetTokens(ctx)
	if err != nil {
		t.Fatalf("GetTokens failed: %v", err)
	}

	if len(tokens) == 0 {
		t.Skip("No tokens found for timestamp test")
	}

	t.Logf("✓ Retrieved %d tokens", len(tokens))

	// Analyze timestamp distribution
	now := time.Now()
	last24h := 0
	lastWeek := 0
	older := 0

	for _, token := range tokens {
		age := now.Sub(token.Timestamp)
		if age < 24*time.Hour {
			last24h++
		} else if age < 7*24*time.Hour {
			lastWeek++
		} else {
			older++
		}
	}

	t.Logf("  Timestamp distribution:")
	t.Logf("    Last 24 hours: %d", last24h)
	t.Logf("    Last week: %d", lastWeek)
	t.Logf("    Older: %d", older)

	// Show age of newest and oldest
	var newest, oldest time.Time
	for i, token := range tokens {
		if i == 0 {
			newest = token.Timestamp
			oldest = token.Timestamp
		} else {
			if token.Timestamp.After(newest) {
				newest = token.Timestamp
			}
			if token.Timestamp.Before(oldest) {
				oldest = token.Timestamp
			}
		}
	}

	t.Logf("  Newest token: %s (age: %s)", newest.Format(time.RFC3339), now.Sub(newest))
	t.Logf("  Oldest token: %s (age: %s)", oldest.Format(time.RFC3339), now.Sub(oldest))

	t.Log("=== ✓ Timestamp analysis complete ===")
}
