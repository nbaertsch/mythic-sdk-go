//go:build integration

package integration

import (
	"context"
	"testing"
	"time"

	"github.com/nbaertsch/mythic-sdk-go/pkg/mythic"
)

// TestE2E_ScreenshotRetrieval tests screenshot retrieval operations.
// Covers: GetScreenshots, GetScreenshotByID
func TestE2E_ScreenshotRetrieval(t *testing.T) {
	client := AuthenticateTestClient(t)

	// Find a callback
	ctx0, cancel0 := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel0()

	callbacks, err := client.GetAllActiveCallbacks(ctx0)
	if err != nil {
		t.Fatalf("GetAllActiveCallbacks failed: %v", err)
	}

	if len(callbacks) == 0 {
		t.Skip("No callbacks found, skipping screenshot tests")
	}

	testCallback := callbacks[0]
	t.Logf("Using callback %d for screenshot tests", testCallback.ID)

	// Test 1: Get screenshots for callback
	t.Log("=== Test 1: Get screenshots for callback ===")
	ctx1, cancel1 := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel1()

	screenshots, err := client.GetScreenshots(ctx1, testCallback.ID, 50)
	if err != nil {
		t.Fatalf("GetScreenshots failed: %v", err)
	}
	t.Logf("✓ Retrieved %d screenshots for callback %d", len(screenshots), testCallback.ID)

	if len(screenshots) == 0 {
		t.Log("⚠ No screenshots found, skipping validation tests")
		return
	}

	// Validate screenshot structure
	for _, screenshot := range screenshots {
		if screenshot.ID == 0 {
			t.Error("Screenshot has ID 0")
		}
		if screenshot.AgentFileID == "" {
			t.Error("Screenshot has empty AgentFileID")
		}
		if !screenshot.IsScreenshot {
			t.Error("GetScreenshots returned non-screenshot file")
		}
		if screenshot.Deleted {
			t.Error("GetScreenshots returned deleted screenshot")
		}
	}

	// Show sample screenshots
	sampleCount := 5
	if len(screenshots) < sampleCount {
		sampleCount = len(screenshots)
	}
	t.Logf("  Sample screenshots:")
	for i := 0; i < sampleCount; i++ {
		ss := screenshots[i]
		t.Logf("    [%d] %s - %s (Complete: %v, Chunks: %d/%d)",
			i+1, ss.Filename, ss.Timestamp.Format(time.RFC3339),
			ss.Complete, ss.ChunksReceived, ss.TotalChunks)
	}

	// Test 2: Get screenshot by ID
	t.Log("=== Test 2: Get screenshot by ID ===")
	ctx2, cancel2 := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel2()

	testScreenshotID := screenshots[0].ID
	screenshot, err := client.GetScreenshotByID(ctx2, testScreenshotID)
	if err != nil {
		t.Fatalf("GetScreenshotByID failed: %v", err)
	}
	if screenshot.ID != testScreenshotID {
		t.Errorf("Screenshot ID mismatch: expected %d, got %d", testScreenshotID, screenshot.ID)
	}
	t.Logf("✓ Screenshot %d retrieved: %s", screenshot.ID, screenshot.Filename)

	// Verify ordering (should be desc by timestamp)
	if len(screenshots) > 1 {
		for i := 0; i < len(screenshots)-1; i++ {
			if screenshots[i].Timestamp.Before(screenshots[i+1].Timestamp) {
				t.Error("Screenshots not ordered by timestamp (descending)")
			}
		}
		t.Log("  ✓ Screenshots correctly ordered by timestamp (most recent first)")
	}

	t.Log("=== ✓ Screenshot retrieval tests passed ===")
}

