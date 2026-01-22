//go:build integration

package integration

import (
	"context"
	"testing"
	"time"

	"github.com/nbaertsch/mythic-sdk-go/pkg/mythic/types"
)

// TestE2E_ProcessRetrieval tests all process retrieval operations.
// Covers: GetProcesses, GetProcessesByOperation, GetProcessesByHost, GetProcessesByCallback
func TestE2E_ProcessRetrieval(t *testing.T) {
	client := AuthenticateTestClient(t)

	// Test 1: Get all processes for current operation
	t.Log("=== Test 1: Get all processes ===")
	ctx1, cancel1 := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel1()

	processes, err := client.GetProcesses(ctx1)
	if err != nil {
		// Process table may not exist in all Mythic versions
		if isSchemaError(err) {
			t.Skipf("Process table not available in this Mythic version: %v", err)
		}
		t.Fatalf("GetProcesses failed: %v", err)
	}
	t.Logf("✓ Retrieved %d processes for current operation", len(processes))

	if len(processes) == 0 {
		t.Log("⚠ No processes found, skipping validation tests")
		return
	}

	// Validate process structure
	for _, proc := range processes {
		if proc.ID == 0 {
			t.Error("Process has ID 0")
		}
		if proc.ProcessID == 0 {
			t.Error("Process has ProcessID 0")
		}
		if proc.OperationID == 0 {
			t.Error("Process has OperationID 0")
		}
		if proc.HostID == 0 {
			t.Error("Process has HostID 0")
		}
		if proc.Host == nil {
			t.Error("Process has nil Host")
		}
		if proc.Deleted {
			t.Error("GetProcesses returned deleted process")
		}
	}

	// Show sample processes
	sampleCount := 5
	if len(processes) < sampleCount {
		sampleCount = len(processes)
	}
	t.Logf("  Sample processes:")
	for i := 0; i < sampleCount; i++ {
		proc := processes[i]
		t.Logf("    [%d] PID %d: %s (User: %s, Host: %s)",
			i+1, proc.ProcessID, proc.Name, proc.User, proc.Host.Host)
	}

	// Store for later tests
	testOperationID := processes[0].OperationID
	testHostID := processes[0].HostID
	var testCallbackID *int
	if processes[0].CallbackID != nil {
		testCallbackID = processes[0].CallbackID
	}

	// Test 2: Get processes by operation
	t.Log("=== Test 2: Get processes by operation ===")
	ctx2, cancel2 := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel2()

	opProcesses, err := client.GetProcessesByOperation(ctx2, testOperationID)
	if err != nil {
		t.Fatalf("GetProcessesByOperation failed: %v", err)
	}
	t.Logf("✓ Retrieved %d processes for operation %d", len(opProcesses), testOperationID)

	// Verify all processes belong to operation
	for _, proc := range opProcesses {
		if proc.OperationID != testOperationID {
			t.Errorf("Process %d has wrong OperationID: expected %d, got %d",
				proc.ID, testOperationID, proc.OperationID)
		}
	}

	// Test 3: Get processes by host
	t.Log("=== Test 3: Get processes by host ===")
	ctx3, cancel3 := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel3()

	hostProcesses, err := client.GetProcessesByHost(ctx3, testHostID)
	if err != nil {
		t.Fatalf("GetProcessesByHost failed: %v", err)
	}
	t.Logf("✓ Retrieved %d processes for host %d", len(hostProcesses), testHostID)

	// Verify all processes belong to host
	for _, proc := range hostProcesses {
		if proc.HostID != testHostID {
			t.Errorf("Process %d has wrong HostID: expected %d, got %d",
				proc.ID, testHostID, proc.HostID)
		}
	}

	// Verify ordering (should be by process_id ascending)
	if len(hostProcesses) > 1 {
		for i := 0; i < len(hostProcesses)-1; i++ {
			if hostProcesses[i].ProcessID > hostProcesses[i+1].ProcessID {
				t.Error("Host processes not ordered by ProcessID (ascending)")
			}
		}
		t.Log("  ✓ Processes correctly ordered by ProcessID")
	}

	// Test 4: Get processes by callback
	if testCallbackID != nil {
		t.Log("=== Test 4: Get processes by callback ===")
		ctx4, cancel4 := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel4()

		cbProcesses, err := client.GetProcessesByCallback(ctx4, *testCallbackID)
		if err != nil {
			t.Fatalf("GetProcessesByCallback failed: %v", err)
		}
		t.Logf("✓ Retrieved %d processes for callback %d", len(cbProcesses), *testCallbackID)

		// Verify all processes belong to callback
		for _, proc := range cbProcesses {
			if proc.CallbackID == nil || *proc.CallbackID != *testCallbackID {
				t.Errorf("Process %d has wrong or nil CallbackID", proc.ID)
			}
		}
	} else {
		t.Log("⚠ No callback ID available, skipping callback process test")
	}

	t.Log("=== ✓ All process retrieval tests passed ===")
}

