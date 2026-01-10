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

// CreateTag creates a tag on an object.
func (c *Client) CreateTag(ctx context.Context, req *types.CreateTagRequest) (*types.Tag, error) {
	if err := c.EnsureAuthenticated(ctx); err != nil {
		return nil, err
	}

	if req == nil || req.TagTypeID == 0 || req.SourceType == "" || req.SourceID == 0 {
		return nil, WrapError("CreateTag", ErrInvalidInput, "tag type ID, source type, and source ID are required")
	}

	// Mythic uses specific parameters for each object type
	// Execute the appropriate mutation based on source type
	var tagID int
	var err error

	switch req.SourceType {
	case types.TagSourceArtifact:
		var mutation struct {
			CreateTag struct {
				Status string `graphql:"status"`
				Error  string `graphql:"error"`
				TagID  int    `graphql:"id"`
			} `graphql:"createTag(tagtype_id: $tagtype_id, source: $source, url: $url, data: $data, taskartifact_id: $taskartifact_id)"`
		}
		variables := map[string]interface{}{
			"tagtype_id":      req.TagTypeID,
			"source":          req.SourceType,
			"url":             "",
			"data":            map[string]interface{}{},
			"taskartifact_id": req.SourceID,
		}
		err = c.executeMutation(ctx, &mutation, variables)
		if err == nil && mutation.CreateTag.Status == "success" {
			tagID = mutation.CreateTag.TagID
		} else if err == nil {
			err = WrapError("CreateTag", ErrOperationFailed, mutation.CreateTag.Error)
		}

	case types.TagSourceTask:
		var mutation struct {
			CreateTag struct {
				Status string `graphql:"status"`
				Error  string `graphql:"error"`
				TagID  int    `graphql:"id"`
			} `graphql:"createTag(tagtype_id: $tagtype_id, source: $source, url: $url, data: $data, task_id: $task_id)"`
		}
		variables := map[string]interface{}{
			"tagtype_id": req.TagTypeID,
			"source":     req.SourceType,
			"url":        "",
			"data":       map[string]interface{}{},
			"task_id":    req.SourceID,
		}
		err = c.executeMutation(ctx, &mutation, variables)
		if err == nil && mutation.CreateTag.Status == "success" {
			tagID = mutation.CreateTag.TagID
		} else if err == nil {
			err = WrapError("CreateTag", ErrOperationFailed, mutation.CreateTag.Error)
		}

	case types.TagSourceCallback:
		var mutation struct {
			CreateTag struct {
				Status string `graphql:"status"`
				Error  string `graphql:"error"`
				TagID  int    `graphql:"id"`
			} `graphql:"createTag(tagtype_id: $tagtype_id, source: $source, url: $url, data: $data, callback_id: $callback_id)"`
		}
		variables := map[string]interface{}{
			"tagtype_id":  req.TagTypeID,
			"source":      req.SourceType,
			"data":        map[string]interface{}{},
			"callback_id": req.SourceID,
		}
		err = c.executeMutation(ctx, &mutation, variables)
		if err == nil && mutation.CreateTag.Status == "success" {
			tagID = mutation.CreateTag.TagID
		} else if err == nil {
			err = WrapError("CreateTag", ErrOperationFailed, mutation.CreateTag.Error)
		}

	case types.TagSourceFile:
		var mutation struct {
			CreateTag struct {
				Status string `graphql:"status"`
				Error  string `graphql:"error"`
				TagID  int    `graphql:"id"`
			} `graphql:"createTag(tagtype_id: $tagtype_id, source: $source, url: $url, data: $data, filemeta_id: $filemeta_id)"`
		}
		variables := map[string]interface{}{
			"tagtype_id":  req.TagTypeID,
			"source":      req.SourceType,
			"data":        map[string]interface{}{},
			"filemeta_id": req.SourceID,
		}
		err = c.executeMutation(ctx, &mutation, variables)
		if err == nil && mutation.CreateTag.Status == "success" {
			tagID = mutation.CreateTag.TagID
		} else if err == nil {
			err = WrapError("CreateTag", ErrOperationFailed, mutation.CreateTag.Error)
		}

	case types.TagSourcePayload:
		var mutation struct {
			CreateTag struct {
				Status string `graphql:"status"`
				Error  string `graphql:"error"`
				TagID  int    `graphql:"id"`
			} `graphql:"createTag(tagtype_id: $tagtype_id, source: $source, url: $url, data: $data, payload_id: $payload_id)"`
		}
		variables := map[string]interface{}{
			"tagtype_id": req.TagTypeID,
			"source":     req.SourceType,
			"data":       map[string]interface{}{},
			"payload_id": req.SourceID,
		}
		err = c.executeMutation(ctx, &mutation, variables)
		if err == nil && mutation.CreateTag.Status == "success" {
			tagID = mutation.CreateTag.TagID
		} else if err == nil {
			err = WrapError("CreateTag", ErrOperationFailed, mutation.CreateTag.Error)
		}

	case types.TagSourceProcess:
		var mutation struct {
			CreateTag struct {
				Status string `graphql:"status"`
				Error  string `graphql:"error"`
				TagID  int    `graphql:"id"`
			} `graphql:"createTag(tagtype_id: $tagtype_id, source: $source, url: $url, data: $data, mythictree_id: $mythictree_id)"`
		}
		variables := map[string]interface{}{
			"tagtype_id":    req.TagTypeID,
			"source":        req.SourceType,
			"data":          map[string]interface{}{},
			"mythictree_id": req.SourceID,
		}
		err = c.executeMutation(ctx, &mutation, variables)
		if err == nil && mutation.CreateTag.Status == "success" {
			tagID = mutation.CreateTag.TagID
		} else if err == nil {
			err = WrapError("CreateTag", ErrOperationFailed, mutation.CreateTag.Error)
		}

	case types.TagSourceKeylog:
		var mutation struct {
			CreateTag struct {
				Status string `graphql:"status"`
				Error  string `graphql:"error"`
				TagID  int    `graphql:"id"`
			} `graphql:"createTag(tagtype_id: $tagtype_id, source: $source, url: $url, data: $data, keylog_id: $keylog_id)"`
		}
		variables := map[string]interface{}{
			"tagtype_id": req.TagTypeID,
			"source":     req.SourceType,
			"data":       map[string]interface{}{},
			"keylog_id":  req.SourceID,
		}
		err = c.executeMutation(ctx, &mutation, variables)
		if err == nil && mutation.CreateTag.Status == "success" {
			tagID = mutation.CreateTag.TagID
		} else if err == nil {
			err = WrapError("CreateTag", ErrOperationFailed, mutation.CreateTag.Error)
		}

	case "credential":
		var mutation struct {
			CreateTag struct {
				Status string `graphql:"status"`
				Error  string `graphql:"error"`
				TagID  int    `graphql:"id"`
			} `graphql:"createTag(tagtype_id: $tagtype_id, source: $source, url: $url, data: $data, credential_id: $credential_id)"`
		}
		variables := map[string]interface{}{
			"tagtype_id":    req.TagTypeID,
			"source":        req.SourceType,
			"data":          map[string]interface{}{},
			"credential_id": req.SourceID,
		}
		err = c.executeMutation(ctx, &mutation, variables)
		if err == nil && mutation.CreateTag.Status == "success" {
			tagID = mutation.CreateTag.TagID
		} else if err == nil {
			err = WrapError("CreateTag", ErrOperationFailed, mutation.CreateTag.Error)
		}

	case "response":
		var mutation struct {
			CreateTag struct {
				Status string `graphql:"status"`
				Error  string `graphql:"error"`
				TagID  int    `graphql:"id"`
			} `graphql:"createTag(tagtype_id: $tagtype_id, source: $source, url: $url, data: $data, response_id: $response_id)"`
		}
		variables := map[string]interface{}{
			"tagtype_id":  req.TagTypeID,
			"source":      req.SourceType,
			"data":        map[string]interface{}{},
			"response_id": req.SourceID,
		}
		err = c.executeMutation(ctx, &mutation, variables)
		if err == nil && mutation.CreateTag.Status == "success" {
			tagID = mutation.CreateTag.TagID
		} else if err == nil {
			err = WrapError("CreateTag", ErrOperationFailed, mutation.CreateTag.Error)
		}

	default:
		return nil, WrapError("CreateTag", ErrInvalidInput, "unsupported source type: "+req.SourceType)
	}
	if err != nil {
		return nil, err
	}

	// Fetch the created tag
	return c.GetTagByID(ctx, tagID)
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
			ID          int       `graphql:"id"`
			TagTypeID   int       `graphql:"tagtype_id"`
			SourceType  string    `graphql:"source"`
			SourceID    int       `graphql:"source_id"`
			OperatorID  int       `graphql:"operator_id"`
			OperationID int       `graphql:"operation_id"`
			Timestamp   time.Time `graphql:"timestamp"`
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
	return &types.Tag{
		ID:          t.ID,
		TagTypeID:   t.TagTypeID,
		SourceType:  t.SourceType,
		SourceID:    t.SourceID,
		OperatorID:  t.OperatorID,
		OperationID: t.OperationID,
		Timestamp:   t.Timestamp,
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

	var query struct {
		Tag []struct {
			ID          int       `graphql:"id"`
			TagTypeID   int       `graphql:"tagtype_id"`
			SourceType  string    `graphql:"source"`
			SourceID    int       `graphql:"source_id"`
			OperatorID  int       `graphql:"operator_id"`
			OperationID int       `graphql:"operation_id"`
			Timestamp   time.Time `graphql:"timestamp"`
		} `graphql:"tag(where: {source: {_eq: $source}, source_id: {_eq: $source_id}}, order_by: {timestamp: desc})"`
	}

	variables := map[string]interface{}{
		"source":    sourceType,
		"source_id": sourceID,
	}

	err := c.executeQuery(ctx, &query, variables)
	if err != nil {
		return nil, WrapError("GetTags", err, "failed to query tags")
	}

	tags := make([]*types.Tag, len(query.Tag))
	for i, t := range query.Tag {
		tags[i] = &types.Tag{
			ID:          t.ID,
			TagTypeID:   t.TagTypeID,
			SourceType:  t.SourceType,
			SourceID:    t.SourceID,
			OperatorID:  t.OperatorID,
			OperationID: t.OperationID,
			Timestamp:   t.Timestamp,
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
			ID          int       `graphql:"id"`
			TagTypeID   int       `graphql:"tagtype_id"`
			SourceType  string    `graphql:"source"`
			SourceID    int       `graphql:"source_id"`
			OperatorID  int       `graphql:"operator_id"`
			OperationID int       `graphql:"operation_id"`
			Timestamp   time.Time `graphql:"timestamp"`
		} `graphql:"tag(where: {operation_id: {_eq: $operation_id}}, order_by: {timestamp: desc})"`
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
		tags[i] = &types.Tag{
			ID:          t.ID,
			TagTypeID:   t.TagTypeID,
			SourceType:  t.SourceType,
			SourceID:    t.SourceID,
			OperatorID:  t.OperatorID,
			OperationID: t.OperationID,
			Timestamp:   t.Timestamp,
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
