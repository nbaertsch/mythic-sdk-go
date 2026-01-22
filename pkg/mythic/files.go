package mythic

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"strings"
	"time"
)

// Timestamp is a custom type to handle Mythic's timestamp format
// which may not include timezone information.
type Timestamp struct {
	time.Time
}

// UnmarshalJSON handles timestamp parsing for both RFC3339 and Mythic's format.
func (mt *Timestamp) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), `"`)
	if s == "null" || s == "" {
		mt.Time = time.Time{}
		return nil
	}

	// Try RFC3339 first (standard format with timezone)
	t, err := time.Parse(time.RFC3339, s)
	if err == nil {
		mt.Time = t
		return nil
	}

	// Try RFC3339 with nanoseconds
	t, err = time.Parse(time.RFC3339Nano, s)
	if err == nil {
		mt.Time = t
		return nil
	}

	// Try Mythic's format without timezone (treat as UTC)
	formats := []string{
		"2006-01-02T15:04:05.999999",
		"2006-01-02T15:04:05",
		"2006-01-02 15:04:05.999999",
		"2006-01-02 15:04:05",
	}

	for _, format := range formats {
		t, err = time.Parse(format, s)
		if err == nil {
			mt.Time = t.UTC()
			return nil
		}
	}

	return fmt.Errorf("unable to parse timestamp: %s", s)
}

// MarshalJSON implements json.Marshaler.
func (mt Timestamp) MarshalJSON() ([]byte, error) {
	if mt.Time.IsZero() {
		return []byte("null"), nil
	}
	return json.Marshal(mt.Time.Format(time.RFC3339))
}

// FileMeta represents a file in Mythic.
type FileMeta struct {
	ID                  int       `json:"id"`
	AgentFileID         string    `json:"agent_file_id"`
	TotalChunks         int       `json:"total_chunks"`
	ChunksReceived      int       `json:"chunks_received"`
	Complete            bool      `json:"complete"`
	Path                string    `json:"path"`
	FullRemotePath      string    `json:"full_remote_path"`
	Host                string    `json:"host"`
	IsPayload           bool      `json:"is_payload"`
	IsScreenshot        bool      `json:"is_screenshot"`
	IsDownloadFromAgent bool      `json:"is_download_from_agent"`
	Filename            string    `json:"filename"`
	MD5                 string    `json:"md5"`
	SHA1                string    `json:"sha1"`
	Size                int64     `json:"size"`
	Comment             string    `json:"comment"`
	OperatorID          int       `json:"operator_id"`
	Timestamp           time.Time `json:"timestamp"`
	Deleted             bool      `json:"deleted"`
	TaskID              *int      `json:"task_id,omitempty"`
}

// FileUploadResponse represents the response from uploading a file.
type FileUploadResponse struct {
	AgentFileID string `json:"agent_file_id"`
	Status      string `json:"status"`
}

// GetFiles retrieves all files for the current operation.
func (c *Client) GetFiles(ctx context.Context, limit int) ([]*FileMeta, error) {
	if err := c.EnsureAuthenticated(ctx); err != nil {
		return nil, err
	}

	if limit <= 0 {
		limit = 100
	}

	var query struct {
		FileMeta []struct {
			ID                  int    `graphql:"id"`
			AgentFileID         string `graphql:"agent_file_id"`
			TotalChunks         int    `graphql:"total_chunks"`
			ChunksReceived      int    `graphql:"chunks_received"`
			Complete            bool   `graphql:"complete"`
			Path                string `graphql:"path"`
			FullRemotePath      string `graphql:"full_remote_path"`
			Host                string `graphql:"host"`
			IsPayload           bool   `graphql:"is_payload"`
			IsScreenshot        bool   `graphql:"is_screenshot"`
			IsDownloadFromAgent bool   `graphql:"is_download_from_agent"`
			Filename            string `graphql:"filename_text"`
			MD5                 string `graphql:"md5"`
			SHA1                string `graphql:"sha1"`
			Size                int    `graphql:"size"`
			Comment             string `graphql:"comment"`
			OperatorID          int    `graphql:"operator_id"`
			Timestamp           string `graphql:"timestamp"`
			Deleted             bool   `graphql:"deleted"`
			TaskID              *int   `graphql:"task_id"`
		} `graphql:"filemeta(order_by: {id: desc}, limit: $limit)"`
	}

	variables := map[string]interface{}{
		"limit": limit,
	}

	err := c.executeQuery(ctx, &query, variables)
	if err != nil {
		return nil, WrapError("GetFiles", err, "failed to query files")
	}

	files := make([]*FileMeta, 0, len(query.FileMeta))
	for _, f := range query.FileMeta {
		// Parse timestamp
		var timestamp time.Time
		if f.Timestamp != "" {
			var mt Timestamp
			if err := mt.UnmarshalJSON([]byte(`"` + f.Timestamp + `"`)); err == nil {
				timestamp = mt.Time
			}
		}

		files = append(files, &FileMeta{
			ID:                  f.ID,
			AgentFileID:         f.AgentFileID,
			TotalChunks:         f.TotalChunks,
			ChunksReceived:      f.ChunksReceived,
			Complete:            f.Complete,
			Path:                f.Path,
			FullRemotePath:      f.FullRemotePath,
			Host:                f.Host,
			IsPayload:           f.IsPayload,
			IsScreenshot:        f.IsScreenshot,
			IsDownloadFromAgent: f.IsDownloadFromAgent,
			Filename:            decodeFilename(f.Filename),
			MD5:                 f.MD5,
			SHA1:                f.SHA1,
			Size:                int64(f.Size),
			Comment:             f.Comment,
			OperatorID:          f.OperatorID,
			Timestamp:           timestamp,
			Deleted:             f.Deleted,
			TaskID:              f.TaskID,
		})
	}

	return files, nil
}

