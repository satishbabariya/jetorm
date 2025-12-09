# JetORM Packages Overview

This document provides an overview of all packages in the JetORM project.

## Core Packages

### `core/`
Core functionality including repository pattern, entity management, transactions, and specifications.

**Key Files:**
- `repository.go` - Generic repository interface
- `base_repository.go` - Base repository implementation
- `entity.go` - Entity metadata extraction
- `database.go` - Database connection and pooling
- `transaction.go` - Transaction management
- `specification.go` - Specification/Criteria API

## Supporting Packages

### `hooks/`
Lifecycle hooks for entity operations.

**Features:**
- Before/After hooks for Create, Update, Delete, Save operations
- Auditing hooks for automatic timestamp and user tracking
- Soft delete hooks

**Example:**
```go
hooks := hooks.NewHooks[User]()
hooks.RegisterBeforeCreate(hooks.CreateAuditHook[User]())
hooks.RegisterBeforeUpdate(hooks.AuditHook[User]())
```

### `migration/`
Database migration management.

**Features:**
- Migration version tracking
- Apply/Rollback migrations
- Schema generation from entity definitions
- Migration history

**Example:**
```go
migrator := migration.NewMigrator(db)
migrator.Initialize(ctx)
migration := migration.Migration{
    Version: 1,
    Name: "create_users_table",
    UpSQL: "CREATE TABLE users...",
    DownSQL: "DROP TABLE users",
}
migrator.Apply(ctx, migration)
```

### `query/`
Dynamic query building utilities.

**Features:**
- Fluent query builder API
- WHERE, ORDER BY, LIMIT, OFFSET support
- GROUP BY and HAVING clauses
- Type-safe query construction

**Example:**
```go
qb := query.NewQueryBuilder("users")
qb.WhereEqual("status", "active")
qb.OrderBy("created_at", "DESC")
qb.Limit(10)
query, args := qb.Build()
```

### `testing/`
Testing utilities and mocks.

**Features:**
- Mock repository implementation
- Test fixtures management
- Time utilities for testing

**Example:**
```go
mockRepo := testing.NewMockRepository[User, int64]()
mockRepo.FindByIDFunc = func(ctx context.Context, id int64) (*User, error) {
    return &User{ID: id, Email: "test@example.com"}, nil
}
```

### `tx/`
Advanced transaction support with propagation.

**Features:**
- Transaction propagation (REQUIRED, REQUIRES_NEW, etc.)
- Nested transaction support
- Context-based transaction management

**Example:**
```go
tm := tx.NewTransactionManager(db)
err := tm.Execute(ctx, tx.PropagationRequired, func(tx *sql.Tx) error {
    // Transaction logic
    return nil
})
```

### `logging/`
SQL query logging and monitoring.

**Features:**
- SQL query logging
- Slow query detection
- Transaction event logging
- Query formatting utilities

**Example:**
```go
sqlLogger := logging.NewSQLLogger(logger)
sqlLogger.SetSlowThreshold(100 * time.Millisecond)
sqlLogger.LogQuery(ctx, query, args, duration)
```

### `generator/`
Code generation for repository implementations.

**Features:**
- Method name analyzer (FindByX patterns)
- Interface parser
- Code generator
- CLI tool (jetorm-gen)

**Example:**
```go
//go:generate jetorm-gen -type=User -interface=UserRepository -input=user.go -output=user_repository_gen.go
```

## Package Dependencies

```
core/
├── hooks/ (optional)
├── migration/ (optional)
├── query/ (optional)
├── tx/ (optional)
└── logging/ (optional)

generator/
└── core/ (for type information)

testing/
└── core/ (for mocks)
```

## Usage Patterns

### Basic Repository
```go
db := core.MustConnect(config)
repo := core.NewBaseRepository[User, int64](db)
```

### With Hooks
```go
hooks := hooks.NewHooks[User]()
hooks.RegisterBeforeSave(hooks.CreateAuditHook[User]())
// Integrate hooks into repository
```

### With Migrations
```go
migrator := migration.NewMigrator(db)
migrator.ApplyAll(ctx, migrations)
```

### With Query Builder
```go
qb := query.NewQueryBuilder("users")
qb.WhereEqual("status", "active").Limit(10)
query, args := qb.Build()
```

## Testing

All packages include comprehensive test coverage. Run tests with:

```bash
go test ./...
```

## Examples

See the `examples/` directory for complete usage examples:
- `basic/` - Basic CRUD operations
- `codegen/` - Code generation workflow
- `advanced/` - Advanced features

