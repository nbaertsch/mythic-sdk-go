//go:build integration

package integration

import (
	"context"
	"testing"
	"time"

	"github.com/nbaertsch/mythic-sdk-go/pkg/mythic/types"
)

func TestTags_GetTagTypes(t *testing.T) {
	SkipIfNoMythic(t)

	client := AuthenticateTestClient(t)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	tagTypes, err := client.GetTagTypes(ctx)
	if err != nil {
		t.Fatalf("GetTagTypes failed: %v", err)
	}

	if tagTypes == nil {
		t.Fatal("GetTagTypes returned nil")
	}

	t.Logf("Found %d tag type(s)", len(tagTypes))

	// If there are tag types, verify structure
	if len(tagTypes) > 0 {
		tt := tagTypes[0]
		if tt.ID == 0 {
			t.Error("TagType ID should not be 0")
		}
		if tt.Name == "" {
			t.Error("TagType should have a name")
		}
		t.Logf("First tag type: %s", tt.String())
		t.Logf("  - ID: %d", tt.ID)
		t.Logf("  - Name: %s", tt.Name)
		t.Logf("  - Description: %s", tt.Description)
		t.Logf("  - Color: %s", tt.Color)
		t.Logf("  - Deleted: %v", tt.IsDeleted())
	}
}

func TestTags_CreateAndRetrieveTagType(t *testing.T) {
	SkipIfNoMythic(t)

	client := AuthenticateTestClient(t)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Create a test tag type
	tagName := "Test-Tag-" + time.Now().Format("20060102150405")
	description := "Test tag type created by integration test"
	color := "#FF5500"

	req := &types.CreateTagTypeRequest{
		Name:        tagName,
		Description: &description,
		Color:       &color,
	}

	tagType, err := client.CreateTagType(ctx, req)
	if err != nil {
		t.Fatalf("CreateTagType failed: %v", err)
	}

	if tagType == nil {
		t.Fatal("CreateTagType returned nil")
	}

	t.Logf("Created tag type: %s", tagType.String())
	t.Logf("  - ID: %d", tagType.ID)
	t.Logf("  - Name: %s", tagType.Name)
	t.Logf("  - Description: %s", tagType.Description)
	t.Logf("  - Color: %s", tagType.Color)

	// Verify created tag type
	if tagType.Name != tagName {
		t.Errorf("Expected name %q, got %q", tagName, tagType.Name)
	}
	if tagType.Description != description {
		t.Errorf("Expected description %q, got %q", description, tagType.Description)
	}
	if tagType.Color != color {
		t.Errorf("Expected color %q, got %q", color, tagType.Color)
	}

	// Retrieve the tag type by ID
	retrieved, err := client.GetTagTypeByID(ctx, tagType.ID)
	if err != nil {
		t.Fatalf("GetTagTypeByID failed: %v", err)
	}

	if retrieved.ID != tagType.ID {
		t.Errorf("Expected tag type ID %d, got %d", tagType.ID, retrieved.ID)
	}
	if retrieved.Name != tagType.Name {
		t.Errorf("Expected name %q, got %q", tagType.Name, retrieved.Name)
	}

	// Clean up: delete the tag type
	err = client.DeleteTagType(ctx, tagType.ID)
	if err != nil {
		t.Logf("Warning: Failed to delete test tag type: %v", err)
	} else {
		t.Logf("Successfully deleted test tag type")
	}
}

func TestTags_CreateTagType_InvalidInput(t *testing.T) {
	SkipIfNoMythic(t)

	client := AuthenticateTestClient(t)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Test with nil request
	_, err := client.CreateTagType(ctx, nil)
	if err == nil {
		t.Fatal("Expected error for nil request, got nil")
	}
	t.Logf("Nil request error: %v", err)

	// Test with empty name
	_, err = client.CreateTagType(ctx, &types.CreateTagTypeRequest{
		Name: "",
	})
	if err == nil {
		t.Fatal("Expected error for empty name, got nil")
	}
	t.Logf("Empty name error: %v", err)
}

func TestTags_GetTagTypeByID_InvalidID(t *testing.T) {
	SkipIfNoMythic(t)

	client := AuthenticateTestClient(t)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Test with zero ID
	_, err := client.GetTagTypeByID(ctx, 0)
	if err == nil {
		t.Fatal("Expected error for zero tag type ID, got nil")
	}
	t.Logf("Zero ID error: %v", err)

	// Test with non-existent ID
	_, err = client.GetTagTypeByID(ctx, 999999)
	if err == nil {
		t.Fatal("Expected error for non-existent tag type ID, got nil")
	}
	t.Logf("Non-existent ID error: %v", err)
}

