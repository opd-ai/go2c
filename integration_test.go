package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestLLVM2CIntegration(t *testing.T) {
	// Build the llvm2c binary
	buildCmd := exec.Command("go", "build", "-o", "/tmp/llvm2c_test", "./cmd/llvm2c")
	if err := buildCmd.Run(); err != nil {
		t.Fatalf("Failed to build llvm2c: %v", err)
	}
	defer os.Remove("/tmp/llvm2c_test")

	tests := []struct {
		name       string
		inputFile  string
		checkStrings []string
	}{
		{
			name:      "hello world",
			inputFile: "testdata/input/hello.ll",
			checkStrings: []string{
				"#include <stdio.h>",
				"void main(void)",
				"int _start(void)",
				"return 0;",
			},
		},
		{
			name:      "arithmetic operations",
			inputFile: "testdata/input/simple.ll",
			checkStrings: []string{
				"int add(int a, int b)",
				"int result = a + b;",
				"return result;",
			},
		},
		{
			name:      "control flow patterns",
			inputFile: "testdata/input/controlflow.ll",
			checkStrings: []string{
				"if (cmp)",
				"while (cmp)",
				"bool cmp = a > b;",
				"bool cmp = i_val < n;",
				"return a;",
				"return b;",
			},
		},
		{
			name:      "switch statement",
			inputFile: "testdata/input/switch.ll",
			checkStrings: []string{
				"switch (x)",
				"case 0:",
				"case 1:",
				"case 2:",
				"default:",
				"return 10;",
				"return 20;",
				"return 30;",
				"return -1;",
			},
		},
		{
			name:      "for loop",
			inputFile: "testdata/input/forloop.ll",
			checkStrings: []string{
				"int sum_range(int start, int end)",
				"int count_to_n(int n)",
				"for (;",
				"bool cmp =",
			},
		},
		{
			name:      "multiple return values",
			inputFile: "testdata/input/multiret.ll",
			checkStrings: []string{
				"typedef struct {",
				"multi_return_1_t foo(int x)",
				"multi_return_1_t result;",
				"result.field0 = x;",
				"result2.field1 = 42;",
				"return result2;",
				"multi_return_1_t var_0 = foo(10);",
				"int val1 = var_0.field0;",
				"int val2 = var_0.field1;",
			},
		},
		{
			name:      "do-while loop",
			inputFile: "testdata/input/dowhile.ll",
			checkStrings: []string{
				"int sum_do_while(int n)",
				"do {",
				"} while (cmp);",
				"void print_countdown(int n)",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			outputFile := filepath.Join("/tmp", "test_output_"+strings.ReplaceAll(tt.name, " ", "_")+".c")
			defer os.Remove(outputFile)

			// Run llvm2c
			cmd := exec.Command("/tmp/llvm2c_test", "-input", tt.inputFile, "-output", outputFile)
			output, err := cmd.CombinedOutput()
			if err != nil {
				t.Fatalf("llvm2c failed: %v\nOutput: %s", err, output)
			}

			// Read the generated C code
			cCode, err := os.ReadFile(outputFile)
			if err != nil {
				t.Fatalf("Failed to read output file: %v", err)
			}

			cCodeStr := string(cCode)

			// Check that all expected strings are present
			for _, expected := range tt.checkStrings {
				if !strings.Contains(cCodeStr, expected) {
					t.Errorf("Generated C code does not contain expected string: %q", expected)
				}
			}
		})
	}
}

// TestRuntimeLibraryCompilation tests that generated C code with runtime functions compiles
func TestRuntimeLibraryCompilation(t *testing.T) {
	// Build the llvm2c binary
	buildCmd := exec.Command("go", "build", "-o", "/tmp/llvm2c_test_rt", "./cmd/llvm2c")
	if err := buildCmd.Run(); err != nil {
		t.Fatalf("Failed to build llvm2c: %v", err)
	}
	defer os.Remove("/tmp/llvm2c_test_rt")

	// Generate C code from simple.ll (which uses runtime_printint)
	inputFile := "testdata/input/simple.ll"
	outputFile := "/tmp/test_runtime_compilation.c"
	defer os.Remove(outputFile)

	cmd := exec.Command("/tmp/llvm2c_test_rt", "-input", inputFile, "-output", outputFile)
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("llvm2c failed: %v\nOutput: %s", err, output)
	}

	// Read the generated C code
	cCode, err := os.ReadFile(outputFile)
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}

	cCodeStr := string(cCode)

	// Verify runtime library is included
	if !strings.Contains(cCodeStr, "Go Runtime Library") {
		t.Error("Generated C code should include runtime library")
	}

	if !strings.Contains(cCodeStr, "runtime_printint") {
		t.Error("Generated C code should include runtime_printint function")
	}

	// Create a wrapper to avoid _start conflict
	wrapperFile := "/tmp/test_runtime_wrapper.c"
	wrapperCode := `#define main generated_main
#define _start generated_start
#include "` + outputFile + `"
#undef main
#undef _start

int main(void) {
    generated_main();
    return 0;
}
`
	if err := os.WriteFile(wrapperFile, []byte(wrapperCode), 0644); err != nil {
		t.Fatalf("Failed to write wrapper file: %v", err)
	}
	defer os.Remove(wrapperFile)

	// Try to compile the generated C code
	exeFile := "/tmp/test_runtime_exe"
	defer os.Remove(exeFile)

	compileCmd := exec.Command("gcc", "-o", exeFile, wrapperFile)
	compileOutput, err := compileCmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Failed to compile generated C code: %v\nOutput: %s", err, compileOutput)
	}

	// Run the compiled program
	runCmd := exec.Command(exeFile)
	runOutput, err := runCmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Failed to run compiled program: %v\nOutput: %s", err, runOutput)
	}

	// Check that the output is correct (5 + 3 = 8)
	outputStr := strings.TrimSpace(string(runOutput))
	if outputStr != "8" {
		t.Errorf("Expected output '8', got '%s'", outputStr)
	}

	t.Logf("Successfully compiled and ran generated C code with runtime library")
	t.Logf("Output: %s", outputStr)
}
