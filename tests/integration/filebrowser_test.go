//go:build integration

package integration

import (
	"context"
	"testing"

	"github.com/nbaertsch/mythic-sdk-go/pkg/mythic/types"
)

// TestFileBrowser_GetFileBrowserObjects tests retrieving all file browser objects.
func TestFileBrowser_GetFileBrowserObjects(t *testing.T) {
	client := getTestClient(t)
	ctx := context.Background()

	objects, err := client.GetFileBrowserObjects(ctx)
	if err != nil {
		t.Fatalf("GetFileBrowserObjects() failed: %v", err)
	}

	t.Logf("Found %d file browser objects", len(objects))

	if len(objects) == 0 {
		t.Log("No file browser objects found - this is expected if no file browsing has been performed")
		return
	}

	// Verify object structure
	for i, obj := range objects {
		if i < 5 { // Log first 5 for debugging
			t.Logf("Object %d: %s", i+1, obj.String())
		}

		if obj.ID == 0 {
			t.Errorf("Object %d has zero ID", i)
		}
		if obj.Host == "" {
			t.Errorf("Object %d has empty Host", i)
		}
		if obj.Name == "" {
			t.Errorf("Object %d has empty Name", i)
		}
		if obj.OperationID == 0 {
			t.Errorf("Object %s has zero OperationID", obj.Name)
		}
		if obj.CallbackID == 0 {
			t.Errorf("Object %s has zero CallbackID", obj.Name)
		}

		// Test helper methods
		_ = obj.IsDirectory()
		_ = obj.IsDeleted()
		fullPath := obj.GetFullPath()
		if fullPath == "" {
			t.Errorf("Object %s has empty full path", obj.Name)
		}
		_ = obj.String()

		// Verify no objects are deleted (we filter them out)
		if obj.Deleted {
			t.Errorf("Object %s should not be deleted (we filter deleted objects)", obj.Name)
		}
	}

	// Verify objects are sorted by full path
	for i := 1; i < len(objects); i++ {
		if objects[i-1].FullPathText > objects[i].FullPathText {
			t.Errorf("Objects not sorted by path: %s comes after %s",
				objects[i-1].FullPathText, objects[i].FullPathText)
			break
		}
	}
}

// TestFileBrowser_GetFileBrowserObjectsStructure tests object field values.
func TestFileBrowser_GetFileBrowserObjectsStructure(t *testing.T) {
	client := getTestClient(t)
	ctx := context.Background()

	objects, err := client.GetFileBrowserObjects(ctx)
	if err != nil {
		t.Fatalf("GetFileBrowserObjects() failed: %v", err)
	}

	if len(objects) == 0 {
		t.Skip("No file browser objects available for structure testing")
	}

	// Find a file and a directory
	var file, directory *types.FileBrowserObject
	for _, obj := range objects {
		if obj.IsFile && file == nil {
			file = obj
		}
		if obj.IsDirectory() && directory == nil {
			directory = obj
		}
		if file != nil && directory != nil {
			break
		}
	}

	if file != nil {
		t.Logf("Testing file structure: %s", file.String())
		if !file.IsFile {
			t.Error("File object should have IsFile=true")
		}
		if file.IsDirectory() {
			t.Error("File object should not be identified as directory")
		}
		if file.Size < 0 {
			t.Error("File should have non-negative size")
		}
	}

	if directory != nil {
		t.Logf("Testing directory structure: %s", directory.String())
		if directory.IsFile {
			t.Error("Directory object should have IsFile=false")
		}
		if !directory.IsDirectory() {
			t.Error("Directory object should be identified as directory")
		}
	}
}

// TestFileBrowser_GetFileBrowserObjectsByHostInvalid tests error handling for invalid host.
func TestFileBrowser_GetFileBrowserObjectsByHostInvalid(t *testing.T) {
	client := getTestClient(t)
	ctx := context.Background()

	// Test with empty host
	_, err := client.GetFileBrowserObjectsByHost(ctx, "")
	if err == nil {
		t.Fatal("GetFileBrowserObjectsByHost(\"\") should return an error")
	}
	t.Logf("Empty host error: %v", err)
}

