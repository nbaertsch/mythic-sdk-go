package integration

import (
	"context"
	"testing"
)

// TestBuildParameters_GetBuildParameters tests retrieving all build parameter definitions.
func TestBuildParameters_GetBuildParameters(t *testing.T) {
	client := getTestClient(t)
	ctx := context.Background()

	parameters, err := client.GetBuildParameters(ctx)
	if err != nil {
		t.Fatalf("GetBuildParameters() failed: %v", err)
	}

	t.Logf("Found %d build parameter definitions", len(parameters))

	if len(parameters) == 0 {
		t.Log("No build parameters found - this is expected if no payload types are installed")
		return
	}

	// Verify parameter structure
	for i, param := range parameters {
		if i < 5 { // Log first 5 for debugging
			t.Logf("Parameter %d: %s", i+1, param.String())
		}

		if param.ID == 0 {
			t.Errorf("Parameter %d has zero ID", i)
		}
		if param.Name == "" {
			t.Errorf("Parameter %d has empty Name", i)
		}
		if param.PayloadTypeID == 0 {
			t.Errorf("Parameter %s has zero PayloadTypeID", param.Name)
		}
		if param.ParameterType == "" {
			t.Errorf("Parameter %s has empty ParameterType", param.Name)
		}

		// Test helper methods
		_ = param.IsRequired()
		_ = param.IsCrypto()
		_ = param.ShouldRandomize()
		_ = param.IsDeleted()
		_ = param.String()

		// Verify not deleted (we filter deleted parameters)
		if param.Deleted {
			t.Errorf("Parameter %s should not be deleted (we filter deleted parameters)", param.Name)
		}
	}

	// Verify parameters are sorted by payload type ID, then name
	for i := 1; i < len(parameters); i++ {
		if parameters[i-1].PayloadTypeID > parameters[i].PayloadTypeID {
			t.Errorf("Parameters not sorted by PayloadTypeID: %d comes after %d",
				parameters[i-1].PayloadTypeID, parameters[i].PayloadTypeID)
			break
		}
		if parameters[i-1].PayloadTypeID == parameters[i].PayloadTypeID {
			if parameters[i-1].Name > parameters[i].Name {
				t.Errorf("Parameters not sorted by name within PayloadTypeID %d: %s comes after %s",
					parameters[i].PayloadTypeID, parameters[i-1].Name, parameters[i].Name)
				break
			}
		}
	}
}

// TestBuildParameters_GetBuildParametersStructure tests parameter field values.
func TestBuildParameters_GetBuildParametersStructure(t *testing.T) {
	client := getTestClient(t)
	ctx := context.Background()

	parameters, err := client.GetBuildParameters(ctx)
	if err != nil {
		t.Fatalf("GetBuildParameters() failed: %v", err)
	}

	if len(parameters) == 0 {
		t.Skip("No build parameters available for structure testing")
	}

	// Find required and optional parameters
	var requiredParam, optionalParam *types.BuildParameterType
	for _, param := range parameters {
		if param.Required && requiredParam == nil {
			requiredParam = param
		}
		if !param.Required && optionalParam == nil {
			optionalParam = param
		}
		if requiredParam != nil && optionalParam != nil {
			break
		}
	}

	if requiredParam != nil {
		t.Logf("Testing required parameter: %s", requiredParam.String())
		if !requiredParam.IsRequired() {
			t.Error("Required parameter IsRequired() returned false")
		}
		str := requiredParam.String()
		if !contains(str, "required") {
			t.Errorf("Required parameter String() should contain 'required', got %q", str)
		}
	}

	if optionalParam != nil {
		t.Logf("Testing optional parameter: %s", optionalParam.String())
		if optionalParam.IsRequired() {
			t.Error("Optional parameter IsRequired() returned true")
		}
	}
}

