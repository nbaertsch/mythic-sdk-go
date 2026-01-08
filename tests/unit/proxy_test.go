package unit

import (
	"strings"
	"testing"

	"github.com/nbaertsch/mythic-sdk-go/pkg/mythic/types"
)

func TestProxyInfo_String(t *testing.T) {
	tests := []struct {
		name     string
		proxy    types.ProxyInfo
		contains []string
	}{
		{
			name: "active socks proxy",
			proxy: types.ProxyInfo{
				ID:       1,
				Port:     1080,
				PortType: "socks",
				Active:   true,
				Deleted:  false,
			},
			contains: []string{"socks", "1080", "active"},
		},
		{
			name: "inactive proxy",
			proxy: types.ProxyInfo{
				ID:       2,
				Port:     9050,
				PortType: "socks",
				Active:   false,
				Deleted:  false,
			},
			contains: []string{"socks", "9050", "inactive"},
		},
		{
			name: "deleted proxy",
			proxy: types.ProxyInfo{
				ID:       3,
				Port:     1080,
				PortType: "socks",
				Active:   true,
				Deleted:  true,
			},
			contains: []string{"socks", "1080", "deleted"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.proxy.String()

			for _, substr := range tt.contains {
				if !strings.Contains(result, substr) {
					t.Errorf("String() = %q, should contain %q", result, substr)
				}
			}
		})
	}
}

