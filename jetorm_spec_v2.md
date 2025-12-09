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
    Update(ctx context.Context, entity *T) (*T, error)
    UpdateAll(ctx context.Context, entities []*T) ([]*T, error)
    FindByID(ctx context.Context, id ID) (*T, error)
    FindAll(ctx context.Context) ([]*T, error)
    FindAllByIDs(ctx context.Context, ids []ID) ([]*T, error)
    Delete(ctx context.Context, entity *T) error
    DeleteByID(ctx context.Context, id ID) error
    DeleteAll(ctx context.Context, entities []*T) error
    DeleteAllByIDs(ctx context.Context, ids []ID) error
    Count(ctx context.Context) (int64, error)
    ExistsById(ctx context.Context, id ID) (bool, error)
    
    // Pagination
    FindAllPaged(ctx context.Context, pageable Pageable) (*Page[T], error)
    
    // Specification
    FindOne(ctx context.Context, spec Specification[T]) (*T, error)
    FindAll(ctx context.Context, spec Specification[T]) ([]*T, error)
    FindAllPaged(ctx context.Context, spec Specification[T], pageable Pageable) (*Page[T], error)
    Count(ctx context.Context, spec Specification[T]) (int64, error)
    Exists(ctx context.Context, spec Specification[T]) (bool, error)
    Delete(ctx context.Context, spec Specification[T]) (int64, error)
    
    // Batch Operations
    SaveBatch(ctx context.Context, entities []*T, batchSize int) error
    
    // Transaction
    WithTx(tx *Tx) Repository[T, ID]
    
    // Raw Query Support
    Query(ctx context.Context, query string, args ...interface{}) ([]*T, error)
    QueryOne(ctx context.Context, query string, args ...interface{}) (*T, error)
    Exec(ctx context.Context, query string, args ...interface{}) (int64, error)
}
```

#### 5.1.2 Query Method Naming Convention

**Supported Keywords and Patterns:**

| Category | Keyword | Example | Generated SQL |
|----------|---------|---------|---------------|
| **Find Operations** |
| Simple Find | FindBy{Field} | FindByEmail | WHERE email = $1 |
| Multiple Fields | FindBy{Field}And{Field} | FindByNameAndAge | WHERE name = $1 AND age = $2 |
| | FindBy{Field}Or{Field} | FindByEmailOrUsername | WHERE email = $1 OR username = $2 |
| **Comparison Operators** |
| Greater Than | FindBy{Field}GreaterThan | FindByAgeGreaterThan | WHERE age > $1 |
| Greater or Equal | FindBy{Field}GreaterThanEqual | FindByAgeGreaterThanEqual | WHERE age >= $1 |
| Less Than | FindBy{Field}LessThan | FindByAgeLessThan | WHERE age < $1 |
| Less or Equal | FindBy{Field}LessThanEqual | FindByAgeLessThanEqual | WHERE age <= $1 |
| Between | FindBy{Field}Between | FindByAgeBetween | WHERE age BETWEEN $1 AND $2 |
| **String Operations** |
| Like | FindBy{Field}Like | FindByNameLike | WHERE name LIKE $1 |
| Contains | FindBy{Field}Containing | FindByNameContaining | WHERE name LIKE '%' || $1 || '%' |
| Starts With | FindBy{Field}StartingWith | FindByNameStartingWith | WHERE name LIKE $1 || '%' |
| Ends With | FindBy{Field}EndingWith | FindByNameEndingWith | WHERE name LIKE '%' || $1 |
| Not Like | FindBy{Field}NotLike | FindByNameNotLike | WHERE name NOT LIKE $1 |
| Ignore Case | FindBy{Field}IgnoreCase | FindByEmailIgnoreCase | WHERE LOWER(email) = LOWER($1) |
| **Collection Operations** |
| In | FindBy{Field}In | FindByStatusIn | WHERE status IN ($1, $2, ...) |
| Not In | FindBy{Field}NotIn | FindByStatusNotIn | WHERE status NOT IN ($1, $2, ...) |
| **Null Checks** |
| Is Null | FindBy{Field}IsNull | FindByDeletedAtIsNull | WHERE deleted_at IS NULL |
| Is Not Null | FindBy{Field}IsNotNull | FindByDeletedAtIsNotNull | WHERE deleted_at IS NOT NULL |
| **Boolean Operations** |
| True | FindBy{Field}True | FindByActiveTrue | WHERE active = true |
| False | FindBy{Field}False | FindByActiveFalse | WHERE active = false |
| **Sorting** |
| Order By Asc | FindBy{Field}OrderBy{Field}Asc | FindByStatusOrderByNameAsc | WHERE status = $1 ORDER BY name ASC |
| Order By Desc | FindBy{Field}OrderBy{Field}Desc | FindByStatusOrderByNameDesc | WHERE status = $1 ORDER BY name DESC |
| Multiple Order | OrderBy{Field}Asc{Field}Desc | OrderByNameAscAgeDesc | ORDER BY name ASC, age DESC |
| **Limiting Results** |
| First | FindFirstBy{Field} | FindFirstByStatus | WHERE status = $1 LIMIT 1 |
| Top N | FindTop{N}By{Field} | FindTop10ByStatus | WHERE status = $1 LIMIT 10 |
| **Count Operations** |
| Count | CountBy{Field} | CountByStatus | SELECT COUNT(*) WHERE status = $1 |
| Count Distinct | CountDistinctBy{Field} | CountDistinctByCity | SELECT COUNT(DISTINCT city) |
| **Existence Checks** |
| Exists | ExistsBy{Field} | ExistsByUsername | SELECT EXISTS(SELECT 1 WHERE username = $1) |
| **Delete Operations** |
| Delete | DeleteBy{Field} | DeleteByEmail | DELETE WHERE email = $1 |
| Delete Multiple | DeleteBy{Field}And{Field} | DeleteByStatusAndAge | DELETE WHERE status = $1 AND age = $2 |
| **Distinct** |
| Distinct | FindDistinctBy{Field} | FindDistinctByCity | SELECT DISTINCT * WHERE city = $1 |

**Complex Examples:**

```go
// Find users by multiple conditions with sorting
FindByAgeGreaterThanAndStatusInOrderByCreatedAtDesc(ctx, 18, []string{"active", "pending"})
// WHERE age > $1 AND status IN ($2, $3) ORDER BY created_at DESC

