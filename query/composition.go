package query

import (
	"fmt"
	"strings"

	"github.com/satishbabariya/jetorm/core"
)

// ComposableQuery represents a query that can be composed with specifications
type ComposableQuery[T any] struct {
	builder     *QueryBuilder
	spec        core.Specification[T]
	tableName   string
	entityType  string
}

// NewComposableQuery creates a new composable query
func NewComposableQuery[T any](tableName string) *ComposableQuery[T] {
	return &ComposableQuery[T]{
		builder:   NewQueryBuilder(tableName),
		tableName: tableName,
	}
}

// WithSpecification sets a specification for the query
func (cq *ComposableQuery[T]) WithSpecification(spec core.Specification[T]) *ComposableQuery[T] {
	cq.spec = spec
	return cq
}

// Select sets the columns to select
func (cq *ComposableQuery[T]) Select(cols ...string) *ComposableQuery[T] {
	cq.builder.Select(cols...)
	return cq
}

// Where adds a WHERE clause
func (cq *ComposableQuery[T]) Where(condition string, args ...interface{}) *ComposableQuery[T] {
	cq.builder.Where(condition, args...)
	return cq
}

// WhereEqual adds an equality WHERE clause
func (cq *ComposableQuery[T]) WhereEqual(column string, value interface{}) *ComposableQuery[T] {
	cq.builder.WhereEqual(column, value)
	return cq
}

// WhereSpecification adds a WHERE clause from a specification
func (cq *ComposableQuery[T]) WhereSpecification(spec core.Specification[T]) *ComposableQuery[T] {
	if spec != nil {
		whereClause, args := spec.ToSQL()
		if whereClause != "" {
			cq.builder.Where(whereClause, args...)
		}
	}
	return cq
}

// OrderBy adds an ORDER BY clause
func (cq *ComposableQuery[T]) OrderBy(column string, direction string) *ComposableQuery[T] {
	cq.builder.OrderBy(column, direction)
	return cq
}

// Limit sets the LIMIT clause
func (cq *ComposableQuery[T]) Limit(limit int) *ComposableQuery[T] {
	cq.builder.Limit(limit)
	return cq
}

// Offset sets the OFFSET clause
func (cq *ComposableQuery[T]) Offset(offset int) *ComposableQuery[T] {
	cq.builder.Offset(offset)
	return cq
}

// GroupBy adds a GROUP BY clause
func (cq *ComposableQuery[T]) GroupBy(columns ...string) *ComposableQuery[T] {
	cq.builder.GroupBy(columns...)
	return cq
}

// Having adds a HAVING clause
func (cq *ComposableQuery[T]) Having(condition string, args ...interface{}) *ComposableQuery[T] {
	cq.builder.Having(condition, args...)
	return cq
}

// Build builds the final SQL query
func (cq *ComposableQuery[T]) Build() (string, []interface{}) {
	// Apply specification if set
	if cq.spec != nil {
		whereClause, args := cq.spec.ToSQL()
		if whereClause != "" {
			cq.builder.Where(whereClause, args...)
		}
	}
	return cq.builder.Build()
}

// BuildCount builds a COUNT query
func (cq *ComposableQuery[T]) BuildCount() (string, []interface{}) {
	// Apply specification if set
	if cq.spec != nil {
		whereClause, args := cq.spec.ToSQL()
		if whereClause != "" {
			cq.builder.Where(whereClause, args...)
		}
	}
	return cq.builder.BuildCount()
}

// Join represents a JOIN clause
type Join struct {
	Type      string // "INNER", "LEFT", "RIGHT", "FULL"
	Table     string
	Condition string
	Args      []interface{}
}

// JoinQuery represents a query with joins
type JoinQuery[T any] struct {
	*ComposableQuery[T]
	joins []Join
}

// NewJoinQuery creates a new query with join support
func NewJoinQuery[T any](tableName string) *JoinQuery[T] {
	return &JoinQuery[T]{
		ComposableQuery: NewComposableQuery[T](tableName),
		joins:           make([]Join, 0),
	}
}

// InnerJoin adds an INNER JOIN
func (jq *JoinQuery[T]) InnerJoin(table, condition string, args ...interface{}) *JoinQuery[T] {
	jq.joins = append(jq.joins, Join{
		Type:      "INNER",
		Table:     table,
		Condition: condition,
		Args:      args,
	})
	return jq
}

// WhereEqual adds an equality WHERE clause
func (jq *JoinQuery[T]) WhereEqual(column string, value interface{}) *JoinQuery[T] {
	jq.ComposableQuery.WhereEqual(column, value)
	return jq
}

// LeftJoin adds a LEFT JOIN
func (jq *JoinQuery[T]) LeftJoin(table, condition string, args ...interface{}) *JoinQuery[T] {
	jq.joins = append(jq.joins, Join{
		Type:      "LEFT",
		Table:     table,
		Condition: condition,
		Args:      args,
	})
	return jq
}