// GetFileByID retrieves a specific file by its agent file ID.
func (c *Client) GetFileByID(ctx context.Context, agentFileID string) (*FileMeta, error) {
	if err := c.EnsureAuthenticated(ctx); err != nil {
		return nil, err
	}

	var query struct {
		FileMeta []struct {
			ID                  int    `graphql:"id"`
			AgentFileID         string `graphql:"agent_file_id"`
			TotalChunks         int    `graphql:"total_chunks"`
			ChunksReceived      int    `graphql:"chunks_received"`
			Complete            bool   `graphql:"complete"`
			Path                string `graphql:"path"`
			FullRemotePath      string `graphql:"full_remote_path"`
			Host                string `graphql:"host"`
			IsPayload           bool   `graphql:"is_payload"`
			IsScreenshot        bool   `graphql:"is_screenshot"`
			IsDownloadFromAgent bool   `graphql:"is_download_from_agent"`
			Filename            string `graphql:"filename_text"`
			MD5                 string `graphql:"md5"`
			SHA1                string `graphql:"sha1"`
			Size                int    `graphql:"size"`
			Comment             string `graphql:"comment"`
			OperatorID          int    `graphql:"operator_id"`
			Timestamp           string `graphql:"timestamp"`
			Deleted             bool   `graphql:"deleted"`
			TaskID              *int   `graphql:"task_id"`
		} `graphql:"filemeta(where: {agent_file_id: {_eq: $agent_file_id}}, limit: 1)"`
	}

	variables := map[string]interface{}{
		"agent_file_id": agentFileID,
	}

	err := c.executeQuery(ctx, &query, variables)
	if err != nil {
		return nil, WrapError("GetFileByID", err, "failed to query file")
	}

	if len(query.FileMeta) == 0 {
		return nil, WrapError("GetFileByID", ErrNotFound, fmt.Sprintf("file with agent_file_id %s not found", agentFileID))
	}

	f := query.FileMeta[0]

	// Parse timestamp
	var timestamp time.Time
	if f.Timestamp != "" {
		var mt Timestamp
		if err := mt.UnmarshalJSON([]byte(`"` + f.Timestamp + `"`)); err == nil {
			timestamp = mt.Time
		}
	}

	return &FileMeta{
		ID:                  f.ID,
		AgentFileID:         f.AgentFileID,
		TotalChunks:         f.TotalChunks,
		ChunksReceived:      f.ChunksReceived,
		Complete:            f.Complete,
		Path:                f.Path,
		FullRemotePath:      f.FullRemotePath,
		Host:                f.Host,
		IsPayload:           f.IsPayload,
		IsScreenshot:        f.IsScreenshot,
		IsDownloadFromAgent: f.IsDownloadFromAgent,
		Filename:            decodeFilename(f.Filename),
		MD5:                 f.MD5,
		SHA1:                f.SHA1,
		Size:                int64(f.Size),
		Comment:             f.Comment,
		OperatorID:          f.OperatorID,
		Timestamp:           timestamp,
		Deleted:             f.Deleted,
		TaskID:              f.TaskID,
	}, nil
}

