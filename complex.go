package gsv

// Note (01/21/25): the complex Go primitives get their own schema type and do not use the
// generic NumberSchema[T] due to Go's cmp.Orderd generic not supporting complex
// numbers with comparison operators: "< <= >= >"
// https://pkg.go.dev/cmp@master#Ordered
//
// TODO - implement the complex number types
//
//type Complex64Schema = NumberSchema[complex64]
//type Complex128Schema = NumberSchema[complex128]
