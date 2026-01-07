package unit

import (
	"testing"
	"time"

	"github.com/nbaertsch/mythic-sdk-go/pkg/mythic/types"
)

// TestTokenString tests the Token.String() method
func TestTokenString(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name     string
		token    types.Token
		contains []string
	}{
		{
			name: "with user and host",
			token: types.Token{
				ID:        1,
				User:      "DOMAIN\\Administrator",
				Host:      "SERVER01",
				Timestamp: now,
			},
			contains: []string{"DOMAIN\\Administrator", "SERVER01"},
		},
		{
			name: "with user only",
			token: types.Token{
				ID:        2,
				User:      "admin",
				Timestamp: now,
			},
			contains: []string{"admin"},
		},
		{
			name: "with TokenID only",
			token: types.Token{
				ID:        3,
				TokenID:   "abc123",
				Timestamp: now,
			},
			contains: []string{"Token abc123"},
		},
		{
			name: "with ID only",
			token: types.Token{
				ID:        4,
				Timestamp: now,
			},
			contains: []string{"Token 4"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.token.String()
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

// TestTokenIsDeleted tests the Token.IsDeleted() method
func TestTokenIsDeleted(t *testing.T) {
	tests := []struct {
		name     string
		deleted  bool
		expected bool
	}{
		{"deleted token", true, true},
		{"active token", false, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token := types.Token{Deleted: tt.deleted}
			result := token.IsDeleted()
			if result != tt.expected {
				t.Errorf("IsDeleted() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

// TestTokenHasTask tests the Token.HasTask() method
func TestTokenHasTask(t *testing.T) {
	taskID := 42
	zeroTaskID := 0

	tests := []struct {
		name     string
		taskID   *int
		expected bool
	}{
		{"with task", &taskID, true},
		{"without task", nil, false},
		{"with zero task ID", &zeroTaskID, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token := types.Token{TaskID: tt.taskID}
			result := token.HasTask()
			if result != tt.expected {
				t.Errorf("HasTask() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

// TestTokenGetIntegrityLevelString tests integrity level conversion
func TestTokenGetIntegrityLevelString(t *testing.T) {
	tests := []struct {
		level    int
		expected string
	}{
		{0, "Untrusted"},
		{1, "Low"},
		{2, "Medium"},
		{3, "High"},
		{4, "System"},
		{99, "Unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			token := types.Token{IntegrityLevel: tt.level}
			result := token.GetIntegrityLevelString()
			if result != tt.expected {
				t.Errorf("GetIntegrityLevelString() = %q, expected %q", result, tt.expected)
			}
		})
	}
}

// TestTokenTypes tests the Token type structure
func TestTokenTypes(t *testing.T) {
	now := time.Now()
	taskID := 42

	token := types.Token{
		ID:              1,
		TokenID:         "token123",
		User:            "WORKSTATION\\user",
		Groups:          "Users, Administrators",
		Privileges:      "SeDebugPrivilege",
		ThreadID:        1234,
		ProcessID:       5678,
		SessionID:       1,
		LogonSID:        "S-1-5-21-...",
		IntegrityLevel:  3,
		Restricted:      false,
		DefaultDACL:     "D:(A;;GA;;;BA)",
		Handle:          "0x1234",
		Capabilities:    "cap1,cap2",
		AppContainerSID: "S-1-15-...",
		AppContainerNum: 1,
		TaskID:          &taskID,
		OperationID:     5,
		Timestamp:       now,
		Host:            "SERVER01",
		Deleted:         false,
	}

	if token.ID != 1 {
		t.Errorf("Expected ID 1, got %d", token.ID)
	}
	if token.User != "WORKSTATION\\user" {
		t.Errorf("Expected User 'WORKSTATION\\user', got %q", token.User)
	}
	if !token.HasTask() {
		t.Error("Expected token to have task")
	}
	if token.IsDeleted() {
		t.Error("Expected token to not be deleted")
	}
	if token.GetIntegrityLevelString() != "High" {
		t.Errorf("Expected integrity level 'High', got %q", token.GetIntegrityLevelString())
	}
}

// TestCallbackTokenString tests the CallbackToken.String() method
func TestCallbackTokenString(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name          string
		callbackToken types.CallbackToken
		contains      []string
	}{
		{
			name: "with token loaded",
			callbackToken: types.CallbackToken{
				ID:         1,
				CallbackID: 10,
				TokenID:    5,
				Timestamp:  now,
				Token: &types.Token{
					User: "admin",
					Host: "SERVER01",
				},
			},
			contains: []string{"admin", "SERVER01", "10"},
		},
		{
			name: "without loaded references",
			callbackToken: types.CallbackToken{
				ID:         2,
				CallbackID: 20,
				TokenID:    10,
				Timestamp:  now,
			},
			contains: []string{"CallbackToken 2"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.callbackToken.String()
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

// TestCallbackTokenTypes tests the CallbackToken type structure
func TestCallbackTokenTypes(t *testing.T) {
	now := time.Now()

	ct := types.CallbackToken{
		ID:         1,
		CallbackID: 10,
		TokenID:    5,
		Timestamp:  now,
	}

	if ct.ID != 1 {
		t.Errorf("Expected ID 1, got %d", ct.ID)
	}
	if ct.CallbackID != 10 {
		t.Errorf("Expected CallbackID 10, got %d", ct.CallbackID)
	}
	if ct.TokenID != 5 {
		t.Errorf("Expected TokenID 5, got %d", ct.TokenID)
	}
}

// TestAPITokenString tests the APIToken.String() method
func TestAPITokenString(t *testing.T) {
	tests := []struct {
		name     string
		apiToken types.APIToken
		contains []string
	}{
		{
			name: "with name",
			apiToken: types.APIToken{
				ID:        1,
				Name:      "MyToken",
				TokenType: "User",
			},
			contains: []string{"MyToken", "User"},
		},
		{
			name: "without name",
			apiToken: types.APIToken{
				ID:        2,
				TokenType: "C2",
			},
			contains: []string{"APIToken 2", "C2"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.apiToken.String()
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

// TestAPITokenIsActive tests the APIToken.IsActive() method
func TestAPITokenIsActive(t *testing.T) {
	tests := []struct {
		name     string
		active   bool
		deleted  bool
		expected bool
	}{
		{"active and not deleted", true, false, true},
		{"inactive", false, false, false},
		{"deleted", true, true, false},
		{"inactive and deleted", false, true, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token := types.APIToken{Active: tt.active, Deleted: tt.deleted}
			result := token.IsActive()
			if result != tt.expected {
				t.Errorf("IsActive() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

// TestAPITokenIsDeleted tests the APIToken.IsDeleted() method
func TestAPITokenIsDeleted(t *testing.T) {
	tests := []struct {
		name     string
		deleted  bool
		expected bool
	}{
		{"deleted token", true, true},
		{"active token", false, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token := types.APIToken{Deleted: tt.deleted}
			result := token.IsDeleted()
			if result != tt.expected {
				t.Errorf("IsDeleted() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

// TestAPITokenTypes tests the APIToken type structure
func TestAPITokenTypes(t *testing.T) {
	now := time.Now()
	operationID := 5

	token := types.APIToken{
		ID:           1,
		TokenValue:   "abc123...",
		TokenType:    "User",
		Active:       true,
		CreationTime: now,
		OperatorID:   10,
		OperationID:  &operationID,
		Name:         "TestToken",
		Deleted:      false,
	}

	if token.ID != 1 {
		t.Errorf("Expected ID 1, got %d", token.ID)
	}
	if token.TokenType != "User" {
		t.Errorf("Expected TokenType 'User', got %q", token.TokenType)
	}
	if !token.IsActive() {
		t.Error("Expected token to be active")
	}
	if token.IsDeleted() {
		t.Error("Expected token to not be deleted")
	}
}

// TestTokenTimestamp tests timestamp handling
func TestTokenTimestamp(t *testing.T) {
	specificTime := time.Date(2024, 1, 15, 14, 30, 0, 0, time.UTC)

	token := types.Token{
		ID:        1,
		Timestamp: specificTime,
	}

	if !token.Timestamp.Equal(specificTime) {
		t.Errorf("Expected timestamp %v, got %v", specificTime, token.Timestamp)
	}
}

// TestTokenWithoutOptionalFields tests Token without optional fields
func TestTokenWithoutOptionalFields(t *testing.T) {
	token := types.Token{
		ID:        1,
		User:      "user",
		Timestamp: time.Now(),
	}

	if token.TaskID != nil {
		t.Error("TaskID should be nil")
	}
	if token.Operation != nil {
		t.Error("Operation should be nil")
	}

	str := token.String()
	if str == "" {
		t.Error("String() should not return empty string even without optional fields")
	}
}
