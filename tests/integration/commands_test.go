//go:build integration

package integration

import (
	"context"
	"testing"
	"time"

	"github.com/nbaertsch/mythic-sdk-go/pkg/mythic/types"
)

// TestE2E_CommandDiscovery tests comprehensive command discovery and metadata retrieval.
// Covers: GetCommands, GetCommandParameters, GetCommandWithParameters, GetLoadedCommands
func TestE2E_CommandDiscovery(t *testing.T) {
	client := AuthenticateTestClient(t)

	var testCommand *types.Command
	var testCallback *types.Callback

	// Test 1: Get all commands
	t.Log("=== Test 1: Get all commands ===")
	ctx1, cancel1 := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel1()

	commands, err := client.GetCommands(ctx1)
	if err != nil {
		t.Fatalf("GetCommands failed: %v", err)
	}
	if len(commands) == 0 {
		t.Fatal("No commands found")
	}
	t.Logf("✓ Found %d commands", len(commands))

	// Validate command structure
	foundShellCommand := false
	for _, cmd := range commands {
		if cmd.ID == 0 {
			t.Error("Command has ID 0")
		}
		if cmd.Cmd == "" {
			t.Error("Command has empty name")
		}
		if cmd.PayloadTypeID == 0 {
			t.Error("Command has PayloadTypeID 0")
		}

		// Track a common command for later tests
		if cmd.Cmd == "shell" || cmd.Cmd == "run" || cmd.Cmd == "execute" {
			testCommand = cmd
			foundShellCommand = true
			t.Logf("  Found raw string command: %s (ID: %d, PayloadType: %d)", cmd.Cmd, cmd.ID, cmd.PayloadTypeID)
		}
	}

	if !foundShellCommand {
		// Use first command as fallback
		testCommand = commands[0]
		t.Logf("  Using first command as test command: %s (ID: %d)", testCommand.Cmd, testCommand.ID)
	}

	// Test 2: Get all command parameters
	t.Log("=== Test 2: Get all command parameters ===")
	ctx2, cancel2 := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel2()

	parameters, err := client.GetCommandParameters(ctx2)
	if err != nil {
		t.Fatalf("GetCommandParameters failed: %v", err)
	}
	t.Logf("✓ Found %d command parameters across all commands", len(parameters))

	// Validate parameter structure
	for _, param := range parameters {
		if param.ID == 0 {
			t.Error("Parameter has ID 0")
		}
		if param.CommandID == 0 {
			t.Error("Parameter has CommandID 0")
		}
		if param.Name == "" {
			t.Error("Parameter has empty name")
		}
		if param.Type == "" {
			t.Error("Parameter has empty type")
		}
	}

	// Test 3: Get command with parameters
	t.Log("=== Test 3: Get command with parameters ===")
	ctx3, cancel3 := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel3()

	cmdWithParams, err := client.GetCommandWithParameters(ctx3, testCommand.PayloadTypeID, testCommand.Cmd)
	if err != nil {
		t.Fatalf("GetCommandWithParameters failed: %v", err)
	}
	if cmdWithParams.Command == nil {
		t.Fatal("GetCommandWithParameters returned nil Command")
	}
	if cmdWithParams.Command.Cmd != testCommand.Cmd {
		t.Errorf("Command name mismatch: expected %s, got %s", testCommand.Cmd, cmdWithParams.Command.Cmd)
	}
	t.Logf("✓ Command '%s' retrieved with %d parameters", cmdWithParams.Command.Cmd, len(cmdWithParams.Parameters))

	// Test 4: Validate command parameter metadata
	t.Log("=== Test 4: Validate command parameter metadata ===")
	if len(cmdWithParams.Parameters) > 0 {
		t.Logf("  Command '%s' has %d parameters:", cmdWithParams.Command.Cmd, len(cmdWithParams.Parameters))
		for _, param := range cmdWithParams.Parameters {
			t.Logf("    - %s (%s) [required: %v]", param.Name, param.Type, param.Required)
			if param.Name == "" {
				t.Error("Parameter has empty name")
			}
			if param.Type == "" {
				t.Error("Parameter has empty type")
			}
		}
	} else {
		t.Logf("  Command '%s' is a raw string command (no parameters)", cmdWithParams.Command.Cmd)
	}
	t.Log("✓ Parameter metadata validated")

	// Test 5: Get active callbacks for loaded commands test
	t.Log("=== Test 5: Get active callback for loaded commands test ===")
	ctx5, cancel5 := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel5()

	callbacks, err := client.GetAllActiveCallbacks(ctx5)
	if err != nil {
		t.Fatalf("GetCallbacks failed: %v", err)
	}

	// Find an active callback
	for _, cb := range callbacks {
		if cb.Active {
			testCallback = cb
			break
		}
	}

	if testCallback == nil {
		t.Log("⚠ No active callbacks found, skipping loaded commands test")
		t.Log("=== ✓ Command discovery tests passed (partial) ===")
		return
	}
	t.Logf("✓ Found active callback: %d (Host: %s, User: %s)", testCallback.ID, testCallback.Host, testCallback.User)

	// Test 6: Get loaded commands for callback
	t.Log("=== Test 6: Get loaded commands for callback ===")
	ctx6, cancel6 := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel6()

	loadedCommands, err := client.GetLoadedCommands(ctx6, testCallback.ID)
	if err != nil {
		t.Fatalf("GetLoadedCommands failed: %v", err)
	}
	t.Logf("✓ Callback %d has %d loaded commands", testCallback.ID, len(loadedCommands))

	// Validate loaded commands
	for _, lc := range loadedCommands {
		if lc.ID == 0 {
			t.Error("LoadedCommand has ID 0")
		}
		if lc.CommandID == 0 {
			t.Error("LoadedCommand has CommandID 0")
		}
		if lc.CallbackID != testCallback.ID {
			t.Errorf("LoadedCommand has wrong CallbackID: expected %d, got %d", testCallback.ID, lc.CallbackID)
		}
		if lc.Command == nil {
			t.Error("LoadedCommand has nil Command")
		} else if lc.Command.Cmd == "" {
			t.Error("LoadedCommand has Command with empty name")
		}
	}

	if len(loadedCommands) > 0 {
		t.Logf("  Example loaded commands:")
		count := 5
		if len(loadedCommands) < count {
			count = len(loadedCommands)
		}
		for i := 0; i < count; i++ {
			lc := loadedCommands[i]
			t.Logf("    - %s (Version: %d)", lc.Command.Cmd, lc.Version)
		}
	}

	t.Log("=== ✓ All command discovery tests passed ===")
}

