# JetORM Spec V2 - Key Enhancements Analysis

This document highlights the major enhancements and additions in `jetorm_spec_v2.md` compared to the original specification and current implementation.

## 📊 Overview

The V2 specification is significantly more comprehensive, with **1,663 lines** (vs 791 in v1) and includes many production-ready features and advanced patterns.

## 🔑 Key Enhancements

### 1. Enhanced Repository Interface

**Current Implementation:**
- Basic CRUD operations
- Pagination support
- Transaction support

**V2 Spec Adds:**
- ✅ `Update()` and `UpdateAll()` - Explicit update methods (separate from Save)
- ✅ `DeleteAllByIDs()` - Batch delete by IDs
- ✅ **Specification API** - `FindOne()`, `FindAll()`, `Count()`, `Exists()`, `Delete()` with specifications
- ✅ **Batch Operations** - `SaveBatch()` with configurable batch size
- ✅ **Raw Query Support** - `Query()`, `QueryOne()`, `Exec()` for raw SQL

```go
// V2 adds these methods:
Update(ctx context.Context, entity *T) (*T, error)
UpdateAll(ctx context.Context, entities []*T) ([]*T, error)
DeleteAllByIDs(ctx context.Context, ids []ID) error

// Specification-based queries
FindOne(ctx context.Context, spec Specification[T]) (*T, error)
FindAll(ctx context.Context, spec Specification[T]) ([]*T, error)
FindAllPaged(ctx context.Context, spec Specification[T], pageable Pageable) (*Page[T], error)
Count(ctx context.Context, spec Specification[T]) (int64, error)
Exists(ctx context.Context, spec Specification[T]) (bool, error)
Delete(ctx context.Context, spec Specification[T]) (int64, error)

// Batch operations
SaveBatch(ctx context.Context, entities []*T, batchSize int) error

// Raw SQL
Query(ctx context.Context, query string, args ...interface{}) ([]*T, error)
QueryOne(ctx context.Context, query string, args ...interface{}) (*T, error)
Exec(ctx context.Context, query string, args ...interface{}) (int64, error)
```

### 2. Expanded Query Method Patterns

**V1 Spec:** ~10 basic patterns  
**V2 Spec:** ~30+ patterns with comprehensive coverage

**New Patterns in V2:**
- ✅ `GreaterThanEqual`, `LessThanEqual` - Range comparisons
- ✅ `Between` - Range queries
- ✅ `Containing`, `StartingWith`, `EndingWith` - String operations
- ✅ `NotLike`, `IgnoreCase` - Advanced string matching
- ✅ `NotIn` - Exclusion queries
- ✅ `IsNotNull` - Null checks
- ✅ `True`/`False` - Boolean shortcuts
- ✅ `FindFirstBy`, `FindTop{N}By` - Limiting results
- ✅ `CountDistinctBy` - Distinct counting
- ✅ `FindDistinctBy` - Distinct results
- ✅ Complex combinations with multiple `And`/`Or` conditions

**Example from V2:**
```go
// Complex query method
FindByAgeGreaterThanAndStatusInOrderByCreatedAtDesc(ctx, 18, []string{"active", "pending"})
// WHERE age > $1 AND status IN ($2, $3) ORDER BY created_at DESC
```

### 3. Enhanced Entity Tags

**V1 Spec:** Basic tags (primary_key, unique, not_null, index, size, default, auto_now_add, auto_now)  
**V2 Spec:** Comprehensive tag system

**New Tags in V2:**
- ✅ `type:text`, `type:decimal(p,s)`, `type:jsonb` - Explicit type specification
- ✅ `check:expr` - Check constraints
- ✅ `composite_index:name:order` - Multi-column indexes
- ✅ `unique_index` - Unique indexes
- ✅ `foreign_key:table.col` - Foreign key relationships
- ✅ `on_delete:action`, `on_update:action` - FK cascade actions
- ✅ `unique_composite` - Multi-column unique constraints
- ✅ `where:condition` - Partial indexes (PostgreSQL)
- ✅ `generated:expr,stored` - Generated/computed columns
- ✅ Field ignoring with `db:"-" jet:"-"`

**Example from V2:**
```go
type User struct {
    ID        int64      `db:"id" jet:"primary_key,auto_increment"`
    Email     string     `db:"email" jet:"unique,not_null,index:idx_email,size:255"`
    CompanyID *int64     `db:"company_id" jet:"foreign_key:companies.id,on_delete:cascade"`
    Balance   float64    `db:"balance" jet:"type:decimal(10,2),default:0.00"`
    Metadata  types.JSON `db:"metadata" jet:"type:jsonb"`
    Total     float64    `db:"total" jet:"generated:subtotal * (1 + tax_rate),stored"`
}
```

### 4. Advanced Transaction Management

**V1 Spec:** Basic transaction support  
**V2 Spec:** Production-grade transaction features

**New Features:**
- ✅ **Savepoints** - `SavePoint()`, `RollbackTo()`, `ReleaseSavePoint()`
- ✅ **Transaction Timeout** - Configurable timeout per transaction
- ✅ **Read-only Transactions** - Optimized for read operations
- ✅ **Transaction Middleware** - HTTP middleware pattern
- ✅ **Deadlock Retry Logic** - Automatic retry on deadlocks
- ✅ **Nested Transaction Support** - Transaction propagation

