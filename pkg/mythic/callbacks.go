package mythic

import (
	"context"
	"fmt"
	"time"

	"github.com/nbaertsch/mythic-sdk-go/pkg/mythic/types"
)

// GetAllCallbacks retrieves all callbacks (active and inactive).
func (c *Client) GetAllCallbacks(ctx context.Context) ([]*types.Callback, error) {
	if err := c.EnsureAuthenticated(ctx); err != nil {
		return nil, err
	}

	var query struct {
		Callback []struct {
			ID              int       `graphql:"id"`
			DisplayID       int       `graphql:"display_id"`
			AgentCallbackID string    `graphql:"agent_callback_id"`
			InitCallback    time.Time `graphql:"init_callback"`
			LastCheckin     time.Time `graphql:"last_checkin"`
			User            string    `graphql:"user"`
			Host            string    `graphql:"host"`
			PID             int       `graphql:"pid"`
			IP              string    `graphql:"ip"`
			ExternalIP      string    `graphql:"external_ip"`
			ProcessName     string    `graphql:"process_name"`
			Description     string    `graphql:"description"`
			Active          bool      `graphql:"active"`
			IntegrityLevel  int       `graphql:"integrity_level"`
			Locked          bool      `graphql:"locked"`
			OS              string    `graphql:"os"`
			Architecture    string    `graphql:"architecture"`
			Domain          string    `graphql:"domain"`
			ExtraInfo       string    `graphql:"extra_info"`
			SleepInfo       string    `graphql:"sleep_info"`
			OperationID     int       `graphql:"operation_id"`
			OperatorID      int       `graphql:"operator_id"`
			Payload         struct {
				ID          int    `graphql:"id"`
				UUID        string `graphql:"uuid"`
				Description string `graphql:"description"`
				OS          string `graphql:"os"`
				PayloadType struct {
					ID   int    `graphql:"id"`
					Name string `graphql:"name"`
				} `graphql:"payloadtype"`
			} `graphql:"payload"`
			Operator struct {
				ID       int    `graphql:"id"`
				Username string `graphql:"username"`
			} `graphql:"operator"`
		} `graphql:"callback(order_by: {id: desc})"`
	}

	err := c.executeQuery(ctx, &query, nil)
	if err != nil {
		return nil, WrapError("GetAllCallbacks", err, "failed to query callbacks")
	}

	callbacks := make([]*types.Callback, len(query.Callback))
	for i, cb := range query.Callback {
		callbacks[i] = &types.Callback{
			ID:              cb.ID,
			DisplayID:       cb.DisplayID,
			AgentCallbackID: cb.AgentCallbackID,
			InitCallback:    cb.InitCallback,
			LastCheckin:     cb.LastCheckin,
			User:            cb.User,
			Host:            cb.Host,
			PID:             cb.PID,
			IP:              parseIPString(cb.IP),
			ExternalIP:      cb.ExternalIP,
			ProcessName:     cb.ProcessName,
			Description:     cb.Description,
			Active:          cb.Active,
			IntegrityLevel:  types.CallbackIntegrityLevel(cb.IntegrityLevel),
			Locked:          cb.Locked,
			OS:              cb.OS,
			Architecture:    cb.Architecture,
			Domain:          cb.Domain,
			ExtraInfo:       cb.ExtraInfo,
			SleepInfo:       cb.SleepInfo,
			OperationID:     cb.OperationID,
			PayloadTypeID:   cb.Payload.PayloadType.ID,
			OperatorID:      cb.OperatorID,
			Payload: &types.CallbackPayload{
				ID:          cb.Payload.ID,
				UUID:        cb.Payload.UUID,
				Description: cb.Payload.Description,
				OS:          cb.Payload.OS,
			},
			PayloadType: &types.CallbackPayloadType{
				ID:   cb.Payload.PayloadType.ID,
				Name: cb.Payload.PayloadType.Name,
			},
			Operator: &types.CallbackOperator{
				ID:       cb.Operator.ID,
				Username: cb.Operator.Username,
			},
		}
	}

	return callbacks, nil
}

