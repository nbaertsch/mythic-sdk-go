package unit

import (
	"testing"
	"time"

	"github.com/nbaertsch/mythic-sdk-go/pkg/mythic/types"
)

// TestOperatorString tests the Operator.String() method
func TestOperatorString(t *testing.T) {
	tests := []struct {
		name     string
		operator types.Operator
		contains []string
	}{
		{
			name: "admin operator",
			operator: types.Operator{
				ID:       1,
				Username: "admin",
				Admin:    true,
				Active:   true,
			},
			contains: []string{"admin", "Admin"},
		},
		{
			name: "regular operator",
			operator: types.Operator{
				ID:       2,
				Username: "operator1",
				Admin:    false,
				Active:   true,
			},
			contains: []string{"operator1", "Operator"},
		},
		{
			name: "deleted operator",
			operator: types.Operator{
				ID:       3,
				Username: "olduser",
				Admin:    false,
				Active:   false,
				Deleted:  true,
			},
			contains: []string{"olduser", "deleted"},
		},
		{
			name: "inactive operator",
			operator: types.Operator{
				ID:       4,
				Username: "suspended",
				Admin:    false,
				Active:   false,
				Deleted:  false,
			},
			contains: []string{"suspended", "inactive"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.operator.String()
			if result == "" {
				t.Error("String() should not return empty string")
			}
			for _, substr := range tt.contains {
				if !contains(result, substr) {
					t.Errorf("String() = %q, should contain %q", result, substr)
				}
			}
		})
	}
}

