# Limitation Analysis Report

## Selected Limitation
**Title**: Switch Statements with Fallthrough (Control Flow Complexity)  
**Category**: Control Flow / Code Generation Quality  
**Severity**: Minor to Moderate  
**Selection Rationale**: Switch statements are a common control flow pattern in programming. The current implementation converts all switch statements to label-based goto logic, making generated C code less readable and idiomatic. This limitation has a high impact-to-effort ratio because: (1) it improves code readability significantly, (2) it leverages existing control flow infrastructure, (3) it has clear success criteria, and (4) it requires no external dependencies.

## Root Cause Analysis

The limitation exists due to the current implementation strategy in the code generator:

1. **Current Approach**: The code generator (`internal/codegen/emitter.go`) focuses on converting LLVM IR instructions one-by-one without recognizing higher-level patterns like switch statements.

2. **LLVM Representation**: LLVM IR uses a `switch` instruction that maps perfectly to C's `switch` statement:
   ```llvm
   switch i32 %value, label %default [
     i32 0, label %case0
     i32 1, label %case1
     i32 2, label %case2
   ]
   ```

3. **Missing Pattern Recognition**: While the codebase recently added control flow pattern detection for if-else and while loops (`internal/llvm/controlflow.go`), it doesn't yet recognize switch patterns.

4. **Label-Based Fallback**: Without switch recognition, the code generator falls back to generating labels and goto statements, which is functional but not idiomatic C.

## Solution Design

### Architecture Overview

The solution extends the existing control flow analysis infrastructure to:
1. Detect LLVM `switch` instructions in the IR
2. Parse switch cases and default labels
3. Generate idiomatic C `switch` statements
4. Handle fallthrough behavior correctly

### Components

1. **Switch Instruction Parser** (`internal/llvm/controlflow.go`)
   - Purpose: Parse LLVM switch instructions
   - Function: Extract switch value, cases, and labels
   - Integration: Extends existing Terminator struct

2. **Switch Pattern Detector** (`internal/llvm/controlflow.go`)
   - Purpose: Identify switch control flow patterns
   - Function: Recognize switch blocks and case blocks
   - Integration: Adds to existing pattern detection

3. **C Switch Generator** (`internal/codegen/emitter.go`)
   - Purpose: Generate C switch statements
   - Function: Convert LLVM switch patterns to C code
   - Integration: Extends generateFunctionWithControlFlow

### Third-Party Dependencies

| Library | Version | Purpose | Integration |
|---------|---------|---------|-------------|
| None    | N/A     | N/A     | Uses only standard Go libraries |

### Data Flow

```
LLVM IR switch instruction
    ↓
[Parser] Extract switch value and cases
    ↓
[Analyzer] Build control flow pattern
    ↓
[Code Generator] Generate C switch statement
    ↓
Idiomatic C code with switch/case/break
```

## Implementation Plan

### Step 1: Extend Terminator Struct for Switch
**Objective**: Add switch instruction support to the Terminator data structure
**Actions**:
- [ ] Add switch-specific fields to Terminator struct (SwitchValue, SwitchCases)
- [ ] Update parseTerminator to handle switch instructions
- [ ] Add unit tests for switch parsing
**Verification**: Unit tests pass for switch terminator parsing
**Estimated Time**: 30 minutes

### Step 2: Implement Switch Pattern Detection
**Objective**: Add switch pattern detection to control flow analyzer
**Actions**:
- [ ] Create detectSwitchPattern function
- [ ] Add "switch" pattern type to ControlFlowPattern
- [ ] Update detectPatterns to include switch detection
- [ ] Add unit tests for switch pattern detection
**Verification**: Control flow analyzer correctly identifies switch patterns
**Estimated Time**: 45 minutes

### Step 3: Implement C Switch Code Generation
**Objective**: Generate idiomatic C switch statements
**Actions**:
- [ ] Create generateSwitch function in emitter.go
- [ ] Handle switch value and case labels
- [ ] Generate proper break statements (unless fallthrough is needed)
- [ ] Handle default case
- [ ] Update generateBlockWithPatterns to handle switch patterns
**Verification**: Generated C code contains proper switch/case syntax
**Estimated Time**: 1 hour

### Step 4: Create Test Cases
**Objective**: Comprehensive testing of switch statement support
**Actions**:
- [ ] Create LLVM IR test file with switch statement
- [ ] Add integration test for switch transpilation
- [ ] Test edge cases (no default, single case, multiple cases)
- [ ] Test fallthrough behavior
**Verification**: All tests pass, coverage increases
**Estimated Time**: 45 minutes

### Step 5: Documentation and Cleanup
**Objective**: Update documentation and finalize implementation
**Actions**:
- [ ] Update LIMITATIONS.md to remove/update switch statement limitation
- [ ] Update README.md with switch support in features list
- [ ] Add code comments explaining switch generation
- [ ] Run full test suite and verify no regressions
**Verification**: Documentation is accurate, all tests pass
**Estimated Time**: 30 minutes

## Success Validation

### Before Solution
- ❌ Switch statements converted to labels and goto
- ❌ Generated C code is not idiomatic
- ❌ Harder to read and maintain generated code
- ❌ Listed as limitation in documentation

### After Solution
- ✅ Switch statements generate C switch/case
- ✅ Idiomatic, readable C code
- ✅ Matches hand-written C switch patterns
- ✅ Documentation updated to reflect support

### Metrics
- Code readability: Significant improvement for switch-heavy code
- Test coverage: Additional test cases for switch patterns
- Feature completeness: One more control flow pattern supported

## Next Steps
1. Implement the solution following the plan above
2. Consider adding support for more complex switch patterns (ranges, etc.)
3. Consider detecting and warning about C-incompatible fallthrough patterns
