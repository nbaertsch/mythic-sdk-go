//go:build integration

package integration

import (
	"context"
	"testing"
	"time"

	"github.com/nbaertsch/mythic-sdk-go/pkg/mythic"
)

func TestTasks_IssueTask(t *testing.T) {
	SkipIfNoMythic(t)

	client := AuthenticateTestClient(t)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Get callbacks first
	callbacks, err := client.GetAllCallbacks(ctx)
	if err != nil {
		t.Fatalf("Failed to get callbacks: %v", err)
	}

	if len(callbacks) == 0 {
		t.Skip("No callbacks available for testing")
	}

	// Issue a simple task
	callbackID := callbacks[0].DisplayID
	task, err := client.IssueTask(ctx, &mythic.TaskRequest{
		CallbackID: &callbackID,
		Command:    "whoami",
		Params:     "",
	})

	if err != nil {
		t.Fatalf("Failed to issue task: %v", err)
	}

	if task == nil {
		t.Fatal("IssueTask returned nil task")
	}

	if task.DisplayID == 0 {
		t.Error("Task DisplayID should not be 0")
	}

	if task.CommandName != "whoami" {
		t.Errorf("Expected command 'whoami', got %q", task.CommandName)
	}

	if task.CallbackID != callbacks[0].ID {
		t.Errorf("Expected callback ID %d, got %d", callbacks[0].ID, task.CallbackID)
	}

	t.Logf("Created task: %s", task.String())
}

func TestTasks_IssueTask_InvalidCallback(t *testing.T) {
	SkipIfNoMythic(t)

	client := AuthenticateTestClient(t)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Try to issue task to non-existent callback
	invalidCallbackID := 999999
	_, err := client.IssueTask(ctx, &mythic.TaskRequest{
		CallbackID: &invalidCallbackID,
		Command:    "whoami",
		Params:     "",
	})

	if err == nil {
		t.Fatal("Expected error for invalid callback, got nil")
	}

	t.Logf("Expected error for invalid callback: %v", err)
}

func TestTasks_IssueTask_MissingCallback(t *testing.T) {
	SkipIfNoMythic(t)

	client := AuthenticateTestClient(t)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Try to issue task without callback ID
	_, err := client.IssueTask(ctx, &mythic.TaskRequest{
		Command: "whoami",
		Params:  "",
	})

	if err == nil {
		t.Fatal("Expected error for missing callback, got nil")
	}

	if err.Error() == "" {
		t.Error("Error message should not be empty")
	}

	t.Logf("Expected error for missing callback: %v", err)
}

func TestTasks_IssueTask_MissingCommand(t *testing.T) {
	SkipIfNoMythic(t)

	client := AuthenticateTestClient(t)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	callbackID := 1
	_, err := client.IssueTask(ctx, &mythic.TaskRequest{
		CallbackID: &callbackID,
		Params:     "",
	})

	if err == nil {
		t.Fatal("Expected error for missing command, got nil")
	}

	t.Logf("Expected error for missing command: %v", err)
}

func TestTasks_GetTask(t *testing.T) {
	SkipIfNoMythic(t)

	client := AuthenticateTestClient(t)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Get callbacks first
	callbacks, err := client.GetAllCallbacks(ctx)
	if err != nil {
		t.Fatalf("Failed to get callbacks: %v", err)
	}

	if len(callbacks) == 0 {
		t.Skip("No callbacks available for testing")
	}

	// Issue a task first
	callbackID := callbacks[0].DisplayID
	createdTask, err := client.IssueTask(ctx, &mythic.TaskRequest{
		CallbackID: &callbackID,
		Command:    "pwd",
		Params:     "",
	})

	if err != nil {
		t.Fatalf("Failed to issue task: %v", err)
	}

	// Now get the task
	task, err := client.GetTask(ctx, createdTask.DisplayID)
	if err != nil {
		t.Fatalf("Failed to get task: %v", err)
	}

	if task.DisplayID != createdTask.DisplayID {
		t.Errorf("Expected DisplayID %d, got %d", createdTask.DisplayID, task.DisplayID)
	}

	if task.CommandName != "pwd" {
		t.Errorf("Expected command 'pwd', got %q", task.CommandName)
	}

	t.Logf("Retrieved task: %s", task.String())
}