// TestOperatorIsAdmin tests the Operator.IsAdmin() method
func TestOperatorIsAdmin(t *testing.T) {
	tests := []struct {
		name     string
		operator types.Operator
		want     bool
	}{
		{
			name: "admin operator",
			operator: types.Operator{
				ID:       1,
				Username: "admin",
				Admin:    true,
			},
			want: true,
		},
		{
			name: "regular operator",
			operator: types.Operator{
				ID:       2,
				Username: "user",
				Admin:    false,
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.operator.IsAdmin(); got != tt.want {
				t.Errorf("IsAdmin() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestOperatorIsActive tests the Operator.IsActive() method
func TestOperatorIsActive(t *testing.T) {
	tests := []struct {
		name     string
		operator types.Operator
		want     bool
	}{
		{
			name: "active operator",
			operator: types.Operator{
				ID:       1,
				Username: "user1",
				Active:   true,
				Deleted:  false,
			},
			want: true,
		},
		{
			name: "inactive operator",
			operator: types.Operator{
				ID:       2,
				Username: "user2",
				Active:   false,
				Deleted:  false,
			},
			want: false,
		},
		{
			name: "deleted operator",
			operator: types.Operator{
				ID:       3,
				Username: "user3",
				Active:   true,
				Deleted:  true,
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.operator.IsActive(); got != tt.want {
				t.Errorf("IsActive() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestOperatorIsDeleted tests the Operator.IsDeleted() method
func TestOperatorIsDeleted(t *testing.T) {
	tests := []struct {
		name     string
		operator types.Operator
		want     bool
	}{
		{
			name: "deleted operator",
			operator: types.Operator{
				ID:       1,
				Username: "deleted",
				Deleted:  true,
			},
			want: true,
		},
		{
			name: "active operator",
			operator: types.Operator{
				ID:       2,
				Username: "active",
				Deleted:  false,
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.operator.IsDeleted(); got != tt.want {
				t.Errorf("IsDeleted() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestOperatorIsLocked tests the Operator.IsLocked() method
func TestOperatorIsLocked(t *testing.T) {
	tests := []struct {
		name     string
		operator types.Operator
		want     bool
	}{
		{
			name: "locked operator (10 failures)",
			operator: types.Operator{
				ID:               1,
				Username:         "locked",
				FailedLoginCount: 10,
			},
			want: true,
		},
		{
			name: "locked operator (more than 10)",
			operator: types.Operator{
				ID:               2,
				Username:         "locked2",
				FailedLoginCount: 15,
			},
			want: true,
		},
		{
			name: "not locked operator",
			operator: types.Operator{
				ID:               3,
				Username:         "unlocked",
				FailedLoginCount: 5,
			},
			want: false,
		},
		{
			name: "no failures",
			operator: types.Operator{
				ID:               4,
				Username:         "clean",
				FailedLoginCount: 0,
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.operator.IsLocked(); got != tt.want {
				t.Errorf("IsLocked() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestOperatorIsBotAccount tests the Operator.IsBotAccount() method
func TestOperatorIsBotAccount(t *testing.T) {
	tests := []struct {
		name     string
		operator types.Operator
		want     bool
	}{
		{
			name: "bot account",
			operator: types.Operator{
				ID:          1,
				Username:    "bot_account",
				AccountType: types.AccountTypeBot,
			},
			want: true,
		},
		{
			name: "user account",
			operator: types.Operator{
				ID:          2,
				Username:    "human_user",
				AccountType: types.AccountTypeUser,
			},
			want: false,
		},
		{
			name: "empty account type",
			operator: types.Operator{
				ID:          3,
				Username:    "unknown",
				AccountType: "",
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.operator.IsBotAccount(); got != tt.want {
				t.Errorf("IsBotAccount() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestInviteLinkString tests the InviteLink.String() method
func TestInviteLinkString(t *testing.T) {
	tests := []struct {
		name     string
		link     types.InviteLink
		contains []string
	}{
		{
			name: "basic invite link",
			link: types.InviteLink{
				ID:          1,
				Code:        "ABC123",
				MaxUses:     10,
				CurrentUses: 3,
			},
			contains: []string{"ABC123", "3/10"},
		},
		{
			name: "fully used link",
			link: types.InviteLink{
				ID:          2,
				Code:        "XYZ789",
				MaxUses:     5,
				CurrentUses: 5,
			},
			contains: []string{"XYZ789", "5/5"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.link.String()
			if result == "" {
				t.Error("String() should not return empty string")
			}
			for _, substr := range tt.contains {
				if !contains(result, substr) {
					t.Errorf("String() = %q, should contain %q", result, substr)
				}
			}
		})
	}
}

// TestInviteLinkIsExpired tests the InviteLink.IsExpired() method
func TestInviteLinkIsExpired(t *testing.T) {
	now := time.Now()
	tests := []struct {
		name string
		link types.InviteLink
		want bool
	}{
		{
			name: "expired link",
			link: types.InviteLink{
				ID:        1,
				Code:      "EXPIRED",
				ExpiresAt: now.Add(-24 * time.Hour),
			},
			want: true,
		},
		{
			name: "valid link",
			link: types.InviteLink{
				ID:        2,
				Code:      "VALID",
				ExpiresAt: now.Add(24 * time.Hour),
			},
			want: false,
		},
		{
			name: "link expiring soon",
			link: types.InviteLink{
				ID:        3,
				Code:      "SOON",
				ExpiresAt: now.Add(1 * time.Hour),
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.link.IsExpired(); got != tt.want {
				t.Errorf("IsExpired() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestInviteLinkIsActive tests the InviteLink.IsActive() method
func TestInviteLinkIsActive(t *testing.T) {
	now := time.Now()
	tests := []struct {
		name string
		link types.InviteLink
		want bool
	}{
		{
			name: "active and not expired",
			link: types.InviteLink{
				ID:        1,
				Code:      "ACTIVE",
				Active:    true,
				ExpiresAt: now.Add(24 * time.Hour),
			},
			want: true,
		},
		{
			name: "active but expired",
			link: types.InviteLink{
				ID:        2,
				Code:      "EXPIRED",
				Active:    true,
				ExpiresAt: now.Add(-24 * time.Hour),
			},
			want: false,
		},
		{
			name: "inactive",
			link: types.InviteLink{
				ID:        3,
				Code:      "INACTIVE",
				Active:    false,
				ExpiresAt: now.Add(24 * time.Hour),
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.link.IsActive(); got != tt.want {
				t.Errorf("IsActive() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestInviteLinkHasUsesRemaining tests the InviteLink.HasUsesRemaining() method
func TestInviteLinkHasUsesRemaining(t *testing.T) {
	tests := []struct {
		name string
		link types.InviteLink
		want bool
	}{
		{
			name: "has uses remaining",
			link: types.InviteLink{
				ID:          1,
				Code:        "AVAILABLE",
				MaxUses:     10,
				CurrentUses: 5,
			},
			want: true,
		},
		{
			name: "fully used",
			link: types.InviteLink{
				ID:          2,
				Code:        "FULL",
				MaxUses:     10,
				CurrentUses: 10,
			},
			want: false,
		},
		{
			name: "exceeded max uses",
			link: types.InviteLink{
				ID:          3,
				Code:        "EXCEEDED",
				MaxUses:     10,
				CurrentUses: 15,
			},
			want: false,
		},
		{
			name: "no uses yet",
			link: types.InviteLink{
				ID:          4,
				Code:        "NEW",
				MaxUses:     10,
				CurrentUses: 0,
			},
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.link.HasUsesRemaining(); got != tt.want {
				t.Errorf("HasUsesRemaining() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestOperatorFields tests the Operator type structure
func TestOperatorFields(t *testing.T) {
	currentOpID := 5

	operator := types.Operator{
		ID:                 1,
		Username:           "testuser",
		Admin:              true,
		Active:             true,
		Deleted:            false,
		CurrentOperationID: &currentOpID,
		AccountType:        types.AccountTypeUser,
		FailedLoginCount:   0,
	}

	if operator.ID != 1 {
		t.Errorf("Expected ID 1, got %d", operator.ID)
	}
	if operator.Username != "testuser" {
		t.Errorf("Expected Username 'testuser', got %q", operator.Username)
	}
	if !operator.Admin {
		t.Error("Expected Admin to be true")
	}
	if !operator.Active {
		t.Error("Expected Active to be true")
	}
	if operator.Deleted {
		t.Error("Expected Deleted to be false")
	}
	if operator.CurrentOperationID == nil || *operator.CurrentOperationID != 5 {
		t.Error("Expected CurrentOperationID to be 5")
	}
	if operator.AccountType != types.AccountTypeUser {
		t.Errorf("Expected AccountType 'user', got %q", operator.AccountType)
	}
}

// TestCreateOperatorRequestTypes tests the CreateOperatorRequest type
func TestCreateOperatorRequestTypes(t *testing.T) {
	req := types.CreateOperatorRequest{
		Username: "newuser",
		Password: "securepassword123",
	}

	if req.Username != "newuser" {
		t.Errorf("Expected Username 'newuser', got %q", req.Username)
	}
	if req.Password != "securepassword123" {
		t.Errorf("Expected Password 'securepassword123', got %q", req.Password)
	}
}

// TestUpdateOperatorStatusRequestTypes tests the UpdateOperatorStatusRequest type
func TestUpdateOperatorStatusRequestTypes(t *testing.T) {
	active := true
	admin := false
	deleted := true

	req := types.UpdateOperatorStatusRequest{
		OperatorID: 5,
		Active:     &active,
		Admin:      &admin,
		Deleted:    &deleted,
	}

	if req.OperatorID != 5 {
		t.Errorf("Expected OperatorID 5, got %d", req.OperatorID)
	}
	if req.Active == nil || *req.Active != true {
		t.Error("Expected Active to be true")
	}
	if req.Admin == nil || *req.Admin != false {
		t.Error("Expected Admin to be false")
	}
	if req.Deleted == nil || *req.Deleted != true {
		t.Error("Expected Deleted to be true")
	}
}

// TestUpdatePasswordAndEmailRequestTypes tests the UpdatePasswordAndEmailRequest type
func TestUpdatePasswordAndEmailRequestTypes(t *testing.T) {
	newPassword := "newpassword123"
	email := "user@example.com"

	req := types.UpdatePasswordAndEmailRequest{
		OperatorID:  5,
		OldPassword: "oldpassword",
		NewPassword: &newPassword,
		Email:       &email,
	}

	if req.OperatorID != 5 {
		t.Errorf("Expected OperatorID 5, got %d", req.OperatorID)
	}
	if req.OldPassword != "oldpassword" {
		t.Errorf("Expected OldPassword 'oldpassword', got %q", req.OldPassword)
	}
	if req.NewPassword == nil || *req.NewPassword != "newpassword123" {
		t.Error("Expected NewPassword to be 'newpassword123'")
	}
	if req.Email == nil || *req.Email != "user@example.com" {
		t.Error("Expected Email to be 'user@example.com'")
	}
}

// TestOperatorPreferencesTypes tests the OperatorPreferences type
func TestOperatorPreferencesTypes(t *testing.T) {
	prefs := types.OperatorPreferences{
		OperatorID:      5,
		PreferencesJSON: `{"theme":"dark"}`,
		Preferences: map[string]interface{}{
			"theme": "dark",
		},
		InteractType: "browser",
		ConsoleSize:  80,
		FontSize:     14,
	}

	if prefs.OperatorID != 5 {
		t.Errorf("Expected OperatorID 5, got %d", prefs.OperatorID)
	}
	if prefs.PreferencesJSON == "" {
		t.Error("Expected PreferencesJSON to not be empty")
	}
	if len(prefs.Preferences) == 0 {
		t.Error("Expected Preferences to not be empty")
	}
}

// TestOperatorSecretsTypes tests the OperatorSecrets type
func TestOperatorSecretsTypes(t *testing.T) {
	secrets := types.OperatorSecrets{
		OperatorID:  5,
		SecretsJSON: `{"api_key":"secret123"}`,
		Secrets: map[string]interface{}{
			"api_key": "secret123",
		},
	}

	if secrets.OperatorID != 5 {
		t.Errorf("Expected OperatorID 5, got %d", secrets.OperatorID)
	}
	if secrets.SecretsJSON == "" {
		t.Error("Expected SecretsJSON to not be empty")
	}
	if len(secrets.Secrets) == 0 {
		t.Error("Expected Secrets to not be empty")
	}
}

// TestInviteLinkTypes tests the InviteLink type
func TestInviteLinkTypes(t *testing.T) {
	now := time.Now()
	link := types.InviteLink{
		ID:          1,
		Code:        "INVITE123",
		ExpiresAt:   now.Add(24 * time.Hour),
		CreatedBy:   5,
		CreatedAt:   now,
		MaxUses:     10,
		CurrentUses: 3,
		Active:      true,
	}

	if link.ID != 1 {
		t.Errorf("Expected ID 1, got %d", link.ID)
	}
	if link.Code != "INVITE123" {
		t.Errorf("Expected Code 'INVITE123', got %q", link.Code)
	}
	if link.MaxUses != 10 {
		t.Errorf("Expected MaxUses 10, got %d", link.MaxUses)
	}
	if link.CurrentUses != 3 {
		t.Errorf("Expected CurrentUses 3, got %d", link.CurrentUses)
	}
	if !link.Active {
		t.Error("Expected Active to be true")
	}
}

// TestAccountTypeConstants tests account type constants
func TestAccountTypeConstants(t *testing.T) {
	if types.AccountTypeUser != "user" {
		t.Errorf("Expected AccountTypeUser 'user', got %q", types.AccountTypeUser)
	}
	if types.AccountTypeBot != "bot" {
		t.Errorf("Expected AccountTypeBot 'bot', got %q", types.AccountTypeBot)
	}
}

// TestViewModeConstants tests view mode constants
func TestViewModeConstants(t *testing.T) {
	if types.ViewModeOperator != "operator" {
		t.Errorf("Expected ViewModeOperator 'operator', got %q", types.ViewModeOperator)
	}
	if types.ViewModeSpectator != "spectator" {
		t.Errorf("Expected ViewModeSpectator 'spectator', got %q", types.ViewModeSpectator)
	}
}
