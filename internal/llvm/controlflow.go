package llvm

import (
	"regexp"
	"strings"
)

// BasicBlock represents a basic block in a control flow graph
type BasicBlock struct {
	Label        string
	Instructions []string
	Terminator   *Terminator
}

// SwitchCase represents a case in a switch statement
type SwitchCase struct {
	Value string // The constant value for this case
	Label string // The target label for this case
}

// Terminator represents the terminating instruction of a basic block
type Terminator struct {
	Type         string       // "br_cond", "br_uncon", "ret", "switch", etc.
	Condition    string       // For conditional branches
	TrueLabel    string       // For conditional branches
	FalseLabel   string       // For conditional branches
	UnconLabel   string       // For unconditional branches
	SwitchValue  string       // For switch: the value being switched on
	SwitchCases  []SwitchCase // For switch: the case values and labels
	DefaultLabel string       // For switch: the default label
}

// ControlFlowPattern represents a detected control flow pattern
type ControlFlowPattern struct {
	Type         string       // "if-else", "if-then", "while", "for", "switch"
	StartLabel   string       // The starting block label
	EndLabel     string       // The ending/merge block label (if any)
	CondBlock    string       // For if/while: condition block
	ThenBlock    string       // For if: then block
	ElseBlock    string       // For if: else block
	BodyBlock    string       // For loops: body block
	SwitchValue  string       // For switch: the value being switched on
	SwitchCases  []SwitchCase // For switch: the case values and labels
	DefaultLabel string       // For switch: the default label
	Blocks       []string     // All blocks involved in this pattern
}

// ControlFlowAnalyzer analyzes control flow in LLVM IR functions
type ControlFlowAnalyzer struct {
	function *Function
	blocks   map[string]*BasicBlock
}

// NewControlFlowAnalyzer creates a new control flow analyzer
func NewControlFlowAnalyzer(fn *Function) *ControlFlowAnalyzer {
	return &ControlFlowAnalyzer{
		function: fn,
		blocks:   make(map[string]*BasicBlock),
	}
}

// AnalyzeControlFlow analyzes the control flow of a function
func (cfa *ControlFlowAnalyzer) AnalyzeControlFlow() ([]*ControlFlowPattern, error) {
	// Build basic blocks from function body
	cfa.buildBasicBlocks()
	
	// Detect control flow patterns
	patterns := cfa.detectPatterns()
	
	return patterns, nil
}

// buildBasicBlocks constructs basic blocks from function body
func (cfa *ControlFlowAnalyzer) buildBasicBlocks() {
	var currentBlock *BasicBlock
	var currentLabel string
	var switchAccumulator string
	var inSwitch bool
	
	for _, line := range cfa.function.Body {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		
		// Handle multi-line switch statements
		if inSwitch {
			switchAccumulator += " " + line
			// Check if this is the end of the switch statement
			if strings.HasSuffix(line, "]") {
				inSwitch = false
				if currentBlock != nil {
					currentBlock.Instructions = append(currentBlock.Instructions, switchAccumulator)
					currentBlock.Terminator = cfa.parseTerminator(switchAccumulator)
					cfa.blocks[currentLabel] = currentBlock
					currentBlock = nil
				}
				switchAccumulator = ""
			}
			continue
		}
		
		// Check if this is a label (basic block start)
		if strings.HasSuffix(line, ":") {
			// Save previous block
			if currentBlock != nil {
				cfa.blocks[currentLabel] = currentBlock
			}
			
			// Start new block
			currentLabel = strings.TrimSuffix(line, ":")
			currentBlock = &BasicBlock{
				Label:        currentLabel,
				Instructions: []string{},
			}
			continue
		}
		
		// Check if this is a switch terminator (which may span multiple lines)
		if strings.HasPrefix(line, "switch") {
			inSwitch = true
			switchAccumulator = line
			// Check if it's a single-line switch (unlikely but possible)
			if strings.HasSuffix(line, "]") {
				inSwitch = false
				if currentBlock != nil {
					currentBlock.Instructions = append(currentBlock.Instructions, switchAccumulator)
					currentBlock.Terminator = cfa.parseTerminator(switchAccumulator)
					cfa.blocks[currentLabel] = currentBlock
					currentBlock = nil
				}
				switchAccumulator = ""
			}
			continue
		}
		
		// Check if this is a terminator instruction
		if cfa.isTerminator(line) {
			if currentBlock != nil {
				// Store the full terminator instruction in the block
				currentBlock.Instructions = append(currentBlock.Instructions, line)
				currentBlock.Terminator = cfa.parseTerminator(line)
				cfa.blocks[currentLabel] = currentBlock
				currentBlock = nil
			}
			continue
		}
		
		// Regular instruction
		if currentBlock != nil {
			currentBlock.Instructions = append(currentBlock.Instructions, line)
		}
	}
	
	// Save last block if any
	if currentBlock != nil {
		cfa.blocks[currentLabel] = currentBlock
	}
}