// GetAllActiveCallbacks retrieves only currently active callbacks.
func (c *Client) GetAllActiveCallbacks(ctx context.Context) ([]*types.Callback, error) {
	if err := c.EnsureAuthenticated(ctx); err != nil {
		return nil, err
	}

	var query struct {
		Callback []struct {
			ID              int       `graphql:"id"`
			DisplayID       int       `graphql:"display_id"`
			AgentCallbackID string    `graphql:"agent_callback_id"`
			InitCallback    time.Time `graphql:"init_callback"`
			LastCheckin     time.Time `graphql:"last_checkin"`
			User            string    `graphql:"user"`
			Host            string    `graphql:"host"`
			PID             int       `graphql:"pid"`
			IP              string    `graphql:"ip"`
			ExternalIP      string    `graphql:"external_ip"`
			ProcessName     string    `graphql:"process_name"`
			Description     string    `graphql:"description"`
			Active          bool      `graphql:"active"`
			IntegrityLevel  int       `graphql:"integrity_level"`
			Locked          bool      `graphql:"locked"`
			OS              string    `graphql:"os"`
			Architecture    string    `graphql:"architecture"`
			Domain          string    `graphql:"domain"`
			ExtraInfo       string    `graphql:"extra_info"`
			SleepInfo       string    `graphql:"sleep_info"`
			OperationID     int       `graphql:"operation_id"`
			OperatorID      int       `graphql:"operator_id"`
			Payload         struct {
				ID          int    `graphql:"id"`
				UUID        string `graphql:"uuid"`
				Description string `graphql:"description"`
				OS          string `graphql:"os"`
				PayloadType struct {
					ID   int    `graphql:"id"`
					Name string `graphql:"name"`
				} `graphql:"payloadtype"`
			} `graphql:"payload"`
			Operator struct {
				ID       int    `graphql:"id"`
				Username string `graphql:"username"`
			} `graphql:"operator"`
		} `graphql:"callback(where: {active: {_eq: true}}, order_by: {id: desc})"`
	}

	err := c.executeQuery(ctx, &query, nil)
	if err != nil {
		return nil, WrapError("GetAllActiveCallbacks", err, "failed to query active callbacks")
	}

	callbacks := make([]*types.Callback, len(query.Callback))
	for i, cb := range query.Callback {
		callbacks[i] = &types.Callback{
			ID:              cb.ID,
			DisplayID:       cb.DisplayID,
			AgentCallbackID: cb.AgentCallbackID,
			InitCallback:    cb.InitCallback,
			LastCheckin:     cb.LastCheckin,
			User:            cb.User,
			Host:            cb.Host,
			PID:             cb.PID,
			IP:              parseIPString(cb.IP),
			ExternalIP:      cb.ExternalIP,
			ProcessName:     cb.ProcessName,
			Description:     cb.Description,
			Active:          cb.Active,
			IntegrityLevel:  types.CallbackIntegrityLevel(cb.IntegrityLevel),
			Locked:          cb.Locked,
			OS:              cb.OS,
			Architecture:    cb.Architecture,
			Domain:          cb.Domain,
			ExtraInfo:       cb.ExtraInfo,
			SleepInfo:       cb.SleepInfo,
			OperationID:     cb.OperationID,
			PayloadTypeID:   cb.Payload.PayloadType.ID,
			OperatorID:      cb.OperatorID,
			Payload: &types.CallbackPayload{
				ID:          cb.Payload.ID,
				UUID:        cb.Payload.UUID,
				Description: cb.Payload.Description,
				OS:          cb.Payload.OS,
			},
			PayloadType: &types.CallbackPayloadType{
				ID:   cb.Payload.PayloadType.ID,
				Name: cb.Payload.PayloadType.Name,
			},
			Operator: &types.CallbackOperator{
				ID:       cb.Operator.ID,
				Username: cb.Operator.Username,
			},
		}
	}

	return callbacks, nil
}

