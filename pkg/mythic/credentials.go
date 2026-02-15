package mythic

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/nbaertsch/mythic-sdk-go/pkg/mythic/types"
)

// credentialFieldGraphQLTypes maps credential _set field names to their GraphQL type declarations.
var credentialFieldGraphQLTypes = map[string]string{
	"account":        "String",
	"comment":        "String",
	"credential_raw": "String",
	"deleted":        "Boolean",
	"metadata":       "String",
	"realm":          "String",
	"type":           "String",
}

// buildCredentialVarDeclarations returns the GraphQL variable declarations for a set of fields.
// e.g. "$comment: String, $deleted: Boolean"
func buildCredentialVarDeclarations(fields map[string]interface{}) string {
	parts := make([]string, 0, len(fields))
	// Sort for deterministic output
	keys := make([]string, 0, len(fields))
	for k := range fields {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		gqlType, ok := credentialFieldGraphQLTypes[k]
		if !ok {
			gqlType = "String" // safe default
		}
		parts = append(parts, fmt.Sprintf("$%s: %s", k, gqlType))
	}
	return strings.Join(parts, ", ")
}

// parseTimestamp parses Mythic's timestamp format (RFC3339 without timezone)
func parseCredentialTimestamp(ts string) (time.Time, error) {
	// Try multiple formats that Mythic might use
	formats := []string{
		"2006-01-02T15:04:05.999999", // Microseconds, no timezone
		"2006-01-02T15:04:05",        // No fractional seconds
		time.RFC3339,                 // With timezone
		time.RFC3339Nano,             // With nanoseconds and timezone
	}

	var lastErr error
	for _, format := range formats {
		t, err := time.Parse(format, ts)
		if err == nil {
			return t, nil
		}
		lastErr = err
	}

	return time.Time{}, lastErr
}

// GetCredentials retrieves all credentials for the current operation.
func (c *Client) GetCredentials(ctx context.Context) ([]*types.Credential, error) {
	if err := c.EnsureAuthenticated(ctx); err != nil {
		return nil, err
	}

	var query struct {
		Credential []struct {
			ID          int    `graphql:"id"`
			Type        string `graphql:"type"`
			Account     string `graphql:"account"`
			Realm       string `graphql:"realm"`
			Credential  string `graphql:"credential_text"`
			Comment     string `graphql:"comment"`
			OperationID int    `graphql:"operation_id"`
			OperatorID  int    `graphql:"operator_id"`
			TaskID      *int   `graphql:"task_id"`
			Timestamp   string `graphql:"timestamp"`
			Deleted     bool   `graphql:"deleted"`
			Metadata    string `graphql:"metadata"`
		} `graphql:"credential(where: {deleted: {_eq: false}}, order_by: {timestamp: desc})"`
	}

	err := c.executeQuery(ctx, &query, nil)
	if err != nil {
		return nil, WrapError("GetCredentials", err, "failed to query credentials")
	}

	credentials := make([]*types.Credential, len(query.Credential))
	for i, cred := range query.Credential {
		timestamp, err := parseCredentialTimestamp(cred.Timestamp)
		if err != nil {
			return nil, WrapError("GetCredentials", err, "failed to parse credential timestamp")
		}

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
			Timestamp:   timestamp,
			Deleted:     cred.Deleted,
			Metadata:    cred.Metadata,
		}
	}

	return credentials, nil
}

