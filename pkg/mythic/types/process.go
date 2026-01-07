package types

import (
	"fmt"
	"time"
)

// Process represents a process tracked by Mythic from a callback.
type Process struct {
	ID              int        `json:"id"`
	Name            string     `json:"name"`
	ProcessID       int        `json:"process_id"`
	ParentProcessID int        `json:"parent_process_id"`
	Architecture    string     `json:"architecture"`
	BinPath         string     `json:"bin_path"`
	User            string     `json:"user"`
	CommandLine     string     `json:"command_line"`
	IntegrityLevel  int        `json:"integrity_level"`
	StartTime       time.Time  `json:"start_time"`
	Description     string     `json:"description"`
	OperationID     int        `json:"operation_id"`
	HostID          int        `json:"host_id"`
	CallbackID      *int       `json:"callback_id,omitempty"`
	TaskID          *int       `json:"task_id,omitempty"`
	Timestamp       time.Time  `json:"timestamp"`
	Deleted         bool       `json:"deleted"`
	Host            *Host      `json:"host,omitempty"`
	Operation       *Operation `json:"operation,omitempty"`
	Children        []*Process `json:"children,omitempty"` // For tree structure
}

// Host represents a host in Mythic (simplified for process context).
type Host struct {
	ID   int    `json:"id"`
	Host string `json:"host"`
}

// ProcessTree represents a hierarchical process tree.
type ProcessTree struct {
	Process  *Process
	Children []*ProcessTree
}

// String returns a string representation of a Process.
func (p *Process) String() string {
	if p.Name != "" && p.ProcessID != 0 {
		return fmt.Sprintf("%s (PID: %d)", p.Name, p.ProcessID)
	}
	if p.ProcessID != 0 {
		return fmt.Sprintf("PID %d", p.ProcessID)
	}
	return fmt.Sprintf("Process %d", p.ID)
}

// IsDeleted returns true if the process is marked as deleted.
func (p *Process) IsDeleted() bool {
	return p.Deleted
}

// HasParent returns true if the process has a parent process.
func (p *Process) HasParent() bool {
	return p.ParentProcessID != 0
}

// GetIntegrityLevelString returns a human-readable integrity level.
func (p *Process) GetIntegrityLevelString() string {
	switch p.IntegrityLevel {
	case 0:
		return "Untrusted"
	case 1:
		return "Low"
	case 2:
		return "Medium"
	case 3:
		return "High"
	case 4:
		return "System"
	default:
		return fmt.Sprintf("Unknown (%d)", p.IntegrityLevel)
	}
}
