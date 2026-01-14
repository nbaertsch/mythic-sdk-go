package mythic

import (
	"context"
	"fmt"
	"time"
)

// Task represents a Mythic task.
type Task struct {
	ID                        int       `json:"id"`
	DisplayID                 int       `json:"display_id"`
	AgentTaskID               string    `json:"agent_task_id"`
	CommandName               string    `json:"command_name"`
	Params                    string    `json:"params"`
	DisplayParams             string    `json:"display_params"`
	OriginalParams            string    `json:"original_params"`
	Status                    string    `json:"status"`
	Completed                 bool      `json:"completed"`
	Comment                   string    `json:"comment"`
	Timestamp                 time.Time `json:"timestamp"`
	CallbackID                int       `json:"callback_id"`
	OperatorID                int       `json:"operator_id"`
	OperationID               int       `json:"operation_id"`
	ParentTaskID              *int      `json:"parent_task_id,omitempty"`
	ResponseCount             int       `json:"response_count"`
	IsInteractiveTask         bool      `json:"is_interactive_task"`
	InteractiveTaskType       *int      `json:"interactive_task_type,omitempty"`
	TaskingLocation           string    `json:"tasking_location"`
	ParameterGroupName        string    `json:"parameter_group_name"`
	Stdout                    string    `json:"stdout"`
	Stderr                    string    `json:"stderr"`
	CompletedCallbackFunction string    `json:"completed_callback_function"`
	SubtaskCallbackFunction   string    `json:"subtask_callback_function"`
	GroupCallbackFunction     string    `json:"group_callback_function"`
	OpsecPreBlocked           *bool     `json:"opsec_pre_blocked,omitempty"`
	OpsecPreBypassed          bool      `json:"opsec_pre_bypassed"`
	OpsecPreMessage           string    `json:"opsec_pre_message"`
	OpsecPostBlocked          *bool     `json:"opsec_post_blocked,omitempty"`
	OpsecPostBypassed         bool      `json:"opsec_post_bypassed"`
	OpsecPostMessage          string    `json:"opsec_post_message"`
}

// TaskResponse represents output from a task.
type TaskResponse struct {
	ID             int       `json:"id"`
	TaskID         int       `json:"task_id"`
	ResponseText   string    `json:"response_text"`
	ResponseRaw    []byte    `json:"response_raw"`
	IsError        bool      `json:"is_error"`
	Timestamp      time.Time `json:"timestamp"`
	SequenceNumber *int      `json:"sequence_number,omitempty"`
}

// TaskRequest represents a request to create a new task.
type TaskRequest struct {
	CallbackID          *int     `json:"callback_id,omitempty"`
	CallbackIDs         []int    `json:"callback_ids,omitempty"`
	Command             string   `json:"command"`
	Params              string   `json:"params"`
	Files               []string `json:"files,omitempty"`
	IsInteractiveTask   bool     `json:"is_interactive_task"`
	InteractiveTaskType *int     `json:"interactive_task_type,omitempty"`
	ParentTaskID        *int     `json:"parent_task_id,omitempty"`
	TaskingLocation     string   `json:"tasking_location,omitempty"`
	ParameterGroupName  string   `json:"parameter_group_name,omitempty"`
	OriginalParams      string   `json:"original_params,omitempty"`
	TokenID             *int     `json:"token_id,omitempty"`
}

// TaskStatus represents the status of a task.
type TaskStatus string

const (
	TaskStatusPreprocessing TaskStatus = "preprocessing"
	TaskStatusSubmitted     TaskStatus = "submitted"
	TaskStatusProcessing    TaskStatus = "processing"
	TaskStatusProcessed     TaskStatus = "processed"
	TaskStatusError         TaskStatus = "error"
	TaskStatusCleared       TaskStatus = "cleared"
	TaskStatusCompleted     TaskStatus = "completed"
)