func TestProxyInfo_IsActive(t *testing.T) {
	tests := []struct {
		name     string
		proxy    types.ProxyInfo
		expected bool
	}{
		{
			name: "active and not deleted",
			proxy: types.ProxyInfo{
				Active:  true,
				Deleted: false,
			},
			expected: true,
		},
		{
			name: "inactive",
			proxy: types.ProxyInfo{
				Active:  false,
				Deleted: false,
			},
			expected: false,
		},
		{
			name: "active but deleted",
			proxy: types.ProxyInfo{
				Active:  true,
				Deleted: true,
			},
			expected: false,
		},
		{
			name: "inactive and deleted",
			proxy: types.ProxyInfo{
				Active:  false,
				Deleted: true,
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.proxy.IsActive()
			if result != tt.expected {
				t.Errorf("IsActive() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestProxyInfo_IsDeleted(t *testing.T) {
	tests := []struct {
		name     string
		proxy    types.ProxyInfo
		expected bool
	}{
		{
			name: "not deleted",
			proxy: types.ProxyInfo{
				Deleted: false,
			},
			expected: false,
		},
		{
			name: "deleted",
			proxy: types.ProxyInfo{
				Deleted: true,
			},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.proxy.IsDeleted()
			if result != tt.expected {
				t.Errorf("IsDeleted() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestToggleProxyRequest_String(t *testing.T) {
	tests := []struct {
		name     string
		request  types.ToggleProxyRequest
		contains []string
	}{
		{
			name: "enable proxy",
			request: types.ToggleProxyRequest{
				TaskID: 123,
				Port:   1080,
				Enable: true,
			},
			contains: []string{"enable", "1080", "123"},
		},
		{
			name: "disable proxy",
			request: types.ToggleProxyRequest{
				TaskID: 456,
				Port:   9050,
				Enable: false,
			},
			contains: []string{"disable", "9050", "456"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.request.String()

			for _, substr := range tt.contains {
				if !strings.Contains(result, substr) {
					t.Errorf("String() = %q, should contain %q", result, substr)
				}
			}
		})
	}
}

func TestTestProxyRequest_String(t *testing.T) {
	tests := []struct {
		name     string
		request  types.TestProxyRequest
		contains []string
	}{
		{
			name: "test proxy to google",
			request: types.TestProxyRequest{
				CallbackID: 42,
				Port:       1080,
				TargetURL:  "https://www.google.com",
			},
			contains: []string{"Test", "42", "1080", "https://www.google.com"},
		},
		{
			name: "test proxy to internal site",
			request: types.TestProxyRequest{
				CallbackID: 10,
				Port:       9050,
				TargetURL:  "http://10.0.0.1",
			},
			contains: []string{"Test", "10", "9050", "http://10.0.0.1"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.request.String()

			for _, substr := range tt.contains {
				if !strings.Contains(result, substr) {
					t.Errorf("String() = %q, should contain %q", result, substr)
				}
			}
		})
	}
}

func TestTestProxyResponse_String(t *testing.T) {
	tests := []struct {
		name     string
		response types.TestProxyResponse
		contains []string
	}{
		{
			name: "successful test",
			response: types.TestProxyResponse{
				Status:  "success",
				Message: "Connection successful",
			},
			contains: []string{"Proxy test", "Connection successful"},
		},
		{
			name: "failed test",
			response: types.TestProxyResponse{
				Status: "error",
				Error:  "Connection timeout",
			},
			contains: []string{"Proxy test failed", "Connection timeout"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.response.String()

			for _, substr := range tt.contains {
				if !strings.Contains(result, substr) {
					t.Errorf("String() = %q, should contain %q", result, substr)
				}
			}
		})
	}
}

func TestTestProxyResponse_IsSuccessful(t *testing.T) {
	tests := []struct {
		name     string
		response types.TestProxyResponse
		expected bool
	}{
		{
			name: "success status",
			response: types.TestProxyResponse{
				Status: "success",
			},
			expected: true,
		},
		{
			name: "error status",
			response: types.TestProxyResponse{
				Status: "error",
			},
			expected: false,
		},
		{
			name: "failed status",
			response: types.TestProxyResponse{
				Status: "failed",
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.response.IsSuccessful()
			if result != tt.expected {
				t.Errorf("IsSuccessful() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestProxyInfo_CompleteStructure(t *testing.T) {
	proxyCallbackID := 5
	proxy := types.ProxyInfo{
		ID:            1,
		CallbackID:    10,
		Port:          1080,
		PortType:      "socks",
		Active:        true,
		LocalPort:     8080,
		RemoteIP:      "192.168.1.100",
		RemotePort:    443,
		OperationID:   1,
		Deleted:       false,
		ProxyCallback: &proxyCallbackID,
	}

	// Test all fields are set correctly
	if proxy.ID != 1 {
		t.Errorf("ID = %d, want 1", proxy.ID)
	}
	if proxy.CallbackID != 10 {
		t.Errorf("CallbackID = %d, want 10", proxy.CallbackID)
	}
	if proxy.Port != 1080 {
		t.Errorf("Port = %d, want 1080", proxy.Port)
	}
	if proxy.PortType != "socks" {
		t.Errorf("PortType = %s, want socks", proxy.PortType)
	}
	if !proxy.Active {
		t.Error("Active should be true")
	}
	if proxy.LocalPort != 8080 {
		t.Errorf("LocalPort = %d, want 8080", proxy.LocalPort)
	}
	if proxy.RemoteIP != "192.168.1.100" {
		t.Errorf("RemoteIP = %s, want 192.168.1.100", proxy.RemoteIP)
	}
	if proxy.RemotePort != 443 {
		t.Errorf("RemotePort = %d, want 443", proxy.RemotePort)
	}
	if proxy.OperationID != 1 {
		t.Errorf("OperationID = %d, want 1", proxy.OperationID)
	}
	if proxy.ProxyCallback == nil || *proxy.ProxyCallback != 5 {
		t.Errorf("ProxyCallback = %v, want 5", proxy.ProxyCallback)
	}

	// Test helper methods
	if !proxy.IsActive() {
		t.Error("IsActive() should return true for active, non-deleted proxy")
	}
	if proxy.IsDeleted() {
		t.Error("IsDeleted() should return false")
	}

	str := proxy.String()
	if !strings.Contains(str, "socks") || !strings.Contains(str, "1080") || !strings.Contains(str, "active") {
		t.Errorf("String() = %q, should contain 'socks', '1080', and 'active'", str)
	}
}
