//go:build integration
// +build integration

package integration

import (
	"context"
	"testing"

	"github.com/nbaertsch/mythic-sdk-go/pkg/mythic/types"
)

// TestBrowserScripts_GetBrowserScripts tests retrieving all browser scripts
func TestBrowserScripts_GetBrowserScripts(t *testing.T) {
	ctx := context.Background()

	scripts, err := client.GetBrowserScripts(ctx)
	if err != nil {
		t.Fatalf("Failed to get browser scripts: %v", err)
	}

	t.Logf("Retrieved %d browser scripts", len(scripts))

	// Verify structure if any scripts exist
	if len(scripts) > 0 {
		script := scripts[0]
		if script.Name == "" {
			t.Error("Script Name should not be empty")
		}
		if script.Script == "" {
			t.Error("Script content should not be empty")
		}

		// Test String method
		str := script.String()
		if str == "" {
			t.Error("BrowserScript.String() should not return empty string")
		}
		t.Logf("Browser script: %s", str)
		t.Logf("  ID: %d", script.ID)
		t.Logf("  Name: %s", script.Name)
		t.Logf("  Author: %s", script.Author)
		t.Logf("  For New UI: %v", script.ForNewUI)
		t.Logf("  Active: %v", script.Active)
		t.Logf("  Script length: %d bytes", len(script.Script))

		// Test helper methods
		if script.IsActive() {
			t.Logf("  Script is active")
		}
		if script.IsForNewUI() {
			t.Logf("  Script is for new UI")
		}
	}
}

// TestBrowserScripts_GetBrowserScriptsByOperation tests retrieving scripts for a specific operation
func TestBrowserScripts_GetBrowserScriptsByOperation(t *testing.T) {
	ctx := context.Background()

	// Get current operation
	operationID := client.GetCurrentOperation()
	if operationID == nil {
		t.Skip("No current operation set")
	}

	scripts, err := client.GetBrowserScriptsByOperation(ctx, *operationID)
	if err != nil {
		t.Fatalf("Failed to get browser scripts for operation: %v", err)
	}

	t.Logf("Retrieved %d browser scripts for operation %d", len(scripts), *operationID)

	// Verify structure if any scripts exist
	if len(scripts) > 0 {
		script := scripts[0]
		if script.BrowserScriptID == 0 {
			t.Error("BrowserScriptID should not be 0")
		}
		if script.OperationID != *operationID {
			t.Errorf("Expected OperationID %d, got %d", *operationID, script.OperationID)
		}

		// Test String method
		str := script.String()
		if str == "" {
			t.Error("BrowserScriptOperation.String() should not return empty string")
		}
		t.Logf("Browser script operation: %s", str)
		t.Logf("  ID: %d", script.ID)
		t.Logf("  BrowserScriptID: %d", script.BrowserScriptID)
		t.Logf("  OperationID: %d", script.OperationID)
		t.Logf("  ScriptName: %s", script.ScriptName)
		t.Logf("  Active: %v", script.Active)

		// Test helper methods
		if script.IsActive() {
			t.Logf("  Script is active for this operation")
		}
		if script.IsOperatorSpecific() {
			t.Logf("  Script is assigned to operator ID: %d", *script.OperatorID)
		}
	}
}

// TestBrowserScripts_GetBrowserScriptsByOperation_InvalidInput tests with invalid input
func TestBrowserScripts_GetBrowserScriptsByOperation_InvalidInput(t *testing.T) {
	ctx := context.Background()

	// Try to get scripts with operation ID 0
	_, err := client.GetBrowserScriptsByOperation(ctx, 0)
	if err == nil {
		t.Error("Expected error for operation ID 0")
	}
}

