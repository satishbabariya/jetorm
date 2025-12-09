# Jet SQL Integration Examples

This directory contains examples demonstrating Jet SQL integration with JetORM.

## Overview

Jet SQL integration allows you to:
- Build type-safe SQL queries using Jet SQL builder
- Combine Jet SQL with JetORM repositories
- Use Jet SQL for complex queries while using JetORM for simple CRUD

## Examples

### 1. Basic SELECT Query

```go
stmt := postgres.SELECT(userTable.AllColumns()).
    FROM(userTable).
    WHERE(userTable.Status.EQ(postgres.String("active")))

users, err := jetRepo.FindWithJet(ctx, stmt)
```

### 2. Complex Query with JOIN

```go
stmt := postgres.SELECT(
    userTable.Email,
    profileTable.Bio,
).FROM(userTable.
    INNER_JOIN(profileTable, userTable.ID.EQ(profileTable.UserID)),
).WHERE(userTable.Status.EQ(postgres.String("active")))
```

### 3. Using Jet Helpers

```go
qb := jet.NewQueryBuilder(userTable)
stmt := qb.SelectAll().
    WHERE(jet.And(
        jet.Equal(userTable.Status, "active"),
        jet.GreaterThan(userTable.Age, 18),
    ))
```

### 4. Combining Jet SQL with JetORM

```go
// Simple CRUD with JetORM
user, err := userRepo.Save(ctx, &User{...})

// Complex queries with Jet SQL
stmt := postgres.SELECT(userTable.AllColumns()).
    FROM(userTable).
    WHERE(userTable.Age.GT(postgres.Int(25)))
users, err := jetRepo.FindWithJet(ctx, stmt)
```

## Code Generation

Generate Jet SQL models from your database:

```bash
jet -dsn="postgres://user:pass@localhost/dbname" -path=./jet/models
```

This generates type-safe table definitions that you can use with JetORM.

## Benefits

1. **Type Safety** - Compile-time query validation
2. **Powerful Queries** - Complex SQL queries with type safety
3. **Integration** - Works seamlessly with JetORM
4. **Flexibility** - Use Jet SQL for complex queries, JetORM for simple CRUD

## Next Steps

- Generate Jet SQL models from your database
- Use Jet SQL for complex queries
- Combine with JetORM repositories
- Leverage type safety throughout your application

