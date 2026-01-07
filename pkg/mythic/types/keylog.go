package types

import (
	"fmt"
	"time"
)

// Keylog represents a keylog entry captured from a callback.
type Keylog struct {
	ID          int        `json:"id"`
	TaskID      int        `json:"task_id"`
	Keystrokes  string     `json:"keystrokes"`
	Window      string     `json:"window"`
	Timestamp   time.Time  `json:"timestamp"`
	OperationID int        `json:"operation_id"`
	User        string     `json:"user"`
	CallbackID  int        `json:"callback_id"`
	Callback    *Callback  `json:"callback,omitempty"`
	Operation   *Operation `json:"operation,omitempty"`
}

// String returns a string representation of a Keylog.
func (k *Keylog) String() string {
	if k.Window != "" {
		return fmt.Sprintf("%s - %s (%s)", k.Timestamp.Format("2006-01-02 15:04:05"), k.Window, k.User)
	}
	if k.User != "" {
		return fmt.Sprintf("%s - %s", k.Timestamp.Format("2006-01-02 15:04:05"), k.User)
	}
	return fmt.Sprintf("Keylog %d (%s)", k.ID, k.Timestamp.Format("2006-01-02 15:04:05"))
}

// HasKeystrokes returns true if the keylog has captured keystrokes.
func (k *Keylog) HasKeystrokes() bool {
	return k.Keystrokes != ""
}
