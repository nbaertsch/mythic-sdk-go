package mythic

import (
	"context"

	"github.com/nbaertsch/mythic-sdk-go/pkg/mythic/types"
)

// EventingTriggerManual manually triggers an event group for execution.
// Event groups define automated workflows that can respond to various triggers.
//
// Manual triggers are useful for:
//   - Testing event workflows before enabling automatic triggers
//   - One-time execution of automation workflows
//   - Running workflows on-demand for specific objects
//   - Debugging and troubleshooting event logic
//
// Parameters:
//   - ctx: Context for cancellation and timeout
//   - eventGroupID: ID of the event group to trigger
//   - objectID: Optional object ID to pass to the workflow (0 if not needed)
//   - parameters: Optional parameters to pass to the workflow
//
// Returns:
//   - *types.EventTriggerResponse: Execution status and execution ID
//   - error: Error if the operation fails
//
// Example:
//
//	// Trigger an event group manually
//	params := map[string]interface{}{
//	    "target": "10.0.0.1",
//	    "action": "scan",
//	}
//	response, err := client.EventingTriggerManual(ctx, eventGroupID, callbackID, params)
//	if err != nil {
//	    return err
//	}
//	if response.IsSuccessful() {
//	    fmt.Printf("Event triggered: %s\n", response.String())
//	    fmt.Printf("Execution ID: %d\n", response.ExecutionID)
//	}
func (c *Client) EventingTriggerManual(ctx context.Context, eventGroupID, objectID int, parameters map[string]interface{}) (*types.EventTriggerResponse, error) {
	if err := c.EnsureAuthenticated(ctx); err != nil {
		return nil, err
	}

	// Validate input
	if eventGroupID <= 0 {
		return nil, WrapError("EventingTriggerManual", ErrInvalidInput, "event group ID must be positive")
	}

	var mutation struct {
		Response struct {
			Status       string `graphql:"status"`
			Message      string `graphql:"message"`
			Error        string `graphql:"error"`
			ExecutionID  int    `graphql:"execution_id"`
			EventGroupID int    `graphql:"event_group_id"`
		} `graphql:"eventingTriggerManual(event_group_id: $event_group_id, object_id: $object_id, parameters: $parameters)"`
	}

	variables := map[string]interface{}{
		"event_group_id": eventGroupID,
		"object_id":      objectID,
		"parameters":     parameters,
	}

	if err := c.executeMutation(ctx, &mutation, variables); err != nil {
		return nil, WrapError("EventingTriggerManual", err, "failed to trigger event")
	}

	response := &types.EventTriggerResponse{
		Status:       mutation.Response.Status,
		Message:      mutation.Response.Message,
		Error:        mutation.Response.Error,
		ExecutionID:  mutation.Response.ExecutionID,
		EventGroupID: mutation.Response.EventGroupID,
	}

	if !response.IsSuccessful() {
		return response, WrapError("EventingTriggerManual", ErrOperationFailed, response.Error)
	}

	return response, nil
}

