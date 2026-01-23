//go:build integration

package integration

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/nbaertsch/mythic-sdk-go/pkg/mythic"
	"github.com/nbaertsch/mythic-sdk-go/pkg/mythic/types"
)

// TestE2E_ResponseRetrieval tests comprehensive response retrieval operations.
// Covers: GetResponsesByTask, GetResponseByID, GetResponsesByCallback, GetLatestResponses
func TestE2E_ResponseRetrieval(t *testing.T) {
	// Ensure at least one callback exists (reuses existing or creates one)
	callbackID := EnsureCallbackExists(t)

	client := AuthenticateTestClient(t)

	// Get the callback details
	ctx0, cancel0 := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel0()
	testCallback, err := client.GetCallbackByID(ctx0, callbackID)
	if err != nil {
		t.Fatalf("Failed to get callback: %v", err)
	}

	// Issue a task to generate responses
	t.Log("=== Setup: Issue test task ===")
	ctx01, cancel01 := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel01()

	taskReq := &mythic.TaskRequest{
		Command:    "whoami",
		Params:     "whoami",
		CallbackID: &testCallback.ID,
	}

	task, err := client.IssueTask(ctx01, taskReq)
	if err != nil {
		t.Fatalf("IssueTask failed: %v", err)
	}
	t.Logf("Issued task %d for response tests", task.DisplayID)

	// Wait a bit for task to process
	time.Sleep(2 * time.Second)

	// Test 1: Get task output
	t.Log("=== Test 1: Get task output ===")
	ctx1, cancel1 := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel1()

	output, err := client.GetTaskOutput(ctx1, task.DisplayID)
	if err != nil {
		t.Fatalf("GetTaskOutput failed: %v", err)
	}
	t.Logf("✓ Retrieved %d output entries for task %d", len(output), task.DisplayID)

	// Show sample output
	for i, out := range output {
		if i < 3 {
			preview := out.ResponseText
			if len(preview) > 100 {
				preview = preview[:100] + "..."
			}
			t.Logf("  [%d] %s", i+1, preview)
		}
	}

	// Test 2: Get responses by task
	t.Log("=== Test 2: Get responses by task ===")
	ctx2, cancel2 := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel2()

	responses, err := client.GetResponsesByTask(ctx2, task.ID)
	if err != nil {
		t.Fatalf("GetResponsesByTask failed: %v", err)
	}
	t.Logf("✓ Retrieved %d responses for task %d", len(responses), task.ID)

	// Validate response structure
	for _, resp := range responses {
		if resp.ID == 0 {
			t.Error("Response has ID 0")
		}
		if resp.TaskID != task.ID {
			t.Errorf("Response has wrong TaskID: expected %d, got %d", task.ID, resp.TaskID)
		}
	}

	var testResponseID int
	if len(responses) > 0 {
		testResponseID = responses[0].ID
	}

	// Test 3: Get response by ID
	if testResponseID > 0 {
		t.Log("=== Test 3: Get response by ID ===")
		ctx3, cancel3 := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel3()

		response, err := client.GetResponseByID(ctx3, testResponseID)
		if err != nil {
			t.Fatalf("GetResponseByID failed: %v", err)
		}
		if response.ID != testResponseID {
			t.Errorf("Response ID mismatch: expected %d, got %d", testResponseID, response.ID)
		}
		t.Logf("✓ Response %d retrieved", response.ID)
		t.Logf("  Command: %s, Status: %s", response.TaskCommand, response.TaskStatus)
	}

	// Test 4: Get responses by callback
	t.Log("=== Test 4: Get responses by callback ===")
	ctx4, cancel4 := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel4()

	callbackResponses, err := client.GetResponsesByCallback(ctx4, testCallback.ID, 50)
	if err != nil {
		t.Fatalf("GetResponsesByCallback failed: %v", err)
	}
	t.Logf("✓ Retrieved %d responses for callback %d", len(callbackResponses), testCallback.ID)

	// Test 5: Get latest responses across operation
	t.Log("=== Test 5: Get latest responses across operation ===")
	ctx5, cancel5 := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel5()

	latestResponses, err := client.GetLatestResponses(ctx5, 0, 20)
	if err != nil {
		t.Fatalf("GetLatestResponses failed: %v", err)
	}
	t.Logf("✓ Retrieved %d latest responses from current operation", len(latestResponses))

	// Validate responses are ordered by time (most recent first)
	if len(latestResponses) > 1 {
		for i := 0; i < len(latestResponses)-1; i++ {
			if latestResponses[i].Timestamp.Before(latestResponses[i+1].Timestamp) {
				t.Error("Latest responses not ordered correctly (should be desc)")
			}
		}
		t.Log("  ✓ Responses ordered by timestamp (most recent first)")
	}

	t.Log("=== ✓ All response retrieval tests passed ===")
}

