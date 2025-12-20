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
				{"Rev1.1", "Some comment", "Some response"},
			},
			expected: []response{
				{
					ReviewerID: "Rev1",
					Records: []record{
						{Header: "ID", Text: "Rev1.1"},
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
				{"Rev3.3", "Third comment", "Third response"},
			},
			expected: []response{
				{
					ReviewerID: "Rev2",
					Records: []record{
						{Header: "ID", Text: "Rev2.1"},
						{Header: "Comment", Text: "Another comment"},
						{Header: "response", Text: "Another response"},
					},
				},
				{
					ReviewerID: "Rev3",
					Records: []record{
						{Header: "ID", Text: "Rev3.3"},
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
				{"Rev2.6", "Another comment", "Another response"},
				{},
			},
			expected: []response{
				{
					ReviewerID: "Rev2",
					Records: []record{
						{Header: "ID", Text: "Rev2.6"},
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
				{"Rev2.8", "Another comment"},
				{},
			},
			expected: []response{
				{
					ReviewerID: "Rev2",
					Records: []record{
						{Header: "ID", Text: "Rev2.8"},
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

func TestAsDocHeaders(t *testing.T) {
	tests := []struct {
		name     string
		headers  []string
		expected []header
	}{
		{
			name:    "Single header",
			headers: []string{"ID"},
			expected: []header{
				{Name: "ID", Idx: 2},
			},
		},
		{
			name:    "Multiple headers",
			headers: []string{"ID", "Comment", "Response"},
			expected: []header{
				{Name: "ID", Idx: 2},
				{Name: "Comment", Idx: 3},
				{Name: "Response", Idx: 4},
			},
		},
		{
			name:     "No headers",
			headers:  []string{},
			expected: []header{},
		},
		{
			name:    "Headers with special characters",
			headers: []string{"ID#", "Comment%", "Response&"},
			expected: []header{
				{Name: "ID#", Idx: 2},
				{Name: "Comment%", Idx: 3},
				{Name: "Response&", Idx: 4},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := asDocHeaders(tt.headers)
			if len(result) != len(tt.expected) {
				t.Errorf("asDocHeaders() length = %d; want %d", len(result), len(tt.expected))
			}
			for i, res := range result {
				if res.Name != tt.expected[i].Name || res.Idx != tt.expected[i].Idx {
					t.Errorf("asDocHeaders()[%d] = %+v; want %+v", i, res, tt.expected[i])
				}
			}
		})
	}
}

func TestLatexFileExtension(t *testing.T) {
	lt := NewLatexTemplate()
	got := lt.FileExtension()
	want := ".tex"
	if got != want {
		t.Errorf("FileExtension() = %q; want %q", got, want)
	}
}

func TestEscape(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"NoSpecialChars", "NoSpecialChars"},
		{"100% sure", "100\\% sure"},
		{"Price is $5", "Price is \\$5"},
		{"Use #hashtag", "Use \\#hashtag"},
		{"A_B_C", "A\\_B\\_C"},
		{"Ampersand & more", "Ampersand \\& more"},
		{"Tilde~Caret^", "Tilde\\textasciitilde{}Caret\\textasciicircum{}"},
		{"", ""},
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

func TestEscapeShouldStaySame(t *testing.T) {
	tests := []struct {
		input string
	}{
		{"A previous study by \\citet{bauer2025} showed that XZY."},
		{"A previous study by Bauer et al. showed that XZY \\cite{bauer2025}."},
		{"Software testing is a systematic process for the verification and validation of software against its specifications~\\cite{myers2012ArtSoftwareTesting, sommerville2016SoftwareEngineering,washizaki2024SWEBOKGuideSoftware}. Verification\\footnote{Verification: \\emph{``Are we building the product right?''}\\cite{boehm1984verifying}} ensures that the software meets its defined requirements, while validation\\footnote{Validation: \\emph{``Are we building the right product?''}~\\cite{boehm1984verifying}} ensures the software aligns with the customer’s expectations. It involves executing software with defined inputs and assessing the resulting outputs against expected outcomes."},
		{"Our artifacts are made publicly available on Zenodo\\footnote{\\href{https://zenodo.org}{https://zenodo.org}}"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := escape(tt.input)
			if got != tt.input {
				t.Errorf("escape(%q) = %q; want %q", tt.input, got, tt.input)
			}
		})
	}
}
