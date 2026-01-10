package mythic

import (
	"context"
	"time"

	"github.com/nbaertsch/mythic-sdk-go/pkg/mythic/types"
)

// GetCredentials retrieves all credentials for the current operation.
func (c *Client) GetCredentials(ctx context.Context) ([]*types.Credential, error) {
	if err := c.EnsureAuthenticated(ctx); err != nil {
		return nil, err
	}

	var query struct {
		Credential []struct {
			ID          int       `graphql:"id"`
			Type        string    `graphql:"type"`
			Account     string    `graphql:"account"`
			Realm       string    `graphql:"realm"`
			Credential  string    `graphql:"credential_text"`
			Comment     string    `graphql:"comment"`
			OperationID int       `graphql:"operation_id"`
			OperatorID  int       `graphql:"operator_id"`
			TaskID      *int      `graphql:"task_id"`
			Timestamp   time.Time `graphql:"timestamp"`
			Deleted     bool      `graphql:"deleted"`
			Metadata    string    `graphql:"metadata"`
		} `graphql:"credential(where: {deleted: {_eq: false}}, order_by: {timestamp: desc})"`
	}

	err := c.executeQuery(ctx, &query, nil)
	if err != nil {
		return nil, WrapError("GetCredentials", err, "failed to query credentials")
	}

	credentials := make([]*types.Credential, len(query.Credential))
	for i, cred := range query.Credential {
		credentials[i] = &types.Credential{
			ID:          cred.ID,
			Type:        cred.Type,
			Account:     cred.Account,
			Realm:       cred.Realm,
			Credential:  cred.Credential,
			Comment:     cred.Comment,
			OperationID: cred.OperationID,
			OperatorID:  cred.OperatorID,
			TaskID:      cred.TaskID,
			Timestamp:   cred.Timestamp,
			Deleted:     cred.Deleted,
			Metadata:    cred.Metadata,
		}
	}

	return credentials, nil
}

// GetCredentialsByOperation retrieves credentials for a specific operation.
func (c *Client) GetCredentialsByOperation(ctx context.Context, operationID int) ([]*types.Credential, error) {
	if err := c.EnsureAuthenticated(ctx); err != nil {
		return nil, err
	}

	if operationID == 0 {
		return nil, WrapError("GetCredentialsByOperation", ErrInvalidInput, "operation ID is required")
	}

	var query struct {
		Credential []struct {
			ID          int       `graphql:"id"`
			Type        string    `graphql:"type"`
			Account     string    `graphql:"account"`
			Realm       string    `graphql:"realm"`
			Credential  string    `graphql:"credential_text"`
			Comment     string    `graphql:"comment"`
			OperationID int       `graphql:"operation_id"`
			OperatorID  int       `graphql:"operator_id"`
			TaskID      *int      `graphql:"task_id"`
			Timestamp   time.Time `graphql:"timestamp"`
			Deleted     bool      `graphql:"deleted"`
			Metadata    string    `graphql:"metadata"`
		} `graphql:"credential(where: {operation_id: {_eq: $operation_id}, deleted: {_eq: false}}, order_by: {timestamp: desc})"`
	}

	variables := map[string]interface{}{
		"operation_id": operationID,
	}

	err := c.executeQuery(ctx, &query, variables)
	if err != nil {
		return nil, WrapError("GetCredentialsByOperation", err, "failed to query credentials")
	}

	credentials := make([]*types.Credential, len(query.Credential))
	for i, cred := range query.Credential {
		credentials[i] = &types.Credential{
			ID:          cred.ID,
			Type:        cred.Type,
			Account:     cred.Account,
			Realm:       cred.Realm,
			Credential:  cred.Credential,
			Comment:     cred.Comment,
			OperationID: cred.OperationID,
			OperatorID:  cred.OperatorID,
			TaskID:      cred.TaskID,
			Timestamp:   cred.Timestamp,
			Deleted:     cred.Deleted,
			Metadata:    cred.Metadata,
		}
	}

	return credentials, nil
}

