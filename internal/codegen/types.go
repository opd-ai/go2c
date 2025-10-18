package codegen

import (
	"fmt"
	"regexp"
	"strings"
)

// TypeMapper handles LLVM to C type conversions
type TypeMapper struct {
	customTypes map[string]string
}

// NewTypeMapper creates a new type mapper
func NewTypeMapper() *TypeMapper {
	return &TypeMapper{
		customTypes: make(map[string]string),
	}
}

// LLVMTypeToC converts an LLVM type to a C type
func (tm *TypeMapper) LLVMTypeToC(llvmType string) string {
	llvmType = strings.TrimSpace(llvmType)

	// Check for custom types first
	if cType, ok := tm.customTypes[llvmType]; ok {
		return cType
	}

	// Handle pointer types
	if strings.HasSuffix(llvmType, "*") {
		baseType := strings.TrimSuffix(llvmType, "*")
		baseType = strings.TrimSpace(baseType)
		return tm.LLVMTypeToC(baseType) + "*"
	}

	// Handle array types
	arrayRegex := regexp.MustCompile(`\[(\d+)\s+x\s+(.+)\]`)
	if matches := arrayRegex.FindStringSubmatch(llvmType); len(matches) == 3 {
		elementType := tm.LLVMTypeToC(matches[2])
		return fmt.Sprintf("%s[%s]", elementType, matches[1])
	}

	// Basic type mappings
	switch llvmType {
	case "void":
		return "void"
	case "i1":
		return "bool"
	case "i8":
		return "char"
	case "i16":
		return "short"
	case "i32":
		return "int"
	case "i64":
		return "long long"
	case "i128":
		return "__int128"
	case "float":
		return "float"
	case "double":
		return "double"
	case "i8*":
		return "char*"
	case "void*":
		return "void*"
	default:
		// Handle struct types
		if strings.HasPrefix(llvmType, "%struct.") {
			structName := strings.TrimPrefix(llvmType, "%struct.")
			return "struct " + structName
		}
		// Handle opaque types
		if strings.HasPrefix(llvmType, "%") {
			return "void*"
		}
		// Default: return as-is with a comment
		return fmt.Sprintf("/* %s */ void*", llvmType)
	}
}

// RegisterCustomType registers a custom type mapping
func (tm *TypeMapper) RegisterCustomType(llvmType, cType string) {
	tm.customTypes[llvmType] = cType
}

// GetCIntType returns the appropriate C integer type for a given bit width
func (tm *TypeMapper) GetCIntType(bitWidth int) string {
	switch bitWidth {
	case 1:
		return "bool"
	case 8:
		return "char"
	case 16:
		return "short"
	case 32:
		return "int"
	case 64:
		return "long long"
	case 128:
		return "__int128"
	default:
		return fmt.Sprintf("/* i%d */ int", bitWidth)
	}
}

// IsPointerType checks if an LLVM type is a pointer
func (tm *TypeMapper) IsPointerType(llvmType string) bool {
	return strings.Contains(llvmType, "*")
}

// IsIntegerType checks if an LLVM type is an integer
func (tm *TypeMapper) IsIntegerType(llvmType string) bool {
	return strings.HasPrefix(llvmType, "i") && !tm.IsPointerType(llvmType)
}

// IsFloatType checks if an LLVM type is a floating point type
func (tm *TypeMapper) IsFloatType(llvmType string) bool {
	return llvmType == "float" || llvmType == "double"
}

// SanitizeName sanitizes LLVM names for C
func (tm *TypeMapper) SanitizeName(name string) string {
	// Remove @ prefix for globals
	name = strings.TrimPrefix(name, "@")
	// Remove % prefix for locals
	name = strings.TrimPrefix(name, "%")
	// Remove quotes
	name = strings.Trim(name, "\"")
	// Replace dots with underscores
	name = strings.ReplaceAll(name, ".", "_")
	// Replace special characters
	name = strings.ReplaceAll(name, "-", "_")
	name = strings.ReplaceAll(name, "$", "_")

	// If name starts with a digit, prefix with underscore
	if len(name) > 0 && name[0] >= '0' && name[0] <= '9' {
		name = "var_" + name
	}

	// If name is empty or just whitespace, use a default
	if strings.TrimSpace(name) == "" {
		name = "unnamed"
	}

	return name
}