// EventingTriggerManualBulk triggers an event group on multiple objects simultaneously.
// This is more efficient than triggering events one at a time for batch operations.
//
// Parameters:
//   - ctx: Context for cancellation and timeout
//   - eventGroupID: ID of the event group to trigger
//   - objectIDs: List of object IDs to pass to the workflow
//   - parameters: Optional parameters to pass to all workflows
//
// Returns:
//   - *types.EventTriggerResponse: Execution status
//   - error: Error if the operation fails
//
// Example:
//
//	// Trigger event on multiple callbacks
//	callbackIDs := []int{1, 2, 3, 4, 5}
//	response, err := client.EventingTriggerManualBulk(ctx, eventGroupID, callbackIDs, nil)
//	if err != nil {
//	    return err
//	}
//	fmt.Printf("Bulk trigger: %s\n", response.String())
func (c *Client) EventingTriggerManualBulk(ctx context.Context, eventGroupID int, objectIDs []int, parameters map[string]interface{}) (*types.EventTriggerResponse, error) {
	if err := c.EnsureAuthenticated(ctx); err != nil {
		return nil, err
	}

	// Validate input
	if eventGroupID <= 0 {
		return nil, WrapError("EventingTriggerManualBulk", ErrInvalidInput, "event group ID must be positive")
	}
	if len(objectIDs) == 0 {
		return nil, WrapError("EventingTriggerManualBulk", ErrInvalidInput, "object IDs list cannot be empty")
	}

	var mutation struct {
		Response struct {
			Status       string `graphql:"status"`
			Message      string `graphql:"message"`
			Error        string `graphql:"error"`
			ExecutionID  int    `graphql:"execution_id"`
			EventGroupID int    `graphql:"event_group_id"`
		} `graphql:"eventingTriggerManualBulk(event_group_id: $event_group_id, object_ids: $object_ids, parameters: $parameters)"`
	}

	variables := map[string]interface{}{
		"event_group_id": eventGroupID,
		"object_ids":     objectIDs,
		"parameters":     parameters,
	}

	if err := c.executeMutation(ctx, &mutation, variables); err != nil {
		return nil, WrapError("EventingTriggerManualBulk", err, "failed to trigger bulk event")
	}

	response := &types.EventTriggerResponse{
		Status:       mutation.Response.Status,
		Message:      mutation.Response.Message,
		Error:        mutation.Response.Error,
		ExecutionID:  mutation.Response.ExecutionID,
		EventGroupID: mutation.Response.EventGroupID,
	}

	if !response.IsSuccessful() {
		return response, WrapError("EventingTriggerManualBulk", ErrOperationFailed, response.Error)
	}

	return response, nil
}

// EventingTriggerKeyword triggers event groups by keyword match.
// Event groups can be tagged with keywords for easy triggering.
//
// Parameters:
//   - ctx: Context for cancellation and timeout
//   - keyword: Keyword to match against event group keywords
//   - objectID: Optional object ID to pass to workflows (0 if not needed)
//   - parameters: Optional parameters to pass to workflows
//
// Returns:
//   - *types.EventTriggerResponse: Execution status
//   - error: Error if the operation fails
//
// Example:
//
//	// Trigger all events tagged with "scan"
//	response, err := client.EventingTriggerKeyword(ctx, "scan", targetID, nil)
//	if err != nil {
//	    return err
//	}
func (c *Client) EventingTriggerKeyword(ctx context.Context, keyword string, objectID int, parameters map[string]interface{}) (*types.EventTriggerResponse, error) {
	if err := c.EnsureAuthenticated(ctx); err != nil {
		return nil, err
	}

	// Validate input
	if keyword == "" {
		return nil, WrapError("EventingTriggerKeyword", ErrInvalidInput, "keyword cannot be empty")
	}

	var mutation struct {
		Response struct {
			Status       string `graphql:"status"`
			Message      string `graphql:"message"`
			Error        string `graphql:"error"`
			ExecutionID  int    `graphql:"execution_id"`
			EventGroupID int    `graphql:"event_group_id"`
		} `graphql:"eventingTriggerKeyword(keyword: $keyword, object_id: $object_id, parameters: $parameters)"`
	}

	variables := map[string]interface{}{
		"keyword":    keyword,
		"object_id":  objectID,
		"parameters": parameters,
	}

	if err := c.executeMutation(ctx, &mutation, variables); err != nil {
		return nil, WrapError("EventingTriggerKeyword", err, "failed to trigger keyword event")
	}

	response := &types.EventTriggerResponse{
		Status:       mutation.Response.Status,
		Message:      mutation.Response.Message,
		Error:        mutation.Response.Error,
		ExecutionID:  mutation.Response.ExecutionID,
		EventGroupID: mutation.Response.EventGroupID,
	}

	if !response.IsSuccessful() {
		return response, WrapError("EventingTriggerKeyword", ErrOperationFailed, response.Error)
	}

	return response, nil
}

