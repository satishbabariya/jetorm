package core

import (
	"fmt"
	"reflect"
	"strings"
)

// Entity represents metadata about a database entity
type Entity struct {
	Type       reflect.Type
	TableName  string
	Fields     []Field
	PrimaryKey *Field
}

// Field represents metadata about an entity field
type Field struct {
	Name            string
	DBName          string
	Type            reflect.Type
	PrimaryKey      bool
	AutoIncrement   bool
	Unique          bool
	NotNull         bool
	Index           string
	UniqueIndex     string
	CompositeIndex  *CompositeIndex
	Size            int
	Default         string
	Check           string
	ForeignKey      string
	OnDelete        string // cascade, set_null, set_default, restrict, no_action
	OnUpdate        string // cascade, set_null, set_default, restrict, no_action
	ExplicitType    string // type:text, type:decimal(10,2), etc.
	AutoNowAdd      bool
	AutoNow         bool
	Ignored         bool // Field is ignored (db:"-")
}

// CompositeIndex represents a composite index definition
type CompositeIndex struct {
	Name  string
	Order int
}

// EntityMetadata extracts metadata from an entity type
func EntityMetadata(entity interface{}) (*Entity, error) {
	t := reflect.TypeOf(entity)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	if t.Kind() != reflect.Struct {
		return nil, ErrInvalidEntity
	}

	meta := &Entity{
		Type:      t,
		TableName: toSnakeCase(t.Name()),
		Fields:    make([]Field, 0),
	}

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		fieldMeta := parseFieldTags(field)
		meta.Fields = append(meta.Fields, fieldMeta)

		if fieldMeta.PrimaryKey {
			meta.PrimaryKey = &fieldMeta
		}
	}

	return meta, nil
}

// parseFieldTags parses struct tags for a field
func parseFieldTags(field reflect.StructField) Field {
	dbTag := field.Tag.Get("db")
	
	// Check if field is ignored
	if dbTag == "-" {
		return Field{
			Name:    field.Name,
			DBName:  "-",
			Type:    field.Type,
			Ignored: true,
		}
	}

	f := Field{
		Name:   field.Name,
		DBName: dbTag,
		Type:   field.Type,
	}

	// Default db name to snake_case of field name
	if f.DBName == "" {
		f.DBName = toSnakeCase(field.Name)
	}

	// Parse jet tags
	jetTag := field.Tag.Get("jet")
	if jetTag == "-" {
		f.Ignored = true
		return f
	}

	if jetTag != "" {
		parseTags := parseTag(jetTag)
		for _, tag := range parseTags {
			switch tag.Key {
			case "primary_key":
				f.PrimaryKey = true
			case "auto_increment":
				f.AutoIncrement = true
			case "unique":
				f.Unique = true
			case "not_null":
				f.NotNull = true
			case "index":
				f.Index = tag.Value
				if f.Index == "" {
					f.Index = "idx_" + f.DBName
				}
			case "unique_index":
				f.UniqueIndex = tag.Value
				if f.UniqueIndex == "" {
					f.UniqueIndex = "idx_unique_" + f.DBName
				}
			case "composite_index":
				// Format: composite_index:name:order
				// Example: composite_index:idx_sku_store:1
				// The value is already in format "name:order" from parseTag
				parts := strings.Split(tag.Value, ":")
				if len(parts) >= 2 {
					var order int
					fmt.Sscanf(parts[1], "%d", &order)
					f.CompositeIndex = &CompositeIndex{
						Name:  parts[0],
						Order: order,
					}
				}
			case "size":
				if tag.Value != "" {
					// Parse size value
					var size int
					_, _ = fmt.Sscanf(tag.Value, "%d", &size)
					f.Size = size
				}
			case "type":
				// Explicit type specification
				// Examples: type:text, type:decimal(10,2), type:jsonb
				f.ExplicitType = tag.Value
			case "default":
				f.Default = tag.Value
			case "check":
				f.Check = tag.Value
			case "foreign_key":
				// Format: foreign_key:table.column
				// Example: foreign_key:companies.id
				f.ForeignKey = tag.Value
			case "on_delete":
				// Cascade actions: cascade, set_null, set_default, restrict, no_action
				f.OnDelete = tag.Value
			case "on_update":
				// Cascade actions: cascade, set_null, set_default, restrict, no_action
				f.OnUpdate = tag.Value
			case "auto_now_add":
				f.AutoNowAdd = true
			case "auto_now":
				f.AutoNow = true
			}
		}
	}

	return f
}

type tagPair struct {
	Key   string
	Value string
}

// parseTag parses a comma-separated tag string
func parseTag(tag string) []tagPair {
	var pairs []tagPair
	parts := splitTag(tag)

	for _, part := range parts {
		if idx := strings.Index(part, ":"); idx > 0 {
			pairs = append(pairs, tagPair{
				Key:   part[:idx],
				Value: part[idx+1:],
			})
		} else {
			pairs = append(pairs, tagPair{
				Key: part,
			})
		}
	}

	return pairs
}

// splitTag splits a tag string by commas, respecting quotes and parentheses
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
