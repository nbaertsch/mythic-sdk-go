package unit

import (
	"testing"
	"time"

	"github.com/nbaertsch/mythic-sdk-go/pkg/mythic/types"
)

// TestProcessString tests the Process.String() method
func TestProcessString(t *testing.T) {
	tests := []struct {
		name     string
		process  types.Process
		expected string
	}{
		{
			name: "with name and PID",
			process: types.Process{
				Name:      "chrome.exe",
				ProcessID: 1234,
			},
			expected: "chrome.exe (PID: 1234)",
		},
		{
			name: "with PID only",
			process: types.Process{
				ProcessID: 5678,
			},
			expected: "PID 5678",
		},
		{
			name: "with ID only",
			process: types.Process{
				ID: 42,
			},
			expected: "Process 42",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.process.String()
			if result != tt.expected {
				t.Errorf("Process.String() = %q, expected %q", result, tt.expected)
			}
		})
	}
}

// TestProcessIsDeleted tests the Process.IsDeleted() method
func TestProcessIsDeleted(t *testing.T) {
	tests := []struct {
		name     string
		deleted  bool
		expected bool
	}{
		{"deleted process", true, true},
		{"active process", false, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			process := types.Process{Deleted: tt.deleted}
			result := process.IsDeleted()
			if result != tt.expected {
				t.Errorf("IsDeleted() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

// TestProcessHasParent tests the Process.HasParent() method
func TestProcessHasParent(t *testing.T) {
	tests := []struct {
		name            string
		parentProcessID int
		expected        bool
	}{
		{"process with parent", 1000, true},
		{"root process", 0, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			process := types.Process{ParentProcessID: tt.parentProcessID}
			result := process.HasParent()
			if result != tt.expected {
				t.Errorf("HasParent() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

// TestProcessGetIntegrityLevelString tests the Process.GetIntegrityLevelString() method
func TestProcessGetIntegrityLevelString(t *testing.T) {
	tests := []struct {
		name           string
		integrityLevel int
		expected       string
	}{
		{"untrusted", 0, "Untrusted"},
		{"low", 1, "Low"},
		{"medium", 2, "Medium"},
		{"high", 3, "High"},
		{"system", 4, "System"},
		{"unknown", 99, "Unknown (99)"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			process := types.Process{IntegrityLevel: tt.integrityLevel}
			result := process.GetIntegrityLevelString()
			if result != tt.expected {
				t.Errorf("GetIntegrityLevelString() = %q, expected %q", result, tt.expected)
			}
		})
	}
}

// TestProcessTypes tests the Process type structure
func TestProcessTypes(t *testing.T) {
	now := time.Now()
	callbackID := 10
	taskID := 42

	process := types.Process{
		ID:              1,
		Name:            "explorer.exe",
		ProcessID:       1234,
		ParentProcessID: 4,
		Architecture:    "x64",
		BinPath:         "C:\\Windows\\explorer.exe",
		User:            "DOMAIN\\user",
		CommandLine:     "C:\\Windows\\explorer.exe",
		IntegrityLevel:  2,
		StartTime:       now,
		Description:     "Windows Explorer",
		OperationID:     5,
		HostID:          3,
		CallbackID:      &callbackID,
		TaskID:          &taskID,
		Timestamp:       now,
		Deleted:         false,
	}

	if process.ID != 1 {
		t.Errorf("Expected ID 1, got %d", process.ID)
	}
	if process.Name != "explorer.exe" {
		t.Errorf("Expected Name 'explorer.exe', got %q", process.Name)
	}
	if process.ProcessID != 1234 {
		t.Errorf("Expected ProcessID 1234, got %d", process.ProcessID)
	}
	if !process.HasParent() {
		t.Error("Expected process to have a parent")
	}
	if process.IsDeleted() {
		t.Error("Expected process not to be deleted")
	}
	if process.GetIntegrityLevelString() != "Medium" {
		t.Errorf("Expected Medium integrity, got %q", process.GetIntegrityLevelString())
	}
}

// TestHostTypes tests the Host type structure
func TestHostTypes(t *testing.T) {
	host := types.Host{
		ID:   1,
		Host: "WORKSTATION-01",
	}

	if host.ID != 1 {
		t.Errorf("Expected ID 1, got %d", host.ID)
	}
	if host.Host != "WORKSTATION-01" {
		t.Errorf("Expected Host 'WORKSTATION-01', got %q", host.Host)
	}
}

// TestProcessTreeStructure tests the ProcessTree type structure
func TestProcessTreeStructure(t *testing.T) {
	parent := &types.Process{
		ID:        1,
		Name:      "parent.exe",
		ProcessID: 100,
	}

	child1 := &types.Process{
		ID:              2,
		Name:            "child1.exe",
		ProcessID:       200,
		ParentProcessID: 100,
	}

	child2 := &types.Process{
		ID:              3,
		Name:            "child2.exe",
		ProcessID:       300,
		ParentProcessID: 100,
	}

	tree := &types.ProcessTree{
		Process: parent,
		Children: []*types.ProcessTree{
			{Process: child1, Children: []*types.ProcessTree{}},
			{Process: child2, Children: []*types.ProcessTree{}},
		},
	}

	if tree.Process.ProcessID != 100 {
		t.Errorf("Expected root PID 100, got %d", tree.Process.ProcessID)
	}

	if len(tree.Children) != 2 {
		t.Errorf("Expected 2 children, got %d", len(tree.Children))
	}

	if tree.Children[0].Process.ParentProcessID != 100 {
		t.Error("Expected child to have parent PID 100")
	}
}

// TestProcessWithNilFields tests handling of nil fields
func TestProcessWithNilFields(t *testing.T) {
	process := types.Process{
		ID:         1,
		Name:       "test.exe",
		ProcessID:  1234,
		CallbackID: nil,
		TaskID:     nil,
	}

	if process.CallbackID != nil {
		t.Error("Expected CallbackID to be nil")
	}

	if process.TaskID != nil {
		t.Error("Expected TaskID to be nil")
	}

	str := process.String()
	if str == "" {
		t.Error("String() should not return empty string")
	}
}

// TestProcessIntegrityLevels tests all integrity levels
func TestProcessIntegrityLevels(t *testing.T) {
	levels := []struct {
		level    int
		expected string
	}{
		{0, "Untrusted"},
		{1, "Low"},
		{2, "Medium"},
		{3, "High"},
		{4, "System"},
	}

	for _, l := range levels {
		process := types.Process{IntegrityLevel: l.level}
		result := process.GetIntegrityLevelString()
		if result != l.expected {
			t.Errorf("Integrity level %d: expected %q, got %q", l.level, l.expected, result)
		}
	}
}

// TestProcessArchitectures tests common architectures
func TestProcessArchitectures(t *testing.T) {
	architectures := []string{"x86", "x64", "arm", "arm64"}

	for _, arch := range architectures {
		process := types.Process{
			Name:         "test.exe",
			Architecture: arch,
		}

		if process.Architecture != arch {
			t.Errorf("Expected architecture %q, got %q", arch, process.Architecture)
		}
	}
}

// TestProcessCommandLineHandling tests command line storage
func TestProcessCommandLineHandling(t *testing.T) {
	longCommandLine := "C:\\Windows\\System32\\cmd.exe /c \"echo Hello World && dir C:\\ && whoami\""

	process := types.Process{
		Name:        "cmd.exe",
		ProcessID:   1234,
		CommandLine: longCommandLine,
	}

	if process.CommandLine != longCommandLine {
		t.Error("Command line should be preserved exactly")
	}

	if len(process.CommandLine) == 0 {
		t.Error("Command line should not be empty")
	}
}

// TestProcessParentChildRelationship tests parent-child relationships
func TestProcessParentChildRelationship(t *testing.T) {
	parent := types.Process{
		ProcessID:       1000,
		ParentProcessID: 0,
	}

	child := types.Process{
		ProcessID:       2000,
		ParentProcessID: 1000,
	}

	grandchild := types.Process{
		ProcessID:       3000,
		ParentProcessID: 2000,
	}

	if parent.HasParent() {
		t.Error("Root process should not have a parent")
	}

	if !child.HasParent() {
		t.Error("Child process should have a parent")
	}

	if !grandchild.HasParent() {
		t.Error("Grandchild process should have a parent")
	}

	if child.ParentProcessID != parent.ProcessID {
		t.Error("Child's parent should match parent's PID")
	}
}
