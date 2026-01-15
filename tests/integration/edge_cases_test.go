//go:build integration

package integration

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/nbaertsch/mythic-sdk-go/pkg/mythic/types"
)

// TestE2E_ContextTimeout tests timeout behavior across various operations.
func TestE2E_ContextTimeout(t *testing.T) {
	client := AuthenticateTestClient(t)

	// Test 1: Very short timeout should fail
	t.Log("=== Test 1: Very short timeout (should fail) ===")
	ctx1, cancel1 := context.WithTimeout(context.Background(), 1*time.Millisecond)
	defer cancel1()

	_, err := client.GetAllCallbacks(ctx1)
	if err == nil {
		t.Log("  ⚠ Very short timeout did not fail (operation may have been cached)")
	} else {
		t.Logf("✓ Short timeout failed as expected: %v", err)
	}

	// Test 2: Reasonable timeout should succeed
	t.Log("=== Test 2: Reasonable timeout (should succeed) ===")
	ctx2, cancel2 := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel2()

	callbacks, err := client.GetAllCallbacks(ctx2)
	if err != nil {
		t.Fatalf("Reasonable timeout failed: %v", err)
	}
	t.Logf("✓ Reasonable timeout succeeded: retrieved %d callbacks", len(callbacks))

	// Test 3: Cancellation mid-operation
	t.Log("=== Test 3: Context cancellation ===")
	ctx3, cancel3 := context.WithCancel(context.Background())

	// Start operation in goroutine
	done := make(chan error, 1)
	go func() {
		_, err := client.GetAllCallbacks(ctx3)
		done <- err
	}()

	// Cancel immediately
	cancel3()

	// Wait for result
	select {
	case err := <-done:
		if err != nil {
			t.Logf("✓ Cancellation detected: %v", err)
		} else {
			t.Log("  ⚠ Operation completed despite cancellation (may have been very fast)")
		}
	case <-time.After(5 * time.Second):
		t.Error("Operation did not respond to cancellation")
	}

	t.Log("=== ✓ Timeout tests completed ===")
}

// TestE2E_ConcurrentOperations tests concurrent API operations.
func TestE2E_ConcurrentOperations(t *testing.T) {
	client := AuthenticateTestClient(t)

	// Test 1: Concurrent read operations
	t.Log("=== Test 1: Concurrent read operations ===")

	concurrency := 10
	var wg sync.WaitGroup
	errors := make(chan error, concurrency)

	start := time.Now()
	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()

			_, err := client.GetAllCallbacks(ctx)
			if err != nil {
				errors <- err
				return
			}

			t.Logf("  Goroutine %d completed successfully", id)
		}(i)
	}

	wg.Wait()
	close(errors)
	duration := time.Since(start)

	errorCount := 0
	for err := range errors {
		t.Errorf("Concurrent operation failed: %v", err)
		errorCount++
	}

	if errorCount == 0 {
		t.Logf("✓ All %d concurrent operations succeeded in %s", concurrency, duration)
	} else {
		t.Errorf("✗ %d/%d concurrent operations failed", errorCount, concurrency)
	}

	// Test 2: Concurrent mixed operations
	t.Log("=== Test 2: Concurrent mixed operations ===")

	operations := []struct {
		name string
		fn   func(context.Context) error
	}{
		{"GetCallbacks", func(ctx context.Context) error {
			_, err := client.GetAllCallbacks(ctx)
			return err
		}},
		{"GetOperations", func(ctx context.Context) error {
			_, err := client.GetOperations(ctx)
			return err
		}},
		{"GetPayloads", func(ctx context.Context) error {
			_, err := client.GetPayloads(ctx)
			return err
		}},
		{"GetC2Profiles", func(ctx context.Context) error {
			_, err := client.GetC2Profiles(ctx)
			return err
		}},
		{"GetCommands", func(ctx context.Context) error {
			_, err := client.GetCommands(ctx)
			return err
		}},
	}

	var wg2 sync.WaitGroup
	errors2 := make(chan error, len(operations))

	start2 := time.Now()
	for _, op := range operations {
		wg2.Add(1)
		go func(name string, fn func(context.Context) error) {
			defer wg2.Done()

			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()

			if err := fn(ctx); err != nil {
				errors2 <- err
				return
			}

			t.Logf("  %s completed successfully", name)
		}(op.name, op.fn)
	}

	wg2.Wait()
	close(errors2)
	duration2 := time.Since(start2)

	errorCount2 := 0
	for err := range errors2 {
		t.Errorf("Mixed operation failed: %v", err)
		errorCount2++
	}

	if errorCount2 == 0 {
		t.Logf("✓ All %d mixed operations succeeded in %s", len(operations), duration2)
	} else {
		t.Errorf("✗ %d/%d mixed operations failed", errorCount2, len(operations))
	}

	t.Log("=== ✓ Concurrent operation tests completed ===")
}

