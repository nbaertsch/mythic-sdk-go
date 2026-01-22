//go:build integration

package integration

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/nbaertsch/mythic-sdk-go/pkg/mythic/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestE2E_GetCommands_SchemaValidation validates that GetCommands returns properly
// populated command objects with all required fields and correct types.
// This test goes beyond just checking err != nil - it validates the response data.
func TestE2E_Commands_GetAll_SchemaValidation(t *testing.T) {
	client := AuthenticateTestClient(t)

	t.Log("=== Test: GetCommands comprehensive field validation ===")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	commands, err := client.GetCommands(ctx)
	require.NoError(t, err, "GetCommands should succeed")
	require.NotEmpty(t, commands, "Should return at least one command")

	t.Logf("✓ Retrieved %d commands", len(commands))

	// Validate EVERY command has required fields populated
	for i, cmd := range commands {
		// Critical fields that MUST be present
		assert.NotZero(t, cmd.ID, "Command[%d] has ID 0", i)
		assert.NotEmpty(t, cmd.Cmd, "Command[%d] has empty name", i)
		assert.NotZero(t, cmd.PayloadTypeID, "Command[%d] has PayloadTypeID 0", i)

		// Description and Help should exist (may be empty for some commands)
		assert.NotNil(t, cmd.Description, "Command[%d] has nil Description", i)
		assert.NotNil(t, cmd.Help, "Command[%d] has nil Help", i)

		// Version should be >= 0 (0 is valid for initial version)
		assert.GreaterOrEqual(t, cmd.Version, 0, "Command[%d] has negative Version", i)

		// Note: Author is optional - some commands don't have an author

		// Log first few commands for debugging
		if i < 5 {
			t.Logf("  Command[%d]: %s (ID:%d, PayloadType:%d, Author:%s, Version:%d)",
				i, cmd.Cmd, cmd.ID, cmd.PayloadTypeID, cmd.Author, cmd.Version)
		}
	}

	// Ensure we have a reasonable number of commands (Poseidon should have 40+)
	assert.Greater(t, len(commands), 10, "Expected at least 10 commands across all payload types")

	t.Log("✓ All commands have valid field data")
	t.Log("=== ✓ GetCommands schema validation passed ===")
}

// TestE2E_GetCommandParameters_Complete validates that GetCommandParameters returns
// complete parameter metadata including types, descriptions, and required flags.
func TestE2E_Commands_GetParameters_Complete(t *testing.T) {
	client := AuthenticateTestClient(t)

	t.Log("=== Test: GetCommandParameters comprehensive validation ===")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	parameters, err := client.GetCommandParameters(ctx)
	require.NoError(t, err, "GetCommandParameters should succeed")

	// May be zero if no commands have parameters defined
	t.Logf("✓ Retrieved %d command parameters", len(parameters))

	if len(parameters) == 0 {
		t.Log("⚠ No command parameters found - this may be expected for simple agents")
		return
	}

	// Validate parameter metadata
	paramsByCommand := make(map[int]int) // command_id -> count
	requiredCount := 0
	typesSeen := make(map[string]int)

	for i, param := range parameters {
		// Critical fields
		assert.NotZero(t, param.ID, "Parameter[%d] has ID 0", i)
		assert.NotZero(t, param.CommandID, "Parameter[%d] has CommandID 0", i)
		assert.NotEmpty(t, param.Name, "Parameter[%d] has empty Name", i)
		assert.NotEmpty(t, param.Type, "Parameter[%d] has empty Type", i)

		// Track statistics
		paramsByCommand[param.CommandID]++
		typesSeen[param.Type]++
		if param.Required {
			requiredCount++
		}

		// Log first few parameters
		if i < 5 {
			t.Logf("  Parameter[%d]: %s (Type:%s, Required:%v, CommandID:%d)",
				i, param.Name, param.Type, param.Required, param.CommandID)
		}
	}

	// Statistics
	t.Logf("✓ Statistics:")
	t.Logf("  - Commands with parameters: %d", len(paramsByCommand))
	t.Logf("  - Required parameters: %d", requiredCount)
	t.Logf("  - Parameter types seen: %v", typesSeen)

	// Validate parameter types are valid
	validTypes := map[string]bool{
		"String": true, "Boolean": true, "Number": true, "Choice": true,
		"Array": true, "ChooseOne": true, "ChooseMultiple": true,
		"File": true, "LinkInfo": true, "AgentConnect": true,
		"PayloadList": true, "ConnectionInfo": true, "TypedArray": true,
		"CredentialJson": true, // Added for credential-based parameters
	}

	for paramType := range typesSeen {
		assert.True(t, validTypes[paramType], "Unknown parameter type: %s", paramType)
	}

	t.Log("✓ All parameter metadata is valid")
	t.Log("=== ✓ GetCommandParameters validation passed ===")
}

