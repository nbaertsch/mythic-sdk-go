//go:build integration

package integration

import (
	"context"
	"testing"
	"time"

	"github.com/nbaertsch/mythic-sdk-go/pkg/mythic/types"
)

func TestArtifacts_GetArtifacts(t *testing.T) {

	client := AuthenticateTestClient(t)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	artifacts, err := client.GetArtifacts(ctx)
	if err != nil {
		t.Fatalf("GetArtifacts failed: %v", err)
	}

	if artifacts == nil {
		t.Fatal("GetArtifacts returned nil")
	}

	t.Logf("Found %d artifact(s)", len(artifacts))

	// If there are artifacts, verify structure
	if len(artifacts) > 0 {
		a := artifacts[0]
		if a.ID == 0 {
			t.Error("Artifact ID should not be 0")
		}
		if a.Artifact == "" {
			t.Error("Artifact should have a value")
		}
		t.Logf("First artifact: %s", a.String())
		t.Logf("  - ID: %d", a.ID)
		t.Logf("  - Artifact: %s", a.Artifact)
		t.Logf("  - Host: %s", a.Host)
		t.Logf("  - Has Task: %v", a.HasTask())
	}
}

func TestArtifacts_CreateAndRetrieve(t *testing.T) {

	client := AuthenticateTestClient(t)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Create a test artifact
	artifactValue := "C:\\Windows\\Temp\\test-artifact-" + time.Now().Format("20060102150405") + ".exe"
	baseArtifact := "C:\\Windows\\Temp"
	host := "TEST-WORKSTATION"

	req := &types.CreateArtifactRequest{
		Artifact:     artifactValue,
		BaseArtifact: &baseArtifact,
		Host:         &host,
	}

	artifact, err := client.CreateArtifact(ctx, req)
	if err != nil {
		t.Fatalf("CreateArtifact failed: %v", err)
	}

	if artifact == nil {
		t.Fatal("CreateArtifact returned nil")
	}

	t.Logf("Created artifact: %s", artifact.String())
	t.Logf("  - ID: %d", artifact.ID)
	t.Logf("  - Artifact: %s", artifact.Artifact)
	t.Logf("  - BaseArtifact: %s", artifact.BaseArtifact)
	t.Logf("  - Host: %s", artifact.Host)

	// Verify created artifact
	if artifact.Artifact != artifactValue {
		t.Errorf("Expected artifact %q, got %q", artifactValue, artifact.Artifact)
	}
	if artifact.Host != host {
		t.Errorf("Expected host %q, got %q", host, artifact.Host)
	}

	// Retrieve the artifact by ID
	retrieved, err := client.GetArtifactByID(ctx, artifact.ID)
	if err != nil {
		t.Fatalf("GetArtifactByID failed: %v", err)
	}

	if retrieved.ID != artifact.ID {
		t.Errorf("Expected artifact ID %d, got %d", artifact.ID, retrieved.ID)
	}
	if retrieved.Artifact != artifact.Artifact {
		t.Errorf("Expected artifact value %q, got %q", artifact.Artifact, retrieved.Artifact)
	}

	// Clean up: delete the artifact
	err = client.DeleteArtifact(ctx, artifact.ID)
	if err != nil {
		t.Logf("Warning: Failed to delete test artifact: %v", err)
	} else {
		t.Logf("Successfully deleted test artifact")
	}
}

func TestArtifacts_CreateArtifact_InvalidInput(t *testing.T) {

	client := AuthenticateTestClient(t)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Test with nil request
	_, err := client.CreateArtifact(ctx, nil)
	if err == nil {
		t.Fatal("Expected error for nil request, got nil")
	}
	t.Logf("Nil request error: %v", err)

	// Test with empty artifact value
	_, err = client.CreateArtifact(ctx, &types.CreateArtifactRequest{
		Artifact: "",
	})
	if err == nil {
		t.Fatal("Expected error for empty artifact value, got nil")
	}
	t.Logf("Empty artifact error: %v", err)
}