// IssueTask creates a new task for a callback.
// If CallbackIDs is provided, the task will be issued to multiple callbacks.
//
// Note: This function uses the Hasura webhook endpoint directly instead of the GraphQL
// mutation. This is necessary because the GraphQL client library requires all variables
// to be present in the variables map, and serializes nil values as "null" in JSON.
// When Hasura receives explicit null values for optional array parameters, it validates
// them and rejects with "null value found for non-nullable type" errors, even though
// the GraphQL schema defines these as nullable ([Int] not [Int!]!).
//
// The webhook approach allows us to omit parameters entirely from the JSON request,
// which is the correct way to indicate "not provided" vs explicitly passing null.
// This is a documented, stable API endpoint that Hasura would call anyway.
func (c *Client) IssueTask(ctx context.Context, req *TaskRequest) (*Task, error) {
	if err := c.EnsureAuthenticated(ctx); err != nil {
		return nil, err
	}

	// Validate request
	if req.CallbackID == nil && len(req.CallbackIDs) == 0 {
		return nil, WrapError("IssueTask", ErrInvalidInput, "either callback_id or callback_ids must be provided")
	}
	if req.Command == "" {
		return nil, WrapError("IssueTask", ErrInvalidInput, "command is required")
	}

	// Build request payload - only include non-nil/non-empty values
	// This allows the webhook to properly distinguish "not provided" from "empty"
	payload := map[string]interface{}{
		"input": map[string]interface{}{
			"command":             req.Command,
			"params":              req.Params,
			"is_interactive_task": req.IsInteractiveTask,
		},
	}

	input, ok := payload["input"].(map[string]interface{})
	if !ok {
		return nil, WrapError("IssueTask", fmt.Errorf("unexpected payload structure"), "failed to construct request")
	}

	// Only include callback_id OR callback_ids, not both
	if req.CallbackID != nil {
		input["callback_id"] = req.CallbackID
	}
	if len(req.CallbackIDs) > 0 {
		input["callback_ids"] = req.CallbackIDs
	}

	// Only include optional fields if they're set
	if len(req.Files) > 0 {
		input["files"] = req.Files
	}
	if req.InteractiveTaskType != nil {
		input["interactive_task_type"] = req.InteractiveTaskType
	}
	if req.ParentTaskID != nil {
		input["parent_task_id"] = req.ParentTaskID
	}
	if req.TaskingLocation != "" {
		input["tasking_location"] = req.TaskingLocation
	}
	if req.ParameterGroupName != "" {
		input["parameter_group_name"] = req.ParameterGroupName
	}
	if req.OriginalParams != "" {
		input["original_params"] = req.OriginalParams
	}
	if req.TokenID != nil {
		input["token_id"] = req.TokenID
	}

	// Call the Hasura webhook endpoint
	var response struct {
		Status    string `json:"status"`
		Error     string `json:"error"`
		ID        int    `json:"id"`
		DisplayID int    `json:"display_id"`
	}

	err := c.executeRESTWebhook(ctx, "api/v1.4/create_task_webhook", payload, &response)
	if err != nil {
		return nil, WrapError("IssueTask", err, "failed to create task")
	}

	// Check for error in response
	if response.Status != "success" {
		return nil, WrapError("IssueTask", ErrOperationFailed, fmt.Sprintf("task creation failed: %s", response.Error))
	}

	// Get the full task details
	return c.GetTask(ctx, response.DisplayID)
}

