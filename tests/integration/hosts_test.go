//go:build integration

package integration

import (
	"context"
	"testing"
	"time"
)

func TestHosts_GetHosts(t *testing.T) {
	client := AuthenticateTestClient(t)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	hosts, err := client.GetHosts(ctx, 0)
	if err != nil {
		t.Fatalf("GetHosts failed: %v", err)
	}

	if hosts == nil {
		t.Fatal("GetHosts returned nil")
	}

	t.Logf("Found %d host(s)", len(hosts))

	for _, host := range hosts {
		if host.ID == 0 {
			t.Error("Host ID should not be 0")
		}
		t.Logf("  - %s: %s (%s %s)", host.Hostname, host.IP, host.OS, host.Architecture)
	}
}

func TestHosts_GetHostByID(t *testing.T) {
	client := AuthenticateTestClient(t)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Get hosts first
	hosts, err := client.GetHosts(ctx, 0)
	if err != nil {
		t.Fatalf("GetHosts failed: %v", err)
	}
	if len(hosts) == 0 {
		t.Skip("No hosts available for testing")
	}

	hostID := hosts[0].ID
	host, err := client.GetHostByID(ctx, hostID)
	if err != nil {
		t.Fatalf("GetHostByID failed: %v", err)
	}

	if host == nil {
		t.Fatal("GetHostByID returned nil")
	}

	if host.ID != hostID {
		t.Errorf("Expected host ID %d, got %d", hostID, host.ID)
	}

	t.Logf("Retrieved host %d: %s", hostID, host.String())
	t.Logf("  - Hostname: %s", host.Hostname)
	t.Logf("  - IP: %s", host.IP)
	t.Logf("  - OS: %s", host.OS)
	t.Logf("  - Architecture: %s", host.Architecture)
	t.Logf("  - Domain: %s", host.Domain)
}

func TestHosts_GetHostByHostname(t *testing.T) {
	client := AuthenticateTestClient(t)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Get hosts first to find a valid hostname
	hosts, err := client.GetHosts(ctx, 0)
	if err != nil {
		t.Fatalf("GetHosts failed: %v", err)
	}
	if len(hosts) == 0 {
		t.Skip("No hosts available for testing")
	}

	hostname := hosts[0].Hostname
	host, err := client.GetHostByHostname(ctx, hostname)
	if err != nil {
		t.Fatalf("GetHostByHostname failed: %v", err)
	}

	if host == nil {
		t.Fatal("GetHostByHostname returned nil")
	}

	if host.Hostname != hostname {
		t.Errorf("Expected hostname %s, got %s", hostname, host.Hostname)
	}

	t.Logf("Found host by hostname '%s': ID %d", hostname, host.ID)
}

func TestHosts_GetCallbacksForHost(t *testing.T) {
	client := AuthenticateTestClient(t)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Get hosts first
	hosts, err := client.GetHosts(ctx, 0)
	if err != nil {
		t.Fatalf("GetHosts failed: %v", err)
	}
	if len(hosts) == 0 {
		t.Skip("No hosts available for testing")
	}

	hostID := hosts[0].ID
	callbacks, err := client.GetCallbacksForHost(ctx, hostID)
	if err != nil {
		t.Fatalf("GetCallbacksForHost failed: %v", err)
	}

	if callbacks == nil {
		t.Fatal("GetCallbacksForHost returned nil")
	}

	t.Logf("Found %d callback(s) for host %d (%s)", len(callbacks), hostID, hosts[0].Hostname)

	activeCount := 0
	for _, cb := range callbacks {
		if cb.Active {
			activeCount++
			t.Logf("  - Active: Callback %d - %s@%s (PID: %d)",
				cb.ID, cb.User, cb.Host, cb.PID)
		}
	}

	t.Logf("Active callbacks: %d of %d", activeCount, len(callbacks))

	// Verify callbacks belong to the host
	for _, cb := range callbacks {
		if cb.Host != hosts[0].Hostname {
			t.Logf("Warning: Callback hostname '%s' doesn't match host '%s'",
				cb.Host, hosts[0].Hostname)
		}
	}
}

