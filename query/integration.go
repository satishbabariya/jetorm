package query

import (
	"context"
	"fmt"

	"github.com/satishbabariya/jetorm/core"
)

// RepositoryQuery integrates query building with repository pattern
type RepositoryQuery[T any, ID comparable] struct {
	repo      core.Repository[T, ID]
	query     *ComposableQuery[T]
	tableName string
}

// NewRepositoryQuery creates a new repository query
func NewRepositoryQuery[T any, ID comparable](repo core.Repository[T, ID], tableName string) *RepositoryQuery[T, ID] {
	return &RepositoryQuery[T, ID]{
		repo:      repo,
		query:     NewComposableQuery[T](tableName),
		tableName: tableName,
	}
}

// WithSpecification sets a specification for the query
func (rq *RepositoryQuery[T, ID]) WithSpecification(spec core.Specification[T]) *RepositoryQuery[T, ID] {
	rq.query.WithSpecification(spec)
	return rq
}

// Select sets the columns to select
func (rq *RepositoryQuery[T, ID]) Select(cols ...string) *RepositoryQuery[T, ID] {
	rq.query.Select(cols...)
	return rq
}

// Where adds a WHERE clause
func (rq *RepositoryQuery[T, ID]) Where(condition string, args ...interface{}) *RepositoryQuery[T, ID] {
	rq.query.Where(condition, args...)
	return rq
}

// WhereEqual adds an equality WHERE clause
func (rq *RepositoryQuery[T, ID]) WhereEqual(column string, value interface{}) *RepositoryQuery[T, ID] {
	rq.query.WhereEqual(column, value)
	return rq
}

// OrderBy adds an ORDER BY clause
func (rq *RepositoryQuery[T, ID]) OrderBy(column string, direction string) *RepositoryQuery[T, ID] {
	rq.query.OrderBy(column, direction)
	return rq
}

// Limit sets the LIMIT clause
func (rq *RepositoryQuery[T, ID]) Limit(limit int) *RepositoryQuery[T, ID] {
	rq.query.Limit(limit)
	return rq
}

// Offset sets the OFFSET clause
func (rq *RepositoryQuery[T, ID]) Offset(offset int) *RepositoryQuery[T, ID] {
	rq.query.Offset(offset)
	return rq
}

// Find executes the query and returns results
func (rq *RepositoryQuery[T, ID]) Find(ctx context.Context) ([]*T, error) {
	query, args := rq.query.Build()
	return rq.repo.Query(ctx, query, args...)
}

// FindOne executes the query and returns a single result
func (rq *RepositoryQuery[T, ID]) FindOne(ctx context.Context) (*T, error) {
	rq.query.Limit(1)
	query, args := rq.query.Build()
	return rq.repo.QueryOne(ctx, query, args...)
}

// Count executes a COUNT query
func (rq *RepositoryQuery[T, ID]) Count(ctx context.Context) (int64, error) {
	query, args := rq.query.BuildCount()
	
	// Execute COUNT query - this is a simplified version
	// In a real implementation, we'd need to handle the COUNT result properly
	_, err := rq.repo.QueryOne(ctx, query, args...)
	if err != nil {
		return 0, err
	}
	// Placeholder - would need proper COUNT result handling
	return 0, fmt.Errorf("Count not fully implemented - use repository.CountWithSpec instead")
}

// Exists checks if any rows match the query
func (rq *RepositoryQuery[T, ID]) Exists(ctx context.Context) (bool, error) {
	rq.query.Select("1")
	rq.query.Limit(1)
	query, args := rq.query.Build()
	
	results, err := rq.repo.Query(ctx, query, args...)
	if err != nil {
		return false, err
	}
	return len(results) > 0, nil
}

// Paginate executes the query with pagination
func (rq *RepositoryQuery[T, ID]) Paginate(ctx context.Context, pageable core.Pageable) (*core.Page[T], error) {
	// Use default if nil or unpaged
	if pageable.Size < 0 {
		pageable = core.PageRequest(0, 20)
	}
	
	// Calculate offset
	offset := pageable.Page * pageable.Size
	
	// Apply pagination
	rq.query.Offset(offset)
	rq.query.Limit(pageable.Size)
	
	// Apply sorting
	if len(pageable.Sort.Orders) > 0 {
		for _, order := range pageable.Sort.Orders {
			direction := "ASC"
			if order.Direction == core.Desc {
				direction = "DESC"
			}
			rq.query.OrderBy(order.Field, direction)
		}
	}
	
	// Get results
	results, err := rq.Find(ctx)
	if err != nil {
		return nil, err
	}
	
	// Get total count - simplified version
	countQuery, countArgs := rq.query.BuildCount()
	countResults, err := rq.repo.Query(ctx, countQuery, countArgs...)
	if err != nil {
		return nil, err
	}
	
	total := int64(len(countResults))
	
	// Build page
	page := &core.Page[T]{
		Content:          results,
		TotalElements:    total,
		TotalPages:       int((total + int64(pageable.Size) - 1) / int64(pageable.Size)),
		Size:             pageable.Size,
		Number:           pageable.Page,
		NumberOfElements: len(results),
		First:            pageable.Page == 0,
		Last:             int64(len(results)) < int64(pageable.Size),
		Empty:            len(results) == 0,
		Sort:             pageable.Sort,
		Pageable:         pageable,
	}
	
	return page, nil
}

// QueryBuilderHelper provides helper functions for building queries
type QueryBuilderHelper struct{}

// NewQueryBuilderHelper creates a new query builder helper
func NewQueryBuilderHelper() *QueryBuilderHelper {
	return &QueryBuilderHelper{}
}

// BuildSelectQuery builds a SELECT query with all clauses
func (h *QueryBuilderHelper) BuildSelectQuery(tableName string, options ...QueryOption) (string, []interface{}) {
	qb := NewQueryBuilder(tableName)
	
	for _, option := range options {
		option(qb)
	}
	
	return qb.Build()
}

// QueryOption is a function that modifies a QueryBuilder
type QueryOption func(*QueryBuilder)

// WithSelect sets the SELECT columns
func WithSelect(cols ...string) QueryOption {
	return func(qb *QueryBuilder) {
		qb.Select(cols...)
	}
}

// WithWhere adds a WHERE clause
func WithWhere(condition string, args ...interface{}) QueryOption {
	return func(qb *QueryBuilder) {
		qb.Where(condition, args...)
	}
}

// WithOrderBy adds an ORDER BY clause
func WithOrderBy(column, direction string) QueryOption {
	return func(qb *QueryBuilder) {
		qb.OrderBy(column, direction)
	}
}

// WithLimit sets the LIMIT
func WithLimit(limit int) QueryOption {
	return func(qb *QueryBuilder) {
		qb.Limit(limit)
	}
}

// WithOffset sets the OFFSET
func WithOffset(offset int) QueryOption {
	return func(qb *QueryBuilder) {
		qb.Offset(offset)
	}
}

// WithGroupBy adds a GROUP BY clause
func WithGroupBy(columns ...string) QueryOption {
	return func(qb *QueryBuilder) {
		qb.GroupBy(columns...)
	}
}

// Example usage:
// query, args := helper.BuildSelectQuery("users",
//     WithSelect("id", "email", "name"),
//     WithWhere("status = $1", "active"),
//     WithOrderBy("created_at", "DESC"),
//     WithLimit(10),
// )

