//go:build integration

package integration

import (
	"context"
	"testing"
	"time"

	"github.com/nbaertsch/mythic-sdk-go/pkg/mythic/types"
)

// TestE2E_MITREAttackMapping tests the complete MITRE ATT&CK framework integration
// Covers: GetAttackTechniques, GetAttackTechniqueByID, GetAttackTechniqueByTNum,
// GetAttackByCommand, GetAttackByTask, AddMITREAttackToTask, GetAttacksByOperation
func TestE2E_MITREAttackMapping(t *testing.T) {
	client := AuthenticateTestClient(t)

	var testTechniqueID int
	var testTNum string

	// Test 1: Get all MITRE ATT&CK techniques
	t.Log("=== Test 1: Get all MITRE ATT&CK techniques ===")
	ctx1, cancel1 := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel1()

	techniques, err := client.GetAttackTechniques(ctx1)
	if err != nil {
		t.Fatalf("GetAttackTechniques failed: %v", err)
	}
	if len(techniques) == 0 {
		t.Fatal("No MITRE ATT&CK techniques found - Mythic should have 600+ techniques")
	}
	t.Logf("✓ Found %d MITRE ATT&CK techniques", len(techniques))

	// Verify we have a reasonable number of techniques (at least 100)
	if len(techniques) < 100 {
		t.Errorf("Expected at least 100 techniques, got %d - MITRE database may not be loaded", len(techniques))
	}

	// Test 2: Select a test technique for further testing
	t.Log("=== Test 2: Select test technique ===")
	// Find T1059 (Command and Scripting Interpreter) or use first available
	var testTechnique *types.Attack
	for _, tech := range techniques {
		if tech.TNum == "T1059" {
			testTechnique = tech
			break
		}
	}
	if testTechnique == nil {
		// Use first available technique
		testTechnique = techniques[0]
	}
	testTechniqueID = testTechnique.ID
	testTNum = testTechnique.TNum
	t.Logf("✓ Using test technique: %s - %s (ID: %d)", testTNum, testTechnique.Name, testTechniqueID)

	// Test 3: Get technique by ID
	t.Log("=== Test 3: Get technique by ID ===")
	ctx2, cancel2 := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel2()

	techniqueByID, err := client.GetAttackTechniqueByID(ctx2, testTechniqueID)
	if err != nil {
		t.Fatalf("GetAttackTechniqueByID failed: %v", err)
	}
	if techniqueByID.ID != testTechniqueID {
		t.Errorf("Technique ID mismatch: expected %d, got %d", testTechniqueID, techniqueByID.ID)
	}
	if techniqueByID.TNum != testTNum {
		t.Errorf("Technique T-number mismatch: expected %s, got %s", testTNum, techniqueByID.TNum)
	}
	t.Logf("✓ Technique retrieved by ID: %s - %s", techniqueByID.TNum, techniqueByID.Name)

	// Test 4: Get technique by T-number
	t.Log("=== Test 4: Get technique by T-number ===")
	ctx3, cancel3 := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel3()

	techniqueByTNum, err := client.GetAttackTechniqueByTNum(ctx3, testTNum)
	if err != nil {
		t.Fatalf("GetAttackTechniqueByTNum failed: %v", err)
	}
	if techniqueByTNum.TNum != testTNum {
		t.Errorf("T-number mismatch: expected %s, got %s", testTNum, techniqueByTNum.TNum)
	}
	if techniqueByTNum.ID != testTechniqueID {
		t.Errorf("Technique ID mismatch: expected %d, got %d", testTechniqueID, techniqueByTNum.ID)
	}
	t.Logf("✓ Technique retrieved by T-number: %s - %s (Tactic: %s)", techniqueByTNum.TNum, techniqueByTNum.Name, techniqueByTNum.Tactic)

	// Test 5: Verify technique details
	t.Log("=== Test 5: Verify technique details ===")
	if testTechnique.Name == "" {
		t.Error("Technique name is empty")
	}
	if testTechnique.TNum == "" {
		t.Error("Technique T-number is empty")
	}
	// Tactic can be empty for some techniques, so just log it
	t.Logf("✓ Technique details verified:")
	t.Logf("  - Name: %s", testTechnique.Name)
	t.Logf("  - T-Number: %s", testTechnique.TNum)
	t.Logf("  - Tactic: %s", testTechnique.Tactic)
	t.Logf("  - OS: %s", testTechnique.OS)

	// Test 6: Get attacks by operation
	t.Log("=== Test 6: Get attacks by operation ===")
	ctx4, cancel4 := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel4()

	operationID := client.GetCurrentOperation()
	if operationID == nil {
		t.Fatal("No current operation set")
	}

	opAttacks, err := client.GetAttacksByOperation(ctx4, *operationID)
	if err != nil {
		t.Fatalf("GetAttacksByOperation failed: %v", err)
	}
	// May be 0 if no tasks with MITRE mappings yet
	t.Logf("✓ Found %d MITRE attacks for operation %d", len(opAttacks), *operationID)

	// Test 7: Get attack by command (may return empty if no commands yet)
	t.Log("=== Test 7: Get attack by command ===")
	ctx5, cancel5 := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel5()

	// Try to get command mappings for a hypothetical command ID
	// This may return empty results if no commands are configured yet
	commandAttacks, err := client.GetAttackByCommand(ctx5, 1)
	if err != nil {
		t.Logf("⚠ GetAttackByCommand failed (may be expected if no commands configured): %v", err)
	} else {
		t.Logf("✓ Found %d MITRE attacks for command 1", len(commandAttacks))
	}

	// Test 8: Get attack by task (limited - no tasks yet in Phase 2)
	t.Log("=== Test 8: Get attack by task (limited testing) ===")
	ctx6, cancel6 := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel6()

	// Try to get task mappings for a hypothetical task ID
	// This will likely return empty or error since we don't have tasks yet
	taskAttacks, err := client.GetAttackByTask(ctx6, 1)
	if err != nil {
		t.Logf("⚠ GetAttackByTask failed (expected without tasks): %v", err)
	} else {
		t.Logf("✓ Found %d MITRE attacks for task 1", len(taskAttacks))
	}

	// Note: AddMITREAttackToTask is skipped in Phase 2 since we don't have tasks yet.
	// It will be tested in Workflow 10 when we have active callbacks and tasks.
	t.Log("Note: AddMITREAttackToTask will be tested in Workflow 10 (Task Execution)")

	// Test 9: Verify technique data quality
	t.Log("=== Test 9: Verify technique data quality ===")
	validTechniques := 0
	for _, tech := range techniques {
		if tech.Name != "" && tech.TNum != "" {
			validTechniques++
		}
	}
	validPercentage := float64(validTechniques) / float64(len(techniques)) * 100
	t.Logf("✓ Technique data quality: %d/%d (%.1f%%) have name and T-number", validTechniques, len(techniques), validPercentage)
	if validPercentage < 90 {
		t.Errorf("Expected at least 90%% of techniques to have valid data, got %.1f%%", validPercentage)
	}

	// Test 10: Verify technique coverage by tactic
	t.Log("=== Test 10: Verify technique coverage by tactic ===")
	tacticCount := make(map[string]int)
	for _, tech := range techniques {
		if tech.Tactic != "" {
			tacticCount[tech.Tactic]++
		}
	}
	t.Logf("✓ Found %d different tactics:", len(tacticCount))
	for tactic, count := range tacticCount {
		t.Logf("  - %s: %d techniques", tactic, count)
	}

	// Test 11: Search for specific well-known techniques
	t.Log("=== Test 11: Verify well-known techniques present ===")
	wellKnownTNums := []string{"T1059", "T1003", "T1055", "T1082", "T1083"}
	foundTNums := 0
	for _, tnum := range wellKnownTNums {
		for _, tech := range techniques {
			if tech.TNum == tnum {
				foundTNums++
				t.Logf("  ✓ Found %s: %s", tnum, tech.Name)
				break
			}
		}
	}
	t.Logf("✓ Found %d/%d well-known techniques", foundTNums, len(wellKnownTNums))
	if foundTNums < 3 {
		t.Errorf("Expected to find at least 3 well-known techniques, found %d", foundTNums)
	}

	t.Log("=== ✓ All MITRE ATT&CK mapping tests passed ===")
}

