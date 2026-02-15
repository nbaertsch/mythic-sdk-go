package types

import (
	"fmt"
	"time"
)

// Credential represents a compromised credential in Mythic.
type Credential struct {
	ID          int        `json:"id"`
	Type        string     `json:"type"`
	Account     string     `json:"account"`
	Realm       string     `json:"realm"`
	Credential  string     `json:"credential_text"`
	Comment     string     `json:"comment"`
	OperationID int        `json:"operation_id"`
	OperatorID  int        `json:"operator_id"`
	TaskID      *int       `json:"task_id,omitempty"`
	Timestamp   time.Time  `json:"timestamp"`
	Deleted     bool       `json:"deleted"`
	Metadata    string     `json:"metadata"`
	Operation   *Operation `json:"operation,omitempty"`
	Operator    *Operator  `json:"operator,omitempty"`
}

// CreateCredentialRequest represents a request to create a new credential.
type CreateCredentialRequest struct {
	Type       string `json:"type"`               // e.g., "plaintext", "hash", "key", "ticket", "cookie"
	Account    string `json:"account"`            // Username or account name
	Realm      string `json:"realm"`              // Domain, hostname, or scope
	Credential string `json:"credential_text"`    // The actual credential value
	Comment    string `json:"comment,omitempty"`  // Optional comment
	Metadata   string `json:"metadata,omitempty"` // Optional JSON metadata
	// Note: task_id is populated automatically by Mythic when credentials are
	// reported by C2 agents through task responses. It cannot be set manually
	// via the createCredential mutation.
}

// UpdateCredentialRequest represents a request to update a credential.
type UpdateCredentialRequest struct {
	ID         int     `json:"id"`
	Type       *string `json:"type,omitempty"`
	Account    *string `json:"account,omitempty"`
	Realm      *string `json:"realm,omitempty"`
	Credential *string `json:"credential_text,omitempty"`
	Comment    *string `json:"comment,omitempty"`
	Deleted    *bool   `json:"deleted,omitempty"`
	Metadata   *string `json:"metadata,omitempty"`
}

// String returns a string representation of a Credential.
func (c *Credential) String() string {
	if c.Account != "" && c.Realm != "" {
		return fmt.Sprintf("%s\\%s (%s)", c.Realm, c.Account, c.Type)
	}
	if c.Account != "" {
		return fmt.Sprintf("%s (%s)", c.Account, c.Type)
	}
	return fmt.Sprintf("Credential %d (%s)", c.ID, c.Type)
}

// IsDeleted returns true if the credential is marked as deleted.
func (c *Credential) IsDeleted() bool {
	return c.Deleted
}
