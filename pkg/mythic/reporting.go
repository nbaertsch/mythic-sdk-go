package mythic

import (
	"context"

	"github.com/nbaertsch/mythic-sdk-go/pkg/mythic/types"
)

// GenerateReport generates an operation report with optional MITRE coverage and filters.
func (c *Client) GenerateReport(ctx context.Context, req *types.GenerateReportRequest) (*types.GenerateReportResponse, error) {
	if err := c.EnsureAuthenticated(ctx); err != nil {
		return nil, err
	}

	if req == nil || req.OperationID == 0 {
		return nil, WrapError("GenerateReport", ErrInvalidInput, "operation ID is required")
	}

	// Build the input object for the mutation
	input := map[string]interface{}{
		"operation_id": req.OperationID,
	}

	if req.OutputFormat != "" {
		input["output_format"] = req.OutputFormat
	}
	if req.IncludeMITRE {
		input["include_mitre"] = req.IncludeMITRE
	}
	if req.IncludeCallbacks {
		input["include_callbacks"] = req.IncludeCallbacks
	}
	if req.IncludeTasks {
		input["include_tasks"] = req.IncludeTasks
	}
	if req.IncludeFiles {
		input["include_files"] = req.IncludeFiles
	}
	if req.IncludeCredentials {
		input["include_credentials"] = req.IncludeCredentials
	}
	if req.IncludeArtifacts {
		input["include_artifacts"] = req.IncludeArtifacts
	}
	if req.StartDate != nil {
		input["start_date"] = *req.StartDate
	}
	if req.EndDate != nil {
		input["end_date"] = *req.EndDate
	}
	if len(req.CallbackIDs) > 0 {
		input["callback_ids"] = req.CallbackIDs
	}

	var mutation struct {
		GenerateReport struct {
			Status     string `graphql:"status"`
			Error      string `graphql:"error"`
			ReportData string `graphql:"report_data"`
		} `graphql:"generateReport(input: $input)"`
	}

	variables := map[string]interface{}{
		"input": input,
	}

	err := c.executeMutation(ctx, &mutation, variables)
	if err != nil {
		return nil, WrapError("GenerateReport", err, "failed to generate report")
	}

	if mutation.GenerateReport.Status != "success" {
		return nil, WrapError("GenerateReport", ErrOperationFailed, mutation.GenerateReport.Error)
	}

	return &types.GenerateReportResponse{
		Status:     mutation.GenerateReport.Status,
		Error:      mutation.GenerateReport.Error,
		ReportData: mutation.GenerateReport.ReportData,
	}, nil
}

// GetRedirectRules retrieves C2 redirect rules for payloads.
func (c *Client) GetRedirectRules(ctx context.Context, payloadUUID string) ([]*types.RedirectRule, error) {
	if err := c.EnsureAuthenticated(ctx); err != nil {
		return nil, err
	}

	if payloadUUID == "" {
		return nil, WrapError("GetRedirectRules", ErrInvalidInput, "payload UUID is required")
	}

	var query struct {
		RedirectRules struct {
			Status string `graphql:"status"`
			Error  string `graphql:"error"`
			Rules  []struct {
				PayloadType  string `graphql:"payload_type"`
				C2Profile    string `graphql:"c2_profile"`
				RedirectType string `graphql:"redirect_type"`
				Rule         string `graphql:"rule"`
				Description  string `graphql:"description"`
			} `graphql:"rules"`
		} `graphql:"redirect_rules(uuid: $uuid)"`
	}

	variables := map[string]interface{}{
		"uuid": payloadUUID,
	}

	err := c.executeQuery(ctx, &query, variables)
	if err != nil {
		return nil, WrapError("GetRedirectRules", err, "failed to query redirect rules")
	}

	if query.RedirectRules.Status != "success" {
		return nil, WrapError("GetRedirectRules", ErrOperationFailed, query.RedirectRules.Error)
	}

	rules := make([]*types.RedirectRule, len(query.RedirectRules.Rules))
	for i, r := range query.RedirectRules.Rules {
		rules[i] = &types.RedirectRule{
			ID:           i + 1, // Generate sequential ID
			PayloadType:  r.PayloadType,
			C2Profile:    r.C2Profile,
			RedirectType: r.RedirectType,
			Rule:         r.Rule,
			Description:  r.Description,
		}
	}

	return rules, nil
}
