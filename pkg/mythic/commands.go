package mythic

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/nbaertsch/mythic-sdk-go/pkg/mythic/types"
)

// GetCommands retrieves all available commands from all payload types.
func (c *Client) GetCommands(ctx context.Context) ([]*types.Command, error) {
	if err := c.EnsureAuthenticated(ctx); err != nil {
		return nil, err
	}

	var query struct {
		Commands []struct {
			ID             int    `graphql:"id"`
			Cmd            string `graphql:"cmd"`
			PayloadTypeID  int    `graphql:"payload_type_id"`
			Description    string `graphql:"description"`
			Help           string `graphql:"help_cmd"`
			Version        int    `graphql:"version"`
			Author         string `graphql:"author"`
			ScriptOnly     bool   `graphql:"script_only"`
		} `graphql:"command(order_by: {cmd: asc})"`
	}

	err := c.executeQuery(ctx, &query, nil)
	if err != nil {
		return nil, WrapError("GetCommands", err, "failed to query commands")
	}

	commands := make([]*types.Command, len(query.Commands))
	for i, cmd := range query.Commands {
		commands[i] = &types.Command{
			ID:            cmd.ID,
			Cmd:           cmd.Cmd,
			PayloadTypeID: cmd.PayloadTypeID,
			Description:   cmd.Description,
			Help:          cmd.Help,
			Version:       cmd.Version,
			Author:        cmd.Author,
			ScriptOnly:    cmd.ScriptOnly,
		}
	}

	return commands, nil
}

// GetCommandParameters retrieves all parameters for all commands.
func (c *Client) GetCommandParameters(ctx context.Context) ([]*types.CommandParameter, error) {
	if err := c.EnsureAuthenticated(ctx); err != nil {
		return nil, err
	}

	var query struct {
		Parameters []struct {
			ID                          int    `graphql:"id"`
			CommandID                   int    `graphql:"command_id"`
			Name                        string `graphql:"name"`
			Type                        string `graphql:"type"`
			Description                 string `graphql:"description"`
			Required                    bool   `graphql:"required"`
			DefaultValue                string `graphql:"default_value"`
			Choices                     string `graphql:"choices"`
			SupportedAgents             string `graphql:"supported_agents"`
			SupportedAgentBuildParams   string `graphql:"supported_agent_build_parameters"`
			ChoicesAreAllCommands       bool   `graphql:"choices_are_all_commands"`
			ChoicesAreLoadedCommands    bool   `graphql:"choices_are_loaded_commands"`
			ChoiceFilterByCommandAttrib string `graphql:"choice_filter_by_command_attributes"`
			DynamicQueryFunction        string `graphql:"dynamic_query_function"`
		} `graphql:"commandparameters(order_by: {command_id: asc})"`
	}

	err := c.executeQuery(ctx, &query, nil)
	if err != nil {
		return nil, WrapError("GetCommandParameters", err, "failed to query command parameters")
	}

	parameters := make([]*types.CommandParameter, len(query.Parameters))
	for i, param := range query.Parameters {
		parameters[i] = &types.CommandParameter{
			ID:                          param.ID,
			CommandID:                   param.CommandID,
			Name:                        param.Name,
			Type:                        param.Type,
			Description:                 param.Description,
			Required:                    param.Required,
			DefaultValue:                param.DefaultValue,
			Choices:                     param.Choices,
			SupportedAgents:             param.SupportedAgents,
			SupportedAgentBuildParams:   param.SupportedAgentBuildParams,
			ChoicesAreAllCommands:       param.ChoicesAreAllCommands,
			ChoicesAreLoadedCommands:    param.ChoicesAreLoadedCommands,
			ChoiceFilterByCommandAttrib: param.ChoiceFilterByCommandAttrib,
			DynamicQueryFunction:        param.DynamicQueryFunction,
		}
	}

	return parameters, nil
}

