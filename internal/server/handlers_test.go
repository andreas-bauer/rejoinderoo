package server

import (
	"net/url"
	"testing"
)

func TestFileNameWithoutExtension(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"document.tex", "document"},
		{"document.typ", "document"},
		{"document.csv", "document"},
		{"archive.tar.gz", "archive.tar"},
		{"noextension", "noextension"},
		{".hiddenfile", ""},
		{"multi.part.name.txt", "multi.part.name"},
		{"trailingdot.", "trailingdot"},
		{"", ""},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := fileNameWithoutExtension(tt.input)
			if got != tt.expected {
				t.Errorf("fileNameWithoutExtension(%q) = %q; want %q", tt.input, got, tt.expected)
			}
		})
	}
}
func TestIsAllowedContentType(t *testing.T) {
	tests := []struct {
		contentType string
		expected    bool
	}{
		{"text/csv", true},
		{"application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", true},
		{"application/pdf", false},
		{"image/png", false},
		{"", false},
		{"text/plain", false},
	}

	for _, tt := range tests {
		t.Run(tt.contentType, func(t *testing.T) {
			got := isAllowedContentType(tt.contentType)
			if got != tt.expected {
				t.Errorf("isAllowedContentType(%q) = %v; want %v", tt.contentType, got, tt.expected)
			}
		})
	}
}
func TestGetFormValuesWithPrefix(t *testing.T) {
	headers := url.Values{
		"header-id":      []string{"on"},
		"header-comment": []string{"on"},
		"other-ignore":   []string{"on"},
	}

	got := getFormValuesWithPrefix(headers, "header-")
	if len(got) != 2 {
		t.Errorf("getFormValuesWithPrefix(headers, \"header-\") = %v; want 2 items", got)
	}

	containsID := false
	containsComment := false
	for _, v := range got {
		if v == "id" {
			containsID = true
		} else if v == "comment" {
			containsComment = true
		}
	}
	if !containsID || !containsComment {
		t.Errorf("getFormValuesWithPrefix(headers, \"header-\") = %v; want both id and comment", got)
	}
}
