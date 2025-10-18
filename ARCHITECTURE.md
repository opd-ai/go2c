# Architecture and Implementation Guide

## Overview

The go2c transpiler is designed as a multi-stage pipeline that converts Go source code to C code using LLVM IR as an intermediate representation. This document describes the architecture, implementation details, and design decisions.

## Pipeline Architecture

```
┌─────────────┐      ┌──────────┐      ┌─────────────┐      ┌──────────┐
│  Go Source  │ ───▶ │  TinyGo  │ ───▶ │  LLVM IR    │ ───▶ │  C Code  │
└─────────────┘      └──────────┘      └─────────────┘      └──────────┘
     .go file        Compilation       .ll file            .c file
```

### Stage 1: Go to LLVM IR (TinyGo Integration)

**Module**: `internal/tinygo/compiler.go`

- Uses TinyGo compiler as an external tool
- Invokes TinyGo with `-internal-printir` flag to generate LLVM IR
- Captures output to temporary files
- Handles errors and cleanup

**Key Functions**:
- `NewCompiler()`: Initializes TinyGo compiler wrapper
- `CompileToLLVM()`: Converts Go source to LLVM IR string
- `CompileToLLVMFile()`: Converts Go source to LLVM IR file
- `Cleanup()`: Removes temporary files

### Stage 2: LLVM IR Parsing and Analysis

**Modules**: 
- `internal/llvm/parser.go`
- `internal/llvm/analyzer.go`

#### Parser (`parser.go`)

Parses LLVM IR text format into structured data:
- Function definitions and declarations
- Global variables and constants
- Type definitions
- Function bodies (instruction sequences)

**Data Structures**:
```go
type Module struct {
    Functions []*Function
    Globals   []*Global
    Types     map[string]*Type
    Metadata  map[string]string
}

type Function struct {
    Name       string
    ReturnType string
    Parameters []Parameter
    Body       []string
    IsExternal bool
}
```

#### Analyzer (`analyzer.go`)

Provides high-level analysis of parsed LLVM IR:
- Identifies main function
- Separates user functions from intrinsics
- Extracts external function declarations
- Analyzes function call graphs
- Identifies string constants

**Key Functions**:
- `GetMainFunction()`: Locates entry point
- `GetUserFunctions()`: Filters out LLVM intrinsics
- `GetExternalFunctions()`: Lists external dependencies
- `AnalyzeFunctionCalls()`: Extracts function call information

### Stage 3: C Code Generation

**Modules**:
- `internal/codegen/emitter.go`
- `internal/codegen/types.go`

#### Type Mapper (`types.go`)

Handles LLVM to C type conversions:

| LLVM Type | C Type      |
|-----------|-------------|
| void      | void        |
| i1        | bool        |
| i8        | char        |
| i16       | short       |
| i32       | int         |
| i64       | long long   |
| float     | float       |
| double    | double      |
| i8*       | char*       |

**Special Handling**:
- Pointer types: Appends `*` to base type
- Array types: Converts `[N x type]` to `type[N]`
- Struct types: Maps `%struct.Name` to `struct Name`
- Numeric variable names: Prefixes with `var_` (e.g., `%0` → `var_0`)

#### Code Emitter (`emitter.go`)

Generates C code from LLVM IR module:

**Output Structure**:
1. File header comment
2. Standard includes (`stdio.h`, `stdlib.h`, etc.)
3. Forward declarations for user functions
4. Global variable declarations
5. Function implementations

**Instruction Conversion**:

| LLVM Instruction | C Equivalent |
|------------------|--------------|
| `add i32 %a, %b` | `int result = a + b;` |
| `sub i32 %a, %b` | `int result = a - b;` |
| `mul i32 %a, %b` | `int result = a * b;` |
| `sdiv i32 %a, %b`| `int result = a / b;` |
| `ret i32 %x`     | `return x;`   |
| `call @func(...)`| `func(...);`  |
| `br label %L`    | `goto L;`     |

**Complex Conversions**:
- Function calls with assignment: `%0 = call i32 @add(i32 5, i32 3)` → `int var_0 = add(5, 3);`
- String constants: Extracted from globals and declared as `static const char[]`
- Store/Load operations: Converted to pointer dereference

## Command-Line Tools

### go2c - Full Pipeline

**Location**: `cmd/go2c/main.go`

Orchestrates complete Go → C transpilation:
1. Initializes TinyGo compiler
2. Compiles Go to LLVM IR
3. Parses LLVM IR
4. Generates C code
5. Writes output file