func TestHosts_GetHostNetworkMap(t *testing.T) {
	client := AuthenticateTestClient(t)
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	networkMap, err := client.GetHostNetworkMap(ctx, 0)
	if err != nil {
		t.Fatalf("GetHostNetworkMap failed: %v", err)
	}

	if networkMap == nil {
		t.Fatal("GetHostNetworkMap returned nil")
	}

	t.Logf("Network Map:")
	t.Logf("  - Total hosts: %d", len(networkMap.Hosts))

	// Count hosts by OS
	osCounts := make(map[string]int)
	totalCallbacks := 0

	for _, host := range networkMap.Hosts {
		osCounts[host.OS]++
		if host.Callbacks != nil {
			activeCallbacks := 0
			for _, cb := range host.Callbacks {
				if cb.Active {
					activeCallbacks++
					totalCallbacks++
				}
			}
			t.Logf("  - %s (%s): %d active callback(s)",
				host.Hostname, host.IP, activeCallbacks)
		}
	}

	t.Logf("  - Total active callbacks: %d", totalCallbacks)
	t.Logf("  - Hosts by OS:")
	for os, count := range osCounts {
		if os != "" {
			t.Logf("    - %s: %d", os, count)
		} else {
			t.Logf("    - Unknown: %d", count)
		}
	}

	// Verify metadata
	if networkMap.Metadata != nil {
		if hostCount, ok := networkMap.Metadata["host_count"].(int); ok {
			if hostCount != len(networkMap.Hosts) {
				t.Errorf("Metadata host_count %d doesn't match actual count %d",
					hostCount, len(networkMap.Hosts))
			}
		}

		if totalCbs, ok := networkMap.Metadata["total_active_callbacks"].(int); ok {
			if totalCbs != totalCallbacks {
				t.Errorf("Metadata callback count %d doesn't match counted %d",
					totalCbs, totalCallbacks)
			}
		}
	}

	// Test helper methods
	for _, host := range networkMap.Hosts {
		t.Logf("Host %s:", host.Hostname)
		t.Logf("  - Has active callbacks: %v", host.HasActiveCallbacks())
		t.Logf("  - Callback count: %d", host.GetCallbackCount())
		t.Logf("  - Is Windows: %v", host.IsWindows())
		t.Logf("  - Is Linux: %v", host.IsLinux())
		t.Logf("  - Is macOS: %v", host.IsMacOS())
	}
}

func TestHosts_HelperMethods(t *testing.T) {
	client := AuthenticateTestClient(t)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Get network map with callbacks
	networkMap, err := client.GetHostNetworkMap(ctx, 0)
	if err != nil {
		t.Fatalf("GetHostNetworkMap failed: %v", err)
	}

	if len(networkMap.Hosts) == 0 {
		t.Skip("No hosts available for testing helper methods")
	}

	for _, host := range networkMap.Hosts {
		// Test String() method
		str := host.String()
		if str == "" {
			t.Error("String() should not return empty string")
		}
		t.Logf("Host string: %s", str)

		// Test callback count methods
		count := host.GetCallbackCount()
		hasActive := host.HasActiveCallbacks()

		if count > 0 && !hasActive {
			t.Error("HasActiveCallbacks() should return true when count > 0")
		}
		if count == 0 && hasActive {
			t.Error("HasActiveCallbacks() should return false when count == 0")
		}

		t.Logf("  - Active callbacks: %d (has active: %v)", count, hasActive)

		// Test OS detection methods
		osDetected := false
		if host.IsWindows() {
			t.Logf("  - Detected as Windows")
			osDetected = true
		}
		if host.IsLinux() {
			t.Logf("  - Detected as Linux")
			osDetected = true
		}
		if host.IsMacOS() {
			t.Logf("  - Detected as macOS")
			osDetected = true
		}

		if !osDetected && host.OS != "" {
			t.Logf("  - OS '%s' not detected by helper methods", host.OS)
		}
	}
}

func TestHosts_InvalidInputs(t *testing.T) {
	client := AuthenticateTestClient(t)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Test zero host ID
	_, err := client.GetHostByID(ctx, 0)
	if err == nil {
		t.Error("Expected error for zero host ID")
	}

	// Test empty hostname
	_, err = client.GetHostByHostname(ctx, "")
	if err == nil {
		t.Error("Expected error for empty hostname")
	}

	// Test zero host ID for callbacks
	_, err = client.GetCallbacksForHost(ctx, 0)
	if err == nil {
		t.Error("Expected error for zero host ID in GetCallbacksForHost")
	}

	t.Log("All invalid input tests passed")
}