// TestE2E_ResponseSearch tests response search functionality.
// Covers: SearchResponses with various filters
func TestE2E_ResponseSearch(t *testing.T) {
	// Ensure at least one callback exists (reuses existing or creates one)
	callbackID := EnsureCallbackExists(t)

	client := AuthenticateTestClient(t)

	// Get the callback details
	ctx0, cancel0 := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel0()
	testCallback, err := client.GetCallbackByID(ctx0, callbackID)
	if err != nil {
		t.Fatalf("Failed to get callback: %v", err)
	}

	// Issue a task with known output
	t.Log("=== Setup: Issue task with known output ===")
	ctx02, cancel0 := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel0()

	searchMarker := fmt.Sprintf("SEARCH_TEST_%d", time.Now().Unix())
	taskReq := &mythic.TaskRequest{
		Command:    "whoami",
		Params:     searchMarker,
		CallbackID: &testCallback.ID,
	}

	task, err := client.IssueTask(ctx02, taskReq)
	if err != nil {
		t.Fatalf("IssueTask failed: %v", err)
	}
	t.Logf("Issued task %d with marker: %s", task.DisplayID, searchMarker)

	// Wait for task to process
	time.Sleep(3 * time.Second)

	// Test 1: Search for specific text
	t.Log("=== Test 1: Search responses for text ===")
	ctx1, cancel1 := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel1()

	searchReq := &types.ResponseSearchRequest{
		Query: "whoami",
		Limit: 10,
	}

	results, err := client.SearchResponses(ctx1, searchReq)
	if err != nil {
		t.Fatalf("SearchResponses failed: %v", err)
	}
	t.Logf("✓ Found %d responses matching 'whoami'", len(results))

	// Test 2: Search within specific callback
	t.Log("=== Test 2: Search responses within callback ===")
	ctx2, cancel2 := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel2()

	callbackSearchReq := &types.ResponseSearchRequest{
		Query:      "test",
		CallbackID: &testCallback.ID,
		Limit:      10,
	}

	callbackResults, err := client.SearchResponses(ctx2, callbackSearchReq)
	if err != nil {
		t.Fatalf("SearchResponses (callback) failed: %v", err)
	}
	t.Logf("✓ Found %d responses matching 'test' in callback %d", len(callbackResults), testCallback.ID)

	// Test 3: Search within specific task
	t.Log("=== Test 3: Search responses within task ===")
	ctx3, cancel3 := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel3()

	taskSearchReq := &types.ResponseSearchRequest{
		Query:  searchMarker,
		TaskID: &task.ID,
		Limit:  10,
	}

	taskResults, err := client.SearchResponses(ctx3, taskSearchReq)
	if err != nil {
		t.Fatalf("SearchResponses (task) failed: %v", err)
	}
	t.Logf("✓ Found %d responses matching marker in task %d", len(taskResults), task.ID)

	// Test 4: Search with time range
	t.Log("=== Test 4: Search with time range ===")
	ctx4, cancel4 := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel4()

	startTime := time.Now().Add(-1 * time.Hour)
	endTime := time.Now()

	timeSearchReq := &types.ResponseSearchRequest{
		Query:     "whoami",
		StartTime: &startTime,
		EndTime:   &endTime,
		Limit:     10,
	}

	timeResults, err := client.SearchResponses(ctx4, timeSearchReq)
	if err != nil {
		t.Fatalf("SearchResponses (time range) failed: %v", err)
	}
	t.Logf("✓ Found %d responses in last hour", len(timeResults))

	// Test 5: Search with pagination
	t.Log("=== Test 5: Search with pagination ===")
	ctx5, cancel5 := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel5()

	paginationReq := &types.ResponseSearchRequest{
		Query:  "whoami",
		Limit:  5,
		Offset: 0,
	}

	page1, err := client.SearchResponses(ctx5, paginationReq)
	if err != nil {
		t.Fatalf("SearchResponses (page 1) failed: %v", err)
	}
	t.Logf("✓ Page 1: %d responses", len(page1))

	// Get second page
	ctx6, cancel6 := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel6()

	paginationReq.Offset = 5
	page2, err := client.SearchResponses(ctx6, paginationReq)
	if err != nil {
		t.Fatalf("SearchResponses (page 2) failed: %v", err)
	}
	t.Logf("✓ Page 2: %d responses", len(page2))

	// Verify pages don't overlap
	if len(page1) > 0 && len(page2) > 0 {
		for _, r1 := range page1 {
			for _, r2 := range page2 {
				if r1.ID == r2.ID {
					t.Error("Pages contain duplicate responses")
				}
			}
		}
		t.Log("  ✓ Pages have no overlapping responses")
	}

	t.Log("=== ✓ All search tests passed ===")
}

