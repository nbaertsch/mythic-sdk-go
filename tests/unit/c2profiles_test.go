package unit

import (
	"testing"
	"time"

	"github.com/nbaertsch/mythic-sdk-go/pkg/mythic/types"
)

// TestC2ProfileString tests the C2Profile.String() method
func TestC2ProfileString(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name     string
		profile  types.C2Profile
		contains []string
	}{
		{
			name: "running profile",
			profile: types.C2Profile{
				ID:           1,
				Name:         "http",
				Running:      true,
				CreationTime: now,
			},
			contains: []string{"http", "running"},
		},
		{
			name: "stopped profile",
			profile: types.C2Profile{
				ID:           2,
				Name:         "https",
				Running:      false,
				CreationTime: now,
			},
			contains: []string{"https", "stopped"},
		},
		{
			name: "profile with special name",
			profile: types.C2Profile{
				ID:           3,
				Name:         "websocket-secure",
				Running:      true,
				CreationTime: now,
			},
			contains: []string{"websocket-secure", "running"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.profile.String()
			if result == "" {
				t.Error("String() should not return empty string")
			}
			for _, substr := range tt.contains {
				if !stringContains(result, substr) {
					t.Errorf("String() = %q, should contain %q", result, substr)
				}
			}
		})
	}
}