// TestE2E_ProcessTree tests process tree construction and validation.
// Covers: GetProcessTree, tree structure validation
func TestE2E_ProcessTree(t *testing.T) {
	client := AuthenticateTestClient(t)

	// Find a callback with processes
	t.Log("=== Setup: Find callback with processes ===")
	ctx0, cancel0 := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel0()

	callbacks, err := client.GetAllActiveCallbacks(ctx0)
	if err != nil {
		t.Fatalf("GetAllActiveCallbacks failed: %v", err)
	}

	var testCallback *types.Callback
	for _, cb := range callbacks {
		if cb.Active {
			testCallback = cb
			break
		}
	}

	if testCallback == nil {
		t.Skip("No active callbacks found, skipping process tree test")
	}

	// Check if callback has processes
	ctx1, cancel1 := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel1()

	processes, err := client.GetProcessesByCallback(ctx1, testCallback.ID)
	if err != nil {
		if isSchemaError(err) {
			t.Skipf("Process table not available in this Mythic version: %v", err)
		}
		t.Fatalf("GetProcessesByCallback failed: %v", err)
	}

	if len(processes) == 0 {
		t.Skip("Callback has no processes, skipping process tree test")
	}

	t.Logf("Using callback %d with %d processes", testCallback.ID, len(processes))

	// Test: Get process tree
	t.Log("=== Test: Get process tree ===")
	ctx2, cancel2 := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel2()

	tree, err := client.GetProcessTree(ctx2, testCallback.ID)
	if err != nil {
		t.Fatalf("GetProcessTree failed: %v", err)
	}

	if len(tree) == 0 {
		t.Fatal("GetProcessTree returned empty tree but processes exist")
	}

	t.Logf("✓ Process tree has %d root nodes", len(tree))

	// Validate tree structure
	nodeCount := 0
	var validateTree func([]*types.ProcessTree, int)
	validateTree = func(nodes []*types.ProcessTree, depth int) {
		for _, node := range nodes {
			nodeCount++
			if node.Process == nil {
				t.Error("Process tree node has nil Process")
				continue
			}

			// Validate process data
			if node.Process.ProcessID == 0 {
				t.Error("Process in tree has ProcessID 0")
			}

			// Log structure (first few levels only)
			if depth < 3 {
				indent := ""
				for i := 0; i < depth; i++ {
					indent += "  "
				}
				t.Logf("%s- PID %d: %s (Parent: %d, Children: %d)",
					indent, node.Process.ProcessID, node.Process.Name,
					node.Process.ParentProcessID, len(node.Children))
			}

			// Recursively validate children
			if len(node.Children) > 0 {
				validateTree(node.Children, depth+1)
			}
		}
	}

	validateTree(tree, 0)
	t.Logf("  ✓ Tree contains %d total nodes", nodeCount)

	// Verify node count matches process count
	if nodeCount != len(processes) {
		t.Errorf("Tree node count (%d) doesn't match process count (%d)",
			nodeCount, len(processes))
	}

	// Validate parent-child relationships
	t.Log("  Validating parent-child relationships...")
	processMap := make(map[int]*types.Process)
	for _, proc := range processes {
		processMap[proc.ProcessID] = proc
	}

	var validateRelationships func([]*types.ProcessTree)
	validateRelationships = func(nodes []*types.ProcessTree) {
		for _, node := range nodes {
			for _, child := range node.Children {
				if child.Process.ParentProcessID != node.Process.ProcessID {
					t.Errorf("Invalid parent-child relationship: child %d has ParentProcessID %d but is child of %d",
						child.Process.ProcessID, child.Process.ParentProcessID, node.Process.ProcessID)
				}
			}
			if len(node.Children) > 0 {
				validateRelationships(node.Children)
			}
		}
	}

	validateRelationships(tree)
	t.Log("  ✓ All parent-child relationships valid")

	t.Log("=== ✓ Process tree tests passed ===")
}

