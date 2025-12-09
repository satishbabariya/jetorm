package core

import "time"

// Config holds database configuration
type Config struct {
	// Connection
	Driver   string // "pgx" (default), future: "mysql", "sqlite"
	Host     string // Database host
	Port     int    // Database port
	Database string // Database name
	User     string // Database user
	Password string // Database password
	SSLMode  string // SSL mode: disable, require, verify-ca, verify-full

	// Connection Pool
	MaxOpenConns    int           // Maximum open connections (default: 25)
	MaxIdleConns    int           // Maximum idle connections (default: 5)
	ConnMaxLifetime time.Duration // Connection max lifetime (default: 5m)
	ConnMaxIdleTime time.Duration // Connection max idle time (default: 5m)

	// Migrations
	MigrationsPath string // Path to migration files
	AutoMigrate    bool   // Auto-run migrations on startup
	MigrationTable string // Migration version table (default: "schema_migrations")

	// Jet Code Generation
	JetGenPath    string // Path for generated Jet code
	JetGenPackage string // Package name for Jet code

	// Logging
	Logger         Logger        // Custom logger implementation
	LogLevel       LogLevel      // Log level: Debug, Info, Warn, Error
	LogSQL         bool          // Log SQL queries
	LogSlowQueries time.Duration // Log queries slower than threshold

	// Performance
	PreparedStmts bool          // Use prepared statements (default: true)
	QueryTimeout  time.Duration // Default query timeout (default: 30s)

	// Behavior
	SoftDelete     bool   // Enable soft delete globally
	CreatedAtField string // Custom created_at field name
	UpdatedAtField string // Custom updated_at field name
	DeletedAtField string // Custom deleted_at field name
}

// DefaultConfig returns a Config with sensible defaults
func DefaultConfig() Config {
	return Config{
		Driver:          "pgx",
		Host:            "localhost",
		Port:            5432,
		SSLMode:         "disable",
		MaxOpenConns:    25,
		MaxIdleConns:    5,
		ConnMaxLifetime: 5 * time.Minute,
		ConnMaxIdleTime: 5 * time.Minute,
		MigrationTable:  "schema_migrations",
		PreparedStmts:   true,
		QueryTimeout:    30 * time.Second,
		LogLevel:        InfoLevel,
		CreatedAtField:  "created_at",
		UpdatedAtField:  "updated_at",
		DeletedAtField:  "deleted_at",
	}
}

// LogLevel represents logging level
type LogLevel int

const (
	DebugLevel LogLevel = iota
	InfoLevel
	WarnLevel
	ErrorLevel
)

// Logger interface for custom logging implementations
type Logger interface {
	Debug(msg string, args ...interface{})
	Info(msg string, args ...interface{})
	Warn(msg string, args ...interface{})
	Error(msg string, args ...interface{})
}