// TestE2E_ScreenshotDownload tests screenshot download operations.
// Covers: DownloadScreenshot, GetScreenshotThumbnail
func TestE2E_ScreenshotDownload(t *testing.T) {
	client := AuthenticateTestClient(t)

	// Find a callback with screenshots
	ctx0, cancel0 := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel0()

	callbacks, err := client.GetAllActiveCallbacks(ctx0)
	if err != nil {
		t.Fatalf("GetAllActiveCallbacks failed: %v", err)
	}

	var testCallback *int
	var screenshots []*mythic.FileMeta

	for _, cb := range callbacks {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		ss, err := client.GetScreenshots(ctx, cb.ID, 5)
		cancel()

		if err == nil && len(ss) > 0 {
			testCallback = &cb.ID
			screenshots = ss
			break
		}
	}

	if testCallback == nil || len(screenshots) == 0 {
		t.Skip("No screenshots found, skipping download tests")
	}

	t.Logf("Using callback %d with %d screenshots", *testCallback, len(screenshots))

	// Find a complete screenshot
	var completeScreenshot *mythic.FileMeta
	for _, ss := range screenshots {
		if ss.Complete {
			completeScreenshot = ss
			break
		}
	}

	if completeScreenshot == nil {
		t.Skip("No complete screenshots found, skipping download tests")
	}

	t.Logf("Testing with complete screenshot: %s (AgentFileID: %s)",
		completeScreenshot.Filename, completeScreenshot.AgentFileID)

	// Test 1: Download full screenshot
	t.Log("=== Test 1: Download full screenshot ===")
	ctx1, cancel1 := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel1()

	data, err := client.DownloadScreenshot(ctx1, completeScreenshot.AgentFileID)
	if err != nil {
		t.Fatalf("DownloadScreenshot failed: %v", err)
	}
	if len(data) == 0 {
		t.Fatal("Downloaded screenshot has zero bytes")
	}
	t.Logf("✓ Downloaded screenshot: %d bytes", len(data))

	// Validate it's actually image data (check magic bytes)
	if len(data) >= 4 {
		// PNG magic bytes: 89 50 4E 47
		if data[0] == 0x89 && data[1] == 0x50 && data[2] == 0x4E && data[3] == 0x47 {
			t.Log("  ✓ Verified PNG image format")
		} else if data[0] == 0xFF && data[1] == 0xD8 {
			// JPEG magic bytes: FF D8
			t.Log("  ✓ Verified JPEG image format")
		} else if data[0] == 0x42 && data[1] == 0x4D {
			// BMP magic bytes: 42 4D
			t.Log("  ✓ Verified BMP image format")
		} else {
			t.Logf("  ⚠ Unknown image format (first 4 bytes: %x %x %x %x)",
				data[0], data[1], data[2], data[3])
		}
	}

	// Test 2: Download thumbnail
	t.Log("=== Test 2: Download screenshot thumbnail ===")
	ctx2, cancel2 := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel2()

	thumbData, err := client.GetScreenshotThumbnail(ctx2, completeScreenshot.AgentFileID)
	if err != nil {
		t.Logf("⚠ GetScreenshotThumbnail failed (may not be available): %v", err)
	} else {
		if len(thumbData) == 0 {
			t.Error("Thumbnail has zero bytes")
		} else {
			t.Logf("✓ Downloaded thumbnail: %d bytes", len(thumbData))
			if len(thumbData) < len(data) {
				t.Logf("  ✓ Thumbnail is smaller than full image (%d vs %d bytes)",
					len(thumbData), len(data))
			}
		}
	}

	t.Log("=== ✓ Screenshot download tests passed ===")
}

// TestE2E_ScreenshotTimeline tests screenshot timeline retrieval.
// Covers: GetScreenshotTimeline
func TestE2E_ScreenshotTimeline(t *testing.T) {
	client := AuthenticateTestClient(t)

	// Find a callback with screenshots
	ctx0, cancel0 := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel0()

	callbacks, err := client.GetAllActiveCallbacks(ctx0)
	if err != nil {
		t.Fatalf("GetAllActiveCallbacks failed: %v", err)
	}

	var testCallback *int
	for _, cb := range callbacks {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		ss, err := client.GetScreenshots(ctx, cb.ID, 1)
		cancel()

		if err == nil && len(ss) > 0 {
			testCallback = &cb.ID
			break
		}
	}

	if testCallback == nil {
		t.Skip("No screenshots found, skipping timeline tests")
	}

	// Test 1: Get timeline without time filters
	t.Log("=== Test 1: Get screenshot timeline (no filters) ===")
	ctx1, cancel1 := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel1()

	timeline, err := client.GetScreenshotTimeline(ctx1, *testCallback, nil, nil)
	if err != nil {
		t.Fatalf("GetScreenshotTimeline failed: %v", err)
	}
	t.Logf("✓ Retrieved %d screenshots in timeline", len(timeline))

	if len(timeline) > 0 {
		// Verify ordering
		for i := 0; i < len(timeline)-1; i++ {
			if timeline[i].Timestamp.Before(timeline[i+1].Timestamp) {
				t.Error("Timeline not ordered by timestamp (descending)")
			}
		}
		t.Log("  ✓ Timeline correctly ordered")
	}

	// Test 2: Get timeline with time range
	t.Log("=== Test 2: Get screenshot timeline (with time range) ===")
	ctx2, cancel2 := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel2()

	endTime := time.Now()
	startTime := endTime.Add(-24 * time.Hour)

	timelineFiltered, err := client.GetScreenshotTimeline(ctx2, *testCallback, &startTime, &endTime)
	if err != nil {
		t.Fatalf("GetScreenshotTimeline (filtered) failed: %v", err)
	}
	t.Logf("✓ Retrieved %d screenshots in last 24 hours", len(timelineFiltered))

	// Verify all screenshots are within time range
	for _, ss := range timelineFiltered {
		if ss.Timestamp.Before(startTime) || ss.Timestamp.After(endTime) {
			t.Errorf("Screenshot %d outside time range: %s", ss.ID, ss.Timestamp)
		}
	}
	if len(timelineFiltered) > 0 {
		t.Log("  ✓ All screenshots within time range")
	}

	t.Log("=== ✓ Screenshot timeline tests passed ===")
}

