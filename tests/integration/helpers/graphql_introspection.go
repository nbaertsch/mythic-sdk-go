//go:build integration

package helpers

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// SchemaField represents a GraphQL field from introspection
type SchemaField struct {
	Name string
	Type SchemaType
}

// SchemaType represents a GraphQL type from introspection
type SchemaType struct {
	Kind   string      // SCALAR, OBJECT, LIST, NON_NULL, etc.
	Name   *string     // Type name (null for NON_NULL wrappers)
	OfType *SchemaType // Nested type for NON_NULL and LIST
}

// SchemaTypeInfo represents complete type information from introspection
type SchemaTypeInfo struct {
	Name   string
	Fields []SchemaField
}

// GraphQLIntrospectionClient provides methods for querying GraphQL schema
type GraphQLIntrospectionClient interface {
	ExecuteRawGraphQL(ctx context.Context, query string, variables map[string]interface{}) (map[string]interface{}, error)
}

// QuerySchemaType queries the GraphQL schema for a specific type and returns its fields
func QuerySchemaType(t *testing.T, client GraphQLIntrospectionClient, typeName string) SchemaTypeInfo {
	t.Helper()

	query := `query IntrospectType($typeName: String!) {
		__type(name: $typeName) {
			name
			fields {
				name
				type {
					kind
					name
					ofType {
						kind
						name
						ofType {
							kind
							name
						}
					}
				}
			}
		}
	}`

	variables := map[string]interface{}{
		"typeName": typeName,
	}

	ctx := context.Background()
	result, err := client.ExecuteRawGraphQL(ctx, query, variables)
	require.NoError(t, err, "Failed to execute introspection query for type %s", typeName)

	// Parse response
	data, ok := result["data"].(map[string]interface{})
	require.True(t, ok, "Response missing 'data' field")

	typeInfo, ok := data["__type"].(map[string]interface{})
	require.True(t, ok, "Response missing '__type' field")

	if typeInfo == nil {
		t.Fatalf("Type '%s' not found in schema", typeName)
	}

	schemaType := SchemaTypeInfo{
		Name:   typeName,
		Fields: []SchemaField{},
	}

	// Parse name
	if name, ok := typeInfo["name"].(string); ok {
		schemaType.Name = name
	}

	// Parse fields
	if fields, ok := typeInfo["fields"].([]interface{}); ok {
		for _, fieldInterface := range fields {
			field, ok := fieldInterface.(map[string]interface{})
			if !ok {
				continue
			}

			fieldName, _ := field["name"].(string)
			fieldTypeMap, _ := field["type"].(map[string]interface{})

			schemaField := SchemaField{
				Name: fieldName,
				Type: parseSchemaType(fieldTypeMap),
			}

			schemaType.Fields = append(schemaType.Fields, schemaField)
		}
	}

	return schemaType
}

// parseSchemaType recursively parses GraphQL type information
func parseSchemaType(typeMap map[string]interface{}) SchemaType {
	if typeMap == nil {
		return SchemaType{}
	}

	kind, _ := typeMap["kind"].(string)

	var name *string
	if n, ok := typeMap["name"].(string); ok {
		name = &n
	}

	var ofType *SchemaType
	if ofTypeMap, ok := typeMap["ofType"].(map[string]interface{}); ok && ofTypeMap != nil {
		parsed := parseSchemaType(ofTypeMap)
		ofType = &parsed
	}

	return SchemaType{
		Kind:   kind,
		Name:   name,
		OfType: ofType,
	}
}

// GetBaseTypeName returns the base type name, unwrapping NON_NULL and LIST wrappers
func (st SchemaType) GetBaseTypeName() string {
	if st.Name != nil {
		return *st.Name
	}
	if st.OfType != nil {
		return st.OfType.GetBaseTypeName()
	}
	return ""
}

// GetBaseTypeKind returns the base type kind, unwrapping NON_NULL and LIST wrappers
func (st SchemaType) GetBaseTypeKind() string {
	// If this is a scalar or object, return its kind
	if st.Kind == "SCALAR" || st.Kind == "OBJECT" || st.Kind == "ENUM" {
		return st.Kind
	}

	// If this is a wrapper (NON_NULL, LIST), recurse
	if st.OfType != nil {
		return st.OfType.GetBaseTypeKind()
	}

	return st.Kind
}

// IsArray returns true if this type is a LIST (array)
func (st SchemaType) IsArray() bool {
	if st.Kind == "LIST" {
		return true
	}
	// Check wrapped type for NON_NULL lists
	if st.Kind == "NON_NULL" && st.OfType != nil {
		return st.OfType.IsArray()
	}
	return false
}

// IsScalar returns true if the base type is a SCALAR
func (st SchemaType) IsScalar() bool {
	return st.GetBaseTypeKind() == "SCALAR"
}

// IsObject returns true if the base type is an OBJECT
func (st SchemaType) IsObject() bool {
	return st.GetBaseTypeKind() == "OBJECT"
}

// AssertFieldExists asserts that a field exists in the schema type with the expected type
func AssertFieldExists(t *testing.T, schema SchemaTypeInfo, fieldName string, expectedTypeName string) {
	t.Helper()

	field, found := findField(schema, fieldName)
	if !found {
		t.Errorf("Schema type '%s' missing expected field: %s", schema.Name, fieldName)
		t.Logf("Available fields: %v", getFieldNames(schema))
		return
	}

	actualTypeName := field.Type.GetBaseTypeName()
	assert.Equal(t, expectedTypeName, actualTypeName,
		"Field '%s' has wrong type. Expected: %s, Got: %s",
		fieldName, expectedTypeName, actualTypeName)
}

