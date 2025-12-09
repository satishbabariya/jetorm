package core

import (
	"context"
	"fmt"
	"time"
)

// HealthStatus represents the health status of a component
type HealthStatus string

const (
	HealthStatusUp      HealthStatus = "UP"
	HealthStatusDown    HealthStatus = "DOWN"
	HealthStatusUnknown HealthStatus = "UNKNOWN"
)

// HealthCheck represents a health check result
type HealthCheck struct {
	Status    HealthStatus
	Message   string
	Timestamp time.Time
	Details   map[string]interface{}
}

// HealthChecker checks the health of database connections
type HealthChecker struct {
	db *Database
}

// NewHealthChecker creates a new health checker
func NewHealthChecker(db *Database) *HealthChecker {
	return &HealthChecker{
		db: db,
	}
}

// Check performs a health check
func (hc *HealthChecker) Check(ctx context.Context) HealthCheck {
	check := HealthCheck{
		Status:    HealthStatusUnknown,
		Timestamp: time.Now(),
		Details:   make(map[string]interface{}),
	}

	// Check database connection
	if hc.db == nil || hc.db.pool == nil {
		check.Status = HealthStatusDown
		check.Message = "Database connection not initialized"
		return check
	}

	// Ping database
	pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	err := hc.db.pool.Ping(pingCtx)
	if err != nil {
		check.Status = HealthStatusDown
		check.Message = fmt.Sprintf("Database ping failed: %v", err)
		check.Details["error"] = err.Error()
		return check
	}

	// Get connection pool stats
	stats := hc.db.pool.Stat()
	check.Status = HealthStatusUp
	check.Message = "Database is healthy"
	check.Details["max_connections"] = stats.MaxConns()
	check.Details["acquired_connections"] = stats.AcquiredConns()
	check.Details["idle_connections"] = stats.IdleConns()
	check.Details["total_connections"] = stats.TotalConns()
	check.Details["constructing_connections"] = stats.ConstructingConns()

	return check
}

// CheckWithQuery performs a health check with a test query
func (hc *HealthChecker) CheckWithQuery(ctx context.Context, query string) HealthCheck {
	check := hc.Check(ctx)
	if check.Status != HealthStatusUp {
		return check
	}

	// Execute test query
	queryCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	start := time.Now()
	_, err := hc.db.pool.Query(queryCtx, query)
	duration := time.Since(start)

	if err != nil {
		check.Status = HealthStatusDown
		check.Message = fmt.Sprintf("Health check query failed: %v", err)
		check.Details["query_error"] = err.Error()
		return check
	}

	check.Details["query_duration"] = duration.String()
	check.Details["query_duration_ms"] = duration.Milliseconds()

	return check
}

// IsHealthy returns true if the database is healthy
func (hc *HealthChecker) IsHealthy(ctx context.Context) bool {
	check := hc.Check(ctx)
	return check.Status == HealthStatusUp
}

// HealthMetrics provides database health metrics
type HealthMetrics struct {
	TotalConns        int32
	AcquiredConns     int32
	IdleConns         int32
	MaxConns          int32
	ConstructingConns int32
	AcquireDuration   time.Duration
	AcquireCount      int64
	CanceledAcquireCount int64
	EmptyAcquireCount int64
}

// GetMetrics returns current database metrics
func (hc *HealthChecker) GetMetrics() HealthMetrics {
	if hc.db == nil || hc.db.pool == nil {
		return HealthMetrics{}
	}

	stats := hc.db.pool.Stat()
	return HealthMetrics{
		TotalConns:        stats.TotalConns(),
		AcquiredConns:     stats.AcquiredConns(),
		IdleConns:         stats.IdleConns(),
		MaxConns:          stats.MaxConns(),
		ConstructingConns: stats.ConstructingConns(),
		AcquireDuration:   stats.AcquireDuration(),
		AcquireCount:      stats.AcquireCount(),
		CanceledAcquireCount: stats.CanceledAcquireCount(),
		EmptyAcquireCount:    stats.EmptyAcquireCount(),
	}
}

// ConnectionHealth provides connection health information
type ConnectionHealth struct {
	Status      HealthStatus
	PoolStats   HealthMetrics
	LastCheck   time.Time
	Uptime      time.Duration
}

// GetConnectionHealth returns comprehensive connection health
func (hc *HealthChecker) GetConnectionHealth(ctx context.Context) ConnectionHealth {
	check := hc.Check(ctx)
	metrics := hc.GetMetrics()

	return ConnectionHealth{
		Status:    check.Status,
		PoolStats: metrics,
		LastCheck: check.Timestamp,
		Uptime:    time.Since(check.Timestamp), // Simplified
	}
}

