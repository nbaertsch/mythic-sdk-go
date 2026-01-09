package mythic

import (
	"context"
	"fmt"
	"time"

	"github.com/nbaertsch/mythic-sdk-go/pkg/mythic/types"
)

// GetRPFWDs retrieves all reverse port forward tunnels for a callback.
//
// Reverse port forwards expose internal network services to the operator
// by tunneling traffic through a compromised callback.
//
// Parameters:
//   - ctx: Context for cancellation and timeout
//   - callbackID: ID of the callback providing tunnels
//
// Returns:
//   - []*types.RPFWD: List of RPFWD tunnels
//   - error: Error if callback ID is invalid or query fails
//
// Example:
//
//	tunnels, err := client.GetRPFWDs(ctx, 5)
//	if err != nil {
//	    return err
//	}
//	for _, tunnel := range tunnels {
//	    if tunnel.Active {
//	        fmt.Printf("Active tunnel: %s\n", tunnel.String())
//	    }
//	}
func (c *Client) GetRPFWDs(ctx context.Context, callbackID int) ([]*types.RPFWD, error) {
	if err := c.EnsureAuthenticated(ctx); err != nil {
		return nil, err
	}

	if callbackID == 0 {
		return nil, WrapError("GetRPFWDs", ErrInvalidInput, "callback ID is required")
	}

	var query struct {
		RPFWD []struct {
			ID         int       `graphql:"id"`
			CallbackID int       `graphql:"callback_id"`
			LocalPort  int       `graphql:"local_port"`
			RemoteHost string    `graphql:"remote_host"`
			RemotePort int       `graphql:"remote_port"`
			Active     bool      `graphql:"active"`
			Timestamp  time.Time `graphql:"timestamp"`
		} `graphql:"rpfwd(where: {callback_id: {_eq: $callback_id}}, order_by: {timestamp: desc})"`
	}

	variables := map[string]interface{}{
		"callback_id": callbackID,
	}

	err := c.executeQuery(ctx, &query, variables)
	if err != nil {
		return nil, WrapError("GetRPFWDs", err, "failed to query RPFWDs")
	}

	rpfwds := make([]*types.RPFWD, len(query.RPFWD))
	for i, rpfwdData := range query.RPFWD {
		rpfwds[i] = &types.RPFWD{
			ID:         rpfwdData.ID,
			CallbackID: rpfwdData.CallbackID,
			LocalPort:  rpfwdData.LocalPort,
			RemoteHost: rpfwdData.RemoteHost,
			RemotePort: rpfwdData.RemotePort,
			Active:     rpfwdData.Active,
			Timestamp:  rpfwdData.Timestamp,
		}
	}

	return rpfwds, nil
}

// CreateRPFWD creates a new reverse port forward tunnel.
//
// This establishes a tunnel from the target network through the callback
// back to the Mythic server, exposing an internal service on localhost.
//
// Parameters:
//   - ctx: Context for cancellation and timeout
//   - req: RPFWD creation request with callback, ports, and target
//
// Returns:
//   - *types.RPFWD: Created RPFWD tunnel
//   - error: Error if request is invalid or creation fails
//
// Example:
//
//	// Expose internal RDP service on localhost:13389
//	req := &types.CreateRPFWDRequest{
//	    CallbackID: 5,
//	    LocalPort:  13389,
//	    RemoteHost: "10.10.10.50",
//	    RemotePort: 3389,
//	}
//	tunnel, err := client.CreateRPFWD(ctx, req)
//	if err != nil {
//	    return err
//	}
//	fmt.Printf("Tunnel created: %s\n", tunnel.String())
//	fmt.Printf("Connect to: %s\n", tunnel.GetLocalEndpoint())
func (c *Client) CreateRPFWD(ctx context.Context, req *types.CreateRPFWDRequest) (*types.RPFWD, error) {
	if err := c.EnsureAuthenticated(ctx); err != nil {
		return nil, err
	}

	if req == nil {
		return nil, WrapError("CreateRPFWD", ErrInvalidInput, "request is required")
	}

	if err := req.Validate(); err != nil {
		return nil, WrapError("CreateRPFWD", ErrInvalidInput, err.Error())
	}

	// Create RPFWD tunnel
	var mutation struct {
		InsertRpfwdOne struct {
			ID         int       `graphql:"id"`
			CallbackID int       `graphql:"callback_id"`
			LocalPort  int       `graphql:"local_port"`
			RemoteHost string    `graphql:"remote_host"`
			RemotePort int       `graphql:"remote_port"`
			Active     bool      `graphql:"active"`
			Timestamp  time.Time `graphql:"timestamp"`
		} `graphql:"insert_rpfwd_one(object: {callback_id: $callback_id, local_port: $local_port, remote_host: $remote_host, remote_port: $remote_port, active: true})"`
	}

	variables := map[string]interface{}{
		"callback_id": req.CallbackID,
		"local_port":  req.LocalPort,
		"remote_host": req.RemoteHost,
		"remote_port": req.RemotePort,
	}

	err := c.executeMutation(ctx, &mutation, variables)
	if err != nil {
		return nil, WrapError("CreateRPFWD", err, "failed to create RPFWD")
	}

	return &types.RPFWD{
		ID:         mutation.InsertRpfwdOne.ID,
		CallbackID: mutation.InsertRpfwdOne.CallbackID,
		LocalPort:  mutation.InsertRpfwdOne.LocalPort,
		RemoteHost: mutation.InsertRpfwdOne.RemoteHost,
		RemotePort: mutation.InsertRpfwdOne.RemotePort,
		Active:     mutation.InsertRpfwdOne.Active,
		Timestamp:  mutation.InsertRpfwdOne.Timestamp,
	}, nil
}

