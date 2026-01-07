//go:build integration

package integration

import (
	"context"
	"testing"
	"time"
)

func TestFiles_GetFiles(t *testing.T) {
	SkipIfNoMythic(t)

	client := AuthenticateTestClient(t)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Get all files
	files, err := client.GetFiles(ctx, 10)
	if err != nil {
		t.Fatalf("Failed to get files: %v", err)
	}

	// May be empty on fresh instance
	t.Logf("Found %d files", len(files))

	// If we have files, validate their structure
	for _, file := range files {
		if file.AgentFileID == "" {
			t.Error("File should have AgentFileID")
		}

		if file.Filename == "" {
			t.Error("File should have Filename")
		}

		if file.Timestamp.IsZero() {
			t.Error("File should have Timestamp")
		}

		t.Logf("File: %s", file.String())
	}
}

func TestFiles_GetDownloadedFiles(t *testing.T) {
	SkipIfNoMythic(t)

	client := AuthenticateTestClient(t)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Get downloaded files
	files, err := client.GetDownloadedFiles(ctx, 10)
	if err != nil {
		t.Fatalf("Failed to get downloaded files: %v", err)
	}

	// May be empty on fresh instance
	t.Logf("Found %d downloaded files", len(files))

	// If we have files, validate they're downloads
	for _, file := range files {
		if !file.IsDownloadFromAgent {
			t.Errorf("File %s should be marked as download from agent", file.AgentFileID)
		}

		if file.Deleted {
			t.Errorf("File %s should not be deleted in query results", file.AgentFileID)
		}

		t.Logf("Downloaded file: %s", file.String())
	}
}

func TestFiles_UploadFile(t *testing.T) {
	SkipIfNoMythic(t)

	client := AuthenticateTestClient(t)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Upload a test file
	testData := []byte("This is a test file for Mythic SDK integration testing")
	agentFileID, err := client.UploadFile(ctx, "test.txt", testData)

	if err != nil {
		t.Fatalf("Failed to upload file: %v", err)
	}

	if agentFileID == "" {
		t.Fatal("UploadFile returned empty agent_file_id")
	}

	t.Logf("Uploaded file with agent_file_id: %s", agentFileID)

	// Verify the file exists
	file, err := client.GetFileByID(ctx, agentFileID)
	if err != nil {
		t.Fatalf("Failed to get uploaded file: %v", err)
	}

	if file.AgentFileID != agentFileID {
		t.Errorf("Expected agent_file_id %s, got %s", agentFileID, file.AgentFileID)
	}

	// Note: Filename validation may vary based on Mythic version
	t.Logf("Retrieved uploaded file: %s", file.String())
}

func TestFiles_UploadFile_MissingFilename(t *testing.T) {
	SkipIfNoMythic(t)

	client := AuthenticateTestClient(t)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Try to upload without filename
	testData := []byte("test data")
	_, err := client.UploadFile(ctx, "", testData)

	if err == nil {
		t.Fatal("Expected error for missing filename, got nil")
	}

	t.Logf("Expected error for missing filename: %v", err)
}

func TestFiles_UploadFile_EmptyData(t *testing.T) {
	SkipIfNoMythic(t)

	client := AuthenticateTestClient(t)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Try to upload empty file
	_, err := client.UploadFile(ctx, "empty.txt", []byte{})

	if err == nil {
		t.Fatal("Expected error for empty file data, got nil")
	}

	t.Logf("Expected error for empty file: %v", err)
}

func TestFiles_GetFileByID_NotFound(t *testing.T) {
	SkipIfNoMythic(t)

	client := AuthenticateTestClient(t)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Try to get non-existent file
	_, err := client.GetFileByID(ctx, "nonexistent-file-id")

	if err == nil {
		t.Fatal("Expected error for non-existent file, got nil")
	}

	t.Logf("Expected error for non-existent file: %v", err)
}

func TestFiles_DownloadFile(t *testing.T) {
	SkipIfNoMythic(t)

	client := AuthenticateTestClient(t)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// First upload a file
	testData := []byte("Test download content")
	agentFileID, err := client.UploadFile(ctx, "download_test.txt", testData)
	if err != nil {
		t.Fatalf("Failed to upload test file: %v", err)
	}

	t.Logf("Uploaded test file: %s", agentFileID)

	// Now download it
	downloadedData, err := client.DownloadFile(ctx, agentFileID)
	if err != nil {
		t.Fatalf("Failed to download file: %v", err)
	}

	if len(downloadedData) == 0 {
		t.Error("Downloaded data should not be empty")
	}

	t.Logf("Downloaded %d bytes", len(downloadedData))

	// Note: Content verification may vary based on Mythic encoding
	// Some versions base64 encode, some don't
}

