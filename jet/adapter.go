package jet

import (
	"context"
	"fmt"

	"github.com/go-jet/jet/v2/postgres"
	"github.com/go-jet/jet/v2/qrm"
	"github.com/satishbabariya/jetorm/core"
)

// SpecificationAdapter adapts core.Specification to Jet SQL
type SpecificationAdapter struct {
	table postgres.Table
}

// NewSpecificationAdapter creates a new specification adapter
func NewSpecificationAdapter(table postgres.Table) *SpecificationAdapter {
	return &SpecificationAdapter{
		table: table,
	}
}

// ToJet converts a core.Specification to Jet SQL BoolExpression
func (sa *SpecificationAdapter) ToJet(spec core.Specification[interface{}]) (postgres.BoolExpression, error) {
	if spec == nil {
		return postgres.Bool(true), nil
	}

	// This is a simplified adapter
	// Full implementation would parse the specification tree and convert to Jet expressions
	// For now, return a placeholder
	return postgres.Bool(true), fmt.Errorf("specification conversion not yet implemented")
}

// RepositoryAdapter adapts Jet SQL to work with JetORM repositories
type RepositoryAdapter[T any, ID comparable] struct {
	repo  core.Repository[T, ID]
	db    qrm.DB
	table postgres.Table
}

// NewRepositoryAdapter creates a new repository adapter
func NewRepositoryAdapter[T any, ID comparable](
	repo core.Repository[T, ID],
	db qrm.DB,
	table postgres.Table,
) *RepositoryAdapter[T, ID] {
	return &RepositoryAdapter[T, ID]{
		repo:  repo,
		db:    db,
		table: table,
	}
}

// FindWithJet finds entities using Jet SQL query
func (ra *RepositoryAdapter[T, ID]) FindWithJet(ctx context.Context, stmt postgres.SelectStatement) ([]*T, error) {
	var entities []*T
	err := stmt.QueryContext(ctx, ra.db, &entities)
	if err != nil {
		return nil, fmt.Errorf("jet query failed: %w", err)
	}

	// Convert to pointers
	result := make([]*T, len(entities))
	for i := range entities {
		result[i] = &entities[i]
	}
	return result, nil
}

// FindOneWithJet finds one entity using Jet SQL query
func (ra *RepositoryAdapter[T, ID]) FindOneWithJet(ctx context.Context, stmt postgres.SelectStatement) (*T, error) {
	var entity T
	err := stmt.QueryContext(ctx, ra.db, &entity)
	if err != nil {
		return nil, fmt.Errorf("jet query failed: %w", err)
	}
	return &entity, nil
}

// CountWithJet counts entities using Jet SQL query
func (ra *RepositoryAdapter[T, ID]) CountWithJet(ctx context.Context, table postgres.Table, where postgres.BoolExpression) (int64, error) {
	// Create count query
	countStmt := postgres.SELECT(postgres.COUNT(postgres.Int(1))).
		FROM(table)

	if where != nil {
		countStmt = countStmt.WHERE(where)
	}

	var count int64
	err := countStmt.QueryContext(ctx, ra.db, &count)
	return count, err
}

// ExecuteWithJet executes a Jet SQL statement
func (ra *RepositoryAdapter[T, ID]) ExecuteWithJet(ctx context.Context, stmt postgres.Statement) error {
	_, err := stmt.ExecContext(ctx, ra.db)
	return err
}

// QueryBuilderAdapter adapts Jet SQL query builder to JetORM patterns
// Note: This is a simplified adapter. For full functionality, use Jet SQL directly with generated tables
type QueryBuilderAdapter struct {
	table postgres.Table
}

// NewQueryBuilderAdapter creates a new query builder adapter
func NewQueryBuilderAdapter(table postgres.Table) *QueryBuilderAdapter {
	return &QueryBuilderAdapter{
		table: table,
	}
}

// BuildSelect builds a SELECT statement
// columns should be projections (columns or expressions from generated tables)
// For generated tables, use table.AllColumns or specific columns
func (qba *QueryBuilderAdapter) BuildSelect(columns ...postgres.Projection) postgres.SelectStatement {
	if len(columns) == 0 {
		// Note: For generated tables, use table.AllColumns field
		// This is a placeholder - actual usage requires generated table types
		return nil
	}
	if len(columns) == 1 {
		return postgres.SELECT(columns[0]).FROM(qba.table)
	}
	return postgres.SELECT(columns[0], columns[1:]...).FROM(qba.table)
}

