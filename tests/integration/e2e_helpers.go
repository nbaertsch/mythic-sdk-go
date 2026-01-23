//go:build integration

package integration

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"

	"github.com/nbaertsch/mythic-sdk-go/pkg/mythic"
	"github.com/nbaertsch/mythic-sdk-go/pkg/mythic/types"
)

// E2ETestSetup contains all resources needed for E2E testing
type E2ETestSetup struct {
	Client       *mythic.Client
	OperationID  int
	PayloadUUID  string
	PayloadPath  string
	CallbackID   int
	AgentProcess *os.Process
	TempFiles    []string
	ctx          context.Context
	cancel       context.CancelFunc
	t            *testing.T
}

// SetupE2ETest creates a full E2E test environment
// This is used by agent-dependent tests (Workflows 9-12)
func SetupE2ETest(t *testing.T) *E2ETestSetup {
	t.Helper()

	// Create authenticated client
	client := AuthenticateTestClient(t)

	// Get current operation
	operationID := client.GetCurrentOperation()
	if operationID == nil {
		t.Fatal("No current operation set")
	}

	// Create context with extended timeout for E2E operations
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)

	setup := &E2ETestSetup{
		Client:      client,
		OperationID: *operationID,
		TempFiles:   make([]string, 0),
		ctx:         ctx,
		cancel:      cancel,
		t:           t,
	}

	return setup
}

// Cleanup performs full cleanup of E2E test resources
func (s *E2ETestSetup) Cleanup() {
	s.t.Helper()

	// Kill agent process if running
	if s.AgentProcess != nil {
		s.t.Logf("Killing agent process (PID: %d)", s.AgentProcess.Pid)
		_ = s.AgentProcess.Kill()
		_, _ = s.AgentProcess.Wait()
	}

	// Delete callback if exists
	if s.CallbackID > 0 {
		s.t.Logf("Deleting callback ID: %d", s.CallbackID)
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		_ = s.Client.DeleteCallback(ctx, []int{s.CallbackID})
	}

	// Delete payload if exists
	if s.PayloadUUID != "" {
		s.t.Logf("Deleting payload UUID: %s", s.PayloadUUID)
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		_ = s.Client.DeletePayload(ctx, s.PayloadUUID)
	}

	// Remove payload file
	if s.PayloadPath != "" {
		s.t.Logf("Removing payload file: %s", s.PayloadPath)
		_ = os.Remove(s.PayloadPath)
	}

	// Remove temp files
	for _, file := range s.TempFiles {
		s.t.Logf("Removing temp file: %s", file)
		_ = os.Remove(file)
	}

	// Cancel context
	s.cancel()
}

// WaitForCallback polls for a callback to be established
func (s *E2ETestSetup) WaitForCallback(timeout time.Duration) (int, error) {
	s.t.Helper()

	deadline := time.Now().Add(timeout)
	s.t.Logf("Waiting for callback (timeout: %v)", timeout)

	initialCount := 0
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	callbacks, err := s.Client.GetAllActiveCallbacks(ctx)
	cancel()
	if err == nil {
		initialCount = len(callbacks)
		s.t.Logf("Initial active callbacks: %d", initialCount)
	}

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for time.Now().Before(deadline) {
		<-ticker.C

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		callbacks, err := s.Client.GetAllActiveCallbacks(ctx)
		cancel()

		if err != nil {
			s.t.Logf("Error checking callbacks: %v", err)
			continue
		}

		s.t.Logf("Current active callbacks: %d", len(callbacks))

		// Look for new callback
		if len(callbacks) > initialCount {
			// Get the newest callback
			newestCallback := callbacks[len(callbacks)-1]
			s.t.Logf("Found new callback! ID: %d, Agent: %s", newestCallback.ID, newestCallback.AgentCallbackID)
			s.CallbackID = newestCallback.ID
			return newestCallback.ID, nil
		}
	}

	return 0, fmt.Errorf("timeout waiting for callback after %v", timeout)
}

