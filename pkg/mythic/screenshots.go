package mythic

import (
	"context"
	"fmt"
	"time"
)

// GetScreenshots retrieves screenshots from a specific callback with optional filters.
//
// Screenshots are stored in the filemeta table with is_screenshot=true and require
// specialized handling for display and batch operations.
//
// Parameters:
//   - ctx: Context for cancellation and timeout
//   - callbackID: ID of the callback to retrieve screenshots from
//   - limit: Maximum number of screenshots to return (0 for default: 100)
//
// Returns:
//   - []*FileMeta: List of screenshot metadata (most recent first)
//   - error: Error if callback ID is invalid or query fails
//
// Example:
//
//	// Get last 20 screenshots from callback
//	screenshots, err := client.GetScreenshots(ctx, 5, 20)
//	if err != nil {
//	    return err
//	}
//	for _, screenshot := range screenshots {
//	    fmt.Printf("Screenshot: %s (taken at %s)\n",
//	        screenshot.Filename, screenshot.Timestamp)
//	}
func (c *Client) GetScreenshots(ctx context.Context, callbackID int, limit int) ([]*FileMeta, error) {
	if err := c.EnsureAuthenticated(ctx); err != nil {
		return nil, err
	}

	if callbackID == 0 {
		return nil, WrapError("GetScreenshots", ErrInvalidInput, "callback ID is required")
	}

	if limit <= 0 {
		limit = 100 // Default limit
	}

	var query struct {
		FileMeta []struct {
			ID                  int       `graphql:"id"`
			AgentFileID         string    `graphql:"agent_file_id"`
			TotalChunks         int       `graphql:"total_chunks"`
			ChunksReceived      int       `graphql:"chunks_received"`
			Complete            bool      `graphql:"complete"`
			FullRemotePath      string    `graphql:"full_remote_path"`
			Host                string    `graphql:"host"`
			IsPayload           bool      `graphql:"is_payload"`
			IsScreenshot        bool      `graphql:"is_screenshot"`
			IsDownloadFromAgent bool      `graphql:"is_download_from_agent"`
			Filename            string    `graphql:"filename_text"`
			MD5                 string    `graphql:"md5"`
			SHA1                string    `graphql:"sha1"`
			Comment             string    `graphql:"comment"`
			OperatorID          int       `graphql:"operator_id"`
			Timestamp           time.Time `graphql:"timestamp"`
			Deleted             bool      `graphql:"deleted"`
			TaskID              *int      `graphql:"task_id"`
			CallbackID          *int      `graphql:"callback_id"`
		} `graphql:"filemeta(where: {is_screenshot: {_eq: true}, deleted: {_eq: false}}, order_by: {timestamp: desc})"`
	}

	variables := map[string]interface{}{}

	err := c.executeQuery(ctx, &query, variables)
	if err != nil {
		return nil, WrapError("GetScreenshots", err, "failed to query screenshots")
	}

	// Filter by callback_id in code (Mythic v3.4.20 doesn't support callback_id in filemeta where clause)
	var filtered []*FileMeta
	for _, file := range query.FileMeta {
		if file.CallbackID != nil && *file.CallbackID == callbackID {
			filtered = append(filtered, &FileMeta{
				ID:                  file.ID,
				AgentFileID:         file.AgentFileID,
				TotalChunks:         file.TotalChunks,
				ChunksReceived:      file.ChunksReceived,
				Complete:            file.Complete,
				FullRemotePath:      file.FullRemotePath,
				Host:                file.Host,
				IsPayload:           file.IsPayload,
				IsScreenshot:        file.IsScreenshot,
				IsDownloadFromAgent: file.IsDownloadFromAgent,
				Filename:            decodeFilename(file.Filename),
				MD5:                 file.MD5,
				SHA1:                file.SHA1,
				Comment:             file.Comment,
				OperatorID:          file.OperatorID,
				Timestamp:           file.Timestamp,
				Deleted:             file.Deleted,
				TaskID:              file.TaskID,
			})

			// Apply limit
			if len(filtered) >= limit {
				break
			}
		}
	}

	return filtered, nil
}

