package migration

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"time"
)

// Generator generates migration files from entity definitions
type Generator struct {
	schemaGen *SchemaGenerator
}

// NewGenerator creates a new migration generator
func NewGenerator() *Generator {
	return &Generator{
		schemaGen: NewSchemaGenerator(),
	}
}

// GenerateCreateTableMigration generates a CREATE TABLE migration from an entity type
func (g *Generator) GenerateCreateTableMigration(entityType reflect.Type, tableName string, migrationsDir string) error {
	if tableName == "" {
		tableName = toSnakeCase(entityType.Name())
	}

	// Generate CREATE TABLE SQL
	createSQL, err := g.schemaGen.GenerateCreateTable(entityType, tableName)
	if err != nil {
		return fmt.Errorf("failed to generate CREATE TABLE: %w", err)
	}

	// Generate DROP TABLE SQL for down migration
	dropSQL := fmt.Sprintf("DROP TABLE IF EXISTS %s;", tableName)

	// Create migration files
	version := time.Now().Format("20060102150405")
	sanitizedName := strings.ToLower(toSnakeCase(entityType.Name()))
	
	upFileName := fmt.Sprintf("%s_create_%s_table.up.sql", version, sanitizedName)
	downFileName := fmt.Sprintf("%s_create_%s_table.down.sql", version, sanitizedName)
	
	upPath := filepath.Join(migrationsDir, upFileName)
	downPath := filepath.Join(migrationsDir, downFileName)

	// Ensure directory exists
	if err := os.MkdirAll(migrationsDir, 0755); err != nil {
		return fmt.Errorf("failed to create migrations directory: %w", err)
	}

	// Write up migration
	upContent := fmt.Sprintf("-- Create table: %s\n-- Generated: %s\n\n%s\n", tableName, time.Now().Format(time.RFC3339), createSQL)
	if err := os.WriteFile(upPath, []byte(upContent), 0644); err != nil {
		return fmt.Errorf("failed to write up migration: %w", err)
	}

	// Write down migration
	downContent := fmt.Sprintf("-- Drop table: %s\n-- Generated: %s\n\n%s\n", tableName, time.Now().Format(time.RFC3339), dropSQL)
	if err := os.WriteFile(downPath, []byte(downContent), 0644); err != nil {
		return fmt.Errorf("failed to write down migration: %w", err)
	}

	return nil
}

// GenerateAlterTableMigration generates an ALTER TABLE migration
func (g *Generator) GenerateAlterTableMigration(tableName string, alterSQL string, migrationsDir string) error {
	version := time.Now().Format("20060102150405")
	sanitizedName := strings.ToLower(strings.ReplaceAll(tableName, " ", "_"))
	
	upFileName := fmt.Sprintf("%s_alter_%s.up.sql", version, sanitizedName)
	downFileName := fmt.Sprintf("%s_alter_%s.down.sql", version, sanitizedName)
	
	upPath := filepath.Join(migrationsDir, upFileName)
	downPath := filepath.Join(migrationsDir, downFileName)

	// Ensure directory exists
	if err := os.MkdirAll(migrationsDir, 0755); err != nil {
		return fmt.Errorf("failed to create migrations directory: %w", err)
	}

	// Write up migration
	upContent := fmt.Sprintf("-- Alter table: %s\n-- Generated: %s\n\n%s\n", tableName, time.Now().Format(time.RFC3339), alterSQL)
	if err := os.WriteFile(upPath, []byte(upContent), 0644); err != nil {
		return fmt.Errorf("failed to write up migration: %w", err)
	}

	// Write down migration (placeholder - would need reverse SQL)
	downContent := fmt.Sprintf("-- Rollback alter table: %s\n-- Generated: %s\n\n-- TODO: Add rollback SQL\n", tableName, time.Now().Format(time.RFC3339))
	if err := os.WriteFile(downPath, []byte(downContent), 0644); err != nil {
		return fmt.Errorf("failed to write down migration: %w", err)
	}

	return nil
}

