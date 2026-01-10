package types

import (
	"fmt"
	"time"
)

// Operation represents a Mythic operation (engagement/campaign).
type Operation struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Complete    bool      `json:"complete"`
	Webhook     string    `json:"webhook"`
	Channel     string    `json:"channel"`
	AdminID     int       `json:"admin_id"`
	BannerText  string    `json:"banner_text"`
	BannerColor string    `json:"banner_color"`
	Admin       *Operator `json:"admin,omitempty"`
}

// OperationOperator represents the relationship between an operation and an operator.
type OperationOperator struct {
	ID          int        `json:"id"`
	OperationID int        `json:"operation_id"`
	OperatorID  int        `json:"operator_id"`
	Operation   *Operation `json:"operation,omitempty"`
	Operator    *Operator  `json:"operator,omitempty"`
}

// Operator represents a Mythic operator (user).
type Operator struct {
	ID                 int        `json:"id"`
	Username           string     `json:"username"`
	Admin              bool       `json:"admin"`
	Active             bool       `json:"active"`
	CreationTime       time.Time  `json:"creation_time"`
	LastLogin          time.Time  `json:"last_login"`
	CurrentOperationID *int       `json:"current_operation_id"`
	ViewUTCTime        bool       `json:"view_utc_time"`
	Deleted            bool       `json:"deleted"`
	CurrentOperation   *Operation `json:"current_operation,omitempty"`
	AccountType        string     `json:"account_type,omitempty"`
	FailedLoginCount   int        `json:"failed_login_count"`
}

// OperationEventLog represents an event log entry for an operation.
type OperationEventLog struct {
	ID          int        `json:"id"`
	OperatorID  int        `json:"operator_id"`
	OperationID int        `json:"operation_id"`
	Message     string     `json:"message"`
	Timestamp   time.Time  `json:"timestamp"`
	Level       string     `json:"level"`
	Source      string     `json:"source"`
	Deleted     bool       `json:"deleted"`
	Operator    *Operator  `json:"operator,omitempty"`
	Operation   *Operation `json:"operation,omitempty"`
}

// CreateOperationRequest represents a request to create a new operation.
type CreateOperationRequest struct {
	Name    string `json:"name"`
	AdminID *int   `json:"admin_id,omitempty"`
	Channel string `json:"channel,omitempty"`
	Webhook string `json:"webhook,omitempty"`
}

// UpdateOperationRequest represents a request to update an operation.
type UpdateOperationRequest struct {
	OperationID int     `json:"operation_id"`
	Name        *string `json:"name,omitempty"`
	Channel     *string `json:"channel,omitempty"`
	Complete    *bool   `json:"complete,omitempty"`
	Webhook     *string `json:"webhook,omitempty"`
	AdminID     *int    `json:"admin_id,omitempty"`
	BannerText  *string `json:"banner_text,omitempty"`
	BannerColor *string `json:"banner_color,omitempty"`
}

// UpdateOperatorOperationRequest represents a request to update an operator's role in an operation.
type UpdateOperatorOperationRequest struct {
	OperatorID  int  `json:"operator_id"`
	OperationID int  `json:"operation_id"`
	Remove      bool `json:"remove,omitempty"` // If true, remove operator from operation
}

// CreateOperationEventLogRequest represents a request to create an event log entry.
type CreateOperationEventLogRequest struct {
	OperationID int    `json:"operation_id"`
	Message     string `json:"message"`
	Level       string `json:"level,omitempty"`  // info, warning, error
	Source      string `json:"source,omitempty"` // Source of the event
}

// String returns a string representation of an Operation.
func (o *Operation) String() string {
	status := "active"
	if o.Complete {
		status = "complete"
	}
	return fmt.Sprintf("%s (%s)", o.Name, status)
}

// IsComplete returns true if the operation is marked as complete.
func (o *Operation) IsComplete() bool {
	return o.Complete
}

// String returns a string representation of an OperationEventLog.
func (e *OperationEventLog) String() string {
	return fmt.Sprintf("[%s] %s: %s", e.Level, e.Timestamp.Format("2006-01-02 15:04:05"), e.Message)
}

// CreateOperatorRequest represents a request to create a new operator.
type CreateOperatorRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// UpdateOperatorStatusRequest represents a request to update operator status.
type UpdateOperatorStatusRequest struct {
	OperatorID int   `json:"operator_id"`
	Active     *bool `json:"active,omitempty"`
	Admin      *bool `json:"admin,omitempty"`
	Deleted    *bool `json:"deleted,omitempty"`
}