// Find first active user with name starting with prefix
FindFirstByActiveTrueAndNameStartingWith(ctx, "John")
// WHERE active = true AND name LIKE $1 || '%' LIMIT 1

// Count users in multiple cities
CountByStatusAndCityIn(ctx, "active", []string{"NYC", "LA", "SF"})
// SELECT COUNT(*) WHERE status = $1 AND city IN ($2, $3, $4)

// Find users created between dates
FindByCreatedAtBetweenAndDeletedAtIsNull(ctx, startDate, endDate)
// WHERE created_at BETWEEN $1 AND $2 AND deleted_at IS NULL

// Delete inactive users older than date
DeleteByActiveFalseAndCreatedAtLessThan(ctx, cutoffDate)
// DELETE WHERE active = false AND created_at < $1
```

**Naming Rules:**
1. Method names are case-sensitive and must follow exact patterns
2. Field names must match struct field names exactly (case-sensitive)
3. Multiple conditions are combined with `And` or `Or`
4. `OrderBy` clause comes after all `By` conditions
5. `First` and `Top{N}` must come at the beginning
6. Return types must match operation:
   - Find operations: `(*T, error)` or `([]*T, error)`
   - Count operations: `(int64, error)`
   - Exists operations: `(bool, error)`
   - Delete operations: `(int64, error)` returns rows affected

### 5.2 Entity Definition

#### 5.2.1 Struct Tags

```go
type User struct {
    // Primary Key
    ID        int64      `db:"id" jet:"primary_key,auto_increment"`
    
    // Unique constraints
    Email     string     `db:"email" jet:"unique,not_null,index:idx_email,size:255"`
    Username  string     `db:"username" jet:"unique,not_null,index:idx_username,size:50"`
    
    // Regular fields
    FullName  string     `db:"full_name" jet:"size:255"`
    Bio       string     `db:"bio" jet:"type:text"`
    Age       int        `db:"age" jet:"check:age >= 0 AND age <= 150"`
    Status    string     `db:"status" jet:"default:'active',check:status IN ('active','inactive','suspended'),index"`
    IsActive  bool       `db:"is_active" jet:"default:true,index"`
    Balance   float64    `db:"balance" jet:"type:decimal(10,2),default:0.00"`
    
    // JSON field
    Metadata  types.JSON `db:"metadata" jet:"type:jsonb"`
    
    // Timestamps (auto-managed)
    CreatedAt time.Time  `db:"created_at" jet:"auto_now_add,not_null,index"`
    UpdatedAt time.Time  `db:"updated_at" jet:"auto_now,not_null"`
    DeletedAt *time.Time `db:"deleted_at" jet:"index"` // For soft delete
    
    // Foreign keys
    CompanyID *int64     `db:"company_id" jet:"foreign_key:companies.id,on_delete:cascade,on_update:cascade"`
    RoleID    int64      `db:"role_id" jet:"foreign_key:roles.id,not_null,index:idx_role"`
    
    // Ignored field (not persisted)
    TempData  string     `db:"-" jet:"-"`
}

