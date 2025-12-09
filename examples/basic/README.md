# JetORM Basic Example

This example demonstrates the basic usage of JetORM with CRUD operations.

## Prerequisites

1. PostgreSQL installed and running
2. Go 1.24+ installed

## Setup

1. Create a test database:

```bash
createdb jetorm_test
```

2. Create the users table:

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
```

3. Update the database configuration in `main.go` if needed:

```go
config := core.Config{
    Host:     "localhost",
    Port:     5432,
    Database: "jetorm_test",
    User:     "postgres",
    Password: "postgres",
    SSLMode:  "disable",
}
```

## Run

```bash
cd examples/basic
go run main.go
```

## What This Example Demonstrates

1. **Database Connection**: Connecting to PostgreSQL using JetORM
2. **Create**: Inserting new records
3. **Read**: Finding records by ID and fetching all records
4. **Update**: Modifying existing records
5. **Delete**: Removing records
6. **Count**: Counting total records
7. **Pagination**: Fetching paginated results with sorting
8. **Existence Check**: Checking if a record exists
9. **Transactions**: Executing operations within a transaction

## Expected Output

```
✓ Connected to database

--- Example 1: Create User ---
✓ Created user: ID=1, Email=john.doe@example.com

--- Example 2: Find by ID ---
✓ Found user: johndoe (john.doe@example.com)

--- Example 3: Update User ---
✓ Updated user: Age=31, Status=premium

--- Example 4: Count Users ---
✓ Total users: 1

--- Example 5: Find All Users ---
✓ Found 1 users:
  - johndoe (john.doe@example.com)

--- Example 6: Pagination ---
✓ Page 1 of 1 (Total: 1)
  - johndoe

--- Example 7: Check Existence ---
✓ User exists: true

--- Example 8: Transaction ---
✓ Created user in transaction: janedoe

--- Example 9: Delete User ---
✓ Deleted user with ID: 1

✓ All examples completed!
```

