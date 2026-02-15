package unit

import (
	"testing"

	"github.com/nbaertsch/mythic-sdk-go/pkg/mythic/types"
)

// TestCommandString tests the Command.String() method
func TestCommandString(t *testing.T) {
	tests := []struct {
		name     string
		command  types.Command
		contains []string
	}{
		{
			name: "supported command",
			command: types.Command{
				ID:        1,
				Cmd:       "ls",
				Version:   2,
				Supported: true,
			},
			contains: []string{"ls", "v2"},
		},
		{
			name: "unsupported command",
			command: types.Command{
				ID:        2,
				Cmd:       "old_cmd",
				Version:   1,
				Supported: false,
			},
			contains: []string{"old_cmd", "v1", "unsupported"},
		},
		{
			name: "script-only command",
			command: types.Command{
				ID:         3,
				Cmd:        "script_cmd",
				Version:    1,
				Supported:  true,
				ScriptOnly: true,
			},
			contains: []string{"script_cmd", "v1", "script-only"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.command.String()
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

// TestCommandIsSupported tests the Command.IsSupported() method
func TestCommandIsSupported(t *testing.T) {
	tests := []struct {
		name    string
		command types.Command
		want    bool
	}{
		{
			name: "supported command",
			command: types.Command{
				ID:        1,
				Cmd:       "ls",
				Supported: true,
			},
			want: true,
		},
		{
			name: "unsupported command",
			command: types.Command{
				ID:        2,
				Cmd:       "old",
				Supported: false,
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.command.IsSupported(); got != tt.want {
				t.Errorf("IsSupported() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestCommandIsScriptOnly tests the Command.IsScriptOnly() method
func TestCommandIsScriptOnly(t *testing.T) {
	tests := []struct {
		name    string
		command types.Command
		want    bool
	}{
		{
			name: "script-only command",
			command: types.Command{
				ID:         1,
				Cmd:        "script",
				ScriptOnly: true,
			},
			want: true,
		},
		{
			name: "regular command",
			command: types.Command{
				ID:         2,
				Cmd:        "regular",
				ScriptOnly: false,
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.command.IsScriptOnly(); got != tt.want {
				t.Errorf("IsScriptOnly() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestCommandParameterString tests the CommandParameter.String() method
func TestCommandParameterString(t *testing.T) {
	tests := []struct {
		name     string
		param    types.CommandParameter
		contains []string
	}{
		{
			name: "required string parameter",
			param: types.CommandParameter{
				ID:       1,
				Name:     "path",
				Type:     types.ParameterTypeString,
				Required: true,
			},
			contains: []string{"path", "String", "required"},
		},
		{
			name: "optional boolean parameter",
			param: types.CommandParameter{
				ID:       2,
				Name:     "recursive",
				Type:     types.ParameterTypeBoolean,
				Required: false,
			},
			contains: []string{"recursive", "Boolean"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.param.String()
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

// TestCommandParameterIsRequired tests the CommandParameter.IsRequired() method
func TestCommandParameterIsRequired(t *testing.T) {
	tests := []struct {
		name  string
		param types.CommandParameter
		want  bool
	}{
		{
			name: "required parameter",
			param: types.CommandParameter{
				ID:       1,
				Name:     "path",
				Required: true,
			},
			want: true,
		},
		{
			name: "optional parameter",
			param: types.CommandParameter{
				ID:       2,
				Name:     "flag",
				Required: false,
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.param.IsRequired(); got != tt.want {
				t.Errorf("IsRequired() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestCommandParameterHasChoices tests the CommandParameter.HasChoices() method
func TestCommandParameterHasChoices(t *testing.T) {
	tests := []struct {
		name  string
		param types.CommandParameter
		want  bool
	}{
		{
			name: "parameter with static choices",
			param: types.CommandParameter{
				ID:      1,
				Name:    "option",
				Choices: "[\"A\", \"B\", \"C\"]",
			},
			want: true,
		},
		{
			name: "parameter with all commands",
			param: types.CommandParameter{
				ID:                    2,
				Name:                  "command",
				ChoicesAreAllCommands: true,
			},
			want: true,
		},
		{
			name: "parameter with loaded commands",
			param: types.CommandParameter{
				ID:                       3,
				Name:                     "loaded",
				ChoicesAreLoadedCommands: true,
			},
			want: true,
		},
		{
			name: "parameter without choices",
			param: types.CommandParameter{
				ID:   4,
				Name: "free",
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.param.HasChoices(); got != tt.want {
				t.Errorf("HasChoices() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestCommandParameterIsDynamic tests the CommandParameter.IsDynamic() method
func TestCommandParameterIsDynamic(t *testing.T) {
	tests := []struct {
		name  string
		param types.CommandParameter
		want  bool
	}{
		{
			name: "dynamic parameter",
			param: types.CommandParameter{
				ID:                   1,
				Name:                 "file",
				DynamicQueryFunction: "get_files",
			},
			want: true,
		},
		{
			name: "static parameter",
			param: types.CommandParameter{
				ID:   2,
				Name: "static",
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.param.IsDynamic(); got != tt.want {
				t.Errorf("IsDynamic() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestLoadedCommandString tests the LoadedCommand.String() method
func TestLoadedCommandString(t *testing.T) {
	tests := []struct {
		name     string
		loaded   types.LoadedCommand
		contains []string
	}{
		{
			name: "loaded command with details",
			loaded: types.LoadedCommand{
				ID:         1,
				CommandID:  10,
				CallbackID: 5,
				Version:    2,
				Command: &types.Command{
					Cmd:     "ls",
					Version: 2,
				},
			},
			contains: []string{"ls", "v2", "Callback 5"},
		},
		{
			name: "loaded command without details",
			loaded: types.LoadedCommand{
				ID:         2,
				CommandID:  15,
				CallbackID: 8,
				Version:    1,
			},
			contains: []string{"Command ID 15", "v1", "Callback 8"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.loaded.String()
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

// TestCommandTypes tests the Command type structure
func TestCommandTypes(t *testing.T) {
	command := types.Command{
		ID:                  1,
		Cmd:                 "download",
		PayloadTypeID:       5,
		Description:         "Download a file from the target",
		Help:                "Usage: download <path>",
		Version:             3,
		Supported:           true,
		Author:              "author",
		Attributes:          `{"builtin": true}`,
		ScriptOnly:          false,
		MitreAttackMappings: "T1005",
	}

	if command.ID != 1 {
		t.Errorf("Expected ID 1, got %d", command.ID)
	}
	if command.Cmd != "download" {
		t.Errorf("Expected Cmd 'download', got %q", command.Cmd)
	}
	if command.PayloadTypeID != 5 {
		t.Errorf("Expected PayloadTypeID 5, got %d", command.PayloadTypeID)
	}
	if command.Version != 3 {
		t.Errorf("Expected Version 3, got %d", command.Version)
	}
	if !command.Supported {
		t.Error("Expected Supported to be true")
	}
}

// TestCommandParameterTypes tests the CommandParameter type structure
func TestCommandParameterTypes(t *testing.T) {
	param := types.CommandParameter{
		ID:                          1,
		CommandID:                   10,
		Name:                        "path",
		Type:                        types.ParameterTypeString,
		Description:                 "File path to download",
		Required:                    true,
		DefaultValue:                "/tmp/file.txt",
		Choices:                     "",
		SupportedAgents:             "all",
		ChoicesAreAllCommands:       false,
		ChoicesAreLoadedCommands:    false,
		ChoiceFilterByCommandAttrib: "",
		DynamicQueryFunction:        "",
	}

	if param.ID != 1 {
		t.Errorf("Expected ID 1, got %d", param.ID)
	}
	if param.CommandID != 10 {
		t.Errorf("Expected CommandID 10, got %d", param.CommandID)
	}
	if param.Name != "path" {
		t.Errorf("Expected Name 'path', got %q", param.Name)
	}
	if param.Type != types.ParameterTypeString {
		t.Errorf("Expected Type 'String', got %q", param.Type)
	}
	if !param.Required {
		t.Error("Expected Required to be true")
	}
}

// TestLoadedCommandTypes tests the LoadedCommand type structure
func TestLoadedCommandTypes(t *testing.T) {
	loaded := types.LoadedCommand{
		ID:         1,
		CommandID:  10,
		CallbackID: 5,
		OperatorID: 2,
		Version:    3,
		Command: &types.Command{
			Cmd:     "ls",
			Version: 3,
		},
	}

	if loaded.ID != 1 {
		t.Errorf("Expected ID 1, got %d", loaded.ID)
	}
	if loaded.CommandID != 10 {
		t.Errorf("Expected CommandID 10, got %d", loaded.CommandID)
	}
	if loaded.CallbackID != 5 {
		t.Errorf("Expected CallbackID 5, got %d", loaded.CallbackID)
	}
	if loaded.OperatorID != 2 {
		t.Errorf("Expected OperatorID 2, got %d", loaded.OperatorID)
	}
	if loaded.Version != 3 {
		t.Errorf("Expected Version 3, got %d", loaded.Version)
	}
	if loaded.Command == nil {
		t.Error("Expected Command to not be nil")
	}
}

// TestLoadedCommandWithPayloadType tests loaded commands with enriched metadata
func TestLoadedCommandWithPayloadType(t *testing.T) {
	loaded := types.LoadedCommand{
		ID:         1,
		CommandID:  105,
		CallbackID: 5,
		OperatorID: 1,
		Version:    1,
		Command: &types.Command{
			ID:              105,
			Cmd:             "forge_collections",
			PayloadTypeID:   3,
			PayloadTypeName: "forge",
			Description:     "List available collections",
			Help:            "forge_collections",
			Version:         1,
			Author:          "@nbaertsch",
			ScriptOnly:      true,
			Supported:       true,
		},
	}

	if loaded.Command.PayloadTypeName != "forge" {
		t.Errorf("Expected PayloadTypeName 'forge', got %q", loaded.Command.PayloadTypeName)
	}
	if loaded.Command.PayloadTypeID != 3 {
		t.Errorf("Expected PayloadTypeID 3, got %d", loaded.Command.PayloadTypeID)
	}
	if !loaded.Command.ScriptOnly {
		t.Error("Expected ScriptOnly to be true for forge command")
	}
	if loaded.Command.Cmd != "forge_collections" {
		t.Errorf("Expected Cmd 'forge_collections', got %q", loaded.Command.Cmd)
	}
	if loaded.Command.Author != "@nbaertsch" {
		t.Errorf("Expected Author '@nbaertsch', got %q", loaded.Command.Author)
	}
}

// TestCommandPayloadTypeName tests the PayloadTypeName field on Command
func TestCommandPayloadTypeName(t *testing.T) {
	tests := []struct {
		name            string
		payloadTypeName string
		payloadTypeID   int
	}{
		{"poseidon", "poseidon", 1},
		{"xenon", "xenon", 2},
		{"forge", "forge", 3},
		{"empty", "", 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := types.Command{
				Cmd:             "test_cmd",
				PayloadTypeName: tt.payloadTypeName,
				PayloadTypeID:   tt.payloadTypeID,
			}
			if cmd.PayloadTypeName != tt.payloadTypeName {
				t.Errorf("Expected PayloadTypeName %q, got %q", tt.payloadTypeName, cmd.PayloadTypeName)
			}
			if cmd.PayloadTypeID != tt.payloadTypeID {
				t.Errorf("Expected PayloadTypeID %d, got %d", tt.payloadTypeID, cmd.PayloadTypeID)
			}
		})
	}
}

// TestParameterTypeConstants tests parameter type constants
func TestParameterTypeConstants(t *testing.T) {
	paramTypes := map[string]string{
		types.ParameterTypeString:         "String",
		types.ParameterTypeBoolean:        "Boolean",
		types.ParameterTypeNumber:         "Number",
		types.ParameterTypeChooseOne:      "ChooseOne",
		types.ParameterTypeChooseMultiple: "ChooseMultiple",
		types.ParameterTypeFile:           "File",
		types.ParameterTypeArray:          "Array",
		types.ParameterTypeCredential:     "Credential",
		types.ParameterTypeLinkInfo:       "LinkInfo",
	}

	for constant, expected := range paramTypes {
		if constant != expected {
			t.Errorf("Expected parameter type %q, got %q", expected, constant)
		}
	}
}
