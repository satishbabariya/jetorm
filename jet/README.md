# Jet SQL Integration

This package provides integration with [go-jet/jet](https://github.com/go-jet/jet), a type-safe SQL builder for Go.

## Overview

Jet SQL integration allows you to:
- Build type-safe SQL queries using Jet SQL builder
- Combine Jet SQL with JetORM repositories
- Use Jet SQL's powerful query building capabilities
- Maintain type safety throughout your queries

## Installation

```bash
go get github.com/go-jet/jet/v2
```

## Prerequisites

Jet SQL requires code generation from your database schema. First, install the Jet generator:

```bash
go install github.com/go-jet/jet/v2/cmd/jet@latest
```

Then generate Jet SQL models from your database:

```bash
jet -source=postgres -dsn="user=your_user password=your_password host=localhost port=5432 dbname=your_dbname" -schema=public -path=./generated
```

This generates type-safe table definitions that you can use with JetORM.

## Usage

### Basic Integration

```go
import (
    "github.com/go-jet/jet/v2/postgres"
    "github.com/go-jet/jet/v2/qrm"
    "github.com/satishbabariya/jetorm/jet"
    "path/to/your/generated/table"
    "path/to/your/generated/model"
)

// Create Jet repository
jetRepo := jet.NewJetRepository(repo, db)

// Use Jet SQL queries with generated tables
stmt := postgres.SELECT(table.Users.AllColumns).
    FROM(table.Users).
    WHERE(table.Users.Status.EQ(postgres.String("active")))

users, err := jetRepo.FindWithJetQuery(ctx, stmt)
```

### Query Builder

```go
// Use Jet SQL directly for query building
stmt := postgres.SELECT(
    table.Users.Email,
    table.Users.Username,
).FROM(table.Users).
    WHERE(table.Users.Status.EQ(postgres.String("active")))

// Execute using Jet repository
users, err := jetRepo.FindWithJetQuery(ctx, stmt)
```

### Helper Functions

```go
// Conditions
stmt := postgres.SELECT(table.Users.AllColumns).
    FROM(table.Users).
    WHERE(jet.And(
        jet.Equal(table.Users.Status, "active"),
        jet.GreaterThan(table.Users.Age, 18),
    ))

// Joins
stmt := postgres.SELECT(
    table.Users.Email,
    table.Profiles.Bio,
).FROM(table.Users.
    INNER_JOIN(table.Profiles, table.Users.ID.EQ(table.Profiles.UserID)),
).WHERE(table.Users.Status.EQ(postgres.String("active")))
```

### Combining Jet SQL with JetORM

```go
// Use JetORM for simple CRUD
user, err := userRepo.Save(ctx, &User{...})

// Use Jet SQL for complex queries
stmt := postgres.SELECT(table.Users.AllColumns).
    FROM(table.Users).
    WHERE(table.Users.Age.GT(postgres.Int(25)))
users, err := jetRepo.FindWithJetQuery(ctx, stmt)
```

## Benefits

1. **Type Safety** - Compile-time query validation
2. **Powerful Queries** - Complex SQL queries with type safety
3. **Integration** - Works seamlessly with JetORM
4. **Flexibility** - Use Jet SQL for complex queries, JetORM for simple CRUD

## Examples

See `example.go` for comprehensive examples of Jet SQL integration.

## Next Steps

- Generate Jet SQL models from your database
- Use Jet SQL for complex queries
- Combine with JetORM repositories for simple operations
- Leverage type safety throughout your application

## Note

Jet SQL integration requires generated table definitions from your database schema. The integration provides utilities and helpers, but the main query building should be done using Jet SQL's generated types directly.
