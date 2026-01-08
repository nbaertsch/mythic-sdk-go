package integration

import (
	"context"
	"testing"

	"github.com/nbaertsch/mythic-sdk-go/pkg/mythic/types"
)

// TestCommands_GetCommands tests retrieving all available commands.
func TestCommands_GetCommands(t *testing.T) {
	client := getTestClient(t)
	ctx := context.Background()

	commands, err := client.GetCommands(ctx)
	if err != nil {
		t.Fatalf("GetCommands() failed: %v", err)
	}

	if len(commands) == 0 {
		t.Log("No commands found - this is expected if no payload types are loaded")
		return
	}

	t.Logf("Found %d commands", len(commands))

	// Verify command structure
	for i, cmd := range commands {
		if i < 5 { // Log first 5 for debugging
			t.Logf("Command %d: %s", i+1, cmd.String())
		}

		if cmd.ID == 0 {
			t.Errorf("Command %d has zero ID", i)
		}
		if cmd.Cmd == "" {
			t.Errorf("Command %d has empty Cmd field", i)
		}
		if cmd.PayloadTypeID == 0 {
			t.Errorf("Command %s has zero PayloadTypeID", cmd.Cmd)
		}
		if cmd.Version == 0 {
			t.Errorf("Command %s has zero Version", cmd.Cmd)
		}

		// Test helper methods
		_ = cmd.IsSupported()
		_ = cmd.IsScriptOnly()
		_ = cmd.String()
	}

	// Verify commands are sorted by name
	for i := 1; i < len(commands); i++ {
		if commands[i-1].Cmd > commands[i].Cmd {
			t.Errorf("Commands not sorted: %s comes after %s", commands[i-1].Cmd, commands[i].Cmd)
			break
		}
	}
}

// TestCommands_GetCommandsStructure tests command field values.
func TestCommands_GetCommandsStructure(t *testing.T) {
	client := getTestClient(t)
	ctx := context.Background()

	commands, err := client.GetCommands(ctx)
	if err != nil {
		t.Fatalf("GetCommands() failed: %v", err)
	}

	if len(commands) == 0 {
		t.Skip("No commands available for structure testing")
	}

	// Find a supported command
	var supportedCmd *types.Command
	for _, cmd := range commands {
		if cmd.Supported {
			supportedCmd = cmd
			break
		}
	}

	if supportedCmd == nil {
		t.Skip("No supported commands found")
	}

	t.Logf("Testing structure with command: %s", supportedCmd.Cmd)

	if supportedCmd.Description == "" {
		t.Log("Command has empty description - this may be intentional")
	}
	if supportedCmd.Help == "" {
		t.Log("Command has empty help text - this may be intentional")
	}
	if supportedCmd.Author == "" {
		t.Log("Command has empty author - this may be intentional")
	}

	// Test IsSupported method
	if !supportedCmd.IsSupported() {
		t.Error("IsSupported() returned false for a supported command")
	}
}

// TestCommands_GetCommandsScriptOnly tests script-only command filtering.
func TestCommands_GetCommandsScriptOnly(t *testing.T) {
	client := getTestClient(t)
	ctx := context.Background()

	commands, err := client.GetCommands(ctx)
	if err != nil {
		t.Fatalf("GetCommands() failed: %v", err)
	}

	if len(commands) == 0 {
		t.Skip("No commands available")
	}

	scriptOnlyCount := 0
	regularCount := 0

	for _, cmd := range commands {
		if cmd.ScriptOnly {
			scriptOnlyCount++
			if !cmd.IsScriptOnly() {
				t.Errorf("Command %s has ScriptOnly=true but IsScriptOnly() returned false", cmd.Cmd)
			}
		} else {
			regularCount++
			if cmd.IsScriptOnly() {
				t.Errorf("Command %s has ScriptOnly=false but IsScriptOnly() returned true", cmd.Cmd)
			}
		}
	}

	t.Logf("Found %d script-only commands and %d regular commands", scriptOnlyCount, regularCount)
}

