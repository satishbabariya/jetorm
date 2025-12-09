# Code Generation Example

This example demonstrates how to use JetORM's code generation feature to automatically generate repository implementations from interface definitions.

## Setup

1. Define your entity struct with appropriate tags
2. Define a repository interface with query methods following the naming convention
3. Use `go:generate` directive to trigger code generation
4. Run `go generate` to generate the repository implementation

## Example Usage

```bash
# Generate repository code
go generate ./examples/codegen

# Or run jetorm-gen directly
jetorm-gen -type=User -interface=UserRepository -input=user.go -output=user_repository_gen.go
```

## Query Method Naming Convention

JetORM supports various query method patterns:

- `FindBy{Field}` - Simple equality query
- `FindBy{Field}And{Field}` - Multiple conditions with AND
- `FindBy{Field}Or{Field}` - Multiple conditions with OR
- `FindBy{Field}GreaterThan` - Comparison operators
- `FindBy{Field}In` - Collection operations
- `FindBy{Field}OrderBy{Field}Desc` - Sorting
- `CountBy{Field}` - Count operations
- `DeleteBy{Field}` - Delete operations

See `jetorm_spec_v2.md` for complete documentation of supported patterns.

