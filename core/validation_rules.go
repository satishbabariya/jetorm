package core

import (
	"fmt"
	"regexp"
	"strings"
	"unicode"
)

// Additional validation rules

// MinLength validates minimum string length
func MinLength(min int) ValidationRule {
	return func(value interface{}) error {
		str, ok := value.(string)
		if !ok {
			return nil
		}
		if len(str) < min {
			return fmt.Errorf("must be at least %d characters", min)
		}
		return nil
	}
}

// MaxLength validates maximum string length
func MaxLength(max int) ValidationRule {
	return func(value interface{}) error {
		str, ok := value.(string)
		if !ok {
			return nil
		}
		if len(str) > max {
			return fmt.Errorf("must be at most %d characters", max)
		}
		return nil
	}
}

// Length validates exact string length
func Length(exact int) ValidationRule {
	return func(value interface{}) error {
		str, ok := value.(string)
		if !ok {
			return nil
		}
		if len(str) != exact {
			return fmt.Errorf("must be exactly %d characters", exact)
		}
		return nil
	}
}

// Range validates numeric range
func Range(min, max float64) ValidationRule {
	return func(value interface{}) error {
		var num float64
		switch v := value.(type) {
		case int:
			num = float64(v)
		case int8:
			num = float64(v)
		case int16:
			num = float64(v)
		case int32:
			num = float64(v)
		case int64:
			num = float64(v)
		case float32:
			num = float64(v)
		case float64:
			num = v
		default:
			return nil
		}
		if num < min || num > max {
			return fmt.Errorf("must be between %v and %v", min, max)
		}
		return nil
	}
}

// Pattern validates string against regex pattern
func Pattern(pattern string) ValidationRule {
	regex := regexp.MustCompile(pattern)
	return func(value interface{}) error {
		str, ok := value.(string)
		if !ok {
			return nil
		}
		if !regex.MatchString(str) {
			return fmt.Errorf("does not match pattern")
		}
		return nil
	}
}

// Alpha validates alphabetic characters only
func Alpha() ValidationRule {
	return Pattern("^[a-zA-Z]+$")
}

// Alphanumeric validates alphanumeric characters only
func Alphanumeric() ValidationRule {
	return Pattern("^[a-zA-Z0-9]+$")
}

// Numeric validates numeric characters only
func Numeric() ValidationRule {
	return Pattern("^[0-9]+$")
}

// Lowercase validates lowercase characters only
func Lowercase() ValidationRule {
	return func(value interface{}) error {
		str, ok := value.(string)
		if !ok {
			return nil
		}
		if str != strings.ToLower(str) {
			return fmt.Errorf("must be lowercase")
		}
		return nil
	}
}

// Uppercase validates uppercase characters only
func Uppercase() ValidationRule {
	return func(value interface{}) error {
		str, ok := value.(string)
		if !ok {
			return nil
		}
		if str != strings.ToUpper(str) {
			return fmt.Errorf("must be uppercase")
		}
		return nil
	}
}

// HasLetter validates that string contains at least one letter
func HasLetter() ValidationRule {
	return func(value interface{}) error {
		str, ok := value.(string)
		if !ok {
			return nil
		}
		for _, r := range str {
			if unicode.IsLetter(r) {
				return nil
			}
		}
		return fmt.Errorf("must contain at least one letter")
	}
}

// HasDigit validates that string contains at least one digit
func HasDigit() ValidationRule {
	return func(value interface{}) error {
		str, ok := value.(string)
		if !ok {
			return nil
		}
		for _, r := range str {
			if unicode.IsDigit(r) {
				return nil
			}
		}
		return fmt.Errorf("must contain at least one digit")
	}
}

// HasSpecialChar validates that string contains at least one special character
func HasSpecialChar() ValidationRule {
	return func(value interface{}) error {
		str, ok := value.(string)
		if !ok {
			return nil
		}
		for _, r := range str {
			if !unicode.IsLetter(r) && !unicode.IsDigit(r) && !unicode.IsSpace(r) {
				return nil
			}
		}
		return fmt.Errorf("must contain at least one special character")
	}
}

// InList validates that value is in a list of allowed values
func InList(allowed ...interface{}) ValidationRule {
	return func(value interface{}) error {
		for _, allowedValue := range allowed {
			if value == allowedValue {
				return nil
			}
		}
		return fmt.Errorf("must be one of: %v", allowed)
	}
}

// NotInList validates that value is not in a list of disallowed values
func NotInList(disallowed ...interface{}) ValidationRule {
	return func(value interface{}) error {
		for _, disallowedValue := range disallowed {
			if value == disallowedValue {
				return fmt.Errorf("must not be one of: %v", disallowed)
			}
		}
		return nil
	}
}

// Positive validates that number is positive
func Positive() ValidationRule {
	return func(value interface{}) error {
		var num float64
		switch v := value.(type) {
		case int:
			num = float64(v)
		case int8:
			num = float64(v)
		case int16:
			num = float64(v)
		case int32:
			num = float64(v)
		case int64:
			num = float64(v)
		case float32:
			num = float64(v)
		case float64:
			num = v
		default:
			return nil
		}
		if num <= 0 {
			return fmt.Errorf("must be positive")
		}
		return nil
	}
}

// Negative validates that number is negative
func Negative() ValidationRule {
	return func(value interface{}) error {
		var num float64
		switch v := value.(type) {
		case int:
			num = float64(v)
		case int8:
			num = float64(v)
		case int16:
			num = float64(v)
		case int32:
			num = float64(v)
		case int64:
			num = float64(v)
		case float32:
			num = float64(v)
		case float64:
			num = v
		default:
			return nil
		}
		if num >= 0 {
			return fmt.Errorf("must be negative")
		}
		return nil
	}
}

// NonZero validates that value is not zero
func NonZero() ValidationRule {
	return func(value interface{}) error {
		if isEmpty(value) {
			return fmt.Errorf("must not be zero")
		}
		return nil
	}
}

// Custom creates a custom validation rule
func Custom(fn func(interface{}) error) ValidationRule {
	return fn
}

// All validates that all rules pass
func All(rules ...ValidationRule) ValidationRule {
	return func(value interface{}) error {
		for _, rule := range rules {
			if err := rule(value); err != nil {
				return err
			}
		}
		return nil
	}
}

// AnyRule validates that at least one rule passes
func AnyRule(rules ...ValidationRule) ValidationRule {
	return func(value interface{}) error {
		var lastErr error
		for _, rule := range rules {
			if ruleErr := rule(value); ruleErr == nil {
				return nil
			} else {
				lastErr = ruleErr
			}
		}
		if lastErr != nil {
			return fmt.Errorf("none of the validation rules passed: %v", lastErr)
		}
		return nil
	}
}

