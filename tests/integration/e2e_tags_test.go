//go:build integration

package integration

import (
	"context"
	"testing"
	"time"

	"github.com/nbaertsch/mythic-sdk-go/pkg/mythic/types"
)

// TestE2E_TagManagement tests the complete tag and tag type management workflow
// Covers: CreateTag, GetTags, GetTagsByOperation, DeleteTag,
// CreateTagType, GetTagTypes, GetTagTypesByOperation, UpdateTagType, DeleteTagType
func TestE2E_TagManagement(t *testing.T) {
	client := AuthenticateTestClient(t)

	var createdTagTypeIDs []int
	var createdTagIDs []int
	var testArtifactID int

	// Register cleanup
	defer func() {
		// Delete tags first (they depend on tag types)
		for _, tagID := range createdTagIDs {
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			_ = client.DeleteTag(ctx, tagID)
			cancel()
			t.Logf("Cleaned up tag ID: %d", tagID)
		}
		// Then delete tag types
		for _, tagTypeID := range createdTagTypeIDs {
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			_ = client.DeleteTagType(ctx, tagTypeID)
			cancel()
			t.Logf("Cleaned up tag type ID: %d", tagTypeID)
		}
		// Delete test artifact
		if testArtifactID > 0 {
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			_ = client.DeleteArtifact(ctx, testArtifactID)
			cancel()
			t.Logf("Cleaned up artifact ID: %d", testArtifactID)
		}
	}()

	// Prerequisite: Create a test artifact to tag
	t.Log("=== Prerequisite: Create test artifact for tagging ===")
	ctx0, cancel0 := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel0()

	testHost := "tag-test-host"
	artifactType := types.ArtifactTypeFile
	testArtifact := &types.CreateArtifactRequest{
		Artifact:     "/tmp/test_tagged_file.txt",
		Host:         &testHost,
		ArtifactType: &artifactType,
	}

	createdArtifact, err := client.CreateArtifact(ctx0, testArtifact)
	if err != nil {
		t.Fatalf("Failed to create test artifact: %v", err)
	}
	testArtifactID = createdArtifact.ID
	t.Logf("✓ Test artifact created: ID %d", testArtifactID)

	// Test 1: Get tag types baseline
	t.Log("=== Test 1: Get tag types baseline ===")
	ctx1, cancel1 := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel1()

	baselineTagTypes, err := client.GetTagTypes(ctx1)
	if err != nil {
		t.Fatalf("GetTagTypes baseline failed: %v", err)
	}
	baselineTagTypeCount := len(baselineTagTypes)
	t.Logf("✓ Baseline tag type count: %d", baselineTagTypeCount)

	// Test 2: Create tag type - "Priority"
	t.Log("=== Test 2: Create tag type - Priority ===")
	ctx2, cancel2 := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel2()

	priorityDesc := "Operational priority level"
	priorityColor := "#FF0000"
	priorityTagType := &types.CreateTagTypeRequest{
		Name:        "E2E-Priority",
		Description: &priorityDesc,
		Color:       &priorityColor,
	}

	createdPriority, err := client.CreateTagType(ctx2, priorityTagType)
	if err != nil {
		t.Fatalf("CreateTagType (Priority) failed: %v", err)
	}
	if createdPriority.ID == 0 {
		t.Fatal("Created tag type has ID 0")
	}
	createdTagTypeIDs = append(createdTagTypeIDs, createdPriority.ID)
	t.Logf("✓ Priority tag type created: ID %d, Name %s", createdPriority.ID, createdPriority.Name)

	// Test 3: Create tag type - "Target Type"
	t.Log("=== Test 3: Create tag type - Target Type ===")
	ctx3, cancel3 := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel3()

	targetDesc := "Type of target system"
	targetColor := "#00FF00"
	targetTagType := &types.CreateTagTypeRequest{
		Name:        "E2E-TargetType",
		Description: &targetDesc,
		Color:       &targetColor,
	}

	createdTarget, err := client.CreateTagType(ctx3, targetTagType)
	if err != nil {
		t.Fatalf("CreateTagType (Target) failed: %v", err)
	}
	createdTagTypeIDs = append(createdTagTypeIDs, createdTarget.ID)
	t.Logf("✓ Target Type tag type created: ID %d, Name %s", createdTarget.ID, createdTarget.Name)

	// Test 4: Get all tag types after creation
	t.Log("=== Test 4: Get all tag types after creation ===")
	ctx4, cancel4 := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel4()

	allTagTypes, err := client.GetTagTypes(ctx4)
	if err != nil {
		t.Fatalf("GetTagTypes after creation failed: %v", err)
	}
	newTagTypeCount := len(allTagTypes)
	if newTagTypeCount < baselineTagTypeCount+2 {
		t.Errorf("Expected at least %d tag types, got %d", baselineTagTypeCount+2, newTagTypeCount)
	}
	t.Logf("✓ Total tag types now: %d (added %d)", newTagTypeCount, newTagTypeCount-baselineTagTypeCount)

	// Test 5: Get tag types by operation
	t.Log("=== Test 5: Get tag types by operation ===")
	ctx5, cancel5 := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel5()

	operationID := client.GetCurrentOperation()
	if operationID == nil {
		t.Fatal("No current operation set")
	}

	opTagTypes, err := client.GetTagTypesByOperation(ctx5, *operationID)
	if err != nil {
		t.Fatalf("GetTagTypesByOperation failed: %v", err)
	}
	t.Logf("✓ Found %d tag types for operation %d", len(opTagTypes), *operationID)

	// Verify our created tag types are in the operation
	found := 0
	for _, tagType := range opTagTypes {
		for _, createdID := range createdTagTypeIDs {
			if tagType.ID == createdID {
				found++
				t.Logf("  ✓ Found tag type %d: %s", tagType.ID, tagType.Name)
			}
		}
	}
	if found != len(createdTagTypeIDs) {
		t.Errorf("Expected to find %d tag types in operation, found %d", len(createdTagTypeIDs), found)
	}

	// Test 6: Update tag type (add/modify description)
	t.Log("=== Test 6: Update tag type (modify description) ===")
	ctx6, cancel6 := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel6()

	newDesc := "Updated: Indicates operational priority for tasks and targets"
	updateReq := &types.UpdateTagTypeRequest{
		ID:          createdPriority.ID,
		Description: &newDesc,
	}

	updatedTagType, err := client.UpdateTagType(ctx6, updateReq)
	if err != nil {
		t.Fatalf("UpdateTagType failed: %v", err)
	}
	if updatedTagType.Description != newDesc {
		t.Errorf("Description not updated: expected %s, got %s", newDesc, updatedTagType.Description)
	}
	t.Logf("✓ Tag type %d updated: description = %s", updatedTagType.ID, updatedTagType.Description)

	// Test 7: Create tags with Priority type
	t.Log("=== Test 7: Create tags - Priority levels ===")
	ctx7, cancel7 := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel7()

	// Tag the artifact as "High Priority"
	highPriorityTag := &types.CreateTagRequest{
		TagTypeID:  createdPriority.ID,
		SourceType: types.TagSourceArtifact,
		SourceID:   testArtifactID,
	}

	createdHighTag, err := client.CreateTag(ctx7, highPriorityTag)
	if err != nil {
		t.Fatalf("CreateTag (Priority on artifact) failed: %v", err)
	}
	if createdHighTag.ID == 0 {
		t.Fatal("Created tag has ID 0")
	}
	createdTagIDs = append(createdTagIDs, createdHighTag.ID)
	t.Logf("✓ Tag created: ID %d (Priority → Artifact %d)", createdHighTag.ID, testArtifactID)

	// Test 8: Create tag with Target Type
	t.Log("=== Test 8: Create tag - Target Type ===")
	ctx8, cancel8 := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel8()

	targetTypeTag := &types.CreateTagRequest{
		TagTypeID:  createdTarget.ID,
		SourceType: types.TagSourceArtifact,
		SourceID:   testArtifactID,
	}

	createdTargetTag, err := client.CreateTag(ctx8, targetTypeTag)
	if err != nil {
		t.Fatalf("CreateTag (Target Type on artifact) failed: %v", err)
	}
	createdTagIDs = append(createdTagIDs, createdTargetTag.ID)
	t.Logf("✓ Tag created: ID %d (TargetType → Artifact %d)", createdTargetTag.ID, testArtifactID)

	// Test 9: Get tags for the artifact
	t.Log("=== Test 9: Get tags for artifact ===")
	ctx9, cancel9 := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel9()

	artifactTags, err := client.GetTags(ctx9, types.TagSourceArtifact, testArtifactID)
	if err != nil {
		t.Fatalf("GetTags failed: %v", err)
	}
	if len(artifactTags) < 2 {
		t.Errorf("Expected at least 2 tags on artifact, got %d", len(artifactTags))
	}
	t.Logf("✓ Found %d tags on artifact %d", len(artifactTags), testArtifactID)

	// Verify our tags are present
	foundTags := 0
	for _, tag := range artifactTags {
		for _, createdID := range createdTagIDs {
			if tag.ID == createdID {
				foundTags++
				tagTypeName := "unknown"
				if tag.TagType != nil {
					tagTypeName = tag.TagType.Name
				}
				t.Logf("  ✓ Found tag %d: Type %s", tag.ID, tagTypeName)
			}
		}
	}
	if foundTags != len(createdTagIDs) {
		t.Errorf("Expected to find %d tags on artifact, found %d", len(createdTagIDs), foundTags)
	}

	// Test 10: Get tags by operation
	t.Log("=== Test 10: Get tags by operation ===")
	ctx10, cancel10 := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel10()

	opTags, err := client.GetTagsByOperation(ctx10, *operationID)
	if err != nil {
		t.Fatalf("GetTagsByOperation failed: %v", err)
	}
	t.Logf("✓ Found %d tags for operation %d", len(opTags), *operationID)

	// Verify our tags are in the operation
	foundOpTags := 0
	for _, tag := range opTags {
		for _, createdID := range createdTagIDs {
			if tag.ID == createdID {
				foundOpTags++
			}
		}
	}
	if foundOpTags != len(createdTagIDs) {
		t.Errorf("Expected to find %d tags in operation, found %d", len(createdTagIDs), foundOpTags)
	}

	// Test 11: Delete tags
	t.Log("=== Test 11: Delete tags ===")
	for _, tagID := range createdTagIDs {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		err := client.DeleteTag(ctx, tagID)
		cancel()
		if err != nil {
			t.Errorf("DeleteTag failed for ID %d: %v", tagID, err)
		} else {
			t.Logf("✓ Tag %d deleted", tagID)
		}
	}

	// Test 12: Delete tag types
	t.Log("=== Test 12: Delete tag types ===")
	for _, tagTypeID := range createdTagTypeIDs {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		err := client.DeleteTagType(ctx, tagTypeID)
		cancel()
		if err != nil {
			t.Errorf("DeleteTagType failed for ID %d: %v", tagTypeID, err)
		} else {
			t.Logf("✓ Tag type %d deleted", tagTypeID)
		}
	}

	// Test 13: Verify deletion
	t.Log("=== Test 13: Verify deletion ===")
	ctx11, cancel11 := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel11()

	finalTagTypes, err := client.GetTagTypes(ctx11)
	if err != nil {
		t.Fatalf("GetTagTypes after delete failed: %v", err)
	}

	// Check that deleted tag types are not in the active list (or marked deleted)
	for _, tagType := range finalTagTypes {
		for _, deletedID := range createdTagTypeIDs {
			if tagType.ID == deletedID && !tagType.Deleted {
				t.Errorf("Tag type %d still active after deletion", deletedID)
			}
		}
	}
	t.Logf("✓ Verified deletion of %d tags and %d tag types", len(createdTagIDs), len(createdTagTypeIDs))

	t.Log("=== ✓ All tag management tests passed ===")
}