// GetCallbackByID retrieves a specific callback by its display ID.
func (c *Client) GetCallbackByID(ctx context.Context, displayID int) (*types.Callback, error) {
	if err := c.EnsureAuthenticated(ctx); err != nil {
		return nil, err
	}

	var query struct {
		Callback []struct {
			ID              int       `graphql:"id"`
			DisplayID       int       `graphql:"display_id"`
			AgentCallbackID string    `graphql:"agent_callback_id"`
			InitCallback    time.Time `graphql:"init_callback"`
			LastCheckin     time.Time `graphql:"last_checkin"`
			User            string    `graphql:"user"`
			Host            string    `graphql:"host"`
			PID             int       `graphql:"pid"`
			IP              string    `graphql:"ip"`
			ExternalIP      string    `graphql:"external_ip"`
			ProcessName     string    `graphql:"process_name"`
			Description     string    `graphql:"description"`
			Active          bool      `graphql:"active"`
			IntegrityLevel  int       `graphql:"integrity_level"`
			Locked          bool      `graphql:"locked"`
			OS              string    `graphql:"os"`
			Architecture    string    `graphql:"architecture"`
			Domain          string    `graphql:"domain"`
			ExtraInfo       string    `graphql:"extra_info"`
			SleepInfo       string    `graphql:"sleep_info"`
			OperationID     int       `graphql:"operation_id"`
			OperatorID      int       `graphql:"operator_id"`
			Payload         struct {
				ID          int    `graphql:"id"`
				UUID        string `graphql:"uuid"`
				Description string `graphql:"description"`
				OS          string `graphql:"os"`
				PayloadType struct {
					ID   int    `graphql:"id"`
					Name string `graphql:"name"`
				} `graphql:"payloadtype"`
			} `graphql:"payload"`
			Operator struct {
				ID       int    `graphql:"id"`
				Username string `graphql:"username"`
			} `graphql:"operator"`
		} `graphql:"callback(where: {display_id: {_eq: $displayID}}, limit: 1)"`
	}

	variables := map[string]interface{}{
		"displayID": displayID,
	}

	err := c.executeQuery(ctx, &query, variables)
	if err != nil {
		return nil, WrapError("GetCallbackByID", err, "failed to query callback")
	}

	if len(query.Callback) == 0 {
		return nil, WrapError("GetCallbackByID", ErrNotFound, fmt.Sprintf("callback with display_id %d not found", displayID))
	}

	cb := query.Callback[0]
	callback := &types.Callback{
		ID:              cb.ID,
		DisplayID:       cb.DisplayID,
		AgentCallbackID: cb.AgentCallbackID,
		InitCallback:    cb.InitCallback,
		LastCheckin:     cb.LastCheckin,
		User:            cb.User,
		Host:            cb.Host,
		PID:             cb.PID,
		IP:              parseIPString(cb.IP),
		ExternalIP:      cb.ExternalIP,
		ProcessName:     cb.ProcessName,
		Description:     cb.Description,
		Active:          cb.Active,
		IntegrityLevel:  types.CallbackIntegrityLevel(cb.IntegrityLevel),
		Locked:          cb.Locked,
		OS:              cb.OS,
		Architecture:    cb.Architecture,
		Domain:          cb.Domain,
		ExtraInfo:       cb.ExtraInfo,
		SleepInfo:       cb.SleepInfo,
		OperationID:     cb.OperationID,
		PayloadTypeID:   cb.Payload.PayloadType.ID,
		OperatorID:      cb.OperatorID,
		Payload: &types.CallbackPayload{
			ID:          cb.Payload.ID,
			UUID:        cb.Payload.UUID,
			Description: cb.Payload.Description,
			OS:          cb.Payload.OS,
		},
		PayloadType: &types.CallbackPayloadType{
			ID:   cb.Payload.PayloadType.ID,
			Name: cb.Payload.PayloadType.Name,
		},
		Operator: &types.CallbackOperator{
			ID:       cb.Operator.ID,
			Username: cb.Operator.Username,
		},
	}

	return callback, nil
}

