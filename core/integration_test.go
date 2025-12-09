package core

import (
	"context"
	"testing"
	"time"
)

// Integration tests for feature combinations

func TestCachedRepositoryWithHooks(t *testing.T) {
	// This would require a real database connection
	// Skipping for now but structure is here
	t.Skip("Requires database setup")
}

func TestRepositoryWithValidation(t *testing.T) {
	// Mock repository for testing
	type MockRepo struct {
		saved bool
	}

	validator := NewValidator()
	validator.RegisterRule("Email", Email())
	validator.RegisterRule("Name", Required())

	// Test validation passes
	entity := struct {
		Email string `validate:"email"`
		Name  string `validate:"required"`
	}{
		Email: "test@example.com",
		Name:  "Test",
	}

	err := validator.Validate(entity)
	if err != nil {
		t.Errorf("Validation should pass: %v", err)
	}

	// Test validation fails
	invalidEntity := struct {
		Email string `validate:"email"`
		Name  string `validate:"required"`
	}{
		Email: "invalid-email",
		Name:  "",
	}

	err = validator.Validate(invalidEntity)
	if err == nil {
		t.Error("Validation should fail")
	}
}

func TestRepositoryWithMetrics(t *testing.T) {
	monitor := NewPerformanceMonitor(100 * time.Millisecond)
	profiler := NewQueryProfiler(monitor)

	// Record some queries
	profiler.Profile(context.Background(), "test_query", func(ctx context.Context) error {
		time.Sleep(50 * time.Millisecond)
		return nil
	})

	metrics := monitor.GetMetrics("test_query")
	if metrics == nil {
		t.Error("Metrics should be recorded")
		return
	}

	if metrics.Count != 1 {
		t.Errorf("Expected count 1, got %d", metrics.Count)
	}
}

func TestFullFeaturedRepository_HealthCheck(t *testing.T) {
	// This would require a real database connection
	t.Skip("Requires database setup")
}

func TestHelperFunctions(t *testing.T) {
	// Test IsZero
	if !IsZero(0) {
		t.Error("0 should be zero")
	}
	if IsZero(1) {
		t.Error("1 should not be zero")
	}
	if !IsZero("") {
		t.Error("Empty string should be zero")
	}
	if IsZero("test") {
		t.Error("Non-empty string should not be zero")
	}

	// Test Coalesce
	if Coalesce(0, 0, 5) != 5 {
		t.Error("Coalesce should return first non-zero value")
	}

	// Test DefaultIfZero
	if DefaultIfZero(0, 10) != 10 {
		t.Error("DefaultIfZero should return default for zero")
	}
	if DefaultIfZero(5, 10) != 5 {
		t.Error("DefaultIfZero should return value if not zero")
	}

	// Test SliceContains
	slice := []int{1, 2, 3, 4, 5}
	if !SliceContains(slice, 3) {
		t.Error("SliceContains should find existing element")
	}
	if SliceContains(slice, 6) {
		t.Error("SliceContains should not find non-existing element")
	}

	// Test SliceMap
	doubled := SliceMap([]int{1, 2, 3}, func(x int) int {
		return x * 2
	})
	if len(doubled) != 3 || doubled[0] != 2 || doubled[1] != 4 || doubled[2] != 6 {
		t.Error("SliceMap should map correctly")
	}

	// Test SliceFilter
	filtered := SliceFilter([]int{1, 2, 3, 4, 5}, func(x int) bool {
		return x%2 == 0
	})
	if len(filtered) != 2 || filtered[0] != 2 || filtered[1] != 4 {
		t.Error("SliceFilter should filter correctly")
	}

	// Test SliceUnique
	unique := SliceUnique([]int{1, 2, 2, 3, 3, 3})
	if len(unique) != 3 {
		t.Errorf("SliceUnique should remove duplicates, got %d elements", len(unique))
	}
}

func TestValidationRules(t *testing.T) {
	// Test MinLength
	rule := MinLength(5)
	if err := rule("test"); err == nil {
		t.Error("MinLength should fail for short string")
	}
	if err := rule("testing"); err != nil {
		t.Error("MinLength should pass for long string")
	}

	// Test MaxLength
	rule = MaxLength(5)
	if err := rule("testing"); err == nil {
		t.Error("MaxLength should fail for long string")
	}
	if err := rule("test"); err != nil {
		t.Error("MaxLength should pass for short string")
	}

	// Test Range
	rule = Range(1, 10)
	if err := rule(5); err != nil {
		t.Error("Range should pass for value in range")
	}
	if err := rule(15); err == nil {
		t.Error("Range should fail for value out of range")
	}

	// Test InList
	rule = InList(1, 2, 3)
	if err := rule(2); err != nil {
		t.Error("InList should pass for value in list")
	}
	if err := rule(5); err == nil {
		t.Error("InList should fail for value not in list")
	}

	// Test Positive
	rule = Positive()
	if err := rule(5); err != nil {
		t.Error("Positive should pass for positive number")
	}
	if err := rule(-5); err == nil {
		t.Error("Positive should fail for negative number")
	}
}

func TestEntityHelpers(t *testing.T) {
	type TestEntity struct {
		ID   int64
		Name string
		Age  int
	}

	entity1 := &TestEntity{ID: 1, Name: "Test", Age: 25}

	// Test ExtractField
	name, err := ExtractField(entity1, "Name")
	if err != nil {
		t.Errorf("ExtractField failed: %v", err)
	}
	if name != "Test" {
		t.Errorf("Expected 'Test', got %v", name)
	}

	// Test SetField
	err = SetField(entity1, "Name", "Updated")
	if err != nil {
		t.Errorf("SetField failed: %v", err)
	}
	if entity1.Name != "Updated" {
		t.Errorf("Expected 'Updated', got %s", entity1.Name)
	}

	// Test EntityHelper Copy
	helper := NewEntityHelper()
	dest := &TestEntity{}
	err = helper.Copy(dest, entity1)
	if err != nil {
		t.Errorf("Copy failed: %v", err)
	}
	if dest.Name != entity1.Name {
		t.Errorf("Copy should copy Name field")
	}

	// Test EntityHelper Clone
	cloned := helper.Clone(entity1)
	clonedEntity, ok := cloned.(*TestEntity)
	if !ok {
		t.Error("Clone should return correct type")
	}
	if clonedEntity.Name != entity1.Name {
		t.Errorf("Clone should copy Name field")
	}
}

func TestChunkEntities(t *testing.T) {
	entities := []*struct{ ID int }{
		{ID: 1}, {ID: 2}, {ID: 3}, {ID: 4}, {ID: 5},
	}

	chunks := ChunkEntities(entities, 2)
	if len(chunks) != 3 {
		t.Errorf("Expected 3 chunks, got %d", len(chunks))
	}
	if len(chunks[0]) != 2 {
		t.Errorf("Expected chunk size 2, got %d", len(chunks[0]))
	}
}

func TestGroupEntities(t *testing.T) {
	entities := []*struct {
		ID   int
		Type string
	}{
		{ID: 1, Type: "A"},
		{ID: 2, Type: "B"},
		{ID: 3, Type: "A"},
	}

	groups := GroupEntities(entities, func(e *struct {
		ID   int
		Type string
	}) string {
		return e.Type
	})

	if len(groups) != 2 {
		t.Errorf("Expected 2 groups, got %d", len(groups))
	}
	if len(groups["A"]) != 2 {
		t.Errorf("Expected 2 entities in group A, got %d", len(groups["A"]))
	}
}

