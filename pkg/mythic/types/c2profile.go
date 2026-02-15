package types

import (
	"fmt"
	"time"
)

// C2Profile represents a C2 communication profile in Mythic.
type C2Profile struct {
	ID           int                    `json:"id"`
	Name         string                 `json:"name"`
	Description  string                 `json:"description"`
	CreationTime time.Time              `json:"creation_time"`
	Running      bool                   `json:"running"`
	StartTime    *time.Time             `json:"start_time,omitempty"`
	StopTime     *time.Time             `json:"stop_time,omitempty"`
	Output       string                 `json:"output,omitempty"`
	StdErr       string                 `json:"std_err,omitempty"`
	StdOut       string                 `json:"std_out,omitempty"`
	Parameters   map[string]interface{} `json:"parameters,omitempty"`
	Deleted      bool                   `json:"deleted"`
	IsP2P        bool                   `json:"is_p2p"`
}

// String returns a string representation of a C2Profile.
func (c *C2Profile) String() string {
	status := "stopped"
	if c.Running {
		status = "running"
	}
	return fmt.Sprintf("%s (%s)", c.Name, status)
}

// IsRunning returns true if the C2 profile is currently running.
func (c *C2Profile) IsRunning() bool {
	return c.Running
}

// IsDeleted returns true if the C2 profile is marked as deleted.
func (c *C2Profile) IsDeleted() bool {
	return c.Deleted
}

// CreateC2InstanceRequest represents a request to create a C2 profile instance.
type CreateC2InstanceRequest struct {
	C2ProfileID  int    `json:"c2profile_id"`    // ID of the C2 profile to create an instance of
	InstanceName string `json:"instance_name"`   // Name for this instance
	C2Instance   string `json:"c2_instance"`     // JSON string of the instance configuration
}

// ImportC2InstanceRequest represents a request to import a C2 instance configuration.
type ImportC2InstanceRequest struct {
	Config string `json:"config"` // JSON string of configuration
	Name   string `json:"name"`
}

// StartStopProfileRequest represents a request to start or stop a C2 profile.
type StartStopProfileRequest struct {
	ProfileID int  `json:"profile_id"`
	Start     bool `json:"start"` // true to start, false to stop
}

// C2HostFileRequest represents a request to host a file via a C2 profile.
type C2HostFileRequest struct {
	ProfileID int    `json:"profile_id"`
	FileID    string `json:"file_id"` // UUID of the file to host
}

// C2SampleMessageRequest represents a request to generate a sample C2 message.
type C2SampleMessageRequest struct {
	ProfileID   int                    `json:"profile_id"`
	MessageType string                 `json:"message_type,omitempty"`
	Parameters  map[string]interface{} `json:"parameters,omitempty"`
}

// C2GetIOCRequest represents a request to get IOCs for a C2 profile.
type C2GetIOCRequest struct {
	ProfileID int `json:"profile_id"`
}

// C2ProfileOutput represents the output/logs from a C2 profile.
type C2ProfileOutput struct {
	Output string `json:"output"`
	StdOut string `json:"stdout"`
	StdErr string `json:"stderr"`
}

// C2SampleMessage represents a sample C2 message.
type C2SampleMessage struct {
	Message  string                 `json:"message"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

// C2IOC represents indicators of compromise for a C2 profile.
type C2IOC struct {
	ProfileName string `json:"profile_name"`
	Output      string `json:"output"`
}
