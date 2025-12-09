package core

import (
	"context"
	"fmt"
	"reflect"
	"strings"
)

// RelationshipType represents the type of relationship
type RelationshipType int

const (
	OneToOne RelationshipType = iota
	OneToMany
	ManyToOne
	ManyToMany
)

// Relationship represents a relationship between entities
type Relationship struct {
	Type         RelationshipType
	Field        string
	TargetEntity string
	ForeignKey   string
	JoinTable    string // For many-to-many
	JoinColumn   string // For many-to-many
	InverseJoinColumn string // For many-to-many
	OnDelete     string
	OnUpdate     string
	Lazy         bool
}

// RelationshipManager manages entity relationships
type RelationshipManager struct {
	relationships map[string][]Relationship
}

// NewRelationshipManager creates a new relationship manager
func NewRelationshipManager() *RelationshipManager {
	return &RelationshipManager{
		relationships: make(map[string][]Relationship),
	}
}

// RegisterRelationship registers a relationship for an entity
func (rm *RelationshipManager) RegisterRelationship(entityType string, rel Relationship) {
	if rm.relationships[entityType] == nil {
		rm.relationships[entityType] = make([]Relationship, 0)
	}
	rm.relationships[entityType] = append(rm.relationships[entityType], rel)
}

// GetRelationships returns all relationships for an entity type
func (rm *RelationshipManager) GetRelationships(entityType string) []Relationship {
	return rm.relationships[entityType]
}

// LoadRelationships loads relationships from entity tags
func LoadRelationships(entityType reflect.Type) []Relationship {
	var relationships []Relationship

	for i := 0; i < entityType.NumField(); i++ {
		field := entityType.Field(i)
		if !field.IsExported() {
			continue
		}

		jetTag := field.Tag.Get("jet")
		if jetTag == "" {
			continue
		}

		// Parse relationship tags
		if strings.Contains(jetTag, "one_to_one") {
			rel := parseOneToOne(field, jetTag)
			if rel != nil {
				relationships = append(relationships, *rel)
			}
		} else if strings.Contains(jetTag, "one_to_many") {
			rel := parseOneToMany(field, jetTag)
			if rel != nil {
				relationships = append(relationships, *rel)
			}
		} else if strings.Contains(jetTag, "many_to_one") {
			rel := parseManyToOne(field, jetTag)
			if rel != nil {
				relationships = append(relationships, *rel)
			}
		} else if strings.Contains(jetTag, "many_to_many") {
			rel := parseManyToMany(field, jetTag)
			if rel != nil {
				relationships = append(relationships, *rel)
			}
		}
	}

	return relationships
}

// parseOneToOne parses a one-to-one relationship tag
func parseOneToOne(field reflect.StructField, jetTag string) *Relationship {
	rel := &Relationship{
		Type:  OneToOne,
		Field: field.Name,
	}

	// Extract target entity
	if target := extractTagValue(jetTag, "one_to_one"); target != "" {
		rel.TargetEntity = target
	}

	// Extract foreign key
	if fk := extractTagValue(jetTag, "foreign_key"); fk != "" {
		rel.ForeignKey = fk
	}

	// Extract cascade actions
	if onDelete := extractTagValue(jetTag, "on_delete"); onDelete != "" {
		rel.OnDelete = onDelete
	}
	if onUpdate := extractTagValue(jetTag, "on_update"); onUpdate != "" {
		rel.OnUpdate = onUpdate
	}

	return rel
}

// parseOneToMany parses a one-to-many relationship tag
func parseOneToMany(field reflect.StructField, jetTag string) *Relationship {
	rel := &Relationship{
		Type:  OneToMany,
		Field: field.Name,
	}

	if target := extractTagValue(jetTag, "one_to_many"); target != "" {
		rel.TargetEntity = target
	}

	if mappedBy := extractTagValue(jetTag, "mapped_by"); mappedBy != "" {
		rel.ForeignKey = mappedBy
	}

	return rel
}

