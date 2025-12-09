package core

import (
	"errors"
	"testing"
)

func TestErrorWrapping(t *testing.T) {
	originalErr := ErrNotFound
	wrapped := WrapError(originalErr, "failed to find user")

	if !errors.Is(wrapped, ErrNotFound) {
		t.Error("Wrapped error should be unwrappable to original")
	}

	if !IsNotFound(wrapped) {
		t.Error("IsNotFound should detect wrapped error")
	}
}

func TestCodedError(t *testing.T) {
	err := NewCodedError(ErrorCodeNotFound, "User not found", ErrNotFound)

	if err.Code != ErrorCodeNotFound {
		t.Errorf("Expected code %s, got %s", ErrorCodeNotFound, err.Code)
	}

	code := GetErrorCode(err)
	if code != ErrorCodeNotFound {
		t.Errorf("Expected code %s, got %s", ErrorCodeNotFound, code)
	}
}

func TestErrorWithContext(t *testing.T) {
	originalErr := ErrNotFound
	errCtx := ErrorContext{
		Operation: "find_user",
		EntityID:  123,
	}

	err := WithErrorContext(originalErr, errCtx)

	if !errors.Is(err, ErrNotFound) {
		t.Error("ErrorWithContext should be unwrappable")
	}
}

