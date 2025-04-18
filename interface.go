package gsv

import "github.com/agent-api/gsv/pkg/jsonschema"

// Schema is the core interface for gsv.
//
// It defines the various primitives types that can store values, build schemas,
// and validate data.
type Schema interface {
	// Validate calls the schema's various registered validators and returns the
	// validation results.
	Validate() *ValidationResult

	// MarshalJSON implements the encoding/json Marshaler interface and enables
	// struct with Schemas and values to build valid JSON.
	MarshalJSON() ([]byte, error)

	// UnmarshalJSON implements the encoding/json Unmarshaler interface and enables
	// struct with Schemas and values to build valid JSON.
	UnmarshalJSON([]byte) error

	// IsOptional denotes that a given schema is optional.
	IsOptional() bool

	// TODO IsNotOptional
	// CompileJSONSchema for the JSON schema compiling
	CompileJSONSchema(schema *jsonschema.JSONSchema, jsonTag string) error

	// Clone creates a deep copy of the schema, including all validation rules
	// and current value. The new schema instance will be completely independent
	// from the original.
	Clone() Schema

	// Private method for internal use

	setValue(interface{}) error
	getValue() (interface{}, bool)
}
