package types

import (
	"fmt"
	"time"
)

// RPFWD represents a Reverse Port Forward tunnel in Mythic.
// RPFWD exposes internal network services to the operator by tunneling
// traffic from the target network through a compromised callback back
// to the Mythic server. This is distinct from SOCKS proxies.
//
// Difference from SOCKS:
//   - SOCKS: Operator tools → Mythic → Callback → Target (outbound from operator)
//   - RPFWD: Target Service → Callback → Mythic → Operator (inbound to operator)
type RPFWD struct {
	ID         int       `json:"id"`
	CallbackID int       `json:"callback_id"` // Callback providing the tunnel
	LocalPort  int       `json:"local_port"`  // Port on Mythic server
	RemoteHost string    `json:"remote_host"` // Target host in victim network
	RemotePort int       `json:"remote_port"` // Target port
	Active     bool      `json:"active"`      // Tunnel status
	Timestamp  time.Time `json:"timestamp"`   // When tunnel was created

	// Related entities (populated with nested queries)
	Callback *Callback `json:"callback,omitempty"`
}

// CreateRPFWDRequest represents a request to create a reverse port forward tunnel.
type CreateRPFWDRequest struct {
	CallbackID int    `json:"callback_id"` // Callback to create tunnel through (required)
	LocalPort  int    `json:"local_port"`  // Port on Mythic server (required)
	RemoteHost string `json:"remote_host"` // Target host in victim network (required)
	RemotePort int    `json:"remote_port"` // Target port (required)
}

// String returns a string representation of an RPFWD.
func (r *RPFWD) String() string {
	status := "Active"
	if !r.Active {
		status = "Inactive"
	}
	return fmt.Sprintf("RPFWD %d: localhost:%d → %s:%d (%s)",
		r.ID, r.LocalPort, r.RemoteHost, r.RemotePort, status)
}

// IsActive returns true if the RPFWD tunnel is active.
func (r *RPFWD) IsActive() bool {
	return r.Active
}

// GetLocalEndpoint returns the local endpoint string (localhost:port).
func (r *RPFWD) GetLocalEndpoint() string {
	return fmt.Sprintf("localhost:%d", r.LocalPort)
}

// GetRemoteEndpoint returns the remote endpoint string (host:port).
func (r *RPFWD) GetRemoteEndpoint() string {
	return fmt.Sprintf("%s:%d", r.RemoteHost, r.RemotePort)
}

// Validate validates a CreateRPFWDRequest.
func (req *CreateRPFWDRequest) Validate() error {
	if req.CallbackID == 0 {
		return fmt.Errorf("callback_id is required")
	}
	if req.LocalPort == 0 {
		return fmt.Errorf("local_port is required")
	}
	if req.LocalPort < 1 || req.LocalPort > 65535 {
		return fmt.Errorf("local_port must be between 1 and 65535")
	}
	if req.RemoteHost == "" {
		return fmt.Errorf("remote_host is required")
	}
	if req.RemotePort == 0 {
		return fmt.Errorf("remote_port is required")
	}
	if req.RemotePort < 1 || req.RemotePort > 65535 {
		return fmt.Errorf("remote_port must be between 1 and 65535")
	}
	return nil
}
