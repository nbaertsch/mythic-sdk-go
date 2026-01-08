package types

import (
	"fmt"
	"time"
)

// BrowserScript represents a JavaScript browser script for custom UI rendering.
type BrowserScript struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Script      string    `json:"script"`
	Author      string    `json:"author,omitempty"`
	ForNewUI    bool      `json:"for_new_ui"`
	Active      bool      `json:"active"`
	Description string    `json:"description,omitempty"`
	CreatedAt   time.Time `json:"created_at,omitempty"`
	UpdatedAt   time.Time `json:"updated_at,omitempty"`
}

// BrowserScriptOperation represents the association between a browser script and an operation.
type BrowserScriptOperation struct {
	ID              int       `json:"id"`
	BrowserScriptID int       `json:"browser_script_id"`
	OperationID     int       `json:"operation_id"`
	OperatorID      *int      `json:"operator_id,omitempty"`
	ScriptName      string    `json:"script_name,omitempty"`
	Active          bool      `json:"active"`
	CreatedAt       time.Time `json:"created_at,omitempty"`
}

// CustomBrowserExportRequest represents a request for custom browser data export.
type CustomBrowserExportRequest struct {
	OperationID int                    `json:"operation_id"`
	ScriptName  string                 `json:"script_name"`
	Parameters  map[string]interface{} `json:"parameters,omitempty"`
}

// CustomBrowserExportResponse represents the response from custom browser export.
type CustomBrowserExportResponse struct {
	Status       string `json:"status"`
	Error        string `json:"error,omitempty"`
	ExportedData string `json:"exported_data,omitempty"`
}

// String returns a string representation of a BrowserScript.
func (b *BrowserScript) String() string {
	if b.Description != "" {
		return fmt.Sprintf("%s: %s", b.Name, b.Description)
	}
	return b.Name
}

// IsActive returns true if the browser script is active.
func (b *BrowserScript) IsActive() bool {
	return b.Active
}

// IsForNewUI returns true if the browser script is for the new UI.
func (b *BrowserScript) IsForNewUI() bool {
	return b.ForNewUI
}

// String returns a string representation of a BrowserScriptOperation.
func (b *BrowserScriptOperation) String() string {
	if b.ScriptName != "" {
		return fmt.Sprintf("Script '%s' for Operation %d", b.ScriptName, b.OperationID)
	}
	return fmt.Sprintf("Script ID %d for Operation %d", b.BrowserScriptID, b.OperationID)
}

// IsActive returns true if the browser script is active for the operation.
func (b *BrowserScriptOperation) IsActive() bool {
	return b.Active
}

// IsOperatorSpecific returns true if the script is assigned to a specific operator.
func (b *BrowserScriptOperation) IsOperatorSpecific() bool {
	return b.OperatorID != nil
}

// String returns a string representation of a CustomBrowserExportRequest.
func (c *CustomBrowserExportRequest) String() string {
	return fmt.Sprintf("Export using '%s' for Operation %d", c.ScriptName, c.OperationID)
}
