package core

import (
	"errors"
	"fmt"
)

// Error types with enhanced context
var (
	// Database errors
	ErrDatabaseConnection = errors.New("database connection error")
	ErrDatabaseQuery     = errors.New("database query error")
	ErrDatabaseExec      = errors.New("database execution error")
	
	// Entity errors
	ErrEntityNotFound    = errors.New("entity not found")
	ErrEntityInvalid     = errors.New("entity is invalid")
	ErrEntityDuplicate   = errors.New("entity already exists")
	
	// Validation errors
	ErrValidationFailed  = errors.New("validation failed")
	ErrInvalidInput      = errors.New("invalid input")
	
	// Transaction errors (ErrTransactionFailed already defined in errors.go)
	ErrTransactionRollback = errors.New("transaction rollback failed")
	
	// Query errors
	ErrQueryFailed       = errors.New("query failed")
	ErrQueryTimeout      = errors.New("query timeout")
	ErrQueryInvalid      = errors.New("invalid query")
	
	// Relationship errors
	ErrRelationshipNotFound = errors.New("relationship not found")
	ErrRelationshipInvalid  = errors.New("relationship is invalid")
)

// ErrorWithContext provides error with additional context
type ErrorWithContext struct {
	Err     error
	Message string
	Context map[string]interface{}
}

// Error implements error interface
func (e *ErrorWithContext) Error() string {
	if e.Message != "" {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Err.Error()
}

// Unwrap returns the underlying error
func (e *ErrorWithContext) Unwrap() error {
	return e.Err
}

// WithContext adds context to an error
func WithContext(err error, message string, context map[string]interface{}) error {
	return &ErrorWithContext{
		Err:     err,
		Message: message,
		Context: context,
	}
}

// WrapError wraps an error with a message
func WrapError(err error, message string) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("%s: %w", message, err)
}

// WrapErrorf wraps an error with a formatted message
func WrapErrorf(err error, format string, args ...interface{}) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("%s: %w", fmt.Sprintf(format, args...), err)
}

// IsNotFound checks if error is a not found error
func IsNotFound(err error) bool {
	return errors.Is(err, ErrNotFound) || errors.Is(err, ErrEntityNotFound)
}

// IsDuplicate checks if error is a duplicate error
func IsDuplicate(err error) bool {
	return errors.Is(err, ErrEntityDuplicate)
}

// IsValidationError checks if error is a validation error
func IsValidationError(err error) bool {
	return errors.Is(err, ErrValidationFailed) || errors.Is(err, ErrInvalidInput)
}

// IsTransactionError checks if error is a transaction error
func IsTransactionError(err error) bool {
	return errors.Is(err, ErrTransactionFailed) || errors.Is(err, ErrTransactionRollback)
}

// ErrorCode represents error codes for programmatic error handling
type ErrorCode string

const (
	ErrorCodeNotFound      ErrorCode = "NOT_FOUND"
	ErrorCodeDuplicate     ErrorCode = "DUPLICATE"
	ErrorCodeValidation    ErrorCode = "VALIDATION_ERROR"
	ErrorCodeDatabase      ErrorCode = "DATABASE_ERROR"
	ErrorCodeTransaction   ErrorCode = "TRANSACTION_ERROR"
	ErrorCodeQuery         ErrorCode = "QUERY_ERROR"
	ErrorCodeTimeout       ErrorCode = "TIMEOUT"
	ErrorCodeUnauthorized  ErrorCode = "UNAUTHORIZED"
	ErrorCodeInternal      ErrorCode = "INTERNAL_ERROR"
)

// CodedError provides error with error code
type CodedError struct {
	Code    ErrorCode
	Message string
	Err     error
}

// Error implements error interface
func (e *CodedError) Error() string {
	if e.Message != "" {
		return fmt.Sprintf("[%s] %s: %v", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("[%s] %v", e.Code, e.Err)
}

// Unwrap returns the underlying error
func (e *CodedError) Unwrap() error {
	return e.Err
}

// NewCodedError creates a new coded error
func NewCodedError(code ErrorCode, message string, err error) *CodedError {
	return &CodedError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}

// GetErrorCode extracts error code from error
func GetErrorCode(err error) ErrorCode {
	var codedErr *CodedError
	if errors.As(err, &codedErr) {
		return codedErr.Code
	}
	return ErrorCodeInternal
}

