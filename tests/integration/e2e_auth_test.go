//go:build integration

package integration

import (
	"context"
	"testing"
	"time"

	"github.com/nbaertsch/mythic-sdk-go/pkg/mythic"
)

// TestE2E_AuthenticationLifecycle tests the complete authentication workflow
// Covers: Login, Logout, IsAuthenticated, GetMe, CreateAPIToken, DeleteAPIToken, RefreshAccessToken
func TestE2E_AuthenticationLifecycle(t *testing.T) {
	cfg := GetTestConfig(t)

	// Test 1: Login with username/password
	t.Log("=== Test 1: Login with username/password ===")
	client, err := mythic.NewClient(&mythic.Config{
		ServerURL:     cfg.ServerURL,
		Username:      cfg.Username,
		Password:      cfg.Password,
		SSL:           true,
		SkipTLSVerify: cfg.SkipTLSVerify,
		Timeout:       cfg.DefaultTimeout,
	})
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	ctx1, cancel1 := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel1()

	err = client.Login(ctx1)
	if err != nil {
		t.Fatalf("Login failed: %v", err)
	}
	t.Log("✓ Login successful")

	// Test 2: Verify authenticated
	t.Log("=== Test 2: Verify authenticated ===")
	if !client.IsAuthenticated() {
		t.Fatal("Client reports not authenticated after successful login")
	}
	t.Log("✓ IsAuthenticated() returned true")

	// Test 3: Get current user info
	t.Log("=== Test 3: Get current user info (GetMe) ===")
	ctx2, cancel2 := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel2()

	operator, err := client.GetMe(ctx2)
	if err != nil {
		t.Fatalf("GetMe failed: %v", err)
	}
	if operator == nil {
		t.Fatal("GetMe returned nil operator")
	}
	if operator.Username == "" {
		t.Error("Operator username is empty")
	}
	t.Logf("✓ GetMe returned operator: %s (ID: %d)", operator.Username, operator.ID)

	// Test 4: Get current operation
	t.Log("=== Test 4: Get current operation ===")
	operationID := client.GetCurrentOperation()
	if operationID == nil {
		t.Fatal("GetCurrentOperation returned nil")
	}
	t.Logf("✓ Current operation ID: %d", *operationID)

	// Test 5: Create API token
	t.Log("=== Test 5: Create API token ===")
	ctx3, cancel3 := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel3()

	tokenType := "User"
	apiToken, err := client.CreateAPIToken(ctx3, tokenType)
	if err != nil {
		t.Fatalf("CreateAPIToken failed: %v", err)
	}
	if apiToken == nil {
		t.Fatal("CreateAPIToken returned nil token")
	}
	if apiToken.TokenValue == "" {
		t.Fatal("API token value is empty")
	}
	t.Logf("✓ API token created: %s (ID: %d)", apiToken.TokenValue[:20]+"...", apiToken.ID)

	// Store token for later deletion
	tokenID := apiToken.ID
	tokenValue := apiToken.TokenValue

	// Test 6: Create new client with API token
	t.Log("=== Test 6: Create new client with API token ===")
	tokenClient, err := mythic.NewClient(&mythic.Config{
		ServerURL:     cfg.ServerURL,
		APIToken:      tokenValue,
		SSL:           true,
		SkipTLSVerify: cfg.SkipTLSVerify,
		Timeout:       cfg.DefaultTimeout,
	})
	if err != nil {
		t.Fatalf("Failed to create token client: %v", err)
	}
	t.Log("✓ Token client created")

	// Test 7: Login with API token
	t.Log("=== Test 7: Login with API token ===")
	ctx4, cancel4 := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel4()

	err = tokenClient.Login(ctx4)
	if err != nil {
		t.Fatalf("Token login failed: %v", err)
	}
	t.Log("✓ Token login successful")

	// Test 8: Verify token client authenticated
	t.Log("=== Test 8: Verify token client authenticated ===")
	if !tokenClient.IsAuthenticated() {
		t.Fatal("Token client reports not authenticated after successful login")
	}
	t.Log("✓ Token client IsAuthenticated() returned true")

	// Test 9: Call authenticated endpoint with token client
	t.Log("=== Test 9: Call GetMe with token client ===")
	ctx5, cancel5 := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel5()

	tokenOperator, err := tokenClient.GetMe(ctx5)
	if err != nil {
		t.Fatalf("GetMe failed with token client: %v", err)
	}
	if tokenOperator == nil {
		t.Fatal("GetMe returned nil with token client")
	}
	if tokenOperator.Username != operator.Username {
		t.Errorf("Token client operator mismatch: expected %s, got %s", operator.Username, tokenOperator.Username)
	}
	t.Logf("✓ Token client GetMe successful: %s", tokenOperator.Username)

	// Test 10: Refresh access token
	t.Log("=== Test 10: Refresh access token ===")
	ctx6, cancel6 := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel6()

	err = client.RefreshAccessToken(ctx6)
	if err != nil {
		t.Fatalf("RefreshAccessToken failed: %v", err)
	}
	t.Log("✓ Access token refreshed")

	// Test 11: Verify refresh worked by calling GetMe again
	t.Log("=== Test 11: Verify refresh worked ===")
	ctx7, cancel7 := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel7()

	refreshedOperator, err := client.GetMe(ctx7)
	if err != nil {
		t.Fatalf("GetMe failed after refresh: %v", err)
	}
	if refreshedOperator.Username != operator.Username {
		t.Errorf("Operator changed after refresh: expected %s, got %s", operator.Username, refreshedOperator.Username)
	}
	t.Logf("✓ GetMe works after refresh: %s", refreshedOperator.Username)

	// Test 12: Delete API token
	t.Log("=== Test 12: Delete API token ===")
	ctx8, cancel8 := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel8()

	err = client.DeleteAPIToken(ctx8, tokenID)
	if err != nil {
		t.Fatalf("DeleteAPIToken failed: %v", err)
	}
	t.Log("✓ API token deleted")

	// Test 13: Logout password client
	t.Log("=== Test 13: Logout password client ===")
	client.Logout()
	if client.IsAuthenticated() {
		t.Error("Client still reports authenticated after logout")
	}
	t.Log("✓ Password client logged out")

	// Test 14: Logout token client
	t.Log("=== Test 14: Logout token client ===")
	tokenClient.Logout()
	if tokenClient.IsAuthenticated() {
		t.Error("Token client still reports authenticated after logout")
	}
	t.Log("✓ Token client logged out")

	// Test 15: Verify not authenticated
	t.Log("=== Test 15: Verify not authenticated after logout ===")
	if client.IsAuthenticated() {
		t.Fatal("Client reports authenticated after logout")
	}
	if tokenClient.IsAuthenticated() {
		t.Fatal("Token client reports authenticated after logout")
	}
	t.Log("✓ Both clients report not authenticated")

	t.Log("=== ✓ All authentication tests passed ===")
}

