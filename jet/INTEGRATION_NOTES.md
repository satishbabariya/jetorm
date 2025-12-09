# Jet SQL Integration Notes

## Overview

Jet SQL integration provides utilities and helpers for working with [go-jet/jet](https://github.com/go-jet/jet) in JetORM projects.

## Important Notes

1. **Code Generation Required**: Jet SQL requires code generation from your database schema. You cannot use Jet SQL without generating table definitions first.

2. **Generated Types**: All table definitions, columns, and types must be generated using the Jet SQL generator.

3. **API Differences**: Jet SQL's API is designed to work with generated code. Generic wrappers are limited.

## Installation

```bash
go get github.com/go-jet/jet/v2
go install github.com/go-jet/jet/v2/cmd/jet@latest
```

## Code Generation

Generate Jet SQL models from your database:

```bash
jet -source=postgres -dsn="user=your_user password=your_password host=localhost port=5432 dbname=your_dbname" -schema=public -path=./generated
```

This creates:
- `table/` - Type-safe table definitions
- `model/` - Go structs matching your database schema

## Usage Pattern

```go
import (
    "github.com/go-jet/jet/v2/postgres"
    "github.com/go-jet/jet/v2/qrm"
    "path/to/generated/table"
    "path/to/generated/model"
)

// Use generated tables
stmt := postgres.SELECT(table.Users.AllColumns).
    FROM(table.Users).
    WHERE(table.Users.Status.EQ(postgres.String("active")))

var users []model.Users
err := stmt.Query(db, &users)
```

## Integration with JetORM

The `jet` package provides:
- `JetRepository` - Wrapper for executing Jet SQL queries with JetORM
- `JetQueryExecutor` - Executor for Jet SQL statements
- Helper functions for common patterns

## Limitations

- Generic wrappers are limited due to Jet SQL's code generation requirement
- Most query building should be done using Jet SQL's generated types directly
- The integration provides utilities, not full abstraction

## Best Practices

1. Generate Jet SQL models from your database schema
2. Use Jet SQL for complex queries
3. Use JetORM for simple CRUD operations
4. Combine both as needed for your use case

