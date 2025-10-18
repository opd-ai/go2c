package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/opd-ai/go2c/internal/codegen"
	"github.com/opd-ai/go2c/internal/llvm"
)

func main() {
	inputFile := flag.String("input", "", "Input LLVM IR file (.ll)")
	outputFile := flag.String("output", "", "Output C file")
	verbose := flag.Bool("verbose", false, "Enable verbose output")

	flag.Parse()

	if *inputFile == "" {
		fmt.Fprintf(os.Stderr, "Error: -input flag is required\n")
		flag.Usage()
		os.Exit(1)
	}

	output := *outputFile
	if output == "" {
		ext := filepath.Ext(*inputFile)
		output = (*inputFile)[:len(*inputFile)-len(ext)] + ".c"
	}

	if *verbose {
		fmt.Printf("Input LLVM IR:  %s\n", *inputFile)
		fmt.Printf("Output C file:  %s\n", output)
		fmt.Println()
	}

	// Parse LLVM IR
	if *verbose {
		fmt.Println("Parsing LLVM IR...")
	}

	parser := llvm.NewParser()
	module, err := parser.ParseFile(*inputFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing LLVM IR: %v\n", err)
		os.Exit(1)
	}

	if *verbose {
		analyzer := llvm.NewAnalyzer(module)
		summary := analyzer.GetSummary()
		fmt.Printf("  Functions: %d (user: %d, external: %d)\n",
			summary.TotalFunctions, summary.UserFunctions, summary.ExternalFunctions)
		fmt.Printf("  Globals: %d\n", summary.GlobalVariables)
		fmt.Printf("  Types: %d\n", summary.Types)
		fmt.Println()
	}

	// Generate C code
	if *verbose {
		fmt.Println("Generating C code...")
	}

	emitter := codegen.NewEmitter()
	cCode, err := emitter.Emit(module)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error generating C code: %v\n", err)
		os.Exit(1)
	}

	// Write output
	if err := os.WriteFile(output, []byte(cCode), 0644); err != nil {
		fmt.Fprintf(os.Stderr, "Error writing output: %v\n", err)
		os.Exit(1)
	}

	if *verbose {
		fmt.Printf("\nSuccessfully generated %s (%d bytes)\n", output, len(cCode))
	}

	// Print a preview
	lines := strings.Split(cCode, "\n")
	if len(lines) > 30 {
		fmt.Println("\nPreview (first 30 lines):")
		for i := 0; i < 30; i++ {
			fmt.Println(lines[i])
		}
		fmt.Printf("... (%d more lines)\n", len(lines)-30)
	} else {
		fmt.Println("\nGenerated C code:")
		fmt.Println(cCode)
	}
}