// GetDownloadedFiles retrieves all files downloaded from agents.
func (c *Client) GetDownloadedFiles(ctx context.Context, limit int) ([]*FileMeta, error) {
	if err := c.EnsureAuthenticated(ctx); err != nil {
		return nil, err
	}

	if limit <= 0 {
		limit = 100
	}

	var query struct {
		FileMeta []struct {
			ID                  int    `graphql:"id"`
			AgentFileID         string `graphql:"agent_file_id"`
			TotalChunks         int    `graphql:"total_chunks"`
			ChunksReceived      int    `graphql:"chunks_received"`
			Complete            bool   `graphql:"complete"`
			Path                string `graphql:"path"`
			FullRemotePath      string `graphql:"full_remote_path"`
			Host                string `graphql:"host"`
			IsPayload           bool   `graphql:"is_payload"`
			IsScreenshot        bool   `graphql:"is_screenshot"`
			IsDownloadFromAgent bool   `graphql:"is_download_from_agent"`
			Filename            string `graphql:"filename_text"`
			MD5                 string `graphql:"md5"`
			SHA1                string `graphql:"sha1"`
			Size                int    `graphql:"size"`
			Comment             string `graphql:"comment"`
			OperatorID          int    `graphql:"operator_id"`
			Timestamp           string `graphql:"timestamp"`
			Deleted             bool   `graphql:"deleted"`
			TaskID              *int   `graphql:"task_id"`
		} `graphql:"filemeta(where: {is_download_from_agent: {_eq: true}, deleted: {_eq: false}}, order_by: {id: desc}, limit: $limit)"`
	}

	variables := map[string]interface{}{
		"limit": limit,
	}

	err := c.executeQuery(ctx, &query, variables)
	if err != nil {
		return nil, WrapError("GetDownloadedFiles", err, "failed to query downloaded files")
	}

	files := make([]*FileMeta, 0, len(query.FileMeta))
	for _, f := range query.FileMeta {
		// Parse timestamp
		var timestamp time.Time
		if f.Timestamp != "" {
			var mt Timestamp
			if err := mt.UnmarshalJSON([]byte(`"` + f.Timestamp + `"`)); err == nil {
				timestamp = mt.Time
			}
		}

		files = append(files, &FileMeta{
			ID:                  f.ID,
			AgentFileID:         f.AgentFileID,
			TotalChunks:         f.TotalChunks,
			ChunksReceived:      f.ChunksReceived,
			Complete:            f.Complete,
			Path:                f.Path,
			FullRemotePath:      f.FullRemotePath,
			Host:                f.Host,
			IsPayload:           f.IsPayload,
			IsScreenshot:        f.IsScreenshot,
			IsDownloadFromAgent: f.IsDownloadFromAgent,
			Filename:            decodeFilename(f.Filename),
			MD5:                 f.MD5,
			SHA1:                f.SHA1,
			Size:                int64(f.Size),
			Comment:             f.Comment,
			OperatorID:          f.OperatorID,
			Timestamp:           timestamp,
			Deleted:             f.Deleted,
			TaskID:              f.TaskID,
		})
	}

	return files, nil
}