// TestFileBrowser_GetFileBrowserObjectsByHost tests filtering by host.
func TestFileBrowser_GetFileBrowserObjectsByHost(t *testing.T) {
	client := getTestClient(t)
	ctx := context.Background()

	// First get all objects to find available hosts
	allObjects, err := client.GetFileBrowserObjects(ctx)
	if err != nil {
		t.Fatalf("GetFileBrowserObjects() failed: %v", err)
	}

	if len(allObjects) == 0 {
		t.Skip("No file browser objects available for host filtering test")
	}

	// Get unique hosts
	hostMap := make(map[string]bool)
	for _, obj := range allObjects {
		hostMap[obj.Host] = true
	}

	hosts := make([]string, 0, len(hostMap))
	for host := range hostMap {
		hosts = append(hosts, host)
	}

	t.Logf("Found %d unique hosts", len(hosts))

	// Test with first host
	if len(hosts) > 0 {
		testHost := hosts[0]
		t.Logf("Testing with host: %s", testHost)

		objects, err := client.GetFileBrowserObjectsByHost(ctx, testHost)
		if err != nil {
			t.Fatalf("GetFileBrowserObjectsByHost(%s) failed: %v", testHost, err)
		}

		t.Logf("Found %d objects for host %s", len(objects), testHost)

		// Verify all objects belong to the requested host
		for _, obj := range objects {
			if obj.Host != testHost {
				t.Errorf("Object %s has host %s, expected %s", obj.Name, obj.Host, testHost)
			}
		}

		// Verify sorting
		for i := 1; i < len(objects); i++ {
			if objects[i-1].FullPathText > objects[i].FullPathText {
				t.Errorf("Objects not sorted: %s comes after %s",
					objects[i-1].FullPathText, objects[i].FullPathText)
				break
			}
		}
	}
}

// TestFileBrowser_GetFileBrowserObjectsByHostNonexistent tests behavior with nonexistent host.
func TestFileBrowser_GetFileBrowserObjectsByHostNonexistent(t *testing.T) {
	client := getTestClient(t)
	ctx := context.Background()

	// Use a host that likely doesn't exist
	objects, err := client.GetFileBrowserObjectsByHost(ctx, "nonexistent-host-12345")
	if err != nil {
		t.Fatalf("GetFileBrowserObjectsByHost(nonexistent) failed: %v", err)
	}

	// Should return empty list, not an error
	if len(objects) != 0 {
		t.Errorf("Expected empty list for nonexistent host, got %d objects", len(objects))
	}
}

// TestFileBrowser_GetFileBrowserObjectsByCallbackInvalid tests error handling for invalid callback ID.
func TestFileBrowser_GetFileBrowserObjectsByCallbackInvalid(t *testing.T) {
	client := getTestClient(t)
	ctx := context.Background()

	// Test with zero callback ID
	_, err := client.GetFileBrowserObjectsByCallback(ctx, 0)
	if err == nil {
		t.Fatal("GetFileBrowserObjectsByCallback(0) should return an error")
	}
	t.Logf("Zero callback ID error: %v", err)
}

// TestFileBrowser_GetFileBrowserObjectsByCallback tests filtering by callback.
func TestFileBrowser_GetFileBrowserObjectsByCallback(t *testing.T) {
	client := getTestClient(t)
	ctx := context.Background()

	// First get all objects to find available callbacks
	allObjects, err := client.GetFileBrowserObjects(ctx)
	if err != nil {
		t.Fatalf("GetFileBrowserObjects() failed: %v", err)
	}

	if len(allObjects) == 0 {
		t.Skip("No file browser objects available for callback filtering test")
	}

	// Get unique callback IDs
	callbackMap := make(map[int]bool)
	for _, obj := range allObjects {
		callbackMap[obj.CallbackID] = true
	}

	callbackIDs := make([]int, 0, len(callbackMap))
	for id := range callbackMap {
		callbackIDs = append(callbackIDs, id)
	}

	t.Logf("Found %d unique callbacks", len(callbackIDs))

	// Test with first callback
	if len(callbackIDs) > 0 {
		testCallbackID := callbackIDs[0]
		t.Logf("Testing with callback ID: %d", testCallbackID)

		objects, err := client.GetFileBrowserObjectsByCallback(ctx, testCallbackID)
		if err != nil {
			t.Fatalf("GetFileBrowserObjectsByCallback(%d) failed: %v", testCallbackID, err)
		}

		t.Logf("Found %d objects for callback %d", len(objects), testCallbackID)

		// Verify all objects belong to the requested callback
		for _, obj := range objects {
			if obj.CallbackID != testCallbackID {
				t.Errorf("Object %s has callback ID %d, expected %d",
					obj.Name, obj.CallbackID, testCallbackID)
			}
		}

		// Verify sorting
		for i := 1; i < len(objects); i++ {
			if objects[i-1].FullPathText > objects[i].FullPathText {
				t.Errorf("Objects not sorted: %s comes after %s",
					objects[i-1].FullPathText, objects[i].FullPathText)
				break
			}
		}
	}
}

// TestFileBrowser_GetFileBrowserObjectsByCallbackNonexistent tests behavior with nonexistent callback.
func TestFileBrowser_GetFileBrowserObjectsByCallbackNonexistent(t *testing.T) {
	client := getTestClient(t)
	ctx := context.Background()

	// Use a callback ID that likely doesn't exist
	objects, err := client.GetFileBrowserObjectsByCallback(ctx, 999999)
	if err != nil {
		t.Fatalf("GetFileBrowserObjectsByCallback(999999) failed: %v", err)
	}

	// Should return empty list, not an error
	if len(objects) != 0 {
		t.Errorf("Expected empty list for nonexistent callback, got %d objects", len(objects))
	}
}

