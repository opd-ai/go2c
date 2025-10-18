package llvm

import (
	"testing"
)

func TestParser(t *testing.T) {
	llvmIR := `
define i32 @add(i32 %a, i32 %b) {
entry:
  %result = add i32 %a, %b
  ret i32 %result
}

define void @main() {
entry:
  %0 = call i32 @add(i32 5, i32 3)
  ret void
}

declare void @external_func(i32)
`

	parser := NewParser()
	module, err := parser.Parse(llvmIR)
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	if len(module.Functions) != 3 {
		t.Errorf("Expected 3 functions, got %d", len(module.Functions))
	}

	// Check add function
	var addFunc *Function
	for _, fn := range module.Functions {
		if fn.Name == "add" {
			addFunc = fn
			break
		}
	}

	if addFunc == nil {
		t.Fatal("add function not found")
	}

	if addFunc.ReturnType != "i32" {
		t.Errorf("add return type = %q, want %q", addFunc.ReturnType, "i32")
	}

	if len(addFunc.Parameters) != 2 {
		t.Errorf("add parameters count = %d, want 2", len(addFunc.Parameters))
	}

	if len(addFunc.Body) < 2 {
		t.Errorf("add body lines = %d, want at least 2", len(addFunc.Body))
	}

	if addFunc.IsExternal {
		t.Error("add should not be external")
	}

	// Check external function
	var externalFunc *Function
	for _, fn := range module.Functions {
		if fn.Name == "external_func" {
			externalFunc = fn
			break
		}
	}

	if externalFunc == nil {
		t.Fatal("external_func not found")
	}

	if !externalFunc.IsExternal {
		t.Error("external_func should be external")
	}
}

func TestAnalyzer(t *testing.T) {
	module := &Module{
		Functions: []*Function{
			{Name: "main", ReturnType: "void", IsExternal: false},
			{Name: "add", ReturnType: "i32", IsExternal: false},
			{Name: "llvm.memcpy", ReturnType: "void", IsExternal: true},
			{Name: "external", ReturnType: "void", IsExternal: true},
		},
		Globals:  []*Global{},
		Types:    make(map[string]*Type),
		Metadata: make(map[string]string),
	}

	analyzer := NewAnalyzer(module)

	// Test GetUserFunctions
	userFuncs := analyzer.GetUserFunctions()
	if len(userFuncs) != 2 {
		t.Errorf("GetUserFunctions() count = %d, want 2", len(userFuncs))
	}

	// Test GetExternalFunctions
	externalFuncs := analyzer.GetExternalFunctions()
	if len(externalFuncs) != 1 {
		t.Errorf("GetExternalFunctions() count = %d, want 1", len(externalFuncs))
	}

	// Test GetMainFunction
	mainFunc, err := analyzer.GetMainFunction()
	if err != nil {
		t.Errorf("GetMainFunction() error = %v", err)
	}
	if mainFunc == nil || mainFunc.Name != "main" {
		t.Error("GetMainFunction() did not return main function")
	}

	// Test GetSummary
	summary := analyzer.GetSummary()
	if summary.TotalFunctions != 4 {
		t.Errorf("Summary.TotalFunctions = %d, want 4", summary.TotalFunctions)
	}
	if summary.UserFunctions != 2 {
		t.Errorf("Summary.UserFunctions = %d, want 2", summary.UserFunctions)
	}
	if summary.ExternalFunctions != 1 {
		t.Errorf("Summary.ExternalFunctions = %d, want 1", summary.ExternalFunctions)
	}
}
