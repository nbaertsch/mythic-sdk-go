package mythic

import (
	"context"
	"encoding/base64"

	"github.com/nbaertsch/mythic-sdk-go/pkg/mythic/types"
)

// ContainerListFiles lists files in a Docker container directory.
// This is useful for browsing payload type and C2 profile container filesystems
// during development and debugging.
//
// Parameters:
//   - ctx: Context for cancellation and timeout
//   - containerName: Name of the Docker container (e.g., "mythic_athena", "http")
//   - path: Path within the container to list (e.g., "/Mythic", "/srv")
//
// Returns:
//   - []*types.ContainerFileInfo: List of files and directories
//   - error: Error if the operation fails
//
// Example:
//
//	files, err := client.ContainerListFiles(ctx, "mythic_athena", "/Mythic")
func (c *Client) ContainerListFiles(ctx context.Context, containerName, path string) ([]*types.ContainerFileInfo, error) {
	if err := c.EnsureAuthenticated(ctx); err != nil {
		return nil, err
	}

	// Validate input
	if containerName == "" {
		return nil, WrapError("ContainerListFiles", ErrInvalidInput, "container name cannot be empty")
	}
	if path == "" {
		return nil, WrapError("ContainerListFiles", ErrInvalidInput, "path cannot be empty")
	}

	var query struct {
		Response struct {
			Status string `graphql:"status"`
			Error  string `graphql:"error"`
			Files  []struct {
				Name       string `graphql:"name"`
				Size       int    `graphql:"size"`
				IsDir      bool   `graphql:"is_dir"`
				ModTime    string `graphql:"mod_time"`
				Permission string `graphql:"permission"`
			} `graphql:"files"`
		} `graphql:"containerListFiles(container_name: $container_name, path: $path)"`
	}

	variables := map[string]interface{}{
		"container_name": containerName,
		"path":           path,
	}

	if err := c.executeQuery(ctx, &query, variables); err != nil {
		return nil, WrapError("ContainerListFiles", err, "failed to list container files")
	}

	// Check for error in response
	if query.Response.Status != "success" {
		return nil, WrapError("ContainerListFiles", ErrOperationFailed, query.Response.Error)
	}

	// Map to types
	files := make([]*types.ContainerFileInfo, 0, len(query.Response.Files))
	for _, f := range query.Response.Files {
		files = append(files, &types.ContainerFileInfo{
			Name:       f.Name,
			Size:       int64(f.Size),
			IsDir:      f.IsDir,
			ModTime:    f.ModTime,
			Permission: f.Permission,
		})
	}

	return files, nil
}

// ContainerDownloadFile downloads a file from a Docker container.
// This allows retrieving files from payload type and C2 profile containers
// for backup, analysis, or debugging purposes.
//
// Parameters:
//   - ctx: Context for cancellation and timeout
//   - containerName: Name of the Docker container
//   - path: Path to the file within the container
//
// Returns:
//   - []byte: File content
//   - error: Error if the operation fails
//
// Example:
//
//	content, err := client.ContainerDownloadFile(ctx, "mythic_athena", "/Mythic/agent_code/main.go")
func (c *Client) ContainerDownloadFile(ctx context.Context, containerName, path string) ([]byte, error) {
	if err := c.EnsureAuthenticated(ctx); err != nil {
		return nil, err
	}

	// Validate input
	if containerName == "" {
		return nil, WrapError("ContainerDownloadFile", ErrInvalidInput, "container name cannot be empty")
	}
	if path == "" {
		return nil, WrapError("ContainerDownloadFile", ErrInvalidInput, "path cannot be empty")
	}

	var query struct {
		Response struct {
			Status  string `graphql:"status"`
			Error   string `graphql:"error"`
			Content string `graphql:"content"`
		} `graphql:"containerDownloadFile(container_name: $container_name, path: $path)"`
	}

	variables := map[string]interface{}{
		"container_name": containerName,
		"path":           path,
	}

	if err := c.executeQuery(ctx, &query, variables); err != nil {
		return nil, WrapError("ContainerDownloadFile", err, "failed to download container file")
	}

	// Check for error in response
	if query.Response.Status != "success" {
		return nil, WrapError("ContainerDownloadFile", ErrOperationFailed, query.Response.Error)
	}

	// Decode base64 content
	content, err := base64.StdEncoding.DecodeString(query.Response.Content)
	if err != nil {
		return nil, WrapError("ContainerDownloadFile", err, "failed to decode file content")
	}

	return content, nil
}