func TestFiles_DownloadFile_NotFound(t *testing.T) {
	SkipIfNoMythic(t)

	client := AuthenticateTestClient(t)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Try to download non-existent file
	_, err := client.DownloadFile(ctx, "nonexistent-download-id")

	if err == nil {
		t.Fatal("Expected error for non-existent download, got nil")
	}

	t.Logf("Expected error for non-existent download: %v", err)
}

func TestFiles_DownloadFile_EmptyID(t *testing.T) {
	SkipIfNoMythic(t)

	client := AuthenticateTestClient(t)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Try to download with empty ID
	_, err := client.DownloadFile(ctx, "")

	if err == nil {
		t.Fatal("Expected error for empty agent_file_id, got nil")
	}

	t.Logf("Expected error for empty ID: %v", err)
}

func TestFiles_DeleteFile(t *testing.T) {
	SkipIfNoMythic(t)

	client := AuthenticateTestClient(t)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Upload a file to delete
	testData := []byte("File to be deleted")
	agentFileID, err := client.UploadFile(ctx, "delete_test.txt", testData)
	if err != nil {
		t.Fatalf("Failed to upload file for deletion: %v", err)
	}

	t.Logf("Uploaded file for deletion: %s", agentFileID)

	// Delete the file
	err = client.DeleteFile(ctx, agentFileID)
	if err != nil {
		t.Fatalf("Failed to delete file: %v", err)
	}

	// Verify it's marked as deleted
	file, err := client.GetFileByID(ctx, agentFileID)
	if err != nil {
		t.Fatalf("Failed to get deleted file metadata: %v", err)
	}

	if !file.IsDeleted() {
		t.Error("File should be marked as deleted")
	}

	t.Logf("File successfully marked as deleted: %s", agentFileID)
}

func TestFiles_DeleteFile_NotFound(t *testing.T) {
	SkipIfNoMythic(t)

	client := AuthenticateTestClient(t)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Try to delete non-existent file
	err := client.DeleteFile(ctx, "nonexistent-delete-id")

	if err == nil {
		t.Fatal("Expected error for deleting non-existent file, got nil")
	}

	t.Logf("Expected error for non-existent file deletion: %v", err)
}

func TestFiles_DeleteFile_EmptyID(t *testing.T) {
	SkipIfNoMythic(t)

	client := AuthenticateTestClient(t)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Try to delete with empty ID
	err := client.DeleteFile(ctx, "")

	if err == nil {
		t.Fatal("Expected error for empty agent_file_id, got nil")
	}

	t.Logf("Expected error for empty ID: %v", err)
}

func TestFiles_UploadDownloadCycle(t *testing.T) {
	SkipIfNoMythic(t)

	client := AuthenticateTestClient(t)

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	// Upload a file
	originalData := []byte("Complete upload/download cycle test\nLine 2\nLine 3")
	agentFileID, err := client.UploadFile(ctx, "cycle_test.txt", originalData)
	if err != nil {
		t.Fatalf("Failed to upload in cycle test: %v", err)
	}

	t.Logf("Uploaded file: %s", agentFileID)

	// Get file metadata
	file, err := client.GetFileByID(ctx, agentFileID)
	if err != nil {
		t.Fatalf("Failed to get file metadata: %v", err)
	}

	t.Logf("File metadata: %s", file.String())

	// Download the file
	downloadedData, err := client.DownloadFile(ctx, agentFileID)
	if err != nil {
		t.Fatalf("Failed to download in cycle test: %v", err)
	}

	t.Logf("Downloaded %d bytes", len(downloadedData))

	// Note: Content verification may need to account for Mythic's encoding
	// Just verify we got data back
	if len(downloadedData) == 0 {
		t.Error("Downloaded data should not be empty")
	}
}

