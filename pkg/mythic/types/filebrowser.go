package types

import (
	"fmt"
	"time"
)

// FileBrowserObject represents a file or directory in the file browser.
type FileBrowserObject struct {
	ID               int       `json:"id"`
	Host             string    `json:"host"`
	IsFile           bool      `json:"is_file"`
	Permissions      string    `json:"permissions"`
	Name             string    `json:"name"`
	ParentPath       string    `json:"parent_path"`
	Success          bool      `json:"success"`
	AccessTime       time.Time `json:"access_time"`
	ModifyTime       time.Time `json:"modify_time"`
	Size             int64     `json:"size"`
	UpdateDeleted    bool      `json:"update_deleted"`
	TaskID           int       `json:"task_id"`
	OperationID      int       `json:"operation_id"`
	Timestamp        time.Time `json:"timestamp"`
	Comment          string    `json:"comment"`
	Deleted          bool      `json:"deleted"`
	FullPathText     string    `json:"full_path_text"`
	CallbackID       int       `json:"callback_id"`
	OperatorID       int       `json:"operator_id"`
}

// String returns a string representation of a FileBrowserObject.
func (f *FileBrowserObject) String() string {
	itemType := "file"
	if !f.IsFile {
		itemType = "dir"
	}
	deleted := ""
	if f.Deleted {
		deleted = " (deleted)"
	}
	return fmt.Sprintf("[%s] %s%s", itemType, f.FullPathText, deleted)
}

// IsDirectory returns true if the object is a directory.
func (f *FileBrowserObject) IsDirectory() bool {
	return !f.IsFile
}

// IsDeleted returns true if the object has been deleted.
func (f *FileBrowserObject) IsDeleted() bool {
	return f.Deleted
}

// GetFullPath returns the full path of the file or directory.
func (f *FileBrowserObject) GetFullPath() string {
	if f.FullPathText != "" {
		return f.FullPathText
	}
	if f.ParentPath == "" || f.ParentPath == "/" {
		return "/" + f.Name
	}
	return f.ParentPath + "/" + f.Name
}
