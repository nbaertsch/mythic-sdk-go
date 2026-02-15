package mythic

import (
	"context"
	"fmt"
	"strings"

	"github.com/nbaertsch/mythic-sdk-go/pkg/mythic/types"
)

// GetOperations retrieves all operations.
func (c *Client) GetOperations(ctx context.Context) ([]*types.Operation, error) {
	if err := c.EnsureAuthenticated(ctx); err != nil {
		return nil, err
	}

	var query struct {
		Operation []struct {
			ID          int    `graphql:"id"`
			Name        string `graphql:"name"`
			Complete    bool   `graphql:"complete"`
			Webhook     string `graphql:"webhook"`
			Channel     string `graphql:"channel"`
			AdminID     int    `graphql:"admin_id"`
			BannerText  string `graphql:"banner_text"`
			BannerColor string `graphql:"banner_color"`
			Admin       struct {
				ID       int    `graphql:"id"`
				Username string `graphql:"username"`
				Admin    bool   `graphql:"admin"`
			} `graphql:"admin"`
		} `graphql:"operation(order_by: {id: desc})"`
	}

	err := c.executeQuery(ctx, &query, nil)
	if err != nil {
		return nil, WrapError("GetOperations", err, "failed to query operations")
	}

	operations := make([]*types.Operation, len(query.Operation))
	for i, op := range query.Operation {
		operations[i] = &types.Operation{
			ID:          op.ID,
			Name:        op.Name,
			Complete:    op.Complete,
			Webhook:     op.Webhook,
			Channel:     op.Channel,
			AdminID:     op.AdminID,
			BannerText:  op.BannerText,
			BannerColor: op.BannerColor,
			Admin: &types.Operator{
				ID:       op.Admin.ID,
				Username: op.Admin.Username,
				Admin:    op.Admin.Admin,
			},
		}
	}

	return operations, nil
}

// GetOperationByID retrieves a specific operation by ID.
func (c *Client) GetOperationByID(ctx context.Context, operationID int) (*types.Operation, error) {
	if err := c.EnsureAuthenticated(ctx); err != nil {
		return nil, err
	}

	var query struct {
		Operation []struct {
			ID          int    `graphql:"id"`
			Name        string `graphql:"name"`
			Complete    bool   `graphql:"complete"`
			Webhook     string `graphql:"webhook"`
			Channel     string `graphql:"channel"`
			AdminID     int    `graphql:"admin_id"`
			BannerText  string `graphql:"banner_text"`
			BannerColor string `graphql:"banner_color"`
			Admin       struct {
				ID       int    `graphql:"id"`
				Username string `graphql:"username"`
				Admin    bool   `graphql:"admin"`
			} `graphql:"admin"`
		} `graphql:"operation(where: {id: {_eq: $id}})"`
	}

	variables := map[string]interface{}{
		"id": operationID,
	}

	err := c.executeQuery(ctx, &query, variables)
	if err != nil {
		return nil, WrapError("GetOperationByID", err, "failed to query operation")
	}

	if len(query.Operation) == 0 {
		return nil, WrapError("GetOperationByID", ErrNotFound, fmt.Sprintf("operation %d not found", operationID))
	}

	op := query.Operation[0]
	return &types.Operation{
		ID:          op.ID,
		Name:        op.Name,
		Complete:    op.Complete,
		Webhook:     op.Webhook,
		Channel:     op.Channel,
		AdminID:     op.AdminID,
		BannerText:  op.BannerText,
		BannerColor: op.BannerColor,
		Admin: &types.Operator{
			ID:       op.Admin.ID,
			Username: op.Admin.Username,
			Admin:    op.Admin.Admin,
		},
	}, nil
}

