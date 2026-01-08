package unit

import (
	"testing"

	"github.com/nbaertsch/mythic-sdk-go/pkg/mythic/types"
)

// TestGenerateReportRequestString tests the GenerateReportRequest.String() method
func TestGenerateReportRequestString(t *testing.T) {
	tests := []struct {
		name     string
		request  types.GenerateReportRequest
		contains []string
	}{
		{
			name: "basic request",
			request: types.GenerateReportRequest{
				OperationID: 5,
			},
			contains: []string{"Operation 5"},
		},
		{
			name: "request with options",
			request: types.GenerateReportRequest{
				OperationID:  10,
				IncludeMITRE: true,
				OutputFormat: types.ReportFormatJSON,
			},
			contains: []string{"Operation 10"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.request.String()
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

// TestRedirectRuleString tests the RedirectRule.String() method
func TestRedirectRuleString(t *testing.T) {
	tests := []struct {
		name     string
		rule     types.RedirectRule
		contains []string
	}{
		{
			name: "with description",
			rule: types.RedirectRule{
				ID:           1,
				PayloadType:  "apollo",
				C2Profile:    "http",
				RedirectType: types.RedirectTypeApache,
				Description:  "Block non-agent traffic",
			},
			contains: []string{"apache", "apollo", "http", "Block non-agent traffic"},
		},
		{
			name: "without description",
			rule: types.RedirectRule{
				ID:           2,
				PayloadType:  "apfell",
				C2Profile:    "https",
				RedirectType: types.RedirectTypeNginx,
			},
			contains: []string{"nginx", "apfell", "https"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.rule.String()
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

// TestGenerateReportRequestTypes tests the GenerateReportRequest type structure
func TestGenerateReportRequestTypes(t *testing.T) {
	startDate := "2024-01-01"
	endDate := "2024-12-31"

	request := types.GenerateReportRequest{
		OperationID:        5,
		IncludeMITRE:       true,
		OutputFormat:       types.ReportFormatJSON,
		IncludeCallbacks:   true,
		IncludeTasks:       true,
		IncludeFiles:       true,
		IncludeCredentials: true,
		IncludeArtifacts:   true,
		StartDate:          &startDate,
		EndDate:            &endDate,
		CallbackIDs:        []int{1, 2, 3},
	}

	if request.OperationID != 5 {
		t.Errorf("Expected OperationID 5, got %d", request.OperationID)
	}
	if !request.IncludeMITRE {
		t.Error("Expected IncludeMITRE to be true")
	}
	if request.OutputFormat != types.ReportFormatJSON {
		t.Errorf("Expected OutputFormat 'json', got %q", request.OutputFormat)
	}
	if !request.IncludeCallbacks {
		t.Error("Expected IncludeCallbacks to be true")
	}
	if request.StartDate == nil || *request.StartDate != startDate {
		t.Error("Expected StartDate to be set")
	}
	if request.EndDate == nil || *request.EndDate != endDate {
		t.Error("Expected EndDate to be set")
	}
	if len(request.CallbackIDs) != 3 {
		t.Errorf("Expected 3 CallbackIDs, got %d", len(request.CallbackIDs))
	}
}

// TestRedirectRuleTypes tests the RedirectRule type structure
func TestRedirectRuleTypes(t *testing.T) {
	rule := types.RedirectRule{
		ID:           1,
		PayloadType:  "apollo",
		C2Profile:    "http",
		RedirectType: types.RedirectTypeApache,
		Rule:         "RewriteRule ^.*$ http://redirector [L,R=302]",
		Description:  "Redirect all traffic",
	}

	if rule.ID != 1 {
		t.Errorf("Expected ID 1, got %d", rule.ID)
	}
	if rule.PayloadType != "apollo" {
		t.Errorf("Expected PayloadType 'apollo', got %q", rule.PayloadType)
	}
	if rule.C2Profile != "http" {
		t.Errorf("Expected C2Profile 'http', got %q", rule.C2Profile)
	}
	if rule.RedirectType != types.RedirectTypeApache {
		t.Errorf("Expected RedirectType 'apache', got %q", rule.RedirectType)
	}
	if rule.Rule == "" {
		t.Error("Expected Rule to not be empty")
	}
}

// TestReportFormatConstants tests report format constants
func TestReportFormatConstants(t *testing.T) {
	formats := map[string]string{
		"json":     types.ReportFormatJSON,
		"markdown": types.ReportFormatMarkdown,
		"html":     types.ReportFormatHTML,
		"pdf":      types.ReportFormatPDF,
	}

	for expected, actual := range formats {
		if actual != expected {
			t.Errorf("Expected format %q, got %q", expected, actual)
		}
	}
}

// TestRedirectTypeConstants tests redirect type constants
func TestRedirectTypeConstants(t *testing.T) {
	redirectTypes := map[string]string{
		"apache":      types.RedirectTypeApache,
		"nginx":       types.RedirectTypeNginx,
		"mod_rewrite": types.RedirectTypeModRewrite,
	}

	for expected, actual := range redirectTypes {
		if actual != expected {
			t.Errorf("Expected redirect type %q, got %q", expected, actual)
		}
	}
}

// TestGenerateReportRequestWithoutOptionalFields tests request without optional fields
func TestGenerateReportRequestWithoutOptionalFields(t *testing.T) {
	request := types.GenerateReportRequest{
		OperationID: 1,
	}

	if request.OutputFormat != "" {
		t.Error("OutputFormat should be empty")
	}
	if request.IncludeMITRE {
		t.Error("IncludeMITRE should be false")
	}
	if request.StartDate != nil {
		t.Error("StartDate should be nil")
	}
	if request.EndDate != nil {
		t.Error("EndDate should be nil")
	}
	if len(request.CallbackIDs) != 0 {
		t.Error("CallbackIDs should be empty")
	}

	str := request.String()
	if str == "" {
		t.Error("String() should not return empty string even without optional fields")
	}
}

// TestRedirectRuleAllTypes tests all redirect rule types
func TestRedirectRuleAllTypes(t *testing.T) {
	redirectTypes := []string{
		types.RedirectTypeApache,
		types.RedirectTypeNginx,
		types.RedirectTypeModRewrite,
	}

	for _, redirectType := range redirectTypes {
		rule := types.RedirectRule{
			ID:           1,
			PayloadType:  "test",
			C2Profile:    "http",
			RedirectType: redirectType,
			Rule:         "test rule",
		}

		if rule.RedirectType != redirectType {
			t.Errorf("Expected RedirectType %q, got %q", redirectType, rule.RedirectType)
		}

		str := rule.String()
		if !contains(str, redirectType) {
			t.Errorf("String() should contain redirect type %q, got %q", redirectType, str)
		}
	}
}

// TestGenerateReportResponseTypes tests the GenerateReportResponse type
func TestGenerateReportResponseTypes(t *testing.T) {
	response := types.GenerateReportResponse{
		Status:     "success",
		Error:      "",
		ReportData: "{\"operation\":\"test\"}",
	}

	if response.Status != "success" {
		t.Errorf("Expected Status 'success', got %q", response.Status)
	}
	if response.Error != "" {
		t.Errorf("Expected empty Error, got %q", response.Error)
	}
	if response.ReportData == "" {
		t.Error("Expected ReportData to not be empty")
	}
}

// TestReportFormats tests various report formats
func TestReportFormats(t *testing.T) {
	formats := []string{
		types.ReportFormatJSON,
		types.ReportFormatMarkdown,
		types.ReportFormatHTML,
		types.ReportFormatPDF,
	}

	for _, format := range formats {
		request := types.GenerateReportRequest{
			OperationID:  1,
			OutputFormat: format,
		}

		if request.OutputFormat != format {
			t.Errorf("Expected OutputFormat %q, got %q", format, request.OutputFormat)
		}
	}
}
