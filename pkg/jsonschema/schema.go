package jsonschema

// JSONSchema represents the structure of a JSON Schema
type JSONSchema struct {
	// Metadata
	Title       string                 `json:"title,omitempty"`
	Description string                 `json:"description,omitempty"`
	ID          string                 `json:"$id,omitempty"`
	Schema      string                 `json:"$schema,omitempty"`
	Definitions map[string]*JSONSchema `json:"definitions,omitempty"`

	// Core
	Type string `json:"type"`

	// Object validators
	Properties           map[string]*JSONSchema `json:"properties,omitempty"`
	Required             []string               `json:"required,omitempty"`
	AdditionalProperties *JSONSchema            `json:"additionalProperties,omitempty"`
	PatternProperties    map[string]*JSONSchema `json:"patternProperties,omitempty"`
	MinProperties        *int                   `json:"minProperties,omitempty"`
	MaxProperties        *int                   `json:"maxProperties,omitempty"`

	// String validators
	MinLength *int   `json:"minLength,omitempty"`
	MaxLength *int   `json:"maxLength,omitempty"`
	Pattern   string `json:"pattern,omitempty"`
	Format    string `json:"format,omitempty"`

	// Number validators
	Minimum          *float64 `json:"minimum,omitempty"`
	Maximum          *float64 `json:"maximum,omitempty"`
	ExclusiveMinimum *float64 `json:"exclusiveMinimum,omitempty"`
	ExclusiveMaximum *float64 `json:"exclusiveMaximum,omitempty"`
	MultipleOf       *float64 `json:"multipleOf,omitempty"`

	// Array validators
	Items       *JSONSchema `json:"items,omitempty"`
	MinItems    *int        `json:"minItems,omitempty"`
	MaxItems    *int        `json:"maxItems,omitempty"`
	UniqueItems *bool       `json:"uniqueItems,omitempty"`

	// Generic validators
	Enum  []interface{} `json:"enum,omitempty"`
	Const interface{}   `json:"const,omitempty"`
	AllOf []*JSONSchema `json:"allOf,omitempty"`
	AnyOf []*JSONSchema `json:"anyOf,omitempty"`
	OneOf []*JSONSchema `json:"oneOf,omitempty"`
	Not   *JSONSchema   `json:"not,omitempty"`
}

