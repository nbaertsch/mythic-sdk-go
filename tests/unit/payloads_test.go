package unit

import (
	"testing"
	"time"

	"github.com/nbaertsch/mythic-sdk-go/pkg/mythic/types"
)

// TestPayloadString tests the Payload.String() method
func TestPayloadString(t *testing.T) {
	tests := []struct {
		name     string
		payload  types.Payload
		expected string
	}{
		{
			name: "ready payload with filename",
			payload: types.Payload{
				UUID:         "test-uuid-123",
				FilenameText: "agent.exe",
				OS:           "windows",
				BuildPhase:   "success",
			},
			expected: "agent.exe (windows, ready)",
		},
		{
			name: "building payload",
			payload: types.Payload{
				UUID:         "test-uuid-456",
				FilenameText: "agent.bin",
				OS:           "linux",
				BuildPhase:   "building",
			},
			expected: "agent.bin (linux, building)",
		},
		{
			name: "failed payload",
			payload: types.Payload{
				UUID:         "test-uuid-789",
				FilenameText: "agent.app",
				OS:           "macos",
				BuildPhase:   "error",
			},
			expected: "agent.app (macos, failed)",
		},
		{
			name: "payload without filename",
			payload: types.Payload{
				UUID:       "test-uuid-abc",
				OS:         "linux",
				BuildPhase: "success",
			},
			expected: "test-uuid-abc (linux, ready)",
		},
		{
			name: "submitted payload",
			payload: types.Payload{
				UUID:         "test-uuid-def",
				FilenameText: "agent.elf",
				OS:           "linux",
				BuildPhase:   "submitted",
			},
			expected: "agent.elf (linux, building)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.payload.String()
			if result != tt.expected {
				t.Errorf("Payload.String() = %q, expected %q", result, tt.expected)
			}
		})
	}
}

// TestPayloadIsReady tests the Payload.IsReady() method
func TestPayloadIsReady(t *testing.T) {
	tests := []struct {
		name       string
		buildPhase string
		expected   bool
	}{
		{"success phase", "success", true},
		{"building phase", "building", false},
		{"error phase", "error", false},
		{"submitted phase", "submitted", false},
		{"empty phase", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			payload := types.Payload{BuildPhase: tt.buildPhase}
			result := payload.IsReady()
			if result != tt.expected {
				t.Errorf("IsReady() = %v, expected %v for phase %q", result, tt.expected, tt.buildPhase)
			}
		})
	}
}

// TestPayloadIsFailed tests the Payload.IsFailed() method
func TestPayloadIsFailed(t *testing.T) {
	tests := []struct {
		name       string
		buildPhase string
		expected   bool
	}{
		{"error phase", "error", true},
		{"success phase", "success", false},
		{"building phase", "building", false},
		{"submitted phase", "submitted", false},
		{"empty phase", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			payload := types.Payload{BuildPhase: tt.buildPhase}
			result := payload.IsFailed()
			if result != tt.expected {
				t.Errorf("IsFailed() = %v, expected %v for phase %q", result, tt.expected, tt.buildPhase)
			}
		})
	}
}

// TestPayloadIsBuilding tests the Payload.IsBuilding() method
func TestPayloadIsBuilding(t *testing.T) {
	tests := []struct {
		name       string
		buildPhase string
		expected   bool
	}{
		{"building phase", "building", true},
		{"submitted phase", "submitted", true},
		{"success phase", "success", false},
		{"error phase", "error", false},
		{"empty phase", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			payload := types.Payload{BuildPhase: tt.buildPhase}
			result := payload.IsBuilding()
			if result != tt.expected {
				t.Errorf("IsBuilding() = %v, expected %v for phase %q", result, tt.expected, tt.buildPhase)
			}
		})
	}
}

// TestPayloadTypeString tests the PayloadType.String() method
func TestPayloadTypeString(t *testing.T) {
	tests := []struct {
		name        string
		payloadType types.PayloadType
		expected    string
	}{
		{
			name: "windows agent",
			payloadType: types.PayloadType{
				Name:          "poseidon",
				FileExtension: "exe",
			},
			expected: "poseidon (exe)",
		},
		{
			name: "linux agent",
			payloadType: types.PayloadType{
				Name:          "apfell",
				FileExtension: "bin",
			},
			expected: "apfell (bin)",
		},
		{
			name: "empty extension",
			payloadType: types.PayloadType{
				Name:          "service_wrapper",
				FileExtension: "",
			},
			expected: "service_wrapper ()",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.payloadType.String()
			if result != tt.expected {
				t.Errorf("PayloadType.String() = %q, expected %q", result, tt.expected)
			}
		})
	}
}