// isTerminator checks if an instruction is a terminator
func (cfa *ControlFlowAnalyzer) isTerminator(instruction string) bool {
	return strings.HasPrefix(instruction, "ret") ||
		strings.HasPrefix(instruction, "br") ||
		strings.HasPrefix(instruction, "switch") ||
		strings.HasPrefix(instruction, "unreachable")
}

// parseTerminator parses a terminator instruction
func (cfa *ControlFlowAnalyzer) parseTerminator(instruction string) *Terminator {
	term := &Terminator{}
	
	if strings.HasPrefix(instruction, "ret") {
		term.Type = "ret"
		return term
	}
	
	if strings.HasPrefix(instruction, "br") {
		// Check if conditional or unconditional
		if strings.Contains(instruction, ",") {
			// Conditional: br i1 %cmp, label %true, label %false
			term.Type = "br_cond"
			
			// Extract condition
			re := regexp.MustCompile(`br\s+i1\s+(%\S+)`)
			if matches := re.FindStringSubmatch(instruction); len(matches) > 1 {
				term.Condition = strings.TrimSuffix(matches[1], ",")
			}
			
			// Extract labels
			labelRe := regexp.MustCompile(`label\s+(%\S+)`)
			labels := labelRe.FindAllStringSubmatch(instruction, -1)
			if len(labels) >= 2 {
				term.TrueLabel = strings.TrimSuffix(strings.TrimPrefix(labels[0][1], "%"), ",")
				term.FalseLabel = strings.TrimSuffix(strings.TrimPrefix(labels[1][1], "%"), ",")
			}
		} else {
			// Unconditional: br label %target
			term.Type = "br_uncon"
			
			re := regexp.MustCompile(`label\s+(%\S+)`)
			if matches := re.FindStringSubmatch(instruction); len(matches) > 1 {
				term.UnconLabel = strings.TrimPrefix(matches[1], "%")
			}
		}
		return term
	}
	
	if strings.HasPrefix(instruction, "switch") {
		// Switch: switch i32 %x, label %default [ i32 0, label %case0 ... ]
		term.Type = "switch"
		
		// Extract switch value and type
		// Pattern: switch <type> <value>, label <default> [...]
		re := regexp.MustCompile(`switch\s+(\S+)\s+(%?\S+),\s+label\s+(%\S+)`)
		if matches := re.FindStringSubmatch(instruction); len(matches) > 3 {
			term.SwitchValue = matches[2]
			term.DefaultLabel = strings.TrimPrefix(matches[3], "%")
		}
		
		// Extract case values and labels
		// Pattern: i32 <value>, label %<label>
		caseRe := regexp.MustCompile(`i\d+\s+(-?\d+),\s+label\s+(%\S+)`)
		caseMatches := caseRe.FindAllStringSubmatch(instruction, -1)
		term.SwitchCases = make([]SwitchCase, 0, len(caseMatches))
		for _, match := range caseMatches {
			if len(match) > 2 {
				term.SwitchCases = append(term.SwitchCases, SwitchCase{
					Value: match[1],
					Label: strings.TrimPrefix(match[2], "%"),
				})
			}
		}
		
		return term
	}
	
	term.Type = "unknown"
	return term
}

// detectPatterns detects common control flow patterns
func (cfa *ControlFlowAnalyzer) detectPatterns() []*ControlFlowPattern {
	patterns := []*ControlFlowPattern{}
	
	// Track which blocks are already part of a pattern
	usedBlocks := make(map[string]bool)
	
	// Detect if-else patterns
	for label, block := range cfa.blocks {
		if usedBlocks[label] {
			continue
		}
		
		if block.Terminator != nil && block.Terminator.Type == "br_cond" {
			pattern := cfa.detectIfElsePattern(label, block)
			if pattern != nil {
				patterns = append(patterns, pattern)
				for _, b := range pattern.Blocks {
					usedBlocks[b] = true
				}
			}
		}
	}
	
	// Detect while loop patterns
	for label, block := range cfa.blocks {
		if usedBlocks[label] {
			continue
		}
		
		pattern := cfa.detectWhilePattern(label, block)
		if pattern != nil {
			patterns = append(patterns, pattern)
			for _, b := range pattern.Blocks {
				usedBlocks[b] = true
			}
		}
	}
	
	// Detect switch patterns
	for label, block := range cfa.blocks {
		if usedBlocks[label] {
			continue
		}
		
		if block.Terminator != nil && block.Terminator.Type == "switch" {
			pattern := cfa.detectSwitchPattern(label, block)
			if pattern != nil {
				patterns = append(patterns, pattern)
				for _, b := range pattern.Blocks {
					usedBlocks[b] = true
				}
			}
		}
	}
	
	return patterns
}