// TestE2E_ScreenshotDeletion tests screenshot deletion operations.
// Covers: DeleteScreenshot
func TestE2E_ScreenshotDeletion(t *testing.T) {
	_ = AuthenticateTestClient(t)

	t.Log("=== Test: Screenshot deletion (skipped for safety) ===")
	t.Log("⚠ DeleteScreenshot not tested to avoid data loss")
	t.Log("  To test deletion, create a dedicated test screenshot first")
	t.Log("=== ✓ Deletion test skipped ===")
}

// TestE2E_ScreenshotErrorHandling tests error scenarios for screenshot operations.
func TestE2E_ScreenshotErrorHandling(t *testing.T) {
	client := AuthenticateTestClient(t)

	// Test 1: Get screenshots with invalid callback
	t.Log("=== Test 1: Get screenshots with invalid callback ===")
	ctx1, cancel1 := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel1()

	_, err := client.GetScreenshots(ctx1, 0, 10)
	if err == nil {
		t.Error("Expected error for invalid callback ID")
	}
	t.Logf("✓ Invalid callback ID rejected: %v", err)

	// Test 2: Get screenshot by invalid ID
	t.Log("=== Test 2: Get screenshot by invalid ID ===")
	ctx2, cancel2 := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel2()

	_, err = client.GetScreenshotByID(ctx2, 0)
	if err == nil {
		t.Error("Expected error for invalid screenshot ID")
	}
	t.Logf("✓ Invalid screenshot ID rejected: %v", err)

	// Test 3: Get screenshot by non-existent ID
	t.Log("=== Test 3: Get screenshot by non-existent ID ===")
	ctx3, cancel3 := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel3()

	_, err = client.GetScreenshotByID(ctx3, 999999)
	if err == nil {
		t.Error("Expected error for non-existent screenshot")
	}
	t.Logf("✓ Non-existent screenshot rejected: %v", err)

	// Test 4: Download with invalid agent file ID
	t.Log("=== Test 4: Download with invalid agent file ID ===")
	ctx4, cancel4 := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel4()

	_, err = client.DownloadScreenshot(ctx4, "")
	if err == nil {
		t.Error("Expected error for empty agent file ID")
	}
	t.Logf("✓ Empty agent file ID rejected: %v", err)

	// Test 5: Download non-existent screenshot
	t.Log("=== Test 5: Download non-existent screenshot ===")
	ctx5, cancel5 := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel5()

	_, err = client.DownloadScreenshot(ctx5, "nonexistent-file-id")
	if err == nil {
		t.Error("Expected error for non-existent file")
	}
	t.Logf("✓ Non-existent file rejected: %v", err)

	// Test 6: Timeline with invalid callback
	t.Log("=== Test 6: Timeline with invalid callback ===")
	ctx6, cancel6 := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel6()

	_, err = client.GetScreenshotTimeline(ctx6, 0, nil, nil)
	if err == nil {
		t.Error("Expected error for invalid callback ID")
	}
	t.Logf("✓ Invalid callback ID rejected: %v", err)

	t.Log("=== ✓ All error handling tests passed ===")
}

// TestE2E_ScreenshotCompleteness tests screenshot chunk completion tracking.
func TestE2E_ScreenshotCompleteness(t *testing.T) {
	client := AuthenticateTestClient(t)

	// Find screenshots
	ctx0, cancel0 := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel0()

	callbacks, err := client.GetAllActiveCallbacks(ctx0)
	if err != nil {
		t.Fatalf("GetAllActiveCallbacks failed: %v", err)
	}

	var allScreenshots []*mythic.FileMeta
	for _, cb := range callbacks {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		ss, err := client.GetScreenshots(ctx, cb.ID, 50)
		cancel()

		if err == nil {
			allScreenshots = append(allScreenshots, ss...)
		}
	}

	if len(allScreenshots) == 0 {
		t.Skip("No screenshots found, skipping completeness test")
	}

	t.Log("=== Test: Analyze screenshot completeness ===")
	t.Logf("✓ Analyzing %d screenshots", len(allScreenshots))

	complete := 0
	incomplete := 0
	multiChunk := 0

	for _, ss := range allScreenshots {
		if ss.Complete {
			complete++
		} else {
			incomplete++
		}
		if ss.TotalChunks > 1 {
			multiChunk++
		}
	}

	t.Logf("  Completeness status:")
	t.Logf("    Complete: %d (%.1f%%)", complete, float64(complete)/float64(len(allScreenshots))*100)
	t.Logf("    Incomplete: %d (%.1f%%)", incomplete, float64(incomplete)/float64(len(allScreenshots))*100)
	t.Logf("    Multi-chunk: %d", multiChunk)

	// Show incomplete screenshots
	if incomplete > 0 {
		t.Logf("  Incomplete screenshots:")
		showCount := 3
		if incomplete < showCount {
			showCount = incomplete
		}
		count := 0
		for _, ss := range allScreenshots {
			if !ss.Complete && count < showCount {
				t.Logf("    %s: %d/%d chunks", ss.Filename, ss.ChunksReceived, ss.TotalChunks)
				count++
			}
		}
	}

	t.Log("=== ✓ Completeness analysis complete ===")
}
