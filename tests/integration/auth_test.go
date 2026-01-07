//go:build integration

package integration

import (
	"context"
	"testing"
	"time"

	"github.com/nbaertsch/mythic-sdk-go/pkg/mythic"
)

func TestAuthentication_LoginWithPassword(t *testing.T) {
	SkipIfNoMythic(t)

	client := NewTestClient(t)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Test login
	err := client.Login(ctx)
	if err != nil {
		t.Fatalf("Login failed: %v", err)
	}

	// Verify authenticated
	if !client.IsAuthenticated() {
		t.Error("Client should be authenticated after successful login")
	}

	// Test logout
	client.Logout()
	if client.IsAuthenticated() {
		t.Error("Client should not be authenticated after logout")
	}
}

func TestAuthentication_GetMe(t *testing.T) {
	SkipIfNoMythic(t)

	client := AuthenticateTestClient(t)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	operator, err := client.GetMe(ctx)
	if err != nil {
		t.Fatalf("GetMe failed: %v", err)
	}

	// Verify operator data
	if operator == nil {
		t.Fatal("GetMe returned nil operator")
	}

	cfg := GetTestConfig(t)
	if operator.Username != cfg.Username {
		t.Errorf("Expected username %q, got %q", cfg.Username, operator.Username)
	}

	if operator.ID == 0 {
		t.Error("Operator ID should not be 0")
	}

	// Admin user should have admin flag
	if operator.Username == "mythic_admin" && !operator.Admin {
		t.Error("mythic_admin user should have admin flag set")
	}

	if !operator.Active {
		t.Error("Operator should be active")
	}

	t.Logf("Authenticated as: %s", operator.String())
}

func TestAuthentication_CreateAPIToken(t *testing.T) {
	SkipIfNoMythic(t)

	client := AuthenticateTestClient(t)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Create API token
	token, err := client.CreateAPIToken(ctx)
	if err != nil {
		t.Fatalf("CreateAPIToken failed: %v", err)
	}

	if token == "" {
		t.Fatal("CreateAPIToken returned empty token")
	}

	t.Logf("Created API token: %s...", token[:20])

	// Test using the API token with a new client
	cfg := GetTestConfig(t)
	tokenClient, err := mythic.NewClient(&mythic.Config{
		ServerURL:     cfg.ServerURL,
		APIToken:      token,
		SSL:           true,
		SkipTLSVerify: cfg.SkipTLSVerify,
		Timeout:       cfg.DefaultTimeout,
	})
	if err != nil {
		t.Fatalf("Failed to create client with API token: %v", err)
	}

	// Login should succeed with API token
	ctx2, cancel2 := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel2()

	if err := tokenClient.Login(ctx2); err != nil {
		t.Fatalf("Login with API token failed: %v", err)
	}

	// Verify we can make authenticated calls
	operator, err := tokenClient.GetMe(ctx2)
	if err != nil {
		t.Fatalf("GetMe with API token failed: %v", err)
	}

	if operator.Username != cfg.Username {
		t.Errorf("Expected username %q with API token, got %q", cfg.Username, operator.Username)
	}
}

func TestAuthentication_EnsureAuthenticated(t *testing.T) {
	SkipIfNoMythic(t)

	client := NewTestClient(t)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Should not be authenticated initially
	if client.IsAuthenticated() {
		t.Error("New client should not be authenticated")
	}

	// EnsureAuthenticated should login
	if err := client.EnsureAuthenticated(ctx); err != nil {
		t.Fatalf("EnsureAuthenticated failed: %v", err)
	}

	if !client.IsAuthenticated() {
		t.Error("Client should be authenticated after EnsureAuthenticated")
	}

	// Second call should be no-op
	if err := client.EnsureAuthenticated(ctx); err != nil {
		t.Fatalf("Second EnsureAuthenticated failed: %v", err)
	}
}

