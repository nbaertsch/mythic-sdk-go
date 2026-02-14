package mythic

import (
	"context"

	"github.com/nbaertsch/mythic-sdk-go/pkg/mythic/types"
)

// GetBuildParameters retrieves all build parameter type definitions.
// These define what parameters are available when building payloads.
func (c *Client) GetBuildParameters(ctx context.Context) ([]*types.BuildParameterType, error) {
	if err := c.EnsureAuthenticated(ctx); err != nil {
		return nil, err
	}

	var query struct {
		BuildParameters []struct {
			ID                   int           `graphql:"id"`
			Name                 string        `graphql:"name"`
			PayloadTypeID        int           `graphql:"payload_type_id"`
			Description          string        `graphql:"description"`
			Required             bool          `graphql:"required"`
			VerifierRegex        string        `graphql:"verifier_regex"`
			DefaultValue         string        `graphql:"default_value"`
			Randomize            bool          `graphql:"randomize"`
			FormatString         string        `graphql:"format_string"`
			ParameterType        string        `graphql:"parameter_type"`
			IsCryptoType         bool          `graphql:"crypto_type"`
			Deleted              bool          `graphql:"deleted"`
			GroupName            string        `graphql:"group_name"`
			Choices              []interface{} `graphql:"choices"`
			SupportedOS          []interface{} `graphql:"supported_os"`
			HideConditions       []interface{} `graphql:"hide_conditions"`
			DynamicQueryFunction string        `graphql:"dynamic_query_function"`
			UIPosition           int           `graphql:"ui_position"`
		} `graphql:"buildparameter(where: {deleted: {_eq: false}}, order_by: {payload_type_id: asc, name: asc})"`
	}

	err := c.executeQuery(ctx, &query, nil)
	if err != nil {
		return nil, WrapError("GetBuildParameters", err, "failed to query build parameters")
	}

	parameters := make([]*types.BuildParameterType, len(query.BuildParameters))
	for i, param := range query.BuildParameters {
		parameters[i] = &types.BuildParameterType{
			ID:                   param.ID,
			Name:                 param.Name,
			PayloadTypeID:        param.PayloadTypeID,
			Description:          param.Description,
			Required:             param.Required,
			VerifierRegex:        param.VerifierRegex,
			DefaultValue:         param.DefaultValue,
			Randomize:            param.Randomize,
			FormatString:         param.FormatString,
			ParameterType:        param.ParameterType,
			IsCryptoType:         param.IsCryptoType,
			Deleted:              param.Deleted,
			GroupName:            param.GroupName,
			Choices:              formatChoices(param.Choices),
			SupportedOS:          formatChoices(param.SupportedOS),
			HideConditions:       formatChoices(param.HideConditions),
			DynamicQueryFunction: param.DynamicQueryFunction,
			UIPosition:           param.UIPosition,
		}
	}

	return parameters, nil
}

// GetBuildParametersByPayloadType retrieves build parameters for a specific payload type.
func (c *Client) GetBuildParametersByPayloadType(ctx context.Context, payloadTypeID int) ([]*types.BuildParameterType, error) {
	if err := c.EnsureAuthenticated(ctx); err != nil {
		return nil, err
	}

	if payloadTypeID == 0 {
		return nil, WrapError("GetBuildParametersByPayloadType", ErrInvalidInput, "payload type ID is required")
	}

	var query struct {
		BuildParameters []struct {
			ID                   int           `graphql:"id"`
			Name                 string        `graphql:"name"`
			PayloadTypeID        int           `graphql:"payload_type_id"`
			Description          string        `graphql:"description"`
			Required             bool          `graphql:"required"`
			VerifierRegex        string        `graphql:"verifier_regex"`
			DefaultValue         string        `graphql:"default_value"`
			Randomize            bool          `graphql:"randomize"`
			FormatString         string        `graphql:"format_string"`
			ParameterType        string        `graphql:"parameter_type"`
			IsCryptoType         bool          `graphql:"crypto_type"`
			Deleted              bool          `graphql:"deleted"`
			GroupName            string        `graphql:"group_name"`
			Choices              []interface{} `graphql:"choices"`
			SupportedOS          []interface{} `graphql:"supported_os"`
			HideConditions       []interface{} `graphql:"hide_conditions"`
			DynamicQueryFunction string        `graphql:"dynamic_query_function"`
			UIPosition           int           `graphql:"ui_position"`
		} `graphql:"buildparameter(where: {payload_type_id: {_eq: $payload_type_id}, deleted: {_eq: false}}, order_by: {name: asc})"`
	}

	variables := map[string]interface{}{
		"payload_type_id": payloadTypeID,
	}

	err := c.executeQuery(ctx, &query, variables)
	if err != nil {
		return nil, WrapError("GetBuildParametersByPayloadType", err, "failed to query build parameters by payload type")
	}

	parameters := make([]*types.BuildParameterType, len(query.BuildParameters))
	for i, param := range query.BuildParameters {
		parameters[i] = &types.BuildParameterType{
			ID:                   param.ID,
			Name:                 param.Name,
			PayloadTypeID:        param.PayloadTypeID,
			Description:          param.Description,
			Required:             param.Required,
			VerifierRegex:        param.VerifierRegex,
			DefaultValue:         param.DefaultValue,
			Randomize:            param.Randomize,
			FormatString:         param.FormatString,
			ParameterType:        param.ParameterType,
			IsCryptoType:         param.IsCryptoType,
			Deleted:              param.Deleted,
			GroupName:            param.GroupName,
			Choices:              formatChoices(param.Choices),
			SupportedOS:          formatChoices(param.SupportedOS),
			HideConditions:       formatChoices(param.HideConditions),
			DynamicQueryFunction: param.DynamicQueryFunction,
			UIPosition:           param.UIPosition,
		}
	}

	return parameters, nil
}

