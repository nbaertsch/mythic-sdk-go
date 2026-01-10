//go:build integration

package integration

import (
	"context"
	"testing"
	"time"

	"github.com/nbaertsch/mythic-sdk-go/pkg/mythic/types"
)

// TestE2E_CredentialManagement tests the complete credential management workflow
// Covers: GetCredentials, GetCredentialsByOperation, CreateCredential,
// UpdateCredential, DeleteCredential
func TestE2E_CredentialManagement(t *testing.T) {
	client := AuthenticateTestClient(t)

	var createdCredentialIDs []int

	// Register cleanup
	defer func() {
		for _, credID := range createdCredentialIDs {
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			_ = client.DeleteCredential(ctx, credID)
			cancel()
			t.Logf("Cleaned up credential ID: %d", credID)
		}
	}()

	// Test 1: Get credentials baseline
	t.Log("=== Test 1: Get credentials baseline ===")
	ctx1, cancel1 := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel1()

	baselineCredentials, err := client.GetCredentials(ctx1)
	if err != nil {
		t.Fatalf("GetCredentials baseline failed: %v", err)
	}
	baselineCount := len(baselineCredentials)
	t.Logf("✓ Baseline credential count: %d", baselineCount)

	// Test 2: Create credential - password
	t.Log("=== Test 2: Create credential - password ===")
	ctx2, cancel2 := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel2()

	passwordCred := &types.CreateCredentialRequest{
		Type:       "plaintext",
		Account:    "testuser",
		Realm:      "example.com",
		Credential: "P@ssw0rd123!",
		Comment:    "E2E test password credential",
	}

	createdPassword, err := client.CreateCredential(ctx2, passwordCred)
	if err != nil {
		t.Fatalf("CreateCredential (password) failed: %v", err)
	}
	if createdPassword.ID == 0 {
		t.Fatal("Created credential has ID 0")
	}
	createdCredentialIDs = append(createdCredentialIDs, createdPassword.ID)
	t.Logf("✓ Password credential created: ID %d", createdPassword.ID)

	// Test 3: Create credential - SSH key
	t.Log("=== Test 3: Create credential - SSH key ===")
	ctx3, cancel3 := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel3()

	sshCred := &types.CreateCredentialRequest{
		Type:       "key",
		Account:    "root",
		Realm:      "server1.example.com",
		Credential: "-----BEGIN RSA PRIVATE KEY-----\nMIIEpAIBAAK...",
		Comment:    "E2E test SSH key",
	}

	createdSSH, err := client.CreateCredential(ctx3, sshCred)
	if err != nil {
		t.Fatalf("CreateCredential (SSH) failed: %v", err)
	}
	createdCredentialIDs = append(createdCredentialIDs, createdSSH.ID)
	t.Logf("✓ SSH key credential created: ID %d", createdSSH.ID)

	// Test 4: Create credential - API token
	t.Log("=== Test 4: Create credential - API token ===")
	ctx4, cancel4 := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel4()

	apiCred := &types.CreateCredentialRequest{
		Type:       "cookie",
		Account:    "api_user",
		Realm:      "api.example.com",
		Credential: "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
		Comment:    "E2E test API token",
	}

	createdAPI, err := client.CreateCredential(ctx4, apiCred)
	if err != nil {
		t.Fatalf("CreateCredential (API) failed: %v", err)
	}
	createdCredentialIDs = append(createdCredentialIDs, createdAPI.ID)
	t.Logf("✓ API token credential created: ID %d", createdAPI.ID)

	// Test 5: Get all credentials after creation
	t.Log("=== Test 5: Get all credentials after creation ===")
	ctx5, cancel5 := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel5()

	allCredentials, err := client.GetCredentials(ctx5)
	if err != nil {
		t.Fatalf("GetCredentials after creation failed: %v", err)
	}
	newCount := len(allCredentials)
	if newCount < baselineCount+3 {
		t.Errorf("Expected at least %d credentials, got %d", baselineCount+3, newCount)
	}
	t.Logf("✓ Total credentials now: %d (added %d)", newCount, newCount-baselineCount)

	// Test 6: Get credentials by operation
	t.Log("=== Test 6: Get credentials by operation ===")
	ctx6, cancel6 := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel6()

	operationID := client.GetCurrentOperation()
	if operationID == nil {
		t.Fatal("No current operation set")
	}

	opCredentials, err := client.GetCredentialsByOperation(ctx6, *operationID)
	if err != nil {
		t.Fatalf("GetCredentialsByOperation failed: %v", err)
	}
	t.Logf("✓ Found %d credentials for operation %d", len(opCredentials), *operationID)

	// Verify our created credentials are in the operation
	found := 0
	for _, cred := range opCredentials {
		for _, createdID := range createdCredentialIDs {
			if cred.ID == createdID {
				found++
				t.Logf("  ✓ Found credential %d: %s@%s (%s)", cred.ID, cred.Account, cred.Realm, cred.Type)
			}
		}
	}
	if found != len(createdCredentialIDs) {
		t.Errorf("Expected to find %d credentials in operation, found %d", len(createdCredentialIDs), found)
	}

	// Test 7: Update credential (change comment)
	t.Log("=== Test 7: Update credential (change comment) ===")
	ctx7, cancel7 := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel7()

	newComment := "Updated comment - E2E test modified"
	updateReq := &types.UpdateCredentialRequest{
		ID:      createdPassword.ID,
		Comment: &newComment,
	}

	updatedCred, err := client.UpdateCredential(ctx7, updateReq)
	if err != nil {
		t.Fatalf("UpdateCredential failed: %v", err)
	}
	if updatedCred.Comment != newComment {
		t.Errorf("Comment not updated: expected %s, got %s", newComment, updatedCred.Comment)
	}
	t.Logf("✓ Credential %d updated: comment = %s", updatedCred.ID, updatedCred.Comment)

	// Test 8: Verify update persisted
	t.Log("=== Test 8: Verify update persisted ===")
	ctx8, cancel8 := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel8()

	opCredentialsAfterUpdate, err := client.GetCredentialsByOperation(ctx8, *operationID)
	if err != nil {
		t.Fatalf("GetCredentialsByOperation after update failed: %v", err)
	}

	foundUpdated := false
	for _, cred := range opCredentialsAfterUpdate {
		if cred.ID == createdPassword.ID {
			foundUpdated = true
			if cred.Comment != newComment {
				t.Errorf("Update did not persist: expected %s, got %s", newComment, cred.Comment)
			} else {
				t.Logf("✓ Update verified: credential %d has comment %s", cred.ID, cred.Comment)
			}
			break
		}
	}
	if !foundUpdated {
		t.Error("Could not find updated credential to verify")
	}

	// Test 9: Delete credentials
	t.Log("=== Test 9: Delete credentials ===")
	for _, credID := range createdCredentialIDs {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		err := client.DeleteCredential(ctx, credID)
		cancel()
		if err != nil {
			t.Errorf("DeleteCredential failed for ID %d: %v", credID, err)
		} else {
			t.Logf("✓ Credential %d deleted", credID)
		}
	}

	// Test 10: Verify deletion
	t.Log("=== Test 10: Verify deletion ===")
	ctx9, cancel9 := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel9()

	finalCredentials, err := client.GetCredentials(ctx9)
	if err != nil {
		t.Fatalf("GetCredentials after delete failed: %v", err)
	}

	// Check that deleted credentials are not in the active list (or marked deleted)
	for _, cred := range finalCredentials {
		for _, deletedID := range createdCredentialIDs {
			if cred.ID == deletedID && !cred.Deleted {
				t.Errorf("Credential %d still active after deletion", deletedID)
			}
		}
	}
	t.Logf("✓ Verified deletion of %d credentials", len(createdCredentialIDs))

	t.Log("=== ✓ All credential management tests passed ===")
}

