package unit

import (
	"testing"
	"time"

	"github.com/nbaertsch/mythic-sdk-go/pkg/mythic"
)

func TestTask_String(t *testing.T) {
	task := &mythic.Task{
		DisplayID:     1,
		CommandName:   "shell",
		DisplayParams: "whoami",
		Status:        "processed",
		Completed:     true,
	}

	expected := "Task 1: shell whoami (Status: processed, Completed: true)"
	if task.String() != expected {
		t.Errorf("Expected %q, got %q", expected, task.String())
	}
}

func TestTask_IsCompleted(t *testing.T) {
	tests := []struct {
		name      string
		completed bool
		expected  bool
	}{
		{"Completed task", true, true},
		{"Incomplete task", false, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			task := &mythic.Task{Completed: tt.completed}
			if got := task.IsCompleted(); got != tt.expected {
				t.Errorf("IsCompleted() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestTask_IsError(t *testing.T) {
	tests := []struct {
		name     string
		status   string
		expected bool
	}{
		{"Error task", "error", true},
		{"Processed task", "processed", false},
		{"Submitted task", "submitted", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			task := &mythic.Task{Status: tt.status}
			if got := task.IsError(); got != tt.expected {
				t.Errorf("IsError() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestTask_HasOutput(t *testing.T) {
	tests := []struct {
		name          string
		responseCount int
		expected      bool
	}{
		{"Task with output", 5, true},
		{"Task without output", 0, false},
		{"Task with one response", 1, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			task := &mythic.Task{ResponseCount: tt.responseCount}
			if got := task.HasOutput(); got != tt.expected {
				t.Errorf("HasOutput() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestTaskRequest_Validation(t *testing.T) {
	tests := []struct {
		name        string
		request     *mythic.TaskRequest
		shouldError bool
	}{
		{
			name: "Valid request with CallbackID",
			request: &mythic.TaskRequest{
				CallbackID: intPtr(1),
				Command:    "shell",
				Params:     "whoami",
			},
			shouldError: false,
		},
		{
			name: "Valid request with CallbackIDs",
			request: &mythic.TaskRequest{
				CallbackIDs: []int{1, 2, 3},
				Command:     "shell",
				Params:      "whoami",
			},
			shouldError: false,
		},
		{
			name: "Missing callback ID",
			request: &mythic.TaskRequest{
				Command: "shell",
				Params:  "whoami",
			},
			shouldError: true,
		},
		{
			name: "Missing command",
			request: &mythic.TaskRequest{
				CallbackID: intPtr(1),
				Params:     "whoami",
			},
			shouldError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hasCallbackID := tt.request.CallbackID != nil || len(tt.request.CallbackIDs) > 0
			hasCommand := tt.request.Command != ""

			if tt.shouldError {
				if hasCallbackID && hasCommand {
					t.Error("Expected validation to fail, but request appears valid")
				}
			} else {
				if !hasCallbackID {
					t.Error("Request should have callback ID")
				}
				if !hasCommand {
					t.Error("Request should have command")
				}
			}
		})
	}
}

func TestTaskStatus_Constants(t *testing.T) {
	tests := []struct {
		name   string
		status mythic.TaskStatus
		want   string
	}{
		{"Preprocessing", mythic.TaskStatusPreprocessing, "preprocessing"},
		{"Submitted", mythic.TaskStatusSubmitted, "submitted"},
		{"Processing", mythic.TaskStatusProcessing, "processing"},
		{"Processed", mythic.TaskStatusProcessed, "processed"},
		{"Error", mythic.TaskStatusError, "error"},
		{"Cleared", mythic.TaskStatusCleared, "cleared"},
		{"Completed", mythic.TaskStatusCompleted, "completed"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if string(tt.status) != tt.want {
				t.Errorf("TaskStatus %s = %v, want %v", tt.name, tt.status, tt.want)
			}
		})
	}
}

func TestTaskResponse_Structure(t *testing.T) {
	now := time.Now()
	seqNum := 1

	response := &mythic.TaskResponse{
		ID:             100,
		TaskID:         50,
		ResponseText:   "test output",
		ResponseRaw:    []byte("raw output"),
		IsError:        false,
		Timestamp:      now,
		SequenceNumber: &seqNum,
	}

	// Verify all fields are set correctly
	if response.ID != 100 {
		t.Errorf("Expected ID 100, got %d", response.ID)
	}
	if response.TaskID != 50 {
		t.Errorf("Expected TaskID 50, got %d", response.TaskID)
	}
	if response.ResponseText != "test output" {
		t.Errorf("Expected ResponseText 'test output', got %q", response.ResponseText)
	}
	if string(response.ResponseRaw) != "raw output" {
		t.Errorf("Expected ResponseRaw 'raw output', got %q", string(response.ResponseRaw))
	}
	if response.IsError {
		t.Error("Expected IsError false, got true")
	}
	if !response.Timestamp.Equal(now) {
		t.Errorf("Expected Timestamp %v, got %v", now, response.Timestamp)
	}
	if response.SequenceNumber == nil || *response.SequenceNumber != 1 {
		t.Errorf("Expected SequenceNumber 1, got %v", response.SequenceNumber)
	}
}

func TestTask_CompleteStructure(t *testing.T) {
	now := time.Now()
	parentID := 5
	interactiveType := 1
	opsecBlocked := true

	task := &mythic.Task{
		ID:                        100,
		DisplayID:                 1,
		AgentTaskID:               "abc-123",
		CommandName:               "shell",
		Params:                    "{\"command\":\"whoami\"}",
		DisplayParams:             "whoami",
		OriginalParams:            "whoami",
		Status:                    "processed",
		Completed:                 true,
		Comment:                   "Test task",
		Timestamp:                 now,
		CallbackID:                10,
		OperatorID:                1,
		OperationID:               1,
		ParentTaskID:              &parentID,
		ResponseCount:             3,
		IsInteractiveTask:         true,
		InteractiveTaskType:       &interactiveType,
		TaskingLocation:           "command_line",
		ParameterGroupName:        "Default",
		Stdout:                    "output",
		Stderr:                    "",
		CompletedCallbackFunction: "test_callback",
		SubtaskCallbackFunction:   "subtask_callback",
		GroupCallbackFunction:     "group_callback",
		OpsecPreBlocked:           &opsecBlocked,
		OpsecPreBypassed:          false,
		OpsecPreMessage:           "Pre-execution check",
		OpsecPostBlocked:          &opsecBlocked,
		OpsecPostBypassed:         false,
		OpsecPostMessage:          "Post-execution check",
	}

	// Verify critical fields
	if task.ID != 100 {
		t.Errorf("Expected ID 100, got %d", task.ID)
	}
	if task.DisplayID != 1 {
		t.Errorf("Expected DisplayID 1, got %d", task.DisplayID)
	}
	if task.CommandName != "shell" {
		t.Errorf("Expected CommandName 'shell', got %q", task.CommandName)
	}
	if !task.Completed {
		t.Error("Expected Completed true, got false")
	}
	if task.ParentTaskID == nil || *task.ParentTaskID != 5 {
		t.Errorf("Expected ParentTaskID 5, got %v", task.ParentTaskID)
	}
	if task.InteractiveTaskType == nil || *task.InteractiveTaskType != 1 {
		t.Errorf("Expected InteractiveTaskType 1, got %v", task.InteractiveTaskType)
	}
	if task.OpsecPreBlocked == nil || !*task.OpsecPreBlocked {
		t.Errorf("Expected OpsecPreBlocked true, got %v", task.OpsecPreBlocked)
	}
}

func TestTaskRequest_OptionalFields(t *testing.T) {
	callbackID := 1
	parentTaskID := 5
	interactiveType := 2
	tokenID := 10

	req := &mythic.TaskRequest{
		CallbackID:          &callbackID,
		CallbackIDs:         []int{1, 2, 3},
		Command:             "shell",
		Params:              "whoami",
		Files:               []string{"file1", "file2"},
		IsInteractiveTask:   true,
		InteractiveTaskType: &interactiveType,
		ParentTaskID:        &parentTaskID,
		TaskingLocation:     "scripting",
		ParameterGroupName:  "Custom",
		OriginalParams:      "original",
		TokenID:             &tokenID,
	}

	// Verify optional fields
	if req.CallbackID == nil || *req.CallbackID != 1 {
		t.Errorf("Expected CallbackID 1, got %v", req.CallbackID)
	}
	if len(req.CallbackIDs) != 3 {
		t.Errorf("Expected 3 CallbackIDs, got %d", len(req.CallbackIDs))
	}
	if len(req.Files) != 2 {
		t.Errorf("Expected 2 Files, got %d", len(req.Files))
	}
	if req.ParentTaskID == nil || *req.ParentTaskID != 5 {
		t.Errorf("Expected ParentTaskID 5, got %v", req.ParentTaskID)
	}
	if req.InteractiveTaskType == nil || *req.InteractiveTaskType != 2 {
		t.Errorf("Expected InteractiveTaskType 2, got %v", req.InteractiveTaskType)
	}
	if req.TokenID == nil || *req.TokenID != 10 {
		t.Errorf("Expected TokenID 10, got %v", req.TokenID)
	}
}

func TestTaskRequest_PayloadType(t *testing.T) {
	callbackID := 2
	pt := "forge"

	req := &mythic.TaskRequest{
		CallbackID:  &callbackID,
		Command:     "forge_collections",
		Params:      "",
		PayloadType: &pt,
	}

	if req.PayloadType == nil {
		t.Fatal("Expected PayloadType to be set")
	}
	if *req.PayloadType != "forge" {
		t.Errorf("Expected PayloadType 'forge', got %q", *req.PayloadType)
	}
	if req.Command != "forge_collections" {
		t.Errorf("Expected Command 'forge_collections', got %q", req.Command)
	}
}

func TestTaskRequest_PayloadTypeNil(t *testing.T) {
	callbackID := 1

	req := &mythic.TaskRequest{
		CallbackID: &callbackID,
		Command:    "shell",
		Params:     "whoami",
	}

	if req.PayloadType != nil {
		t.Errorf("Expected PayloadType nil for normal task, got %q", *req.PayloadType)
	}
}

func TestScriptOnlyTaskRequest(t *testing.T) {
	req := &mythic.ScriptOnlyTaskRequest{
		CallbackID:         2,
		Command:            "forge_collections",
		Params:             "{\"collection_name\": \"test\"}",
		TaskingLocation:    "command_line",
		ParameterGroupName: "Default",
	}

	if req.CallbackID != 2 {
		t.Errorf("Expected CallbackID 2, got %d", req.CallbackID)
	}
	if req.Command != "forge_collections" {
		t.Errorf("Expected Command 'forge_collections', got %q", req.Command)
	}
	if req.Params != "{\"collection_name\": \"test\"}" {
		t.Errorf("Expected Params JSON, got %q", req.Params)
	}
	if req.TaskingLocation != "command_line" {
		t.Errorf("Expected TaskingLocation 'command_line', got %q", req.TaskingLocation)
	}
	if req.ParameterGroupName != "Default" {
		t.Errorf("Expected ParameterGroupName 'Default', got %q", req.ParameterGroupName)
	}
}

func TestScriptOnlyTaskRequest_MinimalFields(t *testing.T) {
	req := &mythic.ScriptOnlyTaskRequest{
		CallbackID: 5,
		Command:    "forge_download",
	}

	if req.CallbackID != 5 {
		t.Errorf("Expected CallbackID 5, got %d", req.CallbackID)
	}
	if req.Command != "forge_download" {
		t.Errorf("Expected Command 'forge_download', got %q", req.Command)
	}
	if req.Params != "" {
		t.Errorf("Expected empty Params, got %q", req.Params)
	}
	if req.TaskingLocation != "" {
		t.Errorf("Expected empty TaskingLocation, got %q", req.TaskingLocation)
	}
	if req.ParameterGroupName != "" {
		t.Errorf("Expected empty ParameterGroupName, got %q", req.ParameterGroupName)
	}
}

// Helper function to create int pointers
func intPtr(i int) *int {
	return &i
}

func TestTaskArtifact_Structure(t *testing.T) {
	now := time.Now()

	artifact := &mythic.TaskArtifact{
		ID:           1,
		TaskID:       100,
		Artifact:     "C:\\Windows\\Temp\\malware.exe",
		BaseArtifact: "C:\\Windows\\Temp",
		Host:         "WORKSTATION01",
		Timestamp:    now,
	}

	if artifact.ID != 1 {
		t.Errorf("Expected ID 1, got %d", artifact.ID)
	}

	if artifact.TaskID != 100 {
		t.Errorf("Expected TaskID 100, got %d", artifact.TaskID)
	}

	if artifact.Artifact != "C:\\Windows\\Temp\\malware.exe" {
		t.Errorf("Expected Artifact path, got %q", artifact.Artifact)
	}

	if artifact.BaseArtifact != "C:\\Windows\\Temp" {
		t.Errorf("Expected BaseArtifact, got %q", artifact.BaseArtifact)
	}

	if artifact.Host != "WORKSTATION01" {
		t.Errorf("Expected Host 'WORKSTATION01', got %q", artifact.Host)
	}

	if !artifact.Timestamp.Equal(now) {
		t.Errorf("Expected Timestamp %v, got %v", now, artifact.Timestamp)
	}
}
