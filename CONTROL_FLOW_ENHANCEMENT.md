# Control Flow Enhancement - Before and After

This document demonstrates the improvement in generated C code quality after implementing control flow pattern detection.

## Problem Statement

Previously, the go2c transpiler converted all LLVM IR control flow to labels and goto statements, resulting in C code that was difficult to read and maintain.

## Solution

Implemented a control flow analyzer that:
1. Builds a control flow graph from LLVM IR basic blocks
2. Detects common patterns (if-else, if-then, while loops)
3. Generates idiomatic C code for detected patterns

## Example 1: If-Else Pattern

### LLVM IR Input
```llvm
define i32 @max(i32 %a, i32 %b) {
entry:
  %cmp = icmp sgt i32 %a, %b
  br i1 %cmp, label %if.then, label %if.else

if.then:
  ret i32 %a

if.else:
  ret i32 %b
}
```

### Before (Label-based)
```c
int max(int a, int b) {
    entry:
    bool cmp = a > b;
    goto if_then;  // or if_else
    
    if_then:
    return a;
    
    if_else:
    return b;
}
```

### After (Structured)
```c
int max(int a, int b) {
    bool cmp = a > b;
    if (cmp) {
        return a;
    } else {
        return b;
    }
}
```

## Example 2: While Loop Pattern

### LLVM IR Input
```llvm
define i32 @sum_to_n(i32 %n) {
entry:
  %sum = alloca i32
  %i = alloca i32
  store i32 0, i32* %sum
  store i32 0, i32* %i
  br label %while.cond

while.cond:
  %i.val = load i32, i32* %i
  %cmp = icmp slt i32 %i.val, %n
  br i1 %cmp, label %while.body, label %while.end

while.body:
  %sum.val = load i32, i32* %sum
  %i.val2 = load i32, i32* %i
  %add = add i32 %sum.val, %i.val2
  store i32 %add, i32* %sum
  %inc = add i32 %i.val2, 1
  store i32 %inc, i32* %i
  br label %while.cond

while.end:
  %result = load i32, i32* %sum
  ret i32 %result
}
```

### Before (Label-based)
```c
int sum_to_n(int n) {
    int sum;
    int i;
    *sum = 0;
    *i = 0;
    goto while_cond;
    
    while_cond:
    int i_val = *i;
    bool cmp = i_val < n;
    if (cmp) goto while_body; else goto while_end;
    
    while_body:
    int sum_val = *sum;
    int i_val2 = *i;
    int add = sum_val + i_val2;
    *sum = add;
    int inc = i_val2 + 1;
    *i = inc;
    goto while_cond;
    
    while_end:
    int result = *sum;
    return result;
}
```

### After (Structured)
```c
int sum_to_n(int n) {
    int sum;
    int i;
    *sum = 0;
    *i = 0;
    int i_val = *i;
    bool cmp = i_val < n;
    while (cmp) {
        int sum_val = *sum;
        int i_val2 = *i;
        int add = sum_val + i_val2;
        *sum = add;
        int inc = i_val2 + 1;
        *i = inc;
        int i_val = *i;
        bool cmp = i_val < n;
    }
    int result = *sum;
    return result;
}
```

## Example 3: If-Then Pattern (No Else)

### LLVM IR Input
```llvm
define void @print_positive(i32 %x) {
entry:
  %cmp = icmp sgt i32 %x, 0
  br i1 %cmp, label %if.then, label %if.end

if.then:
  call void @runtime.printint(i32 %x)
  br label %if.end

if.end:
  ret void
}
```

### Before (Label-based)
```c
void print_positive(int x) {
    entry:
    bool cmp = x > 0;
    if (cmp) goto if_then; else goto if_end;
    
    if_then:
    runtime_printint(x);
    goto if_end;
    
    if_end:
    return;
}
```

### After (Structured)
```c
void print_positive(int x) {
    bool cmp = x > 0;
    if (cmp) {
        runtime_printint(x);
    }
    return;
}
```

## Comparison Instructions Support

The enhancement also added support for LLVM comparison instructions:

| LLVM Instruction | C Operator | Description |
|------------------|------------|-------------|
| `icmp eq`        | `==`       | Equal |
| `icmp ne`        | `!=`       | Not equal |
| `icmp sgt/ugt`   | `>`        | Greater than (signed/unsigned) |
| `icmp sge/uge`   | `>=`       | Greater or equal |
| `icmp slt/ult`   | `<`        | Less than |
| `icmp sle/ule`   | `<=`       | Less or equal |

## Benefits

1. **Readability**: Generated C code is much easier to read and understand
2. **Maintainability**: Structured code is easier to modify and debug
3. **Compatibility**: Idiomatic C code works better with C analysis tools
4. **Performance**: Modern C compilers can optimize structured code better than goto-heavy code

## Implementation Details

### Control Flow Analyzer (`internal/llvm/controlflow.go`)

- **BasicBlock**: Represents a basic block with instructions and terminator
- **Terminator**: Represents branch/return instructions
- **ControlFlowPattern**: Detected pattern (if-else, if-then, while)
- **ControlFlowAnalyzer**: Analyzes function body to detect patterns

### Enhanced Code Generator (`internal/codegen/emitter.go`)

- `generateFunctionWithControlFlow()`: Entry point for enhanced generation
- `generateIfElse()`: Generates if-else statements
- `generateIfThen()`: Generates if statements without else
- `generateWhile()`: Generates while loops
- `convertInstructionToC()`: Enhanced with icmp support

## Test Coverage

New tests added:
- `TestControlFlowAnalyzer_IfElsePattern`
- `TestControlFlowAnalyzer_IfThenPattern`
- `TestControlFlowAnalyzer_WhilePattern`
- `TestControlFlowAnalyzer_BasicBlocks`
- Integration test for control flow patterns

All existing tests continue to pass, ensuring backward compatibility.

## Future Enhancements

Potential improvements for future work:
1. **For loops**: Detect for-loop patterns
2. **Do-while loops**: Detect post-condition loops
3. **Switch statements**: Convert LLVM switch to C switch
4. **Nested patterns**: Better handling of nested control flow
5. **Dead code elimination**: Remove unreachable labels
6. **Variable scoping**: Better scoping for loop variables
