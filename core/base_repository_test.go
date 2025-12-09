package core

import (
	"testing"
	"time"
)

// TestUser is a test entity
type TestUser struct {
	ID        int64     `db:"id" jet:"primary_key,auto_increment"`
	Email     string    `db:"email" jet:"unique,not_null"`
	Username  string    `db:"username" jet:"unique,not_null"`
	Age       int       `db:"age"`
	CreatedAt time.Time `db:"created_at" jet:"auto_now_add"`
	UpdatedAt time.Time `db:"updated_at" jet:"auto_now"`
}

func TestNewBaseRepository(t *testing.T) {
	// This is a unit test that doesn't require a database connection
	// We're just testing the repository creation logic
	
	t.Run("should create repository with valid entity", func(t *testing.T) {
		// Note: This will fail without a real DB connection
		// This is a placeholder to show the test structure
		// In a real scenario, we'd use testcontainers or a mock
		
		// For now, we just test entity metadata extraction
		entity, err := EntityMetadata(TestUser{})
		if err != nil {
			t.Fatalf("Failed to extract entity metadata: %v", err)
		}
		
		if entity.TableName != "test_user" {
			t.Errorf("Expected table name 'test_user', got '%s'", entity.TableName)
		}
		
		if entity.PrimaryKey == nil {
			t.Error("Expected primary key to be set")
		} else if entity.PrimaryKey.Name != "ID" {
			t.Errorf("Expected primary key field 'ID', got '%s'", entity.PrimaryKey.Name)
		}
	})
}

func TestEntityMetadata(t *testing.T) {
	t.Run("should extract entity metadata", func(t *testing.T) {
		entity, err := EntityMetadata(TestUser{})
		if err != nil {
			t.Fatalf("Failed to extract metadata: %v", err)
		}
		
		// Check table name
		if entity.TableName != "test_user" {
			t.Errorf("Expected table name 'test_user', got '%s'", entity.TableName)
		}
		
		// Check fields
		expectedFields := []string{"id", "email", "username", "age", "created_at", "updated_at"}
		if len(entity.Fields) != len(expectedFields) {
			t.Errorf("Expected %d fields, got %d", len(expectedFields), len(entity.Fields))
		}
		
		// Check primary key
		if entity.PrimaryKey == nil {
			t.Fatal("Expected primary key to be set")
		}
		if entity.PrimaryKey.DBName != "id" {
			t.Errorf("Expected primary key 'id', got '%s'", entity.PrimaryKey.DBName)
		}
		if !entity.PrimaryKey.PrimaryKey {
			t.Error("Expected primary key flag to be true")
		}
		if !entity.PrimaryKey.AutoIncrement {
			t.Error("Expected auto_increment flag to be true")
		}
	})
	
	t.Run("should handle entity without primary key", func(t *testing.T) {
		type NoPKEntity struct {
			Name string `db:"name"`
		}
		
		entity, err := EntityMetadata(NoPKEntity{})
		if err != nil {
			t.Fatalf("Failed to extract metadata: %v", err)
		}
		
		if entity.PrimaryKey != nil {
			t.Error("Expected no primary key")
		}
	})
	
	t.Run("should handle non-struct type", func(t *testing.T) {
		_, err := EntityMetadata("not a struct")
		if err != ErrInvalidEntity {
			t.Errorf("Expected ErrInvalidEntity, got %v", err)
		}
	})
}

func TestPageable(t *testing.T) {
	t.Run("should create pageable with orders", func(t *testing.T) {
		pageable := PageRequest(0, 10,
			Order{Field: "created_at", Direction: Desc},
			Order{Field: "username", Direction: Asc},
		)
		
		if pageable.Page != 0 {
			t.Errorf("Expected page 0, got %d", pageable.Page)
		}
		if pageable.Size != 10 {
			t.Errorf("Expected size 10, got %d", pageable.Size)
		}
		if len(pageable.Sort.Orders) != 2 {
			t.Errorf("Expected 2 orders, got %d", len(pageable.Sort.Orders))
		}
	})
	
	t.Run("should create unpaged", func(t *testing.T) {
		pageable := Unpaged()
		
		if pageable.Page != 0 {
			t.Errorf("Expected page 0, got %d", pageable.Page)
		}
		if pageable.Size != -1 {
			t.Errorf("Expected size -1, got %d", pageable.Size)
		}
	})
}

func TestUtils(t *testing.T) {
	t.Run("toSnakeCase", func(t *testing.T) {
		tests := []struct {
			input    string
			expected string
		}{
			{"TestUser", "test_user"},
			{"UserProfile", "user_profile"},
			{"ID", "i_d"},
			{"HTTPServer", "h_t_t_p_server"},
			{"lowercase", "lowercase"},
		}
		
		for _, tt := range tests {
			result := toSnakeCase(tt.input)
			if result != tt.expected {
				t.Errorf("toSnakeCase(%s) = %s, expected %s", tt.input, result, tt.expected)
			}
		}
	})
	
	t.Run("toCamelCase", func(t *testing.T) {
		tests := []struct {
			input    string
			expected string
		}{
			{"test_user", "testUser"},
			{"user_profile", "userProfile"},
			{"id", "id"},
			{"first_name", "firstName"},
		}
		
		for _, tt := range tests {
			result := toCamelCase(tt.input)
			if result != tt.expected {
				t.Errorf("toCamelCase(%s) = %s, expected %s", tt.input, result, tt.expected)
			}
		}
	})
	
	t.Run("toPascalCase", func(t *testing.T) {
		tests := []struct {
			input    string
			expected string
		}{
			{"test_user", "TestUser"},
			{"user_profile", "UserProfile"},
			{"id", "Id"},
			{"first_name", "FirstName"},
		}
		
		for _, tt := range tests {
			result := toPascalCase(tt.input)
			if result != tt.expected {
				t.Errorf("toPascalCase(%s) = %s, expected %s", tt.input, result, tt.expected)
			}
		}
	})
}

// Integration tests would go here
// They would require a real database connection (using testcontainers)
// Example structure:

/*
func TestBaseRepositoryIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}
	
	// Setup test database using testcontainers
	ctx := context.Background()
	container, db := setupTestDatabase(t, ctx)
	defer container.Terminate(ctx)
	defer db.Close()
	
	// Create repository
	repo, err := NewBaseRepository[TestUser, int64](db)
	if err != nil {
		t.Fatalf("Failed to create repository: %v", err)
	}
	
	t.Run("Save and FindByID", func(t *testing.T) {
		user := &TestUser{
			Email:    "test@example.com",
			Username: "testuser",
			Age:      25,
		}
		
		saved, err := repo.Save(ctx, user)
		if err != nil {
			t.Fatalf("Failed to save user: %v", err)
		}
		
		if saved.ID == 0 {
			t.Error("Expected non-zero ID")
		}
		
		found, err := repo.FindByID(ctx, saved.ID)
		if err != nil {
			t.Fatalf("Failed to find user: %v", err)
		}
		
		if found.Email != user.Email {
			t.Errorf("Expected email %s, got %s", user.Email, found.Email)
		}
	})
	
	// More integration tests...
}
*/

