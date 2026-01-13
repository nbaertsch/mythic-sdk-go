package mythic

import (
	"context"
	"encoding/json"
	"time"

	"github.com/nbaertsch/mythic-sdk-go/pkg/mythic/types"
)

// GetOperators retrieves all operators (users) in the system.
func (c *Client) GetOperators(ctx context.Context) ([]*types.Operator, error) {
	if err := c.EnsureAuthenticated(ctx); err != nil {
		return nil, err
	}

	var query struct {
		Operators []struct {
			ID               int    `graphql:"id"`
			Username         string `graphql:"username"`
			Admin            bool   `graphql:"admin"`
			Active           bool   `graphql:"active"`
			Deleted          bool   `graphql:"deleted"`
			CurrentOperation *int   `graphql:"current_operation_id"`
			AccountType      string `graphql:"account_type"`
			FailedLoginCount int    `graphql:"failed_login_count"`
		} `graphql:"operator(order_by: {username: asc})"`
	}

	err := c.executeQuery(ctx, &query, nil)
	if err != nil {
		return nil, WrapError("GetOperators", err, "failed to query operators")
	}

	operators := make([]*types.Operator, len(query.Operators))
	for i, op := range query.Operators {
		operators[i] = &types.Operator{
			ID:                 op.ID,
			Username:           op.Username,
			Admin:              op.Admin,
			Active:             op.Active,
			Deleted:            op.Deleted,
			CurrentOperationID: op.CurrentOperation,
			AccountType:        op.AccountType,
			FailedLoginCount:   op.FailedLoginCount,
		}
	}

	return operators, nil
}

// GetOperatorByID retrieves a specific operator by ID.
func (c *Client) GetOperatorByID(ctx context.Context, operatorID int) (*types.Operator, error) {
	if err := c.EnsureAuthenticated(ctx); err != nil {
		return nil, err
	}

	if operatorID == 0 {
		return nil, WrapError("GetOperatorByID", ErrInvalidInput, "operator ID is required")
	}

	var query struct {
		Operators []struct {
			ID               int    `graphql:"id"`
			Username         string `graphql:"username"`
			Admin            bool   `graphql:"admin"`
			Active           bool   `graphql:"active"`
			Deleted          bool   `graphql:"deleted"`
			CurrentOperation *int   `graphql:"current_operation_id"`
			AccountType      string `graphql:"account_type"`
			FailedLoginCount int    `graphql:"failed_login_count"`
		} `graphql:"operator(where: {id: {_eq: $operator_id}})"`
	}

	variables := map[string]interface{}{
		"operator_id": operatorID,
	}

	err := c.executeQuery(ctx, &query, variables)
	if err != nil {
		return nil, WrapError("GetOperatorByID", err, "failed to query operator")
	}

	if len(query.Operators) == 0 {
		return nil, WrapError("GetOperatorByID", ErrNotFound, "operator not found")
	}

	op := query.Operators[0]
	return &types.Operator{
		ID:                 op.ID,
		Username:           op.Username,
		Admin:              op.Admin,
		Active:             op.Active,
		Deleted:            op.Deleted,
		CurrentOperationID: op.CurrentOperation,
		AccountType:        op.AccountType,
		FailedLoginCount:   op.FailedLoginCount,
	}, nil
}

// CreateOperator creates a new operator account.
// Password must be at least 12 characters long.
func (c *Client) CreateOperator(ctx context.Context, req *types.CreateOperatorRequest) (*types.Operator, error) {
	if err := c.EnsureAuthenticated(ctx); err != nil {
		return nil, err
	}

	if req == nil || req.Username == "" || req.Password == "" {
		return nil, WrapError("CreateOperator", ErrInvalidInput, "username and password are required")
	}

	if len(req.Password) < 12 {
		return nil, WrapError("CreateOperator", ErrInvalidInput, "password must be at least 12 characters long")
	}

	var mutation struct {
		CreateOperator struct {
			Status   string `graphql:"status"`
			Error    string `graphql:"error"`
			ID       int    `graphql:"id"`
			Username string `graphql:"username"`
		} `graphql:"createOperator(username: $username, password: $password)"`
	}

	variables := map[string]interface{}{
		"username": req.Username,
		"password": req.Password,
	}

	err := c.executeMutation(ctx, &mutation, variables)
	if err != nil {
		return nil, WrapError("CreateOperator", err, "failed to create operator")
	}

	if mutation.CreateOperator.Status != "success" {
		return nil, WrapError("CreateOperator", ErrOperationFailed, mutation.CreateOperator.Error)
	}

	// Fetch the created operator details
	return c.GetOperatorByID(ctx, mutation.CreateOperator.ID)
}

