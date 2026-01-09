package types

import (
	"fmt"
	"time"
)

// Alert represents an operational security alert in Mythic.
// Alerts provide automated monitoring for suspicious activities,
// policy violations, and security events during operations.
type Alert struct {
	ID          int       `json:"id"`
	Message     string    `json:"message"`      // Alert description
	AlertType   string    `json:"alert"` // Category: opsec, error, warning, info
	Source      string    `json:"source"`       // Component that triggered alert
	Severity    int       `json:"severity"`     // Severity level (1-5, 5 being highest)
	Resolved    bool      `json:"resolved"`     // Whether alert has been acknowledged
	OperationID int       `json:"operation_id"` // Associated operation
	CallbackID  *int      `json:"callback_id,omitempty"` // Associated callback (if applicable)
	Timestamp   time.Time `json:"timestamp"`    // When alert was created

	// Related entities (populated with nested queries)
	OperatorID  *int   `json:"operator_id,omitempty"` // Operator who resolved (if resolved)
	ResolvedAt  *time.Time `json:"resolved_at,omitempty"` // When alert was resolved
	Notes       string `json:"notes,omitempty"`       // Resolution notes
}

// AlertType constants for common alert categories.
const (
	AlertTypeOPSEC   = "opsec"   // OPSEC violations
	AlertTypeError   = "error"   // Error conditions
	AlertTypeWarning = "warning" // Warnings
	AlertTypeInfo    = "info"    // Informational alerts
)

// AlertSeverity constants for alert severity levels.
const (
	AlertSeverityInfo     = 1 // Informational
	AlertSeverityLow      = 2 // Low severity
	AlertSeverityMedium   = 3 // Medium severity
	AlertSeverityHigh     = 4 // High severity
	AlertSeverityCritical = 5 // Critical severity
)

// CreateAlertRequest represents a request to create a new alert.
type CreateAlertRequest struct {
	Message     string `json:"message"`               // Alert description (required)
	AlertType   string `json:"alert"`        // Alert category (required)
	Source      string `json:"source"`                // Alert source (required)
	Severity    int    `json:"severity"`              // Severity level 1-5 (required)
	OperationID int    `json:"operation_id"`          // Operation ID (required)
	CallbackID  *int   `json:"callback_id,omitempty"` // Optional callback association
}

// ResolveAlertRequest represents a request to resolve/acknowledge an alert.
type ResolveAlertRequest struct {
	AlertID int    `json:"alert_id"` // Alert ID to resolve
	Notes   string `json:"notes,omitempty"` // Optional resolution notes
}

// AlertFilter represents filtering options for alert queries.
type AlertFilter struct {
	AlertType   *string `json:"alert_type,omitempty"`   // Filter by alert type
	Severity    *int    `json:"severity,omitempty"`     // Filter by severity level
	MinSeverity *int    `json:"min_severity,omitempty"` // Filter by minimum severity
	Resolved    *bool   `json:"resolved,omitempty"`     // Filter by resolution status
	CallbackID  *int    `json:"callback_id,omitempty"`  // Filter by callback
	Limit       int     `json:"limit,omitempty"`        // Maximum results (default: 100)
}

// String returns a string representation of an Alert.
func (a *Alert) String() string {
	status := "Active"
	if a.Resolved {
		status = "Resolved"
	}

	severityStr := fmt.Sprintf("Severity %d", a.Severity)
	switch a.Severity {
	case AlertSeverityCritical:
		severityStr = "CRITICAL"
	case AlertSeverityHigh:
		severityStr = "HIGH"
	case AlertSeverityMedium:
		severityStr = "MEDIUM"
	case AlertSeverityLow:
		severityStr = "LOW"
	case AlertSeverityInfo:
		severityStr = "INFO"
	}

	return fmt.Sprintf("[%s] %s - %s: %s (%s)",
		severityStr, a.AlertType, status, a.Message, a.Timestamp.Format(time.RFC3339))
}

// IsResolved returns true if the alert has been acknowledged/resolved.
func (a *Alert) IsResolved() bool {
	return a.Resolved
}

// IsCritical returns true if the alert severity is critical (5).
func (a *Alert) IsCritical() bool {
	return a.Severity == AlertSeverityCritical
}

// IsOPSEC returns true if this is an OPSEC-related alert.
func (a *Alert) IsOPSEC() bool {
	return a.AlertType == AlertTypeOPSEC
}

// Validate validates a CreateAlertRequest.
func (r *CreateAlertRequest) Validate() error {
	if r.Message == "" {
		return fmt.Errorf("message is required")
	}
	if r.AlertType == "" {
		return fmt.Errorf("alert_type is required")
	}
	if r.Source == "" {
		return fmt.Errorf("source is required")
	}
	if r.Severity < 1 || r.Severity > 5 {
		return fmt.Errorf("severity must be between 1 and 5")
	}
	if r.OperationID == 0 {
		return fmt.Errorf("operation_id is required")
	}
	return nil
}

// Validate validates a ResolveAlertRequest.
func (r *ResolveAlertRequest) Validate() error {
	if r.AlertID == 0 {
		return fmt.Errorf("alert_id is required")
	}
	return nil
}

// SetDefaults sets default values for unspecified AlertFilter fields.
func (f *AlertFilter) SetDefaults() {
	if f.Limit == 0 {
		f.Limit = 100
	}
}
