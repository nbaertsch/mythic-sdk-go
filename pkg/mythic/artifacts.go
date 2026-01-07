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
		Artifact []struct {
			ID           int       `graphql:"id"`
			Artifact     string    `graphql:"artifact"`
			BaseArtifact string    `graphql:"base_artifact"`
			Host         string    `graphql:"host"`
			ArtifactType string    `graphql:"type"`
			OperationID  int       `graphql:"operation_id"`
			OperatorID   int       `graphql:"operator_id"`
			TaskID       *int      `graphql:"task_id"`
			Timestamp    time.Time `graphql:"timestamp"`
			Deleted      bool      `graphql:"deleted"`
			Metadata     string    `graphql:"metadata"`
		} `graphql:"artifact(where: {operation_id: {_eq: $operation_id}, deleted: {_eq: false}}, order_by: {timestamp: desc})"`
	}

	variables := map[string]interface{}{
		"operation_id": operationID,
	}

	err := c.executeQuery(ctx, &query, variables)
	if err != nil {
		return nil, WrapError("GetArtifactsByOperation", err, "failed to query artifacts")
	}

	artifacts := make([]*types.Artifact, len(query.Artifact))
	for i, a := range query.Artifact {
		artifacts[i] = &types.Artifact{
			ID:           a.ID,
			Artifact:     a.Artifact,
			BaseArtifact: a.BaseArtifact,
			Host:         a.Host,
			ArtifactType: a.ArtifactType,
			OperationID:  a.OperationID,
			OperatorID:   a.OperatorID,
			TaskID:       a.TaskID,
			Timestamp:    a.Timestamp,
			Deleted:      a.Deleted,
			Metadata:     a.Metadata,
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
		} `graphql:"createArtifact(artifact: $artifact, base_artifact: $base_artifact, host: $host, type: $type, operation_id: $operation_id, task_id: $task_id, metadata: $metadata)"`
	}

	baseArtifact := ""
	if req.BaseArtifact != nil {
		baseArtifact = *req.BaseArtifact
	}

	host := ""
	if req.Host != nil {
		host = *req.Host
	}

	artifactType := types.ArtifactTypeOther
	if req.ArtifactType != nil {
		artifactType = *req.ArtifactType
	}

	var taskID *int
	if req.TaskID != nil {
		taskID = req.TaskID
	}

	metadata := ""
	if req.Metadata != nil {
		metadata = *req.Metadata
	}

	variables := map[string]interface{}{
		"artifact":      req.Artifact,
		"base_artifact": baseArtifact,
		"host":          host,
		"type":          artifactType,
		"operation_id":  *operationID,
		"task_id":       taskID,
		"metadata":      metadata,
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
		Artifact []struct {
			ID           int       `graphql:"id"`
			Artifact     string    `graphql:"artifact"`
			BaseArtifact string    `graphql:"base_artifact"`
			Host         string    `graphql:"host"`
			ArtifactType string    `graphql:"type"`
			OperationID  int       `graphql:"operation_id"`
			OperatorID   int       `graphql:"operator_id"`
			TaskID       *int      `graphql:"task_id"`
			Timestamp    time.Time `graphql:"timestamp"`
			Deleted      bool      `graphql:"deleted"`
			Metadata     string    `graphql:"metadata"`
		} `graphql:"artifact(where: {id: {_eq: $artifact_id}})"`
	}

	variables := map[string]interface{}{
		"artifact_id": artifactID,
	}

	err := c.executeQuery(ctx, &query, variables)
	if err != nil {
		return nil, WrapError("GetArtifactByID", err, "failed to query artifact")
	}

	if len(query.Artifact) == 0 {
		return nil, WrapError("GetArtifactByID", ErrNotFound, "artifact not found")
	}

	a := query.Artifact[0]
	return &types.Artifact{
		ID:           a.ID,
		Artifact:     a.Artifact,
		BaseArtifact: a.BaseArtifact,
		Host:         a.Host,
		ArtifactType: a.ArtifactType,
		OperationID:  a.OperationID,
		OperatorID:   a.OperatorID,
		TaskID:       a.TaskID,
		Timestamp:    a.Timestamp,
		Deleted:      a.Deleted,
		Metadata:     a.Metadata,
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
	if req.Deleted != nil {
		updates["deleted"] = *req.Deleted
	}
	if req.Metadata != nil {
		updates["metadata"] = *req.Metadata
	}

	if len(updates) == 0 {
		return nil, WrapError("UpdateArtifact", ErrInvalidInput, "no fields to update")
	}

	var mutation struct {
		UpdateArtifact struct {
			Affected int `graphql:"affected_rows"`
		} `graphql:"update_artifact(where: {id: {_eq: $artifact_id}}, _set: $updates)"`
	}

	variables := map[string]interface{}{
		"artifact_id": req.ID,
		"updates":     updates,
	}

	err := c.executeMutation(ctx, &mutation, variables)
	if err != nil {
		return nil, WrapError("UpdateArtifact", err, "failed to update artifact")
	}

	if mutation.UpdateArtifact.Affected == 0 {
		return nil, WrapError("UpdateArtifact", ErrNotFound, "artifact not found or not updated")
	}

	// Fetch the updated artifact
	return c.GetArtifactByID(ctx, req.ID)
}