// UpdateOperatorStatus updates an operator's status (active, admin, deleted).
func (c *Client) UpdateOperatorStatus(ctx context.Context, req *types.UpdateOperatorStatusRequest) error {
	if err := c.EnsureAuthenticated(ctx); err != nil {
		return err
	}

	if req == nil || req.OperatorID == 0 {
		return WrapError("UpdateOperatorStatus", ErrInvalidInput, "operator ID is required")
	}

	// Mythic does not provide a direct GraphQL mutation for updating operator status
	// Operator status updates require using the REST API or admin interface
	return WrapError("UpdateOperatorStatus", ErrOperationFailed, "operator status updates not supported via GraphQL API")
}

// UpdatePasswordAndEmail updates an operator's password and/or email.
// Old password is required for verification. New password must be at least 12 characters if provided.
func (c *Client) UpdatePasswordAndEmail(ctx context.Context, req *types.UpdatePasswordAndEmailRequest) error {
	if err := c.EnsureAuthenticated(ctx); err != nil {
		return err
	}

	if req == nil || req.OperatorID == 0 || req.OldPassword == "" {
		return WrapError("UpdatePasswordAndEmail", ErrInvalidInput, "operator ID and old password are required")
	}

	if req.NewPassword == nil && req.Email == nil {
		return WrapError("UpdatePasswordAndEmail", ErrInvalidInput, "at least new password or email must be provided")
	}

	if req.NewPassword != nil && len(*req.NewPassword) < 12 {
		return WrapError("UpdatePasswordAndEmail", ErrInvalidInput, "new password must be at least 12 characters long")
	}

	// Build the input object
	input := map[string]interface{}{
		"operator_id":  req.OperatorID,
		"old_password": req.OldPassword,
	}

	if req.NewPassword != nil {
		input["new_password"] = *req.NewPassword
	}
	if req.Email != nil {
		input["email"] = *req.Email
	}

	var mutation struct {
		UpdatePasswordAndEmail struct {
			Status string `graphql:"status"`
			Error  string `graphql:"error"`
		} `graphql:"updatePasswordAndEmail(input: $input)"`
	}

	variables := map[string]interface{}{
		"input": input,
	}

	err := c.executeMutation(ctx, &mutation, variables)
	if err != nil {
		return WrapError("UpdatePasswordAndEmail", err, "failed to update password and email")
	}

	if mutation.UpdatePasswordAndEmail.Status != "success" {
		return WrapError("UpdatePasswordAndEmail", ErrOperationFailed, mutation.UpdatePasswordAndEmail.Error)
	}

	return nil
}

// GetOperatorPreferences retrieves UI preferences for the currently authenticated operator.
// Note: This function returns preferences for the current operator, regardless of the operatorID parameter.
// This is a Mythic API limitation - preferences are retrieved based on the JWT token.
func (c *Client) GetOperatorPreferences(ctx context.Context, operatorID int) (*types.OperatorPreferences, error) {
	if err := c.EnsureAuthenticated(ctx); err != nil {
		return nil, err
	}

	if operatorID == 0 {
		return nil, WrapError("GetOperatorPreferences", ErrInvalidInput, "operator ID is required")
	}

	// Use REST webhook (GraphQL getOperatorPreferences has issues with jsonb decoding)
	var response struct {
		Status      string                 `json:"status"`
		Error       string                 `json:"error"`
		Preferences map[string]interface{} `json:"preferences"`
	}

	// Empty request body - operator is determined from JWT
	requestData := map[string]interface{}{}

	err := c.executeRESTWebhook(ctx, "api/v1.4/operator_get_preferences_webhook", requestData, &response)
	if err != nil {
		return nil, WrapError("GetOperatorPreferences", err, "failed to execute webhook")
	}

	if response.Status != "success" {
		return nil, WrapError("GetOperatorPreferences", ErrOperationFailed, response.Error)
	}

	// Marshal preferences map to JSON string for storage
	prefsJSON, err := json.Marshal(response.Preferences)
	if err != nil {
		return nil, WrapError("GetOperatorPreferences", err, "failed to marshal preferences")
	}

	return &types.OperatorPreferences{
		OperatorID:      operatorID,
		PreferencesJSON: string(prefsJSON),
	}, nil
}

