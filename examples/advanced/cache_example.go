package main

import (
	"context"
	"fmt"
	"time"

	"github.com/satishbabariya/jetorm/core"
)

type Product struct {
	ID    int64  `db:"id" jet:"primary_key"`
	Name  string `db:"name"`
	Price float64 `db:"price"`
}

func exampleCaching(ctx context.Context, productRepo core.Repository[Product, int64]) {
	fmt.Println("Caching Examples")
	fmt.Println("================")
	
	// Create cache
	cache := core.NewInMemoryCache()
	
	// Create cached repository
	cachedRepo := core.NewCachedRepository(
		productRepo,
		cache,
		"Product",
		5*time.Minute, // TTL
	)
	
	// First call - cache miss
	start := time.Now()
	product1, err := cachedRepo.FindByID(ctx, 1)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("First call (cache miss): %v (took %v)\n", product1, time.Since(start))
	
	// Second call - cache hit
	start = time.Now()
	product2, err := cachedRepo.FindByID(ctx, 1)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("Second call (cache hit): %v (took %v)\n", product2, time.Since(start))
	
	// Save - invalidates cache
	product1.Price = 99.99
	_, err = cachedRepo.Save(ctx, product1)
	if err != nil {
		fmt.Printf("Error saving: %v\n", err)
		return
	}
	fmt.Println("Cache invalidated after save")
}

// Uncomment to run example:
// func main() {
// 	ctx := context.Background()
// 	db := core.MustConnect(config)
// 	repo := core.NewBaseRepository[Product, int64](db)
// 	exampleCaching(ctx, repo)
// }

