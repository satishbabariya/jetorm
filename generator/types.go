package generator

import (
	"fmt"
	"go/types"
	"reflect"
	"strings"
)

// TypeLoader loads type information using go/types
type TypeLoader struct {
	pkg *types.Package
}

// NewTypeLoader creates a new type loader
func NewTypeLoader(pkgPath string) (*TypeLoader, error) {
	// This is a placeholder - full implementation would use go/types
	// For now, we'll use reflect-based approach
	return &TypeLoader{}, nil
}

// LoadEntityType loads entity type information
func (tl *TypeLoader) LoadEntityType(typeName string) (*EntityTypeInfo, error) {
	// Placeholder - would use go/types in full implementation
	return nil, fmt.Errorf("full type loading not yet implemented")
}

// EntityTypeInfo contains information about an entity type
type EntityTypeInfo struct {
	Name       string
	Package    string
	Fields     []FieldInfo
	PrimaryKey *FieldInfo
	TableName  string
}

// FieldInfo contains information about a struct field
type FieldInfo struct {
	Name         string
	DBName       string
	Type         types.Type
	IsPrimaryKey bool
	IsAutoInc    bool
	Tags         map[string]string
}

// GetIDType returns the ID type for an entity
func (eti *EntityTypeInfo) GetIDType() string {
	if eti.PrimaryKey != nil {
		return eti.PrimaryKey.Type.String()
	}
	return "int64" // Default
}

// ReflectTypeLoader uses reflection to load type information
type ReflectTypeLoader struct{}

// NewReflectTypeLoader creates a new reflection-based type loader
func NewReflectTypeLoader() *ReflectTypeLoader {
	return &ReflectTypeLoader{}
}

// LoadEntityTypeFromReflect loads entity type from reflect.Type
func (rtl *ReflectTypeLoader) LoadEntityTypeFromReflect(entityType reflect.Type) (*EntityTypeInfo, error) {
	if entityType.Kind() == reflect.Ptr {
		entityType = entityType.Elem()
	}
	if entityType.Kind() != reflect.Struct {
		return nil, fmt.Errorf("entity type must be a struct")
	}

	info := &EntityTypeInfo{
		Name:      entityType.Name(),
		Package:   entityType.PkgPath(),
		Fields:    make([]FieldInfo, 0),
		TableName: toSnakeCase(entityType.Name()),
	}

	for i := 0; i < entityType.NumField(); i++ {
		field := entityType.Field(i)
		if !field.IsExported() {
			continue
		}

		dbTag := field.Tag.Get("db")
		if dbTag == "" || dbTag == "-" {
			continue
		}

		jetTag := field.Tag.Get("jet")
		fieldInfo := FieldInfo{
			Name:  field.Name,
			DBName: dbTag,
			Type:  nil, // Would be set with go/types
			Tags:  parseTags(jetTag),
		}

		// Check for primary key
		if _, ok := fieldInfo.Tags["primary_key"]; ok {
			fieldInfo.IsPrimaryKey = true
			info.PrimaryKey = &fieldInfo
		}

		// Check for auto increment
		if _, ok := fieldInfo.Tags["auto_increment"]; ok {
			fieldInfo.IsAutoInc = true
		}

		info.Fields = append(info.Fields, fieldInfo)
	}

	return info, nil
}

// parseTags parses a tag string into a map
func parseTags(tag string) map[string]string {
	result := make(map[string]string)
	if tag == "" {
		return result
	}

	parts := splitTag(tag)
	for _, part := range parts {
		if idx := strings.Index(part, ":"); idx > 0 {
			key := part[:idx]
			value := part[idx+1:]
			result[key] = value
		} else {
			result[part] = ""
		}
	}

	return result
}

// splitTag splits a tag string (would be imported from entity.go)
func splitTag(tag string) []string {
	var parts []string
	var current strings.Builder
	inQuote := false
	parenDepth := 0

	for _, r := range tag {
		switch r {
		case '\'':
			inQuote = !inQuote
		case '(':
			if !inQuote {
				parenDepth++
			}
		case ')':
			if !inQuote {
				parenDepth--
			}
		case ',':
			if !inQuote && parenDepth == 0 {
				if current.Len() > 0 {
					parts = append(parts, strings.TrimSpace(current.String()))
					current.Reset()
				}
				continue
			}
		}
		current.WriteRune(r)
	}

	if current.Len() > 0 {
		parts = append(parts, strings.TrimSpace(current.String()))
	}

	return parts
}

// toSnakeCaseHelper converts a string to snake_case
func toSnakeCaseHelper(s string) string {
	var result strings.Builder
	for i, r := range s {
		if i > 0 && r >= 'A' && r <= 'Z' {
			result.WriteRune('_')
		}
		if r >= 'A' && r <= 'Z' {
			result.WriteRune(r + 32)
		} else {
			result.WriteRune(r)
		}
	}
	return result.String()
}