func TestTags_UpdateTagType(t *testing.T) {
	SkipIfNoMythic(t)

	client := AuthenticateTestClient(t)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Create a test tag type first
	tagName := "Test-Update-" + time.Now().Format("20060102150405")
	description := "Original description"
	color := "#0000FF"

	tagType, err := client.CreateTagType(ctx, &types.CreateTagTypeRequest{
		Name:        tagName,
		Description: &description,
		Color:       &color,
	})
	if err != nil {
		t.Fatalf("CreateTagType failed: %v", err)
	}

	t.Logf("Created tag type for update test: ID %d", tagType.ID)

	// Update the tag type
	newDescription := "Updated description"
	newColor := "#00FF00"
	updated, err := client.UpdateTagType(ctx, &types.UpdateTagTypeRequest{
		ID:          tagType.ID,
		Description: &newDescription,
		Color:       &newColor,
	})
	if err != nil {
		t.Fatalf("UpdateTagType failed: %v", err)
	}

	if updated.Description != newDescription {
		t.Errorf("Expected description %q, got %q", newDescription, updated.Description)
	}
	if updated.Color != newColor {
		t.Errorf("Expected color %q, got %q", newColor, updated.Color)
	}

	t.Logf("Successfully updated tag type")
	t.Logf("  - Old description: %s", description)
	t.Logf("  - New description: %s", updated.Description)
	t.Logf("  - Old color: %s", color)
	t.Logf("  - New color: %s", updated.Color)

	// Clean up
	err = client.DeleteTagType(ctx, tagType.ID)
	if err != nil {
		t.Logf("Warning: Failed to delete test tag type: %v", err)
	}
}

func TestTags_DeleteTagType(t *testing.T) {
	SkipIfNoMythic(t)

	client := AuthenticateTestClient(t)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Create a test tag type
	tagName := "Test-Delete-" + time.Now().Format("20060102150405")

	tagType, err := client.CreateTagType(ctx, &types.CreateTagTypeRequest{
		Name: tagName,
	})
	if err != nil {
		t.Fatalf("CreateTagType failed: %v", err)
	}

	t.Logf("Created tag type for delete test: ID %d", tagType.ID)

	// Delete the tag type
	err = client.DeleteTagType(ctx, tagType.ID)
	if err != nil {
		t.Fatalf("DeleteTagType failed: %v", err)
	}

	t.Log("Successfully deleted tag type")

	// Verify deletion - tag type should not be found or marked as deleted
	deleted, err := client.GetTagTypeByID(ctx, tagType.ID)
	if err != nil {
		// If it's not found, that's acceptable
		t.Logf("Tag type not found after deletion (expected): %v", err)
		return
	}

	// If we can still retrieve it, verify it's marked as deleted
	if deleted != nil && !deleted.IsDeleted() {
		t.Error("Tag type should be marked as deleted")
	} else if deleted != nil {
		t.Logf("Tag type marked as deleted: %v", deleted.IsDeleted())
	}
}

func TestTags_CreateTag(t *testing.T) {
	SkipIfNoMythic(t)

	client := AuthenticateTestClient(t)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// First, create a tag type
	tagName := "Test-Tag-" + time.Now().Format("20060102150405")
	tagType, err := client.CreateTagType(ctx, &types.CreateTagTypeRequest{
		Name: tagName,
	})
	if err != nil {
		t.Fatalf("CreateTagType failed: %v", err)
	}
	defer client.DeleteTagType(ctx, tagType.ID)

	// Get a task to tag (create one if needed or get existing)
	tasks, err := client.GetTasksByStatus(ctx, types.TaskStatusCompleted)
	if err != nil {
		t.Fatalf("GetTasksByStatus failed: %v", err)
	}

	if len(tasks) == 0 {
		t.Skip("No tasks available to tag")
	}

	taskID := tasks[0].ID

	// Create a tag on the task
	tag, err := client.CreateTag(ctx, &types.CreateTagRequest{
		TagTypeID:  tagType.ID,
		SourceType: types.TagSourceTask,
		SourceID:   taskID,
	})
	if err != nil {
		t.Fatalf("CreateTag failed: %v", err)
	}

	if tag == nil {
		t.Fatal("CreateTag returned nil")
	}

	t.Logf("Created tag: %s", tag.String())
	t.Logf("  - ID: %d", tag.ID)
	t.Logf("  - TagTypeID: %d", tag.TagTypeID)
	t.Logf("  - SourceType: %s", tag.SourceType)
	t.Logf("  - SourceID: %d", tag.SourceID)

	// Verify tag
	if tag.TagTypeID != tagType.ID {
		t.Errorf("Expected TagTypeID %d, got %d", tagType.ID, tag.TagTypeID)
	}
	if tag.SourceType != types.TagSourceTask {
		t.Errorf("Expected SourceType %q, got %q", types.TagSourceTask, tag.SourceType)
	}
	if tag.SourceID != taskID {
		t.Errorf("Expected SourceID %d, got %d", taskID, tag.SourceID)
	}

	// Clean up: delete the tag
	err = client.DeleteTag(ctx, tag.ID)
	if err != nil {
		t.Logf("Warning: Failed to delete test tag: %v", err)
	} else {
		t.Logf("Successfully deleted test tag")
	}
}

