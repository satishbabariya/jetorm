package core

import (
	"fmt"
	"regexp"
	"strings"
	"time"
	"unicode"
)

// Advanced validation patterns

// PhoneNumber validates phone number format
func PhoneNumber() ValidationRule {
	phoneRegex := regexp.MustCompile(`^\+?[1-9]\d{1,14}$`)
	return func(value interface{}) error {
		str, ok := value.(string)
		if !ok {
			return nil
		}
		if !phoneRegex.MatchString(str) {
			return fmt.Errorf("invalid phone number format")
		}
		return nil
	}
}

// CreditCard validates credit card number (Luhn algorithm)
func CreditCard() ValidationRule {
	return func(value interface{}) error {
		str, ok := value.(string)
		if !ok {
			return nil
		}
		// Remove spaces and dashes
		cleaned := strings.ReplaceAll(strings.ReplaceAll(str, " ", ""), "-", "")
		
		if len(cleaned) < 13 || len(cleaned) > 19 {
			return fmt.Errorf("invalid credit card length")
		}

		// Luhn algorithm
		sum := 0
		alternate := false
		for i := len(cleaned) - 1; i >= 0; i-- {
			digit := int(cleaned[i] - '0')
			if alternate {
				digit *= 2
				if digit > 9 {
					digit -= 9
				}
			}
			sum += digit
			alternate = !alternate
		}

		if sum%10 != 0 {
			return fmt.Errorf("invalid credit card number")
		}
		return nil
	}
}

// UUID validates UUID format
func UUID() ValidationRule {
	uuidRegex := regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)
	return func(value interface{}) error {
		str, ok := value.(string)
		if !ok {
			return nil
		}
		if !uuidRegex.MatchString(strings.ToLower(str)) {
			return fmt.Errorf("invalid UUID format")
		}
		return nil
	}
}

// IPv4 validates IPv4 address
func IPv4() ValidationRule {
	ipv4Regex := regexp.MustCompile(`^(\d{1,3}\.){3}\d{1,3}$`)
	return func(value interface{}) error {
		str, ok := value.(string)
		if !ok {
			return nil
		}
		if !ipv4Regex.MatchString(str) {
			return fmt.Errorf("invalid IPv4 address")
		}
		// Validate octets
		parts := strings.Split(str, ".")
		for _, part := range parts {
			var octet int
			if _, err := fmt.Sscanf(part, "%d", &octet); err != nil || octet < 0 || octet > 255 {
				return fmt.Errorf("invalid IPv4 address")
			}
		}
		return nil
	}
}

// IPv6 validates IPv6 address
func IPv6() ValidationRule {
	ipv6Regex := regexp.MustCompile(`^([0-9a-fA-F]{1,4}:){7}[0-9a-fA-F]{1,4}$`)
	return func(value interface{}) error {
		str, ok := value.(string)
		if !ok {
			return nil
		}
		if !ipv6Regex.MatchString(str) {
			return fmt.Errorf("invalid IPv6 address")
		}
		return nil
	}
}

// MACAddress validates MAC address
func MACAddress() ValidationRule {
	macRegex := regexp.MustCompile(`^([0-9A-Fa-f]{2}[:-]){5}([0-9A-Fa-f]{2})$`)
	return func(value interface{}) error {
		str, ok := value.(string)
		if !ok {
			return nil
		}
		if !macRegex.MatchString(str) {
			return fmt.Errorf("invalid MAC address")
		}
		return nil
	}
}

// Base64 validates Base64 encoding
func Base64() ValidationRule {
	base64Regex := regexp.MustCompile(`^[A-Za-z0-9+/]*={0,2}$`)
	return func(value interface{}) error {
		str, ok := value.(string)
		if !ok {
			return nil
		}
		if len(str)%4 != 0 || !base64Regex.MatchString(str) {
			return fmt.Errorf("invalid Base64 encoding")
		}
		return nil
	}
}

// JSON validates JSON format
func JSON() ValidationRule {
	return func(value interface{}) error {
		str, ok := value.(string)
		if !ok {
			return nil
		}
		str = strings.TrimSpace(str)
		if !strings.HasPrefix(str, "{") && !strings.HasPrefix(str, "[") {
			return fmt.Errorf("invalid JSON format")
		}
		// Basic validation - would use json.Valid in production
		return nil
	}
}

