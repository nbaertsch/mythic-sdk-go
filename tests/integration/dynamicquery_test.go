//go:build integration

package integration

import (
	"context"
	"testing"
)

// TestDynamicQueryFunction_InvalidInput tests input validation for DynamicQueryFunction.
func TestDynamicQueryFunction_InvalidInput(t *testing.T) {
	client := getTestClient(t)
	ctx := context.Background()

	// Test with empty command
	_, err := client.DynamicQueryFunction(ctx, "", map[string]interface{}{"key": "value"}, 0)
	if err == nil {
		t.Fatal("DynamicQueryFunction with empty command should return error")
	}
	t.Logf("Empty command error: %v", err)
}

// TestDynamicQueryFunction_NonexistentCommand tests querying a command that doesn't exist.
func TestDynamicQueryFunction_NonexistentCommand(t *testing.T) {
	client := getTestClient(t)
	ctx := context.Background()

	// Try to query a command that likely doesn't exist
	params := map[string]interface{}{
		"test_param": "value",
	}
	result, err := client.DynamicQueryFunction(ctx, "nonexistent_command_12345", params, 0)
	if err != nil {
		t.Logf("Nonexistent command error (expected): %v", err)
	}

	if result != nil {
		t.Logf("Response: %s", result.String())
		if result.IsSuccessful() {
			t.Log("Query succeeded (command might exist)")
			if result.HasChoices() {
				t.Logf("Returned %d choices", len(result.Choices))
			}
		}
	}
}

// TestDynamicQueryFunction_WithCallbackID tests dynamic query with callback context.
func TestDynamicQueryFunction_WithCallbackID(t *testing.T) {
	client := getTestClient(t)
	ctx := context.Background()

	// Test with a callback ID (will likely fail since callback doesn't exist)
	params := map[string]interface{}{
		"path": "/tmp",
	}
	result, err := client.DynamicQueryFunction(ctx, "ls", params, 999999)
	if err != nil {
		t.Logf("Query with callback ID error (expected if command/callback doesn't exist): %v", err)
	}

	if result != nil {
		t.Logf("Response: %s", result.String())
	}
}

// TestDynamicQueryFunction_EmptyParameters tests dynamic query with empty parameters.
func TestDynamicQueryFunction_EmptyParameters(t *testing.T) {
	client := getTestClient(t)
	ctx := context.Background()

	// Test with empty parameters map
	result, err := client.DynamicQueryFunction(ctx, "test_command", map[string]interface{}{}, 0)
	if err != nil {
		t.Logf("Empty parameters error (expected if command doesn't exist): %v", err)
	}

	if result != nil {
		t.Logf("Response: %s", result.String())
	}

	// Test with nil parameters
	result, err = client.DynamicQueryFunction(ctx, "test_command", nil, 0)
	if err != nil {
		t.Logf("Nil parameters error (expected if command doesn't exist): %v", err)
	}

	if result != nil {
		t.Logf("Response: %s", result.String())
	}
}

// TestDynamicBuildParameter_InvalidInput tests input validation for DynamicBuildParameter.
func TestDynamicBuildParameter_InvalidInput(t *testing.T) {
	client := getTestClient(t)
	ctx := context.Background()

	// Test with empty payload type
	_, err := client.DynamicBuildParameter(ctx, "", "parameter", nil)
	if err == nil {
		t.Fatal("DynamicBuildParameter with empty payload type should return error")
	}
	t.Logf("Empty payload type error: %v", err)

	// Test with empty parameter name
	_, err = client.DynamicBuildParameter(ctx, "apollo", "", nil)
	if err == nil {
		t.Fatal("DynamicBuildParameter with empty parameter should return error")
	}
	t.Logf("Empty parameter error: %v", err)
}

// TestDynamicBuildParameter_NonexistentPayloadType tests querying a payload type that doesn't exist.
func TestDynamicBuildParameter_NonexistentPayloadType(t *testing.T) {
	client := getTestClient(t)
	ctx := context.Background()

	// Try to query a payload type that doesn't exist
	result, err := client.DynamicBuildParameter(ctx, "nonexistent_payload_12345", "test_param", nil)
	if err != nil {
		t.Logf("Nonexistent payload type error (expected): %v", err)
	}

	if result != nil {
		t.Logf("Response: %s", result.String())
		if result.IsSuccessful() {
			t.Error("Should not succeed for nonexistent payload type")
		}
	}
}

// TestDynamicBuildParameter_NonexistentParameter tests querying a parameter that doesn't exist.
func TestDynamicBuildParameter_NonexistentParameter(t *testing.T) {
	client := getTestClient(t)
	ctx := context.Background()

	// Try to query a parameter that doesn't exist on a valid payload type
	// Note: "apollo" might not be installed in test environment
	result, err := client.DynamicBuildParameter(ctx, "apollo", "nonexistent_param_12345", nil)
	if err != nil {
		t.Logf("Nonexistent parameter error (expected): %v", err)
	}

	if result != nil {
		t.Logf("Response: %s", result.String())
	}
}

// TestDynamicBuildParameter_WithParameters tests build parameter query with context parameters.
func TestDynamicBuildParameter_WithParameters(t *testing.T) {
	client := getTestClient(t)
	ctx := context.Background()

	// Test with additional parameters (might be used by the dynamic query)
	params := map[string]interface{}{
		"operation_id": 1,
		"other_param":  "value",
	}
	result, err := client.DynamicBuildParameter(ctx, "test_payload", "test_param", params)
	if err != nil {
		t.Logf("Query with parameters error (expected if payload doesn't exist): %v", err)
	}

	if result != nil {
		t.Logf("Response: %s", result.String())
		if result.IsSuccessful() && result.HasChoices() {
			t.Logf("Returned %d choices", len(result.Choices))
		}
	}
}

