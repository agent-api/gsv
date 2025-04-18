package gsv

import (
	"encoding/json"
	"fmt"

	"github.com/agent-api/gsv/pkg/jsonschema"
)

const (
	ArraySchemaType                              = "array"
	MinItemsError            ValidationErrorType = "min_items"
	MaxItemsError                                = "max_items"
	RequiredArrayError                           = "required_array"
	InvalidElementTypeError                      = "invalid_element_type"
	MissingElementValueError                     = "missing_element_value"
)

type ArraySchema struct {
	schemaType    string
	elementSchema Schema
	minItems      *int
	maxItems      *int
	value         []interface{}
	isOptional    bool
	description   *string
	result        *ValidationResult
}

func Array(elementSchema Schema) *ArraySchema {
	if elementSchema == nil {
		panic("elementSchema cannot be nil")
	}

	return &ArraySchema{
		schemaType:    ArraySchemaType,
		elementSchema: elementSchema,
		result:        &ValidationResult{},
	}
}

// IsOptional implements Schema.IsOptional
func (a *ArraySchema) IsOptional() bool {
	return a.isOptional
}

func (a *ArraySchema) MinItems(min int, opts ...ValidationOptions) *ArraySchema {
	if min < 0 {
		panic("minItems cannot be negative")
	}
	a.minItems = &min
	return a
}

func (a *ArraySchema) MaxItems(max int, opts ...ValidationOptions) *ArraySchema {
	if max < 0 {
		panic("maxItems cannot be negative")
	}
	a.maxItems = &max
	return a
}

func (a *ArraySchema) Description(desc string) *ArraySchema {
	a.description = &desc
	return a
}

func (a *ArraySchema) Optional() *ArraySchema {
	a.isOptional = true
	return a
}

func (a *ArraySchema) Clone() Schema {
	clone := &ArraySchema{
		schemaType:    a.schemaType,
		elementSchema: a.elementSchema.Clone(), // Clone the element schema
		isOptional:    a.isOptional,
		result:        &ValidationResult{},
	}

	// Deep copy pointers
	if a.minItems != nil {
		min := *a.minItems
		clone.minItems = &min
	}
	if a.maxItems != nil {
		max := *a.maxItems
		clone.maxItems = &max
	}
	if a.description != nil {
		desc := *a.description
		clone.description = &desc
	}

	// Deep copy the value slice if it exists
	if a.value != nil {
		clone.value = make([]interface{}, len(a.value))
		copy(clone.value, a.value)
	}

	return clone
}

func (a *ArraySchema) Validate() *ValidationResult {
	a.result = &ValidationResult{}

	if a.value == nil {
		if !a.isOptional {
			a.result.AddError(&ValidationError{
				Type:    RequiredArrayError,
				Message: "array is required",
			})
		}
		return a.result
	}

	// Check minItems
	if a.minItems != nil && len(a.value) < *a.minItems {
		a.result.AddError(&ValidationError{
			Type:     MinItemsError,
			Message:  fmt.Sprintf("minimum %d items required", *a.minItems),
			Expected: *a.minItems,
			Actual:   len(a.value),
		})
	}

	// Check maxItems
	if a.maxItems != nil && len(a.value) > *a.maxItems {
		a.result.AddError(&ValidationError{
			Type:     MaxItemsError,
			Message:  fmt.Sprintf("maximum %d items allowed", *a.maxItems),
			Expected: *a.maxItems,
			Actual:   len(a.value),
		})
	}

	// Validate each element
	for i, elem := range a.value {
		cloned := a.elementSchema.Clone()
		if err := cloned.setValue(elem); err != nil {
			a.result.AddError(&ValidationError{
				Type:    InvalidElementTypeError,
				Message: fmt.Sprintf("element %d: %v", i, err),
			})
			continue
		}

		if res := cloned.Validate(); res.HasErrors() {
			for _, err := range res.Errors {
				err.Message = fmt.Sprintf("element %d: %s", i, err.Message)
				a.result.AddError(err)
			}
		}
	}

	return a.result
}