// BuildInsert builds an INSERT statement
// Note: Use Jet SQL INSERT directly with generated tables
func (qba *QueryBuilderAdapter) BuildInsert() interface{} {
	// Placeholder - use Jet SQL INSERT directly
	return nil
}

// BuildUpdate builds an UPDATE statement
// Note: Use Jet SQL UPDATE directly with generated tables
func (qba *QueryBuilderAdapter) BuildUpdate() interface{} {
	// Placeholder - use Jet SQL UPDATE directly
	return nil
}

// BuildDelete builds a DELETE statement
// Note: Use Jet SQL DELETE directly with generated tables
func (qba *QueryBuilderAdapter) BuildDelete() interface{} {
	// Placeholder - use Jet SQL DELETE directly
	return nil
}

// ConditionBuilder builds WHERE conditions
type ConditionBuilder struct {
	conditions []postgres.BoolExpression
}

// NewConditionBuilder creates a new condition builder
func NewConditionBuilder() *ConditionBuilder {
	return &ConditionBuilder{
		conditions: make([]postgres.BoolExpression, 0),
	}
}

// Add adds a condition
func (cb *ConditionBuilder) Add(condition postgres.BoolExpression) *ConditionBuilder {
	cb.conditions = append(cb.conditions, condition)
	return cb
}

// And combines conditions with AND
func (cb *ConditionBuilder) And() postgres.BoolExpression {
	if len(cb.conditions) == 0 {
		return postgres.Bool(true)
	}
	if len(cb.conditions) == 1 {
		return cb.conditions[0]
	}
	result := cb.conditions[0]
	for i := 1; i < len(cb.conditions); i++ {
		result = result.AND(cb.conditions[i])
	}
	return result
}

// Or combines conditions with OR
func (cb *ConditionBuilder) Or() postgres.BoolExpression {
	if len(cb.conditions) == 0 {
		return postgres.Bool(false)
	}
	if len(cb.conditions) == 1 {
		return cb.conditions[0]
	}
	result := cb.conditions[0]
	for i := 1; i < len(cb.conditions); i++ {
		result = result.OR(cb.conditions[i])
	}
	return result
}

// QueryComposer composes complex queries
type QueryComposer struct {
	table postgres.Table
}

// NewQueryComposer creates a new query composer
func NewQueryComposer(table postgres.Table) *QueryComposer {
	return &QueryComposer{
		table: table,
	}
}

// ComposeSelect composes a SELECT query
func (qc *QueryComposer) ComposeSelect(
	columns []postgres.Projection,
	where postgres.BoolExpression,
	orderBy []postgres.OrderByClause,
	limit *int64,
	offset *int64,
) postgres.SelectStatement {
	if len(columns) == 0 {
		return nil
	}
	
	var stmt postgres.SelectStatement
	if len(columns) == 1 {
		stmt = postgres.SELECT(columns[0]).FROM(qc.table)
	} else {
		stmt = postgres.SELECT(columns[0], columns[1:]...).FROM(qc.table)
	}

	if where != nil {
		stmt = stmt.WHERE(where)
	}

	if len(orderBy) > 0 {
		stmt = stmt.ORDER_BY(orderBy...)
	}

	if limit != nil {
		stmt = stmt.LIMIT(*limit)
	}

	if offset != nil {
		stmt = stmt.OFFSET(*offset)
	}

	return stmt
}

// ComposeJoin composes a JOIN query
func (qc *QueryComposer) ComposeJoin(
	columns []postgres.Projection,
	joinTable postgres.Table,
	joinCondition postgres.BoolExpression,
	joinType string,
	where postgres.BoolExpression,
) postgres.SelectStatement {
	if len(columns) == 0 {
		return nil
	}

	var stmt postgres.SelectStatement
	if len(columns) == 1 {
		stmt = postgres.SELECT(columns[0])
	} else {
		stmt = postgres.SELECT(columns[0], columns[1:]...)
	}

	var joinedTable postgres.ReadableTable
	switch joinType {
	case "INNER":
		joinedTable = qc.table.INNER_JOIN(joinTable, joinCondition)
	case "LEFT":
		joinedTable = qc.table.LEFT_JOIN(joinTable, joinCondition)
	case "RIGHT":
		joinedTable = qc.table.RIGHT_JOIN(joinTable, joinCondition)
	case "FULL":
		joinedTable = qc.table.FULL_JOIN(joinTable, joinCondition)
	default:
		joinedTable = qc.table.INNER_JOIN(joinTable, joinCondition)
	}

	stmt = stmt.FROM(joinedTable)

	if where != nil {
		stmt = stmt.WHERE(where)
	}

	return stmt
}

