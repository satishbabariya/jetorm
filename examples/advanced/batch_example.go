package main

import (
	"context"
	"fmt"
	"time"

	"github.com/satishbabariya/jetorm/core"
)

type Order struct {
	ID        int64     `db:"id" jet:"primary_key,auto_increment"`
	UserID    int64     `db:"user_id"`
	Amount    float64   `db:"amount"`
	CreatedAt time.Time `db:"created_at"`
}

func exampleBatchWriter(ctx context.Context, orderRepo core.Repository[Order, int64]) {
	fmt.Println("Batch Writer Example")
	fmt.Println("====================")
	
	// Create batch writer
	config := core.DefaultBatchConfig()
	config.Size = 50 // Batch size
	config.FlushInterval = 2 * time.Second
	
	batchWriter := core.NewBatchWriter(orderRepo, config)
	defer batchWriter.Close(ctx)
	
	// Write orders in batches
	for i := 0; i < 150; i++ {
		order := &Order{
			UserID: int64(i % 10),
			Amount: float64(i * 10),
			CreatedAt: time.Now(),
		}
		
		if err := batchWriter.Write(ctx, order); err != nil {
			fmt.Printf("Error writing: %v\n", err)
			return
		}
	}
	
	fmt.Println("All orders written in batches")
}

func exampleOptimizedBatchSave(ctx context.Context, orderRepo core.Repository[Order, int64]) {
	fmt.Println("Optimized Batch Save Example")
	fmt.Println("============================")
	
	// Create orders
	orders := make([]*Order, 500)
	for i := 0; i < 500; i++ {
		orders[i] = &Order{
			UserID: int64(i % 10),
			Amount: float64(i * 10),
			CreatedAt: time.Now(),
		}
	}
	
	// Save in optimized batches
	start := time.Now()
	err := core.OptimizedBatchSave(ctx, orderRepo, orders, 100)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	
	fmt.Printf("Saved %d orders in %v\n", len(orders), time.Since(start))
}

func exampleBatchProcessor() {
	fmt.Println("Batch Processor Example")
	fmt.Println("========================")
	
	// Create entities
	entities := make([]*Order, 1000)
	for i := 0; i < 1000; i++ {
		entities[i] = &Order{ID: int64(i)}
	}
	
	// Process in batches
	processor := core.NewBatchProcessor(100, func(batch []*Order) error {
		fmt.Printf("Processing batch of %d orders\n", len(batch))
		// Process batch...
		return nil
	})
	
	err := processor.Process(entities)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	
	fmt.Println("All entities processed")
}

// Uncomment to run example:
// func main() {
// 	ctx := context.Background()
// 	db := core.MustConnect(config)
// 	repo := core.NewBaseRepository[Order, int64](db)
// 	exampleBatchWriter(ctx, repo)
// 	exampleOptimizedBatchSave(ctx, repo)
// 	exampleBatchProcessor()
// }