// EventingTriggerCancel cancels a running event execution.
//
// Parameters:
//   - ctx: Context for cancellation and timeout
//   - executionID: ID of the execution to cancel
//
// Returns:
//   - *types.EventTriggerResponse: Cancellation status
//   - error: Error if the operation fails
//
// Example:
//
//	response, err := client.EventingTriggerCancel(ctx, executionID)
//	if err != nil {
//	    return err
//	}
func (c *Client) EventingTriggerCancel(ctx context.Context, executionID int) (*types.EventTriggerResponse, error) {
	if err := c.EnsureAuthenticated(ctx); err != nil {
		return nil, err
	}

	// Validate input
	if executionID <= 0 {
		return nil, WrapError("EventingTriggerCancel", ErrInvalidInput, "execution ID must be positive")
	}

	var mutation struct {
		Response struct {
			Status  string `graphql:"status"`
			Message string `graphql:"message"`
			Error   string `graphql:"error"`
		} `graphql:"eventingTriggerCancel(execution_id: $execution_id)"`
	}

	variables := map[string]interface{}{
		"execution_id": executionID,
	}

	if err := c.executeMutation(ctx, &mutation, variables); err != nil {
		return nil, WrapError("EventingTriggerCancel", err, "failed to cancel event")
	}

	response := &types.EventTriggerResponse{
		Status:  mutation.Response.Status,
		Message: mutation.Response.Message,
		Error:   mutation.Response.Error,
	}

	if !response.IsSuccessful() {
		return response, WrapError("EventingTriggerCancel", ErrOperationFailed, response.Error)
	}

	return response, nil
}

// EventingTriggerRetry retries a failed event execution from the beginning.
//
// Parameters:
//   - ctx: Context for cancellation and timeout
//   - executionID: ID of the execution to retry
//
// Returns:
//   - *types.EventTriggerResponse: Retry status with new execution ID
//   - error: Error if the operation fails
//
// Example:
//
//	response, err := client.EventingTriggerRetry(ctx, failedExecutionID)
//	if err != nil {
//	    return err
//	}
func (c *Client) EventingTriggerRetry(ctx context.Context, executionID int) (*types.EventTriggerResponse, error) {
	if err := c.EnsureAuthenticated(ctx); err != nil {
		return nil, err
	}

	// Validate input
	if executionID <= 0 {
		return nil, WrapError("EventingTriggerRetry", ErrInvalidInput, "execution ID must be positive")
	}

	var mutation struct {
		Response struct {
			Status      string `graphql:"status"`
			Message     string `graphql:"message"`
			Error       string `graphql:"error"`
			ExecutionID int    `graphql:"execution_id"`
		} `graphql:"eventingTriggerRetry(execution_id: $execution_id)"`
	}

	variables := map[string]interface{}{
		"execution_id": executionID,
	}

	if err := c.executeMutation(ctx, &mutation, variables); err != nil {
		return nil, WrapError("EventingTriggerRetry", err, "failed to retry event")
	}

	response := &types.EventTriggerResponse{
		Status:      mutation.Response.Status,
		Message:     mutation.Response.Message,
		Error:       mutation.Response.Error,
		ExecutionID: mutation.Response.ExecutionID,
	}

	if !response.IsSuccessful() {
		return response, WrapError("EventingTriggerRetry", ErrOperationFailed, response.Error)
	}

	return response, nil
}

