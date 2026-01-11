package mythic

import (
	"context"
	"time"

	"github.com/nbaertsch/mythic-sdk-go/pkg/mythic/types"
)

// GetTagTypes retrieves all tag types for the current operation.
func (c *Client) GetTagTypes(ctx context.Context) ([]*types.TagType, error) {
	if err := c.EnsureAuthenticated(ctx); err != nil {
		return nil, err
	}

	// Use current operation if set
	operationID := c.GetCurrentOperation()
	if operationID == nil {
		return nil, WrapError("GetTagTypes", ErrInvalidInput, "no current operation set")
	}

	return c.GetTagTypesByOperation(ctx, *operationID)
}

// GetTagTypesByOperation retrieves all tag types for a specific operation.
func (c *Client) GetTagTypesByOperation(ctx context.Context, operationID int) ([]*types.TagType, error) {
	if err := c.EnsureAuthenticated(ctx); err != nil {
		return nil, err
	}

	if operationID == 0 {
		return nil, WrapError("GetTagTypesByOperation", ErrInvalidInput, "operation ID is required")
	}

	var query struct {
		TagType []struct {
			ID          int    `graphql:"id"`
			Name        string `graphql:"name"`
			Description string `graphql:"description"`
			Color       string `graphql:"color"`
			OperationID int    `graphql:"operation_id"`
		} `graphql:"tagtype(where: {operation_id: {_eq: $operation_id}}, order_by: {name: asc})"`
	}

	variables := map[string]interface{}{
		"operation_id": operationID,
	}

	err := c.executeQuery(ctx, &query, variables)
	if err != nil {
		return nil, WrapError("GetTagTypesByOperation", err, "failed to query tag types")
	}

	tagTypes := make([]*types.TagType, len(query.TagType))
	for i, tt := range query.TagType {
		tagTypes[i] = &types.TagType{
			ID:          tt.ID,
			Name:        tt.Name,
			Description: tt.Description,
			Color:       tt.Color,
			OperationID: tt.OperationID,
		}
	}

	return tagTypes, nil
}

// GetTagTypeByID retrieves a specific tag type by ID.
func (c *Client) GetTagTypeByID(ctx context.Context, tagTypeID int) (*types.TagType, error) {
	if err := c.EnsureAuthenticated(ctx); err != nil {
		return nil, err
	}

	if tagTypeID == 0 {
		return nil, WrapError("GetTagTypeByID", ErrInvalidInput, "tag type ID is required")
	}

	var query struct {
		TagType []struct {
			ID          int    `graphql:"id"`
			Name        string `graphql:"name"`
			Description string `graphql:"description"`
			Color       string `graphql:"color"`
			OperationID int    `graphql:"operation_id"`
		} `graphql:"tagtype(where: {id: {_eq: $tagtype_id}})"`
	}

	variables := map[string]interface{}{
		"tagtype_id": tagTypeID,
	}

	err := c.executeQuery(ctx, &query, variables)
	if err != nil {
		return nil, WrapError("GetTagTypeByID", err, "failed to query tag type")
	}

	if len(query.TagType) == 0 {
		return nil, WrapError("GetTagTypeByID", ErrNotFound, "tag type not found")
	}

	tt := query.TagType[0]
	return &types.TagType{
		ID:          tt.ID,
		Name:        tt.Name,
		Description: tt.Description,
		Color:       tt.Color,
		OperationID: tt.OperationID,
	}, nil
}

// CreateTagType creates a new tag type.
func (c *Client) CreateTagType(ctx context.Context, req *types.CreateTagTypeRequest) (*types.TagType, error) {
	if err := c.EnsureAuthenticated(ctx); err != nil {
		return nil, err
	}

	if req == nil || req.Name == "" {
		return nil, WrapError("CreateTagType", ErrInvalidInput, "tag type name is required")
	}

	// Use current operation
	operationID := c.GetCurrentOperation()
	if operationID == nil {
		return nil, WrapError("CreateTagType", ErrInvalidInput, "no current operation set")
	}

	// Mythic does not provide a createTagType mutation in GraphQL
	// Tag types must be managed through the admin interface
	return nil, WrapError("CreateTagType", ErrOperationFailed, "tag type creation not supported via GraphQL API")
}