// CreateCredential creates a new credential entry.
func (c *Client) CreateCredential(ctx context.Context, req *types.CreateCredentialRequest) (*types.Credential, error) {
	if err := c.EnsureAuthenticated(ctx); err != nil {
		return nil, err
	}

	if req == nil {
		return nil, WrapError("CreateCredential", ErrInvalidInput, "request is required")
	}

	if req.Type == "" {
		return nil, WrapError("CreateCredential", ErrInvalidInput, "credential type is required")
	}

	if req.Account == "" {
		return nil, WrapError("CreateCredential", ErrInvalidInput, "account is required")
	}

	if req.Credential == "" {
		return nil, WrapError("CreateCredential", ErrInvalidInput, "credential value is required")
	}

	// Get current operation if not specified
	currentOp := c.GetCurrentOperation()
	if currentOp == nil {
		return nil, WrapError("CreateCredential", ErrInvalidInput, "no current operation set")
	}

	var mutation struct {
		CreateCredential struct {
			ID          int       `graphql:"id"`
			Type        string    `graphql:"type"`
			Account     string    `graphql:"account"`
			Realm       string    `graphql:"realm"`
			Credential  string    `graphql:"credential_text"`
			Comment     string    `graphql:"comment"`
			OperationID int       `graphql:"operation_id"`
			OperatorID  int       `graphql:"operator_id"`
			TaskID      *int      `graphql:"task_id"`
			Timestamp   time.Time `graphql:"timestamp"`
			Deleted     bool      `graphql:"deleted"`
			Metadata    string    `graphql:"metadata"`
		} `graphql:"createCredential(type: $type, account: $account, realm: $realm, credential_text: $credential_text, comment: $comment, metadata: $metadata)"`
	}

	variables := map[string]interface{}{
		"type":            req.Type,
		"account":         req.Account,
		"realm":           req.Realm,
		"credential_text": req.Credential,
		"comment":         req.Comment,
		"metadata":        req.Metadata,
	}

	err := c.executeMutation(ctx, &mutation, variables)
	if err != nil {
		return nil, WrapError("CreateCredential", err, "failed to create credential")
	}

	credential := &types.Credential{
		ID:          mutation.CreateCredential.ID,
		Type:        mutation.CreateCredential.Type,
		Account:     mutation.CreateCredential.Account,
		Realm:       mutation.CreateCredential.Realm,
		Credential:  mutation.CreateCredential.Credential,
		Comment:     mutation.CreateCredential.Comment,
		OperationID: mutation.CreateCredential.OperationID,
		OperatorID:  mutation.CreateCredential.OperatorID,
		TaskID:      mutation.CreateCredential.TaskID,
		Timestamp:   mutation.CreateCredential.Timestamp,
		Deleted:     mutation.CreateCredential.Deleted,
		Metadata:    mutation.CreateCredential.Metadata,
	}

	return credential, nil
}

// UpdateCredential updates an existing credential.
func (c *Client) UpdateCredential(ctx context.Context, req *types.UpdateCredentialRequest) (*types.Credential, error) {
	if err := c.EnsureAuthenticated(ctx); err != nil {
		return nil, err
	}

	if req == nil {
		return nil, WrapError("UpdateCredential", ErrInvalidInput, "request is required")
	}

	if req.ID == 0 {
		return nil, WrapError("UpdateCredential", ErrInvalidInput, "credential ID is required")
	}

	// Build update map with only non-nil fields
	updates := make(map[string]interface{})
	if req.Type != nil {
		updates["type"] = *req.Type
	}
	if req.Account != nil {
		updates["account"] = *req.Account
	}
	if req.Realm != nil {
		updates["realm"] = *req.Realm
	}
	if req.Credential != nil {
		updates["credential_text"] = *req.Credential
	}
	if req.Comment != nil {
		updates["comment"] = *req.Comment
	}
	if req.Deleted != nil {
		updates["deleted"] = *req.Deleted
	}
	if req.Metadata != nil {
		updates["metadata"] = *req.Metadata
	}

	if len(updates) == 0 {
		return nil, WrapError("UpdateCredential", ErrInvalidInput, "no fields to update")
	}

	var mutation struct {
		UpdateCredential struct {
			Returning []struct {
				ID          int       `graphql:"id"`
				Type        string    `graphql:"type"`
				Account     string    `graphql:"account"`
				Realm       string    `graphql:"realm"`
				Credential  string    `graphql:"credential_text"`
				Comment     string    `graphql:"comment"`
				OperationID int       `graphql:"operation_id"`
				OperatorID  int       `graphql:"operator_id"`
				TaskID      *int      `graphql:"task_id"`
				Timestamp   time.Time `graphql:"timestamp"`
				Deleted     bool      `graphql:"deleted"`
				Metadata    string    `graphql:"metadata"`
			} `graphql:"returning"`
		} `graphql:"update_credential(where: {id: {_eq: $id}}, _set: $updates)"`
	}

	variables := map[string]interface{}{
		"id":      req.ID,
		"updates": updates,
	}

	err := c.executeMutation(ctx, &mutation, variables)
	if err != nil {
		return nil, WrapError("UpdateCredential", err, "failed to update credential")
	}

	if len(mutation.UpdateCredential.Returning) == 0 {
		return nil, WrapError("UpdateCredential", ErrNotFound, "credential not found")
	}

	cred := mutation.UpdateCredential.Returning[0]
	credential := &types.Credential{
		ID:          cred.ID,
		Type:        cred.Type,
		Account:     cred.Account,
		Realm:       cred.Realm,
		Credential:  cred.Credential,
		Comment:     cred.Comment,
		OperationID: cred.OperationID,
		OperatorID:  cred.OperatorID,
		TaskID:      cred.TaskID,
		Timestamp:   cred.Timestamp,
		Deleted:     cred.Deleted,
		Metadata:    cred.Metadata,
	}

	return credential, nil
}

// DeleteCredential marks a credential as deleted.
func (c *Client) DeleteCredential(ctx context.Context, id int) error {
	if err := c.EnsureAuthenticated(ctx); err != nil {
		return err
	}

	if id == 0 {
		return WrapError("DeleteCredential", ErrInvalidInput, "credential ID is required")
	}

	deleted := true
	req := &types.UpdateCredentialRequest{
		ID:      id,
		Deleted: &deleted,
	}

	_, err := c.UpdateCredential(ctx, req)
	if err != nil {
		return WrapError("DeleteCredential", err, "failed to delete credential")
	}

	return nil
}
