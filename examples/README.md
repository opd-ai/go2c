# Examples

This directory contains example Go programs demonstrating various features and their corresponding transpiled C code.

## Basic Examples

### 1. Hello World
**File**: `../testdata/input/hello.go`

Simple hello world program demonstrating basic function calls and string constants.

### 2. Arithmetic Operations
**File**: `../testdata/input/simple.go`

Demonstrates basic arithmetic operations (addition) and function calls.

## Advanced Examples

### 3. Fibonacci (fibonacci.go)
Recursive function demonstrating:
- Function recursion
- Conditional logic
- Integer arithmetic

### 4. Struct Operations (struct.go)
Demonstrates:
- Struct definition
- Method receivers
- Struct field access

### 5. Conditionals (conditionals.go)
Demonstrates:
- If-else statements
- Comparison operators
- Multiple function definitions

## Running Examples

### Using llvm2c (LLVM IR → C)
```bash
# First, you need LLVM IR files
# If you have TinyGo installed:
tinygo build -o example.ll -internal-printir example.go > example.ll

# Then transpile LLVM IR to C:
./bin/llvm2c -input example.ll -output example.c -verbose
```

### Using go2c (Go → C via TinyGo)
**Requires TinyGo to be installed**

```bash
# Direct Go to C transpilation:
./bin/go2c -input example.go -output example.c -verbose
```

## Limitations

Current implementation has the following limitations:
- Goroutines are not supported
- Channels are not supported
- Interfaces have limited support
- Complex pointer operations may not work correctly
- Package imports are limited
- Garbage collection in C requires manual memory management

## Testing Generated C Code

To compile and test generated C code:

```bash
# Generate C code
./bin/llvm2c -input testdata/input/simple.ll -output /tmp/simple.c

# Compile with GCC
gcc /tmp/simple.c -o /tmp/simple

# Note: You may need to implement runtime functions like runtime_printint
```
