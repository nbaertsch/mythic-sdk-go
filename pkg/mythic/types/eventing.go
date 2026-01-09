package types

import "fmt"

// EventGroup represents an event group that can be triggered.
type EventGroup struct {
	ID               int                    `json:"id"`
	Name             string                 `json:"name"`
	Description      string                 `json:"description"`
	TriggerType      string                 `json:"trigger_type"`
	OperationID      int                    `json:"operation_id"`
	Active           bool                   `json:"active"`
	Deleted          bool                   `json:"deleted"`
	RequiresApproval bool                   `json:"requires_approval"`
	Approved         bool                   `json:"approved"`
	Keywords         []string               `json:"keywords,omitempty"`
	Conditions       map[string]interface{} `json:"conditions,omitempty"`
	Actions          []interface{}          `json:"actions,omitempty"`
	Status           string                 `json:"status,omitempty"`
	LastTriggered    string                 `json:"last_triggered,omitempty"`
	ExecutionCount   int                    `json:"execution_count"`
	FailureCount     int                    `json:"failure_count"`
	CurrentStep      int                    `json:"current_step,omitempty"`
	TotalSteps       int                    `json:"total_steps,omitempty"`
}

// String returns a human-readable representation of the event group.
func (e *EventGroup) String() string {
	status := "active"
	if e.Deleted {
		status = "deleted"
	} else if !e.Active {
		status = "inactive"
	}
	return fmt.Sprintf("EventGroup '%s' (%s, trigger: %s)", e.Name, status, e.TriggerType)
}

// IsActive returns true if the event group is active and not deleted.
func (e *EventGroup) IsActive() bool {
	return e.Active && !e.Deleted
}

// IsDeleted returns true if the event group has been deleted.
func (e *EventGroup) IsDeleted() bool {
	return e.Deleted
}

// NeedsApproval returns true if the event group requires approval before execution.
func (e *EventGroup) NeedsApproval() bool {
	return e.RequiresApproval && !e.Approved
}

// EventTriggerManualRequest represents a manual event trigger request.
type EventTriggerManualRequest struct {
	EventGroupID int                    `json:"event_group_id"`
	ObjectID     int                    `json:"object_id,omitempty"`
	Parameters   map[string]interface{} `json:"parameters,omitempty"`
}

// String returns a human-readable representation of the request.
func (e *EventTriggerManualRequest) String() string {
	if e.ObjectID > 0 {
		return fmt.Sprintf("Trigger event group %d on object %d", e.EventGroupID, e.ObjectID)
	}
	return fmt.Sprintf("Trigger event group %d", e.EventGroupID)
}

// EventTriggerResponse represents the response from triggering an event.
type EventTriggerResponse struct {
	Status       string `json:"status"`
	Message      string `json:"message,omitempty"`
	Error        string `json:"error,omitempty"`
	ExecutionID  int    `json:"execution_id,omitempty"`
	EventGroupID int    `json:"event_group_id,omitempty"`
}

// String returns a human-readable representation of the response.
func (e *EventTriggerResponse) String() string {
	if e.Status == "success" {
		if e.Message != "" {
			return e.Message
		}
		if e.ExecutionID > 0 {
			return fmt.Sprintf("Event triggered successfully (execution ID: %d)", e.ExecutionID)
		}
		return "Event triggered successfully"
	}
	return fmt.Sprintf("Failed to trigger event: %s", e.Error)
}

// IsSuccessful returns true if the trigger succeeded.
func (e *EventTriggerResponse) IsSuccessful() bool {
	return e.Status == "success"
}

// EventTriggerBulkRequest represents a bulk event trigger request.
type EventTriggerBulkRequest struct {
	EventGroupID int                    `json:"event_group_id"`
	ObjectIDs    []int                  `json:"object_ids"`
	Parameters   map[string]interface{} `json:"parameters,omitempty"`
}

