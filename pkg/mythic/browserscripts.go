package mythic

import (
	"context"

	"github.com/nbaertsch/mythic-sdk-go/pkg/mythic/types"
)

// GetBrowserScripts retrieves all browser scripts available in the system.
// Browser scripts are JavaScript files used for custom UI rendering in the Mythic web interface.
func (c *Client) GetBrowserScripts(ctx context.Context) ([]*types.BrowserScript, error) {
	if err := c.EnsureAuthenticated(ctx); err != nil {
		return nil, err
	}

	var query struct {
		BrowserScripts []struct {
			ID          int    `graphql:"id"`
			Name        string `graphql:"name"`
			Script      string `graphql:"script"`
			Author      string `graphql:"author"`
			ForNewUI    bool   `graphql:"for_new_ui"`
			Active      bool   `graphql:"active"`
			Description string `graphql:"description"`
		} `graphql:"browserscript"`
	}

	err := c.executeQuery(ctx, &query, nil)
	if err != nil {
		return nil, WrapError("GetBrowserScripts", err, "failed to query browser scripts")
	}

	scripts := make([]*types.BrowserScript, len(query.BrowserScripts))
	for i, s := range query.BrowserScripts {
		scripts[i] = &types.BrowserScript{
			ID:          s.ID,
			Name:        s.Name,
			Script:      s.Script,
			Author:      s.Author,
			ForNewUI:    s.ForNewUI,
			Active:      s.Active,
			Description: s.Description,
		}
	}

	return scripts, nil
}

// GetBrowserScriptsByOperation retrieves browser scripts associated with a specific operation.
// This allows filtering scripts that have been enabled or customized for a particular operation.
func (c *Client) GetBrowserScriptsByOperation(ctx context.Context, operationID int) ([]*types.BrowserScriptOperation, error) {
	if err := c.EnsureAuthenticated(ctx); err != nil {
		return nil, err
	}

	if operationID == 0 {
		return nil, WrapError("GetBrowserScriptsByOperation", ErrInvalidInput, "operation ID is required")
	}

	var query struct {
		BrowserScriptOperations []struct {
			ID              int  `graphql:"id"`
			BrowserScriptID int  `graphql:"browserscript_id"`
			OperationID     int  `graphql:"operation_id"`
			OperatorID      *int `graphql:"operator_id"`
			Active          bool `graphql:"active"`
			Script          struct {
				Name string `graphql:"name"`
			} `graphql:"browserscript"`
		} `graphql:"browserscriptoperation(where: {operation_id: {_eq: $operation_id}})"`
	}

	variables := map[string]interface{}{
		"operation_id": operationID,
	}

	err := c.executeQuery(ctx, &query, variables)
	if err != nil {
		return nil, WrapError("GetBrowserScriptsByOperation", err, "failed to query browser scripts for operation")
	}

	scripts := make([]*types.BrowserScriptOperation, len(query.BrowserScriptOperations))
	for i, s := range query.BrowserScriptOperations {
		scripts[i] = &types.BrowserScriptOperation{
			ID:              s.ID,
			BrowserScriptID: s.BrowserScriptID,
			OperationID:     s.OperationID,
			OperatorID:      s.OperatorID,
			ScriptName:      s.Script.Name,
			Active:          s.Active,
		}
	}

	return scripts, nil
}

// CustomBrowserExport executes a custom browser export function to generate specialized data exports.
// This allows browser scripts to provide custom export functionality for operation data.
func (c *Client) CustomBrowserExport(ctx context.Context, req *types.CustomBrowserExportRequest) (*types.CustomBrowserExportResponse, error) {
	if err := c.EnsureAuthenticated(ctx); err != nil {
		return nil, err
	}

	if req == nil || req.OperationID == 0 || req.ScriptName == "" {
		return nil, WrapError("CustomBrowserExport", ErrInvalidInput, "operation ID and script name are required")
	}

	// Build the input object for the mutation
	input := map[string]interface{}{
		"operation_id": req.OperationID,
		"script_name":  req.ScriptName,
	}

	if len(req.Parameters) > 0 {
		input["parameters"] = req.Parameters
	}

	var mutation struct {
		CustomBrowserExport struct {
			Status       string `graphql:"status"`
			Error        string `graphql:"error"`
			ExportedData string `graphql:"exported_data"`
		} `graphql:"custombrowserExportFunction(input: $input)"`
	}

	variables := map[string]interface{}{
		"input": input,
	}

	err := c.executeMutation(ctx, &mutation, variables)
	if err != nil {
		return nil, WrapError("CustomBrowserExport", err, "failed to execute custom browser export")
	}

	if mutation.CustomBrowserExport.Status != "success" {
		return nil, WrapError("CustomBrowserExport", ErrOperationFailed, mutation.CustomBrowserExport.Error)
	}

	return &types.CustomBrowserExportResponse{
		Status:       mutation.CustomBrowserExport.Status,
		Error:        mutation.CustomBrowserExport.Error,
		ExportedData: mutation.CustomBrowserExport.ExportedData,
	}, nil
}