// UploadFile uploads a file to Mythic for use in tasks.
// Returns the agent_file_id that can be used to reference the file.
func (c *Client) UploadFile(ctx context.Context, filename string, fileData []byte) (string, error) {
	if err := c.EnsureAuthenticated(ctx); err != nil {
		return "", err
	}

	if filename == "" {
		return "", WrapError("UploadFile", ErrInvalidInput, "filename is required")
	}

	if len(fileData) == 0 {
		return "", WrapError("UploadFile", ErrInvalidInput, "file data is required")
	}

	// Construct upload endpoint URL
	scheme := "https"
	if !c.config.SSL {
		scheme = "http"
	}
	uploadURL := fmt.Sprintf("%s://%s/api/v1.4/task_upload_file_webhook", scheme, stripScheme(c.config.ServerURL))

	// Create multipart form
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Add file field
	part, err := writer.CreateFormFile("file", filename)
	if err != nil {
		return "", WrapError("UploadFile", err, "failed to create form file")
	}

	_, err = part.Write(fileData)
	if err != nil {
		return "", WrapError("UploadFile", err, "failed to write file data")
	}

	err = writer.Close()
	if err != nil {
		return "", WrapError("UploadFile", err, "failed to close multipart writer")
	}

	// Create HTTP request
	req, err := http.NewRequestWithContext(ctx, "POST", uploadURL, body)
	if err != nil {
		return "", WrapError("UploadFile", err, "failed to create upload request")
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())

	// Add authentication headers
	authHeaders := c.getAuthHeaders()
	for key, value := range authHeaders {
		req.Header.Set(key, value)
	}

	// Execute request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", WrapError("UploadFile", err, "failed to execute upload request")
	}
	defer resp.Body.Close() //nolint:errcheck // Response body close error not critical

	// Read response
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", WrapError("UploadFile", err, "failed to read upload response")
	}

	if resp.StatusCode != http.StatusOK {
		return "", WrapError("UploadFile", ErrInvalidResponse, fmt.Sprintf("upload failed with status %d: %s", resp.StatusCode, string(respBody)))
	}

	// Parse response - Mythic returns {"agent_file_id": "...", "status": "success"}
	var uploadResp FileUploadResponse
	if err := parseJSON(respBody, &uploadResp); err != nil {
		return "", WrapError("UploadFile", err, "failed to parse upload response")
	}

	if uploadResp.AgentFileID == "" {
		return "", WrapError("UploadFile", ErrInvalidResponse, "no agent_file_id in response")
	}

	return uploadResp.AgentFileID, nil
}

// DownloadFile downloads a file's content from Mythic.
func (c *Client) DownloadFile(ctx context.Context, agentFileID string) ([]byte, error) {
	if err := c.EnsureAuthenticated(ctx); err != nil {
		return nil, err
	}

	if agentFileID == "" {
		return nil, WrapError("DownloadFile", ErrInvalidInput, "agent_file_id is required")
	}

	// Construct download endpoint URL
	scheme := "https"
	if !c.config.SSL {
		scheme = "http"
	}
	downloadURL := fmt.Sprintf("%s://%s/api/v1.4/files/download/%s", scheme, stripScheme(c.config.ServerURL), agentFileID)

	// Create HTTP request
	req, err := http.NewRequestWithContext(ctx, "GET", downloadURL, nil)
	if err != nil {
		return nil, WrapError("DownloadFile", err, "failed to create download request")
	}

	// Add authentication headers
	authHeaders := c.getAuthHeaders()
	for key, value := range authHeaders {
		req.Header.Set(key, value)
	}

	// Execute request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, WrapError("DownloadFile", err, "failed to execute download request")
	}
	defer resp.Body.Close() //nolint:errcheck // Response body close error not critical

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body) //nolint:errcheck // Best effort to read error message
		return nil, WrapError("DownloadFile", ErrInvalidResponse, fmt.Sprintf("download failed with status %d: %s", resp.StatusCode, string(body)))
	}

	// Read file content
	fileData, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, WrapError("DownloadFile", err, "failed to read file data")
	}

	// Check if response is JSON (Mythic returns JSON for errors and base64-encoded files)
	if len(fileData) > 0 && fileData[0] == '{' {
		// First check for error response
		var errorResp struct {
			Status string `json:"status"`
			Error  string `json:"error"`
		}
		if err := parseJSON(fileData, &errorResp); err == nil {
			if errorResp.Status == "error" {
				return nil, WrapError("DownloadFile", ErrNotFound, errorResp.Error)
			}
		}

		// Check if response is base64-encoded JSON (Mythic sometimes returns {file: base64data})
		var jsonResp struct {
			File string `json:"file"`
		}
		if err := parseJSON(fileData, &jsonResp); err == nil && jsonResp.File != "" {
			// Decode base64
			decoded, err := base64.StdEncoding.DecodeString(jsonResp.File)
			if err == nil {
				return decoded, nil
			}
		}
	}

	return fileData, nil
}

// DeleteFile marks a file as deleted in Mythic.
func (c *Client) DeleteFile(ctx context.Context, agentFileID string) error {
	if err := c.EnsureAuthenticated(ctx); err != nil {
		return err
	}

	if agentFileID == "" {
		return WrapError("DeleteFile", ErrInvalidInput, "agent_file_id is required")
	}

	var mutation struct {
		UpdateFileMeta struct {
			Affected int `graphql:"affected_rows"`
		} `graphql:"update_filemeta(where: {agent_file_id: {_eq: $agent_file_id}}, _set: {deleted: true})"`
	}

	variables := map[string]interface{}{
		"agent_file_id": agentFileID,
	}

	err := c.executeMutation(ctx, &mutation, variables)
	if err != nil {
		return WrapError("DeleteFile", err, "failed to delete file")
	}

	if mutation.UpdateFileMeta.Affected == 0 {
		return WrapError("DeleteFile", ErrNotFound, fmt.Sprintf("file with agent_file_id %s not found", agentFileID))
	}

	return nil
}

