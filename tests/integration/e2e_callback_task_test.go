//go:build integration

package integration

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/nbaertsch/mythic-sdk-go/pkg/mythic"
	"github.com/nbaertsch/mythic-sdk-go/pkg/mythic/types"
)

// TestE2E_CallbackTaskLifecycle tests the complete callback and task execution workflow
// This is the most complex E2E test, requiring agent deployment and task execution.
// Covers: Callbacks, Tasks, Task Responses, Processes, Hosts
func TestE2E_CallbackTaskLifecycle(t *testing.T) {
	// This test requires Poseidon agent deployment
	// Skip if running in environments without Mythic server
	if os.Getenv("SKIP_AGENT_TESTS") == "true" {
		t.Skip("Skipping agent-dependent test (SKIP_AGENT_TESTS=true)")
	}

	setup := SetupE2ETest(t)
	defer setup.Cleanup()

	// PART 1: Build and Deploy Payload
	t.Log("========== PART 1: Build and Deploy Payload ==========")

	// Test 1: Get payload types and find Poseidon
	t.Log("=== Test 1: Find Poseidon payload type ===")
	ctx1, cancel1 := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel1()

	payloadTypes, err := setup.Client.GetPayloadTypes(ctx1)
	if err != nil {
		t.Fatalf("GetPayloadTypes failed: %v", err)
	}

	var poseidonType *types.PayloadType
	for _, pt := range payloadTypes {
		if pt.Name == "poseidon" {
			poseidonType = pt
			break
		}
	}

	if poseidonType == nil {
		t.Skip("Poseidon payload type not found - skipping callback/task tests")
	}

	if !poseidonType.ContainerRunning {
		t.Skipf("Poseidon container not running - skipping callback/task tests")
	}

	t.Logf("✓ Found Poseidon payload type (ID: %d)", poseidonType.ID)

	// Test 2: Get C2 profile for payload
	t.Log("=== Test 2: Get C2 profile for payload ===")
	ctx2, cancel2 := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel2()

	c2Profiles, err := setup.Client.GetC2Profiles(ctx2)
	if err != nil {
		t.Fatalf("GetC2Profiles failed: %v", err)
	}

	var c2Profile *types.C2Profile
	for _, profile := range c2Profiles {
		if profile.Name == "http" {
			c2Profile = profile
			break
		}
	}
	if c2Profile == nil {
		c2Profile = c2Profiles[0]
	}

	t.Logf("✓ Using C2 profile: %s", c2Profile.Name)

	// Test 3: Create payload for agent deployment
	t.Log("=== Test 3: Create Poseidon payload ===")
	ctx3, cancel3 := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel3()

	payloadReq := &types.CreatePayloadRequest{
		PayloadType: "poseidon",
		OS:          "linux",
		Description: "E2E Callback Test Payload",
		Filename:    "e2e_callback_test",
		Commands: []string{
			"shell", "ps", "whoami",
		},
		BuildParameters: map[string]interface{}{
			"mode": "default",
		},
		C2Profiles: []types.C2ProfileConfig{
			{
				Name: c2Profile.Name,
				Parameters: map[string]interface{}{
					"callback_host": "http://127.0.0.1",
					"callback_port": 80,
				},
			},
		},
	}

	payload, err := setup.Client.CreatePayload(ctx3, payloadReq)
	if err != nil {
		t.Fatalf("CreatePayload failed: %v", err)
	}
	setup.PayloadUUID = payload.UUID
	t.Logf("✓ Payload created: UUID %s", payload.UUID)

	// Test 4: Wait for payload build
	t.Log("=== Test 4: Wait for payload build ===")
	ctx4, cancel4 := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel4()

	err = setup.Client.WaitForPayloadComplete(ctx4, payload.UUID, 90)
	if err != nil {
		t.Fatalf("Payload build failed: %v", err)
	}
	t.Log("✓ Payload build completed")

	// Test 5: Download payload binary
	t.Log("=== Test 5: Download payload binary ===")
	ctx5, cancel5 := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel5()

	payloadBytes, err := setup.Client.DownloadPayload(ctx5, payload.UUID)
	if err != nil {
		t.Fatalf("DownloadPayload failed: %v", err)
	}

	// Save to temp file
	tmpDir := os.TempDir()
	payloadPath := filepath.Join(tmpDir, "e2e_test_agent_"+payload.UUID[:8])
	err = os.WriteFile(payloadPath, payloadBytes, 0755)
	if err != nil {
		t.Fatalf("Failed to write payload: %v", err)
	}
	setup.PayloadPath = payloadPath
	setup.AddTempFile(payloadPath)

	t.Logf("✓ Payload downloaded and saved: %s (%d bytes)", payloadPath, len(payloadBytes))

	// PART 2: Callback Establishment
	t.Log("========== PART 2: Callback Establishment ==========")

	// Test 6: Get baseline callback count
	t.Log("=== Test 6: Get baseline callback count ===")
	ctx6, cancel6 := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel6()

	baselineCallbacks, err := setup.Client.GetAllActiveCallbacks(ctx6)
	if err != nil {
		t.Fatalf("GetAllActiveCallbacks failed: %v", err)
	}
	t.Logf("✓ Baseline active callbacks: %d", len(baselineCallbacks))

	// Test 7: Start agent
	t.Log("=== Test 7: Start agent ===")
	err = setup.StartAgent()
	if err != nil {
		t.Fatalf("Failed to start agent: %v", err)
	}
	t.Log("✓ Agent started")

	// Test 8: Wait for callback (up to 90 seconds - increased for CI reliability)
	t.Log("=== Test 8: Wait for callback (up to 90 seconds) ===")
	callbackID, err := setup.WaitForCallback(90 * time.Second)
	if err != nil {
		t.Fatalf("Failed to establish callback: %v", err)
	}
	t.Logf("✓ Callback established: ID %d", callbackID)

	// Test 9: Get callback details
	t.Log("=== Test 9: Get callback by ID ===")
	ctx7, cancel7 := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel7()

	callback, err := setup.Client.GetCallbackByID(ctx7, callbackID)
	if err != nil {
		t.Fatalf("GetCallbackByID failed: %v", err)
	}
	t.Logf("✓ Callback details: Host=%s, User=%s, PID=%d", callback.Host, callback.User, callback.PID)

	// Test 10: Update callback description
	t.Log("=== Test 10: Update callback description ===")
	ctx8, cancel8 := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel8()

	newDesc := "E2E Test Callback - Updated"
	updateReq := &types.CallbackUpdateRequest{
		CallbackDisplayID: callbackID,
		Description:       &newDesc,
	}

	err = setup.Client.UpdateCallback(ctx8, updateReq)
	if err != nil {
		t.Fatalf("UpdateCallback failed: %v", err)
	}
	t.Log("✓ Callback description updated")

	// Test 11: Get loaded commands
	t.Log("=== Test 11: Get loaded commands ===")
	ctx9, cancel9 := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel9()

	commands, err := setup.Client.GetLoadedCommands(ctx9, callbackID)
	if err != nil {
		t.Logf("⚠ GetLoadedCommands failed: %v", err)
	} else {
		t.Logf("✓ Loaded commands: %d", len(commands))
		for i, cmd := range commands {
			if i < 5 { // Show first 5
				t.Logf("  - %s", cmd)
			}
		}
	}

	// PART 3: Task Execution - Shell Command
	t.Log("========== PART 3: Task Execution - Shell Command ==========")

	// Test 12: Issue shell task (whoami)
	t.Log("=== Test 12: Issue shell task (whoami) ===")
	// Note: shell command expects params as plain command string, not JSON
	taskDisplayID, err := setup.ExecuteCommand("shell", "whoami")
	if err != nil {
		t.Fatalf("Failed to issue shell task: %v", err)
	}
	t.Logf("✓ Shell task issued: Display ID %d", taskDisplayID)

	// Test 13: Get task immediately
	t.Log("=== Test 13: Get task details ===")
	ctx10, cancel10 := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel10()

	task, err := setup.Client.GetTask(ctx10, taskDisplayID)
	if err != nil {
		t.Fatalf("GetTask failed: %v", err)
	}
	t.Logf("✓ Task status: %s (Completed: %v)", task.Status, task.Completed)

	// Test 14: Wait for task completion
	t.Log("=== Test 14: Wait for task completion ===")
	output, err := setup.WaitForTaskComplete(taskDisplayID, 30*time.Second)
	if err != nil {
		t.Fatalf("Task did not complete: %v", err)
	}
	t.Logf("✓ Task completed with output length: %d bytes", len(output))

	// Test 15: Get task output responses
	t.Log("=== Test 15: Get task output responses ===")
	ctx11, cancel11 := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel11()

	responses, err := setup.Client.GetTaskOutput(ctx11, taskDisplayID)
	if err != nil {
		t.Fatalf("GetTaskOutput failed: %v", err)
	}
	t.Logf("✓ Task has %d responses", len(responses))

	// Test 16: Get responses by task
	t.Log("=== Test 16: Get responses by task ===")
	ctx12, cancel12 := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel12()

	taskResponses, err := setup.Client.GetResponsesByTask(ctx12, task.ID)
	if err != nil {
		t.Fatalf("GetResponsesByTask failed: %v", err)
	}
	t.Logf("✓ Found %d responses for task", len(taskResponses))

	// PART 4: Process List
	t.Log("========== PART 4: Process List ==========")

	// Test 17: Issue ps task
	t.Log("=== Test 17: Issue ps task ===")
	psTaskID, err := setup.ExecuteCommand("ps", "")
	if err != nil {
		t.Fatalf("Failed to issue ps task: %v", err)
	}
	t.Logf("✓ PS task issued: Display ID %d", psTaskID)

	// Test 18: Wait for ps task completion
	t.Log("=== Test 18: Wait for ps task completion ===")
	_, err = setup.WaitForTaskComplete(psTaskID, 30*time.Second)
	if err != nil {
		t.Logf("⚠ PS task did not complete: %v", err)
	} else {
		t.Log("✓ PS task completed")
	}

	// Test 19: Get processes for callback
	t.Log("=== Test 19: Get processes by callback ===")
	ctx13, cancel13 := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel13()

	processes, err := setup.Client.GetProcessesByCallback(ctx13, callbackID)
	if err != nil {
		t.Logf("⚠ GetProcessesByCallback failed: %v", err)
	} else {
		t.Logf("✓ Found %d processes for callback", len(processes))
	}

	// Test 20: Get process tree
	t.Log("=== Test 20: Get process tree for callback ===")
	ctx14, cancel14 := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel14()

	processTree, err := setup.Client.GetProcessTree(ctx14, callbackID)
	if err != nil {
		t.Logf("⚠ GetProcessTree failed: %v", err)
	} else {
		t.Logf("✓ Found %d processes in tree", len(processTree))
	}

	// PART 5: Host Enumeration
	t.Log("========== PART 5: Host Enumeration ==========")

	// Test 21: Get all hosts
	t.Log("=== Test 21: Get all hosts ===")
	ctx15, cancel15 := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel15()

	hosts, err := setup.Client.GetHosts(ctx15, setup.OperationID)
	if err != nil {
		// Host table may not exist in all Mythic versions
		if isSchemaError(err) {
			t.Logf("⚠ GetHosts failed (host table not available in this Mythic version): %v", err)
			t.Log("Skipping host-related tests")
			goto skipHostTests
		}
		t.Fatalf("GetHosts failed: %v", err)
	}
	t.Logf("✓ Found %d hosts", len(hosts))

	// Test 22: Find callback's host
	t.Log("=== Test 22: Find callback's host ===")
	var callbackHost *types.HostInfo
	for _, h := range hosts {
		if h.Hostname == callback.Host {
			callbackHost = h
			break
		}
	}

	if callbackHost != nil {
		t.Logf("✓ Found callback host: %s (ID: %d)", callbackHost.Hostname, callbackHost.ID)

		// Test 23: Get host by ID
		t.Log("=== Test 23: Get host by ID ===")
		ctx16, cancel16 := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel16()

		hostByID, err := setup.Client.GetHostByID(ctx16, callbackHost.ID)
		if err != nil {
			t.Errorf("GetHostByID failed: %v", err)
		} else {
			t.Logf("✓ Host retrieved: %s", hostByID.Hostname)
		}

		// Test 24: Get callbacks for host
		t.Log("=== Test 24: Get callbacks for host ===")
		ctx17, cancel17 := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel17()

		hostCallbacks, err := setup.Client.GetCallbacksForHost(ctx17, callbackHost.ID)
		if err != nil {
			t.Errorf("GetCallbacksForHost failed: %v", err)
		} else {
			t.Logf("✓ Found %d callbacks for host", len(hostCallbacks))
		}
	} else {
		t.Log("⚠ Could not find callback's host in hosts list")
	}

skipHostTests:
	// PART 6: Task Management
	t.Log("========== PART 6: Task Management ==========")

	// Test 25: Get tasks for callback
	t.Log("=== Test 25: Get tasks for callback ===")
	ctx18, cancel18 := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel18()

	callbackTasks, err := setup.Client.GetTasksForCallback(ctx18, callbackID, 100)
	if err != nil {
		t.Fatalf("GetTasksForCallback failed: %v", err)
	}
	t.Logf("✓ Found %d tasks for callback", len(callbackTasks))

	// Test 26: Get tasks by status (completed)
	t.Log("=== Test 26: Get completed tasks ===")
	ctx19, cancel19 := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel19()

	completedTasks, err := setup.Client.GetTasksByStatus(ctx19, callbackID, mythic.TaskStatusCompleted, 100)
	if err != nil {
		t.Fatalf("GetTasksByStatus failed: %v", err)
	}
	t.Logf("✓ Found %d completed tasks", len(completedTasks))

	// Test 27: Update task comment
	t.Log("=== Test 27: Update task comment ===")
	ctx20, cancel20 := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel20()

	taskUpdates := map[string]interface{}{
		"comment": "E2E Test Comment",
	}

	err = setup.Client.UpdateTask(ctx20, taskDisplayID, taskUpdates)
	if err != nil {
		t.Fatalf("UpdateTask failed: %v", err)
	}
	t.Log("✓ Task comment updated")

	// PART 7: Callback Cleanup
	t.Log("========== PART 7: Callback Cleanup ==========")

	// Test 28: Get callback tokens
	t.Log("=== Test 28: Get callback tokens ===")
	ctx21, cancel21 := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel21()

	tokens, err := setup.Client.GetCallbackTokensByCallback(ctx21, callbackID)
	if err != nil {
		t.Logf("⚠ GetCallbackTokensByCallback failed: %v", err)
	} else {
		t.Logf("✓ Found %d tokens for callback", len(tokens))
	}

	// Test 29: Verify callback in list
	t.Log("=== Test 29: Verify callback in all callbacks list ===")
	ctx22, cancel22 := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel22()

	allCallbacks, err := setup.Client.GetAllCallbacks(ctx22)
	if err != nil {
		t.Fatalf("GetAllCallbacks failed: %v", err)
	}

	found := false
	for _, cb := range allCallbacks {
		if cb.ID == callbackID {
			found = true
			break
		}
	}
	if !found {
		t.Error("Callback not found in all callbacks list")
	} else {
		t.Logf("✓ Callback found in list (total callbacks: %d)", len(allCallbacks))
	}

	// Test 30: Delete callback (will be cleaned up in defer as well)
	t.Log("=== Test 30: Delete callback ===")
	ctx23, cancel23 := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel23()

	err = setup.Client.DeleteCallback(ctx23, []int{callbackID})
	if err != nil {
		t.Errorf("DeleteCallback failed: %v", err)
	} else {
		t.Log("✓ Callback deleted successfully")
		setup.CallbackID = 0 // Prevent double-delete in cleanup
	}

	t.Log("========== ✓ All callback/task tests passed ==========")
}