func TestArtifacts_GetArtifactByID_InvalidID(t *testing.T) {

	client := AuthenticateTestClient(t)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Test with zero ID
	_, err := client.GetArtifactByID(ctx, 0)
	if err == nil {
		t.Fatal("Expected error for zero artifact ID, got nil")
	}
	t.Logf("Zero ID error: %v", err)

	// Test with non-existent ID
	_, err = client.GetArtifactByID(ctx, 999999)
	if err == nil {
		t.Fatal("Expected error for non-existent artifact ID, got nil")
	}
	t.Logf("Non-existent ID error: %v", err)
}

func TestArtifacts_UpdateArtifact(t *testing.T) {

	client := AuthenticateTestClient(t)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Create a test artifact first
	artifactValue := "C:\\Windows\\Temp\\test-update-" + time.Now().Format("20060102150405") + ".exe"
	host := "OLD-HOST"

	artifact, err := client.CreateArtifact(ctx, &types.CreateArtifactRequest{
		Artifact: artifactValue,
		Host:     &host,
	})
	if err != nil {
		t.Fatalf("CreateArtifact failed: %v", err)
	}

	t.Logf("Created artifact for update test: ID %d", artifact.ID)

	// Update the artifact
	newHost := "NEW-HOST"
	updated, err := client.UpdateArtifact(ctx, &types.UpdateArtifactRequest{
		ID:   artifact.ID,
		Host: &newHost,
	})
	if err != nil {
		t.Fatalf("UpdateArtifact failed: %v", err)
	}

	if updated.Host != newHost {
		t.Errorf("Expected host %q, got %q", newHost, updated.Host)
	}

	t.Logf("Successfully updated artifact")
	t.Logf("  - Old host: %s", host)
	t.Logf("  - New host: %s", updated.Host)

	// Clean up
	err = client.DeleteArtifact(ctx, artifact.ID)
	if err != nil {
		t.Logf("Warning: Failed to delete test artifact: %v", err)
	}
}

func TestArtifacts_DeleteArtifact(t *testing.T) {

	client := AuthenticateTestClient(t)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Create a test artifact
	artifactValue := "C:\\Windows\\Temp\\test-delete-" + time.Now().Format("20060102150405") + ".exe"

	artifact, err := client.CreateArtifact(ctx, &types.CreateArtifactRequest{
		Artifact: artifactValue,
	})
	if err != nil {
		t.Fatalf("CreateArtifact failed: %v", err)
	}

	t.Logf("Created artifact for delete test: ID %d", artifact.ID)

	// Delete the artifact
	err = client.DeleteArtifact(ctx, artifact.ID)
	if err != nil {
		t.Fatalf("DeleteArtifact failed: %v", err)
	}

	t.Log("Successfully deleted artifact")

	// Verify deletion - artifact should not be retrievable after hard delete
	_, err = client.GetArtifactByID(ctx, artifact.ID)
	if err != nil {
		// If it's not found, that's expected for hard delete
		t.Logf("Artifact not found after deletion (expected): %v", err)
		return
	}

	t.Error("Artifact should not be retrievable after deletion")
}

func TestArtifacts_GetArtifactsByHost(t *testing.T) {

	client := AuthenticateTestClient(t)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Create test artifacts on a specific host
	testHost := "TEST-HOST-" + time.Now().Format("20060102150405")

	artifact1, err := client.CreateArtifact(ctx, &types.CreateArtifactRequest{
		Artifact: "C:\\test1.exe",
		Host:     &testHost,
	})
	if err != nil {
		t.Fatalf("CreateArtifact 1 failed: %v", err)
	}

	artifact2, err := client.CreateArtifact(ctx, &types.CreateArtifactRequest{
		Artifact: "C:\\test2.dll",
		Host:     &testHost,
	})
	if err != nil {
		t.Fatalf("CreateArtifact 2 failed: %v", err)
	}

	t.Logf("Created 2 artifacts on host %s", testHost)

	// Get artifacts by host
	artifacts, err := client.GetArtifactsByHost(ctx, testHost)
	if err != nil {
		t.Fatalf("GetArtifactsByHost failed: %v", err)
	}

	if len(artifacts) < 2 {
		t.Errorf("Expected at least 2 artifacts, got %d", len(artifacts))
	}

	// Verify all artifacts are from the correct host
	for _, a := range artifacts {
		if a.Host != testHost {
			t.Errorf("Expected host %q, got %q", testHost, a.Host)
		}
	}

	t.Logf("Found %d artifacts for host %s", len(artifacts), testHost)

	// Clean up
	_ = client.DeleteArtifact(ctx, artifact1.ID)
	_ = client.DeleteArtifact(ctx, artifact2.ID)
}

