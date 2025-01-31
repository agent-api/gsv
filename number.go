package gsv

import (
	"cmp"
	"encoding/json"
	"fmt"

	"github.com/agent-api/gsv/pkg/jsonschema"
)

const (
	MinNumberError         ValidationErrorType = "min_number"
	MaxNumberError         ValidationErrorType = "max_number"
	RequiredNumberError    ValidationErrorType = "required_number"
	InvalidNumberTypeError ValidationErrorType = "invalid_number_type"
)

// NumberSchema implements the Schema interface. It represents a generic "number"
// with types implemented in int.go, uint.go, float.go, and rune.go.
type NumberSchema[T cmp.Ordered] struct {
	min   *T
	max   *T
	value *T

	description *string

	validators []func(T)
	isOptional bool
	result     *ValidationResult
}

// Number creates a new number validator for a specific type
func Number[T cmp.Ordered]() *NumberSchema[T] {
	return &NumberSchema[T]{
		validators: make([]func(T), 0),
		isOptional: false,
	}
}

func (n *NumberSchema[T]) Min(min T, opts ...ValidationOptions) *NumberSchema[T] {
	n.min = &min

	validationMessage := fmt.Sprintf("must be at least %v", min)
	if len(opts) > 0 && opts[0].Message != "" {
		validationMessage = opts[0].Message
	}

	n.validators = append(n.validators, func(v T) {
		if v < min {
			n.result.AddError(&ValidationError{
				Type:     MinNumberError,
				Message:  validationMessage,
				Expected: min,
				Actual:   v,
			})
		}
	})

	return n
}

func (n *NumberSchema[T]) Max(max T, opts ...ValidationOptions) *NumberSchema[T] {
	n.max = &max

	validationMessage := fmt.Sprintf("must not exceed: %v", max)
	if len(opts) > 0 && opts[0].Message != "" {
		validationMessage = opts[0].Message
	}

	n.validators = append(n.validators, func(v T) {
		if v > max {
			n.result.AddError(&ValidationError{
				Type:     MaxNumberError,
				Message:  validationMessage,
				Expected: max,
				Actual:   v,
			})
		}
	})

	return n
}

// Optional marks the int field as optional
func (n *NumberSchema[T]) Optional() *NumberSchema[T] {
	n.isOptional = true
	return n
}

func (n *NumberSchema[T]) IsOptional() bool {
	return n.isOptional
}

func (n *NumberSchema[T]) Description(val string) *NumberSchema[T] {
	n.description = &val
	return n
}

func (n *NumberSchema[T]) Set(v T) *NumberSchema[T] {
	n.value = &v
	return n
}

func (n *NumberSchema[T]) Value() (T, bool) {
	if n.value == nil {
		var zero T
		return zero, false
	}

	return *n.value, true
}

// Validate performs the validation
func (n *NumberSchema[T]) Validate() *ValidationResult {
	n.result = &ValidationResult{}

	val, ok := n.Value()
	if !ok && n.isOptional {
		n.result.AddError(&ValidationError{
			Type:    RequiredNumberError,
			Message: "value has not been set",
		})
	}

	for _, validator := range n.validators {
		validator(val)
	}

	return n.result
}

func (n *NumberSchema[T]) MarshalJSON() ([]byte, error) {
	if n.value == nil {
		if n.isOptional {
			return json.Marshal(nil)
		}
		return nil, fmt.Errorf("required field has no value")
	}
	return json.Marshal(*n.value)
}

func (n *NumberSchema[T]) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		if !n.isOptional {
			return fmt.Errorf("validation failed: field is required")
		}
		n.value = nil
		return nil
	}

	var v T
	if err := json.Unmarshal(data, &v); err != nil {
		return fmt.Errorf("invalid numeric value: %w", err)
	}

	n.value = &v

	if result := n.Validate(); result.HasErrors() {
		return result.Error()
	}

	return nil
}

func (n *NumberSchema[T]) CompileJSONSchema(schema *jsonschema.JSONSchema, jsonTag string) error {
	if n == nil {
		return fmt.Errorf("found nil schema interface with JSON tag: %s", jsonTag)
	}

	propertySchema := &jsonschema.JSONSchema{
		Type: "number",
	}

	// Add description if present
	if n.description != nil {
		propertySchema.Description = *n.description
	}

	// Add to required fields if not optional
	if !n.IsOptional() {
		schema.Required = append(schema.Required, jsonTag)
	}

	schema.Properties[jsonTag] = propertySchema
	return nil
}