// String returns a string representation of the file.
func (f *FileMeta) String() string {
	status := "incomplete"
	if f.Complete {
		status = "complete"
	}

	return fmt.Sprintf("File %s: %s (%d bytes, %s)", f.AgentFileID, f.Filename, f.Size, status)
}

// IsComplete returns whether the file has been fully received.
func (f *FileMeta) IsComplete() bool {
	return f.Complete
}

// IsDeleted returns whether the file has been marked as deleted.
func (f *FileMeta) IsDeleted() bool {
	return f.Deleted
}

// FilePreview represents a preview of a file's contents.
type FilePreview struct {
	Size           int    `json:"size"`
	Host           string `json:"host"`
	FullRemotePath string `json:"full_remote_path"`
	Filename       string `json:"filename"`
	Contents       string `json:"contents"`
}

// BulkDownloadFiles creates a ZIP archive of multiple files and returns the file ID for download.
// The returned file_id can be used with DownloadFile() to retrieve the ZIP archive.
func (c *Client) BulkDownloadFiles(ctx context.Context, agentFileIDs []string) (string, error) {
	if err := c.EnsureAuthenticated(ctx); err != nil {
		return "", err
	}

	if len(agentFileIDs) == 0 {
		return "", WrapError("BulkDownloadFiles", ErrInvalidInput, "at least one file ID required")
	}

	var mutation struct {
		DownloadBulk struct {
			Status string `graphql:"status"`
			Error  string `graphql:"error"`
			FileID string `graphql:"file_id"`
		} `graphql:"download_bulk(files: $files)"`
	}

	variables := map[string]interface{}{
		"files": agentFileIDs,
	}

	err := c.executeMutation(ctx, &mutation, variables)
	if err != nil {
		return "", WrapError("BulkDownloadFiles", err, "failed to create bulk download")
	}

	if mutation.DownloadBulk.Status != "success" {
		return "", WrapError("BulkDownloadFiles", ErrInvalidResponse, fmt.Sprintf("bulk download failed: %s", mutation.DownloadBulk.Error))
	}

	if mutation.DownloadBulk.FileID == "" {
		return "", WrapError("BulkDownloadFiles", ErrInvalidResponse, "no file_id returned")
	}

	return mutation.DownloadBulk.FileID, nil
}

// PreviewFile retrieves file metadata and a preview of its contents without downloading the full file.
func (c *Client) PreviewFile(ctx context.Context, agentFileID string) (*FilePreview, error) {
	if err := c.EnsureAuthenticated(ctx); err != nil {
		return nil, err
	}

	if agentFileID == "" {
		return nil, WrapError("PreviewFile", ErrInvalidInput, "agent_file_id is required")
	}

	var mutation struct {
		PreviewFile struct {
			Status         string `graphql:"status"`
			Error          string `graphql:"error"`
			Size           int    `graphql:"size"`
			Host           string `graphql:"host"`
			FullRemotePath string `graphql:"full_remote_path"`
			Filename       string `graphql:"filename"`
			Contents       string `graphql:"contents"`
		} `graphql:"previewFile(file_id: $file_id)"`
	}

	variables := map[string]interface{}{
		"file_id": agentFileID,
	}

	err := c.executeMutation(ctx, &mutation, variables)
	if err != nil {
		return nil, WrapError("PreviewFile", err, "failed to preview file")
	}

	if mutation.PreviewFile.Status != "success" {
		return nil, WrapError("PreviewFile", ErrInvalidResponse, fmt.Sprintf("preview failed: %s", mutation.PreviewFile.Error))
	}

	return &FilePreview{
		Size:           mutation.PreviewFile.Size,
		Host:           mutation.PreviewFile.Host,
		FullRemotePath: mutation.PreviewFile.FullRemotePath,
		Filename:       mutation.PreviewFile.Filename,
		Contents:       mutation.PreviewFile.Contents,
	}, nil
}
