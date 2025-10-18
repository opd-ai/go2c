package llvm

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
)

// Module represents an LLVM module
type Module struct {
	Functions []*Function
	Globals   []*Global
	Types     map[string]*Type
	Metadata  map[string]string
}

// Function represents an LLVM function
type Function struct {
	Name       string
	ReturnType string
	Parameters []Parameter
	Body       []string
	IsExternal bool
}

// Parameter represents a function parameter
type Parameter struct {
	Type string
	Name string
}

// Global represents a global variable
type Global struct {
	Name  string
	Type  string
	Value string
}

// Type represents an LLVM type definition
type Type struct {
	Name       string
	Definition string
}

// Parser handles LLVM IR parsing
type Parser struct {
	module *Module
}

// NewParser creates a new LLVM IR parser
func NewParser() *Parser {
	return &Parser{
		module: &Module{
			Functions: []*Function{},
			Globals:   []*Global{},
			Types:     make(map[string]*Type),
			Metadata:  make(map[string]string),
		},
	}
}

// Parse parses LLVM IR from a string
func (p *Parser) Parse(llvmIR string) (*Module, error) {
	scanner := bufio.NewScanner(strings.NewReader(llvmIR))
	
	var currentFunction *Function
	var functionBody []string
	inFunction := false

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		
		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, ";") {
			continue
		}

		// Parse function definitions
		if strings.HasPrefix(line, "define") {
			if currentFunction != nil && inFunction {
				currentFunction.Body = functionBody
				p.module.Functions = append(p.module.Functions, currentFunction)
			}
			
			fn, err := p.parseFunction(line)
			if err != nil {
				return nil, fmt.Errorf("failed to parse function: %w", err)
			}
			currentFunction = fn
			functionBody = []string{}
			inFunction = true
			continue
		}

		// Parse function declarations
		if strings.HasPrefix(line, "declare") {
			fn, err := p.parseFunctionDeclaration(line)
			if err != nil {
				return nil, fmt.Errorf("failed to parse function declaration: %w", err)
			}
			p.module.Functions = append(p.module.Functions, fn)
			continue
		}

		// Check for function end
		if inFunction && line == "}" {
			if currentFunction != nil {
				currentFunction.Body = functionBody
				p.module.Functions = append(p.module.Functions, currentFunction)
				currentFunction = nil
				functionBody = []string{}
			}
			inFunction = false
			continue
		}

		// Add lines to function body
		if inFunction {
			functionBody = append(functionBody, line)
		}

		// Parse global variables
		if strings.HasPrefix(line, "@") && strings.Contains(line, "=") {
			global := p.parseGlobal(line)
			if global != nil {
				p.module.Globals = append(p.module.Globals, global)
			}
		}

		// Parse type definitions
		if strings.Contains(line, "= type") {
			typ := p.parseType(line)
			if typ != nil {
				p.module.Types[typ.Name] = typ
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading LLVM IR: %w", err)
	}

	return p.module, nil
}

// ParseFile parses LLVM IR from a file
func (p *Parser) ParseFile(filePath string) (*Module, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}
	return p.Parse(string(data))
}

// parseFunction parses a function definition line
func (p *Parser) parseFunction(line string) (*Function, error) {
	// Example: define i32 @main() {
	// Example with aggregate: define {i32, i32} @foo(i32 %x) {
	
	// First extract the function name and parameters
	// Handle aggregate return types by matching everything before @
	re := regexp.MustCompile(`define\s+(.*?)\s+@([^\(]+)\((.*?)\)`)
	matches := re.FindStringSubmatch(line)
	
	if len(matches) < 3 {
		return nil, fmt.Errorf("invalid function definition: %s", line)
	}

	returnType := strings.TrimSpace(matches[1])
	
	fn := &Function{
		Name:       matches[2],
		ReturnType: returnType,
		Parameters: []Parameter{},
		IsExternal: false,
	}

	// Parse parameters if present
	if len(matches) > 3 && matches[3] != "" {
		params := p.parseParameters(matches[3])
		fn.Parameters = params
	}

	return fn, nil
}

// parseFunctionDeclaration parses a function declaration line
func (p *Parser) parseFunctionDeclaration(line string) (*Function, error) {
	// Example: declare void @llvm.memcpy.p0i8.p0i8.i32(i8*, i8*, i32, i1)
	// Handle aggregate return types
	re := regexp.MustCompile(`declare\s+(.*?)\s+@([^\(]+)\((.*?)\)`)
	matches := re.FindStringSubmatch(line)
	
	if len(matches) < 3 {
		return nil, fmt.Errorf("invalid function declaration: %s", line)
	}

	returnType := strings.TrimSpace(matches[1])
	
	fn := &Function{
		Name:       matches[2],
		ReturnType: returnType,
		Parameters: []Parameter{},
		IsExternal: true,
	}

	// Parse parameters if present
	if len(matches) > 3 && matches[3] != "" {
		params := p.parseParameters(matches[3])
		fn.Parameters = params
	}

	return fn, nil
}

// parseParameters parses function parameters
func (p *Parser) parseParameters(paramStr string) []Parameter {
	params := []Parameter{}
	if paramStr == "" {
		return params
	}

	parts := strings.Split(paramStr, ",")
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		
		// Split type and name
		tokens := strings.Fields(part)
		if len(tokens) >= 1 {
			param := Parameter{
				Type: tokens[0],
			}
			if len(tokens) >= 2 {
				param.Name = tokens[1]
			}
			params = append(params, param)
		}
	}

	return params
}

// parseGlobal parses a global variable definition
func (p *Parser) parseGlobal(line string) *Global {
	// Example: @.str = private unnamed_addr constant [14 x i8] c"Hello, World!\00"
	parts := strings.SplitN(line, "=", 2)
	if len(parts) < 2 {
		return nil
	}

	name := strings.TrimSpace(parts[0])
	rest := strings.TrimSpace(parts[1])

	return &Global{
		Name:  name,
		Type:  "",
		Value: rest,
	}
}

// parseType parses a type definition
func (p *Parser) parseType(line string) *Type {
	// Example: %struct.MyStruct = type { i32, i8* }
	parts := strings.SplitN(line, "=", 2)
	if len(parts) < 2 {
		return nil
	}

	name := strings.TrimSpace(parts[0])
	definition := strings.TrimSpace(parts[1])
	definition = strings.TrimPrefix(definition, "type ")

	return &Type{
		Name:       name,
		Definition: definition,
	}
}
