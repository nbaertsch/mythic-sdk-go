package mythic

import (
	"context"

	"github.com/nbaertsch/mythic-sdk-go/pkg/mythic/types"
)

// GetCommands retrieves all available commands from all payload types.
func (c *Client) GetCommands(ctx context.Context) ([]*types.Command, error) {
	if err := c.EnsureAuthenticated(ctx); err != nil {
		return nil, err
	}

	var query struct {
		Commands []struct {
			ID              int    `graphql:"id"`
			Cmd             string `graphql:"cmd"`
			PayloadTypeID   int    `graphql:"payloadtype_id"`
			Description     string `graphql:"description"`
			Help            string `graphql:"help_cmd"`
			Version         int    `graphql:"version"`
			Supported       bool   `graphql:"supported_ui_features"`
			Author          string `graphql:"author"`
			Attributes      string `graphql:"attributes"`
			ScriptOnly      bool   `graphql:"script_only"`
			AttackMappings  string `graphql:"attack"`
		} `graphql:"command(order_by: {cmd: asc})"`
	}

	err := c.executeQuery(ctx, &query, nil)
	if err != nil {
		return nil, WrapError("GetCommands", err, "failed to query commands")
	}

	commands := make([]*types.Command, len(query.Commands))
	for i, cmd := range query.Commands {
		commands[i] = &types.Command{
			ID:                  cmd.ID,
			Cmd:                 cmd.Cmd,
			PayloadTypeID:       cmd.PayloadTypeID,
			Description:         cmd.Description,
			Help:                cmd.Help,
			Version:             cmd.Version,
			Supported:           cmd.Supported,
			Author:              cmd.Author,
			Attributes:          cmd.Attributes,
			ScriptOnly:          cmd.ScriptOnly,
			MitreAttackMappings: cmd.AttackMappings,
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
