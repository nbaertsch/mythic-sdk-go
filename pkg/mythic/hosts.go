package mythic

import (
	"context"
	"fmt"
	"time"

	"github.com/nbaertsch/mythic-sdk-go/pkg/mythic/types"
)

// GetHosts retrieves all hosts tracked in an operation.
//
// Hosts represent compromised or discovered systems across the network.
// This is distinct from callbacks - a host may have multiple callbacks
// or be discovered through reconnaissance before compromise.
//
// Parameters:
//   - ctx: Context for cancellation and timeout
//   - operationID: ID of the operation (0 for current operation)
//
// Returns:
//   - []*types.HostInfo: List of hosts in the operation
//   - error: Error if operation ID is invalid or query fails
//
// Example:
//
//	hosts, err := client.GetHosts(ctx, 0)
//	if err != nil {
//	    return err
//	}
//	for _, host := range hosts {
//	    fmt.Printf("Host: %s (%s) - %d active callbacks\n",
//	        host.Hostname, host.IP, host.GetCallbackCount())
//	}
func (c *Client) GetHosts(ctx context.Context, operationID int) ([]*types.HostInfo, error) {
	if err := c.EnsureAuthenticated(ctx); err != nil {
		return nil, err
	}

	// Use current operation if not specified
	if operationID == 0 {
		currentOp := c.GetCurrentOperation()
		if currentOp == nil {
			return nil, WrapError("GetHosts", ErrNotAuthenticated, "no current operation set")
		}
		operationID = *currentOp
	}

	// Mythic has no dedicated "host" table. Hosts are derived from unique
	// host values in the callback table and enriched with mythictree data.
	var query struct {
		Callback []struct {
			Host         string    `graphql:"host"`
			Os           string    `graphql:"os"`
			Domain       string    `graphql:"domain"`
			Architecture string    `graphql:"architecture"`
			IP           string    `graphql:"ip"`
			OperationID  int       `graphql:"operation_id"`
			LastCheckin  string `graphql:"last_checkin"`
		} `graphql:"callback(where: {operation_id: {_eq: $operation_id}}, distinct_on: host, order_by: {host: asc})"`
	}

	variables := map[string]interface{}{
		"operation_id": operationID,
	}

	err := c.executeQuery(ctx, &query, variables)
	if err != nil {
		return nil, WrapError("GetHosts", err, "failed to query hosts")
	}

	hosts := make([]*types.HostInfo, len(query.Callback))
	for i, cbData := range query.Callback {
		ts, _ := parseTimestamp(cbData.LastCheckin)
		hosts[i] = &types.HostInfo{
			ID:           i + 1, // synthetic ID since hosts are derived
			Hostname:     cbData.Host,
			IP:           cbData.IP,
			Domain:       cbData.Domain,
			OS:           cbData.Os,
			Architecture: cbData.Architecture,
			OperationID:  cbData.OperationID,
			Timestamp:    ts,
		}
	}

	return hosts, nil
}

// GetHostByID retrieves a specific host by its database ID.
//
// Parameters:
//   - ctx: Context for cancellation and timeout
//   - hostID: Database ID of the host
//
// Returns:
//   - *types.HostInfo: Host information
//   - error: Error if host ID is invalid or not found
//
// Example:
//
//	host, err := client.GetHostByID(ctx, 42)
//	if err != nil {
//	    return err
//	}
//	fmt.Printf("Host: %s\nOS: %s\nArchitecture: %s\n",
//	    host.Hostname, host.OS, host.Architecture)
func (c *Client) GetHostByID(ctx context.Context, hostID int) (*types.HostInfo, error) {
	if err := c.EnsureAuthenticated(ctx); err != nil {
		return nil, err
	}

	if hostID == 0 {
		return nil, WrapError("GetHostByID", ErrInvalidInput, "host ID is required")
	}

	// Mythic has no dedicated host table. Look up a callback by ID and
	// return its host information as a HostInfo. This provides backward
	// compatibility while using the actual Mythic schema.
	var query struct {
		Callback []struct {
			ID           int       `graphql:"id"`
			Host         string    `graphql:"host"`
			IP           string    `graphql:"ip"`
			Domain       string    `graphql:"domain"`
			Os           string    `graphql:"os"`
			Architecture string    `graphql:"architecture"`
			OperationID  int       `graphql:"operation_id"`
			LastCheckin  string `graphql:"last_checkin"`
		} `graphql:"callback(where: {id: {_eq: $host_id}})"`
	}

	variables := map[string]interface{}{
		"host_id": hostID,
	}

	err := c.executeQuery(ctx, &query, variables)
	if err != nil {
		return nil, WrapError("GetHostByID", err, "failed to query host")
	}

	if len(query.Callback) == 0 {
		return nil, WrapError("GetHostByID", ErrNotFound, fmt.Sprintf("host %d not found", hostID))
	}

	cbData := query.Callback[0]
	ts, _ := parseTimestamp(cbData.LastCheckin)
	return &types.HostInfo{
		ID:           cbData.ID,
		Hostname:     cbData.Host,
		IP:           cbData.IP,
		Domain:       cbData.Domain,
		OS:           cbData.Os,
		Architecture: cbData.Architecture,
		OperationID:  cbData.OperationID,
		Timestamp:    ts,
	}, nil
}

