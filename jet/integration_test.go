package jet

import (
	"context"
	"testing"

	"github.com/go-jet/jet/v2/postgres"
	"github.com/stretchr/testify/assert"
)

// TestJetRepository tests Jet SQL repository integration
func TestJetRepository(t *testing.T) {
	// Note: These tests require actual database connection
	// They are placeholder tests for now

	t.Run("NewJetRepository", func(t *testing.T) {
		// This would require actual repo, db, and table
		// For now, just test that the function exists
		assert.NotNil(t, NewJetRepository)
	})

	t.Run("NewQueryBuilder", func(t *testing.T) {
		// Create a mock table
		table := postgres.NewTable("public", "users", "")
		qb := NewQueryBuilder(table)
		assert.NotNil(t, qb)
		assert.Equal(t, table, qb.table)
	})

	t.Run("QueryBuilder_SelectAll", func(t *testing.T) {
		table := postgres.NewTable("public", "users", "")
		qb := NewQueryBuilder(table)
		stmt := qb.SelectAll()
		assert.NotNil(t, stmt)
	})

	t.Run("QueryBuilder_Select", func(t *testing.T) {
		table := postgres.NewTable("public", "users", "")
		col := postgres.NewStringColumn("email")
		qb := NewQueryBuilder(table)
		stmt := qb.Select(col)
		assert.NotNil(t, stmt)
	})

	t.Run("QueryBuilder_Insert", func(t *testing.T) {
		table := postgres.NewTable("public", "users", "")
		qb := NewQueryBuilder(table)
		stmt := qb.Insert()
		assert.NotNil(t, stmt)
	})

	t.Run("QueryBuilder_Update", func(t *testing.T) {
		table := postgres.NewTable("public", "users", "")
		qb := NewQueryBuilder(table)
		stmt := qb.Update()
		assert.NotNil(t, stmt)
	})

	t.Run("QueryBuilder_Delete", func(t *testing.T) {
		table := postgres.NewTable("public", "users", "")
		qb := NewQueryBuilder(table)
		stmt := qb.Delete()
		assert.NotNil(t, stmt)
	})
}

// TestHelpers tests helper functions
func TestHelpers(t *testing.T) {
	col := postgres.NewStringColumn("email")

	t.Run("Equal", func(t *testing.T) {
		expr := Equal(col, "test@example.com")
		assert.NotNil(t, expr)
	})

	t.Run("NotEqual", func(t *testing.T) {
		expr := NotEqual(col, "test@example.com")
		assert.NotNil(t, expr)
	})

	t.Run("Like", func(t *testing.T) {
		expr := Like(col, "%@example.com")
		assert.NotNil(t, expr)
	})

	t.Run("And", func(t *testing.T) {
		expr1 := Equal(col, "test@example.com")
		expr2 := NotEqual(col, "admin@example.com")
		combined := And(expr1, expr2)
		assert.NotNil(t, combined)
	})

	t.Run("Or", func(t *testing.T) {
		expr1 := Equal(col, "test@example.com")
		expr2 := Equal(col, "admin@example.com")
		combined := Or(expr1, expr2)
		assert.NotNil(t, combined)
	})

	t.Run("In", func(t *testing.T) {
		expr := In(col, "test@example.com", "admin@example.com")
		assert.NotNil(t, expr)
	})

	t.Run("IsNull", func(t *testing.T) {
		expr := IsNull(col)
		assert.NotNil(t, expr)
	})

	t.Run("IsNotNull", func(t *testing.T) {
		expr := IsNotNull(col)
		assert.NotNil(t, expr)
	})
}

// TestAdapters tests adapter functions
func TestAdapters(t *testing.T) {
	table := postgres.NewTable("public", "users", "")

	t.Run("NewSpecificationAdapter", func(t *testing.T) {
		adapter := NewSpecificationAdapter(table)
		assert.NotNil(t, adapter)
		assert.Equal(t, table, adapter.table)
	})

	t.Run("NewQueryBuilderAdapter", func(t *testing.T) {
		adapter := NewQueryBuilderAdapter(table)
		assert.NotNil(t, adapter)
		assert.Equal(t, table, adapter.table)
	})

	t.Run("QueryBuilderAdapter_BuildSelect", func(t *testing.T) {
		adapter := NewQueryBuilderAdapter(table)
		stmt := adapter.BuildSelect()
		assert.NotNil(t, stmt)
	})

	t.Run("NewConditionBuilder", func(t *testing.T) {
		cb := NewConditionBuilder()
		assert.NotNil(t, cb)
		assert.Empty(t, cb.conditions)
	})

	t.Run("ConditionBuilder_Add", func(t *testing.T) {
		cb := NewConditionBuilder()
		col := postgres.NewStringColumn("email")
		expr := Equal(col, "test@example.com")
		cb.Add(expr)
		assert.Len(t, cb.conditions, 1)
	})

	t.Run("ConditionBuilder_And", func(t *testing.T) {
		cb := NewConditionBuilder()
		col := postgres.NewStringColumn("email")
		cb.Add(Equal(col, "test@example.com"))
		cb.Add(NotEqual(col, "admin@example.com"))
		result := cb.And()
		assert.NotNil(t, result)
	})

	t.Run("NewQueryComposer", func(t *testing.T) {
		composer := NewQueryComposer(table)
		assert.NotNil(t, composer)
		assert.Equal(t, table, composer.table)
	})
}

// TestJetQueryExecutor tests query executor
func TestJetQueryExecutor(t *testing.T) {
	t.Run("NewJetQueryExecutor", func(t *testing.T) {
		// This would require actual db connection
		// For now, just test that the function exists
		assert.NotNil(t, NewJetQueryExecutor)
	})
}

// BenchmarkHelpers benchmarks helper functions
func BenchmarkHelpers(b *testing.B) {
	col := postgres.NewStringColumn("email")

	b.Run("Equal", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = Equal(col, "test@example.com")
		}
	})

	b.Run("And", func(b *testing.B) {
		expr1 := Equal(col, "test@example.com")
		expr2 := NotEqual(col, "admin@example.com")
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = And(expr1, expr2)
		}
	})

	b.Run("In", func(b *testing.B) {
		values := []interface{}{"test@example.com", "admin@example.com", "user@example.com"}
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = In(col, values...)
		}
	})
}

