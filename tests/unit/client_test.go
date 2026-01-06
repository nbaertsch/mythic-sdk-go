package unit

import (
	"testing"

	"github.com/your-org/mythic-sdk-go/pkg/mythic"
)

func TestNewClient(t *testing.T) {
	tests := []struct {
		name    string
		config  *mythic.Config
		wantErr bool
	}{
		{
			name: "valid configuration",
			config: &mythic.Config{
				ServerURL: "https://mythic.example.com:7443",
				APIToken:  "test-token",
				SSL:       true,
			},
			wantErr: false,
		},
		{
			name: "nil config uses defaults",
			config: &mythic.Config{
				ServerURL: "mythic.local:7443",
				APIToken:  "token",
			},
			wantErr: false,
		},
		{
			name: "invalid config",
			config: &mythic.Config{
				// Missing ServerURL
				APIToken: "test-token",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := mythic.NewClient(tt.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewClient() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && client == nil {
				t.Error("NewClient() returned nil client without error")
			}
			if client != nil {
				defer client.Close()
			}
		})
	}
}

func TestClientAuthentication(t *testing.T) {
	config := &mythic.Config{
		ServerURL: "https://mythic.example.com:7443",
		APIToken:  "test-token",
		SSL:       true,
	}

	client, err := mythic.NewClient(config)
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}
	defer client.Close()

	// Client with API token should be authenticated immediately
	if !client.IsAuthenticated() {
		t.Error("Client with API token should be authenticated")
	}
}

func TestClientNotAuthenticated(t *testing.T) {
	config := &mythic.Config{
		ServerURL: "https://mythic.example.com:7443",
		Username:  "admin",
		Password:  "password",
		SSL:       true,
	}

	client, err := mythic.NewClient(config)
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}
	defer client.Close()

	// Client with username/password should not be authenticated until Login() is called
	if client.IsAuthenticated() {
		t.Error("Client with username/password should not be authenticated before Login()")
	}
}

func TestClientOperationManagement(t *testing.T) {
	config := &mythic.Config{
		ServerURL: "https://mythic.example.com:7443",
		APIToken:  "test-token",
		SSL:       true,
	}

	client, err := mythic.NewClient(config)
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}
	defer client.Close()

	// Initially no operation is set
	if op := client.GetCurrentOperation(); op != nil {
		t.Errorf("GetCurrentOperation() = %v, want nil", op)
	}

	// Set operation
	client.SetCurrentOperation(42)

	// Should return the operation ID
	if op := client.GetCurrentOperation(); op == nil || *op != 42 {
		t.Errorf("GetCurrentOperation() = %v, want 42", op)
	}
}

func TestClientURLStripping(t *testing.T) {
	tests := []struct {
		name      string
		serverURL string
		wantErr   bool
	}{
		{
			name:      "with https scheme",
			serverURL: "https://mythic.example.com:7443",
			wantErr:   false,
		},
		{
			name:      "with http scheme",
			serverURL: "http://mythic.example.com:7443",
			wantErr:   false,
		},
		{
			name:      "without scheme",
			serverURL: "mythic.example.com:7443",
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := &mythic.Config{
				ServerURL: tt.serverURL,
				APIToken:  "test-token",
				SSL:       true,
			}

			client, err := mythic.NewClient(config)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewClient() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if client != nil {
				defer client.Close()
			}
		})
	}
}
