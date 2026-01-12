package unit

import (
	"testing"

	"github.com/nbaertsch/mythic-sdk-go/pkg/mythic/types"
)

// TestTagTypeString tests the TagType.String() method
func TestTagTypeString(t *testing.T) {

	tests := []struct {
		name     string
		tagType  types.TagType
		contains []string
	}{
		{
			name: "with color",
			tagType: types.TagType{
				ID:        1,
				Name:      "Critical",
				Color:     "#FF0000",
			},
			contains: []string{"Critical", "#FF0000"},
		},
		{
			name: "without color",
			tagType: types.TagType{
				ID:        2,
				Name:      "Low Priority",
			},
			contains: []string{"Low Priority"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.tagType.String()
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

// TestTagTypeIsDeleted tests the TagType.IsDeleted() method
func TestTagTypeIsDeleted(t *testing.T) {
	tests := []struct {
		name     string
		deleted  bool
		expected bool
	}{
		{"deleted tag type", true, true},
		{"active tag type", false, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tagType := types.TagType{Deleted: tt.deleted}
			result := tagType.IsDeleted()
			if result != tt.expected {
				t.Errorf("IsDeleted() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

// TestTagString tests the Tag.String() method
func TestTagString(t *testing.T) {

	tests := []struct {
		name     string
		tag      types.Tag
		contains []string
	}{
		{
			name: "with tag type",
			tag: types.Tag{
				ID:         1,
				TagTypeID:  5,
				SourceType: types.TagSourceTask,
				SourceID:   42,
				TagType: &types.TagType{
					Name: "Important",
				},
			},
			contains: []string{"Important", "task", "42"},
		},
		{
			name: "without tag type loaded",
			tag: types.Tag{
				ID:         2,
				TagTypeID:  10,
				SourceType: types.TagSourceCallback,
				SourceID:   100,
			},
			contains: []string{"10", "callback", "100"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.tag.String()
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

// TestTagTypeTypes tests the TagType type structure
func TestTagTypeTypes(t *testing.T) {

	tagType := types.TagType{
		ID:          1,
		Name:        "High Priority",
		Description: "Tasks that need immediate attention",
		Color:       "#FF5500",
		OperationID: 5,
		Deleted:     false,
	}

	if tagType.ID != 1 {
		t.Errorf("Expected ID 1, got %d", tagType.ID)
	}
	if tagType.Name != "High Priority" {
		t.Errorf("Expected Name 'High Priority', got %q", tagType.Name)
	}
	if tagType.IsDeleted() {
		t.Error("Expected tag type to not be deleted")
	}
	if tagType.Color != "#FF5500" {
		t.Errorf("Expected Color '#FF5500', got %q", tagType.Color)
	}
}

// TestTagTypes tests the Tag type structure
func TestTagTypes(t *testing.T) {

	tag := types.Tag{
		ID:          1,
		TagTypeID:   5,
		SourceType:  types.TagSourceTask,
		SourceID:    42,
		OperationID: 3,
	}

	if tag.ID != 1 {
		t.Errorf("Expected ID 1, got %d", tag.ID)
	}
	if tag.TagTypeID != 5 {
		t.Errorf("Expected TagTypeID 5, got %d", tag.TagTypeID)
	}
	if tag.SourceType != types.TagSourceTask {
		t.Errorf("Expected SourceType 'task', got %q", tag.SourceType)
	}
	if tag.SourceID != 42 {
		t.Errorf("Expected SourceID 42, got %d", tag.SourceID)
	}
}

// TestCreateTagTypeRequest tests CreateTagTypeRequest structure
func TestCreateTagTypeRequest(t *testing.T) {
	description := "Test tag type"
	color := "#0000FF"

	req := types.CreateTagTypeRequest{
		Name:        "Test Tag",
		Description: &description,
		Color:       &color,
	}

	if req.Name != "Test Tag" {
		t.Errorf("Expected Name 'Test Tag', got %q", req.Name)
	}
	if req.Description == nil || *req.Description != description {
		t.Error("Expected Description to be 'Test tag type'")
	}
	if req.Color == nil || *req.Color != color {
		t.Error("Expected Color to be '#0000FF'")
	}
}

// TestUpdateTagTypeRequest tests UpdateTagTypeRequest structure
func TestUpdateTagTypeRequest(t *testing.T) {
	name := "Updated Name"
	description := "Updated description"
	color := "#00FF00"
	deleted := true

	req := types.UpdateTagTypeRequest{
		ID:          5,
		Name:        &name,
		Description: &description,
		Color:       &color,
		Deleted:     &deleted,
	}

	if req.ID != 5 {
		t.Errorf("Expected ID 5, got %d", req.ID)
	}
	if req.Name == nil || *req.Name != name {
		t.Error("Expected Name to be 'Updated Name'")
	}
	if req.Description == nil || *req.Description != description {
		t.Error("Expected Description to be updated")
	}
	if req.Color == nil || *req.Color != color {
		t.Error("Expected Color to be updated")
	}
	if req.Deleted == nil || !*req.Deleted {
		t.Error("Expected Deleted to be true")
	}
}

// TestCreateTagRequest tests CreateTagRequest structure
func TestCreateTagRequest(t *testing.T) {
	req := types.CreateTagRequest{
		TagTypeID:  5,
		SourceType: types.TagSourceTask,
		SourceID:   42,
	}

	if req.TagTypeID != 5 {
		t.Errorf("Expected TagTypeID 5, got %d", req.TagTypeID)
	}
	if req.SourceType != types.TagSourceTask {
		t.Errorf("Expected SourceType 'task', got %q", req.SourceType)
	}
	if req.SourceID != 42 {
		t.Errorf("Expected SourceID 42, got %d", req.SourceID)
	}
}

// TestTagSourceConstants tests tag source constants
func TestTagSourceConstants(t *testing.T) {
	sources := map[string]string{
		"task":     types.TagSourceTask,
		"callback": types.TagSourceCallback,
		"filemeta": types.TagSourceFile,
		"payload":  types.TagSourcePayload,
		"artifact": types.TagSourceArtifact,
		"process":  types.TagSourceProcess,
		"keylog":   types.TagSourceKeylog,
	}

	for expected, actual := range sources {
		if actual != expected {
			t.Errorf("Expected source %q, got %q", expected, actual)
		}
	}
}

// TestTagTypeWithoutOptionalFields tests TagType without optional fields
func TestTagTypeWithoutOptionalFields(t *testing.T) {
	tagType := types.TagType{
		ID:        1,
		Name:      "minimal",
	}

	if tagType.Description != "" {
		t.Error("Description should be empty")
	}
	if tagType.Color != "" {
		t.Error("Color should be empty")
	}
	if tagType.Operation != nil {
		t.Error("Operation should be nil")
	}

	str := tagType.String()
	if str == "" {
		t.Error("String() should not return empty string even without optional fields")
	}
	if str != "minimal" {
		t.Errorf("Expected String() to return 'minimal', got %q", str)
	}
}

// TestTagWithoutOptionalFields tests Tag without optional fields
func TestTagWithoutOptionalFields(t *testing.T) {
	tag := types.Tag{
		ID:         1,
		TagTypeID:  5,
		SourceType: types.TagSourceCallback,
		SourceID:   10,
	}

	if tag.TagType != nil {
		t.Error("TagType should be nil")
	}
	if tag.Operation != nil {
		t.Error("Operation should be nil")
	}

	str := tag.String()
	if str == "" {
		t.Error("String() should not return empty string even without optional fields")
	}
}

// TestTagTypeColors tests various color formats
func TestTagTypeColors(t *testing.T) {
	colors := []string{
		"#FF0000", // Red
		"#00FF00", // Green
		"#0000FF", // Blue
		"#FFFF00", // Yellow
		"#FF00FF", // Magenta
		"#00FFFF", // Cyan
		"#000000", // Black
		"#FFFFFF", // White
		"#808080", // Gray
	}

	for _, color := range colors {
		tagType := types.TagType{
			ID:        1,
			Name:      "Test",
			Color:     color,
		}

		if tagType.Color != color {
			t.Errorf("Expected color %q, got %q", color, tagType.Color)
		}

		str := tagType.String()
		if !contains(str, color) {
			t.Errorf("String() should contain color %q, got %q", color, str)
		}
	}
}

// TestTagAllSources tests tags on all supported source types
func TestTagAllSources(t *testing.T) {
	sources := []struct {
		sourceType string
		sourceID   int
	}{
		{types.TagSourceTask, 1},
		{types.TagSourceCallback, 2},
		{types.TagSourceFile, 3},
		{types.TagSourcePayload, 4},
		{types.TagSourceArtifact, 5},
		{types.TagSourceProcess, 6},
		{types.TagSourceKeylog, 7},
	}

	for _, src := range sources {
		tag := types.Tag{
			ID:         1,
			TagTypeID:  10,
			SourceType: src.sourceType,
			SourceID:   src.sourceID,
		}

		if tag.SourceType != src.sourceType {
			t.Errorf("Expected SourceType %q, got %q", src.sourceType, tag.SourceType)
		}
		if tag.SourceID != src.sourceID {
			t.Errorf("Expected SourceID %d, got %d", src.sourceID, tag.SourceID)
		}

		str := tag.String()
		if !contains(str, src.sourceType) {
			t.Errorf("String() should contain source type %q, got %q", src.sourceType, str)
		}
	}
}

// Note: Tag and TagType timestamp tests removed as these fields don't exist in Mythic's database schema
