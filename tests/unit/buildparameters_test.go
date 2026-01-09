package unit

import (
	"testing"
	"time"

	"github.com/nbaertsch/mythic-sdk-go/pkg/mythic/types"
)

// TestBuildParameterTypeString tests the BuildParameterType.String() method
func TestBuildParameterTypeString(t *testing.T) {
	tests := []struct {
		name     string
		param    types.BuildParameterType
		contains []string
	}{
		{
			name: "required string parameter",
			param: types.BuildParameterType{
				ID:            1,
				Name:          "callback_host",
				ParameterType: types.BuildParameterTypeString,
				Required:      true,
			},
			contains: []string{"callback_host", "String", "required"},
		},
		{
			name: "optional boolean parameter",
			param: types.BuildParameterType{
				ID:            2,
				Name:          "debug_mode",
				ParameterType: types.BuildParameterTypeBoolean,
				Required:      false,
			},
			contains: []string{"debug_mode", "Boolean"},
		},
		{
			name: "required number parameter",
			param: types.BuildParameterType{
				ID:            3,
				Name:          "port",
				ParameterType: types.BuildParameterTypeNumber,
				Required:      true,
			},
			contains: []string{"port", "Number", "required"},
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

// TestBuildParameterTypeIsRequired tests the BuildParameterType.IsRequired() method
func TestBuildParameterTypeIsRequired(t *testing.T) {
	tests := []struct {
		name  string
		param types.BuildParameterType
		want  bool
	}{
		{
			name: "required parameter",
			param: types.BuildParameterType{
				ID:       1,
				Name:     "callback_host",
				Required: true,
			},
			want: true,
		},
		{
			name: "optional parameter",
			param: types.BuildParameterType{
				ID:       2,
				Name:     "optional_param",
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

// TestBuildParameterTypeIsCrypto tests the BuildParameterType.IsCrypto() method
func TestBuildParameterTypeIsCrypto(t *testing.T) {
	tests := []struct {
		name  string
		param types.BuildParameterType
		want  bool
	}{
		{
			name: "crypto parameter",
			param: types.BuildParameterType{
				ID:           1,
				Name:         "encryption_key",
				IsCryptoType: true,
			},
			want: true,
		},
		{
			name: "non-crypto parameter",
			param: types.BuildParameterType{
				ID:           2,
				Name:         "callback_host",
				IsCryptoType: false,
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.param.IsCrypto(); got != tt.want {
				t.Errorf("IsCrypto() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestBuildParameterTypeShouldRandomize tests the BuildParameterType.ShouldRandomize() method
func TestBuildParameterTypeShouldRandomize(t *testing.T) {
	tests := []struct {
		name  string
		param types.BuildParameterType
		want  bool
	}{
		{
			name: "randomized parameter",
			param: types.BuildParameterType{
				ID:        1,
				Name:      "session_id",
				Randomize: true,
			},
			want: true,
		},
		{
			name: "non-randomized parameter",
			param: types.BuildParameterType{
				ID:        2,
				Name:      "callback_host",
				Randomize: false,
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.param.ShouldRandomize(); got != tt.want {
				t.Errorf("ShouldRandomize() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestBuildParameterTypeIsDeleted tests the BuildParameterType.IsDeleted() method
func TestBuildParameterTypeIsDeleted(t *testing.T) {
	tests := []struct {
		name  string
		param types.BuildParameterType
		want  bool
	}{
		{
			name: "deleted parameter",
			param: types.BuildParameterType{
				ID:      1,
				Name:    "old_param",
				Deleted: true,
			},
			want: true,
		},
		{
			name: "active parameter",
			param: types.BuildParameterType{
				ID:      2,
				Name:    "active_param",
				Deleted: false,
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.param.IsDeleted(); got != tt.want {
				t.Errorf("IsDeleted() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestBuildParameterInstanceString tests the BuildParameterInstance.String() method
func TestBuildParameterInstanceString(t *testing.T) {
	tests := []struct {
		name     string
		instance types.BuildParameterInstance
		contains []string
	}{
		{
			name: "instance with parameter details",
			instance: types.BuildParameterInstance{
				ID:               1,
				PayloadID:        10,
				BuildParameterID: 5,
				Value:            "https://127.0.0.1:443",
				BuildParameter: &types.BuildParameterType{
					ID:   5,
					Name: "callback_host",
				},
			},
			contains: []string{"callback_host", "https://127.0.0.1:443"},
		},
		{
			name: "instance without parameter details",
			instance: types.BuildParameterInstance{
				ID:               2,
				PayloadID:        11,
				BuildParameterID: 6,
				Value:            "true",
			},
			contains: []string{"BuildParameter 6", "true"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.instance.String()
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

// TestBuildParameterInstanceIsEncrypted tests the BuildParameterInstance.IsEncrypted() method
func TestBuildParameterInstanceIsEncrypted(t *testing.T) {
	encValue := "encrypted_data"
	emptyEncValue := ""

	tests := []struct {
		name     string
		instance types.BuildParameterInstance
		want     bool
	}{
		{
			name: "encrypted instance",
			instance: types.BuildParameterInstance{
				ID:       1,
				Value:    "plaintext",
				EncValue: &encValue,
			},
			want: true,
		},
		{
			name: "non-encrypted instance",
			instance: types.BuildParameterInstance{
				ID:       2,
				Value:    "plaintext",
				EncValue: nil,
			},
			want: false,
		},
		{
			name: "empty encrypted value",
			instance: types.BuildParameterInstance{
				ID:       3,
				Value:    "plaintext",
				EncValue: &emptyEncValue,
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.instance.IsEncrypted(); got != tt.want {
				t.Errorf("IsEncrypted() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestBuildParameterInstanceGetValue tests the BuildParameterInstance.GetValue() method
func TestBuildParameterInstanceGetValue(t *testing.T) {
	encValue := "encrypted_data"

	tests := []struct {
		name     string
		instance types.BuildParameterInstance
		want     string
	}{
		{
			name: "get encrypted value",
			instance: types.BuildParameterInstance{
				ID:       1,
				Value:    "plaintext",
				EncValue: &encValue,
			},
			want: "encrypted_data",
		},
		{
			name: "get plain value",
			instance: types.BuildParameterInstance{
				ID:       2,
				Value:    "plaintext",
				EncValue: nil,
			},
			want: "plaintext",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.instance.GetValue(); got != tt.want {
				t.Errorf("GetValue() = %q, want %q", got, tt.want)
			}
		})
	}
}

// TestBuildParameterTypeFields tests the BuildParameterType structure
func TestBuildParameterTypeFields(t *testing.T) {
	now := time.Now()
	param := types.BuildParameterType{
		ID:                 1,
		Name:               "callback_host",
		PayloadTypeID:      5,
		Description:        "The callback host for the agent",
		Parameter:          `{"type": "string"}`,
		Required:           true,
		VerifierRegex:      `^https?://.*`,
		DefaultValue:       "https://127.0.0.1:443",
		Randomize:          false,
		FormatString:       "",
		ParameterType:      types.BuildParameterTypeString,
		IsCryptoType:       false,
		Deleted:            false,
		CreationTime:       now,
		ParameterGroupName: "C2 Configuration",
	}

	if param.ID != 1 {
		t.Errorf("Expected ID 1, got %d", param.ID)
	}
	if param.Name != "callback_host" {
		t.Errorf("Expected Name 'callback_host', got %q", param.Name)
	}
	if param.PayloadTypeID != 5 {
		t.Errorf("Expected PayloadTypeID 5, got %d", param.PayloadTypeID)
	}
	if !param.Required {
		t.Error("Expected Required to be true")
	}
	if param.ParameterType != types.BuildParameterTypeString {
		t.Errorf("Expected ParameterType 'String', got %q", param.ParameterType)
	}
	if param.DefaultValue != "https://127.0.0.1:443" {
		t.Errorf("Expected DefaultValue 'https://127.0.0.1:443', got %q", param.DefaultValue)
	}
	if param.ParameterGroupName != "C2 Configuration" {
		t.Errorf("Expected ParameterGroupName 'C2 Configuration', got %q", param.ParameterGroupName)
	}
}

// TestBuildParameterInstanceFields tests the BuildParameterInstance structure
func TestBuildParameterInstanceFields(t *testing.T) {
	now := time.Now()
	encValue := "encrypted_value"

	instance := types.BuildParameterInstance{
		ID:               1,
		PayloadID:        10,
		BuildParameterID: 5,
		Value:            "https://attacker.com:443",
		EncValue:         &encValue,
		CreationTime:     now,
		BuildParameter: &types.BuildParameterType{
			ID:   5,
			Name: "callback_host",
		},
	}

	if instance.ID != 1 {
		t.Errorf("Expected ID 1, got %d", instance.ID)
	}
	if instance.PayloadID != 10 {
		t.Errorf("Expected PayloadID 10, got %d", instance.PayloadID)
	}
	if instance.BuildParameterID != 5 {
		t.Errorf("Expected BuildParameterID 5, got %d", instance.BuildParameterID)
	}
	if instance.Value != "https://attacker.com:443" {
		t.Errorf("Expected Value 'https://attacker.com:443', got %q", instance.Value)
	}
	if !instance.IsEncrypted() {
		t.Error("Expected instance to be encrypted")
	}
	if instance.BuildParameter == nil {
		t.Error("Expected BuildParameter to not be nil")
	}
	if instance.BuildParameter.Name != "callback_host" {
		t.Errorf("Expected BuildParameter.Name 'callback_host', got %q", instance.BuildParameter.Name)
	}
}

// TestBuildParameterTypeConstants tests build parameter type constants
func TestBuildParameterTypeConstants(t *testing.T) {
	paramTypes := map[string]string{
		types.BuildParameterTypeString:         "String",
		types.BuildParameterTypeBoolean:        "Boolean",
		types.BuildParameterTypeNumber:         "Number",
		types.BuildParameterTypeChooseOne:      "ChooseOne",
		types.BuildParameterTypeChooseMultiple: "ChooseMultiple",
		types.BuildParameterTypeFile:           "File",
		types.BuildParameterTypeArray:          "Array",
		types.BuildParameterTypeDate:           "Date",
	}

	for constant, expected := range paramTypes {
		if constant != expected {
			t.Errorf("Expected parameter type %q, got %q", expected, constant)
		}
	}
}

// TestBuildParameterHelperMethods tests all helper methods together
func TestBuildParameterHelperMethods(t *testing.T) {
	param := types.BuildParameterType{
		ID:           1,
		Name:         "encryption_key",
		Required:     true,
		IsCryptoType: true,
		Randomize:    true,
		Deleted:      false,
	}

	if !param.IsRequired() {
		t.Error("Expected IsRequired() to return true")
	}
	if !param.IsCrypto() {
		t.Error("Expected IsCrypto() to return true")
	}
	if !param.ShouldRandomize() {
		t.Error("Expected ShouldRandomize() to return true")
	}
	if param.IsDeleted() {
		t.Error("Expected IsDeleted() to return false")
	}

	str := param.String()
	if str == "" {
		t.Error("Expected non-empty String()")
	}
}

// TestBuildParameterInstanceHelperMethods tests all instance helper methods together
func TestBuildParameterInstanceHelperMethods(t *testing.T) {
	encValue := "encrypted"

	instance := types.BuildParameterInstance{
		ID:       1,
		Value:    "plain",
		EncValue: &encValue,
		BuildParameter: &types.BuildParameterType{
			Name: "test_param",
		},
	}

	if !instance.IsEncrypted() {
		t.Error("Expected IsEncrypted() to return true")
	}
	if got := instance.GetValue(); got != "encrypted" {
		t.Errorf("Expected GetValue() to return 'encrypted', got %q", got)
	}

	str := instance.String()
	if str == "" {
		t.Error("Expected non-empty String()")
	}
	if !contains(str, "test_param") {
		t.Errorf("Expected String() to contain 'test_param', got %q", str)
	}
}