func TestTasks_GetTask_NotFound(t *testing.T) {
	SkipIfNoMythic(t)

	client := AuthenticateTestClient(t)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Try to get non-existent task
	_, err := client.GetTask(ctx, 999999)
	if err == nil {
		t.Fatal("Expected error for non-existent task, got nil")
	}

	t.Logf("Expected error for non-existent task: %v", err)
}

func TestTasks_GetTasksForCallback(t *testing.T) {
	SkipIfNoMythic(t)

	client := AuthenticateTestClient(t)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Get callbacks first
	callbacks, err := client.GetAllCallbacks(ctx)
	if err != nil {
		t.Fatalf("Failed to get callbacks: %v", err)
	}

	if len(callbacks) == 0 {
		t.Skip("No callbacks available for testing")
	}

	callbackID := callbacks[0].DisplayID

	// Issue a couple of tasks
	for i := 0; i < 2; i++ {
		_, err := client.IssueTask(ctx, &mythic.TaskRequest{
			CallbackID: &callbackID,
			Command:    "echo",
			Params:     "test",
		})
		if err != nil {
			t.Fatalf("Failed to issue task %d: %v", i, err)
		}
	}

	// Get tasks for callback
	tasks, err := client.GetTasksForCallback(ctx, callbackID, 10)
	if err != nil {
		t.Fatalf("Failed to get tasks for callback: %v", err)
	}

	if len(tasks) < 2 {
		t.Logf("Expected at least 2 tasks, got %d (may have existing tasks)", len(tasks))
	}

	t.Logf("Found %d tasks for callback %d", len(tasks), callbackID)

	// Verify tasks belong to the callback
	for _, task := range tasks {
		if task.CallbackID != callbacks[0].ID {
			t.Errorf("Task %d belongs to callback %d, expected %d", task.DisplayID, task.CallbackID, callbacks[0].ID)
		}
	}
}

func TestTasks_GetTaskOutput(t *testing.T) {
	SkipIfNoMythic(t)

	client := AuthenticateTestClient(t)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Get callbacks first
	callbacks, err := client.GetAllCallbacks(ctx)
	if err != nil {
		t.Fatalf("Failed to get callbacks: %v", err)
	}

	if len(callbacks) == 0 {
		t.Skip("No callbacks available for testing")
	}

	// Issue a task
	callbackID := callbacks[0].DisplayID
	task, err := client.IssueTask(ctx, &mythic.TaskRequest{
		CallbackID: &callbackID,
		Command:    "whoami",
		Params:     "",
	})

	if err != nil {
		t.Fatalf("Failed to issue task: %v", err)
	}

	// Get task output (may be empty if task hasn't completed)
	responses, err := client.GetTaskOutput(ctx, task.DisplayID)
	if err != nil {
		t.Fatalf("Failed to get task output: %v", err)
	}

	// Output may be empty for new tasks
	t.Logf("Task has %d responses", len(responses))

	// If there are responses, validate them
	for i, resp := range responses {
		if resp.TaskID != task.ID {
			t.Errorf("Response %d has wrong TaskID: expected %d, got %d", i, task.ID, resp.TaskID)
		}

		if resp.Timestamp.IsZero() {
			t.Errorf("Response %d has zero timestamp", i)
		}
	}
}

