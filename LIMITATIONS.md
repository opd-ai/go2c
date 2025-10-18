# Known Limitations and Workarounds

This document describes the current limitations of the go2c transpiler and potential workarounds.

## Critical Limitations

### 1. Goroutines Not Supported

**Issue**: Go's goroutines cannot be directly translated to C without a runtime scheduler.

**Workaround Options**:
- Use pthread library for threading
- Implement a cooperative scheduler in C
- Refactor code to use callbacks instead of goroutines
- Use a coroutine library like libco

**Example**:
```go
// Go code (NOT SUPPORTED)
go func() {
    fmt.Println("Hello from goroutine")
}()

// C workaround with pthreads
void* thread_func(void* arg) {
    printf("Hello from thread\n");
    return NULL;
}

pthread_t thread;
pthread_create(&thread, NULL, thread_func, NULL);
```

### 2. Channels Not Supported

**Issue**: Go channels require complex synchronization that doesn't exist in C.

**Workaround Options**:
- Use pthread condition variables and mutexes
- Use message queues (POSIX mqueue)
- Implement a simple channel library in C

**Example**:
```c
// Simple channel-like structure
typedef struct {
    int buffer[100];
    int read_pos;
    int write_pos;
    pthread_mutex_t mutex;
    pthread_cond_t cond;
} channel_t;
```

### 3. Garbage Collection

**Issue**: Go uses garbage collection, C requires manual memory management.

**Workaround Options**:
- Use Boehm GC (conservative garbage collector for C)
- Implement reference counting
- Manual memory management with malloc/free
- Use arena allocators for specific use cases

**Integration**:
```c
// Using Boehm GC
#include <gc.h>

int main() {
    GC_INIT();
    char* ptr = GC_MALLOC(100);  // Garbage collected
    // No need to free
}
```

### 4. Interface Dispatch

**Issue**: Go interfaces use dynamic dispatch with type assertions.

**Workaround Options**:
- Implement vtables manually
- Use function pointers in structs
- Code generation for interface methods

**Example**:
```c
// Manual vtable implementation
typedef struct {
    void (*method1)(void*);
    void (*method2)(void*);
} interface_vtable;

typedef struct {
    void* data;
    interface_vtable* vtable;
} interface_t;
```

## Moderate Limitations

### 5. Package Imports

**Issue**: Go package system not fully supported.

**Status**: Single-file programs work best.

**Workaround**: Combine multiple Go files before transpilation or manually merge generated C files.

### 6. Complex Pointer Operations

**Issue**: Some pointer arithmetic and unsafe operations may not convert correctly.

**Examples of Issues**:
- `unsafe.Pointer` conversions
- Pointer arithmetic beyond basic operations
- Type punning

**Workaround**: Avoid complex pointer operations or manually adjust generated C code.

### 7. Maps and Slices

**Issue**: Dynamic data structures require runtime support.

**Status**: Basic arrays work, dynamic slices/maps need runtime.

**Workarounds**:
- Use fixed-size arrays where possible
- Implement hash table for maps
- Use dynamic arrays with manual memory management

**Example**:
```c
// Simple dynamic array
typedef struct {
    int* data;
    size_t length;
    size_t capacity;
} slice_int;

void slice_append(slice_int* s, int value) {
    if (s->length >= s->capacity) {
        s->capacity = s->capacity * 2;
        s->data = realloc(s->data, s->capacity * sizeof(int));
    }
    s->data[s->length++] = value;
}
```

### 8. Defer Statements

**Issue**: Go's defer mechanism requires stack unwinding support.

**Workaround Options**:
- Convert defer to explicit cleanup code
- Use `goto` for cleanup paths
- Implement using setjmp/longjmp

### 9. Error Handling (panic/recover)

**Issue**: Go's panic/recover mechanism requires exception handling.

**Workaround Options**:
- Use error return values
- Implement using setjmp/longjmp
- Manual error checking

### 10. Reflection

**Issue**: Runtime type information not preserved.

**Status**: Not supported, no practical workaround.

**Solution**: Refactor code to avoid reflection or use code generation.

## Minor Limitations

### 11. String Operations

