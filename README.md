# JetORM

[![Go Version](https://img.shields.io/badge/Go-1.24+-00ADD8?style=flat&logo=go)](https://golang.org)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

**JetORM** is a next-generation Go database library that combines the type safety of Jet SQL builder, the performance of pgx driver, the convenience of sqlx, and integrated migration management. It provides a Spring Data JPA-like developer experience with zero-boilerplate CRUD operations, automatic query generation, and compile-time type safety.

## ✨ Features

- 🔒 **Type Safe**: Compile-time query validation
- ⚡ **High Performance**: Native pgx PostgreSQL driver
- 🎯 **Zero Boilerplate**: Auto-generated repository implementations
- 🔄 **Transactions**: Declarative transaction management
- 📄 **Pagination**: Built-in pagination and sorting
- 🔍 **Query Builder**: Type-safe query construction
- 🗄️ **Migrations**: Integrated migration support
- 📊 **Logging**: SQL query logging and slow query detection
- 🧪 **Testable**: Easy mocking and test utilities

## 🚀 Quick Start

### Installation

```bash
go get github.com/satishbabariya/jetorm
```

### Define Your Entity

```go
type User struct {
    ID        int64     `db:"id" jet:"primary_key,auto_increment"`
    Email     string    `db:"email" jet:"unique,not_null"`
    Username  string    `db:"username" jet:"unique,not_null"`
    FullName  string    `db:"full_name"`
    Age       int       `db:"age"`
    Status    string    `db:"status" jet:"default:'active'"`
    CreatedAt time.Time `db:"created_at" jet:"auto_now_add,not_null"`
    UpdatedAt time.Time `db:"updated_at" jet:"auto_now,not_null"`
}
```

### Connect and Use

```go
package main

import (
    "context"
    "log"
    
    "github.com/satishbabariya/jetorm/core"
)

func main() {
    // Connect to database
    db, err := core.Connect(core.Config{
        Host:     "localhost",
        Port:     5432,
        Database: "myapp",
        User:     "postgres",
        Password: "secret",
        SSLMode:  "disable",
    })
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    // Create repository
    userRepo, err := core.NewBaseRepository[User, int64](db)
    if err != nil {
        log.Fatal(err)
    }

    ctx := context.Background()

    // Create
    user := &User{
        Email:    "john@example.com",
        Username: "johndoe",
        FullName: "John Doe",
        Age:      30,
        Status:   "active",
    }
    saved, err := userRepo.Save(ctx, user)

    // Read
    found, err := userRepo.FindByID(ctx, saved.ID)

    // Update
    found.Age = 31
    updated, err := userRepo.Save(ctx, found)

    // Delete
    err = userRepo.DeleteByID(ctx, saved.ID)

    // Pagination
    page, err := userRepo.FindAllPaged(ctx, core.PageRequest(0, 10,
        core.Order{Field: "created_at", Direction: core.Desc},
    ))

    // Transaction
    err = db.Transaction(ctx, func(tx *core.Tx) error {
        txRepo := userRepo.WithTx(tx)
        _, err := txRepo.Save(ctx, user)
        return err
    })
}
```

## 📖 Documentation

### Configuration

```go
config := core.Config{
    // Connection
    Host:     "localhost",
    Port:     5432,
    Database: "myapp",
    User:     "postgres",
    Password: "secret",
    SSLMode:  "disable",
    
    // Connection Pool
    MaxOpenConns:    25,
    MaxIdleConns:    5,
    ConnMaxLifetime: 5 * time.Minute,
    ConnMaxIdleTime: 5 * time.Minute,
    
    // Logging
    LogSQL:          true,
    LogLevel:        core.DebugLevel,
    LogSlowQueries:  100 * time.Millisecond,
    
    // Performance
    PreparedStmts:   true,
    QueryTimeout:    30 * time.Second,
}
```

### Entity Tags

| Tag | Description | Example |
|-----|-------------|---------|
| `primary_key` | Primary key field | `jet:"primary_key"` |
| `auto_increment` | Auto-increment integer | `jet:"auto_increment"` |
| `unique` | Unique constraint | `jet:"unique"` |
| `not_null` | NOT NULL constraint | `jet:"not_null"` |
| `index` | Create index | `jet:"index"` |
| `size:n` | VARCHAR size | `jet:"size:255"` |
| `default:value` | Default value | `jet:"default:'active'"` |
| `auto_now_add` | Set on insert | `jet:"auto_now_add"` |
| `auto_now` | Update on save | `jet:"auto_now"` |

### Repository Interface

```go
type Repository[T any, ID comparable] interface {
    // Basic CRUD
    Save(ctx context.Context, entity *T) (*T, error)
    SaveAll(ctx context.Context, entities []*T) ([]*T, error)
    FindByID(ctx context.Context, id ID) (*T, error)
    FindAll(ctx context.Context) ([]*T, error)
    FindAllByIDs(ctx context.Context, ids []ID) ([]*T, error)
    Delete(ctx context.Context, entity *T) error
    DeleteByID(ctx context.Context, id ID) error
    DeleteAll(ctx context.Context, entities []*T) error
    Count(ctx context.Context) (int64, error)
    ExistsById(ctx context.Context, id ID) (bool, error)
    
    // Pagination
    FindAllPaged(ctx context.Context, pageable Pageable) (*Page[T], error)
    
    // Transaction
    WithTx(tx *Tx) Repository[T, ID]
}
```

### Pagination

```go
// Create pageable request
pageable := core.PageRequest(
    0,    // page number (zero-based)
    20,   // page size
    core.Order{Field: "created_at", Direction: core.Desc},
    core.Order{Field: "username", Direction: core.Asc},
)

// Get page
page, err := userRepo.FindAllPaged(ctx, pageable)

// Access page data
fmt.Printf("Page %d of %d\n", page.Number+1, page.TotalPages)
fmt.Printf("Total elements: %d\n", page.TotalElements)
for _, user := range page.Content {
    fmt.Println(user.Username)
}
```

### Transactions

```go
// Simple transaction
err := db.Transaction(ctx, func(tx *core.Tx) error {
    txRepo := userRepo.WithTx(tx)
    
    user1, err := txRepo.Save(ctx, &User{...})
    if err != nil {
        return err // automatic rollback
    }
    
    user2, err := txRepo.Save(ctx, &User{...})
    if err != nil {
        return err // automatic rollback
    }
    
    return nil // automatic commit
})

// Transaction with options
err := db.TransactionWithOptions(ctx, core.TxOptions{
    Isolation:  core.Serializable,
    ReadOnly:   false,
    Deferrable: false,
}, func(tx *core.Tx) error {
    // Your transaction logic
    return nil
})
```

## 🎯 Roadmap

### Phase 1: Foundation ✅
- [x] Core abstractions (Repository, Entity, Config)
- [x] Database connection management (pgx integration)
- [x] Basic CRUD operations
- [x] Transaction support
- [x] Pagination and sorting

### Phase 2: Code Generation (In Progress)
- [ ] Interface parser
- [ ] Method name analyzer (FindByX, DeleteByY)
- [ ] Code generation templates
- [ ] CLI tool (jetorm-gen)

### Phase 3: Query Building
- [ ] Jet integration wrapper
- [ ] Specification/Criteria API
- [ ] Dynamic query composition
- [ ] Raw SQL escape hatch

### Phase 4: Migrations
- [ ] Migration runner (goose integration)
- [ ] Schema generator from entities
- [ ] Version tracking
- [ ] Migration validation

### Phase 5: Advanced Features
- [ ] Lifecycle hooks (Before/After operations)
- [ ] Auditing (created_at, updated_at)
- [ ] Soft delete support
- [ ] Relationship handling (one-to-many, many-to-many)
- [ ] Caching layer

## 📊 Comparison with Alternatives

| Feature | JetORM | GORM | ent | sqlc |
|---------|--------|------|-----|------|
| Type Safety | ✅ Full | ⚠️ Partial | ✅ Full | ✅ Full |
| Code Generation | ✅ Yes | ❌ No | ✅ Yes | ✅ Yes |
| Repository Pattern | ✅ Built-in | ⚠️ Manual | ⚠️ Manual | ❌ No |
| Spring JPA-like | ✅ Yes | ⚠️ Partial | ❌ No | ❌ No |
| Migration Support | 🚧 Planned | ✅ Yes | ✅ Yes | ⚠️ External |
| Transaction Support | ✅ Declarative | ✅ Yes | ✅ Yes | ⚠️ Manual |
| Performance | ⭐⭐⭐⭐⭐ | ⭐⭐⭐⭐ | ⭐⭐⭐⭐⭐ | ⭐⭐⭐⭐⭐ |

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- [pgx](https://github.com/jackc/pgx) - PostgreSQL driver
- [Jet](https://github.com/go-jet/jet) - Type-safe SQL builder
- [Spring Data JPA](https://spring.io/projects/spring-data-jpa) - API inspiration

## 📧 Contact

For questions and support, please open an issue on GitHub.

---

**Note**: JetORM is currently in active development. The API may change before the 1.0 release.