// GetLoadedCommands retrieves all commands loaded in a specific callback.
func (c *Client) GetLoadedCommands(ctx context.Context, callbackID int) ([]*types.LoadedCommand, error) {
	if err := c.EnsureAuthenticated(ctx); err != nil {
		return nil, err
	}

	if callbackID == 0 {
		return nil, WrapError("GetLoadedCommands", ErrInvalidInput, "callback ID is required")
	}

	var query struct {
		LoadedCommands []struct {
			ID         int `graphql:"id"`
			CommandID  int `graphql:"command_id"`
			CallbackID int `graphql:"callback_id"`
			OperatorID int `graphql:"operator_id"`
			Version    int `graphql:"version"`
			Command    struct {
				Cmd         string `graphql:"cmd"`
				Description string `graphql:"description"`
				Version     int    `graphql:"version"`
			} `graphql:"command"`
		} `graphql:"loadedcommands(where: {callback_id: {_eq: $callback_id}}, order_by: {command: {cmd: asc}})"`
	}

	variables := map[string]interface{}{
		"callback_id": callbackID,
	}

	err := c.executeQuery(ctx, &query, variables)
	if err != nil {
		return nil, WrapError("GetLoadedCommands", err, "failed to query loaded commands")
	}

	loadedCommands := make([]*types.LoadedCommand, len(query.LoadedCommands))
	for i, lc := range query.LoadedCommands {
		loadedCommands[i] = &types.LoadedCommand{
			ID:         lc.ID,
			CommandID:  lc.CommandID,
			CallbackID: lc.CallbackID,
			OperatorID: lc.OperatorID,
			Version:    lc.Version,
			Command: &types.Command{
				Cmd:         lc.Command.Cmd,
				Description: lc.Command.Description,
				Version:     lc.Command.Version,
			},
		}
	}

	return loadedCommands, nil
}

// CommandWithParameters represents a command with its parameter definitions.
type CommandWithParameters struct {
	Command    *types.Command
	Parameters []*types.CommandParameter
}

// GetCommandWithParameters retrieves a specific command by name with all its parameters.
// This is useful for building task parameters dynamically.
func (c *Client) GetCommandWithParameters(ctx context.Context, payloadTypeID int, commandName string) (*CommandWithParameters, error) {
	if err := c.EnsureAuthenticated(ctx); err != nil {
		return nil, err
	}

	if commandName == "" {
		return nil, WrapError("GetCommandWithParameters", ErrInvalidInput, "command name is required")
	}

	var query struct {
		Command []struct {
			ID                int    `graphql:"id"`
			Cmd               string `graphql:"cmd"`
			PayloadTypeID     int    `graphql:"payload_type_id"`
			Description       string `graphql:"description"`
			Help              string `graphql:"help_cmd"`
			Version           int    `graphql:"version"`
			Author            string `graphql:"author"`
			ScriptOnly        bool   `graphql:"script_only"`
			CommandParameters []struct {
				ID                          int    `graphql:"id"`
				CommandID                   int    `graphql:"command_id"`
				Name                        string `graphql:"name"`
				Type                        string `graphql:"type"`
				Description                 string `graphql:"description"`
				Required                    bool   `graphql:"required"`
				DefaultValue                string `graphql:"default_value"`
				Choices                     string `graphql:"choices"`
				SupportedAgents             string `graphql:"supported_agents"`
				SupportedAgentBuildParams   string `graphql:"supported_agent_build_parameters"`
				ChoicesAreAllCommands       bool   `graphql:"choices_are_all_commands"`
				ChoicesAreLoadedCommands    bool   `graphql:"choices_are_loaded_commands"`
				ChoiceFilterByCommandAttrib string `graphql:"choice_filter_by_command_attributes"`
				DynamicQueryFunction        string `graphql:"dynamic_query_function"`
			} `graphql:"commandparameters(order_by: {name: asc})"`
		} `graphql:"command(where: {cmd: {_eq: $cmd}, payload_type_id: {_eq: $payload_type_id}}, limit: 1)"`
	}

	variables := map[string]interface{}{
		"cmd":             commandName,
		"payload_type_id": payloadTypeID,
	}

	err := c.executeQuery(ctx, &query, variables)
	if err != nil {
		return nil, WrapError("GetCommandWithParameters", err, "failed to query command with parameters")
	}

	if len(query.Command) == 0 {
		return nil, WrapError("GetCommandWithParameters", ErrNotFound, "command not found")
	}

	cmd := query.Command[0]
	command := &types.Command{
		ID:            cmd.ID,
		Cmd:           cmd.Cmd,
		PayloadTypeID: cmd.PayloadTypeID,
		Description:   cmd.Description,
		Help:          cmd.Help,
		Version:       cmd.Version,
		Author:        cmd.Author,
		ScriptOnly:    cmd.ScriptOnly,
	}

	parameters := make([]*types.CommandParameter, len(cmd.CommandParameters))
	for i, param := range cmd.CommandParameters {
		parameters[i] = &types.CommandParameter{
			ID:                          param.ID,
			CommandID:                   param.CommandID,
			Name:                        param.Name,
			Type:                        param.Type,
			Description:                 param.Description,
			Required:                    param.Required,
			DefaultValue:                param.DefaultValue,
			Choices:                     param.Choices,
			SupportedAgents:             param.SupportedAgents,
			SupportedAgentBuildParams:   param.SupportedAgentBuildParams,
			ChoicesAreAllCommands:       param.ChoicesAreAllCommands,
			ChoicesAreLoadedCommands:    param.ChoicesAreLoadedCommands,
			ChoiceFilterByCommandAttrib: param.ChoiceFilterByCommandAttrib,
			DynamicQueryFunction:        param.DynamicQueryFunction,
		}
	}

	return &CommandWithParameters{
		Command:    command,
		Parameters: parameters,
	}, nil
}

