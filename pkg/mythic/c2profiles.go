package mythic

import (
	"context"
	"encoding/json"
	"fmt"

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
func (c *Client) CreateC2Instance(ctx context.Context, req *types.CreateC2InstanceRequest) error {
	if err := c.EnsureAuthenticated(ctx); err != nil {
		return err
	}

	if req == nil || req.InstanceName == "" || req.C2Instance == "" {
		return WrapError("CreateC2Instance", ErrInvalidInput, "instance_name and c2_instance are required")
	}

	if req.C2ProfileID <= 0 {
		return WrapError("CreateC2Instance", ErrInvalidInput, "c2profile_id must be positive")
	}

	var mutation struct {
		CreateC2Instance struct {
			Status string `graphql:"status"`
			Error  string `graphql:"error"`
		} `graphql:"create_c2_instance(c2_instance: $c2_instance, c2profile_id: $c2profile_id, instance_name: $instance_name)"`
	}

	variables := map[string]interface{}{
		"c2_instance":   req.C2Instance,
		"c2profile_id":  req.C2ProfileID,
		"instance_name": req.InstanceName,
	}

	err := c.executeMutation(ctx, &mutation, variables)
	if err != nil {
		return WrapError("CreateC2Instance", err, "failed to create C2 instance")
	}

	if mutation.CreateC2Instance.Status != "success" {
		return WrapError("CreateC2Instance", ErrOperationFailed, mutation.CreateC2Instance.Error)
	}

	return nil
}

// ImportC2Instance imports a C2 instance configuration.
// The mutation takes c2_instance (jsonb), c2profile_name (String), and instance_name (String).
func (c *Client) ImportC2Instance(ctx context.Context, req *types.ImportC2InstanceRequest) error {
	if err := c.EnsureAuthenticated(ctx); err != nil {
		return err
	}

	if req == nil || req.C2ProfileName == "" || req.InstanceName == "" || req.C2Instance == "" {
		return WrapError("ImportC2Instance", ErrInvalidInput, "c2profile_name, instance_name, and c2_instance are required")
	}

	// Parse the c2_instance JSON string into a map for jsonb param
	var c2InstanceData interface{}
	if err := json.Unmarshal([]byte(req.C2Instance), &c2InstanceData); err != nil {
		return WrapError("ImportC2Instance", ErrInvalidInput, fmt.Sprintf("c2_instance must be valid JSON: %v", err))
	}

	mutation := `mutation ImportC2Instance($c2_instance: jsonb!, $c2profile_name: String!, $instance_name: String!) {
		import_c2_instance(c2_instance: $c2_instance, c2profile_name: $c2profile_name, instance_name: $instance_name) {
			status
			error
		}
	}`

	variables := map[string]interface{}{
		"c2_instance":    c2InstanceData,
		"c2profile_name": req.C2ProfileName,
		"instance_name":  req.InstanceName,
	}

	result, err := c.ExecuteRawGraphQL(ctx, mutation, variables)
	if err != nil {
		return WrapError("ImportC2Instance", err, "failed to import C2 instance")
	}

	// Parse the response
	importResult, ok := result["import_c2_instance"].(map[string]interface{})
	if !ok {
		return WrapError("ImportC2Instance", ErrOperationFailed, "unexpected response format")
	}

	if status, _ := importResult["status"].(string); status != "success" {
		errMsg, _ := importResult["error"].(string)
		return WrapError("ImportC2Instance", ErrOperationFailed, errMsg)
	}

	return nil
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
func (c *Client) C2HostFile(ctx context.Context, c2ID int, fileUUID string, hostURL string) error {
	if err := c.EnsureAuthenticated(ctx); err != nil {
		return err
	}

	if c2ID == 0 || fileUUID == "" || hostURL == "" {
		return WrapError("C2HostFile", ErrInvalidInput, "C2 profile ID, file UUID, and host URL are required")
	}

	var mutation struct {
		C2HostFile struct {
			Status string `graphql:"status"`
			Error  string `graphql:"error"`
		} `graphql:"c2HostFile(c2_id: $c2_id, file_uuid: $file_uuid, host_url: $host_url)"`
	}

	variables := map[string]interface{}{
		"c2_id":     c2ID,
		"file_uuid": fileUUID,
		"host_url":  hostURL,
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
// The profileName parameter is the C2 profile's name (e.g., "http", "httpx"),
// which Mythic's c2SampleMessage query accepts as the "uuid" argument.
func (c *Client) C2SampleMessage(ctx context.Context, profileName string) (*types.C2SampleMessage, error) {
	if err := c.EnsureAuthenticated(ctx); err != nil {
		return nil, err
	}

	if profileName == "" {
		return nil, WrapError("C2SampleMessage", ErrInvalidInput, "profile name is required")
	}

	query := `query C2SampleMessage($uuid: String!) {
		c2SampleMessage(uuid: $uuid) {
			status
			error
			output
		}
	}`

	variables := map[string]interface{}{
		"uuid": profileName,
	}

	result, err := c.ExecuteRawGraphQL(ctx, query, variables)
	if err != nil {
		return nil, WrapError("C2SampleMessage", err, "failed to generate sample C2 message")
	}

	sampleResult, ok := result["c2SampleMessage"].(map[string]interface{})
	if !ok {
		return nil, WrapError("C2SampleMessage", ErrInvalidResponse, "unexpected response format")
	}

	if status, ok := sampleResult["status"].(string); ok && status != "success" {
		errMsg := ""
		if e, ok := sampleResult["error"].(string); ok {
			errMsg = e
		}
		return nil, WrapError("C2SampleMessage", ErrOperationFailed, errMsg)
	}

	output, _ := sampleResult["output"].(string)
	return &types.C2SampleMessage{
		Message: output,
	}, nil
}

// C2GetIOC retrieves indicators of compromise for a C2 profile.
// The profileName parameter is the C2 profile's name (e.g., "http", "httpx"),
// which Mythic's c2GetIOC query accepts as the "uuid" argument.
func (c *Client) C2GetIOC(ctx context.Context, profileName string) (*types.C2IOC, error) {
	if err := c.EnsureAuthenticated(ctx); err != nil {
		return nil, err
	}

	if profileName == "" {
		return nil, WrapError("C2GetIOC", ErrInvalidInput, "profile name is required")
	}

	query := `query C2GetIOC($uuid: String!) {
		c2GetIOC(uuid: $uuid) {
			status
			error
			output
		}
	}`

	variables := map[string]interface{}{
		"uuid": profileName,
	}

	result, err := c.ExecuteRawGraphQL(ctx, query, variables)
	if err != nil {
		return nil, WrapError("C2GetIOC", err, "failed to get C2 IOCs")
	}

	iocResult, ok := result["c2GetIOC"].(map[string]interface{})
	if !ok {
		return nil, WrapError("C2GetIOC", ErrInvalidResponse, "unexpected response format")
	}

	if status, ok := iocResult["status"].(string); ok && status != "success" {
		errMsg := ""
		if e, ok := iocResult["error"].(string); ok {
			errMsg = e
		}
		return nil, WrapError("C2GetIOC", ErrOperationFailed, errMsg)
	}

	output, _ := iocResult["output"].(string)
	return &types.C2IOC{
		ProfileName: profileName,
		Output:      output,
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
		C2ProfileParameters []struct {
			ID            int             `graphql:"id"`
			C2ProfileID   int             `graphql:"c2_profile_id"`
			Name          string          `graphql:"name"`
			Description   string          `graphql:"description"`
			DefaultValue  string          `graphql:"default_value"`
			ParameterType string          `graphql:"parameter_type"`
			Required      bool            `graphql:"required"`
			Randomize     bool            `graphql:"randomize"`
			FormatString  string          `graphql:"format_string"`
			VerifierRegex string          `graphql:"verifier_regex"`
			IsCryptoType  bool            `graphql:"crypto_type"`
			Deleted       bool            `graphql:"deleted"`
			Choices       json.RawMessage `graphql:"choices"`
		} `graphql:"c2profileparameters(where: {c2_profile_id: {_eq: $profile_id}, deleted: {_eq: false}}, order_by: {name: asc})"`
	}

	variables := map[string]interface{}{
		"profile_id": profileID,
	}

	err := c.executeQuery(ctx, &query, variables)
	if err != nil {
		return nil, WrapError("GetC2ProfileParameters", err, "failed to query C2 profile parameters")
	}

	parameters := make([]*types.C2ProfileParameter, len(query.C2ProfileParameters))
	for i, p := range query.C2ProfileParameters {
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
			Choices:       formatRawJSON(p.Choices),
		}
	}

	return parameters, nil
}