// GetScreenshotByID retrieves a specific screenshot's metadata by its database ID.
//
// Parameters:
//   - ctx: Context for cancellation and timeout
//   - screenshotID: Database ID of the screenshot (filemeta.id)
//
// Returns:
//   - *FileMeta: Screenshot metadata
//   - error: Error if screenshot ID is invalid or not found
//
// Example:
//
//	screenshot, err := client.GetScreenshotByID(ctx, 42)
//	if err != nil {
//	    return err
//	}
//	fmt.Printf("Screenshot: %s (%d bytes)\n", screenshot.Filename, screenshot.Size)
func (c *Client) GetScreenshotByID(ctx context.Context, screenshotID int) (*FileMeta, error) {
	if err := c.EnsureAuthenticated(ctx); err != nil {
		return nil, err
	}

	if screenshotID == 0 {
		return nil, WrapError("GetScreenshotByID", ErrInvalidInput, "screenshot ID is required")
	}

	var query struct {
		FileMeta []struct {
			ID                  int       `graphql:"id"`
			AgentFileID         string    `graphql:"agent_file_id"`
			TotalChunks         int       `graphql:"total_chunks"`
			ChunksReceived      int       `graphql:"chunks_received"`
			Complete            bool      `graphql:"complete"`
			FullRemotePath      string    `graphql:"full_remote_path"`
			Host                string    `graphql:"host"`
			IsPayload           bool      `graphql:"is_payload"`
			IsScreenshot        bool      `graphql:"is_screenshot"`
			IsDownloadFromAgent bool      `graphql:"is_download_from_agent"`
			Filename            string    `graphql:"filename_text"`
			MD5                 string    `graphql:"md5"`
			SHA1                string    `graphql:"sha1"`
			Comment             string    `graphql:"comment"`
			OperatorID          int       `graphql:"operator_id"`
			Timestamp           time.Time `graphql:"timestamp"`
			Deleted             bool      `graphql:"deleted"`
			TaskID              *int      `graphql:"task_id"`
		} `graphql:"filemeta(where: {id: {_eq: $id}, is_screenshot: {_eq: true}})"`
	}

	variables := map[string]interface{}{
		"id": screenshotID,
	}

	err := c.executeQuery(ctx, &query, variables)
	if err != nil {
		return nil, WrapError("GetScreenshotByID", err, "failed to query screenshot")
	}

	if len(query.FileMeta) == 0 {
		return nil, WrapError("GetScreenshotByID", ErrNotFound, fmt.Sprintf("screenshot %d not found", screenshotID))
	}

	file := query.FileMeta[0]
	return &FileMeta{
		ID:                  file.ID,
		AgentFileID:         file.AgentFileID,
		TotalChunks:         file.TotalChunks,
		ChunksReceived:      file.ChunksReceived,
		Complete:            file.Complete,
		FullRemotePath:      file.FullRemotePath,
		Host:                file.Host,
		IsPayload:           file.IsPayload,
		IsScreenshot:        file.IsScreenshot,
		IsDownloadFromAgent: file.IsDownloadFromAgent,
		Filename:            decodeFilename(file.Filename),
		MD5:                 file.MD5,
		SHA1:                file.SHA1,
		Comment:             file.Comment,
		OperatorID:          file.OperatorID,
		Timestamp:           file.Timestamp,
		Deleted:             file.Deleted,
		TaskID:              file.TaskID,
	}, nil
}

// DownloadScreenshot downloads a screenshot file by its agent_file_id.
//
// This is a convenience wrapper around DownloadFile() that ensures the
// file is actually a screenshot.
//
// Parameters:
//   - ctx: Context for cancellation and timeout
//   - agentFileID: The agent_file_id of the screenshot
//
// Returns:
//   - []byte: Screenshot file data
//   - error: Error if screenshot not found or download fails
//
// Example:
//
//	data, err := client.DownloadScreenshot(ctx, "abc123-screenshot")
//	if err != nil {
//	    return err
//	}
//	err = os.WriteFile("screenshot.png", data, 0644)
func (c *Client) DownloadScreenshot(ctx context.Context, agentFileID string) ([]byte, error) {
	if err := c.EnsureAuthenticated(ctx); err != nil {
		return nil, err
	}

	if agentFileID == "" {
		return nil, WrapError("DownloadScreenshot", ErrInvalidInput, "agent_file_id is required")
	}

	// First verify it's actually a screenshot
	file, err := c.GetFileByID(ctx, agentFileID)
	if err != nil {
		return nil, WrapError("DownloadScreenshot", err, "failed to get file metadata")
	}

	if !file.IsScreenshot {
		return nil, WrapError("DownloadScreenshot", ErrInvalidInput, fmt.Sprintf("file %s is not a screenshot", agentFileID))
	}

	// Download using existing file download method
	data, err := c.DownloadFile(ctx, agentFileID)
	if err != nil {
		return nil, WrapError("DownloadScreenshot", err, "failed to download screenshot")
	}

	return data, nil
}

