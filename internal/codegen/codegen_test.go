package codegen

import (
	"testing"

	"github.com/opd-ai/go2c/internal/llvm"
)

func TestTypeMapper(t *testing.T) {
	tm := NewTypeMapper()

	tests := []struct {
		name     string
		llvmType string
		expected string
	}{
		{"void type", "void", "void"},
		{"i1 type", "i1", "bool"},
		{"i8 type", "i8", "char"},
		{"i16 type", "i16", "short"},
		{"i32 type", "i32", "int"},
		{"i64 type", "i64", "long long"},
		{"float type", "float", "float"},
		{"double type", "double", "double"},
		{"i8* type", "i8*", "char*"},
		{"void* type", "void*", "void*"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tm.LLVMTypeToC(tt.llvmType)
			if result != tt.expected {
				t.Errorf("LLVMTypeToC(%q) = %q, want %q", tt.llvmType, result, tt.expected)
			}
		})
	}
}

func TestSanitizeName(t *testing.T) {
	tm := NewTypeMapper()

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"global with @", "@main", "main"},
		{"local with %", "%result", "result"},
		{"quoted name", "\"main$string\"", "main_string"},
		{"name with dots", "main.func", "main_func"},
		{"name with dash", "my-var", "my_var"},
		{"numeric start", "%0", "var_0"},
		{"numeric name", "123", "var_123"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tm.SanitizeName(tt.input)
			if result != tt.expected {
				t.Errorf("SanitizeName(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestEmitter(t *testing.T) {
	module := &llvm.Module{
		Functions: []*llvm.Function{
			{
				Name:       "add",
				ReturnType: "i32",
				Parameters: []llvm.Parameter{
					{Type: "i32", Name: "%a"},
					{Type: "i32", Name: "%b"},
				},
				Body:       []string{"%result = add i32 %a, %b", "ret i32 %result"},
				IsExternal: false,
			},
		},
		Globals:  []*llvm.Global{},
		Types:    make(map[string]*llvm.Type),
		Metadata: make(map[string]string),
	}

	emitter := NewEmitter()
	cCode, err := emitter.Emit(module)
	if err != nil {
		t.Fatalf("Emit() error = %v", err)
	}

	// Check that the C code contains expected elements
	expectedStrings := []string{
		"#include <stdio.h>",
		"#include <stdint.h>",
		"int add(int a, int b)",
		"int result = a + b;",
		"return result;",
	}

	for _, expected := range expectedStrings {
		if !contains(cCode, expected) {
			t.Errorf("Generated C code missing expected string: %q", expected)
		}
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) &&
		(findSubstring(s, substr) >= 0))
}

func findSubstring(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}