// TestE2E_TagsErrorHandling tests error scenarios for tag operations
func TestE2E_TagsErrorHandling(t *testing.T) {
	client := AuthenticateTestClient(t)

	// Test 1: Create tag with non-existent tag type
	t.Log("=== Test 1: Create tag with non-existent tag type ===")
	ctx1, cancel1 := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel1()

	invalidTag := &types.CreateTagRequest{
		TagTypeID:  999999,
		SourceType: types.TagSourceArtifact,
		SourceID:   1,
	}

	_, err := client.CreateTag(ctx1, invalidTag)
	if err == nil {
		t.Error("Expected error for non-existent tag type")
	}
	t.Logf("✓ Invalid tag type rejected: %v", err)

	// Test 2: Delete non-existent tag
	t.Log("=== Test 2: Delete non-existent tag ===")
	ctx2, cancel2 := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel2()

	err = client.DeleteTag(ctx2, 999999)
	if err == nil {
		t.Error("Expected error for non-existent tag delete")
	}
	t.Logf("✓ Non-existent tag delete rejected: %v", err)

	// Test 3: Delete non-existent tag type
	t.Log("=== Test 3: Delete non-existent tag type ===")
	ctx3, cancel3 := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel3()

	err = client.DeleteTagType(ctx3, 999999)
	if err == nil {
		t.Error("Expected error for non-existent tag type delete")
	}
	t.Logf("✓ Non-existent tag type delete rejected: %v", err)

	// Test 4: Update non-existent tag type
	t.Log("=== Test 4: Update non-existent tag type ===")
	ctx4, cancel4 := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel4()

	desc := "test"
	updateReq := &types.UpdateTagTypeRequest{
		ID:          999999,
		Description: &desc,
	}

	_, err = client.UpdateTagType(ctx4, updateReq)
	if err == nil {
		t.Error("Expected error for non-existent tag type update")
	}
	t.Logf("✓ Non-existent tag type update rejected: %v", err)

	t.Log("=== ✓ All error handling tests passed ===")
}
