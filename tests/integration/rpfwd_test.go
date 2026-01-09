//go:build integration

package integration

import (
	"context"
	"testing"
	"time"

	"github.com/nbaertsch/mythic-sdk-go/pkg/mythic/types"
)

func TestRPFWD_GetRPFWDs(t *testing.T) {
	SkipIfNoMythic(t)
	client := AuthenticateTestClient(t)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Get active callbacks
	callbacks, err := client.GetAllActiveCallbacks(ctx)
	if err != nil {
		t.Fatalf("GetAllActiveCallbacks failed: %v", err)
	}
	if len(callbacks) == 0 {
		t.Skip("No active callbacks for testing")
	}

	callbackID := callbacks[0].ID
	rpfwds, err := client.GetRPFWDs(ctx, callbackID)
	if err != nil {
		t.Fatalf("GetRPFWDs failed: %v", err)
	}

	if rpfwds == nil {
		t.Fatal("GetRPFWDs returned nil")
	}

	t.Logf("Found %d RPFWD tunnel(s) for callback %d", len(rpfwds), callbackID)

	for _, rpfwd := range rpfwds {
		if rpfwd.ID == 0 {
			t.Error("RPFWD ID should not be 0")
		}
		if rpfwd.CallbackID != callbackID {
			t.Errorf("Expected callback ID %d, got %d", callbackID, rpfwd.CallbackID)
		}
		t.Logf("  - %s", rpfwd.String())
		t.Logf("    Local: %s → Remote: %s",
			rpfwd.GetLocalEndpoint(), rpfwd.GetRemoteEndpoint())
	}
}

func TestRPFWD_CreateAndDeleteRPFWD(t *testing.T) {
	t.Skip("Skipping create/delete test to avoid network changes")

	SkipIfNoMythic(t)
	client := AuthenticateTestClient(t)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Get active callbacks
	callbacks, err := client.GetAllActiveCallbacks(ctx)
	if err != nil {
		t.Fatalf("GetAllActiveCallbacks failed: %v", err)
	}
	if len(callbacks) == 0 {
		t.Skip("No active callbacks for testing")
	}

	// Create RPFWD tunnel
	req := &types.CreateRPFWDRequest{
		CallbackID: callbacks[0].ID,
		LocalPort:  19999,
		RemoteHost: "127.0.0.1",
		RemotePort: 80,
	}

	rpfwd, err := client.CreateRPFWD(ctx, req)
	if err != nil {
		t.Fatalf("CreateRPFWD failed: %v", err)
	}

	if rpfwd == nil {
		t.Fatal("CreateRPFWD returned nil")
	}

	t.Logf("Created RPFWD tunnel %d: %s", rpfwd.ID, rpfwd.String())
	t.Logf("  - Connect to: %s", rpfwd.GetLocalEndpoint())
	t.Logf("  - Forwards to: %s", rpfwd.GetRemoteEndpoint())

	// Verify it was created
	if !rpfwd.IsActive() {
		t.Error("Created RPFWD should be active")
	}

	// Delete the tunnel
	err = client.DeleteRPFWD(ctx, rpfwd.ID)
	if err != nil {
		t.Fatalf("DeleteRPFWD failed: %v", err)
	}

	t.Logf("Successfully deleted RPFWD tunnel %d", rpfwd.ID)

	// Verify it was deleted (marked as inactive)
	status, err := client.GetRPFWDStatus(ctx, rpfwd.ID)
	if err != nil {
		t.Fatalf("GetRPFWDStatus failed: %v", err)
	}

	if status.Active {
		t.Error("Deleted RPFWD should be inactive")
	}
}

func TestRPFWD_GetRPFWDStatus(t *testing.T) {
	SkipIfNoMythic(t)
	client := AuthenticateTestClient(t)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Get callbacks and their RPFWDs
	callbacks, err := client.GetAllActiveCallbacks(ctx)
	if err != nil {
		t.Fatalf("GetAllActiveCallbacks failed: %v", err)
	}
	if len(callbacks) == 0 {
		t.Skip("No active callbacks for testing")
	}

	rpfwds, err := client.GetRPFWDs(ctx, callbacks[0].ID)
	if err != nil {
		t.Fatalf("GetRPFWDs failed: %v", err)
	}
	if len(rpfwds) == 0 {
		t.Skip("No RPFWD tunnels available for testing")
	}

	rpfwdID := rpfwds[0].ID
	status, err := client.GetRPFWDStatus(ctx, rpfwdID)
	if err != nil {
		t.Fatalf("GetRPFWDStatus failed: %v", err)
	}

	if status == nil {
		t.Fatal("GetRPFWDStatus returned nil")
	}

	if status.ID != rpfwdID {
		t.Errorf("Expected RPFWD ID %d, got %d", rpfwdID, status.ID)
	}

	t.Logf("RPFWD %d status:", rpfwdID)
	t.Logf("  - %s", status.String())
	t.Logf("  - Active: %v", status.Active)
	t.Logf("  - Local endpoint: %s", status.GetLocalEndpoint())
	t.Logf("  - Remote endpoint: %s", status.GetRemoteEndpoint())
	t.Logf("  - Timestamp: %s", status.Timestamp.Format("2006-01-02 15:04:05"))
}