// TestE2E_CallbackTaskErrorHandling tests error scenarios for callback/task operations
func TestE2E_CallbackTaskErrorHandling(t *testing.T) {
	client := AuthenticateTestClient(t)

	// Test 1: Get non-existent callback
	t.Log("=== Test 1: Get non-existent callback ===")
	ctx1, cancel1 := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel1()

	_, err := client.GetCallbackByID(ctx1, 999999)
	if err == nil {
		t.Error("Expected error for non-existent callback ID")
	}
	t.Logf("✓ Non-existent callback rejected: %v", err)

	// Test 2: Get non-existent task
	t.Log("=== Test 2: Get non-existent task ===")
	ctx2, cancel2 := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel2()

	_, err = client.GetTask(ctx2, 999999)
	if err == nil {
		t.Error("Expected error for non-existent task display ID")
	}
	t.Logf("✓ Non-existent task rejected: %v", err)

	// Test 3: Issue task to non-existent callback
	t.Log("=== Test 3: Issue task to non-existent callback ===")
	ctx3, cancel3 := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel3()

	nonExistentCallbackID := 999999
	taskReq := &mythic.TaskRequest{
		CallbackID: &nonExistentCallbackID,
		Command:    "shell",
		Params:     "whoami",
	}

	_, err = client.IssueTask(ctx3, taskReq)
	if err == nil {
		t.Error("Expected error for task on non-existent callback")
	}
	t.Logf("✓ Task to non-existent callback rejected: %v", err)

	// Test 4: Get tasks for non-existent callback
	t.Log("=== Test 4: Get tasks for non-existent callback ===")
	ctx4, cancel4 := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel4()

	_, err = client.GetTasksForCallback(ctx4, 999999, 100)
	if err == nil {
		t.Error("Expected error for tasks of non-existent callback")
	}
	t.Logf("✓ Non-existent callback tasks rejected: %v", err)

	t.Log("=== ✓ All error handling tests passed ===")
}
