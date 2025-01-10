package gsv

import (
	"testing"
)

func TestStringValidator_Min(t *testing.T) {
	tests := []struct {
		name     string
		value    string
		min      int
		opts     []ValidationOptions
		wantErrs bool
		errKey   string
		errMsg   string
	}{
		{
			name:     "valid length",
			value:    "hello",
			min:      3,
			wantErrs: false,
		},
		{
			name:     "too short",
			value:    "hi",
			min:      3,
			wantErrs: true,
			errKey:   "min",
			errMsg:   "must be at least 3 characters long",
		},
		{
			name:     "custom error message",
			value:    "hi",
			min:      3,
			opts:     []ValidationOptions{{Message: "too short!"}},
			wantErrs: true,
			errKey:   "min",
			errMsg:   "too short!",
		},
		{
			name:     "empty string",
			value:    "",
			min:      1,
			wantErrs: true,
			errKey:   "min",
			errMsg:   "must be at least 1 characters long",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := String().Min(tt.min, tt.opts...)
			errMap := v.Validate(tt.value)

			if tt.wantErrs && len(errMap) == 0 {
				t.Error("expected validation errors, got none")
			}
			if !tt.wantErrs && len(errMap) > 0 {
				t.Errorf("expected no validation errors, got %v", errMap)
			}
			if tt.wantErrs && len(errMap) > 0 && errMap[tt.errKey].Message != tt.errMsg {
				t.Errorf("expected error message %q, got %q", tt.errMsg, errMap[tt.errKey].Message)
			}
		})
	}
}

func TestStringValidator_Max(t *testing.T) {
	tests := []struct {
		name     string
		value    string
		max      int
		opts     []ValidationOptions
		wantErrs bool
		errKey   string
		errMsg   string
	}{
		{
			name:     "valid length",
			value:    "hello",
			max:      10,
			wantErrs: false,
		},
		{
			name:     "too long",
			value:    "hello world",
			max:      5,
			wantErrs: true,
			errKey:   "max",
			errMsg:   "must be at most 5 characters long",
		},
		{
			name:     "custom error message",
			value:    "hello world",
			max:      5,
			opts:     []ValidationOptions{{Message: "too long!"}},
			wantErrs: true,
			errKey:   "max",
			errMsg:   "too long!",
		},
		{
			name:     "exact length",
			value:    "hello",
			max:      5,
			wantErrs: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := String().Max(tt.max, tt.opts...)
			errs := v.Validate(tt.value)

			if tt.wantErrs && len(errs) == 0 {
				t.Error("expected validation errors, got none")
			}
			if !tt.wantErrs && len(errs) > 0 {
				t.Errorf("expected no validation errors, got %v", errs)
			}
			if tt.wantErrs && len(errs) > 0 && errs[tt.errKey].Message != tt.errMsg {
				t.Errorf("expected error message %q, got %q", tt.errMsg, errs[tt.errKey].Message)
			}
		})
	}
}

func TestStringValidator_MinMax(t *testing.T) {
	tests := []struct {
		name     string
		value    string
		min      int
		max      int
		wantErrs bool
		errMap   map[string]string
	}{
		{
			name:     "valid length",
			value:    "hello",
			min:      3,
			max:      10,
			wantErrs: false,
		},
		{
			name:     "too short",
			value:    "hi",
			min:      3,
			max:      10,
			wantErrs: true,
			errMap:   map[string]string{"min": "must be at least 3 characters long"},
		},
		{
			name:     "too long",
			value:    "hello world",
			min:      3,
			max:      5,
			wantErrs: true,
			errMap:   map[string]string{"max": "must be at most 5 characters long"},
		},
		{
			name:     "exact minimum",
			value:    "hey",
			min:      3,
			max:      5,
			wantErrs: false,
		},
		{
			name:     "exact maximum",
			value:    "hello",
			min:      3,
			max:      5,
			wantErrs: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := String().Min(tt.min).Max(tt.max)
			errs := v.Validate(tt.value)

			if tt.wantErrs && len(errs) == 0 {
				t.Error("expected validation errors, got none")
			}
			if !tt.wantErrs && len(errs) > 0 {
				t.Errorf("expected no validation errors, got %v", errs)
			}
			if tt.wantErrs && len(errs) > 0 {
				for key, expectedMsg := range tt.errMap {
					if errs[key].Message != expectedMsg {
						t.Errorf("expected error message %q, got %q", expectedMsg, errs[key].Message)
					}
				}
			}
		})
	}
}
