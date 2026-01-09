//go:build integration

package integration

import (
	"context"
	"testing"
	"time"

	"github.com/nbaertsch/mythic-sdk-go/pkg/mythic/types"
)

func TestResponses_GetResponsesByTask(t *testing.T) {
	SkipIfNoMythic(t)
	client := AuthenticateTestClient(t)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Get tasks first to find a valid task ID
	tasks, err := client.GetAllTasks(ctx, 10)
	if err != nil {
		t.Fatalf("GetAllTasks failed: %v", err)
	}
	if len(tasks) == 0 {
		t.Skip("No tasks available for testing")
	}

	taskID := tasks[0].ID
	responses, err := client.GetResponsesByTask(ctx, taskID)
	if err != nil {
		t.Fatalf("GetResponsesByTask failed: %v", err)
	}

	if responses == nil {
		t.Fatal("GetResponsesByTask returned nil")
	}

	t.Logf("Found %d response(s) for task %d", len(responses), taskID)

	// Verify responses belong to the task
	for _, resp := range responses {
		if resp.ID == 0 {
			t.Error("Response ID should not be 0")
		}
		if resp.Task != nil && resp.Task.ID != taskID {
			t.Errorf("Expected task ID %d, got %d", taskID, resp.Task.ID)
		}
	}

	if len(responses) > 0 {
		t.Logf("First response: %s", responses[0].String())
	}
}

func TestResponses_GetResponseByID(t *testing.T) {
	SkipIfNoMythic(t)
	client := AuthenticateTestClient(t)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Get latest responses to find a valid response ID
	responses, err := client.GetLatestResponses(ctx, 0, 5)
	if err != nil {
		t.Fatalf("GetLatestResponses failed: %v", err)
	}
	if len(responses) == 0 {
		t.Skip("No responses available for testing")
	}

	responseID := responses[0].ID
	response, err := client.GetResponseByID(ctx, responseID)
	if err != nil {
		t.Fatalf("GetResponseByID failed: %v", err)
	}

	if response == nil {
		t.Fatal("GetResponseByID returned nil")
	}

	if response.ID != responseID {
		t.Errorf("Expected response ID %d, got %d", responseID, response.ID)
	}

	t.Logf("Retrieved response %d: %s", responseID, response.String())
	if response.HasOutput() {
		t.Logf("  - Output length: %d characters", len(response.ResponseText))
	}
}

func TestResponses_GetResponsesByCallback(t *testing.T) {
	SkipIfNoMythic(t)
	client := AuthenticateTestClient(t)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Get active callbacks
	callbacks, err := client.GetAllActiveCallbacks(ctx)
	if err != nil {
		t.Fatalf("GetAllActiveCallbacks failed: %v", err)
	}
	if len(callbacks) == 0 {
		t.Skip("No active callbacks for testing")
	}

	callbackID := callbacks[0].ID
	responses, err := client.GetResponsesByCallback(ctx, callbackID, 10)
	if err != nil {
		t.Fatalf("GetResponsesByCallback failed: %v", err)
	}

	if responses == nil {
		t.Fatal("GetResponsesByCallback returned nil")
	}

	t.Logf("Found %d response(s) for callback %d", len(responses), callbackID)

	// Verify responses are ordered by timestamp (descending)
	if len(responses) > 1 {
		for i := 1; i < len(responses); i++ {
			if responses[i].Timestamp.After(responses[i-1].Timestamp) {
				t.Error("Responses should be sorted by timestamp (descending)")
				break
			}
		}
	}
}

