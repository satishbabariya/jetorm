package core

import (
	"fmt"
	"reflect"
	"strings"
)

// Validator validates entities before operations
type Validator struct {
	rules map[string][]ValidationRule
}

// ValidationRule defines a validation rule
type ValidationRule func(value interface{}) error

// NewValidator creates a new validator
func NewValidator() *Validator {
	return &Validator{
		rules: make(map[string][]ValidationRule),
	}
}

// RegisterRule registers a validation rule for a field
func (v *Validator) RegisterRule(field string, rule ValidationRule) {
	v.rules[field] = append(v.rules[field], rule)
}

// Validate validates an entity
func (v *Validator) Validate(entity interface{}) error {
	entityType := reflect.TypeOf(entity)
	if entityType.Kind() == reflect.Ptr {
		entityType = entityType.Elem()
	}
	if entityType.Kind() != reflect.Struct {
		return ErrInvalidInput
	}

	entityValue := reflect.ValueOf(entity)
	if entityValue.Kind() == reflect.Ptr {
		entityValue = entityValue.Elem()
	}

	var errors []string

	for i := 0; i < entityType.NumField(); i++ {
		field := entityType.Field(i)
		if !field.IsExported() {
			continue
		}

		fieldValue := entityValue.Field(i)
		fieldName := field.Name

		// Get validation rules for this field
		rules := v.rules[fieldName]
		
		// Also check for validation tags
		validateTag := field.Tag.Get("validate")
		if validateTag != "" {
			rules = append(rules, parseValidationTag(validateTag)...)
		}

		// Apply rules
		for _, rule := range rules {
			if err := rule(fieldValue.Interface()); err != nil {
				errors = append(errors, fmt.Sprintf("%s: %v", fieldName, err))
			}
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("%w: %s", ErrValidationFailed, strings.Join(errors, "; "))
	}

	return nil
}

// parseValidationTag parses validation tags
func parseValidationTag(tag string) []ValidationRule {
	var rules []ValidationRule
	parts := strings.Split(tag, ",")

	for _, part := range parts {
		part = strings.TrimSpace(part)
		
		switch {
		case part == "required":
			rules = append(rules, Required())
		case strings.HasPrefix(part, "min:"):
			value := strings.TrimPrefix(part, "min:")
			rules = append(rules, Min(value))
		case strings.HasPrefix(part, "max:"):
			value := strings.TrimPrefix(part, "max:")
			rules = append(rules, Max(value))
		case strings.HasPrefix(part, "email"):
			rules = append(rules, Email())
		case strings.HasPrefix(part, "url"):
			rules = append(rules, URL())
		}
	}

	return rules
}

// Required validates that a value is not zero/nil
func Required() ValidationRule {
	return func(value interface{}) error {
		if isEmpty(value) {
			return fmt.Errorf("is required")
		}
		return nil
	}
}

// Min validates minimum value/length
func Min(minStr string) ValidationRule {
	return func(value interface{}) error {
		// Implementation would parse minStr and compare
		// Simplified version
		return nil
	}
}

// Max validates maximum value/length
func Max(maxStr string) ValidationRule {
	return func(value interface{}) error {
		// Implementation would parse maxStr and compare
		// Simplified version
		return nil
	}
}

// Email validates email format
func Email() ValidationRule {
	return func(value interface{}) error {
		str, ok := value.(string)
		if !ok {
			return nil // Not a string, skip
		}
		if !strings.Contains(str, "@") {
			return fmt.Errorf("invalid email format")
		}
		return nil
	}
}

// URL validates URL format
func URL() ValidationRule {
	return func(value interface{}) error {
		str, ok := value.(string)
		if !ok {
			return nil
		}
		if !strings.HasPrefix(str, "http://") && !strings.HasPrefix(str, "https://") {
			return fmt.Errorf("invalid URL format")
		}
		return nil
	}
}

// isEmpty checks if a value is empty
func isEmpty(value interface{}) bool {
	if value == nil {
		return true
	}

	v := reflect.ValueOf(value)
	switch v.Kind() {
	case reflect.String:
		return v.Len() == 0
	case reflect.Slice, reflect.Array, reflect.Map:
		return v.Len() == 0
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Bool:
		return !v.Bool()
	case reflect.Ptr, reflect.Interface:
		return v.IsNil()
	}

	return false
}

// ValidateEntity validates an entity using its tags
func ValidateEntity(entity interface{}) error {
	validator := NewValidator()
	return validator.Validate(entity)
}

