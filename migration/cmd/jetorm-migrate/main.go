package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/satishbabariya/jetorm/migration"
)

func main() {
	var (
		command      = flag.String("command", "", "Migration command: up, down, down-to, status, create, validate")
		dbURL        = flag.String("db", "", "Database connection string")
		migrationsDir = flag.String("dir", "./migrations", "Migrations directory")
		targetVersion = flag.Int64("to", 0, "Target version for down-to command")
		migrationName = flag.String("name", "", "Migration name for create command")
	)
	flag.Parse()

	if *command == "" {
		fmt.Fprintf(os.Stderr, "Error: -command is required\n")
		flag.Usage()
		os.Exit(1)
	}

	if *dbURL == "" && *command != "create" && *command != "validate" {
		fmt.Fprintf(os.Stderr, "Error: -db is required for command: %s\n", *command)
		os.Exit(1)
	}

	ctx := context.Background()

	switch *command {
	case "create":
		if *migrationName == "" {
			fmt.Fprintf(os.Stderr, "Error: -name is required for create command\n")
			os.Exit(1)
		}
		
		// Create migration files (no DB needed)
		runner := migration.NewRunner(nil, *migrationsDir)
		if err := runner.CreateMigration(*migrationName); err != nil {
			fmt.Fprintf(os.Stderr, "Error creating migration: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Created migration: %s\n", *migrationName)

	case "up":
		db, err := sql.Open("pgx", *dbURL)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error connecting to database: %v\n", err)
			os.Exit(1)
		}
		defer db.Close()

		runner := migration.NewRunner(db, *migrationsDir)
		if err := runner.Up(ctx); err != nil {
			fmt.Fprintf(os.Stderr, "Error applying migrations: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Migrations applied successfully")

	case "down":
		db, err := sql.Open("pgx", *dbURL)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error connecting to database: %v\n", err)
			os.Exit(1)
		}
		defer db.Close()

		runner := migration.NewRunner(db, *migrationsDir)
		if err := runner.Down(ctx); err != nil {
			fmt.Fprintf(os.Stderr, "Error rolling back migration: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Migration rolled back successfully")

	case "down-to":
		if *targetVersion == 0 {
			fmt.Fprintf(os.Stderr, "Error: -to is required for down-to command\n")
			os.Exit(1)
		}

		db, err := sql.Open("pgx", *dbURL)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error connecting to database: %v\n", err)
			os.Exit(1)
		}
		defer db.Close()

		runner := migration.NewRunner(db, *migrationsDir)
		if err := runner.DownTo(ctx, *targetVersion); err != nil {
			fmt.Fprintf(os.Stderr, "Error rolling back migrations: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Migrations rolled back to version %d\n", *targetVersion)

	case "status":
		db, err := sql.Open("pgx", *dbURL)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error connecting to database: %v\n", err)
			os.Exit(1)
		}
		defer db.Close()

		runner := migration.NewRunner(db, *migrationsDir)
		statuses, err := runner.Status(ctx)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting migration status: %v\n", err)
			os.Exit(1)
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

	case "validate":
		db, err := sql.Open("pgx", *dbURL)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error connecting to database: %v\n", err)
			os.Exit(1)
		}
		defer db.Close()

		runner := migration.NewRunner(db, *migrationsDir)
		if err := runner.ValidateMigrations(ctx); err != nil {
			fmt.Fprintf(os.Stderr, "Validation failed: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Migrations validated successfully")

	default:
		fmt.Fprintf(os.Stderr, "Error: unknown command: %s\n", *command)
		fmt.Println("Available commands: up, down, down-to, status, create, validate")
		os.Exit(1)
	}
}

