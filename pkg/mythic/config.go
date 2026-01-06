package mythic

import (
	"fmt"
	"time"
)

// Config holds the configuration for the Mythic client.
type Config struct {
	// ServerURL is the base URL of the Mythic server (e.g., "https://mythic.example.com:7443")
	ServerURL string

	// APIToken is the Mythic API token for authentication (preferred method)
	APIToken string

	// Username for username/password authentication (alternative to APIToken)
	Username string

	// Password for username/password authentication (alternative to APIToken)
	Password string

	// AccessToken is the JWT access token (if already authenticated)
	AccessToken string

	// RefreshToken is the JWT refresh token (if already authenticated)
	RefreshToken string

	// SSL enables HTTPS/WSS connections. Set to false for HTTP/WS.
	SSL bool

	// Timeout is the global timeout for API requests. Zero means no timeout.
	Timeout time.Duration

	// SkipTLSVerify skips TLS certificate verification (use for self-signed certs)
	SkipTLSVerify bool
}

// Validate checks if the configuration is valid.
func (c *Config) Validate() error {
	if c.ServerURL == "" {
		return fmt.Errorf("ServerURL is required")
	}

	// Must have either APIToken, Username/Password, or Access/Refresh tokens
	hasAPIToken := c.APIToken != ""
	hasUserPass := c.Username != "" && c.Password != ""
	hasTokens := c.AccessToken != "" && c.RefreshToken != ""

	if !hasAPIToken && !hasUserPass && !hasTokens {
		return fmt.Errorf("authentication required: provide APIToken, Username/Password, or AccessToken/RefreshToken")
	}

	return nil
}

// DefaultConfig returns a Config with sensible defaults.
func DefaultConfig() *Config {
	return &Config{
		SSL:     true,
		Timeout: 120 * time.Second,
	}
}