// GetTask retrieves a task by its display ID.
func (c *Client) GetTask(ctx context.Context, displayID int) (*Task, error) {
	if err := c.EnsureAuthenticated(ctx); err != nil {
		return nil, err
	}

	var query struct {
		Task []struct {
			ID                        int    `graphql:"id"`
			DisplayID                 int    `graphql:"display_id"`
			AgentTaskID               string `graphql:"agent_task_id"`
			CommandName               string `graphql:"command_name"`
			Params                    string `graphql:"params"`
			DisplayParams             string `graphql:"display_params"`
			OriginalParams            string `graphql:"original_params"`
			Status                    string `graphql:"status"`
			Completed                 bool   `graphql:"completed"`
			Comment                   string `graphql:"comment"`
			Timestamp                 string `graphql:"timestamp"` // Use string to handle Mythic's timestamp format
			CallbackID                int    `graphql:"callback_id"`
			OperatorID                int    `graphql:"operator_id"`
			OperationID               int    `graphql:"operation_id"`
			ParentTaskID              *int   `graphql:"parent_task_id"`
			ResponseCount             int    `graphql:"response_count"`
			IsInteractiveTask         bool   `graphql:"is_interactive_task"`
			InteractiveTaskType       *int   `graphql:"interactive_task_type"`
			TaskingLocation           string `graphql:"tasking_location"`
			ParameterGroupName        string `graphql:"parameter_group_name"`
			Stdout                    string `graphql:"stdout"`
			Stderr                    string `graphql:"stderr"`
			CompletedCallbackFunction string `graphql:"completed_callback_function"`
			SubtaskCallbackFunction   string `graphql:"subtask_callback_function"`
			GroupCallbackFunction     string `graphql:"group_callback_function"`
			OpsecPreBlocked           *bool  `graphql:"opsec_pre_blocked"`
			OpsecPreBypassed          bool   `graphql:"opsec_pre_bypassed"`
			OpsecPreMessage           string `graphql:"opsec_pre_message"`
			OpsecPostBlocked          *bool  `graphql:"opsec_post_blocked"`
			OpsecPostBypassed         bool   `graphql:"opsec_post_bypassed"`
			OpsecPostMessage          string `graphql:"opsec_post_message"`
		} `graphql:"task(where: {display_id: {_eq: $display_id}}, limit: 1)"`
	}

	variables := map[string]interface{}{
		"display_id": displayID,
	}

	err := c.executeQuery(ctx, &query, variables)
	if err != nil {
		return nil, WrapError("GetTask", err, "failed to query task")
	}

	if len(query.Task) == 0 {
		return nil, WrapError("GetTask", ErrNotFound, fmt.Sprintf("task with display_id %d not found", displayID))
	}

	t := query.Task[0]

	// Parse timestamp string (Mythic returns timestamps without timezone)
	timestamp, err := parseTimestamp(t.Timestamp)
	if err != nil {
		// If parsing fails, use zero time but don't fail the entire request
		timestamp = time.Time{}
	}

	return &Task{
		ID:                        t.ID,
		DisplayID:                 t.DisplayID,
		AgentTaskID:               t.AgentTaskID,
		CommandName:               t.CommandName,
		Params:                    t.Params,
		DisplayParams:             t.DisplayParams,
		OriginalParams:            t.OriginalParams,
		Status:                    t.Status,
		Completed:                 t.Completed,
		Comment:                   t.Comment,
		Timestamp:                 timestamp,
		CallbackID:                t.CallbackID,
		OperatorID:                t.OperatorID,
		OperationID:               t.OperationID,
		ParentTaskID:              t.ParentTaskID,
		ResponseCount:             t.ResponseCount,
		IsInteractiveTask:         t.IsInteractiveTask,
		InteractiveTaskType:       t.InteractiveTaskType,
		TaskingLocation:           t.TaskingLocation,
		ParameterGroupName:        t.ParameterGroupName,
		Stdout:                    t.Stdout,
		Stderr:                    t.Stderr,
		CompletedCallbackFunction: t.CompletedCallbackFunction,
		SubtaskCallbackFunction:   t.SubtaskCallbackFunction,
		GroupCallbackFunction:     t.GroupCallbackFunction,
		OpsecPreBlocked:           t.OpsecPreBlocked,
		OpsecPreBypassed:          t.OpsecPreBypassed,
		OpsecPreMessage:           t.OpsecPreMessage,
		OpsecPostBlocked:          t.OpsecPostBlocked,
		OpsecPostBypassed:         t.OpsecPostBypassed,
		OpsecPostMessage:          t.OpsecPostMessage,
	}, nil
}

