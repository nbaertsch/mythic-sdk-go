package unit

import (
	"testing"
	"time"

	"github.com/nbaertsch/mythic-sdk-go/pkg/mythic"
	"github.com/nbaertsch/mythic-sdk-go/pkg/mythic/types"
)

func TestCallbackString(t *testing.T) {
	cb := &types.Callback{
		User:   "admin",
		Host:   "DC01",
		OS:     "Windows",
		Active: true,
	}

	expected := "admin@DC01 (Windows, active)"
	if cb.String() != expected {
		t.Errorf("Callback.String() = %q, want %q", cb.String(), expected)
	}

	cb.Active = false
	expected = "admin@DC01 (Windows, inactive)"
	if cb.String() != expected {
		t.Errorf("Callback.String() = %q, want %q", cb.String(), expected)
	}
}

func TestCallbackIsHigh(t *testing.T) {
	tests := []struct {
		name           string
		integrityLevel types.CallbackIntegrityLevel
		want           bool
	}{
		{
			name:           "low integrity",
			integrityLevel: types.IntegrityLevelLow,
			want:           false,
		},
		{
			name:           "medium integrity",
			integrityLevel: types.IntegrityLevelMedium,
			want:           false,
		},
		{
			name:           "high integrity",
			integrityLevel: types.IntegrityLevelHigh,
			want:           true,
		},
		{
			name:           "system integrity",
			integrityLevel: types.IntegrityLevelSystem,
			want:           true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cb := &types.Callback{
				IntegrityLevel: tt.integrityLevel,
			}
			if got := cb.IsHigh(); got != tt.want {
				t.Errorf("Callback.IsHigh() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCallbackIsSystem(t *testing.T) {
	tests := []struct {
		name           string
		integrityLevel types.CallbackIntegrityLevel
		want           bool
	}{
		{
			name:           "low integrity",
			integrityLevel: types.IntegrityLevelLow,
			want:           false,
		},
		{
			name:           "medium integrity",
			integrityLevel: types.IntegrityLevelMedium,
			want:           false,
		},
		{
			name:           "high integrity",
			integrityLevel: types.IntegrityLevelHigh,
			want:           false,
		},
		{
			name:           "system integrity",
			integrityLevel: types.IntegrityLevelSystem,
			want:           true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cb := &types.Callback{
				IntegrityLevel: tt.integrityLevel,
			}
			if got := cb.IsSystem(); got != tt.want {
				t.Errorf("Callback.IsSystem() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCallbackUpdateRequest(t *testing.T) {
	// Test creating update request with various fields
	active := true
	locked := false
	description := "High-value target"
	user := "SYSTEM"
	host := "DC01"
	os := "Windows Server 2019"
	arch := "x64"
	extraInfo := "Domain Controller"
	sleepInfo := "60s/10%"
	pid := 1234
	processName := "lsass.exe"
	integrityLevel := types.IntegrityLevelSystem
	domain := "CORP.LOCAL"

	req := &types.CallbackUpdateRequest{
		CallbackDisplayID: 1,
		Active:            &active,
		Locked:            &locked,
		Description:       &description,
		IPs:               []string{"10.0.0.1", "192.168.1.100"},
		User:              &user,
		Host:              &host,
		OS:                &os,
		Architecture:      &arch,
		ExtraInfo:         &extraInfo,
		SleepInfo:         &sleepInfo,
		PID:               &pid,
		ProcessName:       &processName,
		IntegrityLevel:    &integrityLevel,
		Domain:            &domain,
	}

	// Verify all fields are set
	if req.CallbackDisplayID != 1 {
		t.Error("CallbackDisplayID not set correctly")
	}
	if req.Active == nil || *req.Active != true {
		t.Error("Active not set correctly")
	}
	if req.Locked == nil || *req.Locked != false {
		t.Error("Locked not set correctly")
	}
	if req.Description == nil || *req.Description != description {
		t.Error("Description not set correctly")
	}
	if len(req.IPs) != 2 {
		t.Error("IPs not set correctly")
	}
	if req.User == nil || *req.User != user {
		t.Error("User not set correctly")
	}
	if req.Host == nil || *req.Host != host {
		t.Error("Host not set correctly")
	}
	if req.IntegrityLevel == nil || *req.IntegrityLevel != integrityLevel {
		t.Error("IntegrityLevel not set correctly")
	}
}

func TestCallbackTypes(t *testing.T) {
	// Test that all callback-related types can be created
	now := time.Now()

	callback := &types.Callback{
		ID:              1,
		DisplayID:       1,
		AgentCallbackID: "abc-123",
		InitCallback:    now,
		LastCheckin:     now,
		User:            "admin",
		Host:            "WORKSTATION01",
		PID:             1234,
		IP:              []string{"10.0.0.1", "192.168.1.100"},
		ExternalIP:      "1.2.3.4",
		ProcessName:     "explorer.exe",
		Description:     "Test callback",
		Active:          true,
		IntegrityLevel:  types.IntegrityLevelHigh,
		Locked:          false,
		OS:              "Windows 10",
		Architecture:    "x64",
		Domain:          "WORKGROUP",
		ExtraInfo:       "Additional info",
		SleepInfo:       "30s/5%",
		OperationID:     1,
		PayloadTypeID:   1,
		OperatorID:      1,
		Payload: &types.CallbackPayload{
			ID:          1,
			UUID:        "payload-uuid",
			Description: "Test payload",
			OS:          "Windows",
		},
		PayloadType: &types.CallbackPayloadType{
			ID:   1,
			Name: "apollo",
		},
		Operator: &types.CallbackOperator{
			ID:       1,
			Username: "operator1",
		},
	}

	// Verify the callback was created successfully
	if callback.ID != 1 {
		t.Error("Callback ID not set")
	}
	if callback.User != "admin" {
		t.Error("Callback User not set")
	}
	if callback.Payload == nil {
		t.Error("Callback Payload not set")
	}
	if callback.PayloadType == nil {
		t.Error("Callback PayloadType not set")
	}
	if callback.Operator == nil {
		t.Error("Callback Operator not set")
	}
}

func TestCallbackIntegrityLevels(t *testing.T) {
	// Test that integrity level constants are defined correctly
	levels := []types.CallbackIntegrityLevel{
		types.IntegrityLevelLow,
		types.IntegrityLevelMedium,
		types.IntegrityLevelHigh,
		types.IntegrityLevelSystem,
	}

	// Verify they are in ascending order
	for i := 1; i < len(levels); i++ {
		if levels[i] <= levels[i-1] {
			t.Errorf("Integrity levels not in ascending order: %d <= %d", levels[i], levels[i-1])
		}
	}

	// Verify specific values
	if types.IntegrityLevelLow != 2 {
		t.Errorf("IntegrityLevelLow = %d, want 2", types.IntegrityLevelLow)
	}
	if types.IntegrityLevelMedium != 3 {
		t.Errorf("IntegrityLevelMedium = %d, want 3", types.IntegrityLevelMedium)
	}
	if types.IntegrityLevelHigh != 4 {
		t.Errorf("IntegrityLevelHigh = %d, want 4", types.IntegrityLevelHigh)
	}
	if types.IntegrityLevelSystem != 5 {
		t.Errorf("IntegrityLevelSystem = %d, want 5", types.IntegrityLevelSystem)
	}
}

func TestCreateCallbackInput_Structure(t *testing.T) {
	// Test CreateCallbackInput with all optional fields
	ip := "192.168.1.100"
	externalIP := "1.2.3.4"
	user := "testuser"
	host := "TESTHOST"
	domain := "WORKGROUP"
	description := "Test callback"
	processName := "explorer.exe"
	sleepInfo := "60s/10%"
	extraInfo := "Additional info"

	input := &mythic.CreateCallbackInput{
		PayloadUUID: "test-uuid-123",
		IP:          &ip,
		ExternalIP:  &externalIP,
		User:        &user,
		Host:        &host,
		Domain:      &domain,
		Description: &description,
		ProcessName: &processName,
		SleepInfo:   &sleepInfo,
		ExtraInfo:   &extraInfo,
	}

	// Verify all fields are set correctly
	if input.PayloadUUID != "test-uuid-123" {
		t.Errorf("Expected PayloadUUID 'test-uuid-123', got %q", input.PayloadUUID)
	}
	if input.IP == nil || *input.IP != ip {
		t.Errorf("Expected IP %q, got %v", ip, input.IP)
	}
	if input.ExternalIP == nil || *input.ExternalIP != externalIP {
		t.Errorf("Expected ExternalIP %q, got %v", externalIP, input.ExternalIP)
	}
	if input.User == nil || *input.User != user {
		t.Errorf("Expected User %q, got %v", user, input.User)
	}
	if input.Host == nil || *input.Host != host {
		t.Errorf("Expected Host %q, got %v", host, input.Host)
	}
	if input.Domain == nil || *input.Domain != domain {
		t.Errorf("Expected Domain %q, got %v", domain, input.Domain)
	}
	if input.Description == nil || *input.Description != description {
		t.Errorf("Expected Description %q, got %v", description, input.Description)
	}
	if input.ProcessName == nil || *input.ProcessName != processName {
		t.Errorf("Expected ProcessName %q, got %v", processName, input.ProcessName)
	}
	if input.SleepInfo == nil || *input.SleepInfo != sleepInfo {
		t.Errorf("Expected SleepInfo %q, got %v", sleepInfo, input.SleepInfo)
	}
	if input.ExtraInfo == nil || *input.ExtraInfo != extraInfo {
		t.Errorf("Expected ExtraInfo %q, got %v", extraInfo, input.ExtraInfo)
	}
}

func TestCreateCallbackInput_MinimalRequired(t *testing.T) {
	// Test CreateCallbackInput with only required field
	input := &mythic.CreateCallbackInput{
		PayloadUUID: "minimal-uuid",
	}

	if input.PayloadUUID != "minimal-uuid" {
		t.Errorf("Expected PayloadUUID 'minimal-uuid', got %q", input.PayloadUUID)
	}

	// Verify optional fields are nil
	if input.IP != nil {
		t.Error("Expected IP to be nil")
	}
	if input.ExternalIP != nil {
		t.Error("Expected ExternalIP to be nil")
	}
	if input.User != nil {
		t.Error("Expected User to be nil")
	}
	if input.Host != nil {
		t.Error("Expected Host to be nil")
	}
}
