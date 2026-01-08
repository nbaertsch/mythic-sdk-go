//go:build integration

package integration

import (
	"context"
	"testing"
)

// TestContainers_ListFilesInvalidContainer tests error handling with invalid container name.
func TestContainers_ListFilesInvalidContainer(t *testing.T) {
	client := getTestClient(t)
	ctx := context.Background()

	// Test with empty container name
	_, err := client.ContainerListFiles(ctx, "", "/tmp")
	if err == nil {
		t.Fatal("ContainerListFiles with empty container name should return error")
	}
	t.Logf("Empty container name error: %v", err)

	// Test with empty path
	_, err = client.ContainerListFiles(ctx, "mythic_server", "")
	if err == nil {
		t.Fatal("ContainerListFiles with empty path should return error")
	}
	t.Logf("Empty path error: %v", err)
}

// TestContainers_ListFilesNonexistentContainer tests behavior with nonexistent container.
func TestContainers_ListFilesNonexistentContainer(t *testing.T) {
	client := getTestClient(t)
	ctx := context.Background()

	// Use a container that likely doesn't exist
	_, err := client.ContainerListFiles(ctx, "nonexistent_container_12345", "/tmp")
	if err == nil {
		t.Fatal("ContainerListFiles with nonexistent container should return error")
	}
	t.Logf("Nonexistent container error: %v", err)
}

// TestContainers_ListFilesInvalidPath tests behavior with invalid path.
func TestContainers_ListFilesInvalidPath(t *testing.T) {
	client := getTestClient(t)
	ctx := context.Background()

	// Try to list files in a path that likely doesn't exist
	// Using mythic_server container which should exist
	files, err := client.ContainerListFiles(ctx, "mythic_server", "/nonexistent_path_12345")
	if err != nil {
		// It's okay if this returns an error - path doesn't exist
		t.Logf("Invalid path error (expected): %v", err)
		return
	}

	// Or it might return an empty list
	if len(files) > 0 {
		t.Errorf("Expected empty list or error for invalid path, got %d files", len(files))
	}
	t.Logf("Invalid path returned empty list (expected)")
}

// TestContainers_DownloadFileInvalidInput tests error handling for download.
func TestContainers_DownloadFileInvalidInput(t *testing.T) {
	client := getTestClient(t)
	ctx := context.Background()

	// Test with empty container name
	_, err := client.ContainerDownloadFile(ctx, "", "/tmp/test.txt")
	if err == nil {
		t.Fatal("ContainerDownloadFile with empty container name should return error")
	}
	t.Logf("Empty container name error: %v", err)

	// Test with empty path
	_, err = client.ContainerDownloadFile(ctx, "mythic_server", "")
	if err == nil {
		t.Fatal("ContainerDownloadFile with empty path should return error")
	}
	t.Logf("Empty path error: %v", err)
}

// TestContainers_DownloadFileNonexistent tests downloading nonexistent file.
func TestContainers_DownloadFileNonexistent(t *testing.T) {
	client := getTestClient(t)
	ctx := context.Background()

	// Try to download a file that doesn't exist
	_, err := client.ContainerDownloadFile(ctx, "mythic_server", "/nonexistent_file_12345.txt")
	if err == nil {
		t.Fatal("ContainerDownloadFile with nonexistent file should return error")
	}
	t.Logf("Nonexistent file error: %v", err)
}

// TestContainers_WriteFileInvalidInput tests error handling for write.
func TestContainers_WriteFileInvalidInput(t *testing.T) {
	client := getTestClient(t)
	ctx := context.Background()

	content := []byte("test content")

	// Test with empty container name
	err := client.ContainerWriteFile(ctx, "", "/tmp/test.txt", content)
	if err == nil {
		t.Fatal("ContainerWriteFile with empty container name should return error")
	}
	t.Logf("Empty container name error: %v", err)

	// Test with empty path
	err = client.ContainerWriteFile(ctx, "mythic_server", "", content)
	if err == nil {
		t.Fatal("ContainerWriteFile with empty path should return error")
	}
	t.Logf("Empty path error: %v", err)
}

// TestContainers_RemoveFileInvalidInput tests error handling for remove.
func TestContainers_RemoveFileInvalidInput(t *testing.T) {
	client := getTestClient(t)
	ctx := context.Background()

	// Test with empty container name
	err := client.ContainerRemoveFile(ctx, "", "/tmp/test.txt")
	if err == nil {
		t.Fatal("ContainerRemoveFile with empty container name should return error")
	}
	t.Logf("Empty container name error: %v", err)

	// Test with empty path
	err = client.ContainerRemoveFile(ctx, "mythic_server", "")
	if err == nil {
		t.Fatal("ContainerRemoveFile with empty path should return error")
	}
	t.Logf("Empty path error: %v", err)
}