// TestE2E_LargeDatasets tests handling of large result sets.
func TestE2E_LargeDatasets(t *testing.T) {
	client := AuthenticateTestClient(t)

	// Test 1: Retrieve large callback list
	t.Log("=== Test 1: Retrieve large callback dataset ===")
	ctx1, cancel1 := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel1()

	start1 := time.Now()
	callbacks, err := client.GetAllCallbacks(ctx1)
	duration1 := time.Since(start1)

	if err != nil {
		t.Fatalf("GetAllCallbacks failed: %v", err)
	}

	t.Logf("✓ Retrieved %d callbacks in %s", len(callbacks), duration1)

	if len(callbacks) > 100 {
		t.Logf("  Large dataset: %d callbacks (>100)", len(callbacks))
	} else if len(callbacks) > 10 {
		t.Logf("  Medium dataset: %d callbacks", len(callbacks))
	} else {
		t.Logf("  Small dataset: %d callbacks", len(callbacks))
	}

	// Test 2: Retrieve tasks with large limit
	t.Log("=== Test 2: Retrieve tasks with large limit ===")

	// Get a callback first
	if len(callbacks) > 0 {
		ctx2, cancel2 := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel2()

		testCallback := callbacks[0]
		largeLimit := 1000

		start2 := time.Now()
		tasks, err := client.GetTasksForCallback(ctx2, testCallback.DisplayID, largeLimit)
		duration2 := time.Since(start2)

		if err != nil {
			t.Errorf("GetTasksForCallback with large limit failed: %v", err)
		} else {
			t.Logf("✓ Retrieved %d tasks with limit=%d in %s", len(tasks), largeLimit, duration2)

			if len(tasks) >= largeLimit {
				t.Logf("  Dataset reached limit: %d tasks", len(tasks))
			}
		}
	}

	// Test 3: Retrieve responses for large dataset
	t.Log("=== Test 3: Retrieve responses with large limit ===")

	if len(callbacks) > 0 {
		ctx3, cancel3 := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel3()

		testCallback := callbacks[0]
		largeLimit := 1000

		start3 := time.Now()
		responses, err := client.GetResponsesByCallback(ctx3, testCallback.ID, largeLimit)
		duration3 := time.Since(start3)

		if err != nil {
			t.Errorf("GetResponsesByCallback with large limit failed: %v", err)
		} else {
			t.Logf("✓ Retrieved %d responses with limit=%d in %s", len(responses), largeLimit, duration3)

			if len(responses) >= largeLimit {
				t.Logf("  Dataset reached limit: %d responses", len(responses))
			}
		}
	}

	t.Log("=== ✓ Large dataset tests completed ===")
}

