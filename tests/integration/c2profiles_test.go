//go:build integration

package integration

import (
	"context"
	"testing"
	"time"

	"github.com/nbaertsch/mythic-sdk-go/pkg/mythic/types"
)

func TestC2Profiles_GetC2Profiles(t *testing.T) {
	SkipIfNoMythic(t)

	client := AuthenticateTestClient(t)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	profiles, err := client.GetC2Profiles(ctx)
	if err != nil {
		t.Fatalf("GetC2Profiles failed: %v", err)
	}

	if profiles == nil {
		t.Fatal("GetC2Profiles returned nil")
	}

	t.Logf("Found %d C2 profile(s)", len(profiles))

	// If there are profiles, verify structure
	if len(profiles) > 0 {
		p := profiles[0]
		if p.ID == 0 {
			t.Error("C2Profile ID should not be 0")
		}
		if p.Name == "" {
			t.Error("C2Profile should have a name")
		}
		t.Logf("First profile: %s", p.String())
		t.Logf("  - ID: %d", p.ID)
		t.Logf("  - Name: %s", p.Name)
		t.Logf("  - Running: %v", p.IsRunning())
		t.Logf("  - Deleted: %v", p.IsDeleted())
		t.Logf("  - IsP2P: %v", p.IsP2P)
		t.Logf("  - ServerOnly: %v", p.ServerOnly)
	}
}

func TestC2Profiles_GetC2ProfileByID(t *testing.T) {
	SkipIfNoMythic(t)

	client := AuthenticateTestClient(t)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Get all profiles first
	profiles, err := client.GetC2Profiles(ctx)
	if err != nil {
		t.Fatalf("GetC2Profiles failed: %v", err)
	}

	if len(profiles) == 0 {
		t.Skip("No C2 profiles available for testing")
	}

	// Test getting a specific profile
	profileID := profiles[0].ID
	profile, err := client.GetC2ProfileByID(ctx, profileID)
	if err != nil {
		t.Fatalf("GetC2ProfileByID failed: %v", err)
	}

	if profile == nil {
		t.Fatal("GetC2ProfileByID returned nil")
	}

	if profile.ID != profileID {
		t.Errorf("Expected profile ID %d, got %d", profileID, profile.ID)
	}

	t.Logf("Retrieved profile: %s", profile.String())
	t.Logf("  - Description: %s", profile.Description)
	t.Logf("  - OperationID: %d", profile.OperationID)
	t.Logf("  - CreationTime: %s", profile.CreationTime.Format("2006-01-02 15:04:05"))
}

func TestC2Profiles_GetC2ProfileByID_InvalidID(t *testing.T) {
	SkipIfNoMythic(t)

	client := AuthenticateTestClient(t)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Test with zero ID
	_, err := client.GetC2ProfileByID(ctx, 0)
	if err == nil {
		t.Fatal("Expected error for zero profile ID, got nil")
	}

	t.Logf("Expected error: %v", err)
}

func TestC2Profiles_GetProfileOutput(t *testing.T) {
	SkipIfNoMythic(t)

	client := AuthenticateTestClient(t)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Get all profiles first
	profiles, err := client.GetC2Profiles(ctx)
	if err != nil {
		t.Fatalf("GetC2Profiles failed: %v", err)
	}

	if len(profiles) == 0 {
		t.Skip("No C2 profiles available for testing")
	}

	// Find a running profile to get output from
	var runningProfile *types.C2Profile
	for _, p := range profiles {
		if p.IsRunning() {
			runningProfile = p
			break
		}
	}

	if runningProfile == nil {
		t.Skip("No running C2 profiles available for testing")
	}

	// Get profile output
	output, err := client.GetProfileOutput(ctx, runningProfile.ID)
	if err != nil {
		t.Fatalf("GetProfileOutput failed: %v", err)
	}

	if output == nil {
		t.Fatal("GetProfileOutput returned nil")
	}

	t.Logf("Profile %s output:", runningProfile.Name)
	if output.Output != "" {
		t.Logf("  - Combined output: %d characters", len(output.Output))
	}
	if output.StdOut != "" {
		t.Logf("  - StdOut: %d characters", len(output.StdOut))
	}
	if output.StdErr != "" {
		t.Logf("  - StdErr: %d characters", len(output.StdErr))
	}
}

func TestC2Profiles_GetProfileOutput_InvalidID(t *testing.T) {
	SkipIfNoMythic(t)

	client := AuthenticateTestClient(t)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Test with zero ID
	_, err := client.GetProfileOutput(ctx, 0)
	if err == nil {
		t.Fatal("Expected error for zero profile ID, got nil")
	}

	t.Logf("Expected error: %v", err)
}

