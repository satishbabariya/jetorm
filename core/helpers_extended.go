package core

import (
	"context"
	"fmt"
	"reflect"
	"time"
)

// Extended helper functions for common patterns

// FindOrCreateBy finds an entity by a condition or creates it
func FindOrCreateBy[T any, ID comparable](
	ctx context.Context,
	repo Repository[T, ID],
	spec Specification[T],
	creator func(context.Context) (*T, error),
) (*T, error) {
	entity, err := repo.FindOne(ctx, spec)
	if err == nil && entity != nil {
		return entity, nil
	}
	if err != nil && !IsNotFound(err) {
		return nil, err
	}

	return creator(ctx)
}

// UpdateOrCreate updates an entity if it exists, otherwise creates it
func UpdateOrCreate[T any, ID comparable](
	ctx context.Context,
	repo Repository[T, ID],
	spec Specification[T],
	updater func(*T) error,
	creator func(context.Context) (*T, error),
) (*T, error) {
	entity, err := repo.FindOne(ctx, spec)
	if err == nil && entity != nil {
		// Update existing
		if err := updater(entity); err != nil {
			return nil, err
		}
		return repo.Update(ctx, entity)
	}
	if err != nil && !IsNotFound(err) {
		return nil, err
	}

	// Create new
	return creator(ctx)
}

// Upsert saves an entity (insert or update)
func Upsert[T any, ID comparable](
	ctx context.Context,
	repo Repository[T, ID],
	entity *T,
	id ID,
) (*T, error) {
	existing, err := repo.FindByID(ctx, id)
	if err == nil && existing != nil {
		// Update
		return repo.Update(ctx, entity)
	}
	if err != nil && !IsNotFound(err) {
		return nil, err
	}

	// Insert
	return repo.Save(ctx, entity)
}

// DeleteIf deletes entities matching a specification
func DeleteIf[T any, ID comparable](
	ctx context.Context,
	repo Repository[T, ID],
	spec Specification[T],
) (int64, error) {
	return repo.DeleteWithSpec(ctx, spec)
}

// CountIf counts entities matching a specification
func CountIf[T any, ID comparable](
	ctx context.Context,
	repo Repository[T, ID],
	spec Specification[T],
) (int64, error) {
	return repo.CountWithSpec(ctx, spec)
}

