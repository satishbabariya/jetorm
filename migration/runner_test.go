package migration

import (
	"context"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"
)

func TestRunner_LoadMigrations(t *testing.T) {
	tmpDir := t.TempDir()
	migrationsDir := filepath.Join(tmpDir, "migrations")
	os.MkdirAll(migrationsDir, 0755)

	// Create test migration files
	upContent := "CREATE TABLE test (id BIGINT PRIMARY KEY);"
	downContent := "DROP TABLE test;"

	version := time.Now().Format("20060102150405")
	upFile := filepath.Join(migrationsDir, version+"_create_test_table.up.sql")
	downFile := filepath.Join(migrationsDir, version+"_create_test_table.down.sql")

	os.WriteFile(upFile, []byte(upContent), 0644)
	os.WriteFile(downFile, []byte(downContent), 0644)

	// Create runner with nil DB (for file loading only)
	runner := &Runner{
		migrator:      nil, // No DB needed for loading files
		migrationsDir: migrationsDir,
	}
	
	migrations, err := runner.LoadMigrations(context.Background())
	if err != nil {
		t.Fatalf("Failed to load migrations: %v", err)
	}

	if len(migrations) != 1 {
		t.Errorf("Expected 1 migration, got %d", len(migrations))
	}

	if migrations[0].UpSQL == "" {
		t.Error("Migration should have UpSQL")
	}
	if migrations[0].DownSQL == "" {
		t.Error("Migration should have DownSQL")
	}
}

func TestRunner_CreateMigration(t *testing.T) {
	tmpDir := t.TempDir()
	migrationsDir := filepath.Join(tmpDir, "migrations")

	runner := NewRunner(nil, migrationsDir)
	err := runner.CreateMigration("test_migration")
	if err != nil {
		t.Fatalf("Failed to create migration: %v", err)
	}

	// Check that files were created
	files, err := os.ReadDir(migrationsDir)
	if err != nil {
		t.Fatalf("Failed to read migrations directory: %v", err)
	}

	if len(files) != 2 {
		t.Errorf("Expected 2 files (up and down), got %d", len(files))
	}

	hasUp := false
	hasDown := false
	for _, file := range files {
		if strings.Contains(file.Name(), ".up.sql") {
			hasUp = true
		}
		if strings.Contains(file.Name(), ".down.sql") {
			hasDown = true
		}
	}

	if !hasUp {
		t.Error("Up migration file not created")
	}
	if !hasDown {
		t.Error("Down migration file not created")
	}
}

func TestRunner_ValidateMigrations(t *testing.T) {
	tmpDir := t.TempDir()
	migrationsDir := filepath.Join(tmpDir, "migrations")
	os.MkdirAll(migrationsDir, 0755)

	// Create valid migration
	version := time.Now().Format("20060102150405")
	upFile := filepath.Join(migrationsDir, version+"_test.up.sql")
	os.WriteFile(upFile, []byte("CREATE TABLE test (id BIGINT);"), 0644)

	runner := NewRunner(nil, migrationsDir)
	err := runner.ValidateMigrations(context.Background())
	if err != nil {
		t.Errorf("Validation should pass for valid migration: %v", err)
	}
}

func TestSchemaGenerator_GenerateCreateTable(t *testing.T) {
	type TestUser struct {
		ID       int64  `db:"id" jet:"primary_key,auto_increment"`
		Email    string `db:"email" jet:"unique,not_null"`
		Username string `db:"username" jet:"not_null"`
		Age      int    `db:"age"`
	}

	sg := NewSchemaGenerator()
	sql, err := sg.GenerateCreateTable(reflect.TypeOf(TestUser{}), "test_users")
	if err != nil {
		t.Fatalf("Failed to generate CREATE TABLE: %v", err)
	}

	if !strings.Contains(sql, "CREATE TABLE") {
		t.Error("SQL should contain CREATE TABLE")
	}
	if !strings.Contains(sql, "test_users") {
		t.Error("SQL should contain table name")
	}
	if !strings.Contains(sql, "id") {
		t.Error("SQL should contain id column")
	}
	if !strings.Contains(sql, "email") {
		t.Error("SQL should contain email column")
	}
	if !strings.Contains(sql, "PRIMARY KEY") {
		t.Error("SQL should contain PRIMARY KEY")
	}
}

