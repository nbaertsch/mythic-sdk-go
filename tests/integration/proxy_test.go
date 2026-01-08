//go:build integration

package integration

import (
	"context"
	"testing"
)

// TestProxy_ToggleInvalidInput tests input validation for ToggleProxy.
func TestProxy_ToggleInvalidInput(t *testing.T) {
	client := getTestClient(t)
	ctx := context.Background()

	// Test with zero task ID
	_, err := client.ToggleProxy(ctx, 0, 1080, true)
	if err == nil {
		t.Fatal("ToggleProxy with zero task ID should return error")
	}
	t.Logf("Zero task ID error: %v", err)

	// Test with negative task ID
	_, err = client.ToggleProxy(ctx, -1, 1080, true)
	if err == nil {
		t.Fatal("ToggleProxy with negative task ID should return error")
	}
	t.Logf("Negative task ID error: %v", err)

	// Test with zero port
	_, err = client.ToggleProxy(ctx, 123, 0, true)
	if err == nil {
		t.Fatal("ToggleProxy with zero port should return error")
	}
	t.Logf("Zero port error: %v", err)

	// Test with negative port
	_, err = client.ToggleProxy(ctx, 123, -1, true)
	if err == nil {
		t.Fatal("ToggleProxy with negative port should return error")
	}
	t.Logf("Negative port error: %v", err)

	// Test with port too high
	_, err = client.ToggleProxy(ctx, 123, 65536, true)
	if err == nil {
		t.Fatal("ToggleProxy with port > 65535 should return error")
	}
	t.Logf("Port too high error: %v", err)
}

// TestProxy_TestInvalidInput tests input validation for TestProxy.
func TestProxy_TestInvalidInput(t *testing.T) {
	client := getTestClient(t)
	ctx := context.Background()

	// Test with zero callback ID
	_, err := client.TestProxy(ctx, 0, 1080, "https://www.google.com")
	if err == nil {
		t.Fatal("TestProxy with zero callback ID should return error")
	}
	t.Logf("Zero callback ID error: %v", err)

	// Test with negative callback ID
	_, err = client.TestProxy(ctx, -1, 1080, "https://www.google.com")
	if err == nil {
		t.Fatal("TestProxy with negative callback ID should return error")
	}
	t.Logf("Negative callback ID error: %v", err)

	// Test with zero port
	_, err = client.TestProxy(ctx, 123, 0, "https://www.google.com")
	if err == nil {
		t.Fatal("TestProxy with zero port should return error")
	}
	t.Logf("Zero port error: %v", err)

	// Test with port too high
	_, err = client.TestProxy(ctx, 123, 99999, "https://www.google.com")
	if err == nil {
		t.Fatal("TestProxy with port > 65535 should return error")
	}
	t.Logf("Port too high error: %v", err)

	// Test with empty target URL
	_, err = client.TestProxy(ctx, 123, 1080, "")
	if err == nil {
		t.Fatal("TestProxy with empty target URL should return error")
	}
	t.Logf("Empty target URL error: %v", err)
}

// TestProxy_ToggleNonexistentTask tests ToggleProxy with task that doesn't exist.
func TestProxy_ToggleNonexistentTask(t *testing.T) {
	client := getTestClient(t)
	ctx := context.Background()

	// Try to toggle proxy with a task ID that likely doesn't exist
	_, err := client.ToggleProxy(ctx, 999999, 1080, true)
	if err == nil {
		t.Fatal("ToggleProxy with nonexistent task should return error")
	}
	t.Logf("Nonexistent task error (expected): %v", err)
}

// TestProxy_TestNonexistentCallback tests TestProxy with callback that doesn't exist.
func TestProxy_TestNonexistentCallback(t *testing.T) {
	client := getTestClient(t)
	ctx := context.Background()

	// Try to test proxy with a callback ID that likely doesn't exist
	result, err := client.TestProxy(ctx, 999999, 1080, "https://www.google.com")
	if err == nil && !result.IsSuccessful() {
		t.Logf("Test returned unsuccessful result (expected): %s", result.String())
		return
	}
	if err != nil {
		t.Logf("Test returned error (expected): %v", err)
		return
	}

	// If we got a successful result, something is wrong
	if result.IsSuccessful() {
		t.Error("TestProxy with nonexistent callback should not succeed")
	}
}

