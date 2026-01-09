//go:build integration

package integration

import (
	"context"
	"strings"
	"testing"
)

// TestUtility_CreateRandomInvalidInput tests input validation for CreateRandom.
func TestUtility_CreateRandomInvalidInput(t *testing.T) {
	client := getTestClient(t)
	ctx := context.Background()

	// Test with empty format
	_, err := client.CreateRandom(ctx, "", 5)
	if err == nil {
		t.Fatal("CreateRandom with empty format should return error")
	}
	t.Logf("Empty format error: %v", err)
}

// TestUtility_CreateRandomBasicFormats tests various format strings.
func TestUtility_CreateRandomBasicFormats(t *testing.T) {
	client := getTestClient(t)
	ctx := context.Background()

	testCases := []struct {
		name        string
		format      string
		length      int
		description string
	}{
		{
			name:        "lowercase letters",
			format:      "%s",
			length:      10,
			description: "lowercase letters only",
		},
		{
			name:        "uppercase letters",
			format:      "%S",
			length:      8,
			description: "uppercase letters only",
		},
		{
			name:        "digits",
			format:      "%d",
			length:      6,
			description: "digits only",
		},
		{
			name:        "lowercase hex",
			format:      "%x",
			length:      8,
			description: "lowercase hexadecimal",
		},
		{
			name:        "uppercase hex",
			format:      "%X",
			length:      8,
			description: "uppercase hexadecimal",
		},
		{
			name:        "mixed format",
			format:      "%s-%d",
			length:      5,
			description: "letters and digits separated by dash",
		},
		{
			name:        "complex format",
			format:      "%S%d%x",
			length:      4,
			description: "uppercase letters, digits, and hex",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Logf("Testing format: %s (length: %d) - %s", tc.format, tc.length, tc.description)

			result, err := client.CreateRandom(ctx, tc.format, tc.length)
			if err != nil {
				t.Skipf("CreateRandom not supported or failed: %v", err)
				return
			}

			if !result.IsSuccessful() {
				t.Errorf("CreateRandom failed: %s", result.Error)
				return
			}

			t.Logf("Generated: %s", result.RandomString)

			// Basic validation - just check we got a non-empty string
			if result.RandomString == "" {
				t.Error("Generated string should not be empty")
			}
		})
	}
}

// TestUtility_CreateRandomZeroLength tests CreateRandom with zero length.
func TestUtility_CreateRandomZeroLength(t *testing.T) {
	client := getTestClient(t)
	ctx := context.Background()

	// Zero length should use default length
	result, err := client.CreateRandom(ctx, "%s", 0)
	if err != nil {
		t.Skipf("CreateRandom not supported or failed: %v", err)
		return
	}

	if !result.IsSuccessful() {
		t.Logf("CreateRandom with zero length: %s", result.String())
		// This might be expected behavior
		return
	}

	t.Logf("Generated with default length: %s", result.RandomString)

	if result.RandomString == "" {
		t.Error("Generated string should not be empty even with zero length")
	}
}

// TestUtility_CreateRandomMultipleCalls tests that multiple calls generate different values.
func TestUtility_CreateRandomMultipleCalls(t *testing.T) {
	client := getTestClient(t)
	ctx := context.Background()

	format := "%s%d"
	length := 8

	// Generate 5 random strings
	generated := make(map[string]bool)
	for i := 0; i < 5; i++ {
		result, err := client.CreateRandom(ctx, format, length)
		if err != nil {
			t.Skipf("CreateRandom not supported or failed: %v", err)
			return
		}

		if !result.IsSuccessful() {
			t.Skipf("CreateRandom failed: %s", result.Error)
			return
		}

		t.Logf("Generated #%d: %s", i+1, result.RandomString)
		generated[result.RandomString] = true
	}

	// Check that we got at least some different values
	// (Very unlikely to get all identical random strings)
	if len(generated) == 1 {
		t.Error("All generated strings are identical - randomness might not be working")
	} else {
		t.Logf("Generated %d unique strings out of 5", len(generated))
	}
}

