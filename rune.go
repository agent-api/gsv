package gsv

// RuneSchema implements the Schema interface.
// In Go, a "rune" is an alias for int32 and represents a Unicode code point.
type RuneSchema = NumberSchema[int32]

// Byte creates a new ByteSchema for uint8 (byte) values
func Rune() *RuneSchema {
	return (*RuneSchema)(Number[int32]())
}