// TestPayloadOnHostString tests the PayloadOnHost.String() method
func TestPayloadOnHostString(t *testing.T) {
	tests := []struct {
		name          string
		payloadOnHost types.PayloadOnHost
		expected      string
	}{
		{
			name: "with host name",
			payloadOnHost: types.PayloadOnHost{
				PayloadID: 42,
				Host:      "workstation-01",
			},
			expected: "Payload 42 on workstation-01",
		},
		{
			name: "without host name",
			payloadOnHost: types.PayloadOnHost{
				PayloadID: 123,
				HostID:    456,
			},
			expected: "Payload 123 on Host ID 456",
		},
		{
			name: "empty host string",
			payloadOnHost: types.PayloadOnHost{
				PayloadID: 789,
				HostID:    101,
				Host:      "",
			},
			expected: "Payload 789 on Host ID 101",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.payloadOnHost.String()
			if result != tt.expected {
				t.Errorf("PayloadOnHost.String() = %q, expected %q", result, tt.expected)
			}
		})
	}
}

// TestPayloadTypes tests the Payload type structure
func TestPayloadTypes(t *testing.T) {
	now := time.Now()

	payload := types.Payload{
		ID:             1,
		UUID:           "payload-uuid-123",
		Description:    "Test payload",
		OperatorID:     10,
		OperationID:    5,
		CreationTime:   now,
		PayloadTypeID:  3,
		OS:             "linux",
		BuildContainer: "poseidon:latest",
		BuildPhase:     "success",
		BuildMessage:   "Build completed successfully",
		CallbackAlert:  true,
		AutoGenerated:  false,
		Deleted:        false,
		CallbacksCount: 2,
		FilenameText:   "agent.bin",
		TagStr:         "test-tag",
	}

	if payload.ID != 1 {
		t.Errorf("Expected ID 1, got %d", payload.ID)
	}
	if payload.UUID != "payload-uuid-123" {
		t.Errorf("Expected UUID 'payload-uuid-123', got %q", payload.UUID)
	}
	if !payload.IsReady() {
		t.Error("Expected payload to be ready")
	}
	if payload.IsFailed() {
		t.Error("Expected payload not to be failed")
	}
	if payload.IsBuilding() {
		t.Error("Expected payload not to be building")
	}
}

// TestPayloadTypeTypes tests the PayloadType type structure
func TestPayloadTypeTypes(t *testing.T) {
	now := time.Now()

	payloadType := types.PayloadType{
		ID:                  1,
		Name:                "poseidon",
		FileExtension:       "bin",
		Author:              "Test Author",
		Supported:           true,
		WrapperMode:         false,
		Wrapped:             false,
		Note:                "Test note",
		SupportsDynamicLoad: true,
		BuildURL:            "https://example.com/build",
		ExternalURL:         "https://example.com/docs",
		CreationTime:        now,
		Deleted:             false,
		ContainerRunning:    true,
		OS:                  "linux",
	}

	if payloadType.ID != 1 {
		t.Errorf("Expected ID 1, got %d", payloadType.ID)
	}
	if payloadType.Name != "poseidon" {
		t.Errorf("Expected Name 'poseidon', got %q", payloadType.Name)
	}
	if !payloadType.Supported {
		t.Error("Expected payload type to be supported")
	}
	if payloadType.WrapperMode {
		t.Error("Expected wrapper mode to be false")
	}
}

// TestPayloadCommandTypes tests the PayloadCommand type structure
func TestPayloadCommandTypes(t *testing.T) {
	command := types.PayloadCommand{
		ID:          1,
		CommandID:   42,
		PayloadID:   10,
		Version:     1,
		CommandName: "shell",
	}

	if command.ID != 1 {
		t.Errorf("Expected ID 1, got %d", command.ID)
	}
	if command.CommandName != "shell" {
		t.Errorf("Expected CommandName 'shell', got %q", command.CommandName)
	}
}

// TestPayloadC2ProfileTypes tests the PayloadC2Profile type structure
func TestPayloadC2ProfileTypes(t *testing.T) {
	profile := types.PayloadC2Profile{
		ID:            1,
		PayloadID:     10,
		C2ProfileID:   5,
		C2ProfileName: "http",
		Parameters: map[string]interface{}{
			"callback_host": "10.0.0.1",
			"callback_port": 443,
		},
	}

	if profile.ID != 1 {
		t.Errorf("Expected ID 1, got %d", profile.ID)
	}
	if profile.C2ProfileName != "http" {
		t.Errorf("Expected C2ProfileName 'http', got %q", profile.C2ProfileName)
	}
	if len(profile.Parameters) != 2 {
		t.Errorf("Expected 2 parameters, got %d", len(profile.Parameters))
	}
}

