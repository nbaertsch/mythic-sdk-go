package mythic

import (
	"context"
	"time"

	"github.com/nbaertsch/mythic-sdk-go/pkg/mythic/types"
)

// GetTokens retrieves all tokens (non-deleted) for the current operation.
func (c *Client) GetTokens(ctx context.Context) ([]*types.Token, error) {
	if err := c.EnsureAuthenticated(ctx); err != nil {
		return nil, err
	}

	// Use current operation if set
	operationID := c.GetCurrentOperation()
	if operationID == nil {
		return nil, WrapError("GetTokens", ErrInvalidInput, "no current operation set")
	}

	return c.GetTokensByOperation(ctx, *operationID)
}

// GetTokensByOperation retrieves all tokens for a specific operation.
func (c *Client) GetTokensByOperation(ctx context.Context, operationID int) ([]*types.Token, error) {
	if err := c.EnsureAuthenticated(ctx); err != nil {
		return nil, err
	}

	if operationID == 0 {
		return nil, WrapError("GetTokensByOperation", ErrInvalidInput, "operation ID is required")
	}

	var query struct {
		Token []struct {
			ID                 int       `graphql:"id"`
			TokenID            string    `graphql:"token_id"`
			User               string    `graphql:"user"`
			Groups             string    `graphql:"groups"`
			Privileges         string    `graphql:"privileges"`
			ThreadID           int       `graphql:"thread_id"`
			ProcessID          int       `graphql:"process_id"`
			SessionID          int       `graphql:"session_id"`
			LogonSID           string    `graphql:"logon_sid"`
			IntegrityLevelInt  int       `graphql:"integrity_level_int"`
			Restricted         bool      `graphql:"restricted"`
			DefaultDACL        string    `graphql:"default_dacl"`
			Handle             string    `graphql:"handle"`
			Capabilities       string    `graphql:"capabilities"`
			AppContainerSID    string    `graphql:"app_container_sid"`
			AppContainerNumber int       `graphql:"app_container_number"`
			TaskID             *int      `graphql:"task_id"`
			OperationID        int       `graphql:"operation_id"`
			Timestamp          time.Time `graphql:"timestamp"`
			Host               string    `graphql:"host"`
			Deleted            bool      `graphql:"deleted"`
		} `graphql:"token(where: {operation_id: {_eq: $operation_id}, deleted: {_eq: false}}, order_by: {timestamp: desc})"`
	}

	variables := map[string]interface{}{
		"operation_id": operationID,
	}

	err := c.executeQuery(ctx, &query, variables)
	if err != nil {
		return nil, WrapError("GetTokensByOperation", err, "failed to query tokens")
	}

	tokens := make([]*types.Token, len(query.Token))
	for i, t := range query.Token {
		tokens[i] = &types.Token{
			ID:              t.ID,
			TokenID:         t.TokenID,
			User:            t.User,
			Groups:          t.Groups,
			Privileges:      t.Privileges,
			ThreadID:        t.ThreadID,
			ProcessID:       t.ProcessID,
			SessionID:       t.SessionID,
			LogonSID:        t.LogonSID,
			IntegrityLevel:  t.IntegrityLevelInt,
			Restricted:      t.Restricted,
			DefaultDACL:     t.DefaultDACL,
			Handle:          t.Handle,
			Capabilities:    t.Capabilities,
			AppContainerSID: t.AppContainerSID,
			AppContainerNum: t.AppContainerNumber,
			TaskID:          t.TaskID,
			OperationID:     t.OperationID,
			Timestamp:       t.Timestamp,
			Host:            t.Host,
			Deleted:         t.Deleted,
		}
	}

	return tokens, nil
}

