package unit

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/nbaertsch/mythic-sdk-go/pkg/mythic"
)

// newTestServer returns an httptest.Server that mimics the Mythic /auth endpoint.
// The handler validates the posted username/password and returns a fake token
// response. The returned server URL already has the scheme stripped (host:port).
func newAuthServer(t *testing.T, expectedUser, expectedPass string) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/auth":
			var req map[string]string
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				http.Error(w, "bad json", http.StatusBadRequest)
				return
			}
			if req["username"] != expectedUser || req["password"] != expectedPass {
				w.WriteHeader(http.StatusForbidden)
				json.NewEncoder(w).Encode(map[string]string{"error": "bad creds"})
				return
			}
			json.NewEncoder(w).Encode(map[string]interface{}{
				"access_token":  "fake-access-token",
				"refresh_token": "fake-refresh-token",
				"user": map[string]interface{}{
					"id":                   1,
					"user_id":              1,
					"username":             expectedUser,
					"current_operation_id": 1,
					"current_operation":    "Default",
				},
			})

		case "/refresh":
			json.NewEncoder(w).Encode(map[string]interface{}{
				"access_token":  "refreshed-access-token",
				"refresh_token": "refreshed-refresh-token",
				"user": map[string]interface{}{
					"id":                   1,
					"user_id":              1,
					"username":             expectedUser,
					"current_operation_id": 1,
				},
			})

		default:
			http.Error(w, "not found", http.StatusNotFound)
		}
	}))
}

func TestLogin_UsernamePassword(t *testing.T) {
	srv := newAuthServer(t, "operator1", "s3cretpass")
	defer srv.Close()

	client, err := mythic.NewClient(&mythic.Config{
		ServerURL: srv.URL,
		Username:  "operator1",
		Password:  "s3cretpass",
		SSL:       false,
	})
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	defer client.Close()

	if client.IsAuthenticated() {
		t.Fatal("should not be authenticated before Login()")
	}

	if err := client.Login(context.Background()); err != nil {
		t.Fatalf("Login() failed: %v", err)
	}

	if !client.IsAuthenticated() {
		t.Fatal("should be authenticated after Login()")
	}

	// Verify tokens were stored
	cfg := client.GetConfig()
	if cfg.AccessToken != "fake-access-token" {
		t.Errorf("AccessToken = %q, want %q", cfg.AccessToken, "fake-access-token")
	}
	if cfg.RefreshToken != "fake-refresh-token" {
		t.Errorf("RefreshToken = %q, want %q", cfg.RefreshToken, "fake-refresh-token")
	}
}

func TestLogin_WrongCredentials(t *testing.T) {
	srv := newAuthServer(t, "admin", "correct-password")
	defer srv.Close()

	client, err := mythic.NewClient(&mythic.Config{
		ServerURL: srv.URL,
		Username:  "admin",
		Password:  "wrong-password",
		SSL:       false,
	})
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	defer client.Close()

	if err := client.Login(context.Background()); err == nil {
		t.Fatal("Login() should have failed with wrong password")
	}

	if client.IsAuthenticated() {
		t.Fatal("should not be authenticated after failed Login()")
	}
}

func TestLogin_MissingCredentials(t *testing.T) {
	client, err := mythic.NewClient(&mythic.Config{
		ServerURL: "http://localhost:9999",
		Username:  "",
		Password:  "",
		SSL:       false,
	})
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	defer client.Close()

	if err := client.Login(context.Background()); err == nil {
		t.Fatal("Login() should fail when username/password are empty")
	}
}

func TestLogout(t *testing.T) {
	srv := newAuthServer(t, "operator1", "pass123")
	defer srv.Close()

	client, err := mythic.NewClient(&mythic.Config{
		ServerURL: srv.URL,
		Username:  "operator1",
		Password:  "pass123",
		SSL:       false,
	})
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	defer client.Close()

	// Login first
	if err := client.Login(context.Background()); err != nil {
		t.Fatalf("Login: %v", err)
	}
	if !client.IsAuthenticated() {
		t.Fatal("should be authenticated after Login()")
	}

	// Logout
	client.Logout()

	if client.IsAuthenticated() {
		t.Fatal("should not be authenticated after Logout()")
	}

	// Tokens should be cleared
	cfg := client.GetConfig()
	if cfg.AccessToken != "" {
		t.Errorf("AccessToken should be empty after Logout, got %q", cfg.AccessToken)
	}
	if cfg.RefreshToken != "" {
		t.Errorf("RefreshToken should be empty after Logout, got %q", cfg.RefreshToken)
	}
}

func TestSetCredentials(t *testing.T) {
	srv := newAuthServer(t, "new_user", "new_pass")
	defer srv.Close()

	// Start with one set of credentials
	client, err := mythic.NewClient(&mythic.Config{
		ServerURL: srv.URL,
		Username:  "old_user",
		Password:  "old_pass",
		SSL:       false,
	})
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	defer client.Close()

	// SetCredentials should update config and clear auth state
	client.SetCredentials("new_user", "new_pass")

	cfg := client.GetConfig()
	if cfg.Username != "new_user" {
		t.Errorf("Username = %q, want %q", cfg.Username, "new_user")
	}
	if cfg.Password != "new_pass" {
		t.Errorf("Password = %q, want %q", cfg.Password, "new_pass")
	}
	if cfg.AccessToken != "" {
		t.Errorf("AccessToken should be cleared after SetCredentials, got %q", cfg.AccessToken)
	}
	if cfg.APIToken != "" {
		t.Errorf("APIToken should be cleared after SetCredentials, got %q", cfg.APIToken)
	}
	if client.IsAuthenticated() {
		t.Error("should not be authenticated after SetCredentials()")
	}

	// Now Login should use the new credentials and succeed
	if err := client.Login(context.Background()); err != nil {
		t.Fatalf("Login() with new credentials failed: %v", err)
	}

	if !client.IsAuthenticated() {
		t.Fatal("should be authenticated after Login() with new credentials")
	}
}

