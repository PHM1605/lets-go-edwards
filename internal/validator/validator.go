package validator

import (
	"regexp"
	"slices"
	"strings"
	"unicode/utf8"
)

type Validator struct {
	NonFieldErrors []string          // General errors for the Login page
	FieldErrors    map[string]string // Field-specific errors for the Signup page
}

// email matching standard; MustCompile() panics if there is a syntax error in string argument
// convert string to "*regexp.Regexp type"
var EmailRX = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

// True if FieldErrors doesn't have any entry
func (v *Validator) Valid() bool {
	return len(v.FieldErrors) == 0 && len(v.NonFieldErrors) == 0
}

// Add error message to a Field in FieldErrors
func (v *Validator) AddFieldError(key, message string) {
	// make map first if not exists
	if v.FieldErrors == nil {
		v.FieldErrors = make(map[string]string)
	}
	// add entry if not yet exists
	if _, exists := v.FieldErrors[key]; !exists {
		v.FieldErrors[key] = message
	}
}

// Add error message to NonFieldErrors
func (v *Validator) AddNonFieldError(message string) {
	v.NonFieldErrors = append(v.NonFieldErrors, message)
}

func (v *Validator) CheckField(ok bool, key, message string) {
	if !ok {
		v.AddFieldError(key, message)
	}
}

// --------- Auxiliary functions -------------
// Check if empty string
func NotBlank(value string) bool {
	return strings.TrimSpace(value) != ""
}

// Check if string has no more than n characters
func MaxChars(value string, n int) bool {
	return utf8.RuneCountInString(value) <= n
}

// Check if "value" is in a list of permitted values "permittedValues"
func PermittedValues[T comparable](value T, permittedValues ...T) bool {
	return slices.Contains(permittedValues, value)
}

// Check string has at least "n" characters
func MinChars(value string, n int) bool {
	return utf8.RuneCountInString(value) >= n
}

// Check if a string matches a regex (for Email)
func Matches(value string, rx *regexp.Regexp) bool {
	return rx.MatchString(value)
}