// UpdateCallback updates properties of a callback.
func (c *Client) UpdateCallback(ctx context.Context, req *types.CallbackUpdateRequest) error {
	if err := c.EnsureAuthenticated(ctx); err != nil {
		return err
	}

	if req.CallbackDisplayID <= 0 {
		return WrapError("UpdateCallback", ErrInvalidConfig, "callback display ID is required")
	}

	// Build the set clause dynamically based on what fields are provided
	setClause := make(map[string]interface{})

	if req.Active != nil {
		setClause["active"] = *req.Active
	}
	if req.Locked != nil {
		setClause["locked"] = *req.Locked
	}
	if req.Description != nil {
		setClause["description"] = *req.Description
	}
	if req.IPs != nil {
		setClause["ip"] = formatIPString(req.IPs)
	}
	if req.User != nil {
		setClause["user"] = *req.User
	}
	if req.Host != nil {
		setClause["host"] = *req.Host
	}
	if req.OS != nil {
		setClause["os"] = *req.OS
	}
	if req.Architecture != nil {
		setClause["architecture"] = *req.Architecture
	}
	if req.ExtraInfo != nil {
		setClause["extra_info"] = *req.ExtraInfo
	}
	if req.SleepInfo != nil {
		setClause["sleep_info"] = *req.SleepInfo
	}
	if req.PID != nil {
		setClause["pid"] = *req.PID
	}
	if req.ProcessName != nil {
		setClause["process_name"] = *req.ProcessName
	}
	if req.IntegrityLevel != nil {
		setClause["integrity_level"] = int(*req.IntegrityLevel)
	}
	if req.Domain != nil {
		setClause["domain"] = *req.Domain
	}

	if len(setClause) == 0 {
		return WrapError("UpdateCallback", ErrInvalidConfig, "no fields to update")
	}

	var mutation struct {
		UpdateCallback struct {
			Returning []struct {
				ID int `graphql:"id"`
			} `graphql:"returning"`
		} `graphql:"update_callback(where: {display_id: {_eq: $displayID}}, _set: $set)"`
	}

	variables := map[string]interface{}{
		"displayID": req.CallbackDisplayID,
		"set":       setClause,
	}

	err := c.executeMutation(ctx, &mutation, variables)
	if err != nil {
		return WrapError("UpdateCallback", err, "failed to update callback")
	}

	if len(mutation.UpdateCallback.Returning) == 0 {
		return WrapError("UpdateCallback", ErrNotFound, fmt.Sprintf("callback with display_id %d not found", req.CallbackDisplayID))
	}

	return nil
}

// CreateCallbackInput represents the input for manually creating a callback.
type CreateCallbackInput struct {
	PayloadUUID string  `json:"payloadUuid"`
	IP          *string `json:"ip,omitempty"`
	ExternalIP  *string `json:"externalIp,omitempty"`
	User        *string `json:"user,omitempty"`
	Host        *string `json:"host,omitempty"`
	Domain      *string `json:"domain,omitempty"`
	Description *string `json:"description,omitempty"`
	ProcessName *string `json:"processName,omitempty"`
	SleepInfo   *string `json:"sleepInfo,omitempty"`
	ExtraInfo   *string `json:"extraInfo,omitempty"`
}

// CreateCallback manually registers a new callback session.
func (c *Client) CreateCallback(ctx context.Context, input *CreateCallbackInput) error {
	if err := c.EnsureAuthenticated(ctx); err != nil {
		return err
	}

	if input == nil || input.PayloadUUID == "" {
		return WrapError("CreateCallback", ErrInvalidInput, "payload UUID is required")
	}

	var mutation struct {
		CreateCallback struct {
			Status string `graphql:"status"`
			Error  string `graphql:"error"`
		} `graphql:"createCallback(payloadUuid: $payloadUuid, newCallback: $newCallback)"`
	}

	// Build newCallback input
	newCallback := make(map[string]interface{})
	if input.IP != nil {
		newCallback["ip"] = *input.IP
	}
	if input.ExternalIP != nil {
		newCallback["externalIp"] = *input.ExternalIP
	}
	if input.User != nil {
		newCallback["user"] = *input.User
	}
	if input.Host != nil {
		newCallback["host"] = *input.Host
	}
	if input.Domain != nil {
		newCallback["domain"] = *input.Domain
	}
	if input.Description != nil {
		newCallback["description"] = *input.Description
	}
	if input.ProcessName != nil {
		newCallback["processName"] = *input.ProcessName
	}
	if input.SleepInfo != nil {
		newCallback["sleepInfo"] = *input.SleepInfo
	}
	if input.ExtraInfo != nil {
		newCallback["extraInfo"] = *input.ExtraInfo
	}

	variables := map[string]interface{}{
		"payloadUuid": input.PayloadUUID,
		"newCallback": newCallback,
	}

	err := c.executeMutation(ctx, &mutation, variables)
	if err != nil {
		return WrapError("CreateCallback", err, "failed to create callback")
	}

	if mutation.CreateCallback.Status != "success" {
		return WrapError("CreateCallback", ErrInvalidResponse, fmt.Sprintf("callback creation failed: %s", mutation.CreateCallback.Error))
	}

	return nil
}

