package query

import (
	"context"
	"fmt"
	"strings"
)

// AdvancedQueryBuilder provides advanced query building features
type AdvancedQueryBuilder struct {
	*QueryBuilder
	subqueries []*AdvancedSubquery
	unions     []*UnionQuery
	ctes       []*CTE
	window     *WindowFunction
}

// AdvancedSubquery represents a subquery in advanced builder
type AdvancedSubquery struct {
	Alias   string
	Builder *QueryBuilder
}

// UnionQuery represents a UNION query
type UnionQuery struct {
	Type   string // UNION, UNION ALL, INTERSECT, EXCEPT
	Builder *QueryBuilder
}

// CTE represents a Common Table Expression
type CTE struct {
	Name    string
	Builder *QueryBuilder
}

// WindowFunction represents a window function
type WindowFunction struct {
	Function string
	Over     string
}

// NewAdvancedQueryBuilder creates a new advanced query builder
func NewAdvancedQueryBuilder(tableName string) *AdvancedQueryBuilder {
	return &AdvancedQueryBuilder{
		QueryBuilder: NewQueryBuilder(tableName),
		subqueries:   make([]*AdvancedSubquery, 0),
		unions:      make([]*UnionQuery, 0),
		ctes:        make([]*CTE, 0),
	}
}

// WithCTE adds a Common Table Expression
func (aqb *AdvancedQueryBuilder) WithCTE(name string, builder *QueryBuilder) *AdvancedQueryBuilder {
	aqb.ctes = append(aqb.ctes, &CTE{
		Name:    name,
		Builder: builder,
	})
	return aqb
}

// Subquery adds a subquery
func (aqb *AdvancedQueryBuilder) Subquery(alias string, builder *QueryBuilder) *AdvancedQueryBuilder {
	aqb.subqueries = append(aqb.subqueries, &AdvancedSubquery{
		Alias:   alias,
		Builder: builder,
	})
	return aqb
}

// Union adds a UNION query
func (aqb *AdvancedQueryBuilder) Union(builder *QueryBuilder) *AdvancedQueryBuilder {
	aqb.unions = append(aqb.unions, &UnionQuery{
		Type:    "UNION",
		Builder: builder,
	})
	return aqb
}

// UnionAll adds a UNION ALL query
func (aqb *AdvancedQueryBuilder) UnionAll(builder *QueryBuilder) *AdvancedQueryBuilder {
	aqb.unions = append(aqb.unions, &UnionQuery{
		Type:    "UNION ALL",
		Builder: builder,
	})
	return aqb
}

// Window adds a window function
func (aqb *AdvancedQueryBuilder) Window(function, over string) *AdvancedQueryBuilder {
	aqb.window = &WindowFunction{
		Function: function,
		Over:     over,
	}
	return aqb
}

// BuildAdvanced builds the advanced query
func (aqb *AdvancedQueryBuilder) BuildAdvanced() (string, []interface{}) {
	var parts []string
	var args []interface{}

	// Build CTEs
	if len(aqb.ctes) > 0 {
		cteParts := make([]string, 0, len(aqb.ctes))
		for _, cte := range aqb.ctes {
			cteQuery, cteArgs := cte.Builder.Build()
			cteParts = append(cteParts, fmt.Sprintf("%s AS (%s)", cte.Name, cteQuery))
			args = append(args, cteArgs...)
		}
		parts = append(parts, "WITH "+strings.Join(cteParts, ", "))
	}

	// Build main query
	mainQuery, mainArgs := aqb.QueryBuilder.Build()
	parts = append(parts, mainQuery)
	args = append(args, mainArgs...)

	// Build unions
	for _, union := range aqb.unions {
		unionQuery, unionArgs := union.Builder.Build()
		parts = append(parts, union.Type, unionQuery)
		args = append(args, unionArgs...)
	}

	return strings.Join(parts, " "), args
}

// QueryComposer provides fluent query composition
type QueryComposer struct {
	queries []*QueryBuilder
}

// NewQueryComposer creates a new query composer
func NewQueryComposer() *QueryComposer {
	return &QueryComposer{
		queries: make([]*QueryBuilder, 0),
	}
}

// AddQuery adds a query to the composition
func (qc *QueryComposer) AddQuery(builder *QueryBuilder) *QueryComposer {
	qc.queries = append(qc.queries, builder)
	return qc
}

// Compose composes all queries
func (qc *QueryComposer) Compose(operator string) (string, []interface{}) {
	if len(qc.queries) == 0 {
		return "", nil
	}

	var parts []string
	var args []interface{}

	for i, query := range qc.queries {
		if i > 0 {
			parts = append(parts, operator)
		}
		queryStr, queryArgs := query.Build()
		parts = append(parts, "("+queryStr+")")
		args = append(args, queryArgs...)
	}

	return strings.Join(parts, " "), args
}

// ConditionalBuilder provides conditional query building
type ConditionalBuilder struct {
	conditions []*ConditionalClause
}

// ConditionalClause represents a conditional clause
type ConditionalClause struct {
	Condition bool
	Builder   *QueryBuilder
}