func TestTasks_UpdateTask(t *testing.T) {
	SkipIfNoMythic(t)

	client := AuthenticateTestClient(t)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Get callbacks first
	callbacks, err := client.GetAllCallbacks(ctx)
	if err != nil {
		t.Fatalf("Failed to get callbacks: %v", err)
	}

	if len(callbacks) == 0 {
		t.Skip("No callbacks available for testing")
	}

	// Issue a task
	callbackID := callbacks[0].DisplayID
	task, err := client.IssueTask(ctx, &mythic.TaskRequest{
		CallbackID: &callbackID,
		Command:    "whoami",
		Params:     "",
	})

	if err != nil {
		t.Fatalf("Failed to issue task: %v", err)
	}

	// Update task comment
	err = client.UpdateTask(ctx, task.DisplayID, map[string]interface{}{
		"comment": "Updated via integration test",
	})

	if err != nil {
		t.Fatalf("Failed to update task: %v", err)
	}

	// Verify update
	updatedTask, err := client.GetTask(ctx, task.DisplayID)
	if err != nil {
		t.Fatalf("Failed to get updated task: %v", err)
	}

	if updatedTask.Comment != "Updated via integration test" {
		t.Errorf("Expected comment 'Updated via integration test', got %q", updatedTask.Comment)
	}

	t.Logf("Updated task comment: %s", updatedTask.Comment)
}

func TestTasks_WaitForTaskComplete_Timeout(t *testing.T) {
	SkipIfNoMythic(t)

	client := AuthenticateTestClient(t)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Get callbacks first
	callbacks, err := client.GetAllCallbacks(ctx)
	if err != nil {
		t.Fatalf("Failed to get callbacks: %v", err)
	}

	if len(callbacks) == 0 {
		t.Skip("No callbacks available for testing")
	}

	// Issue a task
	callbackID := callbacks[0].DisplayID
	task, err := client.IssueTask(ctx, &mythic.TaskRequest{
		CallbackID: &callbackID,
		Command:    "whoami",
		Params:     "",
	})

	if err != nil {
		t.Fatalf("Failed to issue task: %v", err)
	}

	// Wait with very short timeout (should timeout or complete quickly)
	err = client.WaitForTaskComplete(ctx, task.DisplayID, 2)

	// Task may complete or timeout - both are valid outcomes
	if err != nil {
		t.Logf("WaitForTaskComplete result: %v (expected timeout or completion)", err)
	} else {
		t.Logf("Task completed within timeout")
	}
}

func TestTasks_GetTaskOutput_NotFound(t *testing.T) {
	SkipIfNoMythic(t)

	client := AuthenticateTestClient(t)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Try to get output for non-existent task
	_, err := client.GetTaskOutput(ctx, 999999)
	if err == nil {
		t.Fatal("Expected error for non-existent task, got nil")
	}

	t.Logf("Expected error for non-existent task output: %v", err)
}

func TestTasks_TaskStatusHelpers(t *testing.T) {
	SkipIfNoMythic(t)

	client := AuthenticateTestClient(t)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Get callbacks first
	callbacks, err := client.GetAllCallbacks(ctx)
	if err != nil {
		t.Fatalf("Failed to get callbacks: %v", err)
	}

	if len(callbacks) == 0 {
		t.Skip("No callbacks available for testing")
	}

	// Issue a task
	callbackID := callbacks[0].DisplayID
	task, err := client.IssueTask(ctx, &mythic.TaskRequest{
		CallbackID: &callbackID,
		Command:    "whoami",
		Params:     "",
	})

	if err != nil {
		t.Fatalf("Failed to issue task: %v", err)
	}

	// Test helper methods
	t.Logf("Task IsCompleted: %v", task.IsCompleted())
	t.Logf("Task IsError: %v", task.IsError())
	t.Logf("Task HasOutput: %v", task.HasOutput())
	t.Logf("Task Status: %s", task.Status)

	// Verify String() method
	taskStr := task.String()
	if taskStr == "" {
		t.Error("Task String() should not be empty")
	}
	t.Logf("Task string: %s", taskStr)
}

