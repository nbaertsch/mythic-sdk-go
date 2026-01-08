package types

import (
	"fmt"
)

// GenerateReportRequest represents a request to generate an operation report.
type GenerateReportRequest struct {
	OperationID        int     `json:"operation_id"`
	IncludeMITRE       bool    `json:"include_mitre,omitempty"`
	OutputFormat       string  `json:"output_format,omitempty"`
	IncludeCallbacks   bool    `json:"include_callbacks,omitempty"`
	IncludeTasks       bool    `json:"include_tasks,omitempty"`
	IncludeFiles       bool    `json:"include_files,omitempty"`
	IncludeCredentials bool    `json:"include_credentials,omitempty"`
	IncludeArtifacts   bool    `json:"include_artifacts,omitempty"`
	StartDate          *string `json:"start_date,omitempty"`
	EndDate            *string `json:"end_date,omitempty"`
	CallbackIDs        []int   `json:"callback_ids,omitempty"`
}

// GenerateReportResponse represents the response from generating a report.
type GenerateReportResponse struct {
	Status     string `json:"status"`
	Error      string `json:"error,omitempty"`
	ReportData string `json:"report_data,omitempty"`
}

// String returns a string representation of a GenerateReportRequest.
func (r *GenerateReportRequest) String() string {
	return fmt.Sprintf("Report for Operation %d", r.OperationID)
}

// RedirectRule represents a C2 redirect rule configuration.
type RedirectRule struct {
	ID           int    `json:"id"`
	PayloadType  string `json:"payload_type"`
	C2Profile    string `json:"c2_profile"`
	RedirectType string `json:"redirect_type"`
	Rule         string `json:"rule"`
	Description  string `json:"description,omitempty"`
}

// String returns a string representation of a RedirectRule.
func (r *RedirectRule) String() string {
	if r.Description != "" {
		return fmt.Sprintf("%s (%s - %s): %s", r.RedirectType, r.PayloadType, r.C2Profile, r.Description)
	}
	return fmt.Sprintf("%s (%s - %s)", r.RedirectType, r.PayloadType, r.C2Profile)
}

// Report output format constants
const (
	ReportFormatJSON     = "json"
	ReportFormatMarkdown = "markdown"
	ReportFormatHTML     = "html"
	ReportFormatPDF      = "pdf"
)

// Redirect rule type constants
const (
	RedirectTypeApache     = "apache"
	RedirectTypeNginx      = "nginx"
	RedirectTypeModRewrite = "mod_rewrite"
)
