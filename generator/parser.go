package generator

import (
	"go/ast"
	"go/parser"
	"go/token"
	"reflect"
	"strings"
)

// InterfaceInfo represents information about a repository interface
type InterfaceInfo struct {
	Name       string
	EntityType reflect.Type
	Methods    []MethodInfo
}

// MethodInfo represents information about a method in an interface
type MethodInfo struct {
	Name       string
	Parameters []ParameterInfo
	Returns    []ReturnInfo
}

// ParameterInfo represents a method parameter
type ParameterInfo struct {
	Name string
	Type string
}

// ReturnInfo represents a method return value
type ReturnInfo struct {
	Type string
}

// Parser parses Go source files to extract interface definitions
type Parser struct {
	fset *token.FileSet
}

// NewParser creates a new parser
func NewParser() *Parser {
	return &Parser{
		fset: token.NewFileSet(),
	}
}

// ParseInterface parses a Go source file and extracts interface information
func (p *Parser) ParseInterface(filePath string, interfaceName string) (*InterfaceInfo, error) {
	// Parse the Go source file
	f, err := parser.ParseFile(p.fset, filePath, nil, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	// Find the interface declaration
	var interfaceInfo *InterfaceInfo
	ast.Inspect(f, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.GenDecl:
			if x.Tok == token.TYPE {
				for _, spec := range x.Specs {
					if ts, ok := spec.(*ast.TypeSpec); ok {
						if ts.Name.Name == interfaceName {
							if it, ok := ts.Type.(*ast.InterfaceType); ok {
								interfaceInfo = p.extractInterface(ts.Name.Name, it)
								return false // Stop inspection
							}
						}
					}
				}
			}
		}
		return true
	})

	if interfaceInfo == nil {
		return nil, nil // Interface not found
	}

	return interfaceInfo, nil
}

// extractInterface extracts interface information from AST
func (p *Parser) extractInterface(name string, it *ast.InterfaceType) *InterfaceInfo {
	info := &InterfaceInfo{
		Name:    name,
		Methods: make([]MethodInfo, 0),
	}

	for _, method := range it.Methods.List {
		if fn, ok := method.Type.(*ast.FuncType); ok {
			methodInfo := MethodInfo{
				Name:       method.Names[0].Name,
				Parameters: p.extractParameters(fn.Params),
				Returns:    p.extractReturns(fn.Results),
			}
			info.Methods = append(info.Methods, methodInfo)
		}
	}

	return info
}

// extractParameters extracts parameter information
func (p *Parser) extractParameters(params *ast.FieldList) []ParameterInfo {
	if params == nil {
		return nil
	}

	var parameters []ParameterInfo
	for _, param := range params.List {
		typeStr := p.typeToString(param.Type)
		if len(param.Names) > 0 {
			for _, name := range param.Names {
				parameters = append(parameters, ParameterInfo{
					Name: name.Name,
					Type: typeStr,
				})
			}
		} else {
			parameters = append(parameters, ParameterInfo{
				Name: "",
				Type: typeStr,
			})
		}
	}

	return parameters
}

// extractReturns extracts return type information
func (p *Parser) extractReturns(results *ast.FieldList) []ReturnInfo {
	if results == nil {
		return nil
	}

	var returns []ReturnInfo
	for _, result := range results.List {
		typeStr := p.typeToString(result.Type)
		returns = append(returns, ReturnInfo{
			Type: typeStr,
		})
	}

	return returns
}

// typeToString converts an AST type to a string representation
func (p *Parser) typeToString(expr ast.Expr) string {
	switch x := expr.(type) {
	case *ast.Ident:
		return x.Name
	case *ast.SelectorExpr:
		return p.typeToString(x.X) + "." + x.Sel.Name
	case *ast.ArrayType:
		return "[]" + p.typeToString(x.Elt)
	case *ast.StarExpr:
		return "*" + p.typeToString(x.X)
	case *ast.MapType:
		return "map[" + p.typeToString(x.Key) + "]" + p.typeToString(x.Value)
	case *ast.ChanType:
		return "chan " + p.typeToString(x.Value)
	case *ast.InterfaceType:
		return "interface{}"
	case *ast.FuncType:
		return "func(...)" // Simplified
	case *ast.Ellipsis:
		return "..." + p.typeToString(x.Elt)
	default:
		return "unknown"
	}
}

// FindCustomMethods finds methods in an interface that are not part of the base Repository interface
func (info *InterfaceInfo) FindCustomMethods() []MethodInfo {
	baseMethods := map[string]bool{
		"Save":           true,
		"SaveAll":        true,
		"Update":         true,
		"UpdateAll":      true,
		"FindByID":       true,
		"FindAll":        true,
		"FindAllByIDs":   true,
		"Delete":         true,
		"DeleteByID":     true,
		"DeleteAll":      true,
		"DeleteAllByIDs": true,
		"Count":          true,
		"ExistsById":     true,
		"FindAllPaged":   true,
		"SaveBatch":      true,
		"WithTx":         true,
		"Query":          true,
		"QueryOne":       true,
		"Exec":           true,
		"FindOne":        true,
		"FindAllWithSpec": true,
		"FindAllPagedWithSpec": true,
		"CountWithSpec":  true,
		"ExistsWithSpec": true,
		"DeleteWithSpec": true,
	}

	var customMethods []MethodInfo
	for _, method := range info.Methods {
		if !baseMethods[method.Name] {
			customMethods = append(customMethods, method)
		}
	}

	return customMethods
}

// IsQueryMethod checks if a method name follows the query method naming convention
func IsQueryMethod(methodName string) bool {
	queryPrefixes := []string{
		"FindBy", "FindFirstBy", "FindTop",
		"CountBy", "CountDistinctBy",
		"ExistsBy",
		"DeleteBy",
		"FindDistinctBy",
	}

	for _, prefix := range queryPrefixes {
		if strings.HasPrefix(methodName, prefix) {
			return true
		}
	}

	return false
}

