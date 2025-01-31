package gsv

// TODO -fix the comments
//
// Float schema type aliases for the generic NumberSchema.
// These provide concrete implementations for Go's built-in floating-point types.
type (
	// Float32Schema validates float32 values
	Float32Schema = NumberSchema[float32]

	// Float64Schema validates float64 values
	Float64Schema = NumberSchema[float64]
)

// Float32 creates a new schema for validating float32 values
func Float32() *Float32Schema {
	return Number[float32]()
}

// Float64 creates a new schema for validating float64 values
func Float64() *Float64Schema {
	return Number[float64]()
}

