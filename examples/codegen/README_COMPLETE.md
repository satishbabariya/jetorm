# Complete Code Generation Example

This example demonstrates the complete code generation workflow for JetORM.

## Overview

This example shows:
1. Entity definition with comprehensive tags
2. Repository interface with query methods
3. Code generation using `go generate`
4. Usage of generated repository

## Entity Definition

The `Product` entity demonstrates:
- Primary key with auto-increment
- Unique constraints
- Indexes (single and composite)
- Foreign keys with cascade actions
- Check constraints
- Default values
- Custom column types
- Timestamp fields (created_at, updated_at)

## Repository Interface

The `ProductRepository` interface includes:
- Find operations with various conditions
- Count operations
- Exists operations
- Delete operations
- Custom ordering
- Limit operations (First, TopN)

## Code Generation

### Step 1: Define Entity and Interface

```go
type Product struct {
    // ... fields with tags
}

type ProductRepository interface {
    // ... query methods
}
```

### Step 2: Add go:generate Directive

```go
//go:generate jetorm-gen -type=Product -interface=ProductRepository -input=complete_example.go -output=product_repository_gen.go -comments=true
```

### Step 3: Generate Code

```bash
go generate ./examples/codegen
```

Or manually:
```bash
jetorm-gen -type=Product -interface=ProductRepository \
  -input=complete_example.go \
  -output=product_repository_gen.go \
  -comments=true
```

### Step 4: Use Generated Repository

```go
db := core.MustConnect(config)
repo := NewProductRepository(db)

// Use generated methods
product, err := repo.FindBySKU(ctx, "SKU-123")
products, err := repo.FindByStatusOrderByPriceAsc(ctx, "active")
count, err := repo.CountByStatus(ctx, "active")
```

## Supported Query Patterns

### Find Operations
- `FindBy{Field}` - Simple equality
- `FindBy{Field}And{Field}` - Multiple conditions
- `FindBy{Field}GreaterThan` - Comparison operators
- `FindBy{Field}Between` - Range queries
- `FindBy{Field}In` - Collection queries
- `FindBy{Field}OrderBy{Field}Asc/Desc` - With ordering
- `FindFirstBy{Field}` - Limit 1
- `FindTop{N}By{Field}` - Limit N

### Count Operations
- `CountBy{Field}` - Count with condition

### Exists Operations
- `ExistsBy{Field}` - Check existence

### Delete Operations
- `DeleteBy{Field}` - Delete with condition

## Configuration File

You can also use a configuration file:

```bash
# Create config
jetorm-gen init jetorm-gen.json

# Edit jetorm-gen.json
{
  "entity_type": "Product",
  "interface_name": "ProductRepository",
  "input_file": "complete_example.go",
  "output_file": "product_repository_gen.go",
  "generate_comments": true,
  "generate_tests": false
}

# Generate
jetorm-gen -config=jetorm-gen.json
```

## Next Steps

1. Run code generation
2. Implement custom methods (like `FindLowStockProducts`)
3. Write tests for your repository
4. Use the repository in your application

