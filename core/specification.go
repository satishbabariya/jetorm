package core

import (
	"fmt"
	"regexp"
	"strings"
)

// Specification represents a composable query criteria
type Specification[T any] interface {
	// ToSQL converts the specification to SQL WHERE clause and arguments
	ToSQL() (whereClause string, args []interface{})
	
	// And combines this specification with another using AND
	And(other Specification[T]) Specification[T]
	
	// Or combines this specification with another using OR
	Or(other Specification[T]) Specification[T]
	
	// Not negates this specification
	Not() Specification[T]
}

// baseSpecification is the base implementation of Specification
type baseSpecification[T any] struct {
	whereClause string
	args        []interface{}
	operator    string // "AND", "OR", "NOT"
	left        Specification[T]
	right       Specification[T]
}

var placeholderRegex = regexp.MustCompile(`\$(\d+)`)

// ToSQL converts the specification to SQL WHERE clause and arguments
func (s *baseSpecification[T]) ToSQL() (string, []interface{}) {
	if s.operator == "" {
		// Simple specification
		return s.whereClause, s.args
	}
	
	// Composite specification
	var parts []string
	var allArgs []interface{}
	
	var leftArgs []interface{}
	if s.left != nil {
		leftSQL, args := s.left.ToSQL()
		if leftSQL != "" {
			// Renumber placeholders in left SQL starting from 1
			leftSQL = renumberPlaceholders(leftSQL, 1)
			parts = append(parts, fmt.Sprintf("(%s)", leftSQL))
			leftArgs = args
			allArgs = append(allArgs, leftArgs...)
		}
	}
	
	if s.right != nil {
		rightSQL, rightArgs := s.right.ToSQL()
		if rightSQL != "" {
			// Renumber placeholders in right SQL starting after left args
			rightSQL = renumberPlaceholders(rightSQL, len(leftArgs)+1)
			parts = append(parts, fmt.Sprintf("(%s)", rightSQL))
			allArgs = append(allArgs, rightArgs...)
		}
	}
	
	if len(parts) == 0 {
		return "", nil
	}
	
	if s.operator == "NOT" {
		if len(parts) > 0 {
			return fmt.Sprintf("NOT %s", parts[0]), allArgs
		}
		return "", nil
	}
	
	return strings.Join(parts, fmt.Sprintf(" %s ", s.operator)), allArgs
}

// renumberPlaceholders renumbers SQL placeholders starting from startNum
// For example, if sql is "field = $1 AND other = $2" and startNum is 3,
// it becomes "field = $3 AND other = $4"
func renumberPlaceholders(sql string, startNum int) string {
	if startNum == 1 {
		// No renumbering needed
		return sql
	}
	
	offset := startNum - 1
	return placeholderRegex.ReplaceAllStringFunc(sql, func(match string) string {
		// Extract the number from $1, $2, etc.
		var num int
		fmt.Sscanf(match, "$%d", &num)
		// Renumber: if original was $1 and startNum is 3, result is $3
		// So: newNum = oldNum + (startNum - 1)
		return fmt.Sprintf("$%d", num+offset)
	})
}

// And combines this specification with another using AND
func (s *baseSpecification[T]) And(other Specification[T]) Specification[T] {
	return &baseSpecification[T]{
		operator: "AND",
		left:     s,
		right:    other,
	}
}

// Or combines this specification with another using OR
func (s *baseSpecification[T]) Or(other Specification[T]) Specification[T] {
	return &baseSpecification[T]{
		operator: "OR",
		left:     s,
		right:    other,
	}
}

// Not negates this specification
func (s *baseSpecification[T]) Not() Specification[T] {
	return &baseSpecification[T]{
		operator: "NOT",
		left:     s,
	}
}

// Where creates a specification from a SQL WHERE clause
func Where[T any](whereClause string, args ...interface{}) Specification[T] {
	return &baseSpecification[T]{
		whereClause: whereClause,
		args:        args,
	}
}

// And combines multiple specifications using AND
func And[T any](specs ...Specification[T]) Specification[T] {
	if len(specs) == 0 {
		return nil
	}
	if len(specs) == 1 {
		return specs[0]
	}
	
	result := specs[0]
	for i := 1; i < len(specs); i++ {
		result = result.And(specs[i])
	}
	return result
}

