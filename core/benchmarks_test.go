package core

import (
	"testing"
)

// Benchmark helpers for performance testing

func BenchmarkRepository_Save(b *testing.B) {
	// Setup would create test database and repository
	b.Skip("Requires test database setup")
}

func BenchmarkRepository_FindByID(b *testing.B) {
	b.Skip("Requires test database setup")
}

func BenchmarkRepository_FindAll(b *testing.B) {
	b.Skip("Requires test database setup")
}

func BenchmarkCachedRepository_FindByID(b *testing.B) {
	b.Skip("Requires test database setup")
}

func BenchmarkBatchWriter_Write(b *testing.B) {
	b.Skip("Requires test database setup")
}

// Benchmark utility functions
func BenchmarkExtractID(b *testing.B) {
	type TestEntity struct {
		ID int64 `db:"id" jet:"primary_key"`
	}

	entity := &TestEntity{ID: 123}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = ExtractID[TestEntity, int64](entity)
	}
}

func BenchmarkIsZero(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = IsZero(123)
		_ = IsZero("")
		_ = IsZero(nil)
	}
}

func BenchmarkSliceContains(b *testing.B) {
	slice := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = SliceContains(slice, 5)
	}
}

func BenchmarkSliceMap(b *testing.B) {
	slice := make([]int, 1000)
	for i := range slice {
		slice[i] = i
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = SliceMap(slice, func(x int) int {
			return x * 2
		})
	}
}