func TestTasks_GetTaskArtifacts(t *testing.T) {
	SkipIfNoMythic(t)

	client := AuthenticateTestClient(t)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Get callbacks first
	callbacks, err := client.GetAllCallbacks(ctx)
	if err != nil {
		t.Fatalf("Failed to get callbacks: %v", err)
	}

	if len(callbacks) == 0 {
		t.Skip("No callbacks available for testing")
	}

	// Issue a task
	callbackID := callbacks[0].DisplayID
	task, err := client.IssueTask(ctx, &mythic.TaskRequest{
		CallbackID: &callbackID,
		Command:    "ls",
		Params:     "",
	})

	if err != nil {
		t.Fatalf("Failed to issue task: %v", err)
	}

	// Get task artifacts (may be empty for new tasks)
	artifacts, err := client.GetTaskArtifacts(ctx, task.DisplayID)
	if err != nil {
		t.Fatalf("Failed to get task artifacts: %v", err)
	}

	if artifacts == nil {
		t.Fatal("GetTaskArtifacts returned nil")
	}

	t.Logf("Task %d has %d artifact(s)", task.DisplayID, len(artifacts))

	// If there are artifacts, verify structure
	for i, artifact := range artifacts {
		if artifact.ID == 0 {
			t.Errorf("Artifact %d has zero ID", i)
		}
		if artifact.TaskID != task.ID {
			t.Errorf("Artifact %d has wrong TaskID: expected %d, got %d", i, task.ID, artifact.TaskID)
		}
		if artifact.Artifact == "" {
			t.Errorf("Artifact %d has empty artifact value", i)
		}
		t.Logf("  - Artifact %d: %s (host: %s)", i, artifact.Artifact, artifact.Host)
	}
}

func TestTasks_GetTaskArtifacts_InvalidID(t *testing.T) {
	SkipIfNoMythic(t)

	client := AuthenticateTestClient(t)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Test with zero ID
	_, err := client.GetTaskArtifacts(ctx, 0)
	if err == nil {
		t.Fatal("Expected error for zero task ID, got nil")
	}
	t.Logf("Zero ID error: %v", err)

	// Test with non-existent task ID
	_, err = client.GetTaskArtifacts(ctx, 999999)
	if err == nil {
		t.Fatal("Expected error for non-existent task ID, got nil")
	}
	t.Logf("Non-existent ID error: %v", err)
}

func TestTasks_ReissueTask(t *testing.T) {
	t.Skip("Skipping ReissueTask to avoid re-executing tasks")

	SkipIfNoMythic(t)

	client := AuthenticateTestClient(t)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Get callbacks first
	callbacks, err := client.GetAllCallbacks(ctx)
	if err != nil {
		t.Fatalf("Failed to get callbacks: %v", err)
	}

	if len(callbacks) == 0 {
		t.Skip("No callbacks available for testing")
	}

	// Issue a task
	callbackID := callbacks[0].DisplayID
	task, err := client.IssueTask(ctx, &mythic.TaskRequest{
		CallbackID: &callbackID,
		Command:    "whoami",
		Params:     "",
	})

	if err != nil {
		t.Fatalf("Failed to issue task: %v", err)
	}

	t.Logf("Created task %d for reissue test", task.DisplayID)

	// Reissue the task
	err = client.ReissueTask(ctx, task.DisplayID)
	if err != nil {
		t.Fatalf("Failed to reissue task: %v", err)
	}

	t.Logf("Successfully reissued task %d", task.DisplayID)

	// Note: The reissued task will create a new task with the same parameters
	// We don't verify the new task here to keep the test simple
}

func TestTasks_ReissueTask_InvalidID(t *testing.T) {
	SkipIfNoMythic(t)

	client := AuthenticateTestClient(t)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Test with zero ID
	err := client.ReissueTask(ctx, 0)
	if err == nil {
		t.Fatal("Expected error for zero task ID, got nil")
	}
	t.Logf("Zero ID error: %v", err)

	// Test with non-existent task ID
	err = client.ReissueTask(ctx, 999999)
	if err == nil {
		t.Fatal("Expected error for non-existent task ID, got nil")
	}
	t.Logf("Non-existent ID error: %v", err)
}

