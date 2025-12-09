package core

import (
	"testing"
)

func TestSpecification_ToSQL(t *testing.T) {
	t.Run("simple specification", func(t *testing.T) {
		spec := Where[TestUser]("age > $1", 18)
		where, args := spec.ToSQL()

		if where != "age > $1" {
			t.Errorf("Expected 'age > $1', got '%s'", where)
		}
		if len(args) != 1 || args[0] != 18 {
			t.Errorf("Expected args [18], got %v", args)
		}
	})

	t.Run("AND specification", func(t *testing.T) {
		spec1 := Where[TestUser]("age > $1", 18)
		spec2 := Where[TestUser]("status = $1", "active")
		combined := spec1.And(spec2)

		where, args := combined.ToSQL()
		expected := "(age > $1) AND (status = $2)"
		if where != expected {
			t.Errorf("Expected '%s', got '%s'", expected, where)
		}
		if len(args) != 2 {
			t.Errorf("Expected 2 args, got %d", len(args))
		}
		if args[0] != 18 || args[1] != "active" {
			t.Errorf("Expected args [18, 'active'], got %v", args)
		}
	})

	t.Run("OR specification", func(t *testing.T) {
		spec1 := Where[TestUser]("status = $1", "active")
		spec2 := Where[TestUser]("status = $1", "pending")
		combined := spec1.Or(spec2)

		where, args := combined.ToSQL()
		expected := "(status = $1) OR (status = $2)"
		if where != expected {
			t.Errorf("Expected '%s', got '%s'", expected, where)
		}
		if len(args) != 2 {
			t.Errorf("Expected 2 args, got %d", len(args))
		}
	})

	t.Run("NOT specification", func(t *testing.T) {
		spec := Where[TestUser]("active = $1", true)
		negated := spec.Not()

		where, args := negated.ToSQL()
		expected := "NOT (active = $1)"
		if where != expected {
			t.Errorf("Expected '%s', got '%s'", expected, where)
		}
		if len(args) != 1 {
			t.Errorf("Expected 1 arg, got %d", len(args))
		}
	})

	t.Run("complex nested specification", func(t *testing.T) {
		spec1 := Where[TestUser]("age > $1", 18)
		spec2 := Where[TestUser]("status = $1", "active")
		spec3 := Where[TestUser]("email LIKE $1", "%@example.com")

		combined := And(spec1, spec2).Or(spec3)
		where, args := combined.ToSQL()

		// Should be: ((age > $1) AND (status = $2)) OR (email LIKE $3)
		if len(args) != 3 {
			t.Errorf("Expected 3 args, got %d", len(args))
		}
		if !contains(where, "age > $1") || !contains(where, "status = $2") || !contains(where, "email LIKE $3") {
			t.Errorf("Unexpected SQL: %s", where)
		}
	})
}

func TestSpecification_HelperFunctions(t *testing.T) {
	t.Run("Equal", func(t *testing.T) {
		spec := Equal[TestUser]("status", "active")
		where, args := spec.ToSQL()

		if where != "status = $1" {
			t.Errorf("Expected 'status = $1', got '%s'", where)
		}
		if args[0] != "active" {
			t.Errorf("Expected 'active', got %v", args[0])
		}
	})

	t.Run("In", func(t *testing.T) {
		spec := In[TestUser]("status", "active", "pending", "suspended")
		where, args := spec.ToSQL()

		expected := "status IN ($1, $2, $3)"
		if where != expected {
			t.Errorf("Expected '%s', got '%s'", expected, where)
		}
		if len(args) != 3 {
			t.Errorf("Expected 3 args, got %d", len(args))
		}
	})

	t.Run("Between", func(t *testing.T) {
		spec := Between[TestUser]("age", 18, 65)
		where, args := spec.ToSQL()

		expected := "age BETWEEN $1 AND $2"
		if where != expected {
			t.Errorf("Expected '%s', got '%s'", expected, where)
		}
		if len(args) != 2 {
			t.Errorf("Expected 2 args, got %d", len(args))
		}
	})

	t.Run("Contains", func(t *testing.T) {
		spec := Contains[TestUser]("name", "john")
		where, args := spec.ToSQL()

		if where != "name LIKE $1" {
			t.Errorf("Expected 'name LIKE $1', got '%s'", where)
		}
		if args[0] != "%john%" {
			t.Errorf("Expected '%%john%%', got %v", args[0])
		}
	})

	t.Run("StartsWith", func(t *testing.T) {
		spec := StartsWith[TestUser]("email", "john")
		_, args := spec.ToSQL()

		if args[0] != "john%" {
			t.Errorf("Expected 'john%%', got %v", args[0])
		}
	})

	t.Run("EndsWith", func(t *testing.T) {
		spec := EndsWith[TestUser]("email", ".com")
		_, args := spec.ToSQL()

		if args[0] != "%.com" {
			t.Errorf("Expected '%%.com', got %v", args[0])
		}
	})

	t.Run("IsNull", func(t *testing.T) {
		spec := IsNull[TestUser]("deleted_at")
		where, args := spec.ToSQL()

		if where != "deleted_at IS NULL" {
			t.Errorf("Expected 'deleted_at IS NULL', got '%s'", where)
		}
		if len(args) != 0 {
			t.Errorf("Expected 0 args, got %d", len(args))
		}
	})
}

func TestSpecification_AndOr(t *testing.T) {
	t.Run("And with multiple specs", func(t *testing.T) {
		spec1 := Equal[TestUser]("status", "active")
		spec2 := GreaterThan[TestUser]("age", 18)
		spec3 := Equal[TestUser]("email", "test@example.com")

		combined := And(spec1, spec2, spec3)
		where, args := combined.ToSQL()

		if len(args) != 3 {
			t.Errorf("Expected 3 args, got %d", len(args))
		}
		if !contains(where, "AND") {
			t.Errorf("Expected AND in SQL, got: %s", where)
		}
	})

	t.Run("Or with multiple specs", func(t *testing.T) {
		spec1 := Equal[TestUser]("status", "active")
		spec2 := Equal[TestUser]("status", "pending")

		combined := Or(spec1, spec2)
		where, args := combined.ToSQL()

		if len(args) != 2 {
			t.Errorf("Expected 2 args, got %d", len(args))
		}
		if !contains(where, "OR") {
			t.Errorf("Expected OR in SQL, got: %s", where)
		}
	})
}

// contains checks if substr is in s
func contains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

