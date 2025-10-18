# Multiple Return Values Support

## Overview

The go2c transpiler now supports Go's multiple return values pattern by automatically converting them to C struct types. This is a fundamental Go idiom, especially for error handling, and is now fully supported in the transpilation process.

## Feature Description

### What It Does

When a Go function returns multiple values (represented as aggregate types in LLVM IR), the transpiler:

1. **Detects** aggregate return types like `{i32, i32}` in function signatures
2. **Generates** a C struct definition to hold all return values
3. **Converts** `insertvalue` instructions to struct field assignments
4. **Converts** `extractvalue` instructions to struct field accesses
5. **Handles** function calls that return aggregate types

### LLVM IR Pattern

In LLVM IR, multiple return values are represented as aggregate (struct) types:

```llvm
define {i32, i32} @foo(i32 %x) {
entry:
  %result = insertvalue {i32, i32} undef, i32 %x, 0
  %result2 = insertvalue {i32, i32} %result, i32 42, 1
  ret {i32, i32} %result2
}
```

### Generated C Code

The transpiler generates clean, type-safe C code:

```c
// Automatically generated struct type
typedef struct {
    int field0;
    int field1;
} multi_return_1_t;

// Function with struct return type
multi_return_1_t foo(int x) {
    multi_return_1_t result; 
    result.field0 = x;
    multi_return_1_t result2 = result; 
    result2.field1 = 42;
    return result2;
}

// Usage
multi_return_1_t ret = foo(10);
int val1 = ret.field0;
int val2 = ret.field1;
```

## Use Cases

### 1. Error Handling Pattern

**Go Code:**
```go
func divide(a, b int) (int, bool) {
    if b == 0 {
        return 0, false  // error
    }
    return a / b, true  // success
}

result, ok := divide(10, 2)
if !ok {
    // handle error
}
```

**Generated C Code:**
```c
typedef struct {
    int field0;      // result
    bool field1;     // success flag
} divide_return_t;

divide_return_t divide(int a, int b) {
    divide_return_t ret;
    if (b == 0) {
        ret.field0 = 0;
        ret.field1 = false;
        return ret;
    }
    ret.field0 = a / b;
    ret.field1 = true;
    return ret;
}

// Usage
divide_return_t ret = divide(10, 2);
if (!ret.field1) {
    // handle error
}
```

### 2. Multiple Values Pattern

**Go Code:**
```go
func minMax(a, b int) (int, int) {
    if a < b {
        return a, b
    }
    return b, a
}

min, max := minMax(5, 3)
```

**Generated C Code:**
```c
typedef struct {
    int field0;      // min
    int field1;      // max
} minMax_return_t;

minMax_return_t minMax(int a, int b) {
    minMax_return_t ret;
    if (a < b) {
        ret.field0 = a;
        ret.field1 = b;
    } else {
        ret.field0 = b;
        ret.field1 = a;
    }
    return ret;
}

// Usage
minMax_return_t result = minMax(5, 3);
int min = result.field0;
int max = result.field1;
```

### 3. Coordinate/Position Pattern

**Go Code:**
```go
func getPosition() (int, int) {
    return 10, 20
}

x, y := getPosition()
```

**Generated C Code:**
```c
typedef struct {
    int field0;      // x
    int field1;      // y
} getPosition_return_t;

getPosition_return_t getPosition(void) {
    getPosition_return_t ret;
    ret.field0 = 10;
    ret.field1 = 20;
    return ret;
}

// Usage
getPosition_return_t pos = getPosition();
int x = pos.field0;
int y = pos.field1;
```

## Technical Details

### Type Registration

The transpiler maintains a registry of aggregate types to avoid duplicate struct definitions:

- Each unique aggregate type signature gets a unique struct name
- Struct names are generated as `multi_return_N_t` where N is incremented
- The same aggregate type signature reuses the same struct definition

### LLVM Instructions Supported

1. **insertvalue**: Assigns a value to a specific field in an aggregate
   ```llvm
   %result = insertvalue {i32, i32} %aggregate, i32 %value, 0
   ```
   Converts to:
   ```c
   multi_return_1_t result = aggregate;
   result.field0 = value;
   ```

2. **extractvalue**: Extracts a value from a specific field in an aggregate
   ```llvm
   %value = extractvalue {i32, i32} %aggregate, 0
   ```
   Converts to:
   ```c
   int value = aggregate.field0;
   ```

### Return Statement Handling

Return statements with aggregate types are properly converted:

```llvm
ret {i32, i32} %result
```

Converts to:
```c
return result;
```

### Function Call Handling

Function calls that return aggregate types are properly typed:

```llvm
%ret = call {i32, i32} @foo(i32 10)
```

Converts to:
```c
multi_return_1_t ret = foo(10);
```

## Limitations

1. **Generic Field Names**: Fields are named `field0`, `field1`, etc., rather than semantic names like `result`, `error`
2. **No Named Returns**: Go's named return values are not preserved in C
3. **Error Type Simplification**: Complex Go error types become simple integers or bools
4. **Performance**: Large structs may be less efficient than pass-by-reference patterns

## Best Practices

### 1. Use for Simple Value Returns

Multiple return values work best for small, simple types:

```c
// Good: Two integers
typedef struct {
    int field0;
    int field1;
} point_t;

// Less ideal: Large structs
typedef struct {
    char field0[1000];
    int field1[100];
} large_return_t;  // Consider pass-by-reference instead
```

### 2. Document Field Meanings

Add comments to clarify what each field represents:

```c
typedef struct {
    int field0;      // x coordinate
    int field1;      // y coordinate
} position_return_t;
```

### 3. Use Consistent Error Patterns

Adopt a consistent pattern for error handling:

```c
typedef struct {
    int field0;       // result value
    bool field1;      // success flag (true = success, false = error)
} operation_return_t;
```

## Testing

The feature includes comprehensive test coverage:

1. **Unit Tests**: Type mapping for aggregate types
2. **Integration Tests**: End-to-end transpilation with multiple return values
3. **Validation**: Generated C code compiles and executes correctly

See `testdata/input/multiret.ll` for test examples.

## Implementation Details

### Files Modified

1. **internal/llvm/parser.go**: Updated to parse aggregate return types in function signatures
2. **internal/codegen/types.go**: Added aggregate type registration and struct generation
3. **internal/codegen/emitter.go**: Added support for:
   - Struct definitions
   - insertvalue instruction conversion
   - extractvalue instruction conversion
   - Aggregate return type handling

### Key Data Structures

```go
// AggregateType represents a multi-return value struct
type AggregateType struct {
    Name       string   // C struct name (e.g., "multi_return_1_t")
    LLVMType   string   // Original LLVM type (e.g., "{i32, i32}")
    FieldTypes []string // Field types (e.g., ["i32", "i32"])
}
```

## Future Enhancements

Potential improvements for future versions:

1. **Named Fields**: Extract semantic names from Go source when available
2. **Error Type Support**: Better representation of Go error types
3. **Optimization**: Detect when pass-by-reference would be more efficient
4. **Tuple Syntax**: Generate helper macros for tuple-like unpacking
5. **Documentation Generation**: Auto-generate field documentation from Go comments

## Examples

See the `examples/` directory for complete working examples of multiple return values in action.

## Conclusion

Multiple return values support brings the go2c transpiler closer to supporting idiomatic Go code. This feature enables a wide range of Go programs to be successfully transpiled to C while maintaining type safety and code clarity.
