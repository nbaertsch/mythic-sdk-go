package types

import (
	"fmt"
	"time"
)

// Command represents a command available in a payload type.
type Command struct {
	ID                  int       `json:"id"`
	Cmd                 string    `json:"cmd"`
	PayloadTypeID       int       `json:"payload_type_id"`
	PayloadTypeName     string    `json:"payload_type_name,omitempty"`
	Description         string    `json:"description"`
	Help                string    `json:"help_cmd"`
	Version             int       `json:"version"`
	Supported           bool      `json:"supported_ui_features"`
	Author              string    `json:"author"`
	Attributes          string    `json:"attributes"`
	ScriptOnly          bool      `json:"script_only"`
	MitreAttackMappings string    `json:"attack"`
	CreationTime        time.Time `json:"creation_time,omitempty"`
}

// CommandParameter represents a parameter for a command.
type CommandParameter struct {
	ID                          int       `json:"id"`
	CommandID                   int       `json:"command_id"`
	Name                        string    `json:"name"`
	Type                        string    `json:"type"`
	Description                 string    `json:"description"`
	Required                    bool      `json:"required"`
	DefaultValue                string    `json:"default_value,omitempty"`
	Choices                     string    `json:"choices,omitempty"`
	SupportedAgents             string    `json:"supported_agents,omitempty"`
	SupportedAgentBuildParams   string    `json:"supported_agent_build_parameters,omitempty"`
	ChoicesAreAllCommands       bool      `json:"choices_are_all_commands"`
	ChoicesAreLoadedCommands    bool      `json:"choices_are_loaded_commands"`
	ChoiceFilterByCommandAttrib string    `json:"choice_filter_by_command_attributes,omitempty"`
	DynamicQueryFunction        string    `json:"dynamic_query_function,omitempty"`
	CreationTime                time.Time `json:"creation_time,omitempty"`
}

// LoadedCommand represents a command loaded in a callback.
type LoadedCommand struct {
	ID         int       `json:"id"`
	CommandID  int       `json:"command_id"`
	CallbackID int       `json:"callback_id"`
	OperatorID int       `json:"operator_id"`
	Version    int       `json:"version"`
	Timestamp  time.Time `json:"timestamp"`
	Command    *Command  `json:"command,omitempty"`
}

// String returns a string representation of a Command.
func (c *Command) String() string {
	status := ""
	if !c.Supported {
		status = " (unsupported)"
	}
	if c.ScriptOnly {
		status += " (script-only)"
	}
	return fmt.Sprintf("%s v%d%s", c.Cmd, c.Version, status)
}

// IsSupported returns true if the command is supported.
func (c *Command) IsSupported() bool {
	return c.Supported
}

// IsScriptOnly returns true if the command is script-only.
func (c *Command) IsScriptOnly() bool {
	return c.ScriptOnly
}

// String returns a string representation of a CommandParameter.
func (cp *CommandParameter) String() string {
	required := ""
	if cp.Required {
		required = " (required)"
	}
	return fmt.Sprintf("%s (%s)%s", cp.Name, cp.Type, required)
}

// IsRequired returns true if the parameter is required.
func (cp *CommandParameter) IsRequired() bool {
	return cp.Required
}

// HasChoices returns true if the parameter has predefined choices.
func (cp *CommandParameter) HasChoices() bool {
	return cp.Choices != "" || cp.ChoicesAreAllCommands || cp.ChoicesAreLoadedCommands
}

// IsDynamic returns true if the parameter uses dynamic query function.
func (cp *CommandParameter) IsDynamic() bool {
	return cp.DynamicQueryFunction != ""
}

// String returns a string representation of a LoadedCommand.
func (lc *LoadedCommand) String() string {
	if lc.Command != nil {
		return fmt.Sprintf("%s v%d (Callback %d)", lc.Command.Cmd, lc.Version, lc.CallbackID)
	}
	return fmt.Sprintf("Command ID %d v%d (Callback %d)", lc.CommandID, lc.Version, lc.CallbackID)
}

// Command parameter type constants
const (
	ParameterTypeString         = "String"
	ParameterTypeBoolean        = "Boolean"
	ParameterTypeNumber         = "Number"
	ParameterTypeChooseOne      = "ChooseOne"
	ParameterTypeChooseMultiple = "ChooseMultiple"
	ParameterTypeFile           = "File"
	ParameterTypeArray          = "Array"
	ParameterTypeCredential     = "Credential"
	ParameterTypeLinkInfo       = "LinkInfo"
)
