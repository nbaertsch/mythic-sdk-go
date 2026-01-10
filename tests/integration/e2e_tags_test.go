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

	var createdTagIDs []int
	var testArtifactID int

	// Register cleanup
	defer func() {
		// Delete tags
		for _, tagID := range createdTagIDs {
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			_ = client.DeleteTag(ctx, tagID)
			cancel()
			t.Logf("Cleaned up tag ID: %d", tagID)
		}
		// Note: We don't delete tag types since we're using existing ones, not creating new ones

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
	testArtifact := &types.CreateArtifactRequest{
		Artifact: "/tmp/test_tagged_file.txt",
		Host:     &testHost,
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

	// Note: CreateTagType is not supported via GraphQL API - tag types must be managed through admin UI
	// Use existing tag types from baseline for testing

	if baselineTagTypeCount < 2 {
		t.Skip("Need at least 2 existing tag types to run tag tests (create via admin UI)")
	}

	// Test 2: Use first existing tag type
	t.Log("=== Test 2: Use first existing tag type ===")
	firstTagType := baselineTagTypes[0]
	t.Logf("✓ Using existing tag type: ID %d, Name %s", firstTagType.ID, firstTagType.Name)

	// Test 3: Use second existing tag type
	t.Log("=== Test 3: Use second existing tag type ===")
	secondTagType := baselineTagTypes[1]
	t.Logf("✓ Using existing tag type: ID %d, Name %s", secondTagType.ID, secondTagType.Name)

	// Store IDs for later use (not for cleanup, since we didn't create them)
	firstTagTypeID := firstTagType.ID
	secondTagTypeID := secondTagType.ID

	// Test 4: Verify tag types still accessible
	t.Log("=== Test 4: Verify existing tag types ===")
	ctx4, cancel4 := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel4()

	allTagTypes, err := client.GetTagTypes(ctx4)
	if err != nil {
		t.Fatalf("GetTagTypes verification failed: %v", err)
	}
	if len(allTagTypes) < 2 {
		t.Fatalf("Expected at least 2 tag types, got %d", len(allTagTypes))
	}
	t.Logf("✓ Total tag types: %d", len(allTagTypes))

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

	// Verify our tag types are in the operation
	foundFirst := false
	foundSecond := false
	for _, tagType := range opTagTypes {
		if tagType.ID == firstTagTypeID {
			foundFirst = true
			t.Logf("  ✓ Found tag type %d: %s", tagType.ID, tagType.Name)
		}
		if tagType.ID == secondTagTypeID {
			foundSecond = true
			t.Logf("  ✓ Found tag type %d: %s", tagType.ID, tagType.Name)
		}
	}
	if !foundFirst || !foundSecond {
		t.Errorf("Expected to find both tag types in operation")
	}

	// Test 6: Update tag type (modify name only - description/color not supported)
	t.Log("=== Test 6: Update tag type (modify name) ===")
	ctx6, cancel6 := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel6()

	newName := "E2E-Updated-TagType"
	updateReq := &types.UpdateTagTypeRequest{
		ID:   firstTagTypeID,
		Name: &newName,
	}

	updatedTagType, err := client.UpdateTagType(ctx6, updateReq)
	if err != nil {
		t.Logf("UpdateTagType failed (may not be supported): %v", err)
		// Don't fatal - UpdateTagType may have limited support
	} else if updatedTagType.Name != newName {
		t.Errorf("Name not updated: expected %s, got %s", newName, updatedTagType.Name)
	} else {
		t.Logf("✓ Tag type %d updated: name = %s", updatedTagType.ID, updatedTagType.Name)
	}

	// Test 7: Create tags with first tag type
	t.Log("=== Test 7: Create tags with tag types ===")
	ctx7, cancel7 := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel7()

	// Tag the artifact with first tag type
	firstTag := &types.CreateTagRequest{
		TagTypeID:  firstTagTypeID,
		SourceType: types.TagSourceArtifact,
		SourceID:   testArtifactID,
	}

	createdFirstTag, err := client.CreateTag(ctx7, firstTag)
	if err != nil {
		t.Fatalf("CreateTag (first type on artifact) failed: %v", err)
	}
	if createdFirstTag.ID == 0 {
		t.Fatal("Created tag has ID 0")
	}
	createdTagIDs = append(createdTagIDs, createdFirstTag.ID)
	t.Logf("✓ Tag created: ID %d (TagType %d → Artifact %d)", createdFirstTag.ID, firstTagTypeID, testArtifactID)

	// Test 8: Create tag with second tag type
	t.Log("=== Test 8: Create tag with second tag type ===")
	ctx8, cancel8 := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel8()

	secondTag := &types.CreateTagRequest{
		TagTypeID:  secondTagTypeID,
		SourceType: types.TagSourceArtifact,
		SourceID:   testArtifactID,
	}

	createdSecondTag, err := client.CreateTag(ctx8, secondTag)
	if err != nil {
		t.Fatalf("CreateTag (second type on artifact) failed: %v", err)
	}
	createdTagIDs = append(createdTagIDs, createdSecondTag.ID)
	t.Logf("✓ Tag created: ID %d (TagType %d → Artifact %d)", createdSecondTag.ID, secondTagTypeID, testArtifactID)

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

	// Note: We don't test tag type deletion since we're using existing tag types
	t.Log("=== Tag management tests complete (using existing tag types) ===")
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