// GetScreenshotThumbnail retrieves a thumbnail version of a screenshot.
//
// Note: This is currently a placeholder as Mythic's thumbnail generation
// capabilities may vary by deployment. This method downloads the full
// screenshot. Future implementations may add server-side thumbnail support.
//
// Parameters:
//   - ctx: Context for cancellation and timeout
//   - agentFileID: The agent_file_id of the screenshot
//
// Returns:
//   - []byte: Screenshot data (currently full-size, may be thumbnail in future)
//   - error: Error if screenshot not found or download fails
//
// Example:
//
//	thumbnail, err := client.GetScreenshotThumbnail(ctx, "abc123-screenshot")
//	if err != nil {
//	    return err
//	}
//	// Display thumbnail in UI
func (c *Client) GetScreenshotThumbnail(ctx context.Context, agentFileID string) ([]byte, error) {
	// TODO: Implement server-side thumbnail generation if Mythic adds support
	// For now, return full screenshot
	return c.DownloadScreenshot(ctx, agentFileID)
}

// DeleteScreenshot marks a screenshot as deleted.
//
// This is a convenience wrapper around DeleteFile() that ensures the
// file is actually a screenshot before deletion.
//
// Parameters:
//   - ctx: Context for cancellation and timeout
//   - agentFileID: The agent_file_id of the screenshot to delete
//
// Returns:
//   - error: Error if screenshot not found or deletion fails
//
// Example:
//
//	err := client.DeleteScreenshot(ctx, "abc123-screenshot")
//	if err != nil {
//	    return err
//	}
//	fmt.Println("Screenshot deleted successfully")
func (c *Client) DeleteScreenshot(ctx context.Context, agentFileID string) error {
	if err := c.EnsureAuthenticated(ctx); err != nil {
		return err
	}

	if agentFileID == "" {
		return WrapError("DeleteScreenshot", ErrInvalidInput, "agent_file_id is required")
	}

	// First verify it's actually a screenshot
	file, err := c.GetFileByID(ctx, agentFileID)
	if err != nil {
		return WrapError("DeleteScreenshot", err, "failed to get file metadata")
	}

	if !file.IsScreenshot {
		return WrapError("DeleteScreenshot", ErrInvalidInput, fmt.Sprintf("file %s is not a screenshot", agentFileID))
	}

	// Delete using existing file deletion method
	err = c.DeleteFile(ctx, agentFileID)
	if err != nil {
		return WrapError("DeleteScreenshot", err, "failed to delete screenshot")
	}

	return nil
}