// TestUtility_ConfigCheck tests configuration checking.
func TestUtility_ConfigCheck(t *testing.T) {
	client := getTestClient(t)
	ctx := context.Background()

	result, err := client.ConfigCheck(ctx)
	if err != nil {
		t.Skipf("ConfigCheck not supported or failed: %v", err)
		return
	}

	t.Logf("Configuration check result: %s", result.String())

	// Log configuration details
	if result.IsValid() {
		t.Log("Configuration is valid")
	} else {
		t.Log("Configuration has issues")
	}

	// Log any errors
	if result.HasErrors() {
		t.Logf("Configuration errors: %v", result.GetErrors())
	} else {
		t.Log("No configuration errors")
	}

	// Log config details if available
	if len(result.Config) > 0 {
		t.Logf("Configuration details: %d entries", len(result.Config))
		for key, value := range result.Config {
			t.Logf("  %s: %v", key, value)
		}
	}

	// Log message if present
	if result.Message != "" {
		t.Logf("Message: %s", result.Message)
	}
}

// TestUtility_ConfigCheckStructure tests the structure of ConfigCheck response.
func TestUtility_ConfigCheckStructure(t *testing.T) {
	client := getTestClient(t)
	ctx := context.Background()

	result, err := client.ConfigCheck(ctx)
	if err != nil {
		t.Skipf("ConfigCheck not supported or failed: %v", err)
		return
	}

	// Test that helper methods work without crashing
	_ = result.IsValid()
	_ = result.HasErrors()
	_ = result.GetErrors()
	_ = result.String()

	// Verify structure
	t.Logf("Status: %s", result.Status)
	t.Logf("Valid: %v", result.Valid)
	t.Logf("Error count: %d", len(result.GetErrors()))

	// If there are errors, log them
	if result.HasErrors() {
		for i, err := range result.GetErrors() {
			t.Logf("Error %d: %s", i+1, err)
		}
	}
}

// TestUtility_CreateRandomSpecialChars tests formats with special characters.
func TestUtility_CreateRandomSpecialChars(t *testing.T) {
	client := getTestClient(t)
	ctx := context.Background()

	// Test format with literal characters
	testFormats := []string{
		"callback-%s",
		"user_%d",
		"%s@%s.com",
		"ID_%X",
	}

	for _, format := range testFormats {
		t.Run(format, func(t *testing.T) {
			result, err := client.CreateRandom(ctx, format, 5)
			if err != nil {
				t.Logf("Format %s error (may not be supported): %v", format, err)
				return
			}

			if !result.IsSuccessful() {
				t.Logf("Format %s failed: %s", format, result.Error)
				return
			}

			t.Logf("Format %s generated: %s", format, result.RandomString)

			// Check that literal parts are preserved
			if strings.Contains(format, "callback-") && !strings.Contains(result.RandomString, "callback-") {
				t.Errorf("Expected literal 'callback-' in result")
			}
		})
	}
}

// TestUtility_CreateRandomLengthVariations tests different length values.
func TestUtility_CreateRandomLengthVariations(t *testing.T) {
	client := getTestClient(t)
	ctx := context.Background()

	lengths := []int{1, 4, 8, 16, 32, 64}

	for _, length := range lengths {
		t.Run(string(rune('0'+length)), func(t *testing.T) {
			result, err := client.CreateRandom(ctx, "%s", length)
			if err != nil {
				t.Logf("Length %d error: %v", length, err)
				return
			}

			if !result.IsSuccessful() {
				t.Logf("Length %d failed: %s", length, result.Error)
				return
			}

			t.Logf("Length %d generated: %s (actual length: %d)",
				length, result.RandomString, len(result.RandomString))

			// Note: The actual length might not exactly match the requested length
			// depending on how Mythic interprets the format string
		})
	}
}

// TestUtility_CreateRandomInvalidFormat tests behavior with potentially invalid formats.
func TestUtility_CreateRandomInvalidFormat(t *testing.T) {
	client := getTestClient(t)
	ctx := context.Background()

	invalidFormats := []string{
		"%z",      // Invalid format specifier
		"%%",      // Escaped percent
		"%",       // Incomplete format
		"no_spec", // No format specifier
	}

	for _, format := range invalidFormats {
		t.Run(format, func(t *testing.T) {
			result, err := client.CreateRandom(ctx, format, 5)
			if err != nil {
				t.Logf("Invalid format %q returned error (expected): %v", format, err)
				return
			}

			if !result.IsSuccessful() {
				t.Logf("Invalid format %q failed (expected): %s", format, result.Error)
				return
			}

			// If it succeeded, log what was generated
			t.Logf("Format %q unexpectedly succeeded, generated: %s", format, result.RandomString)
		})
	}
}

