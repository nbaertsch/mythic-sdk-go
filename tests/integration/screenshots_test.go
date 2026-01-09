//go:build integration

package integration

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/nbaertsch/mythic-sdk-go/pkg/mythic/types"
)

func TestScreenshots_GetScreenshots(t *testing.T) {
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
	screenshots, err := client.GetScreenshots(ctx, callbackID, &types.ScreenshotFilters{
		Limit: 10,
	})
	if err != nil {
		t.Fatalf("GetScreenshots failed: %v", err)
	}

	if screenshots == nil {
		t.Fatal("GetScreenshots returned nil")
	}

	t.Logf("Found %d screenshot(s) for callback %d", len(screenshots), callbackID)

	for _, screenshot := range screenshots {
		if screenshot.ID == 0 {
			t.Error("Screenshot ID should not be 0")
		}
		if !screenshot.IsScreenshot {
			t.Error("File should be marked as screenshot")
		}
		t.Logf("  - %s (%s) - Complete: %v",
			screenshot.Filename, screenshot.Host, screenshot.Complete)
	}
}

func TestScreenshots_GetScreenshotByID(t *testing.T) {
	SkipIfNoMythic(t)
	client := AuthenticateTestClient(t)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Get callbacks and their screenshots
	callbacks, err := client.GetAllActiveCallbacks(ctx)
	if err != nil {
		t.Fatalf("GetAllActiveCallbacks failed: %v", err)
	}
	if len(callbacks) == 0 {
		t.Skip("No active callbacks for testing")
	}

	screenshots, err := client.GetScreenshots(ctx, callbacks[0].ID, nil)
	if err != nil {
		t.Fatalf("GetScreenshots failed: %v", err)
	}
	if len(screenshots) == 0 {
		t.Skip("No screenshots available for testing")
	}

	screenshotID := screenshots[0].ID
	screenshot, err := client.GetScreenshotByID(ctx, screenshotID)
	if err != nil {
		t.Fatalf("GetScreenshotByID failed: %v", err)
	}

	if screenshot == nil {
		t.Fatal("GetScreenshotByID returned nil")
	}

	if screenshot.ID != screenshotID {
		t.Errorf("Expected screenshot ID %d, got %d", screenshotID, screenshot.ID)
	}

	t.Logf("Retrieved screenshot %d: %s", screenshotID, screenshot.Filename)
	t.Logf("  - Host: %s", screenshot.Host)
	t.Logf("  - Timestamp: %s", screenshot.Timestamp.Format("2006-01-02 15:04:05"))
	t.Logf("  - Complete: %v", screenshot.Complete)
}

func TestScreenshots_DownloadScreenshot(t *testing.T) {
	SkipIfNoMythic(t)
	client := AuthenticateTestClient(t)
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	// Get callbacks and screenshots
	callbacks, err := client.GetAllActiveCallbacks(ctx)
	if err != nil {
		t.Fatalf("GetAllActiveCallbacks failed: %v", err)
	}
	if len(callbacks) == 0 {
		t.Skip("No active callbacks for testing")
	}

	screenshots, err := client.GetScreenshots(ctx, callbacks[0].ID, &types.ScreenshotFilters{
		CompleteOnly: true,
		Limit:        1,
	})
	if err != nil {
		t.Fatalf("GetScreenshots failed: %v", err)
	}
	if len(screenshots) == 0 {
		t.Skip("No complete screenshots available for testing")
	}

	screenshot := screenshots[0]
	tempFile := "/tmp/test_screenshot_" + screenshot.AgentFileID + ".png"
	defer os.Remove(tempFile)

	err = client.DownloadScreenshot(ctx, screenshot.ID, tempFile)
	if err != nil {
		t.Fatalf("DownloadScreenshot failed: %v", err)
	}

	// Verify file exists
	info, err := os.Stat(tempFile)
	if err != nil {
		t.Fatalf("Downloaded file not found: %v", err)
	}

	t.Logf("Downloaded screenshot to %s (%d bytes)", tempFile, info.Size())

	if info.Size() == 0 {
		t.Error("Downloaded screenshot is empty")
	}
}

