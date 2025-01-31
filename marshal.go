package gsv

import (
	"encoding/json"
	"fmt"
)

// SafeMarshal takes any struct type, calls "Ensure" with it to validate
// the stored schema values, and then marshals it to json, returning a byte slice
// and any error that occurred
func SafeMarshal(v any) ([]byte, error) {
	// First validate the struct
	result := ensure(&v)
	if result.HasErrors() {
		return nil, result.Error()
	}

	// Then marshal to JSON
	data, err := json.Marshal(v)
	if err != nil {
		return nil, fmt.Errorf("could not marshal to json: %w", err)
	}

	return data, nil
}