// Composite index example
type Product struct {
    ID         int64  `db:"id" jet:"primary_key,auto_increment"`
    SKU        string `db:"sku" jet:"not_null,composite_index:idx_sku_store:1"`
    StoreID    int64  `db:"store_id" jet:"not_null,composite_index:idx_sku_store:2"`
    CategoryID int64  `db:"category_id" jet:"index"`
}
```

#### 5.2.2 Tag Specifications

**Core Tags:**

| Tag | Description | Example | Generated DDL |
|-----|-------------|---------|---------------|
| **Identity** |
| primary_key | Primary key field | `jet:"primary_key"` | PRIMARY KEY |
| auto_increment | Auto-increment (serial) | `jet:"auto_increment"` | SERIAL/BIGSERIAL |
| **Constraints** |
| unique | Unique constraint | `jet:"unique"` | UNIQUE |
| not_null | NOT NULL constraint | `jet:"not_null"` | NOT NULL |
| check:expr | Check constraint | `jet:"check:age >= 0"` | CHECK (age >= 0) |
| default:value | Default value | `jet:"default:'active'"` | DEFAULT 'active' |
| **Data Types** |
| size:n | VARCHAR size | `jet:"size:255"` | VARCHAR(255) |
| type:name | Explicit type | `jet:"type:text"` | TEXT |
| type:decimal(p,s) | Decimal type | `jet:"type:decimal(10,2)"` | DECIMAL(10,2) |
| **Indexes** |
| index | Create basic index | `jet:"index"` | CREATE INDEX idx_{table}_{field} |
| index:name | Named index | `jet:"index:idx_email"` | CREATE INDEX idx_email |
| composite_index:name:order | Composite index | `jet:"composite_index:idx_name:1"` | CREATE INDEX idx_name ON (field1, field2) |
| unique_index | Unique index | `jet:"unique_index:idx_unique"` | CREATE UNIQUE INDEX |
| **Foreign Keys** |
| foreign_key:table.col | Foreign key reference | `jet:"foreign_key:users.id"` | FOREIGN KEY REFERENCES users(id) |
| on_delete:action | Delete action | `jet:"on_delete:cascade"` | ON DELETE CASCADE |
| on_update:action | Update action | `jet:"on_update:set_null"` | ON UPDATE SET NULL |
| **Timestamps** |
| auto_now_add | Set on insert only | `jet:"auto_now_add"` | Triggers on INSERT |
| auto_now | Update on every save | `jet:"auto_now"` | Triggers on INSERT/UPDATE |
| **Special** |
| - | Ignore field | `db:"-" jet:"-"` | Not persisted |

**Foreign Key Actions:**
- `cascade` - CASCADE
- `set_null` - SET NULL
- `set_default` - SET DEFAULT
- `restrict` - RESTRICT
- `no_action` - NO ACTION

**Advanced Examples:**

```go
// Multi-column unique constraint
type UserProfile struct {
    UserID   int64  `db:"user_id" jet:"unique_composite:user_email:1"`
    Email    string `db:"email" jet:"unique_composite:user_email:2"`
}

// Partial index (PostgreSQL)
type Document struct {
    ID        int64      `db:"id" jet:"primary_key"`
    Status    string     `db:"status" jet:"index:idx_active_docs,where:status='active'"`
    DeletedAt *time.Time `db:"deleted_at" jet:"index,where:deleted_at IS NULL"`
}

// JSON fields
type Settings struct {
    ID       int64      `db:"id" jet:"primary_key"`
    Config   types.JSON `db:"config" jet:"type:jsonb,default:'{}'"`
    Tags     []string   `db:"tags" jet:"type:text[]"` // PostgreSQL array
}

