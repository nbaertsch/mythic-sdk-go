package mythic

import (
	"context"

	"github.com/nbaertsch/mythic-sdk-go/pkg/mythic/types"
)

// ToggleProxy enables or disables a SOCKS proxy on a callback.
// SOCKS proxies allow routing network traffic through compromised systems,
// enabling lateral movement and access to internal networks.
//
// Parameters:
//   - ctx: Context for cancellation and timeout
//   - taskID: The task ID that started/will stop the proxy
//   - port: The port number for the SOCKS proxy
//   - enable: true to start the proxy, false to stop it
//
// Returns:
//   - *types.ProxyInfo: Information about the proxy state
//   - error: Error if the operation fails
//
// Example:
//
//	// Enable SOCKS proxy on port 1080
//	proxy, err := client.ToggleProxy(ctx, taskID, 1080, true)
//	if err != nil {
//	    return err
//	}
//	fmt.Printf("Proxy started: %s\n", proxy.String())
//
//	// Later, disable the proxy
//	proxy, err = client.ToggleProxy(ctx, taskID, 1080, false)
func (c *Client) ToggleProxy(ctx context.Context, taskID, port int, enable bool) (*types.ProxyInfo, error) {
	if err := c.EnsureAuthenticated(ctx); err != nil {
		return nil, err
	}

	// Validate input
	if taskID <= 0 {
		return nil, WrapError("ToggleProxy", ErrInvalidInput, "task ID must be positive")
	}
	if port <= 0 || port > 65535 {
		return nil, WrapError("ToggleProxy", ErrInvalidInput, "port must be between 1 and 65535")
	}

	var mutation struct {
		Response struct {
			Status string `graphql:"status"`
			Error  string `graphql:"error"`
			Proxy  *struct {
				ID            int    `graphql:"id"`
				CallbackID    int    `graphql:"callback_id"`
				Port          int    `graphql:"port"`
				PortType      string `graphql:"port_type"`
				Active        bool   `graphql:"active"`
				LocalPort     int    `graphql:"local_port"`
				RemoteIP      string `graphql:"remote_ip"`
				RemotePort    int    `graphql:"remote_port"`
				OperationID   int    `graphql:"operation_id"`
				Deleted       bool   `graphql:"deleted"`
				ProxyCallback *int   `graphql:"proxy_callback_id"`
			} `graphql:"proxy"`
		} `graphql:"toggleProxy(task_id: $task_id, port: $port, enable: $enable)"`
	}

	variables := map[string]interface{}{
		"task_id": taskID,
		"port":    port,
		"enable":  enable,
	}

	if err := c.executeMutation(ctx, &mutation, variables); err != nil {
		return nil, WrapError("ToggleProxy", err, "failed to toggle proxy")
	}

	// Check for error in response
	if mutation.Response.Status != "success" {
		return nil, WrapError("ToggleProxy", ErrOperationFailed, mutation.Response.Error)
	}

	// If no proxy info returned, return nil (might happen on disable)
	if mutation.Response.Proxy == nil {
		return nil, nil
	}

	// Map to types.ProxyInfo
	proxyInfo := &types.ProxyInfo{
		ID:            mutation.Response.Proxy.ID,
		CallbackID:    mutation.Response.Proxy.CallbackID,
		Port:          mutation.Response.Proxy.Port,
		PortType:      mutation.Response.Proxy.PortType,
		Active:        mutation.Response.Proxy.Active,
		LocalPort:     mutation.Response.Proxy.LocalPort,
		RemoteIP:      mutation.Response.Proxy.RemoteIP,
		RemotePort:    mutation.Response.Proxy.RemotePort,
		OperationID:   mutation.Response.Proxy.OperationID,
		Deleted:       mutation.Response.Proxy.Deleted,
		ProxyCallback: mutation.Response.Proxy.ProxyCallback,
	}

	return proxyInfo, nil
}

// TestProxy tests a SOCKS proxy connection by attempting to connect
// to a target URL through the proxy. This validates that the proxy
// is functioning correctly and can route traffic.
//
// Parameters:
//   - ctx: Context for cancellation and timeout
//   - callbackID: The callback that hosts the SOCKS proxy
//   - port: The SOCKS proxy port to test
//   - targetURL: The URL to test connectivity to (e.g., "https://www.google.com")
//
// Returns:
//   - *types.TestProxyResponse: Result of the proxy test
//   - error: Error if the operation fails
//
// Example:
//
//	// Test if the proxy can reach an external site
//	result, err := client.TestProxy(ctx, callbackID, 1080, "https://www.google.com")
//	if err != nil {
//	    return err
//	}
//	if result.IsSuccessful() {
//	    fmt.Printf("Proxy is working: %s\n", result.Message)
//	} else {
//	    fmt.Printf("Proxy test failed: %s\n", result.Error)
//	}
func (c *Client) TestProxy(ctx context.Context, callbackID, port int, targetURL string) (*types.TestProxyResponse, error) {
	if err := c.EnsureAuthenticated(ctx); err != nil {
		return nil, err
	}

	// Validate input
	if callbackID <= 0 {
		return nil, WrapError("TestProxy", ErrInvalidInput, "callback ID must be positive")
	}
	if port <= 0 || port > 65535 {
		return nil, WrapError("TestProxy", ErrInvalidInput, "port must be between 1 and 65535")
	}
	if targetURL == "" {
		return nil, WrapError("TestProxy", ErrInvalidInput, "target URL cannot be empty")
	}

	var mutation struct {
		Response struct {
			Status  string `graphql:"status"`
			Message string `graphql:"message"`
			Error   string `graphql:"error"`
		} `graphql:"testProxy(callback_id: $callback_id, port: $port, target_url: $target_url)"`
	}

	variables := map[string]interface{}{
		"callback_id": callbackID,
		"port":        port,
		"target_url":  targetURL,
	}

	if err := c.executeMutation(ctx, &mutation, variables); err != nil {
		return nil, WrapError("TestProxy", err, "failed to test proxy")
	}

	// Map to types.TestProxyResponse
	response := &types.TestProxyResponse{
		Status:  mutation.Response.Status,
		Message: mutation.Response.Message,
		Error:   mutation.Response.Error,
	}

	return response, nil
}
