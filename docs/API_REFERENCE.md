# JetORM API Reference

Complete API reference for JetORM.

## Core Package

### Repository Interface

```go
type Repository[T any, ID comparable] interface {
    // Basic CRUD
    Save(ctx context.Context, entity *T) (*T, error)
    FindByID(ctx context.Context, id ID) (*T, error)
    FindAll(ctx context.Context) ([]*T, error)
    Update(ctx context.Context, entity *T) (*T, error)
    Delete(ctx context.Context, entity *T) error
    DeleteByID(ctx context.Context, id ID) error
    
    // Batch operations
    SaveAll(ctx context.Context, entities []*T) ([]*T, error)
    FindAllByIDs(ctx context.Context, ids []ID) ([]*T, error)
    DeleteAll(ctx context.Context, entities []*T) error
    SaveBatch(ctx context.Context, entities []*T, batchSize int) error
    
    // Query operations
    FindOne(ctx context.Context, spec Specification[T]) (*T, error)
    FindAllWithSpec(ctx context.Context, spec Specification[T]) ([]*T, error)
    CountWithSpec(ctx context.Context, spec Specification[T]) (int64, error)
    ExistsWithSpec(ctx context.Context, spec Specification[T]) (bool, error)
    DeleteWithSpec(ctx context.Context, spec Specification[T]) (int64, error)
    
    // Pagination
    FindAllPaged(ctx context.Context, pageable Pageable) (*Page[T], error)
    FindAllPagedWithSpec(ctx context.Context, spec Specification[T], pageable Pageable) (*Page[T], error)
    
    // Raw queries
    Query(ctx context.Context, query string, args ...interface{}) ([]*T, error)
    QueryOne(ctx context.Context, query string, args ...interface{}) (*T, error)
    Exec(ctx context.Context, query string, args ...interface{}) (int64, error)
    
    // Transaction support
    WithTx(tx *Tx) Repository[T, ID]
}
```

### Database Connection

```go
// Connect to database
func Connect(config Config) (*Database, error)

// ConnectURL connects using connection string
func ConnectURL(url string) (*Database, error)

// MustConnect panics on error
func MustConnect(config Config) *Database
```

### Specification API

```go
// Create specifications
func Equal[T any](field string, value interface{}) Specification[T]
func NotEqual[T any](field string, value interface{}) Specification[T]
func GreaterThan[T any](field string, value interface{}) Specification[T]
func LessThan[T any](field string, value interface{}) Specification[T]
func Like[T any](field string, pattern string) Specification[T]
func In[T any](field string, values ...interface{}) Specification[T]
func NotIn[T any](field string, values ...interface{}) Specification[T]
func IsNull[T any](field string) Specification[T]
func IsNotNull[T any](field string) Specification[T]

// Combine specifications
func And[T any](specs ...Specification[T]) Specification[T]
func Or[T any](specs ...Specification[T]) Specification[T]
func Not[T any](spec Specification[T]) Specification[T]
```

### Pagination

```go
// Create pageable
func PageRequest(page, size int, orders ...Order) Pageable

// Create order
func OrderBy(property string, direction Direction) Order

// Page result
type Page[T any] struct {
    Content       []*T
    Pageable      Pageable
    TotalElements int64
    TotalPages    int
    First         bool
    Last          bool
}
```

### Validation

```go
// Create validator
func NewValidator() *Validator

// Register rules
validator.RegisterRule("Email", Email())
validator.RegisterRule("Age", Range(18, 120))

// Validate entity
err := validator.Validate(entity)

// Available rules
Required()
Email()
URL()
MinLength(min int)
MaxLength(max int)
Range(min, max float64)
Pattern(pattern string)
In(allowed ...interface{})
Positive()
Negative()
```

### Caching

```go
// Create cache
cache := NewInMemoryCache()

// Create cached repository
cachedRepo := NewCachedRepository(repo, cache, "User", 5*time.Minute)
```

### Health Monitoring

```go
// Create health checker
checker := NewHealthChecker(db)

// Check health
health := checker.Check(ctx)

// Get metrics
metrics := checker.GetMetrics()
```

### Performance Monitoring

```go
// Create monitor
monitor := NewPerformanceMonitor(100 * time.Millisecond)

// Create profiler
profiler := NewQueryProfiler(monitor)

// Profile query
profiler.Profile(ctx, "FindByID", func(ctx context.Context) error {
    return repo.FindByID(ctx, id)
})
```

### Helper Functions

```go
// Find or create
FindOrCreate[T, ID](ctx, repo, finder, creator)

// Batch operations
OptimizedBatchSave[T, ID](ctx, repo, entities, batchSize)
BatchFind[T, ID](ctx, repo, ids, batchSize)

// Entity helpers
ExtractID[T, ID](entity)
SetID[T, ID](entity, id)
Clone(entity)
Copy(dest, src)
Merge(dest, src)

// Slice utilities
SliceMap[T, U](slice, fn)
SliceFilter[T](slice, fn)
SliceUnique[T](slice)
SliceContains[T](slice, value)

// Retry
Retry(ctx, maxAttempts, fn)
RetryWithBackoff(ctx, maxAttempts, backoff, fn)

// Timeout
Timeout(ctx, timeout, fn)
```

## Query Package

### Query Builder

```go
// Create builder
builder := NewQueryBuilder("users")

// Build query
builder.WhereEqual("status", "active")
       .OrderBy("created_at", Desc)
       .Limit(10)
query, args := builder.Build()
```

### Composable Query

```go
// Create composable query
query := NewComposableQuery[User]("users")

// Add conditions
query.WhereSpecification(spec)
     .OrderBy("name", Asc)
     .Limit(20)
```

## Migration Package

### Migration Runner

```go
// Create runner
runner := migration.NewRunner(db, "./migrations")

// Apply migrations
err := runner.Up(ctx)

// Rollback
err := runner.Down(ctx)

// Status
statuses, err := runner.Status(ctx)
```

## Generator Package

### Code Generation

```go
// Parse interface
parser := generator.NewParser()
methods, err := parser.ParseInterface(file, interfaceName)

// Analyze methods
analyzer := generator.NewAnalyzer()
queryMethods, err := analyzer.AnalyzeMethods(methods)

// Generate code
codegen := generator.NewCodeGenerator()
code, err := codegen.Generate(config, queryMethods)
```

## Hooks Package

### Lifecycle Hooks

```go
// Create hooks
hooks := hooks.NewHooks[User]()

// Register hooks
hooks.RegisterBeforeCreate(func(ctx context.Context, entity *User) error {
    // Before create logic
    return nil
})

hooks.RegisterAfterUpdate(func(ctx context.Context, entity *User) error {
    // After update logic
    return nil
})
```

## Examples

See `examples/` directory for complete usage examples.

