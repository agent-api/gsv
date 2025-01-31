package gsv

import (
	"encoding/json"
	"fmt"

	"github.com/agent-api/gsv/pkg/jsonschema"
)

// stringValidatorFunc is a validation function that expects a string:
type stringValidatorFunc func(string)

const (
	MinStringLengthError   ValidationErrorType = "min_string_length"
	MaxStringLengthError   ValidationErrorType = "max_string_length"
	RequiredStringError    ValidationErrorType = "required_string"
	InvalidStringTypeError ValidationErrorType = "invalid_string_type"
)

const (
	StringSchemaType string = "string"
)

// StringSchema implements the Schema interface for strings.
type StringSchema struct {
	schemaType string

	// minLength is the optional field to denote the minimum length of the string
	minLength *int

	// maxLength is the optional field to denote the maximum length of the string
	maxLength *int

	// validators are the registered functions to validate the string against
	validators []stringValidatorFunc

	// value is the actual string value
	value *string

	description *string

	// isOptional denotes if the string value in the schema is optional
	isOptional bool

	// the result of the last validation
	result *ValidationResult
}

// String creates a new string validator
func String() *StringSchema {
	return &StringSchema{
		schemaType: StringSchemaType,
		validators: make([]stringValidatorFunc, 0),
		isOptional: false,
	}
}

// Optional marks the string field as optional
func (s *StringSchema) Description(val string) *StringSchema {
	s.description = &val
	return s
}

// Optional marks the string field as optional
func (s *StringSchema) Optional() *StringSchema {
	s.isOptional = true
	return s
}

// Optional marks the string field as optional
func (s *StringSchema) IsOptional() bool {
	return s.isOptional
}

// Min adds minimum length validation
func (v *StringSchema) Min(length int, opts ...ValidationOptions) *StringSchema {
	v.minLength = &length

	validationMessage := fmt.Sprintf("must be at least %d characters long", *v.minLength)
	if len(opts) > 0 && opts[0].Message != "" {
		validationMessage = opts[0].Message
	}

	v.validators = append(v.validators, func(s string) {
		if *v.minLength > len(s) {
			v.result.AddError(&ValidationError{
				Type:     MinStringLengthError,
				Message:  validationMessage,
				Expected: *v.minLength,
				Actual:   len(s),
			})
		}
	})

	return v
}

// Max adds maximum length validation
func (v *StringSchema) Max(length int, opts ...ValidationOptions) *StringSchema {
	v.maxLength = &length

	validationMessage := fmt.Sprintf("must be at most %d characters long", *v.maxLength)
	if len(opts) > 0 && opts[0].Message != "" {
		validationMessage = opts[0].Message
	}

	v.validators = append(v.validators, func(s string) {
		if *v.maxLength < len(s) {
			v.result.AddError(&ValidationError{
				Type:     MaxStringLengthError,
				Message:  validationMessage,
				Expected: *v.maxLength,
				Actual:   len(s),
			})
		}
	})

	return v
}

// Validate performs the validation
func (v *StringSchema) Validate() *ValidationResult {
	v.result = &ValidationResult{}

	val, ok := v.Value()
	if !ok && !v.isOptional {
		v.result.AddError(&ValidationError{
			Type:    RequiredStringError,
			Message: "value has not been set",
		})
		return v.result
	}

	for _, validator := range v.validators {
		validator(val)
	}

	return v.result
}

func (v *StringSchema) Set(s string) *StringSchema {
	v.value = &s
	return v
}

// UnmarshalJSON implements json.Unmarshaler
func (s *StringSchema) MarshalJSON() ([]byte, error) {
	if s.value == nil {
		if s.isOptional {
			return json.Marshal(nil)
		}
		return nil, fmt.Errorf("required field has no value")
	}

	return json.Marshal(*s.value)
}

// UnmarshalJSON implements json.Unmarshaler
func (s *StringSchema) UnmarshalJSON(data []byte) error {

	// Handle null values
	if string(data) == "null" {
		if !s.isOptional {
			return fmt.Errorf("validation failed: field is required")
		}
		s.value = nil
		return nil
	}

	// Handle missing fields (empty string in JSON)
	if len(data) == 0 {
		if !s.isOptional {
			return fmt.Errorf("validation failed: field is required")
		}
		s.value = nil
		return nil
	}

	// Unmarshal the string value
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return fmt.Errorf("invalid string value: %w", err)
	}

	// Store the value
	s.value = &str

	// Run validation and return errors if there are any
	if result := s.Validate(); result.HasErrors() {
		return result.Error()
	}

	return nil
}

func (s *StringSchema) CompileJSONSchema(schema *jsonschema.JSONSchema, jsonTag string) error {
	if s == nil {
		return fmt.Errorf("found nil schema interface with JSON tag: %s", jsonTag)
	}

	propertySchema := &jsonschema.JSONSchema{
		Type: "string",
	}

	// Add description if present
	if s.description != nil {
		propertySchema.Description = *s.description
	}

	// Add min length if present
	if s.minLength != nil {
		propertySchema.MinLength = s.minLength
	}

	// Add max length if present
	if s.maxLength != nil {
		propertySchema.MaxLength = s.maxLength
	}

	// Add to required fields if not optional
	if !s.IsOptional() {
		schema.Required = append(schema.Required, jsonTag)
	}

	schema.Properties[jsonTag] = propertySchema
	return nil
}

// Value returns the validated string value. This method returns ("", false) if
// the string value is a null pointer due to it not being set
func (sv *StringSchema) Value() (string, bool) {
	if sv.value == nil {
		return "", false
	}

	return *sv.value, true
}
