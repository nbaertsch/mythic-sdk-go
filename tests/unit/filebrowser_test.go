package unit

import (
	"testing"
	"time"

	"github.com/nbaertsch/mythic-sdk-go/pkg/mythic/types"
)

// TestFileBrowserObjectString tests the FileBrowserObject.String() method
func TestFileBrowserObjectString(t *testing.T) {
	tests := []struct {
		name     string
		obj      types.FileBrowserObject
		contains []string
	}{
		{
			name: "file object",
			obj: types.FileBrowserObject{
				ID:           1,
				IsFile:       true,
				Name:         "test.txt",
				FullPathText: "/home/user/test.txt",
			},
			contains: []string{"file", "/home/user/test.txt"},
		},
		{
			name: "directory object",
			obj: types.FileBrowserObject{
				ID:           2,
				IsFile:       false,
				Name:         "Documents",
				FullPathText: "/home/user/Documents",
			},
			contains: []string{"dir", "/home/user/Documents"},
		},
		{
			name: "deleted file",
			obj: types.FileBrowserObject{
				ID:           3,
				IsFile:       true,
				Name:         "deleted.txt",
				FullPathText: "/tmp/deleted.txt",
				Deleted:      true,
			},
			contains: []string{"file", "/tmp/deleted.txt", "deleted"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.obj.String()
			if result == "" {
				t.Error("String() should not return empty string")
			}
			for _, substr := range tt.contains {
				if !contains(result, substr) {
					t.Errorf("String() = %q, should contain %q", result, substr)
				}
			}
		})
	}
}