// TestContainers_WriteDownloadRemoveFile tests full file lifecycle.
// This test is conditional - it will skip if container operations aren't supported.
func TestContainers_WriteDownloadRemoveFile(t *testing.T) {
	client := getTestClient(t)
	ctx := context.Background()

	containerName := "mythic_server"
	testPath := "/tmp/mythic_sdk_test_file.txt"
	testContent := []byte("This is a test file created by mythic-sdk-go integration tests\n")

	// Try to write a test file
	t.Logf("Writing test file to %s:%s", containerName, testPath)
	err := client.ContainerWriteFile(ctx, containerName, testPath, testContent)
	if err != nil {
		t.Skipf("Container write not supported or container not available: %v", err)
		return
	}
	t.Logf("Successfully wrote %d bytes", len(testContent))

	// Try to download the file we just wrote
	t.Logf("Downloading test file from %s:%s", containerName, testPath)
	downloadedContent, err := client.ContainerDownloadFile(ctx, containerName, testPath)
	if err != nil {
		t.Errorf("Failed to download file we just wrote: %v", err)
	} else {
		t.Logf("Successfully downloaded %d bytes", len(downloadedContent))

		// Verify content matches
		if string(downloadedContent) != string(testContent) {
			t.Errorf("Downloaded content doesn't match written content")
			t.Errorf("Expected: %q", string(testContent))
			t.Errorf("Got: %q", string(downloadedContent))
		} else {
			t.Log("Content verification successful")
		}
	}

	// Clean up - remove the test file
	t.Logf("Removing test file from %s:%s", containerName, testPath)
	err = client.ContainerRemoveFile(ctx, containerName, testPath)
	if err != nil {
		t.Errorf("Failed to remove test file: %v", err)
	} else {
		t.Log("Successfully removed test file")
	}

	// Verify file was removed by trying to download it again
	t.Log("Verifying file was removed")
	_, err = client.ContainerDownloadFile(ctx, containerName, testPath)
	if err == nil {
		t.Error("File should not exist after removal")
	} else {
		t.Logf("Confirmed file was removed (error expected): %v", err)
	}
}

// TestContainers_ListFilesStructure tests the structure of returned file info.
// This test is conditional and will skip if no containers are accessible.
func TestContainers_ListFilesStructure(t *testing.T) {
	client := getTestClient(t)
	ctx := context.Background()

	// Try to list files in a common directory
	containerName := "mythic_server"
	path := "/tmp"

	t.Logf("Listing files in %s:%s", containerName, path)
	files, err := client.ContainerListFiles(ctx, containerName, path)
	if err != nil {
		t.Skipf("Container list not supported or container not available: %v", err)
		return
	}

	t.Logf("Found %d files/directories", len(files))

	// Log first few entries for debugging
	for i, file := range files {
		if i < 10 { // Log first 10
			t.Logf("File %d: %s", i+1, file.String())

			// Verify basic structure
			if file.Name == "" {
				t.Errorf("File %d has empty name", i)
			}

			// Test helper methods
			_ = file.IsDirectory()
			_ = file.String()
		}
	}
}

// TestContainers_WriteEmptyFile tests writing an empty file.
func TestContainers_WriteEmptyFile(t *testing.T) {
	client := getTestClient(t)
	ctx := context.Background()

	containerName := "mythic_server"
	testPath := "/tmp/mythic_sdk_test_empty.txt"
	emptyContent := []byte{}

	t.Logf("Writing empty file to %s:%s", containerName, testPath)
	err := client.ContainerWriteFile(ctx, containerName, testPath, emptyContent)
	if err != nil {
		t.Skipf("Container write not supported or container not available: %v", err)
		return
	}
	t.Log("Successfully wrote empty file")

	// Download and verify it's empty
	content, err := client.ContainerDownloadFile(ctx, containerName, testPath)
	if err != nil {
		t.Errorf("Failed to download empty file: %v", err)
	} else {
		if len(content) != 0 {
			t.Errorf("Expected empty file, got %d bytes", len(content))
		} else {
			t.Log("Verified file is empty")
		}
	}

	// Clean up
	err = client.ContainerRemoveFile(ctx, containerName, testPath)
	if err != nil {
		t.Errorf("Failed to remove empty test file: %v", err)
	}
}

// TestContainers_WriteLargeFile tests writing a larger file.
func TestContainers_WriteLargeFile(t *testing.T) {
	client := getTestClient(t)
	ctx := context.Background()

	containerName := "mythic_server"
	testPath := "/tmp/mythic_sdk_test_large.bin"

	// Create a 1MB test file
	largeContent := make([]byte, 1024*1024)
	for i := range largeContent {
		largeContent[i] = byte(i % 256)
	}

	t.Logf("Writing large file (%d bytes) to %s:%s", len(largeContent), containerName, testPath)
	err := client.ContainerWriteFile(ctx, containerName, testPath, largeContent)
	if err != nil {
		t.Skipf("Container write not supported or container not available: %v", err)
		return
	}
	t.Log("Successfully wrote large file")

	// Download and verify size
	content, err := client.ContainerDownloadFile(ctx, containerName, testPath)
	if err != nil {
		t.Errorf("Failed to download large file: %v", err)
	} else {
		if len(content) != len(largeContent) {
			t.Errorf("Size mismatch: expected %d bytes, got %d bytes", len(largeContent), len(content))
		} else {
			t.Logf("Verified file size: %d bytes", len(content))
		}
	}

	// Clean up
	err = client.ContainerRemoveFile(ctx, containerName, testPath)
	if err != nil {
		t.Errorf("Failed to remove large test file: %v", err)
	}
}