// IsRawStringCommand returns true if the command expects raw string parameters
// rather than a JSON object. Commands WITHOUT CommandParameters defined expect
// raw string params (e.g., shell, run, execute).
func (cwp *CommandWithParameters) IsRawStringCommand() bool {
	return len(cwp.Parameters) == 0
}

// HasRequiredParameters returns true if the command has any required parameters.
func (cwp *CommandWithParameters) HasRequiredParameters() bool {
	for _, param := range cwp.Parameters {
		if param.Required {
			return true
		}
	}
	return false
}

// BuildTaskParams constructs the params string for a task based on command definition.
// For raw string commands (no parameters defined), returns the input directly as a string.
// For parameterized commands, builds a JSON object from the input map and returns it as a string.
func (cwp *CommandWithParameters) BuildTaskParams(input interface{}) (string, error) {
	// Raw string command - return input as-is (convert to string if needed)
	if cwp.IsRawStringCommand() {
		switch v := input.(type) {
		case string:
			return v, nil
		case map[string]interface{}:
			// If user provided a map but command expects raw string,
			// check if there's a "raw" or "command" key
			if rawVal, ok := v["raw"]; ok {
				return fmt.Sprintf("%v", rawVal), nil
			}
			if cmdVal, ok := v["command"]; ok {
				return fmt.Sprintf("%v", cmdVal), nil
			}
			return "", WrapError("BuildTaskParams", ErrInvalidInput,
				"command expects raw string params but got map without 'raw' or 'command' key")
		default:
			return fmt.Sprintf("%v", v), nil
		}
	}

	// Parameterized command - build JSON object
	inputMap, ok := input.(map[string]interface{})
	if !ok {
		return "", WrapError("BuildTaskParams", ErrInvalidInput,
			"parameterized command expects map[string]interface{} input")
	}

	// Validate required parameters
	for _, param := range cwp.Parameters {
		if param.Required {
			if _, exists := inputMap[param.Name]; !exists {
				// Check if there's a default value
				if param.DefaultValue != "" {
					inputMap[param.Name] = param.DefaultValue
				} else {
					return "", WrapError("BuildTaskParams", ErrInvalidInput,
						fmt.Sprintf("required parameter '%s' is missing", param.Name))
				}
			}
		}
	}

	// Marshal to JSON
	paramsJSON, err := json.Marshal(inputMap)
	if err != nil {
		return "", WrapError("BuildTaskParams", err, "failed to marshal parameters to JSON")
	}

	return string(paramsJSON), nil
}
