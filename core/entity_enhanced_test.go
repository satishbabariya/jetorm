package core

import (
	"testing"
	"time"
)

// Test entities for enhanced tag parsing

type EnhancedUser struct {
	ID        int64      `db:"id" jet:"primary_key,auto_increment"`
	Email     string     `db:"email" jet:"unique,not_null,index:idx_email,size:255"`
	Username  string     `db:"username" jet:"unique,not_null,unique_index:idx_unique_username,size:50"`
	FullName  string     `db:"full_name" jet:"size:255"`
	Bio       string     `db:"bio" jet:"type:text"`
	Age       int        `db:"age" jet:"check:age >= 0 AND age <= 150"`
	Status    string     `db:"status" jet:"default:'active',check:status IN ('active','inactive'),index"`
	Balance   float64    `db:"balance" jet:"type:decimal(10,2),default:0.00"`
	CompanyID *int64     `db:"company_id" jet:"foreign_key:companies.id,on_delete:cascade,on_update:set_null"`
	RoleID    int64      `db:"role_id" jet:"foreign_key:roles.id,not_null,index:idx_role"`
	CreatedAt time.Time  `db:"created_at" jet:"auto_now_add,not_null,index"`
	UpdatedAt time.Time  `db:"updated_at" jet:"auto_now,not_null"`
	DeletedAt *time.Time `db:"deleted_at" jet:"index"`
	TempData  string     `db:"-" jet:"-"`
}

type ProductWithCompositeIndex struct {
	ID         int64  `db:"id" jet:"primary_key,auto_increment"`
	SKU        string `db:"sku" jet:"not_null,composite_index:idx_sku_store:1"`
	StoreID    int64  `db:"store_id" jet:"not_null,composite_index:idx_sku_store:2"`
	CategoryID int64  `db:"category_id" jet:"index"`
}