// GetTasksForCallback retrieves all tasks for a specific callback.
func (c *Client) GetTasksForCallback(ctx context.Context, callbackDisplayID int, limit int) ([]*Task, error) {
	if err := c.EnsureAuthenticated(ctx); err != nil {
		return nil, err
	}

	if limit <= 0 {
		limit = 100 // Default limit
	}

	// First get the callback's actual ID
	callback, err := c.GetCallbackByID(ctx, callbackDisplayID)
	if err != nil {
		return nil, WrapError("GetTasksForCallback", err, "failed to get callback")
	}

	var query struct {
		Task []struct {
			ID                int       `graphql:"id"`
			DisplayID         int       `graphql:"display_id"`
			AgentTaskID       string    `graphql:"agent_task_id"`
			CommandName       string    `graphql:"command_name"`
			Params            string    `graphql:"params"`
			DisplayParams     string    `graphql:"display_params"`
			Status            string    `graphql:"status"`
			Completed         bool      `graphql:"completed"`
			Comment           string    `graphql:"comment"`
			Timestamp         time.Time `graphql:"timestamp"`
			CallbackID        int       `graphql:"callback_id"`
			ResponseCount     int       `graphql:"response_count"`
			IsInteractiveTask bool      `graphql:"is_interactive_task"`
		} `graphql:"task(where: {callback_id: {_eq: $callback_id}}, order_by: {id: desc}, limit: $limit)"`
	}

	variables := map[string]interface{}{
		"callback_id": callback.ID,
		"limit":       limit,
	}

	err = c.executeQuery(ctx, &query, variables)
	if err != nil {
		return nil, WrapError("GetTasksForCallback", err, "failed to query tasks")
	}

	tasks := make([]*Task, 0, len(query.Task))
	for _, t := range query.Task {
		tasks = append(tasks, &Task{
			ID:                t.ID,
			DisplayID:         t.DisplayID,
			AgentTaskID:       t.AgentTaskID,
			CommandName:       t.CommandName,
			Params:            t.Params,
			DisplayParams:     t.DisplayParams,
			Status:            t.Status,
			Completed:         t.Completed,
			Comment:           t.Comment,
			Timestamp:         t.Timestamp,
			CallbackID:        t.CallbackID,
			ResponseCount:     t.ResponseCount,
			IsInteractiveTask: t.IsInteractiveTask,
		})
	}

	return tasks, nil
}

// GetTaskOutput retrieves all responses (output) for a task.
func (c *Client) GetTaskOutput(ctx context.Context, taskDisplayID int) ([]*TaskResponse, error) {
	if err := c.EnsureAuthenticated(ctx); err != nil {
		return nil, err
	}

	// First get the task's actual ID
	task, err := c.GetTask(ctx, taskDisplayID)
	if err != nil {
		return nil, WrapError("GetTaskOutput", err, "failed to get task")
	}

	var query struct {
		Response []struct {
			ID             int       `graphql:"id"`
			TaskID         int       `graphql:"task_id"`
			ResponseText   string    `graphql:"response_text"`
			IsError        bool      `graphql:"is_error"`
			Timestamp      time.Time `graphql:"timestamp"`
			SequenceNumber *int      `graphql:"sequence_number"`
		} `graphql:"response(where: {task_id: {_eq: $task_id}}, order_by: {id: asc})"`
	}

	variables := map[string]interface{}{
		"task_id": task.ID,
	}

	err = c.executeQuery(ctx, &query, variables)
	if err != nil {
		return nil, WrapError("GetTaskOutput", err, "failed to query responses")
	}

	responses := make([]*TaskResponse, 0, len(query.Response))
	for _, r := range query.Response {
		responses = append(responses, &TaskResponse{
			ID:             r.ID,
			TaskID:         r.TaskID,
			ResponseText:   r.ResponseText,
			IsError:        r.IsError,
			Timestamp:      r.Timestamp,
			SequenceNumber: r.SequenceNumber,
		})
	}

	return responses, nil
}

