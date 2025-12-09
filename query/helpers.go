package query

import (
	"fmt"
	"strings"
	"time"
)

// ConditionBuilder helps build WHERE conditions
type ConditionBuilder struct {
	conditions []string
	args       []interface{}
}

// NewConditionBuilder creates a new condition builder
func NewConditionBuilder() *ConditionBuilder {
	return &ConditionBuilder{
		conditions: make([]string, 0),
		args:       make([]interface{}, 0),
	}
}

// Equal adds an equality condition
func (cb *ConditionBuilder) Equal(column string, value interface{}) *ConditionBuilder {
	argIndex := len(cb.args) + 1
	cb.conditions = append(cb.conditions, fmt.Sprintf("%s = $%d", column, argIndex))
	cb.args = append(cb.args, value)
	return cb
}

// NotEqual adds a not-equal condition
func (cb *ConditionBuilder) NotEqual(column string, value interface{}) *ConditionBuilder {
	argIndex := len(cb.args) + 1
	cb.conditions = append(cb.conditions, fmt.Sprintf("%s != $%d", column, argIndex))
	cb.args = append(cb.args, value)
	return cb
}

// GreaterThan adds a greater-than condition
func (cb *ConditionBuilder) GreaterThan(column string, value interface{}) *ConditionBuilder {
	argIndex := len(cb.args) + 1
	cb.conditions = append(cb.conditions, fmt.Sprintf("%s > $%d", column, argIndex))
	cb.args = append(cb.args, value)
	return cb
}

// GreaterThanEqual adds a greater-than-or-equal condition
func (cb *ConditionBuilder) GreaterThanEqual(column string, value interface{}) *ConditionBuilder {
	argIndex := len(cb.args) + 1
	cb.conditions = append(cb.conditions, fmt.Sprintf("%s >= $%d", column, argIndex))
	cb.args = append(cb.args, value)
	return cb
}

// LessThan adds a less-than condition
func (cb *ConditionBuilder) LessThan(column string, value interface{}) *ConditionBuilder {
	argIndex := len(cb.args) + 1
	cb.conditions = append(cb.conditions, fmt.Sprintf("%s < $%d", column, argIndex))
	cb.args = append(cb.args, value)
	return cb
}

// LessThanEqual adds a less-than-or-equal condition
func (cb *ConditionBuilder) LessThanEqual(column string, value interface{}) *ConditionBuilder {
	argIndex := len(cb.args) + 1
	cb.conditions = append(cb.conditions, fmt.Sprintf("%s <= $%d", column, argIndex))
	cb.args = append(cb.args, value)
	return cb
}

// Like adds a LIKE condition
func (cb *ConditionBuilder) Like(column string, pattern string) *ConditionBuilder {
	argIndex := len(cb.args) + 1
	cb.conditions = append(cb.conditions, fmt.Sprintf("%s LIKE $%d", column, argIndex))
	cb.args = append(cb.args, pattern)
	return cb
}

// ILike adds a case-insensitive LIKE condition (PostgreSQL)
func (cb *ConditionBuilder) ILike(column string, pattern string) *ConditionBuilder {
	argIndex := len(cb.args) + 1
	cb.conditions = append(cb.conditions, fmt.Sprintf("%s ILIKE $%d", column, argIndex))
	cb.args = append(cb.args, pattern)
	return cb
}

// In adds an IN condition
func (cb *ConditionBuilder) In(column string, values []interface{}) *ConditionBuilder {
	if len(values) == 0 {
		return cb
	}
	placeholders := make([]string, len(values))
	for i := range values {
		argIndex := len(cb.args) + i + 1
		placeholders[i] = fmt.Sprintf("$%d", argIndex)
	}
	cb.conditions = append(cb.conditions, fmt.Sprintf("%s IN (%s)", column, strings.Join(placeholders, ", ")))
	cb.args = append(cb.args, values...)
	return cb
}

// NotIn adds a NOT IN condition
func (cb *ConditionBuilder) NotIn(column string, values []interface{}) *ConditionBuilder {
	if len(values) == 0 {
		return cb
	}
	placeholders := make([]string, len(values))
	for i := range values {
		argIndex := len(cb.args) + i + 1
		placeholders[i] = fmt.Sprintf("$%d", argIndex)
	}
	cb.conditions = append(cb.conditions, fmt.Sprintf("%s NOT IN (%s)", column, strings.Join(placeholders, ", ")))
	cb.args = append(cb.args, values...)
	return cb
}

