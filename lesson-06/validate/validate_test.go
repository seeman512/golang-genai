// Homework — Task 2: extend emailCases and/or phoneCases below to at
// least 8 cases each (your mentor may ask for just one of the two
// functions), then implement validate.go until every subtest passes.
//
// Like todo_test.go, this file uses t.Run per case and t.Errorf (not
// t.Fatalf) for the actual assertions, so one wrong case never hides
// the others. Run `go test -v ./validate/...` and read every FAIL line.
package validate

import (
	"strings"
	"testing"
)

const minCases = 8

// emailCases is the table of test cases for ValidateEmail.
//
// The table includes malformed addresses and boundary cases in addition
// to the basic valid examples.
var emailCases = []struct {
	name  string
	input string
	want  bool
}{
	{"valid simple", "student@softserve.academy", true},
	{"missing at sign", "student-softserve.academy", false},
	{"empty string", "", false},
	{"whitespace inside", "student name@example.com", false},
	{"missing domain", "student@", false},
	{"trailing dot", "student@example.com.", false},
	{"consecutive dots", "student..name@example.com", false},
	{"unicode character", "élève@example.com", false},
	{"valid plus tag", "student+tag@example.com", true},
	{"valid single-label domain", "student@localhost", true},
	{"local part too long", strings.Repeat("a", 65) + "@example.com", false},
}

func TestValidateEmail(t *testing.T) {
	if len(emailCases) < minCases {
		t.Fatalf(
			"emailCases has %d case(s), need at least %d — add more edge cases to validate_test.go before this test can pass",
			len(emailCases), minCases,
		)
	}

	for _, tc := range emailCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			got := ValidateEmail(tc.input)
			if got != tc.want {
				t.Errorf("ValidateEmail(%q) = %v, want %v", tc.input, got, tc.want)
			}
		})
	}
}

// phoneCases is the table of test cases for ValidatePhone.
// This is only required if your mentor asked you to validate phone
// numbers instead of (or in addition to) email addresses.
//
// The table includes malformed numbers and examples of each accepted
// format.
var phoneCases = []struct {
	name  string
	input string
	want  bool
}{
	{"valid with plus", "+380501234567", true},
	{"contains letters", "050-abc-4567", false},
	{"empty string", "", false},
	{"valid national format", "050-123-4567", true},
	{"valid national digits", "0501234567", true},
	{"international too short", "+380501234", false},
	{"international too long", "+3805012345678901", false},
	{"spaces are not accepted", "+380 50 123 4567", false},
	{"parentheses are not accepted", "+380(50)1234567", false},
	{"dots are not accepted", "050.123.4567", false},
	{"wrong national digit count", "050-123-456", false},
}

func TestValidatePhone(t *testing.T) {
	if len(phoneCases) < minCases {
		t.Fatalf(
			"phoneCases has %d case(s), need at least %d — add more edge cases to validate_test.go before this test can pass",
			len(phoneCases), minCases,
		)
	}

	for _, tc := range phoneCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			got := ValidatePhone(tc.input)
			if got != tc.want {
				t.Errorf("ValidatePhone(%q) = %v, want %v", tc.input, got, tc.want)
			}
		})
	}
}
