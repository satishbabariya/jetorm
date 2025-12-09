package logging

import (
	"context"
	"fmt"
	"log/slog"
	"time"
)

// SQLLogger logs SQL queries and their execution details
type SQLLogger struct {
	logger    *slog.Logger
	logSlow   bool
	slowThreshold time.Duration
}

// NewSQLLogger creates a new SQL logger
func NewSQLLogger(logger *slog.Logger) *SQLLogger {
	return &SQLLogger{
		logger:        logger,
		logSlow:       true,
		slowThreshold: 100 * time.Millisecond,
	}
}

// SetSlowThreshold sets the threshold for slow query logging
func (sl *SQLLogger) SetSlowThreshold(threshold time.Duration) {
	sl.slowThreshold = threshold
}

// LogQuery logs a SQL query
func (sl *SQLLogger) LogQuery(ctx context.Context, query string, args []interface{}, duration time.Duration) {
	attrs := []any{
		slog.String("query", query),
		slog.Duration("duration", duration),
	}
	
	if len(args) > 0 {
		attrs = append(attrs, slog.Any("args", args))
	}
	
	if sl.logSlow && duration > sl.slowThreshold {
		sl.logger.Warn("Slow query detected", slog.Group("sql", attrs...))
	} else {
		sl.logger.Debug("SQL query executed", slog.Group("sql", attrs...))
	}
}

// LogError logs a SQL error
func (sl *SQLLogger) LogError(ctx context.Context, query string, err error) {
	sl.logger.Error("SQL query error",
		slog.String("query", query),
		slog.String("error", err.Error()),
	)
}

// LogTransaction logs transaction events
func (sl *SQLLogger) LogTransaction(ctx context.Context, event string, txID string) {
	sl.logger.Info("Transaction event",
		slog.String("event", event),
		slog.String("tx_id", txID),
	)
}

// FormatQuery formats a query with arguments for logging
func FormatQuery(query string, args []interface{}) string {
	if len(args) == 0 {
		return query
	}
	
	formatted := query
	for i, arg := range args {
		placeholder := fmt.Sprintf("$%d", i+1)
		value := fmt.Sprintf("%v", arg)
		formatted = fmt.Sprintf("%s -- %s: %s", formatted, placeholder, value)
	}
	
	return formatted
}