// CreateOperation creates a new operation.
func (c *Client) CreateOperation(ctx context.Context, req *types.CreateOperationRequest) (*types.Operation, error) {
	if err := c.EnsureAuthenticated(ctx); err != nil {
		return nil, err
	}

	if req == nil || req.Name == "" {
		return nil, WrapError("CreateOperation", ErrInvalidInput, "operation name is required")
	}

	var mutation struct {
		CreateOperation struct {
			Status      string `graphql:"status"`
			Error       string `graphql:"error"`
			OperationID int    `graphql:"operation_id"`
		} `graphql:"createOperation(name: $name)"`
	}

	variables := map[string]interface{}{
		"name": req.Name,
	}

	err := c.executeMutation(ctx, &mutation, variables)
	if err != nil {
		return nil, WrapError("CreateOperation", err, "failed to create operation")
	}

	if mutation.CreateOperation.Status != "success" {
		return nil, WrapError("CreateOperation", ErrOperationFailed, mutation.CreateOperation.Error)
	}

	// Fetch the full operation details
	return c.GetOperationByID(ctx, mutation.CreateOperation.OperationID)
}

// UpdateOperation updates an existing operation using the GraphQL mutation.
func (c *Client) UpdateOperation(ctx context.Context, req *types.UpdateOperationRequest) (*types.Operation, error) {
	if err := c.EnsureAuthenticated(ctx); err != nil {
		return nil, err
	}

	if req == nil || req.OperationID == 0 {
		return nil, WrapError("UpdateOperation", ErrInvalidInput, "operation ID is required")
	}

	// Check if there are any fields to update
	hasUpdates := req.Name != nil || req.Channel != nil || req.Complete != nil ||
		req.Webhook != nil || req.AdminID != nil || req.BannerText != nil || req.BannerColor != nil

	if !hasUpdates {
		return nil, WrapError("UpdateOperation", ErrInvalidInput, "no fields to update")
	}

	// Build dynamic GraphQL mutation using updateOperation custom mutation
	varDecls := []string{"$operation_id: Int!"}
	variables := map[string]interface{}{
		"operation_id": req.OperationID,
	}
	argParts := []string{"operation_id: $operation_id"}

	if req.Name != nil {
		varDecls = append(varDecls, "$name: String")
		variables["name"] = *req.Name
		argParts = append(argParts, "name: $name")
	}
	if req.Channel != nil {
		varDecls = append(varDecls, "$channel: String")
		variables["channel"] = *req.Channel
		argParts = append(argParts, "channel: $channel")
	}
	if req.Complete != nil {
		varDecls = append(varDecls, "$complete: Boolean")
		variables["complete"] = *req.Complete
		argParts = append(argParts, "complete: $complete")
	}
	if req.Webhook != nil {
		varDecls = append(varDecls, "$webhook: String")
		variables["webhook"] = *req.Webhook
		argParts = append(argParts, "webhook: $webhook")
	}
	if req.AdminID != nil {
		varDecls = append(varDecls, "$admin_id: Int")
		variables["admin_id"] = *req.AdminID
		argParts = append(argParts, "admin_id: $admin_id")
	}
	if req.BannerText != nil {
		varDecls = append(varDecls, "$banner_text: String")
		variables["banner_text"] = *req.BannerText
		argParts = append(argParts, "banner_text: $banner_text")
	}
	if req.BannerColor != nil {
		varDecls = append(varDecls, "$banner_color: String")
		variables["banner_color"] = *req.BannerColor
		argParts = append(argParts, "banner_color: $banner_color")
	}

	query := fmt.Sprintf(`mutation UpdateOperation(%s) {
		updateOperation(%s) {
			status
			error
		}
	}`, strings.Join(varDecls, ", "), strings.Join(argParts, ", "))

	result, err := c.ExecuteRawGraphQL(ctx, query, variables)
	if err != nil {
		return nil, WrapError("UpdateOperation", err, "failed to update operation")
	}

	if updateResult, ok := result["updateOperation"].(map[string]interface{}); ok {
		if status, ok := updateResult["status"].(string); ok && status != "success" {
			errMsg := ""
			if e, ok := updateResult["error"].(string); ok {
				errMsg = e
			}
			return nil, WrapError("UpdateOperation", ErrOperationFailed, errMsg)
		}
	}

	// Fetch the updated operation
	return c.GetOperationByID(ctx, req.OperationID)
}

