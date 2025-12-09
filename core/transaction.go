package core

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
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
	Timeout    time.Duration  // Transaction timeout
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
	ctx      context.Context
	tx       pgx.Tx
	savepoints map[string]bool // Track savepoints
}

// Commit commits the transaction
func (t *Tx) Commit() error {
	if t.tx == nil {
		return fmt.Errorf("transaction is nil")
	}
	return t.tx.Commit(t.ctx)
}

// Rollback rolls back the transaction
func (t *Tx) Rollback() error {
	if t.tx == nil {
		return fmt.Errorf("transaction is nil")
	}
	return t.tx.Rollback(t.ctx)
}

// SavePoint creates a savepoint with the given name
func (t *Tx) SavePoint(name string) error {
	if t.tx == nil {
		return fmt.Errorf("transaction is nil")
	}
	if t.savepoints == nil {
		t.savepoints = make(map[string]bool)
	}
	
	query := fmt.Sprintf("SAVEPOINT %s", name)
	_, err := t.tx.Exec(t.ctx, query)
	if err != nil {
		return fmt.Errorf("failed to create savepoint %s: %w", name, err)
	}
	
	t.savepoints[name] = true
	return nil
}

// RollbackTo rolls back to a specific savepoint
func (t *Tx) RollbackTo(name string) error {
	if t.tx == nil {
		return fmt.Errorf("transaction is nil")
	}
	if t.savepoints == nil || !t.savepoints[name] {
		return fmt.Errorf("savepoint %s does not exist", name)
	}
	
	query := fmt.Sprintf("ROLLBACK TO SAVEPOINT %s", name)
	_, err := t.tx.Exec(t.ctx, query)
	if err != nil {
		return fmt.Errorf("failed to rollback to savepoint %s: %w", name, err)
	}
	
	return nil
}

// ReleaseSavePoint releases a savepoint
func (t *Tx) ReleaseSavePoint(name string) error {
	if t.tx == nil {
		return fmt.Errorf("transaction is nil")
	}
	if t.savepoints == nil || !t.savepoints[name] {
		return fmt.Errorf("savepoint %s does not exist", name)
	}
	
	query := fmt.Sprintf("RELEASE SAVEPOINT %s", name)
	_, err := t.tx.Exec(t.ctx, query)
	if err != nil {
		return fmt.Errorf("failed to release savepoint %s: %w", name, err)
	}
	
	delete(t.savepoints, name)
	return nil
}

// Context returns the transaction context
func (t *Tx) Context() context.Context {
	return t.ctx
}