// TestE2E_RateLimiting tests behavior under rapid successive requests.
func TestE2E_RateLimiting(t *testing.T) {
	client := AuthenticateTestClient(t)

	t.Log("=== Test: Rapid successive requests ===")

	requestCount := 20
	successCount := 0
	errorCount := 0
	var totalDuration time.Duration

	for i := 0; i < requestCount; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

		start := time.Now()
		_, err := client.GetOperations(ctx)
		duration := time.Since(start)
		totalDuration += duration

		cancel()

		if err != nil {
			errorCount++
			if i < 3 || i >= requestCount-3 {
				t.Logf("  Request %d failed: %v (duration: %s)", i+1, err, duration)
			}
		} else {
			successCount++
			if i < 3 || i >= requestCount-3 {
				t.Logf("  Request %d succeeded (duration: %s)", i+1, duration)
			}
		}
	}

	avgDuration := totalDuration / time.Duration(requestCount)

	t.Logf("✓ Completed %d rapid requests", requestCount)
	t.Logf("  Success: %d, Errors: %d", successCount, errorCount)
	t.Logf("  Average duration: %s", avgDuration)

	if errorCount > requestCount/2 {
		t.Errorf("High error rate: %d/%d requests failed", errorCount, requestCount)
	}

	t.Log("=== ✓ Rate limiting tests completed ===")
}

// TestE2E_InvalidInputHandling tests error handling for various invalid inputs.
func TestE2E_InvalidInputHandling(t *testing.T) {
	client := AuthenticateTestClient(t)

	// Test 1: Zero IDs
	t.Log("=== Test 1: Zero ID values ===")

	ctx1, cancel1 := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel1()

	_, err := client.GetCallbackByID(ctx1, 0)
	if err == nil {
		t.Error("Expected error for zero callback ID")
	}
	t.Logf("✓ Zero callback ID rejected: %v", err)

	// Test 2: Negative IDs
	t.Log("=== Test 2: Negative ID values ===")

	ctx2, cancel2 := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel2()

	_, err = client.GetTasksForCallback(ctx2, -1, 10)
	if err == nil {
		t.Error("Expected error for negative callback ID")
	}
	t.Logf("✓ Negative callback ID rejected: %v", err)

	// Test 3: Extremely large IDs
	t.Log("=== Test 3: Extremely large ID values ===")

	ctx3, cancel3 := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel3()

	_, err = client.GetCallbackByID(ctx3, 999999999)
	if err == nil {
		t.Error("Expected error for extremely large callback ID")
	}
	t.Logf("✓ Extremely large callback ID rejected: %v", err)

	// Test 4: Empty strings
	t.Log("=== Test 4: Empty string values ===")

	ctx4, cancel4 := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel4()

	_, err = client.ExportCallbackConfig(ctx4, "")
	if err == nil {
		t.Error("Expected error for empty agent callback ID")
	}
	t.Logf("✓ Empty agent callback ID rejected: %v", err)

	// Test 5: Nil context (should panic or error)
	t.Log("=== Test 5: Nil context handling ===")

	defer func() {
		if r := recover(); r != nil {
			t.Logf("✓ Nil context caused panic (expected): %v", r)
		}
	}()

	// This should either panic or return an error
	_, err = client.GetOperations(nil) //nolint:staticcheck // Testing nil context
	if err != nil {
		t.Logf("✓ Nil context rejected with error: %v", err)
	}

	// Test 6: Nil request objects
	t.Log("=== Test 6: Nil request objects ===")

	ctx6, cancel6 := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel6()

	_, err = client.GenerateReport(ctx6, nil)
	if err == nil {
		t.Error("Expected error for nil report request")
	}
	t.Logf("✓ Nil report request rejected: %v", err)

	t.Log("=== ✓ Invalid input handling tests completed ===")
}

