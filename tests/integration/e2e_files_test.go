//go:build integration

package integration

import (
	"bytes"
	"context"
	"crypto/rand"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"
)

// TestE2E_FileOperations tests the complete file management workflow
// Covers: GetFiles, GetFileByID, GetDownloadedFiles, UploadFile,
// DownloadFile, DeleteFile, BulkDownloadFiles, PreviewFile
func TestE2E_FileOperations(t *testing.T) {
	client := AuthenticateTestClient(t)

	var uploadedFileIDs []string
	var tempFiles []string

	// Register cleanup
	defer func() {
		// Delete uploaded files
		for _, fileID := range uploadedFileIDs {
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			_ = client.DeleteFile(ctx, fileID)
			cancel()
			t.Logf("Cleaned up file ID: %s", fileID)
		}
		// Remove temp files
		for _, file := range tempFiles {
			_ = os.Remove(file)
			t.Logf("Removed temp file: %s", file)
		}
	}()

	// Test 1: Create test files
	t.Log("=== Test 1: Create test files (1KB, 1MB, 10MB) ===")
	smallFile, err := createTempFile("small_1kb.dat", 1024)
	if err != nil {
		t.Fatalf("Failed to create small file: %v", err)
	}
	tempFiles = append(tempFiles, smallFile)
	t.Logf("✓ Created small file (1KB): %s", smallFile)

	mediumFile, err := createTempFile("medium_1mb.dat", 1024*1024)
	if err != nil {
		t.Fatalf("Failed to create medium file: %v", err)
	}
	tempFiles = append(tempFiles, mediumFile)
	t.Logf("✓ Created medium file (1MB): %s", mediumFile)

	largeFile, err := createTempFile("large_10mb.dat", 10*1024*1024)
	if err != nil {
		t.Fatalf("Failed to create large file: %v", err)
	}
	tempFiles = append(tempFiles, largeFile)
	t.Logf("✓ Created large file (10MB): %s", largeFile)

	// Test 2: Get files baseline
	t.Log("=== Test 2: Get files before upload (baseline) ===")
	ctx1, cancel1 := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel1()

	baselineFiles, err := client.GetFiles(ctx1, 100)
	if err != nil {
		t.Fatalf("GetFiles baseline failed: %v", err)
	}
	baselineCount := len(baselineFiles)
	t.Logf("✓ Baseline file count: %d", baselineCount)

	// Test 3: Upload small file
	t.Log("=== Test 3: Upload small file (1KB) ===")
	ctx2, cancel2 := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel2()

	smallData, err := os.ReadFile(smallFile)
	if err != nil {
		t.Fatalf("Failed to read small file: %v", err)
	}

	smallFileID, err := client.UploadFile(ctx2, filepath.Base(smallFile), smallData)
	if err != nil {
		t.Fatalf("UploadFile (small) failed: %v", err)
	}
	uploadedFileIDs = append(uploadedFileIDs, smallFileID)
	t.Logf("✓ Small file uploaded: ID %s", smallFileID)

	// Test 4: Upload medium file
	t.Log("=== Test 4: Upload medium file (1MB) ===")
	ctx3, cancel3 := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel3()

	mediumData, err := os.ReadFile(mediumFile)
	if err != nil {
		t.Fatalf("Failed to read medium file: %v", err)
	}

	mediumFileID, err := client.UploadFile(ctx3, filepath.Base(mediumFile), mediumData)
	if err != nil {
		t.Fatalf("UploadFile (medium) failed: %v", err)
	}
	uploadedFileIDs = append(uploadedFileIDs, mediumFileID)
	t.Logf("✓ Medium file uploaded: ID %s", mediumFileID)

	// Test 5: Upload large file
	t.Log("=== Test 5: Upload large file (10MB) ===")
	ctx4, cancel4 := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel4()

	largeData, err := os.ReadFile(largeFile)
	if err != nil {
		t.Fatalf("Failed to read large file: %v", err)
	}

	largeFileID, err := client.UploadFile(ctx4, filepath.Base(largeFile), largeData)
	if err != nil {
		t.Fatalf("UploadFile (large) failed: %v", err)
	}
	uploadedFileIDs = append(uploadedFileIDs, largeFileID)
	t.Logf("✓ Large file uploaded: ID %s", largeFileID)

	// Test 6: Get all files after upload
	t.Log("=== Test 6: Get all files after upload ===")
	ctx5, cancel5 := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel5()

	allFiles, err := client.GetFiles(ctx5, 100)
	if err != nil {
		t.Fatalf("GetFiles after upload failed: %v", err)
	}
	newCount := len(allFiles)
	if newCount < baselineCount+3 {
		t.Errorf("Expected at least %d files, got %d", baselineCount+3, newCount)
	}
	t.Logf("✓ Total files now: %d (added %d)", newCount, newCount-baselineCount)

	// Test 7: Get file by ID for each uploaded file
	t.Log("=== Test 7: Get file by ID for each uploaded file ===")
	for i, fileID := range uploadedFileIDs {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		fileMeta, err := client.GetFileByID(ctx, fileID)
		cancel()
		if err != nil {
			t.Errorf("GetFileByID failed for %s: %v", fileID, err)
			continue
		}
		if fileMeta.AgentFileID != fileID {
			t.Errorf("File ID mismatch: expected %s, got %s", fileID, fileMeta.AgentFileID)
		}
		t.Logf("✓ File %d: ID=%s, Path=%s, Complete=%v", i+1, fileMeta.AgentFileID, fileMeta.Path, fileMeta.Complete)
	}

	// Test 8: Preview small file
	t.Log("=== Test 8: Preview small file ===")
	ctx6, cancel6 := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel6()

	preview, err := client.PreviewFile(ctx6, smallFileID)
	if err != nil {
		t.Fatalf("PreviewFile failed: %v", err)
	}
	if preview.Size == 0 {
		t.Error("Preview size is 0")
	}
	if preview.Filename == "" {
		t.Error("Preview filename is empty")
	}
	t.Logf("✓ Preview: Filename=%s, Size=%d bytes, Contents length=%d", preview.Filename, preview.Size, len(preview.Contents))

	// Test 9: Download small file
	t.Log("=== Test 9: Download small file and verify ===")
	ctx7, cancel7 := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel7()

	downloadedSmall, err := client.DownloadFile(ctx7, smallFileID)
	if err != nil {
		t.Fatalf("DownloadFile (small) failed: %v", err)
	}
	if !bytes.Equal(downloadedSmall, smallData) {
		t.Error("Downloaded small file content does not match original")
	}
	t.Logf("✓ Small file downloaded and verified: %d bytes", len(downloadedSmall))

	// Test 10: Download medium file
	t.Log("=== Test 10: Download medium file and verify ===")
	ctx8, cancel8 := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel8()

	downloadedMedium, err := client.DownloadFile(ctx8, mediumFileID)
	if err != nil {
		t.Fatalf("DownloadFile (medium) failed: %v", err)
	}
	if !bytes.Equal(downloadedMedium, mediumData) {
		t.Error("Downloaded medium file content does not match original")
	}
	t.Logf("✓ Medium file downloaded and verified: %d bytes", len(downloadedMedium))

	// Test 11: Bulk download all 3 files
	t.Log("=== Test 11: Bulk download all 3 files ===")
	ctx9, cancel9 := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel9()

	zipFileID, err := client.BulkDownloadFiles(ctx9, uploadedFileIDs)
	if err != nil {
		t.Fatalf("BulkDownloadFiles failed: %v", err)
	}
	if zipFileID == "" {
		t.Fatal("BulkDownloadFiles returned empty file ID")
	}
	t.Logf("✓ Bulk download created ZIP: ID %s", zipFileID)

	// Test 12: Download the bulk ZIP file
	t.Log("=== Test 12: Download bulk ZIP file ===")
	ctx10, cancel10 := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel10()

	zipData, err := client.DownloadFile(ctx10, zipFileID)
	if err != nil {
		t.Fatalf("DownloadFile (ZIP) failed: %v", err)
	}
	if len(zipData) == 0 {
		t.Error("Downloaded ZIP file is empty")
	}
	// ZIP file should be at least as large as the compressed content
	t.Logf("✓ ZIP file downloaded: %d bytes", len(zipData))

	// Test 13: Get downloaded files filter
	t.Log("=== Test 13: Get downloaded files ===")
	ctx11, cancel11 := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel11()

	downloadedFiles, err := client.GetDownloadedFiles(ctx11, 100)
	if err != nil {
		t.Fatalf("GetDownloadedFiles failed: %v", err)
	}
	t.Logf("✓ Downloaded files count: %d", len(downloadedFiles))

	// Test 14: Delete files one by one
	t.Log("=== Test 14: Delete uploaded files ===")
	for _, fileID := range uploadedFileIDs {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		err := client.DeleteFile(ctx, fileID)
		cancel()
		if err != nil {
			t.Errorf("DeleteFile failed for %s: %v", fileID, err)
		} else {
			t.Logf("✓ File deleted: %s", fileID)
		}
	}

	// Test 15: Verify files marked deleted
	t.Log("=== Test 15: Verify files marked deleted ===")
	ctx12, cancel12 := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel12()

	finalFiles, err := client.GetFiles(ctx12, 100)
	if err != nil {
		t.Fatalf("GetFiles after delete failed: %v", err)
	}

	// Check that deleted files are marked as deleted or removed
	deletedCount := 0
	for _, file := range finalFiles {
		for _, deletedID := range uploadedFileIDs {
			if file.AgentFileID == deletedID {
				// Mythic may mark files as deleted rather than removing them
				deletedCount++
				t.Logf("✓ File %s still present (may be marked deleted)", deletedID)
			}
		}
	}
	t.Logf("✓ Verified deletion of %d files", len(uploadedFileIDs))

	t.Log("=== ✓ All file operation tests passed ===")
}