// EventingTriggerRetryFromStep retries a failed event execution from a specific step.
//
// Parameters:
//   - ctx: Context for cancellation and timeout
//   - executionID: ID of the execution to retry
//   - stepNumber: Step number to retry from (1-based)
//
// Returns:
//   - *types.EventTriggerResponse: Retry status with new execution ID
//   - error: Error if the operation fails
//
// Example:
//
//	// Retry from step 3
//	response, err := client.EventingTriggerRetryFromStep(ctx, executionID, 3)
//	if err != nil {
//	    return err
//	}
func (c *Client) EventingTriggerRetryFromStep(ctx context.Context, executionID, stepNumber int) (*types.EventTriggerResponse, error) {
	if err := c.EnsureAuthenticated(ctx); err != nil {
		return nil, err
	}

	// Validate input
	if executionID <= 0 {
		return nil, WrapError("EventingTriggerRetryFromStep", ErrInvalidInput, "execution ID must be positive")
	}
	if stepNumber <= 0 {
		return nil, WrapError("EventingTriggerRetryFromStep", ErrInvalidInput, "step number must be positive")
	}

	var mutation struct {
		Response struct {
			Status      string `graphql:"status"`
			Message     string `graphql:"message"`
			Error       string `graphql:"error"`
			ExecutionID int    `graphql:"execution_id"`
		} `graphql:"eventingTriggerRetryFromStep(execution_id: $execution_id, step_number: $step_number)"`
	}

	variables := map[string]interface{}{
		"execution_id": executionID,
		"step_number":  stepNumber,
	}

	if err := c.executeMutation(ctx, &mutation, variables); err != nil {
		return nil, WrapError("EventingTriggerRetryFromStep", err, "failed to retry event from step")
	}

	response := &types.EventTriggerResponse{
		Status:      mutation.Response.Status,
		Message:     mutation.Response.Message,
		Error:       mutation.Response.Error,
		ExecutionID: mutation.Response.ExecutionID,
	}

	if !response.IsSuccessful() {
		return response, WrapError("EventingTriggerRetryFromStep", ErrOperationFailed, response.Error)
	}

	return response, nil
}

// EventingTriggerRunAgain re-runs a completed event execution with the same parameters.
//
// Parameters:
//   - ctx: Context for cancellation and timeout
//   - executionID: ID of the execution to run again
//
// Returns:
//   - *types.EventTriggerResponse: Execution status with new execution ID
//   - error: Error if the operation fails
//
// Example:
//
//	response, err := client.EventingTriggerRunAgain(ctx, completedExecutionID)
//	if err != nil {
//	    return err
//	}
func (c *Client) EventingTriggerRunAgain(ctx context.Context, executionID int) (*types.EventTriggerResponse, error) {
	if err := c.EnsureAuthenticated(ctx); err != nil {
		return nil, err
	}

	// Validate input
	if executionID <= 0 {
		return nil, WrapError("EventingTriggerRunAgain", ErrInvalidInput, "execution ID must be positive")
	}

	var mutation struct {
		Response struct {
			Status      string `graphql:"status"`
			Message     string `graphql:"message"`
			Error       string `graphql:"error"`
			ExecutionID int    `graphql:"execution_id"`
		} `graphql:"eventingTriggerRunAgain(execution_id: $execution_id)"`
	}

	variables := map[string]interface{}{
		"execution_id": executionID,
	}

	if err := c.executeMutation(ctx, &mutation, variables); err != nil {
		return nil, WrapError("EventingTriggerRunAgain", err, "failed to run event again")
	}

	response := &types.EventTriggerResponse{
		Status:      mutation.Response.Status,
		Message:     mutation.Response.Message,
		Error:       mutation.Response.Error,
		ExecutionID: mutation.Response.ExecutionID,
	}

	if !response.IsSuccessful() {
		return response, WrapError("EventingTriggerRunAgain", ErrOperationFailed, response.Error)
	}

	return response, nil
}