// UpdateCurrentOperationForUser switches the current operation for the authenticated user.
func (c *Client) UpdateCurrentOperationForUser(ctx context.Context, operationID int) error {
	if err := c.EnsureAuthenticated(ctx); err != nil {
		return err
	}

	// Get current user ID
	me, err := c.GetMe(ctx)
	if err != nil {
		return WrapError("UpdateCurrentOperationForUser", err, "failed to get current user")
	}

	var mutation struct {
		UpdateCurrentOperation struct {
			Status string `graphql:"status"`
		} `graphql:"updateCurrentOperation(user_id: $user_id, operation_id: $operation_id)"`
	}

	variables := map[string]interface{}{
		"user_id":      me.ID,
		"operation_id": operationID,
	}

	err = c.executeMutation(ctx, &mutation, variables)
	if err != nil {
		return WrapError("UpdateCurrentOperationForUser", err, "failed to update current operation")
	}

	// Update the client's current operation ID
	c.SetCurrentOperation(operationID)

	return nil
}

// GetOperatorsByOperation lists all operators in an operation.
func (c *Client) GetOperatorsByOperation(ctx context.Context, operationID int) ([]*types.OperationOperator, error) {
	if err := c.EnsureAuthenticated(ctx); err != nil {
		return nil, err
	}

	// Verify operation exists
	_, err := c.GetOperationByID(ctx, operationID)
	if err != nil {
		return nil, WrapError("GetOperatorsByOperation", err, "operation not found")
	}

	var query struct {
		OperatorOperation []struct {
			ID          int `graphql:"id"`
			OperationID int `graphql:"operation_id"`
			OperatorID  int `graphql:"operator_id"`
			Operator    struct {
				ID       int    `graphql:"id"`
				Username string `graphql:"username"`
				Admin    bool   `graphql:"admin"`
			} `graphql:"operator"`
		} `graphql:"operatoroperation(where: {operation_id: {_eq: $operation_id}})"`
	}

	variables := map[string]interface{}{
		"operation_id": operationID,
	}

	err = c.executeQuery(ctx, &query, variables)
	if err != nil {
		return nil, WrapError("GetOperatorsByOperation", err, "failed to query operators")
	}

	operators := make([]*types.OperationOperator, len(query.OperatorOperation))
	for i, opOp := range query.OperatorOperation {
		operators[i] = &types.OperationOperator{
			ID:          opOp.ID,
			OperationID: opOp.OperationID,
			OperatorID:  opOp.OperatorID,
			Operator: &types.Operator{
				ID:       opOp.Operator.ID,
				Username: opOp.Operator.Username,
				Admin:    opOp.Operator.Admin,
			},
		}
	}

	return operators, nil
}

