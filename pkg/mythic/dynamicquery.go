package mythic

import (
	"context"

	"github.com/nbaertsch/mythic-sdk-go/pkg/mythic/types"
)

// DynamicQueryFunction executes a dynamic query function for a command parameter.
// Dynamic queries allow command parameters to populate their choices dynamically at runtime,
// pulling data from the Mythic database, callback information, or other sources.
//
// This is useful for parameters that need current data like:
//   - List of active callbacks
//   - Available files from a specific path
//   - Currently loaded payload types
//   - Operation-specific configurations
//
// Parameters:
//   - ctx: Context for cancellation and timeout
//   - command: The command name to query
//   - parameters: Current parameter values (used by the query function)
//   - callbackID: Optional callback ID for context-specific queries (0 if not needed)
//
// Returns:
//   - *types.DynamicQueryResponse: List of choices for the parameter
//   - error: Error if the operation fails
//
// Example:
//
//	// Get dynamic choices for a file browser parameter
//	params := map[string]interface{}{
//	    "path": "/home/user",
//	}
//	result, err := client.DynamicQueryFunction(ctx, "ls", params, callbackID)
//	if err != nil {
//	    return err
//	}
//	if result.IsSuccessful() {
//	    fmt.Printf("Available choices: %v\n", result.Choices)
//	}
func (c *Client) DynamicQueryFunction(ctx context.Context, command string, parameters map[string]interface{}, callbackID int) (*types.DynamicQueryResponse, error) {
	if err := c.EnsureAuthenticated(ctx); err != nil {
		return nil, err
	}

	// Validate input
	if command == "" {
		return nil, WrapError("DynamicQueryFunction", ErrInvalidInput, "command name cannot be empty")
	}

	var query struct {
		Response struct {
			Status  string        `graphql:"status"`
			Choices []interface{} `graphql:"choices"`
			Error   string        `graphql:"error"`
		} `graphql:"dynamicQueryFunction(command: $command, parameters: $parameters, callback_id: $callback_id)"`
	}

	variables := map[string]interface{}{
		"command":     command,
		"parameters":  parameters,
		"callback_id": callbackID,
	}

	if err := c.executeQuery(ctx, &query, variables); err != nil {
		return nil, WrapError("DynamicQueryFunction", err, "failed to execute dynamic query")
	}

	// Map to types.DynamicQueryResponse
	response := &types.DynamicQueryResponse{
		Status:  query.Response.Status,
		Choices: query.Response.Choices,
		Error:   query.Response.Error,
	}

	// Check for error in response
	if !response.IsSuccessful() {
		return response, WrapError("DynamicQueryFunction", ErrOperationFailed, response.Error)
	}

	return response, nil
}

// DynamicBuildParameter queries dynamic choices for a payload build parameter.
// Build parameters can use dynamic queries to populate their choices based on
// the current payload configuration, available resources, or operation settings.
//
// This is particularly useful for:
//   - Selecting C2 profiles based on operation
//   - Choosing encryption keys from available key stores
//   - Selecting callback intervals from approved ranges
//   - Populating configuration options from templates
//
// Parameters:
//   - ctx: Context for cancellation and timeout
//   - payloadType: The payload type name (e.g., "apollo", "poseidon")
//   - parameter: The build parameter name to query
//   - parameters: Current build parameter values (optional, can be nil)
//
// Returns:
//   - *types.DynamicBuildParameterResponse: List of choices for the parameter
//   - error: Error if the operation fails
//
// Example:
//
//	// Get available C2 profiles for a payload type
//	result, err := client.DynamicBuildParameter(ctx, "apollo", "c2_profile", nil)
//	if err != nil {
//	    return err
//	}
//	if result.IsSuccessful() && result.HasChoices() {
//	    fmt.Printf("Available C2 profiles: %v\n", result.Choices)
//	}
func (c *Client) DynamicBuildParameter(ctx context.Context, payloadType, parameter string, parameters map[string]interface{}) (*types.DynamicBuildParameterResponse, error) {
	if err := c.EnsureAuthenticated(ctx); err != nil {
		return nil, err
	}

	// Validate input
	if payloadType == "" {
		return nil, WrapError("DynamicBuildParameter", ErrInvalidInput, "payload type cannot be empty")
	}
	if parameter == "" {
		return nil, WrapError("DynamicBuildParameter", ErrInvalidInput, "parameter name cannot be empty")
	}

	var query struct {
		Response struct {
			Status  string        `graphql:"status"`
			Choices []interface{} `graphql:"choices"`
			Error   string        `graphql:"error"`
		} `graphql:"dynamicBuildParameter(payload_type: $payload_type, parameter: $parameter, parameters: $parameters)"`
	}

	variables := map[string]interface{}{
		"payload_type": payloadType,
		"parameter":    parameter,
		"parameters":   parameters,
	}

	if err := c.executeQuery(ctx, &query, variables); err != nil {
		return nil, WrapError("DynamicBuildParameter", err, "failed to query build parameter")
	}

	// Map to types.DynamicBuildParameterResponse
	response := &types.DynamicBuildParameterResponse{
		Status:  query.Response.Status,
		Choices: query.Response.Choices,
		Error:   query.Response.Error,
	}

	// Check for error in response
	if !response.IsSuccessful() {
		return response, WrapError("DynamicBuildParameter", ErrOperationFailed, response.Error)
	}

	return response, nil
}