func TestTags_CreateTag_InvalidInput(t *testing.T) {
	SkipIfNoMythic(t)

	client := AuthenticateTestClient(t)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Test with nil request
	_, err := client.CreateTag(ctx, nil)
	if err == nil {
		t.Fatal("Expected error for nil request, got nil")
	}
	t.Logf("Nil request error: %v", err)

	// Test with zero tag type ID
	_, err = client.CreateTag(ctx, &types.CreateTagRequest{
		TagTypeID:  0,
		SourceType: types.TagSourceTask,
		SourceID:   1,
	})
	if err == nil {
		t.Fatal("Expected error for zero tag type ID, got nil")
	}
	t.Logf("Zero TagTypeID error: %v", err)

	// Test with empty source type
	_, err = client.CreateTag(ctx, &types.CreateTagRequest{
		TagTypeID:  1,
		SourceType: "",
		SourceID:   1,
	})
	if err == nil {
		t.Fatal("Expected error for empty source type, got nil")
	}
	t.Logf("Empty SourceType error: %v", err)

	// Test with zero source ID
	_, err = client.CreateTag(ctx, &types.CreateTagRequest{
		TagTypeID:  1,
		SourceType: types.TagSourceTask,
		SourceID:   0,
	})
	if err == nil {
		t.Fatal("Expected error for zero source ID, got nil")
	}
	t.Logf("Zero SourceID error: %v", err)
}

func TestTags_GetTags(t *testing.T) {
	SkipIfNoMythic(t)

	client := AuthenticateTestClient(t)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// First, create a tag type
	tagName := "Test-GetTags-" + time.Now().Format("20060102150405")
	tagType, err := client.CreateTagType(ctx, &types.CreateTagTypeRequest{
		Name: tagName,
	})
	if err != nil {
		t.Fatalf("CreateTagType failed: %v", err)
	}
	defer client.DeleteTagType(ctx, tagType.ID)

	// Get a task to tag
	tasks, err := client.GetTasksByStatus(ctx, types.TaskStatusCompleted)
	if err != nil {
		t.Fatalf("GetTasksByStatus failed: %v", err)
	}

	if len(tasks) == 0 {
		t.Skip("No tasks available to tag")
	}

	taskID := tasks[0].ID

	// Create multiple tags on the task
	tag1, err := client.CreateTag(ctx, &types.CreateTagRequest{
		TagTypeID:  tagType.ID,
		SourceType: types.TagSourceTask,
		SourceID:   taskID,
	})
	if err != nil {
		t.Fatalf("CreateTag 1 failed: %v", err)
	}
	defer client.DeleteTag(ctx, tag1.ID)

	t.Logf("Created tag on task %d", taskID)

	// Get tags for the task
	tags, err := client.GetTags(ctx, types.TagSourceTask, taskID)
	if err != nil {
		t.Fatalf("GetTags failed: %v", err)
	}

	if len(tags) == 0 {
		t.Error("Expected at least one tag")
	}

	// Verify all tags belong to the correct object
	for _, tag := range tags {
		if tag.SourceType != types.TagSourceTask {
			t.Errorf("Expected SourceType %q, got %q", types.TagSourceTask, tag.SourceType)
		}
		if tag.SourceID != taskID {
			t.Errorf("Expected SourceID %d, got %d", taskID, tag.SourceID)
		}
	}

	t.Logf("Found %d tag(s) on task %d", len(tags), taskID)
}

