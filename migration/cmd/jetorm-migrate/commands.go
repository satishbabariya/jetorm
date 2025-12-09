package main

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/satishbabariya/jetorm/migration"
)

// Command represents a migration command
type Command struct {
	Name        string
	Description string
	Execute     func(context.Context, *sql.DB, string, []string) error
}

// Available commands
var migrationCommands = []Command{
	{
		Name:        "create",
		Description: "Create a new migration file",
		Execute:     cmdCreate,
	},
	{
		Name:        "up",
		Description: "Apply all pending migrations",
		Execute:     cmdUp,
	},
	{
		Name:        "down",
		Description: "Rollback the last migration",
		Execute:     cmdDown,
	},
	{
		Name:        "down-to",
		Description: "Rollback migrations to a specific version",
		Execute:     cmdDownTo,
	},
	{
		Name:        "status",
		Description: "Show migration status",
		Execute:     cmdStatus,
	},
	{
		Name:        "validate",
		Description: "Validate migrations",
		Execute:     cmdValidate,
	},
}

// cmdCreate creates a new migration
func cmdCreate(ctx context.Context, db *sql.DB, migrationsDir string, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("migration name is required")
	}

	runner := migration.NewRunner(nil, migrationsDir)
	return runner.CreateMigration(args[0])
}

// cmdUp applies migrations
func cmdUp(ctx context.Context, db *sql.DB, migrationsDir string, args []string) error {
	runner := migration.NewRunner(db, migrationsDir)
	return runner.Up(ctx)
}

// cmdDown rolls back last migration
func cmdDown(ctx context.Context, db *sql.DB, migrationsDir string, args []string) error {
	runner := migration.NewRunner(db, migrationsDir)
	return runner.Down(ctx)
}

// cmdDownTo rolls back to version
func cmdDownTo(ctx context.Context, db *sql.DB, migrationsDir string, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("target version is required")
	}

	var version int64
	if _, err := fmt.Sscanf(args[0], "%d", &version); err != nil {
		return fmt.Errorf("invalid version: %w", err)
	}

	runner := migration.NewRunner(db, migrationsDir)
	return runner.DownTo(ctx, version)
}

// cmdStatus shows migration status
func cmdStatus(ctx context.Context, db *sql.DB, migrationsDir string, args []string) error {
	runner := migration.NewRunner(db, migrationsDir)
	statuses, err := runner.Status(ctx)
	if err != nil {
		return err
	}

	fmt.Println("Migration Status:")
	fmt.Println("=================")
	for _, status := range statuses {
		statusStr := status.Status
		if status.AppliedAt != nil {
			statusStr += fmt.Sprintf(" (%s)", status.AppliedAt.Format("2006-01-02 15:04:05"))
		}
		fmt.Printf("%d - %s: %s\n", status.Version, status.Name, statusStr)
	}

	return nil
}

// cmdValidate validates migrations
func cmdValidate(ctx context.Context, db *sql.DB, migrationsDir string, args []string) error {
	runner := migration.NewRunner(db, migrationsDir)
	return runner.ValidateMigrations(ctx)
}

// printUsage prints migration command usage
func printMigrationUsage() {
	fmt.Println("Usage: jetorm-migrate [command] [options]")
	fmt.Println("\nCommands:")
	for _, cmd := range migrationCommands {
		fmt.Printf("  %-15s %s\n", cmd.Name, cmd.Description)
	}
	fmt.Println("\nOptions:")
	fmt.Println("  -db string        Database connection string")
	fmt.Println("  -dir string       Migrations directory (default: ./migrations)")
	fmt.Println("  -to int64         Target version for down-to command")
	fmt.Println("  -name string      Migration name for create command")
}

// executeMigrationCommand executes a migration command
func executeMigrationCommand(name string, ctx context.Context, db *sql.DB, migrationsDir string, args []string) error {
	for _, cmd := range migrationCommands {
		if cmd.Name == name {
			return cmd.Execute(ctx, db, migrationsDir, args)
		}
	}
	return fmt.Errorf("unknown command: %s", name)
}