// GetTokenByID retrieves a specific token by ID.
func (c *Client) GetTokenByID(ctx context.Context, tokenID int) (*types.Token, error) {
	if err := c.EnsureAuthenticated(ctx); err != nil {
		return nil, err
	}

	if tokenID == 0 {
		return nil, WrapError("GetTokenByID", ErrInvalidInput, "token ID is required")
	}

	var query struct {
		Token []struct {
			ID                 int       `graphql:"id"`
			TokenID            string    `graphql:"token_id"`
			User               string    `graphql:"user"`
			Groups             string    `graphql:"groups"`
			Privileges         string    `graphql:"privileges"`
			ThreadID           int       `graphql:"thread_id"`
			ProcessID          int       `graphql:"process_id"`
			SessionID          int       `graphql:"session_id"`
			LogonSID           string    `graphql:"logon_sid"`
			IntegrityLevelInt  int       `graphql:"integrity_level_int"`
			Restricted         bool      `graphql:"restricted"`
			DefaultDACL        string    `graphql:"default_dacl"`
			Handle             string    `graphql:"handle"`
			Capabilities       string    `graphql:"capabilities"`
			AppContainerSID    string    `graphql:"app_container_sid"`
			AppContainerNumber int       `graphql:"app_container_number"`
			TaskID             *int      `graphql:"task_id"`
			OperationID        int       `graphql:"operation_id"`
			Timestamp          time.Time `graphql:"timestamp"`
			Host               string    `graphql:"host"`
			Deleted            bool      `graphql:"deleted"`
		} `graphql:"token(where: {id: {_eq: $token_id}})"`
	}

	variables := map[string]interface{}{
		"token_id": tokenID,
	}

	err := c.executeQuery(ctx, &query, variables)
	if err != nil {
		return nil, WrapError("GetTokenByID", err, "failed to query token")
	}

	if len(query.Token) == 0 {
		return nil, WrapError("GetTokenByID", ErrNotFound, "token not found")
	}

	t := query.Token[0]
	return &types.Token{
		ID:              t.ID,
		TokenID:         t.TokenID,
		User:            t.User,
		Groups:          t.Groups,
		Privileges:      t.Privileges,
		ThreadID:        t.ThreadID,
		ProcessID:       t.ProcessID,
		SessionID:       t.SessionID,
		LogonSID:        t.LogonSID,
		IntegrityLevel:  t.IntegrityLevelInt,
		Restricted:      t.Restricted,
		DefaultDACL:     t.DefaultDACL,
		Handle:          t.Handle,
		Capabilities:    t.Capabilities,
		AppContainerSID: t.AppContainerSID,
		AppContainerNum: t.AppContainerNumber,
		TaskID:          t.TaskID,
		OperationID:     t.OperationID,
		Timestamp:       t.Timestamp,
		Host:            t.Host,
		Deleted:         t.Deleted,
	}, nil
}

// GetCallbackTokens retrieves all callback tokens for the current operation.
func (c *Client) GetCallbackTokens(ctx context.Context) ([]*types.CallbackToken, error) {
	if err := c.EnsureAuthenticated(ctx); err != nil {
		return nil, err
	}

	// Use current operation if set
	operationID := c.GetCurrentOperation()
	if operationID == nil {
		return nil, WrapError("GetCallbackTokens", ErrInvalidInput, "no current operation set")
	}

	var query struct {
		CallbackToken []struct {
			ID         int       `graphql:"id"`
			CallbackID int       `graphql:"callback_id"`
			TokenID    int       `graphql:"token_id"`
			Timestamp  time.Time `graphql:"timestamp"`
		} `graphql:"callbacktoken(where: {callback: {operation_id: {_eq: $operation_id}}}, order_by: {id: desc})"`
	}

	variables := map[string]interface{}{
		"operation_id": operationID,
	}

	err := c.executeQuery(ctx, &query, variables)
	if err != nil {
		return nil, WrapError("GetCallbackTokens", err, "failed to query callback tokens")
	}

	callbackTokens := make([]*types.CallbackToken, len(query.CallbackToken))
	for i, ct := range query.CallbackToken {
		callbackTokens[i] = &types.CallbackToken{
			ID:         ct.ID,
			CallbackID: ct.CallbackID,
			TokenID:    ct.TokenID,
			Timestamp:  ct.Timestamp,
		}
	}

	return callbackTokens, nil
}