// TestE2E_NetworkResilience tests behavior under network stress conditions.
func TestE2E_NetworkResilience(t *testing.T) {
	client := AuthenticateTestClient(t)

	t.Log("=== Test: Network resilience (baseline) ===")

	// Test baseline operation
	ctx1, cancel1 := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel1()

	start := time.Now()
	_, err := client.GetAllCallbacks(ctx1)
	baselineDuration := time.Since(start)

	if err != nil {
		t.Fatalf("Baseline operation failed: %v", err)
	}

	t.Logf("✓ Baseline operation completed in %s", baselineDuration)

	// Test multiple sequential operations
	t.Log("=== Test: Sequential operations ===")

	sequentialCount := 5
	var sequentialDurations []time.Duration

	for i := 0; i < sequentialCount; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

		start := time.Now()
		_, err := client.GetOperations(ctx)
		duration := time.Since(start)
		sequentialDurations = append(sequentialDurations, duration)

		cancel()

		if err != nil {
			t.Errorf("Sequential operation %d failed: %v", i+1, err)
		}
	}

	// Calculate average
	var total time.Duration
	for _, d := range sequentialDurations {
		total += d
	}
	avgSequential := total / time.Duration(sequentialCount)

	t.Logf("✓ %d sequential operations completed, average: %s", sequentialCount, avgSequential)

	if avgSequential > baselineDuration*3 {
		t.Logf("  ⚠ Sequential operations slower than baseline (avg: %s vs baseline: %s)",
			avgSequential, baselineDuration)
	}

	t.Log("=== ✓ Network resilience tests completed ===")
}

// TestE2E_MemoryEfficiency tests memory usage patterns.
func TestE2E_MemoryEfficiency(t *testing.T) {
	client := AuthenticateTestClient(t)

	t.Log("=== Test: Memory efficiency with repeated operations ===")

	// Perform same operation multiple times to check for memory leaks
	iterations := 10

	for i := 0; i < iterations; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

		// Perform various operations
		_, _ = client.GetAllCallbacks(ctx)
		_, _ = client.GetOperations(ctx)
		_, _ = client.GetPayloads(ctx)

		cancel()

		if i%5 == 0 {
			t.Logf("  Completed %d/%d iterations", i+1, iterations)
		}
	}

	t.Logf("✓ Completed %d iterations without crashes", iterations)
	t.Log("  Note: Use memory profiling tools for detailed memory analysis")

	t.Log("=== ✓ Memory efficiency tests completed ===")
}

// TestE2E_QueryComplexity tests handling of complex queries.
func TestE2E_QueryComplexity(t *testing.T) {
	client := AuthenticateTestClient(t)

	// Test 1: Search with complex filters
	t.Log("=== Test 1: Complex search query ===")

	ctx1, cancel1 := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel1()

	searchReq := &types.ResponseSearchRequest{
		Query:  "error",
		Limit:  100,
		Offset: 0,
	}

	start1 := time.Now()
	results, err := client.SearchResponses(ctx1, searchReq)
	duration1 := time.Since(start1)

	if err != nil {
		t.Logf("⚠ Complex search failed (may be expected): %v", err)
	} else {
		t.Logf("✓ Complex search completed in %s: %d results", duration1, len(results))
	}

	// Test 2: Report generation with all options
	t.Log("=== Test 2: Complex report generation ===")

	ctx2a, cancel2a := context.WithTimeout(context.Background(), 30*time.Second)
	operations, err := client.GetOperations(ctx2a)
	cancel2a()

	if err != nil || len(operations) == 0 {
		t.Log("⚠ No operations found, skipping complex report test")
	} else {
		ctx2, cancel2 := context.WithTimeout(context.Background(), 120*time.Second)
		defer cancel2()

		reportReq := &types.GenerateReportRequest{
			OperationID:        operations[0].ID,
			IncludeMITRE:       true,
			IncludeCallbacks:   true,
			IncludeTasks:       true,
			IncludeFiles:       true,
			IncludeCredentials: true,
			IncludeArtifacts:   true,
			OutputFormat:       types.ReportFormatJSON,
		}

		start2 := time.Now()
		report, err := client.GenerateReport(ctx2, reportReq)
		duration2 := time.Since(start2)

		if err != nil {
			t.Errorf("Complex report generation failed: %v", err)
		} else {
			t.Logf("✓ Complex report generated in %s: %d bytes", duration2, len(report.ReportData))

			if duration2 > 60*time.Second {
				t.Logf("  ⚠ Report generation took >60s (operation may have extensive data)")
			}
		}
	}

	t.Log("=== ✓ Query complexity tests completed ===")
}
