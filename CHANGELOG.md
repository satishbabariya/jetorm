# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added (Phase 1 - Foundation) ✅

- Core abstractions and interfaces
  - `Repository[T, ID]` generic interface
  - `Entity` metadata extraction
  - `Config` with comprehensive database configuration
  - `Pageable`, `Page`, `Sort`, and `Order` types
  - Error types (`ErrNotFound`, `ErrInvalidEntity`, etc.)

- Database connection management
  - PostgreSQL connection via pgx v5
  - Connection pooling with configurable parameters
  - Connection health checks (Ping)
  - SSL mode support

- Base repository implementation
  - Generic `BaseRepository[T, ID]` with full CRUD operations
  - `Save()` - Insert or update based on primary key
  - `SaveAll()` - Batch save operations
  - `FindByID()` - Find single entity by ID
  - `FindAll()` - Find all entities
  - `FindAllByIDs()` - Find multiple entities by IDs
  - `Delete()` / `DeleteByID()` / `DeleteAll()` - Delete operations
  - `Count()` - Count total entities
  - `ExistsById()` - Check entity existence
  - `FindAllPaged()` - Pagination with sorting support

- Transaction support
  - `Transaction()` - Execute function in transaction
  - `TransactionWithOptions()` - Transaction with isolation levels
  - `Begin()` / `BeginWithOptions()` - Manual transaction control
  - `WithTx()` - Repository bound to transaction
  - Support for all isolation levels (Read Uncommitted, Read Committed, Repeatable Read, Serializable)

- Entity metadata system
  - Struct tag parsing (`db` and `jet` tags)
  - Support for: primary_key, auto_increment, unique, not_null, index, size, default, auto_now_add, auto_now
  - Automatic table name generation (snake_case)
  - Field metadata extraction

- Logging
  - Configurable log levels (Debug, Info, Warn, Error)
  - SQL query logging
  - Slow query detection
  - Default logger implementation

- Utilities
  - `toSnakeCase()` - Convert PascalCase to snake_case
  - `toCamelCase()` - Convert snake_case to camelCase
  - `toPascalCase()` - Convert snake_case to PascalCase

- Testing
  - Unit tests for core functionality
  - Entity metadata extraction tests
  - Pagination tests
  - Utility function tests
  - Test coverage for all public APIs

- Documentation
  - Comprehensive README with features and examples
  - Getting Started guide
  - Basic example application
  - API documentation in code
  - Project specification document

- Project infrastructure
  - Go module setup (go.mod)
  - Project structure following specification
  - .gitignore configuration
  - MIT License
  - Example applications structure

## [0.1.0] - 2025-12-09

### Added

- Initial project setup
- Phase 1 (Foundation) implementation complete
- Core CRUD operations
- Transaction management
- Pagination and sorting
- Basic example application
- Comprehensive documentation

### Notes

This is the initial development release. The API is subject to change before 1.0.

## Roadmap

### Phase 2: Code Generation (Next)
- [ ] Interface parser for custom repository methods
- [ ] Method name analyzer (FindByX, DeleteByY patterns)
- [ ] Code generation templates
- [ ] CLI tool (jetorm-gen)
- [ ] Integration tests

### Phase 3: Query Building
- [ ] Jet SQL builder integration
- [ ] Specification/Criteria API
- [ ] Dynamic query composition
- [ ] Complex query support
- [ ] Raw SQL escape hatch

### Phase 4: Migrations
- [ ] Migration runner (goose integration)
- [ ] Schema generator from entities
- [ ] Version tracking
- [ ] Migration validation
- [ ] Up/Down migration support

### Phase 5: Advanced Features
- [ ] Lifecycle hooks (Before/After operations)
- [ ] Auditing support (created_at, updated_at)
- [ ] Soft delete support
- [ ] Relationship handling (one-to-many, many-to-many)
- [ ] Eager/lazy loading
- [ ] Caching layer

### Phase 6: Testing & Polish
- [ ] Integration tests with testcontainers
- [ ] Performance benchmarks
- [ ] Documentation improvements
- [ ] More example applications
- [ ] API stabilization for 1.0

---

## Version History

- **v0.1.0** (2025-12-09) - Initial release with Phase 1 complete