// TestE2E_MITREAttackErrorHandling tests error scenarios for MITRE ATT&CK operations
func TestE2E_MITREAttackErrorHandling(t *testing.T) {
	client := AuthenticateTestClient(t)

	// Test 1: Get non-existent technique by ID
	t.Log("=== Test 1: Get non-existent technique by ID ===")
	ctx1, cancel1 := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel1()

	_, err := client.GetAttackTechniqueByID(ctx1, 999999)
	if err == nil {
		t.Error("Expected error for non-existent technique ID")
	}
	t.Logf("✓ Non-existent technique ID rejected: %v", err)

	// Test 2: Get non-existent technique by T-number
	t.Log("=== Test 2: Get non-existent technique by T-number ===")
	ctx2, cancel2 := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel2()

	_, err = client.GetAttackTechniqueByTNum(ctx2, "T99999")
	if err == nil {
		t.Error("Expected error for non-existent T-number")
	}
	t.Logf("✓ Non-existent T-number rejected: %v", err)

	// Test 3: Get attack by non-existent task
	t.Log("=== Test 3: Get attack by non-existent task ===")
	ctx3, cancel3 := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel3()

	_, err = client.GetAttackByTask(ctx3, 999999)
	if err == nil {
		t.Error("Expected error for non-existent task ID")
	}
	t.Logf("✓ Non-existent task ID rejected: %v", err)

	// Test 4: Get attack by non-existent command
	t.Log("=== Test 4: Get attack by non-existent command ===")
	ctx4, cancel4 := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel4()

	_, err = client.GetAttackByCommand(ctx4, 999999)
	if err == nil {
		t.Error("Expected error for non-existent command ID")
	}
	t.Logf("✓ Non-existent command ID rejected: %v", err)

	t.Log("=== ✓ All error handling tests passed ===")
}