// Computed/Generated columns
type Order struct {
    ID         int64   `db:"id" jet:"primary_key"`
    Subtotal   float64 `db:"subtotal" jet:"not_null"`
    TaxRate    float64 `db:"tax_rate" jet:"not_null"`
    Total      float64 `db:"total" jet:"generated:subtotal * (1 + tax_rate),stored"`
}
```

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
    Isolation   IsolationLevel    // Transaction isolation level
    ReadOnly    bool              // Read-only transaction
    Deferrable  bool              // Deferrable constraint checking (PostgreSQL)
    Timeout     time.Duration     // Transaction timeout
}

type IsolationLevel int
const (
    LevelDefault IsolationLevel = iota
    LevelReadUncommitted
    LevelReadCommitted
    LevelRepeatableRead
    LevelSerializable
    LevelSnapshot // SQL Server
)

type Tx struct {
    db      *Database
    sqlxTx  *sqlx.Tx
    opts    TxOptions
}

// Transaction methods
func (tx *Tx) Commit() error
func (tx *Tx) Rollback() error
func (tx *Tx) SavePoint(name string) error
func (tx *Tx) RollbackTo(name string) error
func (tx *Tx) ReleaseSavePoint(name string) error

// Get repository bound to transaction
func (tx *Tx) Repository(repo interface{}) interface{}

// Execute raw SQL in transaction
func (tx *Tx) Exec(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
func (tx *Tx) Query(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)

// Usage examples

// Simple transaction
err := db.Transaction(ctx, func(tx *jetorm.Tx) error {
    user := &User{Email: "john@example.com"}
    saved, err := userRepo.WithTx(tx).Save(ctx, user)
    if err != nil {
        return err // Automatic rollback
    }
    
    order := &Order{UserID: saved.ID, Total: 100.0}
    _, err = orderRepo.WithTx(tx).Save(ctx, order)
    if err != nil {
        return err // Automatic rollback
    }
    
    return nil // Automatic commit
})

// Transaction with isolation level
err := db.TransactionWithOptions(ctx, jetorm.TxOptions{
    Isolation: jetorm.LevelSerializable,
    ReadOnly:  false,
    Timeout:   30 * time.Second,
}, func(tx *jetorm.Tx) error {
    // Critical operation requiring serializable isolation
    count, _ := userRepo.WithTx(tx).Count(ctx)
    // ... perform operations based on count
    return nil
})

// Manual transaction control
tx, err := db.Begin(ctx)
if err != nil {
    return err
}
defer tx.Rollback() // Rollback if not committed

user, err := userRepo.WithTx(tx).Save(ctx, user)
if err != nil {
    return err
}

order, err := orderRepo.WithTx(tx).Save(ctx, order)
if err != nil {
    return err
}

if err := tx.Commit(); err != nil {
    return err
}

// Savepoints for nested transactions
err := db.Transaction(ctx, func(tx *jetorm.Tx) error {
    user, _ := userRepo.WithTx(tx).Save(ctx, user)
    
    // Create savepoint
    tx.SavePoint("before_orders")
    
    for _, item := range orderItems {
        _, err := orderRepo.WithTx(tx).Save(ctx, item)
        if err != nil {
            // Rollback to savepoint, not entire transaction
            tx.RollbackTo("before_orders")
            break
        }
    }
    
    tx.ReleaseSavePoint("before_orders")
    return nil
})

// Read-only transaction for consistency
err := db.TransactionWithOptions(ctx, jetorm.TxOptions{
    ReadOnly: true,
    Isolation: jetorm.LevelRepeatableRead,
}, func(tx *jetorm.Tx) error {
    users, _ := userRepo.WithTx(tx).FindAll(ctx)
    orders, _ := orderRepo.WithTx(tx).FindAll(ctx)
    // Generate report with consistent snapshot
    return nil
})

// Transaction propagation (nested transactions)
func CreateUserWithProfile(ctx context.Context, user *User, profile *Profile) error {
    return db.Transaction(ctx, func(tx *jetorm.Tx) error {
        // This transaction can be nested
        saved, err := createUser(ctx, tx, user)
        if err != nil {
            return err
        }
        
        profile.UserID = saved.ID
        return createProfile(ctx, tx, profile)
    })
}

func createUser(ctx context.Context, tx *jetorm.Tx, user *User) (*User, error) {
    // Can participate in parent transaction or create new one
    return userRepo.WithTx(tx).Save(ctx, user)
}

// Deadlock retry logic
func withRetry(ctx context.Context, db *Database, fn func(*Tx) error) error {
    maxRetries := 3
    for i := 0; i < maxRetries; i++ {
        err := db.Transaction(ctx, fn)
        if err == nil {
            return nil
        }
        
        // Check if deadlock error
        if isDeadlock(err) && i < maxRetries-1 {
            time.Sleep(time.Duration(i+1) * 100 * time.Millisecond)
            continue
        }
        
        return err
    }
    return errors.New("max retries exceeded")
}

// Transaction middleware (for web handlers)
type TxMiddleware struct {
    db *Database
}

func (m *TxMiddleware) Handle(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        m.db.Transaction(r.Context(), func(tx *Tx) error {
            // Attach transaction to request context
            ctx := context.WithValue(r.Context(), "tx", tx)
            next.ServeHTTP(w, r.WithContext(ctx))
            return nil
        })
    })
}

// Get transaction from context
func getTx(ctx context.Context) *Tx {
    if tx, ok := ctx.Value("tx").(*Tx); ok {
        return tx
    }
    return nil
}
```