// GenerateIndexMigration generates a CREATE INDEX migration
func (g *Generator) GenerateIndexMigration(tableName string, indexName string, columns []string, unique bool, migrationsDir string) error {
	version := time.Now().Format("20060102150405")
	sanitizedName := strings.ToLower(strings.ReplaceAll(indexName, " ", "_"))
	
	upFileName := fmt.Sprintf("%s_create_index_%s.up.sql", version, sanitizedName)
	downFileName := fmt.Sprintf("%s_create_index_%s.down.sql", version, sanitizedName)
	
	upPath := filepath.Join(migrationsDir, upFileName)
	downPath := filepath.Join(migrationsDir, downFileName)

	// Ensure directory exists
	if err := os.MkdirAll(migrationsDir, 0755); err != nil {
		return fmt.Errorf("failed to create migrations directory: %w", err)
	}

	// Build CREATE INDEX SQL
	uniqueClause := ""
	if unique {
		uniqueClause = "UNIQUE "
	}
	columnsStr := strings.Join(columns, ", ")
	createIndexSQL := fmt.Sprintf("CREATE %sINDEX IF NOT EXISTS %s ON %s (%s);", uniqueClause, indexName, tableName, columnsStr)
	dropIndexSQL := fmt.Sprintf("DROP INDEX IF EXISTS %s;", indexName)

	// Write up migration
	upContent := fmt.Sprintf("-- Create index: %s on %s\n-- Generated: %s\n\n%s\n", indexName, tableName, time.Now().Format(time.RFC3339), createIndexSQL)
	if err := os.WriteFile(upPath, []byte(upContent), 0644); err != nil {
		return fmt.Errorf("failed to write up migration: %w", err)
	}

	// Write down migration
	downContent := fmt.Sprintf("-- Drop index: %s\n-- Generated: %s\n\n%s\n", indexName, time.Now().Format(time.RFC3339), dropIndexSQL)
	if err := os.WriteFile(downPath, []byte(downContent), 0644); err != nil {
		return fmt.Errorf("failed to write down migration: %w", err)
	}

	return nil
}

// GenerateForeignKeyMigration generates a FOREIGN KEY migration
func (g *Generator) GenerateForeignKeyMigration(tableName string, columnName string, refTable string, refColumn string, onDelete string, onUpdate string, migrationsDir string) error {
	version := time.Now().Format("20060102150405")
	fkName := fmt.Sprintf("fk_%s_%s", tableName, columnName)
	sanitizedName := strings.ToLower(strings.ReplaceAll(fkName, " ", "_"))
	
	upFileName := fmt.Sprintf("%s_add_foreign_key_%s.up.sql", version, sanitizedName)
	downFileName := fmt.Sprintf("%s_add_foreign_key_%s.down.sql", version, sanitizedName)
	
	upPath := filepath.Join(migrationsDir, upFileName)
	downPath := filepath.Join(migrationsDir, downFileName)

	// Ensure directory exists
	if err := os.MkdirAll(migrationsDir, 0755); err != nil {
		return fmt.Errorf("failed to create migrations directory: %w", err)
	}

	// Build ALTER TABLE ADD FOREIGN KEY SQL
	onDeleteClause := ""
	if onDelete != "" {
		onDeleteClause = " ON DELETE " + strings.ToUpper(onDelete)
	}
	onUpdateClause := ""
	if onUpdate != "" {
		onUpdateClause = " ON UPDATE " + strings.ToUpper(onUpdate)
	}
	
	addFKSQL := fmt.Sprintf("ALTER TABLE %s ADD CONSTRAINT %s FOREIGN KEY (%s) REFERENCES %s (%s)%s%s;",
		tableName, fkName, columnName, refTable, refColumn, onDeleteClause, onUpdateClause)
	dropFKSQL := fmt.Sprintf("ALTER TABLE %s DROP CONSTRAINT IF EXISTS %s;", tableName, fkName)

	// Write up migration
	upContent := fmt.Sprintf("-- Add foreign key: %s.%s -> %s.%s\n-- Generated: %s\n\n%s\n",
		tableName, columnName, refTable, refColumn, time.Now().Format(time.RFC3339), addFKSQL)
	if err := os.WriteFile(upPath, []byte(upContent), 0644); err != nil {
		return fmt.Errorf("failed to write up migration: %w", err)
	}

	// Write down migration
	downContent := fmt.Sprintf("-- Drop foreign key: %s\n-- Generated: %s\n\n%s\n",
		fkName, time.Now().Format(time.RFC3339), dropFKSQL)
	if err := os.WriteFile(downPath, []byte(downContent), 0644); err != nil {
		return fmt.Errorf("failed to write down migration: %w", err)
	}

	return nil
}

// toSnakeCase converts a string to snake_case
func toSnakeCase(s string) string {
	var result strings.Builder
	for i, r := range s {
		if i > 0 && r >= 'A' && r <= 'Z' {
			result.WriteRune('_')
		}
		if r >= 'A' && r <= 'Z' {
			result.WriteRune(r + 32)
		} else {
			result.WriteRune(r)
		}
	}
	return result.String()
}

