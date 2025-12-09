package core

import (
	"context"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Database represents the database connection
type Database struct {
	pool   *pgxpool.Pool
	config Config
	logger Logger
}

// Connect creates a new database connection
func Connect(config Config) (*Database, error) {
	// Apply defaults
	if config.MaxOpenConns == 0 {
		config.MaxOpenConns = 25
	}
	if config.MaxIdleConns == 0 {
		config.MaxIdleConns = 5
	}
	if config.ConnMaxLifetime == 0 {
		config.ConnMaxLifetime = 5 * time.Minute
	}
	if config.ConnMaxIdleTime == 0 {
		config.ConnMaxIdleTime = 5 * time.Minute
	}
	if config.QueryTimeout == 0 {
		config.QueryTimeout = 30 * time.Second
	}

	// Build connection string
	connString := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		config.Host,
		config.Port,
		config.User,
		config.Password,
		config.Database,
		config.SSLMode,
	)

	// Create pool config
	poolConfig, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrConnectionFailed, err)
	}

	// Configure pool
	poolConfig.MaxConns = int32(config.MaxOpenConns)
	poolConfig.MinConns = int32(config.MaxIdleConns)
	poolConfig.MaxConnLifetime = config.ConnMaxLifetime
	poolConfig.MaxConnIdleTime = config.ConnMaxIdleTime

	// Create pool
	pool, err := pgxpool.NewWithConfig(context.Background(), poolConfig)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrConnectionFailed, err)
	}

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("%w: %v", ErrConnectionFailed, err)
	}

	db := &Database{
		pool:   pool,
		config: config,
		logger: config.Logger,
	}

	// Initialize default logger if none provided
	if db.logger == nil {
		db.logger = &defaultLogger{level: config.LogLevel}
	}

	db.logger.Info("database connection established", "host", config.Host, "database", config.Database)

	return db, nil
}

// MustConnect creates a new database connection and panics on error
func MustConnect(config Config) *Database {
	db, err := Connect(config)
	if err != nil {
		panic(err)
	}
	return db
}

// ConnectURL creates a database connection from a connection string
func ConnectURL(connString string, opts ...ConfigOption) (*Database, error) {
	// Parse connection string
	parsedURL, err := url.Parse(connString)
	if err != nil {
		return nil, fmt.Errorf("%w: invalid connection string: %v", ErrConnectionFailed, err)
	}

	// Extract components
	config := DefaultConfig()
	config.Host = parsedURL.Hostname()
	if port := parsedURL.Port(); port != "" {
		fmt.Sscanf(port, "%d", &config.Port)
	}
	config.Database = strings.TrimPrefix(parsedURL.Path, "/")
	config.User = parsedURL.User.Username()
	if password, ok := parsedURL.User.Password(); ok {
		config.Password = password
	}

	// Parse query parameters
	query := parsedURL.Query()
	if sslMode := query.Get("sslmode"); sslMode != "" {
		config.SSLMode = sslMode
	}

	// Apply additional options
	for _, opt := range opts {
		opt(&config)
	}

	return Connect(config)
}

// ConfigOption is a function that modifies Config
type ConfigOption func(*Config)

// WithMaxOpenConns sets the maximum open connections
func WithMaxOpenConns(n int) ConfigOption {
	return func(c *Config) {
		c.MaxOpenConns = n
	}
}

// WithMaxIdleConns sets the maximum idle connections
func WithMaxIdleConns(n int) ConfigOption {
	return func(c *Config) {
		c.MaxIdleConns = n
	}
}

// WithLogger sets a custom logger
func WithLogger(logger Logger) ConfigOption {
	return func(c *Config) {
		c.Logger = logger
	}
}

// WithLogSQL enables SQL logging
func WithLogSQL(enabled bool) ConfigOption {
	return func(c *Config) {
		c.LogSQL = enabled
	}
}

// Close closes the database connection
func (db *Database) Close() {
	if db.pool != nil {
		db.pool.Close()
		db.logger.Info("database connection closed")
	}
}

