package unit

import (
	"testing"
	"time"

	"github.com/nbaertsch/mythic-sdk-go/pkg/mythic/types"
)

// TestCredentialString tests the Credential.String() method
func TestCredentialString(t *testing.T) {
	tests := []struct {
		name       string
		credential types.Credential
		expected   string
	}{
		{
			name: "with account and realm",
			credential: types.Credential{
				ID:      1,
				Account: "admin",
				Realm:   "WORKGROUP",
				Type:    "plaintext",
			},
			expected: "WORKGROUP\\admin (plaintext)",
		},
		{
			name: "with account only",
			credential: types.Credential{
				ID:      2,
				Account: "user123",
				Type:    "hash",
			},
			expected: "user123 (hash)",
		},
		{
			name: "without account",
			credential: types.Credential{
				ID:   3,
				Type: "key",
			},
			expected: "Credential 3 (key)",
		},
		{
			name: "domain account",
			credential: types.Credential{
				ID:      4,
				Account: "sqlserver",
				Realm:   "DOMAIN.LOCAL",
				Type:    "plaintext",
			},
			expected: "DOMAIN.LOCAL\\sqlserver (plaintext)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.credential.String()
			if result != tt.expected {
				t.Errorf("Credential.String() = %q, expected %q", result, tt.expected)
			}
		})
	}
}

