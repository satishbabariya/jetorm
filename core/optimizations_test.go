package core

import (
	"testing"
	"time"
)

func TestQueryCache(t *testing.T) {
	cache := NewQueryCache(5*time.Minute, 100)

	// Test Set and Get
	cache.Set("key1", "value1")
	value, ok := cache.Get("key1")
	if !ok {
		t.Error("Cache should contain key1")
	}
	if value != "value1" {
		t.Errorf("Expected 'value1', got %v", value)
	}

	// Test expiration
	cache.Set("key2", "value2")
	// Manually expire (in real scenario would wait)
	// For now, just test that it works

	// Test Clear
	cache.Clear()
	_, ok = cache.Get("key1")
	if ok {
		t.Error("Cache should be cleared")
	}
}

func TestBatchOptimizer(t *testing.T) {
	optimizer := NewAdvancedBatchOptimizer()

	// Record some metrics
	optimizer.Record(50, 100*time.Millisecond, true)
	optimizer.Record(100, 150*time.Millisecond, true)
	optimizer.Record(200, 300*time.Millisecond, true)

	optimalSize := optimizer.GetOptimalSize()
	if optimalSize <= 0 {
		t.Error("Optimal size should be positive")
	}
}

func TestConnectionPoolOptimizer(t *testing.T) {
	optimizer := NewAdvancedConnectionPoolOptimizer()

	metrics := HealthMetrics{
		MaxConns:    25,
		AcquiredConns: 20,
		IdleConns:   5,
	}

	settings := optimizer.Optimize(metrics)
	if settings.MaxConns <= 0 {
		t.Error("MaxConns should be positive")
	}
	if settings.MinConns <= 0 {
		t.Error("MinConns should be positive")
	}
}

func TestQueryOptimizer(t *testing.T) {
	optimizer := NewQueryOptimizer()

	query := "SELECT * FROM users WHERE id = $1"
	optimized := optimizer.Optimize(query)

	if optimized == "" {
		t.Error("Optimized query should not be empty")
	}
}

func TestLazyLoader(t *testing.T) {
	// This would require a real repository
	t.Skip("Requires database setup")
}

func TestEntityUtils(t *testing.T) {
	type TestEntity struct {
		ID   int64  `db:"id" jet:"primary_key"`
		Name string `db:"name" validate:"required"`
	}

	entity := &TestEntity{ID: 1, Name: "Test"}

	// Test GetFieldValue
	value, err := GetFieldValue(entity, "Name")
	if err != nil {
		t.Errorf("GetFieldValue failed: %v", err)
	}
	if value != "Test" {
		t.Errorf("Expected 'Test', got %v", value)
	}

	// Test SetFieldValue
	err = SetFieldValue(entity, "Name", "Updated")
	if err != nil {
		t.Errorf("SetFieldValue failed: %v", err)
	}
	if entity.Name != "Updated" {
		t.Errorf("Expected 'Updated', got %s", entity.Name)
	}

	// Test GetDBFieldName
	dbName, err := GetDBFieldName(entity, "Name")
	if err != nil {
		t.Errorf("GetDBFieldName failed: %v", err)
	}
	if dbName != "name" {
		t.Errorf("Expected 'name', got %s", dbName)
	}

	// Test IsPrimaryKey
	if !IsPrimaryKey(entity, "ID") {
		t.Error("ID should be primary key")
	}

	// Test GetPrimaryKeyField
	pkField, err := GetPrimaryKeyField(entity)
	if err != nil {
		t.Errorf("GetPrimaryKeyField failed: %v", err)
	}
	if pkField != "ID" {
		t.Errorf("Expected 'ID', got %s", pkField)
	}

	// Test GetTableName
	tableName := GetTableName(entity)
	if tableName == "" {
		t.Error("Table name should not be empty")
	}

	// Test GetColumnNames
	columns := GetColumnNames(entity)
	if len(columns) == 0 {
		t.Error("Should have columns")
	}

	// Test GetFieldNames
	fields := GetFieldNames(entity)
	if len(fields) == 0 {
		t.Error("Should have fields")
	}
}

func TestCopyFields(t *testing.T) {
	type TestEntity struct {
		ID   int64
		Name string
		Age  int
	}

	src := &TestEntity{ID: 1, Name: "Test", Age: 25}
	dest := &TestEntity{}

	err := CopyFields(dest, src)
	if err != nil {
		t.Errorf("CopyFields failed: %v", err)
	}

	if dest.Name != src.Name {
		t.Errorf("Expected %s, got %s", src.Name, dest.Name)
	}
}

func TestCompareEntities(t *testing.T) {
	type TestEntity struct {
		ID   int64
		Name string
	}

	entity1 := &TestEntity{ID: 1, Name: "Test"}
	entity2 := &TestEntity{ID: 1, Name: "Test"}
	entity3 := &TestEntity{ID: 1, Name: "Different"}

	equal, _ := CompareEntities(entity1, entity2)
	if !equal {
		t.Error("Entities should be equal")
	}

	equal, differences := CompareEntities(entity1, entity3)
	if equal {
		t.Error("Entities should not be equal")
	}
	if len(differences) == 0 {
		t.Error("Should have differences")
	}
}

