package core

import (
	"context"
	"fmt"
	"reflect"
	"strings"
	"time"
)

// Extended utility functions

// RetryWithCondition retries with a condition function
func RetryWithCondition(
	ctx context.Context,
	maxAttempts int,
	backoff time.Duration,
	fn func(context.Context) error,
	condition func(error) bool,
) error {
	var lastErr error
	currentBackoff := backoff

	for attempt := 0; attempt < maxAttempts; attempt++ {
		if attempt > 0 {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(currentBackoff):
				currentBackoff *= 2
			}
		}

		err := fn(ctx)
		if err == nil {
			return nil
		}

		// Check if we should retry
		if condition != nil && !condition(err) {
			return err
		}

		lastErr = err
	}

	return fmt.Errorf("retry failed after %d attempts: %w", maxAttempts, lastErr)
}

// ParallelWithLimit executes functions in parallel with a limit
func ParallelWithLimit(limit int, fns ...func() error) error {
	if limit <= 0 {
		limit = len(fns)
	}

	semaphore := make(chan struct{}, limit)
	errChan := make(chan error, len(fns))

	for _, fn := range fns {
		go func(f func() error) {
			semaphore <- struct{}{}
			defer func() { <-semaphore }()
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

// Debounce debounces function calls
func Debounce(duration time.Duration, fn func()) func() {
	var timer *time.Timer
	return func() {
		if timer != nil {
			timer.Stop()
		}
		timer = time.AfterFunc(duration, fn)
	}
}

// Throttle throttles function calls
func Throttle(duration time.Duration, fn func()) func() {
	var lastCall time.Time
	return func() {
		now := time.Now()
		if now.Sub(lastCall) >= duration {
			lastCall = now
			fn()
		}
	}
}

// Memoize memoizes function results
func Memoize[T comparable, R any](fn func(T) R) func(T) R {
	cache := make(map[T]R)
	return func(input T) R {
		if result, ok := cache[input]; ok {
			return result
		}
		result := fn(input)
		cache[input] = result
		return result
	}
}

// Pipeline creates a pipeline of functions
func Pipeline[T any](fns ...func(T) T) func(T) T {
	return func(input T) T {
		result := input
		for _, fn := range fns {
			result = fn(result)
		}
		return result
	}
}

// Compose composes functions
func Compose[T any](fns ...func(T) T) func(T) T {
	return Pipeline(fns...)
}

// Chain chains operations
func Chain[T any](value T, fns ...func(T) T) T {
	return Pipeline(fns...)(value)
}

// Transform transforms a slice
func Transform[T, U any](slice []T, fn func(T) U) []U {
	result := make([]U, len(slice))
	for i, v := range slice {
		result[i] = fn(v)
	}
	return result
}

// Reduce reduces a slice
func Reduce[T, U any](slice []T, initial U, fn func(U, T) U) U {
	result := initial
	for _, v := range slice {
		result = fn(result, v)
	}
	return result
}

// Partition partitions a slice
func Partition[T any](slice []T, fn func(T) bool) ([]T, []T) {
	var truePart, falsePart []T
	for _, v := range slice {
		if fn(v) {
			truePart = append(truePart, v)
		} else {
			falsePart = append(falsePart, v)
		}
	}
	return truePart, falsePart
}

// Zip zips two slices
func Zip[T, U any](slice1 []T, slice2 []U) []struct {
	First  T
	Second U
} {
	minLen := len(slice1)
	if len(slice2) < minLen {
		minLen = len(slice2)
	}

	result := make([]struct {
		First  T
		Second U
	}, minLen)

	for i := 0; i < minLen; i++ {
		result[i] = struct {
			First  T
			Second U
		}{slice1[i], slice2[i]}
	}

	return result
}

// Unzip unzips a slice of pairs
func Unzip[T, U any](pairs []struct {
	First  T
	Second U
}) ([]T, []U) {
	firsts := make([]T, len(pairs))
	seconds := make([]U, len(pairs))

	for i, pair := range pairs {
		firsts[i] = pair.First
		seconds[i] = pair.Second
	}

	return firsts, seconds
}

// Flatten flattens a slice of slices
func Flatten[T any](slices [][]T) []T {
	var result []T
	for _, slice := range slices {
		result = append(result, slice...)
	}
	return result
}

// Chunk chunks a slice
func Chunk[T any](slice []T, size int) [][]T {
	if size <= 0 {
		size = 1
	}

	var chunks [][]T
	for i := 0; i < len(slice); i += size {
		end := i + size
		if end > len(slice) {
			end = len(slice)
		}
		chunks = append(chunks, slice[i:end])
	}

	return chunks
}

// Intersect finds intersection of two slices
func Intersect[T comparable](slice1, slice2 []T) []T {
	set := make(map[T]bool)
	for _, v := range slice1 {
		set[v] = true
	}

	var result []T
	for _, v := range slice2 {
		if set[v] {
			result = append(result, v)
			delete(set, v) // Avoid duplicates
		}
	}

	return result
}

// Difference finds difference of two slices
func Difference[T comparable](slice1, slice2 []T) []T {
	set := make(map[T]bool)
	for _, v := range slice2 {
		set[v] = true
	}

	var result []T
	for _, v := range slice1 {
		if !set[v] {
			result = append(result, v)
		}
	}

	return result
}

// Union finds union of two slices
func Union[T comparable](slice1, slice2 []T) []T {
	set := make(map[T]bool)
	var result []T

	for _, v := range slice1 {
		if !set[v] {
			set[v] = true
			result = append(result, v)
		}
	}

	for _, v := range slice2 {
		if !set[v] {
			set[v] = true
			result = append(result, v)
		}
	}

	return result
}

// AllMatch checks if all elements satisfy a predicate
func AllMatch[T any](slice []T, predicate func(T) bool) bool {
	for _, v := range slice {
		if !predicate(v) {
			return false
		}
	}
	return true
}

// Any checks if any element satisfies a predicate
func Any[T any](slice []T, predicate func(T) bool) bool {
	for _, v := range slice {
		if predicate(v) {
			return true
		}
	}
	return false
}

// None checks if no elements satisfy a predicate
func None[T any](slice []T, predicate func(T) bool) bool {
	return !Any(slice, predicate)
}

// Count counts elements satisfying a predicate
func Count[T any](slice []T, predicate func(T) bool) int {
	count := 0
	for _, v := range slice {
		if predicate(v) {
			count++
		}
	}
	return count
}

// First finds first element satisfying a predicate
func First[T any](slice []T, predicate func(T) bool) (T, bool) {
	var zero T
	for _, v := range slice {
		if predicate(v) {
			return v, true
		}
	}
	return zero, false
}

// Last finds last element satisfying a predicate
func Last[T any](slice []T, predicate func(T) bool) (T, bool) {
	var zero T
	for i := len(slice) - 1; i >= 0; i-- {
		if predicate(slice[i]) {
			return slice[i], true
		}
	}
	return zero, false
}

// Take takes first n elements
func Take[T any](slice []T, n int) []T {
	if n <= 0 {
		return []T{}
	}
	if n >= len(slice) {
		return slice
	}
	return slice[:n]
}

// Drop drops first n elements
func Drop[T any](slice []T, n int) []T {
	if n <= 0 {
		return slice
	}
	if n >= len(slice) {
		return []T{}
	}
	return slice[n:]
}

// TakeWhile takes elements while predicate is true
func TakeWhile[T any](slice []T, predicate func(T) bool) []T {
	var result []T
	for _, v := range slice {
		if !predicate(v) {
			break
		}
		result = append(result, v)
	}
	return result
}

// DropWhile drops elements while predicate is true
func DropWhile[T any](slice []T, predicate func(T) bool) []T {
	var result []T
	drop := true
	for _, v := range slice {
		if drop && !predicate(v) {
			drop = false
		}
		if !drop {
			result = append(result, v)
		}
	}
	return result
}

// Reverse reverses a slice
func Reverse[T any](slice []T) []T {
	result := make([]T, len(slice))
	for i, v := range slice {
		result[len(slice)-1-i] = v
	}
	return result
}

// Shuffle shuffles a slice
func Shuffle[T any](slice []T, rng func() int) []T {
	result := make([]T, len(slice))
	copy(result, slice)

	for i := len(result) - 1; i > 0; i-- {
		j := rng() % (i + 1)
		result[i], result[j] = result[j], result[i]
	}

	return result
}

// GroupBy groups elements by a key function
func GroupBy[T any, K comparable](slice []T, keyFn func(T) K) map[K][]T {
	groups := make(map[K][]T)
	for _, v := range slice {
		key := keyFn(v)
		groups[key] = append(groups[key], v)
	}
	return groups
}

// IndexBy indexes elements by a key function
func IndexBy[T any, K comparable](slice []T, keyFn func(T) K) map[K]T {
	index := make(map[K]T)
	for _, v := range slice {
		key := keyFn(v)
		index[key] = v
	}
	return index
}

// Sum sums numeric values
func Sum[T NumericType](slice []T) T {
	var sum T
	for _, v := range slice {
		sum += v
	}
	return sum
}

// Average calculates average of numeric values
func Average[T NumericType](slice []T) float64 {
	if len(slice) == 0 {
		return 0
	}
	sum := Sum(slice)
	return float64(sum) / float64(len(slice))
}

// MinValue finds minimum value in a slice
func MinValue[T Ordered](slice []T) T {
	if len(slice) == 0 {
		var zero T
		return zero
	}
	min := slice[0]
	for _, v := range slice[1:] {
		if v < min {
			min = v
		}
	}
	return min
}

// MaxValue finds maximum value in a slice
func MaxValue[T Ordered](slice []T) T {
	if len(slice) == 0 {
		var zero T
		return zero
	}
	max := slice[0]
	for _, v := range slice[1:] {
		if v > max {
			max = v
		}
	}
	return max
}

// NumericType represents numeric types
type NumericType interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 |
		~float32 | ~float64
}

// Ordered represents ordered types
type Ordered interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 |
		~float32 | ~float64 | ~string
}

// GetStructTags gets all struct tags for a type
func GetStructTags(entity interface{}, tagName string) map[string]string {
	tags := make(map[string]string)
	entityType := reflect.TypeOf(entity)
	if entityType.Kind() == reflect.Ptr {
		entityType = entityType.Elem()
	}

	for i := 0; i < entityType.NumField(); i++ {
		field := entityType.Field(i)
		if !field.IsExported() {
			continue
		}

		tag := field.Tag.Get(tagName)
		if tag != "" {
			tags[field.Name] = tag
		}
	}

	return tags
}

// GetFieldByTag gets field name by tag value
func GetFieldByTag(entity interface{}, tagName, tagValue string) (string, bool) {
	entityType := reflect.TypeOf(entity)
	if entityType.Kind() == reflect.Ptr {
		entityType = entityType.Elem()
	}

	for i := 0; i < entityType.NumField(); i++ {
		field := entityType.Field(i)
		tag := field.Tag.Get(tagName)
		if strings.Contains(tag, tagValue) {
			return field.Name, true
		}
	}

	return "", false
}

// ValidateStruct validates a struct using tags
func ValidateStruct(entity interface{}) error {
	validator := NewValidator()
	return validator.Validate(entity)
}

// DeepEqual compares two values deeply
func DeepEqual(a, b interface{}) bool {
	return reflect.DeepEqual(a, b)
}

// TypeOf gets the type of a value
func TypeOf(value interface{}) reflect.Type {
	return reflect.TypeOf(value)
}

// ValueOf gets the value of a value
func ValueOf(value interface{}) reflect.Value {
	return reflect.ValueOf(value)
}

// IsNil checks if a value is nil
func IsNil(value interface{}) bool {
	if value == nil {
		return true
	}
	v := reflect.ValueOf(value)
	switch v.Kind() {
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Ptr, reflect.Slice:
		return v.IsNil()
	}
	return false
}

// IsZeroValue checks if a value is its zero value
func IsZeroValue(value interface{}) bool {
	return IsZero(value)
}