// AssertFieldNotExists asserts that a field does NOT exist in the schema type
func AssertFieldNotExists(t *testing.T, schema SchemaTypeInfo, fieldName string) {
	t.Helper()

	_, found := findField(schema, fieldName)
	if found {
		t.Errorf("Schema type '%s' contains unexpected field: %s (this field should not exist)",
			schema.Name, fieldName)
	}
}

// AssertFieldType asserts that a field has a specific type characteristic (scalar, object, array, etc.)
func AssertFieldType(t *testing.T, schema SchemaTypeInfo, fieldName string, expectedTypeKind string) {
	t.Helper()

	field, found := findField(schema, fieldName)
	if !found {
		t.Errorf("Schema type '%s' missing field: %s", schema.Name, fieldName)
		return
	}

	switch expectedTypeKind {
	case "array", "list", "LIST":
		assert.True(t, field.Type.IsArray(),
			"Field '%s' should be an array but is type: %s",
			fieldName, field.Type.GetBaseTypeName())

	case "scalar", "SCALAR":
		assert.True(t, field.Type.IsScalar(),
			"Field '%s' should be a scalar but is type: %s (kind: %s)",
			fieldName, field.Type.GetBaseTypeName(), field.Type.GetBaseTypeKind())

	case "object", "OBJECT":
		assert.True(t, field.Type.IsObject(),
			"Field '%s' should be an object but is type: %s (kind: %s)",
			fieldName, field.Type.GetBaseTypeName(), field.Type.GetBaseTypeKind())

	case "jsonb", "JSONB":
		// JSONB is typically a scalar in GraphQL
		baseName := field.Type.GetBaseTypeName()
		assert.True(t, baseName == "jsonb" || baseName == "JSONB",
			"Field '%s' should be jsonb but is type: %s",
			fieldName, baseName)

	default:
		// Try exact match on base type name
		actualTypeName := field.Type.GetBaseTypeName()
		assert.Equal(t, expectedTypeKind, actualTypeName,
			"Field '%s' has wrong type. Expected: %s, Got: %s",
			fieldName, expectedTypeKind, actualTypeName)
	}
}

// AssertFieldTypeExact asserts exact type name match (case-sensitive)
func AssertFieldTypeExact(t *testing.T, schema SchemaTypeInfo, fieldName string, expectedTypeName string) {
	t.Helper()

	field, found := findField(schema, fieldName)
	if !found {
		t.Errorf("Schema type '%s' missing field: %s", schema.Name, fieldName)
		return
	}

	actualTypeName := field.Type.GetBaseTypeName()
	require.Equal(t, expectedTypeName, actualTypeName,
		"Field '%s' has wrong type. Expected: %s, Got: %s",
		fieldName, expectedTypeName, actualTypeName)
}

// findField searches for a field by name in the schema
func findField(schema SchemaTypeInfo, fieldName string) (SchemaField, bool) {
	for _, field := range schema.Fields {
		if field.Name == fieldName {
			return field, true
		}
	}
	return SchemaField{}, false
}

// getFieldNames returns a list of all field names in the schema (for debugging)
func getFieldNames(schema SchemaTypeInfo) []string {
	names := make([]string, len(schema.Fields))
	for i, field := range schema.Fields {
		names[i] = field.Name
	}
	return names
}

// PrintSchemaFields prints all fields in a schema type (for debugging)
func PrintSchemaFields(t *testing.T, schema SchemaTypeInfo) {
	t.Helper()

	t.Logf("Schema Type: %s", schema.Name)
	t.Logf("Total Fields: %d", len(schema.Fields))
	for _, field := range schema.Fields {
		t.Logf("  - %s: %s (kind: %s, isArray: %v)",
			field.Name,
			field.Type.GetBaseTypeName(),
			field.Type.GetBaseTypeKind(),
			field.Type.IsArray())
	}
}

// AssertSchemaMatchesSDK validates that SDK struct tags match the GraphQL schema
func AssertSchemaMatchesSDK(t *testing.T, schema SchemaTypeInfo, sdkFieldMappings map[string]string) {
	t.Helper()

	for sdkFieldName, graphqlFieldName := range sdkFieldMappings {
		_, found := findField(schema, graphqlFieldName)
		assert.True(t, found,
			"SDK field '%s' maps to GraphQL field '%s' which does not exist in schema",
			sdkFieldName, graphqlFieldName)
	}
}

// ValidateNoDeprecatedFields checks that certain fields are NOT present (removed in newer schema versions)
func ValidateNoDeprecatedFields(t *testing.T, schema SchemaTypeInfo, deprecatedFields []string) {
	t.Helper()

	for _, fieldName := range deprecatedFields {
		AssertFieldNotExists(t, schema, fieldName)
	}
}

// GetFieldTypeDescription returns a human-readable description of a field's type
func GetFieldTypeDescription(field SchemaField) string {
	if field.Type.IsArray() {
		return fmt.Sprintf("array of %s", field.Type.GetBaseTypeName())
	}
	return field.Type.GetBaseTypeName()
}
