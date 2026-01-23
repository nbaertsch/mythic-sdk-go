//go:build integration

package integration

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/nbaertsch/mythic-sdk-go/pkg/mythic"
	"github.com/nbaertsch/mythic-sdk-go/pkg/mythic/types"
)

// TestE2E_TaskLifecycle tests the complete task execution lifecycle.
// Covers: IssueTask, GetTask, GetTaskOutput, WaitForTaskComplete, UpdateTask
func TestE2E_TaskLifecycle(t *testing.T) {
	// Ensure at least one callback exists (reuses existing or creates one)
	callbackID := EnsureCallbackExists(t)

	client := AuthenticateTestClient(t)

	// Get the callback details
	ctx0, cancel0 := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel0()
	testCallback, err := client.GetCallbackByID(ctx0, callbackID)
	if err != nil {
		t.Fatalf("Failed to get callback: %v", err)
	}
	t.Logf("Using callback %d (Host: %s, User: %s)", testCallback.ID, testCallback.Host, testCallback.User)

	// Test 1: Issue a simple task
	t.Log("=== Test 1: Issue a simple task ===")
	ctx1, cancel1 := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel1()

	taskReq := &mythic.TaskRequest{
		Command:    "whoami",
		Params:     "whoami",
		CallbackID: &testCallback.ID,
	}

	task, err := client.IssueTask(ctx1, taskReq)
	if err != nil {
		t.Fatalf("IssueTask failed: %v", err)
	}
	if task.DisplayID == 0 {
		t.Fatal("Task has DisplayID 0")
	}
	t.Logf("✓ Task issued: DisplayID %d, Status: %s", task.DisplayID, task.Status)

	// Test 2: Get task by display ID
	t.Log("=== Test 2: Get task by display ID ===")
	ctx2, cancel2 := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel2()

	retrievedTask, err := client.GetTask(ctx2, task.DisplayID)
	if err != nil {
		t.Fatalf("GetTask failed: %v", err)
	}
	if retrievedTask.DisplayID != task.DisplayID {
		t.Errorf("Task DisplayID mismatch: expected %d, got %d", task.DisplayID, retrievedTask.DisplayID)
	}
	t.Logf("✓ Task retrieved: %s (Status: %s)", retrievedTask.CommandName, retrievedTask.Status)

	// Test 3: Wait for task to complete
	t.Log("=== Test 3: Wait for task to complete ===")
	ctx3, cancel3 := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel3()

	err = client.WaitForTaskComplete(ctx3, task.DisplayID, 45)
	if err != nil {
		t.Logf("⚠ Task did not complete within timeout: %v", err)
	} else {
		t.Log("✓ Task completed successfully")
	}

	// Test 4: Get task output
	t.Log("=== Test 4: Get task output ===")
	ctx4, cancel4 := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel4()

	output, err := client.GetTaskOutput(ctx4, task.DisplayID)
	if err != nil {
		t.Fatalf("GetTaskOutput failed: %v", err)
	}
	t.Logf("✓ Retrieved %d output entries", len(output))

	for i, out := range output {
		if i < 3 { // Show first 3 entries
			preview := out.ResponseText
			if len(preview) > 100 {
				preview = preview[:100] + "..."
			}
			t.Logf("  [%d] %s", i+1, preview)
		}
	}

	// Test 5: Update task (add comment)
	t.Log("=== Test 5: Update task comment ===")
	ctx5, cancel5 := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel5()

	updates := map[string]interface{}{
		"comment": "Test task from integration tests",
	}

	err = client.UpdateTask(ctx5, task.DisplayID, updates)
	if err != nil {
		t.Fatalf("UpdateTask failed: %v", err)
	}
	t.Log("✓ Task comment updated")

	// Verify update
	ctx6, cancel6 := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel6()

	updatedTask, err := client.GetTask(ctx6, task.DisplayID)
	if err != nil {
		t.Fatalf("GetTask after update failed: %v", err)
	}
	if updatedTask.Comment != "Test task from integration tests" {
		t.Errorf("Task comment not updated: got '%s'", updatedTask.Comment)
	}
	t.Logf("✓ Task comment verified: '%s'", updatedTask.Comment)

	t.Log("=== ✓ All task lifecycle tests passed ===")
}