func TestRPFWD_RequestValidation(t *testing.T) {
	SkipIfNoMythic(t)
	client := AuthenticateTestClient(t)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Test nil request
	_, err := client.CreateRPFWD(ctx, nil)
	if err == nil {
		t.Error("Expected error for nil request")
	}

	// Test zero callback ID
	req := &types.CreateRPFWDRequest{
		CallbackID: 0,
		LocalPort:  8080,
		RemoteHost: "localhost",
		RemotePort: 80,
	}
	_, err = client.CreateRPFWD(ctx, req)
	if err == nil {
		t.Error("Expected error for zero callback ID")
	}

	// Test zero local port
	req.CallbackID = 1
	req.LocalPort = 0
	_, err = client.CreateRPFWD(ctx, req)
	if err == nil {
		t.Error("Expected error for zero local port")
	}

	// Test invalid local port (too high)
	req.LocalPort = 70000
	_, err = client.CreateRPFWD(ctx, req)
	if err == nil {
		t.Error("Expected error for invalid local port")
	}

	// Test empty remote host
	req.LocalPort = 8080
	req.RemoteHost = ""
	_, err = client.CreateRPFWD(ctx, req)
	if err == nil {
		t.Error("Expected error for empty remote host")
	}

	// Test zero remote port
	req.RemoteHost = "localhost"
	req.RemotePort = 0
	_, err = client.CreateRPFWD(ctx, req)
	if err == nil {
		t.Error("Expected error for zero remote port")
	}

	// Test invalid remote port (too high)
	req.RemotePort = 70000
	_, err = client.CreateRPFWD(ctx, req)
	if err == nil {
		t.Error("Expected error for invalid remote port")
	}

	t.Log("All validation tests passed")
}

func TestRPFWD_HelperMethods(t *testing.T) {
	SkipIfNoMythic(t)
	client := AuthenticateTestClient(t)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Get callbacks and RPFWDs
	callbacks, err := client.GetAllActiveCallbacks(ctx)
	if err != nil {
		t.Fatalf("GetAllActiveCallbacks failed: %v", err)
	}
	if len(callbacks) == 0 {
		t.Skip("No active callbacks for testing")
	}

	rpfwds, err := client.GetRPFWDs(ctx, callbacks[0].ID)
	if err != nil {
		t.Fatalf("GetRPFWDs failed: %v", err)
	}
	if len(rpfwds) == 0 {
		t.Skip("No RPFWD tunnels for testing helper methods")
	}

	rpfwd := rpfwds[0]

	// Test String() method
	str := rpfwd.String()
	if str == "" {
		t.Error("String() should not return empty string")
	}
	t.Logf("RPFWD string: %s", str)

	// Test IsActive() method
	isActive := rpfwd.IsActive()
	if isActive != rpfwd.Active {
		t.Error("IsActive() should match Active field")
	}
	t.Logf("Is active: %v", isActive)

	// Test GetLocalEndpoint()
	localEndpoint := rpfwd.GetLocalEndpoint()
	if localEndpoint == "" {
		t.Error("GetLocalEndpoint() should not be empty")
	}
	if !contains(localEndpoint, "localhost") {
		t.Errorf("Local endpoint should contain 'localhost': %s", localEndpoint)
	}
	t.Logf("Local endpoint: %s", localEndpoint)

	// Test GetRemoteEndpoint()
	remoteEndpoint := rpfwd.GetRemoteEndpoint()
	if remoteEndpoint == "" {
		t.Error("GetRemoteEndpoint() should not be empty")
	}
	if !contains(remoteEndpoint, ":") {
		t.Errorf("Remote endpoint should be in host:port format: %s", remoteEndpoint)
	}
	t.Logf("Remote endpoint: %s", remoteEndpoint)

	// Verify endpoint formats
	if !contains(localEndpoint, ":") {
		t.Error("Local endpoint should contain port separator ':'")
	}
}

func TestRPFWD_InvalidInputs(t *testing.T) {
	SkipIfNoMythic(t)
	client := AuthenticateTestClient(t)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Test zero callback ID
	_, err := client.GetRPFWDs(ctx, 0)
	if err == nil {
		t.Error("Expected error for zero callback ID")
	}

	// Test zero RPFWD ID for status
	_, err = client.GetRPFWDStatus(ctx, 0)
	if err == nil {
		t.Error("Expected error for zero RPFWD ID in GetRPFWDStatus")
	}

	// Test zero RPFWD ID for delete
	err = client.DeleteRPFWD(ctx, 0)
	if err == nil {
		t.Error("Expected error for zero RPFWD ID in DeleteRPFWD")
	}

	// Test nonexistent RPFWD ID
	_, err = client.GetRPFWDStatus(ctx, 999999)
	if err == nil {
		t.Error("Expected error for nonexistent RPFWD ID")
	}

	t.Log("All invalid input tests passed")
}

func TestRPFWD_PortRangeValidation(t *testing.T) {
	testCases := []struct {
		name       string
		localPort  int
		remotePort int
		shouldFail bool
	}{
		{"Valid standard ports", 8080, 80, false},
		{"Valid high ports", 50000, 60000, false},
		{"Local port too low", 0, 80, true},
		{"Local port too high", 70000, 80, true},
		{"Remote port too low", 8080, 0, true},
		{"Remote port too high", 8080, 70000, true},
		{"Both ports invalid", 0, 0, true},
		{"Port 1 valid", 1, 1, false},
		{"Port 65535 valid", 65535, 65535, false},
		{"Port 65536 invalid", 65536, 80, true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := &types.CreateRPFWDRequest{
				CallbackID: 1,
				LocalPort:  tc.localPort,
				RemoteHost: "localhost",
				RemotePort: tc.remotePort,
			}

			err := req.Validate()

			if tc.shouldFail && err == nil {
				t.Errorf("Expected validation error for %s", tc.name)
			}
			if !tc.shouldFail && err != nil {
				t.Errorf("Unexpected validation error for %s: %v", tc.name, err)
			}
		})
	}
}
