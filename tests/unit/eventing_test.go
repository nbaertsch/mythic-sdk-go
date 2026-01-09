package unit

import (
	"strings"
	"testing"

	"github.com/nbaertsch/mythic-sdk-go/pkg/mythic/types"
)

func TestEventGroup_String(t *testing.T) {
	tests := []struct {
		name     string
		group    types.EventGroup
		contains []string
	}{
		{
			name: "active event group",
			group: types.EventGroup{
				ID:          1,
				Name:        "Auto Scan",
				TriggerType: "callback_new",
				Active:      true,
				Deleted:     false,
			},
			contains: []string{"Auto Scan", "active", "callback_new"},
		},
		{
			name: "inactive event group",
			group: types.EventGroup{
				ID:          2,
				Name:        "Test Workflow",
				TriggerType: "manual",
				Active:      false,
				Deleted:     false,
			},
			contains: []string{"Test Workflow", "inactive", "manual"},
		},
		{
			name: "deleted event group",
			group: types.EventGroup{
				ID:          3,
				Name:        "Old Workflow",
				TriggerType: "keyword",
				Active:      true,
				Deleted:     true,
			},
			contains: []string{"Old Workflow", "deleted", "keyword"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.group.String()

			for _, substr := range tt.contains {
				if !strings.Contains(result, substr) {
					t.Errorf("String() = %q, should contain %q", result, substr)
				}
			}
		})
	}
}