// Or combines multiple specifications using OR
func Or[T any](specs ...Specification[T]) Specification[T] {
	if len(specs) == 0 {
		return nil
	}
	if len(specs) == 1 {
		return specs[0]
	}
	
	result := specs[0]
	for i := 1; i < len(specs); i++ {
		result = result.Or(specs[i])
	}
	return result
}

// Not negates a specification
func Not[T any](spec Specification[T]) Specification[T] {
	if spec == nil {
		return nil
	}
	return spec.Not()
}

// Helper functions for common conditions

// Equal creates a specification for field = value
func Equal[T any](field string, value interface{}) Specification[T] {
	return Where[T](fmt.Sprintf("%s = $1", field), value)
}

// NotEqual creates a specification for field != value
func NotEqual[T any](field string, value interface{}) Specification[T] {
	return Where[T](fmt.Sprintf("%s != $1", field), value)
}

// GreaterThan creates a specification for field > value
func GreaterThan[T any](field string, value interface{}) Specification[T] {
	return Where[T](fmt.Sprintf("%s > $1", field), value)
}

// GreaterThanEqual creates a specification for field >= value
func GreaterThanEqual[T any](field string, value interface{}) Specification[T] {
	return Where[T](fmt.Sprintf("%s >= $1", field), value)
}

// LessThan creates a specification for field < value
func LessThan[T any](field string, value interface{}) Specification[T] {
	return Where[T](fmt.Sprintf("%s < $1", field), value)
}

// LessThanEqual creates a specification for field <= value
func LessThanEqual[T any](field string, value interface{}) Specification[T] {
	return Where[T](fmt.Sprintf("%s <= $1", field), value)
}

// Like creates a specification for field LIKE pattern
func Like[T any](field string, pattern string) Specification[T] {
	return Where[T](fmt.Sprintf("%s LIKE $1", field), pattern)
}

// In creates a specification for field IN (values...)
func In[T any](field string, values ...interface{}) Specification[T] {
	if len(values) == 0 {
		return Where[T]("1 = 0") // Always false
	}
	
	placeholders := make([]string, len(values))
	for i := range values {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
	}
	
	return Where[T](
		fmt.Sprintf("%s IN (%s)", field, strings.Join(placeholders, ", ")),
		values...,
	)
}

// NotIn creates a specification for field NOT IN (values...)
func NotIn[T any](field string, values ...interface{}) Specification[T] {
	if len(values) == 0 {
		return Where[T]("1 = 1") // Always true
	}
	
	placeholders := make([]string, len(values))
	for i := range values {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
	}
	
	return Where[T](
		fmt.Sprintf("%s NOT IN (%s)", field, strings.Join(placeholders, ", ")),
		values...,
	)
}

// IsNull creates a specification for field IS NULL
func IsNull[T any](field string) Specification[T] {
	return Where[T](fmt.Sprintf("%s IS NULL", field))
}

// IsNotNull creates a specification for field IS NOT NULL
func IsNotNull[T any](field string) Specification[T] {
	return Where[T](fmt.Sprintf("%s IS NOT NULL", field))
}

// Between creates a specification for field BETWEEN min AND max
func Between[T any](field string, min, max interface{}) Specification[T] {
	return Where[T](
		fmt.Sprintf("%s BETWEEN $1 AND $2", field),
		min, max,
	)
}

// Contains creates a specification for field LIKE '%value%'
func Contains[T any](field string, value string) Specification[T] {
	return Where[T](fmt.Sprintf("%s LIKE $1", field), "%"+value+"%")
}

// StartsWith creates a specification for field LIKE 'value%'
func StartsWith[T any](field string, value string) Specification[T] {
	return Where[T](fmt.Sprintf("%s LIKE $1", field), value+"%")
}

// EndsWith creates a specification for field LIKE '%value'
func EndsWith[T any](field string, value string) Specification[T] {
	return Where[T](fmt.Sprintf("%s LIKE $1", field), "%"+value)
}

