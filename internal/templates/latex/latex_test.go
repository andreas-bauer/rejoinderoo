package latex

import (
	"testing"
)

func TestAsDocresponses(t *testing.T) {
	tests := []struct {
		name     string
		headers  []string
		records  [][]string
		expected []response
	}{
		{
			name:    "Single record with single header",
			headers: []string{"ID", "Comment", "response"},
			records: [][]string{
				{"Rev1", "Some comment", "Some response"},
			},
			expected: []response{
				{
					ReviewerID: "Rev1",
					Records: []record{
						{Header: "ID", Text: "Rev1"},
						{Header: "Comment", Text: "Some comment"},
						{Header: "response", Text: "Some response"},
					},
				},
			},
		},
		{
			name:    "Multiple records with multiple headers",
			headers: []string{"ID", "Comment", "response"},
			records: [][]string{
				{"Rev2", "Another comment", "Another response"},
				{"Rev3", "Third comment", "Third response"},
			},
			expected: []response{
				{
					ReviewerID: "Rev2",
					Records: []record{
						{Header: "ID", Text: "Rev2"},
						{Header: "Comment", Text: "Another comment"},
						{Header: "response", Text: "Another response"},
					},
				},
				{
					ReviewerID: "Rev3",
					Records: []record{
						{Header: "ID", Text: "Rev3"},
						{Header: "Comment", Text: "Third comment"},
						{Header: "response", Text: "Third response"},
					},
				},
			},
		},
		{
			name:    "Contains empty records",
			headers: []string{"ID", "Comment", "response"},
			records: [][]string{
				{"Rev2", "Another comment", "Another response"},
				{},
			},
			expected: []response{
				{
					ReviewerID: "Rev2",
					Records: []record{
						{Header: "ID", Text: "Rev2"},
						{Header: "Comment", Text: "Another comment"},
						{Header: "response", Text: "Another response"},
					},
				},
				{
					ReviewerID: "",
					Records:    []record{},
				},
			},
		},
		{
			name:    "Contains records with length less than headers",
			headers: []string{"ID", "Comment", "Response"},
			records: [][]string{
				{"Rev2", "Another comment"},
				{},
			},
			expected: []response{
				{
					ReviewerID: "Rev2",
					Records: []record{
						{Header: "ID", Text: "Rev2"},
						{Header: "Comment", Text: "Another comment"},
						{Header: "Response", Text: ""},
					},
				},
				{
					ReviewerID: "",
					Records:    []record{},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := asDocResponses(tt.headers, tt.records)
			if len(result) != len(tt.expected) {
				t.Errorf("asDocresponses() length = %d; want %d", len(result), len(tt.expected))
			}
			for i, res := range result {
				if res.ReviewerID != tt.expected[i].ReviewerID {
					t.Errorf("asDocresponses()[%d].ReviewerID = %q; want %q", i, res.ReviewerID, tt.expected[i].ReviewerID)
				}
				if len(res.Records) != len(tt.expected[i].Records) {
					t.Errorf("asDocresponses()[%d].Records length = %d; want %d", i, len(res.Records), len(tt.expected[i].Records))
				}
				for j, rec := range res.Records {
					if rec.Header != tt.expected[i].Records[j].Header || rec.Text != tt.expected[i].Records[j].Text {
						t.Errorf("asDocresponses()[%d].Records[%d] = %+v; want %+v", i, j, rec, tt.expected[i].Records[j])
					}
				}
			}
		})
	}
}