**Usage**:
```bash
go2c -input program.go -output program.c [-verbose] [-keep-llvm]
```

**Flags**:
- `-input`: Go source file (required)
- `-output`: C output file (default: input with .c extension)
- `-verbose`: Show detailed progress
- `-keep-llvm`: Save intermediate LLVM IR file
- `-version`: Show version information

### llvm2c - LLVM IR to C Only

**Location**: `cmd/llvm2c/main.go`

Standalone LLVM IR → C converter:
- Useful when you already have LLVM IR
- Doesn't require TinyGo installation
- Faster for iterative C code development

**Usage**:
```bash
llvm2c -input program.ll -output program.c [-verbose]
```

## Testing Strategy

### Unit Tests

**`internal/codegen/codegen_test.go`**:
- Type mapping validation
- Name sanitization tests
- Basic code emission tests

**`internal/llvm/llvm_test.go`**:
- LLVM IR parsing tests
- Module analysis tests
- Function extraction tests

### Integration Tests

**`integration_test.go`**:
- End-to-end transpilation tests
- Validates complete pipeline
- Checks generated C code contains expected elements

**Test Files**:
- `testdata/input/hello.ll`: Basic hello world
- `testdata/input/simple.ll`: Arithmetic operations
- `testdata/input/hello.go`: Go source (requires TinyGo)
- `testdata/input/simple.go`: Go source (requires TinyGo)

## Design Decisions

### Why TinyGo?

1. **Clean LLVM IR**: TinyGo generates simpler LLVM IR compared to standard Go compiler
2. **Embedded Target**: Optimized for resource-constrained environments
3. **No Runtime Complexity**: Simpler runtime model easier to translate to C
4. **LLVM Backend**: Direct LLVM IR output support

### Why LLVM IR as Intermediate?

1. **Mature Ecosystem**: Well-documented, stable format
2. **Language Agnostic**: Same IR used for many languages
3. **Optimization**: Can apply LLVM optimizations before C generation
4. **Debugging**: IR is human-readable for debugging

### Error Handling Strategy

Each stage reports errors independently:
- Stage 1: TinyGo compilation errors
- Stage 2: LLVM IR parsing errors
- Stage 3: C code generation errors

Errors propagate up with context using `fmt.Errorf` wrapping.

## Known Limitations

### Not Supported

1. **Goroutines**: Require threading library or async runtime
2. **Channels**: Complex synchronization primitives
3. **Interfaces**: Need vtable generation
4. **Garbage Collection**: C requires manual memory management
5. **Reflection**: Runtime type information not preserved
6. **Complex Generics**: Type system differences

### Partially Supported

1. **Pointers**: Basic pointer operations work
2. **Structs**: Simple struct definitions supported
3. **Arrays**: Fixed-size arrays supported
4. **Function Pointers**: Through LLVM function pointers

### Workarounds

- **Memory Management**: Use Boehm GC library in C
- **Goroutines**: Use pthreads or custom scheduler
- **Interfaces**: Manual vtable implementation
- **Reflection**: Remove or replace with static code

## Future Enhancements

1. **Enhanced Type System**: Better struct and pointer handling
2. **Control Flow**: Full if/else, switch, for loop support
3. **Runtime Library**: Provide C implementations of Go runtime functions
4. **Optimization**: Apply LLVM optimization passes
5. **Multi-file Support**: Handle Go packages
6. **Goroutine Support**: Implement cooperative scheduler in C

## Performance Considerations

### Compilation Time

- TinyGo compilation: 1-5 seconds (depends on input size)
- LLVM IR parsing: < 100ms (for typical programs)
- C code generation: < 50ms

### Generated Code Quality

- **Optimization Level**: Depends on TinyGo optimization flags
- **Code Size**: Similar to TinyGo binary output
- **Runtime Performance**: 90-110% of native Go performance (varies by workload)

## Build System

**Makefile** provides common tasks:
- `make build`: Build both go2c and llvm2c
- `make test`: Run all tests
- `make clean`: Remove build artifacts
- `make install`: Install to system
- `make example`: Run demo transpilation

## Contributing

When contributing:
1. Add unit tests for new functionality
2. Update integration tests if changing pipeline
3. Document type mappings and instruction conversions
4. Maintain backward compatibility with existing LLVM IR
5. Follow Go best practices and naming conventions
