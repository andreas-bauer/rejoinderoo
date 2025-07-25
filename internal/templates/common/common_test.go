package common

import (
	"testing"
)

func TestExtractReviewerID(t *testing.T) {
	tests := []struct {
		name     string
		fullID   string
		expected string
	}{
		{
			name:     "Simple ID with colon",
			fullID:   "R1:123",
			expected: "R1",
		},
		{
			name:     "ID with dash and colon",
			fullID:   "R1-456:789",
			expected: "R1",
		},
		{
			name:     "ID with dot, dash, and colon",
			fullID:   "R2-456.789:123",
			expected: "R2",
		},
		{
			name:     "ID with only dot",
			fullID:   "R3.123",
			expected: "R3",
		},
		{
			name:     "ID with only dash",
			fullID:   "R1-123",
			expected: "R1",
		},
		{
			name:     "ID with no special characters",
			fullID:   "R1",
			expected: "R1",
		},
		{
			name:     "ID with space",
			fullID:   "R5 C1",
			expected: "R5",
		},
		{
			name:     "Empty ID",
			fullID:   "",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ExtractReviewerID(tt.fullID)
			if result != tt.expected {
				t.Errorf("extractReviewerID(%q) = %q; want %q", tt.fullID, result, tt.expected)
			}
		})
	}
}