// TestE2E_GetCommandWithParameters_AllPayloadTypes validates that GetCommandWithParameters
// works correctly across all available payload types and returns complete data.
func TestE2E_Commands_GetWithParameters_AllPayloadTypes(t *testing.T) {
	client := AuthenticateTestClient(t)

	t.Log("=== Test: GetCommandWithParameters across all payload types ===")

	// First get all payload types
	ctx1, cancel1 := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel1()

	payloadTypes, err := client.GetPayloadTypes(ctx1)
	require.NoError(t, err, "GetPayloadTypes should succeed")
	require.NotEmpty(t, payloadTypes, "Should have at least one payload type")

	t.Logf("✓ Found %d payload types", len(payloadTypes))

	// For each payload type, get a command with parameters
	for _, pt := range payloadTypes {
		t.Logf("--- Testing payload type: %s (ID:%d) ---", pt.Name, pt.ID)

		// Get all commands for this payload type
		ctx2, cancel2 := context.WithTimeout(context.Background(), 10*time.Second)
		commands, err := client.GetCommands(ctx2)
		cancel2()
		require.NoError(t, err)

		// Find a command for this payload type
		var testCommand *types.Command
		for _, cmd := range commands {
			if cmd.PayloadTypeID == pt.ID {
				testCommand = cmd
				break
			}
		}

		if testCommand == nil {
			t.Logf("⚠ No commands found for payload type %s", pt.Name)
			continue
		}

		// Get command with parameters
		ctx3, cancel3 := context.WithTimeout(context.Background(), 10*time.Second)
		cwp, err := client.GetCommandWithParameters(ctx3, pt.ID, testCommand.Cmd)
		cancel3()

		require.NoError(t, err, "GetCommandWithParameters should succeed for %s:%s", pt.Name, testCommand.Cmd)
		require.NotNil(t, cwp, "CommandWithParameters should not be nil")
		require.NotNil(t, cwp.Command, "Command should not be nil")

		// Validate command data
		assert.Equal(t, testCommand.Cmd, cwp.Command.Cmd, "Command name mismatch")
		assert.Equal(t, testCommand.ID, cwp.Command.ID, "Command ID mismatch")
		assert.Equal(t, pt.ID, cwp.Command.PayloadTypeID, "PayloadTypeID mismatch")

		// Validate helper methods
		isRaw := cwp.IsRawStringCommand()
		hasRequired := cwp.HasRequiredParameters()

		t.Logf("  ✓ Command: %s (RawString:%v, HasRequired:%v, ParamCount:%d)",
			cwp.Command.Cmd, isRaw, hasRequired, len(cwp.Parameters))

		// Validate parameters if present
		if len(cwp.Parameters) > 0 {
			for j, param := range cwp.Parameters {
				assert.NotEmpty(t, param.Name, "Parameter[%d] has empty Name", j)
				assert.NotEmpty(t, param.Type, "Parameter[%d] has empty Type", j)

				if j < 3 {
					t.Logf("    - Param: %s (Type:%s, Required:%v)", param.Name, param.Type, param.Required)
				}
			}
		}
	}

	t.Log("=== ✓ GetCommandWithParameters validation passed for all payload types ===")
}