func TestResponses_SearchResponses(t *testing.T) {
	SkipIfNoMythic(t)
	client := AuthenticateTestClient(t)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Search for common output patterns
	testQueries := []string{"success", "error", "completed", "task"}

	for _, query := range testQueries {
		responses, err := client.SearchResponses(ctx, query, &types.ResponseFilters{
			Limit: 5,
		})
		if err != nil {
			t.Errorf("SearchResponses failed for query '%s': %v", query, err)
			continue
		}

		t.Logf("Search '%s': found %d result(s)", query, len(responses))

		// Verify search results contain the query
		for _, resp := range responses {
			if !contains(resp.ResponseText, query) {
				t.Errorf("Response %d doesn't contain search term '%s'", resp.ID, query)
			}
		}
	}
}

func TestResponses_GetLatestResponses(t *testing.T) {
	SkipIfNoMythic(t)
	client := AuthenticateTestClient(t)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	limit := 10
	responses, err := client.GetLatestResponses(ctx, 0, limit)
	if err != nil {
		t.Fatalf("GetLatestResponses failed: %v", err)
	}

	if responses == nil {
		t.Fatal("GetLatestResponses returned nil")
	}

	t.Logf("Found %d latest response(s)", len(responses))

	// Verify limit is respected
	if len(responses) > limit {
		t.Errorf("Expected at most %d responses, got %d", limit, len(responses))
	}

	// Verify timestamp ordering (descending)
	if len(responses) > 1 {
		for i := 1; i < len(responses); i++ {
			if responses[i].Timestamp.After(responses[i-1].Timestamp) {
				t.Error("Responses should be sorted by timestamp (descending)")
				break
			}
		}
	}

	// Log details of latest response
	if len(responses) > 0 {
		latest := responses[0]
		t.Logf("Latest response: ID %d at %s",
			latest.ID, latest.Timestamp.Format("2006-01-02 15:04:05"))
	}
}

func TestResponses_GetResponseStatistics(t *testing.T) {
	SkipIfNoMythic(t)
	client := AuthenticateTestClient(t)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Get tasks first
	tasks, err := client.GetAllTasks(ctx, 10)
	if err != nil {
		t.Fatalf("GetAllTasks failed: %v", err)
	}
	if len(tasks) == 0 {
		t.Skip("No tasks available for testing")
	}

	taskID := tasks[0].ID
	stats, err := client.GetResponseStatistics(ctx, taskID)
	if err != nil {
		t.Fatalf("GetResponseStatistics failed: %v", err)
	}

	if stats == nil {
		t.Fatal("GetResponseStatistics returned nil")
	}

	t.Logf("Response statistics for task %d:", taskID)
	t.Logf("  - Total responses: %d", stats.TotalResponses)
	t.Logf("  - Total output size: %d bytes", stats.TotalOutputSize)

	if stats.TotalResponses > 0 {
		t.Logf("  - Average response size: %d bytes",
			stats.TotalOutputSize/int64(stats.TotalResponses))
		t.Logf("  - First response: %s", stats.FirstResponseTime.Format("2006-01-02 15:04:05"))
		t.Logf("  - Last response: %s", stats.LastResponseTime.Format("2006-01-02 15:04:05"))
	}

	// Verify consistency
	if stats.TotalResponses < 0 {
		t.Error("Total responses should not be negative")
	}
	if stats.TotalOutputSize < 0 {
		t.Error("Total output size should not be negative")
	}
}

func TestResponses_InvalidInputs(t *testing.T) {
	SkipIfNoMythic(t)
	client := AuthenticateTestClient(t)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Test zero task ID
	_, err := client.GetResponsesByTask(ctx, 0)
	if err == nil {
		t.Error("Expected error for zero task ID")
	}

	// Test zero response ID
	_, err = client.GetResponseByID(ctx, 0)
	if err == nil {
		t.Error("Expected error for zero response ID")
	}

	// Test zero callback ID
	_, err = client.GetResponsesByCallback(ctx, 0, 10)
	if err == nil {
		t.Error("Expected error for zero callback ID")
	}

	// Test empty search query
	_, err = client.SearchResponses(ctx, "", nil)
	if err == nil {
		t.Error("Expected error for empty search query")
	}

	t.Log("All invalid input tests passed")
}