// TestBuildParameterTypes tests the BuildParameter type structure
func TestBuildParameterTypes(t *testing.T) {
	param := types.BuildParameter{
		Name:        "output_format",
		Value:       "exe",
		Description: "Output file format",
	}

	if param.Name != "output_format" {
		t.Errorf("Expected Name 'output_format', got %q", param.Name)
	}
	if param.Value != "exe" {
		t.Errorf("Expected Value 'exe', got %v", param.Value)
	}
}

// TestPayloadOnHostTypes tests the PayloadOnHost type structure
func TestPayloadOnHostTypes(t *testing.T) {
	now := time.Now()
	taskID := 42

	poh := types.PayloadOnHost{
		ID:          1,
		HostID:      10,
		PayloadID:   5,
		OperationID: 3,
		TaskID:      &taskID,
		Timestamp:   now,
		Deleted:     false,
		Host:        "workstation-01",
	}

	if poh.ID != 1 {
		t.Errorf("Expected ID 1, got %d", poh.ID)
	}
	if poh.TaskID == nil || *poh.TaskID != 42 {
		t.Errorf("Expected TaskID 42, got %v", poh.TaskID)
	}
	if poh.Host != "workstation-01" {
		t.Errorf("Expected Host 'workstation-01', got %q", poh.Host)
	}
}

// TestCreatePayloadRequestTypes tests the CreatePayloadRequest type structure
func TestCreatePayloadRequestTypes(t *testing.T) {
	req := types.CreatePayloadRequest{
		PayloadType: "poseidon",
		OS:          "linux",
		Description: "Test payload",
		Filename:    "agent.bin",
		Commands:    []string{"shell", "download", "upload"},
		C2Profiles: []types.C2ProfileConfig{
			{
				Name: "http",
				Parameters: map[string]interface{}{
					"callback_host": "10.0.0.1",
				},
			},
		},
		BuildParameters: map[string]interface{}{
			"output_format": "exe",
		},
		Tag:            "test-tag",
		SelectedOS:     "linux",
		WrapperPayload: "wrapper-uuid-123",
	}

	if req.PayloadType != "poseidon" {
		t.Errorf("Expected PayloadType 'poseidon', got %q", req.PayloadType)
	}
	if len(req.Commands) != 3 {
		t.Errorf("Expected 3 commands, got %d", len(req.Commands))
	}
	if len(req.C2Profiles) != 1 {
		t.Errorf("Expected 1 C2 profile, got %d", len(req.C2Profiles))
	}
}

// TestC2ProfileConfigTypes tests the C2ProfileConfig type structure
func TestC2ProfileConfigTypes(t *testing.T) {
	config := types.C2ProfileConfig{
		Name: "http",
		Parameters: map[string]interface{}{
			"callback_host": "10.0.0.1",
			"callback_port": 443,
			"useSSL":        true,
		},
	}

	if config.Name != "http" {
		t.Errorf("Expected Name 'http', got %q", config.Name)
	}
	if len(config.Parameters) != 3 {
		t.Errorf("Expected 3 parameters, got %d", len(config.Parameters))
	}
}

// TestUpdatePayloadRequestTypes tests the UpdatePayloadRequest type structure
func TestUpdatePayloadRequestTypes(t *testing.T) {
	description := "Updated description"
	callbackAlert := true
	deleted := false

	req := types.UpdatePayloadRequest{
		UUID:          "payload-uuid-123",
		Description:   &description,
		CallbackAlert: &callbackAlert,
		Deleted:       &deleted,
	}

	if req.UUID != "payload-uuid-123" {
		t.Errorf("Expected UUID 'payload-uuid-123', got %q", req.UUID)
	}
	if req.Description == nil || *req.Description != "Updated description" {
		t.Error("Expected Description to be set")
	}
	if req.CallbackAlert == nil || !*req.CallbackAlert {
		t.Error("Expected CallbackAlert to be true")
	}
	if req.Deleted == nil || *req.Deleted {
		t.Error("Expected Deleted to be false")
	}
}

// TestRebuildPayloadRequestTypes tests the RebuildPayloadRequest type structure
func TestRebuildPayloadRequestTypes(t *testing.T) {
	req := types.RebuildPayloadRequest{
		PayloadUUID: "payload-uuid-123",
	}

	if req.PayloadUUID != "payload-uuid-123" {
		t.Errorf("Expected PayloadUUID 'payload-uuid-123', got %q", req.PayloadUUID)
	}
}

// TestExportPayloadConfigResponseTypes tests the ExportPayloadConfigResponse type structure
func TestExportPayloadConfigResponseTypes(t *testing.T) {
	response := types.ExportPayloadConfigResponse{
		Config: `{"payload_type":"poseidon","os":"linux"}`,
	}

	if response.Config == "" {
		t.Error("Expected Config to be set")
	}
}
