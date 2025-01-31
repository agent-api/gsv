package gsv

// Integer schema type aliases for the generic NumberSchema.
// These provide concrete implementations for Go's built-in integer types.
type (
	// IntSchema validates int values
	IntSchema = NumberSchema[int]

	// Int8Schema validates int8 values
	Int8Schema = NumberSchema[int8]

	// Int16Schema validates int16 values
	Int16Schema = NumberSchema[int16]

	// Int32Schema validates int32 values
	Int32Schema = NumberSchema[int32]

	// Int64Schema validates int64 values
	Int64Schema = NumberSchema[int64]
)

// Int creates a new schema for validating int values
func Int() *IntSchema {
	return Number[int]()
}

// Int8 creates a new schema for validating int8 values
func Int8() *Int8Schema {
	return Number[int8]()
}

// Int16 creates a new schema for validating int16 values
func Int16() *Int16Schema {
	return Number[int16]()
}

// Int32 creates a new schema for validating int32 values
func Int32() *Int32Schema {
	return Number[int32]()
}

// Int64 creates a new schema for validating int64 values
func Int64() *Int64Schema {
	return Number[int64]()
}