// TestBuildParameters_GetBuildParametersByPayloadTypeInvalid tests error handling.
func TestBuildParameters_GetBuildParametersByPayloadTypeInvalid(t *testing.T) {
	client := getTestClient(t)
	ctx := context.Background()

	// Test with zero payload type ID
	_, err := client.GetBuildParametersByPayloadType(ctx, 0)
	if err == nil {
		t.Fatal("GetBuildParametersByPayloadType(0) should return an error")
	}
	t.Logf("Zero payload type ID error: %v", err)
}

// TestBuildParameters_GetBuildParametersByPayloadType tests filtering by payload type.
func TestBuildParameters_GetBuildParametersByPayloadType(t *testing.T) {
	client := getTestClient(t)
	ctx := context.Background()

	// First get all parameters to find available payload types
	allParameters, err := client.GetBuildParameters(ctx)
	if err != nil {
		t.Fatalf("GetBuildParameters() failed: %v", err)
	}

	if len(allParameters) == 0 {
		t.Skip("No build parameters available for payload type filtering test")
	}

	// Get unique payload type IDs
	payloadTypeMap := make(map[int]bool)
	for _, param := range allParameters {
		payloadTypeMap[param.PayloadTypeID] = true
	}

	payloadTypeIDs := make([]int, 0, len(payloadTypeMap))
	for id := range payloadTypeMap {
		payloadTypeIDs = append(payloadTypeIDs, id)
	}

	t.Logf("Found %d unique payload types", len(payloadTypeIDs))

	// Test with first payload type
	if len(payloadTypeIDs) > 0 {
		testPayloadTypeID := payloadTypeIDs[0]
		t.Logf("Testing with payload type ID: %d", testPayloadTypeID)

		parameters, err := client.GetBuildParametersByPayloadType(ctx, testPayloadTypeID)
		if err != nil {
			t.Fatalf("GetBuildParametersByPayloadType(%d) failed: %v", testPayloadTypeID, err)
		}

		t.Logf("Found %d parameters for payload type %d", len(parameters), testPayloadTypeID)

		// Verify all parameters belong to the requested payload type
		for _, param := range parameters {
			if param.PayloadTypeID != testPayloadTypeID {
				t.Errorf("Parameter %s has PayloadTypeID %d, expected %d",
					param.Name, param.PayloadTypeID, testPayloadTypeID)
			}
		}

		// Verify sorting by name
		for i := 1; i < len(parameters); i++ {
			if parameters[i-1].Name > parameters[i].Name {
				t.Errorf("Parameters not sorted by name: %s comes after %s",
					parameters[i-1].Name, parameters[i].Name)
				break
			}
		}
	}
}

// TestBuildParameters_GetBuildParametersByPayloadTypeNonexistent tests behavior with nonexistent payload type.
func TestBuildParameters_GetBuildParametersByPayloadTypeNonexistent(t *testing.T) {
	client := getTestClient(t)
	ctx := context.Background()

	// Use a payload type ID that likely doesn't exist
	parameters, err := client.GetBuildParametersByPayloadType(ctx, 999999)
	if err != nil {
		t.Fatalf("GetBuildParametersByPayloadType(999999) failed: %v", err)
	}

	// Should return empty list, not an error
	if len(parameters) != 0 {
		t.Errorf("Expected empty list for nonexistent payload type, got %d parameters", len(parameters))
	}
}

