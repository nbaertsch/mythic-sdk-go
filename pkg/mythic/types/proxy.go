package types

import "fmt"

// ProxyInfo represents a SOCKS proxy configuration for a callback.
type ProxyInfo struct {
	ID            int    `json:"id"`
	CallbackID    int    `json:"callback_id"`
	Port          int    `json:"port"`
	PortType      string `json:"port_type"` // "socks" or other types
	Active        bool   `json:"active"`
	LocalPort     int    `json:"local_port"`
	RemoteIP      string `json:"remote_ip"`
	RemotePort    int    `json:"remote_port"`
	OperationID   int    `json:"operation_id"`
	Deleted       bool   `json:"deleted"`
	ProxyCallback *int   `json:"proxy_callback_id,omitempty"` // Callback that provides the proxy
}

// String returns a human-readable representation of the proxy info.
func (p *ProxyInfo) String() string {
	status := "inactive"
	if p.Active {
		status = "active"
	}
	if p.Deleted {
		status = "deleted"
	}
	return fmt.Sprintf("%s proxy on port %d (%s)", p.PortType, p.Port, status)
}

// IsActive returns true if the proxy is active.
func (p *ProxyInfo) IsActive() bool {
	return p.Active && !p.Deleted
}

// IsDeleted returns true if the proxy has been deleted.
func (p *ProxyInfo) IsDeleted() bool {
	return p.Deleted
}

// ToggleProxyRequest represents a request to enable or disable a SOCKS proxy.
type ToggleProxyRequest struct {
	TaskID int  `json:"task_id"`
	Port   int  `json:"port"`
	Enable bool `json:"enable"`
}

// String returns a human-readable representation of the request.
func (t *ToggleProxyRequest) String() string {
	action := "disable"
	if t.Enable {
		action = "enable"
	}
	return fmt.Sprintf("%s proxy on port %d (task %d)", action, t.Port, t.TaskID)
}

// TestProxyRequest represents a request to test a proxy connection.
type TestProxyRequest struct {
	CallbackID int    `json:"callback_id"`
	Port       int    `json:"port"`
	TargetURL  string `json:"target_url"`
}

// String returns a human-readable representation of the request.
func (t *TestProxyRequest) String() string {
	return fmt.Sprintf("Test proxy on callback %d port %d -> %s", t.CallbackID, t.Port, t.TargetURL)
}

// TestProxyResponse represents the response from testing a proxy connection.
type TestProxyResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Error   string `json:"error,omitempty"`
}

// String returns a human-readable representation of the response.
func (t *TestProxyResponse) String() string {
	if t.Status == "success" {
		return fmt.Sprintf("Proxy test: %s", t.Message)
	}
	return fmt.Sprintf("Proxy test failed: %s", t.Error)
}

// IsSuccessful returns true if the proxy test succeeded.
func (t *TestProxyResponse) IsSuccessful() bool {
	return t.Status == "success"
}