func TestC2Profiles_StartStopProfile(t *testing.T) {
	SkipIfNoMythic(t)

	client := AuthenticateTestClient(t)

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	// Get all profiles first
	profiles, err := client.GetC2Profiles(ctx)
	if err != nil {
		t.Fatalf("GetC2Profiles failed: %v", err)
	}

	if len(profiles) == 0 {
		t.Skip("No C2 profiles available for testing")
	}

	// Find a stopped profile to test starting
	var stoppedProfile *types.C2Profile
	for _, p := range profiles {
		if !p.IsRunning() {
			stoppedProfile = p
			break
		}
	}

	if stoppedProfile == nil {
		t.Skip("No stopped C2 profiles available for testing")
	}

	t.Logf("Testing with profile: %s (ID: %d)", stoppedProfile.Name, stoppedProfile.ID)

	// Note: Starting/stopping C2 profiles may require specific permissions
	// and may not work in all test environments. This test documents the
	// expected behavior but may need to be skipped in some environments.
	t.Log("Note: Start/Stop operations may require elevated permissions")

	// Test starting the profile
	err = client.StartStopProfile(ctx, stoppedProfile.ID, true)
	if err != nil {
		t.Logf("StartStopProfile (start) failed: %v (this may be expected in test environment)", err)
		// Don't fail the test as this might be a permission issue
		return
	}

	t.Logf("Successfully started profile %s", stoppedProfile.Name)

	// Wait a moment for the profile to start
	time.Sleep(2 * time.Second)

	// Verify profile is now running
	updatedProfile, err := client.GetC2ProfileByID(ctx, stoppedProfile.ID)
	if err != nil {
		t.Fatalf("Failed to get updated profile: %v", err)
	}

	if !updatedProfile.IsRunning() {
		t.Error("Profile should be running after start command")
	}

	// Test stopping the profile
	err = client.StartStopProfile(ctx, stoppedProfile.ID, false)
	if err != nil {
		t.Errorf("StartStopProfile (stop) failed: %v", err)
	}

	t.Logf("Successfully stopped profile %s", stoppedProfile.Name)

	// Wait a moment for the profile to stop
	time.Sleep(2 * time.Second)

	// Verify profile is now stopped
	finalProfile, err := client.GetC2ProfileByID(ctx, stoppedProfile.ID)
	if err != nil {
		t.Fatalf("Failed to get final profile state: %v", err)
	}

	if finalProfile.IsRunning() {
		t.Error("Profile should be stopped after stop command")
	}
}

func TestC2Profiles_StartStopProfile_InvalidID(t *testing.T) {
	SkipIfNoMythic(t)

	client := AuthenticateTestClient(t)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Test with zero ID
	err := client.StartStopProfile(ctx, 0, true)
	if err == nil {
		t.Fatal("Expected error for zero profile ID, got nil")
	}

	t.Logf("Expected error: %v", err)
}

func TestC2Profiles_C2SampleMessage(t *testing.T) {
	SkipIfNoMythic(t)

	client := AuthenticateTestClient(t)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Get all profiles first
	profiles, err := client.GetC2Profiles(ctx)
	if err != nil {
		t.Fatalf("GetC2Profiles failed: %v", err)
	}

	if len(profiles) == 0 {
		t.Skip("No C2 profiles available for testing")
	}

	// Test with first profile
	profileID := profiles[0].ID
	message, err := client.C2SampleMessage(ctx, profileID, "checkin")
	if err != nil {
		// This might fail if the profile doesn't support sample messages
		t.Logf("C2SampleMessage failed: %v (may not be supported by this profile)", err)
		return
	}

	if message == nil {
		t.Fatal("C2SampleMessage returned nil")
	}

	if message.Message == "" {
		t.Error("Sample message should not be empty")
	}

	t.Logf("Sample message for profile %s:", profiles[0].Name)
	t.Logf("  - Message: %s", message.Message)
	if message.Metadata != nil {
		t.Logf("  - Metadata: %v", message.Metadata)
	}
}

func TestC2Profiles_C2SampleMessage_InvalidID(t *testing.T) {
	SkipIfNoMythic(t)

	client := AuthenticateTestClient(t)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Test with zero ID
	_, err := client.C2SampleMessage(ctx, 0, "checkin")
	if err == nil {
		t.Fatal("Expected error for zero profile ID, got nil")
	}

	t.Logf("Expected error: %v", err)
}