// TestBuildParameters_GetBuildParameterInstances tests retrieving parameter instances.
func TestBuildParameters_GetBuildParameterInstances(t *testing.T) {
	client := getTestClient(t)
	ctx := context.Background()

	instances, err := client.GetBuildParameterInstances(ctx)
	if err != nil {
		t.Fatalf("GetBuildParameterInstances() failed: %v", err)
	}

	t.Logf("Found %d build parameter instances", len(instances))

	if len(instances) == 0 {
		t.Log("No build parameter instances found - this is expected if no payloads have been created")
		return
	}

	// Verify instance structure
	for i, inst := range instances {
		if i < 5 { // Log first 5 for debugging
			t.Logf("Instance %d: %s", i+1, inst.String())
		}

		if inst.ID == 0 {
			t.Errorf("Instance %d has zero ID", i)
		}
		if inst.PayloadID == 0 {
			t.Errorf("Instance %d has zero PayloadID", i)
		}
		if inst.BuildParameterID == 0 {
			t.Errorf("Instance %d has zero BuildParameterID", i)
		}
		if inst.Value == "" && !inst.IsEncrypted() {
			t.Errorf("Instance %d has empty value and is not encrypted", i)
		}

		// Test helper methods
		_ = inst.IsEncrypted()
		_ = inst.GetValue()
		_ = inst.String()
	}

	// Verify instances are sorted by payload ID, then build parameter ID
	for i := 1; i < len(instances); i++ {
		if instances[i-1].PayloadID > instances[i].PayloadID {
			t.Errorf("Instances not sorted by PayloadID: %d comes after %d",
				instances[i-1].PayloadID, instances[i].PayloadID)
			break
		}
		if instances[i-1].PayloadID == instances[i].PayloadID {
			if instances[i-1].BuildParameterID > instances[i].BuildParameterID {
				t.Errorf("Instances not sorted by BuildParameterID within PayloadID %d: %d comes after %d",
					instances[i].PayloadID, instances[i-1].BuildParameterID, instances[i].BuildParameterID)
				break
			}
		}
	}
}

// TestBuildParameters_GetBuildParameterInstancesByPayloadInvalid tests error handling.
func TestBuildParameters_GetBuildParameterInstancesByPayloadInvalid(t *testing.T) {
	client := getTestClient(t)
	ctx := context.Background()

	// Test with zero payload ID
	_, err := client.GetBuildParameterInstancesByPayload(ctx, 0)
	if err == nil {
		t.Fatal("GetBuildParameterInstancesByPayload(0) should return an error")
	}
	t.Logf("Zero payload ID error: %v", err)
}

// TestBuildParameters_GetBuildParameterInstancesByPayload tests filtering by payload.
func TestBuildParameters_GetBuildParameterInstancesByPayload(t *testing.T) {
	client := getTestClient(t)
	ctx := context.Background()

	// First get payloads to find one with build parameter instances
	payloads, err := client.GetPayloads(ctx)
	if err != nil {
		t.Fatalf("GetPayloads() failed: %v", err)
	}

	if len(payloads) == 0 {
		t.Skip("No payloads available for build parameter instance filtering test")
	}

	// Test with first payload
	testPayloadID := payloads[0].ID
	t.Logf("Testing with payload ID: %d", testPayloadID)

	instances, err := client.GetBuildParameterInstancesByPayload(ctx, testPayloadID)
	if err != nil {
		t.Fatalf("GetBuildParameterInstancesByPayload(%d) failed: %v", testPayloadID, err)
	}

	t.Logf("Found %d parameter instances for payload %d", len(instances), testPayloadID)

	// Verify all instances belong to the requested payload
	for _, inst := range instances {
		if inst.PayloadID != testPayloadID {
			t.Errorf("Instance has PayloadID %d, expected %d", inst.PayloadID, testPayloadID)
		}

		// Verify nested BuildParameter is populated
		if inst.BuildParameter == nil {
			t.Error("Instance BuildParameter should not be nil")
		} else {
			if inst.BuildParameter.Name == "" {
				t.Error("Instance BuildParameter.Name should not be empty")
			}
			if inst.BuildParameter.ParameterType == "" {
				t.Error("Instance BuildParameter.ParameterType should not be empty")
			}
		}
	}

	// Verify sorting by build parameter ID
	for i := 1; i < len(instances); i++ {
		if instances[i-1].BuildParameterID > instances[i].BuildParameterID {
			t.Errorf("Instances not sorted by BuildParameterID: %d comes after %d",
				instances[i-1].BuildParameterID, instances[i].BuildParameterID)
			break
		}
	}
}