// UpdateTagType updates a tag type's properties.
func (c *Client) UpdateTagType(ctx context.Context, req *types.UpdateTagTypeRequest) (*types.TagType, error) {
	if err := c.EnsureAuthenticated(ctx); err != nil {
		return nil, err
	}

	if req == nil || req.ID == 0 {
		return nil, WrapError("UpdateTagType", ErrInvalidInput, "tag type ID is required")
	}

	// Check if there are any fields to update
	hasUpdates := req.Name != nil || req.Description != nil || req.Color != nil

	if !hasUpdates {
		return nil, WrapError("UpdateTagType", ErrInvalidInput, "no fields to update")
	}

	// Simplified: only support updating name for now
	// Note: Full multi-field updates require a different GraphQL approach
	if req.Name == nil {
		return nil, WrapError("UpdateTagType", ErrInvalidInput, "currently only name field updates are supported")
	}

	var mutation struct {
		UpdateTagtype struct {
			Affected int `graphql:"affected_rows"`
		} `graphql:"update_tagtype(where: {id: {_eq: $tagtype_id}}, _set: {name: $name})"`
	}

	variables := map[string]interface{}{
		"tagtype_id": req.ID,
		"name":       *req.Name,
	}

	err := c.executeMutation(ctx, &mutation, variables)
	if err != nil {
		return nil, WrapError("UpdateTagType", err, "failed to update tag type")
	}

	if mutation.UpdateTagtype.Affected == 0 {
		return nil, WrapError("UpdateTagType", ErrNotFound, "tag type not found or not updated")
	}

	// Fetch the updated tag type
	return c.GetTagTypeByID(ctx, req.ID)
}

// DeleteTagType marks a tag type as deleted (soft delete).
func (c *Client) DeleteTagType(ctx context.Context, tagTypeID int) error {
	if err := c.EnsureAuthenticated(ctx); err != nil {
		return err
	}

	if tagTypeID == 0 {
		return WrapError("DeleteTagType", ErrInvalidInput, "tag type ID is required")
	}

	var mutation struct {
		DeleteTagtype struct {
			Status string `graphql:"status"`
			Error  string `graphql:"error"`
		} `graphql:"deleteTagtype(id: $id)"`
	}

	variables := map[string]interface{}{
		"id": tagTypeID,
	}

	err := c.executeMutation(ctx, &mutation, variables)
	if err != nil {
		return WrapError("DeleteTagType", err, "failed to delete tag type")
	}

	if mutation.DeleteTagtype.Status != "success" {
		return WrapError("DeleteTagType", ErrOperationFailed, mutation.DeleteTagtype.Error)
	}

	return nil
}

// CreateTag creates a tag on an object using the REST API webhook.
func (c *Client) CreateTag(ctx context.Context, req *types.CreateTagRequest) (*types.Tag, error) {
	if err := c.EnsureAuthenticated(ctx); err != nil {
		return nil, err
	}

	if req == nil || req.TagTypeID == 0 || req.SourceType == "" || req.SourceID == 0 {
		return nil, WrapError("CreateTag", ErrInvalidInput, "tag type ID, source type, and source ID are required")
	}

	// Build REST API request using Mythic's webhook format
	input := map[string]interface{}{
		"tagtype_id": req.TagTypeID,
		"source":     req.SourceType,
		"data":       map[string]interface{}{},
		"url":        "",
	}

	// Add the appropriate object ID parameter based on source type
	switch req.SourceType {
	case types.TagSourceArtifact:
		input["taskartifact_id"] = req.SourceID
	case types.TagSourceTask:
		input["task_id"] = req.SourceID
	case types.TagSourceCallback:
		input["callback_id"] = req.SourceID
	case types.TagSourceFile:
		input["filemeta_id"] = req.SourceID
	case types.TagSourcePayload:
		input["payload_id"] = req.SourceID
	case types.TagSourceProcess:
		input["mythictree_id"] = req.SourceID
	case types.TagSourceKeylog:
		input["keylog_id"] = req.SourceID
	case "credential":
		input["credential_id"] = req.SourceID
	case "response":
		input["response_id"] = req.SourceID
	default:
		return nil, WrapError("CreateTag", ErrInvalidInput, "unsupported source type: "+req.SourceType)
	}

	// Call REST API webhook
	requestData := map[string]interface{}{
		"input": input,
	}

	var response struct {
		Status string `json:"status"`
		Error  string `json:"error"`
		ID     int    `json:"id"`
	}

	err := c.executeRESTWebhook(ctx, "api/v1.4/tag_create_webhook", requestData, &response)
	if err != nil {
		return nil, WrapError("CreateTag", err, "failed to execute webhook")
	}

	if response.Status != "success" {
		return nil, WrapError("CreateTag", ErrOperationFailed, response.Error)
	}

	// Fetch the created tag
	return c.GetTagByID(ctx, response.ID)
}

