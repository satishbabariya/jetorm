package migration

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

// Migration represents a database migration
type Migration struct {
	Version   int64
	Name      string
	UpSQL     string
	DownSQL   string
	AppliedAt *time.Time
}

// Migrator manages database migrations
type Migrator struct {
	db        *sql.DB
	tableName string
}

// NewMigrator creates a new migrator instance
func NewMigrator(db *sql.DB) *Migrator {
	return &Migrator{
		db:        db,
		tableName: "schema_migrations",
	}
}

// SetTableName sets the name of the migrations tracking table
func (m *Migrator) SetTableName(name string) {
	m.tableName = name
}

// Initialize creates the migrations tracking table if it doesn't exist
func (m *Migrator) Initialize(ctx context.Context) error {
	query := fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %s (
			version BIGINT PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			applied_at TIMESTAMP NOT NULL DEFAULT NOW()
		)
	`, m.tableName)

	_, err := m.db.ExecContext(ctx, query)
	return err
}

// GetAppliedMigrations returns a list of applied migrations
func (m *Migrator) GetAppliedMigrations(ctx context.Context) ([]Migration, error) {
	if err := m.Initialize(ctx); err != nil {
		return nil, err
	}

	query := fmt.Sprintf("SELECT version, name, applied_at FROM %s ORDER BY version", m.tableName)
	rows, err := m.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var migrations []Migration
	for rows.Next() {
		var m Migration
		if err := rows.Scan(&m.Version, &m.Name, &m.AppliedAt); err != nil {
			return nil, err
		}
		migrations = append(migrations, m)
	}

	return migrations, rows.Err()
}

// IsApplied checks if a migration has been applied
func (m *Migrator) IsApplied(ctx context.Context, version int64) (bool, error) {
	if err := m.Initialize(ctx); err != nil {
		return false, err
	}

	query := fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE version = $1", m.tableName)
	var count int
	err := m.db.QueryRowContext(ctx, query, version).Scan(&count)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

// Apply applies a migration
func (m *Migrator) Apply(ctx context.Context, migration Migration) error {
	if err := m.Initialize(ctx); err != nil {
		return err
	}

	// Check if already applied
	applied, err := m.IsApplied(ctx, migration.Version)
	if err != nil {
		return err
	}
	if applied {
		return fmt.Errorf("migration %d (%s) already applied", migration.Version, migration.Name)
	}

	// Begin transaction
	tx, err := m.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Execute up migration
	if _, err := tx.ExecContext(ctx, migration.UpSQL); err != nil {
		return fmt.Errorf("failed to apply migration %d (%s): %w", migration.Version, migration.Name, err)
	}

	// Record migration
	recordQuery := fmt.Sprintf("INSERT INTO %s (version, name, applied_at) VALUES ($1, $2, NOW())", m.tableName)
	if _, err := tx.ExecContext(ctx, recordQuery, migration.Version, migration.Name); err != nil {
		return fmt.Errorf("failed to record migration %d (%s): %w", migration.Version, migration.Name, err)
	}

	return tx.Commit()
}

// Rollback rolls back a migration
func (m *Migrator) Rollback(ctx context.Context, migration Migration) error {
	if err := m.Initialize(ctx); err != nil {
		return err
	}

	// Check if applied
	applied, err := m.IsApplied(ctx, migration.Version)
	if err != nil {
		return err
	}
	if !applied {
		return fmt.Errorf("migration %d (%s) not applied", migration.Version, migration.Name)
	}

	// Begin transaction
	tx, err := m.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Execute down migration
	if migration.DownSQL != "" {
		if _, err := tx.ExecContext(ctx, migration.DownSQL); err != nil {
			return fmt.Errorf("failed to rollback migration %d (%s): %w", migration.Version, migration.Name, err)
		}
	}

	// Remove migration record
	recordQuery := fmt.Sprintf("DELETE FROM %s WHERE version = $1", m.tableName)
	if _, err := tx.ExecContext(ctx, recordQuery, migration.Version); err != nil {
		return fmt.Errorf("failed to remove migration record %d (%s): %w", migration.Version, migration.Name, err)
	}

	return tx.Commit()
}

// ApplyAll applies all pending migrations
func (m *Migrator) ApplyAll(ctx context.Context, migrations []Migration) error {
	for _, migration := range migrations {
		applied, err := m.IsApplied(ctx, migration.Version)
		if err != nil {
			return err
		}
		if !applied {
			if err := m.Apply(ctx, migration); err != nil {
				return err
			}
		}
	}
	return nil
}

// GetCurrentVersion returns the highest applied migration version
func (m *Migrator) GetCurrentVersion(ctx context.Context) (int64, error) {
	if err := m.Initialize(ctx); err != nil {
		return 0, err
	}

	query := fmt.Sprintf("SELECT MAX(version) FROM %s", m.tableName)
	var version sql.NullInt64
	err := m.db.QueryRowContext(ctx, query).Scan(&version)
	if err != nil {
		return 0, err
	}

	if !version.Valid {
		return 0, nil
	}

	return version.Int64, nil
}
