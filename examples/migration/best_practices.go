package main

import (
	"context"
	"fmt"

	"github.com/satishbabariya/jetorm/migration"
)

// Migration best practices example

func exampleMigrationBestPractices() {
	fmt.Println("Migration Best Practices")
	fmt.Println("========================")

	// 1. Create migrations with descriptive names
	// jetorm-migrate create -name add_user_email_index

	// 2. Always include both UP and DOWN migrations
	// UP migration: CREATE INDEX idx_user_email ON users(email);
	// DOWN migration: DROP INDEX idx_user_email;

	// 3. Use transactions for data migrations
	// BEGIN;
	// UPDATE users SET status = 'active' WHERE status IS NULL;
	// COMMIT;

	// 4. Test migrations on staging first
	// Always test migrations before applying to production

	// 5. Keep migrations small and focused
	// One migration per logical change

	// 6. Never modify applied migrations
	// Create new migrations instead

	// 7. Use schema generator for initial schema
	// generator := migration.NewSchemaGenerator()
	// sql := generator.GenerateCreateTable(User{})

	// 8. Validate migrations before applying
	// runner := migration.NewRunner(db, "./migrations")
	// err := runner.ValidateMigrations(ctx)

	// 9. Use version control for migrations
	// All migration files should be in version control

	// 10. Document complex migrations
	// Add comments explaining complex logic
}

func exampleMigrationWorkflow() {
	fmt.Println("Migration Workflow")
	fmt.Println("==================")

	// 1. Create migration
	// runner := migration.NewRunner(nil, "./migrations")
	// err := runner.CreateMigration("add_user_table")

	// 2. Write migration SQL
	// Edit the generated migration files

	// 3. Validate migration
	// runner := migration.NewRunner(db, "./migrations")
	// err := runner.ValidateMigrations(ctx)

	// 4. Apply migration
	// err := runner.Up(ctx)

	// 5. Check status
	// statuses, err := runner.Status(ctx)
	// for _, status := range statuses {
	//     fmt.Printf("%d - %s: %s\n", status.Version, status.Name, status.Status)
	// }

	// 6. Rollback if needed
	// err := runner.Down(ctx)
}

// Uncomment to run example:
// func main() {
// 	exampleMigrationBestPractices()
// 	exampleMigrationWorkflow()
// }