// TestE2E_ProcessAttributes tests process attribute validation and analysis.
func TestE2E_ProcessAttributes(t *testing.T) {
	client := AuthenticateTestClient(t)

	t.Log("=== Test: Analyze process attributes ===")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	processes, err := client.GetProcesses(ctx)
	if err != nil {
		if isSchemaError(err) {
			t.Skipf("Process table not available in this Mythic version: %v", err)
		}
		t.Fatalf("GetProcesses failed: %v", err)
	}

	if len(processes) == 0 {
		t.Skip("No processes found for attribute analysis")
	}

	t.Logf("✓ Analyzing %d processes", len(processes))

	// Analyze attributes
	architectures := make(map[string]int)
	users := make(map[string]int)
	integrityLevels := make(map[int]int)
	withCallback := 0
	withTask := 0
	withCommandLine := 0

	for _, proc := range processes {
		if proc.Architecture != "" {
			architectures[proc.Architecture]++
		}
		if proc.User != "" {
			users[proc.User]++
		}
		integrityLevels[proc.IntegrityLevel]++

		if proc.CallbackID != nil {
			withCallback++
		}
		if proc.TaskID != nil {
			withTask++
		}
		if proc.CommandLine != "" {
			withCommandLine++
		}
	}

	t.Logf("  Architecture distribution:")
	for arch, count := range architectures {
		t.Logf("    %s: %d", arch, count)
	}

	t.Logf("  Top users:")
	userCount := 0
	for user, count := range users {
		if userCount < 5 {
			t.Logf("    %s: %d", user, count)
			userCount++
		}
	}

	t.Logf("  Integrity levels:")
	for level, count := range integrityLevels {
		t.Logf("    Level %d: %d", level, count)
	}

	t.Logf("  Metadata:")
	t.Logf("    With callback: %d", withCallback)
	t.Logf("    With task: %d", withTask)
	t.Logf("    With command line: %d", withCommandLine)

	// Find interesting processes
	var privilegedProcesses []*types.Process
	var systemProcesses []*types.Process

	for _, proc := range processes {
		// High integrity level or SYSTEM user
		if proc.IntegrityLevel >= 3 || proc.User == "SYSTEM" || proc.User == "NT AUTHORITY\\SYSTEM" {
			privilegedProcesses = append(privilegedProcesses, proc)
		}
		// Common system processes
		if proc.Name == "lsass.exe" || proc.Name == "winlogon.exe" || proc.Name == "services.exe" {
			systemProcesses = append(systemProcesses, proc)
		}
	}

	if len(privilegedProcesses) > 0 {
		t.Logf("  Found %d privileged processes", len(privilegedProcesses))
		showCount := 3
		if len(privilegedProcesses) < showCount {
			showCount = len(privilegedProcesses)
		}
		for i := 0; i < showCount; i++ {
			proc := privilegedProcesses[i]
			t.Logf("    PID %d: %s (User: %s, Integrity: %d)",
				proc.ProcessID, proc.Name, proc.User, proc.IntegrityLevel)
		}
	}

	if len(systemProcesses) > 0 {
		t.Logf("  Found %d critical system processes", len(systemProcesses))
		for _, proc := range systemProcesses {
			t.Logf("    PID %d: %s", proc.ProcessID, proc.Name)
		}
	}

	t.Log("=== ✓ Process attribute analysis complete ===")
}

