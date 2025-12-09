package core

import "errors"

var (
	// ErrInvalidEntity is returned when an entity type is invalid
	ErrInvalidEntity = errors.New("jetorm: invalid entity type, must be a struct")
	
	// ErrNotFound is returned when a record is not found
	ErrNotFound = errors.New("jetorm: record not found")
	
	// ErrInvalidID is returned when an ID is invalid
	ErrInvalidID = errors.New("jetorm: invalid ID")
	
	// ErrNoPrimaryKey is returned when an entity has no primary key
	ErrNoPrimaryKey = errors.New("jetorm: entity has no primary key")
	
	// ErrConnectionFailed is returned when database connection fails
	ErrConnectionFailed = errors.New("jetorm: database connection failed")
	
	// ErrTransactionFailed is returned when a transaction fails
	ErrTransactionFailed = errors.New("jetorm: transaction failed")
)