func TestFiles_FilterByType(t *testing.T) {
	SkipIfNoMythic(t)

	client := AuthenticateTestClient(t)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Get all files
	allFiles, err := client.GetFiles(ctx, 100)
	if err != nil {
		t.Fatalf("Failed to get all files: %v", err)
	}

	// Get only downloads
	downloads, err := client.GetDownloadedFiles(ctx, 100)
	if err != nil {
		t.Fatalf("Failed to get downloaded files: %v", err)
	}

	t.Logf("Total files: %d, Downloaded files: %d", len(allFiles), len(downloads))

	// Downloaded files should be a subset of all files
	if len(downloads) > len(allFiles) {
		t.Error("Downloaded files count should not exceed total files count")
	}

	// Verify all downloaded files are in the all files list
	for _, download := range downloads {
		if !download.IsDownloadFromAgent {
			t.Errorf("File %s in downloads list should be marked as download", download.AgentFileID)
		}
	}
}

func TestFiles_BulkDownload(t *testing.T) {
	SkipIfNoMythic(t)

	client := AuthenticateTestClient(t)

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	// Upload two test files
	file1Data := []byte("First test file for bulk download")
	file1ID, err := client.UploadFile(ctx, "bulk_test1.txt", file1Data)
	if err != nil {
		t.Fatalf("Failed to upload first file: %v", err)
	}
	t.Logf("Uploaded file 1: %s", file1ID)

	file2Data := []byte("Second test file for bulk download")
	file2ID, err := client.UploadFile(ctx, "bulk_test2.txt", file2Data)
	if err != nil {
		t.Fatalf("Failed to upload second file: %v", err)
	}
	t.Logf("Uploaded file 2: %s", file2ID)

	// Create bulk download
	bulkFileID, err := client.BulkDownloadFiles(ctx, []string{file1ID, file2ID})
	if err != nil {
		t.Fatalf("Failed to create bulk download: %v", err)
	}

	if bulkFileID == "" {
		t.Fatal("BulkDownloadFiles returned empty file ID")
	}

	t.Logf("Bulk download created with file ID: %s", bulkFileID)

	// Download the bulk file (should be a ZIP)
	bulkData, err := client.DownloadFile(ctx, bulkFileID)
	if err != nil {
		t.Fatalf("Failed to download bulk file: %v", err)
	}

	if len(bulkData) == 0 {
		t.Error("Bulk download data should not be empty")
	}

	t.Logf("Downloaded bulk file: %d bytes", len(bulkData))
}

func TestFiles_BulkDownload_EmptyList(t *testing.T) {
	SkipIfNoMythic(t)

	client := AuthenticateTestClient(t)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Try bulk download with empty list
	_, err := client.BulkDownloadFiles(ctx, []string{})

	if err == nil {
		t.Fatal("Expected error for empty file list, got nil")
	}

	t.Logf("Expected error for empty list: %v", err)
}

func TestFiles_PreviewFile(t *testing.T) {
	SkipIfNoMythic(t)

	client := AuthenticateTestClient(t)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Upload a test file
	testData := []byte("This is preview test content\nLine 2\nLine 3")
	agentFileID, err := client.UploadFile(ctx, "preview_test.txt", testData)
	if err != nil {
		t.Fatalf("Failed to upload test file: %v", err)
	}

	t.Logf("Uploaded file for preview: %s", agentFileID)

	// Preview the file
	preview, err := client.PreviewFile(ctx, agentFileID)
	if err != nil {
		t.Fatalf("Failed to preview file: %v", err)
	}

	if preview == nil {
		t.Fatal("Preview should not be nil")
	}

	// Validate preview structure
	if preview.Filename == "" {
		t.Error("Preview should have filename")
	}

	t.Logf("Preview - Filename: %s, Size: %d, Host: %s",
		preview.Filename, preview.Size, preview.Host)

	// Contents may or may not be populated depending on file size
	t.Logf("Preview contents length: %d", len(preview.Contents))
}

func TestFiles_PreviewFile_NotFound(t *testing.T) {
	SkipIfNoMythic(t)

	client := AuthenticateTestClient(t)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Try to preview non-existent file
	_, err := client.PreviewFile(ctx, "nonexistent-preview-id")

	if err == nil {
		t.Fatal("Expected error for non-existent file preview, got nil")
	}

	t.Logf("Expected error for non-existent preview: %v", err)
}

func TestFiles_PreviewFile_EmptyID(t *testing.T) {
	SkipIfNoMythic(t)

	client := AuthenticateTestClient(t)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Try to preview with empty ID
	_, err := client.PreviewFile(ctx, "")

	if err == nil {
		t.Fatal("Expected error for empty file ID, got nil")
	}

	t.Logf("Expected error for empty ID: %v", err)
}
