//go:build integration

package integration

import (
	"context"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/nbaertsch/mythic-sdk-go/pkg/mythic"
)

// TestConfig holds configuration for integration tests.
type TestConfig struct {
	ServerURL      string
	Username       string
	Password       string
	SkipTLSVerify  bool
	DefaultTimeout time.Duration
}

// GetTestConfig loads test configuration from environment variables.
func GetTestConfig(t *testing.T) *TestConfig {
	t.Helper()

	url := os.Getenv("MYTHIC_URL")
	if url == "" {
		url = "https://localhost:7443"
	}

	username := os.Getenv("MYTHIC_USERNAME")
	if username == "" {
		username = "mythic_admin"
	}

	password := os.Getenv("MYTHIC_PASSWORD")
	if password == "" {
		password = "mythic_password"
	}

	skipTLS := true
	if skip := os.Getenv("MYTHIC_SKIP_TLS_VERIFY"); skip != "" {
		if parsed, err := strconv.ParseBool(skip); err == nil {
			skipTLS = parsed
		}
	}

	return &TestConfig{
		ServerURL:      url,
		Username:       username,
		Password:       password,
		SkipTLSVerify:  skipTLS,
		DefaultTimeout: 180 * time.Second, // Increased for slow webhooks in CI environments
	}
}

// NewTestClient creates a new Mythic client for testing.
func NewTestClient(t *testing.T) *mythic.Client {
	t.Helper()

	cfg := GetTestConfig(t)

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

	return client
}

// AuthenticateTestClient creates and authenticates a client for testing.
func AuthenticateTestClient(t *testing.T) *mythic.Client {
	t.Helper()

	client := NewTestClient(t)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := client.Login(ctx); err != nil {
		t.Fatalf("Failed to authenticate: %v", err)
	}

	if !client.IsAuthenticated() {
		t.Fatal("Client reports not authenticated after successful login")
	}

	return client
}

// SkipIfNoMythic skips the test if Mythic server is not available.
func SkipIfNoMythic(t *testing.T) {
	t.Helper()

	client := NewTestClient(t)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Login(ctx); err != nil {
		t.Skipf("Mythic server not available: %v", err)
	}
}

// RequireMythicVersion skips the test if Mythic version doesn't meet requirements.
// This is a placeholder for future version checking functionality.
func RequireMythicVersion(t *testing.T, minVersion string) {
	t.Helper()
	// TODO: Implement version checking when API supports it
}

// getTestClient is an alias for AuthenticateTestClient for consistency with newer tests.
func getTestClient(t *testing.T) *mythic.Client {
	return AuthenticateTestClient(t)
}

// contains is a helper function to check if a string contains a substring.
func contains(s, substr string) bool {
	return strings.Contains(s, substr)
}