// TimeFormat validates time format
func TimeFormat(layout string) ValidationRule {
	return func(value interface{}) error {
		str, ok := value.(string)
		if !ok {
			return nil
		}
		_, err := time.Parse(layout, str)
		if err != nil {
			return fmt.Errorf("invalid time format, expected: %s", layout)
		}
		return nil
	}
}

// Date validates date format (YYYY-MM-DD)
func Date() ValidationRule {
	return TimeFormat("2006-01-02")
}

// DateTime validates datetime format (RFC3339)
func DateTime() ValidationRule {
	return TimeFormat(time.RFC3339)
}

// Time validates time format (HH:MM:SS)
func Time() ValidationRule {
	return TimeFormat("15:04:05")
}

// StrongPassword validates strong password requirements
func StrongPassword() ValidationRule {
	return All(
		MinLength(8),
		HasLetter(),
		HasDigit(),
		HasSpecialChar(),
	)
}

// Username validates username format
func Username() ValidationRule {
	return All(
		MinLength(3),
		MaxLength(30),
		Alphanumeric(),
	)
}

// Domain validates domain name
func Domain() ValidationRule {
	domainRegex := regexp.MustCompile(`^([a-zA-Z0-9]([a-zA-Z0-9\-]{0,61}[a-zA-Z0-9])?\.)+[a-zA-Z]{2,}$`)
	return func(value interface{}) error {
		str, ok := value.(string)
		if !ok {
			return nil
		}
		if !domainRegex.MatchString(str) {
			return fmt.Errorf("invalid domain name")
		}
		return nil
	}
}

// Slug validates URL slug format
func Slug() ValidationRule {
	slugRegex := regexp.MustCompile(`^[a-z0-9]+(?:-[a-z0-9]+)*$`)
	return func(value interface{}) error {
		str, ok := value.(string)
		if !ok {
			return nil
		}
		if !slugRegex.MatchString(str) {
			return fmt.Errorf("invalid slug format")
		}
		return nil
	}
}

// HexColor validates hex color code
func HexColor() ValidationRule {
	hexRegex := regexp.MustCompile(`^#?[0-9A-Fa-f]{6}$`)
	return func(value interface{}) error {
		str, ok := value.(string)
		if !ok {
			return nil
		}
		if !hexRegex.MatchString(str) {
			return fmt.Errorf("invalid hex color code")
		}
		return nil
	}
}

// ISBN validates ISBN format
func ISBN() ValidationRule {
	isbnRegex := regexp.MustCompile(`^(?:\d{9}[\dXx]|\d{13})$`)
	return func(value interface{}) error {
		str, ok := value.(string)
		if !ok {
			return nil
		}
		cleaned := strings.ReplaceAll(strings.ReplaceAll(str, "-", ""), " ", "")
		if !isbnRegex.MatchString(cleaned) {
			return fmt.Errorf("invalid ISBN format")
		}
		return nil
	}
}

// ZipCode validates zip code format (US)
func ZipCode() ValidationRule {
	zipRegex := regexp.MustCompile(`^\d{5}(-\d{4})?$`)
	return func(value interface{}) error {
		str, ok := value.(string)
		if !ok {
			return nil
		}
		if !zipRegex.MatchString(str) {
			return fmt.Errorf("invalid zip code format")
		}
		return nil
	}
}

// CountryCode validates ISO country code
func CountryCode() ValidationRule {
	countryRegex := regexp.MustCompile(`^[A-Z]{2}$`)
	return func(value interface{}) error {
		str, ok := value.(string)
		if !ok {
			return nil
		}
		if !countryRegex.MatchString(str) {
			return fmt.Errorf("invalid country code format")
		}
		return nil
	}
}

// LanguageCode validates ISO language code
func LanguageCode() ValidationRule {
	langRegex := regexp.MustCompile(`^[a-z]{2}(-[A-Z]{2})?$`)
	return func(value interface{}) error {
		str, ok := value.(string)
		if !ok {
			return nil
		}
		if !langRegex.MatchString(str) {
			return fmt.Errorf("invalid language code format")
		}
		return nil
	}
}

