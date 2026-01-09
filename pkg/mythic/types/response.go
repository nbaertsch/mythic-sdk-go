package types

import (
	"fmt"
	"time"
)

// Response represents task output/response data in Mythic.
// Responses are stored separately from tasks to allow efficient output streaming.
// A single task may have multiple response entries as output is generated.
type Response struct {
	ID             int       `json:"id"`
	Response       string    `json:"response"`      // The actual output content
	Timestamp      time.Time `json:"timestamp"`     // When output was generated
	TaskID         int       `json:"task_id"`       // Associated task
	SequenceNumber *int      `json:"sequence_number,omitempty"` // For ordering multi-part responses

	// Task details (populated inline in API responses when available)
	TaskCommand    string `json:"task_command,omitempty"`     // Command name
	TaskStatus     string `json:"task_status,omitempty"`      // Task status
	TaskCallbackID int    `json:"task_callback_id,omitempty"` // Callback ID
}

// ResponseSearchRequest represents parameters for searching responses.
type ResponseSearchRequest struct {
	Query         string                 `json:"query,omitempty"`          // Full-text search query
	TaskID        *int                   `json:"task_id,omitempty"`        // Filter by specific task
	CallbackID    *int                   `json:"callback_id,omitempty"`    // Filter by callback
	OperationID   *int                   `json:"operation_id,omitempty"`   // Filter by operation
	StartTime     *time.Time             `json:"start_time,omitempty"`     // Filter by time range (start)
	EndTime       *time.Time             `json:"end_time,omitempty"`       // Filter by time range (end)
	Limit         int                    `json:"limit,omitempty"`          // Maximum results to return
	Offset        int                    `json:"offset,omitempty"`         // Pagination offset
	SortBy        string                 `json:"sort_by,omitempty"`        // Sort field (timestamp, id)
	SortOrder     string                 `json:"sort_order,omitempty"`     // Sort direction (asc, desc)
	CustomFilters map[string]interface{} `json:"custom_filters,omitempty"` // Additional GraphQL filters
}

// ResponseStatistics represents aggregated response statistics for a task.
type ResponseStatistics struct {
	TaskID         int       `json:"task_id"`
	ResponseCount  int       `json:"response_count"`
	TotalSize      int       `json:"total_size"`       // Total bytes of response data
	FirstResponse  time.Time `json:"first_response"`   // Timestamp of first response
	LatestResponse time.Time `json:"latest_response"`  // Timestamp of latest response
	IsComplete     bool      `json:"is_complete"`      // Whether task is completed
}

// String returns a string representation of a Response.
func (r *Response) String() string {
	preview := r.Response
	if len(preview) > 100 {
		preview = preview[:97] + "..."
	}
	return fmt.Sprintf("Response %d (Task %d): %s", r.ID, r.TaskID, preview)
}

// IsEmpty returns true if the response has no content.
func (r *Response) IsEmpty() bool {
	return len(r.Response) == 0
}

// Size returns the byte size of the response content.
func (r *Response) Size() int {
	return len(r.Response)
}

// Validate validates a ResponseSearchRequest.
func (r *ResponseSearchRequest) Validate() error {
	if r.Limit < 0 {
		return fmt.Errorf("limit must be non-negative")
	}
	if r.Offset < 0 {
		return fmt.Errorf("offset must be non-negative")
	}
	if r.SortBy != "" && r.SortBy != "timestamp" && r.SortBy != "id" {
		return fmt.Errorf("invalid sort_by value: must be 'timestamp' or 'id'")
	}
	if r.SortOrder != "" && r.SortOrder != "asc" && r.SortOrder != "desc" {
		return fmt.Errorf("invalid sort_order value: must be 'asc' or 'desc'")
	}
	if r.StartTime != nil && r.EndTime != nil && r.EndTime.Before(*r.StartTime) {
		return fmt.Errorf("end_time must be after start_time")
	}
	return nil
}

// SetDefaults sets default values for unspecified fields.
func (r *ResponseSearchRequest) SetDefaults() {
	if r.Limit == 0 {
		r.Limit = 100 // Default limit
	}
	if r.SortBy == "" {
		r.SortBy = "timestamp"
	}
	if r.SortOrder == "" {
		r.SortOrder = "asc"
	}
}