// GetHostByHostname finds a host by its hostname.
//
// This performs a case-insensitive search for the hostname.
//
// Parameters:
//   - ctx: Context for cancellation and timeout
//   - hostname: Hostname to search for
//
// Returns:
//   - *types.HostInfo: Host information (first match if multiple)
//   - error: Error if hostname is empty or not found
//
// Example:
//
//	host, err := client.GetHostByHostname(ctx, "WORKSTATION-01")
//	if err != nil {
//	    return err
//	}
//	fmt.Printf("Found host: %s (%s)\n", host.Hostname, host.IP)
func (c *Client) GetHostByHostname(ctx context.Context, hostname string) (*types.HostInfo, error) {
	if err := c.EnsureAuthenticated(ctx); err != nil {
		return nil, err
	}

	if hostname == "" {
		return nil, WrapError("GetHostByHostname", ErrInvalidInput, "hostname is required")
	}

	// Hosts are derived from callbacks — find the most recent callback
	// matching this hostname.
	var query struct {
		Callback []struct {
			ID           int       `graphql:"id"`
			Host         string    `graphql:"host"`
			IP           string    `graphql:"ip"`
			Domain       string    `graphql:"domain"`
			Os           string    `graphql:"os"`
			Architecture string    `graphql:"architecture"`
			OperationID  int       `graphql:"operation_id"`
			LastCheckin  string `graphql:"last_checkin"`
		} `graphql:"callback(where: {host: {_ilike: $hostname}}, order_by: {last_checkin: desc}, limit: 1)"`
	}

	variables := map[string]interface{}{
		"hostname": hostname,
	}

	err := c.executeQuery(ctx, &query, variables)
	if err != nil {
		return nil, WrapError("GetHostByHostname", err, "failed to query host")
	}

	if len(query.Callback) == 0 {
		return nil, WrapError("GetHostByHostname", ErrNotFound, fmt.Sprintf("host '%s' not found", hostname))
	}

	cbData := query.Callback[0]
	ts, _ := parseTimestamp(cbData.LastCheckin)
	return &types.HostInfo{
		ID:           cbData.ID,
		Hostname:     cbData.Host,
		IP:           cbData.IP,
		Domain:       cbData.Domain,
		OS:           cbData.Os,
		Architecture: cbData.Architecture,
		OperationID:  cbData.OperationID,
		Timestamp:    ts,
	}, nil
}