// TestFileBrowser_FileTypes tests distinguishing files from directories.
func TestFileBrowser_FileTypes(t *testing.T) {
	client := getTestClient(t)
	ctx := context.Background()

	objects, err := client.GetFileBrowserObjects(ctx)
	if err != nil {
		t.Fatalf("GetFileBrowserObjects() failed: %v", err)
	}

	if len(objects) == 0 {
		t.Skip("No file browser objects available")
	}

	fileCount := 0
	dirCount := 0

	for _, obj := range objects {
		if obj.IsFile {
			fileCount++
			if obj.IsDirectory() {
				t.Errorf("Object %s has IsFile=true but IsDirectory() returned true", obj.Name)
			}
		} else {
			dirCount++
			if !obj.IsDirectory() {
				t.Errorf("Object %s has IsFile=false but IsDirectory() returned false", obj.Name)
			}
		}
	}

	t.Logf("Found %d files and %d directories", fileCount, dirCount)
}

// TestFileBrowser_PathHandling tests path construction and handling.
func TestFileBrowser_PathHandling(t *testing.T) {
	client := getTestClient(t)
	ctx := context.Background()

	objects, err := client.GetFileBrowserObjects(ctx)
	if err != nil {
		t.Fatalf("GetFileBrowserObjects() failed: %v", err)
	}

	if len(objects) == 0 {
		t.Skip("No file browser objects available")
	}

	for i, obj := range objects {
		// GetFullPath() should always return a non-empty path
		fullPath := obj.GetFullPath()
		if fullPath == "" {
			t.Errorf("Object %d (%s) GetFullPath() returned empty string", i, obj.Name)
		}

		// If FullPathText is set, GetFullPath() should return it
		if obj.FullPathText != "" && fullPath != obj.FullPathText {
			t.Logf("Object %s: GetFullPath()=%q, FullPathText=%q (may be constructing path)",
				obj.Name, fullPath, obj.FullPathText)
		}

		// Verify String() includes the path
		str := obj.String()
		if !contains(str, fullPath) && !contains(str, obj.Name) {
			t.Errorf("Object String() should contain path or name: got %q", str)
		}
	}
}

// TestFileBrowser_MultipleHosts tests objects across multiple hosts.
func TestFileBrowser_MultipleHosts(t *testing.T) {
	client := getTestClient(t)
	ctx := context.Background()

	objects, err := client.GetFileBrowserObjects(ctx)
	if err != nil {
		t.Fatalf("GetFileBrowserObjects() failed: %v", err)
	}

	if len(objects) == 0 {
		t.Skip("No file browser objects available")
	}

	// Group objects by host
	hostObjects := make(map[string]int)
	for _, obj := range objects {
		hostObjects[obj.Host]++
	}

	t.Logf("File browser objects by host:")
	for host, count := range hostObjects {
		t.Logf("  %s: %d objects", host, count)

		// Verify GetFileBrowserObjectsByHost returns the same count
		hostObjs, err := client.GetFileBrowserObjectsByHost(ctx, host)
		if err != nil {
			t.Errorf("GetFileBrowserObjectsByHost(%s) failed: %v", host, err)
			continue
		}

		if len(hostObjs) != count {
			t.Errorf("GetFileBrowserObjectsByHost(%s) returned %d objects, expected %d",
				host, len(hostObjs), count)
		}
	}
}

// TestFileBrowser_MultipleCallbacks tests objects across multiple callbacks.
func TestFileBrowser_MultipleCallbacks(t *testing.T) {
	client := getTestClient(t)
	ctx := context.Background()

	objects, err := client.GetFileBrowserObjects(ctx)
	if err != nil {
		t.Fatalf("GetFileBrowserObjects() failed: %v", err)
	}

	if len(objects) == 0 {
		t.Skip("No file browser objects available")
	}

	// Group objects by callback
	callbackObjects := make(map[int]int)
	for _, obj := range objects {
		callbackObjects[obj.CallbackID]++
	}

	t.Logf("File browser objects by callback:")
	for callbackID, count := range callbackObjects {
		t.Logf("  Callback %d: %d objects", callbackID, count)

		// Verify GetFileBrowserObjectsByCallback returns the same count
		callbackObjs, err := client.GetFileBrowserObjectsByCallback(ctx, callbackID)
		if err != nil {
			t.Errorf("GetFileBrowserObjectsByCallback(%d) failed: %v", callbackID, err)
			continue
		}

		if len(callbackObjs) != count {
			t.Errorf("GetFileBrowserObjectsByCallback(%d) returned %d objects, expected %d",
				callbackID, len(callbackObjs), count)
		}
	}
}
