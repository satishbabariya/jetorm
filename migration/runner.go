package migration

import (
	"context"
	"database/sql"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"
)

// Runner manages and executes migrations
type Runner struct {
	migrator *Migrator
	migrationsDir string
}

// NewRunner creates a new migration runner
func NewRunner(db *sql.DB, migrationsDir string) *Runner {
	return &Runner{
		migrator:      NewMigrator(db),
		migrationsDir: migrationsDir,
	}
}

// LoadMigrations loads migrations from the migrations directory
func (r *Runner) LoadMigrations(ctx context.Context) ([]Migration, error) {
	// Initialize migrator if database is available
	if r.migrator != nil && r.migrator.db != nil {
		if err := r.migrator.Initialize(ctx); err != nil {
			return nil, fmt.Errorf("failed to initialize migrator: %w", err)
		}
	}

	var migrations []Migration

	err := filepath.WalkDir(r.migrationsDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		// Check if file matches migration pattern: YYYYMMDDHHMMSS_name.up.sql or YYYYMMDDHHMMSS_name.down.sql
		baseName := filepath.Base(path)
		if !strings.HasSuffix(baseName, ".sql") {
			return nil
		}

		// Parse migration file name
		migration, err := r.parseMigrationFile(path, baseName)
		if err != nil {
			return err
		}

		if migration != nil {
			// Check if we already have this migration
			found := false
			for i, m := range migrations {
				if m.Version == migration.Version {
					// Update existing migration with up/down SQL
					if strings.HasSuffix(baseName, ".up.sql") {
						migrations[i].UpSQL = migration.UpSQL
					} else if strings.HasSuffix(baseName, ".down.sql") {
						migrations[i].DownSQL = migration.DownSQL
					}
					found = true
					break
				}
			}
			if !found {
				migrations = append(migrations, *migration)
			}
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to walk migrations directory: %w", err)
	}

	// Sort migrations by version
	sort.Slice(migrations, func(i, j int) bool {
		return migrations[i].Version < migrations[j].Version
	})

	return migrations, nil
}

// parseMigrationFile parses a migration file and returns a Migration
func (r *Runner) parseMigrationFile(path, fileName string) (*Migration, error) {
	// Parse file name: YYYYMMDDHHMMSS_name.up.sql or YYYYMMDDHHMMSS_name.down.sql
	parts := strings.Split(fileName, "_")
	if len(parts) < 2 {
		return nil, nil // Not a migration file
	}

	// Parse version (timestamp)
	versionStr := parts[0]
	version, err := strconv.ParseInt(versionStr, 10, 64)
	if err != nil {
		return nil, nil // Not a valid migration file
	}

	// Extract name and direction
	nameAndExt := strings.Join(parts[1:], "_")
	nameParts := strings.Split(nameAndExt, ".")
	if len(nameParts) < 3 {
		return nil, nil
	}

	name := nameParts[0]
	direction := nameParts[1] // "up" or "down"

	// Read file content
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read migration file %s: %w", path, err)
	}

	migration := &Migration{
		Version: version,
		Name:    name,
	}

	if direction == "up" {
		migration.UpSQL = string(content)
	} else if direction == "down" {
		migration.DownSQL = string(content)
	}

	return migration, nil
}

// Up applies all pending migrations
func (r *Runner) Up(ctx context.Context) error {
	migrations, err := r.LoadMigrations(ctx)
	if err != nil {
		return fmt.Errorf("failed to load migrations: %w", err)
	}

	appliedMigrations, err := r.migrator.GetAppliedMigrations(ctx)
	if err != nil {
		return fmt.Errorf("failed to get applied migrations: %w", err)
	}

	appliedVersions := make(map[int64]bool)
	for _, m := range appliedMigrations {
		appliedVersions[m.Version] = true
	}

	for _, migration := range migrations {
		if appliedVersions[migration.Version] {
			continue // Already applied
		}

		if migration.UpSQL == "" {
			return fmt.Errorf("migration %d (%s) has no up SQL", migration.Version, migration.Name)
		}

		if err := r.migrator.Apply(ctx, migration); err != nil {
			return fmt.Errorf("failed to apply migration %d (%s): %w", migration.Version, migration.Name, err)
		}
	}

	return nil
}

// Down rolls back the last migration
func (r *Runner) Down(ctx context.Context) error {
	appliedMigrations, err := r.migrator.GetAppliedMigrations(ctx)
	if err != nil {
		return fmt.Errorf("failed to get applied migrations: %w", err)
	}

	if len(appliedMigrations) == 0 {
		return fmt.Errorf("no migrations to rollback")
	}

	// Get the last applied migration
	lastMigration := appliedMigrations[len(appliedMigrations)-1]

	// Load migrations to get down SQL
	migrations, err := r.LoadMigrations(ctx)
	if err != nil {
		return fmt.Errorf("failed to load migrations: %w", err)
	}

	// Find the migration
	var migration *Migration
	for i := range migrations {
		if migrations[i].Version == lastMigration.Version {
			migration = &migrations[i]
			break
		}
	}

	if migration == nil {
		return fmt.Errorf("migration %d (%s) not found in migrations directory", lastMigration.Version, lastMigration.Name)
	}

	if migration.DownSQL == "" {
		return fmt.Errorf("migration %d (%s) has no down SQL", migration.Version, migration.Name)
	}

	return r.migrator.Rollback(ctx, *migration)
}

