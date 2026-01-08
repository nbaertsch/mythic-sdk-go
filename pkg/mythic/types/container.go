package types

import "fmt"

// ContainerFileInfo represents a file or directory within a Docker container.
type ContainerFileInfo struct {
	Name       string `json:"name"`
	Size       int64  `json:"size"`
	IsDir      bool   `json:"is_dir"`
	ModTime    string `json:"mod_time"`
	Permission string `json:"permission"`
}

// String returns a human-readable representation of the container file.
func (c *ContainerFileInfo) String() string {
	fileType := "file"
	if c.IsDir {
		fileType = "dir"
	}
	return fmt.Sprintf("%s (%s, %d bytes)", c.Name, fileType, c.Size)
}

// IsDirectory returns true if the file info represents a directory.
func (c *ContainerFileInfo) IsDirectory() bool {
	return c.IsDir
}

// ContainerListFilesRequest represents a request to list files in a container.
type ContainerListFilesRequest struct {
	ContainerName string `json:"container_name"`
	Path          string `json:"path"`
}

// String returns a human-readable representation of the request.
func (c *ContainerListFilesRequest) String() string {
	return fmt.Sprintf("List files in %s:%s", c.ContainerName, c.Path)
}

// ContainerDownloadFileRequest represents a request to download a file from a container.
type ContainerDownloadFileRequest struct {
	ContainerName string `json:"container_name"`
	Path          string `json:"path"`
}

// String returns a human-readable representation of the request.
func (c *ContainerDownloadFileRequest) String() string {
	return fmt.Sprintf("Download %s from container %s", c.Path, c.ContainerName)
}

// ContainerWriteFileRequest represents a request to write a file to a container.
type ContainerWriteFileRequest struct {
	ContainerName string `json:"container_name"`
	Path          string `json:"path"`
	Content       []byte `json:"content"`
}

// String returns a human-readable representation of the request.
func (c *ContainerWriteFileRequest) String() string {
	return fmt.Sprintf("Write %d bytes to %s in container %s", len(c.Content), c.Path, c.ContainerName)
}

// ContainerRemoveFileRequest represents a request to remove a file from a container.
type ContainerRemoveFileRequest struct {
	ContainerName string `json:"container_name"`
	Path          string `json:"path"`
}

// String returns a human-readable representation of the request.
func (c *ContainerRemoveFileRequest) String() string {
	return fmt.Sprintf("Remove %s from container %s", c.Path, c.ContainerName)
}
