package validator

import (
	"net/mail"
	"strings"
	"unicode/utf8"
)

type Validator struct {
	FieldErrors map[string]string
	// a generic map for errors that don't relate to a specific form field.
	NonFieldErrors []string
}

// Valid() returns true if the FieldErrors and NonFieldErrors map doesn't contain any entries.
func (v *Validator) Valid() bool {
	return len(v.FieldErrors) == 0 && len(v.NonFieldErrors) == 0
}

// AddFieldError() adds an error message to the FieldErrors map
// as long as no  entry already exists for the given key.
func (v *Validator) AddFieldError(key, message string) {
	//  initialize the map, if it isn't already initialized.
	if v.FieldErrors == nil {
		v.FieldErrors = make(map[string]string)
	}
	if _, exists := v.FieldErrors[key]; !exists {
		v.FieldErrors[key] = message
	}
}

func (v *Validator) AddNonFieldError(message string) {
	v.NonFieldErrors = append(v.NonFieldErrors, message)
}

// add an error message to the FieldErrors map if a
// validation check is not 'ok'.
func (v *Validator) CheckField(ok bool, key, message string) {
	if !ok {
		v.AddFieldError(key, message)
	}
}

// returns true if a value is not an empty string.
func NotBlank(value string) bool {
	return strings.TrimSpace(value) != ""
}

// return true if a string characters are less than or equal to n.
func MaxChars(value string, n int) bool {
	return utf8.RuneCountInString(value) <= n
}

// return true if a string characters are more than or equal to n.
func MinChars(value string, n int) bool {
	return utf8.RuneCountInString(value) >= n
}

func ValidateEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

func IsBlank(s string) bool {
	return strings.TrimSpace(s) == ""
}
