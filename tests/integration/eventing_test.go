//go:build integration

package integration

import (
	"context"
	"testing"
)

// TestEventingTriggerManual_InvalidInput tests input validation for EventingTriggerManual.
func TestEventingTriggerManual_InvalidInput(t *testing.T) {
	client := getTestClient(t)
	ctx := context.Background()

	// Test with zero event group ID
	_, err := client.EventingTriggerManual(ctx, 0, 0, nil)
	if err == nil {
		t.Fatal("EventingTriggerManual with zero event group ID should return error")
	}
	t.Logf("Zero event group ID error: %v", err)

	// Test with negative event group ID
	_, err = client.EventingTriggerManual(ctx, -1, 0, nil)
	if err == nil {
		t.Fatal("EventingTriggerManual with negative event group ID should return error")
	}
	t.Logf("Negative event group ID error: %v", err)
}

// TestEventingTriggerManual_NonexistentEventGroup tests triggering non-existent event group.
func TestEventingTriggerManual_NonexistentEventGroup(t *testing.T) {
	client := getTestClient(t)
	ctx := context.Background()

	// Try to trigger a non-existent event group
	result, err := client.EventingTriggerManual(ctx, 999999, 0, nil)
	if err != nil {
		t.Logf("Nonexistent event group error (expected): %v", err)
	}

	if result != nil {
		t.Logf("Response: %s", result.String())
		if result.IsSuccessful() {
			t.Error("Should not succeed for nonexistent event group")
		}
	}
}

// TestEventingTriggerManualBulk_InvalidInput tests input validation for bulk trigger.
func TestEventingTriggerManualBulk_InvalidInput(t *testing.T) {
	client := getTestClient(t)
	ctx := context.Background()

	// Test with zero event group ID
	_, err := client.EventingTriggerManualBulk(ctx, 0, []int{1, 2, 3}, nil)
	if err == nil {
		t.Fatal("EventingTriggerManualBulk with zero event group ID should return error")
	}
	t.Logf("Zero event group ID error: %v", err)

	// Test with empty object IDs list
	_, err = client.EventingTriggerManualBulk(ctx, 1, []int{}, nil)
	if err == nil {
		t.Fatal("EventingTriggerManualBulk with empty object IDs should return error")
	}
	t.Logf("Empty object IDs error: %v", err)
}

// TestEventingTriggerKeyword_InvalidInput tests keyword trigger validation.
func TestEventingTriggerKeyword_InvalidInput(t *testing.T) {
	client := getTestClient(t)
	ctx := context.Background()

	// Test with empty keyword
	_, err := client.EventingTriggerKeyword(ctx, "", 0, nil)
	if err == nil {
		t.Fatal("EventingTriggerKeyword with empty keyword should return error")
	}
	t.Logf("Empty keyword error: %v", err)
}

// TestEventingTriggerKeyword_NonexistentKeyword tests keyword that doesn't match any events.
func TestEventingTriggerKeyword_NonexistentKeyword(t *testing.T) {
	client := getTestClient(t)
	ctx := context.Background()

	// Try to trigger with a keyword that likely doesn't exist
	result, err := client.EventingTriggerKeyword(ctx, "nonexistent_keyword_12345", 0, nil)
	if err != nil {
		t.Logf("Nonexistent keyword error (expected): %v", err)
	}

	if result != nil {
		t.Logf("Response: %s", result.String())
	}
}

// TestEventingTriggerCancel_InvalidInput tests cancel validation.
func TestEventingTriggerCancel_InvalidInput(t *testing.T) {
	client := getTestClient(t)
	ctx := context.Background()

	// Test with zero execution ID
	_, err := client.EventingTriggerCancel(ctx, 0)
	if err == nil {
		t.Fatal("EventingTriggerCancel with zero execution ID should return error")
	}
	t.Logf("Zero execution ID error: %v", err)

	// Test with negative execution ID
	_, err = client.EventingTriggerCancel(ctx, -1)
	if err == nil {
		t.Fatal("EventingTriggerCancel with negative execution ID should return error")
	}
	t.Logf("Negative execution ID error: %v", err)
}

// TestEventingTriggerRetry_InvalidInput tests retry validation.
func TestEventingTriggerRetry_InvalidInput(t *testing.T) {
	client := getTestClient(t)
	ctx := context.Background()

	// Test with zero execution ID
	_, err := client.EventingTriggerRetry(ctx, 0)
	if err == nil {
		t.Fatal("EventingTriggerRetry with zero execution ID should return error")
	}
	t.Logf("Zero execution ID error: %v", err)
}