// TestCredentialIsDeleted tests the Credential.IsDeleted() method
func TestCredentialIsDeleted(t *testing.T) {
	tests := []struct {
		name     string
		deleted  bool
		expected bool
	}{
		{"deleted credential", true, true},
		{"active credential", false, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			credential := types.Credential{Deleted: tt.deleted}
			result := credential.IsDeleted()
			if result != tt.expected {
				t.Errorf("IsDeleted() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

// TestCredentialTypes tests the Credential type structure
func TestCredentialTypes(t *testing.T) {
	now := time.Now()
	taskID := 42

	credential := types.Credential{
		ID:          1,
		Type:        "plaintext",
		Account:     "administrator",
		Realm:       "WORKGROUP",
		Credential:  "P@ssw0rd123",
		Comment:     "Found during enumeration",
		OperationID: 5,
		OperatorID:  10,
		TaskID:      &taskID,
		Timestamp:   now,
		Deleted:     false,
		Metadata:    `{"source":"mimikatz"}`,
	}

	if credential.ID != 1 {
		t.Errorf("Expected ID 1, got %d", credential.ID)
	}
	if credential.Type != "plaintext" {
		t.Errorf("Expected Type 'plaintext', got %q", credential.Type)
	}
	if credential.Account != "administrator" {
		t.Errorf("Expected Account 'administrator', got %q", credential.Account)
	}
	if credential.TaskID == nil || *credential.TaskID != 42 {
		t.Errorf("Expected TaskID 42, got %v", credential.TaskID)
	}
	if credential.IsDeleted() {
		t.Error("Expected credential not to be deleted")
	}
}

// TestCreateCredentialRequestTypes tests the CreateCredentialRequest type structure
func TestCreateCredentialRequestTypes(t *testing.T) {
	req := types.CreateCredentialRequest{
		Type:       "hash",
		Account:    "user1",
		Realm:      "DOMAIN",
		Credential: "aad3b435b51404eeaad3b435b51404ee",
		Comment:    "NTLM hash",
		Metadata:   `{"hash_type":"ntlm"}`,
	}

	if req.Type != "hash" {
		t.Errorf("Expected Type 'hash', got %q", req.Type)
	}
	if req.Account != "user1" {
		t.Errorf("Expected Account 'user1', got %q", req.Account)
	}
	if req.Credential != "aad3b435b51404eeaad3b435b51404ee" {
		t.Errorf("Expected Credential to be set, got %q", req.Credential)
	}
}

// TestUpdateCredentialRequestTypes tests the UpdateCredentialRequest type structure
func TestUpdateCredentialRequestTypes(t *testing.T) {
	newType := "plaintext"
	newComment := "Updated comment"
	deleted := true

	req := types.UpdateCredentialRequest{
		ID:      1,
		Type:    &newType,
		Comment: &newComment,
		Deleted: &deleted,
	}

	if req.ID != 1 {
		t.Errorf("Expected ID 1, got %d", req.ID)
	}
	if req.Type == nil || *req.Type != "plaintext" {
		t.Error("Expected Type to be 'plaintext'")
	}
	if req.Comment == nil || *req.Comment != "Updated comment" {
		t.Error("Expected Comment to be set")
	}
	if req.Deleted == nil || !*req.Deleted {
		t.Error("Expected Deleted to be true")
	}
}

// TestCredentialTypeConstants tests the CredentialType constants
func TestCredentialTypeConstants(t *testing.T) {
	credTypes := []types.CredentialType{
		types.CredentialTypePlaintext,
		types.CredentialTypeHash,
		types.CredentialTypeKey,
		types.CredentialTypeTicket,
		types.CredentialTypeCookie,
		types.CredentialTypeCertificate,
	}

	for _, credType := range credTypes {
		if credType == "" {
			t.Errorf("Credential type should not be empty")
		}
		t.Logf("Credential type: %s", credType)
	}

	// Verify specific values
	if types.CredentialTypePlaintext != "plaintext" {
		t.Errorf("Expected 'plaintext', got %q", types.CredentialTypePlaintext)
	}
	if types.CredentialTypeHash != "hash" {
		t.Errorf("Expected 'hash', got %q", types.CredentialTypeHash)
	}
	if types.CredentialTypeKey != "key" {
		t.Errorf("Expected 'key', got %q", types.CredentialTypeKey)
	}
}

// TestCredentialWithNilFields tests handling of nil fields
func TestCredentialWithNilFields(t *testing.T) {
	credential := types.Credential{
		ID:      1,
		Type:    "plaintext",
		Account: "user",
		Realm:   "domain",
		TaskID:  nil, // No associated task
	}

	if credential.TaskID != nil {
		t.Error("Expected TaskID to be nil")
	}

	str := credential.String()
	if str == "" {
		t.Error("String() should not return empty string")
	}
}

// TestCreateCredentialRequestValidation tests request validation
func TestCreateCredentialRequestValidation(t *testing.T) {
	tests := []struct {
		name  string
		req   types.CreateCredentialRequest
		valid bool
	}{
		{
			name: "valid request",
			req: types.CreateCredentialRequest{
				Type:       "plaintext",
				Account:    "admin",
				Realm:      "WORKGROUP",
				Credential: "password123",
			},
			valid: true,
		},
		{
			name: "missing type",
			req: types.CreateCredentialRequest{
				Account:    "admin",
				Credential: "password123",
			},
			valid: false,
		},
		{
			name: "missing account",
			req: types.CreateCredentialRequest{
				Type:       "plaintext",
				Credential: "password123",
			},
			valid: false,
		},
		{
			name: "missing credential",
			req: types.CreateCredentialRequest{
				Type:    "plaintext",
				Account: "admin",
			},
			valid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hasRequiredFields := tt.req.Type != "" && tt.req.Account != "" && tt.req.Credential != ""
			if hasRequiredFields != tt.valid {
				t.Errorf("Request validity mismatch: expected %v, got %v", tt.valid, hasRequiredFields)
			}
		})
	}
}

// TestUpdateCredentialRequestPartial tests partial updates
func TestUpdateCredentialRequestPartial(t *testing.T) {
	// Test that we can update just one field
	newComment := "Only updating comment"
	req := types.UpdateCredentialRequest{
		ID:      1,
		Comment: &newComment,
	}

	if req.ID != 1 {
		t.Errorf("Expected ID 1, got %d", req.ID)
	}
	if req.Comment == nil {
		t.Error("Comment should be set")
	}
	if req.Type != nil {
		t.Error("Type should be nil")
	}
	if req.Account != nil {
		t.Error("Account should be nil")
	}
}

// TestCredentialMetadata tests metadata handling
func TestCredentialMetadata(t *testing.T) {
	metadata := `{"source":"mimikatz","module":"sekurlsa::logonpasswords"}`

	credential := types.Credential{
		ID:       1,
		Type:     "plaintext",
		Account:  "admin",
		Metadata: metadata,
	}

	if credential.Metadata != metadata {
		t.Errorf("Expected metadata %q, got %q", metadata, credential.Metadata)
	}
}