// TestE2E_GetLoadedCommands_WithCallback validates that loaded command tracking works
// correctly after commands are loaded into a callback.
func TestE2E_Commands_GetLoaded_WithCallback(t *testing.T) {
	client := AuthenticateTestClient(t)

	t.Log("=== Test: GetLoadedCommands with active callback ===")

	// Find an active callback
	ctx1, cancel1 := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel1()

	callbacks, err := client.GetAllActiveCallbacks(ctx1)
	require.NoError(t, err, "GetAllActiveCallbacks should succeed")

	if len(callbacks) == 0 {
		t.Skip("⚠ No active callbacks found - skipping loaded commands test")
		return
	}

	testCallback := callbacks[0]
	t.Logf("✓ Using callback: %d (Host:%s, User:%s)", testCallback.ID, testCallback.Host, testCallback.User)

	// Get loaded commands for this callback
	ctx2, cancel2 := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel2()

	loadedCommands, err := client.GetLoadedCommands(ctx2, testCallback.ID)
	require.NoError(t, err, "GetLoadedCommands should succeed")

	t.Logf("✓ Callback has %d loaded commands", len(loadedCommands))

	if len(loadedCommands) == 0 {
		t.Log("⚠ No commands loaded yet - this is expected for new callbacks")
		return
	}

	// Validate loaded command data
	for i, lc := range loadedCommands {
		assert.NotZero(t, lc.ID, "LoadedCommand[%d] has ID 0", i)
		assert.NotZero(t, lc.CommandID, "LoadedCommand[%d] has CommandID 0", i)
		assert.Equal(t, testCallback.ID, lc.CallbackID, "LoadedCommand[%d] has wrong CallbackID", i)
		assert.NotZero(t, lc.OperatorID, "LoadedCommand[%d] has OperatorID 0", i)
		assert.GreaterOrEqual(t, lc.Version, 1, "LoadedCommand[%d] has invalid Version", i)

		// Command metadata should be populated
		require.NotNil(t, lc.Command, "LoadedCommand[%d] has nil Command", i)
		assert.NotEmpty(t, lc.Command.Cmd, "LoadedCommand[%d] Command has empty name", i)

		if i < 5 {
			t.Logf("  LoadedCommand[%d]: %s (Version:%d)", i, lc.Command.Cmd, lc.Version)
		}
	}

	t.Log("✓ All loaded commands have valid data")
	t.Log("=== ✓ GetLoadedCommands validation passed ===")
}

