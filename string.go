package gsv

import (
	"encoding/json"
	"fmt"
	"strings"
)

type StringValidatorFunc func(string)

// StringSchema represents a string validation chain
type StringSchema struct {
	minLength *int
	maxLength *int
	//pattern       *string
	//email         bool
	//customChecks  []func(string) error
	errors ValidationErrorMap

	// registered functions to validate against
	validators []StringValidatorFunc

	value *string // Using pointer to handle null values
}

// String creates a new string validator
func String() *StringSchema {
	return &StringSchema{
		validators: make([]StringValidatorFunc, 0),
		errors:     make(ValidationErrorMap),
	}
}

// Min adds minimum length validation
func (v *StringSchema) Min(length int, opts ...ValidationOptions) *StringSchema {
	v.minLength = &length

	v.validators = append(v.validators, func(s string) {
		if *v.minLength > len(s) {
			if len(opts) > 0 && opts[0].Message != "" {
				v.errors["min"] = &ValidationError{
					Message: opts[0].Message,
				}
			} else {
				v.errors["min"] = &ValidationError{
					Message: fmt.Sprintf("must be at least %d characters long", *v.minLength),
				}
			}
		}
	})

	return v
}

// Max adds maximum length validation
func (v *StringSchema) Max(length int, opts ...ValidationOptions) *StringSchema {
	v.maxLength = &length

	v.validators = append(v.validators, func(s string) {
		if *v.maxLength < len(s) {
			if len(opts) > 0 && opts[0].Message != "" {
				v.errors["max"] = &ValidationError{
					Message: opts[0].Message,
				}
			} else {
				v.errors["max"] = &ValidationError{
					Message: fmt.Sprintf("must be at most %d characters long", *v.maxLength),
				}
			}
		}
	})

	return v
}

// Validate performs the validation
func (v *StringSchema) Validate(value string) ValidationErrorMap {
	for _, validator := range v.validators {
		validator(value)
	}

	return v.errors
}

// UnmarshalJSON implements json.Unmarshaler
func (sv *StringSchema) UnmarshalJSON(data []byte) error {
	// Handle null values as per convention
	if string(data) == "null" {
		sv.value = nil
		return nil
	}

	// Unmarshal the string value
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return fmt.Errorf("invalid string value: %w", err)
	}

	// Store the value
	sv.value = &str

	// Run validation
	if errMap := sv.Validate(str); len(errMap) > 0 {
		// Convert validation errors to a single error
		var errMsgs []string
		for key, err := range errMap {
			errMsgs = append(errMsgs, fmt.Sprintf("%s: %s", key, err.Message))
		}
		return fmt.Errorf("validation failed: %s", strings.Join(errMsgs, "; "))
	}

	return nil
}

// Value returns the validated string value
func (sv *StringSchema) Value() (string, bool) {
	if sv.value == nil {
		return "", false
	}
	return *sv.value, true
}
