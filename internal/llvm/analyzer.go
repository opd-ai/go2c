package llvm

import (
	"fmt"
	"strings"
)

// Analyzer analyzes LLVM IR modules
type Analyzer struct {
	module *Module
}

// NewAnalyzer creates a new LLVM IR analyzer
func NewAnalyzer(module *Module) *Analyzer {
	return &Analyzer{
		module: module,
	}
}

// GetMainFunction finds the main function in the module
func (a *Analyzer) GetMainFunction() (*Function, error) {
	for _, fn := range a.module.Functions {
		if fn.Name == "main" {
			return fn, nil
		}
	}
	return nil, fmt.Errorf("main function not found")
}

// GetUserFunctions returns all user-defined functions (excluding LLVM intrinsics)
func (a *Analyzer) GetUserFunctions() []*Function {
	userFunctions := []*Function{}
	for _, fn := range a.module.Functions {
		// Skip LLVM intrinsics and runtime functions
		if !strings.HasPrefix(fn.Name, "llvm.") &&
			!strings.HasPrefix(fn.Name, "runtime.") &&
			!fn.IsExternal {
			userFunctions = append(userFunctions, fn)
		}
	}
	return userFunctions
}

// GetExternalFunctions returns all external function declarations
func (a *Analyzer) GetExternalFunctions() []*Function {
	externalFunctions := []*Function{}
	for _, fn := range a.module.Functions {
		if fn.IsExternal && !strings.HasPrefix(fn.Name, "llvm.") {
			externalFunctions = append(externalFunctions, fn)
		}
	}
	return externalFunctions
}

// AnalyzeFunctionCalls extracts function calls from a function body
func (a *Analyzer) AnalyzeFunctionCalls(fn *Function) []string {
	calls := []string{}
	for _, line := range fn.Body {
		if strings.Contains(line, "call") {
			// Extract function name from call instruction
			// Example: call void @puts(i8* %1)
			if strings.Contains(line, "@") {
				parts := strings.Split(line, "@")
				if len(parts) > 1 {
					funcName := strings.Split(parts[1], "(")[0]
					calls = append(calls, funcName)
				}
			}
		}
	}
	return calls
}

// GetStringConstants extracts string constants from globals
func (a *Analyzer) GetStringConstants() []string {
	constants := []string{}
	for _, global := range a.module.Globals {
		if strings.Contains(global.Value, "constant") && strings.Contains(global.Value, "c\"") {
			constants = append(constants, global.Name)
		}
	}
	return constants
}

// Summary provides a summary of the module
type Summary struct {
	TotalFunctions    int
	UserFunctions     int
	ExternalFunctions int
	GlobalVariables   int
	Types             int
}

// GetSummary returns a summary of the analyzed module
func (a *Analyzer) GetSummary() Summary {
	return Summary{
		TotalFunctions:    len(a.module.Functions),
		UserFunctions:     len(a.GetUserFunctions()),
		ExternalFunctions: len(a.GetExternalFunctions()),
		GlobalVariables:   len(a.module.Globals),
		Types:             len(a.module.Types),
	}
}
