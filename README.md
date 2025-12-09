# JetORM

[![Go Version](https://img.shields.io/badge/go-1.21+-00ADD8?style=flat&logo=go)](https://golang.org)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)
[![Status](https://img.shields.io/badge/status-production%20ready-success)]()

**JetORM** is a next-generation Go database library that combines type safety, performance, and developer productivity. Inspired by Spring Data JPA, JetORM provides a powerful, type-safe repository pattern for PostgreSQL with code generation, advanced query building, migrations, and comprehensive features.

## ✨ Features

### 🎯 Core Features
- ✅ **Generic Repository Pattern** - Type-safe CRUD operations
- ✅ **Specification API** - Composable query criteria
- ✅ **Transaction Management** - Declarative transactions with isolation levels
- ✅ **Pagination & Sorting** - Advanced pagination with cursor support
- ✅ **Raw Query Support** - Escape hatch for complex queries
- ✅ **Batch Operations** - Optimized bulk operations

### 🔧 Code Generation
- ✅ **30+ Query Patterns** - Automatic method name parsing
- ✅ **Interface Parsing** - Parse Go interfaces
- ✅ **Code Generation** - Generate repository implementations
- ✅ **CLI Tool** - `jetorm-gen` for code generation
- ✅ **Configuration System** - Flexible configuration

### 🔍 Query Building
- ✅ **Fluent Query Builder** - Dynamic SQL construction
- ✅ **Composable Queries** - Build complex queries
- ✅ **Join Support** - Inner, left, right, full outer joins
- ✅ **Subqueries** - Nested query support
- ✅ **CTEs** - Common Table Expressions
- ✅ **UNION Support** - Union queries
- ✅ **Window Functions** - Advanced SQL features

### 📦 Migrations
- ✅ **Migration Runner** - File-based migrations
- ✅ **Schema Generator** - Generate SQL from entities
- ✅ **Migration Validation** - Validate migration SQL
- ✅ **Rollback Support** - Rollback migrations
- ✅ **CLI Tool** - `jetorm-migrate` for migration management

### 🚀 Advanced Features
- ✅ **Relationships** - One-to-one, one-to-many, many-to-many
- ✅ **Caching Layer** - Pluggable cache interface with in-memory implementation
- ✅ **Lifecycle Hooks** - Before/After operations
- ✅ **Soft Delete** - Soft delete support
- ✅ **Auditing** - Created/Updated timestamps
- ✅ **Health Monitoring** - Database health checks
- ✅ **Performance Monitoring** - Query profiling and metrics
- ✅ **Metrics Collection** - Counter, Gauge, Histogram, Timer

### 🛡️ Validation & Error Handling
- ✅ **50+ Validation Rules** - Comprehensive validation patterns
- ✅ **Custom Validators** - Create custom validation rules
- ✅ **Enhanced Error Handling** - Rich error context
- ✅ **Error Codes** - Programmatic error handling

### 🧰 Utilities
- ✅ **100+ Helper Functions** - Common operations
- ✅ **Entity Utilities** - Reflection-based helpers
- ✅ **Collection Operations** - Transform, reduce, filter, map
- ✅ **Set Operations** - Intersect, difference, union
- ✅ **Retry Logic** - Exponential backoff
- ✅ **Debounce/Throttle** - Rate limiting utilities

## 📦 Installation

```bash
go get github.com/satishbabariya/jetorm
```

## 🚀 Quick Start

### 1. Define Your Entity

```go
package main

import (
    "time"
    "github.com/satishbabariya/jetorm/core"
)

type User struct {
    ID        int64     `db:"id" jet:"primary_key,auto_increment"`
    Email     string    `db:"email" jet:"unique,not_null" validate:"required,email"`
    Name      string    `db:"name" validate:"required,min:3"`
    Age       int       `db:"age" validate:"min:18"`
    CreatedAt time.Time `db:"created_at" jet:"auto_now_add"`
    UpdatedAt time.Time `db:"updated_at" jet:"auto_now"`
}
```

### 2. Connect to Database

```go
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

### 3. Create Repository

```go
repo, err := core.NewBaseRepository[User, int64](db)
if err != nil {
    log.Fatal(err)
}
```

### 4. Use Repository

```go
ctx := context.Background()

// Create
user := &User{
    Email: "john@example.com",
    Name:  "John Doe",
    Age:   25,
}
saved, err := repo.Save(ctx, user)

// Read
user, err := repo.FindByID(ctx, 1)

// Update
user.Name = "Jane Doe"
updated, err := repo.Update(ctx, user)

// Delete
err := repo.DeleteByID(ctx, 1)
```

## 📚 Usage Examples

### Query with Specifications

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

### Pagination

```go
pageable := core.PageRequest(0, 10, core.Order{
    Field:     "created_at",
    Direction: core.Desc,
})
page, err := repo.FindAllPaged(ctx, pageable)

fmt.Printf("Total: %d, Page: %d\n", page.TotalElements, page.Number)
for _, user := range page.Content {
    fmt.Println(user.Email)
}
```

### Transactions

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

### Validation

```go
validator := core.NewValidator()
validator.RegisterRule("Email", core.Email())
validator.RegisterRule("Password", core.All(
    core.MinLength(8),
    core.HasLetter(),
    core.HasDigit(),
    core.HasSpecialChar(),
))
validator.RegisterRule("Phone", core.PhoneNumber())
validator.RegisterRule("UUID", core.UUID())

err := validator.Validate(user)
```

### Caching

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

### Advanced Query Building

```go
import "github.com/satishbabariya/jetorm/query"

// Basic query builder
builder := query.NewQueryBuilder("users")
builder.WhereEqual("status", "active")
       .OrderBy("created_at", query.Desc)
       .Limit(10)
query, args := builder.Build()

// Advanced query builder with CTEs
advancedBuilder := query.NewAdvancedQueryBuilder("users")
cteBuilder := query.NewQueryBuilder("active_users")
cteBuilder.WhereEqual("status", "active")
advancedBuilder.WithCTE("active", cteBuilder)
query, args := advancedBuilder.BuildAdvanced()
```

### Batch Operations

```go
users := []*User{...}
err := core.OptimizedBatchSave(ctx, repo, users, 100)

// Batch find
ids := []int64{1, 2, 3, 4, 5}
users, err := core.BatchFind(ctx, repo, ids, 100)
```

### Metrics Collection

```go
// Repository metrics
metrics := core.NewRepositoryMetrics()
metrics.RecordOperation("Save", duration, err)
stats := metrics.GetOperationStats("Save")

// Counter
counter := core.NewCounter("requests")
counter.Inc()

// Timer
timer := core.NewTimer("operation")
timer.Time(func() {
    // Operation
})
```

### Helper Functions

```go
// Find or create
user, err := core.FindOrCreate(ctx, repo, finder, creator)

// Exists check
exists, err := core.Exists(ctx, repo, id)

// Collection operations
doubled := core.Transform(slice, func(x int) int { return x * 2 })
sum := core.Reduce(slice, 0, func(acc, val int) int { return acc + val })
filtered := core.FilterEntities(entities, func(e *User) bool { return e.Active })
```

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

// Use with all features enabled
user, err := fullRepo.Save(ctx, &User{...})
health := fullRepo.HealthCheck(ctx)
```

## 🛠️ Code Generation

### Generate Repository Code

```bash
# Using CLI
jetorm-gen -type=User -interface=UserRepository \
  -input=user.go -output=user_repository_gen.go \
  -package=repository

# Or use config file
jetorm-gen init
# Edit jetorm-gen.json
jetorm-gen generate
```

### Supported Query Patterns

JetORM supports 30+ query method patterns:

- `FindBy{Field}` - Find by field
- `FindBy{Field}And{Field}` - Find with AND conditions
- `FindBy{Field}Or{Field}` - Find with OR conditions
- `FindBy{Field}In` - Find by field in list
- `FindBy{Field}Like` - Find by field like pattern
- `FindBy{Field}GreaterThan` - Find by field greater than
- `FindBy{Field}LessThan` - Find by field less than
- `FindBy{Field}Between` - Find by field between values
- `FindBy{Field}IsNull` - Find by field is null
- `FindBy{Field}IsNotNull` - Find by field is not null
- `DeleteBy{Field}` - Delete by field
- `CountBy{Field}` - Count by field
- `ExistsBy{Field}` - Check existence by field
- `OrderBy{Field}` - Order by field
- `OrderBy{Field}Desc` - Order by field descending
- And many more...

## 📦 Migrations

### Create Migration

```bash
jetorm-migrate create -name add_user_email_index
```

### Apply Migrations

```bash
jetorm-migrate up -db="postgres://user:pass@localhost/dbname" -dir=./migrations
```

### Rollback Migration

```bash
jetorm-migrate down -db="postgres://user:pass@localhost/dbname" -dir=./migrations
```

### Check Status

```bash
jetorm-migrate status -db="postgres://user:pass@localhost/dbname" -dir=./migrations
```

## 📖 Documentation

- **[Getting Started](GETTING_STARTED.md)** - Detailed getting started guide
- **[Tutorial](docs/TUTORIAL.md)** - Comprehensive tutorial
- **[API Reference](docs/API_REFERENCE.md)** - Complete API documentation
- **[Packages](PACKAGES.md)** - Package overview
- **[Examples](examples/)** - Working code examples

## 🎯 Key Features in Detail

### Type Safety
- Full compile-time type checking
- Generic repository pattern
- Type-safe query building

### Performance
- Connection pooling
- Query caching
- Batch operations
- Performance monitoring
- Query optimization

### Developer Experience
- Zero boilerplate
- 100+ helper functions
- Comprehensive validation
- Rich error messages
- Extensive documentation

### Production Ready
- Error handling
- Health monitoring
- Metrics collection
- Transaction management
- Migration support

## 🏗️ Architecture

```
jetorm/
├── core/          # Core functionality
├── generator/     # Code generation
├── query/         # Query building
├── migration/     # Migrations
├── hooks/         # Lifecycle hooks
├── testing/       # Test utilities
├── tx/            # Advanced transactions
├── logging/       # Logging
└── examples/      # Examples
```

## 📊 Statistics

- **83+ Go Files** - Comprehensive implementation
- **16,000+ Lines of Code** - Production-ready codebase
- **17+ Test Files** - Extensive test coverage
- **100+ Test Cases** - Comprehensive testing
- **50+ Validation Rules** - Advanced validation
- **100+ Utility Functions** - Helper functions
- **30+ Query Patterns** - Code generation support
- **29 Documentation Files** - Extensive documentation

## 🔄 Comparison with Alternatives

| Feature | JetORM | GORM | ent | sqlc |
|---------|--------|------|-----|------|
| Type Safety | ✅ Full | ⚠️ Partial | ✅ Full | ✅ Full |
| Code Generation | ✅ Yes | ❌ No | ✅ Yes | ✅ Yes |
| Repository Pattern | ✅ Built-in | ⚠️ Manual | ⚠️ Manual | ❌ No |
| Spring JPA-like | ✅ Yes | ⚠️ Partial | ❌ No | ❌ No |
| Query Builder | ✅ Yes | ✅ Yes | ✅ Yes | ❌ No |
| Migrations | ✅ Built-in | ✅ Yes | ✅ Yes | ⚠️ External |
| Validation | ✅ 50+ Rules | ⚠️ Basic | ❌ No | ❌ No |
| Caching | ✅ Built-in | ⚠️ Plugin | ❌ No | ❌ No |
| Metrics | ✅ Built-in | ❌ No | ❌ No | ❌ No |
| Hooks | ✅ Built-in | ✅ Yes | ⚠️ Limited | ❌ No |

## 🎓 Learning Resources

1. **[Tutorial](docs/TUTORIAL.md)** - Step-by-step tutorial
2. **[Examples](examples/)** - Working code examples
3. **[API Reference](docs/API_REFERENCE.md)** - Complete API docs
4. **[Packages](PACKAGES.md)** - Package overview

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- Inspired by Spring Data JPA
- Built with PostgreSQL and pgx
- Uses Jet SQL for query building

## 🚀 Roadmap

- [x] Jet SQL builder integration
- [ ] Additional database drivers (MySQL, SQLite)
- [ ] Redis cache implementation
- [ ] Query optimization
- [ ] Performance benchmarks
- [ ] IDE language server

## 📞 Support

For questions, issues, or contributions, please open an issue on GitHub.

---

**Made with ❤️ for the Go community**