// TestEventingTriggerRetryFromStep_InvalidInput tests retry from step validation.
func TestEventingTriggerRetryFromStep_InvalidInput(t *testing.T) {
	client := getTestClient(t)
	ctx := context.Background()

	// Test with zero execution ID
	_, err := client.EventingTriggerRetryFromStep(ctx, 0, 1)
	if err == nil {
		t.Fatal("EventingTriggerRetryFromStep with zero execution ID should return error")
	}
	t.Logf("Zero execution ID error: %v", err)

	// Test with zero step number
	_, err = client.EventingTriggerRetryFromStep(ctx, 1, 0)
	if err == nil {
		t.Fatal("EventingTriggerRetryFromStep with zero step number should return error")
	}
	t.Logf("Zero step number error: %v", err)

	// Test with negative step number
	_, err = client.EventingTriggerRetryFromStep(ctx, 1, -1)
	if err == nil {
		t.Fatal("EventingTriggerRetryFromStep with negative step number should return error")
	}
	t.Logf("Negative step number error: %v", err)
}

// TestEventingTriggerRunAgain_InvalidInput tests run again validation.
func TestEventingTriggerRunAgain_InvalidInput(t *testing.T) {
	client := getTestClient(t)
	ctx := context.Background()

	// Test with zero execution ID
	_, err := client.EventingTriggerRunAgain(ctx, 0)
	if err == nil {
		t.Fatal("EventingTriggerRunAgain with zero execution ID should return error")
	}
	t.Logf("Zero execution ID error: %v", err)
}

// TestEventingTriggerUpdate_InvalidInput tests update validation.
func TestEventingTriggerUpdate_InvalidInput(t *testing.T) {
	client := getTestClient(t)
	ctx := context.Background()

	// Test with zero event group ID
	_, err := client.EventingTriggerUpdate(ctx, 0, "", "", nil, nil, nil, nil, nil)
	if err == nil {
		t.Fatal("EventingTriggerUpdate with zero event group ID should return error")
	}
	t.Logf("Zero event group ID error: %v", err)
}

// TestEventingExportWorkflow_InvalidInput tests export validation.
func TestEventingExportWorkflow_InvalidInput(t *testing.T) {
	client := getTestClient(t)
	ctx := context.Background()

	// Test with zero workflow ID
	_, err := client.EventingExportWorkflow(ctx, 0)
	if err == nil {
		t.Fatal("EventingExportWorkflow with zero workflow ID should return error")
	}
	t.Logf("Zero workflow ID error: %v", err)

	// Test with negative workflow ID
	_, err = client.EventingExportWorkflow(ctx, -1)
	if err == nil {
		t.Fatal("EventingExportWorkflow with negative workflow ID should return error")
	}
	t.Logf("Negative workflow ID error: %v", err)
}

// TestEventingExportWorkflow_NonexistentWorkflow tests exporting non-existent workflow.
func TestEventingExportWorkflow_NonexistentWorkflow(t *testing.T) {
	client := getTestClient(t)
	ctx := context.Background()

	// Try to export a non-existent workflow
	result, err := client.EventingExportWorkflow(ctx, 999999)
	if err != nil {
		t.Logf("Nonexistent workflow error (expected): %v", err)
	}

	if result != nil {
		t.Logf("Response: %s", result.String())
		if result.IsSuccessful() {
			t.Error("Should not succeed for nonexistent workflow")
		}
	}
}

// TestEventingImportContainerWorkflow_InvalidInput tests import validation.
func TestEventingImportContainerWorkflow_InvalidInput(t *testing.T) {
	client := getTestClient(t)
	ctx := context.Background()

	// Test with empty container name
	_, err := client.EventingImportContainerWorkflow(ctx, "", "workflow.json", 0)
	if err == nil {
		t.Fatal("EventingImportContainerWorkflow with empty container name should return error")
	}
	t.Logf("Empty container name error: %v", err)

	// Test with empty workflow file
	_, err = client.EventingImportContainerWorkflow(ctx, "apollo", "", 0)
	if err == nil {
		t.Fatal("EventingImportContainerWorkflow with empty workflow file should return error")
	}
	t.Logf("Empty workflow file error: %v", err)
}

// TestEventingImportContainerWorkflow_NonexistentContainer tests importing from non-existent container.
func TestEventingImportContainerWorkflow_NonexistentContainer(t *testing.T) {
	client := getTestClient(t)
	ctx := context.Background()

	// Try to import from a non-existent container
	result, err := client.EventingImportContainerWorkflow(ctx, "nonexistent_container_12345", "workflow.json", 0)
	if err != nil {
		t.Logf("Nonexistent container error (expected): %v", err)
	}

	if result != nil {
		t.Logf("Response: %s", result.String())
		if result.IsSuccessful() {
			t.Error("Should not succeed for nonexistent container")
		}
	}
}

// TestEventingTestFile_InvalidInput tests test file validation.
func TestEventingTestFile_InvalidInput(t *testing.T) {
	client := getTestClient(t)
	ctx := context.Background()

	// Test with empty workflow file
	_, err := client.EventingTestFile(ctx, "", nil)
	if err == nil {
		t.Fatal("EventingTestFile with empty workflow file should return error")
	}
	t.Logf("Empty workflow file error: %v", err)
}

