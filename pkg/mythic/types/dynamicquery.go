package types

import "fmt"

// DynamicQueryRequest represents a request to execute a dynamic query function.
type DynamicQueryRequest struct {
	Command    string                 `json:"command"`
	Parameters map[string]interface{} `json:"parameters"`
	CallbackID int                    `json:"callback_id,omitempty"`
}

// String returns a human-readable representation of the request.
func (d *DynamicQueryRequest) String() string {
	if d.CallbackID > 0 {
		return fmt.Sprintf("Dynamic query for command '%s' (callback %d)", d.Command, d.CallbackID)
	}
	return fmt.Sprintf("Dynamic query for command '%s'", d.Command)
}

// DynamicQueryResponse represents the response from a dynamic query.
type DynamicQueryResponse struct {
	Status  string        `json:"status"`
	Choices []interface{} `json:"choices"`
	Error   string        `json:"error,omitempty"`
}

// String returns a human-readable representation of the response.
func (d *DynamicQueryResponse) String() string {
	if d.Status == "success" {
		return fmt.Sprintf("Query returned %d choices", len(d.Choices))
	}
	return fmt.Sprintf("Query failed: %s", d.Error)
}

// IsSuccessful returns true if the query succeeded.
func (d *DynamicQueryResponse) IsSuccessful() bool {
	return d.Status == "success"
}

// HasChoices returns true if the query returned any choices.
func (d *DynamicQueryResponse) HasChoices() bool {
	return len(d.Choices) > 0
}

// DynamicBuildParameterRequest represents a request for build parameter dynamic query.
type DynamicBuildParameterRequest struct {
	PayloadType string                 `json:"payload_type"`
	Parameter   string                 `json:"parameter"`
	Parameters  map[string]interface{} `json:"parameters,omitempty"`
}

// String returns a human-readable representation of the request.
func (d *DynamicBuildParameterRequest) String() string {
	return fmt.Sprintf("Dynamic build parameter query for %s.%s", d.PayloadType, d.Parameter)
}

// DynamicBuildParameterResponse represents the response from a build parameter query.
type DynamicBuildParameterResponse struct {
	Status  string        `json:"status"`
	Choices []interface{} `json:"choices"`
	Error   string        `json:"error,omitempty"`
}

// String returns a human-readable representation of the response.
func (d *DynamicBuildParameterResponse) String() string {
	if d.Status == "success" {
		return fmt.Sprintf("Query returned %d choices", len(d.Choices))
	}
	return fmt.Sprintf("Query failed: %s", d.Error)
}

// IsSuccessful returns true if the query succeeded.
func (d *DynamicBuildParameterResponse) IsSuccessful() bool {
	return d.Status == "success"
}

// HasChoices returns true if the query returned any choices.
func (d *DynamicBuildParameterResponse) HasChoices() bool {
	return len(d.Choices) > 0
}

// TypedArrayParseRequest represents a request to parse a typed array.
type TypedArrayParseRequest struct {
	InputArray    string `json:"input_array"`
	ParameterType string `json:"parameter_type"`
}

// String returns a human-readable representation of the request.
func (t *TypedArrayParseRequest) String() string {
	return fmt.Sprintf("Parse typed array for parameter type '%s'", t.ParameterType)
}

// TypedArrayParseResponse represents the response from parsing a typed array.
type TypedArrayParseResponse struct {
	Status      string        `json:"status"`
	ParsedArray []interface{} `json:"parsed_array"`
	Error       string        `json:"error,omitempty"`
}

// String returns a human-readable representation of the response.
func (t *TypedArrayParseResponse) String() string {
	if t.Status == "success" {
		return fmt.Sprintf("Parsed %d array elements", len(t.ParsedArray))
	}
	return fmt.Sprintf("Parse failed: %s", t.Error)
}

// IsSuccessful returns true if the parse succeeded.
func (t *TypedArrayParseResponse) IsSuccessful() bool {
	return t.Status == "success"
}

// HasElements returns true if the parsed array has elements.
func (t *TypedArrayParseResponse) HasElements() bool {
	return len(t.ParsedArray) > 0
}