// TestE2E_BuildTaskParams_RawString validates BuildTaskParams for raw string commands
// (shell, run, execute) that don't have parameter definitions.
func TestE2E_Commands_BuildTaskParams_RawString(t *testing.T) {
	client := AuthenticateTestClient(t)

	t.Log("=== Test: BuildTaskParams for raw string commands ===")

	// Get commands and find a raw string command
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	commands, err := client.GetCommands(ctx)
	require.NoError(t, err)

	// Find ANY raw string command by checking if it has no parameters
	// Don't assume specific names like "shell" or "run" - check actual parameter count
	var cwp *mythic.CommandWithParameters
	for _, cmd := range commands {
		ctx2, cancel2 := context.WithTimeout(context.Background(), 10*time.Second)
		testCwp, err := client.GetCommandWithParameters(ctx2, cmd.PayloadTypeID, cmd.Cmd)
		cancel2()

		if err == nil && testCwp.IsRawStringCommand() {
			cwp = testCwp
			t.Logf("✓ Testing with command: %s (no parameters)", cmd.Cmd)
			break
		}
	}

	if cwp == nil {
		t.Skip("⚠ No raw string commands found - skipping test")
		return
	}

	// cwp is now a confirmed raw string command

	// Test 1: Simple string input
	params1, err := cwp.BuildTaskParams("whoami")
	require.NoError(t, err, "BuildTaskParams should accept simple string")
	assert.Equal(t, "whoami", params1, "Should return string as-is")
	t.Log("✓ Test 1: Simple string input works")

	// Test 2: String with spaces and special characters
	params2, err := cwp.BuildTaskParams("ls -la /tmp")
	require.NoError(t, err, "BuildTaskParams should accept string with spaces")
	assert.Equal(t, "ls -la /tmp", params2, "Should preserve spaces")
	t.Log("✓ Test 2: String with spaces works")

	// Test 3: Map input with "raw" key
	params3, err := cwp.BuildTaskParams(map[string]interface{}{
		"raw": "cat /etc/passwd",
	})
	require.NoError(t, err, "BuildTaskParams should accept map with 'raw' key")
	assert.Equal(t, "cat /etc/passwd", params3, "Should extract 'raw' key")
	t.Log("✓ Test 3: Map with 'raw' key works")

	// Test 4: Map input with "command" key
	params4, err := cwp.BuildTaskParams(map[string]interface{}{
		"command": "pwd",
	})
	require.NoError(t, err, "BuildTaskParams should accept map with 'command' key")
	assert.Equal(t, "pwd", params4, "Should extract 'command' key")
	t.Log("✓ Test 4: Map with 'command' key works")

	// Test 5: Map input without 'raw' or 'command' key should error
	_, err = cwp.BuildTaskParams(map[string]interface{}{
		"invalid": "test",
	})
	assert.Error(t, err, "BuildTaskParams should error for map without 'raw'/'command' key")
	t.Log("✓ Test 5: Invalid map input errors correctly")

	t.Log("=== ✓ BuildTaskParams raw string validation passed ===")
}

// TestE2E_BuildTaskParams_JSONParams validates BuildTaskParams for parameterized commands
// that expect structured JSON input.
func TestE2E_Commands_BuildTaskParams_JSONParams(t *testing.T) {
	client := AuthenticateTestClient(t)

	t.Log("=== Test: BuildTaskParams for parameterized commands ===")

	// Get commands and find a parameterized command
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	commands, err := client.GetCommands(ctx)
	require.NoError(t, err)

	// Find a command with parameters
	var paramCommand *types.Command
	for _, cmd := range commands {
		ctx2, cancel2 := context.WithTimeout(context.Background(), 5*time.Second)
		cwp, err := client.GetCommandWithParameters(ctx2, cmd.PayloadTypeID, cmd.Cmd)
		cancel2()

		if err == nil && len(cwp.Parameters) > 0 {
			paramCommand = cmd
			break
		}
	}

	if paramCommand == nil {
		t.Skip("⚠ No parameterized commands found - skipping test")
		return
	}

	t.Logf("✓ Testing with command: %s", paramCommand.Cmd)

	// Get command with parameters
	ctx3, cancel3 := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel3()

	cwp, err := client.GetCommandWithParameters(ctx3, paramCommand.PayloadTypeID, paramCommand.Cmd)
	require.NoError(t, err)
	require.False(t, cwp.IsRawStringCommand(), "Command should be parameterized type")
	require.Greater(t, len(cwp.Parameters), 0, "Should have parameters")

	t.Logf("  Command has %d parameters", len(cwp.Parameters))

	// Build input map with all parameters
	input := make(map[string]interface{})
	for _, param := range cwp.Parameters {
		// Provide test values based on type
		switch param.Type {
		case "String":
			input[param.Name] = "test_value"
		case "Boolean":
			input[param.Name] = true
		case "Number":
			input[param.Name] = 42
		case "Array":
			input[param.Name] = []string{"item1", "item2"}
		default:
			// Use default value if available
			if param.DefaultValue != "" {
				input[param.Name] = param.DefaultValue
			} else {
				input[param.Name] = "test"
			}
		}
		t.Logf("    - %s: %s = %v", param.Name, param.Type, input[param.Name])
	}

	// Test 1: Build params with all parameters
	paramsJSON, err := cwp.BuildTaskParams(input)
	require.NoError(t, err, "BuildTaskParams should succeed with valid input")
	assert.NotEmpty(t, paramsJSON, "Should return non-empty JSON")
	t.Log("✓ Test 1: BuildTaskParams with all parameters works")

	// Validate JSON structure
	var parsedParams map[string]interface{}
	err = json.Unmarshal([]byte(paramsJSON), &parsedParams)
	require.NoError(t, err, "Should produce valid JSON")
	assert.Equal(t, len(input), len(parsedParams), "All parameters should be in JSON")
	t.Logf("✓ Test 2: Produced valid JSON with %d fields", len(parsedParams))

	// Test 3: String input should error for parameterized commands
	_, err = cwp.BuildTaskParams("raw string")
	assert.Error(t, err, "BuildTaskParams should error for string input to parameterized command")
	t.Log("✓ Test 3: String input correctly errors for parameterized command")

	t.Log("=== ✓ BuildTaskParams JSON validation passed ===")
}