// TestEventingTestFile_NonexistentFile tests testing non-existent file.
func TestEventingTestFile_NonexistentFile(t *testing.T) {
	client := getTestClient(t)
	ctx := context.Background()

	// Try to test a non-existent workflow file
	result, err := client.EventingTestFile(ctx, "nonexistent_workflow_12345.json", nil)
	if err != nil {
		t.Logf("Nonexistent file error: %v", err)
	}

	if result != nil {
		t.Logf("Response: %s", result.String())
		t.Logf("Valid: %v", result.IsValid())
		t.Logf("Has errors: %v", result.HasErrors())

		if result.HasErrors() {
			t.Logf("Validation errors: %v", result.Errors)
		}
	}
}

// TestUpdateEventGroupApproval_InvalidInput tests approval validation.
func TestUpdateEventGroupApproval_InvalidInput(t *testing.T) {
	client := getTestClient(t)
	ctx := context.Background()

	// Test with zero event group ID
	_, err := client.UpdateEventGroupApproval(ctx, 0, true, "")
	if err == nil {
		t.Fatal("UpdateEventGroupApproval with zero event group ID should return error")
	}
	t.Logf("Zero event group ID error: %v", err)

	// Test with negative event group ID
	_, err = client.UpdateEventGroupApproval(ctx, -1, false, "Test reason")
	if err == nil {
		t.Fatal("UpdateEventGroupApproval with negative event group ID should return error")
	}
	t.Logf("Negative event group ID error: %v", err)
}

// TestSendExternalWebhook_InvalidInput tests webhook validation.
func TestSendExternalWebhook_InvalidInput(t *testing.T) {
	client := getTestClient(t)
	ctx := context.Background()

	// Test with empty webhook URL
	_, err := client.SendExternalWebhook(ctx, "", "POST", nil, nil)
	if err == nil {
		t.Fatal("SendExternalWebhook with empty URL should return error")
	}
	t.Logf("Empty URL error: %v", err)
}

// TestSendExternalWebhook_InvalidURL tests webhook with invalid URL.
func TestSendExternalWebhook_InvalidURL(t *testing.T) {
	client := getTestClient(t)
	ctx := context.Background()

	// Try to send webhook to invalid URL
	result, err := client.SendExternalWebhook(ctx, "http://invalid-url-that-doesnt-exist-12345.com", "POST", nil, nil)
	if err != nil {
		t.Logf("Invalid URL error (expected): %v", err)
	}

	if result != nil {
		t.Logf("Response: %s", result.String())
		if result.IsSuccessful() {
			t.Error("Should not succeed for invalid URL")
		}
	}
}

// TestConsumingServicesTestWebhook_InvalidInput tests webhook service test validation.
func TestConsumingServicesTestWebhook_InvalidInput(t *testing.T) {
	client := getTestClient(t)
	ctx := context.Background()

	// Test with empty service name
	_, err := client.ConsumingServicesTestWebhook(ctx, "", nil)
	if err == nil {
		t.Fatal("ConsumingServicesTestWebhook with empty service name should return error")
	}
	t.Logf("Empty service name error: %v", err)
}

// TestConsumingServicesTestWebhook_NonexistentService tests non-existent service.
func TestConsumingServicesTestWebhook_NonexistentService(t *testing.T) {
	client := getTestClient(t)
	ctx := context.Background()

	// Try to test a non-existent service
	result, err := client.ConsumingServicesTestWebhook(ctx, "nonexistent_service_12345", nil)
	if err != nil {
		t.Logf("Nonexistent service error (expected): %v", err)
	}

	if result != nil {
		t.Logf("Response: %s", result.String())
		if result.IsSuccessful() {
			t.Error("Should not succeed for nonexistent service")
		}
	}
}

// TestConsumingServicesTestLog_InvalidInput tests log service test validation.
func TestConsumingServicesTestLog_InvalidInput(t *testing.T) {
	client := getTestClient(t)
	ctx := context.Background()

	// Test with empty service name
	_, err := client.ConsumingServicesTestLog(ctx, "", nil)
	if err == nil {
		t.Fatal("ConsumingServicesTestLog with empty service name should return error")
	}
	t.Logf("Empty service name error: %v", err)
}

// TestConsumingServicesTestLog_NonexistentService tests non-existent log service.
func TestConsumingServicesTestLog_NonexistentService(t *testing.T) {
	client := getTestClient(t)
	ctx := context.Background()

	// Try to test a non-existent log service
	result, err := client.ConsumingServicesTestLog(ctx, "nonexistent_log_service_12345", nil)
	if err != nil {
		t.Logf("Nonexistent service error (expected): %v", err)
	}

	if result != nil {
		t.Logf("Response: %s", result.String())
		if result.IsSuccessful() {
			t.Error("Should not succeed for nonexistent service")
		}
	}
}