// TestFileBrowserObjectIsDirectory tests the FileBrowserObject.IsDirectory() method
func TestFileBrowserObjectIsDirectory(t *testing.T) {
	tests := []struct {
		name string
		obj  types.FileBrowserObject
		want bool
	}{
		{
			name: "directory",
			obj: types.FileBrowserObject{
				ID:     1,
				IsFile: false,
				Name:   "Documents",
			},
			want: true,
		},
		{
			name: "file",
			obj: types.FileBrowserObject{
				ID:     2,
				IsFile: true,
				Name:   "test.txt",
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.obj.IsDirectory(); got != tt.want {
				t.Errorf("IsDirectory() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestFileBrowserObjectIsDeleted tests the FileBrowserObject.IsDeleted() method
func TestFileBrowserObjectIsDeleted(t *testing.T) {
	tests := []struct {
		name string
		obj  types.FileBrowserObject
		want bool
	}{
		{
			name: "deleted object",
			obj: types.FileBrowserObject{
				ID:      1,
				Name:    "deleted.txt",
				Deleted: true,
			},
			want: true,
		},
		{
			name: "active object",
			obj: types.FileBrowserObject{
				ID:      2,
				Name:    "active.txt",
				Deleted: false,
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.obj.IsDeleted(); got != tt.want {
				t.Errorf("IsDeleted() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestFileBrowserObjectGetFullPath tests the FileBrowserObject.GetFullPath() method
func TestFileBrowserObjectGetFullPath(t *testing.T) {
	tests := []struct {
		name string
		obj  types.FileBrowserObject
		want string
	}{
		{
			name: "with full_path_text",
			obj: types.FileBrowserObject{
				ID:           1,
				Name:         "test.txt",
				ParentPath:   "/home/user",
				FullPathText: "/home/user/test.txt",
			},
			want: "/home/user/test.txt",
		},
		{
			name: "root parent path",
			obj: types.FileBrowserObject{
				ID:         2,
				Name:       "etc",
				ParentPath: "/",
			},
			want: "/etc",
		},
		{
			name: "empty parent path",
			obj: types.FileBrowserObject{
				ID:         3,
				Name:       "tmp",
				ParentPath: "",
			},
			want: "/tmp",
		},
		{
			name: "normal path",
			obj: types.FileBrowserObject{
				ID:         4,
				Name:       "file.txt",
				ParentPath: "/home/user/Documents",
			},
			want: "/home/user/Documents/file.txt",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.obj.GetFullPath(); got != tt.want {
				t.Errorf("GetFullPath() = %q, want %q", got, tt.want)
			}
		})
	}
}

// TestFileBrowserObjectTypes tests the FileBrowserObject type structure
func TestFileBrowserObjectTypes(t *testing.T) {
	now := time.Now()
	obj := types.FileBrowserObject{
		ID:            1,
		Host:          "target-host",
		IsFile:        true,
		Permissions:   "-rw-r--r--",
		Name:          "test.txt",
		ParentPath:    "/home/user",
		Success:       true,
		AccessTime:    now,
		ModifyTime:    now,
		Size:          1024,
		UpdateDeleted: false,
		TaskID:        5,
		OperationID:   2,
		Timestamp:     now,
		Comment:       "Test file",
		Deleted:       false,
		FullPathText:  "/home/user/test.txt",
		CallbackID:    3,
		OperatorID:    1,
	}

	if obj.ID != 1 {
		t.Errorf("Expected ID 1, got %d", obj.ID)
	}
	if obj.Host != "target-host" {
		t.Errorf("Expected Host 'target-host', got %q", obj.Host)
	}
	if !obj.IsFile {
		t.Error("Expected IsFile to be true")
	}
	if obj.Permissions != "-rw-r--r--" {
		t.Errorf("Expected Permissions '-rw-r--r--', got %q", obj.Permissions)
	}
	if obj.Name != "test.txt" {
		t.Errorf("Expected Name 'test.txt', got %q", obj.Name)
	}
	if obj.ParentPath != "/home/user" {
		t.Errorf("Expected ParentPath '/home/user', got %q", obj.ParentPath)
	}
	if !obj.Success {
		t.Error("Expected Success to be true")
	}
	if obj.Size != 1024 {
		t.Errorf("Expected Size 1024, got %d", obj.Size)
	}
	if obj.TaskID != 5 {
		t.Errorf("Expected TaskID 5, got %d", obj.TaskID)
	}
	if obj.OperationID != 2 {
		t.Errorf("Expected OperationID 2, got %d", obj.OperationID)
	}
	if obj.CallbackID != 3 {
		t.Errorf("Expected CallbackID 3, got %d", obj.CallbackID)
	}
	if obj.OperatorID != 1 {
		t.Errorf("Expected OperatorID 1, got %d", obj.OperatorID)
	}
	if obj.Comment != "Test file" {
		t.Errorf("Expected Comment 'Test file', got %q", obj.Comment)
	}
	if obj.Deleted {
		t.Error("Expected Deleted to be false")
	}
	if obj.FullPathText != "/home/user/test.txt" {
		t.Errorf("Expected FullPathText '/home/user/test.txt', got %q", obj.FullPathText)
	}
}

// TestFileBrowserObjectFileVsDirectory tests distinguishing files from directories
func TestFileBrowserObjectFileVsDirectory(t *testing.T) {
	file := types.FileBrowserObject{
		ID:     1,
		IsFile: true,
		Name:   "document.pdf",
		Size:   2048,
	}

	directory := types.FileBrowserObject{
		ID:     2,
		IsFile: false,
		Name:   "Downloads",
		Size:   0,
	}

	if file.IsDirectory() {
		t.Error("File should not be identified as directory")
	}

	if !directory.IsDirectory() {
		t.Error("Directory should be identified as directory")
	}

	// Verify String() output differs
	fileStr := file.String()
	dirStr := directory.String()

	if !contains(fileStr, "file") {
		t.Errorf("File String() should contain 'file', got %q", fileStr)
	}

	if !contains(dirStr, "dir") {
		t.Errorf("Directory String() should contain 'dir', got %q", dirStr)
	}
}

// TestFileBrowserObjectDeletionState tests deletion state tracking
func TestFileBrowserObjectDeletionState(t *testing.T) {
	active := types.FileBrowserObject{
		ID:      1,
		Name:    "active.txt",
		Deleted: false,
	}

	deleted := types.FileBrowserObject{
		ID:      2,
		Name:    "removed.txt",
		Deleted: true,
	}

	if active.IsDeleted() {
		t.Error("Active object should not be deleted")
	}

	if !deleted.IsDeleted() {
		t.Error("Deleted object should be deleted")
	}

	// Verify String() reflects deletion state
	deletedStr := deleted.String()
	if !contains(deletedStr, "deleted") {
		t.Errorf("Deleted object String() should contain 'deleted', got %q", deletedStr)
	}

	activeStr := active.String()
	if contains(activeStr, "deleted") {
		t.Errorf("Active object String() should not contain 'deleted', got %q", activeStr)
	}
}

// TestFileBrowserObjectPathConstruction tests path construction logic
func TestFileBrowserObjectPathConstruction(t *testing.T) {
	tests := []struct {
		name         string
		objName      string
		parentPath   string
		fullPathText string
		expected     string
	}{
		{
			name:         "uses full_path_text if available",
			objName:      "file.txt",
			parentPath:   "/incorrect",
			fullPathText: "/correct/path/file.txt",
			expected:     "/correct/path/file.txt",
		},
		{
			name:         "constructs from root",
			objName:      "bin",
			parentPath:   "/",
			fullPathText: "",
			expected:     "/bin",
		},
		{
			name:         "constructs from empty parent",
			objName:      "usr",
			parentPath:   "",
			fullPathText: "",
			expected:     "/usr",
		},
		{
			name:         "constructs normal path",
			objName:      "config.txt",
			parentPath:   "/etc/app",
			fullPathText: "",
			expected:     "/etc/app/config.txt",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			obj := types.FileBrowserObject{
				Name:         tt.objName,
				ParentPath:   tt.parentPath,
				FullPathText: tt.fullPathText,
			}

			got := obj.GetFullPath()
			if got != tt.expected {
				t.Errorf("GetFullPath() = %q, expected %q", got, tt.expected)
			}
		})
	}
}
