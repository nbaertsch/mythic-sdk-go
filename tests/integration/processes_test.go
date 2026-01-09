//go:build integration

package integration

import (
	"context"
	"testing"
	"time"

	"github.com/nbaertsch/mythic-sdk-go/pkg/mythic/types"
)

func TestProcesses_GetProcesses(t *testing.T) {

	client := AuthenticateTestClient(t)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	processes, err := client.GetProcesses(ctx)
	if err != nil {
		t.Fatalf("GetProcesses failed: %v", err)
	}

	if processes == nil {
		t.Fatal("GetProcesses returned nil")
	}

	t.Logf("Found %d process(es)", len(processes))

	// If there are processes, verify structure
	if len(processes) > 0 {
		proc := processes[0]
		if proc.ID == 0 {
			t.Error("Process ID should not be 0")
		}
		if proc.ProcessID == 0 {
			t.Error("Process PID should not be 0")
		}
		t.Logf("First process: %s", proc.String())
		t.Logf("  - Architecture: %s", proc.Architecture)
		t.Logf("  - User: %s", proc.User)
		t.Logf("  - Integrity: %s", proc.GetIntegrityLevelString())
	}
}

func TestProcesses_GetProcessesByOperation(t *testing.T) {

	client := AuthenticateTestClient(t)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Get current operation
	currentOpID := client.GetCurrentOperation()
	if currentOpID == nil {
		t.Skip("No current operation set")
	}

	processes, err := client.GetProcessesByOperation(ctx, *currentOpID)
	if err != nil {
		t.Fatalf("GetProcessesByOperation failed: %v", err)
	}

	if processes == nil {
		t.Fatal("GetProcessesByOperation returned nil")
	}

	t.Logf("Found %d process(es) for operation %d", len(processes), *currentOpID)

	// Verify all processes belong to the operation
	for _, proc := range processes {
		if proc.OperationID != *currentOpID {
			t.Errorf("Expected operation ID %d, got %d", *currentOpID, proc.OperationID)
		}
	}
}

func TestProcesses_GetProcessesByOperation_InvalidID(t *testing.T) {

	client := AuthenticateTestClient(t)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Test with zero ID
	_, err := client.GetProcessesByOperation(ctx, 0)
	if err == nil {
		t.Fatal("Expected error for zero operation ID, got nil")
	}

	t.Logf("Expected error: %v", err)
}

func TestProcesses_GetProcessesByCallback(t *testing.T) {

	client := AuthenticateTestClient(t)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Get active callbacks first
	callbacks, err := client.GetAllActiveCallbacks(ctx)
	if err != nil {
		t.Fatalf("GetAllActiveCallbacks failed: %v", err)
	}

	if len(callbacks) == 0 {
		t.Skip("No active callbacks available for testing")
	}

	// Get processes for first callback
	callbackID := callbacks[0].ID
	processes, err := client.GetProcessesByCallback(ctx, callbackID)
	if err != nil {
		t.Fatalf("GetProcessesByCallback failed: %v", err)
	}

	if processes == nil {
		t.Fatal("GetProcessesByCallback returned nil")
	}

	t.Logf("Found %d process(es) for callback %d", len(processes), callbackID)

	// Verify all processes belong to the callback (if callback ID is set)
	for _, proc := range processes {
		if proc.CallbackID != nil && *proc.CallbackID != callbackID {
			t.Errorf("Expected callback ID %d, got %d", callbackID, *proc.CallbackID)
		}
	}

	// Verify processes are sorted by PID
	if len(processes) > 1 {
		for i := 1; i < len(processes); i++ {
			if processes[i].ProcessID < processes[i-1].ProcessID {
				t.Error("Processes should be sorted by PID (ascending)")
				break
			}
		}
	}
}

func TestProcesses_GetProcessesByCallback_InvalidID(t *testing.T) {

	client := AuthenticateTestClient(t)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Test with zero ID
	_, err := client.GetProcessesByCallback(ctx, 0)
	if err == nil {
		t.Fatal("Expected error for zero callback ID, got nil")
	}

	t.Logf("Expected error: %v", err)
}

func TestProcesses_GetProcessTree(t *testing.T) {

	client := AuthenticateTestClient(t)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Get active callbacks first
	callbacks, err := client.GetAllActiveCallbacks(ctx)
	if err != nil {
		t.Fatalf("GetAllActiveCallbacks failed: %v", err)
	}

	if len(callbacks) == 0 {
		t.Skip("No active callbacks available for testing")
	}

	// Get process tree for first callback
	callbackID := callbacks[0].ID
	tree, err := client.GetProcessTree(ctx, callbackID)
	if err != nil {
		t.Fatalf("GetProcessTree failed: %v", err)
	}

	if tree == nil {
		t.Fatal("GetProcessTree returned nil")
	}

	t.Logf("Process tree has %d root(s) for callback %d", len(tree), callbackID)

	// Verify tree structure
	totalProcesses := countProcessesInTree(tree)
	t.Logf("Total processes in tree: %d", totalProcesses)

	// Walk the tree and verify structure
	for i, root := range tree {
		t.Logf("Root %d: %s", i+1, root.Process.String())
		if root.Process.HasParent() {
			// Check if parent exists in our process list
			t.Logf("  - Has parent PID: %d (parent may be outside our process list)", root.Process.ParentProcessID)
		}
		if len(root.Children) > 0 {
			t.Logf("  - Has %d child process(es)", len(root.Children))
			verifyTreeChildren(t, root, 1)
		}
	}
}

