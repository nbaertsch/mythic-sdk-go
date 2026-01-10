package unit

import (
	"testing"
	"time"

	"github.com/nbaertsch/mythic-sdk-go/pkg/mythic/types"
)

func TestOperationString(t *testing.T) {
	// Test active operation
	op := &types.Operation{
		Name:     "Test Operation",
		Complete: false,
	}

	expected := "Test Operation (active)"
	if op.String() != expected {
		t.Errorf("Operation.String() = %q, want %q", op.String(), expected)
	}

	// Test completed operation
	op.Complete = true
	expected = "Test Operation (complete)"
	if op.String() != expected {
		t.Errorf("Operation.String() = %q, want %q", op.String(), expected)
	}
}

func TestOperationIsComplete(t *testing.T) {
	tests := []struct {
		name     string
		complete bool
		want     bool
	}{
		{
			name:     "active operation",
			complete: false,
			want:     false,
		},
		{
			name:     "completed operation",
			complete: true,
			want:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			op := &types.Operation{
				Complete: tt.complete,
			}
			if got := op.IsComplete(); got != tt.want {
				t.Errorf("Operation.IsComplete() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOperationEventLogString(t *testing.T) {
	timestamp := time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC)
	log := &types.OperationEventLog{
		Level:     "info",
		Timestamp: timestamp,
		Message:   "Test event",
	}

	expected := "[info] 2024-01-15 10:30:00: Test event"
	if log.String() != expected {
		t.Errorf("OperationEventLog.String() = %q, want %q", log.String(), expected)
	}

	// Test warning level
	log.Level = "warning"
	expected = "[warning] 2024-01-15 10:30:00: Test event"
	if log.String() != expected {
		t.Errorf("OperationEventLog.String() = %q, want %q", log.String(), expected)
	}
}

func TestOperationTypes(t *testing.T) {
	// Test Operation structure
	op := &types.Operation{
		ID:          1,
		Name:        "Red Team 2024",
		Complete:    false,
		Webhook:     "https://hooks.example.com/webhook",
		Channel:     "#redteam",
		AdminID:     1,
		BannerText:  "Active Engagement",
		BannerColor: "#ff0000",
		Created:     time.Now(),
	}

	if op.ID != 1 {
		t.Errorf("Operation.ID = %d, want 1", op.ID)
	}
	if op.Name != "Red Team 2024" {
		t.Errorf("Operation.Name = %q, want %q", op.Name, "Red Team 2024")
	}
	if op.Complete {
		t.Error("Operation.Complete = true, want false")
	}
}

func TestOperationOperatorTypes(t *testing.T) {
	// Test OperationOperator structure
	opOp := &types.OperationOperator{
		ID:          1,
		OperationID: 1,
		OperatorID:  2,
	}

	if opOp.OperationID != 1 {
		t.Errorf("OperationOperator.OperationID = %d, want 1", opOp.OperationID)
	}
	if opOp.OperatorID != 2 {
		t.Errorf("OperationOperator.OperatorID = %d, want 2", opOp.OperatorID)
	}
}

func TestCreateOperationRequest(t *testing.T) {
	adminID := 1
	req := &types.CreateOperationRequest{
		Name:    "New Operation",
		AdminID: &adminID,
		Channel: "#ops",
		Webhook: "https://example.com/hook",
	}

	if req.Name != "New Operation" {
		t.Errorf("CreateOperationRequest.Name = %q, want %q", req.Name, "New Operation")
	}
	if req.AdminID == nil || *req.AdminID != 1 {
		t.Errorf("CreateOperationRequest.AdminID = %v, want 1", req.AdminID)
	}
}

func TestUpdateOperationRequest(t *testing.T) {
	name := "Updated Name"
	complete := true
	req := &types.UpdateOperationRequest{
		OperationID: 1,
		Name:        &name,
		Complete:    &complete,
	}

	if req.OperationID != 1 {
		t.Errorf("UpdateOperationRequest.OperationID = %d, want 1", req.OperationID)
	}
	if req.Name == nil || *req.Name != "Updated Name" {
		t.Errorf("UpdateOperationRequest.Name = %v, want %q", req.Name, "Updated Name")
	}
	if req.Complete == nil || !*req.Complete {
		t.Errorf("UpdateOperationRequest.Complete = %v, want true", req.Complete)
	}
}

func TestUpdateOperatorOperationRequest(t *testing.T) {
	req := &types.UpdateOperatorOperationRequest{
		OperatorID:  1,
		OperationID: 2,
		Remove:      false,
	}

	if req.OperatorID != 1 {
		t.Errorf("UpdateOperatorOperationRequest.OperatorID = %d, want 1", req.OperatorID)
	}
	if req.OperationID != 2 {
		t.Errorf("UpdateOperatorOperationRequest.OperationID = %d, want 2", req.OperationID)
	}
	if req.Remove {
		t.Error("UpdateOperatorOperationRequest.Remove = true, want false")
	}
}

func TestCreateOperationEventLogRequest(t *testing.T) {
	req := &types.CreateOperationEventLogRequest{
		OperationID: 1,
		Message:     "Test event",
		Level:       "info",
		Source:      "sdk",
	}

	if req.OperationID != 1 {
		t.Errorf("CreateOperationEventLogRequest.OperationID = %d, want 1", req.OperationID)
	}
	if req.Message != "Test event" {
		t.Errorf("CreateOperationEventLogRequest.Message = %q, want %q", req.Message, "Test event")
	}
	if req.Level != "info" {
		t.Errorf("CreateOperationEventLogRequest.Level = %q, want %q", req.Level, "info")
	}
	if req.Source != "sdk" {
		t.Errorf("CreateOperationEventLogRequest.Source = %q, want %q", req.Source, "sdk")
	}
}

func TestOperatorTypes(t *testing.T) {
	currentOpID := 1
	operator := &types.Operator{
		ID:                 1,
		Username:           "testuser",
		Admin:              true,
		Active:             true,
		CurrentOperationID: &currentOpID,
		ViewUTCTime:        false,
		Deleted:            false,
	}

	if operator.ID != 1 {
		t.Errorf("Operator.ID = %d, want 1", operator.ID)
	}
	if operator.Username != "testuser" {
		t.Errorf("Operator.Username = %q, want %q", operator.Username, "testuser")
	}
	if !operator.Admin {
		t.Error("Operator.Admin = false, want true")
	}
	if !operator.Active {
		t.Error("Operator.Active = false, want true")
	}
	if operator.CurrentOperationID == nil || *operator.CurrentOperationID != 1 {
		t.Errorf("Operator.CurrentOperationID = %v, want 1", operator.CurrentOperationID)
	}
}

func TestOperationEventLogTypes(t *testing.T) {
	timestamp := time.Now()
	log := &types.OperationEventLog{
		ID:          1,
		OperatorID:  2,
		OperationID: 3,
		Message:     "Task completed successfully",
		Timestamp:   timestamp,
		Level:       "info",
		Source:      "automation",
		Deleted:     false,
	}

	if log.ID != 1 {
		t.Errorf("OperationEventLog.ID = %d, want 1", log.ID)
	}
	if log.OperatorID != 2 {
		t.Errorf("OperationEventLog.OperatorID = %d, want 2", log.OperatorID)
	}
	if log.OperationID != 3 {
		t.Errorf("OperationEventLog.OperationID = %d, want 3", log.OperationID)
	}
	if log.Message != "Task completed successfully" {
		t.Errorf("OperationEventLog.Message = %q, want %q", log.Message, "Task completed successfully")
	}
	if log.Level != "info" {
		t.Errorf("OperationEventLog.Level = %q, want %q", log.Level, "info")
	}
	if log.Source != "automation" {
		t.Errorf("OperationEventLog.Source = %q, want %q", log.Source, "automation")
	}
	if log.Deleted {
		t.Error("OperationEventLog.Deleted = true, want false")
	}
}
