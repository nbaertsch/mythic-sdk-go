package types

import (
	"fmt"
	"time"
)

// HostInfo represents a compromised or discovered host in Mythic.
// Hosts are distinct from callbacks - a single host may have multiple
// callbacks or be discovered through reconnaissance before compromise.
type HostInfo struct {
	ID          int       `json:"id"`
	Hostname    string    `json:"host"`         // Host identifier/hostname
	IP          string    `json:"ip"`           // IP address(es)
	Domain      string    `json:"domain"`       // Active Directory domain
	OS          string    `json:"os"`           // Operating system details
	Architecture string   `json:"architecture"` // x64, x86, ARM, etc.
	OperationID int       `json:"operation_id"` // Associated operation
	Timestamp   time.Time `json:"timestamp"`    // When host was discovered/added

	// Related entities (populated with nested queries)
	Callbacks []*Callback `json:"callbacks,omitempty"` // Active callbacks on this host
}

// HostNetworkMap represents network topology information.
type HostNetworkMap struct {
	Hosts       []*HostInfo        `json:"hosts"`
	Connections []HostConnection   `json:"connections"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

// HostConnection represents a network connection or relationship between hosts.
type HostConnection struct {
	SourceHostID int    `json:"source_host_id"`
	TargetHostID int    `json:"target_host_id"`
	Type         string `json:"type"` // e.g., "direct", "pivot", "discovered"
	Protocol     string `json:"protocol,omitempty"`
}

// String returns a string representation of a HostInfo.
func (h *HostInfo) String() string {
	if h.Hostname != "" && h.IP != "" {
		return fmt.Sprintf("%s (%s)", h.Hostname, h.IP)
	}
	if h.Hostname != "" {
		return h.Hostname
	}
	if h.IP != "" {
		return h.IP
	}
	return fmt.Sprintf("Host %d", h.ID)
}

// GetCallbackCount returns the number of active callbacks on this host.
func (h *HostInfo) GetCallbackCount() int {
	if h.Callbacks == nil {
		return 0
	}
	count := 0
	for _, cb := range h.Callbacks {
		if cb.Active {
			count++
		}
	}
	return count
}

// HasActiveCallbacks returns true if the host has any active callbacks.
func (h *HostInfo) HasActiveCallbacks() bool {
	return h.GetCallbackCount() > 0
}

// IsWindows returns true if the host is running Windows.
func (h *HostInfo) IsWindows() bool {
	return contains(h.OS, "Windows") || contains(h.OS, "windows")
}

// IsLinux returns true if the host is running Linux.
func (h *HostInfo) IsLinux() bool {
	return contains(h.OS, "Linux") || contains(h.OS, "linux")
}

// IsMacOS returns true if the host is running macOS.
func (h *HostInfo) IsMacOS() bool {
	return contains(h.OS, "Darwin") || contains(h.OS, "macOS") || contains(h.OS, "Mac OS")
}

// contains checks if a string contains a substring (case-sensitive).
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) &&
		(s[:len(substr)] == substr || s[len(s)-len(substr):] == substr ||
		findSubstring(s, substr)))
}

// findSubstring performs a simple substring search.
func findSubstring(s, substr string) bool {
	if len(substr) == 0 {
		return true
	}
	if len(s) < len(substr) {
		return false
	}
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
