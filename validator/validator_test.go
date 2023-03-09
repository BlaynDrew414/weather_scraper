package validator

import (
	"testing"

)

func TestIsValidInput(t *testing.T) {
	cases := []struct {
		name       string
		city       string
		state      string
		wantResult bool
	}{
		{"empty city and state", "", "", false},
		{"empty city", "", "New York", false},
		{"empty state", "New York", "", false},
		{"valid input", "New York", "NY", true},
		{"input with numbers", "San Francisco1", "CA", false},
		{"input with special characters", "Washington*", "DC", false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := IsValidInput(tc.city, tc.state)
			if got != tc.wantResult {
				t.Errorf("IsValidInput(%q, %q) = %v, want %v", tc.city, tc.state, got, tc.wantResult)
			}
		})
	}
}