// EventingTriggerUpdate updates an event group's configuration.
//
// Parameters:
//   - ctx: Context for cancellation and timeout
//   - eventGroupID: ID of the event group to update
//   - name: Optional new name
//   - description: Optional new description
//   - active: Optional active status
//   - requiresApproval: Optional approval requirement
//   - conditions: Optional new conditions
//   - actions: Optional new actions
//   - keywords: Optional new keywords
//
// Returns:
//   - *types.EventTriggerResponse: Update status
//   - error: Error if the operation fails
//
// Example:
//
//	active := true
//	response, err := client.EventingTriggerUpdate(ctx, eventGroupID, "New Name", "", &active, nil, nil, nil, nil)
//	if err != nil {
//	    return err
//	}
func (c *Client) EventingTriggerUpdate(ctx context.Context, eventGroupID int, name, description string, active, requiresApproval *bool, conditions map[string]interface{}, actions []interface{}, keywords []string) (*types.EventTriggerResponse, error) {
	if err := c.EnsureAuthenticated(ctx); err != nil {
		return nil, err
	}

	// Validate input
	if eventGroupID <= 0 {
		return nil, WrapError("EventingTriggerUpdate", ErrInvalidInput, "event group ID must be positive")
	}

	var mutation struct {
		Response struct {
			Status  string `graphql:"status"`
			Message string `graphql:"message"`
			Error   string `graphql:"error"`
		} `graphql:"eventingTriggerUpdate(event_group_id: $event_group_id, name: $name, description: $description, active: $active, requires_approval: $requires_approval, conditions: $conditions, actions: $actions, keywords: $keywords)"`
	}

	variables := map[string]interface{}{
		"event_group_id":    eventGroupID,
		"name":              name,
		"description":       description,
		"active":            active,
		"requires_approval": requiresApproval,
		"conditions":        conditions,
		"actions":           actions,
		"keywords":          keywords,
	}

	if err := c.executeMutation(ctx, &mutation, variables); err != nil {
		return nil, WrapError("EventingTriggerUpdate", err, "failed to update event group")
	}

	response := &types.EventTriggerResponse{
		Status:  mutation.Response.Status,
		Message: mutation.Response.Message,
		Error:   mutation.Response.Error,
	}

	if !response.IsSuccessful() {
		return response, WrapError("EventingTriggerUpdate", ErrOperationFailed, response.Error)
	}

	return response, nil
}

// EventingExportWorkflow exports a workflow definition for backup or sharing.
//
// Parameters:
//   - ctx: Context for cancellation and timeout
//   - workflowID: ID of the workflow to export
//
// Returns:
//   - *types.WorkflowExportResponse: Exported workflow definition
//   - error: Error if the operation fails
//
// Example:
//
//	response, err := client.EventingExportWorkflow(ctx, workflowID)
//	if err != nil {
//	    return err
//	}
//	if response.IsSuccessful() {
//	    // Save response.Definition to file
//	}
func (c *Client) EventingExportWorkflow(ctx context.Context, workflowID int) (*types.WorkflowExportResponse, error) {
	if err := c.EnsureAuthenticated(ctx); err != nil {
		return nil, err
	}

	// Validate input
	if workflowID <= 0 {
		return nil, WrapError("EventingExportWorkflow", ErrInvalidInput, "workflow ID must be positive")
	}

	var query struct {
		Response struct {
			Status       string                 `graphql:"status"`
			WorkflowName string                 `graphql:"workflow_name"`
			Definition   map[string]interface{} `graphql:"definition"`
			Error        string                 `graphql:"error"`
		} `graphql:"eventingExportWorkflow(workflow_id: $workflow_id)"`
	}

	variables := map[string]interface{}{
		"workflow_id": workflowID,
	}

	if err := c.executeQuery(ctx, &query, variables); err != nil {
		return nil, WrapError("EventingExportWorkflow", err, "failed to export workflow")
	}

	response := &types.WorkflowExportResponse{
		Status:       query.Response.Status,
		WorkflowName: query.Response.WorkflowName,
		Definition:   query.Response.Definition,
		Error:        query.Response.Error,
	}

	if !response.IsSuccessful() {
		return response, WrapError("EventingExportWorkflow", ErrOperationFailed, response.Error)
	}

	return response, nil
}