// TestCommands_GetCommandParameters tests retrieving all command parameters.
func TestCommands_GetCommandParameters(t *testing.T) {
	client := getTestClient(t)
	ctx := context.Background()

	parameters, err := client.GetCommandParameters(ctx)
	if err != nil {
		t.Fatalf("GetCommandParameters() failed: %v", err)
	}

	if len(parameters) == 0 {
		t.Log("No command parameters found - this is expected if no payload types are loaded")
		return
	}

	t.Logf("Found %d command parameters", len(parameters))

	// Verify parameter structure
	for i, param := range parameters {
		if i < 5 { // Log first 5 for debugging
			t.Logf("Parameter %d: %s", i+1, param.String())
		}

		if param.ID == 0 {
			t.Errorf("Parameter %d has zero ID", i)
		}
		if param.CommandID == 0 {
			t.Errorf("Parameter %s has zero CommandID", param.Name)
		}
		if param.Name == "" {
			t.Errorf("Parameter %d has empty Name field", i)
		}
		if param.Type == "" {
			t.Errorf("Parameter %s has empty Type field", param.Name)
		}

		// Test helper methods
		_ = param.IsRequired()
		_ = param.HasChoices()
		_ = param.IsDynamic()
		_ = param.String()
	}

	// Verify parameters are sorted by command_id
	for i := 1; i < len(parameters); i++ {
		if parameters[i-1].CommandID > parameters[i].CommandID {
			t.Errorf("Parameters not sorted by command_id: %d comes after %d",
				parameters[i-1].CommandID, parameters[i].CommandID)
			break
		}
	}
}

// TestCommands_GetCommandParametersTypes tests parameter type validation.
func TestCommands_GetCommandParametersTypes(t *testing.T) {
	client := getTestClient(t)
	ctx := context.Background()

	parameters, err := client.GetCommandParameters(ctx)
	if err != nil {
		t.Fatalf("GetCommandParameters() failed: %v", err)
	}

	if len(parameters) == 0 {
		t.Skip("No command parameters available")
	}

	validTypes := map[string]bool{
		types.ParameterTypeString:         true,
		types.ParameterTypeBoolean:        true,
		types.ParameterTypeNumber:         true,
		types.ParameterTypeChooseOne:      true,
		types.ParameterTypeChooseMultiple: true,
		types.ParameterTypeFile:           true,
		types.ParameterTypeArray:          true,
		types.ParameterTypeCredential:     true,
		types.ParameterTypeLinkInfo:       true,
	}

	typeCount := make(map[string]int)
	for _, param := range parameters {
		typeCount[param.Type]++
		if !validTypes[param.Type] {
			t.Logf("Warning: Parameter %s has non-standard type: %s", param.Name, param.Type)
		}
	}

	t.Logf("Parameter type distribution:")
	for typ, count := range typeCount {
		t.Logf("  %s: %d", typ, count)
	}
}

// TestCommands_GetCommandParametersRequired tests required parameter identification.
func TestCommands_GetCommandParametersRequired(t *testing.T) {
	client := getTestClient(t)
	ctx := context.Background()

	parameters, err := client.GetCommandParameters(ctx)
	if err != nil {
		t.Fatalf("GetCommandParameters() failed: %v", err)
	}

	if len(parameters) == 0 {
		t.Skip("No command parameters available")
	}

	requiredCount := 0
	optionalCount := 0

	for _, param := range parameters {
		if param.Required {
			requiredCount++
			if !param.IsRequired() {
				t.Errorf("Parameter %s has Required=true but IsRequired() returned false", param.Name)
			}
		} else {
			optionalCount++
			if param.IsRequired() {
				t.Errorf("Parameter %s has Required=false but IsRequired() returned true", param.Name)
			}
		}
	}

	t.Logf("Found %d required parameters and %d optional parameters", requiredCount, optionalCount)
}

// TestCommands_GetCommandParametersChoices tests parameter choice types.
func TestCommands_GetCommandParametersChoices(t *testing.T) {
	client := getTestClient(t)
	ctx := context.Background()

	parameters, err := client.GetCommandParameters(ctx)
	if err != nil {
		t.Fatalf("GetCommandParameters() failed: %v", err)
	}

	if len(parameters) == 0 {
		t.Skip("No command parameters available")
	}

	staticChoices := 0
	allCommands := 0
	loadedCommands := 0
	noChoices := 0

	for _, param := range parameters {
		hasChoices := param.HasChoices()
		actualHasChoices := param.Choices != "" || param.ChoicesAreAllCommands || param.ChoicesAreLoadedCommands

		if hasChoices != actualHasChoices {
			t.Errorf("Parameter %s: HasChoices()=%v but actual=%v", param.Name, hasChoices, actualHasChoices)
		}

		switch {
		case param.Choices != "":
			staticChoices++
		case param.ChoicesAreAllCommands:
			allCommands++
		case param.ChoicesAreLoadedCommands:
			loadedCommands++
		default:
			noChoices++
		}
	}

	t.Logf("Choice type distribution:")
	t.Logf("  Static choices: %d", staticChoices)
	t.Logf("  All commands: %d", allCommands)
	t.Logf("  Loaded commands: %d", loadedCommands)
	t.Logf("  No choices: %d", noChoices)
}

