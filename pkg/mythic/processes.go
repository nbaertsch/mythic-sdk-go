package mythic

import (
	"context"
	"encoding/json"
	"time"

	"github.com/nbaertsch/mythic-sdk-go/pkg/mythic/types"
)


// mythictreeProcess is the internal GraphQL shape for a process entry in mythictree.
type mythictreeProcess struct {
	ID          int             `graphql:"id"`
	Name        string          `graphql:"name"`
	FullPath    string          `graphql:"full_path"`
	Host        string          `graphql:"host"`
	Metadata    json.RawMessage `graphql:"metadata"`
	CallbackID  *int            `graphql:"callback_id"`
	TaskID      *int            `graphql:"task_id"`
	OperationID int             `graphql:"operation_id"`
	Timestamp   string          `graphql:"timestamp"`
	Deleted     bool            `graphql:"deleted"`
	Os          string          `graphql:"os"`
	ParentPath  string          `graphql:"parent_path"`
}

// processMetadata captures the dynamic fields Mythic stores in the metadata JSONB column.
type processMetadata struct {
	ProcessID       int    `json:"process_id"`
	ParentProcessID int    `json:"parent_process_id"`
	Architecture    string `json:"architecture"`
	User            string `json:"user"`
	CommandLine     string `json:"command_line"`
	BinPath         string `json:"bin_path"`
	IntegrityLevel  int    `json:"integrity_level"`
	StartTime       int64  `json:"start_time"` // unix epoch millis
	Description     string `json:"description"`
}

// toProcess converts a mythictree row into the SDK Process type.
func (m *mythictreeProcess) toProcess() *types.Process {
	var meta processMetadata
	if len(m.Metadata) > 0 {
		_ = json.Unmarshal(m.Metadata, &meta)
	}

	var startTime time.Time
	if meta.StartTime > 0 {
		startTime = time.UnixMilli(meta.StartTime)
	}

	binPath := meta.BinPath
	if binPath == "" {
		binPath = m.FullPath
	}
	name := m.Name
	if name == "" && meta.BinPath != "" {
		name = meta.BinPath
	}

	var ts time.Time
	if m.Timestamp != "" {
		parsed, err := parseTimestamp(m.Timestamp)
		if err == nil {
			ts = parsed
		}
	}

	return &types.Process{
		ID:              m.ID,
		Name:            name,
		ProcessID:       meta.ProcessID,
		ParentProcessID: meta.ParentProcessID,
		Architecture:    meta.Architecture,
		BinPath:         binPath,
		User:            meta.User,
		CommandLine:     meta.CommandLine,
		IntegrityLevel:  meta.IntegrityLevel,
		StartTime:       startTime,
		Description:     meta.Description,
		OperationID:     m.OperationID,
		CallbackID:      m.CallbackID,
		TaskID:          m.TaskID,
		Timestamp:       ts,
		Deleted:         m.Deleted,
		Host: &types.Host{
			Host: m.Host,
		},
	}
}

// GetProcesses retrieves all processes for the current operation.
// Processes in Mythic are stored in the mythictree table with tree_type = "process".
func (c *Client) GetProcesses(ctx context.Context) ([]*types.Process, error) {
	if err := c.EnsureAuthenticated(ctx); err != nil {
		return nil, err
	}

	var query struct {
		Mythictree []mythictreeProcess `graphql:"mythictree(where: {tree_type: {_eq: \"process\"}, deleted: {_eq: false}}, order_by: {timestamp: desc})"`
	}

	err := c.executeQuery(ctx, &query, nil)
	if err != nil {
		return nil, WrapError("GetProcesses", err, "failed to query processes")
	}

	processes := make([]*types.Process, len(query.Mythictree))
	for i, mt := range query.Mythictree {
		processes[i] = mt.toProcess()
	}

	return processes, nil
}

// GetProcessesByOperation retrieves processes for a specific operation.
func (c *Client) GetProcessesByOperation(ctx context.Context, operationID int) ([]*types.Process, error) {
	if err := c.EnsureAuthenticated(ctx); err != nil {
		return nil, err
	}

	if operationID == 0 {
		return nil, WrapError("GetProcessesByOperation", ErrInvalidInput, "operation ID is required")
	}

	var query struct {
		Mythictree []mythictreeProcess `graphql:"mythictree(where: {tree_type: {_eq: \"process\"}, operation_id: {_eq: $operation_id}, deleted: {_eq: false}}, order_by: {timestamp: desc})"`
	}

	variables := map[string]interface{}{
		"operation_id": operationID,
	}

	err := c.executeQuery(ctx, &query, variables)
	if err != nil {
		return nil, WrapError("GetProcessesByOperation", err, "failed to query processes")
	}

	processes := make([]*types.Process, len(query.Mythictree))
	for i, mt := range query.Mythictree {
		processes[i] = mt.toProcess()
	}

	return processes, nil
}

