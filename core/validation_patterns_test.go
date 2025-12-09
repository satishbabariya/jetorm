package core

import (
	"testing"
)

func TestPhoneNumber(t *testing.T) {
	rule := PhoneNumber()

	testCases := []struct {
		value string
		valid bool
	}{
		{"+1234567890", true},
		{"1234567890", true},
		{"invalid", false},
	}

	for _, tc := range testCases {
		err := rule(tc.value)
		if tc.valid && err != nil {
			t.Errorf("Phone %s should be valid, got error: %v", tc.value, err)
		}
		if !tc.valid && err == nil {
			t.Errorf("Phone %s should be invalid", tc.value)
		}
	}
}

func TestUUID(t *testing.T) {
	rule := UUID()

	testCases := []struct {
		value string
		valid bool
	}{
		{"550e8400-e29b-41d4-a716-446655440000", true},
		{"550E8400-E29B-41D4-A716-446655440000", true},
		{"invalid-uuid", false},
	}

	for _, tc := range testCases {
		err := rule(tc.value)
		if tc.valid && err != nil {
			t.Errorf("UUID %s should be valid, got error: %v", tc.value, err)
		}
		if !tc.valid && err == nil {
			t.Errorf("UUID %s should be invalid", tc.value)
		}
	}
}

func TestIPv4(t *testing.T) {
	rule := IPv4()

	testCases := []struct {
		value string
		valid bool
	}{
		{"192.168.1.1", true},
		{"255.255.255.255", true},
		{"256.1.1.1", false},
		{"invalid", false},
	}

	for _, tc := range testCases {
		err := rule(tc.value)
		if tc.valid && err != nil {
			t.Errorf("IPv4 %s should be valid, got error: %v", tc.value, err)
		}
		if !tc.valid && err == nil {
			t.Errorf("IPv4 %s should be invalid", tc.value)
		}
	}
}

func TestDate(t *testing.T) {
	rule := Date()

	testCases := []struct {
		value string
		valid bool
	}{
		{"2023-12-25", true},
		{"2023-01-01", true},
		{"invalid-date", false},
		{"2023/12/25", false},
	}

	for _, tc := range testCases {
		err := rule(tc.value)
		if tc.valid && err != nil {
			t.Errorf("Date %s should be valid, got error: %v", tc.value, err)
		}
		if !tc.valid && err == nil {
			t.Errorf("Date %s should be invalid", tc.value)
		}
	}
}

func TestStrongPassword(t *testing.T) {
	rule := StrongPassword()

	testCases := []struct {
		value string
		valid bool
	}{
		{"SecurePass123!", true},
		{"weak", false},
		{"12345678", false},
		{"password", false},
	}

	for _, tc := range testCases {
		err := rule(tc.value)
		if tc.valid && err != nil {
			t.Errorf("Password %s should be valid, got error: %v", tc.value, err)
		}
		if !tc.valid && err == nil {
			t.Errorf("Password %s should be invalid", tc.value)
		}
	}
}

func TestSlug(t *testing.T) {
	rule := Slug()

	testCases := []struct {
		value string
		valid bool
	}{
		{"hello-world", true},
		{"hello123", true},
		{"Hello-World", false},
		{"hello world", false},
	}

	for _, tc := range testCases {
		err := rule(tc.value)
		if tc.valid && err != nil {
			t.Errorf("Slug %s should be valid, got error: %v", tc.value, err)
		}
		if !tc.valid && err == nil {
			t.Errorf("Slug %s should be invalid", tc.value)
		}
	}
}

func TestLatitude(t *testing.T) {
	rule := Latitude()

	testCases := []struct {
		value float64
		valid bool
	}{
		{45.0, true},
		{90.0, true},
		{-90.0, true},
		{91.0, false},
		{-91.0, false},
	}

	for _, tc := range testCases {
		err := rule(tc.value)
		if tc.valid && err != nil {
			t.Errorf("Latitude %f should be valid, got error: %v", tc.value, err)
		}
		if !tc.valid && err == nil {
			t.Errorf("Latitude %f should be invalid", tc.value)
		}
	}
}

func TestLongitude(t *testing.T) {
	rule := Longitude()

	testCases := []struct {
		value float64
		valid bool
	}{
		{120.0, true},
		{180.0, true},
		{-180.0, true},
		{181.0, false},
		{-181.0, false},
	}

	for _, tc := range testCases {
		err := rule(tc.value)
		if tc.valid && err != nil {
			t.Errorf("Longitude %f should be valid, got error: %v", tc.value, err)
		}
		if !tc.valid && err == nil {
			t.Errorf("Longitude %f should be invalid", tc.value)
		}
	}
}

func TestSemVer(t *testing.T) {
	rule := SemVer()

	testCases := []struct {
		value string
		valid bool
	}{
		{"1.0.0", true},
		{"1.2.3", true},
		{"1.0.0-alpha", true},
		{"1.0.0-alpha.1", true},
		{"1.0.0+20130313144700", true},
		{"invalid", false},
		{"1.0", false},
	}

	for _, tc := range testCases {
		err := rule(tc.value)
		if tc.valid && err != nil {
			t.Errorf("SemVer %s should be valid, got error: %v", tc.value, err)
		}
		if !tc.valid && err == nil {
			t.Errorf("SemVer %s should be invalid", tc.value)
		}
	}
}

