//go:build integration

package integration

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/nbaertsch/mythic-sdk-go/pkg/mythic/types"
)

// TestE2E_GenerateReport tests operation report generation.
// Covers: GenerateReport with various filters and formats
func TestE2E_GenerateReport(t *testing.T) {
	client := AuthenticateTestClient(t)

	// Get current operation
	ctx0, cancel0 := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel0()

	operations, err := client.GetOperations(ctx0)
	if err != nil {
		t.Fatalf("GetOperations failed: %v", err)
	}

	if len(operations) == 0 {
		t.Skip("No operations found, skipping report generation tests")
	}

	testOperation := operations[0]
	t.Logf("Using operation: %s (ID: %d)", testOperation.Name, testOperation.ID)

	// Test 1: Basic report generation (minimal)
	t.Log("=== Test 1: Generate basic report ===")
	ctx1, cancel1 := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel1()

	basicReq := &types.GenerateReportRequest{
		OperationID: testOperation.ID,
	}

	report, err := client.GenerateReport(ctx1, basicReq)
	if err != nil {
		t.Fatalf("GenerateReport (basic) failed: %v", err)
	}

	if report.Status != "success" {
		t.Errorf("Report generation failed: %s", report.Error)
	}

	if len(report.ReportData) == 0 {
		t.Error("Report data is empty")
	}
	t.Logf("✓ Basic report generated: %d bytes", len(report.ReportData))

	// Test 2: Report with all sections included
	t.Log("=== Test 2: Generate comprehensive report ===")
	ctx2, cancel2 := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel2()

	comprehensiveReq := &types.GenerateReportRequest{
		OperationID:        testOperation.ID,
		IncludeCallbacks:   true,
		IncludeTasks:       true,
		IncludeFiles:       true,
		IncludeCredentials: true,
		IncludeArtifacts:   true,
		IncludeMITRE:       true,
	}

	compReport, err := client.GenerateReport(ctx2, comprehensiveReq)
	if err != nil {
		t.Fatalf("GenerateReport (comprehensive) failed: %v", err)
	}

	if compReport.Status != "success" {
		t.Errorf("Comprehensive report generation failed: %s", compReport.Error)
	}

	if len(compReport.ReportData) == 0 {
		t.Error("Comprehensive report data is empty")
	}
	t.Logf("✓ Comprehensive report generated: %d bytes", len(compReport.ReportData))

	// Verify comprehensive report is larger than basic
	if len(compReport.ReportData) <= len(report.ReportData) {
		t.Log("  ⚠ Comprehensive report not significantly larger than basic report")
	} else {
		t.Logf("  ✓ Comprehensive report is %d bytes larger", len(compReport.ReportData)-len(report.ReportData))
	}

	// Test 3: Report with different output formats
	t.Log("=== Test 3: Generate reports in different formats ===")

	formats := []string{
		types.ReportFormatJSON,
		types.ReportFormatMarkdown,
	}

	for _, format := range formats {
		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)

		formatReq := &types.GenerateReportRequest{
			OperationID:  testOperation.ID,
			OutputFormat: format,
		}

		formatReport, err := client.GenerateReport(ctx, formatReq)
		cancel()

		if err != nil {
			t.Errorf("GenerateReport (%s format) failed: %v", format, err)
			continue
		}

		if formatReport.Status != "success" {
			t.Errorf("Report generation (%s format) failed: %s", format, formatReport.Error)
			continue
		}

		if len(formatReport.ReportData) == 0 {
			t.Errorf("Report data is empty for %s format", format)
			continue
		}

		t.Logf("✓ Report generated in %s format: %d bytes", format, len(formatReport.ReportData))

		// Basic format validation
		switch format {
		case types.ReportFormatJSON:
			if !strings.HasPrefix(strings.TrimSpace(formatReport.ReportData), "{") {
				t.Errorf("JSON report doesn't start with '{'")
			}
		case types.ReportFormatMarkdown:
			if !strings.Contains(formatReport.ReportData, "#") {
				t.Log("  ⚠ Markdown report doesn't contain '#' headers")
			}
		}
	}

	t.Log("=== ✓ Report generation tests passed ===")
}