// NewConditionalBuilder creates a new conditional builder
func NewConditionalBuilder() *ConditionalBuilder {
	return &ConditionalBuilder{
		conditions: make([]*ConditionalClause, 0),
	}
}

// When adds a conditional clause
func (cb *ConditionalBuilder) When(condition bool, builder *QueryBuilder) *ConditionalBuilder {
	cb.conditions = append(cb.conditions, &ConditionalClause{
		Condition: condition,
		Builder:   builder,
	})
	return cb
}

// BuildConditional builds the conditional query
func (cb *ConditionalBuilder) BuildConditional() (string, []interface{}) {
	for _, clause := range cb.conditions {
		if clause.Condition {
			return clause.Builder.Build()
		}
	}
	return "", nil
}

// QueryAnalyzer analyzes queries for optimization
type QueryAnalyzer struct{}

// NewQueryAnalyzer creates a new query analyzer
func NewQueryAnalyzer() *QueryAnalyzer {
	return &QueryAnalyzer{}
}

// Analyze analyzes a query
func (qa *QueryAnalyzer) Analyze(query string) AdvancedQueryAnalysis {
	// Simplified analysis - would use SQL parser in production
	return AdvancedQueryAnalysis{
		HasJoins:            strings.Contains(query, "JOIN"),
		HasSubqueries:       strings.Contains(query, "SELECT") && strings.Count(query, "SELECT") > 1,
		HasAggregations:     strings.Contains(query, "COUNT") || strings.Contains(query, "SUM") || strings.Contains(query, "AVG"),
		HasOrderBy:          strings.Contains(query, "ORDER BY"),
		HasGroupBy:          strings.Contains(query, "GROUP BY"),
		EstimatedComplexity: estimateComplexity(query),
	}
}

// AdvancedQueryAnalysis represents query analysis results
type AdvancedQueryAnalysis struct {
	HasJoins            bool
	HasSubqueries       bool
	HasAggregations     bool
	HasOrderBy          bool
	HasGroupBy          bool
	EstimatedComplexity int
}

// estimateComplexity estimates query complexity
func estimateComplexity(query string) int {
	complexity := 1
	
	if strings.Contains(query, "JOIN") {
		complexity += strings.Count(query, "JOIN")
	}
	if strings.Contains(query, "SELECT") {
		complexity += strings.Count(query, "SELECT") - 1
	}
	if strings.Contains(query, "WHERE") {
		complexity += strings.Count(query, "AND") + strings.Count(query, "OR")
	}
	
	return complexity
}

// QueryFormatter formats queries for readability
type QueryFormatter struct{}

// NewQueryFormatter creates a new query formatter
func NewQueryFormatter() *QueryFormatter {
	return &QueryFormatter{}
}

// Format formats a query string
func (qf *QueryFormatter) Format(query string) string {
	// Basic formatting - would use SQL formatter in production
	formatted := strings.ReplaceAll(query, "SELECT", "\nSELECT")
	formatted = strings.ReplaceAll(formatted, "FROM", "\nFROM")
	formatted = strings.ReplaceAll(formatted, "WHERE", "\nWHERE")
	formatted = strings.ReplaceAll(formatted, "JOIN", "\nJOIN")
	formatted = strings.ReplaceAll(formatted, "ORDER BY", "\nORDER BY")
	formatted = strings.ReplaceAll(formatted, "GROUP BY", "\nGROUP BY")
	return formatted
}

// QueryValidator validates queries
type QueryValidator struct{}

// NewQueryValidator creates a new query validator
func NewQueryValidator() *QueryValidator {
	return &QueryValidator{}
}

// Validate validates a query
func (qv *QueryValidator) Validate(query string) error {
	// Basic validation - would use SQL parser in production
	if !strings.HasPrefix(strings.TrimSpace(query), "SELECT") {
		return fmt.Errorf("query must start with SELECT")
	}
	
	if strings.Contains(query, ";") && !strings.HasSuffix(strings.TrimSpace(query), ";") {
		return fmt.Errorf("query contains semicolon but doesn't end with one")
	}
	
	return nil
}

// QueryExecutor provides advanced query execution
type QueryExecutor struct {
	builder *QueryBuilder
}

// NewQueryExecutor creates a new query executor
func NewQueryExecutor(builder *QueryBuilder) *QueryExecutor {
	return &QueryExecutor{
		builder: builder,
	}
}

// Execute executes the query
func (qe *QueryExecutor) Execute(ctx context.Context, executor Executor) (Rows, error) {
	query, args := qe.builder.Build()
	return executor.Query(ctx, query, args...)
}

// ExecuteOne executes the query and returns one row
func (qe *QueryExecutor) ExecuteOne(ctx context.Context, executor Executor) Row {
	query, args := qe.builder.Build()
	return executor.QueryRow(ctx, query, args...)
}

// Explain generates EXPLAIN query
func (qe *QueryExecutor) Explain(ctx context.Context, executor Executor) (Rows, error) {
	query, args := qe.builder.Build()
	explainQuery := "EXPLAIN " + query
	return executor.Query(ctx, explainQuery, args...)
}

