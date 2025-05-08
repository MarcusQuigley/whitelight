package data

import (
	"encoding/json"
	"testing"
)

func TestRuntime_MarshalJSON(t *testing.T) {
	tests := []struct {
		name     string
		runtime  Runtime
		expected string
	}{
		{
			name:     "Zero minutes",
			runtime:  Runtime(0),
			expected: `"0 mins"`,
		},
		{
			name:     "Single minute",
			runtime:  Runtime(1),
			expected: `"1 mins"`,
		},
		{
			name:     "Multiple minutes",
			runtime:  Runtime(120),
			expected: `"120 mins"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.runtime.MarshalJSON()
			if err != nil {
				t.Errorf("MarshalJSON() error = %v", err)
				return
			}

			gotStr := string(got)

			if gotStr != tt.expected {
				t.Errorf("MarshalJSON() = %v, want %v", gotStr, tt.expected)
			}

			var unmarshaled string
			if err := json.Unmarshal(got, &unmarshaled); err != nil {
				t.Errorf("Failed to unmarshal JSON: %v", err)
			}
		})
	}
}
