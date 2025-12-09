package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/satishbabariya/jetorm/core"
	"github.com/satishbabariya/jetorm/hooks"
)

// User entity with comprehensive features
type User struct {
	ID        int64     `db:"id" jet:"primary_key,auto_increment" validate:"required"`
	Email     string    `db:"email" jet:"unique,not_null" validate:"required,email"`
	Username  string    `db:"username" jet:"unique,not_null" validate:"required,min:3"`
	Age       int       `db:"age" validate:"min:18"`
	Status    string    `db:"status" jet:"default:'active'"`
	CreatedAt time.Time `db:"created_at" jet:"auto_now_add"`
	UpdatedAt time.Time `db:"updated_at" jet:"auto_now"`
}

func exampleCompleteWorkflow() {
	fmt.Println("Complete JetORM Example")
	fmt.Println("======================")

	// 1. Connect to database
	config := core.Config{
		Host:     "localhost",
		Port:     5432,
		Database: "jetorm_test",
		User:     "postgres",
		Password: "secret",
	}

	db, err := core.Connect(config)
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer db.Close()

	// 2. Create base repository
	baseRepo, err := core.NewBaseRepository[User, int64](db)
	if err != nil {
		log.Fatalf("Failed to create repository: %v", err)
	}

	// 3. Set up hooks (simplified - User doesn't implement Auditable)
	userHooks := hooks.NewHooks[User]()
	// Note: For auditing, User would need to implement hooks.Auditable interface

	// 4. Set up cache
	cache := core.NewInMemoryCache()

	// 5. Set up validator
	validator := core.NewValidator()
	validator.RegisterRule("Email", core.Email())
	validator.RegisterRule("Username", core.Required())

	// 6. Set up performance monitor
	monitor := core.NewPerformanceMonitor(100 * time.Millisecond)
	profiler := core.NewQueryProfiler(monitor)

	// 7. Create full-featured repository
	fullRepo := core.NewFullFeaturedRepository(
		baseRepo,
		cache,
		"User",
		5*time.Minute,
		userHooks,
		validator,
		profiler,
		db,
	)

	ctx := context.Background()

	// 8. Health check
	health := fullRepo.HealthCheck(ctx)
	fmt.Printf("Database Health: %s - %s\n", health.Status, health.Message)

	// 9. Create user
	user := &User{
		Email:    "john@example.com",
		Username: "johndoe",
		Age:      25,
		Status:   "active",
	}

	saved, err := fullRepo.Save(ctx, user)
	if err != nil {
		log.Fatalf("Failed to save: %v", err)
	}
	fmt.Printf("Created user: ID=%d, Email=%s\n", saved.ID, saved.Email)

	// 10. Find user (cached)
	start := time.Now()
	found, err := fullRepo.FindByID(ctx, saved.ID)
	if err != nil {
		log.Fatalf("Failed to find: %v", err)
	}
	fmt.Printf("Found user (took %v): %s\n", time.Since(start), found.Email)

	// 11. Find again (should be cached)
	start = time.Now()
	found2, err := fullRepo.FindByID(ctx, saved.ID)
	if err != nil {
		log.Fatalf("Failed to find: %v", err)
	}
	fmt.Printf("Found user again (cached, took %v): %s\n", time.Since(start), found2.Email)

	// 12. Update user
	found.Status = "inactive"
	updated, err := fullRepo.Save(ctx, found) // Use Save for update
	if err != nil {
		log.Fatalf("Failed to update: %v", err)
	}
	fmt.Printf("Updated user: Status=%s\n", updated.Status)

	// 13. Use helpers with base repo
	exists, err := core.Exists(ctx, baseRepo, saved.ID)
	if err != nil {
		log.Fatalf("Failed to check existence: %v", err)
	}
	fmt.Printf("User exists: %v\n", exists)

	// 14. Batch operations with base repo
	users := []*User{
		{Email: "user1@example.com", Username: "user1", Age: 20},
		{Email: "user2@example.com", Username: "user2", Age: 21},
		{Email: "user3@example.com", Username: "user3", Age: 22},
	}

	err = core.OptimizedBatchSave(ctx, baseRepo, users, 10)
	if err != nil {
		log.Fatalf("Failed batch save: %v", err)
	}
	fmt.Printf("Saved %d users in batch\n", len(users))

	// 15. Get metrics
	metrics := monitor.GetAllMetrics()
	fmt.Printf("Performance metrics: %d queries tracked\n", len(metrics))

	// 16. Cleanup
	for _, u := range users {
		if u.ID > 0 {
			baseRepo.DeleteByID(ctx, u.ID)
		}
	}
	baseRepo.DeleteByID(ctx, saved.ID)
	fmt.Println("Cleanup complete")
}

func main() {
	exampleCompleteWorkflow()
}

