package codegen

import (
	"strings"
	"testing"

	"github.com/opd-ai/go2c/internal/llvm"
)

func TestGetRuntimeLibrary(t *testing.T) {
	lib := GetRuntimeLibrary()
	
	// Check that the library is non-empty
	if lib == "" {
		t.Fatal("Runtime library should not be empty")
	}
	
	// Check for expected content
	if !strings.Contains(lib, "Go Runtime Library") {
		t.Error("Runtime library should contain header comment")
	}
	
	if !strings.Contains(lib, "runtime_printint") {
		t.Error("Runtime library should contain runtime_printint")
	}
	
	if !strings.Contains(lib, "runtime_printstring") {
		t.Error("Runtime library should contain runtime_printstring")
	}
}

func TestRuntimeLibraryFunctions(t *testing.T) {
	lib := GetRuntimeLibrary()
	
	tests := []struct {
		name     string
		function string
	}{
		{"printint", "runtime_printint"},
		{"printstring", "runtime_printstring"},
		{"printuint32", "runtime_printuint32"},
		{"printuint64", "runtime_printuint64"},
		{"printint64", "runtime_printint64"},
		{"printfloat32", "runtime_printfloat32"},
		{"printfloat64", "runtime_printfloat64"},
		{"printbool", "runtime_printbool"},
		{"printpointer", "runtime_printpointer"},
		{"printnl", "runtime_printnl"},
		{"printspace", "runtime_printspace"},
		{"alloc", "runtime_alloc"},
		{"free", "runtime_free"},
		{"memcpy", "runtime_memcpy"},
		{"memset", "runtime_memset"},
		{"slicecopy", "runtime_slicecopy"},
		{"strcmp", "runtime_strcmp"},
		{"strlen", "runtime_strlen"},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !strings.Contains(lib, tt.function) {
				t.Errorf("Runtime library should contain %s", tt.function)
			}
		})
	}
}

func TestDetectRuntimeUsage(t *testing.T) {
	tests := []struct {
		name           string
		llvmIR         string
		expectRuntime  bool
	}{
		{
			name: "module_with_runtime_printint",
			llvmIR: `
define void @main() {
entry:
  call void @runtime.printint(i32 42)
  ret void
}

declare void @runtime.printint(i32)
`,
			expectRuntime: true,
		},
		{
			name: "module_with_runtime_printstring",
			llvmIR: `
define void @main() {
entry:
  call void @runtime_printstring(i8* @str, i32 5)
  ret void
}

declare void @runtime_printstring(i8*, i32)
`,
			expectRuntime: false,  // printstring is handled by string library, not runtime library
		},
		{
			name: "module_without_runtime",
			llvmIR: `
define i32 @add(i32 %a, i32 %b) {
entry:
  %result = add i32 %a, %b
  ret i32 %result
}
`,
			expectRuntime: false,
		},
		{
			name: "module_with_runtime_alloc",
			llvmIR: `
define i8* @allocate(i32 %size) {
entry:
  %ptr = call i8* @runtime.alloc(i32 %size)
  ret i8* %ptr
}

declare i8* @runtime.alloc(i32)
`,
			expectRuntime: true,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser := llvm.NewParser()
			module, err := parser.Parse(tt.llvmIR)
			if err != nil {
				t.Fatalf("Failed to parse LLVM IR: %v", err)
			}
			
			emitter := NewEmitter()
			emitter.detectRuntimeUsage(module)
			
			if emitter.needsRuntimeSupport != tt.expectRuntime {
				t.Errorf("Expected needsRuntimeSupport=%v, got %v", tt.expectRuntime, emitter.needsRuntimeSupport)
			}
		})
	}
}

func TestEmitWithRuntime(t *testing.T) {
	llvmIR := `
define void @test() {
entry:
  call void @runtime.printint(i32 42)
  ret void
}

declare void @runtime.printint(i32)
`
	
	parser := llvm.NewParser()
	module, err := parser.Parse(llvmIR)
	if err != nil {
		t.Fatalf("Failed to parse LLVM IR: %v", err)
	}
	
	emitter := NewEmitter()
	cCode, err := emitter.Emit(module)
	if err != nil {
		t.Fatalf("Failed to emit C code: %v", err)
	}
	
	// Check that runtime library is included
	if !strings.Contains(cCode, "Go Runtime Library") {
		t.Error("Generated C code should include runtime library")
	}
	
	if !strings.Contains(cCode, "runtime_printint") {
		t.Error("Generated C code should include runtime_printint function")
	}
	
	// Check that the function call is generated
	if !strings.Contains(cCode, "runtime_printint(42)") {
		t.Error("Generated C code should contain runtime_printint call")
	}
}

func TestEmitWithoutRuntime(t *testing.T) {
	llvmIR := `
define i32 @add(i32 %a, i32 %b) {
entry:
  %result = add i32 %a, %b
  ret i32 %result
}
`
	
	parser := llvm.NewParser()
	module, err := parser.Parse(llvmIR)
	if err != nil {
		t.Fatalf("Failed to parse LLVM IR: %v", err)
	}
	
	emitter := NewEmitter()
	cCode, err := emitter.Emit(module)
	if err != nil {
		t.Fatalf("Failed to emit C code: %v", err)
	}
	
	// Check that runtime library is NOT included
	if strings.Contains(cCode, "Go Runtime Library") {
		t.Error("Generated C code should not include runtime library when not needed")
	}
}