// TestE2E_BuildTaskParams_MissingRequired validates that BuildTaskParams correctly
// handles missing required parameters and provides helpful error messages.
func TestE2E_Commands_BuildTaskParams_MissingRequired(t *testing.T) {
	client := AuthenticateTestClient(t)

	t.Log("=== Test: BuildTaskParams error handling for missing required params ===")

	// Find a command with required parameters
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	commands, err := client.GetCommands(ctx)
	require.NoError(t, err)

	// Find a command with required parameters (no defaults)
	var requiredCommand *types.Command
	var requiredParam *types.CommandParameter

	for _, cmd := range commands {
		ctx2, cancel2 := context.WithTimeout(context.Background(), 5*time.Second)
		cwp, err := client.GetCommandWithParameters(ctx2, cmd.PayloadTypeID, cmd.Cmd)
		cancel2()

		if err == nil && cwp.HasRequiredParameters() {
			// Check if any required param lacks a default
			for _, param := range cwp.Parameters {
				if param.Required && param.DefaultValue == "" {
					requiredCommand = cmd
					requiredParam = param
					break
				}
			}
		}

		if requiredCommand != nil {
			break
		}
	}

	if requiredCommand == nil {
		t.Skip("⚠ No commands with required parameters (no defaults) found - skipping test")
		return
	}

	t.Logf("✓ Testing with command: %s (required param: %s)", requiredCommand.Cmd, requiredParam.Name)

	// Get command with parameters
	ctx3, cancel3 := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel3()

	cwp, err := client.GetCommandWithParameters(ctx3, requiredCommand.PayloadTypeID, requiredCommand.Cmd)
	require.NoError(t, err)

	// Test 1: Empty map should error
	_, err = cwp.BuildTaskParams(map[string]interface{}{})
	assert.Error(t, err, "BuildTaskParams should error with empty input")
	assert.Contains(t, err.Error(), "required parameter", "Error should mention required parameter")
	assert.Contains(t, err.Error(), requiredParam.Name, "Error should mention parameter name")
	t.Logf("✓ Test 1: Empty input errors: %v", err)

	// Test 2: Map with wrong parameters should error
	_, err = cwp.BuildTaskParams(map[string]interface{}{
		"wrong_param": "value",
	})
	assert.Error(t, err, "BuildTaskParams should error when required param missing")
	t.Logf("✓ Test 2: Wrong params error: %v", err)

	// Test 3: Providing required param should succeed
	validInput := map[string]interface{}{
		requiredParam.Name: "test_value",
	}
	// Add all other required params
	for _, param := range cwp.Parameters {
		if param.Required && param.Name != requiredParam.Name {
			if param.DefaultValue != "" {
				validInput[param.Name] = param.DefaultValue
			} else {
				validInput[param.Name] = "test"
			}
		}
	}

	params, err := cwp.BuildTaskParams(validInput)
	require.NoError(t, err, "BuildTaskParams should succeed with required params")
	assert.NotEmpty(t, params, "Should return valid JSON")
	t.Log("✓ Test 3: Valid input with required params succeeds")

	t.Log("=== ✓ BuildTaskParams error handling validation passed ===")
}
