package core

import (
	"fmt"
	"reflect"
	"strings"
	"unicode"
)

// EntityUtils provides utilities for entity manipulation

// GetFieldValue gets a field value from an entity
func GetFieldValue(entity interface{}, fieldName string) (interface{}, error) {
	entityValue := reflect.ValueOf(entity)
	if entityValue.Kind() == reflect.Ptr {
		entityValue = entityValue.Elem()
	}

	field := entityValue.FieldByName(fieldName)
	if !field.IsValid() {
		return nil, fmt.Errorf("field %s not found", fieldName)
	}

	return field.Interface(), nil
}

// SetFieldValue sets a field value on an entity
func SetFieldValue(entity interface{}, fieldName string, value interface{}) error {
	entityValue := reflect.ValueOf(entity)
	if entityValue.Kind() != reflect.Ptr {
		return fmt.Errorf("entity must be a pointer")
	}

	entityValue = entityValue.Elem()
	field := entityValue.FieldByName(fieldName)
	if !field.IsValid() {
		return fmt.Errorf("field %s not found", fieldName)
	}

	if !field.CanSet() {
		return fmt.Errorf("field %s cannot be set", fieldName)
	}

	valueValue := reflect.ValueOf(value)
	if !valueValue.Type().AssignableTo(field.Type()) {
		return fmt.Errorf("value type %v is not assignable to field type %v", valueValue.Type(), field.Type())
	}

	field.Set(valueValue)
	return nil
}

// GetFieldTag gets a field's tag value
func GetFieldTag(entity interface{}, fieldName, tagName string) (string, bool) {
	entityType := reflect.TypeOf(entity)
	if entityType.Kind() == reflect.Ptr {
		entityType = entityType.Elem()
	}

	field, found := entityType.FieldByName(fieldName)
	if !found {
		return "", false
	}

	tag := field.Tag.Get(tagName)
	return tag, tag != ""
}

// GetDBFieldName gets the database field name for a struct field
func GetDBFieldName(entity interface{}, fieldName string) (string, error) {
	tag, found := GetFieldTag(entity, fieldName, "db")
	if found && tag != "" {
		// Extract first part before comma
		parts := strings.Split(tag, ",")
		return parts[0], nil
	}
	return fieldName, nil
}

// GetJetTag gets the jet tag value for a field
func GetJetTag(entity interface{}, fieldName string) (string, bool) {
	return GetFieldTag(entity, fieldName, "jet")
}

// HasTag checks if a field has a specific tag
func HasTag(entity interface{}, fieldName, tagName, tagValue string) bool {
	tag, found := GetFieldTag(entity, fieldName, tagName)
	if !found {
		return false
	}

	parts := strings.Split(tag, ",")
	for _, part := range parts {
		if strings.TrimSpace(part) == tagValue {
			return true
		}
	}
	return false
}

// IsPrimaryKey checks if a field is a primary key
func IsPrimaryKey(entity interface{}, fieldName string) bool {
	return HasTag(entity, fieldName, "jet", "primary_key")
}

// IsRequired checks if a field is required
func IsRequired(entity interface{}, fieldName string) bool {
	return HasTag(entity, fieldName, "jet", "not_null") || HasTag(entity, fieldName, "validate", "required")
}

// GetPrimaryKeyField gets the primary key field name
func GetPrimaryKeyField(entity interface{}) (string, error) {
	entityType := reflect.TypeOf(entity)
	if entityType.Kind() == reflect.Ptr {
		entityType = entityType.Elem()
	}

	for i := 0; i < entityType.NumField(); i++ {
		field := entityType.Field(i)
		if IsPrimaryKey(entity, field.Name) {
			return field.Name, nil
		}
	}

	return "", fmt.Errorf("no primary key field found")
}

