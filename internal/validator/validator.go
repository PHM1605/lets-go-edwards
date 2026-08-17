package validator

import (
	"slices"
	"strings"
	"unicode/utf8"
)

type Validator struct {
	FieldErrors map[string]string
}

// True if FieldErrors doesn't have any entry
func (v *Validator) Valid() bool {
	return len(v.FieldErrors) == 0
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