// ContainerWriteFile writes a file to a Docker container.
// This allows updating configuration files, adding scripts, or modifying
// payload type and C2 profile containers during development.
//
// Parameters:
//   - ctx: Context for cancellation and timeout
//   - containerName: Name of the Docker container
//   - path: Path where the file should be written in the container
//   - content: File content as bytes
//
// Returns:
//   - error: Error if the operation fails
//
// Example:
//
//	content := []byte("#!/bin/bash\necho 'Hello World'\n")
//	err := client.ContainerWriteFile(ctx, "mythic_athena", "/tmp/test.sh", content)
func (c *Client) ContainerWriteFile(ctx context.Context, containerName, path string, content []byte) error {
	if err := c.EnsureAuthenticated(ctx); err != nil {
		return err
	}

	// Validate input
	if containerName == "" {
		return WrapError("ContainerWriteFile", ErrInvalidInput, "container name cannot be empty")
	}
	if path == "" {
		return WrapError("ContainerWriteFile", ErrInvalidInput, "path cannot be empty")
	}

	// Encode content as base64
	encodedContent := base64.StdEncoding.EncodeToString(content)

	var mutation struct {
		Response struct {
			Status string `graphql:"status"`
			Error  string `graphql:"error"`
		} `graphql:"containerWriteFile(container_name: $container_name, path: $path, content: $content)"`
	}

	variables := map[string]interface{}{
		"container_name": containerName,
		"path":           path,
		"content":        encodedContent,
	}

	if err := c.executeMutation(ctx, &mutation, variables); err != nil {
		return WrapError("ContainerWriteFile", err, "failed to write container file")
	}

	// Check for error in response
	if mutation.Response.Status != "success" {
		return WrapError("ContainerWriteFile", ErrOperationFailed, mutation.Response.Error)
	}

	return nil
}

// ContainerRemoveFile removes a file from a Docker container.
// This allows cleaning up temporary files, removing old configurations,
// or deleting files from payload type and C2 profile containers.
//
// Parameters:
//   - ctx: Context for cancellation and timeout
//   - containerName: Name of the Docker container
//   - path: Path to the file to remove in the container
//
// Returns:
//   - error: Error if the operation fails
//
// Example:
//
//	err := client.ContainerRemoveFile(ctx, "mythic_athena", "/tmp/test.sh")
func (c *Client) ContainerRemoveFile(ctx context.Context, containerName, path string) error {
	if err := c.EnsureAuthenticated(ctx); err != nil {
		return err
	}

	// Validate input
	if containerName == "" {
		return WrapError("ContainerRemoveFile", ErrInvalidInput, "container name cannot be empty")
	}
	if path == "" {
		return WrapError("ContainerRemoveFile", ErrInvalidInput, "path cannot be empty")
	}

	var mutation struct {
		Response struct {
			Status string `graphql:"status"`
			Error  string `graphql:"error"`
		} `graphql:"containerRemoveFile(container_name: $container_name, path: $path)"`
	}

	variables := map[string]interface{}{
		"container_name": containerName,
		"path":           path,
	}

	if err := c.executeMutation(ctx, &mutation, variables); err != nil {
		return WrapError("ContainerRemoveFile", err, "failed to remove container file")
	}

	// Check for error in response
	if mutation.Response.Status != "success" {
		return WrapError("ContainerRemoveFile", ErrOperationFailed, mutation.Response.Error)
	}

	return nil
}
