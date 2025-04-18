package gsv

import (
	"encoding/json"
	"fmt"

	"github.com/agent-api/gsv/pkg/jsonschema"
)

// boolValidatorFunc is a validation function that expects a boolean for validation
type boolValidatorFunc func(bool)

const (
	BoolRequiredError    ValidationErrorType = "required"
	BoolInvalidTypeError ValidationErrorType = "invalid_type"
)

// BoolSchema implements the Schema interface for booleans.
type BoolSchema struct {
	// TODO need schema type

	// validators are the registered functions to validate the string against
	validators []boolValidatorFunc

	value *bool // Using pointer to handle null values

	description *string

	// isOptional denotes if the bool value in the schema is optional
	isOptional bool

	// the result of the last validation
	result *ValidationResult
}

// Optional marks the bool field as optional
func (b *BoolSchema) Optional() *BoolSchema {
	b.isOptional = true
	return b
}

func (b *BoolSchema) IsOptional() bool {
	return b.isOptional
}

// NewBool creates a new string validator
func Bool() *BoolSchema {
	return &BoolSchema{
		validators: make([]boolValidatorFunc, 0),
		isOptional: false,
	}
}

func (b *BoolSchema) Set(v bool) *BoolSchema {
	b.value = &v
	return b
}

func (b *BoolSchema) setValue(val interface{}) error {
	s, ok := val.(bool)
	if !ok {
		return fmt.Errorf("expected bool value, got %T", val)
	}
	b.value = &s
	return nil
}

func (b *BoolSchema) Description(val string) *BoolSchema {
	b.description = &val
	return b
}

// Validate performs the validation
func (b *BoolSchema) Validate() *ValidationResult {
	b.result = &ValidationResult{}

	val, ok := b.Value()
	if !ok && b.isOptional {
		b.result.AddError(&ValidationError{
			Type:    BoolRequiredError,
			Message: "bool has not been set",
		})
	}

	for _, validator := range b.validators {
		validator(val)
	}

	return b.result
}

// UnmarshalJSON implements json.Unmarshaler
func (b *BoolSchema) MarshalJSON() ([]byte, error) {
	if b.value == nil {
		if b.isOptional {
			return json.Marshal(nil)
		}
		return nil, fmt.Errorf("required field has no value")
	}

	return json.Marshal(*b.value)
}

// UnmarshalJSON implements json.Unmarshaler
func (b *BoolSchema) UnmarshalJSON(data []byte) error {
	// Handle null values
	if string(data) == "null" {
		if !b.isOptional {
			return fmt.Errorf("validation failed: field is required")
		}
		b.value = nil
		return nil
	}

	// Handle missing fields (empty string in JSON)
	if len(data) == 0 {
		if !b.isOptional {
			return fmt.Errorf("validation failed: field is required")
		}
		b.value = nil
		return nil
	}

	// Unmarshal the string value
	var v bool
	if err := json.Unmarshal(data, &v); err != nil {
		return fmt.Errorf("invalid int value: %w", err)
	}

	// Store the value
	b.value = &v

	// Run validation and return errors if there are any
	if result := b.Validate(); result.HasErrors() {
		return result.Error()
	}

	return nil
}

// Clone implements Schema.Clone by creating a deep copy of the StringSchema
func (b *BoolSchema) Clone() Schema {
	// Create new instance
	clone := &BoolSchema{
		isOptional: b.isOptional,
		validators: make([]boolValidatorFunc, len(b.validators)),
	}

	// Deep copy the validators slice
	copy(clone.validators, b.validators)

	if b.description != nil {
		desc := *b.description
		clone.description = &desc
	}

	if b.value != nil {
		val := *b.value
		clone.value = &val
	}

	// Initialize a new validation result
	clone.result = &ValidationResult{}

	return clone
}

// Value returns the validated string value
func (b *BoolSchema) Value() (bool, bool) {
	val, ok := b.getValue()
	if !ok {
		return false, false
	}
	boolVal, ok := val.(bool)
	if !ok {
		// who watches the watchers? - this should never happen
		panic(fmt.Sprintf("BoolSchema: invalid internal value type %T, expected bool", val))
	}
	return boolVal, true
}

func (b *BoolSchema) getValue() (interface{}, bool) {
	if b.value == nil {
		return nil, false
	}
	return *b.value, true
}

func (b *BoolSchema) CompileJSONSchema(schema *jsonschema.JSONSchema, jsonTag string) error {
	if b == nil {
		return fmt.Errorf("found nil schema interface with JSON tag: %s", jsonTag)
	}

	propertySchema := &jsonschema.JSONSchema{
		Type: "boolean",
	}

	// Add description if present
	if b.description != nil {
		propertySchema.Description = *b.description
	}

	// Add to required fields if not optional
	if !b.IsOptional() {
		schema.Required = append(schema.Required, jsonTag)
	}

	schema.Properties[jsonTag] = propertySchema
	return nil
}
