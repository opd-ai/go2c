package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/opd-ai/go2c/internal/codegen"
	"github.com/opd-ai/go2c/internal/llvm"
	"github.com/opd-ai/go2c/internal/tinygo"
)

var (
	version = "0.1.0"
)

func main() {
	// Command-line flags
	inputFile := flag.String("input", "", "Input Go source file (required)")
	outputFile := flag.String("output", "", "Output C file (default: input file with .c extension)")
	verbose := flag.Bool("verbose", false, "Enable verbose output")
	showVersion := flag.Bool("version", false, "Show version information")
	keepLLVM := flag.Bool("keep-llvm", false, "Keep intermediate LLVM IR file")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: go2c [options]\n\n")
		fmt.Fprintf(os.Stderr, "go2c - Go to C transpiler using TinyGo and LLVM\n\n")
		fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nExample:\n")
		fmt.Fprintf(os.Stderr, "  go2c -input main.go -output main.c\n")
	}

	flag.Parse()

	if *showVersion {
		fmt.Printf("go2c version %s\n", version)
		os.Exit(0)
	}

	if *inputFile == "" {
		fmt.Fprintf(os.Stderr, "Error: -input flag is required\n\n")
		flag.Usage()
		os.Exit(1)
	}

	// Determine output file
	output := *outputFile
	if output == "" {
		ext := filepath.Ext(*inputFile)
		output = (*inputFile)[:len(*inputFile)-len(ext)] + ".c"
	}

	if *verbose {
		fmt.Printf("go2c v%s\n", version)
		fmt.Printf("Input file:  %s\n", *inputFile)
		fmt.Printf("Output file: %s\n", output)
		fmt.Println()
	}

	// Run the transpiler
	if err := transpile(*inputFile, output, *verbose, *keepLLVM); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	if *verbose {
		fmt.Println("\nTranspilation completed successfully!")
	}
}

func transpile(inputFile, outputFile string, verbose, keepLLVM bool) error {
	// Step 1: Compile Go to LLVM IR using TinyGo
	if verbose {
		fmt.Println("Step 1: Compiling Go to LLVM IR using TinyGo...")
	}

	compiler, err := tinygo.NewCompiler()
	if err != nil {
		return fmt.Errorf("failed to initialize TinyGo compiler: %w", err)
	}
	defer compiler.Cleanup()

	llvmIRPath, err := compiler.CompileToLLVMFile(inputFile)
	if err != nil {
		return fmt.Errorf("TinyGo compilation failed: %w", err)
	}

	if verbose {
		fmt.Printf("  Generated LLVM IR: %s\n", llvmIRPath)
	}

	// Keep LLVM IR if requested
	if keepLLVM {
		destPath := outputFile[:len(outputFile)-len(filepath.Ext(outputFile))] + ".ll"
		if err := copyFile(llvmIRPath, destPath); err != nil {
			fmt.Fprintf(os.Stderr, "Warning: failed to copy LLVM IR file: %v\n", err)
		} else if verbose {
			fmt.Printf("  Saved LLVM IR to: %s\n", destPath)
		}
	}

	// Step 2: Parse LLVM IR
	if verbose {
		fmt.Println("\nStep 2: Parsing LLVM IR...")
	}

	parser := llvm.NewParser()
	module, err := parser.ParseFile(llvmIRPath)
	if err != nil {
		return fmt.Errorf("LLVM IR parsing failed: %w", err)
	}

	if verbose {
		analyzer := llvm.NewAnalyzer(module)
		summary := analyzer.GetSummary()
		fmt.Printf("  Functions: %d (user: %d, external: %d)\n",
			summary.TotalFunctions, summary.UserFunctions, summary.ExternalFunctions)
		fmt.Printf("  Globals: %d\n", summary.GlobalVariables)
		fmt.Printf("  Types: %d\n", summary.Types)
	}

	// Step 3: Generate C code
	if verbose {
		fmt.Println("\nStep 3: Generating C code...")
	}

	emitter := codegen.NewEmitter()
	cCode, err := emitter.Emit(module)
	if err != nil {
		return fmt.Errorf("C code generation failed: %w", err)
	}

	// Step 4: Write C code to output file
	if verbose {
		fmt.Println("\nStep 4: Writing C code to output file...")
	}

	if err := os.WriteFile(outputFile, []byte(cCode), 0644); err != nil {
		return fmt.Errorf("failed to write output file: %w", err)
	}

	if verbose {
		fmt.Printf("  Written %d bytes to %s\n", len(cCode), outputFile)
	}

	return nil
}

func copyFile(src, dst string) error {
	data, err := os.ReadFile(src)
	if err != nil {
		return err
	}
	return os.WriteFile(dst, data, 0644)
}