**Example from V2:**
```go
// Savepoints for partial rollback
err := db.Transaction(ctx, func(tx *Tx) error {
    tx.SavePoint("sp1")
    // ... operations ...
    if err != nil {
        tx.RollbackTo("sp1") // Rollback to savepoint, continue transaction
    }
    tx.ReleaseSavePoint("sp1")
    return nil
})

// Transaction with timeout
err := db.TransactionWithOptions(ctx, TxOptions{
    Isolation: LevelSerializable,
    Timeout:   30 * time.Second,
}, func(tx *Tx) error {
    // Operations
})
```

### 5. Specification/Criteria API

**V1 Spec:** Mentioned but not detailed  
**V2 Spec:** Complete implementation with examples

**Features:**
- ✅ Composable specifications
- ✅ Reusable specification factories
- ✅ Fluent API for building queries
- ✅ Integration with Jet expressions
- ✅ Dynamic query building

**Example from V2:**
```go
type UserSpecifications struct{}

func (s UserSpecifications) IsActive() Specification[User] {
    return Where(User.IsActive.EQ(Bool(true)))
}

func (s UserSpecifications) AgeGreaterThan(age int) Specification[User] {
    return Where(User.Age.GT(Int(age)))
}

// Usage
spec := And(
    specs.IsActive(),
    specs.AgeGreaterThan(18),
    specs.StatusIn("premium", "gold"),
)
users := userRepo.FindAll(ctx, spec)
```

### 6. Enhanced Pagination

**V1 Spec:** Basic pagination  
**V2 Spec:** Advanced pagination features

**New Features:**
- ✅ `Pageable` interface with `Next()`, `Previous()`, `First()` methods
- ✅ `NullHandling` - NULLS FIRST/LAST support
- ✅ `OffsetRequest` - Offset-based pagination
- ✅ `Slice[T]` - In-memory pagination
- ✅ `MapPage()` - Transform page content
- ✅ Cursor-based pagination support
- ✅ Window functions for advanced pagination

**Example from V2:**
```go
// Null handling
order := Order{
    Property: "last_login",
    Direction: Desc,
    NullHandling: NullsLast, // Users who never logged in appear last
}

// Page iteration
pageable := PageRequest(0, 50)
for {
    page, _ := userRepo.FindAllPaged(ctx, pageable)
    // Process page
    if page.Last {
        break
    }
    pageable = page.Pageable.Next()
}
```

### 7. Additional Features

**Connection Management:**
- ✅ Connection string support (`ConnectURL()`)
- ✅ Multiple database connections (read/write splitting)
- ✅ Read replica support

**Entity Features:**
- ✅ JSON field support (`types.JSON`, `type:jsonb`)
- ✅ Array field support (`type:text[]`)
- ✅ Computed/generated columns
- ✅ Partial indexes

**Advanced Patterns:**
- ✅ Caching layer examples
- ✅ Testing with mocks
- ✅ Bulk operations
- ✅ Relationship loading (eager/lazy)

## 📋 Implementation Priority

### High Priority (Core Functionality)
1. ✅ **Update/UpdateAll methods** - Separate from Save
2. ✅ **Specification API** - Essential for complex queries
3. ✅ **Enhanced entity tags** - Foreign keys, composite indexes
4. ✅ **Savepoints** - Important for complex transactions

### Medium Priority (Enhanced Features)
5. ✅ **Raw query support** - Escape hatch for complex queries
6. ✅ **Batch operations** - Performance optimization
7. ✅ **Enhanced pagination** - Null handling, page iteration
8. ✅ **Expanded query patterns** - More method name patterns

### Low Priority (Nice to Have)
9. ✅ **Connection string support** - Convenience feature
10. ✅ **Read/write splitting** - Advanced deployment pattern
11. ✅ **Cursor-based pagination** - Alternative pagination strategy

## 🔄 Migration Path

### Phase 1 → V2 Compatibility

**Current Implementation Status:**
- ✅ Basic CRUD - Complete
- ✅ Pagination - Basic version complete
- ✅ Transactions - Basic version complete
- ⚠️ Specification API - Not implemented
- ⚠️ Enhanced query patterns - Not implemented
- ⚠️ Advanced entity tags - Partial (basic tags only)

### Recommended Implementation Order

1. **Add Missing Repository Methods** (Week 1)
   - `Update()` / `UpdateAll()`
   - `DeleteAllByIDs()`
   - `SaveBatch()`
   - Raw query methods

2. **Implement Specification API** (Week 2-3)
   - Core specification interface
   - Jet integration
   - Builder functions

3. **Enhance Entity Tags** (Week 3-4)
   - Foreign keys
   - Composite indexes
   - Check constraints
   - Generated columns

4. **Advanced Transactions** (Week 4-5)
   - Savepoints
   - Timeout support
   - Read-only transactions

5. **Query Method Generation** (Week 5-6)
   - Expanded pattern support
   - Code generation for all patterns

## 📝 Notes

### Breaking Changes
The V2 spec introduces some changes that may require updates:
- `FindAll()` method signature changes when using specifications
- `Delete()` returns `int64` (rows affected) when using specifications
- `Page` structure may need additional fields

### Backward Compatibility
Most V2 features are additive and won't break existing code:
- New methods are additions, not replacements
- Existing methods maintain same signatures
- Specification API is optional

## 🎯 Conclusion

The V2 specification represents a **production-ready, enterprise-grade** ORM with:
- Comprehensive query capabilities
- Advanced transaction management
- Rich entity modeling
- Flexible pagination
- Type-safe query building

**Recommendation:** Use V2 spec as the target for full implementation, implementing features incrementally while maintaining backward compatibility with Phase 1 code.

