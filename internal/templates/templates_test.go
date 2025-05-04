package templates

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
			result := extractReviewerID(tt.fullID)
			if result != tt.expected {
				t.Errorf("extractReviewerID(%q) = %q; want %q", tt.fullID, result, tt.expected)
			}
		})
	}
}
func TestAsDocResponses(t *testing.T) {
	tests := []struct {
		name     string
		headers  []string
		records  [][]string
		expected []Response
	}{
		{
			name:    "Single record with single header",
			headers: []string{"ID", "Comment", "Response"},
			records: [][]string{
				{"Rev1", "Some comment", "Some response"},
			},
			expected: []Response{
				{
					ReviewerID: "Rev1",
					Records: []Record{
						{Header: "ID", Text: "Rev1"},
						{Header: "Comment", Text: "Some comment"},
						{Header: "Response", Text: "Some response"},
					},
				},
			},
		},
		{
			name:    "Multiple records with multiple headers",
			headers: []string{"ID", "Comment", "Response"},
			records: [][]string{
				{"Rev2", "Another comment", "Another response"},
				{"Rev3", "Third comment", "Third response"},
			},
			expected: []Response{
				{
					ReviewerID: "Rev2",
					Records: []Record{
						{Header: "ID", Text: "Rev2"},
						{Header: "Comment", Text: "Another comment"},
						{Header: "Response", Text: "Another response"},
					},
				},
				{
					ReviewerID: "Rev3",
					Records: []Record{
						{Header: "ID", Text: "Rev3"},
						{Header: "Comment", Text: "Third comment"},
						{Header: "Response", Text: "Third response"},
					},
				},
			},
		},
		{
			name:    "Contains empty records",
			headers: []string{"ID", "Comment", "Response"},
			records: [][]string{
				{"Rev2", "Another comment", "Another response"},
				{},
			},
			expected: []Response{
				{
					ReviewerID: "Rev2",
					Records: []Record{
						{Header: "ID", Text: "Rev2"},
						{Header: "Comment", Text: "Another comment"},
						{Header: "Response", Text: "Another response"},
					},
				},
				{
					ReviewerID: "",
					Records: []Record{},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := asDocResponses(tt.headers, tt.records)
			if len(result) != len(tt.expected) {
				t.Errorf("asDocResponses() length = %d; want %d", len(result), len(tt.expected))
			}
			for i, res := range result {
				if res.ReviewerID != tt.expected[i].ReviewerID {
					t.Errorf("asDocResponses()[%d].ReviewerID = %q; want %q", i, res.ReviewerID, tt.expected[i].ReviewerID)
				}
				if len(res.Records) != len(tt.expected[i].Records) {
					t.Errorf("asDocResponses()[%d].Records length = %d; want %d", i, len(res.Records), len(tt.expected[i].Records))
				}
				for j, rec := range res.Records {
					if rec.Header != tt.expected[i].Records[j].Header || rec.Text != tt.expected[i].Records[j].Text {
						t.Errorf("asDocResponses()[%d].Records[%d] = %+v; want %+v", i, j, rec, tt.expected[i].Records[j])
					}
				}
			}
		})
	}
}

