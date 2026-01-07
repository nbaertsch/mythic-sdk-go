package unit

import (
	"testing"
	"time"

	"github.com/nbaertsch/mythic-sdk-go/pkg/mythic/types"
)

// TestKeylogString tests the Keylog.String() method
func TestKeylogString(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name     string
		keylog   types.Keylog
		contains []string
	}{
		{
			name: "with window and user",
			keylog: types.Keylog{
				ID:        1,
				Window:    "Chrome - Gmail",
				User:      "DOMAIN\\user",
				Timestamp: now,
			},
			contains: []string{"Chrome - Gmail", "DOMAIN\\user"},
		},
		{
			name: "with user only",
			keylog: types.Keylog{
				ID:        2,
				User:      "admin",
				Timestamp: now,
			},
			contains: []string{"admin"},
		},
		{
			name: "with ID only",
			keylog: types.Keylog{
				ID:        3,
				Timestamp: now,
			},
			contains: []string{"Keylog 3"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.keylog.String()
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

// TestKeylogHasKeystrokes tests the Keylog.HasKeystrokes() method
func TestKeylogHasKeystrokes(t *testing.T) {
	tests := []struct {
		name       string
		keystrokes string
		expected   bool
	}{
		{"with keystrokes", "password123", true},
		{"empty keystrokes", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			keylog := types.Keylog{Keystrokes: tt.keystrokes}
			result := keylog.HasKeystrokes()
			if result != tt.expected {
				t.Errorf("HasKeystrokes() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

// TestKeylogTypes tests the Keylog type structure
func TestKeylogTypes(t *testing.T) {
	now := time.Now()

	keylog := types.Keylog{
		ID:          1,
		TaskID:      42,
		Keystrokes:  "username: admin\npassword: secret123",
		Window:      "Remote Desktop Connection",
		Timestamp:   now,
		OperationID: 5,
		User:        "WORKSTATION\\user",
		CallbackID:  10,
	}

	if keylog.ID != 1 {
		t.Errorf("Expected ID 1, got %d", keylog.ID)
	}
	if keylog.TaskID != 42 {
		t.Errorf("Expected TaskID 42, got %d", keylog.TaskID)
	}
	if !keylog.HasKeystrokes() {
		t.Error("Expected keylog to have keystrokes")
	}
	if keylog.Window != "Remote Desktop Connection" {
		t.Errorf("Expected Window 'Remote Desktop Connection', got %q", keylog.Window)
	}
	if keylog.User != "WORKSTATION\\user" {
		t.Errorf("Expected User 'WORKSTATION\\user', got %q", keylog.User)
	}
}

// TestKeylogWithoutKeystrokes tests keylog without captured keystrokes
func TestKeylogWithoutKeystrokes(t *testing.T) {
	keylog := types.Keylog{
		ID:        1,
		Window:    "Notepad",
		Timestamp: time.Now(),
	}

	if keylog.HasKeystrokes() {
		t.Error("Expected keylog to have no keystrokes")
	}

	str := keylog.String()
	if str == "" {
		t.Error("String() should not return empty string even without keystrokes")
	}
}

// TestKeylogTimestamp tests timestamp handling
func TestKeylogTimestamp(t *testing.T) {
	specificTime := time.Date(2024, 1, 15, 14, 30, 0, 0, time.UTC)

	keylog := types.Keylog{
		ID:        1,
		Timestamp: specificTime,
		Window:    "Test Window",
	}

	if !keylog.Timestamp.Equal(specificTime) {
		t.Errorf("Expected timestamp %v, got %v", specificTime, keylog.Timestamp)
	}

	str := keylog.String()
	// Should contain the date in some format
	if str == "" {
		t.Error("String() should include timestamp information")
	}
}

// TestKeylogMultilineKeystrokes tests handling of multiline keystrokes
func TestKeylogMultilineKeystrokes(t *testing.T) {
	multiline := `Line 1: username
Line 2: password
Line 3: some command
Line 4: more text`

	keylog := types.Keylog{
		ID:         1,
		Keystrokes: multiline,
		Window:     "Terminal",
	}

	if !keylog.HasKeystrokes() {
		t.Error("Should detect multiline keystrokes")
	}

	if keylog.Keystrokes != multiline {
		t.Error("Multiline keystrokes should be preserved exactly")
	}
}

// TestKeylogWindowNames tests various window name formats
func TestKeylogWindowNames(t *testing.T) {
	windows := []string{
		"Chrome - example.com",
		"Microsoft Word - Document1.docx",
		"cmd.exe",
		"Outlook - Inbox",
		"Remote Desktop - 10.0.0.1",
	}

	for _, window := range windows {
		keylog := types.Keylog{
			ID:        1,
			Window:    window,
			Timestamp: time.Now(),
		}

		if keylog.Window != window {
			t.Errorf("Expected window %q, got %q", window, keylog.Window)
		}

		str := keylog.String()
		if !contains(str, window) {
			t.Errorf("String() should contain window name %q, got %q", window, str)
		}
	}
}

// TestKeylogUserFormats tests various user format handling
func TestKeylogUserFormats(t *testing.T) {
	users := []string{
		"DOMAIN\\user",
		"user@domain.com",
		"administrator",
		"WORKSTATION-01\\localuser",
	}

	for _, user := range users {
		keylog := types.Keylog{
			ID:        1,
			User:      user,
			Timestamp: time.Now(),
		}

		if keylog.User != user {
			t.Errorf("Expected user %q, got %q", user, keylog.User)
		}

		str := keylog.String()
		if user != "" && !contains(str, user) {
			t.Errorf("String() should contain user %q, got %q", user, str)
		}
	}
}

// TestKeylogFields tests all keylog fields
func TestKeylogFields(t *testing.T) {
	now := time.Now()

	keylog := types.Keylog{
		ID:          123,
		TaskID:      456,
		Keystrokes:  "test keystrokes",
		Window:      "Test Window",
		Timestamp:   now,
		OperationID: 789,
		User:        "testuser",
		CallbackID:  101,
	}

	// Verify all fields are set correctly
	if keylog.ID != 123 {
		t.Errorf("Expected ID 123, got %d", keylog.ID)
	}
	if keylog.TaskID != 456 {
		t.Errorf("Expected TaskID 456, got %d", keylog.TaskID)
	}
	if keylog.Keystrokes != "test keystrokes" {
		t.Errorf("Expected keystrokes 'test keystrokes', got %q", keylog.Keystrokes)
	}
	if keylog.Window != "Test Window" {
		t.Errorf("Expected window 'Test Window', got %q", keylog.Window)
	}
	if keylog.OperationID != 789 {
		t.Errorf("Expected OperationID 789, got %d", keylog.OperationID)
	}
	if keylog.User != "testuser" {
		t.Errorf("Expected user 'testuser', got %q", keylog.User)
	}
	if keylog.CallbackID != 101 {
		t.Errorf("Expected CallbackID 101, got %d", keylog.CallbackID)
	}
}
