package core

import (
	"context"
	"time"
)

// PerformanceMonitor monitors query performance
type PerformanceMonitor struct {
	slowQueryThreshold time.Duration
	metrics            map[string]*QueryMetrics
}

// QueryMetrics tracks metrics for a query
type QueryMetrics struct {
	Count         int64
	TotalDuration time.Duration
	MinDuration   time.Duration
	MaxDuration   time.Duration
	AvgDuration   time.Duration
	SlowQueries   int64
}

// NewPerformanceMonitor creates a new performance monitor
func NewPerformanceMonitor(slowQueryThreshold time.Duration) *PerformanceMonitor {
	return &PerformanceMonitor{
		slowQueryThreshold: slowQueryThreshold,
		metrics:            make(map[string]*QueryMetrics),
	}
}

// RecordQuery records a query execution
func (pm *PerformanceMonitor) RecordQuery(query string, duration time.Duration) {
	metrics, exists := pm.metrics[query]
	if !exists {
		metrics = &QueryMetrics{
			MinDuration: duration,
			MaxDuration: duration,
		}
		pm.metrics[query] = metrics
	}

	metrics.Count++
	metrics.TotalDuration += duration

	if duration < metrics.MinDuration {
		metrics.MinDuration = duration
	}
	if duration > metrics.MaxDuration {
		metrics.MaxDuration = duration
	}

	metrics.AvgDuration = metrics.TotalDuration / time.Duration(metrics.Count)

	if duration > pm.slowQueryThreshold {
		metrics.SlowQueries++
	}
}

// GetMetrics returns metrics for a query
func (pm *PerformanceMonitor) GetMetrics(query string) *QueryMetrics {
	return pm.metrics[query]
}

// GetAllMetrics returns all metrics
func (pm *PerformanceMonitor) GetAllMetrics() map[string]*QueryMetrics {
	return pm.metrics
}

// Reset resets all metrics
func (pm *PerformanceMonitor) Reset() {
	pm.metrics = make(map[string]*QueryMetrics)
}

// QueryProfiler profiles query execution
type QueryProfiler struct {
	monitor *PerformanceMonitor
}

// NewQueryProfiler creates a new query profiler
func NewQueryProfiler(monitor *PerformanceMonitor) *QueryProfiler {
	return &QueryProfiler{
		monitor: monitor,
	}
}

// Profile profiles a query execution
func (qp *QueryProfiler) Profile(ctx context.Context, query string, fn func(context.Context) error) error {
	start := time.Now()
	err := fn(ctx)
	duration := time.Since(start)

	if qp.monitor != nil {
		qp.monitor.RecordQuery(query, duration)
	}

	return err
}

// BatchOptimizer optimizes batch operations
type BatchOptimizer struct {
	optimalBatchSize int
	maxBatchSize     int
	minBatchSize     int
}

// NewBatchOptimizer creates a new batch optimizer
func NewBatchOptimizer() *BatchOptimizer {
	return &BatchOptimizer{
		optimalBatchSize: 100,
		maxBatchSize:     1000,
		minBatchSize:     10,
	}
}

// OptimizeBatchSize optimizes batch size based on performance
func (bo *BatchOptimizer) OptimizeBatchSize(currentSize int, avgDuration time.Duration) int {
	// Simple optimization: adjust based on duration
	if avgDuration > 1*time.Second {
		// Too slow, reduce batch size
		newSize := currentSize / 2
		if newSize < bo.minBatchSize {
			return bo.minBatchSize
		}
		return newSize
	} else if avgDuration < 100*time.Millisecond {
		// Fast, can increase batch size
		newSize := currentSize * 2
		if newSize > bo.maxBatchSize {
			return bo.maxBatchSize
		}
		return newSize
	}

	return currentSize
}

// ConnectionPoolOptimizer optimizes connection pool settings
type ConnectionPoolOptimizer struct {
	currentMaxConns int32
	currentMinConns int32
}

// NewConnectionPoolOptimizer creates a new pool optimizer
func NewConnectionPoolOptimizer() *ConnectionPoolOptimizer {
	return &ConnectionPoolOptimizer{
		currentMaxConns: 25,
		currentMinConns: 5,
	}
}

// OptimizePoolSize optimizes pool size based on metrics
func (cpo *ConnectionPoolOptimizer) OptimizePoolSize(metrics HealthMetrics) (maxConns, minConns int32) {
	// Simple optimization based on usage
	usageRatio := float64(metrics.AcquiredConns) / float64(metrics.MaxConns)

	if usageRatio > 0.8 {
		// High usage, increase pool
		return cpo.currentMaxConns * 2, cpo.currentMinConns * 2
	} else if usageRatio < 0.3 {
		// Low usage, decrease pool
		newMax := cpo.currentMaxConns / 2
		if newMax < 10 {
			newMax = 10
		}
		return newMax, cpo.currentMinConns / 2
	}

	return cpo.currentMaxConns, cpo.currentMinConns
}

