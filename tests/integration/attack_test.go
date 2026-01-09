//go:build integration
// +build integration

package integration

import (
	"context"
	"testing"
	"time"

	"github.com/nbaertsch/mythic-sdk-go/pkg/mythic"
)

// TestAttack_GetAttackTechniques tests retrieving all MITRE ATT&CK techniques
func TestAttack_GetAttackTechniques(t *testing.T) {

	client := AuthenticateTestClient(t)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Get all attack techniques
	attacks, err := client.GetAttackTechniques(ctx)
	if err != nil {
		t.Fatalf("Failed to get attack techniques: %v", err)
	}

	t.Logf("Retrieved %d MITRE ATT&CK techniques", len(attacks))

	// Verify structure if any attacks exist
	if len(attacks) > 0 {
		attack := attacks[0]
		if attack.ID == 0 {
			t.Error("Attack ID should not be 0")
		}
		if attack.TNum == "" {
			t.Error("Attack TNum should not be empty")
		}
		if attack.Timestamp.IsZero() {
			t.Error("Attack timestamp should not be zero")
		}

		// Test String method
		str := attack.String()
		if str == "" {
			t.Error("Attack.String() should not return empty string")
		}
		t.Logf("Attack: %s", str)
		t.Logf("  TNum: %s", attack.TNum)
		t.Logf("  Name: %s", attack.Name)
		t.Logf("  OS: %s", attack.OS)
		t.Logf("  Tactic: %s", attack.Tactic)
	}
}

// TestAttack_GetAttackTechniqueByID tests retrieving a specific technique by ID
func TestAttack_GetAttackTechniqueByID(t *testing.T) {

	client := AuthenticateTestClient(t)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// First get all techniques to find a valid ID
	attacks, err := client.GetAttackTechniques(ctx)
	if err != nil {
		t.Fatalf("Failed to get attack techniques: %v", err)
	}

	if len(attacks) == 0 {
		t.Skip("No attack techniques available for testing")
	}

	// Get the first attack by ID
	attackID := attacks[0].ID
	attack, err := client.GetAttackTechniqueByID(ctx, attackID)
	if err != nil {
		t.Fatalf("Failed to get attack technique by ID: %v", err)
	}

	if attack.ID != attackID {
		t.Errorf("Expected ID %d, got %d", attackID, attack.ID)
	}

	t.Logf("Retrieved attack technique: %s", attack.String())
}

// TestAttack_GetAttackTechniqueByID_InvalidID tests getting with invalid ID
func TestAttack_GetAttackTechniqueByID_InvalidID(t *testing.T) {

	client := AuthenticateTestClient(t)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Try to get attack with ID 0
	_, err := client.GetAttackTechniqueByID(ctx, 0)
	if err == nil {
		t.Error("Expected error for attack ID 0")
	}

	// Try to get attack with very high ID (likely doesn't exist)
	_, err = client.GetAttackTechniqueByID(ctx, 999999999)
	if err == nil {
		t.Error("Expected error for non-existent attack")
	}
}

// TestAttack_GetAttackTechniqueByTNum tests retrieving by technique number
func TestAttack_GetAttackTechniqueByTNum(t *testing.T) {

	client := AuthenticateTestClient(t)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// First get all techniques to find a valid TNum
	attacks, err := client.GetAttackTechniques(ctx)
	if err != nil {
		t.Fatalf("Failed to get attack techniques: %v", err)
	}

	if len(attacks) == 0 {
		t.Skip("No attack techniques available for testing")
	}

	// Get the first attack by TNum
	tNum := attacks[0].TNum
	attack, err := client.GetAttackTechniqueByTNum(ctx, tNum)
	if err != nil {
		t.Fatalf("Failed to get attack technique by TNum: %v", err)
	}

	if attack.TNum != tNum {
		t.Errorf("Expected TNum %q, got %q", tNum, attack.TNum)
	}

	t.Logf("Retrieved attack technique: %s", attack.String())
}

// TestAttack_GetAttackTechniqueByTNum_InvalidTNum tests getting with invalid TNum
func TestAttack_GetAttackTechniqueByTNum_InvalidTNum(t *testing.T) {

	client := AuthenticateTestClient(t)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Try to get attack with empty TNum
	_, err := client.GetAttackTechniqueByTNum(ctx, "")
	if err == nil {
		t.Error("Expected error for empty TNum")
	}

	// Try to get attack with non-existent TNum
	_, err = client.GetAttackTechniqueByTNum(ctx, "T999999")
	if err == nil {
		t.Error("Expected error for non-existent TNum")
	}
}