// UpdateOperatorOperation adds/removes/updates operator(s) in an operation.
// Supports bulk operations and view-only permissions.
func (c *Client) UpdateOperatorOperation(ctx context.Context, req *types.UpdateOperatorOperationRequest) error {
	if err := c.EnsureAuthenticated(ctx); err != nil {
		return err
	}

	if req == nil || req.OperationID == 0 {
		return WrapError("UpdateOperatorOperation", ErrInvalidInput, "operation ID is required")
	}

	// Build variables map with all provided parameters
	// GraphQL signature: updateOperatorOperation(operation_id: Int!, add_users: [Int],
	//                    remove_users: [Int], view_mode_operators: [Int], view_mode_spectators: [Int])

	// Handle new array-based fields
	addUsers := req.AddUsers
	removeUsers := req.RemoveUsers

	// Handle legacy single-operator fields for backwards compatibility
	if req.OperatorID != 0 {
		if req.Remove {
			// Legacy: Remove single operator
			if len(removeUsers) == 0 {
				removeUsers = []int{req.OperatorID}
			}
		} else {
			// Legacy: Add single operator
			if len(addUsers) == 0 {
				addUsers = []int{req.OperatorID}
			}
		}
	}

	// Build variables map with empty slices for unset parameters
	// NOTE: We cannot omit parameters because the GraphQL mutation references them all
	// Using empty slices [] instead of nil to avoid "null value" validation errors
	addUsersVar := []int{}
	if len(addUsers) > 0 {
		addUsersVar = addUsers
	}

	removeUsersVar := []int{}
	if len(removeUsers) > 0 {
		removeUsersVar = removeUsers
	}

	viewOpsVar := []int{}
	if len(req.ViewModeOperators) > 0 {
		viewOpsVar = req.ViewModeOperators
	}

	viewSpecsVar := []int{}
	if len(req.ViewModeSpectators) > 0 {
		viewSpecsVar = req.ViewModeSpectators
	}

	variables := map[string]interface{}{
		"operation_id":         req.OperationID,
		"add_users":            addUsersVar,
		"remove_users":         removeUsersVar,
		"view_mode_operators":  viewOpsVar,
		"view_mode_spectators": viewSpecsVar,
	}

	query := `mutation UpdateOperatorOperation(
		$operation_id: Int!,
		$add_users: [Int],
		$remove_users: [Int],
		$view_mode_operators: [Int],
		$view_mode_spectators: [Int]
	) {
		updateOperatorOperation(
			operation_id: $operation_id,
			add_users: $add_users,
			remove_users: $remove_users,
			view_mode_operators: $view_mode_operators,
			view_mode_spectators: $view_mode_spectators
		) {
			status
			error
		}
	}`

	result, err := c.ExecuteRawGraphQL(ctx, query, variables)
	if err != nil {
		return WrapError("UpdateOperatorOperation", err, "failed to update operator in operation")
	}

	if updateResult, ok := result["updateOperatorOperation"].(map[string]interface{}); ok {
		if status, ok := updateResult["status"].(string); ok && status != "success" {
			errMsg := ""
			if e, ok := updateResult["error"].(string); ok {
				errMsg = e
			}
			return WrapError("UpdateOperatorOperation", ErrOperationFailed, errMsg)
		}
	}

	return nil
}

// GetOperationEventLog retrieves event logs for an operation.
func (c *Client) GetOperationEventLog(ctx context.Context, operationID int, limit int) ([]*types.OperationEventLog, error) {
	if err := c.EnsureAuthenticated(ctx); err != nil {
		return nil, err
	}

	if limit <= 0 {
		limit = 100 // Default limit
	}

	var query struct {
		OperationEventLog []struct {
			ID          int    `graphql:"id"`
			OperatorID  int    `graphql:"operator_id"`
			OperationID int    `graphql:"operation_id"`
			Message     string `graphql:"message"`
			Timestamp   string `graphql:"timestamp"`
			Level       string `graphql:"level"`
			Source      string `graphql:"source"`
			Deleted     bool   `graphql:"deleted"`
			Operator    struct {
				ID       int    `graphql:"id"`
				Username string `graphql:"username"`
			} `graphql:"operator"`
		} `graphql:"operationeventlog(where: {operation_id: {_eq: $operation_id}, deleted: {_eq: false}}, order_by: {timestamp: desc}, limit: $limit)"`
	}

	variables := map[string]interface{}{
		"operation_id": operationID,
		"limit":        limit,
	}

	err := c.executeQuery(ctx, &query, variables)
	if err != nil {
		return nil, WrapError("GetOperationEventLog", err, "failed to query operation event log")
	}

	logs := make([]*types.OperationEventLog, len(query.OperationEventLog))
	for i, log := range query.OperationEventLog {
		timestamp, _ := parseTime(log.Timestamp) //nolint:errcheck // Timestamp parse errors not critical
		logs[i] = &types.OperationEventLog{
			ID:          log.ID,
			OperatorID:  log.OperatorID,
			OperationID: log.OperationID,
			Message:     log.Message,
			Timestamp:   timestamp,
			Level:       log.Level,
			Source:      log.Source,
			Deleted:     log.Deleted,
			Operator: &types.Operator{
				ID:       log.Operator.ID,
				Username: log.Operator.Username,
			},
		}
	}

	return logs, nil
}