// DeleteRPFWD closes a reverse port forward tunnel.
//
// This terminates the tunnel and marks it as inactive.
//
// Parameters:
//   - ctx: Context for cancellation and timeout
//   - rpfwdID: ID of the RPFWD tunnel to close
//
// Returns:
//   - error: Error if RPFWD ID is invalid or deletion fails
//
// Example:
//
//	err := client.DeleteRPFWD(ctx, 42)
//	if err != nil {
//	    return err
//	}
//	fmt.Println("Tunnel closed successfully")
func (c *Client) DeleteRPFWD(ctx context.Context, rpfwdID int) error {
	if err := c.EnsureAuthenticated(ctx); err != nil {
		return err
	}

	if rpfwdID == 0 {
		return WrapError("DeleteRPFWD", ErrInvalidInput, "RPFWD ID is required")
	}

	// Mark RPFWD as inactive
	var mutation struct {
		UpdateRpfwd struct {
			AffectedRows int `graphql:"affected_rows"`
		} `graphql:"update_rpfwd(where: {id: {_eq: $rpfwd_id}}, _set: {active: false})"`
	}

	variables := map[string]interface{}{
		"rpfwd_id": rpfwdID,
	}

	err := c.executeMutation(ctx, &mutation, variables)
	if err != nil {
		return WrapError("DeleteRPFWD", err, "failed to delete RPFWD")
	}

	if mutation.UpdateRpfwd.AffectedRows == 0 {
		return WrapError("DeleteRPFWD", ErrNotFound, fmt.Sprintf("RPFWD %d not found", rpfwdID))
	}

	return nil
}

// GetRPFWDStatus retrieves the current status of a reverse port forward tunnel.
//
// This provides detailed information about the tunnel including whether
// it's currently active and its configuration.
//
// Parameters:
//   - ctx: Context for cancellation and timeout
//   - rpfwdID: ID of the RPFWD tunnel
//
// Returns:
//   - *types.RPFWD: RPFWD tunnel information
//   - error: Error if RPFWD ID is invalid or not found
//
// Example:
//
//	status, err := client.GetRPFWDStatus(ctx, 42)
//	if err != nil {
//	    return err
//	}
//	if status.Active {
//	    fmt.Printf("Tunnel is active: %s → %s\n",
//	        status.GetLocalEndpoint(), status.GetRemoteEndpoint())
//	} else {
//	    fmt.Println("Tunnel is inactive")
//	}
func (c *Client) GetRPFWDStatus(ctx context.Context, rpfwdID int) (*types.RPFWD, error) {
	if err := c.EnsureAuthenticated(ctx); err != nil {
		return nil, err
	}

	if rpfwdID == 0 {
		return nil, WrapError("GetRPFWDStatus", ErrInvalidInput, "RPFWD ID is required")
	}

	var query struct {
		RPFWD []struct {
			ID         int       `graphql:"id"`
			CallbackID int       `graphql:"callback_id"`
			LocalPort  int       `graphql:"local_port"`
			RemoteHost string    `graphql:"remote_host"`
			RemotePort int       `graphql:"remote_port"`
			Active     bool      `graphql:"active"`
			Timestamp  time.Time `graphql:"timestamp"`
		} `graphql:"rpfwd(where: {id: {_eq: $rpfwd_id}})"`
	}

	variables := map[string]interface{}{
		"rpfwd_id": rpfwdID,
	}

	err := c.executeQuery(ctx, &query, variables)
	if err != nil {
		return nil, WrapError("GetRPFWDStatus", err, "failed to query RPFWD")
	}

	if len(query.RPFWD) == 0 {
		return nil, WrapError("GetRPFWDStatus", ErrNotFound, fmt.Sprintf("RPFWD %d not found", rpfwdID))
	}

	rpfwdData := query.RPFWD[0]
	return &types.RPFWD{
		ID:         rpfwdData.ID,
		CallbackID: rpfwdData.CallbackID,
		LocalPort:  rpfwdData.LocalPort,
		RemoteHost: rpfwdData.RemoteHost,
		RemotePort: rpfwdData.RemotePort,
		Active:     rpfwdData.Active,
		Timestamp:  rpfwdData.Timestamp,
	}, nil
}
