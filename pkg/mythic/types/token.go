package types

import (
	"fmt"
	"time"
)

// Token represents a process or user security token.
type Token struct {
	ID              int        `json:"id"`
	TokenID         string     `json:"token_id"`
	User            string     `json:"user"`
	Groups          string     `json:"groups"`
	Privileges      string     `json:"privileges"`
	ThreadID        int        `json:"thread_id"`
	ProcessID       int        `json:"process_id"`
	SessionID       int        `json:"session_id"`
	LogonSID        string     `json:"logon_sid"`
	IntegrityLevel  int        `json:"integrity_level_int"`
	Restricted      bool       `json:"restricted"`
	DefaultDACL     string     `json:"default_dacl"`
	Handle          string     `json:"handle"`
	Capabilities    string     `json:"capabilities"`
	AppContainerSID string     `json:"app_container_sid"`
	AppContainerNum int        `json:"app_container_number"`
	TaskID          *int       `json:"task_id,omitempty"`
	OperationID     int        `json:"operation_id"`
	Timestamp       time.Time  `json:"timestamp"`
	Host            string     `json:"host"`
	Deleted         bool       `json:"deleted"`
	Operation       *Operation `json:"operation,omitempty"`
}

// String returns a string representation of a Token.
func (t *Token) String() string {
	if t.User != "" && t.Host != "" {
		return fmt.Sprintf("%s on %s", t.User, t.Host)
	}
	if t.User != "" {
		return t.User
	}
	if t.TokenID != "" {
		return fmt.Sprintf("Token %s", t.TokenID)
	}
	return fmt.Sprintf("Token %d", t.ID)
}

// IsDeleted returns true if the token is marked as deleted.
func (t *Token) IsDeleted() bool {
	return t.Deleted
}

// HasTask returns true if the token is linked to a task.
func (t *Token) HasTask() bool {
	return t.TaskID != nil && *t.TaskID > 0
}

// GetIntegrityLevelString returns a human-readable integrity level.
func (t *Token) GetIntegrityLevelString() string {
	switch t.IntegrityLevel {
	case 0:
		return "Untrusted"
	case 1:
		return "Low"
	case 2:
		return "Medium"
	case 3:
		return "High"
	case 4:
		return "System"
	default:
		return "Unknown"
	}
}

// CallbackToken represents the association between a callback and a token.
type CallbackToken struct {
	ID         int       `json:"id"`
	CallbackID int       `json:"callback_id"`
	TokenID    int       `json:"token_id"`
	Timestamp  time.Time `json:"timestamp"`
	Token      *Token    `json:"token,omitempty"`
}

// String returns a string representation of a CallbackToken.
func (ct *CallbackToken) String() string {
	if ct.Token != nil {
		return fmt.Sprintf("Token %s for Callback %d", ct.Token.String(), ct.CallbackID)
	}
	return fmt.Sprintf("CallbackToken %d (Callback %d, Token %d)", ct.ID, ct.CallbackID, ct.TokenID)
}

// APIToken represents an API authentication token.
type APIToken struct {
	ID           int        `json:"id"`
	TokenValue   string     `json:"token_value"`
	TokenType    string     `json:"token_type"`
	Active       bool       `json:"active"`
	CreationTime time.Time  `json:"creation_time"`
	OperatorID   int        `json:"operator_id"`
	OperationID  *int       `json:"operation_id,omitempty"`
	Name         string     `json:"name"`
	Deleted      bool       `json:"deleted"`
	Operator     *Operator  `json:"operator,omitempty"`
	Operation    *Operation `json:"operation,omitempty"`
}

// String returns a string representation of an APIToken.
func (at *APIToken) String() string {
	if at.Name != "" {
		return fmt.Sprintf("%s (%s)", at.Name, at.TokenType)
	}
	return fmt.Sprintf("APIToken %d (%s)", at.ID, at.TokenType)
}

// IsActive returns true if the token is active.
func (at *APIToken) IsActive() bool {
	return at.Active && !at.Deleted
}

// IsDeleted returns true if the token is marked as deleted.
func (at *APIToken) IsDeleted() bool {
	return at.Deleted
}
