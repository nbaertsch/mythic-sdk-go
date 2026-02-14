package mythic

import (
	"context"

	"github.com/nbaertsch/mythic-sdk-go/pkg/mythic/types"
)

// GetC2Profiles retrieves all C2 profiles.
func (c *Client) GetC2Profiles(ctx context.Context) ([]*types.C2Profile, error) {
	if err := c.EnsureAuthenticated(ctx); err != nil {
		return nil, err
	}

	var query struct {
		C2Profile []struct {
			ID           int    `graphql:"id"`
			Name         string `graphql:"name"`
			Description  string `graphql:"description"`
			CreationTime string `graphql:"creation_time"`
			Running      bool   `graphql:"running"`
			Deleted      bool   `graphql:"deleted"`
			IsP2P        bool   `graphql:"is_p2p"`
		} `graphql:"c2profile(where: {deleted: {_eq: false}}, order_by: {name: asc})"`
	}

	err := c.executeQuery(ctx, &query, nil)
	if err != nil {
		return nil, WrapError("GetC2Profiles", err, "failed to query C2 profiles")
	}

	profiles := make([]*types.C2Profile, len(query.C2Profile))
	for i, p := range query.C2Profile {
		creationTime, _ := parseTime(p.CreationTime) //nolint:errcheck // Timestamp parse errors not critical
		profiles[i] = &types.C2Profile{
			ID:           p.ID,
			Name:         p.Name,
			Description:  p.Description,
			CreationTime: creationTime,
			Running:      p.Running,
			Deleted:      p.Deleted,
			IsP2P:        p.IsP2P,
		}
	}

	return profiles, nil
}

// GetC2ProfileByID retrieves a specific C2 profile by ID.
func (c *Client) GetC2ProfileByID(ctx context.Context, profileID int) (*types.C2Profile, error) {
	if err := c.EnsureAuthenticated(ctx); err != nil {
		return nil, err
	}

	if profileID == 0 {
		return nil, WrapError("GetC2ProfileByID", ErrInvalidInput, "profile ID is required")
	}

	var query struct {
		C2Profile []struct {
			ID           int    `graphql:"id"`
			Name         string `graphql:"name"`
			Description  string `graphql:"description"`
			CreationTime string `graphql:"creation_time"`
			Running      bool   `graphql:"running"`
			Deleted      bool   `graphql:"deleted"`
			IsP2P        bool   `graphql:"is_p2p"`
		} `graphql:"c2profile(where: {id: {_eq: $profile_id}})"`
	}

	variables := map[string]interface{}{
		"profile_id": profileID,
	}

	err := c.executeQuery(ctx, &query, variables)
	if err != nil {
		return nil, WrapError("GetC2ProfileByID", err, "failed to query C2 profile")
	}

	if len(query.C2Profile) == 0 {
		return nil, WrapError("GetC2ProfileByID", ErrNotFound, "C2 profile not found")
	}

	p := query.C2Profile[0]
	creationTime, _ := parseTime(p.CreationTime) //nolint:errcheck // Timestamp parse errors not critical
	return &types.C2Profile{
		ID:           p.ID,
		Name:         p.Name,
		Description:  p.Description,
		CreationTime: creationTime,
		Running:      p.Running,
		Deleted:      p.Deleted,
		IsP2P:        p.IsP2P,
	}, nil
}

// CreateC2Instance creates a new C2 profile instance.
func (c *Client) CreateC2Instance(ctx context.Context, req *types.CreateC2InstanceRequest) (*types.C2Profile, error) {
	if err := c.EnsureAuthenticated(ctx); err != nil {
		return nil, err
	}

	if req == nil || req.Name == "" {
		return nil, WrapError("CreateC2Instance", ErrInvalidInput, "request and profile name are required")
	}

	// Use current operation if not specified
	operationID := c.GetCurrentOperation()
	if req.OperationID != nil {
		operationID = req.OperationID
	}
	if operationID == nil {
		return nil, WrapError("CreateC2Instance", ErrInvalidInput, "operation ID is required")
	}

	var mutation struct {
		CreateC2Instance struct {
			Status string `graphql:"status"`
			Error  string `graphql:"error"`
			ID     int    `graphql:"id"`
		} `graphql:"create_c2_instance(name: $name, description: $description, operation_id: $operation_id, parameters: $parameters)"`
	}

	description := ""
	if req.Description != nil {
		description = *req.Description
	}

	variables := map[string]interface{}{
		"name":         req.Name,
		"description":  description,
		"operation_id": *operationID,
		"parameters":   req.Parameters,
	}

	err := c.executeMutation(ctx, &mutation, variables)
	if err != nil {
		return nil, WrapError("CreateC2Instance", err, "failed to create C2 instance")
	}

	if mutation.CreateC2Instance.Status != "success" {
		return nil, WrapError("CreateC2Instance", ErrOperationFailed, mutation.CreateC2Instance.Error)
	}

	// Fetch the created profile
	return c.GetC2ProfileByID(ctx, mutation.CreateC2Instance.ID)
}

