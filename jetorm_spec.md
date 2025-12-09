# JetORM - Project Specification & Design Document

**Version:** 1.0  
**Date:** December 2025  
**Status:** Draft

---

## Table of Contents

1. [Executive Summary](#1-executive-summary)
2. [Project Overview](#2-project-overview)
3. [Core Objectives](#3-core-objectives)
4. [Technical Architecture](#4-technical-architecture)
5. [Feature Specifications](#5-feature-specifications)
6. [API Design](#6-api-design)
7. [Implementation Plan](#7-implementation-plan)
8. [Testing Strategy](#8-testing-strategy)
9. [Documentation Requirements](#9-documentation-requirements)
10. [Future Roadmap](#10-future-roadmap)

---

## 1. Executive Summary

**JetORM** is a next-generation Go database library that combines the type safety of Jet SQL builder, the performance of pgx driver, the convenience of sqlx, and integrated migration management. It provides a Spring Data JPA-like developer experience with zero-boilerplate CRUD operations, automatic query generation, and compile-time type safety.

### Key Value Propositions

- **Zero Boilerplate**: Auto-generated repository implementations
- **Type Safe**: Compile-time query validation via Jet
- **High Performance**: Native pgx PostgreSQL driver
- **Developer Friendly**: Spring JPA-inspired API design
- **Production Ready**: Built-in migrations, pooling, transactions, logging

---

## 2. Project Overview

### 2.1 Problem Statement

Current Go database solutions require developers to:
- Write repetitive CRUD operations manually
- Choose between type safety OR convenience (not both)
- Manage migrations separately from application code
- Integrate multiple libraries without cohesive abstraction
- Write significant boilerplate for common patterns

### 2.2 Solution

JetORM provides a unified, opinionated framework that:
- Generates type-safe repository implementations from interface definitions
- Integrates Jet (type safety) + pgx (performance) + sqlx (convenience)
- Includes first-class migration support
- Offers Spring Data JPA familiarity for developers from JVM backgrounds

### 2.3 Target Audience

- Go developers building database-driven applications
- Teams migrating from Spring Boot/JPA to Go
- Projects requiring type safety without sacrificing productivity
- Microservices requiring consistent data access patterns

---

## 3. Core Objectives

### 3.1 Functional Objectives

1. **Repository Pattern Implementation**
   - Generic CRUD operations (Create, Read, Update, Delete)
   - Query method name parsing (FindByX, DeleteByY)
   - Custom query support with Jet integration
   - Pagination and sorting built-in

2. **Migration Management**
   - Versioned migration support
   - Up/Down migration execution
   - Schema generation from entity definitions
   - Migration history tracking

3. **Transaction Management**
   - Declarative transaction boundaries
   - Nested transaction support
   - Automatic rollback on errors
   - Transaction propagation

4. **Query Building**
   - Type-safe query construction via Jet
   - Specification/Criteria API
   - Dynamic query composition
   - Raw SQL escape hatch

### 3.2 Non-Functional Objectives

- **Performance**: Minimal overhead over raw pgx/Jet
- **Type Safety**: 100% compile-time query validation
- **Maintainability**: Generated code should be readable
- **Testability**: Easy mocking and test utilities
- **Extensibility**: Plugin architecture for custom behaviors

---

## 4. Technical Architecture

### 4.1 System Architecture

```
┌─────────────────────────────────────────────────────────┐
│                    Application Layer                     │
│              (User Repositories & Services)              │
└────────────────────────┬────────────────────────────────┘
                         │
┌────────────────────────▼────────────────────────────────┐
│                   JetORM Core Layer                      │
├──────────────────────────────────────────────────────────┤
│  Repository Interface  │  Transaction Manager            │
│  Query Builder         │  Migration Engine               │
│  Specification API     │  Connection Pool Manager        │
└────────────────────────┬────────────────────────────────┘
                         │
┌────────────────────────▼────────────────────────────────┐
│              Integration Layer                           │
├──────────────┬──────────────┬──────────────┬────────────┤
│   Jet SQL    │    sqlx      │     pgx      │  Migrations│
│   Builder    │  Extensions  │    Driver    │   (goose)  │
└──────────────┴──────────────┴──────────────┴────────────┘
                         │
┌────────────────────────▼────────────────────────────────┐
│                   PostgreSQL Database                    │
└──────────────────────────────────────────────────────────┘
```

### 4.2 Component Architecture

```
jetorm/
├── core/                           # Core abstractions
│   ├── database.go                 # DB connection & config
│   ├── repository.go               # Base repository interface
│   ├── transaction.go              # Transaction management
│   ├── context.go                  # Context utilities
│   └── entity.go                   # Entity metadata
│
├── generator/                      # Code generation
│   ├── parser.go                   # Parse Go interfaces
│   ├── analyzer.go                 # Method name analysis
│   ├── codegen.go                  # Generate implementations
│   ├── templates/                  # Go templates
│   │   ├── repository.tmpl
│   │   ├── query.tmpl
│   │   └── transaction.tmpl
│   └── cmd/jetorm-gen/            # CLI tool
│
├── query/                          # Query building
│   ├── builder.go                  # Jet wrapper
│   ├── specification.go            # Criteria API
│   ├── pagination.go               # Pagination support
│   ├── sort.go                     # Sorting utilities
│   └── joins.go                    # Relationship handling
│
├── migration/                      # Migration support
│   ├── migrator.go                 # Migration runner
│   ├── generator.go                # Schema generation
│   ├── versioning.go               # Version tracking
│   ├── sql/                        # SQL templates
│   └── embedded.go                 # Embed migrations
│
├── tx/                             # Transaction support
│   ├── manager.go                  # Transaction manager
│   ├── propagation.go              # Propagation rules
│   └── isolation.go                # Isolation levels
│
├── hooks/                          # Lifecycle hooks
│   ├── lifecycle.go                # Before/After hooks
│   ├── auditing.go                 # Audit fields
│   └── softdelete.go               # Soft delete support
│
├── logging/                        # Logging
│   ├── logger.go                   # Logger interface
│   ├── sql_logger.go               # SQL query logging
│   └── adapters/                   # Logger adapters
│
├── testing/                        # Test utilities
│   ├── mock.go                     # Mock repositories
│   ├── fixtures.go                 # Test data fixtures
│   └── assertions.go               # Custom assertions
│
└── examples/                       # Example applications
    ├── basic/
    ├── advanced/
    └── migration/
```

### 4.3 Technology Stack

| Component | Technology | Rationale |
|-----------|-----------|-----------|
| SQL Builder | Jet | Type-safe, generates from schema |
| PostgreSQL Driver | pgx v5 | Best performance, native protocol |
| SQL Extensions | sqlx | Convenience methods, struct scanning |
| Migration Tool | goose | Mature, embedded support |
| Code Generation | text/template | Standard library, no dependencies |
| Logging | slog | Standard library (Go 1.21+) |
| Testing | testify | Rich assertions, mocking |

---

## 5. Feature Specifications

### 5.1 Repository Pattern

#### 5.1.1 Generic Repository Interface

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
    
    // Specification
    FindOne(ctx context.Context, spec Specification[T]) (*T, error)
    FindAll(ctx context.Context, spec Specification[T]) ([]*T, error)
    Count(ctx context.Context, spec Specification[T]) (int64, error)
    
    // Transaction
    WithTx(tx *Tx) Repository[T, ID]
}
```

#### 5.1.2 Query Method Naming Convention

| Pattern | Example | Generated SQL |
|---------|---------|---------------|
| FindBy{Field} | FindByEmail | WHERE email = ? |
| FindBy{Field}And{Field} | FindByNameAndAge | WHERE name = ? AND age = ? |
| FindBy{Field}Or{Field} | FindByEmailOrUsername | WHERE email = ? OR username = ? |
| FindBy{Field}GreaterThan | FindByAgeGreaterThan | WHERE age > ? |
| FindBy{Field}LessThan | FindByAgeLessThan | WHERE age < ? |
| FindBy{Field}Like | FindByNameLike | WHERE name LIKE ? |
| FindBy{Field}In | FindByStatusIn | WHERE status IN (?) |
| FindBy{Field}IsNull | FindByDeletedAtIsNull | WHERE deleted_at IS NULL |
| FindBy{Field}OrderBy{Field}Asc | FindByActiveOrderByNameAsc | WHERE active = ? ORDER BY name ASC |
| CountBy{Field} | CountByStatus | SELECT COUNT(*) WHERE status = ? |
| DeleteBy{Field} | DeleteByEmail | DELETE WHERE email = ? |
| ExistsBy{Field} | ExistsByUsername | SELECT EXISTS WHERE username = ? |

### 5.2 Entity Definition

#### 5.2.1 Struct Tags

```go
type User struct {
    ID        int64      `db:"id" jet:"primary_key,auto_increment"`
    Email     string     `db:"email" jet:"unique,not_null,index"`
    Username  string     `db:"username" jet:"unique,not_null,size:50"`
    FullName  string     `db:"full_name" jet:"size:255"`
    Age       int        `db:"age" jet:"check:age >= 0"`
    Status    string     `db:"status" jet:"default:'active',index"`
    IsActive  bool       `db:"is_active" jet:"default:true"`
    CreatedAt time.Time  `db:"created_at" jet:"auto_now_add,not_null"`
    UpdatedAt time.Time  `db:"updated_at" jet:"auto_now,not_null"`
    DeletedAt *time.Time `db:"deleted_at" jet:"index"`
}
```

#### 5.2.2 Tag Specifications

| Tag | Description | Example |
|-----|-------------|---------|
| primary_key | Primary key field | `jet:"primary_key"` |
| auto_increment | Auto-increment integer | `jet:"auto_increment"` |
| unique | Unique constraint | `jet:"unique"` |
| not_null | NOT NULL constraint | `jet:"not_null"` |
| index | Create index | `jet:"index"` |
| index:name | Named index | `jet:"index:idx_email"` |
| size:n | VARCHAR size | `jet:"size:255"` |
| default:value | Default value | `jet:"default:'active'"` |
| check:expr | Check constraint | `jet:"check:age >= 0"` |
| auto_now_add | Set on insert | `jet:"auto_now_add"` |
| auto_now | Update on save | `jet:"auto_now"` |
| foreign_key | Foreign key | `jet:"foreign_key:users.id"` |

### 5.3 Configuration

```go
type Config struct {
    // Connection
    Driver          string        // "pgx" (default), future: "mysql", "sqlite"
    Host            string        // Database host
    Port            int           // Database port
    Database        string        // Database name
    User            string        // Database user
    Password        string        // Database password
    SSLMode         string        // SSL mode: disable, require, verify-ca, verify-full
    
    // Connection Pool
    MaxOpenConns    int           // Maximum open connections (default: 25)
    MaxIdleConns    int           // Maximum idle connections (default: 5)
    ConnMaxLifetime time.Duration // Connection max lifetime (default: 5m)
    ConnMaxIdleTime time.Duration // Connection max idle time (default: 5m)
    
    // Migrations
    MigrationsPath  string        // Path to migration files
    AutoMigrate     bool          // Auto-run migrations on startup
    MigrationTable  string        // Migration version table (default: "schema_migrations")
    
    // Jet Code Generation
    JetGenPath      string        // Path for generated Jet code
    JetGenPackage   string        // Package name for Jet code
    
    // Logging
    Logger          Logger        // Custom logger implementation
    LogLevel        LogLevel      // Log level: Debug, Info, Warn, Error
    LogSQL          bool          // Log SQL queries
    LogSlowQueries  time.Duration // Log queries slower than threshold
    
    // Performance
    PreparedStmts   bool          // Use prepared statements (default: true)
    QueryTimeout    time.Duration // Default query timeout (default: 30s)
    
    // Behavior
    SoftDelete      bool          // Enable soft delete globally
    CreatedAtField  string        // Custom created_at field name
    UpdatedAtField  string        // Custom updated_at field name
    DeletedAtField  string        // Custom deleted_at field name
}
```

### 5.4 Migration System

#### 5.4.1 Migration File Format

```sql
-- +jetorm Up
-- Description: Create users table
CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    username VARCHAR(50) UNIQUE NOT NULL,
    full_name VARCHAR(255),
    age INTEGER CHECK (age >= 0),
    status VARCHAR(20) DEFAULT 'active',
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP
);

CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_status ON users(status);
CREATE INDEX idx_users_deleted_at ON users(deleted_at);

-- +jetorm Down
DROP TABLE IF EXISTS users;
```

#### 5.4.2 Migration API

```go
type Migrator interface {
    // Run all pending migrations
    Up(ctx context.Context) error
    
    // Rollback last migration
    Down(ctx context.Context) error
    
    // Rollback to specific version
    DownTo(ctx context.Context, version int64) error
    
    // Get current version
    Version(ctx context.Context) (int64, error)
    
    // Generate migration from entity
    Generate(entity interface{}, name string) error
    
    // Validate migrations
    Validate(ctx context.Context) error
}
```

### 5.5 Transaction Management

```go
type TransactionManager interface {
    // Execute function in transaction
    Transaction(ctx context.Context, fn func(tx *Tx) error) error
    
    // Execute with options
    TransactionWithOptions(ctx context.Context, opts TxOptions, fn func(tx *Tx) error) error
    
    // Begin manual transaction
    Begin(ctx context.Context) (*Tx, error)
    BeginWithOptions(ctx context.Context, opts TxOptions) (*Tx, error)
}

type TxOptions struct {
    Isolation   IsolationLevel    // Read Uncommitted, Read Committed, Repeatable Read, Serializable
    ReadOnly    bool              // Read-only transaction
    Deferrable  bool              // Deferrable constraint checking
}

type Tx struct {
    // Commit transaction
    Commit() error
    
    // Rollback transaction
    Rollback() error
    
    // Get repository with transaction
    Repository(repo interface{}) interface{}
}
```

### 5.6 Specification API

```go
type Specification[T any] interface {
    ToSQL() (string, []interface{})
}

// Factory functions
func Where[T any](condition jet.BoolExpression) Specification[T]
func And[T any](specs ...Specification[T]) Specification[T]
func Or[T any](specs ...Specification[T]) Specification[T]
func Not[T any](spec Specification[T]) Specification[T]

// Usage example
spec := And(
    Where(User.Email.Like("%@gmail.com")),
    Where(User.Age.GT(Int(18))),
    Or(
        Where(User.Status.EQ(String("active"))),
        Where(User.Status.EQ(String("pending"))),
    ),
)
users, err := userRepo.FindAll(ctx, spec)
```

### 5.7 Pagination

```go
type Pageable struct {
    Page int         // Zero-based page number
    Size int         // Page size
    Sort Sort        // Sort specification
}

type Sort struct {
    Orders []Order
}

type Order struct {
    Field     string
    Direction Direction  // Asc, Desc
}

type Page[T any] struct {
    Content       []*T   // Page content
    Number        int    // Current page number
    Size          int    // Page size
    TotalElements int64  // Total elements
    TotalPages    int    // Total pages
    First         bool   // Is first page
    Last          bool   // Is last page
}

// Helper functions
func PageRequest(page, size int, sort ...Order) Pageable
func Unpaged() Pageable
```

---

## 6. API Design

### 6.1 Initialization

```go
// Simple connection
db, err := jetorm.Connect(jetorm.Config{
    Host:     "localhost",
    Port:     5432,
    Database: "myapp",
    User:     "postgres",
    Password: "secret",
})
defer db.Close()

// With full configuration
db, err := jetorm.Connect(jetorm.Config{
    Host:            "localhost",
    Port:            5432,
    Database:        "myapp",
    User:            "postgres",
    Password:        "secret",
    MaxOpenConns:    25,
    MaxIdleConns:    5,
    MigrationsPath:  "./db/migrations",
    AutoMigrate:     true,
    JetGenPath:      "./internal/db/jet",
    Logger:          slog.Default(),
    LogLevel:        jetorm.InfoLevel,
    LogSQL:          true,
})
```

### 6.2 Repository Definition

```go
//go:generate jetorm-gen -type=User -output=user_repository_gen.go

type User struct {
    ID        int64     `db:"id" jet:"primary_key,auto_increment"`
    Email     string    `db:"email" jet:"unique,not_null"`
    Username  string    `db:"username" jet:"unique,not_null"`
    Age       int       `db:"age"`
    Status    string    `db:"status" jet:"default:'active'"`
    CreatedAt time.Time `db:"created_at" jet:"auto_now_add"`
    UpdatedAt time.Time `db:"updated_at" jet:"auto_now"`
}

type UserRepository interface {
    jetorm.Repository[User, int64]
    
    // Custom query methods (auto-implemented)
    FindByEmail(ctx context.Context, email string) (*User, error)
    FindByUsername(ctx context.Context, username string) (*User, error)
    FindByAgeGreaterThan(ctx context.Context, age int) ([]*User, error)
    FindByStatusIn(ctx context.Context, statuses []string) ([]*User, error)
    CountByStatus(ctx context.Context, status string) (int64, error)
    DeleteByEmail(ctx context.Context, email string) error
    
    // Custom implementations (manually written)
    FindActiveUsersWithOrders(ctx context.Context) ([]*User, error)
}
```

### 6.3 Usage Examples

```go
func main() {
    db := jetorm.MustConnect(config)
    defer db.Close()
    
    // Create repository
    userRepo := NewUserRepository(db)
    
    // Create
    user := &User{
        Email:    "john@example.com",
        Username: "john_doe",
        Age:      25,
        Status:   "active",
    }
    created, err := userRepo.Save(ctx, user)
    
    // Find by ID
    found, err := userRepo.FindByID(ctx, created.ID)
    
    // Find by custom method
    byEmail, err := userRepo.FindByEmail(ctx, "john@example.com")
    
    // Find with specification
    spec := jetorm.And(
        jetorm.Where(User.Age.GT(jetorm.Int(18))),
        jetorm.Where(User.Status.EQ(jetorm.String("active"))),
    )
    adults, err := userRepo.FindAll(ctx, spec)
    
    // Pagination
    page, err := userRepo.FindAllPaged(ctx, jetorm.PageRequest(0, 20,
        jetorm.Order{Field: "created_at", Direction: jetorm.Desc},
    ))
    
    // Transaction
    err = db.Transaction(ctx, func(tx *jetorm.Tx) error {
        user.Age = 26
        _, err := userRepo.WithTx(tx).Save(ctx, user)
        if err != nil {
            return err
        }
        // Other operations...
        return nil
    })
    
    // Update
    user.Status = "inactive"
    updated, err := userRepo.Save(ctx, user)
    
    // Delete
    err = userRepo.Delete(ctx, user)
    
    // Count
    count, err := userRepo.Count(ctx)
}
```

---

## 7. Implementation Plan

### Phase 1: Foundation (Weeks 1-2)
- [ ] Project setup and structure
- [ ] Core abstractions (Repository, Entity, Config)
- [ ] Database connection management (pgx integration)
- [ ] Basic CRUD operations
- [ ] Unit tests for core components

### Phase 2: Code Generation (Weeks 3-4)
- [ ] Interface parser
- [ ] Method name analyzer
- [ ] Code generation templates
- [ ] CLI tool (jetorm-gen)
- [ ] Integration tests

### Phase 3: Query Building (Weeks 5-6)
- [ ] Jet integration wrapper
- [ ] Specification API
- [ ] Pagination support
- [ ] Sorting implementation
- [ ] Query logging

### Phase 4: Migrations (Weeks 7-8)
- [ ] Migration runner (goose integration)
- [ ] Schema generator from entities
- [ ] Version tracking
- [ ] Migration validation
- [ ] Rollback support

### Phase 5: Advanced Features (Weeks 9-10)
- [ ] Transaction management
- [ ] Lifecycle hooks
- [ ] Auditing (created_at, updated_at)
- [ ] Soft delete support
- [ ] Relationship handling (one-to-many, many-to-many)

### Phase 6: Testing & Polish (Weeks 11-12)
- [ ] Comprehensive test suite
- [ ] Performance benchmarks
- [ ] Documentation
- [ ] Example applications
- [ ] API refinement

---

## 8. Testing Strategy

### 8.1 Unit Tests
- Test each component in isolation
- Mock external dependencies
- Target: 80%+ code coverage

### 8.2 Integration Tests
- Test with real PostgreSQL database (testcontainers)
- End-to-end repository operations
- Migration scenarios
- Transaction behavior

### 8.3 Performance Tests
- Benchmark against raw pgx/Jet
- Connection pool behavior
- Query performance
- Memory usage

### 8.4 Example Test Structure

```go
func TestUserRepository_Save(t *testing.T) {
    db := setupTestDB(t)
    defer teardownTestDB(t, db)
    
    repo := NewUserRepository(db)
    
    user := &User{Email: "test@example.com", Username: "testuser"}
    saved, err := repo.Save(context.Background(), user)
    
    require.NoError(t, err)
    assert.NotZero(t, saved.ID)
    assert.Equal(t, "test@example.com", saved.Email)
}
```

---

## 9. Documentation Requirements

### 9.1 README.md
- Quick start guide
- Installation instructions
- Basic usage examples
- Link to full documentation

### 9.2 User Guide
- Getting started tutorial
- Configuration reference
- Entity definition guide
- Repository patterns
- Migration management
- Transaction handling
- Best practices

### 9.3 API Documentation
- GoDoc comments for all public APIs
- Code examples in documentation
- Migration from other ORMs guide

### 9.4 Examples
- Basic CRUD application
- REST API with pagination
- Complex queries and joins
- Transaction management
- Testing strategies

---

## 10. Future Roadmap

### Version 1.0 (MVP)
- PostgreSQL support only
- Core CRUD operations
- Basic query methods
- Migration support
- Transaction management

### Version 1.1
- MySQL support
- Relationship handling (has-many, belongs-to)
- Eager/lazy loading
- Caching layer

### Version 1.2
- SQLite support
- Optimistic locking
- Event listeners
- Query result streaming

### Version 2.0
- Multi-tenancy support
- Sharding support
- Read replicas
- Advanced caching strategies
- GraphQL integration

---

## Appendix A: Comparison with Alternatives

| Feature | JetORM | GORM | ent | sqlc |
|---------|--------|------|-----|------|
| Type Safety | ✅ Full | ⚠️ Partial | ✅ Full | ✅ Full |
| Code Generation | ✅ Yes | ❌ No | ✅ Yes | ✅ Yes |
| Repository Pattern | ✅ Built-in | ⚠️ Manual | ⚠️ Manual | ❌ No |
| Spring JPA-like | ✅ Yes | ⚠️ Partial | ❌ No | ❌ No |
| Migration Support | ✅ Integrated | ✅ Yes | ✅ Yes | ⚠️ External |
| Query Method Parsing | ✅ Yes | ❌ No | ❌ No | ❌ No |
| Transaction Support | ✅ Declarative | ✅ Yes | ✅ Yes | ⚠️ Manual |
| Performance | ⭐⭐⭐⭐⭐ | ⭐⭐⭐⭐ | ⭐⭐⭐⭐⭐ | ⭐⭐⭐⭐⭐ |

---

## Appendix B: Glossary

- **Entity**: A struct representing a database table row
- **Repository**: Data access interface for an entity
- **Specification**: Composable query criteria
- **Pageable**: Pagination and sorting request
- **Migration**: Database schema version change
- **Transaction**: Atomic database operation unit
- **Hook**: Lifecycle callback (before/after operations)

---

**Document Control**

| Version | Date | Author | Changes |
|---------|------|--------|---------|
| 1.0 | Dec 2025 | Team | Initial draft |