**Issue**: Go strings are immutable and UTF-8, C strings are mutable byte arrays.

**Impact**: String manipulation may differ.

**Workaround**: Be careful with string modifications, consider using a string library.

### 12. Numeric Type Sizes

**Issue**: Go guarantees specific sizes (int32, int64), C types vary by platform.

**Workaround**: Use fixed-width integer types from `<stdint.h>` (int32_t, int64_t).

### 13. Control Flow Complexity

**Issue**: Some complex control flow patterns may not convert perfectly.

**Examples**:
- Labeled break/continue
- ~~Switch with fallthrough~~ **[RESOLVED]** - Switch statements now generate idiomatic C switch/case
- ~~Multiple return values~~ **[RESOLVED]** - Multiple return values now supported via aggregate types/structs

**Status**: Basic if/else, loops, and switch statements work well. Multiple return values are now converted to struct types. Complex patterns like labeled break/continue may need manual adjustment.

**Multiple Return Values Support** (NEW):

Go functions with multiple return values are now automatically converted to C functions that return a struct. This is particularly useful for error handling patterns.

**Example**:
```go
// Go code with multiple return values
func divide(a, b int) (int, error) {
    if b == 0 {
        return 0, errors.New("division by zero")
    }
    return a / b, nil
}

result, err := divide(10, 2)
if err != nil {
    // handle error
}
```

**Generated C Code**:
```c
// Struct definition for return values
typedef struct {
    int field0;      // First return value
    int field1;      // Second return value (error)
} divide_return_t;

// Function returns struct
divide_return_t divide(int a, int b) {
    divide_return_t result;
    if (b == 0) {
        result.field0 = 0;
        result.field1 = 1;  // error indicator
        return result;
    }
    result.field0 = a / b;
    result.field1 = 0;  // no error
    return result;
}

// Usage
divide_return_t ret = divide(10, 2);
int result = ret.field0;
int err = ret.field1;
if (err != 0) {
    // handle error
}
```

**How It Works**:
1. Functions with aggregate return types `{type1, type2, ...}` in LLVM IR are detected
2. A C struct type is automatically generated with fields for each return value
3. `insertvalue` instructions are converted to struct field assignments
4. `extractvalue` instructions are converted to struct field accesses
5. Function calls properly handle the struct return type

**Limitations**:
- Field names are generic (field0, field1, etc.) rather than semantic names
- Error types are represented as integers rather than full error objects
- Large return values may be less efficient than pass-by-reference

### 14. Method Receivers

**Issue**: Go methods with receivers need special handling.

**Status**: Partially supported.

**Workaround**: Methods become regular functions with explicit receiver parameter.

```c
// Go: func (p *Point) Distance() int
// C:  int Point_Distance(Point* p)
```

### 15. Anonymous Functions and Closures

**Issue**: Closures capture variables from outer scope.

**Status**: Not supported.

**Workaround**: Use function pointers with explicit context parameter.

## Platform-Specific Issues

### TinyGo Requirement

**Issue**: The go2c tool requires TinyGo to be installed.

**Impact**: Cannot transpile Go source directly without TinyGo.

**Workarounds**:
1. Install TinyGo from https://tinygo.org/
2. Use `llvm2c` tool with pre-generated LLVM IR
3. Generate LLVM IR using TinyGo separately

### LLVM Version Compatibility

**Issue**: Different LLVM versions may produce different IR formats.

**Status**: Tested with LLVM 18, should work with 14-18.

**Workaround**: Use compatible LLVM/TinyGo versions.

### C Compiler Differences

**Issue**: Generated C code may behave differently on GCC vs Clang.

**Impact**: Minimal for most code.

**Recommendation**: Test with both compilers.

## Code Generation Quality

### Label-Based Control Flow

**Issue**: The transpiler now intelligently detects common control flow patterns and generates idiomatic C code.

**Status**: **IMPROVED** - The code generator now recognizes and converts:
- **If-else patterns**: Conditional branches with both true and false paths
- **If-then patterns**: Conditional branches with only a true path
- **While loops**: Loop patterns with condition check and back-edges
- **For loops**: Loop patterns with initialization, condition, body, and increment blocks