// TestProxy_TogglePortRange tests various valid port numbers.
func TestProxy_TogglePortRange(t *testing.T) {
	client := getTestClient(t)
	ctx := context.Background()

	// Test boundary ports (these should pass validation even if task doesn't exist)
	validPorts := []int{1, 1080, 8080, 9050, 65535}

	for _, port := range validPorts {
		_, err := client.ToggleProxy(ctx, 999999, port, true)
		// We expect an error because the task doesn't exist, not because of validation
		if err != nil {
			errStr := err.Error()
			// Check that it's not a validation error
			if contains(errStr, "must be between") || contains(errStr, "must be positive") {
				t.Errorf("Port %d failed validation but should be valid", port)
			} else {
				t.Logf("Port %d passed validation (task error expected): %v", port, err)
			}
		}
	}
}

// TestProxy_TestURLVariations tests various URL formats.
func TestProxy_TestURLVariations(t *testing.T) {
	client := getTestClient(t)
	ctx := context.Background()

	// Test various URL formats (should pass validation even if callback doesn't exist)
	urls := []string{
		"https://www.google.com",
		"http://example.com",
		"https://10.0.0.1",
		"http://192.168.1.1:8080",
		"https://internal.company.com/test",
	}

	for _, url := range urls {
		_, err := client.TestProxy(ctx, 999999, 1080, url)
		// We expect an error because the callback doesn't exist, not because of validation
		if err != nil {
			errStr := err.Error()
			// Check that it's not a validation error
			if contains(errStr, "cannot be empty") {
				t.Errorf("URL %s failed validation but should be valid", url)
			} else {
				t.Logf("URL %s passed validation (callback error expected): %v", url, err)
			}
		}
	}
}

// TestProxy_ToggleBothDirections tests enable and disable operations.
func TestProxy_ToggleBothDirections(t *testing.T) {
	client := getTestClient(t)
	ctx := context.Background()

	// Test that both enable=true and enable=false pass validation
	testTaskID := 999999
	testPort := 1080

	t.Log("Testing enable=true")
	_, err := client.ToggleProxy(ctx, testTaskID, testPort, true)
	if err != nil {
		t.Logf("Enable proxy error (task doesn't exist, expected): %v", err)
	}

	t.Log("Testing enable=false")
	_, err = client.ToggleProxy(ctx, testTaskID, testPort, false)
	if err != nil {
		t.Logf("Disable proxy error (task doesn't exist, expected): %v", err)
	}
}

// TestProxy_ConcurrentOperations tests that multiple proxy operations can be called.
func TestProxy_ConcurrentOperations(t *testing.T) {
	client := getTestClient(t)
	ctx := context.Background()

	// Test that we can make multiple proxy API calls without issues
	operations := []struct {
		name string
		fn   func() error
	}{
		{
			name: "toggle enable",
			fn: func() error {
				_, err := client.ToggleProxy(ctx, 999999, 1080, true)
				return err
			},
		},
		{
			name: "toggle disable",
			fn: func() error {
				_, err := client.ToggleProxy(ctx, 999999, 1080, false)
				return err
			},
		},
		{
			name: "test proxy",
			fn: func() error {
				_, err := client.TestProxy(ctx, 999999, 1080, "https://www.google.com")
				return err
			},
		},
	}

	for _, op := range operations {
		t.Run(op.name, func(t *testing.T) {
			err := op.fn()
			if err != nil {
				t.Logf("%s returned error (expected): %v", op.name, err)
			}
		})
	}
}

// TestProxy_ValidPortBoundaries tests port boundary conditions.
func TestProxy_ValidPortBoundaries(t *testing.T) {
	client := getTestClient(t)
	ctx := context.Background()

	testCases := []struct {
		port        int
		shouldPass  bool
		description string
	}{
		{0, false, "zero"},
		{1, true, "minimum valid"},
		{1023, true, "below privileged"},
		{1024, true, "first non-privileged"},
		{8080, true, "common proxy port"},
		{65535, true, "maximum valid"},
		{65536, false, "above maximum"},
		{99999, false, "way above maximum"},
		{-1, false, "negative"},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			_, err := client.ToggleProxy(ctx, 999999, tc.port, true)

			if tc.shouldPass {
				// Should NOT get validation error
				if err != nil {
					errStr := err.Error()
					if contains(errStr, "must be between") || contains(errStr, "must be positive") {
						t.Errorf("Port %d should pass validation but got: %v", tc.port, err)
					}
				}
			} else {
				// Should get validation error
				if err == nil {
					t.Errorf("Port %d should fail validation but passed", tc.port)
				} else {
					t.Logf("Port %d correctly failed validation: %v", tc.port, err)
				}
			}
		})
	}
}
