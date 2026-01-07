package mythic

import (
	"context"
	"fmt"
	"time"
)

// Task represents a Mythic task.
type Task struct {
	ID                   int       `json:"id"`
	DisplayID            int       `json:"display_id"`
	AgentTaskID          string    `json:"agent_task_id"`
	CommandName          string    `json:"command_name"`
	Params               string    `json:"params"`
	DisplayParams        string    `json:"display_params"`
	OriginalParams       string    `json:"original_params"`
	Status               string    `json:"status"`
	Completed            bool      `json:"completed"`
	Comment              string    `json:"comment"`
	Timestamp            time.Time `json:"timestamp"`
	CallbackID           int       `json:"callback_id"`
	OperatorID           int       `json:"operator_id"`
	OperationID          int       `json:"operation_id"`
	ParentTaskID         *int      `json:"parent_task_id,omitempty"`
	ResponseCount        int       `json:"response_count"`
	IsInteractiveTask    bool      `json:"is_interactive_task"`
	InteractiveTaskType  *int      `json:"interactive_task_type,omitempty"`
	TaskingLocation      string    `json:"tasking_location"`
	ParameterGroupName   string    `json:"parameter_group_name"`
	Stdout               string    `json:"stdout"`
	Stderr               string    `json:"stderr"`
	CompletedCallbackFunction string `json:"completed_callback_function"`
	SubtaskCallbackFunction   string `json:"subtask_callback_function"`
	GroupCallbackFunction     string `json:"group_callback_function"`
	OpsecPreBlocked      *bool  `json:"opsec_pre_blocked,omitempty"`
	OpsecPreBypassed     bool   `json:"opsec_pre_bypassed"`
	OpsecPreMessage      string `json:"opsec_pre_message"`
	OpsecPostBlocked     *bool  `json:"opsec_post_blocked,omitempty"`
	OpsecPostBypassed    bool   `json:"opsec_post_bypassed"`
	OpsecPostMessage     string `json:"opsec_post_message"`
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

	var mutation struct {
		CreateTask struct {
			ID            int    `graphql:"id"`
			DisplayID     int    `graphql:"display_id"`
			AgentTaskID   string `graphql:"agent_task_id"`
			Status        string `graphql:"status"`
			Error         string `graphql:"error"`
			Stdout        string `graphql:"stdout"`
			Stderr        string `graphql:"stderr"`
		} `graphql:"createTask(callback_id: $callback_id, callback_ids: $callback_ids, command: $command, params: $params, files: $files, is_interactive_task: $is_interactive_task, interactive_task_type: $interactive_task_type, parent_task_id: $parent_task_id, tasking_location: $tasking_location, parameter_group_name: $parameter_group_name, original_params: $original_params, token_id: $token_id)"`
	}

	variables := map[string]interface{}{
		"callback_id":           req.CallbackID,
		"callback_ids":          req.CallbackIDs,
		"command":               req.Command,
		"params":                req.Params,
		"files":                 req.Files,
		"is_interactive_task":   req.IsInteractiveTask,
		"interactive_task_type": req.InteractiveTaskType,
		"parent_task_id":        req.ParentTaskID,
		"tasking_location":      req.TaskingLocation,
		"parameter_group_name":  req.ParameterGroupName,
		"original_params":       req.OriginalParams,
		"token_id":              req.TokenID,
	}

	err := c.executeMutation(ctx, &mutation, variables)
	if err != nil {
		return nil, WrapError("IssueTask", err, "failed to create task")
	}

	// Check for error in response
	if mutation.CreateTask.Error != "" {
		return nil, WrapError("IssueTask", ErrInvalidResponse, fmt.Sprintf("task creation failed: %s", mutation.CreateTask.Error))
	}

	// Get the full task details
	return c.GetTask(ctx, mutation.CreateTask.DisplayID)
}

// GetTask retrieves a task by its display ID.
func (c *Client) GetTask(ctx context.Context, displayID int) (*Task, error) {
	if err := c.EnsureAuthenticated(ctx); err != nil {
		return nil, err
	}

	var query struct {
		Task []struct {
			ID                        int       `graphql:"id"`
			DisplayID                 int       `graphql:"display_id"`
			AgentTaskID               string    `graphql:"agent_task_id"`
			CommandName               string    `graphql:"command_name"`
			Params                    string    `graphql:"params"`
			DisplayParams             string    `graphql:"display_params"`
			OriginalParams            string    `graphql:"original_params"`
			Status                    string    `graphql:"status"`
			Completed                 bool      `graphql:"completed"`
			Comment                   string    `graphql:"comment"`
			Timestamp                 time.Time `graphql:"timestamp"`
			CallbackID                int       `graphql:"callback_id"`
			OperatorID                int       `graphql:"operator_id"`
			OperationID               int       `graphql:"operation_id"`
			ParentTaskID              *int      `graphql:"parent_task_id"`
			ResponseCount             int       `graphql:"response_count"`
			IsInteractiveTask         bool      `graphql:"is_interactive_task"`
			InteractiveTaskType       *int      `graphql:"interactive_task_type"`
			TaskingLocation           string    `graphql:"tasking_location"`
			ParameterGroupName        string    `graphql:"parameter_group_name"`
			Stdout                    string    `graphql:"stdout"`
			Stderr                    string    `graphql:"stderr"`
			CompletedCallbackFunction string    `graphql:"completed_callback_function"`
			SubtaskCallbackFunction   string    `graphql:"subtask_callback_function"`
			GroupCallbackFunction     string    `graphql:"group_callback_function"`
			OpsecPreBlocked           *bool     `graphql:"opsec_pre_blocked"`
			OpsecPreBypassed          bool      `graphql:"opsec_pre_bypassed"`
			OpsecPreMessage           string    `graphql:"opsec_pre_message"`
			OpsecPostBlocked          *bool     `graphql:"opsec_post_blocked"`
			OpsecPostBypassed         bool      `graphql:"opsec_post_bypassed"`
			OpsecPostMessage          string    `graphql:"opsec_post_message"`
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
		Timestamp:                 t.Timestamp,
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
			ID                        int       `graphql:"id"`
			DisplayID                 int       `graphql:"display_id"`
			AgentTaskID               string    `graphql:"agent_task_id"`
			CommandName               string    `graphql:"command_name"`
			Params                    string    `graphql:"params"`
			DisplayParams             string    `graphql:"display_params"`
			Status                    string    `graphql:"status"`
			Completed                 bool      `graphql:"completed"`
			Comment                   string    `graphql:"comment"`
			Timestamp                 time.Time `graphql:"timestamp"`
			CallbackID                int       `graphql:"callback_id"`
			ResponseCount             int       `graphql:"response_count"`
			IsInteractiveTask         bool      `graphql:"is_interactive_task"`
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
