package types

import "fmt"

// C2ProfileParameter represents a configuration parameter for a C2 profile.
// These define what configuration options are available when setting up a C2 profile
// for a payload (e.g., callback_host, callback_port, etc.).
type C2ProfileParameter struct {
	ID            int    `json:"id"`
	C2ProfileID   int    `json:"c2_profile_id"`
	Name          string `json:"name"`
	Description   string `json:"description"`
	DefaultValue  string `json:"default_value"`
	ParameterType string `json:"parameter_type"` // String, Boolean, Number, ChooseOne, etc.
	Required      bool   `json:"required"`
	Randomize     bool   `json:"randomize"`
	FormatString  string `json:"format_string"`
	VerifierRegex string `json:"verifier_regex"`
	IsCryptoType  bool   `json:"crypto_type"`
	Deleted       bool   `json:"deleted"`
	Choices       string `json:"choices"` // JSON array of choices (for ChooseOne/ChooseMultiple)
}

// String returns a string representation of a C2ProfileParameter.
func (p *C2ProfileParameter) String() string {
	required := ""
	if p.Required {
		required = " (required)"
	}
	return fmt.Sprintf("%s (%s)%s", p.Name, p.ParameterType, required)
}