// DeleteCallback deletes callback(s) and their associated tasks.
func (c *Client) DeleteCallback(ctx context.Context, callbackIDs []int) error {
	if err := c.EnsureAuthenticated(ctx); err != nil {
		return err
	}

	if len(callbackIDs) == 0 {
		return WrapError("DeleteCallback", ErrInvalidInput, "at least one callback ID required")
	}

	var mutation struct {
		DeleteTasksAndCallbacks struct {
			Status          string `graphql:"status"`
			Error           string `graphql:"error"`
			FailedTasks     []int  `graphql:"failed_tasks"`
			FailedCallbacks []int  `graphql:"failed_callbacks"`
		} `graphql:"deleteTasksAndCallbacks(callbacks: $callbacks)"`
	}

	variables := map[string]interface{}{
		"callbacks": callbackIDs,
	}

	err := c.executeMutation(ctx, &mutation, variables)
	if err != nil {
		return WrapError("DeleteCallback", err, "failed to delete callbacks")
	}

	if mutation.DeleteTasksAndCallbacks.Status != "success" {
		return WrapError("DeleteCallback", ErrInvalidResponse, fmt.Sprintf("deletion failed: %s", mutation.DeleteTasksAndCallbacks.Error))
	}

	if len(mutation.DeleteTasksAndCallbacks.FailedCallbacks) > 0 {
		return WrapError("DeleteCallback", ErrInvalidResponse, fmt.Sprintf("failed to delete some callbacks: %v", mutation.DeleteTasksAndCallbacks.FailedCallbacks))
	}

	return nil
}

// AddCallbackGraphEdge adds a P2P connection edge between two callbacks.
func (c *Client) AddCallbackGraphEdge(ctx context.Context, sourceID, destinationID int, c2ProfileName string) error {
	if err := c.EnsureAuthenticated(ctx); err != nil {
		return err
	}

	if sourceID <= 0 || destinationID <= 0 {
		return WrapError("AddCallbackGraphEdge", ErrInvalidInput, "source and destination IDs must be positive")
	}

	if c2ProfileName == "" {
		return WrapError("AddCallbackGraphEdge", ErrInvalidInput, "c2 profile name is required")
	}

	var mutation struct {
		CallbackGraphEdgeAdd struct {
			Status string `graphql:"status"`
			Error  string `graphql:"error"`
		} `graphql:"callbackgraphedge_add(source_id: $source_id, destination_id: $destination_id, c2profile: $c2profile)"`
	}

	variables := map[string]interface{}{
		"source_id":      sourceID,
		"destination_id": destinationID,
		"c2profile":      c2ProfileName,
	}

	err := c.executeMutation(ctx, &mutation, variables)
	if err != nil {
		return WrapError("AddCallbackGraphEdge", err, "failed to add callback edge")
	}

	if mutation.CallbackGraphEdgeAdd.Status != "success" {
		return WrapError("AddCallbackGraphEdge", ErrInvalidResponse, fmt.Sprintf("edge creation failed: %s", mutation.CallbackGraphEdgeAdd.Error))
	}

	return nil
}

