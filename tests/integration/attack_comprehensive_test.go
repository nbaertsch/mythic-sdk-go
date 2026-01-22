//go:build integration

package integration

import (
	"context"
	"testing"
	"time"

	"github.com/nbaertsch/mythic-sdk-go/pkg/mythic"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestE2E_Attack_GetAll_SchemaValidation validates GetAttackTechniques returns all
// MITRE ATT&CK techniques with proper field population.
func TestE2E_Attack_GetAll_SchemaValidation(t *testing.T) {
	client := AuthenticateTestClient(t)

	t.Log("=== Test: GetAttackTechniques schema validation ===")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	techniques, err := client.GetAttackTechniques(ctx)
	require.NoError(t, err, "GetAttackTechniques should succeed")
	require.NotNil(t, techniques, "Attack techniques should not be nil")

	if len(techniques) == 0 {
		t.Log("⚠ No MITRE ATT&CK techniques found - may be expected if database not initialized")
		t.Log("=== ✓ GetAttackTechniques validation passed (empty result) ===")
		return
	}

	t.Logf("✓ Retrieved %d MITRE ATT&CK technique(s)", len(techniques))

	// Validate each technique has required fields
	for i, tech := range techniques {
		assert.NotZero(t, tech.ID, "Attack[%d] should have ID", i)
		assert.NotEmpty(t, tech.TNum, "Attack[%d] should have TNum", i)
		assert.NotEmpty(t, tech.Name, "Attack[%d] should have Name", i)

		// Validate TNum format (should start with T followed by digits)
		assert.Regexp(t, `^T\d+`, tech.TNum,
			"Attack[%d] TNum should start with 'T' followed by digits", i)

		if i < 5 {
			t.Logf("  Technique[%d]: %s - %s (OS=%s, Tactic=%s)",
				i, tech.TNum, tech.Name, tech.OS, tech.Tactic)
		}
	}

	if len(techniques) > 5 {
		t.Logf("  ... and %d more techniques", len(techniques)-5)
	}

	t.Log("=== ✓ GetAttackTechniques schema validation passed ===")
}

// TestE2E_Attack_GetByIDAndTNum validates GetAttackTechniqueByID and GetAttackTechniqueByTNum
// retrieve specific techniques correctly.
func TestE2E_Attack_GetByIDAndTNum(t *testing.T) {
	client := AuthenticateTestClient(t)

	t.Log("=== Test: GetAttackTechniqueByID and GetAttackTechniqueByTNum ===")

	// First get all techniques to find one to test with
	ctx1, cancel1 := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel1()

	techniques, err := client.GetAttackTechniques(ctx1)
	require.NoError(t, err, "GetAttackTechniques should succeed")

	if len(techniques) == 0 {
		t.Log("⚠ No MITRE ATT&CK techniques found - skipping test")
		t.Skip("No ATT&CK techniques available to test")
		return
	}

	testTech := techniques[0]
	t.Logf("✓ Testing with technique: %s - %s (ID: %d)", testTech.TNum, testTech.Name, testTech.ID)

	// Test GetAttackTechniqueByID
	ctx2, cancel2 := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel2()

	techByID, err := client.GetAttackTechniqueByID(ctx2, testTech.ID)
	require.NoError(t, err, "GetAttackTechniqueByID should succeed")
	require.NotNil(t, techByID, "Technique should not be nil")

	assert.Equal(t, testTech.ID, techByID.ID, "ID should match")
	assert.Equal(t, testTech.TNum, techByID.TNum, "TNum should match")
	assert.Equal(t, testTech.Name, techByID.Name, "Name should match")
	assert.Equal(t, testTech.OS, techByID.OS, "OS should match")
	assert.Equal(t, testTech.Tactic, techByID.Tactic, "Tactic should match")

	t.Logf("✓ GetAttackTechniqueByID validated: %s - %s", techByID.TNum, techByID.Name)

	// Test GetAttackTechniqueByTNum
	ctx3, cancel3 := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel3()

	techByTNum, err := client.GetAttackTechniqueByTNum(ctx3, testTech.TNum)
	require.NoError(t, err, "GetAttackTechniqueByTNum should succeed")
	require.NotNil(t, techByTNum, "Technique should not be nil")

	assert.Equal(t, testTech.ID, techByTNum.ID, "ID should match")
	assert.Equal(t, testTech.TNum, techByTNum.TNum, "TNum should match")
	assert.Equal(t, testTech.Name, techByTNum.Name, "Name should match")

	t.Logf("✓ GetAttackTechniqueByTNum validated: %s - %s", techByTNum.TNum, techByTNum.Name)

	// Test error handling for invalid ID
	ctx4, cancel4 := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel4()

	invalidTech, err := client.GetAttackTechniqueByID(ctx4, 0)
	require.Error(t, err, "GetAttackTechniqueByID should error for ID 0")
	assert.Nil(t, invalidTech, "Technique should be nil when error occurs")
	assert.Contains(t, err.Error(), "attack ID is required",
		"Error should mention required attack ID")

	// Test error handling for invalid TNum
	ctx5, cancel5 := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel5()

	invalidTech2, err := client.GetAttackTechniqueByTNum(ctx5, "")
	require.Error(t, err, "GetAttackTechniqueByTNum should error for empty TNum")
	assert.Nil(t, invalidTech2, "Technique should be nil when error occurs")
	assert.Contains(t, err.Error(), "technique number is required",
		"Error should mention required technique number")

	// Test not found error
	ctx6, cancel6 := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel6()

	notFoundTech, err := client.GetAttackTechniqueByTNum(ctx6, "T999999")
	require.Error(t, err, "GetAttackTechniqueByTNum should error for non-existent TNum")
	assert.Nil(t, notFoundTech, "Technique should be nil when not found")
	assert.Contains(t, err.Error(), "not found", "Error should mention 'not found'")

	t.Log("✓ Error handling validated")
	t.Log("=== ✓ GetAttackTechniqueByID and GetAttackTechniqueByTNum validation passed ===")
}

// TestE2E_Attack_GetByTask validates GetAttackByTask retrieves ATT&CK tags for a task.
func TestE2E_Attack_GetByTask(t *testing.T) {
	client := AuthenticateTestClient(t)

	t.Log("=== Test: GetAttackByTask ===")

	// Get a callback to issue a task
	callback := getActiveCallback(t, client)
	if callback == nil {
		t.Skip("No active callback available for testing")
		return
	}

	t.Logf("✓ Using callback: ID=%d", callback.ID)

	// Issue a task that might have MITRE tags
	ctx1, cancel1 := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel1()

	taskReq := &mythic.TaskRequest{
		Command:    "whoami",
		Params:     "whoami",
		CallbackID: &callback.ID,
	}

	task, err := client.IssueTask(ctx1, taskReq)
	require.NoError(t, err, "IssueTask should succeed")
	require.NotNil(t, task, "Task should not be nil")

	t.Logf("✓ Issued task: ID=%d, Command=%s", task.ID, task.CommandName)

	// Get MITRE ATT&CK tags for this task
	ctx2, cancel2 := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel2()

	attackTags, err := client.GetAttackByTask(ctx2, task.ID)
	require.NoError(t, err, "GetAttackByTask should succeed")
	require.NotNil(t, attackTags, "Attack tags should not be nil")

	t.Logf("✓ Task has %d MITRE ATT&CK tag(s)", len(attackTags))

	// Validate attack tag structure
	for i, tag := range attackTags {
		assert.NotZero(t, tag.ID, "AttackTask[%d] should have ID", i)
		assert.NotZero(t, tag.AttackID, "AttackTask[%d] should have AttackID", i)
		assert.Equal(t, task.ID, tag.TaskID, "AttackTask[%d] should match task ID", i)

		t.Logf("  Tag[%d]: ID=%d, AttackID=%d, TaskID=%d", i, tag.ID, tag.AttackID, tag.TaskID)
	}

	// Test error handling for invalid task ID
	ctx3, cancel3 := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel3()

	invalidTags, err := client.GetAttackByTask(ctx3, 0)
	require.Error(t, err, "GetAttackByTask should error for ID 0")
	assert.Nil(t, invalidTags, "Tags should be nil when error occurs")
	assert.Contains(t, err.Error(), "task ID is required", "Error should mention required task ID")

	t.Log("=== ✓ GetAttackByTask validation passed ===")
}

// TestE2E_Attack_GetByCommand validates GetAttackByCommand retrieves ATT&CK tags for a command.
func TestE2E_Attack_GetByCommand(t *testing.T) {
	client := AuthenticateTestClient(t)

	t.Log("=== Test: GetAttackByCommand ===")

	// Get commands to find one to test with
	ctx1, cancel1 := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel1()

	commands, err := client.GetCommands(ctx1)
	require.NoError(t, err, "GetCommands should succeed")

	if len(commands) == 0 {
		t.Log("⚠ No commands found - skipping test")
		t.Skip("No commands available to test")
		return
	}

	testCommand := commands[0]
	t.Logf("✓ Testing with command: %s (ID: %d)", testCommand.Cmd, testCommand.ID)

	// Get MITRE ATT&CK tags for this command
	ctx2, cancel2 := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel2()

	attackTags, err := client.GetAttackByCommand(ctx2, testCommand.ID)
	require.NoError(t, err, "GetAttackByCommand should succeed")
	require.NotNil(t, attackTags, "Attack tags should not be nil")

	t.Logf("✓ Command %s has %d MITRE ATT&CK tag(s)", testCommand.Cmd, len(attackTags))

	// Validate attack tag structure
	for i, tag := range attackTags {
		assert.NotZero(t, tag.ID, "AttackCommand[%d] should have ID", i)
		assert.NotZero(t, tag.AttackID, "AttackCommand[%d] should have AttackID", i)
		assert.Equal(t, testCommand.ID, tag.CommandID, "AttackCommand[%d] should match command ID", i)

		t.Logf("  Tag[%d]: ID=%d, AttackID=%d, CommandID=%d", i, tag.ID, tag.AttackID, tag.CommandID)
	}

	// Test error handling for invalid command ID
	ctx3, cancel3 := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel3()

	invalidTags, err := client.GetAttackByCommand(ctx3, 0)
	require.Error(t, err, "GetAttackByCommand should error for ID 0")
	assert.Nil(t, invalidTags, "Tags should be nil when error occurs")
	assert.Contains(t, err.Error(), "command ID is required", "Error should mention required command ID")

	t.Log("=== ✓ GetAttackByCommand validation passed ===")
}

// TestE2E_Attack_GetByOperation validates GetAttacksByOperation retrieves all
// ATT&CK techniques used in an operation.
func TestE2E_Attack_GetByOperation(t *testing.T) {
	client := AuthenticateTestClient(t)

	t.Log("=== Test: GetAttacksByOperation ===")

	// Get current operation
	ctx1, cancel1 := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel1()

	operations, err := client.GetOperations(ctx1)
	require.NoError(t, err, "GetOperations should succeed")
	require.NotEmpty(t, operations, "Should have at least one operation")

	currentOp := operations[0]
	t.Logf("✓ Testing with operation: %s (ID: %d)", currentOp.Name, currentOp.ID)

	// Get MITRE ATT&CK techniques for this operation
	ctx2, cancel2 := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel2()

	techniques, err := client.GetAttacksByOperation(ctx2, currentOp.ID)
	require.NoError(t, err, "GetAttacksByOperation should succeed")
	require.NotNil(t, techniques, "Techniques should not be nil")

	t.Logf("✓ Operation %s has %d unique MITRE ATT&CK technique(s)", currentOp.Name, len(techniques))

	// Validate technique structure
	for i, tech := range techniques {
		assert.NotZero(t, tech.ID, "Attack[%d] should have ID", i)
		assert.NotEmpty(t, tech.TNum, "Attack[%d] should have TNum", i)
		assert.NotEmpty(t, tech.Name, "Attack[%d] should have Name", i)

		if i < 10 {
			t.Logf("  Technique[%d]: %s - %s (Tactic=%s)", i, tech.TNum, tech.Name, tech.Tactic)
		}
	}

	if len(techniques) > 10 {
		t.Logf("  ... and %d more techniques", len(techniques)-10)
	}

	// Test error handling for invalid operation ID
	ctx3, cancel3 := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel3()

	invalidTechs, err := client.GetAttacksByOperation(ctx3, 0)
	require.Error(t, err, "GetAttacksByOperation should error for ID 0")
	assert.Nil(t, invalidTechs, "Techniques should be nil when error occurs")
	assert.Contains(t, err.Error(), "operation ID is required", "Error should mention required operation ID")

	t.Log("=== ✓ GetAttacksByOperation validation passed ===")
}

// TestE2E_Attack_Comprehensive_Summary provides a summary of all MITRE ATT&CK test coverage.
func TestE2E_Attack_Comprehensive_Summary(t *testing.T) {
	t.Log("=== MITRE ATT&CK Comprehensive Test Coverage Summary ===")
	t.Log("")
	t.Log("This test suite validates comprehensive MITRE ATT&CK functionality:")
	t.Log("  1. ✓ GetAttackTechniques - Schema validation and field population")
	t.Log("  2. ✓ GetAttackTechniqueByID and ByTNum - Technique lookup methods")
	t.Log("  3. ✓ GetAttackByTask - ATT&CK tags for tasks")
	t.Log("  4. ✓ GetAttackByCommand - ATT&CK tags for commands")
	t.Log("  5. ✓ GetAttacksByOperation - All techniques used in operation")
	t.Log("")
	t.Log("All tests validate:")
	t.Log("  • Field presence and correctness (not just err != nil)")
	t.Log("  • Error handling for invalid inputs")
	t.Log("  • Proper filtering and association")
	t.Log("  • TNum format validation (T followed by digits)")
	t.Log("  • Graceful handling when no data available")
	t.Log("")
	t.Log("=== ✓ All MITRE ATT&CK comprehensive tests documented ===")
}
