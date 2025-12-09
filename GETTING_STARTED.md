# Getting Started with JetORM

This guide will help you get started with JetORM, a next-generation Go database library.

## Prerequisites

- Go 1.24 or higher
- PostgreSQL 12 or higher
- Basic understanding of Go and SQL

## Installation

```bash
go get github.com/satishbabariya/jetorm
```

## Step 1: Define Your Entity

Create a struct that represents your database table. Use struct tags to define database constraints and behaviors:

```go
package main

import "time"

type User struct {
    ID        int64     `db:"id" jet:"primary_key,auto_increment"`
    Email     string    `db:"email" jet:"unique,not_null"`
    Username  string    `db:"username" jet:"unique,not_null"`
    FullName  string    `db:"full_name"`
    Age       int       `db:"age"`
    Status    string    `db:"status" jet:"default:'active'"`
    IsActive  bool      `db:"is_active" jet:"default:true"`
    CreatedAt time.Time `db:"created_at" jet:"auto_now_add,not_null"`
    UpdatedAt time.Time `db:"updated_at" jet:"auto_now,not_null"`
}
```

### Available Struct Tags

- `db:"column_name"` - Database column name
- `jet:"primary_key"` - Mark as primary key
- `jet:"auto_increment"` - Auto-incrementing integer
- `jet:"unique"` - Unique constraint
- `jet:"not_null"` - NOT NULL constraint
- `jet:"index"` - Create an index
- `jet:"size:255"` - VARCHAR size
- `jet:"default:'value'"` - Default value
- `jet:"auto_now_add"` - Set timestamp on insert
- `jet:"auto_now"` - Update timestamp on save

## Step 2: Create Database Schema

Create the corresponding table in PostgreSQL:

```sql
CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    username VARCHAR(50) UNIQUE NOT NULL,
    full_name VARCHAR(255),
    age INTEGER,
    status VARCHAR(20) DEFAULT 'active',
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_status ON users(status);
```

## Step 3: Connect to Database

```go
package main

import (
    "log"
    
    "github.com/satishbabariya/jetorm/core"
)

func main() {
    // Create configuration
    config := core.Config{
        Host:     "localhost",
        Port:     5432,
        Database: "myapp",
        User:     "postgres",
        Password: "secret",
        SSLMode:  "disable",
        
        // Optional: Enable SQL logging
        LogSQL:   true,
        LogLevel: core.DebugLevel,
    }

    // Connect to database
    db, err := core.Connect(config)
    if err != nil {
        log.Fatalf("Failed to connect: %v", err)
    }
    defer db.Close()

    // Verify connection
    if err := db.Ping(context.Background()); err != nil {
        log.Fatalf("Failed to ping database: %v", err)
    }
    
    log.Println("Connected to database!")
}
```

## Step 4: Create Repository

```go
// Create a repository for the User entity
userRepo, err := core.NewBaseRepository[User, int64](db)
if err != nil {
    log.Fatalf("Failed to create repository: %v", err)
}
```

## Step 5: Perform CRUD Operations

### Create (Insert)

```go
ctx := context.Background()

// Create a new user
newUser := &User{
    Email:    "john.doe@example.com",
    Username: "johndoe",
    FullName: "John Doe",
    Age:      30,
    Status:   "active",
    IsActive: true,
}

savedUser, err := userRepo.Save(ctx, newUser)
if err != nil {
    log.Fatalf("Failed to save user: %v", err)
}

fmt.Printf("Created user with ID: %d\n", savedUser.ID)
```

### Read (Query)

```go
// Find by ID
user, err := userRepo.FindByID(ctx, 1)
if err != nil {
    if err == core.ErrNotFound {
        log.Println("User not found")
    } else {
        log.Fatalf("Failed to find user: %v", err)
    }
}

// Find all users
allUsers, err := userRepo.FindAll(ctx)
if err != nil {
    log.Fatalf("Failed to find users: %v", err)
}

// Find multiple users by IDs
users, err := userRepo.FindAllByIDs(ctx, []int64{1, 2, 3})
if err != nil {
    log.Fatalf("Failed to find users: %v", err)
}

// Check if user exists
exists, err := userRepo.ExistsById(ctx, 1)
if err != nil {
    log.Fatalf("Failed to check existence: %v", err)
}

// Count users
count, err := userRepo.Count(ctx)
if err != nil {
    log.Fatalf("Failed to count users: %v", err)
}
```

### Update

```go
// Fetch user
user, err := userRepo.FindByID(ctx, 1)
if err != nil {
    log.Fatalf("Failed to find user: %v", err)
}

// Modify fields
user.Age = 31
user.Status = "premium"

// Save changes
updatedUser, err := userRepo.Save(ctx, user)
if err != nil {
    log.Fatalf("Failed to update user: %v", err)
}
```

### Delete

```go
// Delete by ID
err := userRepo.DeleteByID(ctx, 1)
if err != nil {
    log.Fatalf("Failed to delete user: %v", err)
}

// Delete entity
err = userRepo.Delete(ctx, user)
if err != nil {
    log.Fatalf("Failed to delete user: %v", err)
}

// Delete multiple entities
err = userRepo.DeleteAll(ctx, []*User{user1, user2})
if err != nil {
    log.Fatalf("Failed to delete users: %v", err)
}
```

