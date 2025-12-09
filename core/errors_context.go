package core

import (
	"fmt"
	"runtime"
	"strings"
)

// ErrorContext provides additional context for errors
type ErrorContext struct {
	Operation   string
	EntityType  string
	EntityID    interface{}
	Field       string
	Value       interface{}
	Query       string
	Args        []interface{}
	Stack       []string
	UserMessage string
}

// ContextualError provides error with rich context
type ContextualError struct {
	Err     error
	Context ErrorContext
}

// Error implements error interface
func (e *ContextualError) Error() string {
	var parts []string

	if e.Context.Operation != "" {
		parts = append(parts, fmt.Sprintf("operation: %s", e.Context.Operation))
	}

	if e.Context.EntityType != "" {
		parts = append(parts, fmt.Sprintf("entity: %s", e.Context.EntityType))
	}

	if e.Context.EntityID != nil {
		parts = append(parts, fmt.Sprintf("id: %v", e.Context.EntityID))
	}

	if e.Context.Field != "" {
		parts = append(parts, fmt.Sprintf("field: %s", e.Context.Field))
	}

	if e.Context.Query != "" {
		parts = append(parts, fmt.Sprintf("query: %s", e.Context.Query))
	}

	msg := e.Err.Error()
	if len(parts) > 0 {
		msg = fmt.Sprintf("%s (%s)", msg, strings.Join(parts, ", "))
	}

	return msg
}

// Unwrap returns the underlying error
func (e *ContextualError) Unwrap() error {
	return e.Err
}

// WithErrorContext creates an error with context
func WithErrorContext(err error, ctx ErrorContext) error {
	if err == nil {
		return nil
	}

	// Capture stack trace
	if len(ctx.Stack) == 0 {
		ctx.Stack = captureStack(3)
	}

	return &ContextualError{
		Err:     err,
		Context: ctx,
	}
}

// captureStack captures stack trace
func captureStack(skip int) []string {
	var stack []string
	for i := skip; i < skip+10; i++ {
		pc, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}
		fn := runtime.FuncForPC(pc)
		if fn != nil {
			stack = append(stack, fmt.Sprintf("%s:%d %s", file, line, fn.Name()))
		}
	}
	return stack
}

// NewOperationError creates an error for a specific operation
func NewOperationError(operation string, err error) error {
	return WithErrorContext(err, ErrorContext{
		Operation: operation,
	})
}

// NewEntityError creates an error for entity operations
func NewEntityError(entityType string, entityID interface{}, err error) error {
	return WithErrorContext(err, ErrorContext{
		Operation:  "entity_operation",
		EntityType: entityType,
		EntityID:   entityID,
	})
}

// NewFieldError creates an error for field validation
func NewFieldError(field string, value interface{}, err error) error {
	return WithErrorContext(err, ErrorContext{
		Operation: "field_validation",
		Field:     field,
		Value:     value,
	})
}

// NewQueryError creates an error for query operations
func NewQueryError(query string, args []interface{}, err error) error {
	return WithErrorContext(err, ErrorContext{
		Operation: "query_execution",
		Query:     query,
		Args:      args,
	})
}

// FormatError formats an error with user-friendly message
func FormatError(err error) string {
	if err == nil {
		return ""
	}

	var contextualErr *ContextualError
	if As(err, &contextualErr) {
		if contextualErr.Context.UserMessage != "" {
			return contextualErr.Context.UserMessage
		}
	}

	return err.Error()
}

// As checks if error can be unwrapped to target type
func As(err error, target interface{}) bool {
	// Simplified version - would use errors.As in production
	return false
}

// ErrorFormatter formats errors for display
type ErrorFormatter struct {
	IncludeStack bool
	IncludeQuery bool
}

// Format formats an error
func (ef *ErrorFormatter) Format(err error) string {
	if err == nil {
		return ""
	}

	var contextualErr *ContextualError
	if As(err, &contextualErr) {
		return ef.formatContextualError(contextualErr)
	}

	return err.Error()
}

// formatContextualError formats a contextual error
func (ef *ErrorFormatter) formatContextualError(err *ContextualError) string {
	var parts []string

	parts = append(parts, err.Err.Error())

	if err.Context.Operation != "" {
		parts = append(parts, fmt.Sprintf("Operation: %s", err.Context.Operation))
	}

	if err.Context.EntityType != "" {
		parts = append(parts, fmt.Sprintf("Entity: %s", err.Context.EntityType))
	}

	if err.Context.EntityID != nil {
		parts = append(parts, fmt.Sprintf("ID: %v", err.Context.EntityID))
	}

	if ef.IncludeQuery && err.Context.Query != "" {
		parts = append(parts, fmt.Sprintf("Query: %s", err.Context.Query))
	}

	if ef.IncludeStack && len(err.Context.Stack) > 0 {
		parts = append(parts, fmt.Sprintf("Stack: %s", strings.Join(err.Context.Stack, "\n")))
	}

	return strings.Join(parts, "\n")
}

