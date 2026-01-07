package unit

import (
	"testing"
	"time"

	"github.com/nbaertsch/mythic-sdk-go/pkg/mythic"
)

func TestFileMeta_String(t *testing.T) {
	file := &mythic.FileMeta{
		AgentFileID: "abc-123",
		Filename:    "test.txt",
		Size:        1024,
		Complete:    true,
	}

	expected := "File abc-123: test.txt (1024 bytes, complete)"
	if file.String() != expected {
		t.Errorf("Expected %q, got %q", expected, file.String())
	}
}

func TestFileMeta_String_Incomplete(t *testing.T) {
	file := &mythic.FileMeta{
		AgentFileID: "abc-123",
		Filename:    "incomplete.bin",
		Size:        2048,
		Complete:    false,
	}

	expected := "File abc-123: incomplete.bin (2048 bytes, incomplete)"
	if file.String() != expected {
		t.Errorf("Expected %q, got %q", expected, file.String())
	}
}

func TestFileMeta_IsComplete(t *testing.T) {
	tests := []struct {
		name     string
		complete bool
		expected bool
	}{
		{"Complete file", true, true},
		{"Incomplete file", false, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			file := &mythic.FileMeta{Complete: tt.complete}
			if got := file.IsComplete(); got != tt.expected {
				t.Errorf("IsComplete() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestFileMeta_IsDeleted(t *testing.T) {
	tests := []struct {
		name     string
		deleted  bool
		expected bool
	}{
		{"Deleted file", true, true},
		{"Active file", false, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			file := &mythic.FileMeta{Deleted: tt.deleted}
			if got := file.IsDeleted(); got != tt.expected {
				t.Errorf("IsDeleted() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestFileMeta_Structure(t *testing.T) {
	now := time.Now()
	taskID := 123

	file := &mythic.FileMeta{
		ID:                  1,
		AgentFileID:         "file-abc-123",
		TotalChunks:         10,
		ChunksReceived:      10,
		Complete:            true,
		Path:                "/tmp",
		FullRemotePath:      "/tmp/test.txt",
		Host:                "workstation01",
		IsPayload:           false,
		IsScreenshot:        false,
		IsDownloadFromAgent: true,
		Filename:            "test.txt",
		MD5:                 "d41d8cd98f00b204e9800998ecf8427e",
		SHA1:                "da39a3ee5e6b4b0d3255bfef95601890afd80709",
		Size:                1024,
		Comment:             "Test file",
		OperatorID:          1,
		Timestamp:           now,
		Deleted:             false,
		TaskID:              &taskID,
	}

	// Verify all fields are set correctly
	if file.ID != 1 {
		t.Errorf("Expected ID 1, got %d", file.ID)
	}
	if file.AgentFileID != "file-abc-123" {
		t.Errorf("Expected AgentFileID 'file-abc-123', got %q", file.AgentFileID)
	}
	if file.TotalChunks != 10 {
		t.Errorf("Expected TotalChunks 10, got %d", file.TotalChunks)
	}
	if file.ChunksReceived != 10 {
		t.Errorf("Expected ChunksReceived 10, got %d", file.ChunksReceived)
	}
	if !file.Complete {
		t.Error("Expected Complete true, got false")
	}
	if file.Path != "/tmp" {
		t.Errorf("Expected Path '/tmp', got %q", file.Path)
	}
	if file.FullRemotePath != "/tmp/test.txt" {
		t.Errorf("Expected FullRemotePath '/tmp/test.txt', got %q", file.FullRemotePath)
	}
	if file.Host != "workstation01" {
		t.Errorf("Expected Host 'workstation01', got %q", file.Host)
	}
	if file.IsPayload {
		t.Error("Expected IsPayload false, got true")
	}
	if file.IsScreenshot {
		t.Error("Expected IsScreenshot false, got true")
	}
	if !file.IsDownloadFromAgent {
		t.Error("Expected IsDownloadFromAgent true, got false")
	}
	if file.Filename != "test.txt" {
		t.Errorf("Expected Filename 'test.txt', got %q", file.Filename)
	}
	if file.MD5 != "d41d8cd98f00b204e9800998ecf8427e" {
		t.Errorf("Expected MD5 hash, got %q", file.MD5)
	}
	if file.SHA1 != "da39a3ee5e6b4b0d3255bfef95601890afd80709" {
		t.Errorf("Expected SHA1 hash, got %q", file.SHA1)
	}
	if file.Size != 1024 {
		t.Errorf("Expected Size 1024, got %d", file.Size)
	}
	if file.Comment != "Test file" {
		t.Errorf("Expected Comment 'Test file', got %q", file.Comment)
	}
	if file.OperatorID != 1 {
		t.Errorf("Expected OperatorID 1, got %d", file.OperatorID)
	}
	if !file.Timestamp.Equal(now) {
		t.Errorf("Expected Timestamp %v, got %v", now, file.Timestamp)
	}
	if file.Deleted {
		t.Error("Expected Deleted false, got true")
	}
	if file.TaskID == nil || *file.TaskID != 123 {
		t.Errorf("Expected TaskID 123, got %v", file.TaskID)
	}
}

func TestFileMeta_DownloadFile(t *testing.T) {
	file := &mythic.FileMeta{
		AgentFileID:         "download-123",
		Filename:            "data.bin",
		IsDownloadFromAgent: true,
		Complete:            true,
		Size:                2048,
	}

	if !file.IsDownloadFromAgent {
		t.Error("Expected IsDownloadFromAgent true for downloaded file")
	}

	if !file.IsComplete() {
		t.Error("Downloaded file should be complete")
	}

	if file.IsPayload {
		t.Error("Downloaded file should not be marked as payload")
	}
}

func TestFileMeta_Payload(t *testing.T) {
	file := &mythic.FileMeta{
		AgentFileID: "payload-456",
		Filename:    "agent.exe",
		IsPayload:   true,
		Complete:    true,
		Size:        512000,
	}

	if !file.IsPayload {
		t.Error("Expected IsPayload true for payload file")
	}

	if !file.IsComplete() {
		t.Error("Payload file should be complete")
	}

	if file.IsDownloadFromAgent {
		t.Error("Payload file should not be marked as download from agent")
	}
}

func TestFileMeta_Screenshot(t *testing.T) {
	file := &mythic.FileMeta{
		AgentFileID:  "screenshot-789",
		Filename:     "screen_2024.png",
		IsScreenshot: true,
		Complete:     true,
		Size:         150000,
	}

	if !file.IsScreenshot {
		t.Error("Expected IsScreenshot true for screenshot file")
	}

	if !file.IsComplete() {
		t.Error("Screenshot should be complete")
	}
}

func TestFileMeta_PartialChunks(t *testing.T) {
	file := &mythic.FileMeta{
		AgentFileID:    "partial-123",
		Filename:       "large.bin",
		TotalChunks:    100,
		ChunksReceived: 45,
		Complete:       false,
		Size:           10240000,
	}

	if file.ChunksReceived >= file.TotalChunks {
		t.Error("Partial file should not have all chunks received")
	}

	if file.IsComplete() {
		t.Error("Partial file should not be complete")
	}

	expectedStr := "File partial-123: large.bin (10240000 bytes, incomplete)"
	if file.String() != expectedStr {
		t.Errorf("Expected %q, got %q", expectedStr, file.String())
	}
}

func TestFileUploadResponse_Structure(t *testing.T) {
	resp := &mythic.FileUploadResponse{
		AgentFileID: "upload-abc-123",
		Status:      "success",
	}

	if resp.AgentFileID != "upload-abc-123" {
		t.Errorf("Expected AgentFileID 'upload-abc-123', got %q", resp.AgentFileID)
	}

	if resp.Status != "success" {
		t.Errorf("Expected Status 'success', got %q", resp.Status)
	}
}

func TestFileMeta_NilTaskID(t *testing.T) {
	file := &mythic.FileMeta{
		AgentFileID: "no-task-123",
		Filename:    "standalone.txt",
		TaskID:      nil,
	}

	if file.TaskID != nil {
		t.Error("Expected TaskID nil for file not associated with task")
	}
}

func TestFileMeta_ZeroSizeFile(t *testing.T) {
	file := &mythic.FileMeta{
		AgentFileID: "empty-123",
		Filename:    "empty.txt",
		Size:        0,
		Complete:    true,
	}

	if file.Size != 0 {
		t.Errorf("Expected Size 0, got %d", file.Size)
	}

	expectedStr := "File empty-123: empty.txt (0 bytes, complete)"
	if file.String() != expectedStr {
		t.Errorf("Expected %q, got %q", expectedStr, file.String())
	}
}

func TestFileMeta_LargeFile(t *testing.T) {
	file := &mythic.FileMeta{
		AgentFileID:    "large-123",
		Filename:       "bigdata.bin",
		Size:           5368709120, // 5GB
		TotalChunks:    5120,
		ChunksReceived: 5120,
		Complete:       true,
	}

	if file.Size != 5368709120 {
		t.Errorf("Expected Size 5368709120, got %d", file.Size)
	}

	if file.TotalChunks != file.ChunksReceived {
		t.Error("Large file should have all chunks received")
	}

	if !file.IsComplete() {
		t.Error("Large file with all chunks should be complete")
	}
}
