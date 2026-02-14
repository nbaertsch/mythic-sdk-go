package types

import (
	"fmt"
	"time"
)

// BuildParameterType represents the definition of a build parameter for a payload type.
// This defines what parameters are available when building payloads of a specific type.
type BuildParameterType struct {
	ID                   int    `json:"id"`
	Name                 string `json:"name"`
	PayloadTypeID        int    `json:"payload_type_id"`
	Description          string `json:"description"`
	Required             bool   `json:"required"`
	VerifierRegex        string `json:"verifier_regex"`
	DefaultValue         string `json:"default_value"`
	Randomize            bool   `json:"randomize"`
	FormatString         string `json:"format_string"`
	ParameterType        string `json:"parameter_type"` // String, Boolean, Number, ChooseOne, etc.
	IsCryptoType         bool   `json:"crypto_type"`
	Deleted              bool   `json:"deleted"`
	GroupName            string `json:"group_name"`
	Choices              string `json:"choices"`         // JSON array of choices
	SupportedOS          string `json:"supported_os"`    // JSON array of supported OS
	HideConditions       string `json:"hide_conditions"` // JSON array of hide conditions
	DynamicQueryFunction string `json:"dynamic_query_function"`
	UIPosition           int    `json:"ui_position"`
}

// BuildParameterInstance represents an actual value set for a build parameter in a specific payload.
// This stores the concrete values used when building a payload.
type BuildParameterInstance struct {
	ID               int                 `json:"id"`
	PayloadID        int                 `json:"payload_id"`
	BuildParameterID int                 `json:"build_parameter_id"`
	Value            string              `json:"value"`
	EncValue         *string             `json:"enc_value,omitempty"` // Encrypted value for sensitive parameters
	CreationTime     time.Time           `json:"creation_time"`
	BuildParameter   *BuildParameterType `json:"buildparameter,omitempty"`
	Payload          *Payload            `json:"payload,omitempty"`
}

// String returns a string representation of a BuildParameterType.
func (b *BuildParameterType) String() string {
	required := ""
	if b.Required {
		required = " (required)"
	}
	return fmt.Sprintf("%s (%s)%s", b.Name, b.ParameterType, required)
}

// IsRequired returns true if the build parameter is required.
func (b *BuildParameterType) IsRequired() bool {
	return b.Required
}

// IsCrypto returns true if the build parameter is a crypto type.
func (b *BuildParameterType) IsCrypto() bool {
	return b.IsCryptoType
}

// ShouldRandomize returns true if the parameter should be randomized.
func (b *BuildParameterType) ShouldRandomize() bool {
	return b.Randomize
}

// IsDeleted returns true if the build parameter has been deleted.
func (b *BuildParameterType) IsDeleted() bool {
	return b.Deleted
}

// String returns a string representation of a BuildParameterInstance.
func (b *BuildParameterInstance) String() string {
	if b.BuildParameter != nil {
		return fmt.Sprintf("%s = %s", b.BuildParameter.Name, b.Value)
	}
	return fmt.Sprintf("BuildParameter %d = %s", b.BuildParameterID, b.Value)
}

// IsEncrypted returns true if the parameter value is encrypted.
func (b *BuildParameterInstance) IsEncrypted() bool {
	return b.EncValue != nil && *b.EncValue != ""
}

// GetValue returns the parameter value (decrypted value takes precedence).
func (b *BuildParameterInstance) GetValue() string {
	if b.IsEncrypted() {
		return *b.EncValue
	}
	return b.Value
}

// Build parameter type constants (same as command parameter types)
const (
	BuildParameterTypeString         = "String"
	BuildParameterTypeBoolean        = "Boolean"
	BuildParameterTypeNumber         = "Number"
	BuildParameterTypeChooseOne      = "ChooseOne"
	BuildParameterTypeChooseMultiple = "ChooseMultiple"
	BuildParameterTypeFile           = "File"
	BuildParameterTypeArray          = "Array"
	BuildParameterTypeDate           = "Date"
)