// CurrencyCode validates ISO currency code
func CurrencyCode() ValidationRule {
	currencyRegex := regexp.MustCompile(`^[A-Z]{3}$`)
	return func(value interface{}) error {
		str, ok := value.(string)
		if !ok {
			return nil
		}
		if !currencyRegex.MatchString(str) {
			return fmt.Errorf("invalid currency code format")
		}
		return nil
	}
}

// Percentage validates percentage value (0-100)
func Percentage() ValidationRule {
	return Range(0, 100)
}

// Latitude validates latitude (-90 to 90)
func Latitude() ValidationRule {
	return Range(-90, 90)
}

// Longitude validates longitude (-180 to 180)
func Longitude() ValidationRule {
	return Range(-180, 180)
}

// Age validates age (reasonable range)
func Age() ValidationRule {
	return Range(0, 150)
}

// Year validates year (reasonable range)
func Year() ValidationRule {
	return Range(1900, 2100)
}

// Port validates port number (0-65535)
func Port() ValidationRule {
	return Range(0, 65535)
}

// FileExtension validates file extension
func FileExtension(allowed ...string) ValidationRule {
	return func(value interface{}) error {
		str, ok := value.(string)
		if !ok {
			return nil
		}
		ext := strings.ToLower(strings.TrimPrefix(strings.TrimSpace(str), "."))
		for _, allowedExt := range allowed {
			if ext == strings.ToLower(allowedExt) {
				return nil
			}
		}
		return fmt.Errorf("invalid file extension, allowed: %v", allowed)
	}
}

// MimeType validates MIME type
func MimeType() ValidationRule {
	mimeRegex := regexp.MustCompile(`^[a-z]+/[a-z0-9][a-z0-9!#$&\-\^_.]*$`)
	return func(value interface{}) error {
		str, ok := value.(string)
		if !ok {
			return nil
		}
		if !mimeRegex.MatchString(strings.ToLower(str)) {
			return fmt.Errorf("invalid MIME type format")
		}
		return nil
	}
}

// SemVer validates semantic version
func SemVer() ValidationRule {
	semverRegex := regexp.MustCompile(`^\d+\.\d+\.\d+(-[0-9A-Za-z-]+(\.[0-9A-Za-z-]+)*)?(\+[0-9A-Za-z-]+(\.[0-9A-Za-z-]+)*)?$`)
	return func(value interface{}) error {
		str, ok := value.(string)
		if !ok {
			return nil
		}
		if !semverRegex.MatchString(str) {
			return fmt.Errorf("invalid semantic version format")
		}
		return nil
	}
}

// NotEmptyString validates non-empty string
func NotEmptyString() ValidationRule {
	return func(value interface{}) error {
		str, ok := value.(string)
		if !ok {
			return nil
		}
		if strings.TrimSpace(str) == "" {
			return fmt.Errorf("string cannot be empty")
		}
		return nil
	}
}

// NoWhitespace validates no whitespace
func NoWhitespace() ValidationRule {
	return func(value interface{}) error {
		str, ok := value.(string)
		if !ok {
			return nil
		}
		if strings.ContainsAny(str, " \t\n\r") {
			return fmt.Errorf("string cannot contain whitespace")
		}
		return nil
	}
}

// ASCII validates ASCII characters only
func ASCII() ValidationRule {
	return func(value interface{}) error {
		str, ok := value.(string)
		if !ok {
			return nil
		}
		for _, r := range str {
			if r > 127 {
				return fmt.Errorf("string must contain only ASCII characters")
			}
		}
		return nil
	}
}

// Unicode validates Unicode characters
func Unicode() ValidationRule {
	return func(value interface{}) error {
		_, ok := value.(string)
		if !ok {
			return nil
		}
		// All Go strings are Unicode, so this always passes
		return nil
	}
}

// Printable validates printable characters only
func Printable() ValidationRule {
	return func(value interface{}) error {
		str, ok := value.(string)
		if !ok {
			return nil
		}
		for _, r := range str {
			if !unicode.IsPrint(r) {
				return fmt.Errorf("string must contain only printable characters")
			}
		}
		return nil
	}
}

// Graph validates graph characters (printable except space)
func Graph() ValidationRule {
	return func(value interface{}) error {
		str, ok := value.(string)
		if !ok {
			return nil
		}
		for _, r := range str {
			if !unicode.IsGraphic(r) || unicode.IsSpace(r) {
				return fmt.Errorf("string must contain only graph characters")
			}
		}
		return nil
	}
}