**Transaction Best Practices:**

```go
// 1. Keep transactions short
err := db.Transaction(ctx, func(tx *Tx) error {
    // ❌ BAD: External API call in transaction
    response := callExternalAPI()
    
    // ✅ GOOD: Only database operations
    return userRepo.WithTx(tx).Save(ctx, user)
})

// 2. Handle errors properly
err := db.Transaction(ctx, func(tx *Tx) error {
    user, err := userRepo.WithTx(tx).Save(ctx, user)
    if err != nil {
        return fmt.Errorf("save user: %w", err) // Preserves error chain
    }
    
    order.UserID = user.ID
    _, err = orderRepo.WithTx(tx).Save(ctx, order)
    if err != nil {
        return fmt.Errorf("save order: %w", err)
    }
    
    return nil
})

// 3. Use appropriate isolation level
// Read Committed (default) - Good for most cases
// Repeatable Read - Prevents non-repeatable reads
// Serializable - Strictest, prevents phantom reads

// 4. Avoid long-running transactions
ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
defer cancel()

err := db.TransactionWithOptions(ctx, jetorm.TxOptions{
    Timeout: 5 * time.Second,
}, func(tx *Tx) error {
    // Operations
    return nil
})
```

### 5.6 Specification API

```go
// Core Specification interface
type Specification[T any] interface {
    ToSQL() (string, []interface{})
    And(spec Specification[T]) Specification[T]
    Or(spec Specification[T]) Specification[T]
    Not() Specification[T]
}

// Builder functions
func Where[T any](condition jet.BoolExpression) Specification[T]
func And[T any](specs ...Specification[T]) Specification[T]
func Or[T any](specs ...Specification[T]) Specification[T]
func Not[T any](spec Specification[T]) Specification[T]

// Reusable specifications (like Spring Data JPA Specifications)
type UserSpecifications struct{}

func (s UserSpecifications) IsActive() Specification[User] {
    return Where(User.IsActive.EQ(Bool(true)))
}

func (s UserSpecifications) HasEmail(email string) Specification[User] {
    return Where(User.Email.EQ(String(email)))
}

func (s UserSpecifications) AgeGreaterThan(age int) Specification[User] {
    return Where(User.Age.GT(Int(age)))
}

func (s UserSpecifications) CreatedAfter(date time.Time) Specification[User] {
    return Where(User.CreatedAt.GT(Timestamp(date)))
}

func (s UserSpecifications) StatusIn(statuses ...string) Specification[User] {
    return Where(User.Status.IN(StringArray(statuses)))
}

func (s UserSpecifications) SearchByName(name string) Specification[User] {
    return Where(User.FullName.LIKE(String("%" + name + "%")))
}

// Usage examples
var specs UserSpecifications

// Simple specification
activeUsers := userRepo.FindAll(ctx, specs.IsActive())

// Combined specifications
spec := And(
    specs.IsActive(),
    specs.AgeGreaterThan(18),
    specs.StatusIn("premium", "gold"),
)
users := userRepo.FindAll(ctx, spec)

// Dynamic query building
func buildUserQuery(filters UserFilters) Specification[User] {
    conditions := []Specification[User]{}
    
    if filters.IsActive {
        conditions = append(conditions, specs.IsActive())
    }
    
    if filters.MinAge > 0 {
        conditions = append(conditions, specs.AgeGreaterThan(filters.MinAge))
    }
    
    if filters.Email != "" {
        conditions = append(conditions, specs.HasEmail(filters.Email))
    }
    
    if len(filters.Statuses) > 0 {
        conditions = append(conditions, specs.StatusIn(filters.Statuses...))
    }
    
    if len(conditions) == 0 {
        return nil
    }
    
    return And(conditions...)
}

// Complex nested specifications
advancedSpec := Or(
    And(
        specs.IsActive(),
        specs.AgeGreaterThan(18),
    ),
    And(
        specs.StatusIn("premium"),
        specs.CreatedAfter(time.Now().AddDate(0, -1, 0)),
    ),
)

// With pagination
page := userRepo.FindAllPaged(ctx, advancedSpec, PageRequest(0, 20, 
    Order{Field: "created_at", Direction: Desc},
))

// Negation
notActiveUsers := userRepo.FindAll(ctx, Not(specs.IsActive()))

// Fluent API
spec := specs.IsActive().
    And(specs.AgeGreaterThan(18)).
    Or(specs.StatusIn("premium"))
```

