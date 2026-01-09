package unit

import (
	"strings"
	"testing"

	"github.com/nbaertsch/mythic-sdk-go/pkg/mythic/types"
)

func TestCreateRandomRequest_String(t *testing.T) {
	tests := []struct {
		name     string
		request  types.CreateRandomRequest
		contains []string
	}{
		{
			name: "format only",
			request: types.CreateRandomRequest{
				Format: "%s-%d",
				Length: 0,
			},
			contains: []string{"Generate", "random", "%s-%d"},
		},
		{
			name: "format with length",
			request: types.CreateRandomRequest{
				Format: "%x",
				Length: 16,
			},
			contains: []string{"Generate", "random", "%x", "16"},
		},
		{
			name: "complex format",
			request: types.CreateRandomRequest{
				Format: "%S%d%x",
				Length: 8,
			},
			contains: []string{"Generate", "%S%d%x", "8"},
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

func TestCreateRandomResponse_String(t *testing.T) {
	tests := []struct {
		name     string
		response types.CreateRandomResponse
		contains []string
	}{
		{
			name: "successful generation",
			response: types.CreateRandomResponse{
				Status:       "success",
				RandomString: "abc123",
			},
			contains: []string{"Generated", "abc123"},
		},
		{
			name: "failed generation",
			response: types.CreateRandomResponse{
				Status: "error",
				Error:  "Invalid format",
			},
			contains: []string{"Failed", "Invalid format"},
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

func TestCreateRandomResponse_IsSuccessful(t *testing.T) {
	tests := []struct {
		name     string
		response types.CreateRandomResponse
		expected bool
	}{
		{
			name: "success with string",
			response: types.CreateRandomResponse{
				Status:       "success",
				RandomString: "abc123",
			},
			expected: true,
		},
		{
			name: "success but empty string",
			response: types.CreateRandomResponse{
				Status:       "success",
				RandomString: "",
			},
			expected: false,
		},
		{
			name: "error status",
			response: types.CreateRandomResponse{
				Status:       "error",
				RandomString: "",
			},
			expected: false,
		},
		{
			name: "failed status",
			response: types.CreateRandomResponse{
				Status:       "failed",
				RandomString: "something",
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

func TestConfigCheckResponse_String(t *testing.T) {
	tests := []struct {
		name     string
		response types.ConfigCheckResponse
		contains []string
	}{
		{
			name: "valid configuration",
			response: types.ConfigCheckResponse{
				Valid:  true,
				Errors: []string{},
			},
			contains: []string{"Configuration", "valid"},
		},
		{
			name: "configuration with errors",
			response: types.ConfigCheckResponse{
				Valid:  false,
				Errors: []string{"Database not connected", "Redis unavailable"},
			},
			contains: []string{"Configuration", "2 error"},
		},
		{
			name: "configuration with message",
			response: types.ConfigCheckResponse{
				Valid:   false,
				Message: "Missing environment variables",
			},
			contains: []string{"Configuration", "Missing environment"},
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

func TestConfigCheckResponse_IsValid(t *testing.T) {
	tests := []struct {
		name     string
		response types.ConfigCheckResponse
		expected bool
	}{
		{
			name: "valid with no errors",
			response: types.ConfigCheckResponse{
				Valid:  true,
				Errors: []string{},
			},
			expected: true,
		},
		{
			name: "valid flag but has errors",
			response: types.ConfigCheckResponse{
				Valid:  true,
				Errors: []string{"Some error"},
			},
			expected: false,
		},
		{
			name: "invalid",
			response: types.ConfigCheckResponse{
				Valid:  false,
				Errors: []string{},
			},
			expected: false,
		},
		{
			name: "invalid with errors",
			response: types.ConfigCheckResponse{
				Valid:  false,
				Errors: []string{"Error 1", "Error 2"},
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

func TestConfigCheckResponse_HasErrors(t *testing.T) {
	tests := []struct {
		name     string
		response types.ConfigCheckResponse
		expected bool
	}{
		{
			name: "no errors",
			response: types.ConfigCheckResponse{
				Errors: []string{},
			},
			expected: false,
		},
		{
			name: "nil errors",
			response: types.ConfigCheckResponse{
				Errors: nil,
			},
			expected: false,
		},
		{
			name: "has errors",
			response: types.ConfigCheckResponse{
				Errors: []string{"Error 1"},
			},
			expected: true,
		},
		{
			name: "multiple errors",
			response: types.ConfigCheckResponse{
				Errors: []string{"Error 1", "Error 2", "Error 3"},
			},
			expected: true,
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

func TestConfigCheckResponse_GetErrors(t *testing.T) {
	tests := []struct {
		name     string
		response types.ConfigCheckResponse
		expected []string
	}{
		{
			name: "no errors",
			response: types.ConfigCheckResponse{
				Errors: []string{},
			},
			expected: []string{},
		},
		{
			name: "single error",
			response: types.ConfigCheckResponse{
				Errors: []string{"Database connection failed"},
			},
			expected: []string{"Database connection failed"},
		},
		{
			name: "multiple errors",
			response: types.ConfigCheckResponse{
				Errors: []string{"Error 1", "Error 2"},
			},
			expected: []string{"Error 1", "Error 2"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.response.GetErrors()

			if len(result) != len(tt.expected) {
				t.Errorf("GetErrors() returned %d errors, want %d", len(result), len(tt.expected))
				return
			}

			for i, err := range result {
				if err != tt.expected[i] {
					t.Errorf("GetErrors()[%d] = %q, want %q", i, err, tt.expected[i])
				}
			}
		})
	}
}

func TestConfigCheckResponse_CompleteStructure(t *testing.T) {
	config := map[string]interface{}{
		"database":  "postgresql",
		"rabbitmq":  "connected",
		"containers": 5,
	}

	response := types.ConfigCheckResponse{
		Status:  "success",
		Valid:   true,
		Errors:  []string{},
		Config:  config,
		Message: "All systems operational",
	}

	// Test all fields
	if response.Status != "success" {
		t.Errorf("Status = %s, want success", response.Status)
	}
	if !response.Valid {
		t.Error("Valid should be true")
	}
	if len(response.Errors) != 0 {
		t.Errorf("Errors should be empty, got %v", response.Errors)
	}
	if response.Message != "All systems operational" {
		t.Errorf("Message = %s, want 'All systems operational'", response.Message)
	}

	// Test helper methods
	if !response.IsValid() {
		t.Error("IsValid() should return true")
	}
	if response.HasErrors() {
		t.Error("HasErrors() should return false")
	}

	errors := response.GetErrors()
	if len(errors) != 0 {
		t.Errorf("GetErrors() should return empty slice, got %v", errors)
	}

	str := response.String()
	if !strings.Contains(str, "valid") {
		t.Errorf("String() = %q, should contain 'valid'", str)
	}
}