// ImportC2Instance imports a C2 instance configuration.
func (c *Client) ImportC2Instance(ctx context.Context, req *types.ImportC2InstanceRequest) (*types.C2Profile, error) {
	if err := c.EnsureAuthenticated(ctx); err != nil {
		return nil, err
	}

	if req == nil || req.Config == "" || req.Name == "" {
		return nil, WrapError("ImportC2Instance", ErrInvalidInput, "config and name are required")
	}

	var mutation struct {
		ImportC2Instance struct {
			Status string `graphql:"status"`
			Error  string `graphql:"error"`
			ID     int    `graphql:"id"`
		} `graphql:"import_c2_instance(name: $name, config: $config)"`
	}

	variables := map[string]interface{}{
		"name":   req.Name,
		"config": req.Config,
	}

	err := c.executeMutation(ctx, &mutation, variables)
	if err != nil {
		return nil, WrapError("ImportC2Instance", err, "failed to import C2 instance")
	}

	if mutation.ImportC2Instance.Status != "success" {
		return nil, WrapError("ImportC2Instance", ErrOperationFailed, mutation.ImportC2Instance.Error)
	}

	// Fetch the imported profile
	return c.GetC2ProfileByID(ctx, mutation.ImportC2Instance.ID)
}

// StartStopProfile starts or stops a C2 profile.
func (c *Client) StartStopProfile(ctx context.Context, profileID int, start bool) error {
	if err := c.EnsureAuthenticated(ctx); err != nil {
		return err
	}

	if profileID == 0 {
		return WrapError("StartStopProfile", ErrInvalidInput, "profile ID is required")
	}

	var mutation struct {
		StartStopProfile struct {
			Status string `graphql:"status"`
			Error  string `graphql:"error"`
		} `graphql:"startStopProfile(id: $id, action: $action)"`
	}

	action := "stop"
	if start {
		action = "start"
	}

	variables := map[string]interface{}{
		"id":     profileID,
		"action": action,
	}

	err := c.executeMutation(ctx, &mutation, variables)
	if err != nil {
		return WrapError("StartStopProfile", err, "failed to start/stop C2 profile")
	}

	if mutation.StartStopProfile.Status != "success" {
		return WrapError("StartStopProfile", ErrOperationFailed, mutation.StartStopProfile.Error)
	}

	return nil
}

// GetProfileOutput retrieves the output/logs from a C2 profile.
func (c *Client) GetProfileOutput(ctx context.Context, profileID int) (*types.C2ProfileOutput, error) {
	if err := c.EnsureAuthenticated(ctx); err != nil {
		return nil, err
	}

	if profileID == 0 {
		return nil, WrapError("GetProfileOutput", ErrInvalidInput, "profile ID is required")
	}

	var query struct {
		GetProfileOutput struct {
			Status string `graphql:"status"`
			Error  string `graphql:"error"`
			Output string `graphql:"output"`
		} `graphql:"getProfileOutput(id: $id)"`
	}

	variables := map[string]interface{}{
		"id": profileID,
	}

	err := c.executeQuery(ctx, &query, variables)
	if err != nil {
		return nil, WrapError("GetProfileOutput", err, "failed to get profile output")
	}

	if query.GetProfileOutput.Status != "success" {
		return nil, WrapError("GetProfileOutput", ErrOperationFailed, query.GetProfileOutput.Error)
	}

	return &types.C2ProfileOutput{
		Output: query.GetProfileOutput.Output,
		StdOut: "",
		StdErr: "",
	}, nil
}

// C2HostFile hosts a file via a C2 profile.
func (c *Client) C2HostFile(ctx context.Context, profileID int, fileUUID string) error {
	if err := c.EnsureAuthenticated(ctx); err != nil {
		return err
	}

	if profileID == 0 || fileUUID == "" {
		return WrapError("C2HostFile", ErrInvalidInput, "profile ID and file UUID are required")
	}

	var mutation struct {
		C2HostFile struct {
			Status string `graphql:"status"`
			Error  string `graphql:"error"`
		} `graphql:"c2HostFile(profile_id: $profile_id, file_uuid: $file_uuid)"`
	}

	variables := map[string]interface{}{
		"profile_id": profileID,
		"file_uuid":  fileUUID,
	}

	err := c.executeMutation(ctx, &mutation, variables)
	if err != nil {
		return WrapError("C2HostFile", err, "failed to host file via C2 profile")
	}

	if mutation.C2HostFile.Status != "success" {
		return WrapError("C2HostFile", ErrOperationFailed, mutation.C2HostFile.Error)
	}

	return nil
}