func TestEnhancedEntityTags(t *testing.T) {
	t.Run("parse enhanced user entity", func(t *testing.T) {
		entity, err := EntityMetadata(EnhancedUser{})
		if err != nil {
			t.Fatalf("Failed to parse entity: %v", err)
		}

		if entity.TableName != "enhanced_user" {
			t.Errorf("Expected table name 'enhanced_user', got '%s'", entity.TableName)
		}

		// Find email field
		var emailField *Field
		for i := range entity.Fields {
			if entity.Fields[i].Name == "Email" {
				emailField = &entity.Fields[i]
				break
			}
		}

		if emailField == nil {
			t.Fatal("Email field not found")
		}

		// Verify email field tags
		if !emailField.Unique {
			t.Error("Email field should be unique")
		}
		if !emailField.NotNull {
			t.Error("Email field should be not null")
		}
		if emailField.Index != "idx_email" {
			t.Errorf("Expected index 'idx_email', got '%s'", emailField.Index)
		}
		if emailField.Size != 255 {
			t.Errorf("Expected size 255, got %d", emailField.Size)
		}
	})

	t.Run("parse foreign key with cascade actions", func(t *testing.T) {
		entity, err := EntityMetadata(EnhancedUser{})
		if err != nil {
			t.Fatalf("Failed to parse entity: %v", err)
		}

		// Find CompanyID field
		var companyField *Field
		for i := range entity.Fields {
			if entity.Fields[i].Name == "CompanyID" {
				companyField = &entity.Fields[i]
				break
			}
		}

		if companyField == nil {
			t.Fatal("CompanyID field not found")
		}

		if companyField.ForeignKey != "companies.id" {
			t.Errorf("Expected foreign key 'companies.id', got '%s'", companyField.ForeignKey)
		}
		if companyField.OnDelete != "cascade" {
			t.Errorf("Expected on_delete 'cascade', got '%s'", companyField.OnDelete)
		}
		if companyField.OnUpdate != "set_null" {
			t.Errorf("Expected on_update 'set_null', got '%s'", companyField.OnUpdate)
		}
	})

	t.Run("parse explicit type", func(t *testing.T) {
		entity, err := EntityMetadata(EnhancedUser{})
		if err != nil {
			t.Fatalf("Failed to parse entity: %v", err)
		}

		// Find Bio field
		var bioField *Field
		for i := range entity.Fields {
			if entity.Fields[i].Name == "Bio" {
				bioField = &entity.Fields[i]
				break
			}
		}

		if bioField == nil {
			t.Fatal("Bio field not found")
		}

		if bioField.ExplicitType != "text" {
			t.Errorf("Expected type 'text', got '%s'", bioField.ExplicitType)
		}

		// Find Balance field
		var balanceField *Field
		for i := range entity.Fields {
			if entity.Fields[i].Name == "Balance" {
				balanceField = &entity.Fields[i]
				break
			}
		}

		if balanceField == nil {
			t.Fatal("Balance field not found")
		}

		if balanceField.ExplicitType != "decimal(10,2)" {
			t.Errorf("Expected type 'decimal(10,2)', got '%s'", balanceField.ExplicitType)
		}
		if balanceField.Default != "0.00" {
			t.Errorf("Expected default '0.00', got '%s'", balanceField.Default)
		}
	})

	t.Run("parse check constraint", func(t *testing.T) {
		entity, err := EntityMetadata(EnhancedUser{})
		if err != nil {
			t.Fatalf("Failed to parse entity: %v", err)
		}

		// Find Age field
		var ageField *Field
		for i := range entity.Fields {
			if entity.Fields[i].Name == "Age" {
				ageField = &entity.Fields[i]
				break
			}
		}

		if ageField == nil {
			t.Fatal("Age field not found")
		}

		expectedCheck := "age >= 0 AND age <= 150"
		if ageField.Check != expectedCheck {
			t.Errorf("Expected check '%s', got '%s'", expectedCheck, ageField.Check)
		}
	})

	t.Run("parse unique index", func(t *testing.T) {
		entity, err := EntityMetadata(EnhancedUser{})
		if err != nil {
			t.Fatalf("Failed to parse entity: %v", err)
		}

		// Find Username field
		var usernameField *Field
		for i := range entity.Fields {
			if entity.Fields[i].Name == "Username" {
				usernameField = &entity.Fields[i]
				break
			}
		}

		if usernameField == nil {
			t.Fatal("Username field not found")
		}

		if usernameField.UniqueIndex != "idx_unique_username" {
			t.Errorf("Expected unique index 'idx_unique_username', got '%s'", usernameField.UniqueIndex)
		}
	})

	t.Run("parse composite index", func(t *testing.T) {
		entity, err := EntityMetadata(ProductWithCompositeIndex{})
		if err != nil {
			t.Fatalf("Failed to parse entity: %v", err)
		}

		// Find SKU field
		var skuField *Field
		for i := range entity.Fields {
			if entity.Fields[i].Name == "SKU" {
				skuField = &entity.Fields[i]
				break
			}
		}

		if skuField == nil {
			t.Fatal("SKU field not found")
		}

		if skuField.CompositeIndex == nil {
			t.Fatal("SKU field should have composite index")
		}

		if skuField.CompositeIndex.Name != "idx_sku_store" {
			t.Errorf("Expected composite index name 'idx_sku_store', got '%s'", skuField.CompositeIndex.Name)
		}
		if skuField.CompositeIndex.Order != 1 {
			t.Errorf("Expected composite index order 1, got %d", skuField.CompositeIndex.Order)
		}

		// Find StoreID field
		var storeField *Field
		for i := range entity.Fields {
			if entity.Fields[i].Name == "StoreID" {
				storeField = &entity.Fields[i]
				break
			}
		}

		if storeField == nil {
			t.Fatal("StoreID field not found")
		}

		if storeField.CompositeIndex == nil {
			t.Fatal("StoreID field should have composite index")
		}

		if storeField.CompositeIndex.Name != "idx_sku_store" {
			t.Errorf("Expected composite index name 'idx_sku_store', got '%s'", storeField.CompositeIndex.Name)
		}
		if storeField.CompositeIndex.Order != 2 {
			t.Errorf("Expected composite index order 2, got %d", storeField.CompositeIndex.Order)
		}
	})

	t.Run("parse ignored field", func(t *testing.T) {
		entity, err := EntityMetadata(EnhancedUser{})
		if err != nil {
			t.Fatalf("Failed to parse entity: %v", err)
		}

		// Find TempData field
		var tempField *Field
		for i := range entity.Fields {
			if entity.Fields[i].Name == "TempData" {
				tempField = &entity.Fields[i]
				break
			}
		}

		if tempField == nil {
			t.Fatal("TempData field not found")
		}

		if !tempField.Ignored {
			t.Error("TempData field should be ignored")
		}
		if tempField.DBName != "-" {
			t.Errorf("Expected DBName '-', got '%s'", tempField.DBName)
		}
	})

	t.Run("parse default value with quotes", func(t *testing.T) {
		entity, err := EntityMetadata(EnhancedUser{})
		if err != nil {
			t.Fatalf("Failed to parse entity: %v", err)
		}

		// Find Status field
		var statusField *Field
		for i := range entity.Fields {
			if entity.Fields[i].Name == "Status" {
				statusField = &entity.Fields[i]
				break
			}
		}

		if statusField == nil {
			t.Fatal("Status field not found")
		}

		if statusField.Default != "'active'" {
			t.Errorf("Expected default ''active'', got '%s'", statusField.Default)
		}
	})
}