## Step 6: Pagination

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
if err != nil {
    log.Fatalf("Failed to get page: %v", err)
}

// Access page data
fmt.Printf("Page %d of %d\n", page.Number+1, page.TotalPages)
fmt.Printf("Total elements: %d\n", page.TotalElements)
fmt.Printf("Is first page: %v\n", page.First)
fmt.Printf("Is last page: %v\n", page.Last)

for _, user := range page.Content {
    fmt.Printf("- %s (%s)\n", user.Username, user.Email)
}
```

## Step 7: Transactions

```go
// Simple transaction
err := db.Transaction(ctx, func(tx *core.Tx) error {
    // Get transaction-aware repository
    txRepo := userRepo.WithTx(tx)
    
    // Create first user
    user1 := &User{
        Email:    "alice@example.com",
        Username: "alice",
        Age:      25,
    }
    _, err := txRepo.Save(ctx, user1)
    if err != nil {
        return err // Automatic rollback
    }
    
    // Create second user
    user2 := &User{
        Email:    "bob@example.com",
        Username: "bob",
        Age:      28,
    }
    _, err = txRepo.Save(ctx, user2)
    if err != nil {
        return err // Automatic rollback
    }
    
    return nil // Automatic commit
})

if err != nil {
    log.Fatalf("Transaction failed: %v", err)
}

// Transaction with options
err = db.TransactionWithOptions(ctx, core.TxOptions{
    Isolation:  core.Serializable,
    ReadOnly:   false,
    Deferrable: false,
}, func(tx *core.Tx) error {
    // Your transaction logic here
    return nil
})
```

## Step 8: Advanced Configuration

```go
config := core.Config{
    // Connection
    Host:     "localhost",
    Port:     5432,
    Database: "myapp",
    User:     "postgres",
    Password: "secret",
    SSLMode:  "require", // disable, require, verify-ca, verify-full
    
    // Connection Pool
    MaxOpenConns:    25,                // Maximum open connections
    MaxIdleConns:    5,                 // Maximum idle connections
    ConnMaxLifetime: 5 * time.Minute,   // Connection max lifetime
    ConnMaxIdleTime: 5 * time.Minute,   // Connection max idle time
    
    // Logging
    LogSQL:          true,                      // Log all SQL queries
    LogLevel:        core.DebugLevel,           // Log level
    LogSlowQueries:  100 * time.Millisecond,    // Log slow queries
    
    // Performance
    PreparedStmts:   true,              // Use prepared statements
    QueryTimeout:    30 * time.Second,  // Default query timeout
    
    // Behavior
    SoftDelete:      false,             // Enable soft delete globally
    CreatedAtField:  "created_at",      // Custom created_at field name
    UpdatedAtField:  "updated_at",      // Custom updated_at field name
    DeletedAtField:  "deleted_at",      // Custom deleted_at field name
}
```

## Complete Example

See the [examples/basic](examples/basic) directory for a complete working example.

```bash
cd examples/basic
go run main.go
```

## Next Steps

1. Explore the [API Documentation](README.md)
2. Check out [Advanced Examples](examples/advanced)
3. Learn about [Code Generation](generator/README.md) (coming soon)
4. Read about [Migrations](migration/README.md) (coming soon)

## Common Patterns

### Repository as a Service Dependency

```go
type UserService struct {
    userRepo core.Repository[User, int64]
}

func NewUserService(db *core.Database) (*UserService, error) {
    repo, err := core.NewBaseRepository[User, int64](db)
    if err != nil {
        return nil, err
    }
    return &UserService{userRepo: repo}, nil
}

func (s *UserService) GetActiveUsers(ctx context.Context) ([]*User, error) {
    // Custom business logic
    allUsers, err := s.userRepo.FindAll(ctx)
    if err != nil {
        return nil, err
    }
    
    // Filter active users
    activeUsers := make([]*User, 0)
    for _, user := range allUsers {
        if user.IsActive {
            activeUsers = append(activeUsers, user)
        }
    }
    
    return activeUsers, nil
}
```

### Error Handling

```go
user, err := userRepo.FindByID(ctx, id)
if err != nil {
    switch err {
    case core.ErrNotFound:
        return nil, fmt.Errorf("user not found")
    case core.ErrInvalidID:
        return nil, fmt.Errorf("invalid user ID")
    default:
        return nil, fmt.Errorf("database error: %w", err)
    }
}
```

## Troubleshooting

### Connection Issues

If you're having trouble connecting:

1. Verify PostgreSQL is running: `pg_isready`
2. Check credentials and database name
3. Verify network connectivity and firewall rules
4. Check PostgreSQL logs for connection errors

### Performance Issues

1. Enable connection pooling (default)
2. Use prepared statements (default)
3. Add appropriate indexes to your tables
4. Monitor slow queries with `LogSlowQueries`

## Getting Help

- Check the [README](README.md) for API documentation
- Look at [examples](examples/) for code samples
- Open an issue on GitHub for bugs or questions

---

Happy coding with JetORM! 🚀

