//go:build integration
// +build integration

package integration

import (
	"context"
	"testing"

	"github.com/nbaertsch/mythic-sdk-go/pkg/mythic/types"
)

// TestReporting_GenerateReport tests generating an operation report
func TestReporting_GenerateReport(t *testing.T) {
	ctx := context.Background()

	// Get current operation
	operationID := client.GetCurrentOperation()
	if operationID == nil {
		t.Skip("No current operation set")
	}

	// Generate a basic report
	request := &types.GenerateReportRequest{
		OperationID:  *operationID,
		OutputFormat: types.ReportFormatJSON,
		IncludeMITRE: true,
	}

	response, err := client.GenerateReport(ctx, request)
	if err != nil {
		t.Fatalf("Failed to generate report: %v", err)
	}

	if response.Status != "success" {
		t.Errorf("Expected status 'success', got %q (error: %s)", response.Status, response.Error)
	}

	t.Logf("Generated report for operation %d", *operationID)
	t.Logf("Report status: %s", response.Status)
	if response.ReportData != "" {
		t.Logf("Report data length: %d bytes", len(response.ReportData))
	}
}

// TestReporting_GenerateReportWithAllOptions tests report generation with all options
func TestReporting_GenerateReportWithAllOptions(t *testing.T) {
	ctx := context.Background()

	// Get current operation
	operationID := client.GetCurrentOperation()
	if operationID == nil {
		t.Skip("No current operation set")
	}

	// Generate a comprehensive report with all options
	request := &types.GenerateReportRequest{
		OperationID:        *operationID,
		OutputFormat:       types.ReportFormatJSON,
		IncludeMITRE:       true,
		IncludeCallbacks:   true,
		IncludeTasks:       true,
		IncludeFiles:       true,
		IncludeCredentials: true,
		IncludeArtifacts:   true,
	}

	response, err := client.GenerateReport(ctx, request)
	if err != nil {
		t.Fatalf("Failed to generate comprehensive report: %v", err)
	}

	if response.Status != "success" {
		t.Errorf("Expected status 'success', got %q (error: %s)", response.Status, response.Error)
	}

	t.Logf("Generated comprehensive report for operation %d", *operationID)
	t.Logf("Report includes: MITRE, Callbacks, Tasks, Files, Credentials, Artifacts")
}

// TestReporting_GenerateReportWithDateRange tests report with date filtering
func TestReporting_GenerateReportWithDateRange(t *testing.T) {
	ctx := context.Background()

	// Get current operation
	operationID := client.GetCurrentOperation()
	if operationID == nil {
		t.Skip("No current operation set")
	}

	startDate := "2024-01-01"
	endDate := "2024-12-31"

	// Generate a report with date range
	request := &types.GenerateReportRequest{
		OperationID:  *operationID,
		OutputFormat: types.ReportFormatJSON,
		StartDate:    &startDate,
		EndDate:      &endDate,
	}

	response, err := client.GenerateReport(ctx, request)
	if err != nil {
		t.Fatalf("Failed to generate report with date range: %v", err)
	}

	if response.Status != "success" {
		t.Errorf("Expected status 'success', got %q (error: %s)", response.Status, response.Error)
	}

	t.Logf("Generated report with date range: %s to %s", startDate, endDate)
}

// TestReporting_GenerateReportWithCallbackFilter tests report with specific callbacks
func TestReporting_GenerateReportWithCallbackFilter(t *testing.T) {
	ctx := context.Background()

	// Get current operation
	operationID := client.GetCurrentOperation()
	if operationID == nil {
		t.Skip("No current operation set")
	}

	// Get callbacks to filter by
	callbacks, err := client.GetAllCallbacks(ctx)
	if err != nil || len(callbacks) == 0 {
		t.Skip("No callbacks available for filtering")
	}

	// Use first two callbacks for filtering
	callbackIDs := []int{callbacks[0].ID}
	if len(callbacks) > 1 {
		callbackIDs = append(callbackIDs, callbacks[1].ID)
	}

	// Generate a report filtered by callbacks
	request := &types.GenerateReportRequest{
		OperationID:  *operationID,
		OutputFormat: types.ReportFormatJSON,
		CallbackIDs:  callbackIDs,
	}

	response, err := client.GenerateReport(ctx, request)
	if err != nil {
		t.Fatalf("Failed to generate report with callback filter: %v", err)
	}

	if response.Status != "success" {
		t.Errorf("Expected status 'success', got %q (error: %s)", response.Status, response.Error)
	}

	t.Logf("Generated report filtered by %d callbacks", len(callbackIDs))
}

// TestReporting_GenerateReport_InvalidInput tests report generation with invalid input
func TestReporting_GenerateReport_InvalidInput(t *testing.T) {
	ctx := context.Background()

	// Try to generate report with nil request
	_, err := client.GenerateReport(ctx, nil)
	if err == nil {
		t.Error("Expected error for nil request")
	}

	// Try to generate report with operation ID 0
	request := &types.GenerateReportRequest{
		OperationID: 0,
	}
	_, err = client.GenerateReport(ctx, request)
	if err == nil {
		t.Error("Expected error for operation ID 0")
	}
}

