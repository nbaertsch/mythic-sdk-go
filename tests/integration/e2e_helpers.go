//go:build integration

package integration

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"testing"
	"time"

	"github.com/nbaertsch/mythic-sdk-go/pkg/mythic"
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

// WaitForTaskComplete polls for task completion and returns concatenated output
func (s *E2ETestSetup) WaitForTaskComplete(taskDisplayID int, timeout time.Duration) (string, error) {
	s.t.Helper()

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	s.t.Logf("Waiting for task %d to complete (timeout: %v)", taskDisplayID, timeout)

	timeoutSeconds := int(timeout.Seconds())
	err := s.Client.WaitForTaskComplete(ctx, taskDisplayID, timeoutSeconds)
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
