package gsv

// ValidationErrorType represents the specific type of validation error
type ValidationErrorType string

// ValidationError represents a single validation error with strong typing
type ValidationError struct {
	Type     ValidationErrorType
	Field    string      // The field path where the error occurred
	Message  string      // Human readable message
	Expected interface{} // The expected value/constraint
	Actual   interface{} // The actual value that failed validation

	// Could add more structured fields like:
	// Constraint interface{} // The specific constraint that failed
	// Metadata   map[string]interface{} // Any additional error context
}