// TestCommands_GetCommandParametersDynamic tests dynamic parameter identification.
func TestCommands_GetCommandParametersDynamic(t *testing.T) {
	client := getTestClient(t)
	ctx := context.Background()

	parameters, err := client.GetCommandParameters(ctx)
	if err != nil {
		t.Fatalf("GetCommandParameters() failed: %v", err)
	}

	if len(parameters) == 0 {
		t.Skip("No command parameters available")
	}

	dynamicCount := 0
	staticCount := 0

	for _, param := range parameters {
		if param.DynamicQueryFunction != "" {
			dynamicCount++
			if !param.IsDynamic() {
				t.Errorf("Parameter %s has DynamicQueryFunction but IsDynamic() returned false", param.Name)
			}
		} else {
			staticCount++
			if param.IsDynamic() {
				t.Errorf("Parameter %s has no DynamicQueryFunction but IsDynamic() returned true", param.Name)
			}
		}
	}

	t.Logf("Found %d dynamic parameters and %d static parameters", dynamicCount, staticCount)
}

// TestCommands_GetLoadedCommandsInvalidCallback tests error handling for invalid callback ID.
func TestCommands_GetLoadedCommandsInvalidCallback(t *testing.T) {
	client := getTestClient(t)
	ctx := context.Background()

	// Test with zero callback ID
	_, err := client.GetLoadedCommands(ctx, 0)
	if err == nil {
		t.Fatal("GetLoadedCommands(0) should return an error")
	}
	t.Logf("Zero callback ID error: %v", err)
}

// TestCommands_GetLoadedCommandsNonexistent tests behavior with nonexistent callback.
func TestCommands_GetLoadedCommandsNonexistent(t *testing.T) {
	client := getTestClient(t)
	ctx := context.Background()

	// Use a very high callback ID that likely doesn't exist
	loadedCommands, err := client.GetLoadedCommands(ctx, 999999)
	if err != nil {
		t.Fatalf("GetLoadedCommands(999999) failed: %v", err)
	}

	// Should return empty list, not an error
	if len(loadedCommands) != 0 {
		t.Errorf("Expected empty list for nonexistent callback, got %d commands", len(loadedCommands))
	}
}

// TestCommands_GetLoadedCommands tests retrieving loaded commands for a callback.
func TestCommands_GetLoadedCommands(t *testing.T) {
	client := getTestClient(t)
	ctx := context.Background()

	// First, get all callbacks to find one with loaded commands
	callbacks, err := client.GetCallbacks(ctx)
	if err != nil {
		t.Fatalf("GetCallbacks() failed: %v", err)
	}

	if len(callbacks) == 0 {
		t.Skip("No callbacks available - cannot test loaded commands")
	}

	// Test with the first callback
	callbackID := callbacks[0].ID
	t.Logf("Testing with callback ID: %d", callbackID)

	loadedCommands, err := client.GetLoadedCommands(ctx, callbackID)
	if err != nil {
		t.Fatalf("GetLoadedCommands(%d) failed: %v", callbackID, err)
	}

	if len(loadedCommands) == 0 {
		t.Log("No loaded commands found - this may be expected for this callback")
		return
	}

	t.Logf("Found %d loaded commands for callback %d", len(loadedCommands), callbackID)

	// Verify loaded command structure
	for i, lc := range loadedCommands {
		if i < 5 { // Log first 5 for debugging
			t.Logf("Loaded command %d: %s", i+1, lc.String())
		}

		if lc.ID == 0 {
			t.Errorf("LoadedCommand %d has zero ID", i)
		}
		if lc.CommandID == 0 {
			t.Errorf("LoadedCommand %d has zero CommandID", i)
		}
		if lc.CallbackID != callbackID {
			t.Errorf("LoadedCommand %d has CallbackID %d, expected %d", i, lc.CallbackID, callbackID)
		}
		if lc.Version == 0 {
			t.Errorf("LoadedCommand %d has zero Version", i)
		}

		// Verify nested command structure
		if lc.Command != nil {
			if lc.Command.Cmd == "" {
				t.Errorf("LoadedCommand %d has Command with empty Cmd", i)
			}
			if lc.Command.Version == 0 {
				t.Errorf("LoadedCommand %d has Command with zero Version", i)
			}
			if lc.Command.Version != lc.Version {
				t.Logf("Warning: LoadedCommand %d version (%d) differs from Command version (%d)",
					i, lc.Version, lc.Command.Version)
			}
		}

		_ = lc.String()
	}

	// Verify loaded commands are sorted by command name
	for i := 1; i < len(loadedCommands); i++ {
		if loadedCommands[i-1].Command != nil && loadedCommands[i].Command != nil {
			if loadedCommands[i-1].Command.Cmd > loadedCommands[i].Command.Cmd {
				t.Errorf("Loaded commands not sorted: %s comes after %s",
					loadedCommands[i-1].Command.Cmd, loadedCommands[i].Command.Cmd)
				break
			}
		}
	}
}

