package tx

import (
	"context"
	"database/sql"
	"fmt"
)

// Propagation defines transaction propagation behavior
type Propagation int

const (
	// PropagationRequired requires a transaction, creates one if none exists
	PropagationRequired Propagation = iota
	
	// PropagationRequiresNew always creates a new transaction
	PropagationRequiresNew
	
	// PropagationSupports supports a transaction but doesn't require one
	PropagationSupports
	
	// PropagationNotSupported doesn't support transactions, suspends if one exists
	PropagationNotSupported
	
	// PropagationNever never uses transactions, throws error if one exists
	PropagationNever
	
	// PropagationMandatory requires an existing transaction, throws error if none exists
	PropagationMandatory
)

// TransactionManager manages transactions with propagation support
type TransactionManager struct {
	db *sql.DB
}

// NewTransactionManager creates a new transaction manager
func NewTransactionManager(db *sql.DB) *TransactionManager {
	return &TransactionManager{
		db: db,
	}
}

// Execute executes a function within a transaction based on propagation
func (tm *TransactionManager) Execute(ctx context.Context, propagation Propagation, fn func(*sql.Tx) error) error {
	existingTx := getTxFromContext(ctx)
	
	switch propagation {
	case PropagationRequired:
		if existingTx != nil {
			return fn(existingTx)
		}
		return tm.executeInNewTx(ctx, fn)
		
	case PropagationRequiresNew:
		return tm.executeInNewTx(ctx, fn)
		
	case PropagationSupports:
		if existingTx != nil {
			return fn(existingTx)
		}
		return fn(nil)
		
	case PropagationNotSupported:
		if existingTx != nil {
			// Suspend transaction - execute without transaction
			return fn(nil)
		}
		return fn(nil)
		
	case PropagationNever:
		if existingTx != nil {
			return fmt.Errorf("transaction propagation NEVER: existing transaction found")
		}
		return fn(nil)
		
	case PropagationMandatory:
		if existingTx == nil {
			return fmt.Errorf("transaction propagation MANDATORY: no existing transaction found")
		}
		return fn(existingTx)
		
	default:
		return fmt.Errorf("unknown propagation: %d", propagation)
	}
}

// executeInNewTx executes a function in a new transaction
func (tm *TransactionManager) executeInNewTx(ctx context.Context, fn func(*sql.Tx) error) error {
	tx, err := tm.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		} else if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()
	
	err = fn(tx)
	return err
}

// Context key for storing transaction
type txKey struct{}

// WithTx adds a transaction to context
func WithTx(ctx context.Context, tx *sql.Tx) context.Context {
	return context.WithValue(ctx, txKey{}, tx)
}

// getTxFromContext retrieves transaction from context
func getTxFromContext(ctx context.Context) *sql.Tx {
	if tx, ok := ctx.Value(txKey{}).(*sql.Tx); ok {
		return tx
	}
	return nil
}

