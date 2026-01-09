package unit

import (
	"strings"
	"testing"

	"github.com/nbaertsch/mythic-sdk-go/pkg/mythic/types"
)

func TestDynamicQueryRequest_String(t *testing.T) {
	tests := []struct {
		name     string
		request  types.DynamicQueryRequest
		contains []string
	}{
		{
			name: "with callback ID",
			request: types.DynamicQueryRequest{
				Command:    "ls",
				CallbackID: 42,
				Parameters: map[string]interface{}{"path": "/home"},
			},
			contains: []string{"ls", "callback", "42"},
		},
		{
			name: "without callback ID",
			request: types.DynamicQueryRequest{
				Command:    "download",
				CallbackID: 0,
				Parameters: map[string]interface{}{"file": "test.txt"},
			},
			contains: []string{"download"},
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

func TestDynamicQueryResponse_String(t *testing.T) {
	tests := []struct {
		name     string
		response types.DynamicQueryResponse
		contains []string
	}{
		{
			name: "successful with choices",
			response: types.DynamicQueryResponse{
				Status:  "success",
				Choices: []interface{}{"choice1", "choice2", "choice3"},
			},
			contains: []string{"3", "choices"},
		},
		{
			name: "successful with no choices",
			response: types.DynamicQueryResponse{
				Status:  "success",
				Choices: []interface{}{},
			},
			contains: []string{"0", "choices"},
		},
		{
			name: "failed query",
			response: types.DynamicQueryResponse{
				Status: "error",
				Error:  "Command not found",
			},
			contains: []string{"failed", "Command not found"},
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

func TestDynamicQueryResponse_IsSuccessful(t *testing.T) {
	tests := []struct {
		name     string
		response types.DynamicQueryResponse
		expected bool
	}{
		{
			name: "success status",
			response: types.DynamicQueryResponse{
				Status: "success",
			},
			expected: true,
		},
		{
			name: "error status",
			response: types.DynamicQueryResponse{
				Status: "error",
			},
			expected: false,
		},
		{
			name: "failed status",
			response: types.DynamicQueryResponse{
				Status: "failed",
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

func TestDynamicQueryResponse_HasChoices(t *testing.T) {
	tests := []struct {
		name     string
		response types.DynamicQueryResponse
		expected bool
	}{
		{
			name: "has choices",
			response: types.DynamicQueryResponse{
				Choices: []interface{}{"choice1", "choice2"},
			},
			expected: true,
		},
		{
			name: "empty choices",
			response: types.DynamicQueryResponse{
				Choices: []interface{}{},
			},
			expected: false,
		},
		{
			name: "nil choices",
			response: types.DynamicQueryResponse{
				Choices: nil,
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.response.HasChoices()
			if result != tt.expected {
				t.Errorf("HasChoices() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestDynamicBuildParameterRequest_String(t *testing.T) {
	request := types.DynamicBuildParameterRequest{
		PayloadType: "apollo",
		Parameter:   "c2_profile",
		Parameters:  map[string]interface{}{"operation": "test-op"},
	}

	result := request.String()

	if !strings.Contains(result, "apollo") {
		t.Errorf("String() = %q, should contain 'apollo'", result)
	}
	if !strings.Contains(result, "c2_profile") {
		t.Errorf("String() = %q, should contain 'c2_profile'", result)
	}
}

func TestDynamicBuildParameterResponse_String(t *testing.T) {
	tests := []struct {
		name     string
		response types.DynamicBuildParameterResponse
		contains []string
	}{
		{
			name: "successful with choices",
			response: types.DynamicBuildParameterResponse{
				Status:  "success",
				Choices: []interface{}{"http", "https", "dns"},
			},
			contains: []string{"3", "choices"},
		},
		{
			name: "failed query",
			response: types.DynamicBuildParameterResponse{
				Status: "error",
				Error:  "Payload type not found",
			},
			contains: []string{"failed", "Payload type not found"},
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

func TestDynamicBuildParameterResponse_IsSuccessful(t *testing.T) {
	tests := []struct {
		name     string
		response types.DynamicBuildParameterResponse
		expected bool
	}{
		{
			name: "success status",
			response: types.DynamicBuildParameterResponse{
				Status: "success",
			},
			expected: true,
		},
		{
			name: "error status",
			response: types.DynamicBuildParameterResponse{
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

func TestDynamicBuildParameterResponse_HasChoices(t *testing.T) {
	tests := []struct {
		name     string
		response types.DynamicBuildParameterResponse
		expected bool
	}{
		{
			name: "has choices",
			response: types.DynamicBuildParameterResponse{
				Choices: []interface{}{"http", "https"},
			},
			expected: true,
		},
		{
			name: "empty choices",
			response: types.DynamicBuildParameterResponse{
				Choices: []interface{}{},
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.response.HasChoices()
			if result != tt.expected {
				t.Errorf("HasChoices() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestTypedArrayParseRequest_String(t *testing.T) {
	request := types.TypedArrayParseRequest{
		InputArray:    "/etc/passwd:read,/etc/shadow:read",
		ParameterType: "file_list",
	}

	result := request.String()

	if !strings.Contains(result, "file_list") {
		t.Errorf("String() = %q, should contain 'file_list'", result)
	}
}

func TestTypedArrayParseResponse_String(t *testing.T) {
	tests := []struct {
		name     string
		response types.TypedArrayParseResponse
		contains []string
	}{
		{
			name: "successful parse",
			response: types.TypedArrayParseResponse{
				Status: "success",
				ParsedArray: []interface{}{
					map[string]interface{}{"path": "/etc/passwd", "permission": "read"},
					map[string]interface{}{"path": "/etc/shadow", "permission": "read"},
				},
			},
			contains: []string{"Parsed", "2", "elements"},
		},
		{
			name: "empty array",
			response: types.TypedArrayParseResponse{
				Status:      "success",
				ParsedArray: []interface{}{},
			},
			contains: []string{"0", "elements"},
		},
		{
			name: "parse error",
			response: types.TypedArrayParseResponse{
				Status: "error",
				Error:  "Invalid format",
			},
			contains: []string{"failed", "Invalid format"},
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

func TestTypedArrayParseResponse_IsSuccessful(t *testing.T) {
	tests := []struct {
		name     string
		response types.TypedArrayParseResponse
		expected bool
	}{
		{
			name: "success status",
			response: types.TypedArrayParseResponse{
				Status: "success",
			},
			expected: true,
		},
		{
			name: "error status",
			response: types.TypedArrayParseResponse{
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

func TestTypedArrayParseResponse_HasElements(t *testing.T) {
	tests := []struct {
		name     string
		response types.TypedArrayParseResponse
		expected bool
	}{
		{
			name: "has elements",
			response: types.TypedArrayParseResponse{
				ParsedArray: []interface{}{
					map[string]interface{}{"key": "value"},
				},
			},
			expected: true,
		},
		{
			name: "empty array",
			response: types.TypedArrayParseResponse{
				ParsedArray: []interface{}{},
			},
			expected: false,
		},
		{
			name: "nil array",
			response: types.TypedArrayParseResponse{
				ParsedArray: nil,
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.response.HasElements()
			if result != tt.expected {
				t.Errorf("HasElements() = %v, want %v", result, tt.expected)
			}
		})
	}
}