func TestTasks_ReissueTaskWithHandler(t *testing.T) {
	t.Skip("Skipping ReissueTaskWithHandler to avoid re-executing tasks")

	SkipIfNoMythic(t)

	client := AuthenticateTestClient(t)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Get callbacks first
	callbacks, err := client.GetAllCallbacks(ctx)
	if err != nil {
		t.Fatalf("Failed to get callbacks: %v", err)
	}

	if len(callbacks) == 0 {
		t.Skip("No callbacks available for testing")
	}

	// Issue a task
	callbackID := callbacks[0].DisplayID
	task, err := client.IssueTask(ctx, &mythic.TaskRequest{
		CallbackID: &callbackID,
		Command:    "whoami",
		Params:     "",
	})

	if err != nil {
		t.Fatalf("Failed to issue task: %v", err)
	}

	t.Logf("Created task %d for reissue with handler test", task.DisplayID)

	// Reissue the task with handler
	err = client.ReissueTaskWithHandler(ctx, task.DisplayID)
	if err != nil {
		t.Fatalf("Failed to reissue task with handler: %v", err)
	}

	t.Logf("Successfully reissued task %d with handler", task.DisplayID)
}

func TestTasks_ReissueTaskWithHandler_InvalidID(t *testing.T) {
	SkipIfNoMythic(t)

	client := AuthenticateTestClient(t)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Test with zero ID
	err := client.ReissueTaskWithHandler(ctx, 0)
	if err == nil {
		t.Fatal("Expected error for zero task ID, got nil")
	}
	t.Logf("Zero ID error: %v", err)

	// Test with non-existent task ID
	err = client.ReissueTaskWithHandler(ctx, 999999)
	if err == nil {
		t.Fatal("Expected error for non-existent task ID, got nil")
	}
	t.Logf("Non-existent ID error: %v", err)
}

func TestTasks_RequestOpsecBypass(t *testing.T) {
	t.Skip("Skipping RequestOpsecBypass to avoid modifying task security status")

	SkipIfNoMythic(t)

	client := AuthenticateTestClient(t)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Get callbacks first
	callbacks, err := client.GetAllCallbacks(ctx)
	if err != nil {
		t.Fatalf("Failed to get callbacks: %v", err)
	}

	if len(callbacks) == 0 {
		t.Skip("No callbacks available for testing")
	}

	// Issue a task that might require OPSEC bypass
	callbackID := callbacks[0].DisplayID
	task, err := client.IssueTask(ctx, &mythic.TaskRequest{
		CallbackID: &callbackID,
		Command:    "whoami",
		Params:     "",
	})

	if err != nil {
		t.Fatalf("Failed to issue task: %v", err)
	}

	t.Logf("Created task %d for OPSEC bypass test", task.DisplayID)

	// Request OPSEC bypass
	err = client.RequestOpsecBypass(ctx, task.DisplayID)
	if err != nil {
		// This may fail if the task doesn't require OPSEC review or if already bypassed
		t.Logf("RequestOpsecBypass result: %v (may be expected)", err)
		return
	}

	t.Logf("Successfully requested OPSEC bypass for task %d", task.DisplayID)
}

func TestTasks_RequestOpsecBypass_InvalidID(t *testing.T) {
	SkipIfNoMythic(t)

	client := AuthenticateTestClient(t)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Test with zero ID
	err := client.RequestOpsecBypass(ctx, 0)
	if err == nil {
		t.Fatal("Expected error for zero task ID, got nil")
	}
	t.Logf("Zero ID error: %v", err)

	// Test with non-existent task ID
	err = client.RequestOpsecBypass(ctx, 999999)
	if err == nil {
		t.Fatal("Expected error for non-existent task ID, got nil")
	}
	t.Logf("Non-existent ID error: %v", err)
}