// countProcessesInTree counts total processes in a tree
func countProcessesInTree(trees []*types.ProcessTree) int {
	count := 0
	for _, tree := range trees {
		count++ // Count this process
		count += countProcessesInTree(tree.Children)
	}
	return count
}

// verifyTreeChildren recursively verifies tree structure
func verifyTreeChildren(t *testing.T, node *types.ProcessTree, depth int) {
	indent := ""
	for i := 0; i < depth; i++ {
		indent += "  "
	}

	for _, child := range node.Children {
		t.Logf("%s- Child: %s (PPID: %d)", indent, child.Process.String(), child.Process.ParentProcessID)

		// Verify parent-child relationship
		if child.Process.ParentProcessID != node.Process.ProcessID {
			t.Errorf("%sChild's parent PID (%d) doesn't match parent's PID (%d)",
				indent, child.Process.ParentProcessID, node.Process.ProcessID)
		}

		// Recurse into children
		if len(child.Children) > 0 {
			verifyTreeChildren(t, child, depth+1)
		}
	}
}

func TestProcesses_GetProcessTree_InvalidID(t *testing.T) {

	client := AuthenticateTestClient(t)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Test with zero ID
	_, err := client.GetProcessTree(ctx, 0)
	if err == nil {
		t.Fatal("Expected error for zero callback ID, got nil")
	}

	t.Logf("Expected error: %v", err)
}

func TestProcesses_GetProcessesByHost(t *testing.T) {

	client := AuthenticateTestClient(t)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Get some processes first to find a valid host ID
	allProcesses, err := client.GetProcesses(ctx)
	if err != nil {
		t.Fatalf("GetProcesses failed: %v", err)
	}

	if len(allProcesses) == 0 {
		t.Skip("No processes available for testing")
	}

	// Use host ID from first process
	hostID := allProcesses[0].HostID
	processes, err := client.GetProcessesByHost(ctx, hostID)
	if err != nil {
		t.Fatalf("GetProcessesByHost failed: %v", err)
	}

	if processes == nil {
		t.Fatal("GetProcessesByHost returned nil")
	}

	t.Logf("Found %d process(es) for host %d", len(processes), hostID)

	// Verify all processes belong to the host
	for _, proc := range processes {
		if proc.HostID != hostID {
			t.Errorf("Expected host ID %d, got %d", hostID, proc.HostID)
		}
	}

	// Verify processes are sorted by PID
	if len(processes) > 1 {
		for i := 1; i < len(processes); i++ {
			if processes[i].ProcessID < processes[i-1].ProcessID {
				t.Error("Processes should be sorted by PID (ascending)")
				break
			}
		}
	}
}

func TestProcesses_GetProcessesByHost_InvalidID(t *testing.T) {

	client := AuthenticateTestClient(t)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Test with zero ID
	_, err := client.GetProcessesByHost(ctx, 0)
	if err == nil {
		t.Fatal("Expected error for zero host ID, got nil")
	}

	t.Logf("Expected error: %v", err)
}

func TestProcesses_ProcessHelperMethods(t *testing.T) {

	client := AuthenticateTestClient(t)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	processes, err := client.GetProcesses(ctx)
	if err != nil {
		t.Fatalf("GetProcesses failed: %v", err)
	}

	if len(processes) == 0 {
		t.Skip("No processes available for testing helper methods")
	}

	proc := processes[0]

	// Test String() method
	str := proc.String()
	if str == "" {
		t.Error("String() should not return empty string")
	}
	t.Logf("Process string: %s", str)

	// Test IsDeleted()
	if proc.IsDeleted() {
		t.Log("Process is marked as deleted")
	} else {
		t.Log("Process is active")
	}

	// Test HasParent()
	if proc.HasParent() {
		t.Logf("Process has parent PID: %d", proc.ParentProcessID)
	} else {
		t.Log("Process is a root process (no parent)")
	}

	// Test GetIntegrityLevelString()
	integrityStr := proc.GetIntegrityLevelString()
	t.Logf("Integrity level: %s", integrityStr)
}

func TestProcesses_ProcessFields(t *testing.T) {

	client := AuthenticateTestClient(t)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	processes, err := client.GetProcesses(ctx)
	if err != nil {
		t.Fatalf("GetProcesses failed: %v", err)
	}

	if len(processes) == 0 {
		t.Skip("No processes available for testing")
	}

	proc := processes[0]

	// Log all process fields
	t.Logf("Process Details:")
	t.Logf("  - ID: %d", proc.ID)
	t.Logf("  - Name: %s", proc.Name)
	t.Logf("  - PID: %d", proc.ProcessID)
	t.Logf("  - PPID: %d", proc.ParentProcessID)
	t.Logf("  - Architecture: %s", proc.Architecture)
	t.Logf("  - Binary Path: %s", proc.BinPath)
	t.Logf("  - User: %s", proc.User)
	t.Logf("  - Command Line: %s", proc.CommandLine)
	t.Logf("  - Integrity Level: %s", proc.GetIntegrityLevelString())
	t.Logf("  - Start Time: %s", proc.StartTime.Format("2006-01-02 15:04:05"))

	if proc.Host != nil {
		t.Logf("  - Host: %s (ID: %d)", proc.Host.Host, proc.Host.ID)
	}

	if proc.CallbackID != nil {
		t.Logf("  - Callback ID: %d", *proc.CallbackID)
	}

	if proc.TaskID != nil {
		t.Logf("  - Task ID: %d", *proc.TaskID)
	}
}