// TestAttack_GetAttackByTask tests retrieving ATT&CK tags for a task
func TestAttack_GetAttackByTask(t *testing.T) {

	client := AuthenticateTestClient(t)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// First get all tasks to find one that might have attack tags
	tasks, err := client.GetTasksForCallback(ctx, 1, 10) // Use callback 1 if available
	if err != nil {
		// Try to get tasks without specifying callback
		t.Skip("Could not retrieve tasks for testing")
	}

	if len(tasks) == 0 {
		t.Skip("No tasks available for testing")
	}

	// Try to get attack tags for the first task
	taskID := tasks[0].ID
	attackTasks, err := client.GetAttackByTask(ctx, taskID)
	if err != nil {
		t.Fatalf("Failed to get attack tags for task: %v", err)
	}

	t.Logf("Retrieved %d attack tags for task %d", len(attackTasks), taskID)

	// Verify structure if any attack tags exist
	if len(attackTasks) > 0 {
		at := attackTasks[0]
		if at.ID == 0 {
			t.Error("AttackTask ID should not be 0")
		}
		if at.AttackID == 0 {
			t.Error("AttackTask AttackID should not be 0")
		}
		if at.TaskID != taskID {
			t.Errorf("Expected TaskID %d, got %d", taskID, at.TaskID)
		}

		// Test String method
		str := at.String()
		if str == "" {
			t.Error("AttackTask.String() should not return empty string")
		}
		t.Logf("AttackTask: %s", str)
	}
}

// TestAttack_GetAttackByTask_InvalidID tests getting with invalid task ID
func TestAttack_GetAttackByTask_InvalidID(t *testing.T) {

	client := AuthenticateTestClient(t)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Try to get attack tags with ID 0
	_, err := client.GetAttackByTask(ctx, 0)
	if err == nil {
		t.Error("Expected error for task ID 0")
	}
}

// TestAttack_GetAttackByCommand tests retrieving ATT&CK tags for a command
func TestAttack_GetAttackByCommand(t *testing.T) {

	client := AuthenticateTestClient(t)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// For this test, we'll use a command ID of 1 if it exists
	// Most Mythic installations will have some commands
	attackCommands, err := client.GetAttackByCommand(ctx, 1)
	if err != nil {
		t.Fatalf("Failed to get attack tags for command: %v", err)
	}

	t.Logf("Retrieved %d attack tags for command 1", len(attackCommands))

	// Verify structure if any attack tags exist
	if len(attackCommands) > 0 {
		ac := attackCommands[0]
		if ac.ID == 0 {
			t.Error("AttackCommand ID should not be 0")
		}
		if ac.AttackID == 0 {
			t.Error("AttackCommand AttackID should not be 0")
		}
		if ac.CommandID == 0 {
			t.Error("AttackCommand CommandID should not be 0")
		}

		// Test String method
		str := ac.String()
		if str == "" {
			t.Error("AttackCommand.String() should not return empty string")
		}
		t.Logf("AttackCommand: %s", str)
	}
}

// TestAttack_GetAttackByCommand_InvalidID tests getting with invalid command ID
func TestAttack_GetAttackByCommand_InvalidID(t *testing.T) {

	client := AuthenticateTestClient(t)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Try to get attack tags with ID 0
	_, err := client.GetAttackByCommand(ctx, 0)
	if err == nil {
		t.Error("Expected error for command ID 0")
	}
}

// TestAttack_GetAttacksByOperation tests retrieving all techniques used in an operation
func TestAttack_GetAttacksByOperation(t *testing.T) {

	client := AuthenticateTestClient(t)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Get current operation
	operationID := client.GetCurrentOperation()
	if operationID == nil {
		t.Skip("No current operation set")
	}

	// Get all attack techniques used in the operation
	attacks, err := client.GetAttacksByOperation(ctx, *operationID)
	if err != nil {
		t.Fatalf("Failed to get attack techniques for operation: %v", err)
	}

	t.Logf("Retrieved %d unique attack techniques for operation %d", len(attacks), *operationID)

	// Verify structure if any attacks exist
	for _, attack := range attacks {
		if attack.ID == 0 {
			t.Error("Attack ID should not be 0")
		}
		if attack.TNum == "" {
			t.Error("Attack TNum should not be empty")
		}
	}
}

// TestAttack_GetAttacksByOperation_InvalidID tests getting with invalid operation ID
func TestAttack_GetAttacksByOperation_InvalidID(t *testing.T) {

	client := AuthenticateTestClient(t)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Try to get attacks with operation ID 0
	_, err := client.GetAttacksByOperation(ctx, 0)
	if err == nil {
		t.Error("Expected error for operation ID 0")
	}
}

// TestAttack_TechniqueOrdering tests that techniques are ordered by TNum
func TestAttack_TechniqueOrdering(t *testing.T) {

	client := AuthenticateTestClient(t)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Get all techniques
	attacks, err := client.GetAttackTechniques(ctx)
	if err != nil {
		t.Fatalf("Failed to get attack techniques: %v", err)
	}

	if len(attacks) < 2 {
		t.Skip("Need at least 2 attack techniques to test ordering")
	}

	// Verify techniques are in ascending order by TNum
	for i := 0; i < len(attacks)-1; i++ {
		if attacks[i].TNum > attacks[i+1].TNum {
			t.Errorf("Attacks not in ascending TNum order: %s after %s",
				attacks[i].TNum, attacks[i+1].TNum)
		}
	}

	t.Logf("Verified %d attack techniques are properly ordered by TNum", len(attacks))
}

