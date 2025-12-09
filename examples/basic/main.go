package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/satishbabariya/jetorm/core"
)

// User represents a user entity
type User struct {
	ID        int64     `db:"id" jet:"primary_key,auto_increment"`
	Email     string    `db:"email" jet:"unique,not_null"`
	Username  string    `db:"username" jet:"unique,not_null"`
	FullName  string    `db:"full_name"`
	Age       int       `db:"age"`
	Status    string    `db:"status" jet:"default:'active'"`
	IsActive  bool      `db:"is_active" jet:"default:true"`
	CreatedAt time.Time `db:"created_at" jet:"auto_now_add,not_null"`
	UpdatedAt time.Time `db:"updated_at" jet:"auto_now,not_null"`
}

func main() {
	// Configure database connection
	config := core.Config{
		Host:     "localhost",
		Port:     5432,
		Database: "jetorm_test",
		User:     "postgres",
		Password: "postgres",
		SSLMode:  "disable",
		LogSQL:   true,
		LogLevel: core.DebugLevel,
	}

	// Connect to database
	db, err := core.Connect(config)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	fmt.Println("✓ Connected to database")

	// Create repository
	userRepo, err := core.NewBaseRepository[User, int64](db)
	if err != nil {
		log.Fatalf("Failed to create repository: %v", err)
	}

	ctx := context.Background()

	// Example 1: Create a new user
	fmt.Println("\n--- Example 1: Create User ---")
	newUser := &User{
		Email:    "john.doe@example.com",
		Username: "johndoe",
		FullName: "John Doe",
		Age:      30,
		Status:   "active",
		IsActive: true,
	}

	savedUser, err := userRepo.Save(ctx, newUser)
	if err != nil {
		log.Printf("Failed to save user: %v", err)
	} else {
		fmt.Printf("✓ Created user: ID=%d, Email=%s\n", savedUser.ID, savedUser.Email)
	}

	// Example 2: Find by ID
	if savedUser != nil {
		fmt.Println("\n--- Example 2: Find by ID ---")
		foundUser, err := userRepo.FindByID(ctx, savedUser.ID)
		if err != nil {
			log.Printf("Failed to find user: %v", err)
		} else {
			fmt.Printf("✓ Found user: %s (%s)\n", foundUser.Username, foundUser.Email)
		}
	}

	// Example 3: Update user
	if savedUser != nil {
		fmt.Println("\n--- Example 3: Update User ---")
		savedUser.Age = 31
		savedUser.Status = "premium"
		updatedUser, err := userRepo.Save(ctx, savedUser)
		if err != nil {
			log.Printf("Failed to update user: %v", err)
		} else {
			fmt.Printf("✓ Updated user: Age=%d, Status=%s\n", updatedUser.Age, updatedUser.Status)
		}
	}

	// Example 4: Count users
	fmt.Println("\n--- Example 4: Count Users ---")
	count, err := userRepo.Count(ctx)
	if err != nil {
		log.Printf("Failed to count users: %v", err)
	} else {
		fmt.Printf("✓ Total users: %d\n", count)
	}

	// Example 5: Find all users
	fmt.Println("\n--- Example 5: Find All Users ---")
	allUsers, err := userRepo.FindAll(ctx)
	if err != nil {
		log.Printf("Failed to find all users: %v", err)
	} else {
		fmt.Printf("✓ Found %d users:\n", len(allUsers))
		for _, u := range allUsers {
			fmt.Printf("  - %s (%s)\n", u.Username, u.Email)
		}
	}

	// Example 6: Pagination
	fmt.Println("\n--- Example 6: Pagination ---")
	page, err := userRepo.FindAllPaged(ctx, core.PageRequest(0, 10, core.Order{
		Field:     "created_at",
		Direction: core.Desc,
	}))
	if err != nil {
		log.Printf("Failed to get page: %v", err)
	} else {
		fmt.Printf("✓ Page %d of %d (Total: %d)\n", page.Number+1, page.TotalPages, page.TotalElements)
		for _, u := range page.Content {
			fmt.Printf("  - %s\n", u.Username)
		}
	}

	// Example 7: Check existence
	if savedUser != nil {
		fmt.Println("\n--- Example 7: Check Existence ---")
		exists, err := userRepo.ExistsById(ctx, savedUser.ID)
		if err != nil {
			log.Printf("Failed to check existence: %v", err)
		} else {
			fmt.Printf("✓ User exists: %v\n", exists)
		}
	}

	// Example 8: Transaction
	fmt.Println("\n--- Example 8: Transaction ---")
	err = db.Transaction(ctx, func(tx *core.Tx) error {
		txRepo := userRepo.WithTx(tx)

		// Create user in transaction
		txUser := &User{
			Email:    "jane.doe@example.com",
			Username: "janedoe",
			FullName: "Jane Doe",
			Age:      28,
			Status:   "active",
			IsActive: true,
		}

		saved, err := txRepo.Save(ctx, txUser)
		if err != nil {
			return err
		}

		fmt.Printf("✓ Created user in transaction: %s\n", saved.Username)
		return nil
	})
	if err != nil {
		log.Printf("Transaction failed: %v", err)
	}

	// Example 9: Delete user
	if savedUser != nil {
		fmt.Println("\n--- Example 9: Delete User ---")
		err = userRepo.DeleteByID(ctx, savedUser.ID)
		if err != nil {
			log.Printf("Failed to delete user: %v", err)
		} else {
			fmt.Printf("✓ Deleted user with ID: %d\n", savedUser.ID)
		}
	}

	fmt.Println("\n✓ All examples completed!")
}