func TestScreenshots_GetScreenshotThumbnail(t *testing.T) {
	SkipIfNoMythic(t)
	client := AuthenticateTestClient(t)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Get complete screenshots
	callbacks, err := client.GetAllActiveCallbacks(ctx)
	if err != nil {
		t.Fatalf("GetAllActiveCallbacks failed: %v", err)
	}
	if len(callbacks) == 0 {
		t.Skip("No active callbacks for testing")
	}

	screenshots, err := client.GetScreenshots(ctx, callbacks[0].ID, &types.ScreenshotFilters{
		CompleteOnly: true,
		Limit:        1,
	})
	if err != nil {
		t.Fatalf("GetScreenshots failed: %v", err)
	}
	if len(screenshots) == 0 {
		t.Skip("No complete screenshots available for testing")
	}

	screenshotID := screenshots[0].ID
	thumbnail, err := client.GetScreenshotThumbnail(ctx, screenshotID)
	if err != nil {
		t.Fatalf("GetScreenshotThumbnail failed: %v", err)
	}

	if thumbnail == "" {
		t.Error("GetScreenshotThumbnail returned empty string")
	}

	t.Logf("Retrieved thumbnail for screenshot %d (base64 length: %d)",
		screenshotID, len(thumbnail))

	// Verify it's base64 encoded
	if len(thumbnail) < 100 {
		t.Error("Thumbnail seems too short to be valid")
	}
}

func TestScreenshots_DeleteScreenshot(t *testing.T) {
	t.Skip("Skipping delete test to preserve test data")

	SkipIfNoMythic(t)
	client := AuthenticateTestClient(t)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Get screenshots
	callbacks, err := client.GetAllActiveCallbacks(ctx)
	if err != nil {
		t.Fatalf("GetAllActiveCallbacks failed: %v", err)
	}
	if len(callbacks) == 0 {
		t.Skip("No active callbacks for testing")
	}

	screenshots, err := client.GetScreenshots(ctx, callbacks[0].ID, nil)
	if err != nil {
		t.Fatalf("GetScreenshots failed: %v", err)
	}
	if len(screenshots) == 0 {
		t.Skip("No screenshots available for testing")
	}

	screenshotID := screenshots[len(screenshots)-1].ID // Delete oldest
	err = client.DeleteScreenshot(ctx, screenshotID)
	if err != nil {
		t.Fatalf("DeleteScreenshot failed: %v", err)
	}

	t.Logf("Successfully deleted screenshot %d", screenshotID)

	// Verify deletion
	_, err = client.GetScreenshotByID(ctx, screenshotID)
	if err == nil {
		t.Error("Screenshot should not be retrievable after deletion")
	}
}

func TestScreenshots_GetScreenshotTimeline(t *testing.T) {
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

	// Get timeline for last 24 hours
	endTime := time.Now()
	startTime := endTime.Add(-24 * time.Hour)

	screenshots, err := client.GetScreenshotTimeline(ctx, callbackID, startTime, endTime)
	if err != nil {
		t.Fatalf("GetScreenshotTimeline failed: %v", err)
	}

	if screenshots == nil {
		t.Fatal("GetScreenshotTimeline returned nil")
	}

	t.Logf("Found %d screenshot(s) in timeline for callback %d", len(screenshots), callbackID)

	// Verify timeline ordering
	if len(screenshots) > 1 {
		for i := 1; i < len(screenshots); i++ {
			if screenshots[i].Timestamp.After(screenshots[i-1].Timestamp) {
				t.Error("Screenshots should be sorted by timestamp (descending)")
				break
			}
		}
	}

	// Verify time range
	for _, screenshot := range screenshots {
		if screenshot.Timestamp.Before(startTime) || screenshot.Timestamp.After(endTime) {
			t.Errorf("Screenshot %d timestamp %s outside range [%s, %s]",
				screenshot.ID,
				screenshot.Timestamp.Format("2006-01-02 15:04:05"),
				startTime.Format("2006-01-02 15:04:05"),
				endTime.Format("2006-01-02 15:04:05"))
		}
	}

	if len(screenshots) > 0 {
		t.Logf("Timeline range: %s to %s",
			screenshots[len(screenshots)-1].Timestamp.Format("2006-01-02 15:04:05"),
			screenshots[0].Timestamp.Format("2006-01-02 15:04:05"))
	}
}

func TestScreenshots_InvalidInputs(t *testing.T) {
	SkipIfNoMythic(t)
	client := AuthenticateTestClient(t)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Test zero callback ID
	_, err := client.GetScreenshots(ctx, 0, nil)
	if err == nil {
		t.Error("Expected error for zero callback ID")
	}

	// Test zero screenshot ID
	_, err = client.GetScreenshotByID(ctx, 0)
	if err == nil {
		t.Error("Expected error for zero screenshot ID")
	}

	// Test empty output path
	_, err = client.DownloadScreenshot(ctx, 1, "")
	if err == nil {
		t.Error("Expected error for empty output path")
	}

	// Test invalid time range
	endTime := time.Now().Add(-48 * time.Hour)
	startTime := time.Now()
	_, err = client.GetScreenshotTimeline(ctx, 1, startTime, endTime)
	if err == nil {
		t.Error("Expected error for invalid time range (start > end)")
	}

	t.Log("All invalid input tests passed")
}
