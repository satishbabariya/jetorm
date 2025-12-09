package migration

import (
	"fmt"
	"reflect"
	"strings"
)

// SchemaGenerator generates SQL schema from Go struct definitions
type SchemaGenerator struct{}

// NewSchemaGenerator creates a new schema generator
func NewSchemaGenerator() *SchemaGenerator {
	return &SchemaGenerator{}
}

// GenerateCreateTable generates a CREATE TABLE statement from a struct type
func (sg *SchemaGenerator) GenerateCreateTable(entityType reflect.Type, tableName string) (string, error) {
	if entityType.Kind() == reflect.Ptr {
		entityType = entityType.Elem()
	}
	if entityType.Kind() != reflect.Struct {
		return "", fmt.Errorf("entity type must be a struct")
	}
	
	var columns []string
	var primaryKeys []string
	
	for i := 0; i < entityType.NumField(); i++ {
		field := entityType.Field(i)
		
		// Skip unexported fields
		if !field.IsExported() {
			continue
		}
		
		dbTag := field.Tag.Get("db")
		if dbTag == "" || dbTag == "-" {
			continue
		}
		
		jetTag := field.Tag.Get("jet")
		columnDef := sg.generateColumnDefinition(field, dbTag, jetTag)
		columns = append(columns, columnDef)
		
		// Check for primary key
		if strings.Contains(jetTag, "primary_key") {
			primaryKeys = append(primaryKeys, dbTag)
		}
	}
	
	if len(columns) == 0 {
		return "", fmt.Errorf("no columns found for table %s", tableName)
	}
	
	query := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (\n", tableName)
	query += strings.Join(columns, ",\n")
	
	if len(primaryKeys) > 0 {
		query += fmt.Sprintf(",\nPRIMARY KEY (%s)", strings.Join(primaryKeys, ", "))
	}
	
	query += "\n);"
	
	return query, nil
}

// generateColumnDefinition generates a column definition from field metadata
func (sg *SchemaGenerator) generateColumnDefinition(field reflect.StructField, dbName, jetTag string) string {
	var parts []string
	
	// Column name
	parts = append(parts, dbName)
	
	// Column type
	columnType := sg.getColumnType(field.Type, jetTag)
	parts = append(parts, columnType)
	
	// Constraints
	if strings.Contains(jetTag, "not_null") {
		parts = append(parts, "NOT NULL")
	}
	
	if strings.Contains(jetTag, "unique") {
		parts = append(parts, "UNIQUE")
	}
	
	// Default value
	if defaultVal := sg.extractTagValue(jetTag, "default"); defaultVal != "" {
		parts = append(parts, fmt.Sprintf("DEFAULT %s", defaultVal))
	}
	
	return strings.Join(parts, " ")
}

// getColumnType maps Go types to PostgreSQL column types
func (sg *SchemaGenerator) getColumnType(goType reflect.Type, jetTag string) string {
	// Check for explicit type in jet tag
	if explicitType := sg.extractTagValue(jetTag, "type"); explicitType != "" {
		return explicitType
	}
	
	// Map Go types to PostgreSQL types
	switch goType.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return "BIGINT"
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return "BIGINT"
	case reflect.Float32:
		return "REAL"
	case reflect.Float64:
		return "DOUBLE PRECISION"
	case reflect.Bool:
		return "BOOLEAN"
	case reflect.String:
		if size := sg.extractTagValue(jetTag, "size"); size != "" {
			return fmt.Sprintf("VARCHAR(%s)", size)
		}
		return "TEXT"
	case reflect.Slice, reflect.Array:
		if goType.Elem().Kind() == reflect.Uint8 {
			return "BYTEA"
		}
		return "TEXT" // JSON array
	case reflect.Struct:
		if goType.String() == "time.Time" {
			return "TIMESTAMP"
		}
		return "TEXT" // JSON object
	default:
		return "TEXT"
	}
}

// extractTagValue extracts a value from a tag string
func (sg *SchemaGenerator) extractTagValue(tag, key string) string {
	parts := strings.Split(tag, ",")
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if strings.HasPrefix(part, key+":") {
			return strings.TrimPrefix(part, key+":")
		}
	}
	return ""
}