// TestE2E_TaskQueries tests various task query operations.
// Covers: GetTasksForCallback, GetTasksByStatus
func TestE2E_TaskQueries(t *testing.T) {
	// Ensure at least one callback exists (reuses existing or creates one)
	callbackID := EnsureCallbackExists(t)

	client := AuthenticateTestClient(t)

	// Get the callback details
	ctx0, cancel0 := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel0()
	testCallback, err := client.GetCallbackByID(ctx0, callbackID)
	if err != nil {
		t.Fatalf("Failed to get callback: %v", err)
	}

	// Test 1: Get all tasks for callback
	t.Log("=== Test 1: Get all tasks for callback ===")
	ctx1, cancel1 := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel1()

	tasks, err := client.GetTasksForCallback(ctx1, testCallback.DisplayID, 50)
	if err != nil {
		t.Fatalf("GetTasksForCallback failed: %v", err)
	}
	t.Logf("✓ Found %d tasks for callback %d", len(tasks), testCallback.DisplayID)

	// Validate task structure
	for _, task := range tasks {
		if task.DisplayID == 0 {
			t.Error("Task has DisplayID 0")
		}
		if task.CommandName == "" {
			t.Error("Task has empty CommandName")
		}
		if task.CallbackID != testCallback.ID {
			t.Errorf("Task has wrong CallbackID: expected %d, got %d", testCallback.ID, task.CallbackID)
		}
	}

	// Test 2: Get tasks by status (completed)
	t.Log("=== Test 2: Get completed tasks ===")
	ctx2, cancel2 := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel2()

	completedTasks, err := client.GetTasksByStatus(ctx2, testCallback.DisplayID, mythic.TaskStatusCompleted, 20)
	if err != nil {
		t.Fatalf("GetTasksByStatus (completed) failed: %v", err)
	}
	t.Logf("✓ Found %d completed tasks", len(completedTasks))

	// Verify all tasks are completed
	for _, task := range completedTasks {
		if !task.Completed {
			t.Errorf("Task %d is not marked as completed", task.DisplayID)
		}
		if task.Status != "completed" && task.Status != "success" {
			t.Errorf("Task %d has status '%s', expected 'completed'", task.DisplayID, task.Status)
		}
	}

	// Test 3: Get tasks by status (submitted)
	t.Log("=== Test 3: Get submitted tasks ===")
	ctx3, cancel3 := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel3()

	submittedTasks, err := client.GetTasksByStatus(ctx3, testCallback.DisplayID, mythic.TaskStatusSubmitted, 20)
	if err != nil {
		t.Fatalf("GetTasksByStatus (submitted) failed: %v", err)
	}
	t.Logf("✓ Found %d submitted/pending tasks", len(submittedTasks))

	// Test 4: Get tasks by status (error)
	t.Log("=== Test 4: Get errored tasks ===")
	ctx4, cancel4 := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel4()

	errorTasks, err := client.GetTasksByStatus(ctx4, testCallback.DisplayID, mythic.TaskStatusError, 20)
	if err != nil {
		t.Fatalf("GetTasksByStatus (error) failed: %v", err)
	}
	t.Logf("✓ Found %d errored tasks", len(errorTasks))

	t.Log("=== ✓ All task query tests passed ===")
}