// WaitForTaskComplete polls a task until it completes or times out.
// Returns an error if the task fails or times out.
func (c *Client) WaitForTaskComplete(ctx context.Context, taskDisplayID int, timeoutSeconds int) error {
	if err := c.EnsureAuthenticated(ctx); err != nil {
		return err
	}

	if timeoutSeconds <= 0 {
		timeoutSeconds = 300 // Default 5 minutes
	}

	timeout := time.After(time.Duration(timeoutSeconds) * time.Second)
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-timeout:
			return WrapError("WaitForTaskComplete", ErrTimeout, fmt.Sprintf("task %d did not complete within %d seconds", taskDisplayID, timeoutSeconds))
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			task, err := c.GetTask(ctx, taskDisplayID)
			if err != nil {
				return WrapError("WaitForTaskComplete", err, "failed to check task status")
			}

			// Check if task completed
			if task.Completed {
				return nil
			}

			// Check if task errored
			if task.Status == string(TaskStatusError) {
				return WrapError("WaitForTaskComplete", ErrTaskFailed, fmt.Sprintf("task %d failed: %s", taskDisplayID, task.Stderr))
			}
		}
	}
}

// UpdateTask updates a task's properties.
func (c *Client) UpdateTask(ctx context.Context, displayID int, updates map[string]interface{}) error {
	if err := c.EnsureAuthenticated(ctx); err != nil {
		return err
	}

	// Get the task ID first
	task, err := c.GetTask(ctx, displayID)
	if err != nil {
		return WrapError("UpdateTask", err, "failed to get task")
	}

	var mutation struct {
		UpdateTask struct {
			Affected int `graphql:"affected_rows"`
		} `graphql:"update_task(where: {id: {_eq: $id}}, _set: $updates)"`
	}

	variables := map[string]interface{}{
		"id":      task.ID,
		"updates": updates,
	}

	err = c.executeMutation(ctx, &mutation, variables)
	if err != nil {
		return WrapError("UpdateTask", err, "failed to update task")
	}

	if mutation.UpdateTask.Affected == 0 {
		return WrapError("UpdateTask", ErrNotFound, fmt.Sprintf("task with display_id %d not found or no changes made", displayID))
	}

	return nil
}

// String returns a string representation of the task.
func (t *Task) String() string {
	return fmt.Sprintf("Task %d: %s %s (Status: %s, Completed: %t)",
		t.DisplayID, t.CommandName, t.DisplayParams, t.Status, t.Completed)
}

// ReissueTask reissues an existing task (creates a copy with same parameters).
func (c *Client) ReissueTask(ctx context.Context, taskID int) error {
	if err := c.EnsureAuthenticated(ctx); err != nil {
		return err
	}

	if taskID <= 0 {
		return WrapError("ReissueTask", ErrInvalidInput, "task_id must be positive")
	}

	var mutation struct {
		ReissueTask struct {
			Status string `graphql:"status"`
			Error  string `graphql:"error"`
		} `graphql:"reissue_task(task_id: $task_id)"`
	}

	variables := map[string]interface{}{
		"task_id": taskID,
	}

	err := c.executeMutation(ctx, &mutation, variables)
	if err != nil {
		return WrapError("ReissueTask", err, "failed to reissue task")
	}

	if mutation.ReissueTask.Status != "success" {
		return WrapError("ReissueTask", ErrInvalidResponse, fmt.Sprintf("reissue failed: %s", mutation.ReissueTask.Error))
	}

	return nil
}