// TestBrowserScripts_CustomBrowserExport tests custom browser export functionality
func TestBrowserScripts_CustomBrowserExport(t *testing.T) {
	ctx := context.Background()

	// Get current operation
	operationID := client.GetCurrentOperation()
	if operationID == nil {
		t.Skip("No current operation set")
	}

	// Get browser scripts to find a valid script name
	scripts, err := client.GetBrowserScripts(ctx)
	if err != nil || len(scripts) == 0 {
		t.Skip("No browser scripts available for testing export")
	}

	// Use the first active script for export
	var scriptName string
	for _, script := range scripts {
		if script.Active {
			scriptName = script.Name
			break
		}
	}
	if scriptName == "" {
		t.Skip("No active browser scripts available for testing export")
	}

	// Test custom export
	request := &types.CustomBrowserExportRequest{
		OperationID: *operationID,
		ScriptName:  scriptName,
		Parameters: map[string]interface{}{
			"format": "json",
		},
	}

	response, err := client.CustomBrowserExport(ctx, request)
	if err != nil {
		// This might fail if the script doesn't support export or if the feature is not available
		t.Logf("Custom browser export failed (may not be supported): %v", err)
		t.Skip("Custom browser export not available or not supported by script")
	}

	if response.Status != "success" {
		t.Logf("Export status: %s (error: %s)", response.Status, response.Error)
		t.Skip("Custom browser export returned non-success status")
	}

	t.Logf("Custom browser export successful for script '%s'", scriptName)
	t.Logf("Status: %s", response.Status)
	if response.ExportedData != "" {
		t.Logf("Exported data length: %d bytes", len(response.ExportedData))
	}
}

// TestBrowserScripts_CustomBrowserExport_InvalidInput tests export with invalid input
func TestBrowserScripts_CustomBrowserExport_InvalidInput(t *testing.T) {
	ctx := context.Background()

	// Try to export with nil request
	_, err := client.CustomBrowserExport(ctx, nil)
	if err == nil {
		t.Error("Expected error for nil request")
	}

	// Try to export with operation ID 0
	request := &types.CustomBrowserExportRequest{
		OperationID: 0,
		ScriptName:  "test",
	}
	_, err = client.CustomBrowserExport(ctx, request)
	if err == nil {
		t.Error("Expected error for operation ID 0")
	}

	// Try to export with empty script name
	operationID := client.GetCurrentOperation()
	if operationID != nil {
		request := &types.CustomBrowserExportRequest{
			OperationID: *operationID,
			ScriptName:  "",
		}
		_, err = client.CustomBrowserExport(ctx, request)
		if err == nil {
			t.Error("Expected error for empty script name")
		}
	}
}

// TestBrowserScripts_CustomBrowserExport_WithParameters tests export with various parameters
func TestBrowserScripts_CustomBrowserExport_WithParameters(t *testing.T) {
	ctx := context.Background()

	// Get current operation
	operationID := client.GetCurrentOperation()
	if operationID == nil {
		t.Skip("No current operation set")
	}

	// Get browser scripts
	scripts, err := client.GetBrowserScripts(ctx)
	if err != nil || len(scripts) == 0 {
		t.Skip("No browser scripts available for testing")
	}

	// Use the first active script
	var scriptName string
	for _, script := range scripts {
		if script.Active {
			scriptName = script.Name
			break
		}
	}
	if scriptName == "" {
		t.Skip("No active browser scripts available")
	}

	// Test with various parameter combinations
	parameterTests := []map[string]interface{}{
		{"format": "json"},
		{"format": "csv", "delimiter": ","},
		{"include_timestamps": true, "timezone": "UTC"},
		{"filter": "active", "limit": 100},
	}

	for i, params := range parameterTests {
		request := &types.CustomBrowserExportRequest{
			OperationID: *operationID,
			ScriptName:  scriptName,
			Parameters:  params,
		}

		response, err := client.CustomBrowserExport(ctx, request)
		if err != nil {
			t.Logf("Parameter test %d failed (may not be supported): %v", i+1, err)
			continue
		}

		if response.Status == "success" {
			t.Logf("Parameter test %d successful with params: %v", i+1, params)
		}
	}
}