// GetTableName gets the table name for an entity
func GetTableName(entity interface{}) string {
	entityType := reflect.TypeOf(entity)
	if entityType.Kind() == reflect.Ptr {
		entityType = entityType.Elem()
	}

	// Check for TableName method
	tableNameMethod := reflect.ValueOf(entity).MethodByName("TableName")
	if tableNameMethod.IsValid() {
		results := tableNameMethod.Call(nil)
		if len(results) > 0 && results[0].Kind() == reflect.String {
			return results[0].String()
		}
	}

	// Default to snake_case of type name
	typeName := entityType.Name()
	return toSnakeCaseHelper(typeName)
}

// GetColumnNames gets all column names for an entity
func GetColumnNames(entity interface{}) []string {
	entityType := reflect.TypeOf(entity)
	if entityType.Kind() == reflect.Ptr {
		entityType = entityType.Elem()
	}

	var columns []string
	for i := 0; i < entityType.NumField(); i++ {
		field := entityType.Field(i)
		if !field.IsExported() {
			continue
		}

		dbName, err := GetDBFieldName(entity, field.Name)
		if err == nil && dbName != "-" {
			columns = append(columns, dbName)
		}
	}

	return columns
}

// GetFieldNames gets all field names for an entity
func GetFieldNames(entity interface{}) []string {
	entityType := reflect.TypeOf(entity)
	if entityType.Kind() == reflect.Ptr {
		entityType = entityType.Elem()
	}

	var fields []string
	for i := 0; i < entityType.NumField(); i++ {
		field := entityType.Field(i)
		if field.IsExported() {
			fields = append(fields, field.Name)
		}
	}

	return fields
}

// CopyFields copies fields from source to destination
func CopyFields(dest, src interface{}) error {
	destValue := reflect.ValueOf(dest)
	srcValue := reflect.ValueOf(src)

	if destValue.Kind() != reflect.Ptr {
		return fmt.Errorf("destination must be a pointer")
	}
	if srcValue.Kind() == reflect.Ptr {
		srcValue = srcValue.Elem()
	}

	destValue = destValue.Elem()
	destType := destValue.Type()

	for i := 0; i < destType.NumField(); i++ {
		destField := destValue.Field(i)
		if !destField.CanSet() {
			continue
		}

		fieldName := destType.Field(i).Name
		srcField := srcValue.FieldByName(fieldName)
		if srcField.IsValid() && srcField.Type() == destField.Type() {
			destField.Set(srcField)
		}
	}

	return nil
}

// CompareEntities compares two entities field by field
func CompareEntities(entity1, entity2 interface{}) (bool, []string) {
	var differences []string

	entity1Type := reflect.TypeOf(entity1)
	entity2Type := reflect.TypeOf(entity2)

	if entity1Type != entity2Type {
		return false, []string{"entity types differ"}
	}

	if entity1Type.Kind() == reflect.Ptr {
		entity1Type = entity1Type.Elem()
		entity2Type = entity2Type.Elem()
	}

	entity1Value := reflect.ValueOf(entity1)
	entity2Value := reflect.ValueOf(entity2)

	if entity1Value.Kind() == reflect.Ptr {
		entity1Value = entity1Value.Elem()
	}
	if entity2Value.Kind() == reflect.Ptr {
		entity2Value = entity2Value.Elem()
	}

	for i := 0; i < entity1Type.NumField(); i++ {
		field := entity1Type.Field(i)
		if !field.IsExported() {
			continue
		}

		field1Value := entity1Value.Field(i)
		field2Value := entity2Value.Field(i)

		if !reflect.DeepEqual(field1Value.Interface(), field2Value.Interface()) {
			differences = append(differences, fmt.Sprintf("field %s differs: %v != %v", field.Name, field1Value.Interface(), field2Value.Interface()))
		}
	}

	return len(differences) == 0, differences
}

// toSnakeCaseHelper converts PascalCase to snake_case (helper to avoid conflict)
func toSnakeCaseHelper(s string) string {
	var result strings.Builder
	for i, r := range s {
		if i > 0 && unicode.IsUpper(r) {
			result.WriteRune('_')
		}
		result.WriteRune(unicode.ToLower(r))
	}
	return result.String()
}

