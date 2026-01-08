package unit

import (
	"testing"

	"github.com/nbaertsch/mythic-sdk-go/pkg/mythic/types"
)

// TestBrowserScriptString tests the BrowserScript.String() method
func TestBrowserScriptString(t *testing.T) {
	tests := []struct {
		name     string
		script   types.BrowserScript
		contains []string
	}{
		{
			name: "with description",
			script: types.BrowserScript{
				ID:          1,
				Name:        "custom_task_button",
				Description: "Adds custom task button to UI",
				Active:      true,
			},
			contains: []string{"custom_task_button", "Adds custom task button to UI"},
		},
		{
			name: "without description",
			script: types.BrowserScript{
				ID:     2,
				Name:   "screenshot_viewer",
				Active: false,
			},
			contains: []string{"screenshot_viewer"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.script.String()
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

// TestBrowserScriptIsActive tests the BrowserScript.IsActive() method
func TestBrowserScriptIsActive(t *testing.T) {
	tests := []struct {
		name   string
		script types.BrowserScript
		want   bool
	}{
		{
			name: "active script",
			script: types.BrowserScript{
				ID:     1,
				Name:   "test",
				Active: true,
			},
			want: true,
		},
		{
			name: "inactive script",
			script: types.BrowserScript{
				ID:     2,
				Name:   "test",
				Active: false,
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.script.IsActive(); got != tt.want {
				t.Errorf("IsActive() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestBrowserScriptIsForNewUI tests the BrowserScript.IsForNewUI() method
func TestBrowserScriptIsForNewUI(t *testing.T) {
	tests := []struct {
		name   string
		script types.BrowserScript
		want   bool
	}{
		{
			name: "for new UI",
			script: types.BrowserScript{
				ID:       1,
				Name:     "test",
				ForNewUI: true,
			},
			want: true,
		},
		{
			name: "for old UI",
			script: types.BrowserScript{
				ID:       2,
				Name:     "test",
				ForNewUI: false,
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.script.IsForNewUI(); got != tt.want {
				t.Errorf("IsForNewUI() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestBrowserScriptOperationString tests the BrowserScriptOperation.String() method
func TestBrowserScriptOperationString(t *testing.T) {
	operatorID := 5

	tests := []struct {
		name     string
		bso      types.BrowserScriptOperation
		contains []string
	}{
		{
			name: "with script name",
			bso: types.BrowserScriptOperation{
				ID:              1,
				BrowserScriptID: 10,
				OperationID:     5,
				ScriptName:      "custom_exporter",
				Active:          true,
			},
			contains: []string{"custom_exporter", "Operation 5"},
		},
		{
			name: "without script name",
			bso: types.BrowserScriptOperation{
				ID:              2,
				BrowserScriptID: 15,
				OperationID:     8,
				Active:          false,
			},
			contains: []string{"Script ID 15", "Operation 8"},
		},
		{
			name: "with operator ID",
			bso: types.BrowserScriptOperation{
				ID:              3,
				BrowserScriptID: 20,
				OperationID:     12,
				OperatorID:      &operatorID,
				ScriptName:      "operator_script",
				Active:          true,
			},
			contains: []string{"operator_script", "Operation 12"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.bso.String()
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

// TestBrowserScriptOperationIsActive tests the BrowserScriptOperation.IsActive() method
func TestBrowserScriptOperationIsActive(t *testing.T) {
	tests := []struct {
		name string
		bso  types.BrowserScriptOperation
		want bool
	}{
		{
			name: "active",
			bso: types.BrowserScriptOperation{
				ID:              1,
				BrowserScriptID: 10,
				OperationID:     5,
				Active:          true,
			},
			want: true,
		},
		{
			name: "inactive",
			bso: types.BrowserScriptOperation{
				ID:              2,
				BrowserScriptID: 15,
				OperationID:     8,
				Active:          false,
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.bso.IsActive(); got != tt.want {
				t.Errorf("IsActive() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestBrowserScriptOperationIsOperatorSpecific tests the BrowserScriptOperation.IsOperatorSpecific() method
func TestBrowserScriptOperationIsOperatorSpecific(t *testing.T) {
	operatorID := 5

	tests := []struct {
		name string
		bso  types.BrowserScriptOperation
		want bool
	}{
		{
			name: "operator specific",
			bso: types.BrowserScriptOperation{
				ID:              1,
				BrowserScriptID: 10,
				OperationID:     5,
				OperatorID:      &operatorID,
				Active:          true,
			},
			want: true,
		},
		{
			name: "not operator specific",
			bso: types.BrowserScriptOperation{
				ID:              2,
				BrowserScriptID: 15,
				OperationID:     8,
				OperatorID:      nil,
				Active:          true,
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.bso.IsOperatorSpecific(); got != tt.want {
				t.Errorf("IsOperatorSpecific() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestCustomBrowserExportRequestString tests the CustomBrowserExportRequest.String() method
func TestCustomBrowserExportRequestString(t *testing.T) {
	tests := []struct {
		name     string
		req      types.CustomBrowserExportRequest
		contains []string
	}{
		{
			name: "basic request",
			req: types.CustomBrowserExportRequest{
				OperationID: 5,
				ScriptName:  "custom_exporter",
			},
			contains: []string{"custom_exporter", "Operation 5"},
		},
		{
			name: "with parameters",
			req: types.CustomBrowserExportRequest{
				OperationID: 10,
				ScriptName:  "advanced_export",
				Parameters: map[string]interface{}{
					"format": "csv",
					"filter": "active",
				},
			},
			contains: []string{"advanced_export", "Operation 10"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.req.String()
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

// TestBrowserScriptTypes tests the BrowserScript type structure
func TestBrowserScriptTypes(t *testing.T) {
	script := types.BrowserScript{
		ID:          1,
		Name:        "screenshot_renderer",
		Script:      "function render(data) { return '<img src=\"' + data.screenshot + '\">'; }",
		Author:      "admin",
		ForNewUI:    true,
		Active:      true,
		Description: "Renders screenshots in custom format",
	}

	if script.ID != 1 {
		t.Errorf("Expected ID 1, got %d", script.ID)
	}
	if script.Name != "screenshot_renderer" {
		t.Errorf("Expected Name 'screenshot_renderer', got %q", script.Name)
	}
	if script.Script == "" {
		t.Error("Expected Script to not be empty")
	}
	if script.Author != "admin" {
		t.Errorf("Expected Author 'admin', got %q", script.Author)
	}
	if !script.ForNewUI {
		t.Error("Expected ForNewUI to be true")
	}
	if !script.Active {
		t.Error("Expected Active to be true")
	}
	if script.Description == "" {
		t.Error("Expected Description to not be empty")
	}
}

// TestBrowserScriptOperationTypes tests the BrowserScriptOperation type structure
func TestBrowserScriptOperationTypes(t *testing.T) {
	operatorID := 5

	bso := types.BrowserScriptOperation{
		ID:              1,
		BrowserScriptID: 10,
		OperationID:     5,
		OperatorID:      &operatorID,
		ScriptName:      "custom_task_button",
		Active:          true,
	}

	if bso.ID != 1 {
		t.Errorf("Expected ID 1, got %d", bso.ID)
	}
	if bso.BrowserScriptID != 10 {
		t.Errorf("Expected BrowserScriptID 10, got %d", bso.BrowserScriptID)
	}
	if bso.OperationID != 5 {
		t.Errorf("Expected OperationID 5, got %d", bso.OperationID)
	}
	if bso.OperatorID == nil || *bso.OperatorID != 5 {
		t.Error("Expected OperatorID to be 5")
	}
	if bso.ScriptName != "custom_task_button" {
		t.Errorf("Expected ScriptName 'custom_task_button', got %q", bso.ScriptName)
	}
	if !bso.Active {
		t.Error("Expected Active to be true")
	}
}

// TestCustomBrowserExportRequestTypes tests the CustomBrowserExportRequest type structure
func TestCustomBrowserExportRequestTypes(t *testing.T) {
	req := types.CustomBrowserExportRequest{
		OperationID: 5,
		ScriptName:  "data_exporter",
		Parameters: map[string]interface{}{
			"format":     "json",
			"include":    []string{"tasks", "callbacks"},
			"start_date": "2024-01-01",
		},
	}

	if req.OperationID != 5 {
		t.Errorf("Expected OperationID 5, got %d", req.OperationID)
	}
	if req.ScriptName != "data_exporter" {
		t.Errorf("Expected ScriptName 'data_exporter', got %q", req.ScriptName)
	}
	if req.Parameters == nil {
		t.Error("Expected Parameters to not be nil")
	}
	if len(req.Parameters) != 3 {
		t.Errorf("Expected 3 parameters, got %d", len(req.Parameters))
	}
}

// TestCustomBrowserExportResponseTypes tests the CustomBrowserExportResponse type structure
func TestCustomBrowserExportResponseTypes(t *testing.T) {
	response := types.CustomBrowserExportResponse{
		Status:       "success",
		Error:        "",
		ExportedData: `{"tasks": 150, "callbacks": 25}`,
	}

	if response.Status != "success" {
		t.Errorf("Expected Status 'success', got %q", response.Status)
	}
	if response.Error != "" {
		t.Errorf("Expected empty Error, got %q", response.Error)
	}
	if response.ExportedData == "" {
		t.Error("Expected ExportedData to not be empty")
	}
}

// TestBrowserScriptWithoutOptionalFields tests BrowserScript without optional fields
func TestBrowserScriptWithoutOptionalFields(t *testing.T) {
	script := types.BrowserScript{
		ID:       1,
		Name:     "minimal_script",
		Script:   "function test() {}",
		ForNewUI: false,
		Active:   false,
	}

	if script.Author != "" {
		t.Error("Author should be empty")
	}
	if script.Description != "" {
		t.Error("Description should be empty")
	}

	str := script.String()
	if str == "" {
		t.Error("String() should not return empty string even without optional fields")
	}
}

// TestBrowserScriptOperationWithoutOptionalFields tests BrowserScriptOperation without optional fields
func TestBrowserScriptOperationWithoutOptionalFields(t *testing.T) {
	bso := types.BrowserScriptOperation{
		ID:              1,
		BrowserScriptID: 10,
		OperationID:     5,
		Active:          true,
	}

	if bso.OperatorID != nil {
		t.Error("OperatorID should be nil")
	}
	if bso.ScriptName != "" {
		t.Error("ScriptName should be empty")
	}

	str := bso.String()
	if str == "" {
		t.Error("String() should not return empty string even without optional fields")
	}
}

// TestCustomBrowserExportRequestWithoutParameters tests request without parameters
func TestCustomBrowserExportRequestWithoutParameters(t *testing.T) {
	req := types.CustomBrowserExportRequest{
		OperationID: 5,
		ScriptName:  "simple_export",
	}

	if req.Parameters != nil && len(req.Parameters) > 0 {
		t.Error("Parameters should be nil or empty")
	}

	str := req.String()
	if str == "" {
		t.Error("String() should not return empty string even without parameters")
	}
}
