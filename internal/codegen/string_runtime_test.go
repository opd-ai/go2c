package codegen

import (
	"strings"
	"testing"
)

func TestGetStringRuntimeLibrary(t *testing.T) {
	library := GetStringRuntimeLibrary()

	// Check that library is not empty
	if library == "" {
		t.Fatal("String runtime library is empty")
	}

	// Check for key components
	requiredComponents := []string{
		"go_string_t",
		"go_string_new",
		"go_string_from_bytes",
		"go_string_compare",
		"go_string_equals",
		"go_string_concat",
		"go_string_substring",
		"go_string_to_cstr",
		"go_string_free",
		"go_string_at",
		"go_string_contains",
	}

	for _, component := range requiredComponents {
		if !strings.Contains(library, component) {
			t.Errorf("String runtime library missing component: %s", component)
		}
	}

	// Check that it has proper C syntax
	if !strings.Contains(library, "typedef struct") {
		t.Error("String runtime library missing struct definition")
	}

	if !strings.Contains(library, "static inline") {
		t.Error("String runtime library missing inline functions")
	}
}

func TestStringRuntimeLibraryStructure(t *testing.T) {
	library := GetStringRuntimeLibrary()

	// Verify structure has data and len fields
	if !strings.Contains(library, "const char* data") {
		t.Error("go_string_t structure missing 'data' field")
	}

	if !strings.Contains(library, "size_t len") {
		t.Error("go_string_t structure missing 'len' field")
	}

	// Verify it includes necessary headers (via comments)
	if !strings.Contains(library, "Go String Runtime Library") {
		t.Error("String runtime library missing header comment")
	}
}

func TestStringRuntimeLibraryFunctions(t *testing.T) {
	library := GetStringRuntimeLibrary()

	tests := []struct {
		name     string
		function string
		params   string
	}{
		{"new", "go_string_new", "const char* cstr"},
		{"from_bytes", "go_string_from_bytes", "const char* data, size_t len"},
		{"compare", "go_string_compare", "go_string_t a, go_string_t b"},
		{"equals", "go_string_equals", "go_string_t a, go_string_t b"},
		{"concat", "go_string_concat", "go_string_t a, go_string_t b"},
		{"substring", "go_string_substring", "go_string_t s, size_t start, size_t end"},
		{"to_cstr", "go_string_to_cstr", "go_string_t s"},
		{"free", "go_string_free", "char* data"},
		{"at", "go_string_at", "go_string_t s, size_t index"},
		{"contains", "go_string_contains", "go_string_t s, go_string_t substr"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !strings.Contains(library, tt.function) {
				t.Errorf("Function %s not found in library", tt.function)
			}
			if !strings.Contains(library, tt.params) {
				t.Errorf("Function %s parameters not correct", tt.function)
			}
		})
	}
}
