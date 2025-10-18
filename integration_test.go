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
