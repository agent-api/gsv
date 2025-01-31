package gsv

import (
	"fmt"
	"strings"
)

// ValidationResult holds all validation errors for a schema
type ValidationResult struct {
	Errors []*ValidationError
}

// Helper methods for ValidationResult
func (vr *ValidationResult) HasErrors() bool {
	return len(vr.Errors) > 0
}

func (vr *ValidationResult) AddError(err *ValidationError) {
	vr.Errors = append(vr.Errors, err)
}

// Error converts the ValidationResult into a single error message
// This implements the error interface and provides a clean way to convert
// structured errors into a single error when needed
func (vr *ValidationResult) Error() error {
	if !vr.HasErrors() {
		return nil
	}

	var errMsgs []string
	for _, err := range vr.Errors {
		// Include field path if it exists
		if err.Field != "" {
			errMsgs = append(errMsgs, fmt.Sprintf("%s: [%s] %s",
				err.Field,
				err.Type,
				err.Message))
		} else {
			errMsgs = append(errMsgs, fmt.Sprintf("[%s] %s",
				err.Type,
				err.Message))
		}
	}

	return fmt.Errorf("validation failed: %s", strings.Join(errMsgs, "; "))
}
