package types

import "time"

// Callback represents a Mythic callback (active agent connection).
type Callback struct {
	// ID is the database ID
	ID int `json:"id"`

	// DisplayID is the user-friendly display ID
	DisplayID int `json:"display_id"`

	// AgentCallbackID is the agent's callback identifier
	AgentCallbackID string `json:"agent_callback_id"`

	// InitCallback is the initial checkin timestamp
	InitCallback time.Time `json:"init_callback"`

	// LastCheckin is the most recent checkin timestamp
	LastCheckin time.Time `json:"last_checkin"`

	// User is the username context the callback is running as
	User string `json:"user"`

	// Host is the hostname of the compromised system
	Host string `json:"host"`

	// PID is the process ID
	PID int `json:"pid"`

	// IP is the list of IP addresses
	IP []string `json:"ip"`

	// ExternalIP is the external/public IP address
	ExternalIP string `json:"external_ip"`

	// ProcessName is the name of the process
	ProcessName string `json:"process_name"`

	// Description is an operator-set description
	Description string `json:"description"`

	// OperatorID is the ID of the operator who received the callback
	OperatorID int `json:"operator_id"`

	// Active indicates if the callback is currently active
	Active bool `json:"active"`

	// RegisteredPayloadID is the ID of the payload that created this callback
	RegisteredPayloadID string `json:"registered_payload_id"`

	// IntegrityLevel is the Windows integrity level (2=low, 3=medium, 4=high, 5=system)
	IntegrityLevel CallbackIntegrityLevel `json:"integrity_level"`

	// Locked indicates if the callback is locked for tasking
	Locked bool `json:"locked"`

	// OperationID is the ID of the operation this callback belongs to
	OperationID int `json:"operation_id"`

	// CryptoType is the encryption type used
	CryptoType string `json:"crypto_type"`

	// DecKey is the decryption key (base64 encoded)
	DecKey *string `json:"dec_key,omitempty"`

	// EncKey is the encryption key (base64 encoded)
	EncKey *string `json:"enc_key,omitempty"`

	// OS is the operating system
	OS string `json:"os"`

	// Architecture is the system architecture (x86, x64, arm64, etc.)
	Architecture string `json:"architecture"`

	// Domain is the domain name (Windows)
	Domain string `json:"domain"`

	// ExtraInfo is additional metadata about the callback
	ExtraInfo string `json:"extra_info"`

	// SleepInfo is the sleep/jitter configuration
	SleepInfo string `json:"sleep_info"`

	// PayloadTypeID is the ID of the payload type
	PayloadTypeID int `json:"payload_type_id"`

	// C2ProfileID is the ID of the C2 profile
	C2ProfileID *int `json:"c2_profile_id,omitempty"`

	// Payload is the associated payload information (if loaded)
	Payload *CallbackPayload `json:"payload,omitempty"`

	// PayloadType is the payload type information (if loaded)
	PayloadType *CallbackPayloadType `json:"payloadtype,omitempty"`

	// Operator is the operator information (if loaded)
	Operator *CallbackOperator `json:"operator,omitempty"`
}

// CallbackPayload represents minimal payload information in a callback.
type CallbackPayload struct {
	ID          int    `json:"id"`
	UUID        string `json:"uuid"`
	Description string `json:"description"`
	OS          string `json:"os"`
}

// CallbackPayloadType represents payload type information in a callback.
type CallbackPayloadType struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// CallbackOperator represents operator information in a callback.
type CallbackOperator struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
}

// CallbackUpdateRequest represents a request to update callback properties.
type CallbackUpdateRequest struct {
	// CallbackDisplayID is the display ID of the callback to update
	CallbackDisplayID int

	// Active sets the active status
	Active *bool

	// Locked sets the locked status
	Locked *bool

	// Description sets the description
	Description *string

	// IPs sets the IP addresses
	IPs []string

	// User sets the username
	User *string

	// Host sets the hostname
	Host *string

	// OS sets the operating system
	OS *string

	// Architecture sets the architecture
	Architecture *string

	// ExtraInfo sets extra information
	ExtraInfo *string

	// SleepInfo sets sleep/jitter info
	SleepInfo *string

	// PID sets the process ID
	PID *int

	// ProcessName sets the process name
	ProcessName *string

	// IntegrityLevel sets the integrity level
	IntegrityLevel *CallbackIntegrityLevel

	// Domain sets the domain
	Domain *string
}

// String returns a string representation of the callback.
func (c *Callback) String() string {
	status := "inactive"
	if c.Active {
		status = "active"
	}
	return c.User + "@" + c.Host + " (" + c.OS + ", " + status + ")"
}

// IsHigh returns true if the callback has high or system integrity level.
func (c *Callback) IsHigh() bool {
	return c.IntegrityLevel >= IntegrityLevelHigh
}

// IsSystem returns true if the callback has system integrity level.
func (c *Callback) IsSystem() bool {
	return c.IntegrityLevel == IntegrityLevelSystem
}