func TestGenerator_GenerateCreateTableMigration(t *testing.T) {
	type TestUser struct {
		ID    int64  `db:"id" jet:"primary_key"`
		Email string `db:"email" jet:"unique"`
	}

	tmpDir := t.TempDir()
	migrationsDir := filepath.Join(tmpDir, "migrations")

	gen := NewGenerator()
	err := gen.GenerateCreateTableMigration(reflect.TypeOf(TestUser{}), "test_users", migrationsDir)
	if err != nil {
		t.Fatalf("Failed to generate migration: %v", err)
	}

	// Check that files were created
	files, err := os.ReadDir(migrationsDir)
	if err != nil {
		t.Fatalf("Failed to read migrations directory: %v", err)
	}

	if len(files) != 2 {
		t.Errorf("Expected 2 files, got %d", len(files))
	}
}

func TestGenerator_GenerateIndexMigration(t *testing.T) {
	tmpDir := t.TempDir()
	migrationsDir := filepath.Join(tmpDir, "migrations")

	gen := NewGenerator()
	err := gen.GenerateIndexMigration("users", "idx_email", []string{"email"}, false, migrationsDir)
	if err != nil {
		t.Fatalf("Failed to generate index migration: %v", err)
	}

	// Check that files were created
	files, err := os.ReadDir(migrationsDir)
	if err != nil {
		t.Fatalf("Failed to read migrations directory: %v", err)
	}

	if len(files) != 2 {
		t.Errorf("Expected 2 files, got %d", len(files))
	}
}

func TestGenerator_GenerateForeignKeyMigration(t *testing.T) {
	tmpDir := t.TempDir()
	migrationsDir := filepath.Join(tmpDir, "migrations")

	gen := NewGenerator()
	err := gen.GenerateForeignKeyMigration("users", "company_id", "companies", "id", "cascade", "set_null", migrationsDir)
	if err != nil {
		t.Fatalf("Failed to generate foreign key migration: %v", err)
	}

	// Check that files were created
	files, err := os.ReadDir(migrationsDir)
	if err != nil {
		t.Fatalf("Failed to read migrations directory: %v", err)
	}

	if len(files) != 2 {
		t.Errorf("Expected 2 files, got %d", len(files))
	}
}

func TestValidator_ValidateSQL(t *testing.T) {
	validator := NewValidator(nil)

	// Valid SQL
	err := validator.ValidateSQL("CREATE TABLE test (id BIGINT);")
	if err != nil {
		t.Errorf("Valid SQL should pass validation: %v", err)
	}

	// Empty SQL
	err = validator.ValidateSQL("")
	if err == nil {
		t.Error("Empty SQL should fail validation")
	}

	// Dangerous SQL (should warn but not fail)
	err = validator.ValidateSQL("DROP DATABASE test;")
	if err == nil {
		t.Error("Dangerous SQL should fail validation")
	}
}

func TestValidator_ValidateMigrationOrder(t *testing.T) {
	validator := NewValidator(nil)

	// Valid order
	migrations := []Migration{
		{Version: 1, Name: "first"},
		{Version: 2, Name: "second"},
		{Version: 3, Name: "third"},
	}
	err := validator.ValidateMigrationOrder(migrations)
	if err != nil {
		t.Errorf("Valid order should pass: %v", err)
	}

	// Invalid order
	migrations = []Migration{
		{Version: 2, Name: "second"},
		{Version: 1, Name: "first"},
	}
	err = validator.ValidateMigrationOrder(migrations)
	if err == nil {
		t.Error("Invalid order should fail validation")
	}

	// Duplicate versions
	migrations = []Migration{
		{Version: 1, Name: "first"},
		{Version: 1, Name: "duplicate"},
	}
	err = validator.ValidateMigrationOrder(migrations)
	if err == nil {
		t.Error("Duplicate versions should fail validation")
	}
}

