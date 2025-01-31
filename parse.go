package gsv

import (
	"encoding/json"
	"fmt"
)

type ParseOptions struct {
	StopOnFirst bool           // Stop validation on first error
	SkipMissing bool           // Skip validation of missing fields
	ErrorMode   ValidationMode // How to handle errors
}

type ValidationMode int

const (
	ReturnAllErrors ValidationMode = iota
	ReturnFirstError
	PanicOnError
)

// TODO - T needs to just be a gsv schema type?

func Parse[T any](data []byte, t *T, opts ...ParseOptions) (*ValidationResult, error) {
	// First unmarshal the JSON
	if err := json.Unmarshal(data, t); err != nil {
		return nil, fmt.Errorf("could not unmarshal json: %w", err)
	}

	// todo - handle OPTS

	// Then validate the struct
	return ensure(*t), nil
}