// TestE2E_GenerateReportWithFilters tests report generation with date and callback filters.
func TestE2E_GenerateReportWithFilters(t *testing.T) {
	client := AuthenticateTestClient(t)

	// Get current operation
	ctx0, cancel0 := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel0()

	operations, err := client.GetOperations(ctx0)
	if err != nil {
		t.Fatalf("GetOperations failed: %v", err)
	}

	if len(operations) == 0 {
		t.Skip("No operations found, skipping filtered report tests")
	}

	testOperation := operations[0]

	// Test 1: Report with date range filter
	t.Log("=== Test 1: Generate report with date range ===")
	ctx1, cancel1 := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel1()

	endDate := time.Now().Format("2006-01-02")
	startDate := time.Now().AddDate(0, 0, -7).Format("2006-01-02") // Last 7 days

	dateReq := &types.GenerateReportRequest{
		OperationID:  testOperation.ID,
		StartDate:    &startDate,
		EndDate:      &endDate,
		IncludeTasks: true,
	}

	dateReport, err := client.GenerateReport(ctx1, dateReq)
	if err != nil {
		t.Fatalf("GenerateReport with date filter failed: %v", err)
	}

	if dateReport.Status != "success" {
		t.Errorf("Date-filtered report generation failed: %s", dateReport.Error)
	}
	t.Logf("✓ Report generated with date range (%s to %s): %d bytes", startDate, endDate, len(dateReport.ReportData))

	// Test 2: Report with callback filter
	t.Log("=== Test 2: Generate report with callback filter ===")

	// Get callbacks for operation
	ctx2a, cancel2a := context.WithTimeout(context.Background(), 30*time.Second)
	callbacks, err := client.GetAllActiveCallbacks(ctx2a)
	cancel2a()

	if err != nil || len(callbacks) == 0 {
		t.Log("⚠ No callbacks found, skipping callback filter test")
	} else {
		ctx2, cancel2 := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel2()

		// Filter by first callback
		callbackIDs := []int{callbacks[0].ID}

		callbackReq := &types.GenerateReportRequest{
			OperationID:      testOperation.ID,
			CallbackIDs:      callbackIDs,
			IncludeCallbacks: true,
			IncludeTasks:     true,
		}

		cbReport, err := client.GenerateReport(ctx2, callbackReq)
		if err != nil {
			t.Fatalf("GenerateReport with callback filter failed: %v", err)
		}

		if cbReport.Status != "success" {
			t.Errorf("Callback-filtered report generation failed: %s", cbReport.Error)
		}
		t.Logf("✓ Report generated for callback %d: %d bytes", callbackIDs[0], len(cbReport.ReportData))
	}

	t.Log("=== ✓ Filtered report tests passed ===")
}

// TestE2E_GenerateReportMITRE tests report generation with MITRE ATT&CK coverage.
func TestE2E_GenerateReportMITRE(t *testing.T) {
	client := AuthenticateTestClient(t)

	// Get current operation
	ctx0, cancel0 := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel0()

	operations, err := client.GetOperations(ctx0)
	if err != nil {
		t.Fatalf("GetOperations failed: %v", err)
	}

	if len(operations) == 0 {
		t.Skip("No operations found, skipping MITRE report tests")
	}

	testOperation := operations[0]

	t.Log("=== Test: Generate report with MITRE ATT&CK coverage ===")
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	mitreReq := &types.GenerateReportRequest{
		OperationID:  testOperation.ID,
		IncludeMITRE: true,
		IncludeTasks: true,
	}

	mitreReport, err := client.GenerateReport(ctx, mitreReq)
	if err != nil {
		t.Fatalf("GenerateReport with MITRE failed: %v", err)
	}

	if mitreReport.Status != "success" {
		t.Errorf("MITRE report generation failed: %s", mitreReport.Error)
	}

	if len(mitreReport.ReportData) == 0 {
		t.Error("MITRE report data is empty")
	}
	t.Logf("✓ MITRE report generated: %d bytes", len(mitreReport.ReportData))

	// Check for MITRE-related content
	if strings.Contains(strings.ToLower(mitreReport.ReportData), "mitre") ||
		strings.Contains(strings.ToLower(mitreReport.ReportData), "att&ck") ||
		strings.Contains(strings.ToLower(mitreReport.ReportData), "tactic") ||
		strings.Contains(strings.ToLower(mitreReport.ReportData), "technique") {
		t.Log("  ✓ Report contains MITRE ATT&CK references")
	} else {
		t.Log("  ⚠ Report may not contain MITRE ATT&CK data (or no techniques used)")
	}

	t.Log("=== ✓ MITRE report tests passed ===")
}

