# Advanced Features Examples

This directory contains examples demonstrating JetORM's advanced features from Phase 5.

## Overview

Phase 5 introduces advanced features:
- Relationship handling (one-to-one, one-to-many, many-to-many)
- Caching layer
- Optimized batch operations
- Lifecycle hooks integration

## Examples

### 1. Relationships

```go
type User struct {
    ID      int64    `db:"id" jet:"primary_key"`
    Profile *Profile `db:"-" jet:"one_to_one:Profile,foreign_key:user_id"`
    Posts   []*Post  `db:"-" jet:"one_to_many:Post,mapped_by:user_id"`
    Roles   []*Role  `db:"-" jet:"many_to_many:Role,join_table:user_roles"`
}

// Load relationships
relationships := core.LoadRelationships(reflect.TypeOf(User{}))
```

### 2. Caching

```go
// Create cache
cache := core.NewInMemoryCache()

// Create cached repository
cachedRepo := core.NewCachedRepository(
    repo,
    cache,
    "Product",
    5*time.Minute, // TTL
)

// Use cached repository
product, err := cachedRepo.FindByID(ctx, 1)
```

### 3. Batch Operations

```go
// Batch writer with auto-flush
config := core.DefaultBatchConfig()
config.Size = 50
config.FlushInterval = 2 * time.Second

batchWriter := core.NewBatchWriter(repo, config)
defer batchWriter.Close(ctx)

// Write entities
for _, entity := range entities {
    batchWriter.Write(ctx, entity)
}
```

### 4. Optimized Batch Save

```go
// Save large number of entities efficiently
err := core.OptimizedBatchSave(ctx, repo, entities, 100)
```

## Relationship Types

### One-to-One
```go
Profile *Profile `jet:"one_to_one:Profile,foreign_key:user_id"`
```

### One-to-Many
```go
Posts []*Post `jet:"one_to_many:Post,mapped_by:user_id"`
```

### Many-to-One
```go
User *User `jet:"many_to_one:User,foreign_key:user_id"`
```

### Many-to-Many
```go
Roles []*Role `jet:"many_to_many:Role,join_table:user_roles,join_column:user_id,inverse_join_column:role_id"`
```

## Caching

### In-Memory Cache
```go
cache := core.NewInMemoryCache()
```

### Cache Key Generation
```go
keyGen := core.NewCacheKeyGenerator[Product, int64]("Product")
key := keyGen.KeyForID(1)
```

### Cached Repository
```go
cachedRepo := core.NewCachedRepository(repo, cache, "Product", 5*time.Minute)
```

## Batch Operations

### Batch Writer
- Automatic batching
- Auto-flush on interval
- Configurable batch size
- Timeout support

### Batch Reader
- Cursor-based reading
- Configurable batch size
- Efficient memory usage

### Batch Processor
- Process entities in batches
- Custom processing logic
- Error handling

## Best Practices

1. **Relationships**: Use eager loading for frequently accessed relationships
2. **Caching**: Set appropriate TTL based on data volatility
3. **Batch Operations**: Use appropriate batch sizes (50-1000)
4. **Memory**: Monitor memory usage with large batches

## Next Steps

- See individual example files for complete code
- Check `PACKAGES.md` for package overview
- Read `GETTING_STARTED.md` for basic usage

