package core

import (
	"context"
	"fmt"
	"reflect"
	"strings"
	"time"
)

// Convenience functions for common operations

// Must panics if error is not nil
func Must[T any](value T, err error) T {
	if err != nil {
		panic(fmt.Sprintf("must: %v", err))
	}
	return value
}

// Retry retries an operation with exponential backoff
func Retry(ctx context.Context, maxAttempts int, fn func() error) error {
	var lastErr error
	backoff := time.Second

	for attempt := 0; attempt < maxAttempts; attempt++ {
		if attempt > 0 {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(backoff):
				backoff *= 2 // Exponential backoff
			}
		}

		err := fn()
		if err == nil {
			return nil
		}
		lastErr = err
	}

	return fmt.Errorf("retry failed after %d attempts: %w", maxAttempts, lastErr)
}

// RetryWithBackoff retries with custom backoff
func RetryWithBackoff(ctx context.Context, maxAttempts int, initialBackoff time.Duration, fn func() error) error {
	var lastErr error
	backoff := initialBackoff

	for attempt := 0; attempt < maxAttempts; attempt++ {
		if attempt > 0 {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(backoff):
				backoff *= 2
			}
		}

		err := fn()
		if err == nil {
			return nil
		}
		lastErr = err
	}

	return fmt.Errorf("retry failed after %d attempts: %w", maxAttempts, lastErr)
}

// Timeout executes a function with timeout
func Timeout(ctx context.Context, timeout time.Duration, fn func(context.Context) error) error {
	timeoutCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	return fn(timeoutCtx)
}

// Parallel executes functions in parallel and returns first error
func Parallel(fns ...func() error) error {
	errChan := make(chan error, len(fns))

	for _, fn := range fns {
		go func(f func() error) {
			errChan <- f()
		}(fn)
	}

	for i := 0; i < len(fns); i++ {
		if err := <-errChan; err != nil {
			return err
		}
	}

	return nil
}

// ExtractID extracts ID from entity using reflection
func ExtractID[T any, ID comparable](entity *T) (ID, error) {
	var zeroID ID
	entityType := reflect.TypeOf(entity)
	if entityType.Kind() == reflect.Ptr {
		entityType = entityType.Elem()
	}

	entityValue := reflect.ValueOf(entity)
	if entityValue.Kind() == reflect.Ptr {
		entityValue = entityValue.Elem()
	}

	// Find primary key field
	for i := 0; i < entityType.NumField(); i++ {
		field := entityType.Field(i)
		jetTag := field.Tag.Get("jet")
		if strings.Contains(jetTag, "primary_key") {
			fieldValue := entityValue.Field(i)
			if fieldValue.CanInterface() {
				if id, ok := fieldValue.Interface().(ID); ok {
					return id, nil
				}
			}
		}
	}

	return zeroID, fmt.Errorf("could not extract ID from entity")
}

// SetID sets ID on entity using reflection
func SetID[T any, ID comparable](entity *T, id ID) error {
	entityType := reflect.TypeOf(entity)
	if entityType.Kind() == reflect.Ptr {
		entityType = entityType.Elem()
	}

	entityValue := reflect.ValueOf(entity)
	if entityValue.Kind() == reflect.Ptr {
		entityValue = entityValue.Elem()
	}

	// Find primary key field
	for i := 0; i < entityType.NumField(); i++ {
		field := entityType.Field(i)
		jetTag := field.Tag.Get("jet")
		if strings.Contains(jetTag, "primary_key") {
			fieldValue := entityValue.Field(i)
			if fieldValue.CanSet() {
				idValue := reflect.ValueOf(id)
				if idValue.Type().AssignableTo(fieldValue.Type()) {
					fieldValue.Set(idValue)
					return nil
				}
			}
		}
	}

	return fmt.Errorf("could not set ID on entity")
}

// IsZero checks if a value is zero
func IsZero(value interface{}) bool {
	if value == nil {
		return true
	}

	v := reflect.ValueOf(value)
	switch v.Kind() {
	case reflect.Ptr, reflect.Interface:
		return v.IsNil()
	case reflect.String:
		return v.Len() == 0
	case reflect.Slice, reflect.Array, reflect.Map:
		return v.Len() == 0
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Bool:
		return !v.Bool()
	}

	return false
}

// Coalesce returns the first non-zero value
func Coalesce[T comparable](values ...T) T {
	var zero T
	for _, v := range values {
		if v != zero {
			return v
		}
	}
	return zero
}

// DefaultIfZero returns default value if value is zero
func DefaultIfZero[T comparable](value, defaultValue T) T {
	var zero T
	if value == zero {
		return defaultValue
	}
	return value
}

// SliceContains checks if slice contains value
func SliceContains[T comparable](slice []T, value T) bool {
	for _, v := range slice {
		if v == value {
			return true
		}
	}
	return false
}

// SliceMap maps a slice to another type
func SliceMap[T, U any](slice []T, fn func(T) U) []U {
	result := make([]U, len(slice))
	for i, v := range slice {
		result[i] = fn(v)
	}
	return result
}

// SliceFilter filters a slice
func SliceFilter[T any](slice []T, fn func(T) bool) []T {
	var result []T
	for _, v := range slice {
		if fn(v) {
			result = append(result, v)
		}
	}
	return result
}

// SliceUnique removes duplicates from slice
func SliceUnique[T comparable](slice []T) []T {
	seen := make(map[T]bool)
	var result []T
	for _, v := range slice {
		if !seen[v] {
			seen[v] = true
			result = append(result, v)
		}
	}
	return result
}

// MapKeys returns keys of a map
func MapKeys[K comparable, V any](m map[K]V) []K {
	keys := make([]K, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

// MapValues returns values of a map
func MapValues[K comparable, V any](m map[K]V) []V {
	values := make([]V, 0, len(m))
	for _, v := range m {
		values = append(values, v)
	}
	return values
}

// MapEntries returns key-value pairs from a map
func MapEntries[K comparable, V any](m map[K]V) []struct {
	Key   K
	Value V
} {
	entries := make([]struct {
		Key   K
		Value V
	}, 0, len(m))
	for k, v := range m {
		entries = append(entries, struct {
			Key   K
			Value V
		}{k, v})
	}
	return entries
}

