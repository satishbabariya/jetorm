# JetORM Tutorial

A comprehensive tutorial for using JetORM.

## Table of Contents

1. [Getting Started](#getting-started)
2. [Basic CRUD Operations](#basic-crud-operations)
3. [Query Building](#query-building)
4. [Transactions](#transactions)
5. [Pagination](#pagination)
6. [Validation](#validation)
7. [Caching](#caching)
8. [Hooks](#hooks)
9. [Advanced Features](#advanced-features)

## Getting Started

### Installation

```bash
go get github.com/satishbabariya/jetorm
```

### Database Connection

```go
import "github.com/satishbabariya/jetorm/core"

config := core.Config{
    Host:     "localhost",
    Port:     5432,
    Database: "mydb",
    User:     "postgres",
    Password: "password",
}

db, err := core.Connect(config)
if err != nil {
    log.Fatal(err)
}
defer db.Close()
```

### Creating a Repository

```go
type User struct {
    ID    int64  `db:"id" jet:"primary_key,auto_increment"`
    Email string `db:"email" jet:"unique,not_null"`
    Name  string `db:"name"`
}

repo, err := core.NewBaseRepository[User, int64](db)
if err != nil {
    log.Fatal(err)
}
```

## Basic CRUD Operations

### Create

```go
user := &User{
    Email: "john@example.com",
    Name:  "John Doe",
}

saved, err := repo.Save(ctx, user)
```

### Read

```go
// Find by ID
user, err := repo.FindByID(ctx, 1)

// Find all
users, err := repo.FindAll(ctx)

// Find by IDs
users, err := repo.FindAllByIDs(ctx, []int64{1, 2, 3})
```

### Update

```go
user.Name = "Jane Doe"
updated, err := repo.Update(ctx, user)
```

### Delete

```go
// Delete entity
err := repo.Delete(ctx, user)

// Delete by ID
err := repo.DeleteByID(ctx, 1)
```

## Query Building

### Using Specifications

```go
// Find by email
spec := core.Equal[User]("email", "john@example.com")
user, err := repo.FindOne(ctx, spec)

// Find with multiple conditions
spec := core.And(
    core.Equal[User]("active", true),
    core.GreaterThan[User]("age", 18),
)
users, err := repo.FindAllWithSpec(ctx, spec)

// Find with OR
spec := core.Or(
    core.Equal[User]("role", "admin"),
    core.Equal[User]("role", "moderator"),
)
users, err := repo.FindAllWithSpec(ctx, spec)
```

### Using Query Builder

```go
import "github.com/satishbabariya/jetorm/query"

builder := query.NewQueryBuilder("users")
builder.WhereEqual("status", "active")
       .OrderBy("created_at", query.Desc)
       .Limit(10)

queryStr, args := builder.Build()
```

## Transactions

### Simple Transaction

```go
err := db.Transaction(ctx, func(tx *core.Tx) error {
    txRepo := repo.WithTx(tx)
    
    user1, err := txRepo.Save(ctx, &User{Email: "user1@example.com"})
    if err != nil {
        return err
    }
    
    user2, err := txRepo.Save(ctx, &User{Email: "user2@example.com"})
    if err != nil {
        return err
    }
    
    return nil
})
```

### Transaction with Options

```go
err := db.TransactionWithOptions(ctx, core.TxOptions{
    Isolation: core.Serializable,
    ReadOnly:  false,
}, func(tx *core.Tx) error {
    // Transaction logic
    return nil
})
```

## Pagination

### Basic Pagination

```go
pageable := core.PageRequest(0, 10)
page, err := repo.FindAllPaged(ctx, pageable)

fmt.Printf("Total: %d, Page: %d\n", page.TotalElements, page.Number)
for _, user := range page.Content {
    fmt.Println(user.Email)
}
```

### Pagination with Sorting

```go
order := core.Order{
    Field:     "created_at",
    Direction: core.Desc,
}
pageable := core.PageRequest(0, 10, order)
page, err := repo.FindAllPaged(ctx, pageable)
```

## Validation

### Basic Validation

```go
validator := core.NewValidator()
validator.RegisterRule("Email", core.Email())
validator.RegisterRule("Age", core.Range(18, 120))

err := validator.Validate(user)
```

### Advanced Validation

```go
validator.RegisterRule("Password", core.All(
    core.MinLength(8),
    core.HasLetter(),
    core.HasDigit(),
    core.HasSpecialChar(),
))
```

## Caching

### Basic Caching

```go
cache := core.NewInMemoryCache()
cachedRepo := core.NewCachedRepository(
    repo,
    cache,
    "User",
    5*time.Minute,
)

user, err := cachedRepo.FindByID(ctx, 1)
```

## Hooks

### Lifecycle Hooks

```go
import "github.com/satishbabariya/jetorm/hooks"

userHooks := hooks.NewHooks[User]()
userHooks.RegisterBeforeCreate(func(ctx context.Context, user *User) error {
    // Before create logic
    return nil
})

userHooks.RegisterAfterUpdate(func(ctx context.Context, user *User) error {
    // After update logic
    return nil
})
```

## Advanced Features

### Full-Featured Repository

```go
fullRepo := core.NewFullFeaturedRepository(
    baseRepo,
    cache,
    "User",
    5*time.Minute,
    hooks,
    validator,
    profiler,
    db,
)
```

### Batch Operations

```go
users := []*User{...}
err := core.OptimizedBatchSave(ctx, repo, users, 100)
```

### Helper Functions

```go
// Find or create
user, err := core.FindOrCreate(ctx, repo, finder, creator)

// Batch find
users, err := core.BatchFind(ctx, repo, ids, 100)

// Exists check
exists, err := core.Exists(ctx, repo, id)
```

## Best Practices

1. **Always use context** - Pass context for cancellation and timeouts
2. **Use transactions** - For multiple related operations
3. **Validate input** - Always validate entities before saving
4. **Use caching** - For frequently accessed data
5. **Monitor performance** - Track query performance
6. **Handle errors** - Use error helpers for better error handling

## Next Steps

- See `examples/` directory for more examples
- Read `API_REFERENCE.md` for complete API documentation
- Check `PACKAGES.md` for package overview