// TestE2E_ResponseStatistics tests response statistics and aggregation.
// Covers: GetResponseStatistics
func TestE2E_ResponseStatistics(t *testing.T) {
	// Ensure at least one callback exists (reuses existing or creates one)
	callbackID := EnsureCallbackExists(t)

	client := AuthenticateTestClient(t)

	// Get the callback details
	ctx0, cancel0 := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel0()
	testCallback, err := client.GetCallbackByID(ctx0, callbackID)
	if err != nil {
		t.Fatalf("Failed to get callback: %v", err)
	}

	// Issue a task
	t.Log("=== Setup: Issue task for statistics ===")
	ctx1, cancel1 := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel1()

	taskReq := &mythic.TaskRequest{
		Command:    "whoami",
		Params:     "whoami",
		CallbackID: &testCallback.ID,
	}

	task, err := client.IssueTask(ctx1, taskReq)
	if err != nil {
		t.Fatalf("IssueTask failed: %v", err)
	}

	// Wait for responses
	time.Sleep(3 * time.Second)

	// Test: Get response statistics
	t.Log("=== Test: Get response statistics ===")
	ctx2, cancel2 := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel2()

	stats, err := client.GetResponseStatistics(ctx2, task.ID)
	if err != nil {
		t.Fatalf("GetResponseStatistics failed: %v", err)
	}

	t.Logf("✓ Statistics for task %d:", task.ID)
	t.Logf("  Response count: %d", stats.ResponseCount)
	t.Logf("  Total size: %d bytes", stats.TotalSize)
	if stats.ResponseCount > 0 {
		t.Logf("  Average size: %d bytes", stats.TotalSize/stats.ResponseCount)
		if !stats.FirstResponse.IsZero() {
			t.Logf("  First response: %s", stats.FirstResponse.Format(time.RFC3339))
		}
		if !stats.LatestResponse.IsZero() {
			t.Logf("  Last response: %s", stats.LatestResponse.Format(time.RFC3339))
			if !stats.FirstResponse.IsZero() {
				duration := stats.LatestResponse.Sub(stats.FirstResponse)
				t.Logf("  Response span: %s", duration)
			}
		}
	}

	// Validate statistics
	if stats.TaskID != task.ID {
		t.Errorf("Statistics TaskID mismatch: expected %d, got %d", task.ID, stats.TaskID)
	}
	if stats.ResponseCount < 0 {
		t.Error("Response count is negative")
	}
	if stats.TotalSize < 0 {
		t.Error("Total size is negative")
	}

	t.Log("=== ✓ Statistics tests passed ===")
}