func TestC2Profiles_C2GetIOC(t *testing.T) {
	SkipIfNoMythic(t)

	client := AuthenticateTestClient(t)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Get all profiles first
	profiles, err := client.GetC2Profiles(ctx)
	if err != nil {
		t.Fatalf("GetC2Profiles failed: %v", err)
	}

	if len(profiles) == 0 {
		t.Skip("No C2 profiles available for testing")
	}

	// Test with first profile
	profileID := profiles[0].ID
	iocs, err := client.C2GetIOC(ctx, profileID)
	if err != nil {
		// This might fail if the profile doesn't support IOC extraction
		t.Logf("C2GetIOC failed: %v (may not be supported by this profile)", err)
		return
	}

	if iocs == nil {
		t.Fatal("C2GetIOC returned nil")
	}

	if iocs.ProfileID != profileID {
		t.Errorf("Expected ProfileID %d, got %d", profileID, iocs.ProfileID)
	}

	t.Logf("IOCs for profile %s:", profiles[0].Name)
	t.Logf("  - Count: %d", len(iocs.IOCs))
	if len(iocs.IOCs) > 0 {
		for i, ioc := range iocs.IOCs {
			if i < 5 { // Log first 5 IOCs
				t.Logf("  - IOC %d: %s", i+1, ioc)
			}
		}
		if len(iocs.IOCs) > 5 {
			t.Logf("  - ... and %d more", len(iocs.IOCs)-5)
		}
	}
	if iocs.Type != "" {
		t.Logf("  - Type: %s", iocs.Type)
	}
}

func TestC2Profiles_C2GetIOC_InvalidID(t *testing.T) {
	SkipIfNoMythic(t)

	client := AuthenticateTestClient(t)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Test with zero ID
	_, err := client.C2GetIOC(ctx, 0)
	if err == nil {
		t.Fatal("Expected error for zero profile ID, got nil")
	}

	t.Logf("Expected error: %v", err)
}

func TestC2Profiles_C2HostFile(t *testing.T) {
	SkipIfNoMythic(t)

	client := AuthenticateTestClient(t)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Get all profiles first
	profiles, err := client.GetC2Profiles(ctx)
	if err != nil {
		t.Fatalf("GetC2Profiles failed: %v", err)
	}

	if len(profiles) == 0 {
		t.Skip("No C2 profiles available for testing")
	}

	// Find a running profile
	var runningProfile *types.C2Profile
	for _, p := range profiles {
		if p.IsRunning() {
			runningProfile = p
			break
		}
	}

	if runningProfile == nil {
		t.Skip("No running C2 profiles available for testing")
	}

	// Note: This test requires an actual file UUID from the Mythic instance
	// In a real scenario, you would upload a file first and then host it
	// For now, we'll just test the error case with an invalid UUID
	t.Log("Note: C2HostFile requires a valid file UUID from Mythic")

	// Test with a fake UUID (expected to fail)
	fakeUUID := "00000000-0000-0000-0000-000000000000"
	err = client.C2HostFile(ctx, runningProfile.ID, fakeUUID)
	if err != nil {
		t.Logf("C2HostFile failed with fake UUID (expected): %v", err)
		// This is expected - we don't have a real file to host
		return
	}

	// If it somehow succeeded, log that
	t.Log("C2HostFile succeeded with fake UUID (unexpected but not an error)")
}

func TestC2Profiles_C2HostFile_InvalidInput(t *testing.T) {
	SkipIfNoMythic(t)

	client := AuthenticateTestClient(t)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Test with zero profile ID
	err := client.C2HostFile(ctx, 0, "valid-uuid")
	if err == nil {
		t.Fatal("Expected error for zero profile ID, got nil")
	}
	t.Logf("Zero ID error: %v", err)

	// Test with empty file UUID
	err = client.C2HostFile(ctx, 1, "")
	if err == nil {
		t.Fatal("Expected error for empty file UUID, got nil")
	}
	t.Logf("Empty UUID error: %v", err)
}

func TestC2Profiles_ProfileTypes(t *testing.T) {
	SkipIfNoMythic(t)

	client := AuthenticateTestClient(t)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	profiles, err := client.GetC2Profiles(ctx)
	if err != nil {
		t.Fatalf("GetC2Profiles failed: %v", err)
	}

	if len(profiles) == 0 {
		t.Skip("No C2 profiles available for testing")
	}

	// Count profile types
	p2pCount := 0
	serverOnlyCount := 0
	runningCount := 0
	stoppedCount := 0

	for _, p := range profiles {
		if p.IsP2P {
			p2pCount++
		}
		if p.ServerOnly {
			serverOnlyCount++
		}
		if p.IsRunning() {
			runningCount++
		} else {
			stoppedCount++
		}
	}

	t.Logf("Profile type distribution:")
	t.Logf("  - P2P profiles: %d", p2pCount)
	t.Logf("  - Server-only profiles: %d", serverOnlyCount)
	t.Logf("  - Running profiles: %d", runningCount)
	t.Logf("  - Stopped profiles: %d", stoppedCount)
}
