package types

import (
	"fmt"
)

// Attack represents a MITRE ATT&CK technique.
type Attack struct {
	ID     int    `json:"id"`
	TNum   string `json:"t_num"`
	Name   string `json:"name"`
	OS     string `json:"os"`
	Tactic string `json:"tactic"`
}

// String returns a string representation of an Attack technique.
func (a *Attack) String() string {
	if a.Name != "" && a.TNum != "" {
		return fmt.Sprintf("%s (%s)", a.TNum, a.Name)
	}
	if a.TNum != "" {
		return a.TNum
	}
	return fmt.Sprintf("Attack %d", a.ID)
}

// AttackTask represents the association between a task and MITRE ATT&CK technique.
type AttackTask struct {
	ID       int     `json:"id"`
	AttackID int     `json:"attack_id"`
	TaskID   int     `json:"task_id"`
	Attack   *Attack `json:"attack,omitempty"`
}

// String returns a string representation of an AttackTask.
func (at *AttackTask) String() string {
	if at.Attack != nil {
		return fmt.Sprintf("%s on Task %d", at.Attack.String(), at.TaskID)
	}
	return fmt.Sprintf("AttackTask %d (Attack %d, Task %d)", at.ID, at.AttackID, at.TaskID)
}

// AttackCommand represents the association between a command and MITRE ATT&CK technique.
type AttackCommand struct {
	ID        int     `json:"id"`
	AttackID  int     `json:"attack_id"`
	CommandID int     `json:"command_id"`
	Attack    *Attack `json:"attack,omitempty"`
}

// String returns a string representation of an AttackCommand.
func (ac *AttackCommand) String() string {
	if ac.Attack != nil {
		return fmt.Sprintf("%s on Command %d", ac.Attack.String(), ac.CommandID)
	}
	return fmt.Sprintf("AttackCommand %d (Attack %d, Command %d)", ac.ID, ac.AttackID, ac.CommandID)
}