// TestE2E_CommandWithParametersHelpers tests helper methods on CommandWithParameters.
func TestE2E_CommandWithParametersHelpers(t *testing.T) {
	client := AuthenticateTestClient(t)

	// Test 1: Find a raw string command (no parameters)
	t.Log("=== Test 1: Test raw string command helpers ===")
	ctx1, cancel1 := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel1()

	commands, err := client.GetCommands(ctx1)
	if err != nil {
		t.Fatalf("GetCommands failed: %v", err)
	}

	var rawStringCommand *types.Command
	var parameterizedCommand *types.Command

	// Find examples of both types
	for _, cmd := range commands {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		cmdWithParams, err := client.GetCommandWithParameters(ctx, cmd.PayloadTypeID, cmd.Cmd)
		cancel()

		if err != nil {
			continue
		}

		if len(cmdWithParams.Parameters) == 0 && rawStringCommand == nil {
			rawStringCommand = cmd
		}
		if len(cmdWithParams.Parameters) > 0 && parameterizedCommand == nil {
			parameterizedCommand = cmd
		}

		if rawStringCommand != nil && parameterizedCommand != nil {
			break
		}
	}

	// Test raw string command
	if rawStringCommand != nil {
		ctx2, cancel2 := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel2()

		cwp, err := client.GetCommandWithParameters(ctx2, rawStringCommand.PayloadTypeID, rawStringCommand.Cmd)
		if err != nil {
			t.Fatalf("GetCommandWithParameters failed for raw string command: %v", err)
		}

		if !cwp.IsRawStringCommand() {
			t.Errorf("Expected %s to be a raw string command", rawStringCommand.Cmd)
		}
		t.Logf("✓ Raw string command detected: %s", rawStringCommand.Cmd)

		// Test BuildTaskParams with raw string input
		params, err := cwp.BuildTaskParams("whoami")
		if err != nil {
			t.Errorf("BuildTaskParams failed for raw string: %v", err)
		}
		if params != "whoami" {
			t.Errorf("BuildTaskParams returned unexpected value: got %s, expected 'whoami'", params)
		}
		t.Logf("✓ BuildTaskParams for raw string: '%s'", params)
	} else {
		t.Log("⚠ No raw string commands found")
	}

	// Test 2: Test parameterized command
	if parameterizedCommand != nil {
		t.Log("=== Test 2: Test parameterized command helpers ===")
		ctx3, cancel3 := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel3()

		cwp, err := client.GetCommandWithParameters(ctx3, parameterizedCommand.PayloadTypeID, parameterizedCommand.Cmd)
		if err != nil {
			t.Fatalf("GetCommandWithParameters failed for parameterized command: %v", err)
		}

		if cwp.IsRawStringCommand() {
			t.Errorf("Expected %s to be a parameterized command", parameterizedCommand.Cmd)
		}
		t.Logf("✓ Parameterized command detected: %s (%d parameters)", parameterizedCommand.Cmd, len(cwp.Parameters))

		// Check for required parameters
		hasRequired := cwp.HasRequiredParameters()
		t.Logf("  Has required parameters: %v", hasRequired)

		// Test BuildTaskParams with valid input
		testParams := make(map[string]interface{})
		for _, param := range cwp.Parameters {
			if param.Required {
				// Provide a test value based on type
				switch param.Type {
				case "String":
					testParams[param.Name] = "test_value"
				case "Boolean":
					testParams[param.Name] = true
				case "Number":
					testParams[param.Name] = 1
				default:
					testParams[param.Name] = "test"
				}
			}
		}

		if len(testParams) > 0 {
			params, err := cwp.BuildTaskParams(testParams)
			if err != nil {
				t.Errorf("BuildTaskParams failed: %v", err)
			} else {
				t.Logf("✓ BuildTaskParams for parameterized command: %s", params)
			}
		}
	} else {
		t.Log("⚠ No parameterized commands found")
	}

	t.Log("=== ✓ Command helper tests passed ===")
}