// GetCallbacksForHost retrieves all callbacks associated with a specific host.
//
// This shows which agents are currently running on the host and their status.
// Mythic has no dedicated host table, so this queries callbacks by hostname
// using a case-insensitive match.
//
// Parameters:
//   - ctx: Context for cancellation and timeout
//   - hostname: Hostname to search for (case-insensitive)
//
// Returns:
//   - []*types.Callback: List of callbacks on the host
//   - error: Error if hostname is empty or query fails
//
// Example:
//
//	callbacks, err := client.GetCallbacksForHost(ctx, "WORKSTATION-01")
//	if err != nil {
//	    return err
//	}
//	fmt.Printf("Host has %d callbacks:\n", len(callbacks))
//	for _, cb := range callbacks {
//	    status := "Inactive"
//	    if cb.Active {
//	        status = "Active"
//	    }
//	    fmt.Printf("  - Callback %d: %s@%s (%s)\n",
//	        cb.ID, cb.User, cb.Host, status)
//	}
func (c *Client) GetCallbacksForHost(ctx context.Context, hostname string) ([]*types.Callback, error) {
	if err := c.EnsureAuthenticated(ctx); err != nil {
		return nil, err
	}

	if hostname == "" {
		return nil, WrapError("GetCallbacksForHost", ErrInvalidInput, "hostname is required")
	}

	// Query callbacks by matching hostname directly
	var query struct {
		Callback []struct {
			ID                  int       `graphql:"id"`
			DisplayID           int       `graphql:"display_id"`
			AgentCallbackID     string    `graphql:"agent_callback_id"`
			InitCallback        string `graphql:"init_callback"`
			LastCheckin         string `graphql:"last_checkin"`
			User                string    `graphql:"user"`
			Host                string    `graphql:"host"`
			PID                 int       `graphql:"pid"`
			IP                  string    `graphql:"ip"`
			ExternalIP          string    `graphql:"external_ip"`
			ProcessName         string    `graphql:"process_name"`
			Description         string    `graphql:"description"`
			OperatorID          int       `graphql:"operator_id"`
			Active              bool      `graphql:"active"`
			RegisteredPayloadID string    `graphql:"registered_payload_id"`
			IntegrityLevel      int       `graphql:"integrity_level"`
			Locked              bool      `graphql:"locked"`
			OperationID         int       `graphql:"operation_id"`
			SleepInfo           string    `graphql:"sleep_info"`
			Architecture        string    `graphql:"architecture"`
			Domain              string    `graphql:"domain"`
			Os                  string    `graphql:"os"`
		} `graphql:"callback(where: {host: {_ilike: $hostname}}, order_by: {last_checkin: desc})"`
	}

	variables := map[string]interface{}{
		"hostname": hostname,
	}

	err := c.executeQuery(ctx, &query, variables)
	if err != nil {
		return nil, WrapError("GetCallbacksForHost", err, "failed to query callbacks")
	}

	callbacks := make([]*types.Callback, len(query.Callback))
	for i, cbData := range query.Callback {
		// Convert IP string to []string for types.Callback
		ips := []string{}
		if cbData.IP != "" {
			ips = []string{cbData.IP}
		}

		initCb, _ := parseTimestamp(cbData.InitCallback)
		lastCb, _ := parseTimestamp(cbData.LastCheckin)
		callbacks[i] = &types.Callback{
			ID:                  cbData.ID,
			DisplayID:           cbData.DisplayID,
			AgentCallbackID:     cbData.AgentCallbackID,
			InitCallback:        initCb,
			LastCheckin:         lastCb,
			User:                cbData.User,
			Host:                cbData.Host,
			PID:                 cbData.PID,
			IP:                  ips,
			ExternalIP:          cbData.ExternalIP,
			ProcessName:         cbData.ProcessName,
			Description:         cbData.Description,
			OperatorID:          cbData.OperatorID,
			Active:              cbData.Active,
			RegisteredPayloadID: cbData.RegisteredPayloadID,
			IntegrityLevel:      types.CallbackIntegrityLevel(cbData.IntegrityLevel),
			Locked:              cbData.Locked,
			OperationID:         cbData.OperationID,
			SleepInfo:           cbData.SleepInfo,
			Architecture:        cbData.Architecture,
			Domain:              cbData.Domain,
			OS:                  cbData.Os,
		}
	}

	return callbacks, nil
}

// GetHostNetworkMap builds a network topology map for an operation.
//
// This provides an overview of all discovered hosts and their relationships,
// useful for lateral movement planning and pivot path identification.
//
// Parameters:
//   - ctx: Context for cancellation and timeout
//   - operationID: ID of the operation (0 for current operation)
//
// Returns:
//   - *types.HostNetworkMap: Network topology information
//   - error: Error if operation ID is invalid or query fails
//
// Example:
//
//	networkMap, err := client.GetHostNetworkMap(ctx, 0)
//	if err != nil {
//	    return err
//	}
//	fmt.Printf("Network topology: %d hosts\n", len(networkMap.Hosts))
//	for _, host := range networkMap.Hosts {
//	    callbackCount := 0
//	    if host.Callbacks != nil {
//	        for _, cb := range host.Callbacks {
//	            if cb.Active {
//	                callbackCount++
//	            }
//	        }
//	    }
//	    fmt.Printf("  - %s: %d active callbacks\n", host.Hostname, callbackCount)
//	}
func (c *Client) GetHostNetworkMap(ctx context.Context, operationID int) (*types.HostNetworkMap, error) {
	if err := c.EnsureAuthenticated(ctx); err != nil {
		return nil, err
	}

	// Use current operation if not specified
	if operationID == 0 {
		currentOp := c.GetCurrentOperation()
		if currentOp == nil {
			return nil, WrapError("GetHostNetworkMap", ErrNotAuthenticated, "no current operation set")
		}
		operationID = *currentOp
	}

	// Get all hosts in the operation
	hosts, err := c.GetHosts(ctx, operationID)
	if err != nil {
		return nil, WrapError("GetHostNetworkMap", err, "failed to get hosts")
	}

	// Enrich each host with callback information
	for _, host := range hosts {
		callbacks, err := c.GetCallbacksForHost(ctx, host.Hostname)
		if err != nil {
			// Log error but continue with other hosts
			continue
		}

		// Callbacks are already types.Callback, just assign directly
		host.Callbacks = callbacks
	}

	// Build network map
	networkMap := &types.HostNetworkMap{
		Hosts:       hosts,
		Connections: []types.HostConnection{}, // TODO: Could be enhanced with actual network connections
		Metadata: map[string]interface{}{
			"operation_id": operationID,
			"host_count":   len(hosts),
			"timestamp":    time.Now().Format(time.RFC3339),
		},
	}

	// Calculate total active callbacks across all hosts
	totalCallbacks := 0
	for _, host := range hosts {
		for _, cb := range host.Callbacks {
			if cb.Active {
				totalCallbacks++
			}
		}
	}
	networkMap.Metadata["total_active_callbacks"] = totalCallbacks

	return networkMap, nil
}