// TestTypedArrayParseFunction_InvalidInput tests input validation for TypedArrayParseFunction.
func TestTypedArrayParseFunction_InvalidInput(t *testing.T) {
	client := getTestClient(t)
	ctx := context.Background()

	// Test with empty input array
	_, err := client.TypedArrayParseFunction(ctx, "", "file_list")
	if err == nil {
		t.Fatal("TypedArrayParseFunction with empty input should return error")
	}
	t.Logf("Empty input error: %v", err)

	// Test with empty parameter type
	_, err = client.TypedArrayParseFunction(ctx, "test:value", "")
	if err == nil {
		t.Fatal("TypedArrayParseFunction with empty parameter type should return error")
	}
	t.Logf("Empty parameter type error: %v", err)
}

// TestTypedArrayParseFunction_InvalidFormat tests parsing with invalid format.
func TestTypedArrayParseFunction_InvalidFormat(t *testing.T) {
	client := getTestClient(t)
	ctx := context.Background()

	// Try to parse with invalid format for the parameter type
	result, err := client.TypedArrayParseFunction(ctx, "this is not valid format !!!", "file_list")
	if err != nil {
		t.Logf("Invalid format error (expected): %v", err)
	}

	if result != nil {
		t.Logf("Response: %s", result.String())
		if result.IsSuccessful() {
			// Might succeed with empty array or partial parse
			t.Logf("Parse completed with %d elements", len(result.ParsedArray))
		}
	}
}

// TestTypedArrayParseFunction_NonexistentParameterType tests parsing with unknown parameter type.
func TestTypedArrayParseFunction_NonexistentParameterType(t *testing.T) {
	client := getTestClient(t)
	ctx := context.Background()

	// Try to parse with a parameter type that doesn't exist
	result, err := client.TypedArrayParseFunction(ctx, "value1,value2,value3", "nonexistent_type_12345")
	if err != nil {
		t.Logf("Nonexistent type error (expected): %v", err)
	}

	if result != nil {
		t.Logf("Response: %s", result.String())
	}
}

// TestTypedArrayParseFunction_SimpleArray tests parsing a simple comma-separated array.
func TestTypedArrayParseFunction_SimpleArray(t *testing.T) {
	client := getTestClient(t)
	ctx := context.Background()

	// Try to parse a simple array
	input := "item1,item2,item3"
	result, err := client.TypedArrayParseFunction(ctx, input, "string_list")
	if err != nil {
		t.Logf("Simple array parse error (expected if type doesn't exist): %v", err)
	}

	if result != nil {
		t.Logf("Response: %s", result.String())
		if result.IsSuccessful() && result.HasElements() {
			t.Logf("Successfully parsed %d elements", len(result.ParsedArray))
			for i, elem := range result.ParsedArray {
				t.Logf("  Element %d: %v", i, elem)
			}
		}
	}
}

// TestTypedArrayParseFunction_ComplexFormat tests parsing complex typed array format.
func TestTypedArrayParseFunction_ComplexFormat(t *testing.T) {
	client := getTestClient(t)
	ctx := context.Background()

	// Try to parse a complex format (file paths with permissions)
	input := "/etc/passwd:read,/etc/shadow:read,/var/log/auth.log:write"
	result, err := client.TypedArrayParseFunction(ctx, input, "file_list")
	if err != nil {
		t.Logf("Complex format parse error (expected if type doesn't exist): %v", err)
	}

	if result != nil {
		t.Logf("Response: %s", result.String())
		if result.IsSuccessful() && result.HasElements() {
			t.Logf("Successfully parsed %d elements", len(result.ParsedArray))
			for i, elem := range result.ParsedArray {
				t.Logf("  Element %d: %v", i, elem)
			}
		}
	}
}

// TestTypedArrayParseFunction_EmptyArray tests parsing an empty array.
func TestTypedArrayParseFunction_EmptyArray(t *testing.T) {
	client := getTestClient(t)
	ctx := context.Background()

	// The empty string check is done in validation, so this should fail
	// But we can test with just whitespace
	input := "   "
	result, err := client.TypedArrayParseFunction(ctx, input, "string_list")
	if err != nil {
		t.Logf("Whitespace-only input error: %v", err)
	}

	if result != nil {
		t.Logf("Response: %s", result.String())
		if result.IsSuccessful() {
			t.Logf("Parsed %d elements from whitespace input", len(result.ParsedArray))
		}
	}
}

// TestTypedArrayParseFunction_JSONFormat tests parsing JSON array format.
func TestTypedArrayParseFunction_JSONFormat(t *testing.T) {
	client := getTestClient(t)
	ctx := context.Background()

	// Try to parse JSON format
	input := `[{"key": "value1"}, {"key": "value2"}, {"key": "value3"}]`
	result, err := client.TypedArrayParseFunction(ctx, input, "json_array")
	if err != nil {
		t.Logf("JSON format parse error (expected if type doesn't exist): %v", err)
	}

	if result != nil {
		t.Logf("Response: %s", result.String())
		if result.IsSuccessful() && result.HasElements() {
			t.Logf("Successfully parsed %d JSON elements", len(result.ParsedArray))
		}
	}
}