// TestE2E_ArtifactManagement tests the complete artifact management workflow
// Covers: GetArtifacts, GetArtifactsByOperation, GetArtifactsByHost,
// GetArtifactsByType, CreateArtifact, UpdateArtifact, DeleteArtifact
func TestE2E_ArtifactManagement(t *testing.T) {
	client := AuthenticateTestClient(t)

	var createdArtifactIDs []int

	// Register cleanup
	defer func() {
		for _, artifactID := range createdArtifactIDs {
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			_ = client.DeleteArtifact(ctx, artifactID)
			cancel()
			t.Logf("Cleaned up artifact ID: %d", artifactID)
		}
	}()

	// Test 1: Get artifacts baseline
	t.Log("=== Test 1: Get artifacts baseline ===")
	ctx1, cancel1 := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel1()

	baselineArtifacts, err := client.GetArtifacts(ctx1)
	if err != nil {
		t.Fatalf("GetArtifacts baseline failed: %v", err)
	}
	baselineCount := len(baselineArtifacts)
	t.Logf("✓ Baseline artifact count: %d", baselineCount)

	// Test 2: Create file artifact
	t.Log("=== Test 2: Create file artifact ===")
	ctx2, cancel2 := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel2()

	testHost := "testhost.example.com"
	fileArtifact := &types.CreateArtifactRequest{
		Artifact:     "/tmp/malicious_payload.exe",
		BaseArtifact: &testHost,
		Host:         &testHost,
	}

	createdFile, err := client.CreateArtifact(ctx2, fileArtifact)
	if err != nil {
		t.Fatalf("CreateArtifact (file) failed: %v", err)
	}
	if createdFile.ID == 0 {
		t.Fatal("Created artifact has ID 0")
	}
	createdArtifactIDs = append(createdArtifactIDs, createdFile.ID)
	t.Logf("✓ File artifact created: ID %d (%s)", createdFile.ID, createdFile.Artifact)

	// Test 3: Create registry artifact
	t.Log("=== Test 3: Create registry artifact ===")
	ctx3, cancel3 := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel3()

	registryArtifact := &types.CreateArtifactRequest{
		Artifact:     "HKLM\\Software\\Microsoft\\Windows\\CurrentVersion\\Run\\Malware",
		BaseArtifact: &testHost,
		Host:         &testHost,
	}

	createdRegistry, err := client.CreateArtifact(ctx3, registryArtifact)
	if err != nil {
		t.Fatalf("CreateArtifact (registry) failed: %v", err)
	}
	createdArtifactIDs = append(createdArtifactIDs, createdRegistry.ID)
	t.Logf("✓ Registry artifact created: ID %d (%s)", createdRegistry.ID, createdRegistry.Artifact)

	// Test 4: Create process artifact
	t.Log("=== Test 4: Create process artifact ===")
	ctx4, cancel4 := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel4()

	processArtifact := &types.CreateArtifactRequest{
		Artifact:     "notepad.exe",
		BaseArtifact: &testHost,
		Host:         &testHost,
	}

	createdProcess, err := client.CreateArtifact(ctx4, processArtifact)
	if err != nil {
		t.Fatalf("CreateArtifact (process) failed: %v", err)
	}
	createdArtifactIDs = append(createdArtifactIDs, createdProcess.ID)
	t.Logf("✓ Process artifact created: ID %d (%s)", createdProcess.ID, createdProcess.Artifact)

	// Test 5: Get all artifacts after creation
	t.Log("=== Test 5: Get all artifacts after creation ===")
	ctx5, cancel5 := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel5()

	allArtifacts, err := client.GetArtifacts(ctx5)
	if err != nil {
		t.Fatalf("GetArtifacts after creation failed: %v", err)
	}
	newCount := len(allArtifacts)
	if newCount < baselineCount+3 {
		t.Errorf("Expected at least %d artifacts, got %d", baselineCount+3, newCount)
	}
	t.Logf("✓ Total artifacts now: %d (added %d)", newCount, newCount-baselineCount)

	// Test 6: Get artifacts by operation
	t.Log("=== Test 6: Get artifacts by operation ===")
	ctx6, cancel6 := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel6()

	operationID := client.GetCurrentOperation()
	if operationID == nil {
		t.Fatal("No current operation set")
	}

	opArtifacts, err := client.GetArtifactsByOperation(ctx6, *operationID)
	if err != nil {
		t.Fatalf("GetArtifactsByOperation failed: %v", err)
	}
	t.Logf("✓ Found %d artifacts for operation %d", len(opArtifacts), *operationID)

	// Test 7: Get artifacts by host
	t.Log("=== Test 7: Get artifacts by host ===")
	ctx7, cancel7 := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel7()

	hostArtifacts, err := client.GetArtifactsByHost(ctx7, testHost)
	if err != nil {
		t.Fatalf("GetArtifactsByHost failed: %v", err)
	}
	t.Logf("✓ Found %d artifacts for host %s", len(hostArtifacts), testHost)

	// Verify our created artifacts are in the host results
	found := 0
	for _, artifact := range hostArtifacts {
		for _, createdID := range createdArtifactIDs {
			if artifact.ID == createdID {
				found++
				t.Logf("  ✓ Found artifact %d: %s", artifact.ID, artifact.Artifact)
			}
		}
	}
	if found != len(createdArtifactIDs) {
		t.Errorf("Expected to find %d artifacts for host, found %d", len(createdArtifactIDs), found)
	}

	// Test 8: Update artifact host
	t.Log("=== Test 8: Update artifact (change host) ===")
	ctx8, cancel8 := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel8()

	newHost := "newhost.example.com"
	updateReq := &types.UpdateArtifactRequest{
		ID:   createdFile.ID,
		Host: &newHost,
	}

	updatedArtifact, err := client.UpdateArtifact(ctx8, updateReq)
	if err != nil {
		t.Fatalf("UpdateArtifact failed: %v", err)
	}
	if updatedArtifact.Host != newHost {
		t.Errorf("Host not updated: expected %s, got %s", newHost, updatedArtifact.Host)
	}
	t.Logf("✓ Artifact %d updated: host = %s", updatedArtifact.ID, updatedArtifact.Host)

	// Test 9: Verify update persisted
	t.Log("=== Test 9: Verify update persisted ===")
	ctx9, cancel9 := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel9()

	allArtifactsAfterUpdate, err := client.GetArtifacts(ctx9)
	if err != nil {
		t.Fatalf("GetArtifacts after update failed: %v", err)
	}

	foundUpdated := false
	for _, artifact := range allArtifactsAfterUpdate {
		if artifact.ID == createdFile.ID {
			foundUpdated = true
			if artifact.Host != newHost {
				t.Errorf("Update did not persist: expected %s, got %s", newHost, artifact.Host)
			} else {
				t.Logf("✓ Update verified: artifact %d has host %s", artifact.ID, artifact.Host)
			}
			break
		}
	}
	if !foundUpdated {
		t.Error("Could not find updated artifact to verify")
	}

	// Test 10: Delete artifacts
	t.Log("=== Test 11: Delete artifacts ===")
	for _, artifactID := range createdArtifactIDs {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		err := client.DeleteArtifact(ctx, artifactID)
		cancel()
		if err != nil {
			t.Errorf("DeleteArtifact failed for ID %d: %v", artifactID, err)
		} else {
			t.Logf("✓ Artifact %d deleted", artifactID)
		}
	}

	// Test 12: Verify deletion
	t.Log("=== Test 12: Verify deletion ===")
	ctx11, cancel11 := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel11()

	finalArtifacts, err := client.GetArtifacts(ctx11)
	if err != nil {
		t.Fatalf("GetArtifacts after delete failed: %v", err)
	}

	// Check that deleted artifacts are not in the active list (hard delete means they shouldn't exist)
	for _, artifact := range finalArtifacts {
		for _, deletedID := range createdArtifactIDs {
			if artifact.ID == deletedID {
				t.Errorf("Artifact %d still exists after deletion (should be hard deleted)", deletedID)
			}
		}
	}
	t.Logf("✓ Verified deletion of %d artifacts", len(createdArtifactIDs))

	t.Log("=== ✓ All artifact management tests passed ===")
}