// UpdateOperatorPreferences updates UI preferences for an operator using the REST API webhook.
func (c *Client) UpdateOperatorPreferences(ctx context.Context, req *types.UpdateOperatorPreferencesRequest) error {
	if err := c.EnsureAuthenticated(ctx); err != nil {
		return err
	}

	if req == nil || req.OperatorID == 0 {
		return WrapError("UpdateOperatorPreferences", ErrInvalidInput, "operator ID is required")
	}

	if len(req.Preferences) == 0 {
		return WrapError("UpdateOperatorPreferences", ErrInvalidInput, "preferences must not be empty")
	}

	// Build REST API request using Mythic's webhook format
	// Note: Mythic webhook expects parameters wrapped in "Input" object
	requestData := map[string]interface{}{
		"Input": map[string]interface{}{
			"preferences": req.Preferences,
		},
	}

	var response struct {
		Status string `json:"status"`
		Error  string `json:"error"`
	}

	err := c.executeRESTWebhook(ctx, "api/v1.4/operator_update_preferences_webhook", requestData, &response)
	if err != nil {
		return WrapError("UpdateOperatorPreferences", err, "failed to execute webhook")
	}

	if response.Status != "success" {
		return WrapError("UpdateOperatorPreferences", ErrOperationFailed, response.Error)
	}

	return nil
}

// GetOperatorSecrets retrieves secrets/keys for an operator.
func (c *Client) GetOperatorSecrets(ctx context.Context, operatorID int) (*types.OperatorSecrets, error) {
	if err := c.EnsureAuthenticated(ctx); err != nil {
		return nil, err
	}

	if operatorID == 0 {
		return nil, WrapError("GetOperatorSecrets", ErrInvalidInput, "operator ID is required")
	}

	var query struct {
		Secrets struct {
			Status  string `graphql:"status"`
			Error   string `graphql:"error"`
			Secrets string `graphql:"secrets"`
		} `graphql:"getOperatorSecrets(operator_id: $operator_id)"`
	}

	variables := map[string]interface{}{
		"operator_id": operatorID,
	}

	err := c.executeQuery(ctx, &query, variables)
	if err != nil {
		return nil, WrapError("GetOperatorSecrets", err, "failed to query operator secrets")
	}

	if query.Secrets.Status != "success" {
		return nil, WrapError("GetOperatorSecrets", ErrOperationFailed, query.Secrets.Error)
	}

	return &types.OperatorSecrets{
		OperatorID:  operatorID,
		SecretsJSON: query.Secrets.Secrets,
	}, nil
}

// UpdateOperatorSecrets updates secrets/keys for an operator.
func (c *Client) UpdateOperatorSecrets(ctx context.Context, req *types.UpdateOperatorSecretsRequest) error {
	if err := c.EnsureAuthenticated(ctx); err != nil {
		return err
	}

	if req == nil || req.OperatorID == 0 {
		return WrapError("UpdateOperatorSecrets", ErrInvalidInput, "operator ID is required")
	}

	if len(req.Secrets) == 0 {
		return WrapError("UpdateOperatorSecrets", ErrInvalidInput, "secrets must not be empty")
	}

	// Build REST API request using Mythic's webhook format
	requestData := map[string]interface{}{
		"Input": map[string]interface{}{
			"operator_id": req.OperatorID,
			"secrets":     req.Secrets,
		},
	}

	var response struct {
		Status string `json:"status"`
		Error  string `json:"error"`
	}

	err := c.executeRESTWebhook(ctx, "api/v1.4/operator_update_secrets_webhook", requestData, &response)
	if err != nil {
		return WrapError("UpdateOperatorSecrets", err, "failed to execute webhook")
	}

	if response.Status != "success" {
		return WrapError("UpdateOperatorSecrets", ErrOperationFailed, response.Error)
	}

	return nil
}

