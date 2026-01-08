package mythic

import (
	"context"

	"github.com/nbaertsch/mythic-sdk-go/pkg/mythic/types"
)

// GetFileBrowserObjects retrieves file browser objects for the current operation.
// These represent files and directories discovered through file browsing commands.
func (c *Client) GetFileBrowserObjects(ctx context.Context) ([]*types.FileBrowserObject, error) {
	if err := c.EnsureAuthenticated(ctx); err != nil {
		return nil, err
	}

	operationID := c.GetCurrentOperation()
	if operationID == nil {
		return nil, WrapError("GetFileBrowserObjects", ErrNotAuthenticated, "no current operation set")
	}

	var query struct {
		FileBrowserObjects []struct {
			ID            int    `graphql:"id"`
			Host          string `graphql:"host"`
			IsFile        bool   `graphql:"is_file"`
			Permissions   string `graphql:"permissions"`
			Name          string `graphql:"name"`
			ParentPath    string `graphql:"parent_path"`
			Success       bool   `graphql:"success"`
			AccessTime    string `graphql:"access_time"`
			ModifyTime    string `graphql:"modify_time"`
			Size          int64  `graphql:"size"`
			UpdateDeleted bool   `graphql:"update_deleted"`
			TaskID        int    `graphql:"task_id"`
			OperationID   int    `graphql:"operation_id"`
			Timestamp     string `graphql:"timestamp"`
			Comment       string `graphql:"comment"`
			Deleted       bool   `graphql:"deleted"`
			FullPathText  string `graphql:"full_path_text"`
			CallbackID    int    `graphql:"callback_id"`
			OperatorID    int    `graphql:"operator_id"`
		} `graphql:"filebrowserobj(where: {operation_id: {_eq: $operation_id}, deleted: {_eq: false}}, order_by: {full_path_text: asc})"`
	}

	variables := map[string]interface{}{
		"operation_id": *operationID,
	}

	err := c.executeQuery(ctx, &query, variables)
	if err != nil {
		return nil, WrapError("GetFileBrowserObjects", err, "failed to query file browser objects")
	}

	objects := make([]*types.FileBrowserObject, len(query.FileBrowserObjects))
	for i, obj := range query.FileBrowserObjects {
		accessTime, _ := parseTime(obj.AccessTime)
		modifyTime, _ := parseTime(obj.ModifyTime)
		timestamp, _ := parseTime(obj.Timestamp)

		objects[i] = &types.FileBrowserObject{
			ID:            obj.ID,
			Host:          obj.Host,
			IsFile:        obj.IsFile,
			Permissions:   obj.Permissions,
			Name:          obj.Name,
			ParentPath:    obj.ParentPath,
			Success:       obj.Success,
			AccessTime:    accessTime,
			ModifyTime:    modifyTime,
			Size:          obj.Size,
			UpdateDeleted: obj.UpdateDeleted,
			TaskID:        obj.TaskID,
			OperationID:   obj.OperationID,
			Timestamp:     timestamp,
			Comment:       obj.Comment,
			Deleted:       obj.Deleted,
			FullPathText:  obj.FullPathText,
			CallbackID:    obj.CallbackID,
			OperatorID:    obj.OperatorID,
		}
	}

	return objects, nil
}

// GetFileBrowserObjectsByHost retrieves file browser objects filtered by host.
func (c *Client) GetFileBrowserObjectsByHost(ctx context.Context, host string) ([]*types.FileBrowserObject, error) {
	if err := c.EnsureAuthenticated(ctx); err != nil {
		return nil, err
	}

	if host == "" {
		return nil, WrapError("GetFileBrowserObjectsByHost", ErrInvalidInput, "host is required")
	}

	operationID := c.GetCurrentOperation()
	if operationID == nil {
		return nil, WrapError("GetFileBrowserObjectsByHost", ErrNotAuthenticated, "no current operation set")
	}

	var query struct {
		FileBrowserObjects []struct {
			ID            int    `graphql:"id"`
			Host          string `graphql:"host"`
			IsFile        bool   `graphql:"is_file"`
			Permissions   string `graphql:"permissions"`
			Name          string `graphql:"name"`
			ParentPath    string `graphql:"parent_path"`
			Success       bool   `graphql:"success"`
			AccessTime    string `graphql:"access_time"`
			ModifyTime    string `graphql:"modify_time"`
			Size          int64  `graphql:"size"`
			UpdateDeleted bool   `graphql:"update_deleted"`
			TaskID        int    `graphql:"task_id"`
			OperationID   int    `graphql:"operation_id"`
			Timestamp     string `graphql:"timestamp"`
			Comment       string `graphql:"comment"`
			Deleted       bool   `graphql:"deleted"`
			FullPathText  string `graphql:"full_path_text"`
			CallbackID    int    `graphql:"callback_id"`
			OperatorID    int    `graphql:"operator_id"`
		} `graphql:"filebrowserobj(where: {operation_id: {_eq: $operation_id}, host: {_eq: $host}, deleted: {_eq: false}}, order_by: {full_path_text: asc})"`
	}

	variables := map[string]interface{}{
		"operation_id": *operationID,
		"host":         host,
	}

	err := c.executeQuery(ctx, &query, variables)
	if err != nil {
		return nil, WrapError("GetFileBrowserObjectsByHost", err, "failed to query file browser objects by host")
	}

	objects := make([]*types.FileBrowserObject, len(query.FileBrowserObjects))
	for i, obj := range query.FileBrowserObjects {
		accessTime, _ := parseTime(obj.AccessTime)
		modifyTime, _ := parseTime(obj.ModifyTime)
		timestamp, _ := parseTime(obj.Timestamp)

		objects[i] = &types.FileBrowserObject{
			ID:            obj.ID,
			Host:          obj.Host,
			IsFile:        obj.IsFile,
			Permissions:   obj.Permissions,
			Name:          obj.Name,
			ParentPath:    obj.ParentPath,
			Success:       obj.Success,
			AccessTime:    accessTime,
			ModifyTime:    modifyTime,
			Size:          obj.Size,
			UpdateDeleted: obj.UpdateDeleted,
			TaskID:        obj.TaskID,
			OperationID:   obj.OperationID,
			Timestamp:     timestamp,
			Comment:       obj.Comment,
			Deleted:       obj.Deleted,
			FullPathText:  obj.FullPathText,
			CallbackID:    obj.CallbackID,
			OperatorID:    obj.OperatorID,
		}
	}

	return objects, nil
}