// TestE2E_CommandErrorHandling tests error scenarios for command operations.
func TestE2E_CommandErrorHandling(t *testing.T) {
	client := AuthenticateTestClient(t)

	// Test 1: Get command with invalid payload type
	t.Log("=== Test 1: Get command with invalid payload type ===")
	ctx1, cancel1 := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel1()

	_, err := client.GetCommandWithParameters(ctx1, 999999, "nonexistent")
	if err == nil {
		t.Error("Expected error for non-existent command")
	}
	t.Logf("✓ Non-existent command rejected: %v", err)

	// Test 2: Get command with empty name
	t.Log("=== Test 2: Get command with empty name ===")
	ctx2, cancel2 := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel2()

	_, err = client.GetCommandWithParameters(ctx2, 1, "")
	if err == nil {
		t.Error("Expected error for empty command name")
	}
	t.Logf("✓ Empty command name rejected: %v", err)

	// Test 3: Get loaded commands with invalid callback ID
	t.Log("=== Test 3: Get loaded commands with invalid callback ID ===")
	ctx3, cancel3 := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel3()

	_, err = client.GetLoadedCommands(ctx3, 0)
	if err == nil {
		t.Error("Expected error for invalid callback ID")
	}
	t.Logf("✓ Invalid callback ID rejected: %v", err)

	// Test 4: BuildTaskParams with missing required parameters
	t.Log("=== Test 4: BuildTaskParams with missing required parameters ===")
	ctx4, cancel4 := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel4()

	commands, err := client.GetCommands(ctx4)
	if err != nil {
		t.Fatalf("GetCommands failed: %v", err)
	}

	// Find a command with required parameters
	for _, cmd := range commands {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		cwp, err := client.GetCommandWithParameters(ctx, cmd.PayloadTypeID, cmd.Cmd)
		cancel()

		if err != nil || !cwp.HasRequiredParameters() {
			continue
		}

		// Try to build params without providing required parameters
		_, err = cwp.BuildTaskParams(map[string]interface{}{})
		if err == nil {
			// Check if all required params have defaults
			allHaveDefaults := true
			for _, param := range cwp.Parameters {
				if param.Required && param.DefaultValue == "" {
					allHaveDefaults = false
					break
				}
			}
			if !allHaveDefaults {
				t.Errorf("Expected error for missing required parameters in command %s", cmd.Cmd)
			}
		} else {
			t.Logf("✓ Missing required parameters rejected for command %s: %v", cmd.Cmd, err)
		}
		break
	}

	t.Log("=== ✓ All error handling tests passed ===")
}

// TestE2E_CommandAttributes tests command attribute parsing and filtering.
func TestE2E_CommandAttributes(t *testing.T) {
	client := AuthenticateTestClient(t)

	t.Log("=== Test: Command attribute analysis ===")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	commands, err := client.GetCommands(ctx)
	if err != nil {
		t.Fatalf("GetCommands failed: %v", err)
	}

	// Analyze command attributes
	attributeStats := make(map[string]int)
	scriptOnlyCount := 0
	withMitreCount := 0

	for _, cmd := range commands {
		if cmd.ScriptOnly {
			scriptOnlyCount++
		}
		if cmd.MitreAttackMappings != "" && cmd.MitreAttackMappings != "[]" {
			withMitreCount++
		}
		if cmd.Attributes != "" {
			attributeStats[cmd.Attributes]++
		}
	}

	t.Logf("✓ Command statistics:")
	t.Logf("  Total commands: %d", len(commands))
	t.Logf("  Script-only commands: %d", scriptOnlyCount)
	t.Logf("  Commands with MITRE ATT&CK mappings: %d", withMitreCount)
	t.Logf("  Unique attribute patterns: %d", len(attributeStats))

	t.Log("=== ✓ Command attribute analysis complete ===")
}
