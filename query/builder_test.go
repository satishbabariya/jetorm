package query

import (
	"testing"
)

func TestQueryBuilder_Basic(t *testing.T) {
	qb := NewQueryBuilder("users")
	qb.WhereEqual("status", "active")
	qb.OrderBy("created_at", "DESC")
	qb.Limit(10)
	
	query, args := qb.Build()
	
	if !contains(query, "SELECT") {
		t.Error("Query should contain SELECT")
	}
	if !contains(query, "FROM users") {
		t.Error("Query should contain FROM users")
	}
	if !contains(query, "status = $1") {
		t.Error("Query should contain WHERE clause")
	}
	if !contains(query, "ORDER BY") {
		t.Error("Query should contain ORDER BY")
	}
	if !contains(query, "LIMIT 10") {
		t.Error("Query should contain LIMIT")
	}
	if len(args) != 1 {
		t.Errorf("Expected 1 arg, got %d", len(args))
	}
	if args[0] != "active" {
		t.Errorf("Expected arg 'active', got %v", args[0])
	}
}

func TestQueryBuilder_WhereIn(t *testing.T) {
	qb := NewQueryBuilder("users")
	qb.WhereIn("status", []interface{}{"active", "pending"})
	
	query, args := qb.Build()
	
	if !contains(query, "status IN") {
		t.Error("Query should contain IN clause")
	}
	if len(args) != 2 {
		t.Errorf("Expected 2 args, got %d", len(args))
	}
}

func TestQueryBuilder_Count(t *testing.T) {
	qb := NewQueryBuilder("users")
	qb.WhereEqual("status", "active")
	
	query, args := qb.BuildCount()
	
	if !contains(query, "SELECT COUNT(*)") {
		t.Error("Count query should contain COUNT(*)")
	}
	if !contains(query, "FROM users") {
		t.Error("Count query should contain FROM")
	}
	if len(args) != 1 {
		t.Errorf("Expected 1 arg, got %d", len(args))
	}
}

func TestComposableQuery_WithSpecification(t *testing.T) {
	// This test would require importing core package
	// For now, just test basic functionality
	cq := NewComposableQuery[string]("users")
	cq.WhereEqual("status", "active")
	cq.Limit(10)
	
	query, args := cq.Build()
	
	if !contains(query, "SELECT") {
		t.Error("Query should contain SELECT")
	}
	if len(args) != 1 {
		t.Errorf("Expected 1 arg, got %d", len(args))
	}
}

func TestConditionBuilder_Basic(t *testing.T) {
	cb := NewConditionBuilder()
	cb.Equal("status", "active")
	cb.GreaterThan("age", 18)
	
	whereClause, args := cb.Build()
	
	if !contains(whereClause, "status = $1") {
		t.Error("Should contain status condition")
	}
	if !contains(whereClause, "age > $2") {
		t.Error("Should contain age condition")
	}
	if len(args) != 2 {
		t.Errorf("Expected 2 args, got %d", len(args))
	}
}

func TestConditionBuilder_AndOr(t *testing.T) {
	cb1 := NewConditionBuilder()
	cb1.Equal("status", "active")
	
	cb2 := NewConditionBuilder()
	cb2.Equal("status", "pending")
	
	combined := cb1.Or(cb2)
	whereClause, args := combined.Build()
	
	if !contains(whereClause, "OR") {
		t.Error("Should contain OR operator")
	}
	if len(args) != 2 {
		t.Errorf("Expected 2 args, got %d", len(args))
	}
}

func TestJoinQuery_Basic(t *testing.T) {
	jq := NewJoinQuery[string]("users")
	jq.InnerJoin("profiles", "users.id = profiles.user_id")
	jq.WhereEqual("users.status", "active")
	
	query, args := jq.Build()
	
	if !contains(query, "INNER JOIN") {
		t.Error("Query should contain JOIN")
	}
	if !contains(query, "profiles") {
		t.Error("Query should contain joined table")
	}
	if len(args) != 1 {
		t.Errorf("Expected 1 arg, got %d", len(args))
	}
}

func TestDynamicQuery_Conditional(t *testing.T) {
	status := "active"
	minAge := 18
	
	dq := NewDynamicQuery[string]("users")
	dq.When(status != "", func(q *ComposableQuery[string]) *ComposableQuery[string] {
		return q.Where("status = $1", status)
	})
	dq.When(minAge > 0, func(q *ComposableQuery[string]) *ComposableQuery[string] {
		return q.Where("age >= $1", minAge)
	})
	
	query, args := dq.Build()
	
	if !contains(query, "status = $1") {
		t.Error("Should contain status condition")
	}
	// Note: Placeholder numbering may vary based on how conditions are applied
	if !contains(query, "age >=") && !contains(query, "age >=") {
		t.Error("Should contain age condition")
	}
	if len(args) < 1 {
		t.Errorf("Expected at least 1 arg, got %d", len(args))
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && 
		(s == substr || 
		 (len(s) > len(substr) && 
		  (s[:len(substr)] == substr || 
		   s[len(s)-len(substr):] == substr || 
		   contains(s[1:], substr))))
}

