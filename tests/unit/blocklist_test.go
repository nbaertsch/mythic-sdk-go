package unit

import (
	"strings"
	"testing"

	"github.com/nbaertsch/mythic-sdk-go/pkg/mythic/types"
)

func TestBlockList_String(t *testing.T) {
	tests := []struct {
		name      string
		blocklist types.BlockList
		contains  []string
	}{
		{
			name: "active block list",
			blocklist: types.BlockList{
				ID:      1,
				Name:    "Security Tools",
				Active:  true,
				Deleted: false,
			},
			contains: []string{"Security Tools", "active"},
		},
		{
			name: "inactive block list",
			blocklist: types.BlockList{
				ID:      2,
				Name:    "Test List",
				Active:  false,
				Deleted: false,
			},
			contains: []string{"Test List", "inactive"},
		},
		{
			name: "deleted block list",
			blocklist: types.BlockList{
				ID:      3,
				Name:    "Old List",
				Active:  true,
				Deleted: true,
			},
			contains: []string{"Old List", "deleted"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.blocklist.String()

			for _, substr := range tt.contains {
				if !strings.Contains(result, substr) {
					t.Errorf("String() = %q, should contain %q", result, substr)
				}
			}
		})
	}
}

func TestBlockList_IsActive(t *testing.T) {
	tests := []struct {
		name      string
		blocklist types.BlockList
		expected  bool
	}{
		{
			name: "active and not deleted",
			blocklist: types.BlockList{
				Active:  true,
				Deleted: false,
			},
			expected: true,
		},
		{
			name: "inactive",
			blocklist: types.BlockList{
				Active:  false,
				Deleted: false,
			},
			expected: false,
		},
		{
			name: "active but deleted",
			blocklist: types.BlockList{
				Active:  true,
				Deleted: true,
			},
			expected: false,
		},
		{
			name: "inactive and deleted",
			blocklist: types.BlockList{
				Active:  false,
				Deleted: true,
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.blocklist.IsActive()
			if result != tt.expected {
				t.Errorf("IsActive() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestBlockList_IsDeleted(t *testing.T) {
	tests := []struct {
		name      string
		blocklist types.BlockList
		expected  bool
	}{
		{
			name: "not deleted",
			blocklist: types.BlockList{
				Deleted: false,
			},
			expected: false,
		},
		{
			name: "deleted",
			blocklist: types.BlockList{
				Deleted: true,
			},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.blocklist.IsDeleted()
			if result != tt.expected {
				t.Errorf("IsDeleted() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestBlockListEntry_String(t *testing.T) {
	tests := []struct {
		name     string
		entry    types.BlockListEntry
		contains []string
	}{
		{
			name: "active IP entry",
			entry: types.BlockListEntry{
				Type:    "ip",
				Value:   "192.168.1.1",
				Active:  true,
				Deleted: false,
			},
			contains: []string{"ip", "192.168.1.1", "active"},
		},
		{
			name: "inactive domain entry",
			entry: types.BlockListEntry{
				Type:    "domain",
				Value:   "scanner.example.com",
				Active:  false,
				Deleted: false,
			},
			contains: []string{"domain", "scanner.example.com", "inactive"},
		},
		{
			name: "deleted user agent entry",
			entry: types.BlockListEntry{
				Type:    "user_agent",
				Value:   "SecurityScanner/1.0",
				Active:  true,
				Deleted: true,
			},
			contains: []string{"user_agent", "SecurityScanner/1.0", "deleted"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.entry.String()

			for _, substr := range tt.contains {
				if !strings.Contains(result, substr) {
					t.Errorf("String() = %q, should contain %q", result, substr)
				}
			}
		})
	}
}

func TestBlockListEntry_IsActive(t *testing.T) {
	tests := []struct {
		name     string
		entry    types.BlockListEntry
		expected bool
	}{
		{
			name: "active and not deleted",
			entry: types.BlockListEntry{
				Active:  true,
				Deleted: false,
			},
			expected: true,
		},
		{
			name: "inactive",
			entry: types.BlockListEntry{
				Active:  false,
				Deleted: false,
			},
			expected: false,
		},
		{
			name: "active but deleted",
			entry: types.BlockListEntry{
				Active:  true,
				Deleted: true,
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.entry.IsActive()
			if result != tt.expected {
				t.Errorf("IsActive() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestDeleteBlockListRequest_String(t *testing.T) {
	request := types.DeleteBlockListRequest{
		BlockListID: 42,
	}

	result := request.String()

	if !strings.Contains(result, "42") {
		t.Errorf("String() = %q, should contain block list ID 42", result)
	}
	if !strings.Contains(result, "Delete") {
		t.Errorf("String() = %q, should contain 'Delete'", result)
	}
}

func TestDeleteBlockListResponse_String(t *testing.T) {
	tests := []struct {
		name     string
		response types.DeleteBlockListResponse
		contains []string
	}{
		{
			name: "successful deletion",
			response: types.DeleteBlockListResponse{
				Status:  "success",
				Message: "Block list removed",
			},
			contains: []string{"Block list removed"},
		},
		{
			name: "successful without message",
			response: types.DeleteBlockListResponse{
				Status: "success",
			},
			contains: []string{"deleted successfully"},
		},
		{
			name: "failed deletion",
			response: types.DeleteBlockListResponse{
				Status: "error",
				Error:  "Block list not found",
			},
			contains: []string{"Failed", "Block list not found"},
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

func TestDeleteBlockListResponse_IsSuccessful(t *testing.T) {
	tests := []struct {
		name     string
		response types.DeleteBlockListResponse
		expected bool
	}{
		{
			name: "success status",
			response: types.DeleteBlockListResponse{
				Status: "success",
			},
			expected: true,
		},
		{
			name: "error status",
			response: types.DeleteBlockListResponse{
				Status: "error",
			},
			expected: false,
		},
		{
			name: "failed status",
			response: types.DeleteBlockListResponse{
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

func TestDeleteBlockListEntryRequest_String(t *testing.T) {
	tests := []struct {
		name     string
		request  types.DeleteBlockListEntryRequest
		contains []string
	}{
		{
			name: "single entry",
			request: types.DeleteBlockListEntryRequest{
				EntryIDs: []int{123},
			},
			contains: []string{"Delete", "1", "entries"},
		},
		{
			name: "multiple entries",
			request: types.DeleteBlockListEntryRequest{
				EntryIDs: []int{123, 456, 789},
			},
			contains: []string{"Delete", "3", "entries"},
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

func TestDeleteBlockListEntryResponse_String(t *testing.T) {
	tests := []struct {
		name     string
		response types.DeleteBlockListEntryResponse
		contains []string
	}{
		{
			name: "successful deletion",
			response: types.DeleteBlockListEntryResponse{
				Status:       "success",
				DeletedCount: 5,
			},
			contains: []string{"Deleted", "5", "entries"},
		},
		{
			name: "failed deletion",
			response: types.DeleteBlockListEntryResponse{
				Status: "error",
				Error:  "Entries not found",
			},
			contains: []string{"Failed", "Entries not found"},
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

func TestDeleteBlockListEntryResponse_IsSuccessful(t *testing.T) {
	tests := []struct {
		name     string
		response types.DeleteBlockListEntryResponse
		expected bool
	}{
		{
			name: "success status",
			response: types.DeleteBlockListEntryResponse{
				Status: "success",
			},
			expected: true,
		},
		{
			name: "error status",
			response: types.DeleteBlockListEntryResponse{
				Status: "error",
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
