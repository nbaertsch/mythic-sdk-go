package unit

import (
	"testing"
	"time"

	"github.com/nbaertsch/mythic-sdk-go/pkg/mythic/types"
)

// TestAttackString tests the Attack.String() method
func TestAttackString(t *testing.T) {
	tests := []struct {
		name     string
		attack   types.Attack
		contains []string
	}{
		{
			name: "with name and technique number",
			attack: types.Attack{
				ID:   1,
				TNum: "T1003",
				Name: "OS Credential Dumping",
			},
			contains: []string{"T1003", "OS Credential Dumping"},
		},
		{
			name: "with technique number only",
			attack: types.Attack{
				ID:   2,
				TNum: "T1059",
			},
			contains: []string{"T1059"},
		},
		{
			name: "with ID only",
			attack: types.Attack{
				ID: 3,
			},
			contains: []string{"Attack 3"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.attack.String()
			if result == "" {
				t.Error("String() should not return empty string")
			}
			for _, substr := range tt.contains {
				if !contains(result, substr) {
					t.Errorf("String() = %q, should contain %q", result, substr)
				}
			}
		})
	}
}

// TestAttackTypes tests the Attack type structure
func TestAttackTypes(t *testing.T) {
	attack := types.Attack{
		ID:     1,
		TNum:   "T1003.001",
		Name:   "LSASS Memory",
		OS:     "Windows",
		Tactic: "Credential Access",
	}

	if attack.ID != 1 {
		t.Errorf("Expected ID 1, got %d", attack.ID)
	}
	if attack.TNum != "T1003.001" {
		t.Errorf("Expected TNum 'T1003.001', got %q", attack.TNum)
	}
	if attack.Name != "LSASS Memory" {
		t.Errorf("Expected Name 'LSASS Memory', got %q", attack.Name)
	}
	if attack.OS != "Windows" {
		t.Errorf("Expected OS 'Windows', got %q", attack.OS)
	}
	if attack.Tactic != "Credential Access" {
		t.Errorf("Expected Tactic 'Credential Access', got %q", attack.Tactic)
	}
}

// TestAttackTaskString tests the AttackTask.String() method
func TestAttackTaskString(t *testing.T) {
	tests := []struct {
		name       string
		attackTask types.AttackTask
		contains   []string
	}{
		{
			name: "with attack loaded",
			attackTask: types.AttackTask{
				ID:       1,
				AttackID: 5,
				TaskID:   42,
				Attack: &types.Attack{
					TNum: "T1003",
					Name: "OS Credential Dumping",
				},
			},
			contains: []string{"T1003", "OS Credential Dumping", "Task 42"},
		},
		{
			name: "without attack loaded",
			attackTask: types.AttackTask{
				ID:       2,
				AttackID: 10,
				TaskID:   100,
			},
			contains: []string{"AttackTask 2", "Attack 10", "Task 100"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.attackTask.String()
			if result == "" {
				t.Error("String() should not return empty string")
			}
			for _, substr := range tt.contains {
				if !contains(result, substr) {
					t.Errorf("String() = %q, should contain %q", result, substr)
				}
			}
		})
	}
}

// TestAttackTaskTypes tests the AttackTask type structure
func TestAttackTaskTypes(t *testing.T) {
	attackTask := types.AttackTask{
		ID:       1,
		AttackID: 5,
		TaskID:   42,
	}

	if attackTask.ID != 1 {
		t.Errorf("Expected ID 1, got %d", attackTask.ID)
	}
	if attackTask.AttackID != 5 {
		t.Errorf("Expected AttackID 5, got %d", attackTask.AttackID)
	}
	if attackTask.TaskID != 42 {
		t.Errorf("Expected TaskID 42, got %d", attackTask.TaskID)
	}
}

// TestAttackCommandString tests the AttackCommand.String() method
func TestAttackCommandString(t *testing.T) {
	tests := []struct {
		name          string
		attackCommand types.AttackCommand
		contains      []string
	}{
		{
			name: "with attack loaded",
			attackCommand: types.AttackCommand{
				ID:        1,
				AttackID:  5,
				CommandID: 10,
				Attack: &types.Attack{
					TNum: "T1059",
					Name: "Command and Scripting Interpreter",
				},
			},
			contains: []string{"T1059", "Command and Scripting Interpreter", "Command 10"},
		},
		{
			name: "without attack loaded",
			attackCommand: types.AttackCommand{
				ID:        2,
				AttackID:  15,
				CommandID: 20,
			},
			contains: []string{"AttackCommand 2", "Attack 15", "Command 20"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.attackCommand.String()
			if result == "" {
				t.Error("String() should not return empty string")
			}
			for _, substr := range tt.contains {
				if !contains(result, substr) {
					t.Errorf("String() = %q, should contain %q", result, substr)
				}
			}
		})
	}
}

// TestAttackCommandTypes tests the AttackCommand type structure
func TestAttackCommandTypes(t *testing.T) {
	attackCommand := types.AttackCommand{
		ID:        1,
		AttackID:  5,
		CommandID: 10,
	}

	if attackCommand.ID != 1 {
		t.Errorf("Expected ID 1, got %d", attackCommand.ID)
	}
	if attackCommand.AttackID != 5 {
		t.Errorf("Expected AttackID 5, got %d", attackCommand.AttackID)
	}
	if attackCommand.CommandID != 10 {
		t.Errorf("Expected CommandID 10, got %d", attackCommand.CommandID)
	}
}

// TestAttackTimestamp removed - attack table doesn't have timestamp field

// TestAttackTechniqueNumbers tests various technique number formats
func TestAttackTechniqueNumbers(t *testing.T) {
	techniques := []string{
		"T1003",     // Basic technique
		"T1003.001", // Sub-technique
		"T1059.001", // PowerShell
		"T1055.012", // Process Hollowing
		"T1078.004", // Cloud Accounts
	}

	for _, tNum := range techniques {
		attack := types.Attack{
			ID:   1,
			TNum: tNum,
			Name: "Test Technique",
		}

		if attack.TNum != tNum {
			t.Errorf("Expected TNum %q, got %q", tNum, attack.TNum)
		}

		str := attack.String()
		if !contains(str, tNum) {
			t.Errorf("String() should contain technique number %q, got %q", tNum, str)
		}
	}
}

// TestAttackTactics tests various MITRE ATT&CK tactics
func TestAttackTactics(t *testing.T) {
	tactics := []string{
		"Initial Access",
		"Execution",
		"Persistence",
		"Privilege Escalation",
		"Defense Evasion",
		"Credential Access",
		"Discovery",
		"Lateral Movement",
		"Collection",
		"Command and Control",
		"Exfiltration",
		"Impact",
	}

	for _, tactic := range tactics {
		attack := types.Attack{
			ID:     1,
			TNum:   "T1000",
			Name:   "Test Technique",
			Tactic: tactic,
		}

		if attack.Tactic != tactic {
			t.Errorf("Expected Tactic %q, got %q", tactic, attack.Tactic)
		}
	}
}

// TestAttackOS tests various operating system values
func TestAttackOS(t *testing.T) {
	operatingSystems := []string{
		"Windows",
		"Linux",
		"macOS",
		"Network",
		"Containers",
		"IaaS",
		"SaaS",
		"Office 365",
		"Azure AD",
		"Google Workspace",
	}

	for _, os := range operatingSystems {
		attack := types.Attack{
			ID:   1,
			TNum: "T1000",
			Name: "Test Technique",
			OS:   os,
		}

		if attack.OS != os {
			t.Errorf("Expected OS %q, got %q", os, attack.OS)
		}
	}
}