// GetProcessesByCallback retrieves processes for a specific callback.
func (c *Client) GetProcessesByCallback(ctx context.Context, callbackID int) ([]*types.Process, error) {
	if err := c.EnsureAuthenticated(ctx); err != nil {
		return nil, err
	}

	if callbackID == 0 {
		return nil, WrapError("GetProcessesByCallback", ErrInvalidInput, "callback ID is required")
	}

	var query struct {
		Mythictree []mythictreeProcess `graphql:"mythictree(where: {tree_type: {_eq: \"process\"}, callback_id: {_eq: $callback_id}, deleted: {_eq: false}}, order_by: {name: asc})"`
	}

	variables := map[string]interface{}{
		"callback_id": callbackID,
	}

	err := c.executeQuery(ctx, &query, variables)
	if err != nil {
		return nil, WrapError("GetProcessesByCallback", err, "failed to query processes")
	}

	processes := make([]*types.Process, len(query.Mythictree))
	for i, mt := range query.Mythictree {
		processes[i] = mt.toProcess()
	}

	return processes, nil
}

// GetProcessTree retrieves processes and organizes them into a tree structure.
// This builds a hierarchical view of processes based on parent-child relationships.
func (c *Client) GetProcessTree(ctx context.Context, callbackID int) ([]*types.ProcessTree, error) {
	if err := c.EnsureAuthenticated(ctx); err != nil {
		return nil, err
	}

	if callbackID == 0 {
		return nil, WrapError("GetProcessTree", ErrInvalidInput, "callback ID is required")
	}

	// Get all processes for the callback
	processes, err := c.GetProcessesByCallback(ctx, callbackID)
	if err != nil {
		return nil, WrapError("GetProcessTree", err, "failed to get processes")
	}

	// Build process tree
	return buildProcessTree(processes), nil
}

// buildProcessTree constructs a hierarchical tree from a flat list of processes.
func buildProcessTree(processes []*types.Process) []*types.ProcessTree {
	// Create maps for quick lookup
	processMap := make(map[int]*types.Process)
	treeMap := make(map[int]*types.ProcessTree)

	// Index all processes
	for _, proc := range processes {
		processMap[proc.ProcessID] = proc
		treeMap[proc.ProcessID] = &types.ProcessTree{
			Process:  proc,
			Children: []*types.ProcessTree{},
		}
	}

	// Build parent-child relationships
	var roots []*types.ProcessTree

	for _, proc := range processes {
		if proc.ParentProcessID == 0 || processMap[proc.ParentProcessID] == nil {
			// This is a root process (no parent or parent not in our list)
			roots = append(roots, treeMap[proc.ProcessID])
		} else {
			// Add to parent's children
			parentTree := treeMap[proc.ParentProcessID]
			if parentTree != nil {
				parentTree.Children = append(parentTree.Children, treeMap[proc.ProcessID])
			}
		}
	}

	return roots
}

// GetProcessesByHost retrieves processes for a specific host.
func (c *Client) GetProcessesByHost(ctx context.Context, hostID int) ([]*types.Process, error) {
	if err := c.EnsureAuthenticated(ctx); err != nil {
		return nil, err
	}

	if hostID == 0 {
		return nil, WrapError("GetProcessesByHost", ErrInvalidInput, "host ID is required")
	}

	// Look up host info first (hosts are derived from callbacks)
	host, err := c.GetHostByID(ctx, hostID)
	if err != nil {
		return nil, WrapError("GetProcessesByHost", err, "failed to resolve host")
	}

	// Query mythictree processes matching this hostname
	var query struct {
		Mythictree []mythictreeProcess `graphql:"mythictree(where: {tree_type: {_eq: \"process\"}, host: {_ilike: $hostname}, deleted: {_eq: false}}, order_by: {name: asc})"`
	}

	variables := map[string]interface{}{
		"hostname": host.Hostname,
	}

	err = c.executeQuery(ctx, &query, variables)
	if err != nil {
		return nil, WrapError("GetProcessesByHost", err, "failed to query processes")
	}

	processes := make([]*types.Process, len(query.Mythictree))
	for i, mt := range query.Mythictree {
		processes[i] = mt.toProcess()
	}

	return processes, nil
}