// GetCallbackTokensByCallback retrieves all tokens associated with a specific callback.
func (c *Client) GetCallbackTokensByCallback(ctx context.Context, callbackID int) ([]*types.CallbackToken, error) {
	if err := c.EnsureAuthenticated(ctx); err != nil {
		return nil, err
	}

	if callbackID == 0 {
		return nil, WrapError("GetCallbackTokensByCallback", ErrInvalidInput, "callback ID is required")
	}

	var query struct {
		CallbackToken []struct {
			ID         int       `graphql:"id"`
			CallbackID int       `graphql:"callback_id"`
			TokenID    int       `graphql:"token_id"`
			Timestamp  time.Time `graphql:"timestamp"`
		} `graphql:"callbacktoken(where: {callback_id: {_eq: $callback_id}}, order_by: {id: desc})"`
	}

	variables := map[string]interface{}{
		"callback_id": callbackID,
	}

	err := c.executeQuery(ctx, &query, variables)
	if err != nil {
		return nil, WrapError("GetCallbackTokensByCallback", err, "failed to query callback tokens")
	}

	callbackTokens := make([]*types.CallbackToken, len(query.CallbackToken))
	for i, ct := range query.CallbackToken {
		callbackTokens[i] = &types.CallbackToken{
			ID:         ct.ID,
			CallbackID: ct.CallbackID,
			TokenID:    ct.TokenID,
			Timestamp:  ct.Timestamp,
		}
	}

	return callbackTokens, nil
}

// GetAPITokens retrieves all API tokens for the authenticated user.
func (c *Client) GetAPITokens(ctx context.Context) ([]*types.APIToken, error) {
	if err := c.EnsureAuthenticated(ctx); err != nil {
		return nil, err
	}

	var query struct {
		APITokens []struct {
			ID           int       `graphql:"id"`
			TokenValue   string    `graphql:"token_value"`
			TokenType    string    `graphql:"token_type"`
			Active       bool      `graphql:"active"`
			CreationTime time.Time `graphql:"creation_time"`
			OperatorID   int       `graphql:"operator_id"`
			OperationID  *int      `graphql:"operation_id"`
			Name         string    `graphql:"name"`
			Deleted      bool      `graphql:"deleted"`
		} `graphql:"apitokens(where: {deleted: {_eq: false}}, order_by: {creation_time: desc})"`
	}

	err := c.executeQuery(ctx, &query, nil)
	if err != nil {
		return nil, WrapError("GetAPITokens", err, "failed to query API tokens")
	}

	apiTokens := make([]*types.APIToken, len(query.APITokens))
	for i, at := range query.APITokens {
		apiTokens[i] = &types.APIToken{
			ID:           at.ID,
			TokenValue:   at.TokenValue,
			TokenType:    at.TokenType,
			Active:       at.Active,
			CreationTime: at.CreationTime,
			OperatorID:   at.OperatorID,
			OperationID:  at.OperationID,
			Name:         at.Name,
			Deleted:      at.Deleted,
		}
	}

	return apiTokens, nil
}

// DeleteAPIToken marks an API token as deleted.
func (c *Client) DeleteAPIToken(ctx context.Context, tokenID int) error {
	if err := c.EnsureAuthenticated(ctx); err != nil {
		return err
	}

	if tokenID == 0 {
		return WrapError("DeleteAPIToken", ErrInvalidInput, "token ID is required")
	}

	var mutation struct {
		DeleteAPIToken struct {
			Status string `graphql:"status"`
			Error  string `graphql:"error"`
		} `graphql:"deleteAPIToken(apitokens_id: $token_id)"`
	}

	variables := map[string]interface{}{
		"token_id": tokenID,
	}

	err := c.executeMutation(ctx, &mutation, variables)
	if err != nil {
		return WrapError("DeleteAPIToken", err, "failed to delete API token")
	}

	if mutation.DeleteAPIToken.Status != "success" {
		return WrapError("DeleteAPIToken", ErrOperationFailed, mutation.DeleteAPIToken.Error)
	}

	return nil
}
