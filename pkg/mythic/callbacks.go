package mythic

import (
	"context"
	"fmt"
	"time"

	"github.com/your-org/mythic-sdk-go/pkg/mythic/types"
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
			PayloadTypeID   int       `graphql:"payload_type_id"`
			OperatorID      int       `graphql:"operator_id"`
			Payload         struct {
				ID          int    `graphql:"id"`
				UUID        string `graphql:"uuid"`
				Description string `graphql:"description"`
				OS          string `graphql:"os"`
			} `graphql:"payload"`
			PayloadType struct {
				ID   int    `graphql:"id"`
				Name string `graphql:"name"`
			} `graphql:"payloadtype"`
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
			PayloadTypeID:   cb.PayloadTypeID,
			OperatorID:      cb.OperatorID,
			Payload: &types.CallbackPayload{
				ID:          cb.Payload.ID,
				UUID:        cb.Payload.UUID,
				Description: cb.Payload.Description,
				OS:          cb.Payload.OS,
			},
			PayloadType: &types.CallbackPayloadType{
				ID:   cb.PayloadType.ID,
				Name: cb.PayloadType.Name,
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
			PayloadTypeID   int       `graphql:"payload_type_id"`
			OperatorID      int       `graphql:"operator_id"`
			Payload         struct {
				ID          int    `graphql:"id"`
				UUID        string `graphql:"uuid"`
				Description string `graphql:"description"`
				OS          string `graphql:"os"`
			} `graphql:"payload"`
			PayloadType struct {
				ID   int    `graphql:"id"`
				Name string `graphql:"name"`
			} `graphql:"payloadtype"`
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
			PayloadTypeID:   cb.PayloadTypeID,
			OperatorID:      cb.OperatorID,
			Payload: &types.CallbackPayload{
				ID:          cb.Payload.ID,
				UUID:        cb.Payload.UUID,
				Description: cb.Payload.Description,
				OS:          cb.Payload.OS,
			},
			PayloadType: &types.CallbackPayloadType{
				ID:   cb.PayloadType.ID,
				Name: cb.PayloadType.Name,
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
			PayloadTypeID   int       `graphql:"payload_type_id"`
			OperatorID      int       `graphql:"operator_id"`
			Payload         struct {
				ID          int    `graphql:"id"`
				UUID        string `graphql:"uuid"`
				Description string `graphql:"description"`
				OS          string `graphql:"os"`
			} `graphql:"payload"`
			PayloadType struct {
				ID   int    `graphql:"id"`
				Name string `graphql:"name"`
			} `graphql:"payloadtype"`
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
		PayloadTypeID:   cb.PayloadTypeID,
		OperatorID:      cb.OperatorID,
		Payload: &types.CallbackPayload{
			ID:          cb.Payload.ID,
			UUID:        cb.Payload.UUID,
			Description: cb.Payload.Description,
			OS:          cb.Payload.OS,
		},
		PayloadType: &types.CallbackPayloadType{
			ID:   cb.PayloadType.ID,
			Name: cb.PayloadType.Name,
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