// C2SampleMessage generates a sample C2 message for testing.
func (c *Client) C2SampleMessage(ctx context.Context, profileID int, messageType string) (*types.C2SampleMessage, error) {
	if err := c.EnsureAuthenticated(ctx); err != nil {
		return nil, err
	}

	if profileID == 0 {
		return nil, WrapError("C2SampleMessage", ErrInvalidInput, "profile ID is required")
	}

	var query struct {
		C2SampleMessage struct {
			Status  string `graphql:"status"`
			Error   string `graphql:"error"`
			Message string `graphql:"message"`
		} `graphql:"c2SampleMessage(profile_id: $profile_id, message_type: $message_type)"`
	}

	variables := map[string]interface{}{
		"profile_id":   profileID,
		"message_type": messageType,
	}

	err := c.executeQuery(ctx, &query, variables)
	if err != nil {
		return nil, WrapError("C2SampleMessage", err, "failed to generate sample C2 message")
	}

	if query.C2SampleMessage.Status != "success" {
		return nil, WrapError("C2SampleMessage", ErrOperationFailed, query.C2SampleMessage.Error)
	}

	return &types.C2SampleMessage{
		Message: query.C2SampleMessage.Message,
	}, nil
}

// C2GetIOC retrieves indicators of compromise for a C2 profile.
func (c *Client) C2GetIOC(ctx context.Context, profileID int) (*types.C2IOC, error) {
	if err := c.EnsureAuthenticated(ctx); err != nil {
		return nil, err
	}

	if profileID == 0 {
		return nil, WrapError("C2GetIOC", ErrInvalidInput, "profile ID is required")
	}

	var query struct {
		C2GetIOC struct {
			Status string   `graphql:"status"`
			Error  string   `graphql:"error"`
			IOCs   []string `graphql:"iocs"`
		} `graphql:"c2GetIOC(profile_id: $profile_id)"`
	}

	variables := map[string]interface{}{
		"profile_id": profileID,
	}

	err := c.executeQuery(ctx, &query, variables)
	if err != nil {
		return nil, WrapError("C2GetIOC", err, "failed to get C2 IOCs")
	}

	if query.C2GetIOC.Status != "success" {
		return nil, WrapError("C2GetIOC", ErrOperationFailed, query.C2GetIOC.Error)
	}

	return &types.C2IOC{
		ProfileID: profileID,
		IOCs:      query.C2GetIOC.IOCs,
	}, nil
}

// GetC2ProfileParameters retrieves all configuration parameters for a specific C2 profile.
// These define what options are available when configuring this C2 profile for a payload
// (e.g., callback_host, callback_port, callback_interval, etc.).
func (c *Client) GetC2ProfileParameters(ctx context.Context, profileID int) ([]*types.C2ProfileParameter, error) {
	if err := c.EnsureAuthenticated(ctx); err != nil {
		return nil, err
	}

	if profileID == 0 {
		return nil, WrapError("GetC2ProfileParameters", ErrInvalidInput, "profile ID is required")
	}

	var query struct {
		Parameters []struct {
			ID            int    `graphql:"id"`
			C2ProfileID   int    `graphql:"c2_profile_id"`
			Name          string `graphql:"name"`
			Description   string `graphql:"description"`
			DefaultValue  string `graphql:"default_value"`
			ParameterType string `graphql:"parameter_type"`
			Required      bool   `graphql:"required"`
			Randomize     bool   `graphql:"randomize"`
			FormatString  string `graphql:"format_string"`
			VerifierRegex string `graphql:"verifier_regex"`
			IsCryptoType  bool   `graphql:"crypto_type"`
			Deleted       bool   `graphql:"deleted"`
			Choices       string `graphql:"choices"`
		} `graphql:"c2profileparameters(where: {c2_profile_id: {_eq: $profile_id}, deleted: {_eq: false}}, order_by: {name: asc})"`
	}

	variables := map[string]interface{}{
		"profile_id": profileID,
	}

	err := c.executeQuery(ctx, &query, variables)
	if err != nil {
		return nil, WrapError("GetC2ProfileParameters", err, "failed to query C2 profile parameters")
	}

	parameters := make([]*types.C2ProfileParameter, len(query.Parameters))
	for i, p := range query.Parameters {
		parameters[i] = &types.C2ProfileParameter{
			ID:            p.ID,
			C2ProfileID:   p.C2ProfileID,
			Name:          p.Name,
			Description:   p.Description,
			DefaultValue:  p.DefaultValue,
			ParameterType: p.ParameterType,
			Required:      p.Required,
			Randomize:     p.Randomize,
			FormatString:  p.FormatString,
			VerifierRegex: p.VerifierRegex,
			IsCryptoType:  p.IsCryptoType,
			Deleted:       p.Deleted,
			Choices:       p.Choices,
		}
	}

	return parameters, nil
}
