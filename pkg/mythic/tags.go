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
			ID          int       `graphql:"id"`
			Name        string    `graphql:"name"`
			Description string    `graphql:"description"`
			Color       string    `graphql:"color"`
			OperationID int       `graphql:"operation_id"`
			Deleted     bool      `graphql:"deleted"`
			Timestamp   time.Time `graphql:"timestamp"`
		} `graphql:"tagtype(where: {operation_id: {_eq: $operation_id}, deleted: {_eq: false}}, order_by: {name: asc})"`
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
			Deleted:     tt.Deleted,
			Timestamp:   tt.Timestamp,
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
			ID          int       `graphql:"id"`
			Name        string    `graphql:"name"`
			Description string    `graphql:"description"`
			Color       string    `graphql:"color"`
			OperationID int       `graphql:"operation_id"`
			Deleted     bool      `graphql:"deleted"`
			Timestamp   time.Time `graphql:"timestamp"`
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
		Deleted:     tt.Deleted,
		Timestamp:   tt.Timestamp,
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

	var mutation struct {
		CreateTagtype struct {
			Status    string `graphql:"status"`
			Error     string `graphql:"error"`
			TagtypeID int    `graphql:"id"`
		} `graphql:"createTagtype(name: $name, description: $description, color: $color, operation_id: $operation_id)"`
	}

	description := ""
	if req.Description != nil {
		description = *req.Description
	}

	color := ""
	if req.Color != nil {
		color = *req.Color
	}

	variables := map[string]interface{}{
		"name":         req.Name,
		"description":  description,
		"color":        color,
		"operation_id": *operationID,
	}

	err := c.executeMutation(ctx, &mutation, variables)
	if err != nil {
		return nil, WrapError("CreateTagType", err, "failed to create tag type")
	}

	if mutation.CreateTagtype.Status != "success" {
		return nil, WrapError("CreateTagType", ErrOperationFailed, mutation.CreateTagtype.Error)
	}

	// Fetch the created tag type
	return c.GetTagTypeByID(ctx, mutation.CreateTagtype.TagtypeID)
}

// UpdateTagType updates a tag type's properties.
func (c *Client) UpdateTagType(ctx context.Context, req *types.UpdateTagTypeRequest) (*types.TagType, error) {
	if err := c.EnsureAuthenticated(ctx); err != nil {
		return nil, err
	}

	if req == nil || req.ID == 0 {
		return nil, WrapError("UpdateTagType", ErrInvalidInput, "tag type ID is required")
	}

	// Build update fields
	updates := make(map[string]interface{})
	if req.Name != nil {
		updates["name"] = *req.Name
	}
	if req.Description != nil {
		updates["description"] = *req.Description
	}
	if req.Color != nil {
		updates["color"] = *req.Color
	}
	if req.Deleted != nil {
		updates["deleted"] = *req.Deleted
	}

	if len(updates) == 0 {
		return nil, WrapError("UpdateTagType", ErrInvalidInput, "no fields to update")
	}

	var mutation struct {
		UpdateTagtype struct {
			Affected int `graphql:"affected_rows"`
		} `graphql:"update_tagtype(where: {id: {_eq: $tagtype_id}}, _set: $updates)"`
	}

	variables := map[string]interface{}{
		"tagtype_id": req.ID,
		"updates":    updates,
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

	var mutation struct {
		CreateTag struct {
			Status string `graphql:"status"`
			Error  string `graphql:"error"`
			TagID  int    `graphql:"id"`
		} `graphql:"createTag(tagtype_id: $tagtype_id, source: $source, object_id: $object_id)"`
	}

	variables := map[string]interface{}{
		"tagtype_id": req.TagTypeID,
		"source":     req.SourceType,
		"object_id":  req.SourceID,
	}

	err := c.executeMutation(ctx, &mutation, variables)
	if err != nil {
		return nil, WrapError("CreateTag", err, "failed to create tag")
	}

	if mutation.CreateTag.Status != "success" {
		return nil, WrapError("CreateTag", ErrOperationFailed, mutation.CreateTag.Error)
	}

	// Fetch the created tag
	return c.GetTagByID(ctx, mutation.CreateTag.TagID)
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
