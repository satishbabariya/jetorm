package core

import (
	"strings"
	"testing"
)

type TestEntity struct {
	ID    int64  `db:"id" validate:"required"`
	Email string `db:"email" validate:"required,email"`
	Name  string `db:"name" validate:"required,min:3"`
}

func TestValidator_Required(t *testing.T) {
	validator := NewValidator()
	validator.RegisterRule("Email", Required())

	entity := &TestEntity{
		ID:    1,
		Email: "", // Empty - should fail
		Name:  "Test",
	}

	err := validator.Validate(entity)
	if err == nil {
		t.Error("Validation should fail for empty email")
	}
}

func TestValidator_Email(t *testing.T) {
	validator := NewValidator()
	validator.RegisterRule("Email", Email())

	testCases := []struct {
		email string
		valid bool
	}{
		{"test@example.com", true},
		{"invalid-email", false},
	}

	for _, tc := range testCases {
		entity := &TestEntity{
			ID:    1, // Set required fields
			Name:  "Test",
			Email: tc.email,
		}
		err := validator.Validate(entity)
		// Check if error is specifically about email format
		if tc.valid {
			// Should not have email format error
			if err != nil && strings.Contains(err.Error(), "invalid email format") {
				t.Errorf("Email %s should be valid, got error: %v", tc.email, err)
			}
		} else {
			// Should have email format error
			if err == nil || !strings.Contains(err.Error(), "invalid email format") {
				t.Errorf("Email %s should be invalid", tc.email)
			}
		}
	}
}

func TestValidateEntity(t *testing.T) {
	entity := &TestEntity{
		ID:    1,
		Email: "test@example.com",
		Name:  "Test",
	}

	err := ValidateEntity(entity)
	if err != nil {
		t.Errorf("Valid entity should pass validation: %v", err)
	}
}