// Pool returns the underlying connection pool
func (db *Database) Pool() *pgxpool.Pool {
	return db.pool
}

// Ping checks if the database is reachable
func (db *Database) Ping(ctx context.Context) error {
	return db.pool.Ping(ctx)
}

// Transaction executes a function within a transaction
func (db *Database) Transaction(ctx context.Context, fn func(tx *Tx) error) error {
	return db.TransactionWithOptions(ctx, TxOptions{}, fn)
}

// TransactionWithOptions executes a function within a transaction with options
func (db *Database) TransactionWithOptions(ctx context.Context, opts TxOptions, fn func(tx *Tx) error) error {
	// Apply timeout if specified
	if opts.Timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, opts.Timeout)
		defer cancel()
	}

	// Begin transaction
	pgxTx, err := db.pool.BeginTx(ctx, pgx.TxOptions{
		IsoLevel:   pgx.TxIsoLevel(opts.Isolation.ToSQLIsolation().String()),
		AccessMode: func() pgx.TxAccessMode {
			if opts.ReadOnly {
				return pgx.ReadOnly
			}
			return pgx.ReadWrite
		}(),
		DeferrableMode: func() pgx.TxDeferrableMode {
			if opts.Deferrable {
				return pgx.Deferrable
			}
			return pgx.NotDeferrable
		}(),
	})
	if err != nil {
		return fmt.Errorf("%w: %v", ErrTransactionFailed, err)
	}

	tx := &Tx{
		ctx:        ctx,
		tx:         pgxTx,
		savepoints: make(map[string]bool),
	}

	// Execute function
	if err := fn(tx); err != nil {
		if rbErr := pgxTx.Rollback(ctx); rbErr != nil {
			db.logger.Error("failed to rollback transaction", "error", rbErr)
		}
		return err
	}

	// Commit transaction
	if err := pgxTx.Commit(ctx); err != nil {
		return fmt.Errorf("%w: %v", ErrTransactionFailed, err)
	}

	return nil
}

// Begin starts a new transaction
func (db *Database) Begin(ctx context.Context) (*Tx, error) {
	return db.BeginWithOptions(ctx, TxOptions{})
}

// BeginWithOptions starts a new transaction with options
func (db *Database) BeginWithOptions(ctx context.Context, opts TxOptions) (*Tx, error) {
	pgxTx, err := db.pool.BeginTx(ctx, pgx.TxOptions{
		IsoLevel:   pgx.TxIsoLevel(opts.Isolation.ToSQLIsolation().String()),
		AccessMode: func() pgx.TxAccessMode {
			if opts.ReadOnly {
				return pgx.ReadOnly
			}
			return pgx.ReadWrite
		}(),
		DeferrableMode: func() pgx.TxDeferrableMode {
			if opts.Deferrable {
				return pgx.Deferrable
			}
			return pgx.NotDeferrable
		}(),
	})
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrTransactionFailed, err)
	}

	return &Tx{
		ctx:        ctx,
		tx:         pgxTx,
		savepoints: make(map[string]bool),
	}, nil
}

// Config returns the database configuration
func (db *Database) Config() Config {
	return db.config
}

// Logger returns the database logger
func (db *Database) Logger() Logger {
	return db.logger
}

// defaultLogger is a simple default logger implementation
type defaultLogger struct {
	level LogLevel
}

func (l *defaultLogger) Debug(msg string, args ...interface{}) {
	if l.level <= DebugLevel {
		fmt.Printf("[DEBUG] %s %v\n", msg, args)
	}
}

func (l *defaultLogger) Info(msg string, args ...interface{}) {
	if l.level <= InfoLevel {
		fmt.Printf("[INFO] %s %v\n", msg, args)
	}
}

func (l *defaultLogger) Warn(msg string, args ...interface{}) {
	if l.level <= WarnLevel {
		fmt.Printf("[WARN] %s %v\n", msg, args)
	}
}

func (l *defaultLogger) Error(msg string, args ...interface{}) {
	if l.level <= ErrorLevel {
		fmt.Printf("[ERROR] %s %v\n", msg, args)
	}
}

