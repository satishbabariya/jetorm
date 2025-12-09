package generator

import (
	"reflect"
	"testing"
)

type TestUser struct {
	ID        int64
	Email     string
	Username  string
	Age       int
	Status    string
	IsActive  bool
	CreatedAt string
}

func TestAnalyzer_AnalyzeMethod(t *testing.T) {
	entityType := reflect.TypeOf(TestUser{})
	analyzer, err := NewAnalyzer(entityType)
	if err != nil {
		t.Fatalf("Failed to create analyzer: %v", err)
	}

	tests := []struct {
		name           string
		methodName     string
		expectedOp     Operation
		expectedFields int
		checkSQL       bool
	}{
		{
			name:           "FindByEmail",
			methodName:     "FindByEmail",
			expectedOp:     OpFind,
			expectedFields: 1,
		},
		{
			name:           "FindByEmailAndAge",
			methodName:     "FindByEmailAndAge",
			expectedOp:     OpFind,
			expectedFields: 2,
		},
		{
			name:           "FindByAgeGreaterThan",
			methodName:     "FindByAgeGreaterThan",
			expectedOp:     OpFind,
			expectedFields: 1,
		},
		{
			name:           "FindByStatusIn",
			methodName:     "FindByStatusIn",
			expectedOp:     OpFind,
			expectedFields: 1,
		},
		{
			name:           "FindFirstByStatus",
			methodName:     "FindFirstByStatus",
			expectedOp:     OpFind,
			expectedFields: 1,
		},
		{
			name:           "CountByStatus",
			methodName:     "CountByStatus",
			expectedOp:     OpCount,
			expectedFields: 1,
		},
		{
			name:           "ExistsByEmail",
			methodName:     "ExistsByEmail",
			expectedOp:     OpExists,
			expectedFields: 1,
		},
		{
			name:           "DeleteByEmail",
			methodName:     "DeleteByEmail",
			expectedOp:     OpDelete,
			expectedFields: 1,
		},
		{
			name:           "FindByStatusOrderByCreatedAtDesc",
			methodName:     "FindByStatusOrderByCreatedAtDesc",
			expectedOp:     OpFind,
			expectedFields: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			method, err := analyzer.AnalyzeMethod(tt.methodName)
			if err != nil {
				t.Fatalf("Failed to analyze method: %v", err)
			}

			if method.Operation != tt.expectedOp {
				t.Errorf("Expected operation %v, got %v", tt.expectedOp, method.Operation)
			}

			if len(method.Fields) != tt.expectedFields {
				t.Errorf("Expected %d fields, got %d", tt.expectedFields, len(method.Fields))
			}

			if method.Name != tt.methodName {
				t.Errorf("Expected method name %s, got %s", tt.methodName, method.Name)
			}
		})
	}
}

func TestAnalyzer_ComplexMethods(t *testing.T) {
	entityType := reflect.TypeOf(TestUser{})
	analyzer, err := NewAnalyzer(entityType)
	if err != nil {
		t.Fatalf("Failed to create analyzer: %v", err)
	}

	t.Run("FindByAgeGreaterThanAndStatusIn", func(t *testing.T) {
		method, err := analyzer.AnalyzeMethod("FindByAgeGreaterThanAndStatusIn")
		if err != nil {
			t.Fatalf("Failed to analyze: %v", err)
		}

		if len(method.Fields) != 2 {
			t.Fatalf("Expected 2 fields, got %d", len(method.Fields))
		}

		if method.Fields[0].FieldName != "Age" {
			t.Errorf("Expected first field 'Age', got '%s'", method.Fields[0].FieldName)
		}
		if method.Fields[0].Operator != OpGreaterThan {
			t.Errorf("Expected operator OpGreaterThan, got %v", method.Fields[0].Operator)
		}

		if method.Fields[1].FieldName != "Status" {
			t.Errorf("Expected second field 'Status', got '%s'", method.Fields[1].FieldName)
		}
		if method.Fields[1].Operator != OpIn {
			t.Errorf("Expected operator OpIn, got %v", method.Fields[1].Operator)
		}
		if method.Fields[1].AndOr != "AND" {
			t.Errorf("Expected AND, got '%s'", method.Fields[1].AndOr)
		}
	})

	t.Run("FindByStatusOrderByCreatedAtDesc", func(t *testing.T) {
		method, err := analyzer.AnalyzeMethod("FindByStatusOrderByCreatedAtDesc")
		if err != nil {
			t.Fatalf("Failed to analyze: %v", err)
		}

		if len(method.SortFields) != 1 {
			t.Fatalf("Expected 1 sort field, got %d", len(method.SortFields))
		}

		if method.SortFields[0].FieldName != "CreatedAt" {
			t.Errorf("Expected sort field 'CreatedAt', got '%s'", method.SortFields[0].FieldName)
		}
		if method.SortFields[0].Direction != "DESC" {
			t.Errorf("Expected direction 'DESC', got '%s'", method.SortFields[0].Direction)
		}
	})

	t.Run("FindFirstByStatus", func(t *testing.T) {
		method, err := analyzer.AnalyzeMethod("FindFirstByStatus")
		if err != nil {
			t.Fatalf("Failed to analyze: %v", err)
		}

		if method.Limit != 1 {
			t.Errorf("Expected limit 1, got %d", method.Limit)
		}
	})
}

func TestAnalyzer_ToSQL(t *testing.T) {
	entityType := reflect.TypeOf(TestUser{})
	analyzer, err := NewAnalyzer(entityType)
	if err != nil {
		t.Fatalf("Failed to create analyzer: %v", err)
	}

	fieldToColumn := func(fieldName string) string {
		// Simple snake_case conversion for test
		return fieldName
	}

	t.Run("FindByEmail SQL", func(t *testing.T) {
		method, err := analyzer.AnalyzeMethod("FindByEmail")
		if err != nil {
			t.Fatalf("Failed to analyze: %v", err)
		}

		sql := method.ToSQL("users", fieldToColumn)
		expected := "SELECT * FROM users WHERE Email = $1"
		if sql != expected {
			t.Errorf("Expected SQL '%s', got '%s'", expected, sql)
		}
	})

	t.Run("FindByAgeGreaterThanAndStatus SQL", func(t *testing.T) {
		method, err := analyzer.AnalyzeMethod("FindByAgeGreaterThanAndStatus")
		if err != nil {
			t.Fatalf("Failed to analyze: %v", err)
		}

		sql := method.ToSQL("users", fieldToColumn)
		if !contains(sql, "Age > $1") {
			t.Errorf("SQL should contain 'Age > $1', got: %s", sql)
		}
		if !contains(sql, "Status = $2") {
			t.Errorf("SQL should contain 'Status = $2', got: %s", sql)
		}
		if !contains(sql, "AND") {
			t.Errorf("SQL should contain 'AND', got: %s", sql)
		}
	})
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && 
		(s == substr || 
		 (len(s) > len(substr) && 
		  (s[:len(substr)] == substr || 
		   s[len(s)-len(substr):] == substr || 
		   contains(s[1:], substr))))
}