// Between adds a BETWEEN condition
func (cb *ConditionBuilder) Between(column string, min, max interface{}) *ConditionBuilder {
	argIndex := len(cb.args) + 1
	cb.conditions = append(cb.conditions, fmt.Sprintf("%s BETWEEN $%d AND $%d", column, argIndex, argIndex+1))
	cb.args = append(cb.args, min, max)
	return cb
}

// IsNull adds an IS NULL condition
func (cb *ConditionBuilder) IsNull(column string) *ConditionBuilder {
	cb.conditions = append(cb.conditions, fmt.Sprintf("%s IS NULL", column))
	return cb
}

// IsNotNull adds an IS NOT NULL condition
func (cb *ConditionBuilder) IsNotNull(column string) *ConditionBuilder {
	cb.conditions = append(cb.conditions, fmt.Sprintf("%s IS NOT NULL", column))
	return cb
}

// Exists adds an EXISTS condition
func (cb *ConditionBuilder) Exists(subquery string, args ...interface{}) *ConditionBuilder {
	cb.conditions = append(cb.conditions, fmt.Sprintf("EXISTS (%s)", subquery))
	cb.args = append(cb.args, args...)
	return cb
}

// NotExists adds a NOT EXISTS condition
func (cb *ConditionBuilder) NotExists(subquery string, args ...interface{}) *ConditionBuilder {
	cb.conditions = append(cb.conditions, fmt.Sprintf("NOT EXISTS (%s)", subquery))
	cb.args = append(cb.args, args...)
	return cb
}

// And combines conditions with AND
func (cb *ConditionBuilder) And(other *ConditionBuilder) *ConditionBuilder {
	cb.conditions = append(cb.conditions, other.conditions...)
	cb.args = append(cb.args, other.args...)
	return cb
}

// Or combines conditions with OR (wraps in parentheses)
func (cb *ConditionBuilder) Or(other *ConditionBuilder) *ConditionBuilder {
	if len(cb.conditions) > 0 && len(other.conditions) > 0 {
		left := "(" + strings.Join(cb.conditions, " AND ") + ")"
		right := "(" + strings.Join(other.conditions, " AND ") + ")"
		cb.conditions = []string{left + " OR " + right}
		cb.args = append(cb.args, other.args...)
	}
	return cb
}

// Build builds the WHERE clause
func (cb *ConditionBuilder) Build() (string, []interface{}) {
	if len(cb.conditions) == 0 {
		return "", nil
	}
	return strings.Join(cb.conditions, " AND "), cb.args
}

// DateRange creates a condition for date range queries
func DateRange(column string, start, end time.Time) *ConditionBuilder {
	cb := NewConditionBuilder()
	return cb.GreaterThanEqual(column, start).And(cb.LessThanEqual(column, end))
}

// TextSearch creates a condition for full-text search (PostgreSQL)
func TextSearch(column string, searchTerm string) *ConditionBuilder {
	cb := NewConditionBuilder()
	// Use PostgreSQL's to_tsvector for full-text search
	cb.conditions = append(cb.conditions, fmt.Sprintf("to_tsvector('english', %s) @@ plainto_tsquery('english', $%d)", column, len(cb.args)+1))
	cb.args = append(cb.args, searchTerm)
	return cb
}

// ArrayContains creates a condition for array containment (PostgreSQL)
func ArrayContains(column string, value interface{}) *ConditionBuilder {
	cb := NewConditionBuilder()
	argIndex := len(cb.args) + 1
	cb.conditions = append(cb.conditions, fmt.Sprintf("$%d = ANY(%s)", argIndex, column))
	cb.args = append(cb.args, value)
	return cb
}

// ArrayOverlaps creates a condition for array overlap (PostgreSQL)
func ArrayOverlaps(column string, values []interface{}) *ConditionBuilder {
	cb := NewConditionBuilder()
	if len(values) == 0 {
		return cb
	}
	placeholders := make([]string, len(values))
	for i := range values {
		argIndex := len(cb.args) + i + 1
		placeholders[i] = fmt.Sprintf("$%d", argIndex)
	}
	arrayLiteral := "ARRAY[" + strings.Join(placeholders, ", ") + "]"
	cb.conditions = append(cb.conditions, fmt.Sprintf("%s && %s", column, arrayLiteral))
	cb.args = append(cb.args, values...)
	return cb
}

