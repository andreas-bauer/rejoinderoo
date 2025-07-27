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
	tests := []struct {
		name     string
		form     url.Values
		prefix   string
		expected []string
	}{
		{
			name:     "single match",
			form:     url.Values{"header-name": {"value1"}, "other": {"x"}},
			prefix:   "header-",
			expected: []string{"value1"},
		},
		{
			name:     "multiple matches",
			form:     url.Values{"header-a": {"v1"}, "header-b": {"v2"}, "header-c": {"v3"}},
			prefix:   "header-",
			expected: []string{"v1", "v2", "v3"},
		},
		{
			name:     "no matches",
			form:     url.Values{"foo": {"bar"}, "baz": {"qux"}},
			prefix:   "header-",
			expected: []string{},
		},
		{
			name:     "empty form",
			form:     url.Values{},
			prefix:   "header-",
			expected: []string{},
		},
		{
			name:     "multiple values per key",
			form:     url.Values{"header-x": {"a", "b"}, "header-y": {"c"}},
			prefix:   "header-",
			expected: []string{"a", "b", "c"},
		},
		{
			name:     "prefix is empty",
			form:     url.Values{"foo": {"bar"}, "baz": {"qux"}},
			prefix:   "",
			expected: []string{"bar", "qux"},
		},
	}

	// Helper to check if two slices have the same elements, regardless of order
	slicesEqualIgnoreOrder := func(a, b []string) bool {
		if len(a) != len(b) {
			return false
		}
		counts := make(map[string]int)
		for _, v := range a {
			counts[v]++
		}
		for _, v := range b {
			counts[v]--
		}
		for _, c := range counts {
			if c != 0 {
				return false
			}
		}
		return true
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getFormValuesWithPrefix(tt.form, tt.prefix)
			if !slicesEqualIgnoreOrder(got, tt.expected) {
				t.Errorf("got %v, want %v (order ignored)", got, tt.expected)
			}
		})
	}
}