// String returns a human-readable representation of the request.
func (e *EventTriggerBulkRequest) String() string {
	return fmt.Sprintf("Trigger event group %d on %d objects", e.EventGroupID, len(e.ObjectIDs))
}

// EventTriggerKeywordRequest represents a keyword-based trigger request.
type EventTriggerKeywordRequest struct {
	Keyword    string                 `json:"keyword"`
	ObjectID   int                    `json:"object_id,omitempty"`
	Parameters map[string]interface{} `json:"parameters,omitempty"`
}

// String returns a human-readable representation of the request.
func (e *EventTriggerKeywordRequest) String() string {
	return fmt.Sprintf("Trigger events with keyword '%s'", e.Keyword)
}

// EventControlRequest represents a request to control event execution (cancel, retry, etc.).
type EventControlRequest struct {
	ExecutionID int `json:"execution_id"`
	StepNumber  int `json:"step_number,omitempty"`
}

// String returns a human-readable representation of the request.
func (e *EventControlRequest) String() string {
	if e.StepNumber > 0 {
		return fmt.Sprintf("Control execution %d at step %d", e.ExecutionID, e.StepNumber)
	}
	return fmt.Sprintf("Control execution %d", e.ExecutionID)
}

// EventGroupUpdateRequest represents a request to update event group configuration.
type EventGroupUpdateRequest struct {
	EventGroupID     int                    `json:"event_group_id"`
	Name             string                 `json:"name,omitempty"`
	Description      string                 `json:"description,omitempty"`
	Active           *bool                  `json:"active,omitempty"`
	RequiresApproval *bool                  `json:"requires_approval,omitempty"`
	Conditions       map[string]interface{} `json:"conditions,omitempty"`
	Actions          []interface{}          `json:"actions,omitempty"`
	Keywords         []string               `json:"keywords,omitempty"`
}

// String returns a human-readable representation of the request.
func (e *EventGroupUpdateRequest) String() string {
	return fmt.Sprintf("Update event group %d", e.EventGroupID)
}

// WorkflowExportResponse represents exported workflow definition.
type WorkflowExportResponse struct {
	Status       string                 `json:"status"`
	WorkflowName string                 `json:"workflow_name,omitempty"`
	Definition   map[string]interface{} `json:"definition,omitempty"`
	Error        string                 `json:"error,omitempty"`
}

// String returns a human-readable representation of the response.
func (w *WorkflowExportResponse) String() string {
	if w.Status == "success" {
		return fmt.Sprintf("Exported workflow: %s", w.WorkflowName)
	}
	return fmt.Sprintf("Failed to export workflow: %s", w.Error)
}

// IsSuccessful returns true if the export succeeded.
func (w *WorkflowExportResponse) IsSuccessful() bool {
	return w.Status == "success"
}

// WorkflowImportRequest represents a workflow import request.
type WorkflowImportRequest struct {
	ContainerName string `json:"container_name"`
	WorkflowFile  string `json:"workflow_file"`
	OperationID   int    `json:"operation_id,omitempty"`
}

// String returns a human-readable representation of the request.
func (w *WorkflowImportRequest) String() string {
	return fmt.Sprintf("Import workflow from %s/%s", w.ContainerName, w.WorkflowFile)
}

// WorkflowImportResponse represents the response from importing a workflow.
type WorkflowImportResponse struct {
	Status     string `json:"status"`
	Message    string `json:"message,omitempty"`
	Error      string `json:"error,omitempty"`
	WorkflowID int    `json:"workflow_id,omitempty"`
}

// String returns a human-readable representation of the response.
func (w *WorkflowImportResponse) String() string {
	if w.Status == "success" {
		if w.Message != "" {
			return w.Message
		}
		return fmt.Sprintf("Workflow imported successfully (ID: %d)", w.WorkflowID)
	}
	return fmt.Sprintf("Failed to import workflow: %s", w.Error)
}

// IsSuccessful returns true if the import succeeded.
func (w *WorkflowImportResponse) IsSuccessful() bool {
	return w.Status == "success"
}