// GetCredentialByID retrieves a specific credential by ID.
func (c *Client) GetCredentialByID(ctx context.Context, credentialID int) (*types.Credential, error) {
	if err := c.EnsureAuthenticated(ctx); err != nil {
		return nil, err
	}

	if credentialID == 0 {
		return nil, WrapError("GetCredentialByID", ErrInvalidInput, "credential ID is required")
	}

	var query struct {
		Credential []struct {
			ID          int    `graphql:"id"`
			Type        string `graphql:"type"`
			Account     string `graphql:"account"`
			Realm       string `graphql:"realm"`
			Credential  string `graphql:"credential_text"`
			Comment     string `graphql:"comment"`
			OperationID int    `graphql:"operation_id"`
			OperatorID  int    `graphql:"operator_id"`
			TaskID      *int   `graphql:"task_id"`
			Timestamp   string `graphql:"timestamp"`
			Deleted     bool   `graphql:"deleted"`
			Metadata    string `graphql:"metadata"`
		} `graphql:"credential(where: {id: {_eq: $id}})"`
	}

	variables := map[string]interface{}{
		"id": credentialID,
	}

	err := c.executeQuery(ctx, &query, variables)
	if err != nil {
		return nil, WrapError("GetCredentialByID", err, "failed to query credential")
	}

	if len(query.Credential) == 0 {
		return nil, WrapError("GetCredentialByID", ErrNotFound, "credential not found")
	}

	cred := query.Credential[0]
	timestamp, err := parseCredentialTimestamp(cred.Timestamp)
	if err != nil {
		return nil, WrapError("GetCredentialByID", err, "failed to parse credential timestamp")
	}

	return &types.Credential{
		ID:          cred.ID,
		Type:        cred.Type,
		Account:     cred.Account,
		Realm:       cred.Realm,
		Credential:  cred.Credential,
		Comment:     cred.Comment,
		OperationID: cred.OperationID,
		OperatorID:  cred.OperatorID,
		TaskID:      cred.TaskID,
		Timestamp:   timestamp,
		Deleted:     cred.Deleted,
		Metadata:    cred.Metadata,
	}, nil
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
			ID          int    `graphql:"id"`
			Type        string `graphql:"type"`
			Account     string `graphql:"account"`
			Realm       string `graphql:"realm"`
			Credential  string `graphql:"credential_text"`
			Comment     string `graphql:"comment"`
			OperationID int    `graphql:"operation_id"`
			OperatorID  int    `graphql:"operator_id"`
			TaskID      *int   `graphql:"task_id"`
			Timestamp   string `graphql:"timestamp"`
			Deleted     bool   `graphql:"deleted"`
			Metadata    string `graphql:"metadata"`
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
		timestamp, err := parseCredentialTimestamp(cred.Timestamp)
		if err != nil {
			return nil, WrapError("GetCredentialsByOperation", err, "failed to parse credential timestamp")
		}

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
			Timestamp:   timestamp,
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
			Status string `graphql:"status"`
			Error  string `graphql:"error"`
			ID     int    `graphql:"id"`
		} `graphql:"createCredential(credential_type: $credential_type, account: $account, realm: $realm, credential: $credential, comment: $comment)"`
	}

	variables := map[string]interface{}{
		"credential_type": req.Type,
		"account":         req.Account,
		"realm":           req.Realm,
		"credential":      req.Credential,
		"comment":         req.Comment,
	}

	err := c.executeMutation(ctx, &mutation, variables)
	if err != nil {
		return nil, WrapError("CreateCredential", err, "failed to create credential")
	}

	if mutation.CreateCredential.Status != "success" {
		return nil, WrapError("CreateCredential", ErrOperationFailed, mutation.CreateCredential.Error)
	}

	// Return credential with the provided data
	credential := &types.Credential{
		ID:         mutation.CreateCredential.ID,
		Type:       req.Type,
		Account:    req.Account,
		Realm:      req.Realm,
		Credential: req.Credential,
		Comment:    req.Comment,
		Metadata:   req.Metadata,
	}

	return credential, nil
}

// UpdateCredential updates an existing credential.
// Supports updating: type, account, realm, credential (text), comment, deleted, metadata.
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

	// Check if there are any fields to update
	hasUpdates := req.Type != nil || req.Account != nil || req.Realm != nil ||
		req.Credential != nil || req.Comment != nil || req.Deleted != nil || req.Metadata != nil

	if !hasUpdates {
		return nil, WrapError("UpdateCredential", ErrInvalidInput, "no fields to update")
	}

	// Build the _set object dynamically based on which fields are provided.
	// The GraphQL schema's credential_set_input supports: account, comment,
	// credential_raw (bytea), deleted, metadata, realm, type.
	// Note: credential_text is a computed column; the writable column is credential_raw.
	setFields := make(map[string]interface{})

	if req.Comment != nil {
		setFields["comment"] = *req.Comment
	}
	if req.Deleted != nil {
		setFields["deleted"] = *req.Deleted
	}
	if req.Type != nil {
		setFields["type"] = *req.Type
	}
	if req.Account != nil {
		setFields["account"] = *req.Account
	}
	if req.Realm != nil {
		setFields["realm"] = *req.Realm
	}
	if req.Credential != nil {
		setFields["credential_raw"] = *req.Credential
	}
	if req.Metadata != nil {
		setFields["metadata"] = *req.Metadata
	}

	// Build the _set clause dynamically since the Go GraphQL client
	// requires struct tags at compile time and can't handle dynamic fields.
	// We use a raw GraphQL mutation instead.
	setParts := make([]string, 0, len(setFields))
	setKeys := make([]string, 0, len(setFields))
	for k := range setFields {
		setKeys = append(setKeys, k)
	}
	sort.Strings(setKeys)
	for _, k := range setKeys {
		setParts = append(setParts, fmt.Sprintf("%s: $%s", k, k))
	}

	query := fmt.Sprintf(`mutation UpdateCredential($id: Int!, %s) {
		update_credential(where: {id: {_eq: $id}}, _set: {%s}) {
			affected_rows
		}
	}`, buildCredentialVarDeclarations(setFields), strings.Join(setParts, ", "))

	variables := map[string]interface{}{
		"id": req.ID,
	}
	for k, v := range setFields {
		variables[k] = v
	}

	result, err := c.ExecuteRawGraphQL(ctx, query, variables)
	if err != nil {
		return nil, WrapError("UpdateCredential", err, "failed to update credential")
	}

	// Check affected rows
	if updateCred, ok := result["update_credential"].(map[string]interface{}); ok {
		if affected, ok := updateCred["affected_rows"].(float64); ok && affected == 0 {
			return nil, WrapError("UpdateCredential", ErrNotFound, "credential not found or not updated")
		}
	}

	// Fetch and return the updated credential
	return c.GetCredentialByID(ctx, req.ID)
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
