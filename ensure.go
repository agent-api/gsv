package gsv

import (
	"fmt"
	"reflect"
)

// A function that takes a generic struct, iterates it members, checks if they're
// implementations of the gsv "schema" interface ... and calls their "validate"
// method

// ensure takes a generic, checks it's a gsv schema, and calls its validators.
// This function will recursively iterate all fields in a struct in order to
// validate all members of a schema, even the ones that haven't been set.
//
// Because JSON marshal and unmarshal won't be called for missing fields, this
// is needed in order to fully validate a schema.
//
// A ValidationResult is returned which wraps all errors and a boolean error signal
func ensure[T any](t T) *ValidationResult {
	return ensureRecursive(reflect.ValueOf(t), "")
}

func ensureRecursive(v reflect.Value, path string) *ValidationResult {
	// Initialize a new ValidationResult to collect all errors
	result := &ValidationResult{}

	// Handle pointers by getting their underlying value
	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return result
		}

		v = v.Elem()
	}

	fmt.Printf("Processing value of type: %v, kind: %v\n", v.Type(), v.Kind())

	// If it's a struct, process each field
	if v.Kind() == reflect.Struct {
		// First check if the struct itself implements GSVSchema
		if v.CanInterface() {
			fmt.Printf("Checking if %v implements GSVSchema\n", v.Type())
			if schema, ok := v.Interface().(Schema); ok {
				fmt.Printf("Found GSVSchema implementation: %+v\n", schema)
				if schemaResult := schema.Validate(); schemaResult.HasErrors() {
					for _, err := range schemaResult.Errors {
						if path != "" {
							err.Field = fmt.Sprintf("%s.%s", path, err.Field)
						}

						result.AddError(err)
					}
				}

			}
		}

		// Then process each field
		for i := 0; i < v.NumField(); i++ {
			field := v.Field(i)
			fieldType := v.Type().Field(i)

			// Build the field path
			fieldPath := fieldType.Name
			if path != "" {
				fieldPath = fmt.Sprintf("%s.%s", path, fieldType.Name)
			}

			fmt.Printf("Processing field: %s of type %v (kind: %v)\n",
				fieldType.Name, field.Type(), field.Kind())

			// Try to validate the field directly first
			if field.CanInterface() {
				fmt.Printf("Checking if field %s implements GSVSchema\n", fieldType.Name)
				if schema, ok := field.Interface().(Schema); ok {
					fmt.Printf("Found GSVSchema implementation on field %s\n", fieldType.Name)
					if fieldResult := schema.Validate(); fieldResult.HasErrors() {
						for _, err := range fieldResult.Errors {
							err.Field = fieldPath
							result.AddError(err)
						}
					}
				}
			}

			// Then recurse on the field
			if fieldResult := ensureRecursive(field, fieldPath); fieldResult.HasErrors() {
				result.Errors = append(result.Errors, fieldResult.Errors...)
			}
		}
	}

	return result
}
