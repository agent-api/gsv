package gsv

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/agent-api/gsv/pkg/jsonschema"
)

type CompileSchemaOpts struct {
	SchemaTitle       string
	SchemaDescription string
}

// CompileSchema converts a gsv schema struct into a JSON Schema
func CompileSchema(schema interface{}, cso *CompileSchemaOpts) ([]byte, error) {
	jsonSchema := &jsonschema.JSONSchema{
		Title:       cso.SchemaTitle,
		Description: cso.SchemaDescription,
		Type:        "object",
		Properties:  make(map[string]*jsonschema.JSONSchema),
		Required:    make([]string, 0),
	}

	if err := compileFields(jsonSchema, schema); err != nil {
		return nil, err
	}

	return json.MarshalIndent(jsonSchema, "", "  ")
}

// compileFields handles recursive field compilation
func compileFields(schema *jsonschema.JSONSchema, value interface{}) error {
	val := reflect.ValueOf(value)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	typ := val.Type()

	// Iterate over struct fields
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldType := typ.Field(i)

		// Get JSON tag
		jsonTag := fieldType.Tag.Get("json")
		if jsonTag == "" || jsonTag == "-" {
			continue
		}

		// Handle different types of fields
		switch {
		case isStringSchema(field):
			stringSchema, _ := field.Interface().(*StringSchema)
			if err := stringSchema.CompileJSONSchema(schema, jsonTag); err != nil {
				return err
			}

		case isIntSchema(field):
			// TODO - need to test
			intSchema, _ := field.Interface().(*IntSchema)
			if err := intSchema.CompileJSONSchema(schema, jsonTag); err != nil {
				return err
			}

		case isStructOrPtrToStruct(field):
			// Create a new nested object schema
			nestedSchema := &jsonschema.JSONSchema{
				Type:       "object",
				Properties: make(map[string]*jsonschema.JSONSchema),
				Required:   make([]string, 0),
			}

			// Recursively compile the nested struct
			if err := compileFields(nestedSchema, field.Interface()); err != nil {
				return err
			}

			schema.Properties[jsonTag] = nestedSchema
			schema.Required = append(schema.Required, jsonTag)

		default:
			return fmt.Errorf("unsupported schema type for field %s", fieldType.Name)
		}
	}

	return nil
}