// EventingImportContainerWorkflow imports a workflow from a container.
//
// Parameters:
//   - ctx: Context for cancellation and timeout
//   - containerName: Name of the container (payload type or C2 profile)
//   - workflowFile: Path to the workflow file in the container
//   - operationID: Optional operation ID (0 to use current)
//
// Returns:
//   - *types.WorkflowImportResponse: Import status and workflow ID
//   - error: Error if the operation fails
//
// Example:
//
//	response, err := client.EventingImportContainerWorkflow(ctx, "apollo", "workflows/scan.json", 0)
//	if err != nil {
//	    return err
//	}
//	if response.IsSuccessful() {
//	    fmt.Printf("Imported workflow ID: %d\n", response.WorkflowID)
//	}
func (c *Client) EventingImportContainerWorkflow(ctx context.Context, containerName, workflowFile string, operationID int) (*types.WorkflowImportResponse, error) {
	if err := c.EnsureAuthenticated(ctx); err != nil {
		return nil, err
	}

	// Validate input
	if containerName == "" {
		return nil, WrapError("EventingImportContainerWorkflow", ErrInvalidInput, "container name cannot be empty")
	}
	if workflowFile == "" {
		return nil, WrapError("EventingImportContainerWorkflow", ErrInvalidInput, "workflow file cannot be empty")
	}

	// Use current operation if not specified
	if operationID == 0 {
		opID := c.GetCurrentOperation()
		if opID == nil {
			return nil, WrapError("EventingImportContainerWorkflow", ErrNotAuthenticated, "no current operation set")
		}
		operationID = *opID
	}

	var mutation struct {
		Response struct {
			Status     string `graphql:"status"`
			Message    string `graphql:"message"`
			Error      string `graphql:"error"`
			WorkflowID int    `graphql:"workflow_id"`
		} `graphql:"eventingImportContainerWorkflow(container_name: $container_name, workflow_file: $workflow_file, operation_id: $operation_id)"`
	}

	variables := map[string]interface{}{
		"container_name": containerName,
		"workflow_file":  workflowFile,
		"operation_id":   operationID,
	}

	if err := c.executeMutation(ctx, &mutation, variables); err != nil {
		return nil, WrapError("EventingImportContainerWorkflow", err, "failed to import workflow")
	}

	response := &types.WorkflowImportResponse{
		Status:     mutation.Response.Status,
		Message:    mutation.Response.Message,
		Error:      mutation.Response.Error,
		WorkflowID: mutation.Response.WorkflowID,
	}

	if !response.IsSuccessful() {
		return response, WrapError("EventingImportContainerWorkflow", ErrOperationFailed, response.Error)
	}

	return response, nil
}

// EventingTestFile tests a workflow file for validity without importing it.
//
// Parameters:
//   - ctx: Context for cancellation and timeout
//   - workflowFile: Path to the workflow file to test
//   - testData: Optional test data to validate against
//
// Returns:
//   - *types.WorkflowTestResponse: Validation result with errors if any
//   - error: Error if the operation fails
//
// Example:
//
//	response, err := client.EventingTestFile(ctx, "workflows/scan.json", nil)
//	if err != nil {
//	    return err
//	}
//	if !response.IsValid() {
//	    for _, err := range response.Errors {
//	        fmt.Printf("Validation error: %s\n", err)
//	    }
//	}
func (c *Client) EventingTestFile(ctx context.Context, workflowFile string, testData map[string]interface{}) (*types.WorkflowTestResponse, error) {
	if err := c.EnsureAuthenticated(ctx); err != nil {
		return nil, err
	}

	// Validate input
	if workflowFile == "" {
		return nil, WrapError("EventingTestFile", ErrInvalidInput, "workflow file cannot be empty")
	}

	var query struct {
		Response struct {
			Status  string   `graphql:"status"`
			Valid   bool     `graphql:"valid"`
			Errors  []string `graphql:"errors"`
			Message string   `graphql:"message"`
		} `graphql:"eventingTestFile(workflow_file: $workflow_file, test_data: $test_data)"`
	}

	variables := map[string]interface{}{
		"workflow_file": workflowFile,
		"test_data":     testData,
	}

	if err := c.executeQuery(ctx, &query, variables); err != nil {
		return nil, WrapError("EventingTestFile", err, "failed to test workflow file")
	}

	response := &types.WorkflowTestResponse{
		Status:  query.Response.Status,
		Valid:   query.Response.Valid,
		Errors:  query.Response.Errors,
		Message: query.Response.Message,
	}

	return response, nil
}