func TestSetCredentials_ClearsExistingAuth(t *testing.T) {
	srv := newAuthServer(t, "user1", "pass1")
	defer srv.Close()

	client, err := mythic.NewClient(&mythic.Config{
		ServerURL: srv.URL,
		Username:  "user1",
		Password:  "pass1",
		SSL:       false,
	})
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	defer client.Close()

	// Login as user1
	if err := client.Login(context.Background()); err != nil {
		t.Fatalf("Login: %v", err)
	}
	if !client.IsAuthenticated() {
		t.Fatal("expected authenticated after Login")
	}

	// SetCredentials should reset everything
	client.SetCredentials("user2", "pass2")

	if client.IsAuthenticated() {
		t.Fatal("should NOT be authenticated after SetCredentials — stale session must be cleared")
	}
	cfg := client.GetConfig()
	if cfg.Username != "user2" || cfg.Password != "pass2" {
		t.Errorf("credentials not updated: got %q / %q", cfg.Username, cfg.Password)
	}
}

func TestSetCredentials_ReLoginAsDifferentUser(t *testing.T) {
	// Server accepts user1/pass1 first, then user2/pass2 after SetCredentials
	mux := http.NewServeMux()
	mux.HandleFunc("/auth", func(w http.ResponseWriter, r *http.Request) {
		var req map[string]string
		json.NewDecoder(r.Body).Decode(&req)

		// Accept both users — just echo back who logged in
		if (req["username"] == "user1" && req["password"] == "pass1") ||
			(req["username"] == "user2" && req["password"] == "pass2") {
			json.NewEncoder(w).Encode(map[string]interface{}{
				"access_token":  "token-for-" + req["username"],
				"refresh_token": "refresh-for-" + req["username"],
				"user": map[string]interface{}{
					"id":       1,
					"username": req["username"],
				},
			})
			return
		}
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(map[string]string{"error": "bad creds"})
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()

	client, err := mythic.NewClient(&mythic.Config{
		ServerURL: srv.URL,
		Username:  "user1",
		Password:  "pass1",
		SSL:       false,
	})
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	defer client.Close()

	// Login as user1
	if err := client.Login(context.Background()); err != nil {
		t.Fatalf("Login as user1: %v", err)
	}
	cfg := client.GetConfig()
	if cfg.AccessToken != "token-for-user1" {
		t.Fatalf("expected token-for-user1, got %q", cfg.AccessToken)
	}

	// Switch to user2
	client.SetCredentials("user2", "pass2")
	if err := client.Login(context.Background()); err != nil {
		t.Fatalf("Login as user2: %v", err)
	}
	cfg = client.GetConfig()
	if cfg.AccessToken != "token-for-user2" {
		t.Fatalf("expected token-for-user2, got %q", cfg.AccessToken)
	}
}

func TestRefreshAccessToken(t *testing.T) {
	srv := newAuthServer(t, "operator1", "pass")
	defer srv.Close()

	client, err := mythic.NewClient(&mythic.Config{
		ServerURL: srv.URL,
		Username:  "operator1",
		Password:  "pass",
		SSL:       false,
	})
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	defer client.Close()

	// Login first to get tokens
	if err := client.Login(context.Background()); err != nil {
		t.Fatalf("Login: %v", err)
	}

	// Refresh
	if err := client.RefreshAccessToken(context.Background()); err != nil {
		t.Fatalf("RefreshAccessToken: %v", err)
	}

	cfg := client.GetConfig()
	if cfg.AccessToken != "refreshed-access-token" {
		t.Errorf("AccessToken = %q, want %q", cfg.AccessToken, "refreshed-access-token")
	}
	if cfg.RefreshToken != "refreshed-refresh-token" {
		t.Errorf("RefreshToken = %q, want %q", cfg.RefreshToken, "refreshed-refresh-token")
	}
}

func TestRefreshAccessToken_NoRefreshToken(t *testing.T) {
	client, err := mythic.NewClient(&mythic.Config{
		ServerURL: "http://localhost:9999",
		Username:  "u",
		Password:  "p",
		SSL:       false,
	})
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	defer client.Close()

	// No login, no refresh token
	if err := client.RefreshAccessToken(context.Background()); err == nil {
		t.Fatal("RefreshAccessToken should fail without a refresh token")
	}
}

func TestLogin_ContextCancellation(t *testing.T) {
	srv := newAuthServer(t, "u", "p")
	defer srv.Close()

	client, err := mythic.NewClient(&mythic.Config{
		ServerURL: srv.URL,
		Username:  "u",
		Password:  "p",
		SSL:       false,
	})
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	defer client.Close()

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // cancel immediately

	if err := client.Login(ctx); err == nil {
		t.Fatal("Login() should fail with cancelled context")
	}
}
