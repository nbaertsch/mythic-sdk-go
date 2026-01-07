package unit

import (
	"testing"
	"time"

	"github.com/nbaertsch/mythic-sdk-go/pkg/mythic"
)

func TestDefaultConfig(t *testing.T) {
	config := mythic.DefaultConfig()

	if !config.SSL {
		t.Error("DefaultConfig should have SSL enabled")
	}

	if config.Timeout != 120*time.Second {
		t.Errorf("DefaultConfig timeout = %v, want %v", config.Timeout, 120*time.Second)
	}
}

func TestConfigValidate(t *testing.T) {
	tests := []struct {
		name    string
		config  *mythic.Config
		wantErr bool
	}{
		{
			name: "valid with API token",
			config: &mythic.Config{
				ServerURL: "https://mythic.example.com:7443",
				APIToken:  "test-token",
				SSL:       true,
			},
			wantErr: false,
		},
		{
			name: "valid with username/password",
			config: &mythic.Config{
				ServerURL: "https://mythic.example.com:7443",
				Username:  "admin",
				Password:  "password",
				SSL:       true,
			},
			wantErr: false,
		},
		{
			name: "valid with access/refresh tokens",
			config: &mythic.Config{
				ServerURL:    "https://mythic.example.com:7443",
				AccessToken:  "access",
				RefreshToken: "refresh",
				SSL:          true,
			},
			wantErr: false,
		},
		{
			name: "missing ServerURL",
			config: &mythic.Config{
				APIToken: "test-token",
			},
			wantErr: true,
		},
		{
			name: "missing authentication",
			config: &mythic.Config{
				ServerURL: "https://mythic.example.com:7443",
			},
			wantErr: true,
		},
		{
			name: "incomplete username/password",
			config: &mythic.Config{
				ServerURL: "https://mythic.example.com:7443",
				Username:  "admin",
				// Missing Password
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Config.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
