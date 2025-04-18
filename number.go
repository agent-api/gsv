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

type NumberValidatorFunc[T cmp.Ordered] func(T)

// NumberSchema implements the Schema interface. It represents a generic "number"
// with types implemented in int.go, uint.go, float.go, and rune.go.
type NumberSchema[T cmp.Ordered] struct {
	min   *T
	max   *T
	value *T

	description *string

	validators []NumberValidatorFunc[T]
	isOptional bool
	result     *ValidationResult
}

// Number creates a new number validator for a specific type
func Number[T cmp.Ordered]() *NumberSchema[T] {
	return &NumberSchema[T]{
		validators: make([]NumberValidatorFunc[T], 0),
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

func (n *NumberSchema[T]) setValue(val interface{}) error {
	num, ok := val.(T)
	if !ok {
		return fmt.Errorf("expected %T value, got %T", *new(T), val)
	}
	n.value = &num
	return nil
}

// NumberSchema
func (n *NumberSchema[T]) Value() (T, bool) {
	val, ok := n.getValue()
	if !ok {
		var zero T
		return zero, false
	}
	numVal, ok := val.(T)
	if !ok {
		panic(fmt.Sprintf("NumberSchema: invalid internal value type %T, expected %T", val, *new(T)))
	}
	return numVal, true
}

func (n *NumberSchema[T]) getValue() (interface{}, bool) {
	if n.value == nil {
		return nil, false
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

func (n *NumberSchema[T]) Clone() Schema {
	// Create new instance
	clone := &NumberSchema[T]{
		isOptional: n.isOptional,
		validators: make([]NumberValidatorFunc[T], len(n.validators)),
	}
	// Deep copy validators slice
	copy(clone.validators, n.validators)

	// Deep copy pointer fields
	if n.min != nil {
		min := *n.min
		clone.min = &min
	}
	if n.max != nil {
		max := *n.max
		clone.max = &max
	}
	if n.value != nil {
		val := *n.value
		clone.value = &val
	}
	if n.description != nil {
		desc := *n.description
		clone.description = &desc
	}

	// Initialize new validation result
	clone.result = &ValidationResult{}

	return clone
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