// TestE2E_GetRedirectRules tests redirect rule retrieval for payloads.
func TestE2E_GetRedirectRules(t *testing.T) {
	client := AuthenticateTestClient(t)

	// Get a payload to test with
	ctx0, cancel0 := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel0()

	payloads, err := client.GetPayloads(ctx0)
	if err != nil {
		t.Fatalf("GetPayloads failed: %v", err)
	}

	if len(payloads) == 0 {
		t.Skip("No payloads found, skipping redirect rule tests")
	}

	testPayload := payloads[0]
	t.Logf("Using payload: %s (UUID: %s)", testPayload.Description, testPayload.UUID)

	// Test: Get redirect rules
	t.Log("=== Test: Get redirect rules for payload ===")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	rules, err := client.GetRedirectRules(ctx, testPayload.UUID)
	if err != nil {
		t.Fatalf("GetRedirectRules failed: %v", err)
	}

	t.Logf("✓ Retrieved %d redirect rules", len(rules))

	if len(rules) == 0 {
		t.Log("  ⚠ No redirect rules found for payload (may be expected)")
		t.Log("=== ✓ Redirect rule tests passed ===")
		return
	}

	// Validate redirect rule structure
	for _, rule := range rules {
		if rule.ID == 0 {
			t.Error("RedirectRule has ID 0")
		}
		if rule.PayloadType == "" {
			t.Error("RedirectRule has empty PayloadType")
		}
		if rule.C2Profile == "" {
			t.Error("RedirectRule has empty C2Profile")
		}
		if rule.RedirectType == "" {
			t.Error("RedirectRule has empty RedirectType")
		}
		if rule.Rule == "" {
			t.Error("RedirectRule has empty Rule")
		}
	}

	// Show sample rules
	sampleCount := 3
	if len(rules) < sampleCount {
		sampleCount = len(rules)
	}

	t.Logf("  Sample redirect rules:")
	for i := 0; i < sampleCount; i++ {
		rule := rules[i]
		t.Logf("    [%d] %s (%s) - %s", i+1, rule.RedirectType, rule.C2Profile, rule.PayloadType)
		if rule.Description != "" {
			t.Logf("        Description: %s", rule.Description)
		}
	}

	// Analyze redirect types
	redirectTypes := make(map[string]int)
	for _, rule := range rules {
		redirectTypes[rule.RedirectType]++
	}

	t.Logf("  Redirect type distribution:")
	for rtype, count := range redirectTypes {
		t.Logf("    %s: %d", rtype, count)
	}

	t.Log("=== ✓ Redirect rule tests passed ===")
}

