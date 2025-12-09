package core

import (
	"context"
	"database/sql"
)

// TransactionManager handles database transactions
type TransactionManager interface {
	// Transaction executes function in a transaction
	Transaction(ctx context.Context, fn func(tx *Tx) error) error

	// TransactionWithOptions executes with options
	TransactionWithOptions(ctx context.Context, opts TxOptions, fn func(tx *Tx) error) error

	// Begin starts a manual transaction
	Begin(ctx context.Context) (*Tx, error)
	BeginWithOptions(ctx context.Context, opts TxOptions) (*Tx, error)
}

// TxOptions represents transaction options
type TxOptions struct {
	Isolation  IsolationLevel // Transaction isolation level
	ReadOnly   bool           // Read-only transaction
	Deferrable bool           // Deferrable constraint checking
}

// IsolationLevel represents transaction isolation level
type IsolationLevel int

const (
	ReadUncommitted IsolationLevel = iota
	ReadCommitted
	RepeatableRead
	Serializable
)

// ToSQLIsolation converts IsolationLevel to sql.IsolationLevel
func (l IsolationLevel) ToSQLIsolation() sql.IsolationLevel {
	switch l {
	case ReadUncommitted:
		return sql.LevelReadUncommitted
	case ReadCommitted:
		return sql.LevelReadCommitted
	case RepeatableRead:
		return sql.LevelRepeatableRead
	case Serializable:
		return sql.LevelSerializable
	default:
		return sql.LevelDefault
	}
}

// Tx represents a database transaction
type Tx struct {
	ctx context.Context
	tx  interface{} // Underlying transaction (pgx or sql.Tx)
}

// Commit commits the transaction
func (t *Tx) Commit() error {
	// Implementation will depend on the underlying driver
	return nil
}

// Rollback rolls back the transaction
func (t *Tx) Rollback() error {
	// Implementation will depend on the underlying driver
	return nil
}

// Context returns the transaction context
func (t *Tx) Context() context.Context {
	return t.ctx
}