// TestC2ProfileIsRunning tests the C2Profile.IsRunning() method
func TestC2ProfileIsRunning(t *testing.T) {
	tests := []struct {
		name     string
		running  bool
		expected bool
	}{
		{"running profile", true, true},
		{"stopped profile", false, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			profile := types.C2Profile{Running: tt.running}
			result := profile.IsRunning()
			if result != tt.expected {
				t.Errorf("IsRunning() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

// TestC2ProfileIsDeleted tests the C2Profile.IsDeleted() method
func TestC2ProfileIsDeleted(t *testing.T) {
	tests := []struct {
		name     string
		deleted  bool
		expected bool
	}{
		{"deleted profile", true, true},
		{"active profile", false, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			profile := types.C2Profile{Deleted: tt.deleted}
			result := profile.IsDeleted()
			if result != tt.expected {
				t.Errorf("IsDeleted() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

// TestC2ProfileTypes tests the C2Profile type structure
func TestC2ProfileTypes(t *testing.T) {
	now := time.Now()
	startTime := now.Add(-1 * time.Hour)
	payloadTypeID := 5

	profile := types.C2Profile{
		ID:            1,
		Name:          "http",
		Description:   "HTTP C2 Profile",
		CreationTime:  now,
		Running:       true,
		StartTime:     &startTime,
		Deleted:       false,
		IsP2P:         false,
		ContainerID:   "container123",
		PayloadTypeID: &payloadTypeID,
		Parameters: map[string]interface{}{
			"port": 8080,
			"host": "0.0.0.0",
		},
	}

	if profile.ID != 1 {
		t.Errorf("Expected ID 1, got %d", profile.ID)
	}
	if profile.Name != "http" {
		t.Errorf("Expected Name 'http', got %q", profile.Name)
	}
	if !profile.IsRunning() {
		t.Error("Expected profile to be running")
	}
	if profile.IsDeleted() {
		t.Error("Expected profile to not be deleted")
	}
	if profile.StartTime == nil {
		t.Error("Expected StartTime to be set")
	}
	if profile.PayloadTypeID == nil || *profile.PayloadTypeID != 5 {
		t.Error("Expected PayloadTypeID to be 5")
	}
	if profile.Parameters == nil {
		t.Error("Expected Parameters to be set")
	}
}

// TestC2ProfileWithoutOptionalFields tests C2Profile without optional fields
func TestC2ProfileWithoutOptionalFields(t *testing.T) {
	profile := types.C2Profile{
		ID:           1,
		Name:         "minimal",
		CreationTime: time.Now(),
	}

	if profile.StartTime != nil {
		t.Error("StartTime should be nil")
	}
	if profile.StopTime != nil {
		t.Error("StopTime should be nil")
	}
	if profile.PayloadTypeID != nil {
		t.Error("PayloadTypeID should be nil")
	}

	str := profile.String()
	if str == "" {
		t.Error("String() should not return empty string even without optional fields")
	}
}

// TestCreateC2InstanceRequest tests CreateC2InstanceRequest structure
func TestCreateC2InstanceRequest(t *testing.T) {
	description := "Test C2 Profile"
	operationID := 5

	req := types.CreateC2InstanceRequest{
		Name:        "test-http",
		Description: &description,
		OperationID: &operationID,
		Parameters: map[string]interface{}{
			"port":     8080,
			"callback": "https://example.com/callback",
		},
	}

	if req.Name != "test-http" {
		t.Errorf("Expected Name 'test-http', got %q", req.Name)
	}
	if req.Description == nil || *req.Description != "Test C2 Profile" {
		t.Error("Expected Description to be 'Test C2 Profile'")
	}
	if req.OperationID == nil || *req.OperationID != 5 {
		t.Error("Expected OperationID to be 5")
	}
	if req.Parameters == nil {
		t.Error("Expected Parameters to be set")
	}
}

// TestImportC2InstanceRequest tests ImportC2InstanceRequest structure
func TestImportC2InstanceRequest(t *testing.T) {
	config := `{"name":"http","parameters":{"port":8080}}`

	req := types.ImportC2InstanceRequest{
		Config: config,
		Name:   "imported-http",
	}

	if req.Config != config {
		t.Errorf("Expected Config %q, got %q", config, req.Config)
	}
	if req.Name != "imported-http" {
		t.Errorf("Expected Name 'imported-http', got %q", req.Name)
	}
}

// TestC2ProfileOutputTypes tests C2ProfileOutput structure
func TestC2ProfileOutputTypes(t *testing.T) {
	output := types.C2ProfileOutput{
		Output: "Full combined output",
		StdOut: "Standard output",
		StdErr: "Error output",
	}

	if output.Output != "Full combined output" {
		t.Errorf("Expected Output 'Full combined output', got %q", output.Output)
	}
	if output.StdOut != "Standard output" {
		t.Errorf("Expected StdOut 'Standard output', got %q", output.StdOut)
	}
	if output.StdErr != "Error output" {
		t.Errorf("Expected StdErr 'Error output', got %q", output.StdErr)
	}
}

// TestC2SampleMessageTypes tests C2SampleMessage structure
func TestC2SampleMessageTypes(t *testing.T) {
	msg := types.C2SampleMessage{
		Message: "GET /api/callback HTTP/1.1",
		Metadata: map[string]interface{}{
			"method": "GET",
			"path":   "/api/callback",
		},
	}

	if msg.Message != "GET /api/callback HTTP/1.1" {
		t.Errorf("Expected Message 'GET /api/callback HTTP/1.1', got %q", msg.Message)
	}
	if msg.Metadata == nil {
		t.Error("Expected Metadata to be set")
	}
}

// TestC2IOCTypes tests C2IOC structure
func TestC2IOCTypes(t *testing.T) {
	iocs := types.C2IOC{
		ProfileID: 1,
		IOCs: []string{
			"example.com",
			"192.168.1.1",
			"/api/callback",
		},
		Type: "network",
	}

	if iocs.ProfileID != 1 {
		t.Errorf("Expected ProfileID 1, got %d", iocs.ProfileID)
	}
	if len(iocs.IOCs) != 3 {
		t.Errorf("Expected 3 IOCs, got %d", len(iocs.IOCs))
	}
	if iocs.Type != "network" {
		t.Errorf("Expected Type 'network', got %q", iocs.Type)
	}
}

// TestC2ProfileFlags tests C2Profile boolean flags
func TestC2ProfileFlags(t *testing.T) {
	profile := types.C2Profile{
		ID:      1,
		Name:    "test",
		IsP2P:   true,
		Running: true,
		Deleted: false,
	}

	if !profile.IsP2P {
		t.Error("Expected IsP2P to be true")
	}
	if !profile.IsRunning() {
		t.Error("Expected profile to be running")
	}
	if profile.IsDeleted() {
		t.Error("Expected profile to not be deleted")
	}
}

// TestC2ProfileLifecycle tests profile lifecycle states
func TestC2ProfileLifecycle(t *testing.T) {
	now := time.Now()
	startTime := now.Add(-2 * time.Hour)
	stopTime := now.Add(-1 * time.Hour)

	// Profile that was started and stopped
	profile := types.C2Profile{
		ID:           1,
		Name:         "lifecycle-test",
		CreationTime: now.Add(-3 * time.Hour),
		Running:      false,
		StartTime:    &startTime,
		StopTime:     &stopTime,
	}

	if profile.IsRunning() {
		t.Error("Expected profile to not be running")
	}
	if profile.StartTime == nil {
		t.Error("Expected StartTime to be set")
	}
	if profile.StopTime == nil {
		t.Error("Expected StopTime to be set")
	}
	if profile.StopTime.Before(*profile.StartTime) {
		t.Error("StopTime should be after StartTime")
	}
}

// TestC2ProfileParameters tests parameter handling
func TestC2ProfileParameters(t *testing.T) {
	profile := types.C2Profile{
		ID:   1,
		Name: "param-test",
		Parameters: map[string]interface{}{
			"port":        8080,
			"host":        "0.0.0.0",
			"ssl":         true,
			"user_agent":  "Mozilla/5.0",
			"headers":     []string{"X-Custom-Header: value"},
			"retry_count": 3,
		},
	}

	if profile.Parameters == nil {
		t.Fatal("Parameters should not be nil")
	}

	// Verify different parameter types
	if port, ok := profile.Parameters["port"].(int); !ok || port != 8080 {
		t.Error("Expected port parameter to be 8080")
	}
	if ssl, ok := profile.Parameters["ssl"].(bool); !ok || !ssl {
		t.Error("Expected ssl parameter to be true")
	}
	if host, ok := profile.Parameters["host"].(string); !ok || host != "0.0.0.0" {
		t.Error("Expected host parameter to be '0.0.0.0'")
	}
}

// stringContains checks if a string contains a substring
func stringContains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && findSubstr(s, substr))
}

func findSubstr(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