func (a *ArraySchema) UnmarshalJSON(data []byte) error {
	a.result = &ValidationResult{}

	if string(data) == "null" {
		if !a.isOptional {
			return fmt.Errorf("array is required")
		}
		a.value = nil
		return nil
	}

	var rawElements []json.RawMessage
	if err := json.Unmarshal(data, &rawElements); err != nil {
		return fmt.Errorf("invalid array format: %w", err)
	}

	a.value = make([]interface{}, 0, len(rawElements))
	for i, elemData := range rawElements {
		elem := a.elementSchema.Clone()
		if err := elem.UnmarshalJSON(elemData); err != nil {
			a.result.AddError(&ValidationError{
				Type:    InvalidElementTypeError,
				Message: fmt.Sprintf("element %d: %v", i, err),
			})
			continue
		}

		if val, ok := elem.getValue(); ok {
			a.value = append(a.value, val)
		} else {
			a.result.AddError(&ValidationError{
				Type:    MissingElementValueError,
				Message: fmt.Sprintf("element %d: missing value", i),
			})
		}
	}

	return a.Validate().Error()
}

func (a *ArraySchema) MarshalJSON() ([]byte, error) {
	if a.value == nil {
		if a.isOptional {
			return json.Marshal(nil)
		}
		return nil, fmt.Errorf("required array has no value")
	}
	return json.Marshal(a.value)
}

func (a *ArraySchema) CompileJSONSchema(schema *jsonschema.JSONSchema, jsonTag string) error {
	itemsSchema := &jsonschema.JSONSchema{}
	if err := a.elementSchema.CompileJSONSchema(itemsSchema, ""); err != nil {
		return fmt.Errorf("failed to compile element schema: %w", err)
	}

	arraySchema := &jsonschema.JSONSchema{
		Type:  ArraySchemaType,
		Items: itemsSchema,
	}

	if a.description != nil {
		arraySchema.Description = *a.description
	}
	if a.minItems != nil {
		arraySchema.MinItems = a.minItems
	}
	if a.maxItems != nil {
		arraySchema.MaxItems = a.maxItems
	}

	schema.Properties[jsonTag] = arraySchema
	if !a.isOptional {
		schema.Required = append(schema.Required, jsonTag)
	}
	return nil
}

func (a *ArraySchema) Value() ([]interface{}, bool) {
	val, ok := a.getValue()
	if !ok {
		return nil, false
	}
	arrayVal, ok := val.([]interface{})
	if !ok {
		panic(fmt.Sprintf("ArraySchema: invalid internal value type %T, expected []interface{}", val))
	}
	return arrayVal, true
}

// Set provides a type-safe way to set array values
func (a *ArraySchema) Set(values ...interface{}) *ArraySchema {
	// reset validation results
	a.result = &ValidationResult{}

	// Create new value slice
	a.value = make([]interface{}, 0, len(values))

	// Validate each value against the element schema
	fmt.Printf("%v - %T\n", values, values)
	for i, val := range values {
		fmt.Printf("%v - %T\n", val, val)
		elem := a.elementSchema.Clone()
		if err := elem.setValue(val); err != nil {
			a.result.AddError(&ValidationError{
				Type:    InvalidElementTypeError,
				Message: fmt.Sprintf("element %d: %v", i, err),
			})
			continue
		}

		// If the value is valid for the element schema, add it
		if elemVal, ok := elem.getValue(); ok {
			a.value = append(a.value, elemVal)
		}
	}

	fmt.Printf("%v - %T\n", a.value, a.value)

	return a
}

func (a *ArraySchema) setValue(value interface{}) error {
	slice, ok := value.([]interface{})
	if !ok {
		return fmt.Errorf("expected array type, got %T", value)
	}
	a.value = slice
	return nil
}

func (a *ArraySchema) getValue() (interface{}, bool) {
	if a.value == nil {
		return nil, false
	}

	return a.value, true
}
