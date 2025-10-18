package llvm

import (
	"testing"
)

func TestControlFlowAnalyzer_IfElsePattern(t *testing.T) {
	// Test function with if-else pattern
	fn := &Function{
		Name:       "max",
		ReturnType: "i32",
		Parameters: []Parameter{
			{Type: "i32", Name: "%a"},
			{Type: "i32", Name: "%b"},
		},
		Body: []string{
			"entry:",
			"  %cmp = icmp sgt i32 %a, %b",
			"  br i1 %cmp, label %if.then, label %if.else",
			"if.then:",
			"  ret i32 %a",
			"if.else:",
			"  ret i32 %b",
		},
	}

	cfa := NewControlFlowAnalyzer(fn)
	patterns, err := cfa.AnalyzeControlFlow()

	if err != nil {
		t.Fatalf("AnalyzeControlFlow failed: %v", err)
	}

	// Debug: print blocks
	blocks := cfa.GetBasicBlocks()
	t.Logf("Found %d blocks", len(blocks))
	for name, block := range blocks {
		t.Logf("Block %s: %d instructions, terminator: %v", name, len(block.Instructions), block.Terminator)
		if block.Terminator != nil {
			t.Logf("  Terminator type: %s, true: %s, false: %s",
				block.Terminator.Type, block.Terminator.TrueLabel, block.Terminator.FalseLabel)
		}
	}

	t.Logf("Found %d patterns", len(patterns))
	for i, p := range patterns {
		t.Logf("Pattern %d: type=%s, start=%s, end=%s", i, p.Type, p.StartLabel, p.EndLabel)
	}

	if len(patterns) != 1 {
		t.Errorf("Expected 1 pattern, got %d", len(patterns))
	}

	if len(patterns) > 0 && patterns[0].Type != "if-else" {
		t.Errorf("Expected if-else pattern, got %s", patterns[0].Type)
	}
}

func TestControlFlowAnalyzer_IfThenPattern(t *testing.T) {
	// Test function with if-then pattern (no else)
	fn := &Function{
		Name:       "print_positive",
		ReturnType: "void",
		Parameters: []Parameter{
			{Type: "i32", Name: "%x"},
		},
		Body: []string{
			"entry:",
			"  %cmp = icmp sgt i32 %x, 0",
			"  br i1 %cmp, label %if.then, label %if.end",
			"if.then:",
			"  call void @runtime.printint(i32 %x)",
			"  br label %if.end",
			"if.end:",
			"  ret void",
		},
	}

	cfa := NewControlFlowAnalyzer(fn)
	patterns, err := cfa.AnalyzeControlFlow()

	if err != nil {
		t.Fatalf("AnalyzeControlFlow failed: %v", err)
	}

	if len(patterns) != 1 {
		t.Errorf("Expected 1 pattern, got %d", len(patterns))
	}

	if len(patterns) > 0 && patterns[0].Type != "if-then" {
		t.Errorf("Expected if-then pattern, got %s", patterns[0].Type)
	}
}

func TestControlFlowAnalyzer_WhilePattern(t *testing.T) {
	// Test function with while loop pattern
	fn := &Function{
		Name:       "sum_to_n",
		ReturnType: "i32",
		Parameters: []Parameter{
			{Type: "i32", Name: "%n"},
		},
		Body: []string{
			"entry:",
			"  %sum = alloca i32",
			"  %i = alloca i32",
			"  store i32 0, i32* %sum",
			"  store i32 0, i32* %i",
			"  br label %while.cond",
			"while.cond:",
			"  %i.val = load i32, i32* %i",
			"  %cmp = icmp slt i32 %i.val, %n",
			"  br i1 %cmp, label %while.body, label %while.end",
			"while.body:",
			"  %sum.val = load i32, i32* %sum",
			"  %i.val2 = load i32, i32* %i",
			"  %add = add i32 %sum.val, %i.val2",
			"  store i32 %add, i32* %sum",
			"  %inc = add i32 %i.val2, 1",
			"  store i32 %inc, i32* %i",
			"  br label %while.cond",
			"while.end:",
			"  %result = load i32, i32* %sum",
			"  ret i32 %result",
		},
	}

	cfa := NewControlFlowAnalyzer(fn)
	patterns, err := cfa.AnalyzeControlFlow()

	if err != nil {
		t.Fatalf("AnalyzeControlFlow failed: %v", err)
	}

	if len(patterns) != 1 {
		t.Errorf("Expected 1 pattern, got %d", len(patterns))
	}

	if len(patterns) > 0 && patterns[0].Type != "while" {
		t.Errorf("Expected while pattern, got %s", patterns[0].Type)
	}
}