func TestTags_GetTagsByOperation(t *testing.T) {
	SkipIfNoMythic(t)

	client := AuthenticateTestClient(t)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Get current operation
	currentOpID := client.GetCurrentOperation()
	if currentOpID == nil {
		t.Skip("No current operation set")
	}

	tags, err := client.GetTagsByOperation(ctx, *currentOpID)
	if err != nil {
		t.Fatalf("GetTagsByOperation failed: %v", err)
	}

	if tags == nil {
		t.Fatal("GetTagsByOperation returned nil")
	}

	t.Logf("Found %d tag(s) for operation %d", len(tags), *currentOpID)

	// Verify all tags belong to the operation
	for _, tag := range tags {
		if tag.OperationID != *currentOpID {
			t.Errorf("Expected operation ID %d, got %d", *currentOpID, tag.OperationID)
		}
	}

	// Verify tags are sorted by timestamp (descending)
	if len(tags) > 1 {
		for i := 1; i < len(tags); i++ {
			if tags[i].Timestamp.After(tags[i-1].Timestamp) {
				t.Error("Tags should be sorted by timestamp (descending/newest first)")
				break
			}
		}
	}
}

func TestTags_DeleteTag(t *testing.T) {
	SkipIfNoMythic(t)

	client := AuthenticateTestClient(t)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// First, create a tag type
	tagName := "Test-DeleteTag-" + time.Now().Format("20060102150405")
	tagType, err := client.CreateTagType(ctx, &types.CreateTagTypeRequest{
		Name: tagName,
	})
	if err != nil {
		t.Fatalf("CreateTagType failed: %v", err)
	}
	defer client.DeleteTagType(ctx, tagType.ID)

	// Get a task to tag
	tasks, err := client.GetTasksByStatus(ctx, types.TaskStatusCompleted)
	if err != nil {
		t.Fatalf("GetTasksByStatus failed: %v", err)
	}

	if len(tasks) == 0 {
		t.Skip("No tasks available to tag")
	}

	taskID := tasks[0].ID

	// Create a tag
	tag, err := client.CreateTag(ctx, &types.CreateTagRequest{
		TagTypeID:  tagType.ID,
		SourceType: types.TagSourceTask,
		SourceID:   taskID,
	})
	if err != nil {
		t.Fatalf("CreateTag failed: %v", err)
	}

	t.Logf("Created tag for delete test: ID %d", tag.ID)

	// Delete the tag
	err = client.DeleteTag(ctx, tag.ID)
	if err != nil {
		t.Fatalf("DeleteTag failed: %v", err)
	}

	t.Log("Successfully deleted tag")

	// Verify deletion - tag should not be found
	_, err = client.GetTagByID(ctx, tag.ID)
	if err == nil {
		t.Error("Expected error when getting deleted tag")
	} else {
		t.Logf("Tag not found after deletion (expected): %v", err)
	}
}

func TestTags_TagSourceTypes(t *testing.T) {
	SkipIfNoMythic(t)

	client := AuthenticateTestClient(t)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Test that all source type constants are valid strings
	sourceTypes := []string{
		types.TagSourceTask,
		types.TagSourceCallback,
		types.TagSourceFile,
		types.TagSourcePayload,
		types.TagSourceArtifact,
		types.TagSourceProcess,
		types.TagSourceKeylog,
	}

	for _, sourceType := range sourceTypes {
		if sourceType == "" {
			t.Errorf("Source type constant should not be empty")
		}
		t.Logf("Source type: %s", sourceType)
	}
}

func TestTags_TimestampOrdering(t *testing.T) {
	SkipIfNoMythic(t)

	client := AuthenticateTestClient(t)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Get current operation
	currentOpID := client.GetCurrentOperation()
	if currentOpID == nil {
		t.Skip("No current operation set")
	}

	tags, err := client.GetTagsByOperation(ctx, *currentOpID)
	if err != nil {
		t.Fatalf("GetTagsByOperation failed: %v", err)
	}

	if len(tags) < 2 {
		t.Skip("Need at least 2 tags to test timestamp ordering")
	}

	// Verify descending order (newest first)
	for i := 1; i < len(tags); i++ {
		if tags[i].Timestamp.After(tags[i-1].Timestamp) {
			t.Errorf("Timestamp ordering broken at index %d: %s > %s",
				i,
				tags[i].Timestamp.Format("2006-01-02 15:04:05"),
				tags[i-1].Timestamp.Format("2006-01-02 15:04:05"))
		}
	}

	t.Log("Timestamp ordering verified (newest first)")
	t.Logf("  - Newest: %s", tags[0].Timestamp.Format("2006-01-02 15:04:05"))
	t.Logf("  - Oldest: %s", tags[len(tags)-1].Timestamp.Format("2006-01-02 15:04:05"))
}
