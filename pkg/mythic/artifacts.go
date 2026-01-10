package mythic

import (
	"context"
	"time"

	"github.com/nbaertsch/mythic-sdk-go/pkg/mythic/types"
)

// GetArtifacts retrieves all artifacts (IOCs) for the current operation.
func (c *Client) GetArtifacts(ctx context.Context) ([]*types.Artifact, error) {
	if err := c.EnsureAuthenticated(ctx); err != nil {
		return nil, err
	}

	// Use current operation if set
	operationID := c.GetCurrentOperation()
	if operationID == nil {
		return nil, WrapError("GetArtifacts", ErrInvalidInput, "no current operation set")
	}

	return c.GetArtifactsByOperation(ctx, *operationID)
}

// GetArtifactsByOperation retrieves all artifacts for a specific operation.
func (c *Client) GetArtifactsByOperation(ctx context.Context, operationID int) ([]*types.Artifact, error) {
	if err := c.EnsureAuthenticated(ctx); err != nil {
		return nil, err
	}

	if operationID == 0 {
		return nil, WrapError("GetArtifactsByOperation", ErrInvalidInput, "operation ID is required")
	}

	var query struct {
		TaskArtifact []struct {
			ID           int       `graphql:"id"`
			Artifact     string    `graphql:"artifact_text"`
			BaseArtifact string    `graphql:"base_artifact"`
			Host         string    `graphql:"host"`
			OperationID  int       `graphql:"operation_id"`
			TaskID       *int      `graphql:"task_id"`
			Timestamp    time.Time `graphql:"timestamp"`
		} `graphql:"taskartifact(where: {operation_id: {_eq: $operation_id}}, order_by: {timestamp: desc})"`
	}

	variables := map[string]interface{}{
		"operation_id": operationID,
	}

	err := c.executeQuery(ctx, &query, variables)
	if err != nil {
		return nil, WrapError("GetArtifactsByOperation", err, "failed to query artifacts")
	}

	artifacts := make([]*types.Artifact, len(query.TaskArtifact))
	for i, a := range query.TaskArtifact {
		artifacts[i] = &types.Artifact{
			ID:           a.ID,
			Artifact:     a.Artifact,
			BaseArtifact: a.BaseArtifact,
			Host:         a.Host,
			OperationID:  a.OperationID,
			TaskID:       a.TaskID,
			Timestamp:    a.Timestamp,
		}
	}

	return artifacts, nil
}

// CreateArtifact creates a new artifact (IOC) entry.
func (c *Client) CreateArtifact(ctx context.Context, req *types.CreateArtifactRequest) (*types.Artifact, error) {
	if err := c.EnsureAuthenticated(ctx); err != nil {
		return nil, err
	}

	if req == nil || req.Artifact == "" {
		return nil, WrapError("CreateArtifact", ErrInvalidInput, "artifact value is required")
	}

	if req.BaseArtifact == nil || *req.BaseArtifact == "" {
		return nil, WrapError("CreateArtifact", ErrInvalidInput, "base artifact is required")
	}

	// Use current operation
	operationID := c.GetCurrentOperation()
	if operationID == nil {
		return nil, WrapError("CreateArtifact", ErrInvalidInput, "no current operation set")
	}

	var mutation struct {
		CreateArtifact struct {
			Status     string `graphql:"status"`
			Error      string `graphql:"error"`
			ArtifactID int    `graphql:"id"`
		} `graphql:"createArtifact(artifact: $artifact, base_artifact: $base_artifact, host: $host, task_id: $task_id)"`
	}

	baseArtifact := *req.BaseArtifact

	host := ""
	if req.Host != nil {
		host = *req.Host
	}

	var taskID *int
	if req.TaskID != nil {
		taskID = req.TaskID
	}

	variables := map[string]interface{}{
		"artifact":      req.Artifact,
		"base_artifact": baseArtifact,
		"host":          host,
		"task_id":       taskID,
	}

	err := c.executeMutation(ctx, &mutation, variables)
	if err != nil {
		return nil, WrapError("CreateArtifact", err, "failed to create artifact")
	}

	if mutation.CreateArtifact.Status != "success" {
		return nil, WrapError("CreateArtifact", ErrOperationFailed, mutation.CreateArtifact.Error)
	}

	// Fetch the created artifact
	return c.GetArtifactByID(ctx, mutation.CreateArtifact.ArtifactID)
}

// GetArtifactByID retrieves a specific artifact by ID.
func (c *Client) GetArtifactByID(ctx context.Context, artifactID int) (*types.Artifact, error) {
	if err := c.EnsureAuthenticated(ctx); err != nil {
		return nil, err
	}

	if artifactID == 0 {
		return nil, WrapError("GetArtifactByID", ErrInvalidInput, "artifact ID is required")
	}

	var query struct {
		TaskArtifact []struct {
			ID           int       `graphql:"id"`
			Artifact     string    `graphql:"artifact_text"`
			BaseArtifact string    `graphql:"base_artifact"`
			Host         string    `graphql:"host"`
			OperationID  int       `graphql:"operation_id"`
			TaskID       *int      `graphql:"task_id"`
			Timestamp    time.Time `graphql:"timestamp"`
		} `graphql:"taskartifact(where: {id: {_eq: $artifact_id}})"`
	}

	variables := map[string]interface{}{
		"artifact_id": artifactID,
	}

	err := c.executeQuery(ctx, &query, variables)
	if err != nil {
		return nil, WrapError("GetArtifactByID", err, "failed to query artifact")
	}

	if len(query.TaskArtifact) == 0 {
		return nil, WrapError("GetArtifactByID", ErrNotFound, "artifact not found")
	}

	a := query.TaskArtifact[0]
	return &types.Artifact{
		ID:           a.ID,
		Artifact:     a.Artifact,
		BaseArtifact: a.BaseArtifact,
		Host:         a.Host,
		OperationID:  a.OperationID,
		TaskID:       a.TaskID,
		Timestamp:    a.Timestamp,
	}, nil
}