**Advanced Specification Patterns:**

```go
// Generic specification factory
type SpecificationFactory[T any] struct{}

func (f SpecificationFactory[T]) Equal(field jet.Column, value interface{}) Specification[T] {
    return Where(field.EQ(value))
}

func (f SpecificationFactory[T]) Like(field jet.StringExpression, pattern string) Specification[T] {
    return Where(field.LIKE(String(pattern)))
}

func (f SpecificationFactory[T]) Between(field jet.NumericExpression, min, max interface{}) Specification[T] {
    return And(
        Where(field.GT_EQ(min)),
        Where(field.LT_EQ(max)),
    )
}

// Reusable cross-entity specifications
type AuditSpecifications[T any] struct{}

func (s AuditSpecifications[T]) CreatedBetween(start, end time.Time) Specification[T] {
    // Assumes entity has CreatedAt field
    return And(
        Where(/* CreatedAt */.GT_EQ(Timestamp(start))),
        Where(/* CreatedAt */.LT_EQ(Timestamp(end))),
    )
}

func (s AuditSpecifications[T]) NotDeleted() Specification[T] {
    return Where(/* DeletedAt */.IS_NULL())
}
```

### 5.7 Pagination & Sorting

```go
// Pageable interface
type Pageable interface {
    GetPageNumber() int
    GetPageSize() int
    GetOffset() int
    GetSort() Sort
    Next() Pageable
    Previous() Pageable
    First() Pageable
}

// Page result
type Page[T any] struct {
    Content          []*T      // Page content
    Pageable         Pageable  // Pageable that produced this page
    TotalElements    int64     // Total elements across all pages
    TotalPages       int       // Total number of pages
    Size             int       // Page size
    Number           int       // Current page number (zero-based)
    NumberOfElements int       // Elements in current page
    First            bool      // Is first page
    Last             bool      // Is last page
    Empty            bool      // Is empty page
    Sort             Sort      // Sort applied
}

// Sort specification
type Sort struct {
    Orders []Order
}

type Order struct {
    Property  string
    Direction Direction
    NullHandling NullHandling // NULLS FIRST or NULLS LAST
}

type Direction int
const (
    Asc Direction = iota
    Desc
)

type NullHandling int
const (
    NullsNative NullHandling = iota // Database default
    NullsFirst
    NullsLast
)

// Builder functions
func PageRequest(page, size int, orders ...Order) Pageable
func Unpaged() Pageable
func OffsetRequest(offset, limit int, orders ...Order) Pageable

// Sort builders
func SortBy(property string, direction Direction) Sort
func (s Sort) And(order Order) Sort
func (s Sort) Ascending(property string) Sort
func (s Sort) Descending(property string) Sort

// Usage examples

// Simple pagination
page1 := userRepo.FindAllPaged(ctx, PageRequest(0, 20))
page2 := userRepo.FindAllPaged(ctx, page1.Pageable.Next())

// With sorting
pageable := PageRequest(0, 20, 
    Order{Property: "created_at", Direction: Desc},
    Order{Property: "name", Direction: Asc},
)
page := userRepo.FindAllPaged(ctx, pageable)

// Using Sort builder
sort := SortBy("name", Asc).
    And(Order{Property: "age", Direction: Desc})
page := userRepo.FindAllPaged(ctx, PageRequest(0, 20, sort.Orders...))

// With null handling
order := Order{
    Property: "last_login",
    Direction: Desc,
    NullHandling: NullsLast, // Users who never logged in appear last
}
page := userRepo.FindAllPaged(ctx, PageRequest(0, 20, order))

// Pagination with specification
spec := And(
    Where(User.Status.EQ(String("active"))),
    Where(User.Age.GT(Int(18))),
)
page := userRepo.FindAllPaged(ctx, spec, PageRequest(0, 20))

// Iterate through all pages
pageable := PageRequest(0, 50)
for {
    page, err := userRepo.FindAllPaged(ctx, pageable)
    if err != nil {
        return err
    }
    
    // Process page.Content
    for _, user := range page.Content {
        fmt.Println(user.Name)
    }
    
    if page.Last {
        break
    }
    
    pageable = page.Pageable.Next()
}

// Unpaged (fetch all)
page := userRepo.FindAllPaged(ctx, Unpaged())

// Page metadata
fmt.Printf("Page %d of %d\n", page.Number+1, page.TotalPages)
fmt.Printf("Showing %d-%d of %d users\n", 
    page.Number*page.Size+1,
    page.Number*page.Size+page.NumberOfElements,
    page.TotalElements,
)

// Dynamic sort from request
func buildSort(sortFields []string, directions []string) Sort {
    orders := make([]Order, 0, len(sortFields))
    for i, field := range sortFields {
        dir := Asc
        if i < len(directions) && directions[i] == "desc" {
            dir = Desc
        }
        orders = append(orders, Order{Property: field, Direction: dir})
    }
    return Sort{Orders: orders}
}

// Cursor-based pagination (for infinite scroll)
type CursorPageable struct {
    Cursor string
    Size   int
    Sort   Sort
}

type CursorPage[T any] struct {
    Content    []*T
    NextCursor string
    HasNext    bool
}

// Window functions for advanced pagination
type WindowPageable struct {
    Pageable Pageable
    Window   WindowSpec
}

type WindowSpec struct {
    PartitionBy []string
    OrderBy     []Order
}
```

