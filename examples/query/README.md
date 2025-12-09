# Query Building Examples

This directory contains examples demonstrating JetORM's query building capabilities.

## Overview

Phase 3 introduces advanced query building features:
- Composable queries with specifications
- Join support (INNER, LEFT, RIGHT, FULL)
- Subquery support
- Dynamic query building
- Condition builders
- Repository-integrated queries
- PostgreSQL-specific features

## Examples

### Basic Query Builder

```go
qb := query.NewQueryBuilder("users")
qb.WhereEqual("status", "active")
qb.OrderBy("created_at", "DESC")
qb.Limit(10)

query, args := qb.Build()
```

### Composable Query with Specification

```go
spec := core.And(
    core.Equal("age", 18),
    core.Equal("status", "active"),
)

cq := query.NewComposableQuery[User]("users")
cq.WithSpecification(spec)
cq.OrderBy("email", "ASC")
cq.Limit(20)

query, args := cq.Build()
```

### Join Queries

```go
jq := query.NewJoinQuery[User]("users")
jq.InnerJoin("profiles", "users.id = profiles.user_id")
jq.WhereEqual("users.status", "active")
jq.Select("users.id", "users.email", "profiles.bio")

query, args := jq.Build()
```

### Condition Builder

```go
cb := query.NewConditionBuilder()
cb.Equal("status", "active")
cb.GreaterThan("age", 18)
cb.Like("email", "%@example.com")

whereClause, args := cb.Build()
```

### Dynamic Query Building

```go
status := "active"
minAge := 18

dq := query.NewDynamicQuery[User]("users")
dq.When(status != "", func(q *query.ComposableQuery[User]) *query.ComposableQuery[User] {
    return q.Where("status = $1", status)
})
dq.When(minAge > 0, func(q *query.ComposableQuery[User]) *query.ComposableQuery[User] {
    return q.Where("age >= $1", minAge)
})

query, args := dq.Build()
```

### Repository-Integrated Queries

```go
rq := query.NewRepositoryQuery(repo, "users")
rq.WhereEqual("status", "active")
rq.OrderBy("created_at", "DESC")
rq.Limit(10)

users, err := rq.Find(ctx)
```

### Pagination

```go
rq := query.NewRepositoryQuery(repo, "users")
rq.WhereEqual("status", "active")

pageable := core.PageRequest(0, 20, core.Order{
    Field:     "created_at",
    Direction: core.Desc,
})

page, err := rq.Paginate(ctx, pageable)
```

### PostgreSQL-Specific Features

```go
// Full-text search
cb := query.TextSearch("description", "search term")

// Array operations
cb = query.ArrayContains("tags", "golang")
cb = query.ArrayOverlaps("categories", []interface{}{"tech", "programming"})
```

## Running Examples

```bash
cd examples/query
go run query_examples.go
```

## Features

### Query Builder
- Fluent API for building SQL queries
- Type-safe parameter handling
- Support for WHERE, ORDER BY, LIMIT, OFFSET
- GROUP BY and HAVING clauses

### Composable Queries
- Integration with Specification API
- Chainable methods
- Reusable query components

### Join Support
- INNER JOIN
- LEFT JOIN
- RIGHT JOIN
- FULL OUTER JOIN

### Dynamic Queries
- Conditional query building
- Runtime query composition
- Flexible query construction

### Condition Builder
- Rich set of comparison operators
- Logical operators (AND, OR)
- PostgreSQL-specific features

### Repository Integration
- Seamless integration with repositories
- Pagination support
- Find, Count, Exists operations

## Best Practices

1. **Use ComposableQuery for complex queries** - Better integration with specifications
2. **Use ConditionBuilder for reusable conditions** - Can be combined and reused
3. **Use DynamicQuery for conditional logic** - Cleaner than if/else statements
4. **Use RepositoryQuery for repository integration** - Automatic pagination and result handling
5. **Prefer helper functions** - More readable and maintainable

## Next Steps

- See `query_examples.go` for complete examples
- Check `PACKAGES.md` for package overview
- Read `GETTING_STARTED.md` for basic usage