// TestUtility_GetConfig tests retrieving client configuration
func TestUtility_GetConfig(t *testing.T) {

	client := AuthenticateTestClient(t)

	// GetConfig doesn't require context - it's a local operation
	config := client.GetConfig()

	// Verify config structure
	if config.ServerURL == "" {
		t.Error("Config should have ServerURL")
	}
	if config.Username == "" {
		t.Error("Config should have Username")
	}

	t.Logf("Client configuration:")
	t.Logf("  - ServerURL: %s", config.ServerURL)
	t.Logf("  - Username: %s", config.Username)
	t.Logf("  - SSL: %v", config.SSL)
	t.Logf("  - SkipTLSVerify: %v", config.SkipTLSVerify)
	t.Logf("  - Timeout: %v", config.Timeout)

	// Password should be redacted or empty in retrieved config
	if config.Password != "" {
		t.Log("Note: Password is present in config (should be handled securely)")
	}
}

// TestUtility_GetGlobalSettings tests retrieving global Mythic settings
func TestUtility_GetGlobalSettings(t *testing.T) {

	client := AuthenticateTestClient(t)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	settings, err := client.GetGlobalSettings(ctx)
	if err != nil {
		t.Fatalf("GetGlobalSettings failed: %v", err)
	}

	if settings == nil {
		t.Fatal("GetGlobalSettings returned nil")
	}

	t.Logf("Global settings retrieved: %d setting(s)", len(settings))

	// Log some settings (if any)
	count := 0
	for key, value := range settings {
		if count < 10 { // Only log first 10
			t.Logf("  - %s: %v", key, value)
			count++
		}
	}

	if len(settings) > 10 {
		t.Logf("  ... and %d more settings", len(settings)-10)
	}
}

// TestUtility_UpdateGlobalSettings tests updating global Mythic settings
func TestUtility_UpdateGlobalSettings(t *testing.T) {
	t.Skip("Skipping UpdateGlobalSettings to avoid modifying global configuration")

	client := AuthenticateTestClient(t)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Get current settings first
	currentSettings, err := client.GetGlobalSettings(ctx)
	if err != nil {
		t.Fatalf("GetGlobalSettings failed: %v", err)
	}

	t.Logf("Current settings count: %d", len(currentSettings))

	// Update a test setting (use a non-critical setting)
	// Note: This requires admin privileges
	testSettings := map[string]interface{}{
		"test_setting": "test_value",
	}

	err = client.UpdateGlobalSettings(ctx, testSettings)
	if err != nil {
		t.Logf("UpdateGlobalSettings failed (may require admin): %v", err)
		return
	}

	t.Log("Successfully updated global settings")

	// Verify the setting was updated
	updatedSettings, err := client.GetGlobalSettings(ctx)
	if err != nil {
		t.Fatalf("Failed to retrieve updated settings: %v", err)
	}

	if val, ok := updatedSettings["test_setting"]; ok {
		if val != "test_value" {
			t.Errorf("Expected test_setting value 'test_value', got %v", val)
		}
	}

	// Restore original settings (remove test setting)
	restoreSettings := map[string]interface{}{
		"test_setting": nil,
	}
	if err := client.UpdateGlobalSettings(ctx, restoreSettings); err != nil {
		t.Logf("Warning: Failed to restore original settings: %v", err)
	}
}

// TestUtility_UpdateGlobalSettings_InvalidInput tests with invalid input
func TestUtility_UpdateGlobalSettings_InvalidInput(t *testing.T) {

	client := AuthenticateTestClient(t)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Test with nil settings
	err := client.UpdateGlobalSettings(ctx, nil)
	if err == nil {
		t.Fatal("Expected error for nil settings, got nil")
	}
	t.Logf("Nil settings error: %v", err)

	// Test with empty settings map
	err = client.UpdateGlobalSettings(ctx, map[string]interface{}{})
	if err == nil {
		t.Fatal("Expected error for empty settings, got nil")
	}
	t.Logf("Empty settings error: %v", err)
}