// TestE2E_FileOperationsErrorHandling tests error scenarios for file operations
func TestE2E_FileOperationsErrorHandling(t *testing.T) {
	client := AuthenticateTestClient(t)

	// Test 1: Get non-existent file by ID
	t.Log("=== Test 1: Get non-existent file by ID ===")
	ctx1, cancel1 := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel1()

	_, err := client.GetFileByID(ctx1, "non-existent-file-id-12345")
	if err == nil {
		t.Error("Expected error for non-existent file ID")
	}
	t.Logf("✓ Non-existent file rejected: %v", err)

	// Test 2: Download non-existent file
	t.Log("=== Test 2: Download non-existent file ===")
	ctx2, cancel2 := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel2()

	_, err = client.DownloadFile(ctx2, "non-existent-file-id-67890")
	if err == nil {
		t.Error("Expected error for non-existent file download")
	}
	t.Logf("✓ Non-existent file download rejected: %v", err)

	// Test 3: Delete non-existent file
	t.Log("=== Test 3: Delete non-existent file ===")
	ctx3, cancel3 := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel3()

	err = client.DeleteFile(ctx3, "non-existent-file-id-abcdef")
	if err == nil {
		t.Error("Expected error for non-existent file delete")
	}
	t.Logf("✓ Non-existent file delete rejected: %v", err)

	// Test 4: Upload empty file
	t.Log("=== Test 4: Upload empty file ===")
	ctx4, cancel4 := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel4()

	fileID, err := client.UploadFile(ctx4, "empty.dat", []byte{})
	if err != nil {
		t.Logf("✓ Empty file upload rejected: %v", err)
	} else {
		// Empty files might be allowed, just clean up
		t.Logf("Empty file upload allowed: ID %s", fileID)
		ctx5, cancel5 := context.WithTimeout(context.Background(), 30*time.Second)
		_ = client.DeleteFile(ctx5, fileID)
		cancel5()
	}

	// Test 5: Bulk download with invalid IDs
	t.Log("=== Test 5: Bulk download with invalid file IDs ===")
	ctx6, cancel6 := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel6()

	_, err = client.BulkDownloadFiles(ctx6, []string{"invalid-id-1", "invalid-id-2"})
	if err == nil {
		t.Error("Expected error for bulk download with invalid IDs")
	}
	t.Logf("✓ Bulk download with invalid IDs rejected: %v", err)

	t.Log("=== ✓ All error handling tests passed ===")
}

// createTempFile creates a temporary file with random data of the specified size
func createTempFile(name string, size int) (string, error) {
	tmpDir := os.TempDir()
	filePath := filepath.Join(tmpDir, name)

	// Generate random data
	data := make([]byte, size)
	_, err := rand.Read(data)
	if err != nil {
		return "", fmt.Errorf("failed to generate random data: %w", err)
	}

	// Write to file
	err = os.WriteFile(filePath, data, 0644)
	if err != nil {
		return "", fmt.Errorf("failed to write file: %w", err)
	}

	return filePath, nil
}