// RemoveCallbackGraphEdge removes a P2P connection edge between callbacks.
func (c *Client) RemoveCallbackGraphEdge(ctx context.Context, edgeID int) error {
	if err := c.EnsureAuthenticated(ctx); err != nil {
		return err
	}

	if edgeID <= 0 {
		return WrapError("RemoveCallbackGraphEdge", ErrInvalidInput, "edge ID must be positive")
	}

	var mutation struct {
		CallbackGraphEdgeRemove struct {
			Status string `graphql:"status"`
			Error  string `graphql:"error"`
		} `graphql:"callbackgraphedge_remove(edge_id: $edge_id)"`
	}

	variables := map[string]interface{}{
		"edge_id": edgeID,
	}

	err := c.executeMutation(ctx, &mutation, variables)
	if err != nil {
		return WrapError("RemoveCallbackGraphEdge", err, "failed to remove callback edge")
	}

	if mutation.CallbackGraphEdgeRemove.Status != "success" {
		return WrapError("RemoveCallbackGraphEdge", ErrInvalidResponse, fmt.Sprintf("edge removal failed: %s", mutation.CallbackGraphEdgeRemove.Error))
	}

	return nil
}

// ExportCallbackConfig exports a callback's configuration.
func (c *Client) ExportCallbackConfig(ctx context.Context, agentCallbackID string) (string, error) {
	if err := c.EnsureAuthenticated(ctx); err != nil {
		return "", err
	}

	if agentCallbackID == "" {
		return "", WrapError("ExportCallbackConfig", ErrInvalidInput, "agent_callback_id is required")
	}

	var query struct {
		ExportCallbackConfig struct {
			Status          string `graphql:"status"`
			Error           string `graphql:"error"`
			AgentCallbackID string `graphql:"agent_callback_id"`
			Config          string `graphql:"config"`
		} `graphql:"exportCallbackConfig(agent_callback_id: $agent_callback_id)"`
	}

	variables := map[string]interface{}{
		"agent_callback_id": agentCallbackID,
	}

	err := c.executeQuery(ctx, &query, variables)
	if err != nil {
		return "", WrapError("ExportCallbackConfig", err, "failed to export callback config")
	}

	if query.ExportCallbackConfig.Status != "success" {
		return "", WrapError("ExportCallbackConfig", ErrInvalidResponse, fmt.Sprintf("export failed: %s", query.ExportCallbackConfig.Error))
	}

	return query.ExportCallbackConfig.Config, nil
}

// ImportCallbackConfig imports a callback configuration.
func (c *Client) ImportCallbackConfig(ctx context.Context, config string) error {
	if err := c.EnsureAuthenticated(ctx); err != nil {
		return err
	}

	if config == "" {
		return WrapError("ImportCallbackConfig", ErrInvalidInput, "config is required")
	}

	// Parse config as JSON to pass as jsonb
	var configMap map[string]interface{}
	if err := parseJSON([]byte(config), &configMap); err != nil {
		return WrapError("ImportCallbackConfig", err, "invalid config JSON")
	}

	var mutation struct {
		ImportCallbackConfig struct {
			Status string `graphql:"status"`
			Error  string `graphql:"error"`
		} `graphql:"importCallbackConfig(config: $config)"`
	}

	variables := map[string]interface{}{
		"config": configMap,
	}

	err := c.executeMutation(ctx, &mutation, variables)
	if err != nil {
		return WrapError("ImportCallbackConfig", err, "failed to import callback config")
	}

	if mutation.ImportCallbackConfig.Status != "success" {
		return WrapError("ImportCallbackConfig", ErrInvalidResponse, fmt.Sprintf("import failed: %s", mutation.ImportCallbackConfig.Error))
	}

	return nil
}

// parseIPString parses the IP string from Mythic into a slice of IP addresses.
// Mythic stores IPs as a comma-separated string.
func parseIPString(ipStr string) []string {
	if ipStr == "" {
		return []string{}
	}

	// Simple split on comma - could be enhanced for more robust parsing
	ips := []string{}
	current := ""
	for _, c := range ipStr {
		if c == ',' {
			if current != "" {
				ips = append(ips, current)
				current = ""
			}
		} else if c != ' ' {
			current += string(c)
		}
	}
	if current != "" {
		ips = append(ips, current)
	}

	return ips
}

// formatIPString formats a slice of IP addresses into Mythic's comma-separated format.
func formatIPString(ips []string) string {
	result := ""
	for i, ip := range ips {
		if i > 0 {
			result += ","
		}
		result += ip
	}
	return result
}