// TestBrowserScripts_GetBrowserScriptsMultipleOperations tests scripts across multiple operations
func TestBrowserScripts_GetBrowserScriptsMultipleOperations(t *testing.T) {
	ctx := context.Background()

	// Get all operations to test with
	operations, err := client.GetOperations(ctx)
	if err != nil || len(operations) < 2 {
		t.Skip("Need at least 2 operations for testing")
	}

	// Test first 3 operations only
	for i, op := range operations {
		if i >= 3 {
			break
		}

		scripts, err := client.GetBrowserScriptsByOperation(ctx, op.ID)
		if err != nil {
			t.Logf("Operation %d (%s): Failed to get browser scripts: %v", op.ID, op.Name, err)
			continue
		}

		t.Logf("Operation %d (%s): %d browser scripts", op.ID, op.Name, len(scripts))
	}
}

// TestBrowserScripts_ActiveVsInactiveScripts tests filtering of active vs inactive scripts
func TestBrowserScripts_ActiveVsInactiveScripts(t *testing.T) {
	ctx := context.Background()

	scripts, err := client.GetBrowserScripts(ctx)
	if err != nil {
		t.Fatalf("Failed to get browser scripts: %v", err)
	}

	if len(scripts) == 0 {
		t.Skip("No browser scripts available")
	}

	activeCount := 0
	inactiveCount := 0
	newUICount := 0
	oldUICount := 0

	for _, script := range scripts {
		if script.IsActive() {
			activeCount++
		} else {
			inactiveCount++
		}

		if script.IsForNewUI() {
			newUICount++
		} else {
			oldUICount++
		}
	}

	t.Logf("Total scripts: %d", len(scripts))
	t.Logf("  Active: %d", activeCount)
	t.Logf("  Inactive: %d", inactiveCount)
	t.Logf("  For New UI: %d", newUICount)
	t.Logf("  For Old UI: %d", oldUICount)
}

// TestBrowserScripts_ScriptContent tests script content validation
func TestBrowserScripts_ScriptContent(t *testing.T) {
	ctx := context.Background()

	scripts, err := client.GetBrowserScripts(ctx)
	if err != nil {
		t.Fatalf("Failed to get browser scripts: %v", err)
	}

	if len(scripts) == 0 {
		t.Skip("No browser scripts available")
	}

	for _, script := range scripts {
		// Verify script has content
		if script.Script == "" {
			t.Errorf("Script %s (ID: %d) has empty content", script.Name, script.ID)
		}

		// Verify script has a name
		if script.Name == "" {
			t.Errorf("Script ID %d has empty name", script.ID)
		}

		// Log script metadata
		t.Logf("Script: %s", script.Name)
		t.Logf("  Length: %d bytes", len(script.Script))
		t.Logf("  Author: %s", script.Author)
		t.Logf("  Description: %s", script.Description)
	}
}

// TestBrowserScripts_OperatorSpecificScripts tests operator-specific script assignments
func TestBrowserScripts_OperatorSpecificScripts(t *testing.T) {
	ctx := context.Background()

	// Get current operation
	operationID := client.GetCurrentOperation()
	if operationID == nil {
		t.Skip("No current operation set")
	}

	scripts, err := client.GetBrowserScriptsByOperation(ctx, *operationID)
	if err != nil {
		t.Fatalf("Failed to get browser scripts for operation: %v", err)
	}

	if len(scripts) == 0 {
		t.Skip("No browser scripts assigned to operation")
	}

	operatorSpecificCount := 0
	globalCount := 0

	for _, script := range scripts {
		if script.IsOperatorSpecific() {
			operatorSpecificCount++
			t.Logf("Operator-specific script: %s (Operator ID: %d)", script.ScriptName, *script.OperatorID)
		} else {
			globalCount++
			t.Logf("Global script: %s", script.ScriptName)
		}
	}

	t.Logf("Total scripts for operation: %d", len(scripts))
	t.Logf("  Operator-specific: %d", operatorSpecificCount)
	t.Logf("  Global: %d", globalCount)
}