// GetInviteLinks retrieves all invitation links for new operators.
func (c *Client) GetInviteLinks(ctx context.Context) ([]*types.InviteLink, error) {
	if err := c.EnsureAuthenticated(ctx); err != nil {
		return nil, err
	}

	var query struct {
		InviteLinks struct {
			Status string `graphql:"status"`
			Error  string `graphql:"error"`
			Links  []struct {
				ID          int    `graphql:"id"`
				Code        string `graphql:"code"`
				ExpiresAt   string `graphql:"expires_at"`
				CreatedBy   int    `graphql:"created_by"`
				CreatedAt   string `graphql:"created_at"`
				MaxUses     int    `graphql:"max_uses"`
				CurrentUses int    `graphql:"current_uses"`
				Active      bool   `graphql:"active"`
			} `graphql:"links"`
		} `graphql:"getInviteLinks"`
	}

	err := c.executeQuery(ctx, &query, nil)
	if err != nil {
		return nil, WrapError("GetInviteLinks", err, "failed to query invite links")
	}

	if query.InviteLinks.Status != "success" {
		return nil, WrapError("GetInviteLinks", ErrOperationFailed, query.InviteLinks.Error)
	}

	links := make([]*types.InviteLink, len(query.InviteLinks.Links))
	for i, link := range query.InviteLinks.Links {
		expiresAt, _ := parseTime(link.ExpiresAt) //nolint:errcheck // Timestamp parse errors not critical
		createdAt, _ := parseTime(link.CreatedAt) //nolint:errcheck // Timestamp parse errors not critical

		links[i] = &types.InviteLink{
			ID:          link.ID,
			Code:        link.Code,
			ExpiresAt:   expiresAt,
			CreatedBy:   link.CreatedBy,
			CreatedAt:   createdAt,
			MaxUses:     link.MaxUses,
			CurrentUses: link.CurrentUses,
			Active:      link.Active,
		}
	}

	return links, nil
}

// CreateInviteLink creates a new invitation link for new operators.
func (c *Client) CreateInviteLink(ctx context.Context, req *types.CreateInviteLinkRequest) (*types.InviteLink, error) {
	if err := c.EnsureAuthenticated(ctx); err != nil {
		return nil, err
	}

	if req == nil || req.MaxUses <= 0 {
		return nil, WrapError("CreateInviteLink", ErrInvalidInput, "max uses must be greater than 0")
	}

	if req.ExpiresAt.Before(time.Now()) {
		return nil, WrapError("CreateInviteLink", ErrInvalidInput, "expiration date must be in the future")
	}

	var mutation struct {
		CreateInviteLink struct {
			Status  string `graphql:"status"`
			Error   string `graphql:"error"`
			ID      int    `graphql:"id"`
			Code    string `graphql:"code"`
			MaxUses int    `graphql:"max_uses"`
		} `graphql:"createInviteLink(max_uses: $max_uses, expires_at: $expires_at)"`
	}

	variables := map[string]interface{}{
		"max_uses":   req.MaxUses,
		"expires_at": req.ExpiresAt.Format(time.RFC3339),
	}

	err := c.executeMutation(ctx, &mutation, variables)
	if err != nil {
		return nil, WrapError("CreateInviteLink", err, "failed to create invite link")
	}

	if mutation.CreateInviteLink.Status != "success" {
		return nil, WrapError("CreateInviteLink", ErrOperationFailed, mutation.CreateInviteLink.Error)
	}

	return &types.InviteLink{
		ID:          mutation.CreateInviteLink.ID,
		Code:        mutation.CreateInviteLink.Code,
		ExpiresAt:   req.ExpiresAt,
		MaxUses:     mutation.CreateInviteLink.MaxUses,
		CurrentUses: 0,
		Active:      true,
	}, nil
}

// parseTime is a helper function to parse time strings from Mythic API.
func parseTime(s string) (time.Time, error) {
	if s == "" || s == "null" {
		return time.Time{}, nil
	}

	// Try RFC3339 first (standard format with timezone)
	t, err := time.Parse(time.RFC3339, s)
	if err == nil {
		return t, nil
	}

	// Try RFC3339 with nanoseconds
	t, err = time.Parse(time.RFC3339Nano, s)
	if err == nil {
		return t, nil
	}

	// Try Mythic's format without timezone (treat as UTC)
	formats := []string{
		"2006-01-02T15:04:05.999999",
		"2006-01-02T15:04:05",
		"2006-01-02 15:04:05.999999",
		"2006-01-02 15:04:05",
	}

	for _, format := range formats {
		t, err = time.Parse(format, s)
		if err == nil {
			return t.UTC(), nil
		}
	}

	return time.Time{}, err
}
