# Complete JetORM Example

This example demonstrates a complete workflow using all JetORM features.

## Features Demonstrated

1. **Database Connection** - Connecting to PostgreSQL
2. **Base Repository** - Creating repository instances
3. **Lifecycle Hooks** - Before/After operation hooks
4. **Caching** - In-memory caching for performance
5. **Validation** - Entity validation before operations
6. **Performance Monitoring** - Query performance tracking
7. **Full-Featured Repository** - Combining all features
8. **Health Checks** - Database health monitoring
9. **Batch Operations** - Efficient bulk operations
10. **Helper Functions** - Convenience utilities

## Usage

```bash
# Set up database
createdb jetorm_test

# Run example
go run examples/complete/complete_example.go
```

## Code Structure

```go
// 1. Connect
db := core.Connect(config)

// 2. Create base repository
baseRepo := core.NewBaseRepository[User, int64](db)

// 3. Set up features
hooks := hooks.NewHooks[User]()
cache := core.NewInMemoryCache()
validator := core.NewValidator()
profiler := core.NewQueryProfiler(monitor)

// 4. Create full-featured repository
fullRepo := core.NewFullFeaturedRepository(...)

// 5. Use repository
user, err := fullRepo.Save(ctx, &User{...})
```

## Features

### Health Monitoring
```go
health := fullRepo.HealthCheck(ctx)
fmt.Printf("Status: %s\n", health.Status)
```

### Performance Metrics
```go
metrics := monitor.GetAllMetrics()
```

### Batch Operations
```go
core.OptimizedBatchSave(ctx, repo, entities, 100)
```

### Helper Functions
```go
exists, _ := core.Exists(ctx, repo, id)
```

## Best Practices

1. **Use Full-Featured Repository** - Combines all features
2. **Monitor Performance** - Track query performance
3. **Validate Input** - Always validate entities
4. **Use Caching** - Cache frequently accessed data
5. **Health Checks** - Monitor database health
6. **Batch Operations** - Use batches for bulk operations

## Next Steps

- See individual feature examples
- Check `PACKAGES.md` for package overview
- Read `GETTING_STARTED.md` for basic usage