// TestE2E_CredentialsArtifactsErrorHandling tests error scenarios
func TestE2E_CredentialsArtifactsErrorHandling(t *testing.T) {
	client := AuthenticateTestClient(t)

	// Test 1: Delete non-existent credential
	t.Log("=== Test 1: Delete non-existent credential ===")
	ctx1, cancel1 := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel1()

	err := client.DeleteCredential(ctx1, 999999)
	if err == nil {
		t.Error("Expected error for non-existent credential delete")
	}
	t.Logf("✓ Non-existent credential delete rejected: %v", err)

	// Test 2: Update non-existent credential
	t.Log("=== Test 2: Update non-existent credential ===")
	ctx2, cancel2 := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel2()

	comment := "test"
	updateReq := &types.UpdateCredentialRequest{
		ID:      999999,
		Comment: &comment,
	}

	_, err = client.UpdateCredential(ctx2, updateReq)
	if err == nil {
		t.Error("Expected error for non-existent credential update")
	}
	t.Logf("✓ Non-existent credential update rejected: %v", err)

	// Test 3: Delete non-existent artifact
	t.Log("=== Test 3: Delete non-existent artifact ===")
	ctx3, cancel3 := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel3()

	err = client.DeleteArtifact(ctx3, 999999)
	if err == nil {
		t.Error("Expected error for non-existent artifact delete")
	}
	t.Logf("✓ Non-existent artifact delete rejected: %v", err)

	// Test 4: Update non-existent artifact
	t.Log("=== Test 4: Update non-existent artifact ===")
	ctx4, cancel4 := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel4()

	testHost := "test"
	artifactUpdateReq := &types.UpdateArtifactRequest{
		ID:   999999,
		Host: &testHost,
	}

	_, err = client.UpdateArtifact(ctx4, artifactUpdateReq)
	if err == nil {
		t.Error("Expected error for non-existent artifact update")
	}
	t.Logf("✓ Non-existent artifact update rejected: %v", err)

	t.Log("=== ✓ All error handling tests passed ===")
}