// TestE2E_ReportErrorHandling tests error scenarios for reporting operations.
func TestE2E_ReportErrorHandling(t *testing.T) {
	client := AuthenticateTestClient(t)

	// Test 1: Generate report with invalid operation ID
	t.Log("=== Test 1: Generate report with invalid operation ID ===")
	ctx1, cancel1 := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel1()

	invalidReq := &types.GenerateReportRequest{
		OperationID: 999999,
	}

	_, err := client.GenerateReport(ctx1, invalidReq)
	if err == nil {
		t.Error("Expected error for invalid operation ID")
	}
	t.Logf("✓ Invalid operation ID rejected: %v", err)

	// Test 2: Generate report with zero operation ID
	t.Log("=== Test 2: Generate report with zero operation ID ===")
	ctx2, cancel2 := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel2()

	zeroReq := &types.GenerateReportRequest{
		OperationID: 0,
	}

	_, err = client.GenerateReport(ctx2, zeroReq)
	if err == nil {
		t.Error("Expected error for zero operation ID")
	}
	t.Logf("✓ Zero operation ID rejected: %v", err)

	// Test 3: Generate report with nil request
	t.Log("=== Test 3: Generate report with nil request ===")
	ctx3, cancel3 := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel3()

	_, err = client.GenerateReport(ctx3, nil)
	if err == nil {
		t.Error("Expected error for nil request")
	}
	t.Logf("✓ Nil request rejected: %v", err)

	// Test 4: Get redirect rules with empty UUID
	t.Log("=== Test 4: Get redirect rules with empty UUID ===")
	ctx4, cancel4 := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel4()

	_, err = client.GetRedirectRules(ctx4, "")
	if err == nil {
		t.Error("Expected error for empty UUID")
	}
	t.Logf("✓ Empty UUID rejected: %v", err)

	// Test 5: Get redirect rules with invalid UUID
	t.Log("=== Test 5: Get redirect rules with invalid UUID ===")
	ctx5, cancel5 := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel5()

	_, err = client.GetRedirectRules(ctx5, "invalid-uuid-123")
	if err == nil {
		t.Error("Expected error for invalid UUID")
	}
	t.Logf("✓ Invalid UUID rejected: %v", err)

	t.Log("=== ✓ All error handling tests passed ===")
}

// TestE2E_ReportContentAnalysis tests report content structure and completeness.
func TestE2E_ReportContentAnalysis(t *testing.T) {
	client := AuthenticateTestClient(t)

	// Get current operation
	ctx0, cancel0 := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel0()

	operations, err := client.GetOperations(ctx0)
	if err != nil {
		t.Fatalf("GetOperations failed: %v", err)
	}

	if len(operations) == 0 {
		t.Skip("No operations found, skipping content analysis")
	}

	testOperation := operations[0]

	t.Log("=== Test: Analyze report content structure ===")
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	req := &types.GenerateReportRequest{
		OperationID:        testOperation.ID,
		IncludeCallbacks:   true,
		IncludeTasks:       true,
		IncludeFiles:       true,
		IncludeCredentials: true,
		IncludeArtifacts:   true,
		IncludeMITRE:       true,
		OutputFormat:       types.ReportFormatJSON,
	}

	report, err := client.GenerateReport(ctx, req)
	if err != nil {
		// Report generation may not be available in all Mythic versions (e.g., v3.4.20)
		if isSchemaError(err) {
			t.Logf("⚠ Report generation not available in this Mythic version: %v", err)
			t.Skip("Skipping report content analysis - feature not available")
		}
		t.Fatalf("GenerateReport failed: %v", err)
	}

	if report.Status != "success" {
		t.Errorf("Report generation failed: %s", report.Error)
	}

	t.Logf("✓ Report generated: %d bytes", len(report.ReportData))

	// Analyze content
	content := report.ReportData

	// Check for operation information
	if strings.Contains(strings.ToLower(content), "operation") {
		t.Log("  ✓ Report contains operation information")
	}

	// Check for requested sections
	sections := map[string]string{
		"callbacks":   "callback",
		"tasks":       "task",
		"files":       "file",
		"credentials": "credential",
		"artifacts":   "artifact",
	}

	for section, keyword := range sections {
		if strings.Contains(strings.ToLower(content), keyword) {
			t.Logf("  ✓ Report contains %s section", section)
		} else {
			t.Logf("  ⚠ Report may not contain %s section (or section is empty)", section)
		}
	}

	// Check size ranges
	if len(content) < 100 {
		t.Error("Report is suspiciously small (< 100 bytes)")
	} else if len(content) < 1000 {
		t.Log("  ⚠ Report is small (< 1KB) - operation may have minimal data")
	} else if len(content) > 1000000 {
		t.Log("  ⚠ Report is large (> 1MB) - may contain extensive operation data")
	} else {
		t.Logf("  ✓ Report size is reasonable (%d bytes)", len(content))
	}

	t.Log("=== ✓ Content analysis complete ===")
}
