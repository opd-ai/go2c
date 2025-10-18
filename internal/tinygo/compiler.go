package tinygo

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

// Compiler handles TinyGo compilation to LLVM IR
type Compiler struct {
	TinyGoPath string
	TempDir    string
}

// NewCompiler creates a new TinyGo compiler instance
func NewCompiler() (*Compiler, error) {
	tinygoPath, err := exec.LookPath("tinygo")
	if err != nil {
		return nil, fmt.Errorf("tinygo not found in PATH: %w", err)
	}

	tempDir, err := os.MkdirTemp("", "go2c-*")
	if err != nil {
		return nil, fmt.Errorf("failed to create temp directory: %w", err)
	}

	return &Compiler{
		TinyGoPath: tinygoPath,
		TempDir:    tempDir,
	}, nil
}

// CompileToLLVM compiles Go source code to LLVM IR
func (c *Compiler) CompileToLLVM(goSourcePath string) (string, error) {
	if _, err := os.Stat(goSourcePath); os.IsNotExist(err) {
		return "", fmt.Errorf("source file does not exist: %s", goSourcePath)
	}

	// Output LLVM IR file path
	llvmIRPath := filepath.Join(c.TempDir, "output.ll")

	// Run TinyGo to generate LLVM IR using -internal-printir
	// We redirect the output to a file
	cmd := exec.Command(c.TinyGoPath, "build", "-internal-printir", "-o", "/dev/null", goSourcePath)

	outFile, err := os.Create(llvmIRPath)
	if err != nil {
		return "", fmt.Errorf("failed to create output file: %w", err)
	}
	defer outFile.Close()

	cmd.Stdout = outFile
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("tinygo compilation failed: %w", err)
	}

	// Read the generated LLVM IR
	llvmIR, err := os.ReadFile(llvmIRPath)
	if err != nil {
		return "", fmt.Errorf("failed to read LLVM IR: %w", err)
	}

	return string(llvmIR), nil
}

// CompileToLLVMFile compiles Go source to LLVM IR and returns the file path
func (c *Compiler) CompileToLLVMFile(goSourcePath string) (string, error) {
	if _, err := os.Stat(goSourcePath); os.IsNotExist(err) {
		return "", fmt.Errorf("source file does not exist: %s", goSourcePath)
	}

	// Output LLVM IR file path
	llvmIRPath := filepath.Join(c.TempDir, "output.ll")

	// Run TinyGo to generate LLVM IR using -internal-printir
	cmd := exec.Command(c.TinyGoPath, "build", "-internal-printir", "-o", "/dev/null", goSourcePath)

	outFile, err := os.Create(llvmIRPath)
	if err != nil {
		return "", fmt.Errorf("failed to create output file: %w", err)
	}
	defer outFile.Close()

	cmd.Stdout = outFile
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("tinygo compilation failed: %w", err)
	}

	return llvmIRPath, nil
}

// Cleanup removes temporary files
func (c *Compiler) Cleanup() error {
	if c.TempDir != "" {
		return os.RemoveAll(c.TempDir)
	}
	return nil
}