// TestE2E_AuthenticationErrorHandling tests authentication error scenarios
func TestE2E_AuthenticationErrorHandling(t *testing.T) {
	cfg := GetTestConfig(t)

	// Test 1: Invalid credentials
	t.Log("=== Test 1: Invalid credentials ===")
	invalidClient, err := mythic.NewClient(&mythic.Config{
		ServerURL:     cfg.ServerURL,
		Username:      "invalid_user",
		Password:      "invalid_password",
		SSL:           true,
		SkipTLSVerify: cfg.SkipTLSVerify,
		Timeout:       cfg.DefaultTimeout,
	})
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	ctx1, cancel1 := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel1()

	err = invalidClient.Login(ctx1)
	if err == nil {
		t.Fatal("Expected error for invalid credentials, got nil")
	}
	t.Logf("✓ Invalid credentials rejected: %v", err)

	// Test 2: Operations without authentication
	t.Log("=== Test 2: Operations without authentication ===")
	unauthClient, err := mythic.NewClient(&mythic.Config{
		ServerURL:     cfg.ServerURL,
		Username:      cfg.Username,
		Password:      cfg.Password,
		SSL:           true,
		SkipTLSVerify: cfg.SkipTLSVerify,
		Timeout:       cfg.DefaultTimeout,
	})
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	ctx2, cancel2 := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel2()

	// Try to call GetMe without logging in
	_, err = unauthClient.GetMe(ctx2)
	if err == nil {
		t.Error("Expected error when calling GetMe without authentication")
	}
	t.Logf("✓ Unauthenticated request rejected: %v", err)

	// Test 3: Invalid API token
	t.Log("=== Test 3: Invalid API token ===")
	invalidTokenClient, err := mythic.NewClient(&mythic.Config{
		ServerURL:     cfg.ServerURL,
		APIToken:      "invalid_token_12345",
		SSL:           true,
		SkipTLSVerify: cfg.SkipTLSVerify,
		Timeout:       cfg.DefaultTimeout,
	})
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	ctx3, cancel3 := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel3()

	err = invalidTokenClient.Login(ctx3)
	if err == nil {
		t.Fatal("Expected error for invalid API token, got nil")
	}
	t.Logf("✓ Invalid API token rejected: %v", err)

	t.Log("=== ✓ All error handling tests passed ===")
}
