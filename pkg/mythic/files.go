package mythic

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"time"
)

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
			ID                  int       `graphql:"id"`
			AgentFileID         string    `graphql:"agent_file_id"`
			TotalChunks         int       `graphql:"total_chunks"`
			ChunksReceived      int       `graphql:"chunks_received"`
			Complete            bool      `graphql:"complete"`
			Path                string    `graphql:"path"`
			FullRemotePath      string    `graphql:"full_remote_path"`
			Host                string    `graphql:"host"`
			IsPayload           bool      `graphql:"is_payload"`
			IsScreenshot        bool      `graphql:"is_screenshot"`
			IsDownloadFromAgent bool      `graphql:"is_download_from_agent"`
			Filename            string    `graphql:"filename_text"`
			MD5                 string    `graphql:"md5"`
			SHA1                string    `graphql:"sha1"`
			Size                int       `graphql:"size"`
			Comment             string    `graphql:"comment"`
			OperatorID          int       `graphql:"operator_id"`
			Timestamp           time.Time `graphql:"timestamp"`
			Deleted             bool      `graphql:"deleted"`
			TaskID              *int      `graphql:"task_id"`
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
			Filename:            f.Filename,
			MD5:                 f.MD5,
			SHA1:                f.SHA1,
			Size:                int64(f.Size),
			Comment:             f.Comment,
			OperatorID:          f.OperatorID,
			Timestamp:           f.Timestamp,
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
			ID                  int       `graphql:"id"`
			AgentFileID         string    `graphql:"agent_file_id"`
			TotalChunks         int       `graphql:"total_chunks"`
			ChunksReceived      int       `graphql:"chunks_received"`
			Complete            bool      `graphql:"complete"`
			Path                string    `graphql:"path"`
			FullRemotePath      string    `graphql:"full_remote_path"`
			Host                string    `graphql:"host"`
			IsPayload           bool      `graphql:"is_payload"`
			IsScreenshot        bool      `graphql:"is_screenshot"`
			IsDownloadFromAgent bool      `graphql:"is_download_from_agent"`
			Filename            string    `graphql:"filename_text"`
			MD5                 string    `graphql:"md5"`
			SHA1                string    `graphql:"sha1"`
			Size                int       `graphql:"size"`
			Comment             string    `graphql:"comment"`
			OperatorID          int       `graphql:"operator_id"`
			Timestamp           time.Time `graphql:"timestamp"`
			Deleted             bool      `graphql:"deleted"`
			TaskID              *int      `graphql:"task_id"`
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
		Filename:            f.Filename,
		MD5:                 f.MD5,
		SHA1:                f.SHA1,
		Size:                int64(f.Size),
		Comment:             f.Comment,
		OperatorID:          f.OperatorID,
		Timestamp:           f.Timestamp,
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
			ID                  int       `graphql:"id"`
			AgentFileID         string    `graphql:"agent_file_id"`
			TotalChunks         int       `graphql:"total_chunks"`
			ChunksReceived      int       `graphql:"chunks_received"`
			Complete            bool      `graphql:"complete"`
			Path                string    `graphql:"path"`
			FullRemotePath      string    `graphql:"full_remote_path"`
			Host                string    `graphql:"host"`
			IsPayload           bool      `graphql:"is_payload"`
			IsScreenshot        bool      `graphql:"is_screenshot"`
			IsDownloadFromAgent bool      `graphql:"is_download_from_agent"`
			Filename            string    `graphql:"filename_text"`
			MD5                 string    `graphql:"md5"`
			SHA1                string    `graphql:"sha1"`
			Size                int       `graphql:"size"`
			Comment             string    `graphql:"comment"`
			OperatorID          int       `graphql:"operator_id"`
			Timestamp           time.Time `graphql:"timestamp"`
			Deleted             bool      `graphql:"deleted"`
			TaskID              *int      `graphql:"task_id"`
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
			Filename:            f.Filename,
			MD5:                 f.MD5,
			SHA1:                f.SHA1,
			Size:                int64(f.Size),
			Comment:             f.Comment,
			OperatorID:          f.OperatorID,
			Timestamp:           f.Timestamp,
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

	// Add authentication header
	if c.config.AccessToken != "" {
		req.Header.Set("Authorization", "Bearer "+c.config.AccessToken)
	} else if c.config.APIToken != "" {
		req.Header.Set("Authorization", "Bearer "+c.config.APIToken)
	}

	// Execute request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", WrapError("UploadFile", err, "failed to execute upload request")
	}
	defer resp.Body.Close()

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

	// Add authentication header
	if c.config.AccessToken != "" {
		req.Header.Set("Authorization", "Bearer "+c.config.AccessToken)
	} else if c.config.APIToken != "" {
		req.Header.Set("Authorization", "Bearer "+c.config.APIToken)
	}

	// Execute request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, WrapError("DownloadFile", err, "failed to execute download request")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body) //nolint:errcheck // Best effort to read error message
		return nil, WrapError("DownloadFile", ErrInvalidResponse, fmt.Sprintf("download failed with status %d: %s", resp.StatusCode, string(body)))
	}

	// Read file content
	fileData, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, WrapError("DownloadFile", err, "failed to read file data")
	}

	// Check if response is base64-encoded JSON (Mythic sometimes returns {file: base64data})
	if len(fileData) > 0 && fileData[0] == '{' {
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