// detectIfElsePattern detects if-then-else or if-then patterns
func (cfa *ControlFlowAnalyzer) detectIfElsePattern(label string, block *BasicBlock) *ControlFlowPattern {
	if block.Terminator == nil || block.Terminator.Type != "br_cond" {
		return nil
	}
	
	trueLabel := block.Terminator.TrueLabel
	falseLabel := block.Terminator.FalseLabel
	
	trueBlock, trueExists := cfa.blocks[trueLabel]
	falseBlock, falseExists := cfa.blocks[falseLabel]
	
	if !trueExists || !falseExists {
		return nil
	}
	
	// Check if both branches converge to a common block (if-else)
	// or if false branch is the merge point (if-then)
	
	// Simple if-then pattern: true branch jumps to false label
	if trueBlock.Terminator != nil && 
	   trueBlock.Terminator.Type == "br_uncon" &&
	   trueBlock.Terminator.UnconLabel == falseLabel {
		return &ControlFlowPattern{
			Type:       "if-then",
			StartLabel: label,
			EndLabel:   falseLabel,
			CondBlock:  label,
			ThenBlock:  trueLabel,
			Blocks:     []string{label, trueLabel, falseLabel},
		}
	}
	
	// Check for if-else pattern: both branches jump to same merge point
	if trueBlock.Terminator != nil && falseBlock.Terminator != nil &&
	   trueBlock.Terminator.Type == "br_uncon" &&
	   falseBlock.Terminator.Type == "br_uncon" &&
	   trueBlock.Terminator.UnconLabel == falseBlock.Terminator.UnconLabel {
		
		mergeLabel := trueBlock.Terminator.UnconLabel
		return &ControlFlowPattern{
			Type:       "if-else",
			StartLabel: label,
			EndLabel:   mergeLabel,
			CondBlock:  label,
			ThenBlock:  trueLabel,
			ElseBlock:  falseLabel,
			Blocks:     []string{label, trueLabel, falseLabel, mergeLabel},
		}
	}
	
	// Check for if-else with immediate return in both branches
	if trueBlock.Terminator != nil && falseBlock.Terminator != nil &&
	   trueBlock.Terminator.Type == "ret" &&
	   falseBlock.Terminator.Type == "ret" {
		
		return &ControlFlowPattern{
			Type:       "if-else",
			StartLabel: label,
			EndLabel:   "", // No merge, both return
			CondBlock:  label,
			ThenBlock:  trueLabel,
			ElseBlock:  falseLabel,
			Blocks:     []string{label, trueLabel, falseLabel},
		}
	}
	
	return nil
}

// detectWhilePattern detects while loop patterns
func (cfa *ControlFlowAnalyzer) detectWhilePattern(label string, block *BasicBlock) *ControlFlowPattern {
	// While pattern: condition block with back-edge
	// while.cond -> (true) while.body -> while.cond
	//            -> (false) while.end
	
	if block.Terminator == nil || block.Terminator.Type != "br_cond" {
		return nil
	}
	
	// Check if this looks like a loop condition block
	if !strings.Contains(label, "while.cond") && !strings.Contains(label, "for.cond") && !strings.Contains(label, "loop") {
		return nil
	}
	
	bodyLabel := block.Terminator.TrueLabel
	endLabel := block.Terminator.FalseLabel
	
	bodyBlock, bodyExists := cfa.blocks[bodyLabel]
	if !bodyExists {
		return nil
	}
	
	// Check if body block jumps back to condition
	if bodyBlock.Terminator != nil &&
	   bodyBlock.Terminator.Type == "br_uncon" &&
	   bodyBlock.Terminator.UnconLabel == label {
		
		return &ControlFlowPattern{
			Type:       "while",
			StartLabel: label,
			EndLabel:   endLabel,
			CondBlock:  label,
			BodyBlock:  bodyLabel,
			Blocks:     []string{label, bodyLabel, endLabel},
		}
	}
	
	return nil
}

// detectSwitchPattern detects switch statement patterns
func (cfa *ControlFlowAnalyzer) detectSwitchPattern(label string, block *BasicBlock) *ControlFlowPattern {
	if block.Terminator == nil || block.Terminator.Type != "switch" {
		return nil
	}
	
	// Collect all case labels and default label
	blocks := []string{label}
	blocks = append(blocks, block.Terminator.DefaultLabel)
	
	for _, switchCase := range block.Terminator.SwitchCases {
		blocks = append(blocks, switchCase.Label)
	}
	
	// Create switch pattern
	return &ControlFlowPattern{
		Type:         "switch",
		StartLabel:   label,
		SwitchValue:  block.Terminator.SwitchValue,
		SwitchCases:  block.Terminator.SwitchCases,
		DefaultLabel: block.Terminator.DefaultLabel,
		Blocks:       blocks,
	}
}

// GetBasicBlocks returns the basic blocks
func (cfa *ControlFlowAnalyzer) GetBasicBlocks() map[string]*BasicBlock {
	return cfa.blocks
}