func TestControlFlowAnalyzer_SwitchPattern(t *testing.T) {
	// Test function with switch pattern
	fn := &Function{
		Name:       "testSwitch",
		ReturnType: "i32",
		Parameters: []Parameter{
			{Type: "i32", Name: "%x"},
		},
		Body: []string{
			"entry:",
			"  switch i32 %x, label %default [",
			"    i32 0, label %case0",
			"    i32 1, label %case1",
			"    i32 2, label %case2",
			"  ]",
			"case0:",
			"  ret i32 10",
			"case1:",
			"  ret i32 20",
			"case2:",
			"  ret i32 30",
			"default:",
			"  ret i32 -1",
		},
	}

	cfa := NewControlFlowAnalyzer(fn)
	patterns, err := cfa.AnalyzeControlFlow()

	if err != nil {
		t.Fatalf("AnalyzeControlFlow failed: %v", err)
	}

	// Debug: print blocks
	blocks := cfa.GetBasicBlocks()
	t.Logf("Found %d blocks", len(blocks))
	for name, block := range blocks {
		t.Logf("Block %s: %d instructions, terminator: %v", name, len(block.Instructions), block.Terminator)
		if block.Terminator != nil && block.Terminator.Type == "switch" {
			t.Logf("  Switch value: %s, default: %s, cases: %d",
				block.Terminator.SwitchValue, block.Terminator.DefaultLabel, len(block.Terminator.SwitchCases))
			for i, c := range block.Terminator.SwitchCases {
				t.Logf("    Case %d: value=%s, label=%s", i, c.Value, c.Label)
			}
		}
	}

	t.Logf("Found %d patterns", len(patterns))
	for i, p := range patterns {
		t.Logf("Pattern %d: type=%s, start=%s", i, p.Type, p.StartLabel)
	}

	if len(patterns) != 1 {
		t.Errorf("Expected 1 pattern, got %d", len(patterns))
	}

	if len(patterns) > 0 {
		if patterns[0].Type != "switch" {
			t.Errorf("Expected switch pattern, got %s", patterns[0].Type)
		}
		if len(patterns[0].SwitchCases) != 3 {
			t.Errorf("Expected 3 switch cases, got %d", len(patterns[0].SwitchCases))
		}
		if patterns[0].DefaultLabel != "default" {
			t.Errorf("Expected default label 'default', got '%s'", patterns[0].DefaultLabel)
		}
	}
}

func TestControlFlowAnalyzer_BasicBlocks(t *testing.T) {
	fn := &Function{
		Name:       "simple",
		ReturnType: "i32",
		Body: []string{
			"entry:",
			"  %x = add i32 1, 2",
			"  br label %next",
			"next:",
			"  ret i32 %x",
		},
	}

	cfa := NewControlFlowAnalyzer(fn)
	cfa.buildBasicBlocks()

	blocks := cfa.GetBasicBlocks()

	if len(blocks) != 2 {
		t.Errorf("Expected 2 basic blocks, got %d", len(blocks))
	}

	if _, exists := blocks["entry"]; !exists {
		t.Errorf("Entry block not found")
	}

	if _, exists := blocks["next"]; !exists {
		t.Errorf("Next block not found")
	}
}