// TestBuildParameters_GetBuildParameterInstancesByPayloadNonexistent tests behavior with nonexistent payload.
func TestBuildParameters_GetBuildParameterInstancesByPayloadNonexistent(t *testing.T) {
	client := getTestClient(t)
	ctx := context.Background()

	// Use a payload ID that likely doesn't exist
	instances, err := client.GetBuildParameterInstancesByPayload(ctx, 999999)
	if err != nil {
		t.Fatalf("GetBuildParameterInstancesByPayload(999999) failed: %v", err)
	}

	// Should return empty list, not an error
	if len(instances) != 0 {
		t.Errorf("Expected empty list for nonexistent payload, got %d instances", len(instances))
	}
}

// TestBuildParameters_ParameterTypes tests different parameter types.
func TestBuildParameters_ParameterTypes(t *testing.T) {
	client := getTestClient(t)
	ctx := context.Background()

	parameters, err := client.GetBuildParameters(ctx)
	if err != nil {
		t.Fatalf("GetBuildParameters() failed: %v", err)
	}

	if len(parameters) == 0 {
		t.Skip("No build parameters available")
	}

	typeCount := make(map[string]int)
	for _, param := range parameters {
		typeCount[param.ParameterType]++
	}

	t.Logf("Parameter type distribution:")
	for typ, count := range typeCount {
		t.Logf("  %s: %d", typ, count)
	}
}

// TestBuildParameters_RequiredVsOptional tests required vs optional parameters.
func TestBuildParameters_RequiredVsOptional(t *testing.T) {
	client := getTestClient(t)
	ctx := context.Background()

	parameters, err := client.GetBuildParameters(ctx)
	if err != nil {
		t.Fatalf("GetBuildParameters() failed: %v", err)
	}

	if len(parameters) == 0 {
		t.Skip("No build parameters available")
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

// TestBuildParameters_EncryptedInstances tests encrypted parameter instances.
func TestBuildParameters_EncryptedInstances(t *testing.T) {
	client := getTestClient(t)
	ctx := context.Background()

	instances, err := client.GetBuildParameterInstances(ctx)
	if err != nil {
		t.Fatalf("GetBuildParameterInstances() failed: %v", err)
	}

	if len(instances) == 0 {
		t.Skip("No build parameter instances available")
	}

	encryptedCount := 0
	plainCount := 0

	for _, inst := range instances {
		if inst.IsEncrypted() {
			encryptedCount++
		} else {
			plainCount++
		}

		// GetValue() should always return a value
		value := inst.GetValue()
		if value == "" {
			t.Errorf("Instance %d GetValue() returned empty string", inst.ID)
		}
	}

	t.Logf("Found %d encrypted instances and %d plain instances", encryptedCount, plainCount)
}

// TestBuildParameters_PayloadTypeGrouping tests parameter grouping by payload type.
func TestBuildParameters_PayloadTypeGrouping(t *testing.T) {
	client := getTestClient(t)
	ctx := context.Background()

	parameters, err := client.GetBuildParameters(ctx)
	if err != nil {
		t.Fatalf("GetBuildParameters() failed: %v", err)
	}

	if len(parameters) == 0 {
		t.Skip("No build parameters available")
	}

	// Group parameters by payload type
	payloadTypeParams := make(map[int][]*types.BuildParameterType)
	for _, param := range parameters {
		payloadTypeParams[param.PayloadTypeID] = append(payloadTypeParams[param.PayloadTypeID], param)
	}

	t.Logf("Build parameters by payload type:")
	for payloadTypeID, params := range payloadTypeParams {
		t.Logf("  Payload Type %d: %d parameters", payloadTypeID, len(params))

		// Verify GetBuildParametersByPayloadType returns the same count
		typeParams, err := client.GetBuildParametersByPayloadType(ctx, payloadTypeID)
		if err != nil {
			t.Errorf("GetBuildParametersByPayloadType(%d) failed: %v", payloadTypeID, err)
			continue
		}

		if len(typeParams) != len(params) {
			t.Errorf("GetBuildParametersByPayloadType(%d) returned %d parameters, expected %d",
				payloadTypeID, len(typeParams), len(params))
		}
	}
}
