package core

import (
	"context"
	"fmt"
	"reflect"
	"strings"
	"time"
)

// RepositoryHelpers provides helper functions for repositories

// FindOrCreate finds an entity or creates it if not found
func FindOrCreate[T any, ID comparable](
	ctx context.Context,
	repo Repository[T, ID],
	finder func(context.Context) (*T, error),
	creator func(context.Context) (*T, error),
) (*T, error) {
	entity, err := finder(ctx)
	if err == nil && entity != nil {
		return entity, nil
	}
	if err != nil && !IsNotFound(err) {
		return nil, err
	}

	// Not found, create
	return creator(ctx)
}

// SaveIfNotExists saves an entity only if it doesn't exist
func SaveIfNotExists[T any, ID comparable](
	ctx context.Context,
	repo Repository[T, ID],
	entity *T,
	id ID,
) (*T, error) {
	existing, err := repo.FindByID(ctx, id)
	if err == nil && existing != nil {
		return existing, nil
	}
	if err != nil && !IsNotFound(err) {
		return nil, err
	}

	return repo.Save(ctx, entity)
}

// UpdateIfExists updates an entity only if it exists
func UpdateIfExists[T any, ID comparable](
	ctx context.Context,
	repo Repository[T, ID],
	entity *T,
	id ID,
) (*T, error) {
	existing, err := repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if existing == nil {
		return nil, ErrNotFound
	}

	return repo.Update(ctx, entity)
}

// DeleteIfExists deletes an entity only if it exists
func DeleteIfExists[T any, ID comparable](
	ctx context.Context,
	repo Repository[T, ID],
	id ID,
) error {
	existing, err := repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if existing == nil {
		return ErrNotFound
	}

	return repo.DeleteByID(ctx, id)
}

// BatchFind finds multiple entities by IDs in batches
func BatchFind[T any, ID comparable](
	ctx context.Context,
	repo Repository[T, ID],
	ids []ID,
	batchSize int,
) ([]*T, error) {
	if batchSize <= 0 {
		batchSize = 100
	}

	var results []*T
	for i := 0; i < len(ids); i += batchSize {
		end := i + batchSize
		if end > len(ids) {
			end = len(ids)
		}

		batch := ids[i:end]
		batchResults, err := repo.FindAllByIDs(ctx, batch)
		if err != nil {
			return nil, fmt.Errorf("batch find failed at offset %d: %w", i, err)
		}

		results = append(results, batchResults...)
	}

	return results, nil
}

// Exists checks if an entity exists
func Exists[T any, ID comparable](
	ctx context.Context,
	repo Repository[T, ID],
	id ID,
) (bool, error) {
	_, err := repo.FindByID(ctx, id)
	if err == nil {
		return true, nil
	}
	if IsNotFound(err) {
		return false, nil
	}
	return false, err
}

// CountByCondition counts entities matching a condition
func CountByCondition[T any, ID comparable](
	ctx context.Context,
	repo Repository[T, ID],
	spec Specification[T],
) (int64, error) {
	return repo.CountWithSpec(ctx, spec)
}

// FindFirst finds the first entity matching a condition
func FindFirst[T any, ID comparable](
	ctx context.Context,
	repo Repository[T, ID],
	spec Specification[T],
) (*T, error) {
	return repo.FindOne(ctx, spec)
}

// FindAllMatching finds all entities matching a condition
func FindAllMatching[T any, ID comparable](
	ctx context.Context,
	repo Repository[T, ID],
	spec Specification[T],
) ([]*T, error) {
	return repo.FindAllWithSpec(ctx, spec)
}

// TimestampHelper provides timestamp utilities
type TimestampHelper struct{}

// NewTimestampHelper creates a new timestamp helper
func NewTimestampHelper() *TimestampHelper {
	return &TimestampHelper{}
}

// SetCreatedAt sets created_at timestamp on entity
func (th *TimestampHelper) SetCreatedAt(entity interface{}) error {
	return th.setTimestamp(entity, "created_at", time.Now())
}

// SetUpdatedAt sets updated_at timestamp on entity
func (th *TimestampHelper) SetUpdatedAt(entity interface{}) error {
	return th.setTimestamp(entity, "updated_at", time.Now())
}

// setTimestamp sets a timestamp field on entity
func (th *TimestampHelper) setTimestamp(entity interface{}, fieldName string, timestamp time.Time) error {
	entityType := reflect.TypeOf(entity)
	if entityType.Kind() == reflect.Ptr {
		entityType = entityType.Elem()
	}

	entityValue := reflect.ValueOf(entity)
	if entityValue.Kind() == reflect.Ptr {
		entityValue = entityValue.Elem()
	}

	// Find field by db tag or name
	for i := 0; i < entityType.NumField(); i++ {
		field := entityType.Field(i)
		dbTag := field.Tag.Get("db")
		
		if dbTag == fieldName || field.Name == toPascalCaseHelper(fieldName) {
			fieldValue := entityValue.Field(i)
			if fieldValue.CanSet() {
				if fieldValue.Type() == reflect.TypeOf(time.Time{}) {
					fieldValue.Set(reflect.ValueOf(timestamp))
					return nil
				}
			}
		}
	}

	return fmt.Errorf("could not set timestamp field %s", fieldName)
}

// toPascalCaseHelper converts snake_case to PascalCase (helper to avoid conflict)
func toPascalCaseHelper(s string) string {
	parts := strings.Split(s, "_")
	var result strings.Builder
	for _, part := range parts {
		if len(part) > 0 {
			result.WriteString(strings.ToUpper(part[:1]) + strings.ToLower(part[1:]))
		}
	}
	return result.String()
}

// EntityHelper provides entity manipulation utilities
type EntityHelper struct{}

// NewEntityHelper creates a new entity helper
func NewEntityHelper() *EntityHelper {
	return &EntityHelper{}
}

// Copy copies fields from source to destination
func (eh *EntityHelper) Copy(dest, src interface{}) error {
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

// Merge merges fields from source to destination (only non-zero values)
func (eh *EntityHelper) Merge(dest, src interface{}) error {
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
			if !IsZero(srcField.Interface()) {
				destField.Set(srcField)
			}
		}
	}

	return nil
}

// Clone creates a deep copy of an entity
func (eh *EntityHelper) Clone(entity interface{}) interface{} {
	entityType := reflect.TypeOf(entity)
	if entityType.Kind() == reflect.Ptr {
		entityType = entityType.Elem()
	}

	entityValue := reflect.ValueOf(entity)
	if entityValue.Kind() == reflect.Ptr {
		entityValue = entityValue.Elem()
	}

	newValue := reflect.New(entityType)
	newElem := newValue.Elem()

	for i := 0; i < entityType.NumField(); i++ {
		field := entityType.Field(i)
		if !field.IsExported() {
			continue
		}

		srcField := entityValue.Field(i)
		destField := newElem.Field(i)

		if destField.CanSet() {
			destField.Set(srcField)
		}
	}

	return newValue.Interface()
}

