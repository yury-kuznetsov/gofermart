package validation

import (
	"testing"
)

func TestCheckLuhn(t *testing.T) {
	testCases := []struct {
		description string
		input       string
		expected    bool
	}{
		{"Case 1", "79927398713", true},
		{"Case 2", "4532015112830366", true},
		{"Case 3", "", false},
		{"Case 4", "79927398710", false},
		{"Case 5", "1234567812345678", false},
	}

	for _, test := range testCases {
		if result := IsValidLuhn(test.input); result != test.expected {
			t.Errorf("Description: %s. Expected %v, got %v", test.description, test.expected, result)
		}
	}
}
