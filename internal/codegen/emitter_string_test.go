package codegen

import (
	"strings"
	"testing"

	"github.com/opd-ai/go2c/internal/llvm"
)

func TestDetectStringUsage(t *testing.T) {
	tests := []struct {
		name     string
		module   *llvm.Module
		expected bool
	}{
		{
			name: "module with string constant",
			module: &llvm.Module{
				Globals: []*llvm.Global{
					{
						Name:  "@main_string",
						Type:  "[14 x i8]",
						Value: `constant [14 x i8] c"Hello, World!\00"`,
					},
				},
				Functions: []*llvm.Function{},
			},
			expected: true,
		},
		{
			name: "module with runtime_printstring call",
			module: &llvm.Module{
				Globals: []*llvm.Global{},
				Functions: []*llvm.Function{
					{
						Name:       "main",
						ReturnType: "void",
						Parameters: []llvm.Parameter{},
						Body: []string{
							`call void @runtime_printstring(i8* @main_string, i32 13)`,
							`ret void`,
						},
					},
				},
			},
			expected: true,
		},
		{
			name: "module without strings",
			module: &llvm.Module{
				Globals: []*llvm.Global{},
				Functions: []*llvm.Function{
					{
						Name:       "add",
						ReturnType: "i32",
						Parameters: []llvm.Parameter{
							{Type: "i32", Name: "%a"},
							{Type: "i32", Name: "%b"},
						},
						Body: []string{
							`%result = add i32 %a, %b`,
							`ret i32 %result`,
						},
					},
				},
			},
			expected: false,
		},
		{
			name: "module with array of i8 (potential string)",
			module: &llvm.Module{
				Globals: []*llvm.Global{
					{
						Name:  "@data",
						Type:  "[10 x i8]",
						Value: `constant [10 x i8] zeroinitializer`,
					},
				},
				Functions: []*llvm.Function{},
			},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			emitter := NewEmitter()
			emitter.detectStringUsage(tt.module)

			if emitter.needsStringSupport != tt.expected {
				t.Errorf("detectStringUsage() = %v, want %v", emitter.needsStringSupport, tt.expected)
			}
		})
	}
}

func TestEmitWithStrings(t *testing.T) {
	module := &llvm.Module{
		Globals: []*llvm.Global{
			{
				Name:  "@hello",
				Type:  "[6 x i8]",
				Value: `constant [6 x i8] c"Hello\00"`,
			},
		},
		Functions: []*llvm.Function{
			{
				Name:       "main",
				ReturnType: "void",
				Parameters: []llvm.Parameter{},
				Body: []string{
					`call void @runtime_printstring(i8* @hello, i32 5)`,
					`ret void`,
				},
			},
		},
	}

	emitter := NewEmitter()
	output, err := emitter.Emit(module)

	if err != nil {
		t.Fatalf("Emit() error = %v", err)
	}

	// Check that string library is included
	if !strings.Contains(output, "Go String Runtime Library") {
		t.Error("Generated C code does not include string runtime library")
	}

	// Check for key string types and functions
	requiredComponents := []string{
		"go_string_t",
		"go_string_new",
		"go_string_compare",
		"go_string_concat",
	}

	for _, component := range requiredComponents {
		if !strings.Contains(output, component) {
			t.Errorf("Generated C code missing: %s", component)
		}
	}

	// Verify basic structure
	if !strings.Contains(output, "#include") {
		t.Error("Generated C code missing includes")
	}

	if !strings.Contains(output, "void main") {
		t.Error("Generated C code missing main function")
	}
}

func TestEmitWithoutStrings(t *testing.T) {
	module := &llvm.Module{
		Globals:   []*llvm.Global{},
		Functions: []*llvm.Function{
			{
				Name:       "add",
				ReturnType: "i32",
				Parameters: []llvm.Parameter{
					{Type: "i32", Name: "%a"},
					{Type: "i32", Name: "%b"},
				},
				Body: []string{
					`%result = add i32 %a, %b`,
					`ret i32 %result`,
				},
			},
		},
	}

	emitter := NewEmitter()
	output, err := emitter.Emit(module)

	if err != nil {
		t.Fatalf("Emit() error = %v", err)
	}

	// Check that string library is NOT included
	if strings.Contains(output, "Go String Runtime Library") {
		t.Error("Generated C code includes string runtime library when it shouldn't")
	}

	// Check that the function is still generated
	if !strings.Contains(output, "int add") {
		t.Error("Generated C code missing add function")
	}
}