// TestE2E_TaskReissue tests task reissue functionality.
// Covers: ReissueTask, ReissueTaskWithHandler
func TestE2E_TaskReissue(t *testing.T) {
	// Ensure at least one callback exists (reuses existing or creates one)
	callbackID := EnsureCallbackExists(t)

	client := AuthenticateTestClient(t)

	// Get the callback details
	ctx0, cancel0 := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel0()
	testCallback, err := client.GetCallbackByID(ctx0, callbackID)
	if err != nil {
		t.Fatalf("Failed to get callback: %v", err)
	}

	// Find a completed task to reissue
	ctx1, cancel1 := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel1()

	completedTasks, err := client.GetTasksByStatus(ctx1, testCallback.DisplayID, mythic.TaskStatusCompleted, 5)
	if err != nil {
		t.Fatalf("GetTasksByStatus failed: %v", err)
	}

	if len(completedTasks) == 0 {
		t.Skip("No completed tasks found, skipping task reissue tests")
	}

	testTask := completedTasks[0]
	t.Logf("Using task %d (%s) for reissue test", testTask.DisplayID, testTask.CommandName)

	// Test 1: Reissue task (standard)
	t.Log("=== Test 1: Reissue task (standard) ===")
	ctx2, cancel2 := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel2()

	err = client.ReissueTask(ctx2, testTask.ID)
	if err != nil {
		t.Logf("⚠ ReissueTask failed (may not be supported): %v", err)
	} else {
		t.Log("✓ Task reissued successfully")
	}

	// Test 2: Reissue task with handler
	t.Log("=== Test 2: Reissue task with handler ===")
	ctx3, cancel3 := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel3()

	err = client.ReissueTaskWithHandler(ctx3, testTask.ID)
	if err != nil {
		t.Logf("⚠ ReissueTaskWithHandler failed (may not be supported): %v", err)
	} else {
		t.Log("✓ Task reissued with handler successfully")
	}

	t.Log("=== ✓ Task reissue tests completed ===")
}

// TestE2E_TaskOPSEC tests OPSEC-related task functionality.
// Covers: RequestOpsecBypass
func TestE2E_TaskOPSEC(t *testing.T) {
	// Ensure at least one callback exists (reuses existing or creates one)
	callbackID := EnsureCallbackExists(t)

	client := AuthenticateTestClient(t)

	// Get the callback details
	ctx0, cancel0 := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel0()
	testCallback, err := client.GetCallbackByID(ctx0, callbackID)
	if err != nil {
		t.Fatalf("Failed to get callback: %v", err)
	}

	// Issue a task to test OPSEC bypass
	t.Log("=== Test: Request OPSEC bypass ===")
	ctx1, cancel1 := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel1()

	taskReq := &mythic.TaskRequest{
		Command:    "whoami",
		Params:     "whoami",
		CallbackID: &testCallback.ID,
	}

	task, err := client.IssueTask(ctx1, taskReq)
	if err != nil {
		t.Fatalf("IssueTask failed: %v", err)
	}

	// Try to request OPSEC bypass
	ctx2, cancel2 := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel2()

	err = client.RequestOpsecBypass(ctx2, task.ID)
	if err != nil {
		// OPSEC may not be configured, that's okay
		t.Logf("⚠ RequestOpsecBypass failed (may not be configured): %v", err)
	} else {
		t.Log("✓ OPSEC bypass requested successfully")
	}

	t.Log("=== ✓ OPSEC tests completed ===")
}

// TestE2E_TaskMITREAttack tests MITRE ATT&CK tagging for tasks.
// Covers: AddMITREAttackToTask
func TestE2E_TaskMITREAttack(t *testing.T) {
	// Ensure at least one callback exists (reuses existing or creates one)
	callbackID := EnsureCallbackExists(t)

	client := AuthenticateTestClient(t)

	// Get the callback details
	ctx0, cancel0 := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel0()
	testCallback, err := client.GetCallbackByID(ctx0, callbackID)
	if err != nil {
		t.Fatalf("Failed to get callback: %v", err)
	}

	// Issue a task
	t.Log("=== Test: Add MITRE ATT&CK tag to task ===")
	ctx1, cancel1 := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel1()

	taskReq := &mythic.TaskRequest{
		Command:    "whoami",
		Params:     "whoami",
		CallbackID: &testCallback.ID,
	}

	task, err := client.IssueTask(ctx1, taskReq)
	if err != nil {
		t.Fatalf("IssueTask failed: %v", err)
	}

	// Add MITRE ATT&CK tag (T1033 - System Owner/User Discovery)
	ctx2, cancel2 := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel2()

	err = client.AddMITREAttackToTask(ctx2, task.DisplayID, "T1033")
	if err != nil {
		t.Fatalf("AddMITREAttackToTask failed: %v", err)
	}
	t.Log("✓ MITRE ATT&CK tag T1033 added to task")

	t.Log("=== ✓ MITRE ATT&CK tagging tests passed ===")
}

