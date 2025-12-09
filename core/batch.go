package core

import (
	"context"
	"fmt"
	"time"
)

// BatchConfig configures batch operations
type BatchConfig struct {
	Size         int           // Batch size
	FlushInterval time.Duration // Auto-flush interval
	Timeout      time.Duration // Operation timeout
}

// DefaultBatchConfig returns default batch configuration
func DefaultBatchConfig() BatchConfig {
	return BatchConfig{
		Size:         100,
		FlushInterval: 5 * time.Second,
		Timeout:      30 * time.Second,
	}
}

// BatchWriter provides optimized batch writing
type BatchWriter[T any, ID comparable] struct {
	repo   Repository[T, ID]
	config BatchConfig
	buffer []*T
	ticker *time.Ticker
	done   chan bool
}

// NewBatchWriter creates a new batch writer
func NewBatchWriter[T any, ID comparable](repo Repository[T, ID], config BatchConfig) *BatchWriter[T, ID] {
	bw := &BatchWriter[T, ID]{
		repo:   repo,
		config: config,
		buffer: make([]*T, 0, config.Size),
		done:   make(chan bool),
	}
	
	// Start auto-flush ticker
	if config.FlushInterval > 0 {
		bw.ticker = time.NewTicker(config.FlushInterval)
		go bw.autoFlush()
	}
	
	return bw
}

// Write adds an entity to the batch
func (bw *BatchWriter[T, ID]) Write(ctx context.Context, entity *T) error {
	bw.buffer = append(bw.buffer, entity)
	
	// Flush if buffer is full
	if len(bw.buffer) >= bw.config.Size {
		return bw.Flush(ctx)
	}
	
	return nil
}

// Flush flushes the batch buffer
func (bw *BatchWriter[T, ID]) Flush(ctx context.Context) error {
	if len(bw.buffer) == 0 {
		return nil
	}
	
	// Create context with timeout
	flushCtx := ctx
	if bw.config.Timeout > 0 {
		var cancel context.CancelFunc
		flushCtx, cancel = context.WithTimeout(ctx, bw.config.Timeout)
		defer cancel()
	}
	
	// Save batch
	err := bw.repo.SaveBatch(flushCtx, bw.buffer, bw.config.Size)
	if err != nil {
		return fmt.Errorf("batch flush failed: %w", err)
	}
	
	// Clear buffer
	bw.buffer = bw.buffer[:0]
	
	return nil
}

// Close closes the batch writer and flushes remaining data
func (bw *BatchWriter[T, ID]) Close(ctx context.Context) error {
	// Stop ticker
	if bw.ticker != nil {
		bw.ticker.Stop()
		bw.done <- true
	}
	
	// Flush remaining
	return bw.Flush(ctx)
}

// autoFlush automatically flushes the buffer at intervals
func (bw *BatchWriter[T, ID]) autoFlush() {
	for {
		select {
		case <-bw.ticker.C:
			ctx := context.Background()
			bw.Flush(ctx)
		case <-bw.done:
			return
		}
	}
}

// BatchReader provides optimized batch reading
type BatchReader[T any, ID comparable] struct {
	repo   Repository[T, ID]
	config BatchConfig
	cursor ID
	limit  int
}

// NewBatchReader creates a new batch reader
func NewBatchReader[T any, ID comparable](repo Repository[T, ID], config BatchConfig) *BatchReader[T, ID] {
	return &BatchReader[T, ID]{
		repo:   repo,
		config: config,
		limit:  config.Size,
	}
}

// ReadBatch reads a batch of entities
func (br *BatchReader[T, ID]) ReadBatch(ctx context.Context) ([]*T, error) {
	// This is a simplified version - would need cursor-based pagination
	// For now, use FindAllPaged
	pageable := PageRequest(0, br.limit)
	page, err := br.repo.FindAllPaged(ctx, pageable)
	if err != nil {
		return nil, err
	}
	
	return page.Content, nil
}

// OptimizedBatchSave performs optimized batch save with batching
func OptimizedBatchSave[T any, ID comparable](
	ctx context.Context,
	repo Repository[T, ID],
	entities []*T,
	batchSize int,
) error {
	if batchSize <= 0 {
		batchSize = 100 // Default
	}
	
	for i := 0; i < len(entities); i += batchSize {
		end := i + batchSize
		if end > len(entities) {
			end = len(entities)
		}
		
		batch := entities[i:end]
		if err := repo.SaveBatch(ctx, batch, batchSize); err != nil {
			return fmt.Errorf("batch save failed at offset %d: %w", i, err)
		}
	}
	
	return nil
}

// BatchProcessor processes entities in batches
type BatchProcessor[T any] struct {
	batchSize int
	processor func([]*T) error
}

// NewBatchProcessor creates a new batch processor
func NewBatchProcessor[T any](batchSize int, processor func([]*T) error) *BatchProcessor[T] {
	return &BatchProcessor[T]{
		batchSize: batchSize,
		processor: processor,
	}
}

// Process processes entities in batches
func (bp *BatchProcessor[T]) Process(entities []*T) error {
	for i := 0; i < len(entities); i += bp.batchSize {
		end := i + bp.batchSize
		if end > len(entities) {
			end = len(entities)
		}
		
		batch := entities[i:end]
		if err := bp.processor(batch); err != nil {
			return fmt.Errorf("batch processing failed at offset %d: %w", i, err)
		}
	}
	
	return nil
}