// UpdateEventGroupApproval approves or rejects an event group execution.
// Event groups can require approval before execution for security.
//
// Parameters:
//   - ctx: Context for cancellation and timeout
//   - eventGroupID: ID of the event group to approve/reject
//   - approved: True to approve, false to reject
//   - reason: Optional reason for approval/rejection
//
// Returns:
//   - *types.EventTriggerResponse: Approval status
//   - error: Error if the operation fails
//
// Example:
//
//	// Approve an event group
//	response, err := client.UpdateEventGroupApproval(ctx, eventGroupID, true, "Verified safe")
//	if err != nil {
//	    return err
//	}
func (c *Client) UpdateEventGroupApproval(ctx context.Context, eventGroupID int, approved bool, reason string) (*types.EventTriggerResponse, error) {
	if err := c.EnsureAuthenticated(ctx); err != nil {
		return nil, err
	}

	// Validate input
	if eventGroupID <= 0 {
		return nil, WrapError("UpdateEventGroupApproval", ErrInvalidInput, "event group ID must be positive")
	}

	var mutation struct {
		Response struct {
			Status  string `graphql:"status"`
			Message string `graphql:"message"`
			Error   string `graphql:"error"`
		} `graphql:"updateEventGroupApproval(event_group_id: $event_group_id, approved: $approved, reason: $reason)"`
	}

	variables := map[string]interface{}{
		"event_group_id": eventGroupID,
		"approved":       approved,
		"reason":         reason,
	}

	if err := c.executeMutation(ctx, &mutation, variables); err != nil {
		return nil, WrapError("UpdateEventGroupApproval", err, "failed to update event approval")
	}

	response := &types.EventTriggerResponse{
		Status:  mutation.Response.Status,
		Message: mutation.Response.Message,
		Error:   mutation.Response.Error,
	}

	if !response.IsSuccessful() {
		return response, WrapError("UpdateEventGroupApproval", ErrOperationFailed, response.Error)
	}

	return response, nil
}

// SendExternalWebhook sends a webhook to an external service.
//
// Parameters:
//   - ctx: Context for cancellation and timeout
//   - webhookURL: URL to send the webhook to
//   - method: HTTP method (GET, POST, PUT, DELETE)
//   - headers: Optional HTTP headers
//   - body: Optional request body
//
// Returns:
//   - *types.WebhookResponse: Webhook response with status code
//   - error: Error if the operation fails
//
// Example:
//
//	headers := map[string]string{"Authorization": "Bearer token123"}
//	body := map[string]interface{}{"event": "callback_created", "callback_id": 42}
//	response, err := client.SendExternalWebhook(ctx, "https://api.example.com/webhook", "POST", headers, body)
//	if err != nil {
//	    return err
//	}
func (c *Client) SendExternalWebhook(ctx context.Context, webhookURL, method string, headers map[string]string, body map[string]interface{}) (*types.WebhookResponse, error) {
	if err := c.EnsureAuthenticated(ctx); err != nil {
		return nil, err
	}

	// Validate input
	if webhookURL == "" {
		return nil, WrapError("SendExternalWebhook", ErrInvalidInput, "webhook URL cannot be empty")
	}
	if method == "" {
		method = "POST" // Default to POST
	}

	var mutation struct {
		Response struct {
			Status     string `graphql:"status"`
			StatusCode int    `graphql:"status_code"`
			Response   string `graphql:"response"`
			Error      string `graphql:"error"`
		} `graphql:"sendExternalWebhook(webhook_url: $webhook_url, method: $method, headers: $headers, body: $body)"`
	}

	variables := map[string]interface{}{
		"webhook_url": webhookURL,
		"method":      method,
		"headers":     headers,
		"body":        body,
	}

	if err := c.executeMutation(ctx, &mutation, variables); err != nil {
		return nil, WrapError("SendExternalWebhook", err, "failed to send webhook")
	}

	response := &types.WebhookResponse{
		Status:     mutation.Response.Status,
		StatusCode: mutation.Response.StatusCode,
		Response:   mutation.Response.Response,
		Error:      mutation.Response.Error,
	}

	if !response.IsSuccessful() {
		return response, WrapError("SendExternalWebhook", ErrOperationFailed, response.Error)
	}

	return response, nil
}