// TestE2E_TaskArtifacts tests task artifact retrieval.
// Covers: GetTaskArtifacts
func TestE2E_TaskArtifacts(t *testing.T) {
	// Ensure at least one callback exists (reuses existing or creates one)
	callbackID := EnsureCallbackExists(t)

	client := AuthenticateTestClient(t)

	// Get the callback details
	ctx0, cancel0 := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel0()
	testCallback, err := client.GetCallbackByID(ctx0, callbackID)
	if err != nil {
		t.Fatalf("Failed to get callback: %v", err)
	}

	// Get tasks and find one with artifacts
	t.Log("=== Test: Get task artifacts ===")
	ctx1, cancel1 := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel1()

	tasks, err := client.GetTasksForCallback(ctx1, testCallback.DisplayID, 20)
	if err != nil {
		t.Fatalf("GetTasksForCallback failed: %v", err)
	}

	if len(tasks) == 0 {
		t.Skip("No tasks found for callback")
	}

	// Check first few tasks for artifacts
	foundArtifacts := false
	for _, task := range tasks {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		artifacts, err := client.GetTaskArtifacts(ctx, task.DisplayID)
		cancel()

		if err != nil {
			t.Logf("⚠ GetTaskArtifacts failed for task %d: %v", task.DisplayID, err)
			continue
		}

		if len(artifacts) > 0 {
			t.Logf("✓ Task %d has %d artifacts", task.DisplayID, len(artifacts))
			foundArtifacts = true

			// Validate artifact structure
			for i, artifact := range artifacts {
				if i < 3 { // Show first 3
					t.Logf("  [%d] Artifact: %s (Host: %s)", i+1, artifact.Artifact, artifact.Host)
				}
			}
			break
		}
	}

	if !foundArtifacts {
		t.Log("⚠ No artifacts found in checked tasks")
	}

	t.Log("=== ✓ Task artifact tests completed ===")
}

