package types

import (
	"fmt"
	"time"
)

// TagType represents a category/label type for tagging objects.
type TagType struct {
	ID          int        `json:"id"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Color       string     `json:"color"`
	OperationID int        `json:"operation_id"`
	Deleted     bool       `json:"deleted"`
	Timestamp   time.Time  `json:"timestamp"`
	Operation   *Operation `json:"operation,omitempty"`
}

// String returns a string representation of a TagType.
func (t *TagType) String() string {
	if t.Color != "" {
		return fmt.Sprintf("%s (%s)", t.Name, t.Color)
	}
	return t.Name
}

// IsDeleted returns true if the tag type is marked as deleted.
func (t *TagType) IsDeleted() bool {
	return t.Deleted
}

// Tag represents a tag instance applied to an object.
type Tag struct {
	ID          int        `json:"id"`
	TagTypeID   int        `json:"tagtype_id"`
	TagType     *TagType   `json:"tagtype,omitempty"`
	SourceType  string     `json:"source"`
	SourceID    int        `json:"source_id"`
	OperatorID  int        `json:"operator_id"`
	OperationID int        `json:"operation_id"`
	Timestamp   time.Time  `json:"timestamp"`
	Operation   *Operation `json:"operation,omitempty"`
	Operator    *Operator  `json:"operator,omitempty"`
}

// String returns a string representation of a Tag.
func (tag *Tag) String() string {
	if tag.TagType != nil {
		return fmt.Sprintf("%s on %s #%d", tag.TagType.Name, tag.SourceType, tag.SourceID)
	}
	return fmt.Sprintf("TagType %d on %s #%d", tag.TagTypeID, tag.SourceType, tag.SourceID)
}

// CreateTagTypeRequest represents a request to create a new tag type.
type CreateTagTypeRequest struct {
	Name        string  `json:"name"`
	Description *string `json:"description,omitempty"`
	Color       *string `json:"color,omitempty"`
}

// UpdateTagTypeRequest represents a request to update a tag type.
type UpdateTagTypeRequest struct {
	ID          int
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
	Color       *string `json:"color,omitempty"`
	Deleted     *bool   `json:"deleted,omitempty"`
}

// CreateTagRequest represents a request to create a tag on an object.
type CreateTagRequest struct {
	TagTypeID  int    `json:"tagtype_id"`
	SourceType string `json:"source"`
	SourceID   int    `json:"source_id"`
}

// Object types that can be tagged
const (
	TagSourceTask     = "task"
	TagSourceCallback = "callback"
	TagSourceFile     = "filemeta"
	TagSourcePayload  = "payload"
	TagSourceArtifact = "artifact"
	TagSourceProcess  = "process"
	TagSourceKeylog   = "keylog"
)