// GetTagByID retrieves a specific tag by ID.
func (c *Client) GetTagByID(ctx context.Context, tagID int) (*types.Tag, error) {
	if err := c.EnsureAuthenticated(ctx); err != nil {
		return nil, err
	}

	if tagID == 0 {
		return nil, WrapError("GetTagByID", ErrInvalidInput, "tag ID is required")
	}

	var query struct {
		Tag []struct {
			ID             int    `graphql:"id"`
			TagTypeID      int    `graphql:"tagtype_id"`
			SourceType     string `graphql:"source"`
			OperationID    int    `graphql:"operation_id"`
			TaskArtifactID *int   `graphql:"taskartifact_id"`
			TaskID         *int   `graphql:"task_id"`
			CallbackID     *int   `graphql:"callback_id"`
			FilemetaID     *int   `graphql:"filemeta_id"`
			PayloadID      *int   `graphql:"payload_id"`
			MythictreeID   *int   `graphql:"mythictree_id"`
			KeylogID       *int   `graphql:"keylog_id"`
			CredentialID   *int   `graphql:"credential_id"`
			ResponseID     *int   `graphql:"response_id"`
		} `graphql:"tag(where: {id: {_eq: $tag_id}})"`
	}

	variables := map[string]interface{}{
		"tag_id": tagID,
	}

	err := c.executeQuery(ctx, &query, variables)
	if err != nil {
		return nil, WrapError("GetTagByID", err, "failed to query tag")
	}

	if len(query.Tag) == 0 {
		return nil, WrapError("GetTagByID", ErrNotFound, "tag not found")
	}

	t := query.Tag[0]

	// Determine which source ID is populated
	var sourceID int
	if t.TaskArtifactID != nil {
		sourceID = *t.TaskArtifactID
	} else if t.TaskID != nil {
		sourceID = *t.TaskID
	} else if t.CallbackID != nil {
		sourceID = *t.CallbackID
	} else if t.FilemetaID != nil {
		sourceID = *t.FilemetaID
	} else if t.PayloadID != nil {
		sourceID = *t.PayloadID
	} else if t.MythictreeID != nil {
		sourceID = *t.MythictreeID
	} else if t.KeylogID != nil {
		sourceID = *t.KeylogID
	} else if t.CredentialID != nil {
		sourceID = *t.CredentialID
	} else if t.ResponseID != nil {
		sourceID = *t.ResponseID
	}

	return &types.Tag{
		ID:          t.ID,
		TagTypeID:   t.TagTypeID,
		SourceType:  t.SourceType,
		SourceID:    sourceID,
		OperationID: t.OperationID,
	}, nil
}

// GetTags retrieves all tags for a specific object.
func (c *Client) GetTags(ctx context.Context, sourceType string, sourceID int) ([]*types.Tag, error) {
	if err := c.EnsureAuthenticated(ctx); err != nil {
		return nil, err
	}

	if sourceType == "" || sourceID == 0 {
		return nil, WrapError("GetTags", ErrInvalidInput, "source type and source ID are required")
	}

	// Build where clause based on source type
	var whereField string
	switch sourceType {
	case types.TagSourceArtifact:
		whereField = "taskartifact_id"
	case types.TagSourceTask:
		whereField = "task_id"
	case types.TagSourceCallback:
		whereField = "callback_id"
	case types.TagSourceFile:
		whereField = "filemeta_id"
	case types.TagSourcePayload:
		whereField = "payload_id"
	case types.TagSourceProcess:
		whereField = "mythictree_id"
	case types.TagSourceKeylog:
		whereField = "keylog_id"
	case "credential":
		whereField = "credential_id"
	case "response":
		whereField = "response_id"
	default:
		return nil, WrapError("GetTags", ErrInvalidInput, "unsupported source type: "+sourceType)
	}

	// Query all tags with this source type and then filter by ID in code
	var query struct {
		Tag []struct {
			ID             int    `graphql:"id"`
			TagTypeID      int    `graphql:"tagtype_id"`
			SourceType     string `graphql:"source"`
			OperationID    int    `graphql:"operation_id"`
			TaskArtifactID *int   `graphql:"taskartifact_id"`
			TaskID         *int   `graphql:"task_id"`
			CallbackID     *int   `graphql:"callback_id"`
			FilemetaID     *int   `graphql:"filemeta_id"`
			PayloadID      *int   `graphql:"payload_id"`
			MythictreeID   *int   `graphql:"mythictree_id"`
			KeylogID       *int   `graphql:"keylog_id"`
			CredentialID   *int   `graphql:"credential_id"`
			ResponseID     *int   `graphql:"response_id"`
		} `graphql:"tag(where: {source: {_eq: $source}}, order_by: {id: desc})"`
	}

	variables := map[string]interface{}{
		"source": sourceType,
	}

	err := c.executeQuery(ctx, &query, variables)
	if err != nil {
		return nil, WrapError("GetTags", err, "failed to query tags")
	}

	// Filter tags by the appropriate ID field
	tags := make([]*types.Tag, 0)
	for _, t := range query.Tag {
		var tagSourceID int
		switch whereField {
		case "taskartifact_id":
			if t.TaskArtifactID != nil {
				tagSourceID = *t.TaskArtifactID
			}
		case "task_id":
			if t.TaskID != nil {
				tagSourceID = *t.TaskID
			}
		case "callback_id":
			if t.CallbackID != nil {
				tagSourceID = *t.CallbackID
			}
		case "filemeta_id":
			if t.FilemetaID != nil {
				tagSourceID = *t.FilemetaID
			}
		case "payload_id":
			if t.PayloadID != nil {
				tagSourceID = *t.PayloadID
			}
		case "mythictree_id":
			if t.MythictreeID != nil {
				tagSourceID = *t.MythictreeID
			}
		case "keylog_id":
			if t.KeylogID != nil {
				tagSourceID = *t.KeylogID
			}
		case "credential_id":
			if t.CredentialID != nil {
				tagSourceID = *t.CredentialID
			}
		case "response_id":
			if t.ResponseID != nil {
				tagSourceID = *t.ResponseID
			}
		}

		// Only include tags that match the requested source ID
		if tagSourceID == sourceID {
			tags = append(tags, &types.Tag{
				ID:          t.ID,
				TagTypeID:   t.TagTypeID,
				SourceType:  t.SourceType,
				SourceID:    tagSourceID,
				OperationID: t.OperationID,
			})
		}
	}

	return tags, nil
}

