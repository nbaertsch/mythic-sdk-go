package mythic

import (
	"context"
	"encoding/base64"
	"encoding/json"

	"github.com/nbaertsch/mythic-sdk-go/pkg/mythic/types"
)

// decodeFilename attempts to decode a base64-encoded filename.
// In Mythic v3.4.20, filename_text fields are stored as base64.
// If decoding fails, returns the original string.
func decodeFilename(encoded string) string {
	if encoded == "" {
		return encoded
	}

	// Try to decode as base64
	decoded, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		// Not valid base64, return as-is
		return encoded
	}

	// Return decoded string
	return string(decoded)
}

// CreateRandom generates a random string based on a format string.
// This is useful for generating random identifiers, callback IDs, payload names,
// or any other random strings with specific patterns.
//
// Format specifiers:
//   - %s: Random lowercase letters
//   - %S: Random uppercase letters
//   - %d: Random digits
//   - %x: Random lowercase hexadecimal characters
//   - %X: Random uppercase hexadecimal characters
//
// The length parameter (if provided) determines how many characters of each type
// to generate. If not provided, a default length is used.
//
// Parameters:
//   - ctx: Context for cancellation and timeout
//   - format: Format string with specifiers (e.g., "%s-%d" for letters-digits)
//   - length: Optional length for random segments (0 for default)
//
// Returns:
//   - *types.CreateRandomResponse: Generated random string
//   - error: Error if the operation fails
//
// Example:
//
//	// Generate random callback ID like "alpha-1234"
//	result, err := client.CreateRandom(ctx, "%s-%d", 5)
//	if err != nil {
//	    return err
//	}
//	fmt.Printf("Random ID: %s\n", result.RandomString)
//
//	// Generate random hex string like "a1b2c3d4"
//	result, err = client.CreateRandom(ctx, "%x", 8)
func (c *Client) CreateRandom(ctx context.Context, format string, length int) (*types.CreateRandomResponse, error) {
	if err := c.EnsureAuthenticated(ctx); err != nil {
		return nil, err
	}

	// Validate input
	if format == "" {
		return nil, WrapError("CreateRandom", ErrInvalidInput, "format cannot be empty")
	}

	var mutation struct {
		Response struct {
			Status       string `graphql:"status"`
			RandomString string `graphql:"random_string"`
			Error        string `graphql:"error"`
		} `graphql:"createRandom(format: $format, length: $length)"`
	}

	variables := map[string]interface{}{
		"format": format,
		"length": length,
	}

	if err := c.executeMutation(ctx, &mutation, variables); err != nil {
		return nil, WrapError("CreateRandom", err, "failed to create random string")
	}

	// Map to types.CreateRandomResponse
	response := &types.CreateRandomResponse{
		Status:       mutation.Response.Status,
		RandomString: mutation.Response.RandomString,
		Error:        mutation.Response.Error,
	}

	// Check for error in response
	if !response.IsSuccessful() {
		return response, WrapError("CreateRandom", ErrOperationFailed, response.Error)
	}

	return response, nil
}

// ConfigCheck checks the Mythic configuration for validity and returns
// configuration details. This is useful for validating settings before
// operations or debugging configuration issues.
//
// The configuration check validates:
//   - Database connectivity
//   - RabbitMQ connectivity
//   - Redis connectivity (if used)
//   - Container status
//   - Required environment variables
//   - Permission settings
//
// Parameters:
//   - ctx: Context for cancellation and timeout
//
// Returns:
//   - *types.ConfigCheckResponse: Configuration validation results
//   - error: Error if the operation fails
//
// Example:
//
//	// Check configuration before starting operations
//	result, err := client.ConfigCheck(ctx)
//	if err != nil {
//	    return err
//	}
//	if !result.IsValid() {
//	    fmt.Printf("Configuration errors: %v\n", result.GetErrors())
//	    return fmt.Errorf("configuration is invalid")
//	}
//	fmt.Println("Configuration is valid")
func (c *Client) ConfigCheck(ctx context.Context) (*types.ConfigCheckResponse, error) {
	if err := c.EnsureAuthenticated(ctx); err != nil {
		return nil, err
	}

	var query struct {
		Response struct {
			Status  string                 `graphql:"status"`
			Valid   bool                   `graphql:"valid"`
			Errors  []string               `graphql:"errors"`
			Config  map[string]interface{} `graphql:"config"`
			Message string                 `graphql:"message"`
		} `graphql:"config_check"`
	}

	if err := c.executeQuery(ctx, &query, nil); err != nil {
		return nil, WrapError("ConfigCheck", err, "failed to check configuration")
	}

	// Map to types.ConfigCheckResponse
	response := &types.ConfigCheckResponse{
		Status:  query.Response.Status,
		Valid:   query.Response.Valid,
		Errors:  query.Response.Errors,
		Config:  query.Response.Config,
		Message: query.Response.Message,
	}

	return response, nil
}

// formatChoices converts a jsonb choices slice ([]interface{}) to a JSON string
// for storage in type structs. Returns "[]" for nil/empty slices.
func formatChoices(choices []interface{}) string {
	if len(choices) == 0 {
		return "[]"
	}
	data, err := json.Marshal(choices)
	if err != nil {
		return "[]"
	}
	return string(data)
}
