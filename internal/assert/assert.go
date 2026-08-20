package assert

import (
	"strings"
	"testing"
)

func Equal[T comparable](t *testing.T, actual, expected T) {
	// Indicate that this function is a Test Helper
	// i.e. t.Errorf triggered here will be reported in the Function that calls Equal()
	t.Helper()

	if actual != expected {
		t.Errorf("got: %v; want: %v", actual, expected)
	}
}

// Check if a substring (2nd) is inside a string (1st)
func StringContains(t *testing.T, actual, expectedSubstring string) {
	// Indicate that this function is a Test Helper
	// i.e. t.Errorf triggered here will be reported in the Function that calls Equal()
	t.Helper()

	if !strings.Contains(actual, expectedSubstring) {
		t.Errorf("got: %q; expected to contain: %q", actual, expectedSubstring)
	}
}

// Check if an Error is "nil" or not
func NilError(t *testing.T, actual error) {
	t.Helper()

	if actual != nil {
		t.Errorf("got: %v; expected: nil", actual)
	}
}
