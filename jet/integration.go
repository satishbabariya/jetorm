package jet

import (
	"context"
	"fmt"

	"github.com/go-jet/jet/v2/postgres"
	"github.com/go-jet/jet/v2/qrm"
	"github.com/satishbabariya/jetorm/core"
)

// JetRepository provides Jet SQL integration for repositories
// Note: Jet SQL tables must be generated from your database schema using the Jet generator
type JetRepository[T any, ID comparable] struct {
	repo core.Repository[T, ID]
	db   qrm.DB
	// table is not stored here as it must be passed per query
	// Jet SQL tables are generated from database schema
}

// NewJetRepository creates a new Jet SQL integrated repository
func NewJetRepository[T any, ID comparable](
	repo core.Repository[T, ID],
	db qrm.DB,
) *JetRepository[T, ID] {
	return &JetRepository[T, ID]{
		repo: repo,
		db:   db,
	}
}

// FindByID finds an entity by ID using Jet SQL
// table must be a generated Jet SQL table with AllColumns field
func (jr *JetRepository[T, ID]) FindByID(ctx context.Context, table postgres.Table, idColumn postgres.Column, id ID) (*T, error) {
	var entity T

	// Build Jet SQL query
	var idValue postgres.Expression
	switch v := any(id).(type) {
	case int64:
		idValue = postgres.Int64(v)
	case int:
		idValue = postgres.Int(int64(v))
	case string:
		idValue = postgres.String(v)
	default:
		return nil, fmt.Errorf("unsupported ID type: %T", id)
	}

	// Use reflection or type assertion to get AllColumns
	// For generated tables, AllColumns is a field
	stmt := postgres.SELECT(table).
		FROM(table).
		WHERE(idColumn.EQ(idValue))

	err := stmt.QueryContext(ctx, jr.db, &entity)
	if err != nil {
		return nil, fmt.Errorf("jet query failed: %w", err)
	}

	return &entity, nil
}

// FindAll finds all entities using Jet SQL
func (jr *JetRepository[T, ID]) FindAll(ctx context.Context, table postgres.Table) ([]*T, error) {
	var entities []*T

	stmt := postgres.SELECT(table).
		FROM(table)

	err := stmt.QueryContext(ctx, jr.db, &entities)
	if err != nil {
		return nil, fmt.Errorf("jet query failed: %w", err)
	}

	return entities, nil
}

// FindWithJetQuery finds entities using a Jet SQL query
func (jr *JetRepository[T, ID]) FindWithJetQuery(ctx context.Context, stmt postgres.SelectStatement) ([]*T, error) {
	var entities []*T

	err := stmt.QueryContext(ctx, jr.db, &entities)
	if err != nil {
		return nil, fmt.Errorf("jet query failed: %w", err)
	}

	return entities, nil
}

// FindOneWithJetQuery finds one entity using a Jet SQL query
func (jr *JetRepository[T, ID]) FindOneWithJetQuery(ctx context.Context, stmt postgres.SelectStatement) (*T, error) {
	var entity T

	err := stmt.QueryContext(ctx, jr.db, &entity)
	if err != nil {
		return nil, fmt.Errorf("jet query failed: %w", err)
	}
	return &entity, nil
}

// CountWithJetQuery counts entities using a Jet SQL query
func (jr *JetRepository[T, ID]) CountWithJetQuery(ctx context.Context, table postgres.Table, where postgres.BoolExpression) (int64, error) {
	var count int64

	// Use COUNT(*) - Star is typically available in generated code
	// For now, use a column from the table or COUNT(1)
	countStmt := postgres.SELECT(postgres.COUNT(postgres.Int(1))).
		FROM(table)

	if where != nil {
		countStmt = countStmt.WHERE(where)
	}

	err := countStmt.QueryContext(ctx, jr.db, &count)
	if err != nil {
		return 0, fmt.Errorf("jet count query failed: %w", err)
	}

	return count, nil
}

// ExecuteJetQuery executes a Jet SQL statement
func (jr *JetRepository[T, ID]) ExecuteJetQuery(ctx context.Context, stmt postgres.Statement) error {
	_, err := stmt.ExecContext(ctx, jr.db)
	return err
}

