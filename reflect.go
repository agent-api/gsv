package gsv

import (
	"reflect"
)

// Helper functions
func isStringSchema(field reflect.Value) bool {
	_, ok := field.Interface().(*StringSchema)

	return ok
}

func isIntSchema(field reflect.Value) bool {
	_, ok := field.Interface().(*IntSchema)

	return ok
}

func isStructOrPtrToStruct(field reflect.Value) bool {
	typ := field.Type()

	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}

	return typ.Kind() == reflect.Struct
}