**Pagination Helpers:**

```go
// Slice for in-memory pagination
type Slice[T any] struct {
    Content []T
    Pageable Pageable
}

func (s Slice[T]) ToPage() Page[T] {
    // Convert to Page structure
}

// Page mapper (transform page content)
func MapPage[T any, R any](page Page[T], mapper func(*T) *R) Page[R] {
    content := make([]*R, len(page.Content))
    for i, item := range page.Content {
        content[i] = mapper(item)
    }
    
    return Page[R]{
        Content:          content,
        Pageable:         page.Pageable,
        TotalElements:    page.TotalElements,
        TotalPages:       page.TotalPages,
        Size:             page.Size,
        Number:           page.Number,
        NumberOfElements: page.NumberOfElements,
        First:            page.First,
        Last:             page.Last,
        Empty:            page.Empty,
        Sort:             page.Sort,
    }
}

// Example: Map User entity to UserDTO
userPage := userRepo.FindAllPaged(ctx, pageable)
dtoPage := MapPage(userPage, func(u *User) *UserDTO {
    return &UserDTO{
        ID:    u.ID,
        Email: u.Email,
        Name:  u.FullName,
    }
})
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

// Connection string (alternative)
db, err := jetorm.ConnectURL("postgres://user:pass@localhost:5432/myapp?sslmode=disable")

// Multiple database connections
primaryDB := jetorm.MustConnect(primaryConfig)
replicaDB := jetorm.MustConnect(replicaConfig)

// Configure read/write splitting
db.SetReadDB(replicaDB)
db.SetWriteDB(primaryDB)
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
    
    // CREATE
    user := &User{
        Email:    "john@example.com",
        Username: "john_doe",
        Age:      25,
        Status:   "active",
    }
    created, err := userRepo.Save(ctx, user)
    
    // READ - Find by ID
    found, err := userRepo.FindByID(ctx, created.ID)
    
    // READ - Find by custom method
    byEmail, err := userRepo.FindByEmail(ctx, "john@example.com")
    
    // READ - Find with specification
    spec := jetorm.And(
        jetorm.Where(User.Age.GT(jetorm.Int(18))),
        jetorm.Where(User.Status.EQ(jetorm.String("active"))),
    )
    adults, err := userRepo.FindAll(ctx, spec)
    
    // READ - Pagination
    page, err := userRepo.FindAllPaged(ctx, jetorm.PageRequest(0, 20,
        jetorm.Order{Property: "created_at", Direction: jetorm.Desc},
    ))
    
    // READ - Complex queries
    users, err := userRepo.FindByAgeGreaterThanAndStatusInOrderByCreatedAtDesc(
        ctx, 18, []string{"active", "premium"},
    )
    
    // UPDATE
    user.Status = "inactive"
    updated, err := userRepo.Update(ctx, user)
    
    // UPDATE - Batch
    users := []*User{user1, user2, user3}
    updated, err := userRepo.UpdateAll(ctx, users)
    
    // DELETE
    err = userRepo.Delete(ctx, user)
    
    // DELETE - By ID
    err = userRepo.DeleteByID(ctx, user.ID)
    
    // DELETE - By custom method
    deleted, err := userRepo.DeleteByEmail(ctx, "old@example.com")
    
    // COUNT
    count, err := userRepo.Count(ctx)
    
    // COUNT - With condition
    activeCount, err := userRepo.CountByStatus(ctx, "active")
    
    // EXISTS
    exists, err := userRepo.ExistsByEmail(ctx, "test@example.com")
    
    // TRANSACTION
    err = db.Transaction(ctx, func(tx *jetorm.Tx) error {
        user.Age = 26
        _, err := userRepo.WithTx(tx).Save(ctx, user)
        if err != nil {
            return err
        }
        
        order := &Order{UserID: user.ID, Total: 100.0}
        _, err = orderRepo.WithTx(tx).Save(ctx, order)
        if err != nil {
            return err
        }
        
        return nil
    })
    
    // RAW QUERY
    var users []*User
    err = userRepo.Query(ctx, 
        "SELECT * FROM users WHERE age > $1 AND status = $2",
        18, "active",
    )
}

// Advanced usage examples

func exampleDynamicQuery(filters UserFilters) {
    specs := UserSpecifications{}
    conditions := []jetorm.Specification[User]{}
    
    if filters.MinAge > 0 {
        conditions = append(conditions, specs.AgeGreaterThan(filters.MinAge))
    }
    if filters.Status != "" {
        conditions = append(conditions, specs.HasStatus(filters.Status))
    }
    if filters.SearchName != "" {
        conditions = append(conditions, specs.SearchByName(filters.SearchName))
    }
    
    var finalSpec jetorm.Specification[User]
    if len(conditions) > 0 {
        finalSpec = jetorm.And(conditions...)
    }
    
    page := userRepo.FindAllPaged(ctx, finalSpec, 
        jetorm.PageRequest(filters.Page, filters.Size),
    )
}

func exampleBulkOperations() {
    // Batch insert
    users := make([]*User, 1000)
    for i := range users {
        users[i] = &User{
            Email: fmt.Sprintf("user%d@example.com", i),
            Username: fmt.Sprintf("user%d", i),
        }
    }
    
    // Save in batches of 100
    err := userRepo.SaveBatch(ctx, users, 100)
    
    // Bulk delete
    spec := jetorm.Where(User.CreatedAt.LT(jetorm.Timestamp(cutoffDate)))
    deleted, err := userRepo.Delete(ctx, spec)
    fmt.Printf("Deleted %d old users\n", deleted)
}

func exampleRelationships() {
    // Load with associations
    user, err := userRepo.FindByID(ctx, 1, 
        jetorm.LoadWith("Orders"),
        jetorm.LoadWith("Profile"),
    )
    
    // Access loaded associations
    for _, order := range user.Orders {
        fmt.Printf("Order: %d - Total: %.2f\n", order.ID, order.Total)
    }
}

func examplePaginationIteration() {
    pageable := jetorm.PageRequest(0, 100)
    
    for {
        page, err := userRepo.FindAllPaged(ctx, pageable)
        if err != nil {
            log.Fatal(err)
        }
        
        // Process page
        for _, user := range page.Content {
            processUser(user)
        }
        
        if page.Last {
            break
        }
        
        pageable = page.Pageable.Next()
    }
}

func exampleComplexTransaction() {
    err := db.TransactionWithOptions(ctx, jetorm.TxOptions{
        Isolation: jetorm.LevelSerializable,
        Timeout: 30 * time.Second,
    }, func(tx *jetorm.Tx) error {
        // Create savepoint
        tx.SavePoint("sp1")
        
        // First operation
        user, err := userRepo.WithTx(tx).Save(ctx, user)
        if err != nil {
            return err
        }
        
        // Try risky operation
        err = riskyOperation(tx, user)
        if err != nil {
            // Rollback to savepoint, continue transaction
            tx.RollbackTo("sp1")
            return handleRiskyFailure(tx)
        }
        
        tx.ReleaseSavePoint("sp1")
        return nil
    })
}

func exampleCaching() {
    // Implement repository with caching
    type CachedUserRepository struct {
        *UserRepository
        cache Cache
    }
    
    func (r *CachedUserRepository) FindByID(ctx context.Context, id int64) (*User, error) {
        key := fmt.Sprintf("user:%d", id)
        
        // Try cache first
        if cached, ok := r.cache.Get(key); ok {
            return cached.(*User), nil
        }
        
        // Fetch from database
        user, err := r.UserRepository.FindByID(ctx, id)
        if err != nil {
            return nil, err
        }
        
        // Store in cache
        r.cache.Set(key, user, 5*time.Minute)
        return user, nil
    }
}

func exampleTestingWithMocks() {
    // Mock repository for testing
    mockRepo := &MockUserRepository{}
    mockRepo.On("FindByID", ctx, 1).Return(&User{
        ID: 1,
        Email: "test@example.com",
    }, nil)
    
    // Use in service
    service := NewUserService(mockRepo)
    user, err := service.GetUser(ctx, 1)
    
    mockRepo.AssertExpectations(t)
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