func TestAuthentication_InvalidCredentials(t *testing.T) {
	SkipIfNoMythic(t)

	cfg := GetTestConfig(t)

	client, err := mythic.NewClient(&mythic.Config{
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

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Login should fail
	err = client.Login(ctx)
	if err == nil {
		t.Fatal("Login should have failed with invalid credentials")
	}

	// Should not be authenticated
	if client.IsAuthenticated() {
		t.Error("Client should not be authenticated after failed login")
	}
}

func TestAuthentication_GetMeWithoutAuth(t *testing.T) {
	client := NewTestClient(t)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// GetMe should fail without authentication
	_, err := client.GetMe(ctx)
	if err == nil {
		t.Fatal("GetMe should fail without authentication")
	}

	if err != mythic.ErrNotAuthenticated {
		t.Errorf("Expected ErrNotAuthenticated, got: %v", err)
	}
}

func TestAuthentication_CurrentOperationID(t *testing.T) {
	SkipIfNoMythic(t)

	client := AuthenticateTestClient(t)

	// Should have operation ID after login (from auth response)
	opID := client.GetCurrentOperation()
	if opID == nil {
		t.Error("Expected operation ID after login, got nil")
	} else {
		t.Logf("Operation ID set from login: %d", *opID)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// GetMe should return operator with current operation
	operator, err := client.GetMe(ctx)
	if err != nil {
		t.Fatalf("GetMe failed: %v", err)
	}

	if operator.CurrentOperation != nil {
		// Operation ID should be set after GetMe
		opID := client.GetCurrentOperation()
		if opID == nil {
			t.Error("Operation ID should be set after GetMe with current operation")
		} else if *opID != operator.CurrentOperation.ID {
			t.Errorf("Expected operation ID %d, got %d", operator.CurrentOperation.ID, *opID)
		}

		// Test SetCurrentOperation
		newOpID := 999
		client.SetCurrentOperation(newOpID)
		opID = client.GetCurrentOperation()
		if opID == nil || *opID != newOpID {
			t.Errorf("Expected operation ID %d after SetCurrentOperation, got %v", newOpID, opID)
		}
	}
}

func TestAuthentication_RefreshAccessToken(t *testing.T) {
	SkipIfNoMythic(t)

	client := NewTestClient(t)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Login to get initial tokens
	err := client.Login(ctx)
	if err != nil {
		t.Fatalf("Login failed: %v", err)
	}

	// Verify authenticated
	if !client.IsAuthenticated() {
		t.Fatal("Client should be authenticated after login")
	}

	// Store the original access token (we can't directly access it, but we can verify behavior changes)
	// Attempt GetMe to verify current token works
	operator1, err := client.GetMe(ctx)
	if err != nil {
		t.Fatalf("GetMe failed before refresh: %v", err)
	}
	if operator1 == nil {
		t.Fatal("GetMe returned nil operator before refresh")
	}

	t.Logf("Before refresh - Authenticated as: %s", operator1.String())

	// Refresh the access token
	err = client.RefreshAccessToken(ctx)
	if err != nil {
		t.Fatalf("RefreshAccessToken failed: %v", err)
	}

	t.Log("Access token refreshed successfully")

	// Verify we're still authenticated
	if !client.IsAuthenticated() {
		t.Error("Client should still be authenticated after token refresh")
	}

	// Verify the new token works by making another authenticated call
	operator2, err := client.GetMe(ctx)
	if err != nil {
		t.Fatalf("GetMe failed after refresh: %v", err)
	}
	if operator2 == nil {
		t.Fatal("GetMe returned nil operator after refresh")
	}

	// Should still be the same user
	if operator2.Username != operator1.Username {
		t.Errorf("Username changed after refresh: expected %q, got %q", operator1.Username, operator2.Username)
	}

	t.Logf("After refresh - Still authenticated as: %s", operator2.String())
}

func TestAuthentication_RefreshAccessToken_NoRefreshToken(t *testing.T) {
	cfg := GetTestConfig(t)

	// Create client with only access token (no refresh token)
	client, err := mythic.NewClient(&mythic.Config{
		ServerURL:     cfg.ServerURL,
		AccessToken:   "fake-access-token",
		SSL:           true,
		SkipTLSVerify: cfg.SkipTLSVerify,
		Timeout:       cfg.DefaultTimeout,
	})
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// RefreshAccessToken should fail without refresh token
	err = client.RefreshAccessToken(ctx)
	if err == nil {
		t.Fatal("RefreshAccessToken should fail without refresh token")
	}

	t.Logf("Expected error for missing refresh token: %v", err)
}
