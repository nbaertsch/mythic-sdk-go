package types

import "time"

// TaskStatus represents the status of a task.
type TaskStatus string

const (
	TaskStatusSubmitted     TaskStatus = "submitted"
	TaskStatusProcessing    TaskStatus = "processing"
	TaskStatusProcessed     TaskStatus = "processed"
	TaskStatusCompleted     TaskStatus = "completed"
	TaskStatusError         TaskStatus = "error"
	TaskStatusPreprocessing TaskStatus = "preprocessing"
	TaskStatusDelegating    TaskStatus = "delegating"
)

// CallbackIntegrityLevel represents Windows integrity levels.
type CallbackIntegrityLevel int

const (
	IntegrityLevelLow    CallbackIntegrityLevel = 2
	IntegrityLevelMedium CallbackIntegrityLevel = 3
	IntegrityLevelHigh   CallbackIntegrityLevel = 4
	IntegrityLevelSystem CallbackIntegrityLevel = 5
)

// OperatorViewMode represents an operator's access level in an operation.
type OperatorViewMode string

const (
	ViewModeOperator  OperatorViewMode = "operator"
	ViewModeSpectator OperatorViewMode = "spectator"
	ViewModeLead      OperatorViewMode = "lead"
)

// CredentialType represents the type of credential.
type CredentialType string

const (
	CredentialTypePlaintext   CredentialType = "plaintext"
	CredentialTypeHash        CredentialType = "hash"
	CredentialTypeKey         CredentialType = "key"
	CredentialTypeTicket      CredentialType = "ticket"
	CredentialTypeCookie      CredentialType = "cookie"
	CredentialTypeCertificate CredentialType = "certificate"
)

// Timestamp is a custom time type that handles Mythic's timestamp format.
type Timestamp struct {
	time.Time
}

// UnmarshalJSON implements json.Unmarshaler for Timestamp.
func (t *Timestamp) UnmarshalJSON(data []byte) error {
	// Remove quotes
	s := string(data[1 : len(data)-1])

	// Try RFC3339 format first
	parsed, err := time.Parse(time.RFC3339, s)
	if err == nil {
		t.Time = parsed
		return nil
	}

	// Try other common formats
	formats := []string{
		time.RFC3339Nano,
		"2006-01-02T15:04:05.999999Z07:00",
		"2006-01-02T15:04:05Z",
	}

	for _, format := range formats {
		parsed, err = time.Parse(format, s)
		if err == nil {
			t.Time = parsed
			return nil
		}
	}

	return err
}

// MarshalJSON implements json.Marshaler for Timestamp.
func (t Timestamp) MarshalJSON() ([]byte, error) {
	return []byte(`"` + t.Time.Format(time.RFC3339) + `"`), nil
}