// ExistsIf checks if any entity matches a specification
func ExistsIf[T any, ID comparable](
	ctx context.Context,
	repo Repository[T, ID],
	spec Specification[T],
) (bool, error) {
	count, err := repo.CountWithSpec(ctx, spec)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// FindFirstN finds the first N entities matching a specification
func FindFirstN[T any, ID comparable](
	ctx context.Context,
	repo Repository[T, ID],
	spec Specification[T],
	n int,
) ([]*T, error) {
	pageable := PageRequest(0, n)
	page, err := repo.FindAllPagedWithSpec(ctx, spec, pageable)
	if err != nil {
		return nil, err
	}
	return page.Content, nil
}

// FindLastN finds the last N entities matching a specification
func FindLastN[T any, ID comparable](
	ctx context.Context,
	repo Repository[T, ID],
	spec Specification[T],
	n int,
	orderBy string,
) ([]*T, error) {
	order := Order{
		Field:     orderBy,
		Direction: Desc,
	}
	pageable := PageRequest(0, n, order)
	page, err := repo.FindAllPagedWithSpec(ctx, spec, pageable)
	if err != nil {
		return nil, err
	}
	return page.Content, nil
}

// BatchUpdate updates entities in batches
func BatchUpdate[T any, ID comparable](
	ctx context.Context,
	repo Repository[T, ID],
	entities []*T,
	batchSize int,
) error {
	if batchSize <= 0 {
		batchSize = 100
	}

	for i := 0; i < len(entities); i += batchSize {
		end := i + batchSize
		if end > len(entities) {
			end = len(entities)
		}

		batch := entities[i:end]
		for _, entity := range batch {
			if _, err := repo.Update(ctx, entity); err != nil {
				return fmt.Errorf("batch update failed at offset %d: %w", i, err)
			}
		}
	}

	return nil
}

// BatchDelete deletes entities in batches
func BatchDelete[T any, ID comparable](
	ctx context.Context,
	repo Repository[T, ID],
	entities []*T,
	batchSize int,
) error {
	if batchSize <= 0 {
		batchSize = 100
	}

	for i := 0; i < len(entities); i += batchSize {
		end := i + batchSize
		if end > len(entities) {
			end = len(entities)
		}

		batch := entities[i:end]
		for _, entity := range batch {
			if err := repo.Delete(ctx, entity); err != nil {
				return fmt.Errorf("batch delete failed at offset %d: %w", i, err)
			}
		}
	}

	return nil
}

// Transactional executes a function within a transaction
func Transactional[T any](
	ctx context.Context,
	db *Database,
	fn func(context.Context, *Tx) (T, error),
) (T, error) {
	var zero T
	var result T
	var err error

	err = db.Transaction(ctx, func(tx *Tx) error {
		result, err = fn(ctx, tx)
		return err
	})

	if err != nil {
		return zero, err
	}

	return result, nil
}

// WithTimeout executes a function with a timeout
func WithTimeout[T any](
	ctx context.Context,
	timeout time.Duration,
	fn func(context.Context) (T, error),
) (T, error) {
	timeoutCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	return fn(timeoutCtx)
}

// RetryWithContext retries a function with context support
func RetryWithContext[T any](
	ctx context.Context,
	maxAttempts int,
	backoff time.Duration,
	fn func(context.Context) (T, error),
) (T, error) {
	var result T
	var lastErr error
	currentBackoff := backoff

	for attempt := 0; attempt < maxAttempts; attempt++ {
		if attempt > 0 {
			select {
			case <-ctx.Done():
				return result, ctx.Err()
			case <-time.After(currentBackoff):
				currentBackoff *= 2
			}
		}

		var err error
		result, err = fn(ctx)
		if err == nil {
			return result, nil
		}
		lastErr = err
	}

	if lastErr != nil {
		return result, fmt.Errorf("retry failed after %d attempts: %w", maxAttempts, lastErr)
	}
	return result, fmt.Errorf("retry failed after %d attempts", maxAttempts)
}

// MapEntities maps entities to another type
func MapEntities[T, U any](
	entities []*T,
	mapper func(*T) *U,
) []*U {
	result := make([]*U, len(entities))
	for i, entity := range entities {
		result[i] = mapper(entity)
	}
	return result
}

// FilterEntities filters entities based on a predicate
func FilterEntities[T any](
	entities []*T,
	predicate func(*T) bool,
) []*T {
	var result []*T
	for _, entity := range entities {
		if predicate(entity) {
			result = append(result, entity)
		}
	}
	return result
}

// GroupEntities groups entities by a key function
func GroupEntities[T any, K comparable](
	entities []*T,
	keyFn func(*T) K,
) map[K][]*T {
	groups := make(map[K][]*T)
	for _, entity := range entities {
		key := keyFn(entity)
		groups[key] = append(groups[key], entity)
	}
	return groups
}

// DistinctEntities returns distinct entities based on a key function
func DistinctEntities[T any, K comparable](
	entities []*T,
	keyFn func(*T) K,
) []*T {
	seen := make(map[K]bool)
	var result []*T
	for _, entity := range entities {
		key := keyFn(entity)
		if !seen[key] {
			seen[key] = true
			result = append(result, entity)
		}
	}
	return result
}

// SortEntities sorts entities using a comparison function
func SortEntities[T any](
	entities []*T,
	less func(*T, *T) bool,
) []*T {
	// Simple bubble sort (for small slices)
	// For production, use sort.Slice
	result := make([]*T, len(entities))
	copy(result, entities)

	for i := 0; i < len(result)-1; i++ {
		for j := i + 1; j < len(result); j++ {
			if less(result[j], result[i]) {
				result[i], result[j] = result[j], result[i]
			}
		}
	}

	return result
}

// ChunkEntities splits entities into chunks of specified size
func ChunkEntities[T any](
	entities []*T,
	chunkSize int,
) [][]*T {
	if chunkSize <= 0 {
		chunkSize = 100
	}

	var chunks [][]*T
	for i := 0; i < len(entities); i += chunkSize {
		end := i + chunkSize
		if end > len(entities) {
			end = len(entities)
		}
		chunks = append(chunks, entities[i:end])
	}

	return chunks
}

// EntityComparator compares two entities
type EntityComparator[T any] func(*T, *T) int

// CompareEntitiesWithComparator compares entities using a comparator function
func CompareEntitiesWithComparator[T any](
	a, b *T,
	comparator EntityComparator[T],
) int {
	return comparator(a, b)
}

// ExtractField extracts a field value from an entity using reflection
func ExtractField(entity interface{}, fieldName string) (interface{}, error) {
	entityType := reflect.TypeOf(entity)
	if entityType.Kind() == reflect.Ptr {
		entityType = entityType.Elem()
	}

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

// SetField sets a field value on an entity using reflection
func SetField(entity interface{}, fieldName string, value interface{}) error {
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