// DeleteArtifact marks an artifact as deleted (soft delete).
func (c *Client) DeleteArtifact(ctx context.Context, artifactID int) error {
	if err := c.EnsureAuthenticated(ctx); err != nil {
		return err
	}

	if artifactID == 0 {
		return WrapError("DeleteArtifact", ErrInvalidInput, "artifact ID is required")
	}

	deleted := true
	_, err := c.UpdateArtifact(ctx, &types.UpdateArtifactRequest{
		ID:      artifactID,
		Deleted: &deleted,
	})
	return err
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
		Artifact []struct {
			ID           int       `graphql:"id"`
			Artifact     string    `graphql:"artifact"`
			BaseArtifact string    `graphql:"base_artifact"`
			Host         string    `graphql:"host"`
			ArtifactType string    `graphql:"type"`
			OperationID  int       `graphql:"operation_id"`
			OperatorID   int       `graphql:"operator_id"`
			TaskID       *int      `graphql:"task_id"`
			Timestamp    time.Time `graphql:"timestamp"`
			Deleted      bool      `graphql:"deleted"`
			Metadata     string    `graphql:"metadata"`
		} `graphql:"artifact(where: {operation_id: {_eq: $operation_id}, host: {_eq: $host}, deleted: {_eq: false}}, order_by: {timestamp: desc})"`
	}

	variables := map[string]interface{}{
		"operation_id": *operationID,
		"host":         host,
	}

	err := c.executeQuery(ctx, &query, variables)
	if err != nil {
		return nil, WrapError("GetArtifactsByHost", err, "failed to query artifacts")
	}

	artifacts := make([]*types.Artifact, len(query.Artifact))
	for i, a := range query.Artifact {
		artifacts[i] = &types.Artifact{
			ID:           a.ID,
			Artifact:     a.Artifact,
			BaseArtifact: a.BaseArtifact,
			Host:         a.Host,
			ArtifactType: a.ArtifactType,
			OperationID:  a.OperationID,
			OperatorID:   a.OperatorID,
			TaskID:       a.TaskID,
			Timestamp:    a.Timestamp,
			Deleted:      a.Deleted,
			Metadata:     a.Metadata,
		}
	}

	return artifacts, nil
}

// GetArtifactsByType retrieves artifacts of a specific type.
func (c *Client) GetArtifactsByType(ctx context.Context, artifactType string) ([]*types.Artifact, error) {
	if err := c.EnsureAuthenticated(ctx); err != nil {
		return nil, err
	}

	if artifactType == "" {
		return nil, WrapError("GetArtifactsByType", ErrInvalidInput, "artifact type is required")
	}

	// Use current operation
	operationID := c.GetCurrentOperation()
	if operationID == nil {
		return nil, WrapError("GetArtifactsByType", ErrInvalidInput, "no current operation set")
	}

	var query struct {
		Artifact []struct {
			ID           int       `graphql:"id"`
			Artifact     string    `graphql:"artifact"`
			BaseArtifact string    `graphql:"base_artifact"`
			Host         string    `graphql:"host"`
			ArtifactType string    `graphql:"type"`
			OperationID  int       `graphql:"operation_id"`
			OperatorID   int       `graphql:"operator_id"`
			TaskID       *int      `graphql:"task_id"`
			Timestamp    time.Time `graphql:"timestamp"`
			Deleted      bool      `graphql:"deleted"`
			Metadata     string    `graphql:"metadata"`
		} `graphql:"artifact(where: {operation_id: {_eq: $operation_id}, type: {_eq: $type}, deleted: {_eq: false}}, order_by: {timestamp: desc})"`
	}

	variables := map[string]interface{}{
		"operation_id": *operationID,
		"type":         artifactType,
	}

	err := c.executeQuery(ctx, &query, variables)
	if err != nil {
		return nil, WrapError("GetArtifactsByType", err, "failed to query artifacts")
	}

	artifacts := make([]*types.Artifact, len(query.Artifact))
	for i, a := range query.Artifact {
		artifacts[i] = &types.Artifact{
			ID:           a.ID,
			Artifact:     a.Artifact,
			BaseArtifact: a.BaseArtifact,
			Host:         a.Host,
			ArtifactType: a.ArtifactType,
			OperationID:  a.OperationID,
			OperatorID:   a.OperatorID,
			TaskID:       a.TaskID,
			Timestamp:    a.Timestamp,
			Deleted:      a.Deleted,
			Metadata:     a.Metadata,
		}
	}

	return artifacts, nil
}
