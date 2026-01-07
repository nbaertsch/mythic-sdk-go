package types

import (
	"fmt"
	"time"
)

// Artifact represents an indicator of compromise (IOC) tracked across an operation.
type Artifact struct {
	ID           int        `json:"id"`
	Artifact     string     `json:"artifact"`
	BaseArtifact string     `json:"base_artifact"`
	Host         string     `json:"host"`
	ArtifactType string     `json:"artifact_type"`
	OperationID  int        `json:"operation_id"`
	OperatorID   int        `json:"operator_id"`
	TaskID       *int       `json:"task_id,omitempty"`
	Timestamp    time.Time  `json:"timestamp"`
	Deleted      bool       `json:"deleted"`
	Metadata     string     `json:"metadata,omitempty"`
	Operation    *Operation `json:"operation,omitempty"`
	Operator     *Operator  `json:"operator,omitempty"`
}

// String returns a string representation of an Artifact.
func (a *Artifact) String() string {
	if a.Host != "" && a.Artifact != "" {
		return fmt.Sprintf("%s on %s (%s)", a.Artifact, a.Host, a.ArtifactType)
	}
	if a.Artifact != "" {
		return fmt.Sprintf("%s (%s)", a.Artifact, a.ArtifactType)
	}
	return fmt.Sprintf("Artifact %d", a.ID)
}

// IsDeleted returns true if the artifact is marked as deleted.
func (a *Artifact) IsDeleted() bool {
	return a.Deleted
}

// HasTask returns true if the artifact is linked to a task.
func (a *Artifact) HasTask() bool {
	return a.TaskID != nil && *a.TaskID != 0
}

// CreateArtifactRequest represents a request to create a new artifact.
type CreateArtifactRequest struct {
	Artifact     string  `json:"artifact"`
	BaseArtifact *string `json:"base_artifact,omitempty"`
	Host         *string `json:"host,omitempty"`
	ArtifactType *string `json:"artifact_type,omitempty"`
	TaskID       *int    `json:"task_id,omitempty"`
	Metadata     *string `json:"metadata,omitempty"`
}

// UpdateArtifactRequest represents a request to update an artifact.
type UpdateArtifactRequest struct {
	ID       int
	Host     *string `json:"host,omitempty"`
	Deleted  *bool   `json:"deleted,omitempty"`
	Metadata *string `json:"metadata,omitempty"`
}

// Common artifact types
const (
	ArtifactTypeFile      = "file"
	ArtifactTypeRegistry  = "registry"
	ArtifactTypeProcess   = "process"
	ArtifactTypeNetwork   = "network"
	ArtifactTypeUser      = "user"
	ArtifactTypeService   = "service"
	ArtifactTypeScheduled = "scheduled_task"
	ArtifactTypeWMI       = "wmi"
	ArtifactTypeOther     = "other"
)
