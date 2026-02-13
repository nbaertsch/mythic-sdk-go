package mythic

import (
	"context"
	"time"

	"github.com/nbaertsch/mythic-sdk-go/pkg/mythic/types"
)

// GetKeylogs retrieves all keylog entries for the current operation.
func (c *Client) GetKeylogs(ctx context.Context) ([]*types.Keylog, error) {
	if err := c.EnsureAuthenticated(ctx); err != nil {
		return nil, err
	}

	var query struct {
		Keylog []struct {
			ID          int       `graphql:"id"`
			TaskID      int       `graphql:"task_id"`
			Keystrokes  string    `graphql:"keystrokes"`
			Window      string    `graphql:"window"`
			Timestamp   time.Time `graphql:"timestamp"`
			OperationID int       `graphql:"operation_id"`
			User        string    `graphql:"user"`
			Task        struct {
				CallbackID int `graphql:"callback_id"`
			} `graphql:"task"`
		} `graphql:"keylog(order_by: {timestamp: desc})"`
	}

	err := c.executeQuery(ctx, &query, nil)
	if err != nil {
		return nil, WrapError("GetKeylogs", err, "failed to query keylogs")
	}

	keylogs := make([]*types.Keylog, len(query.Keylog))
	for i, kl := range query.Keylog {
		keylogs[i] = &types.Keylog{
			ID:          kl.ID,
			TaskID:      kl.TaskID,
			Keystrokes:  kl.Keystrokes,
			Window:      kl.Window,
			Timestamp:   kl.Timestamp,
			OperationID: kl.OperationID,
			User:        kl.User,
			CallbackID:  kl.Task.CallbackID,
		}
	}

	return keylogs, nil
}

// GetKeylogsByOperation retrieves keylog entries for a specific operation.
func (c *Client) GetKeylogsByOperation(ctx context.Context, operationID int) ([]*types.Keylog, error) {
	if err := c.EnsureAuthenticated(ctx); err != nil {
		return nil, err
	}

	if operationID == 0 {
		return nil, WrapError("GetKeylogsByOperation", ErrInvalidInput, "operation ID is required")
	}

	var query struct {
		Keylog []struct {
			ID          int       `graphql:"id"`
			TaskID      int       `graphql:"task_id"`
			Keystrokes  string    `graphql:"keystrokes"`
			Window      string    `graphql:"window"`
			Timestamp   time.Time `graphql:"timestamp"`
			OperationID int       `graphql:"operation_id"`
			User        string    `graphql:"user"`
			Task        struct {
				CallbackID int `graphql:"callback_id"`
			} `graphql:"task"`
		} `graphql:"keylog(where: {operation_id: {_eq: $operation_id}}, order_by: {timestamp: desc})"`
	}

	variables := map[string]interface{}{
		"operation_id": operationID,
	}

	err := c.executeQuery(ctx, &query, variables)
	if err != nil {
		return nil, WrapError("GetKeylogsByOperation", err, "failed to query keylogs")
	}

	keylogs := make([]*types.Keylog, len(query.Keylog))
	for i, kl := range query.Keylog {
		keylogs[i] = &types.Keylog{
			ID:          kl.ID,
			TaskID:      kl.TaskID,
			Keystrokes:  kl.Keystrokes,
			Window:      kl.Window,
			Timestamp:   kl.Timestamp,
			OperationID: kl.OperationID,
			User:        kl.User,
			CallbackID:  kl.Task.CallbackID,
		}
	}

	return keylogs, nil
}

// GetKeylogsByCallback retrieves keylog entries for a specific callback.
func (c *Client) GetKeylogsByCallback(ctx context.Context, callbackID int) ([]*types.Keylog, error) {
	if err := c.EnsureAuthenticated(ctx); err != nil {
		return nil, err
	}

	if callbackID == 0 {
		return nil, WrapError("GetKeylogsByCallback", ErrInvalidInput, "callback ID is required")
	}

	var query struct {
		Keylog []struct {
			ID          int       `graphql:"id"`
			TaskID      int       `graphql:"task_id"`
			Keystrokes  string    `graphql:"keystrokes"`
			Window      string    `graphql:"window"`
			Timestamp   time.Time `graphql:"timestamp"`
			OperationID int       `graphql:"operation_id"`
			User        string    `graphql:"user"`
			Task        struct {
				CallbackID int `graphql:"callback_id"`
			} `graphql:"task"`
		} `graphql:"keylog(where: {task: {callback_id: {_eq: $callback_id}}}, order_by: {timestamp: desc})"`
	}

	variables := map[string]interface{}{
		"callback_id": callbackID,
	}

	err := c.executeQuery(ctx, &query, variables)
	if err != nil {
		return nil, WrapError("GetKeylogsByCallback", err, "failed to query keylogs")
	}

	keylogs := make([]*types.Keylog, len(query.Keylog))
	for i, kl := range query.Keylog {
		keylogs[i] = &types.Keylog{
			ID:          kl.ID,
			TaskID:      kl.TaskID,
			Keystrokes:  kl.Keystrokes,
			Window:      kl.Window,
			Timestamp:   kl.Timestamp,
			OperationID: kl.OperationID,
			User:        kl.User,
			CallbackID:  kl.Task.CallbackID,
		}
	}

	return keylogs, nil
}