// GetTagsByOperation retrieves all tags for a specific operation.
func (c *Client) GetTagsByOperation(ctx context.Context, operationID int) ([]*types.Tag, error) {
	if err := c.EnsureAuthenticated(ctx); err != nil {
		return nil, err
	}

	if operationID == 0 {
		return nil, WrapError("GetTagsByOperation", ErrInvalidInput, "operation ID is required")
	}

	var query struct {
		Tag []struct {
			ID             int    `graphql:"id"`
			TagTypeID      int    `graphql:"tagtype_id"`
			SourceType     string `graphql:"source"`
			OperationID    int    `graphql:"operation_id"`
			TaskArtifactID *int   `graphql:"taskartifact_id"`
			TaskID         *int   `graphql:"task_id"`
			CallbackID     *int   `graphql:"callback_id"`
			FilemetaID     *int   `graphql:"filemeta_id"`
			PayloadID      *int   `graphql:"payload_id"`
			MythictreeID   *int   `graphql:"mythictree_id"`
			KeylogID       *int   `graphql:"keylog_id"`
			CredentialID   *int   `graphql:"credential_id"`
			ResponseID     *int   `graphql:"response_id"`
		} `graphql:"tag(where: {operation_id: {_eq: $operation_id}}, order_by: {id: desc})"`
	}

	variables := map[string]interface{}{
		"operation_id": operationID,
	}

	err := c.executeQuery(ctx, &query, variables)
	if err != nil {
		return nil, WrapError("GetTagsByOperation", err, "failed to query tags")
	}

	tags := make([]*types.Tag, len(query.Tag))
	for i, t := range query.Tag {
		// Determine which source ID is populated
		var sourceID int
		if t.TaskArtifactID != nil {
			sourceID = *t.TaskArtifactID
		} else if t.TaskID != nil {
			sourceID = *t.TaskID
		} else if t.CallbackID != nil {
			sourceID = *t.CallbackID
		} else if t.FilemetaID != nil {
			sourceID = *t.FilemetaID
		} else if t.PayloadID != nil {
			sourceID = *t.PayloadID
		} else if t.MythictreeID != nil {
			sourceID = *t.MythictreeID
		} else if t.KeylogID != nil {
			sourceID = *t.KeylogID
		} else if t.CredentialID != nil {
			sourceID = *t.CredentialID
		} else if t.ResponseID != nil {
			sourceID = *t.ResponseID
		}

		tags[i] = &types.Tag{
			ID:          t.ID,
			TagTypeID:   t.TagTypeID,
			SourceType:  t.SourceType,
			SourceID:    sourceID,
			OperationID: t.OperationID,
		}
	}

	return tags, nil
}

// DeleteTag removes a tag from an object.
func (c *Client) DeleteTag(ctx context.Context, tagID int) error {
	if err := c.EnsureAuthenticated(ctx); err != nil {
		return err
	}

	if tagID == 0 {
		return WrapError("DeleteTag", ErrInvalidInput, "tag ID is required")
	}

	var mutation struct {
		DeleteTag struct {
			Affected int `graphql:"affected_rows"`
		} `graphql:"delete_tag(where: {id: {_eq: $tag_id}})"`
	}

	variables := map[string]interface{}{
		"tag_id": tagID,
	}

	err := c.executeMutation(ctx, &mutation, variables)
	if err != nil {
		return WrapError("DeleteTag", err, "failed to delete tag")
	}

	if mutation.DeleteTag.Affected == 0 {
		return WrapError("DeleteTag", ErrNotFound, "tag not found")
	}

	return nil
}