// ReissueTaskWithHandler reissues a task using the handler (advanced reissue).
func (c *Client) ReissueTaskWithHandler(ctx context.Context, taskID int) error {
	if err := c.EnsureAuthenticated(ctx); err != nil {
		return err
	}

	if taskID <= 0 {
		return WrapError("ReissueTaskWithHandler", ErrInvalidInput, "task_id must be positive")
	}

	var mutation struct {
		ReissueTaskHandler struct {
			Status string `graphql:"status"`
			Error  string `graphql:"error"`
		} `graphql:"reissue_task_handler(task_id: $task_id)"`
	}

	variables := map[string]interface{}{
		"task_id": taskID,
	}

	err := c.executeMutation(ctx, &mutation, variables)
	if err != nil {
		return WrapError("ReissueTaskWithHandler", err, "failed to reissue task with handler")
	}

	if mutation.ReissueTaskHandler.Status != "success" {
		return WrapError("ReissueTaskWithHandler", ErrInvalidResponse, fmt.Sprintf("reissue failed: %s", mutation.ReissueTaskHandler.Error))
	}

	return nil
}

// RequestOpsecBypass requests OPSEC bypass for a blocked task.
func (c *Client) RequestOpsecBypass(ctx context.Context, taskID int) error {
	if err := c.EnsureAuthenticated(ctx); err != nil {
		return err
	}

	if taskID <= 0 {
		return WrapError("RequestOpsecBypass", ErrInvalidInput, "task_id must be positive")
	}

	var mutation struct {
		RequestOpsecBypass struct {
			Status string `graphql:"status"`
			Error  string `graphql:"error"`
		} `graphql:"requestOpsecBypass(task_id: $task_id)"`
	}

	variables := map[string]interface{}{
		"task_id": taskID,
	}

	err := c.executeMutation(ctx, &mutation, variables)
	if err != nil {
		return WrapError("RequestOpsecBypass", err, "failed to request OPSEC bypass")
	}

	if mutation.RequestOpsecBypass.Status != "success" {
		return WrapError("RequestOpsecBypass", ErrInvalidResponse, fmt.Sprintf("bypass request failed: %s", mutation.RequestOpsecBypass.Error))
	}

	return nil
}

// AddMITREAttackToTask tags a task with a MITRE ATT&CK technique.
func (c *Client) AddMITREAttackToTask(ctx context.Context, taskDisplayID int, attackID string) error {
	if err := c.EnsureAuthenticated(ctx); err != nil {
		return err
	}

	if taskDisplayID <= 0 {
		return WrapError("AddMITREAttackToTask", ErrInvalidInput, "task_display_id must be positive")
	}

	if attackID == "" {
		return WrapError("AddMITREAttackToTask", ErrInvalidInput, "attack ID (t_num) is required")
	}

	var mutation struct {
		AddAttackToTask struct {
			Status string `graphql:"status"`
			Error  string `graphql:"error"`
		} `graphql:"addAttackToTask(task_display_id: $task_display_id, t_num: $t_num)"`
	}

	variables := map[string]interface{}{
		"task_display_id": taskDisplayID,
		"t_num":           attackID,
	}

	err := c.executeMutation(ctx, &mutation, variables)
	if err != nil {
		return WrapError("AddMITREAttackToTask", err, "failed to add MITRE ATT&CK tag")
	}

	if mutation.AddAttackToTask.Status != "success" {
		return WrapError("AddMITREAttackToTask", ErrInvalidResponse, fmt.Sprintf("failed to add attack: %s", mutation.AddAttackToTask.Error))
	}

	return nil
}