// TypedArrayParseFunction parses a typed array string into structured elements.
// Typed arrays are special parameter formats that allow users to input structured
// data in a compact string format, which is then parsed into proper data structures.
//
// Common typed array formats include:
//   - JSON arrays: [{"key": "value"}, {"key": "value2"}]
//   - Key-value pairs: key1=value1,key2=value2
//   - File paths with attributes: /path/to/file:read,/other/path:write
//
// The parsing behavior depends on the parameter type definition, which specifies
// the expected format and structure.
//
// Parameters:
//   - ctx: Context for cancellation and timeout
//   - inputArray: The string representation of the typed array
//   - parameterType: The parameter type that defines the parsing rules
//
// Returns:
//   - *types.TypedArrayParseResponse: Parsed array elements
//   - error: Error if the operation fails
//
// Example:
//
//	// Parse a typed array of file paths with permissions
//	input := "/etc/passwd:read,/etc/shadow:read,/var/log/auth.log:write"
//	result, err := client.TypedArrayParseFunction(ctx, input, "file_list")
//	if err != nil {
//	    return err
//	}
//	if result.IsSuccessful() && result.HasElements() {
//	    for _, elem := range result.ParsedArray {
//	        fmt.Printf("Parsed element: %v\n", elem)
//	    }
//	}
func (c *Client) TypedArrayParseFunction(ctx context.Context, inputArray, parameterType string) (*types.TypedArrayParseResponse, error) {
	if err := c.EnsureAuthenticated(ctx); err != nil {
		return nil, err
	}

	// Validate input
	if inputArray == "" {
		return nil, WrapError("TypedArrayParseFunction", ErrInvalidInput, "input array cannot be empty")
	}
	if parameterType == "" {
		return nil, WrapError("TypedArrayParseFunction", ErrInvalidInput, "parameter type cannot be empty")
	}

	var query struct {
		Response struct {
			Status      string        `graphql:"status"`
			ParsedArray []interface{} `graphql:"parsed_array"`
			Error       string        `graphql:"error"`
		} `graphql:"typedArrayParseFunction(input_array: $input_array, parameter_type: $parameter_type)"`
	}

	variables := map[string]interface{}{
		"input_array":    inputArray,
		"parameter_type": parameterType,
	}

	if err := c.executeQuery(ctx, &query, variables); err != nil {
		return nil, WrapError("TypedArrayParseFunction", err, "failed to parse typed array")
	}

	// Map to types.TypedArrayParseResponse
	response := &types.TypedArrayParseResponse{
		Status:      query.Response.Status,
		ParsedArray: query.Response.ParsedArray,
		Error:       query.Response.Error,
	}

	// Check for error in response
	if !response.IsSuccessful() {
		return response, WrapError("TypedArrayParseFunction", ErrOperationFailed, response.Error)
	}

	return response, nil
}
