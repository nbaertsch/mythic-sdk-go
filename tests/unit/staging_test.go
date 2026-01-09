package unit

import (
	"strings"
	"testing"

	"github.com/nbaertsch/mythic-sdk-go/pkg/mythic/types"
)

func TestStagingInfo_String(t *testing.T) {
	tests := []struct {
		name     string
		staging  types.StagingInfo
		contains []string
	}{
		{
			name: "active staging",
			staging: types.StagingInfo{
				ID:          1,
				StagingUUID: "abc-123-def-456",
				Active:      true,
				Deleted:     false,
			},
			contains: []string{"abc-123-def-456", "active"},
		},
		{
			name: "inactive staging",
			staging: types.StagingInfo{
				ID:          2,
				StagingUUID: "xyz-789-uvw-012",
				Active:      false,
				Deleted:     false,
			},
			contains: []string{"xyz-789-uvw-012", "inactive"},
		},
		{
			name: "deleted staging",
			staging: types.StagingInfo{
				ID:          3,
				StagingUUID: "old-staging-uuid",
				Active:      true,
				Deleted:     true,
			},
			contains: []string{"old-staging-uuid", "deleted"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.staging.String()

			for _, substr := range tt.contains {
				if !strings.Contains(result, substr) {
					t.Errorf("String() = %q, should contain %q", result, substr)
				}
			}
		})
	}
}

func TestStagingInfo_IsActive(t *testing.T) {
	tests := []struct {
		name     string
		staging  types.StagingInfo
		expected bool
	}{
		{
			name: "active and not deleted",
			staging: types.StagingInfo{
				Active:  true,
				Deleted: false,
			},
			expected: true,
		},
		{
			name: "inactive",
			staging: types.StagingInfo{
				Active:  false,
				Deleted: false,
			},
			expected: false,
		},
		{
			name: "active but deleted",
			staging: types.StagingInfo{
				Active:  true,
				Deleted: true,
			},
			expected: false,
		},
		{
			name: "inactive and deleted",
			staging: types.StagingInfo{
				Active:  false,
				Deleted: true,
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.staging.IsActive()
			if result != tt.expected {
				t.Errorf("IsActive() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestStagingInfo_IsDeleted(t *testing.T) {
	tests := []struct {
		name     string
		staging  types.StagingInfo
		expected bool
	}{
		{
			name: "not deleted",
			staging: types.StagingInfo{
				Deleted: false,
			},
			expected: false,
		},
		{
			name: "deleted",
			staging: types.StagingInfo{
				Deleted: true,
			},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.staging.IsDeleted()
			if result != tt.expected {
				t.Errorf("IsDeleted() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestStagingInfo_HasEncryption(t *testing.T) {
	tests := []struct {
		name     string
		staging  types.StagingInfo
		expected bool
	}{
		{
			name: "has encryption key",
			staging: types.StagingInfo{
				EncryptionKey: "encryption-key-123",
			},
			expected: true,
		},
		{
			name: "has decryption key",
			staging: types.StagingInfo{
				DecryptionKey: "decryption-key-456",
			},
			expected: true,
		},
		{
			name: "has both keys",
			staging: types.StagingInfo{
				EncryptionKey: "enc-key",
				DecryptionKey: "dec-key",
			},
			expected: true,
		},
		{
			name: "no encryption",
			staging: types.StagingInfo{
				EncryptionKey: "",
				DecryptionKey: "",
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.staging.HasEncryption()
			if result != tt.expected {
				t.Errorf("HasEncryption() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestStagingInfo_HasExpiration(t *testing.T) {
	tests := []struct {
		name     string
		staging  types.StagingInfo
		expected bool
	}{
		{
			name: "has expiration time",
			staging: types.StagingInfo{
				ExpirationTime: "2026-12-31T23:59:59Z",
			},
			expected: true,
		},
		{
			name: "no expiration time",
			staging: types.StagingInfo{
				ExpirationTime: "",
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.staging.HasExpiration()
			if result != tt.expected {
				t.Errorf("HasExpiration() = %v, want %v", result, tt.expected)
			}
		})
	}
}