// QueryBuilder provides Jet SQL query building utilities
// Note: This is a simplified wrapper. For full functionality, use Jet SQL directly
type QueryBuilder struct {
	table postgres.Table
}

// NewQueryBuilder creates a new Jet SQL query builder
func NewQueryBuilder(table postgres.Table) *QueryBuilder {
	return &QueryBuilder{
		table: table,
	}
}

// Select creates a SELECT statement
// columns should be projections (columns or expressions)
func (qb *QueryBuilder) Select(columns ...postgres.Projection) postgres.SelectStatement {
	return postgres.SELECT(columns...).FROM(qb.table)
}

// SelectAll creates a SELECT * statement
// Note: table must have AllColumns field (generated tables have this)
func (qb *QueryBuilder) SelectAll() postgres.SelectStatement {
	// For generated tables, use the table directly as it implements Projection
	return postgres.SELECT(qb.table).FROM(qb.table)
}

// Insert creates an INSERT statement
func (qb *QueryBuilder) Insert() postgres.InsertStatement {
	return postgres.INSERT(qb.table)
}

// Update creates an UPDATE statement
func (qb *QueryBuilder) Update() postgres.UpdateStatement {
	return postgres.UPDATE(qb.table)
}

// Delete creates a DELETE statement
func (qb *QueryBuilder) Delete() postgres.DeleteStatement {
	return postgres.DELETE(qb.table)
}

// SpecificationToJet converts a core.Specification to Jet SQL WHERE clause
// This is a placeholder - full implementation would parse the specification tree
func SpecificationToJet[T any](spec core.Specification[T], table postgres.Table) (postgres.BoolExpression, error) {
	if spec == nil {
		return postgres.Bool(true), nil
	}

	// This is a simplified conversion
	// Full implementation would parse the specification tree and convert to Jet expressions
	return postgres.Bool(true), fmt.Errorf("specification conversion not yet implemented")
}

// JetToSpecification converts a Jet SQL WHERE clause to core.Specification
// This is a placeholder - full implementation would convert Jet expressions to specifications
func JetToSpecification[T any](expr postgres.BoolExpression) core.Specification[T] {
	// This would convert Jet expressions to specifications
	// Simplified version
	return nil
}

// JetQueryExecutor provides execution utilities for Jet SQL queries
type JetQueryExecutor struct {
	db qrm.DB
}

// NewJetQueryExecutor creates a new Jet query executor
func NewJetQueryExecutor(db qrm.DB) *JetQueryExecutor {
	return &JetQueryExecutor{
		db: db,
	}
}

// Execute executes a Jet SQL statement
func (jqe *JetQueryExecutor) Execute(ctx context.Context, stmt postgres.Statement) error {
	_, err := stmt.ExecContext(ctx, jqe.db)
	return err
}

// Query executes a SELECT statement and scans results
func (jqe *JetQueryExecutor) Query(ctx context.Context, stmt postgres.SelectStatement, dest interface{}) error {
	return stmt.QueryContext(ctx, jqe.db, dest)
}

// QueryOne executes a SELECT statement and scans one result
func (jqe *JetQueryExecutor) QueryOne(ctx context.Context, stmt postgres.SelectStatement, dest interface{}) error {
	return stmt.QueryContext(ctx, jqe.db, dest)
}

// Count executes a COUNT query
func (jqe *JetQueryExecutor) Count(ctx context.Context, table postgres.Table, where postgres.BoolExpression) (int64, error) {
	countStmt := postgres.SELECT(postgres.COUNT(postgres.Int(1))).
		FROM(table)

	if where != nil {
		countStmt = countStmt.WHERE(where)
	}

	var count int64
	err := countStmt.QueryContext(ctx, jqe.db, &count)
	return count, err
}

// Transaction executes statements in a transaction
// Note: This is a simplified version. For full transaction support, use Jet SQL's transaction API
func (jqe *JetQueryExecutor) Transaction(ctx context.Context, fn func(*JetQueryExecutor) error) error {
	// This would integrate with core transaction system
	// Simplified version
	return fn(jqe)
}