// parseManyToOne parses a many-to-one relationship tag
func parseManyToOne(field reflect.StructField, jetTag string) *Relationship {
	rel := &Relationship{
		Type:  ManyToOne,
		Field: field.Name,
	}

	if target := extractTagValue(jetTag, "many_to_one"); target != "" {
		rel.TargetEntity = target
	}

	if fk := extractTagValue(jetTag, "foreign_key"); fk != "" {
		rel.ForeignKey = fk
	}

	if onDelete := extractTagValue(jetTag, "on_delete"); onDelete != "" {
		rel.OnDelete = onDelete
	}

	return rel
}

// parseManyToMany parses a many-to-many relationship tag
func parseManyToMany(field reflect.StructField, jetTag string) *Relationship {
	rel := &Relationship{
		Type:  ManyToMany,
		Field: field.Name,
	}

	if target := extractTagValue(jetTag, "many_to_many"); target != "" {
		rel.TargetEntity = target
	}

	if joinTable := extractTagValue(jetTag, "join_table"); joinTable != "" {
		rel.JoinTable = joinTable
	}

	if joinColumn := extractTagValue(jetTag, "join_column"); joinColumn != "" {
		rel.JoinColumn = joinColumn
	}

	if inverseJoinColumn := extractTagValue(jetTag, "inverse_join_column"); inverseJoinColumn != "" {
		rel.InverseJoinColumn = inverseJoinColumn
	}

	return rel
}

// extractTagValue extracts a value from a tag string
func extractTagValue(tag, key string) string {
	parts := strings.Split(tag, ",")
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if strings.HasPrefix(part, key+":") {
			return strings.TrimPrefix(part, key+":")
		}
	}
	return ""
}

// RelationshipRepository provides methods for loading relationships
type RelationshipRepository[T any, ID comparable] interface {
	// LoadOne loads a one-to-one or many-to-one relationship
	LoadOne(ctx context.Context, entity *T, relationship string) error
	
	// LoadMany loads a one-to-many or many-to-many relationship
	LoadMany(ctx context.Context, entity *T, relationship string) error
	
	// LoadAll loads all relationships for an entity
	LoadAll(ctx context.Context, entity *T) error
}

// EagerLoad loads relationships eagerly
func EagerLoad[T any, ID comparable](repo Repository[T, ID], entities []*T, relationships ...string) error {
	// This is a placeholder - full implementation would load relationships
	// based on the relationship configuration
	return nil
}

// LazyLoad loads relationships lazily
func LazyLoad[T any, ID comparable](repo Repository[T, ID], entity *T, relationship string) error {
	// This is a placeholder - full implementation would load relationships
	// on demand
	return nil
}

// JoinQuery builds a query with relationship joins
func JoinQuery[T any](tableName string, relationships []Relationship) string {
	var joins []string
	
	for _, rel := range relationships {
		switch rel.Type {
		case ManyToOne, OneToOne:
			// Add JOIN for foreign key relationship
			join := fmt.Sprintf("LEFT JOIN %s ON %s.%s = %s.id",
				rel.TargetEntity, tableName, rel.ForeignKey, rel.TargetEntity)
			joins = append(joins, join)
		case OneToMany:
			// For one-to-many, we'd typically use a subquery or separate query
			// This is a simplified version
		case ManyToMany:
			// Add JOIN for many-to-many through join table
			if rel.JoinTable != "" {
				join := fmt.Sprintf("LEFT JOIN %s ON %s.%s = %s.id",
					rel.JoinTable, rel.JoinTable, rel.JoinColumn, tableName)
				joins = append(joins, join)
				join2 := fmt.Sprintf("LEFT JOIN %s ON %s.%s = %s.id",
					rel.TargetEntity, rel.JoinTable, rel.InverseJoinColumn, rel.TargetEntity)
				joins = append(joins, join2)
			}
		}
	}
	
	return strings.Join(joins, " ")
}