// GetBuildParameterInstances retrieves all build parameter instances for all payloads in the current operation.
// These are the actual values set for build parameters when payloads were created.
func (c *Client) GetBuildParameterInstances(ctx context.Context) ([]*types.BuildParameterInstance, error) {
	if err := c.EnsureAuthenticated(ctx); err != nil {
		return nil, err
	}

	operationID := c.GetCurrentOperation()
	if operationID == nil {
		return nil, WrapError("GetBuildParameterInstances", ErrNotAuthenticated, "no current operation set")
	}

	var query struct {
		Instances []struct {
			ID               int     `graphql:"id"`
			PayloadID        int     `graphql:"payload_id"`
			BuildParameterID int     `graphql:"build_parameter_id"`
			Value            string  `graphql:"value"`
			EncValue         *string `graphql:"enc_value"`
			CreationTime     string  `graphql:"creation_time"`
		} `graphql:"buildparameterinstance(where: {payload: {operation_id: {_eq: $operation_id}}}, order_by: {payload_id: asc, build_parameter_id: asc})"`
	}

	variables := map[string]interface{}{
		"operation_id": *operationID,
	}

	err := c.executeQuery(ctx, &query, variables)
	if err != nil {
		return nil, WrapError("GetBuildParameterInstances", err, "failed to query build parameter instances")
	}

	instances := make([]*types.BuildParameterInstance, len(query.Instances))
	for i, inst := range query.Instances {
		creationTime, _ := parseTime(inst.CreationTime) //nolint:errcheck // Timestamp parse errors not critical

		instances[i] = &types.BuildParameterInstance{
			ID:               inst.ID,
			PayloadID:        inst.PayloadID,
			BuildParameterID: inst.BuildParameterID,
			Value:            inst.Value,
			EncValue:         inst.EncValue,
			CreationTime:     creationTime,
		}
	}

	return instances, nil
}

// GetBuildParameterInstancesByPayload retrieves build parameter instances for a specific payload.
func (c *Client) GetBuildParameterInstancesByPayload(ctx context.Context, payloadID int) ([]*types.BuildParameterInstance, error) {
	if err := c.EnsureAuthenticated(ctx); err != nil {
		return nil, err
	}

	if payloadID == 0 {
		return nil, WrapError("GetBuildParameterInstancesByPayload", ErrInvalidInput, "payload ID is required")
	}

	var query struct {
		Instances []struct {
			ID               int     `graphql:"id"`
			PayloadID        int     `graphql:"payload_id"`
			BuildParameterID int     `graphql:"build_parameter_id"`
			Value            string  `graphql:"value"`
			EncValue         *string `graphql:"enc_value"`
			CreationTime     string  `graphql:"creation_time"`
			BuildParameter   struct {
				ID            int    `graphql:"id"`
				Name          string `graphql:"name"`
				PayloadTypeID int    `graphql:"payload_type_id"`
				Description   string `graphql:"description"`
				ParameterType string `graphql:"parameter_type"`
				Required      bool   `graphql:"required"`
				DefaultValue  string `graphql:"default_value"`
				GroupName     string `graphql:"group_name"`
			} `graphql:"buildparameter"`
		} `graphql:"buildparameterinstance(where: {payload_id: {_eq: $payload_id}}, order_by: {build_parameter_id: asc})"`
	}

	variables := map[string]interface{}{
		"payload_id": payloadID,
	}

	err := c.executeQuery(ctx, &query, variables)
	if err != nil {
		return nil, WrapError("GetBuildParameterInstancesByPayload", err, "failed to query build parameter instances by payload")
	}

	instances := make([]*types.BuildParameterInstance, len(query.Instances))
	for i, inst := range query.Instances {
		creationTime, _ := parseTime(inst.CreationTime) //nolint:errcheck // Timestamp parse errors not critical

		instances[i] = &types.BuildParameterInstance{
			ID:               inst.ID,
			PayloadID:        inst.PayloadID,
			BuildParameterID: inst.BuildParameterID,
			Value:            inst.Value,
			EncValue:         inst.EncValue,
			CreationTime:     creationTime,
			BuildParameter: &types.BuildParameterType{
				ID:            inst.BuildParameter.ID,
				Name:          inst.BuildParameter.Name,
				PayloadTypeID: inst.BuildParameter.PayloadTypeID,
				Description:   inst.BuildParameter.Description,
				ParameterType: inst.BuildParameter.ParameterType,
				Required:      inst.BuildParameter.Required,
				DefaultValue:  inst.BuildParameter.DefaultValue,
				GroupName:     inst.BuildParameter.GroupName,
			},
		}
	}

	return instances, nil
}