// GetScreenshotTimeline retrieves screenshots from a callback within a time range.
//
// This is useful for building timeline visualizations of user activity
// or identifying surveillance windows.
//
// Parameters:
//   - ctx: Context for cancellation and timeout
//   - callbackID: ID of the callback
//   - startTime: Start of time range (nil for no lower bound)
//   - endTime: End of time range (nil for no upper bound)
//
// Returns:
//   - []*FileMeta: Time-ordered list of screenshots (oldest first)
//   - error: Error if callback ID is invalid or query fails
//
// Example:
//
//	// Get screenshots from last 24 hours
//	endTime := time.Now()
//	startTime := endTime.Add(-24 * time.Hour)
//	screenshots, err := client.GetScreenshotTimeline(ctx, 5, &startTime, &endTime)
//	if err != nil {
//	    return err
//	}
//	fmt.Printf("Found %d screenshots in last 24 hours\n", len(screenshots))
func (c *Client) GetScreenshotTimeline(ctx context.Context, callbackID int, startTime *time.Time, endTime *time.Time) ([]*FileMeta, error) {
	if err := c.EnsureAuthenticated(ctx); err != nil {
		return nil, err
	}

	if callbackID == 0 {
		return nil, WrapError("GetScreenshotTimeline", ErrInvalidInput, "callback ID is required")
	}

	// Build query with time filters
	variables := map[string]interface{}{
		"callback_id": callbackID,
	}

	// GraphQL query structure depends on whether we have time filters
	if startTime != nil && endTime != nil {
		// Both start and end time
		var query struct {
			FileMeta []struct {
				ID                  int       `graphql:"id"`
				AgentFileID         string    `graphql:"agent_file_id"`
				TotalChunks         int       `graphql:"total_chunks"`
				ChunksReceived      int       `graphql:"chunks_received"`
				Complete            bool      `graphql:"complete"`
				FullRemotePath      string    `graphql:"full_remote_path"`
				Host                string    `graphql:"host"`
				IsPayload           bool      `graphql:"is_payload"`
				IsScreenshot        bool      `graphql:"is_screenshot"`
				IsDownloadFromAgent bool      `graphql:"is_download_from_agent"`
				Filename            string    `graphql:"filename_text"`
				MD5                 string    `graphql:"md5"`
				SHA1                string    `graphql:"sha1"`
				Comment             string    `graphql:"comment"`
				OperatorID          int       `graphql:"operator_id"`
				Timestamp           time.Time `graphql:"timestamp"`
				Deleted             bool      `graphql:"deleted"`
				TaskID              *int      `graphql:"task_id"`
			} `graphql:"filemeta(where: {callback_id: {_eq: $callback_id}, is_screenshot: {_eq: true}, deleted: {_eq: false}, timestamp: {_gte: $start_time, _lte: $end_time}}, order_by: {timestamp: asc})"`
		}

		variables["start_time"] = startTime.Format(time.RFC3339)
		variables["end_time"] = endTime.Format(time.RFC3339)

		err := c.executeQuery(ctx, &query, variables)
		if err != nil {
			return nil, WrapError("GetScreenshotTimeline", err, "failed to query screenshots")
		}

		screenshots := make([]*FileMeta, len(query.FileMeta))
		for i, file := range query.FileMeta {
			screenshots[i] = &FileMeta{
				ID:                  file.ID,
				AgentFileID:         file.AgentFileID,
				TotalChunks:         file.TotalChunks,
				ChunksReceived:      file.ChunksReceived,
				Complete:            file.Complete,
				FullRemotePath:      file.FullRemotePath,
				Host:                file.Host,
				IsPayload:           file.IsPayload,
				IsScreenshot:        file.IsScreenshot,
				IsDownloadFromAgent: file.IsDownloadFromAgent,
				Filename:            decodeFilename(file.Filename),
				MD5:                 file.MD5,
				SHA1:                file.SHA1,
				Comment:             file.Comment,
				OperatorID:          file.OperatorID,
				Timestamp:           file.Timestamp,
				Deleted:             file.Deleted,
				TaskID:              file.TaskID,
			}
		}

		return screenshots, nil
	} else if startTime != nil {
		// Only start time
		var query struct {
			FileMeta []struct {
				ID                  int       `graphql:"id"`
				AgentFileID         string    `graphql:"agent_file_id"`
				TotalChunks         int       `graphql:"total_chunks"`
				ChunksReceived      int       `graphql:"chunks_received"`
				Complete            bool      `graphql:"complete"`
				FullRemotePath      string    `graphql:"full_remote_path"`
				Host                string    `graphql:"host"`
				IsPayload           bool      `graphql:"is_payload"`
				IsScreenshot        bool      `graphql:"is_screenshot"`
				IsDownloadFromAgent bool      `graphql:"is_download_from_agent"`
				Filename            string    `graphql:"filename_text"`
				MD5                 string    `graphql:"md5"`
				SHA1                string    `graphql:"sha1"`
				Comment             string    `graphql:"comment"`
				OperatorID          int       `graphql:"operator_id"`
				Timestamp           time.Time `graphql:"timestamp"`
				Deleted             bool      `graphql:"deleted"`
				TaskID              *int      `graphql:"task_id"`
			} `graphql:"filemeta(where: {callback_id: {_eq: $callback_id}, is_screenshot: {_eq: true}, deleted: {_eq: false}, timestamp: {_gte: $start_time}}, order_by: {timestamp: asc})"`
		}

		variables["start_time"] = startTime.Format(time.RFC3339)

		err := c.executeQuery(ctx, &query, variables)
		if err != nil {
			return nil, WrapError("GetScreenshotTimeline", err, "failed to query screenshots")
		}

		screenshots := make([]*FileMeta, len(query.FileMeta))
		for i, file := range query.FileMeta {
			screenshots[i] = &FileMeta{
				ID:                  file.ID,
				AgentFileID:         file.AgentFileID,
				TotalChunks:         file.TotalChunks,
				ChunksReceived:      file.ChunksReceived,
				Complete:            file.Complete,
				FullRemotePath:      file.FullRemotePath,
				Host:                file.Host,
				IsPayload:           file.IsPayload,
				IsScreenshot:        file.IsScreenshot,
				IsDownloadFromAgent: file.IsDownloadFromAgent,
				Filename:            decodeFilename(file.Filename),
				MD5:                 file.MD5,
				SHA1:                file.SHA1,
				Comment:             file.Comment,
				OperatorID:          file.OperatorID,
				Timestamp:           file.Timestamp,
				Deleted:             file.Deleted,
				TaskID:              file.TaskID,
			}
		}

		return screenshots, nil
	} else if endTime != nil {
		// Only end time
		var query struct {
			FileMeta []struct {
				ID                  int       `graphql:"id"`
				AgentFileID         string    `graphql:"agent_file_id"`
				TotalChunks         int       `graphql:"total_chunks"`
				ChunksReceived      int       `graphql:"chunks_received"`
				Complete            bool      `graphql:"complete"`
				FullRemotePath      string    `graphql:"full_remote_path"`
				Host                string    `graphql:"host"`
				IsPayload           bool      `graphql:"is_payload"`
				IsScreenshot        bool      `graphql:"is_screenshot"`
				IsDownloadFromAgent bool      `graphql:"is_download_from_agent"`
				Filename            string    `graphql:"filename_text"`
				MD5                 string    `graphql:"md5"`
				SHA1                string    `graphql:"sha1"`
				Comment             string    `graphql:"comment"`
				OperatorID          int       `graphql:"operator_id"`
				Timestamp           time.Time `graphql:"timestamp"`
				Deleted             bool      `graphql:"deleted"`
				TaskID              *int      `graphql:"task_id"`
			} `graphql:"filemeta(where: {callback_id: {_eq: $callback_id}, is_screenshot: {_eq: true}, deleted: {_eq: false}, timestamp: {_lte: $end_time}}, order_by: {timestamp: asc})"`
		}

		variables["end_time"] = endTime.Format(time.RFC3339)

		err := c.executeQuery(ctx, &query, variables)
		if err != nil {
			return nil, WrapError("GetScreenshotTimeline", err, "failed to query screenshots")
		}

		screenshots := make([]*FileMeta, len(query.FileMeta))
		for i, file := range query.FileMeta {
			screenshots[i] = &FileMeta{
				ID:                  file.ID,
				AgentFileID:         file.AgentFileID,
				TotalChunks:         file.TotalChunks,
				ChunksReceived:      file.ChunksReceived,
				Complete:            file.Complete,
				FullRemotePath:      file.FullRemotePath,
				Host:                file.Host,
				IsPayload:           file.IsPayload,
				IsScreenshot:        file.IsScreenshot,
				IsDownloadFromAgent: file.IsDownloadFromAgent,
				Filename:            decodeFilename(file.Filename),
				MD5:                 file.MD5,
				SHA1:                file.SHA1,
				Comment:             file.Comment,
				OperatorID:          file.OperatorID,
				Timestamp:           file.Timestamp,
				Deleted:             file.Deleted,
				TaskID:              file.TaskID,
			}
		}

		return screenshots, nil
	} else {
		// No time filters - get all screenshots
		return c.GetScreenshots(ctx, callbackID, 1000) // Use high limit
	}
}
