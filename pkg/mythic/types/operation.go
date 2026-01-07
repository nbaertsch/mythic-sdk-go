package types

import (
	"fmt"
	"time"
)

// Operation represents a Mythic operation (engagement/campaign).
type Operation struct {
	ID                  int       `json:"id"`
	Name                string    `json:"name"`
	Complete            bool      `json:"complete"`
	Webhook             string    `json:"webhook"`
	Channel             string    `json:"channel"`
	AdminID             int       `json:"admin_id"`
	BannerText          string    `json:"banner_text"`
	BannerColor         string    `json:"banner_color"`
	DisplayName         string    `json:"display_name"`
	Icon                string    `json:"icon"`
	IconURL             string    `json:"icon_url"`
	IconEmoji           string    `json:"icon_emoji"`
	AESPSK              string    `json:"AESPSK"`
	OperationEventLogID int       `json:"operation_event_log_id"`
	Created             time.Time `json:"created"`
	Admin               *Operator `json:"admin,omitempty"`
}

// OperationOperator represents the relationship between an operation and an operator.
type OperationOperator struct {
	ID          int              `json:"id"`
	OperationID int              `json:"operation_id"`
	OperatorID  int              `json:"operator_id"`
	ViewMode    OperatorViewMode `json:"view_mode"`
	Operation   *Operation       `json:"operation,omitempty"`
	Operator    *Operator        `json:"operator,omitempty"`
}

// Operator represents a Mythic operator (user).
type Operator struct {
	ID                 int        `json:"id"`
	Username           string     `json:"username"`
	Admin              bool       `json:"admin"`
	Active             bool       `json:"active"`
	CreationTime       time.Time  `json:"creation_time"`
	LastLogin          time.Time  `json:"last_login"`
	CurrentOperationID *int       `json:"current_operation_id"`
	ViewUTCTime        bool       `json:"view_utc_time"`
	Deleted            bool       `json:"deleted"`
	ViewMode           string     `json:"view_mode"`
	CurrentOperation   *Operation `json:"current_operation,omitempty"`
}

// OperationEventLog represents an event log entry for an operation.
type OperationEventLog struct {
	ID          int        `json:"id"`
	OperatorID  int        `json:"operator_id"`
	OperationID int        `json:"operation_id"`
	Message     string     `json:"message"`
	Timestamp   time.Time  `json:"timestamp"`
	Level       string     `json:"level"`
	Source      string     `json:"source"`
	Deleted     bool       `json:"deleted"`
	Operator    *Operator  `json:"operator,omitempty"`
	Operation   *Operation `json:"operation,omitempty"`
}

// CreateOperationRequest represents a request to create a new operation.
type CreateOperationRequest struct {
	Name    string `json:"name"`
	AdminID *int   `json:"admin_id,omitempty"`
	Channel string `json:"channel,omitempty"`
	Webhook string `json:"webhook,omitempty"`
}

// UpdateOperationRequest represents a request to update an operation.
type UpdateOperationRequest struct {
	OperationID int     `json:"operation_id"`
	Name        *string `json:"name,omitempty"`
	Channel     *string `json:"channel,omitempty"`
	Complete    *bool   `json:"complete,omitempty"`
	Webhook     *string `json:"webhook,omitempty"`
	AdminID     *int    `json:"admin_id,omitempty"`
	BannerText  *string `json:"banner_text,omitempty"`
	BannerColor *string `json:"banner_color,omitempty"`
}

// UpdateOperatorOperationRequest represents a request to update an operator's role in an operation.
type UpdateOperatorOperationRequest struct {
	OperatorID  int               `json:"operator_id"`
	OperationID int               `json:"operation_id"`
	ViewMode    *OperatorViewMode `json:"view_mode,omitempty"`
	Remove      bool              `json:"remove,omitempty"` // If true, remove operator from operation
}

// CreateOperationEventLogRequest represents a request to create an event log entry.
type CreateOperationEventLogRequest struct {
	OperationID int    `json:"operation_id"`
	Message     string `json:"message"`
	Level       string `json:"level,omitempty"`  // info, warning, error
	Source      string `json:"source,omitempty"` // Source of the event
}

// String returns a string representation of an Operation.
func (o *Operation) String() string {
	status := "active"
	if o.Complete {
		status = "complete"
	}
	return fmt.Sprintf("%s (%s)", o.Name, status)
}

// IsComplete returns true if the operation is marked as complete.
func (o *Operation) IsComplete() bool {
	return o.Complete
}

// String returns a string representation of an OperationEventLog.
func (e *OperationEventLog) String() string {
	return fmt.Sprintf("[%s] %s: %s", e.Level, e.Timestamp.Format("2006-01-02 15:04:05"), e.Message)
}
