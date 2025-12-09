# Migration Examples

This directory contains examples demonstrating JetORM's migration capabilities.

## Overview

Phase 4 introduces comprehensive migration management:
- Migration file management
- Schema generation from entities
- Migration validation
- Rollback support
- CLI tool for migrations

## Migration File Format

Migrations follow the naming convention:
```
YYYYMMDDHHMMSS_name.up.sql
YYYYMMDDHHMMSS_name.down.sql
```

Example:
```
20251209120000_create_users_table.up.sql
20251209120000_create_users_table.down.sql
```

## Examples

### 1. Create Migration Manually

```go
runner := migration.NewRunner(nil, "./migrations")
err := runner.CreateMigration("create_users_table")
```

This creates:
- `YYYYMMDDHHMMSS_create_users_table.up.sql`
- `YYYYMMDDHHMMSS_create_users_table.down.sql`

### 2. Generate Migration from Entity

```go
gen := migration.NewGenerator()
err := gen.GenerateCreateTableMigration(
    reflect.TypeOf(User{}),
    "users",
    "./migrations",
)
```

This automatically generates CREATE TABLE migration from your entity definition.

### 3. Apply Migrations

```go
runner := migration.NewRunner(db, "./migrations")
err := runner.Up(ctx)
```

Applies all pending migrations in order.

### 4. Check Migration Status

```go
runner := migration.NewRunner(db, "./migrations")
statuses, err := runner.Status(ctx)

for _, status := range statuses {
    fmt.Printf("%d - %s: %s\n", status.Version, status.Name, status.Status)
}
```

### 5. Rollback Migrations

```go
// Rollback last migration
runner.Down(ctx)

// Rollback to specific version
runner.DownTo(ctx, targetVersion)
```

### 6. Generate Index Migration

```go
gen := migration.NewGenerator()
err := gen.GenerateIndexMigration(
    "users",
    "idx_email",
    []string{"email"},
    true, // unique
    "./migrations",
)
```

### 7. Generate Foreign Key Migration

```go
gen := migration.NewGenerator()
err := gen.GenerateForeignKeyMigration(
    "users",           // table
    "company_id",      // column
    "companies",       // ref table
    "id",              // ref column
    "cascade",         // on delete
    "set_null",        // on update
    "./migrations",
)
```

## CLI Tool

Use the `jetorm-migrate` CLI tool:

```bash
# Create a new migration
jetorm-migrate -command=create -name=add_user_table -dir=./migrations

# Apply migrations
jetorm-migrate -command=up -db="postgres://..." -dir=./migrations

# Check status
jetorm-migrate -command=status -db="postgres://..." -dir=./migrations

# Rollback last migration
jetorm-migrate -command=down -db="postgres://..." -dir=./migrations

# Rollback to version
jetorm-migrate -command=down-to -db="postgres://..." -dir=./migrations -to=20251209120000

# Validate migrations
jetorm-migrate -command=validate -db="postgres://..." -dir=./migrations
```

## Migration File Structure

### Up Migration
```sql
-- Migration: create_users_table
-- Version: 20251209120000
-- Up migration

CREATE TABLE users (
    id BIGINT PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    username VARCHAR(255) UNIQUE NOT NULL,
    age INTEGER,
    status VARCHAR(50) DEFAULT 'active',
    created_at TIMESTAMP DEFAULT NOW()
);
```

### Down Migration
```sql
-- Migration: create_users_table
-- Version: 20251209120000
-- Down migration

DROP TABLE IF EXISTS users;
```

## Best Practices

1. **Always provide down migrations** - Enables rollback
2. **Use transactions** - Migrations run in transactions automatically
3. **Validate before applying** - Use validation command
4. **Version control** - Commit migration files to version control
5. **Test migrations** - Test both up and down migrations
6. **One change per migration** - Keep migrations focused
7. **Use descriptive names** - Clear migration names help understanding

## Schema Generation

The schema generator automatically:
- Maps Go types to PostgreSQL types
- Handles primary keys
- Applies constraints (unique, not_null)
- Sets default values
- Creates indexes
- Generates foreign keys

## Validation

The validator checks:
- SQL syntax validity
- Migration ordering
- Duplicate versions
- Missing up/down SQL
- Applied migrations consistency

## Integration with JetORM

Migrations work seamlessly with JetORM:
- Schema generation from entities
- Automatic table creation
- Migration tracking
- Rollback support

## Next Steps

1. Create migration files
2. Apply migrations to database
3. Use generated schema in your application
4. Rollback if needed
5. Validate migrations regularly

