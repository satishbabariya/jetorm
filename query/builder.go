package query

import (
	"context"
	"fmt"
	"strings"
)

// QueryBuilder builds SQL queries dynamically
type QueryBuilder struct {
	tableName string
	selectCols []string
	whereClauses []string
	whereArgs []interface{}
	orderBy []string
	limitVal *int
	offsetVal *int
	groupBy []string
	havingClauses []string
	havingArgs []interface{}
}

// NewQueryBuilder creates a new query builder
func NewQueryBuilder(tableName string) *QueryBuilder {
	return &QueryBuilder{
		tableName:     tableName,
		selectCols:    []string{"*"},
		whereClauses:  make([]string, 0),
		whereArgs:     make([]interface{}, 0),
		orderBy:       make([]string, 0),
		groupBy:       make([]string, 0),
		havingClauses: make([]string, 0),
		havingArgs:    make([]interface{}, 0),
	}
}

// Select sets the columns to select
func (qb *QueryBuilder) Select(cols ...string) *QueryBuilder {
	qb.selectCols = cols
	return qb
}

// Where adds a WHERE clause
func (qb *QueryBuilder) Where(condition string, args ...interface{}) *QueryBuilder {
	qb.whereClauses = append(qb.whereClauses, condition)
	qb.whereArgs = append(qb.whereArgs, args...)
	return qb
}

// WhereEqual adds an equality WHERE clause
func (qb *QueryBuilder) WhereEqual(column string, value interface{}) *QueryBuilder {
	argIndex := len(qb.whereArgs) + 1
	qb.whereClauses = append(qb.whereClauses, fmt.Sprintf("%s = $%d", column, argIndex))
	qb.whereArgs = append(qb.whereArgs, value)
	return qb
}

// WhereIn adds an IN clause
func (qb *QueryBuilder) WhereIn(column string, values []interface{}) *QueryBuilder {
	if len(values) == 0 {
		return qb
	}
	
	placeholders := make([]string, len(values))
	for i := range values {
		argIndex := len(qb.whereArgs) + i + 1
		placeholders[i] = fmt.Sprintf("$%d", argIndex)
	}
	
	qb.whereClauses = append(qb.whereClauses, fmt.Sprintf("%s IN (%s)", column, strings.Join(placeholders, ", ")))
	qb.whereArgs = append(qb.whereArgs, values...)
	return qb
}

// OrderBy adds an ORDER BY clause
func (qb *QueryBuilder) OrderBy(column string, direction string) *QueryBuilder {
	qb.orderBy = append(qb.orderBy, fmt.Sprintf("%s %s", column, direction))
	return qb
}

// Limit sets the LIMIT clause
func (qb *QueryBuilder) Limit(limit int) *QueryBuilder {
	qb.limitVal = &limit
	return qb
}

// Offset sets the OFFSET clause
func (qb *QueryBuilder) Offset(offset int) *QueryBuilder {
	qb.offsetVal = &offset
	return qb
}

// GroupBy adds a GROUP BY clause
func (qb *QueryBuilder) GroupBy(columns ...string) *QueryBuilder {
	qb.groupBy = append(qb.groupBy, columns...)
	return qb
}

// Having adds a HAVING clause
func (qb *QueryBuilder) Having(condition string, args ...interface{}) *QueryBuilder {
	qb.havingClauses = append(qb.havingClauses, condition)
	qb.havingArgs = append(qb.havingArgs, args...)
	return qb
}

// Build builds the SQL query string
func (qb *QueryBuilder) Build() (string, []interface{}) {
	var parts []string
	
	// SELECT
	parts = append(parts, "SELECT", strings.Join(qb.selectCols, ", "))
	
	// FROM
	parts = append(parts, "FROM", qb.tableName)
	
	// WHERE
	if len(qb.whereClauses) > 0 {
		parts = append(parts, "WHERE", strings.Join(qb.whereClauses, " AND "))
	}
	
	// GROUP BY
	if len(qb.groupBy) > 0 {
		parts = append(parts, "GROUP BY", strings.Join(qb.groupBy, ", "))
	}
	
	// HAVING
	if len(qb.havingClauses) > 0 {
		parts = append(parts, "HAVING", strings.Join(qb.havingClauses, " AND "))
	}
	
	// ORDER BY
	if len(qb.orderBy) > 0 {
		parts = append(parts, "ORDER BY", strings.Join(qb.orderBy, ", "))
	}
	
	// LIMIT
	if qb.limitVal != nil {
		parts = append(parts, fmt.Sprintf("LIMIT %d", *qb.limitVal))
	}
	
	// OFFSET
	if qb.offsetVal != nil {
		parts = append(parts, fmt.Sprintf("OFFSET %d", *qb.offsetVal))
	}
	
	query := strings.Join(parts, " ")
	args := append(qb.whereArgs, qb.havingArgs...)
	
	return query, args
}

// BuildCount builds a COUNT query
func (qb *QueryBuilder) BuildCount() (string, []interface{}) {
	var parts []string
	
	parts = append(parts, "SELECT COUNT(*)")
	parts = append(parts, "FROM", qb.tableName)
	
	if len(qb.whereClauses) > 0 {
		parts = append(parts, "WHERE", strings.Join(qb.whereClauses, " AND "))
	}
	
	if len(qb.groupBy) > 0 {
		parts = append(parts, "GROUP BY", strings.Join(qb.groupBy, ", "))
	}
	
	if len(qb.havingClauses) > 0 {
		parts = append(parts, "HAVING", strings.Join(qb.havingClauses, " AND "))
	}
	
	query := strings.Join(parts, " ")
	args := append(qb.whereArgs, qb.havingArgs...)
	
	return query, args
}

// Executor executes queries
type Executor interface {
	Query(ctx context.Context, query string, args ...interface{}) (Rows, error)
	QueryRow(ctx context.Context, query string, args ...interface{}) Row
	Exec(ctx context.Context, query string, args ...interface{}) (Result, error)
}

// Rows represents query result rows
type Rows interface {
	Scan(dest ...interface{}) error
	Next() bool
	Close() error
	Err() error
}

// Row represents a single query result row
type Row interface {
	Scan(dest ...interface{}) error
}

// Result represents the result of an Exec operation
type Result interface {
	RowsAffected() int64
}

// Execute executes the query using the provided executor
func (qb *QueryBuilder) Execute(ctx context.Context, executor Executor) (Rows, error) {
	query, args := qb.Build()
	return executor.Query(ctx, query, args...)
}

// ExecuteRow executes the query and returns a single row
func (qb *QueryBuilder) ExecuteRow(ctx context.Context, executor Executor) Row {
	query, args := qb.Build()
	return executor.QueryRow(ctx, query, args...)
}