// GetFileBrowserObjectsByCallback retrieves file browser objects filtered by callback ID.
func (c *Client) GetFileBrowserObjectsByCallback(ctx context.Context, callbackID int) ([]*types.FileBrowserObject, error) {
	if err := c.EnsureAuthenticated(ctx); err != nil {
		return nil, err
	}

	if callbackID == 0 {
		return nil, WrapError("GetFileBrowserObjectsByCallback", ErrInvalidInput, "callback ID is required")
	}

	operationID := c.GetCurrentOperation()
	if operationID == nil {
		return nil, WrapError("GetFileBrowserObjectsByCallback", ErrNotAuthenticated, "no current operation set")
	}

	var query struct {
		FileBrowserObjects []struct {
			ID            int    `graphql:"id"`
			Host          string `graphql:"host"`
			IsFile        bool   `graphql:"is_file"`
			Permissions   string `graphql:"permissions"`
			Name          string `graphql:"name"`
			ParentPath    string `graphql:"parent_path"`
			Success       bool   `graphql:"success"`
			AccessTime    string `graphql:"access_time"`
			ModifyTime    string `graphql:"modify_time"`
			Size          int64  `graphql:"size"`
			UpdateDeleted bool   `graphql:"update_deleted"`
			TaskID        int    `graphql:"task_id"`
			OperationID   int    `graphql:"operation_id"`
			Timestamp     string `graphql:"timestamp"`
			Comment       string `graphql:"comment"`
			Deleted       bool   `graphql:"deleted"`
			FullPathText  string `graphql:"full_path_text"`
			CallbackID    int    `graphql:"callback_id"`
			OperatorID    int    `graphql:"operator_id"`
		} `graphql:"filebrowserobj(where: {operation_id: {_eq: $operation_id}, callback_id: {_eq: $callback_id}, deleted: {_eq: false}}, order_by: {full_path_text: asc})"`
	}

	variables := map[string]interface{}{
		"operation_id": *operationID,
		"callback_id":  callbackID,
	}

	err := c.executeQuery(ctx, &query, variables)
	if err != nil {
		return nil, WrapError("GetFileBrowserObjectsByCallback", err, "failed to query file browser objects by callback")
	}

	objects := make([]*types.FileBrowserObject, len(query.FileBrowserObjects))
	for i, obj := range query.FileBrowserObjects {
		accessTime, _ := parseTime(obj.AccessTime)
		modifyTime, _ := parseTime(obj.ModifyTime)
		timestamp, _ := parseTime(obj.Timestamp)

		objects[i] = &types.FileBrowserObject{
			ID:            obj.ID,
			Host:          obj.Host,
			IsFile:        obj.IsFile,
			Permissions:   obj.Permissions,
			Name:          obj.Name,
			ParentPath:    obj.ParentPath,
			Success:       obj.Success,
			AccessTime:    accessTime,
			ModifyTime:    modifyTime,
			Size:          obj.Size,
			UpdateDeleted: obj.UpdateDeleted,
			TaskID:        obj.TaskID,
			OperationID:   obj.OperationID,
			Timestamp:     timestamp,
			Comment:       obj.Comment,
			Deleted:       obj.Deleted,
			FullPathText:  obj.FullPathText,
			CallbackID:    obj.CallbackID,
			OperatorID:    obj.OperatorID,
		}
	}

	return objects, nil
}