// TestE2E_ResponseErrorHandling tests error scenarios for response operations.
func TestE2E_ResponseErrorHandling(t *testing.T) {
	client := AuthenticateTestClient(t)

	// Test 1: Get responses for non-existent task
	t.Log("=== Test 1: Get responses for non-existent task ===")
	ctx1, cancel1 := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel1()

	_, err := client.GetResponsesByTask(ctx1, 999999)
	// This might not error - empty array is valid
	if err != nil {
		t.Logf("✓ Non-existent task rejected: %v", err)
	} else {
		t.Log("✓ Non-existent task returns empty array (valid behavior)")
	}

	// Test 2: Get response by invalid ID
	t.Log("=== Test 2: Get response by invalid ID ===")
	ctx2, cancel2 := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel2()

	_, err = client.GetResponseByID(ctx2, 999999)
	if err == nil {
		t.Error("Expected error for non-existent response ID")
	}
	t.Logf("✓ Non-existent response rejected: %v", err)

	// Test 3: Get responses for invalid callback
	t.Log("=== Test 3: Get responses for invalid callback ===")
	ctx3, cancel3 := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel3()

	_, err = client.GetResponsesByCallback(ctx3, 999999, 10)
	// This might not error - empty array is valid
	if err != nil {
		t.Logf("✓ Non-existent callback rejected: %v", err)
	} else {
		t.Log("✓ Non-existent callback returns empty array (valid behavior)")
	}

	// Test 4: Search with invalid request
	t.Log("=== Test 4: Search with nil request ===")
	ctx4, cancel4 := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel4()

	_, err = client.SearchResponses(ctx4, nil)
	if err == nil {
		t.Error("Expected error for nil search request")
	}
	t.Logf("✓ Nil search request rejected: %v", err)

	// Test 5: Get statistics for invalid task
	t.Log("=== Test 5: Get statistics for invalid task ===")
	ctx5, cancel5 := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel5()

	_, err = client.GetResponseStatistics(ctx5, 999999)
	if err == nil {
		t.Error("Expected error for non-existent task statistics")
	}
	t.Logf("✓ Non-existent task statistics rejected: %v", err)

	t.Log("=== ✓ All error handling tests passed ===")
}

// TestE2E_ResponseLargeOutput tests handling of large response outputs.
func TestE2E_ResponseLargeOutput(t *testing.T) {
	// Ensure at least one callback exists (reuses existing or creates one)
	callbackID := EnsureCallbackExists(t)

	client := AuthenticateTestClient(t)

	// Get the callback details
	ctx0, cancel0 := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel0()
	testCallback, err := client.GetCallbackByID(ctx0, callbackID)
	if err != nil {
		t.Fatalf("Failed to get callback: %v", err)
	}

	// Test: Get responses from callback with limit
	t.Log("=== Test: Handle large response sets with pagination ===")
	ctx1, cancel1 := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel1()

	// Get a large number of responses in batches
	batchSize := 50
	totalFetched := 0
	var allResponses []*types.Response

	for i := 0; i < 3; i++ { // Fetch 3 batches
		responses, err := client.GetResponsesByCallback(ctx1, testCallback.ID, batchSize)
		if err != nil {
			t.Fatalf("GetResponsesByCallback failed: %v", err)
		}

		totalFetched += len(responses)
		allResponses = append(allResponses, responses...)

		if len(responses) < batchSize {
			break // No more responses
		}
	}

	t.Logf("✓ Fetched %d responses in batches of %d", totalFetched, batchSize)

	// Analyze response sizes
	if len(allResponses) > 0 {
		var totalSize int
		var maxSize int
		var maxResp *types.Response

		for _, resp := range allResponses {
			size := len(resp.Response)
			totalSize += size
			if size > maxSize {
				maxSize = size
				maxResp = resp
			}
		}

		avgSize := totalSize / len(allResponses)
		t.Logf("  Average response size: %d bytes", avgSize)
		t.Logf("  Largest response: %d bytes (Task %d)", maxSize, maxResp.TaskID)

		if maxSize > 1024 {
			preview := maxResp.Response
			if len(preview) > 200 {
				preview = preview[:200] + "..."
			}
			t.Logf("  Preview: %s", preview)
		}
	}

	t.Log("=== ✓ Large output handling tests passed ===")
}

