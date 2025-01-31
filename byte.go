package gsv

// ByteSchema implements the Schema interface.
// In Go, a byte is an alias for uint8.
type ByteSchema = NumberSchema[uint8]

// Byte creates a new ByteSchema for uint8 (byte) values
func Byte() *ByteSchema {
	return (*ByteSchema)(Number[uint8]())
}
