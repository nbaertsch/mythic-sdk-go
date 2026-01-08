package unit

import (
	"strings"
	"testing"

	"github.com/nbaertsch/mythic-sdk-go/pkg/mythic/types"
)

func TestContainerFileInfo_String(t *testing.T) {
	tests := []struct {
		name     string
		file     types.ContainerFileInfo
		contains []string
	}{
		{
			name: "regular file",
			file: types.ContainerFileInfo{
				Name:  "config.json",
				Size:  1024,
				IsDir: false,
			},
			contains: []string{"config.json", "file", "1024"},
		},
		{
			name: "directory",
			file: types.ContainerFileInfo{
				Name:  "agent_code",
				Size:  0,
				IsDir: true,
			},
			contains: []string{"agent_code", "dir"},
		},
		{
			name: "large file",
			file: types.ContainerFileInfo{
				Name:  "payload.exe",
				Size:  15728640, // 15 MB
				IsDir: false,
			},
			contains: []string{"payload.exe", "file", "15728640"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.file.String()

			for _, substr := range tt.contains {
				if !strings.Contains(result, substr) {
					t.Errorf("String() = %q, should contain %q", result, substr)
				}
			}
		})
	}
}

func TestContainerFileInfo_IsDirectory(t *testing.T) {
	tests := []struct {
		name     string
		file     types.ContainerFileInfo
		expected bool
	}{
		{
			name: "file is not directory",
			file: types.ContainerFileInfo{
				Name:  "file.txt",
				IsDir: false,
			},
			expected: false,
		},
		{
			name: "directory",
			file: types.ContainerFileInfo{
				Name:  "folder",
				IsDir: true,
			},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.file.IsDirectory()
			if result != tt.expected {
				t.Errorf("IsDirectory() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestContainerListFilesRequest_String(t *testing.T) {
	tests := []struct {
		name     string
		request  types.ContainerListFilesRequest
		contains []string
	}{
		{
			name: "athena container",
			request: types.ContainerListFilesRequest{
				ContainerName: "mythic_athena",
				Path:          "/Mythic",
			},
			contains: []string{"mythic_athena", "/Mythic", "List"},
		},
		{
			name: "c2 profile container",
			request: types.ContainerListFilesRequest{
				ContainerName: "http",
				Path:          "/srv/config",
			},
			contains: []string{"http", "/srv/config"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.request.String()

			for _, substr := range tt.contains {
				if !strings.Contains(result, substr) {
					t.Errorf("String() = %q, should contain %q", result, substr)
				}
			}
		})
	}
}

func TestContainerDownloadFileRequest_String(t *testing.T) {
	tests := []struct {
		name     string
		request  types.ContainerDownloadFileRequest
		contains []string
	}{
		{
			name: "download config",
			request: types.ContainerDownloadFileRequest{
				ContainerName: "mythic_athena",
				Path:          "/Mythic/config.json",
			},
			contains: []string{"mythic_athena", "/Mythic/config.json", "Download"},
		},
		{
			name: "download source file",
			request: types.ContainerDownloadFileRequest{
				ContainerName: "poseidon",
				Path:          "/go/src/main.go",
			},
			contains: []string{"poseidon", "/go/src/main.go"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.request.String()

			for _, substr := range tt.contains {
				if !strings.Contains(result, substr) {
					t.Errorf("String() = %q, should contain %q", result, substr)
				}
			}
		})
	}
}

func TestContainerWriteFileRequest_String(t *testing.T) {
	tests := []struct {
		name     string
		request  types.ContainerWriteFileRequest
		contains []string
	}{
		{
			name: "write small file",
			request: types.ContainerWriteFileRequest{
				ContainerName: "mythic_athena",
				Path:          "/tmp/test.txt",
				Content:       []byte("hello world"),
			},
			contains: []string{"mythic_athena", "/tmp/test.txt", "11 bytes", "Write"},
		},
		{
			name: "write large file",
			request: types.ContainerWriteFileRequest{
				ContainerName: "http",
				Path:          "/srv/data.bin",
				Content:       make([]byte, 10240),
			},
			contains: []string{"http", "/srv/data.bin", "10240 bytes"},
		},
		{
			name: "write empty file",
			request: types.ContainerWriteFileRequest{
				ContainerName: "poseidon",
				Path:          "/tmp/empty",
				Content:       []byte{},
			},
			contains: []string{"poseidon", "/tmp/empty", "0 bytes"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.request.String()

			for _, substr := range tt.contains {
				if !strings.Contains(result, substr) {
					t.Errorf("String() = %q, should contain %q", result, substr)
				}
			}
		})
	}
}

func TestContainerRemoveFileRequest_String(t *testing.T) {
	tests := []struct {
		name     string
		request  types.ContainerRemoveFileRequest
		contains []string
	}{
		{
			name: "remove temp file",
			request: types.ContainerRemoveFileRequest{
				ContainerName: "mythic_athena",
				Path:          "/tmp/test.txt",
			},
			contains: []string{"mythic_athena", "/tmp/test.txt", "Remove"},
		},
		{
			name: "remove config",
			request: types.ContainerRemoveFileRequest{
				ContainerName: "http",
				Path:          "/srv/old_config.json",
			},
			contains: []string{"http", "/srv/old_config.json"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.request.String()

			for _, substr := range tt.contains {
				if !strings.Contains(result, substr) {
					t.Errorf("String() = %q, should contain %q", result, substr)
				}
			}
		})
	}
}

func TestContainerFileInfo_Permissions(t *testing.T) {
	tests := []struct {
		name       string
		file       types.ContainerFileInfo
		permission string
	}{
		{
			name: "executable file",
			file: types.ContainerFileInfo{
				Name:       "script.sh",
				Permission: "rwxr-xr-x",
			},
			permission: "rwxr-xr-x",
		},
		{
			name: "read-only file",
			file: types.ContainerFileInfo{
				Name:       "config.json",
				Permission: "r--r--r--",
			},
			permission: "r--r--r--",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.file.Permission != tt.permission {
				t.Errorf("Permission = %q, want %q", tt.file.Permission, tt.permission)
			}
		})
	}
}

func TestContainerFileInfo_ZeroSize(t *testing.T) {
	file := types.ContainerFileInfo{
		Name:  "empty.txt",
		Size:  0,
		IsDir: false,
	}

	if file.Size != 0 {
		t.Errorf("Size = %d, want 0", file.Size)
	}

	str := file.String()
	if !strings.Contains(str, "0 bytes") {
		t.Errorf("String() should show 0 bytes for empty file, got %q", str)
	}
}