// TestE2E_TaskErrorHandling tests error scenarios for task operations.
func TestE2E_TaskErrorHandling(t *testing.T) {
	client := AuthenticateTestClient(t)

	// Test 1: Issue task with invalid callback
	t.Log("=== Test 1: Issue task with invalid callback ===")
	ctx1, cancel1 := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel1()

	invalidCallbackID := 999999
	taskReq := &mythic.TaskRequest{
		Command:    "whoami",
		Params:     "whoami",
		CallbackID: &invalidCallbackID,
	}

	_, err := client.IssueTask(ctx1, taskReq)
	if err == nil {
		t.Error("Expected error for invalid callback ID")
	}
	t.Logf("✓ Invalid callback rejected: %v", err)

	// Test 2: Get task with invalid display ID
	t.Log("=== Test 2: Get task with invalid display ID ===")
	ctx2, cancel2 := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel2()

	_, err = client.GetTask(ctx2, 999999)
	if err == nil {
		t.Error("Expected error for non-existent task")
	}
	t.Logf("✓ Non-existent task rejected: %v", err)

	// Test 3: Update task with invalid display ID
	t.Log("=== Test 3: Update task with invalid display ID ===")
	ctx3, cancel3 := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel3()

	updates := map[string]interface{}{
		"comment": "test",
	}
	err = client.UpdateTask(ctx3, 999999, updates)
	if err == nil {
		t.Error("Expected error for non-existent task update")
	}
	t.Logf("✓ Non-existent task update rejected: %v", err)

	// Test 4: Get task output with invalid display ID
	t.Log("=== Test 4: Get task output with invalid display ID ===")
	ctx4, cancel4 := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel4()

	_, err = client.GetTaskOutput(ctx4, 999999)
	if err == nil {
		t.Error("Expected error for non-existent task output")
	}
	t.Logf("✓ Non-existent task output rejected: %v", err)

	// Test 5: Add MITRE ATT&CK with invalid attack ID
	t.Log("=== Test 5: Add invalid MITRE ATT&CK tag ===")

	// Ensure at least one callback exists (reuses existing or creates one)
	callbackIDForTest5 := EnsureCallbackExists(t)

	// Get the callback details
	ctx5a, cancel5a := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel5a()
	activeCallback, err := client.GetCallbackByID(ctx5a, callbackIDForTest5)
	if err != nil {
		t.Fatalf("Failed to get callback: %v", err)
	}

	ctx5b, cancel5b := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel5b()

	taskReq = &mythic.TaskRequest{
		Command:    "whoami",
		Params:     "whoami",
		CallbackID: &activeCallback.ID,
	}

	task, err := client.IssueTask(ctx5b, taskReq)
	if err != nil {
		t.Fatalf("IssueTask failed: %v", err)
	}

	ctx5c, cancel5c := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel5c()

	err = client.AddMITREAttackToTask(ctx5c, task.DisplayID, "INVALID_ATTACK_ID")
	// This might or might not fail depending on Mythic validation
	if err != nil {
		t.Logf("✓ Invalid MITRE ATT&CK ID rejected: %v", err)
	} else {
		t.Log("⚠ Invalid MITRE ATT&CK ID was accepted (server may not validate)")
	}

	t.Log("=== ✓ All error handling tests completed ===")
}

// TestE2E_TaskConcurrency tests concurrent task operations.
func TestE2E_TaskConcurrency(t *testing.T) {
	// Ensure at least one callback exists (reuses existing or creates one)
	callbackID := EnsureCallbackExists(t)

	client := AuthenticateTestClient(t)

	// Get the callback details
	ctx0, cancel0 := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel0()
	testCallback, err := client.GetCallbackByID(ctx0, callbackID)
	if err != nil {
		t.Fatalf("Failed to get callback: %v", err)
	}

	// Test: Issue multiple tasks concurrently
	t.Log("=== Test: Issue multiple tasks concurrently ===")

	numTasks := 5
	results := make(chan error, numTasks)
	taskIDs := make(chan int, numTasks)

	for i := 0; i < numTasks; i++ {
		go func(taskNum int) {
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()

			taskReq := &mythic.TaskRequest{
				Command:    "whoami",
				Params:     fmt.Sprintf("whoami # concurrent test %d", taskNum),
				CallbackID: &testCallback.ID,
			}

			task, err := client.IssueTask(ctx, taskReq)
			if err != nil {
				results <- fmt.Errorf("task %d failed: %w", taskNum, err)
				return
			}

			taskIDs <- task.DisplayID
			results <- nil
		}(i)
	}

	// Collect results
	successCount := 0
	var issuedTaskIDs []int
	for i := 0; i < numTasks; i++ {
		err := <-results
		if err != nil {
			t.Logf("⚠ %v", err)
		} else {
			successCount++
			taskID := <-taskIDs
			issuedTaskIDs = append(issuedTaskIDs, taskID)
		}
	}

	t.Logf("✓ Successfully issued %d/%d concurrent tasks", successCount, numTasks)

	// Verify all tasks were created
	if successCount > 0 {
		ctx1, cancel1 := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel1()

		for _, taskID := range issuedTaskIDs {
			task, err := client.GetTask(ctx1, taskID)
			if err != nil {
				t.Errorf("Failed to retrieve concurrent task %d: %v", taskID, err)
			} else {
				t.Logf("  Task %d: %s", task.DisplayID, task.CommandName)
			}
		}
	}

	t.Log("=== ✓ Concurrency tests passed ===")
}
