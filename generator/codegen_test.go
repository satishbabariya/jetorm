package generator

import (
	"reflect"
	"strings"
	"testing"
)

func TestCodeGenerator_GenerateMethod(t *testing.T) {
	entityType := reflect.TypeOf(TestUser{})
	gen, err := NewCodeGenerator(entityType)
	if err != nil {
		t.Fatalf("Failed to create code generator: %v", err)
	}

	t.Run("generate FindByEmail method", func(t *testing.T) {
		analyzer, _ := NewAnalyzer(entityType)
		method, err := analyzer.AnalyzeMethod("FindByEmail")
		if err != nil {
			t.Fatalf("Failed to analyze method: %v", err)
		}

		code, err := gen.GenerateMethod(method, "User", "int64")
		if err != nil {
			t.Fatalf("Failed to generate method: %v", err)
		}

		if !strings.Contains(code, "FindByEmail") {
			t.Error("Generated code should contain method name")
		}
		if !strings.Contains(code, "email string") {
			t.Error("Generated code should contain parameter")
		}
		if !strings.Contains(code, "*User") || !strings.Contains(code, "error") {
			t.Error("Generated code should contain return type")
		}
		if !strings.Contains(code, "SELECT * FROM") {
			t.Error("Generated code should contain SQL query")
		}
	})

	t.Run("generate FindByAgeGreaterThan method", func(t *testing.T) {
		analyzer, _ := NewAnalyzer(entityType)
		method, err := analyzer.AnalyzeMethod("FindByAgeGreaterThan")
		if err != nil {
			t.Fatalf("Failed to analyze method: %v", err)
		}

		code, err := gen.GenerateMethod(method, "User", "int64")
		if err != nil {
			t.Fatalf("Failed to generate method: %v", err)
		}

		if !strings.Contains(code, "age > $1") {
			t.Error("Generated code should contain correct SQL condition")
		}
	})

	t.Run("generate CountByStatus method", func(t *testing.T) {
		analyzer, _ := NewAnalyzer(entityType)
		method, err := analyzer.AnalyzeMethod("CountByStatus")
		if err != nil {
			t.Fatalf("Failed to analyze method: %v", err)
		}

		code, err := gen.GenerateMethod(method, "User", "int64")
		if err != nil {
			t.Fatalf("Failed to generate method: %v", err)
		}

		if !strings.Contains(code, "SELECT COUNT(*)") {
			t.Error("Generated code should contain COUNT query")
		}
		if !strings.Contains(code, "(int64, error)") {
			t.Error("Generated code should have correct return type for Count")
		}
	})
}