// WorkflowTestRequest represents a workflow test request.
type WorkflowTestRequest struct {
	WorkflowFile string                 `json:"workflow_file"`
	TestData     map[string]interface{} `json:"test_data,omitempty"`
}

// String returns a human-readable representation of the request.
func (w *WorkflowTestRequest) String() string {
	return fmt.Sprintf("Test workflow file: %s", w.WorkflowFile)
}

// WorkflowTestResponse represents the response from testing a workflow.
type WorkflowTestResponse struct {
	Status  string   `json:"status"`
	Valid   bool     `json:"valid"`
	Errors  []string `json:"errors,omitempty"`
	Message string   `json:"message,omitempty"`
}

// String returns a human-readable representation of the response.
func (w *WorkflowTestResponse) String() string {
	if w.Valid {
		return "Workflow is valid"
	}
	return fmt.Sprintf("Workflow validation failed: %d errors", len(w.Errors))
}

// IsValid returns true if the workflow is valid.
func (w *WorkflowTestResponse) IsValid() bool {
	return w.Valid
}

// HasErrors returns true if there are validation errors.
func (w *WorkflowTestResponse) HasErrors() bool {
	return len(w.Errors) > 0
}

// EventApprovalRequest represents a request to approve/reject event execution.
type EventApprovalRequest struct {
	EventGroupID int    `json:"event_group_id"`
	Approved     bool   `json:"approved"`
	Reason       string `json:"reason,omitempty"`
}

// String returns a human-readable representation of the request.
func (e *EventApprovalRequest) String() string {
	action := "Approve"
	if !e.Approved {
		action = "Reject"
	}
	return fmt.Sprintf("%s event group %d", action, e.EventGroupID)
}

// WebhookRequest represents a request to send an external webhook.
type WebhookRequest struct {
	WebhookURL string                 `json:"webhook_url"`
	Method     string                 `json:"method"` // GET, POST, PUT, DELETE
	Headers    map[string]string      `json:"headers,omitempty"`
	Body       map[string]interface{} `json:"body,omitempty"`
}

// String returns a human-readable representation of the request.
func (w *WebhookRequest) String() string {
	return fmt.Sprintf("%s webhook to %s", w.Method, w.WebhookURL)
}

// WebhookResponse represents the response from sending a webhook.
type WebhookResponse struct {
	Status     string `json:"status"`
	StatusCode int    `json:"status_code,omitempty"`
	Response   string `json:"response,omitempty"`
	Error      string `json:"error,omitempty"`
}

// String returns a human-readable representation of the response.
func (w *WebhookResponse) String() string {
	if w.Status == "success" {
		return fmt.Sprintf("Webhook sent successfully (HTTP %d)", w.StatusCode)
	}
	return fmt.Sprintf("Webhook failed: %s", w.Error)
}

// IsSuccessful returns true if the webhook was sent successfully.
func (w *WebhookResponse) IsSuccessful() bool {
	return w.Status == "success"
}

// ConsumingServiceTestRequest represents a request to test a consuming service.
type ConsumingServiceTestRequest struct {
	ServiceName string                 `json:"service_name"`
	TestData    map[string]interface{} `json:"test_data,omitempty"`
}

// String returns a human-readable representation of the request.
func (c *ConsumingServiceTestRequest) String() string {
	return fmt.Sprintf("Test consuming service: %s", c.ServiceName)
}

// ConsumingServiceTestResponse represents the response from testing a service.
type ConsumingServiceTestResponse struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
	Error   string `json:"error,omitempty"`
}

// String returns a human-readable representation of the response.
func (c *ConsumingServiceTestResponse) String() string {
	if c.Status == "success" {
		if c.Message != "" {
			return c.Message
		}
		return "Service test passed"
	}
	return fmt.Sprintf("Service test failed: %s", c.Error)
}

// IsSuccessful returns true if the test passed.
func (c *ConsumingServiceTestResponse) IsSuccessful() bool {
	return c.Status == "success"
}
