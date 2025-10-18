# go2c - Go to C Transpiler

A standalone Go-to-C transpiler that leverages TinyGo's LLVM bytecode generation as an intermediate representation.

## Overview

`go2c` is a tool that converts Go source code to C code by utilizing TinyGo to generate LLVM IR (Intermediate Representation) and then transforming that LLVM IR into equivalent C code. This provides a complete, autonomous conversion pipeline implemented entirely in Go.

## Architecture

The transpiler follows a three-stage pipeline:

```
Go Source → TinyGo → LLVM IR → C Code
```

1. **TinyGo Integration Module**: Compiles Go source to LLVM IR
2. **LLVM IR Processing Module**: Parses and analyzes LLVM bytecode
3. **C Code Generation Module**: Converts LLVM IR to C code

## Prerequisites

- Go 1.18 or later
- TinyGo compiler (https://tinygo.org/getting-started/install/)
- GCC or Clang (for compiling generated C code)

## Installation

### From Source

```bash
git clone https://github.com/opd-ai/go2c.git
cd go2c
go build -o go2c ./cmd/go2c
```

### Using go install

```bash
go install github.com/opd-ai/go2c/cmd/go2c@latest
```

## Usage

Basic usage:

```bash
go2c -input main.go -output main.c
```

### Command-Line Options

- `-input <file>`: Input Go source file (required)
- `-output <file>`: Output C file (default: input file with .c extension)
- `-verbose`: Enable verbose output showing compilation steps
- `-keep-llvm`: Keep intermediate LLVM IR file for debugging
- `-version`: Show version information

### Examples

1. **Simple transpilation:**
   ```bash
   go2c -input hello.go -output hello.c
   ```

2. **Verbose mode with LLVM IR preservation:**
   ```bash
   go2c -input program.go -output program.c -verbose -keep-llvm
   ```

3. **Compile generated C code:**
   ```bash
   go2c -input hello.go -output hello.c
   gcc hello.c -o hello
   ./hello
   ```

## Project Structure

```
go2c/
├── cmd/
│   └── go2c/
│       └── main.go          # CLI entry point
├── internal/
│   ├── tinygo/
│   │   └── compiler.go      # TinyGo integration
│   ├── llvm/
│   │   ├── parser.go        # LLVM IR parsing
│   │   └── analyzer.go      # IR analysis
│   └── codegen/
│       ├── emitter.go       # C code generation
│       └── types.go         # Type mapping
├── testdata/
│   ├── input/               # Test Go files
│   └── expected/            # Expected C outputs
├── go.mod
└── README.md
```

## Supported Features

### Currently Supported
- Basic functions and function calls
- Integer types (int, int8, int16, int32, int64)
- Boolean types
- Simple control flow (return statements)
- Global variables
- String constants

### Limitations
- **Goroutines**: Not yet supported (requires C threading library integration)
- **Channels**: Not supported
- **Interfaces**: Limited support (requires vtable generation)
- **Garbage Collection**: Generated C code uses manual memory management
- **Complex types**: Maps, slices, and channels have limited support
- **Package imports**: Currently focuses on single-file programs

## Development

### Building the Project

```bash
go build ./cmd/go2c
```

### Running Tests

```bash
go test ./...
```

### Testing with Sample Programs

```bash
# Transpile a test program
./go2c -input testdata/input/hello.go -output /tmp/hello.c -verbose

# Compile and run the generated C code
gcc /tmp/hello.c -o /tmp/hello
/tmp/hello
```

## How It Works

### Stage 1: Go to LLVM IR

The transpiler uses TinyGo to compile Go source code into LLVM IR. TinyGo is particularly well-suited for this task as it generates clean, optimized LLVM bytecode.

### Stage 2: LLVM IR Analysis

The LLVM IR is parsed to extract:
- Function definitions and signatures
- Data types and structures
- Control flow patterns
- Memory operations
- Global variables

### Stage 3: C Code Generation

The parsed LLVM IR is transformed into C code by:
- Mapping LLVM types to C types
- Converting LLVM instructions to C statements
- Generating appropriate headers and includes
- Handling name mangling for compatibility

## Quality Criteria

✓ Successfully compiles valid Go programs to LLVM IR  
✓ Generated C code compiles without errors using gcc/clang  
✓ Handles basic Go features (functions, types, control flow)  
✓ Provides clear error messages for unsupported features  
✓ Pure Go implementation  

## Contributing

Contributions are welcome! Please feel free to submit issues or pull requests.

### Development Priorities

1. Expand type system support
2. Improve control flow handling (if/else, loops)
3. Add struct and pointer support
4. Implement interface dispatch
5. Add comprehensive test suite

## License

See [LICENSE](LICENSE) file for details.

## References

- [TinyGo](https://tinygo.org/)
- [LLVM](https://llvm.org/)
- [LLVM Go Bindings](https://pkg.go.dev/tinygo.org/x/go-llvm)

## Acknowledgments

This project builds upon:
- TinyGo for LLVM IR generation
- LLVM for intermediate representation
- The Go community for inspiration and tools
