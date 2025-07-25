package templates

import (
	"reflect"
	"testing"
)

func TestNewTemplate_ReturnsLatexTemplateByDefault(t *testing.T) {
	template := NewTemplate("")
	if template == nil {
		t.Fatal("Expected non-nil template")
	}
	res := reflect.TypeOf(template).String()
	if res != "*latex.Latex" {
		t.Errorf("Expected template type '*latex.Latex', got '%s'", res)
	}
}

func TestNewTemplate_ReturnsLatexTemplate(t *testing.T) {
	template := NewTemplate("LaTeX")
	if template == nil {
		t.Fatal("Expected non-nil template")
	}
	res := reflect.TypeOf(template).String()
	if res != "*latex.Latex" {
		t.Errorf("Expected template type '*latex.Latex', got '%s'", res)
	}
}

func TestNewTemplate_ReturnsTypstTemplate(t *testing.T) {
	template := NewTemplate("Typst")
	if template == nil {
		t.Fatal("Expected non-nil template")
	}
	res := reflect.TypeOf(template).String()
	if res != "*typst.Typst" {
		t.Errorf("Expected template type '*typst.Typst', got '%s'", res)
	}
}

func TestNewTemplate_ReturnsLatexTemplateForUnknownType(t *testing.T) {
	template := NewTemplate("Unknown")
	if template == nil {
		t.Fatal("Expected non-nil template")
	}
	res := reflect.TypeOf(template).String()
	if res != "*latex.Latex" {
		t.Errorf("Expected template type '*latex.Latex', got '%s'", res)
	}
}

func TestAvailable_ReturnsCorrectTemplateNames(t *testing.T) {
	expected := []string{"LaTeX", "Typst"}
	result := Available()

	if len(result) != len(expected) {
		t.Fatalf("Expected %d templates, got %d", len(expected), len(result))
	}

	for i, name := range expected {
		if result[i] != name {
			t.Errorf("Expected template name '%s', got '%s'", name, result[i])
		}
	}
}