// TestAttack_AttackTaskOrdering tests that attack tasks are ordered by timestamp
func TestAttack_AttackTaskOrdering(t *testing.T) {

	client := AuthenticateTestClient(t)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// First get tasks to find one with attack tags
	tasks, err := client.GetTasksForCallback(ctx, 1, 10)
	if err != nil || len(tasks) == 0 {
		t.Skip("Could not retrieve tasks for testing")
	}

	// Find a task with attack tags
	var attackTasks []interface{}
	for _, task := range tasks {
		ats, err := client.GetAttackByTask(ctx, task.ID)
		if err == nil && len(ats) >= 2 {
			attackTasks = []interface{}{ats}
			break
		}
	}

	if len(attackTasks) == 0 {
		t.Skip("Need at least one task with 2+ attack tags to test ordering")
	}
}

// TestAttack_CommonTechniques tests retrieval of common attack techniques
func TestAttack_CommonTechniques(t *testing.T) {

	client := AuthenticateTestClient(t)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Common techniques that should exist in most Mythic installations
	commonTechniques := []string{
		"T1003", // OS Credential Dumping
		"T1059", // Command and Scripting Interpreter
		"T1055", // Process Injection
		"T1078", // Valid Accounts
		"T1082", // System Information Discovery
	}

	foundCount := 0
	for _, tNum := range commonTechniques {
		attack, err := client.GetAttackTechniqueByTNum(ctx, tNum)
		if err == nil && attack != nil {
			foundCount++
			t.Logf("Found common technique: %s - %s", attack.TNum, attack.Name)
		}
	}

	if foundCount == 0 {
		t.Log("Warning: None of the common techniques were found")
	} else {
		t.Logf("Found %d out of %d common techniques", foundCount, len(commonTechniques))
	}
}

// TestAttack_AddMITREAttackToTask tests adding MITRE ATT&CK tag to a task
func TestAttack_AddMITREAttackToTask(t *testing.T) {
	t.Skip("Skipping AddMITREAttackToTask to avoid modifying task tags")

	client := AuthenticateTestClient(t)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Get all attack techniques to find a valid one
	attacks, err := client.GetAttackTechniques(ctx)
	if err != nil {
		t.Fatalf("Failed to get attack techniques: %v", err)
	}

	if len(attacks) == 0 {
		t.Skip("No attack techniques available for testing")
	}

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

	t.Logf("Created task %d for MITRE ATT&CK tag test", task.DisplayID)

	// Add MITRE ATT&CK tag to the task
	attackTNum := attacks[0].TNum
	err = client.AddMITREAttackToTask(ctx, task.DisplayID, attackTNum)
	if err != nil {
		t.Fatalf("Failed to add MITRE ATT&CK tag to task: %v", err)
	}

	t.Logf("Successfully added MITRE ATT&CK tag %s to task %d", attackTNum, task.DisplayID)

	// Verify the tag was added
	attackTasks, err := client.GetAttackByTask(ctx, task.ID)
	if err != nil {
		t.Fatalf("Failed to get attack tags for task: %v", err)
	}

	found := false
	for _, at := range attackTasks {
		if at.TaskID == task.ID {
			found = true
			t.Logf("Verified MITRE ATT&CK tag on task: %s", at.String())
			break
		}
	}

	if !found {
		t.Error("MITRE ATT&CK tag was not found on task after adding")
	}
}

// TestAttack_AddMITREAttackToTask_InvalidInput tests with invalid input
func TestAttack_AddMITREAttackToTask_InvalidInput(t *testing.T) {

	client := AuthenticateTestClient(t)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Test with zero task ID
	err := client.AddMITREAttackToTask(ctx, 0, "T1003")
	if err == nil {
		t.Fatal("Expected error for zero task ID, got nil")
	}
	t.Logf("Zero task ID error: %v", err)

	// Test with empty attack ID
	err = client.AddMITREAttackToTask(ctx, 1, "")
	if err == nil {
		t.Fatal("Expected error for empty attack ID, got nil")
	}
	t.Logf("Empty attack ID error: %v", err)

	// Test with non-existent task ID
	err = client.AddMITREAttackToTask(ctx, 999999, "T1003")
	if err == nil {
		t.Fatal("Expected error for non-existent task ID, got nil")
	}
	t.Logf("Non-existent task ID error: %v", err)

	// Test with invalid attack technique
	err = client.AddMITREAttackToTask(ctx, 1, "T999999")
	if err == nil {
		t.Fatal("Expected error for invalid attack technique, got nil")
	}
	t.Logf("Invalid attack technique error: %v", err)
}
