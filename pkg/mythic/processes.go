package mythic

import (
	"context"
	"time"

	"github.com/nbaertsch/mythic-sdk-go/pkg/mythic/types"
)

// GetProcesses retrieves all processes for the current operation.
func (c *Client) GetProcesses(ctx context.Context) ([]*types.Process, error) {
	if err := c.EnsureAuthenticated(ctx); err != nil {
		return nil, err
	}

	var query struct {
		Process []struct {
			ID              int       `graphql:"id"`
			Name            string    `graphql:"name"`
			ProcessID       int       `graphql:"process_id"`
			ParentProcessID int       `graphql:"parent_process_id"`
			Architecture    string    `graphql:"architecture"`
			BinPath         string    `graphql:"bin_path"`
			User            string    `graphql:"user"`
			CommandLine     string    `graphql:"command_line"`
			IntegrityLevel  int       `graphql:"integrity_level"`
			StartTime       time.Time `graphql:"start_time"`
			Description     string    `graphql:"description"`
			OperationID     int       `graphql:"operation_id"`
			HostID          int       `graphql:"host_id"`
			CallbackID      *int      `graphql:"callback_id"`
			TaskID          *int      `graphql:"task_id"`
			Timestamp       time.Time `graphql:"timestamp"`
			Deleted         bool      `graphql:"deleted"`
			Host            struct {
				ID   int    `graphql:"id"`
				Host string `graphql:"host"`
			} `graphql:"host"`
		} `graphql:"process(where: {deleted: {_eq: false}}, order_by: {timestamp: desc})"`
	}

	err := c.executeQuery(ctx, &query, nil)
	if err != nil {
		return nil, WrapError("GetProcesses", err, "failed to query processes")
	}

	processes := make([]*types.Process, len(query.Process))
	for i, proc := range query.Process {
		processes[i] = &types.Process{
			ID:              proc.ID,
			Name:            proc.Name,
			ProcessID:       proc.ProcessID,
			ParentProcessID: proc.ParentProcessID,
			Architecture:    proc.Architecture,
			BinPath:         proc.BinPath,
			User:            proc.User,
			CommandLine:     proc.CommandLine,
			IntegrityLevel:  proc.IntegrityLevel,
			StartTime:       proc.StartTime,
			Description:     proc.Description,
			OperationID:     proc.OperationID,
			HostID:          proc.HostID,
			CallbackID:      proc.CallbackID,
			TaskID:          proc.TaskID,
			Timestamp:       proc.Timestamp,
			Deleted:         proc.Deleted,
			Host: &types.Host{
				ID:   proc.Host.ID,
				Host: proc.Host.Host,
			},
		}
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
		Process []struct {
			ID              int       `graphql:"id"`
			Name            string    `graphql:"name"`
			ProcessID       int       `graphql:"process_id"`
			ParentProcessID int       `graphql:"parent_process_id"`
			Architecture    string    `graphql:"architecture"`
			BinPath         string    `graphql:"bin_path"`
			User            string    `graphql:"user"`
			CommandLine     string    `graphql:"command_line"`
			IntegrityLevel  int       `graphql:"integrity_level"`
			StartTime       time.Time `graphql:"start_time"`
			Description     string    `graphql:"description"`
			OperationID     int       `graphql:"operation_id"`
			HostID          int       `graphql:"host_id"`
			CallbackID      *int      `graphql:"callback_id"`
			TaskID          *int      `graphql:"task_id"`
			Timestamp       time.Time `graphql:"timestamp"`
			Deleted         bool      `graphql:"deleted"`
			Host            struct {
				ID   int    `graphql:"id"`
				Host string `graphql:"host"`
			} `graphql:"host"`
		} `graphql:"process(where: {operation_id: {_eq: $operation_id}, deleted: {_eq: false}}, order_by: {timestamp: desc})"`
	}

	variables := map[string]interface{}{
		"operation_id": operationID,
	}

	err := c.executeQuery(ctx, &query, variables)
	if err != nil {
		return nil, WrapError("GetProcessesByOperation", err, "failed to query processes")
	}

	processes := make([]*types.Process, len(query.Process))
	for i, proc := range query.Process {
		processes[i] = &types.Process{
			ID:              proc.ID,
			Name:            proc.Name,
			ProcessID:       proc.ProcessID,
			ParentProcessID: proc.ParentProcessID,
			Architecture:    proc.Architecture,
			BinPath:         proc.BinPath,
			User:            proc.User,
			CommandLine:     proc.CommandLine,
			IntegrityLevel:  proc.IntegrityLevel,
			StartTime:       proc.StartTime,
			Description:     proc.Description,
			OperationID:     proc.OperationID,
			HostID:          proc.HostID,
			CallbackID:      proc.CallbackID,
			TaskID:          proc.TaskID,
			Timestamp:       proc.Timestamp,
			Deleted:         proc.Deleted,
			Host: &types.Host{
				ID:   proc.Host.ID,
				Host: proc.Host.Host,
			},
		}
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
		Process []struct {
			ID              int       `graphql:"id"`
			Name            string    `graphql:"name"`
			ProcessID       int       `graphql:"process_id"`
			ParentProcessID int       `graphql:"parent_process_id"`
			Architecture    string    `graphql:"architecture"`
			BinPath         string    `graphql:"bin_path"`
			User            string    `graphql:"user"`
			CommandLine     string    `graphql:"command_line"`
			IntegrityLevel  int       `graphql:"integrity_level"`
			StartTime       time.Time `graphql:"start_time"`
			Description     string    `graphql:"description"`
			OperationID     int       `graphql:"operation_id"`
			HostID          int       `graphql:"host_id"`
			CallbackID      *int      `graphql:"callback_id"`
			TaskID          *int      `graphql:"task_id"`
			Timestamp       time.Time `graphql:"timestamp"`
			Deleted         bool      `graphql:"deleted"`
			Host            struct {
				ID   int    `graphql:"id"`
				Host string `graphql:"host"`
			} `graphql:"host"`
		} `graphql:"process(where: {callback_id: {_eq: $callback_id}, deleted: {_eq: false}}, order_by: {process_id: asc})"`
	}

	variables := map[string]interface{}{
		"callback_id": callbackID,
	}

	err := c.executeQuery(ctx, &query, variables)
	if err != nil {
		return nil, WrapError("GetProcessesByCallback", err, "failed to query processes")
	}

	processes := make([]*types.Process, len(query.Process))
	for i, proc := range query.Process {
		processes[i] = &types.Process{
			ID:              proc.ID,
			Name:            proc.Name,
			ProcessID:       proc.ProcessID,
			ParentProcessID: proc.ParentProcessID,
			Architecture:    proc.Architecture,
			BinPath:         proc.BinPath,
			User:            proc.User,
			CommandLine:     proc.CommandLine,
			IntegrityLevel:  proc.IntegrityLevel,
			StartTime:       proc.StartTime,
			Description:     proc.Description,
			OperationID:     proc.OperationID,
			HostID:          proc.HostID,
			CallbackID:      proc.CallbackID,
			TaskID:          proc.TaskID,
			Timestamp:       proc.Timestamp,
			Deleted:         proc.Deleted,
			Host: &types.Host{
				ID:   proc.Host.ID,
				Host: proc.Host.Host,
			},
		}
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

	var query struct {
		Process []struct {
			ID              int       `graphql:"id"`
			Name            string    `graphql:"name"`
			ProcessID       int       `graphql:"process_id"`
			ParentProcessID int       `graphql:"parent_process_id"`
			Architecture    string    `graphql:"architecture"`
			BinPath         string    `graphql:"bin_path"`
			User            string    `graphql:"user"`
			CommandLine     string    `graphql:"command_line"`
			IntegrityLevel  int       `graphql:"integrity_level"`
			StartTime       time.Time `graphql:"start_time"`
			Description     string    `graphql:"description"`
			OperationID     int       `graphql:"operation_id"`
			HostID          int       `graphql:"host_id"`
			CallbackID      *int      `graphql:"callback_id"`
			TaskID          *int      `graphql:"task_id"`
			Timestamp       time.Time `graphql:"timestamp"`
			Deleted         bool      `graphql:"deleted"`
			Host            struct {
				ID   int    `graphql:"id"`
				Host string `graphql:"host"`
			} `graphql:"host"`
		} `graphql:"process(where: {host_id: {_eq: $host_id}, deleted: {_eq: false}}, order_by: {process_id: asc})"`
	}

	variables := map[string]interface{}{
		"host_id": hostID,
	}

	err := c.executeQuery(ctx, &query, variables)
	if err != nil {
		return nil, WrapError("GetProcessesByHost", err, "failed to query processes")
	}

	processes := make([]*types.Process, len(query.Process))
	for i, proc := range query.Process {
		processes[i] = &types.Process{
			ID:              proc.ID,
			Name:            proc.Name,
			ProcessID:       proc.ProcessID,
			ParentProcessID: proc.ParentProcessID,
			Architecture:    proc.Architecture,
			BinPath:         proc.BinPath,
			User:            proc.User,
			CommandLine:     proc.CommandLine,
			IntegrityLevel:  proc.IntegrityLevel,
			StartTime:       proc.StartTime,
			Description:     proc.Description,
			OperationID:     proc.OperationID,
			HostID:          proc.HostID,
			CallbackID:      proc.CallbackID,
			TaskID:          proc.TaskID,
			Timestamp:       proc.Timestamp,
			Deleted:         proc.Deleted,
			Host: &types.Host{
				ID:   proc.Host.ID,
				Host: proc.Host.Host,
			},
		}
	}

	return processes, nil
}
