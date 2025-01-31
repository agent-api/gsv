package gsv

// Unsigned integer Schema types for generic NumberSchemas.
// These provide concrete implementations for Go's built-in unsigned integer types:
// uint, uint8, uint16, uint32, uint64, uintptr
type (
	// UintSchema for uint
	UintSchema = NumberSchema[uint]

	// Uint8Schema for uint8
	Uint8Schema = NumberSchema[uint8]

	// Uint16Schema for uint16
	Uint16Schema = NumberSchema[uint16]

	// Uint32Schema for uint32
	Uint32Schema = NumberSchema[uint32]

	// Uint64Schema for uint64
	Uint64Schema = NumberSchema[uint64]

	// UintptrSchema for uintptr
	UintptrSchema = NumberSchema[uintptr]
)

// Uint creates a new UintSchema
func Uint() *UintSchema {
	return Number[uint]()
}

// Uint8 creates a new Uint8Schema
func Uint8() *Uint8Schema {
	return Number[uint8]()
}

// Uint16 creates a new Uint16Schema
func Uint16() *Uint16Schema {
	return Number[uint16]()
}

// Uint32 creates a new Uint32Schema
func Uint32() *Uint32Schema {
	return Number[uint32]()
}

// Uint64 creates a new Uint64Schema
func Uint64() *Uint64Schema {
	return Number[uint64]()
}

// Uintptr creates a new UintptrSchema
func Uintptr() *UintptrSchema {
	return Number[uintptr]()
}