// TestCommands_GetLoadedCommandsMultipleCallbacks tests loaded commands across callbacks.
func TestCommands_GetLoadedCommandsMultipleCallbacks(t *testing.T) {
	client := getTestClient(t)
	ctx := context.Background()

	callbacks, err := client.GetCallbacks(ctx)
	if err != nil {
		t.Fatalf("GetCallbacks() failed: %v", err)
	}

	if len(callbacks) < 2 {
		t.Skip("Need at least 2 callbacks to test multiple callback scenario")
	}

	// Test first 3 callbacks (or all if less than 3)
	maxCallbacks := 3
	if len(callbacks) < maxCallbacks {
		maxCallbacks = len(callbacks)
	}

	for i := 0; i < maxCallbacks; i++ {
		callbackID := callbacks[i].ID
		loadedCommands, err := client.GetLoadedCommands(ctx, callbackID)
		if err != nil {
			t.Errorf("GetLoadedCommands(%d) failed: %v", callbackID, err)
			continue
		}

		t.Logf("Callback %d has %d loaded commands", callbackID, len(loadedCommands))

		// Verify all loaded commands belong to this callback
		for _, lc := range loadedCommands {
			if lc.CallbackID != callbackID {
				t.Errorf("LoadedCommand has CallbackID %d but queried for %d", lc.CallbackID, callbackID)
			}
		}
	}
}

// TestCommands_CommandParameterRelationship tests the relationship between commands and parameters.
func TestCommands_CommandParameterRelationship(t *testing.T) {
	client := getTestClient(t)
	ctx := context.Background()

	commands, err := client.GetCommands(ctx)
	if err != nil {
		t.Fatalf("GetCommands() failed: %v", err)
	}

	parameters, err := client.GetCommandParameters(ctx)
	if err != nil {
		t.Fatalf("GetCommandParameters() failed: %v", err)
	}

	if len(commands) == 0 || len(parameters) == 0 {
		t.Skip("Need both commands and parameters for relationship testing")
	}

	// Build command ID map
	commandMap := make(map[int]*types.Command)
	for _, cmd := range commands {
		commandMap[cmd.ID] = cmd
	}

	// Group parameters by command
	paramsByCommand := make(map[int][]*types.CommandParameter)
	for _, param := range parameters {
		paramsByCommand[param.CommandID] = append(paramsByCommand[param.CommandID], param)
	}

	// Verify all parameters reference valid commands
	for _, param := range parameters {
		if _, exists := commandMap[param.CommandID]; !exists {
			t.Errorf("Parameter %s references nonexistent command ID %d", param.Name, param.CommandID)
		}
	}

	// Log statistics
	t.Logf("Total commands: %d", len(commands))
	t.Logf("Total parameters: %d", len(parameters))
	t.Logf("Commands with parameters: %d", len(paramsByCommand))

	commandsWithoutParams := 0
	for _, cmd := range commands {
		if _, hasParams := paramsByCommand[cmd.ID]; !hasParams {
			commandsWithoutParams++
		}
	}
	t.Logf("Commands without parameters: %d", commandsWithoutParams)
}