// TestReporting_GenerateReportFormats tests different report output formats
func TestReporting_GenerateReportFormats(t *testing.T) {
	ctx := context.Background()

	// Get current operation
	operationID := client.GetCurrentOperation()
	if operationID == nil {
		t.Skip("No current operation set")
	}

	formats := []string{
		types.ReportFormatJSON,
		types.ReportFormatMarkdown,
		types.ReportFormatHTML,
		types.ReportFormatPDF,
	}

	for _, format := range formats {
		request := &types.GenerateReportRequest{
			OperationID:  *operationID,
			OutputFormat: format,
		}

		response, err := client.GenerateReport(ctx, request)
		if err != nil {
			t.Logf("Format %s: %v (may not be supported)", format, err)
			continue
		}

		if response.Status == "success" {
			t.Logf("Successfully generated report in %s format", format)
		}
	}
}

// TestReporting_GetRedirectRules tests retrieving redirect rules for a payload
func TestReporting_GetRedirectRules(t *testing.T) {
	ctx := context.Background()

	// Get payloads to find a valid UUID
	payloads, err := client.GetPayloads(ctx)
	if err != nil || len(payloads) == 0 {
		t.Skip("No payloads available for testing redirect rules")
	}

	// Get redirect rules for the first payload
	payloadUUID := payloads[0].UUID
	rules, err := client.GetRedirectRules(ctx, payloadUUID)
	if err != nil {
		t.Fatalf("Failed to get redirect rules: %v", err)
	}

	t.Logf("Retrieved %d redirect rules for payload %s", len(rules), payloadUUID)

	// Verify structure if any rules exist
	if len(rules) > 0 {
		rule := rules[0]
		if rule.PayloadType == "" {
			t.Error("Rule PayloadType should not be empty")
		}
		if rule.C2Profile == "" {
			t.Error("Rule C2Profile should not be empty")
		}
		if rule.RedirectType == "" {
			t.Error("Rule RedirectType should not be empty")
		}
		if rule.Rule == "" {
			t.Error("Rule content should not be empty")
		}

		// Test String method
		str := rule.String()
		if str == "" {
			t.Error("RedirectRule.String() should not return empty string")
		}
		t.Logf("Redirect rule: %s", str)
		t.Logf("  Payload Type: %s", rule.PayloadType)
		t.Logf("  C2 Profile: %s", rule.C2Profile)
		t.Logf("  Redirect Type: %s", rule.RedirectType)
		t.Logf("  Rule length: %d bytes", len(rule.Rule))
	}
}

// TestReporting_GetRedirectRules_InvalidInput tests with invalid input
func TestReporting_GetRedirectRules_InvalidInput(t *testing.T) {
	ctx := context.Background()

	// Try to get redirect rules with empty UUID
	_, err := client.GetRedirectRules(ctx, "")
	if err == nil {
		t.Error("Expected error for empty payload UUID")
	}

	// Try to get redirect rules with invalid UUID
	_, err = client.GetRedirectRules(ctx, "invalid-uuid-12345")
	if err == nil {
		t.Error("Expected error for invalid payload UUID")
	}
}

// TestReporting_GetRedirectRulesMultiplePayloads tests rules for multiple payloads
func TestReporting_GetRedirectRulesMultiplePayloads(t *testing.T) {
	ctx := context.Background()

	// Get payloads
	payloads, err := client.GetPayloads(ctx)
	if err != nil || len(payloads) < 2 {
		t.Skip("Need at least 2 payloads for testing")
	}

	// Get redirect rules for multiple payloads
	for i, payload := range payloads {
		if i >= 3 { // Test first 3 payloads only
			break
		}

		rules, err := client.GetRedirectRules(ctx, payload.UUID)
		if err != nil {
			t.Logf("Payload %s: Failed to get redirect rules: %v", payload.UUID, err)
			continue
		}

		t.Logf("Payload %s (%s): %d redirect rules", payload.UUID, payload.Description, len(rules))
	}
}

// TestReporting_ReportWithMITRECoverage tests report generation with MITRE coverage
func TestReporting_ReportWithMITRECoverage(t *testing.T) {
	ctx := context.Background()

	// Get current operation
	operationID := client.GetCurrentOperation()
	if operationID == nil {
		t.Skip("No current operation set")
	}

	// Generate report with MITRE ATT&CK coverage
	request := &types.GenerateReportRequest{
		OperationID:  *operationID,
		OutputFormat: types.ReportFormatJSON,
		IncludeMITRE: true,
	}

	response, err := client.GenerateReport(ctx, request)
	if err != nil {
		t.Fatalf("Failed to generate report with MITRE coverage: %v", err)
	}

	if response.Status != "success" {
		t.Errorf("Expected status 'success', got %q (error: %s)", response.Status, response.Error)
	}

	t.Logf("Generated report with MITRE ATT&CK coverage for operation %d", *operationID)
}