func TestArtifacts_TimestampOrdering(t *testing.T) {

	client := AuthenticateTestClient(t)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	artifacts, err := client.GetArtifacts(ctx)
	if err != nil {
		t.Fatalf("GetArtifacts failed: %v", err)
	}

	if len(artifacts) < 2 {
		t.Skip("Need at least 2 artifacts to test timestamp ordering")
	}

	// Verify descending order (newest first)
	for i := 1; i < len(artifacts); i++ {
		if artifacts[i].Timestamp.After(artifacts[i-1].Timestamp) {
			t.Errorf("Timestamp ordering broken at index %d: %s > %s",
				i,
				artifacts[i].Timestamp.Format("2006-01-02 15:04:05"),
				artifacts[i-1].Timestamp.Format("2006-01-02 15:04:05"))
		}
	}

	t.Log("Timestamp ordering verified (newest first)")
	if len(artifacts) > 0 {
		t.Logf("  - Newest: %s", artifacts[0].Timestamp.Format("2006-01-02 15:04:05"))
		t.Logf("  - Oldest: %s", artifacts[len(artifacts)-1].Timestamp.Format("2006-01-02 15:04:05"))
	}
}

func TestArtifacts_GetArtifactsByOperation(t *testing.T) {

	client := AuthenticateTestClient(t)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Get current operation
	currentOpID := client.GetCurrentOperation()
	if currentOpID == nil {
		t.Skip("No current operation set")
	}

	// Get artifacts for the operation
	artifacts, err := client.GetArtifactsByOperation(ctx, *currentOpID)
	if err != nil {
		t.Fatalf("GetArtifactsByOperation failed: %v", err)
	}

	if artifacts == nil {
		t.Fatal("GetArtifactsByOperation returned nil")
	}

	t.Logf("Found %d artifact(s) for operation %d", len(artifacts), *currentOpID)

	// Verify all artifacts belong to the operation
	for i, artifact := range artifacts {
		if artifact.OperationID != *currentOpID {
			t.Errorf("Artifact %d belongs to operation %d, expected %d",
				i, artifact.OperationID, *currentOpID)
		}

		// Log first few artifacts
		if i < 5 {
			t.Logf("  - Artifact %d: %s (host: %s)",
				i, artifact.Artifact, artifact.Host)
		}
	}

	if len(artifacts) > 5 {
		t.Logf("  ... and %d more artifacts", len(artifacts)-5)
	}

	// Verify artifacts are sorted by timestamp (descending)
	if len(artifacts) > 1 {
		for i := 1; i < len(artifacts); i++ {
			if artifacts[i].Timestamp.After(artifacts[i-1].Timestamp) {
				t.Error("Artifacts should be sorted by timestamp (descending/newest first)")
				break
			}
		}
	}
}

func TestArtifacts_GetArtifactsByOperation_InvalidID(t *testing.T) {

	client := AuthenticateTestClient(t)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Test with zero ID
	_, err := client.GetArtifactsByOperation(ctx, 0)
	if err == nil {
		t.Fatal("Expected error for zero operation ID, got nil")
	}
	t.Logf("Zero ID error: %v", err)

	// Test with non-existent operation ID
	_, err = client.GetArtifactsByOperation(ctx, 999999)
	if err == nil {
		t.Fatal("Expected error for non-existent operation ID, got nil")
	}
	t.Logf("Non-existent ID error: %v", err)
}
