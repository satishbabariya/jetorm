package main

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/satishbabariya/jetorm/migration"
)

// User entity for migration generation example
type User struct {
	ID        int64  `db:"id" jet:"primary_key,auto_increment"`
	Email     string `db:"email" jet:"unique,not_null"`
	Username  string `db:"username" jet:"unique,not_null"`
	Age       int    `db:"age"`
	Status    string `db:"status" jet:"default:'active'"`
	CreatedAt string `db:"created_at" jet:"type:timestamp,default:now()"`
}

func exampleCreateMigration() {
	// Create a new migration manually
	migrationsDir := "./migrations"
	runner := migration.NewRunner(nil, migrationsDir)
	
	err := runner.CreateMigration("create_users_table")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	
	fmt.Println("Migration files created successfully")
}

func exampleGenerateFromEntity() {
	// Generate migration from entity definition
	migrationsDir := "./migrations"
	gen := migration.NewGenerator()
	
	err := gen.GenerateCreateTableMigration(reflect.TypeOf(User{}), "users", migrationsDir)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	
	fmt.Println("Migration generated from entity")
}

func exampleRunMigrations(db *sql.DB) {
	ctx := context.Background()
	migrationsDir := "./migrations"
	
	runner := migration.NewRunner(db, migrationsDir)
	
	// Apply all pending migrations
	err := runner.Up(ctx)
	if err != nil {
		fmt.Printf("Error applying migrations: %v\n", err)
		return
	}
	
	fmt.Println("Migrations applied successfully")
}

func exampleCheckStatus(db *sql.DB) {
	ctx := context.Background()
	migrationsDir := "./migrations"
	
	runner := migration.NewRunner(db, migrationsDir)
	
	statuses, err := runner.Status(ctx)
	if err != nil {
		fmt.Printf("Error getting status: %v\n", err)
		return
	}
	
	fmt.Println("Migration Status:")
	for _, status := range statuses {
		fmt.Printf("  %d - %s: %s\n", status.Version, status.Name, status.Status)
	}
}

func exampleRollback(db *sql.DB) {
	ctx := context.Background()
	migrationsDir := "./migrations"
	
	runner := migration.NewRunner(db, migrationsDir)
	
	// Rollback last migration
	err := runner.Down(ctx)
	if err != nil {
		fmt.Printf("Error rolling back: %v\n", err)
		return
	}
	
	fmt.Println("Migration rolled back successfully")
}

func exampleGenerateIndexMigration() {
	migrationsDir := "./migrations"
	gen := migration.NewGenerator()
	
	// Generate index migration
	err := gen.GenerateIndexMigration("users", "idx_email", []string{"email"}, true, migrationsDir)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	
	fmt.Println("Index migration generated")
}

func exampleGenerateForeignKeyMigration() {
	migrationsDir := "./migrations"
	gen := migration.NewGenerator()
	
	// Generate foreign key migration
	err := gen.GenerateForeignKeyMigration(
		"users", "company_id", "companies", "id",
		"cascade", "set_null", migrationsDir,
	)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	
	fmt.Println("Foreign key migration generated")
}

func main() {
	fmt.Println("Migration Examples")
	fmt.Println("==================")
	
	exampleCreateMigration()
	fmt.Println()
	
	exampleGenerateFromEntity()
	fmt.Println()
	
	exampleGenerateIndexMigration()
	fmt.Println()
	
	exampleGenerateForeignKeyMigration()
	fmt.Println()
	
	// Note: These require a database connection
	// Uncomment and provide DB connection to use:
	/*
	db, _ := sql.Open("pgx", "postgres://user:pass@localhost/dbname")
	defer db.Close()
	
	exampleRunMigrations(db)
	exampleCheckStatus(db)
	exampleRollback(db)
	*/
}