// ConsumingServicesTestWebhook tests a consuming service webhook configuration.
//
// Parameters:
//   - ctx: Context for cancellation and timeout
//   - serviceName: Name of the consuming service to test
//   - testData: Optional test data to send
//
// Returns:
//   - *types.ConsumingServiceTestResponse: Test result
//   - error: Error if the operation fails
//
// Example:
//
//	testData := map[string]interface{}{"test": "payload"}
//	response, err := client.ConsumingServicesTestWebhook(ctx, "slack-notifications", testData)
//	if err != nil {
//	    return err
//	}
func (c *Client) ConsumingServicesTestWebhook(ctx context.Context, serviceName string, testData map[string]interface{}) (*types.ConsumingServiceTestResponse, error) {
	if err := c.EnsureAuthenticated(ctx); err != nil {
		return nil, err
	}

	// Validate input
	if serviceName == "" {
		return nil, WrapError("ConsumingServicesTestWebhook", ErrInvalidInput, "service name cannot be empty")
	}

	var mutation struct {
		Response struct {
			Status  string `graphql:"status"`
			Message string `graphql:"message"`
			Error   string `graphql:"error"`
		} `graphql:"consumingServicesTestWebhook(service_name: $service_name, test_data: $test_data)"`
	}

	variables := map[string]interface{}{
		"service_name": serviceName,
		"test_data":    testData,
	}

	if err := c.executeMutation(ctx, &mutation, variables); err != nil {
		return nil, WrapError("ConsumingServicesTestWebhook", err, "failed to test webhook service")
	}

	response := &types.ConsumingServiceTestResponse{
		Status:  mutation.Response.Status,
		Message: mutation.Response.Message,
		Error:   mutation.Response.Error,
	}

	if !response.IsSuccessful() {
		return response, WrapError("ConsumingServicesTestWebhook", ErrOperationFailed, response.Error)
	}

	return response, nil
}

// ConsumingServicesTestLog tests a consuming service logging configuration.
//
// Parameters:
//   - ctx: Context for cancellation and timeout
//   - serviceName: Name of the consuming service to test
//   - testData: Optional test data to log
//
// Returns:
//   - *types.ConsumingServiceTestResponse: Test result
//   - error: Error if the operation fails
//
// Example:
//
//	testData := map[string]interface{}{"message": "test log entry"}
//	response, err := client.ConsumingServicesTestLog(ctx, "elasticsearch", testData)
//	if err != nil {
//	    return err
//	}
func (c *Client) ConsumingServicesTestLog(ctx context.Context, serviceName string, testData map[string]interface{}) (*types.ConsumingServiceTestResponse, error) {
	if err := c.EnsureAuthenticated(ctx); err != nil {
		return nil, err
	}

	// Validate input
	if serviceName == "" {
		return nil, WrapError("ConsumingServicesTestLog", ErrInvalidInput, "service name cannot be empty")
	}

	var mutation struct {
		Response struct {
			Status  string `graphql:"status"`
			Message string `graphql:"message"`
			Error   string `graphql:"error"`
		} `graphql:"consumingServicesTestLog(service_name: $service_name, test_data: $test_data)"`
	}

	variables := map[string]interface{}{
		"service_name": serviceName,
		"test_data":    testData,
	}

	if err := c.executeMutation(ctx, &mutation, variables); err != nil {
		return nil, WrapError("ConsumingServicesTestLog", err, "failed to test log service")
	}

	response := &types.ConsumingServiceTestResponse{
		Status:  mutation.Response.Status,
		Message: mutation.Response.Message,
		Error:   mutation.Response.Error,
	}

	if !response.IsSuccessful() {
		return response, WrapError("ConsumingServicesTestLog", ErrOperationFailed, response.Error)
	}

	return response, nil
}