// TestE2E_ResponseOrdering tests response ordering and sequencing.
func TestE2E_ResponseOrdering(t *testing.T) {
	// Ensure at least one callback exists (reuses existing or creates one)
	callbackID := EnsureCallbackExists(t)

	client := AuthenticateTestClient(t)

	// Get the callback details
	ctx0, cancel0 := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel0()
	testCallback, err := client.GetCallbackByID(ctx0, callbackID)
	if err != nil {
		t.Fatalf("Failed to get callback: %v", err)
	}

	// Issue a task
	t.Log("=== Test: Verify response ordering ===")
	ctx1, cancel1 := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel1()

	taskReq := &mythic.TaskRequest{
		Command:    "whoami",
		Params:     "whoami",
		CallbackID: &testCallback.ID,
	}

	task, err := client.IssueTask(ctx1, taskReq)
	if err != nil {
		t.Fatalf("IssueTask failed: %v", err)
	}

	// Wait for responses
	time.Sleep(3 * time.Second)

	// Get responses
	ctx2, cancel2 := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel2()

	responses, err := client.GetResponsesByTask(ctx2, task.ID)
	if err != nil {
		t.Fatalf("GetResponsesByTask failed: %v", err)
	}

	if len(responses) == 0 {
		t.Log("⚠ No responses to test ordering")
		return
	}

	t.Logf("✓ Retrieved %d responses for ordering test", len(responses))

	// Verify responses are ordered by ID (chronological)
	for i := 0; i < len(responses)-1; i++ {
		if responses[i].ID > responses[i+1].ID {
			t.Error("Responses not ordered by ID (ascending)")
		}
		if responses[i].Timestamp.After(responses[i+1].Timestamp) {
			t.Error("Responses not ordered by timestamp (ascending)")
		}

		// Check sequence numbers if present
		if responses[i].SequenceNumber != nil && responses[i+1].SequenceNumber != nil {
			if *responses[i].SequenceNumber > *responses[i+1].SequenceNumber {
				t.Error("Responses not ordered by sequence number (ascending)")
			}
		}
	}

	t.Log("  ✓ Responses correctly ordered by ID, timestamp, and sequence")
	t.Log("=== ✓ Response ordering tests passed ===")
}

// TestE2E_ResponseContentAnalysis tests analysis of response content patterns.
func TestE2E_ResponseContentAnalysis(t *testing.T) {
	client := AuthenticateTestClient(t)

	t.Log("=== Test: Analyze response content patterns ===")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Get recent responses
	responses, err := client.GetLatestResponses(ctx, 0, 100)
	if err != nil {
		t.Fatalf("GetLatestResponses failed: %v", err)
	}

	if len(responses) == 0 {
		t.Skip("No responses found for content analysis")
	}

	t.Logf("✓ Analyzing %d responses", len(responses))

	// Analyze content patterns
	emptyCount := 0
	shortCount := 0  // < 100 bytes
	mediumCount := 0 // 100-1000 bytes
	largeCount := 0  // > 1000 bytes

	commandFreq := make(map[string]int)

	for _, resp := range responses {
		size := len(resp.Response)

		if size == 0 {
			emptyCount++
		} else if size < 100 {
			shortCount++
		} else if size < 1000 {
			mediumCount++
		} else {
			largeCount++
		}

		// Track command frequency
		if resp.TaskCommand != "" {
			commandFreq[resp.TaskCommand]++
		}
	}

	t.Logf("  Response size distribution:")
	t.Logf("    Empty: %d", emptyCount)
	t.Logf("    Short (< 100B): %d", shortCount)
	t.Logf("    Medium (100-1000B): %d", mediumCount)
	t.Logf("    Large (> 1000B): %d", largeCount)

	t.Logf("  Top commands:")
	type cmdCount struct {
		cmd   string
		count int
	}
	var topCommands []cmdCount
	for cmd, count := range commandFreq {
		topCommands = append(topCommands, cmdCount{cmd, count})
	}
	// Sort by count (simple bubble sort for small data)
	for i := 0; i < len(topCommands); i++ {
		for j := i + 1; j < len(topCommands); j++ {
			if topCommands[j].count > topCommands[i].count {
				topCommands[i], topCommands[j] = topCommands[j], topCommands[i]
			}
		}
	}
	showCount := 5
	if len(topCommands) < showCount {
		showCount = len(topCommands)
	}
	for i := 0; i < showCount; i++ {
		t.Logf("    %s: %d", topCommands[i].cmd, topCommands[i].count)
	}

	// Look for common patterns
	errorPatterns := []string{"error", "failed", "exception", "denied"}
	successPatterns := []string{"success", "completed", "done"}

	errorMatches := 0
	successMatches := 0

	for _, resp := range responses {
		lowerResp := strings.ToLower(resp.Response)
		for _, pattern := range errorPatterns {
			if strings.Contains(lowerResp, pattern) {
				errorMatches++
				break
			}
		}
		for _, pattern := range successPatterns {
			if strings.Contains(lowerResp, pattern) {
				successMatches++
				break
			}
		}
	}

	t.Logf("  Content patterns:")
	t.Logf("    Responses with error indicators: %d", errorMatches)
	t.Logf("    Responses with success indicators: %d", successMatches)

	t.Log("=== ✓ Content analysis complete ===")
}
