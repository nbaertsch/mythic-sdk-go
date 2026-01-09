package types

import "fmt"

// CreateRandomRequest represents a request to generate a random string with a format.
type CreateRandomRequest struct {
	Format string `json:"format"`
	Length int    `json:"length,omitempty"`
}

// String returns a human-readable representation of the request.
func (c *CreateRandomRequest) String() string {
	if c.Length > 0 {
		return fmt.Sprintf("Generate random string (format: %s, length: %d)", c.Format, c.Length)
	}
	return fmt.Sprintf("Generate random string (format: %s)", c.Format)
}

// CreateRandomResponse represents the response from creating a random string.
type CreateRandomResponse struct {
	Status       string `json:"status"`
	RandomString string `json:"random_string"`
	Error        string `json:"error,omitempty"`
}

// String returns a human-readable representation of the response.
func (c *CreateRandomResponse) String() string {
	if c.Status == "success" {
		return fmt.Sprintf("Generated: %s", c.RandomString)
	}
	return fmt.Sprintf("Failed to generate: %s", c.Error)
}

// IsSuccessful returns true if the random string generation succeeded.
func (c *CreateRandomResponse) IsSuccessful() bool {
	return c.Status == "success" && c.RandomString != ""
}

// ConfigCheckResponse represents the response from checking configuration.
type ConfigCheckResponse struct {
	Status  string                 `json:"status"`
	Valid   bool                   `json:"valid"`
	Errors  []string               `json:"errors,omitempty"`
	Config  map[string]interface{} `json:"config,omitempty"`
	Message string                 `json:"message,omitempty"`
}

// String returns a human-readable representation of the response.
func (c *ConfigCheckResponse) String() string {
	if c.Valid {
		return "Configuration is valid"
	}
	if len(c.Errors) > 0 {
		return fmt.Sprintf("Configuration has %d error(s): %v", len(c.Errors), c.Errors)
	}
	return fmt.Sprintf("Configuration check: %s", c.Message)
}

// IsValid returns true if the configuration is valid.
func (c *ConfigCheckResponse) IsValid() bool {
	return c.Valid && len(c.Errors) == 0
}

// HasErrors returns true if there are configuration errors.
func (c *ConfigCheckResponse) HasErrors() bool {
	return len(c.Errors) > 0
}

// GetErrors returns the list of configuration errors.
func (c *ConfigCheckResponse) GetErrors() []string {
	return c.Errors
}
