package unit

import (
	"testing"
	"time"

	"github.com/nbaertsch/mythic-sdk-go/pkg/mythic/types"
)

// TestArtifactString tests the Artifact.String() method
func TestArtifactString(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name     string
		artifact types.Artifact
		contains []string
	}{
		{
			name: "with host and artifact",
			artifact: types.Artifact{
				ID:        1,
				Artifact:  "C:\\Windows\\Temp\\malware.exe",
				Host:      "WORKSTATION-01",
				Timestamp: now,
			},
			contains: []string{"malware.exe", "WORKSTATION-01"},
		},
		{
			name: "with artifact only",
			artifact: types.Artifact{
				ID:        2,
				Artifact:  "192.168.1.100:4444",
				Timestamp: now,
			},
			contains: []string{"192.168.1.100:4444"},
		},
		{
			name: "with ID only",
			artifact: types.Artifact{
				ID:        3,
				Timestamp: now,
			},
			contains: []string{"Artifact 3"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.artifact.String()
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

// TestArtifactIsDeleted tests removed - taskartifact table doesn't support soft delete

// TestArtifactHasTask tests the Artifact.HasTask() method
func TestArtifactHasTask(t *testing.T) {
	taskID := 42
	zeroTaskID := 0

	tests := []struct {
		name     string
		taskID   *int
		expected bool
	}{
		{"with task", &taskID, true},
		{"without task", nil, false},
		{"with zero task", &zeroTaskID, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			artifact := types.Artifact{TaskID: tt.taskID}
			result := artifact.HasTask()
			if result != tt.expected {
				t.Errorf("HasTask() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

// TestArtifactTypes tests the Artifact type structure
func TestArtifactTypes(t *testing.T) {
	now := time.Now()
	taskID := 123

	artifact := types.Artifact{
		ID:           1,
		Artifact:     "C:\\Windows\\System32\\evil.dll",
		BaseArtifact: "C:\\Windows\\System32",
		Host:         "SERVER-01",
		OperationID:  5,
		TaskID:       &taskID,
		Timestamp:    now,
	}

	if artifact.ID != 1 {
		t.Errorf("Expected ID 1, got %d", artifact.ID)
	}
	if artifact.Artifact != "C:\\Windows\\System32\\evil.dll" {
		t.Errorf("Expected Artifact 'C:\\Windows\\System32\\evil.dll', got %q", artifact.Artifact)
	}
	if !artifact.HasTask() {
		t.Error("Expected artifact to have a task")
	}
	if artifact.TaskID == nil || *artifact.TaskID != 123 {
		t.Error("Expected TaskID to be 123")
	}
}

// TestCreateArtifactRequest tests CreateArtifactRequest structure
func TestCreateArtifactRequest(t *testing.T) {
	baseArtifact := "C:\\Windows\\Temp"
	host := "WORKSTATION-01"
	taskID := 42

	req := types.CreateArtifactRequest{
		Artifact:     "C:\\Windows\\Temp\\payload.exe",
		BaseArtifact: &baseArtifact,
		Host:         &host,
		TaskID:       &taskID,
	}

	if req.Artifact != "C:\\Windows\\Temp\\payload.exe" {
		t.Errorf("Expected Artifact 'C:\\Windows\\Temp\\payload.exe', got %q", req.Artifact)
	}
	if req.BaseArtifact == nil || *req.BaseArtifact != baseArtifact {
		t.Error("Expected BaseArtifact to be 'C:\\Windows\\Temp'")
	}
	if req.Host == nil || *req.Host != host {
		t.Error("Expected Host to be 'WORKSTATION-01'")
	}
	if req.TaskID == nil || *req.TaskID != taskID {
		t.Error("Expected TaskID to be 42")
	}
}

// TestUpdateArtifactRequest tests UpdateArtifactRequest structure
func TestUpdateArtifactRequest(t *testing.T) {
	host := "NEW-HOST"

	req := types.UpdateArtifactRequest{
		ID:   5,
		Host: &host,
	}

	if req.ID != 5 {
		t.Errorf("Expected ID 5, got %d", req.ID)
	}
	if req.Host == nil || *req.Host != host {
		t.Error("Expected Host to be 'NEW-HOST'")
	}
}

// TestArtifactTypeConstants tests artifact type constants
func TestArtifactTypeConstants(t *testing.T) {
	types := map[string]string{
		"file":           types.ArtifactTypeFile,
		"registry":       types.ArtifactTypeRegistry,
		"process":        types.ArtifactTypeProcess,
		"network":        types.ArtifactTypeNetwork,
		"user":           types.ArtifactTypeUser,
		"service":        types.ArtifactTypeService,
		"scheduled_task": types.ArtifactTypeScheduled,
		"wmi":            types.ArtifactTypeWMI,
		"other":          types.ArtifactTypeOther,
	}

	for expected, actual := range types {
		if actual != expected {
			t.Errorf("Expected artifact type %q, got %q", expected, actual)
		}
	}
}

// TestArtifactWithoutOptionalFields tests Artifact without optional fields
func TestArtifactWithoutOptionalFields(t *testing.T) {
	artifact := types.Artifact{
		ID:        1,
		Artifact:  "minimal-artifact",
		Timestamp: time.Now(),
	}

	if artifact.TaskID != nil {
		t.Error("TaskID should be nil")
	}
	if artifact.Operation != nil {
		t.Error("Operation should be nil")
	}
	if artifact.Operator != nil {
		t.Error("Operator should be nil")
	}
	if artifact.HasTask() {
		t.Error("Should not have task")
	}

	str := artifact.String()
	if str == "" {
		t.Error("String() should not return empty string even without optional fields")
	}
}

// TestArtifactFileTypes tests file artifact types
func TestArtifactFileTypes(t *testing.T) {
	fileArtifacts := []string{
		"C:\\Windows\\System32\\evil.exe",
		"/tmp/suspicious.sh",
		"\\\\SHARE\\malware\\payload.dll",
		"C:\\Users\\victim\\Downloads\\document.pdf",
	}

	for _, file := range fileArtifacts {
		artifact := types.Artifact{
			ID:        1,
			Artifact:  file,
			Timestamp: time.Now(),
		}

		str := artifact.String()
		if !contains(str, file) {
			t.Errorf("String() should contain file path %q, got %q", file, str)
		}
	}
}

// TestArtifactNetworkTypes tests network artifact types
func TestArtifactNetworkTypes(t *testing.T) {
	networkArtifacts := []string{
		"192.168.1.100:4444",
		"example.com",
		"http://malicious-site.com/callback",
		"10.0.0.1",
	}

	for _, network := range networkArtifacts {
		artifact := types.Artifact{
			ID:        1,
			Artifact:  network,
			Timestamp: time.Now(),
		}

		str := artifact.String()
		if !contains(str, network) {
			t.Errorf("String() should contain network indicator %q, got %q", network, str)
		}
	}
}

// TestArtifactRegistryTypes tests registry artifact types
func TestArtifactRegistryTypes(t *testing.T) {
	registryArtifacts := []string{
		"HKLM\\Software\\Microsoft\\Windows\\CurrentVersion\\Run\\Malware",
		"HKCU\\Software\\Classes\\CLSID\\{12345678-1234-1234-1234-123456789012}",
	}

	for _, registry := range registryArtifacts {
		artifact := types.Artifact{
			ID:        1,
			Artifact:  registry,
			Timestamp: time.Now(),
		}

		str := artifact.String()
		if !contains(str, registry) {
			t.Errorf("String() should contain registry key %q, got %q", registry, str)
		}
	}
}

// TestArtifactMetadata removed - taskartifact table doesn't have metadata field

// TestArtifactTimestamp tests timestamp handling
func TestArtifactTimestamp(t *testing.T) {
	specificTime := time.Date(2024, 1, 15, 14, 30, 0, 0, time.UTC)

	artifact := types.Artifact{
		ID:        1,
		Artifact:  "test.exe",
		Timestamp: specificTime,
	}

	if !artifact.Timestamp.Equal(specificTime) {
		t.Errorf("Expected timestamp %v, got %v", specificTime, artifact.Timestamp)
	}
}