// UpdatePasswordAndEmailRequest represents a request to update operator credentials.
type UpdatePasswordAndEmailRequest struct {
	OperatorID  int     `json:"operator_id"`
	OldPassword string  `json:"old_password"`
	NewPassword *string `json:"new_password,omitempty"`
	Email       *string `json:"email,omitempty"`
}

// OperatorPreferences represents UI preferences for an operator.
type OperatorPreferences struct {
	OperatorID      int                    `json:"operator_id"`
	PreferencesJSON string                 `json:"preferences_json,omitempty"`
	Preferences     map[string]interface{} `json:"preferences,omitempty"`
	InteractType    string                 `json:"interact_type,omitempty"`
	ConsoleSize     int                    `json:"console_size,omitempty"`
	FontSize        int                    `json:"font_size,omitempty"`
}

// UpdateOperatorPreferencesRequest represents a request to update operator preferences.
type UpdateOperatorPreferencesRequest struct {
	OperatorID  int                    `json:"operator_id"`
	Preferences map[string]interface{} `json:"preferences"`
}

// OperatorSecrets represents secrets/keys associated with an operator.
type OperatorSecrets struct {
	OperatorID  int                    `json:"operator_id"`
	SecretsJSON string                 `json:"secrets_json,omitempty"`
	Secrets     map[string]interface{} `json:"secrets,omitempty"`
}

// UpdateOperatorSecretsRequest represents a request to update operator secrets.
type UpdateOperatorSecretsRequest struct {
	OperatorID int                    `json:"operator_id"`
	Secrets    map[string]interface{} `json:"secrets"`
}

// InviteLink represents an invitation link for new operators.
type InviteLink struct {
	ID          int       `json:"id"`
	Code        string    `json:"code"`
	ExpiresAt   time.Time `json:"expires_at"`
	CreatedBy   int       `json:"created_by"`
	CreatedAt   time.Time `json:"created_at"`
	MaxUses     int       `json:"max_uses"`
	CurrentUses int       `json:"current_uses"`
	Active      bool      `json:"active"`
}

// CreateInviteLinkRequest represents a request to create an invite link.
type CreateInviteLinkRequest struct {
	MaxUses   int       `json:"max_uses"`
	ExpiresAt time.Time `json:"expires_at"`
}

// String returns a string representation of an Operator.
func (o *Operator) String() string {
	role := "Operator"
	if o.Admin {
		role = "Admin"
	}
	status := ""
	if o.Deleted {
		status = " (deleted)"
	} else if !o.Active {
		status = " (inactive)"
	}
	return fmt.Sprintf("%s (%s)%s", o.Username, role, status)
}

// IsAdmin returns true if the operator has admin privileges.
func (o *Operator) IsAdmin() bool {
	return o.Admin
}

// IsActive returns true if the operator account is active.
func (o *Operator) IsActive() bool {
	return o.Active && !o.Deleted
}

// IsDeleted returns true if the operator has been deleted.
func (o *Operator) IsDeleted() bool {
	return o.Deleted
}

// IsLocked returns true if the operator account is locked due to failed login attempts.
func (o *Operator) IsLocked() bool {
	return o.FailedLoginCount >= 10
}

// IsBotAccount returns true if this is a bot account.
func (o *Operator) IsBotAccount() bool {
	return o.AccountType == AccountTypeBot
}

// String returns a string representation of an InviteLink.
func (i *InviteLink) String() string {
	return fmt.Sprintf("Invite %s (uses: %d/%d)", i.Code, i.CurrentUses, i.MaxUses)
}

// IsExpired returns true if the invite link has expired.
func (i *InviteLink) IsExpired() bool {
	return time.Now().After(i.ExpiresAt)
}

// IsActive returns true if the invite link is active and not expired.
func (i *InviteLink) IsActive() bool {
	return i.Active && !i.IsExpired()
}

// HasUsesRemaining returns true if the invite link has uses remaining.
func (i *InviteLink) HasUsesRemaining() bool {
	return i.CurrentUses < i.MaxUses
}

// Operator account type constants
const (
	AccountTypeUser = "user"
	AccountTypeBot  = "bot"
)