**Generated Code Examples**:

```c
// If-else pattern (LLVM IR with conditional branch)
if (cmp) {
    return a;
} else {
    return b;
}

// While loop pattern (LLVM IR loop structure)
while (cmp) {
    // loop body
}

// For loop pattern (LLVM IR with init/cond/body/incr blocks)
for (; cmp; i = i + 1) {
    // loop body
}
```

**Remaining Issues**:
- Complex loop patterns (do-while) still use goto/labels
- ~~Switch statements not yet converted to C switch~~ **[RESOLVED]** - Now fully supported
- ~~For-loops still use goto/labels~~ **[RESOLVED]** - For-loop patterns now generate idiomatic C for-loops
- Nested patterns may fall back to goto in some cases

**Improvement**: Previously all control flow was converted to labels and goto statements, making generated code difficult to read. The enhanced code generator now produces more readable and maintainable C code for common patterns.

### Unoptimized Output

**Issue**: Generated C code may contain:
- Unnecessary temporary variables
- Redundant operations
- Verbose label usage

**Impact**: Larger code size, potentially slower execution.

**Mitigation**: Rely on C compiler optimization (`-O2`, `-O3`).

### Missing Comments

**Issue**: Generated C code lacks explanatory comments.

**Impact**: Can be harder to understand/debug generated code, though improved control flow readability helps.

**Workaround**: Use `-keep-llvm` flag to inspect LLVM IR, or review the structured control flow (if/while) that is now generated.

### Comparison Operations

**Status**: **SUPPORTED** - The transpiler now supports LLVM icmp instructions:
- `icmp eq` → `==`
- `icmp ne` → `!=`
- `icmp sgt/ugt` → `>`
- `icmp sge/uge` → `>=`
- `icmp slt/ult` → `<`
- `icmp sle/ule` → `<=`

Comparison results are converted to C `bool` type.

### Label-Based Control Flow

**Issue**: Uses labels and goto instead of if/while/for.

**Impact**: Less readable C code.

**Status**: Functional but not idiomatic C.

## Runtime Function Requirements

### Missing Runtime Functions

**Issue**: Generated C code calls runtime functions like:
- `runtime_printstring()`
- `runtime_printint()`
- Memory allocation functions
- Type assertion functions

**Workaround**: Implement these functions in C:

```c
// Example runtime function implementations
void runtime_printstring(const char* str, int len) {
    fwrite(str, 1, len, stdout);
}

void runtime_printint(int value) {
    printf("%d\n", value);
}
```

## Best Practices for Successful Transpilation

### ✓ DO

1. **Use simple, straightforward Go code**
2. **Avoid goroutines and channels**
3. **Use fixed-size arrays instead of slices**
4. **Minimize pointer complexity**
5. **Avoid reflection**
6. **Test with TinyGo first** to ensure compatibility
7. **Keep programs small** and single-file
8. **Use standard types** (int, float, bool)

### ✗ DON'T

1. **Don't use goroutines**
2. **Don't use channels**
3. **Don't use reflection**
4. **Don't use complex interfaces**
5. **Don't rely on garbage collection**
6. **Don't use defer extensively**
7. **Don't use panic/recover**
8. **Don't use complex pointer operations**

## Testing Strategy for Transpiled Code

1. **Compile with warnings enabled**: `gcc -Wall -Wextra`
2. **Test with sanitizers**: `gcc -fsanitize=address,undefined`
3. **Run under valgrind** for memory issues
4. **Compare outputs** with original Go program
5. **Benchmark performance** if needed

## Getting Help

If you encounter issues:
1. Check if limitation is documented here
2. Test with simpler Go code
3. Inspect generated LLVM IR with `-keep-llvm`
4. File an issue with minimal reproduction case
5. Consider manual C code adjustments

## Version Compatibility

**Tested Configurations**:
- TinyGo 0.33.0
- LLVM 18.1.2
- Go 1.19-1.23
- GCC 11+, Clang 14+

## Future Roadmap

Planned improvements:
1. Better control flow detection and generation
2. Runtime library for common Go functions
3. Slice/map support with bundled implementations
4. Interface vtable generation
5. Improved error messages with suggestions