// ExecuteCommand executes a command on the callback and returns the task display ID
func (s *E2ETestSetup) ExecuteCommand(cmd string, params string) (int, error) {
	s.t.Helper()

	if s.CallbackID == 0 {
		return 0, fmt.Errorf("no callback established")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	s.t.Logf("Executing command: %s (params: %s)", cmd, params)

	callbackID := s.CallbackID
	taskReq := &mythic.TaskRequest{
		CallbackID: &callbackID,
		Command:    cmd,
		Params:     params,
	}

	task, err := s.Client.IssueTask(ctx, taskReq)
	if err != nil {
		return 0, fmt.Errorf("failed to issue task: %w", err)
	}

	s.t.Logf("Task issued: Display ID %d (ID: %d)", task.DisplayID, task.ID)
	return task.DisplayID, nil
}

// WaitForTaskComplete polls for task completion and returns concatenated output.
// Automatically requests OPSEC bypass if task is blocked (for automated testing).
func (s *E2ETestSetup) WaitForTaskComplete(taskDisplayID int, timeout time.Duration) (string, error) {
	s.t.Helper()

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	s.t.Logf("Waiting for task %d to complete (timeout: %v, auto-bypass OPSEC: enabled)", taskDisplayID, timeout)

	timeoutSeconds := int(timeout.Seconds())
	// Enable automatic OPSEC bypass for E2E tests
	// This allows tests to run in automated environments without manual approval
	err := s.Client.WaitForTaskCompleteWithOptions(ctx, taskDisplayID, timeoutSeconds, true)
	if err != nil {
		return "", fmt.Errorf("task did not complete: %w", err)
	}

	// Get task output
	outputCtx, outputCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer outputCancel()

	responses, err := s.Client.GetTaskOutput(outputCtx, taskDisplayID)
	if err != nil {
		return "", fmt.Errorf("failed to get task output: %w", err)
	}

	// Concatenate all response outputs
	var output string
	for _, resp := range responses {
		output += resp.ResponseText
	}

	s.t.Logf("Task %d completed with %d responses, total output: %d bytes", taskDisplayID, len(responses), len(output))
	return output, nil
}

// StartAgent starts the payload agent in the background
func (s *E2ETestSetup) StartAgent() error {
	s.t.Helper()

	if s.PayloadPath == "" {
		return fmt.Errorf("no payload path set")
	}

	// Make payload executable
	if err := os.Chmod(s.PayloadPath, 0755); err != nil {
		return fmt.Errorf("failed to chmod payload: %w", err)
	}

	s.t.Logf("Starting agent: %s", s.PayloadPath)

	cmd := exec.Command(s.PayloadPath)
	cmd.Stdout = nil // Suppress output
	cmd.Stderr = nil

	err := cmd.Start()
	if err != nil {
		return fmt.Errorf("failed to start agent: %w", err)
	}

	s.AgentProcess = cmd.Process
	s.t.Logf("Agent started with PID: %d", s.AgentProcess.Pid)

	return nil
}

// AddTempFile tracks a temporary file for cleanup
func (s *E2ETestSetup) AddTempFile(path string) {
	s.TempFiles = append(s.TempFiles, path)
}

// GetContext returns the E2E test context
func (s *E2ETestSetup) GetContext() context.Context {
	return s.ctx
}

// EnsureCallbackExists checks if any active callbacks exist, and if not,
// creates a payload and starts an agent in the background to establish one.
// This should be called at the start of tests that need a callback but don't
// care about creating it themselves. Returns the callback ID.
//
// This function is designed to avoid creating multiple callbacks - it will
// reuse existing callbacks if available to prevent resource exhaustion.
func EnsureCallbackExists(t *testing.T) int {
	t.Helper()

	client := AuthenticateTestClient(t)

	// Check if any active callbacks already exist
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	callbacks, err := client.GetAllActiveCallbacks(ctx)
	if err != nil {
		t.Fatalf("Failed to check for existing callbacks: %v", err)
	}

	if len(callbacks) > 0 {
		t.Logf("Using existing callback ID: %d (no need to create new one)", callbacks[0].ID)
		return callbacks[0].ID
	}

	// No callbacks exist - need to create one
	t.Log("No active callbacks found - creating shared callback for tests")

	// Create E2E setup
	setup := SetupE2ETest(t)

	// Get Poseidon payload type
	t.Log("Finding Poseidon payload type...")
	ctx1, cancel1 := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel1()

	payloadTypes, err := client.GetPayloadTypes(ctx1)
	if err != nil {
		t.Fatalf("GetPayloadTypes failed: %v", err)
	}

	var poseidonType *types.PayloadType
	for i := range payloadTypes {
		if payloadTypes[i].Name == "poseidon" {
			poseidonType = &payloadTypes[i]
			break
		}
	}

	if poseidonType == nil {
		t.Skip("Poseidon payload type not found - cannot create callback")
	}

	if !poseidonType.ContainerRunning {
		t.Skip("Poseidon container not running - cannot create callback")
	}

	// Get C2 profile
	ctx2, cancel2 := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel2()

	c2Profiles, err := client.GetC2Profiles(ctx2)
	if err != nil {
		t.Fatalf("GetC2Profiles failed: %v", err)
	}

	var c2Profile *types.C2Profile
	for i := range c2Profiles {
		if c2Profiles[i].Name == "http" {
			c2Profile = &c2Profiles[i]
			break
		}
	}
	if c2Profile == nil && len(c2Profiles) > 0 {
		c2Profile = &c2Profiles[0]
	}
	if c2Profile == nil {
		t.Fatal("No C2 profiles available")
	}

	// Create payload
	t.Log("Creating Poseidon payload...")
	ctx3, cancel3 := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel3()

	payloadReq := &types.CreatePayloadRequest{
		PayloadType: "poseidon",
		OS:          "linux",
		Description: "Shared Test Callback",
		Filename:    "shared_test_agent",
		Commands: []string{
			"shell", "ps", "whoami",
		},
		BuildParameters: map[string]interface{}{
			"mode": "default",
		},
		C2Profiles: []types.C2ProfileConfig{
			{
				Name: c2Profile.Name,
				Parameters: map[string]interface{}{
					"callback_host": "http://127.0.0.1",
					"callback_port": 80,
				},
			},
		},
	}

	payload, err := client.CreatePayload(ctx3, payloadReq)
	if err != nil {
		t.Fatalf("CreatePayload failed: %v", err)
	}
	setup.PayloadUUID = payload.UUID
	t.Logf("Payload created: %s", payload.UUID)

	// Wait for build
	t.Log("Waiting for payload build...")
	ctx4, cancel4 := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel4()

	err = client.WaitForPayloadComplete(ctx4, payload.UUID, 90)
	if err != nil {
		t.Fatalf("Payload build failed: %v", err)
	}

	// Download payload
	t.Log("Downloading payload...")
	ctx5, cancel5 := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel5()

	payloadBytes, err := client.DownloadPayload(ctx5, payload.UUID)
	if err != nil {
		t.Fatalf("DownloadPayload failed: %v", err)
	}

	// Save payload to temp file
	tmpDir := os.TempDir()
	payloadPath := filepath.Join(tmpDir, "shared_test_agent_"+payload.UUID[:8])
	err = os.WriteFile(payloadPath, payloadBytes, 0755)
	if err != nil {
		t.Fatalf("Failed to write payload: %v", err)
	}
	setup.PayloadPath = payloadPath
	t.Logf("Payload saved: %s (%d bytes)", payloadPath, len(payloadBytes))

	// Register cleanup to run at test end
	t.Cleanup(func() {
		t.Log("Cleaning up shared callback resources...")
		setup.Cleanup()
	})

	// Start agent in BACKGROUND goroutine
	t.Log("Starting agent in background...")
	agentStarted := make(chan error, 1)
	go func() {
		if err := setup.StartAgent(); err != nil {
			agentStarted <- err
		} else {
			agentStarted <- nil
		}
	}()

	// Wait for agent to start
	select {
	case err := <-agentStarted:
		if err != nil {
			t.Fatalf("Failed to start agent in background: %v", err)
		}
		t.Log("Agent started in background")
	case <-time.After(10 * time.Second):
		t.Fatal("Timeout waiting for agent to start")
	}

	// Wait for callback to establish
	t.Log("Waiting for callback to establish...")
	callbackID, err := setup.WaitForCallback(90 * time.Second)
	if err != nil {
		t.Fatalf("Failed to establish callback: %v", err)
	}

	t.Logf("✓ Shared callback established: ID %d", callbackID)
	return callbackID
}