// TestE2E_ProcessErrorHandling tests error scenarios for process operations.
func TestE2E_ProcessErrorHandling(t *testing.T) {
	client := AuthenticateTestClient(t)

	// Quick check if process table exists
	ctxCheck, cancelCheck := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelCheck()
	_, err := client.GetProcesses(ctxCheck)
	if isSchemaError(err) {
		t.Skipf("Process table not available in this Mythic version: %v", err)
	}

	// Test 1: Get processes by invalid operation
	t.Log("=== Test 1: Get processes by invalid operation ===")
	ctx1, cancel1 := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel1()

	_, err = client.GetProcessesByOperation(ctx1, 0)
	if err == nil {
		t.Error("Expected error for invalid operation ID")
	}
	t.Logf("✓ Invalid operation ID rejected: %v", err)

	// Test 2: Get processes by non-existent operation
	t.Log("=== Test 2: Get processes by non-existent operation ===")
	ctx2, cancel2 := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel2()

	processes, err := client.GetProcessesByOperation(ctx2, 999999)
	if err != nil {
		t.Logf("⚠ Query failed (may be valid): %v", err)
	} else {
		if len(processes) > 0 {
			t.Error("Non-existent operation returned processes")
		}
		t.Log("✓ Non-existent operation returns empty array")
	}

	// Test 3: Get processes by invalid host
	t.Log("=== Test 3: Get processes by invalid host ===")
	ctx3, cancel3 := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel3()

	_, err = client.GetProcessesByHost(ctx3, 0)
	if err == nil {
		t.Error("Expected error for invalid host ID")
	}
	t.Logf("✓ Invalid host ID rejected: %v", err)

	// Test 4: Get processes by invalid callback
	t.Log("=== Test 4: Get processes by invalid callback ===")
	ctx4, cancel4 := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel4()

	_, err = client.GetProcessesByCallback(ctx4, 0)
	if err == nil {
		t.Error("Expected error for invalid callback ID")
	}
	t.Logf("✓ Invalid callback ID rejected: %v", err)

	// Test 5: Get process tree with invalid callback
	t.Log("=== Test 5: Get process tree with invalid callback ===")
	ctx5, cancel5 := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel5()

	_, err = client.GetProcessTree(ctx5, 0)
	if err == nil {
		t.Error("Expected error for invalid callback ID")
	}
	t.Logf("✓ Invalid callback ID rejected: %v", err)

	// Test 6: Get process tree for callback with no processes
	t.Log("=== Test 6: Get process tree for callback with no processes ===")
	ctx6, cancel6 := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel6()

	tree, err := client.GetProcessTree(ctx6, 999999)
	if err != nil {
		t.Logf("⚠ Query failed (may be valid): %v", err)
	} else {
		if len(tree) == 0 {
			t.Log("✓ Non-existent callback returns empty tree")
		}
	}

	t.Log("=== ✓ All error handling tests passed ===")
}

// TestE2E_ProcessTimestamps tests process timestamp ordering and filtering.
func TestE2E_ProcessTimestamps(t *testing.T) {
	client := AuthenticateTestClient(t)

	t.Log("=== Test: Process timestamp ordering ===")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	processes, err := client.GetProcesses(ctx)
	if err != nil {
		if isSchemaError(err) {
			t.Skipf("Process table not available in this Mythic version: %v", err)
		}
		t.Fatalf("GetProcesses failed: %v", err)
	}

	if len(processes) == 0 {
		t.Skip("No processes found for timestamp test")
	}

	t.Logf("✓ Retrieved %d processes", len(processes))

	// Verify timestamp ordering (should be desc for GetProcesses)
	if len(processes) > 1 {
		for i := 0; i < len(processes)-1; i++ {
			if processes[i].Timestamp.Before(processes[i+1].Timestamp) {
				t.Error("Processes not ordered by timestamp (descending)")
			}
		}
		t.Log("  ✓ Processes correctly ordered by timestamp (most recent first)")
	}

	// Analyze timestamp distribution
	now := time.Now()
	last24h := 0
	lastWeek := 0
	older := 0

	for _, proc := range processes {
		age := now.Sub(proc.Timestamp)
		if age < 24*time.Hour {
			last24h++
		} else if age < 7*24*time.Hour {
			lastWeek++
		} else {
			older++
		}
	}

	t.Logf("  Timestamp distribution:")
	t.Logf("    Last 24 hours: %d", last24h)
	t.Logf("    Last week: %d", lastWeek)
	t.Logf("    Older: %d", older)

	// Show age of newest and oldest
	if len(processes) > 0 {
		newest := processes[0].Timestamp
		oldest := processes[len(processes)-1].Timestamp
		t.Logf("  Newest process: %s (age: %s)", newest.Format(time.RFC3339), now.Sub(newest))
		t.Logf("  Oldest process: %s (age: %s)", oldest.Format(time.RFC3339), now.Sub(oldest))
	}

	t.Log("=== ✓ Timestamp ordering tests passed ===")
}