// RightJoin adds a RIGHT JOIN
func (jq *JoinQuery[T]) RightJoin(table, condition string, args ...interface{}) *JoinQuery[T] {
	jq.joins = append(jq.joins, Join{
		Type:      "RIGHT",
		Table:     table,
		Condition: condition,
		Args:      args,
	})
	return jq
}

// FullJoin adds a FULL OUTER JOIN
func (jq *JoinQuery[T]) FullJoin(table, condition string, args ...interface{}) *JoinQuery[T] {
	jq.joins = append(jq.joins, Join{
		Type:      "FULL",
		Table:     table,
		Condition: condition,
		Args:      args,
	})
	return jq
}

// Build builds the query with joins
func (jq *JoinQuery[T]) Build() (string, []interface{}) {
	query, args := jq.ComposableQuery.Build()
	
	// Insert JOIN clauses after FROM
	if len(jq.joins) > 0 {
		fromIndex := strings.Index(query, "FROM")
		if fromIndex > 0 {
			beforeFrom := query[:fromIndex+4]
			afterFrom := query[fromIndex+4:]
			
			var joinClauses []string
			joinArgs := make([]interface{}, 0)
			
			for _, join := range jq.joins {
				joinType := join.Type
				if joinType == "FULL" {
					joinType = "FULL OUTER"
				}
				joinClauses = append(joinClauses, fmt.Sprintf("%s JOIN %s ON %s", joinType, join.Table, join.Condition))
				joinArgs = append(joinArgs, join.Args...)
			}
			
			query = beforeFrom + " " + strings.Join(joinClauses, " ") + " " + afterFrom
			args = append(joinArgs, args...)
		}
	}
	
	return query, args
}

// Subquery represents a subquery
type Subquery struct {
	Query string
	Args  []interface{}
	Alias string
}

// SubqueryQuery represents a query with subqueries
type SubqueryQuery[T any] struct {
	*ComposableQuery[T]
	subqueries []Subquery
}

// NewSubqueryQuery creates a new query with subquery support
func NewSubqueryQuery[T any](tableName string) *SubqueryQuery[T] {
	return &SubqueryQuery[T]{
		ComposableQuery: NewComposableQuery[T](tableName),
		subqueries:      make([]Subquery, 0),
	}
}

// WithSubquery adds a subquery to SELECT clause
func (sq *SubqueryQuery[T]) WithSubquery(query string, args []interface{}, alias string) *SubqueryQuery[T] {
	sq.subqueries = append(sq.subqueries, Subquery{
		Query: query,
		Args:  args,
		Alias: alias,
	})
	return sq
}

// Build builds the query with subqueries
func (sq *SubqueryQuery[T]) Build() (string, []interface{}) {
	query, args := sq.ComposableQuery.Build()
	
	// Add subqueries to SELECT clause
	if len(sq.subqueries) > 0 {
		selectIndex := strings.Index(query, "SELECT")
		if selectIndex >= 0 {
			afterSelect := query[selectIndex+6:]
			fromIndex := strings.Index(afterSelect, "FROM")
			
			if fromIndex > 0 {
				beforeFrom := query[:selectIndex+6] + afterSelect[:fromIndex]
				afterFrom := afterSelect[fromIndex:]
				
				var subqueryClauses []string
				subqueryArgs := make([]interface{}, 0)
				
				for _, subq := range sq.subqueries {
					subqueryClauses = append(subqueryClauses, fmt.Sprintf("(%s) AS %s", subq.Query, subq.Alias))
					subqueryArgs = append(subqueryArgs, subq.Args...)
				}
				
				existingCols := strings.TrimSpace(beforeFrom[selectIndex+6:])
				if existingCols == "" || existingCols == "*" {
					query = query[:selectIndex+6] + strings.Join(subqueryClauses, ", ") + " " + afterFrom
				} else {
					query = beforeFrom + ", " + strings.Join(subqueryClauses, ", ") + " " + afterFrom
				}
				
				args = append(subqueryArgs, args...)
			}
		}
	}
	
	return query, args
}

// DynamicQuery allows building queries dynamically based on conditions
type DynamicQuery[T any] struct {
	*ComposableQuery[T]
	conditions []func(*ComposableQuery[T]) *ComposableQuery[T]
}

// NewDynamicQuery creates a new dynamic query builder
func NewDynamicQuery[T any](tableName string) *DynamicQuery[T] {
	return &DynamicQuery[T]{
		ComposableQuery: NewComposableQuery[T](tableName),
		conditions:      make([]func(*ComposableQuery[T]) *ComposableQuery[T], 0),
	}
}

// When adds a conditional clause
func (dq *DynamicQuery[T]) When(condition bool, fn func(*ComposableQuery[T]) *ComposableQuery[T]) *DynamicQuery[T] {
	if condition {
		dq.conditions = append(dq.conditions, fn)
	}
	return dq
}

// Build builds the query with all applied conditions
func (dq *DynamicQuery[T]) Build() (string, []interface{}) {
	// Apply all conditions
	for _, condition := range dq.conditions {
		dq.ComposableQuery = condition(dq.ComposableQuery)
	}
	return dq.ComposableQuery.Build()
}