// UpdateArtifact updates an artifact's properties.
func (c *Client) UpdateArtifact(ctx context.Context, req *types.UpdateArtifactRequest) (*types.Artifact, error) {
	if err := c.EnsureAuthenticated(ctx); err != nil {
		return nil, err
	}

	if req == nil || req.ID == 0 {
		return nil, WrapError("UpdateArtifact", ErrInvalidInput, "artifact ID is required")
	}

	// Build update fields
	updates := make(map[string]interface{})
	if req.Host != nil {
		updates["host"] = *req.Host
	}

	if len(updates) == 0 {
		return nil, WrapError("UpdateArtifact", ErrInvalidInput, "no fields to update")
	}

	var mutation struct {
		UpdateTaskArtifact struct {
			Affected int `graphql:"affected_rows"`
		} `graphql:"update_taskartifact(where: {id: {_eq: $artifact_id}}, _set: $updates)"`
	}

	variables := map[string]interface{}{
		"artifact_id": req.ID,
		"updates":     updates,
	}

	err := c.executeMutation(ctx, &mutation, variables)
	if err != nil {
		return nil, WrapError("UpdateArtifact", err, "failed to update artifact")
	}

	if mutation.UpdateTaskArtifact.Affected == 0 {
		return nil, WrapError("UpdateArtifact", ErrNotFound, "artifact not found or not updated")
	}

	// Fetch the updated artifact
	return c.GetArtifactByID(ctx, req.ID)
}

// DeleteArtifact deletes an artifact by ID.
// Note: Mythic's taskartifact table doesn't support soft delete.
func (c *Client) DeleteArtifact(ctx context.Context, artifactID int) error {
	if err := c.EnsureAuthenticated(ctx); err != nil {
		return err
	}

	if artifactID == 0 {
		return WrapError("DeleteArtifact", ErrInvalidInput, "artifact ID is required")
	}

	var mutation struct {
		DeleteArtifact struct {
			Status string `graphql:"status"`
			Error  string `graphql:"error"`
		} `graphql:"deleteArtifact(artifact_id: $artifact_id)"`
	}

	variables := map[string]interface{}{
		"artifact_id": artifactID,
	}

	err := c.executeMutation(ctx, &mutation, variables)
	if err != nil {
		return WrapError("DeleteArtifact", err, "failed to delete artifact")
	}

	if mutation.DeleteArtifact.Status != "success" {
		return WrapError("DeleteArtifact", ErrOperationFailed, mutation.DeleteArtifact.Error)
	}

	return nil
}

// GetArtifactsByHost retrieves artifacts for a specific host.
func (c *Client) GetArtifactsByHost(ctx context.Context, host string) ([]*types.Artifact, error) {
	if err := c.EnsureAuthenticated(ctx); err != nil {
		return nil, err
	}

	if host == "" {
		return nil, WrapError("GetArtifactsByHost", ErrInvalidInput, "host is required")
	}

	// Use current operation
	operationID := c.GetCurrentOperation()
	if operationID == nil {
		return nil, WrapError("GetArtifactsByHost", ErrInvalidInput, "no current operation set")
	}

	var query struct {
		TaskArtifact []struct {
			ID           int       `graphql:"id"`
			Artifact     string    `graphql:"artifact_text"`
			BaseArtifact string    `graphql:"base_artifact"`
			Host         string    `graphql:"host"`
			OperationID  int       `graphql:"operation_id"`
			TaskID       *int      `graphql:"task_id"`
			Timestamp    time.Time `graphql:"timestamp"`
		} `graphql:"taskartifact(where: {operation_id: {_eq: $operation_id}, host: {_eq: $host}}, order_by: {timestamp: desc})"`
	}

	variables := map[string]interface{}{
		"operation_id": *operationID,
		"host":         host,
	}

	err := c.executeQuery(ctx, &query, variables)
	if err != nil {
		return nil, WrapError("GetArtifactsByHost", err, "failed to query artifacts")
	}

	artifacts := make([]*types.Artifact, len(query.TaskArtifact))
	for i, a := range query.TaskArtifact {
		artifacts[i] = &types.Artifact{
			ID:           a.ID,
			Artifact:     a.Artifact,
			BaseArtifact: a.BaseArtifact,
			Host:         a.Host,
			OperationID:  a.OperationID,
			TaskID:       a.TaskID,
			Timestamp:    a.Timestamp,
		}
	}

	return artifacts, nil
}

// GetArtifactsByType retrieves artifacts of a specific type.
// Note: Mythic's taskartifact table doesn't have a type field, so this function
// simply returns all artifacts for the current operation.
func (c *Client) GetArtifactsByType(ctx context.Context, artifactType string) ([]*types.Artifact, error) {
	// Since taskartifact doesn't have a type field, just return all artifacts
	return c.GetArtifacts(ctx)
}
