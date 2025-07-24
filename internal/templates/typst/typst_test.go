package typst

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
				{"Rev1.11", "Some comment", "Some response"},
			},
			expected: []response{
				{
					ID:         "Rev1.11",
					ReviewerID: "Rev1",
					Records: []record{
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
				{"Rev2.1", "Another comment", "Another response"},
				{"Rev3.4", "Third comment", "Third response"},
			},
			expected: []response{
				{
					ID:         "Rev2.1",
					ReviewerID: "Rev2",
					Records: []record{
						{Header: "Comment", Text: "Another comment"},
						{Header: "response", Text: "Another response"},
					},
				},
				{
					ID:         "Rev3.4",
					ReviewerID: "Rev3",
					Records: []record{
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
				{"Rev2.4", "Another comment", "Another response"},
				{},
			},
			expected: []response{
				{
					ID:         "Rev2.4",
					ReviewerID: "Rev2",
					Records: []record{
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
				{"Rev3.14", "Another comment"},
				{},
			},
			expected: []response{
				{
					ID:         "Rev3.14",
					ReviewerID: "Rev3",
					Records: []record{
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

func TestTypstFileExtension(t *testing.T) {
	lt := NewTypstTemplate()
	got := lt.FileExtension()
	want := ".typ"
	if got != want {
		t.Errorf("FileExtension() = %q; want %q", got, want)
	}
}
func TestEscape(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"", ""},
		{"plain text", "plain text"},
		{"\\", "\\\\"},
		{"{", "\\{"},
		{"}", "\\}"},
		{"[", "\\["},
		{"]", "\\]"},
		{"#", "\\#"},
		{"\\{[#]}", "\\\\\\{\\[\\#\\]\\}"},
		{"Hello #1 [test] {ok}", "Hello \\#1 \\[test\\] \\{ok\\}"},
		{"multiple \\# special {chars}", "multiple \\\\\\# special \\{chars\\}"},
		{"no special chars", "no special chars"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := escape(tt.input)
			if got != tt.expected {
				t.Errorf("escape(%q) = %q; want %q", tt.input, got, tt.expected)
			}
		})
	}
}