func TestEventGroup_IsActive(t *testing.T) {
	tests := []struct {
		name     string
		group    types.EventGroup
		expected bool
	}{
		{
			name: "active and not deleted",
			group: types.EventGroup{
				Active:  true,
				Deleted: false,
			},
			expected: true,
		},
		{
			name: "inactive",
			group: types.EventGroup{
				Active:  false,
				Deleted: false,
			},
			expected: false,
		},
		{
			name: "active but deleted",
			group: types.EventGroup{
				Active:  true,
				Deleted: true,
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.group.IsActive()
			if result != tt.expected {
				t.Errorf("IsActive() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestEventGroup_NeedsApproval(t *testing.T) {
	tests := []struct {
		name     string
		group    types.EventGroup
		expected bool
	}{
		{
			name: "requires approval and not approved",
			group: types.EventGroup{
				RequiresApproval: true,
				Approved:         false,
			},
			expected: true,
		},
		{
			name: "requires approval and approved",
			group: types.EventGroup{
				RequiresApproval: true,
				Approved:         true,
			},
			expected: false,
		},
		{
			name: "does not require approval",
			group: types.EventGroup{
				RequiresApproval: false,
				Approved:         false,
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.group.NeedsApproval()
			if result != tt.expected {
				t.Errorf("NeedsApproval() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestEventTriggerManualRequest_String(t *testing.T) {
	tests := []struct {
		name     string
		request  types.EventTriggerManualRequest
		contains []string
	}{
		{
			name: "with object ID",
			request: types.EventTriggerManualRequest{
				EventGroupID: 1,
				ObjectID:     42,
			},
			contains: []string{"1", "42"},
		},
		{
			name: "without object ID",
			request: types.EventTriggerManualRequest{
				EventGroupID: 2,
				ObjectID:     0,
			},
			contains: []string{"2", "Trigger"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.request.String()

			for _, substr := range tt.contains {
				if !strings.Contains(result, substr) {
					t.Errorf("String() = %q, should contain %q", result, substr)
				}
			}
		})
	}
}

func TestEventTriggerResponse_String(t *testing.T) {
	tests := []struct {
		name     string
		response types.EventTriggerResponse
		contains []string
	}{
		{
			name: "successful with execution ID",
			response: types.EventTriggerResponse{
				Status:      "success",
				ExecutionID: 123,
			},
			contains: []string{"successfully", "123"},
		},
		{
			name: "successful with message",
			response: types.EventTriggerResponse{
				Status:  "success",
				Message: "Event triggered for 5 objects",
			},
			contains: []string{"Event triggered for 5 objects"},
		},
		{
			name: "failed trigger",
			response: types.EventTriggerResponse{
				Status: "error",
				Error:  "Event group not found",
			},
			contains: []string{"Failed", "Event group not found"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.response.String()

			for _, substr := range tt.contains {
				if !strings.Contains(result, substr) {
					t.Errorf("String() = %q, should contain %q", result, substr)
				}
			}
		})
	}
}

func TestEventTriggerResponse_IsSuccessful(t *testing.T) {
	tests := []struct {
		name     string
		response types.EventTriggerResponse
		expected bool
	}{
		{
			name: "success status",
			response: types.EventTriggerResponse{
				Status: "success",
			},
			expected: true,
		},
		{
			name: "error status",
			response: types.EventTriggerResponse{
				Status: "error",
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.response.IsSuccessful()
			if result != tt.expected {
				t.Errorf("IsSuccessful() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestWorkflowExportResponse_String(t *testing.T) {
	tests := []struct {
		name     string
		response types.WorkflowExportResponse
		contains []string
	}{
		{
			name: "successful export",
			response: types.WorkflowExportResponse{
				Status:       "success",
				WorkflowName: "Auto Scanner",
			},
			contains: []string{"Exported", "Auto Scanner"},
		},
		{
			name: "failed export",
			response: types.WorkflowExportResponse{
				Status: "error",
				Error:  "Workflow not found",
			},
			contains: []string{"Failed", "Workflow not found"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.response.String()

			for _, substr := range tt.contains {
				if !strings.Contains(result, substr) {
					t.Errorf("String() = %q, should contain %q", result, substr)
				}
			}
		})
	}
}

func TestWorkflowImportResponse_String(t *testing.T) {
	tests := []struct {
		name     string
		response types.WorkflowImportResponse
		contains []string
	}{
		{
			name: "successful import with ID",
			response: types.WorkflowImportResponse{
				Status:     "success",
				WorkflowID: 42,
			},
			contains: []string{"successfully", "42"},
		},
		{
			name: "successful import with message",
			response: types.WorkflowImportResponse{
				Status:  "success",
				Message: "Workflow imported and activated",
			},
			contains: []string{"Workflow imported and activated"},
		},
		{
			name: "failed import",
			response: types.WorkflowImportResponse{
				Status: "error",
				Error:  "Invalid workflow format",
			},
			contains: []string{"Failed", "Invalid workflow format"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.response.String()

			for _, substr := range tt.contains {
				if !strings.Contains(result, substr) {
					t.Errorf("String() = %q, should contain %q", result, substr)
				}
			}
		})
	}
}

func TestWorkflowTestResponse_String(t *testing.T) {
	tests := []struct {
		name     string
		response types.WorkflowTestResponse
		contains []string
	}{
		{
			name: "valid workflow",
			response: types.WorkflowTestResponse{
				Status: "success",
				Valid:  true,
			},
			contains: []string{"valid"},
		},
		{
			name: "invalid workflow",
			response: types.WorkflowTestResponse{
				Status: "success",
				Valid:  false,
				Errors: []string{"error1", "error2", "error3"},
			},
			contains: []string{"failed", "3", "errors"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.response.String()

			for _, substr := range tt.contains {
				if !strings.Contains(result, substr) {
					t.Errorf("String() = %q, should contain %q", result, substr)
				}
			}
		})
	}
}

func TestWorkflowTestResponse_IsValid(t *testing.T) {
	tests := []struct {
		name     string
		response types.WorkflowTestResponse
		expected bool
	}{
		{
			name: "valid",
			response: types.WorkflowTestResponse{
				Valid: true,
			},
			expected: true,
		},
		{
			name: "invalid",
			response: types.WorkflowTestResponse{
				Valid: false,
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.response.IsValid()
			if result != tt.expected {
				t.Errorf("IsValid() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestWorkflowTestResponse_HasErrors(t *testing.T) {
	tests := []struct {
		name     string
		response types.WorkflowTestResponse
		expected bool
	}{
		{
			name: "has errors",
			response: types.WorkflowTestResponse{
				Errors: []string{"error1", "error2"},
			},
			expected: true,
		},
		{
			name: "no errors",
			response: types.WorkflowTestResponse{
				Errors: []string{},
			},
			expected: false,
		},
		{
			name: "nil errors",
			response: types.WorkflowTestResponse{
				Errors: nil,
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.response.HasErrors()
			if result != tt.expected {
				t.Errorf("HasErrors() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestWebhookResponse_String(t *testing.T) {
	tests := []struct {
		name     string
		response types.WebhookResponse
		contains []string
	}{
		{
			name: "successful webhook",
			response: types.WebhookResponse{
				Status:     "success",
				StatusCode: 200,
			},
			contains: []string{"successfully", "200"},
		},
		{
			name: "failed webhook",
			response: types.WebhookResponse{
				Status: "error",
				Error:  "Connection timeout",
			},
			contains: []string{"failed", "Connection timeout"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.response.String()

			for _, substr := range tt.contains {
				if !strings.Contains(result, substr) {
					t.Errorf("String() = %q, should contain %q", result, substr)
				}
			}
		})
	}
}

func TestConsumingServiceTestResponse_String(t *testing.T) {
	tests := []struct {
		name     string
		response types.ConsumingServiceTestResponse
		contains []string
	}{
		{
			name: "successful test with message",
			response: types.ConsumingServiceTestResponse{
				Status:  "success",
				Message: "Service responded correctly",
			},
			contains: []string{"Service responded correctly"},
		},
		{
			name: "successful test without message",
			response: types.ConsumingServiceTestResponse{
				Status: "success",
			},
			contains: []string{"passed"},
		},
		{
			name: "failed test",
			response: types.ConsumingServiceTestResponse{
				Status: "error",
				Error:  "Service unreachable",
			},
			contains: []string{"failed", "Service unreachable"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.response.String()

			for _, substr := range tt.contains {
				if !strings.Contains(result, substr) {
					t.Errorf("String() = %q, should contain %q", result, substr)
				}
			}
		})
	}
}

func TestEventApprovalRequest_String(t *testing.T) {
	tests := []struct {
		name     string
		request  types.EventApprovalRequest
		contains []string
	}{
		{
			name: "approve",
			request: types.EventApprovalRequest{
				EventGroupID: 1,
				Approved:     true,
			},
			contains: []string{"Approve", "1"},
		},
		{
			name: "reject",
			request: types.EventApprovalRequest{
				EventGroupID: 2,
				Approved:     false,
			},
			contains: []string{"Reject", "2"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.request.String()

			for _, substr := range tt.contains {
				if !strings.Contains(result, substr) {
					t.Errorf("String() = %q, should contain %q", result, substr)
				}
			}
		})
	}
}

func TestEventTriggerBulkRequest_String(t *testing.T) {
	request := types.EventTriggerBulkRequest{
		EventGroupID: 1,
		ObjectIDs:    []int{10, 20, 30, 40, 50},
	}

	result := request.String()

	if !strings.Contains(result, "1") {
		t.Errorf("String() should contain event group ID")
	}
	if !strings.Contains(result, "5") {
		t.Errorf("String() should contain object count")
	}
}

func TestEventTriggerKeywordRequest_String(t *testing.T) {
	request := types.EventTriggerKeywordRequest{
		Keyword: "scan",
	}

	result := request.String()

	if !strings.Contains(result, "scan") {
		t.Errorf("String() should contain keyword")
	}
}

func TestWebhookRequest_String(t *testing.T) {
	request := types.WebhookRequest{
		WebhookURL: "https://example.com/webhook",
		Method:     "POST",
	}

	result := request.String()

	if !strings.Contains(result, "POST") {
		t.Errorf("String() should contain method")
	}
	if !strings.Contains(result, "example.com") {
		t.Errorf("String() should contain URL")
	}
}