// DownTo rolls back migrations to a specific version
func (r *Runner) DownTo(ctx context.Context, targetVersion int64) error {
	appliedMigrations, err := r.migrator.GetAppliedMigrations(ctx)
	if err != nil {
		return fmt.Errorf("failed to get applied migrations: %w", err)
	}

	// Sort by version descending
	sort.Slice(appliedMigrations, func(i, j int) bool {
		return appliedMigrations[i].Version > appliedMigrations[j].Version
	})

	migrations, err := r.LoadMigrations(ctx)
	if err != nil {
		return fmt.Errorf("failed to load migrations: %w", err)
	}

	migrationMap := make(map[int64]*Migration)
	for i := range migrations {
		migrationMap[migrations[i].Version] = &migrations[i]
	}

	for _, applied := range appliedMigrations {
		if applied.Version <= targetVersion {
			break
		}

		migration, ok := migrationMap[applied.Version]
		if !ok {
			return fmt.Errorf("migration %d (%s) not found", applied.Version, applied.Name)
		}

		if migration.DownSQL == "" {
			return fmt.Errorf("migration %d (%s) has no down SQL", migration.Version, migration.Name)
		}

		if err := r.migrator.Rollback(ctx, *migration); err != nil {
			return fmt.Errorf("failed to rollback migration %d (%s): %w", migration.Version, migration.Name, err)
		}
	}

	return nil
}

// Status returns the status of migrations
func (r *Runner) Status(ctx context.Context) ([]MigrationStatus, error) {
	migrations, err := r.LoadMigrations(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to load migrations: %w", err)
	}

	appliedMigrations, err := r.migrator.GetAppliedMigrations(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get applied migrations: %w", err)
	}

	appliedVersions := make(map[int64]time.Time)
	for _, m := range appliedMigrations {
		if m.AppliedAt != nil {
			appliedVersions[m.Version] = *m.AppliedAt
		}
	}

	var statuses []MigrationStatus
	for _, migration := range migrations {
		status := MigrationStatus{
			Version: migration.Version,
			Name:    migration.Name,
			Status:  "pending",
		}

		if appliedAt, ok := appliedVersions[migration.Version]; ok {
			status.Status = "applied"
			status.AppliedAt = &appliedAt
		}

		statuses = append(statuses, status)
	}

	return statuses, nil
}

// MigrationStatus represents the status of a migration
type MigrationStatus struct {
	Version   int64
	Name      string
	Status    string // "applied" or "pending"
	AppliedAt *time.Time
}

// CreateMigration creates a new migration file pair
func (r *Runner) CreateMigration(name string) error {
	// Generate timestamp-based version
	version := time.Now().Format("20060102150405")
	
	// Sanitize name
	sanitizedName := strings.ToLower(strings.ReplaceAll(name, " ", "_"))
	
	// Create up migration file
	upFileName := fmt.Sprintf("%s_%s.up.sql", version, sanitizedName)
	upPath := filepath.Join(r.migrationsDir, upFileName)
	
	// Create down migration file
	downFileName := fmt.Sprintf("%s_%s.down.sql", version, sanitizedName)
	downPath := filepath.Join(r.migrationsDir, downFileName)
	
	// Ensure directory exists
	if err := os.MkdirAll(r.migrationsDir, 0755); err != nil {
		return fmt.Errorf("failed to create migrations directory: %w", err)
	}
	
	// Create up file
	upContent := fmt.Sprintf("-- Migration: %s\n-- Version: %s\n-- Up migration\n\n", name, version)
	if err := os.WriteFile(upPath, []byte(upContent), 0644); err != nil {
		return fmt.Errorf("failed to create up migration file: %w", err)
	}
	
	// Create down file
	downContent := fmt.Sprintf("-- Migration: %s\n-- Version: %s\n-- Down migration\n\n", name, version)
	if err := os.WriteFile(downPath, []byte(downContent), 0644); err != nil {
		return fmt.Errorf("failed to create down migration file: %w", err)
	}
	
	return nil
}

// ValidateMigrations validates that all migrations are properly paired
func (r *Runner) ValidateMigrations(ctx context.Context) error {
	migrations, err := r.LoadMigrations(ctx)
	if err != nil {
		return fmt.Errorf("failed to load migrations: %w", err)
	}

	for _, migration := range migrations {
		if migration.UpSQL == "" {
			return fmt.Errorf("migration %d (%s) is missing up SQL", migration.Version, migration.Name)
		}
		// Down SQL is optional but recommended
		if migration.DownSQL == "" {
			// Warning, not error
			fmt.Printf("Warning: migration %d (%s) is missing down SQL\n", migration.Version, migration.Name)
		}
	}

	return nil
}

// GetCurrentVersion returns the current database version
func (r *Runner) GetCurrentVersion(ctx context.Context) (int64, error) {
	return r.migrator.GetCurrentVersion(ctx)
}

// Version returns the version of a migration file name
func Version(fileName string) (int64, error) {
	parts := strings.Split(fileName, "_")
	if len(parts) == 0 {
		return 0, fmt.Errorf("invalid migration file name: %s", fileName)
	}
	return strconv.ParseInt(parts[0], 10, 64)
}