// GetTasksByStatus retrieves tasks filtered by status.
func (c *Client) GetTasksByStatus(ctx context.Context, callbackDisplayID int, status TaskStatus, limit int) ([]*Task, error) {
	if err := c.EnsureAuthenticated(ctx); err != nil {
		return nil, err
	}

	if callbackDisplayID <= 0 {
		return nil, WrapError("GetTasksByStatus", ErrInvalidInput, "callback_display_id must be positive")
	}

	if limit <= 0 {
		limit = 100
	}

	var query struct {
		Task []struct {
			ID             int    `graphql:"id"`
			DisplayID      int    `graphql:"display_id"`
			CommandName    string `graphql:"command_name"`
			DisplayParams  string `graphql:"display_params"`
			OriginalParams string `graphql:"original_params"`
			Status         string `graphql:"status"`
			Completed      bool   `graphql:"completed"`
			Timestamp      string `graphql:"timestamp"`
			OperatorID     int    `graphql:"operator_id"`
		} `graphql:"task(where: {callback: {display_id: {_eq: $callback_display_id}}, status: {_eq: $status}}, order_by: {id: desc}, limit: $limit)"`
	}

	variables := map[string]interface{}{
		"callback_display_id": callbackDisplayID,
		"status":              string(status),
		"limit":               limit,
	}

	err := c.executeQuery(ctx, &query, variables)
	if err != nil {
		return nil, WrapError("GetTasksByStatus", err, "failed to query tasks by status")
	}

	tasks := make([]*Task, 0, len(query.Task))
	for _, t := range query.Task {
		// Parse timestamp
		var timestamp time.Time
		if t.Timestamp != "" {
			var mt Timestamp
			if err := mt.UnmarshalJSON([]byte(`"` + t.Timestamp + `"`)); err == nil {
				timestamp = mt.Time
			}
		}

		tasks = append(tasks, &Task{
			ID:             t.ID,
			DisplayID:      t.DisplayID,
			CommandName:    t.CommandName,
			DisplayParams:  t.DisplayParams,
			OriginalParams: t.OriginalParams,
			Status:         t.Status,
			Completed:      t.Completed,
			Timestamp:      timestamp,
			OperatorID:     t.OperatorID,
		})
	}

	return tasks, nil
}

// TaskArtifact represents an artifact (IOC) created by a task.
type TaskArtifact struct {
	ID           int       `json:"id"`
	TaskID       int       `json:"task_id"`
	Artifact     string    `json:"artifact"`
	BaseArtifact string    `json:"base_artifact"`
	Host         string    `json:"host"`
	Timestamp    time.Time `json:"timestamp"`
}

// GetTaskArtifacts retrieves artifacts/IOCs created by a task.
func (c *Client) GetTaskArtifacts(ctx context.Context, taskDisplayID int) ([]*TaskArtifact, error) {
	if err := c.EnsureAuthenticated(ctx); err != nil {
		return nil, err
	}

	if taskDisplayID <= 0 {
		return nil, WrapError("GetTaskArtifacts", ErrInvalidInput, "task_display_id must be positive")
	}

	var query struct {
		TaskArtifact []struct {
			ID           int    `graphql:"id"`
			TaskID       int    `graphql:"task_id"`
			Artifact     string `graphql:"artifact"`
			BaseArtifact string `graphql:"base_artifact"`
			Host         string `graphql:"host"`
			Timestamp    string `graphql:"timestamp"`
		} `graphql:"taskartifact(where: {task: {display_id: {_eq: $task_display_id}}}, order_by: {id: desc})"`
	}

	variables := map[string]interface{}{
		"task_display_id": taskDisplayID,
	}

	err := c.executeQuery(ctx, &query, variables)
	if err != nil {
		return nil, WrapError("GetTaskArtifacts", err, "failed to query task artifacts")
	}

	artifacts := make([]*TaskArtifact, 0, len(query.TaskArtifact))
	for _, a := range query.TaskArtifact {
		// Parse timestamp
		var timestamp time.Time
		if a.Timestamp != "" {
			var mt Timestamp
			if err := mt.UnmarshalJSON([]byte(`"` + a.Timestamp + `"`)); err == nil {
				timestamp = mt.Time
			}
		}

		artifacts = append(artifacts, &TaskArtifact{
			ID:           a.ID,
			TaskID:       a.TaskID,
			Artifact:     a.Artifact,
			BaseArtifact: a.BaseArtifact,
			Host:         a.Host,
			Timestamp:    timestamp,
		})
	}

	return artifacts, nil
}

// IsCompleted returns whether the task has completed.
func (t *Task) IsCompleted() bool {
	return t.Completed
}

// IsError returns whether the task has errored.
func (t *Task) IsError() bool {
	return t.Status == string(TaskStatusError)
}

// HasOutput returns whether the task has any responses.
func (t *Task) HasOutput() bool {
	return t.ResponseCount > 0
}
