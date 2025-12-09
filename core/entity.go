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
	Name          string
	DBName        string
	Type          reflect.Type
	PrimaryKey    bool
	AutoIncrement bool
	Unique        bool
	NotNull       bool
	Index         string
	Size          int
	Default       string
	Check         string
	ForeignKey    string
	AutoNowAdd    bool
	AutoNow       bool
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
	f := Field{
		Name:   field.Name,
		DBName: field.Tag.Get("db"),
		Type:   field.Type,
	}

	// Default db name to snake_case of field name
	if f.DBName == "" {
		f.DBName = toSnakeCase(field.Name)
	}

	// Parse jet tags
	jetTag := field.Tag.Get("jet")
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
			case "size":
				if tag.Value != "" {
					// Parse size value
					var size int
					_, _ = fmt.Sscanf(tag.Value, "%d", &size)
					f.Size = size
				}
			case "default":
				f.Default = tag.Value
			case "check":
				f.Check = tag.Value
			case "foreign_key":
				f.ForeignKey = tag.Value
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

// splitTag splits a tag string by commas, respecting quotes
func splitTag(tag string) []string {
	var parts []string
	var current strings.Builder
	inQuote := false

	for _, r := range tag {
		switch r {
		case '\'':
			inQuote = !inQuote
		case ',':
			if !inQuote {
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