// CreateOperationEventLog creates a new event log entry for an operation.
func (c *Client) CreateOperationEventLog(ctx context.Context, req *types.CreateOperationEventLogRequest) (*types.OperationEventLog, error) {
	if err := c.EnsureAuthenticated(ctx); err != nil {
		return nil, err
	}

	if req == nil || req.OperationID == 0 || req.Message == "" {
		return nil, WrapError("CreateOperationEventLog", ErrInvalidInput, "operation ID and message are required")
	}

	// Default level to "info" if not specified
	level := "info"
	if req.Level != "" {
		level = req.Level
	}

	// Default source to "sdk" if not specified
	source := "sdk"
	if req.Source != "" {
		source = req.Source
	}

	// NOTE: In Mythic v3.4.20, operationeventlog insert does not accept operation_id parameter
	// Event logs are automatically created for the operator's current operation
	// To create a log for a specific operation, switch to that operation first using UpdateCurrentOperationForUser

	var mutation struct {
		CreateOperationEventLog struct {
			Returning []struct {
				ID int `graphql:"id"`
			} `graphql:"returning"`
		} `graphql:"insert_operationeventlog(objects: [{message: $message, level: $level, source: $source}])"`
	}

	variables := map[string]interface{}{
		"message": req.Message,
		"level":   level,
		"source":  source,
	}

	err := c.executeMutation(ctx, &mutation, variables)
	if err != nil {
		return nil, WrapError("CreateOperationEventLog", err, "failed to create operation event log")
	}

	if len(mutation.CreateOperationEventLog.Returning) == 0 {
		return nil, WrapError("CreateOperationEventLog", ErrInvalidResponse, "no event log created")
	}

	logID := mutation.CreateOperationEventLog.Returning[0].ID

	// Fetch the created log entry
	var query struct {
		OperationEventLog []struct {
			ID          int    `graphql:"id"`
			OperatorID  int    `graphql:"operator_id"`
			OperationID int    `graphql:"operation_id"`
			Message     string `graphql:"message"`
			Timestamp   string `graphql:"timestamp"`
			Level       string `graphql:"level"`
			Source      string `graphql:"source"`
			Deleted     bool   `graphql:"deleted"`
		} `graphql:"operationeventlog(where: {id: {_eq: $id}})"`
	}

	queryVars := map[string]interface{}{
		"id": logID,
	}

	err = c.executeQuery(ctx, &query, queryVars)
	if err != nil {
		return nil, WrapError("CreateOperationEventLog", err, "failed to fetch created event log")
	}

	if len(query.OperationEventLog) == 0 {
		return nil, WrapError("CreateOperationEventLog", ErrNotFound, "created event log not found")
	}

	log := query.OperationEventLog[0]
	timestamp, _ := parseTime(log.Timestamp) //nolint:errcheck // Timestamp parse errors not critical

	return &types.OperationEventLog{
		ID:          log.ID,
		OperatorID:  log.OperatorID,
		OperationID: log.OperationID,
		Message:     log.Message,
		Timestamp:   timestamp,
		Level:       log.Level,
		Source:      log.Source,
		Deleted:     log.Deleted,
	}, nil
}

// GetGlobalSettings retrieves Mythic global settings.
func (c *Client) GetGlobalSettings(ctx context.Context) (map[string]interface{}, error) {
	if err := c.EnsureAuthenticated(ctx); err != nil {
		return nil, err
	}

	// Note: Global settings table/query not available in all Mythic versions
	// Return empty map for now - this functionality may require admin access
	// or may be version-specific
	return make(map[string]interface{}), nil
}

// UpdateGlobalSettings updates Mythic global settings.
// Note: This feature is not available in the current Mythic GraphQL schema.
// The mutation does not exist or requires admin-level access not exposed via GraphQL.
func (c *Client) UpdateGlobalSettings(ctx context.Context, settings map[string]interface{}) error {
	if err := c.EnsureAuthenticated(ctx); err != nil {
		return err
	}

	if len(settings) == 0 {
		return WrapError("UpdateGlobalSettings", ErrInvalidInput, "settings cannot be empty")
	}

	// UpdateGlobalSettings mutation not available in GraphQL schema
	// Would need REST API endpoint or admin panel access
	return WrapError("UpdateGlobalSettings", ErrOperationFailed, "global settings updates not supported via GraphQL API")
}
