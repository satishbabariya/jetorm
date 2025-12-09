package jet

import (
	"github.com/go-jet/jet/v2/postgres"
)

// Helper functions for Jet SQL integration
// These functions provide convenient wrappers around Jet SQL expressions

// Equal creates an equality condition
func Equal(column postgres.Column, value interface{}) postgres.BoolExpression {
	return column.EQ(postgres.RawValue(value))
}

// NotEqual creates a not-equal condition
func NotEqual(column postgres.Column, value interface{}) postgres.BoolExpression {
	return column.NOT_EQ(postgres.RawValue(value))
}

// GreaterThan creates a greater-than condition
func GreaterThan(column postgres.Column, value interface{}) postgres.BoolExpression {
	return column.GT(postgres.RawValue(value))
}

// GreaterThanOrEqual creates a greater-than-or-equal condition
func GreaterThanOrEqual(column postgres.Column, value interface{}) postgres.BoolExpression {
	return column.GT_EQ(postgres.RawValue(value))
}

// LessThan creates a less-than condition
func LessThan(column postgres.Column, value interface{}) postgres.BoolExpression {
	return column.LT(postgres.RawValue(value))
}

// LessThanOrEqual creates a less-than-or-equal condition
func LessThanOrEqual(column postgres.Column, value interface{}) postgres.BoolExpression {
	return column.LT_EQ(postgres.RawValue(value))
}

// Like creates a LIKE condition
func Like(column postgres.Column, pattern string) postgres.BoolExpression {
	return column.LIKE(postgres.String(pattern))
}

// ILike creates an ILIKE condition (case-insensitive)
func ILike(column postgres.Column, pattern string) postgres.BoolExpression {
	return column.ILIKE(postgres.String(pattern))
}

// In creates an IN condition
func In(column postgres.Column, values ...interface{}) postgres.BoolExpression {
	jetValues := make([]postgres.Expression, len(values))
	for i, v := range values {
		jetValues[i] = postgres.RawValue(v)
	}
	return column.IN(jetValues...)
}

// NotIn creates a NOT IN condition
func NotIn(column postgres.Column, values ...interface{}) postgres.BoolExpression {
	jetValues := make([]postgres.Expression, len(values))
	for i, v := range values {
		jetValues[i] = postgres.RawValue(v)
	}
	return column.NOT_IN(jetValues...)
}

// IsNull creates an IS NULL condition
func IsNull(column postgres.Column) postgres.BoolExpression {
	return column.IS_NULL()
}

// IsNotNull creates an IS NOT NULL condition
func IsNotNull(column postgres.Column) postgres.BoolExpression {
	return column.IS_NOT_NULL()
}

// Between creates a BETWEEN condition
func Between(column postgres.Column, min, max interface{}) postgres.BoolExpression {
	return column.BETWEEN(postgres.RawValue(min), postgres.RawValue(max))
}

// And combines multiple conditions with AND
func And(conditions ...postgres.BoolExpression) postgres.BoolExpression {
	if len(conditions) == 0 {
		return postgres.Bool(true)
	}
	if len(conditions) == 1 {
		return conditions[0]
	}
	result := conditions[0]
	for i := 1; i < len(conditions); i++ {
		result = result.AND(conditions[i])
	}
	return result
}

// Or combines multiple conditions with OR
func Or(conditions ...postgres.BoolExpression) postgres.BoolExpression {
	if len(conditions) == 0 {
		return postgres.Bool(false)
	}
	if len(conditions) == 1 {
		return conditions[0]
	}
	result := conditions[0]
	for i := 1; i < len(conditions); i++ {
		result = result.OR(conditions[i])
	}
	return result
}

// Not negates a condition
func Not(condition postgres.BoolExpression) postgres.BoolExpression {
	return postgres.NOT(condition)
}

// OrderBy creates an ORDER BY clause
// Returns the column with ASC or DESC applied
func OrderBy(column postgres.Column, ascending bool) postgres.OrderByClause {
	if ascending {
		return column.ASC()
	}
	return column.DESC()
}

// Limit creates a LIMIT clause value
func Limit(count int) int64 {
	return int64(count)
}

// Offset creates an OFFSET clause value
func Offset(count int) int64 {
	return int64(count)
}

// Join creates an INNER JOIN clause
// Returns a join that can be used in FROM clause
func Join(leftTable, rightTable postgres.Table, condition postgres.BoolExpression) postgres.Table {
	return leftTable.INNER_JOIN(rightTable, condition)
}

// LeftJoin creates a LEFT JOIN clause
func LeftJoin(leftTable, rightTable postgres.Table, condition postgres.BoolExpression) postgres.Table {
	return leftTable.LEFT_JOIN(rightTable, condition)
}

// RightJoin creates a RIGHT JOIN clause
func RightJoin(leftTable, rightTable postgres.Table, condition postgres.BoolExpression) postgres.Table {
	return leftTable.RIGHT_JOIN(rightTable, condition)
}

// FullJoin creates a FULL OUTER JOIN clause
func FullJoin(leftTable, rightTable postgres.Table, condition postgres.BoolExpression) postgres.Table {
	return leftTable.FULL_JOIN(rightTable, condition)
}

// GroupBy creates a GROUP BY clause
// Returns columns that can be used in GROUP BY
func GroupBy(columns ...postgres.Column) []postgres.Column {
	return columns
}

// Having creates a HAVING clause condition
func Having(condition postgres.BoolExpression) postgres.BoolExpression {
	return condition
}

// Aggregate functions

// Count creates a COUNT expression
func Count(column postgres.Column) postgres.IntegerExpression {
	return postgres.COUNT(column)
}

// CountStar creates a COUNT(*) expression
// Note: Uses COUNT(1) as Star may not be available in all contexts
func CountStar() postgres.IntegerExpression {
	return postgres.COUNT(postgres.Int(1))
}

// Sum creates a SUM expression
func Sum(column postgres.Column) postgres.NumericExpression {
	return postgres.SUM(column)
}

// Avg creates an AVG expression
func Avg(column postgres.Column) postgres.NumericExpression {
	return postgres.AVG(column)
}

// Min creates a MIN expression
func Min(column postgres.Column) postgres.Expression {
	return postgres.MIN(column)
}

// Max creates a MAX expression
func Max(column postgres.Column) postgres.Expression {
	return postgres.MAX(column)
}

// Distinct creates a DISTINCT expression
func Distinct(column postgres.Column) postgres.Expression {
	return postgres.DISTINCT(column)
}

// Window functions

// RowNumber creates a ROW_NUMBER() window function
func RowNumber() postgres.IntegerExpression {
	return postgres.ROW_NUMBER()
}

// Rank creates a RANK() window function
func Rank() postgres.IntegerExpression {
	return postgres.RANK()
}

// DenseRank creates a DENSE_RANK() window function
func DenseRank() postgres.IntegerExpression {
	return postgres.DENSE_RANK()
}

// Subquery helpers

// Exists creates an EXISTS subquery
func Exists(stmt postgres.SelectStatement) postgres.BoolExpression {
	return postgres.EXISTS(stmt)
}

// NotExists creates a NOT EXISTS subquery
func NotExists(stmt postgres.SelectStatement) postgres.BoolExpression {
	return postgres.NOT_EXISTS(stmt)
}

// InSubquery creates an IN subquery condition
func InSubquery(column postgres.Column, stmt postgres.SelectStatement) postgres.BoolExpression {
	return column.IN(stmt)
}

// NotInSubquery creates a NOT IN subquery condition
func NotInSubquery(column postgres.Column, stmt postgres.SelectStatement) postgres.BoolExpression {
	return column.NOT_IN(stmt)
}
