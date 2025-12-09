package migration

import (
	"context"
	"database/sql"
	"fmt"
	"regexp"
	"strings"
)

// Validator validates migrations
type Validator struct {
	db *sql.DB
}

// NewValidator creates a new migration validator
func NewValidator(db *sql.DB) *Validator {
	return &Validator{
		db: db,
	}
}

// ValidateSQL validates SQL syntax
func (v *Validator) ValidateSQL(sql string) error {
	// Basic SQL validation
	if strings.TrimSpace(sql) == "" {
		return fmt.Errorf("SQL is empty")
	}

	// Check for dangerous operations in production
	dangerousPatterns := []struct {
		pattern string
		message string
	}{
		{`(?i)\bDROP\s+DATABASE\b`, "DROP DATABASE is not allowed"},
		{`(?i)\bDROP\s+SCHEMA\b`, "DROP SCHEMA is not allowed"},
		{`(?i)\bTRUNCATE\b`, "TRUNCATE should be used carefully"},
	}

	for _, dp := range dangerousPatterns {
		matched, err := regexp.MatchString(dp.pattern, sql)
		if err != nil {
			continue
		}
		if matched {
			return fmt.Errorf("validation warning: %s", dp.message)
		}
	}

	return nil
}

// ValidateMigration validates a migration before applying
func (v *Validator) ValidateMigration(ctx context.Context, migration Migration) error {
	// Validate up SQL
	if migration.UpSQL != "" {
		if err := v.ValidateSQL(migration.UpSQL); err != nil {
			return fmt.Errorf("up SQL validation failed: %w", err)
		}
	}

	// Validate down SQL
	if migration.DownSQL != "" {
		if err := v.ValidateSQL(migration.DownSQL); err != nil {
			return fmt.Errorf("down SQL validation failed: %w", err)
		}
	}

	// Check for required up SQL
	if migration.UpSQL == "" {
		return fmt.Errorf("migration %d (%s) has no up SQL", migration.Version, migration.Name)
	}

	return nil
}

// ValidateMigrationOrder validates that migrations are in correct order
func (v *Validator) ValidateMigrationOrder(migrations []Migration) error {
	if len(migrations) == 0 {
		return nil
	}

	// Check for duplicate versions
	versions := make(map[int64]bool)
	for _, m := range migrations {
		if versions[m.Version] {
			return fmt.Errorf("duplicate migration version: %d", m.Version)
		}
		versions[m.Version] = true
	}

	// Check ordering
	for i := 1; i < len(migrations); i++ {
		if migrations[i].Version <= migrations[i-1].Version {
			return fmt.Errorf("migrations are not in ascending order: %d <= %d",
				migrations[i].Version, migrations[i-1].Version)
		}
	}

	return nil
}

// ValidateAppliedMigrations validates that applied migrations match files
func (v *Validator) ValidateAppliedMigrations(ctx context.Context, fileMigrations []Migration, appliedMigrations []Migration) error {
	fileVersions := make(map[int64]Migration)
	for _, m := range fileMigrations {
		fileVersions[m.Version] = m
	}

	appliedVersions := make(map[int64]bool)
	for _, m := range appliedMigrations {
		appliedVersions[m.Version] = true
	}

	// Check for applied migrations that don't exist in files
	for _, applied := range appliedMigrations {
		if _, exists := fileVersions[applied.Version]; !exists {
			return fmt.Errorf("applied migration %d (%s) not found in migration files", applied.Version, applied.Name)
		}
	}

	// Check for missing migrations (applied but file missing)
	for version := range appliedVersions {
		if _, exists := fileVersions[version]; !exists {
			return fmt.Errorf("migration file for version %d not found", version)
		}
	}

	return nil
}

// CheckMigrationIntegrity checks the integrity of migrations
func (v *Validator) CheckMigrationIntegrity(ctx context.Context, migrations []Migration) error {
	// Validate order
	if err := v.ValidateMigrationOrder(migrations); err != nil {
		return err
	}

	// Validate each migration
	for _, migration := range migrations {
		if err := v.ValidateMigration(ctx, migration); err != nil {
			return fmt.Errorf("migration %d (%s): %w", migration.Version, migration.Name, err)
		}
	}

	return nil
}

// CheckDatabaseState checks if database is in a valid state for migrations
func (v *Validator) CheckDatabaseState(ctx context.Context) error {
	// Check if migrations table exists
	var exists bool
	err := v.db.QueryRowContext(ctx, `
		SELECT EXISTS (
			SELECT FROM information_schema.tables 
			WHERE table_schema = 'public' 
			AND table_name = 'schema_migrations'
		)
	`).Scan(&exists)
	
	if err != nil {
		return fmt.Errorf("failed to check migrations table: %w", err)
	}

	if !exists {
		// This is okay - table will be created on first migration
		return nil
	}

	// Check for any issues with migrations table
	var count int
	err = v.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM schema_migrations").Scan(&count)
	if err != nil {
		return fmt.Errorf("failed to query migrations table: %w", err)
	}

	return nil
}

