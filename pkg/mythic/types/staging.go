package types

import "fmt"

// StagingInfo represents payload staging information for delivery.
// Staging is used when payloads need to be delivered in multiple parts
// or when the full payload is hosted on the C2 infrastructure.
type StagingInfo struct {
	ID             int    `json:"id"`
	PayloadID      int    `json:"payload_id"`
	StagingUUID    string `json:"staging_uuid"`
	EncryptionKey  string `json:"encryption_key,omitempty"`
	DecryptionKey  string `json:"decryption_key,omitempty"`
	C2ProfileID    int    `json:"c2profile_id"`
	OperationID    int    `json:"operation_id"`
	CreationTime   string `json:"creation_time"`
	ExpirationTime string `json:"expiration_time,omitempty"`
	Active         bool   `json:"active"`
	Deleted        bool   `json:"deleted"`
	OperatorID     int    `json:"operator_id"`

	// Nested relationships
	Payload   *Payload   `json:"payload,omitempty"`
	C2Profile *C2Profile `json:"c2profile,omitempty"`
	Operation *Operation `json:"operation,omitempty"`
}

// String returns a human-readable representation of the staging info.
func (s *StagingInfo) String() string {
	status := "active"
	if s.Deleted {
		status = "deleted"
	} else if !s.Active {
		status = "inactive"
	}

	return fmt.Sprintf("Staging UUID: %s (%s)", s.StagingUUID, status)
}

// IsActive returns true if the staging info is active and not deleted.
func (s *StagingInfo) IsActive() bool {
	return s.Active && !s.Deleted
}

// IsDeleted returns true if the staging info has been deleted.
func (s *StagingInfo) IsDeleted() bool {
	return s.Deleted
}

// HasEncryption returns true if staging uses encryption.
func (s *StagingInfo) HasEncryption() bool {
	return s.EncryptionKey != "" || s.DecryptionKey != ""
}

// IsExpired returns true if the staging has an expiration time set.
// Note: This only checks if expiration time is set, not if it's past.
func (s *StagingInfo) HasExpiration() bool {
	return s.ExpirationTime != ""
}
