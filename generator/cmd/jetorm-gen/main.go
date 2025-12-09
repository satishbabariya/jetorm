package main

import (
	"flag"
	"fmt"
	"go/parser"
	"go/token"
	"os"
	"strings"

	"github.com/satishbabariya/jetorm/generator"
)

func main() {
	var (
		typeName    = flag.String("type", "", "Entity type name (required)")
		output      = flag.String("output", "", "Output file path (required)")
		packageName = flag.String("package", "", "Package name for generated code (default: same as input)")
		inputFile   = flag.String("input", "", "Input Go source file (required)")
		interfaceName = flag.String("interface", "", "Repository interface name (optional)")
	)
	flag.Parse()

	if *typeName == "" {
		fmt.Fprintf(os.Stderr, "Error: -type is required\n")
		flag.Usage()
		os.Exit(1)
	}

	if *output == "" {
		fmt.Fprintf(os.Stderr, "Error: -output is required\n")
		flag.Usage()
		os.Exit(1)
	}

	if *inputFile == "" {
		fmt.Fprintf(os.Stderr, "Error: -input is required\n")
		flag.Usage()
		os.Exit(1)
	}

	// For now, we require the entity type to be passed as a string
	// In a full implementation, we'd parse the Go file and load the package
	// This is a simplified version that requires manual type specification
	if *typeName == "" {
		fmt.Fprintf(os.Stderr, "Error: -type is required\n")
		os.Exit(1)
	}

	// Note: In a production implementation, we'd use go/types to load the actual type
	// For now, this is a placeholder that shows the structure
	// The actual type would be obtained by loading the package
	fmt.Fprintf(os.Stderr, "Note: Full type loading not implemented. Using type name: %s\n", *typeName)
	
	// We'll generate code based on the interface methods instead
	// The entity type will be inferred from the interface

	// Get package name
	pkgName := *packageName
	if pkgName == "" {
		pkgName = extractPackageName(*inputFile)
	}

	// Parse interface to extract methods
	if *interfaceName == "" {
		fmt.Fprintf(os.Stderr, "Error: -interface is required\n")
		os.Exit(1)
	}

	parser := generator.NewParser()
	interfaceInfo, err := parser.ParseInterface(*inputFile, *interfaceName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing interface: %v\n", err)
		os.Exit(1)
	}

	if interfaceInfo == nil {
		fmt.Fprintf(os.Stderr, "Error: interface %s not found in %s\n", *interfaceName, *inputFile)
		os.Exit(1)
	}

	// Extract custom query methods
	customMethods := interfaceInfo.FindCustomMethods()
	if len(customMethods) == 0 {
		fmt.Fprintf(os.Stderr, "Warning: No custom query methods found in interface\n")
	}

	// For each custom method, we need to analyze it
	// Since we don't have the actual entity type loaded, we'll generate
	// code that can be compiled after the entity is available
	// This is a limitation we'll address with go/types in the future
	
	// Generate repository code
	code, err := generateRepositoryCode(pkgName, *typeName, customMethods)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error generating code: %v\n", err)
		os.Exit(1)
	}

	// Write to output file
	if err := os.WriteFile(*output, []byte(code), 0644); err != nil {
		fmt.Fprintf(os.Stderr, "Error writing output file: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Successfully generated repository code: %s\n", *output)
}


// extractPackageName extracts package name from a Go file
func extractPackageName(filePath string) string {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, filePath, nil, parser.ParseComments)
	if err != nil {
		return "main"
	}
	return f.Name.Name
}

// generateRepositoryCode generates the complete repository implementation
func generateRepositoryCode(pkgName, entityName string, customMethods []generator.MethodInfo) (string, error) {
	var buf strings.Builder

	// Write package declaration
	buf.WriteString(fmt.Sprintf("package %s\n\n", pkgName))

	// Write imports
	buf.WriteString(`import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/satishbabariya/jetorm/core"
)
`)

	// Write repository struct
	repoName := fmt.Sprintf("%sRepository", entityName)
	buf.WriteString(fmt.Sprintf(`
// %s is the generated repository implementation
type %s struct {
	*core.BaseRepository[%s, int64]
}

// New%s creates a new %s repository
func New%s(db *core.Database) (*%s, error) {
	baseRepo, err := core.NewBaseRepository[%s, int64](db)
	if err != nil {
		return nil, err
	}
	return &%s{
		BaseRepository: baseRepo,
	}, nil
}
`, repoName, repoName, entityName, repoName, repoName, repoName, repoName, entityName, repoName))

	// Generate custom query methods
	// Note: This is a simplified version that generates method stubs
	// In a full implementation, we'd use go/types to load the entity type
	// and generate complete implementations using the analyzer
	
	for _, methodInfo := range customMethods {
		if generator.IsQueryMethod(methodInfo.Name) {
			// Generate a method stub that will be implemented later
			// or use runtime analysis
			methodCode := generateMethodStub(methodInfo, entityName)
			buf.WriteString("\n")
			buf.WriteString(methodCode)
			buf.WriteString("\n")
		}
	}

	return buf.String(), nil
}

// generateMethodStub generates a method stub for a query method
func generateMethodStub(methodInfo generator.MethodInfo, entityName string) string {
	var buf strings.Builder
	
	// Build parameter list
	var params []string
	for _, param := range methodInfo.Parameters {
		if param.Name != "" {
			params = append(params, fmt.Sprintf("%s %s", param.Name, param.Type))
		} else {
			params = append(params, param.Type)
		}
	}
	paramsStr := strings.Join(params, ", ")
	
	// Build return list
	var returns []string
	for _, ret := range methodInfo.Returns {
		returns = append(returns, ret.Type)
	}
	returnsStr := strings.Join(returns, ", ")
	if len(returns) > 1 {
		returnsStr = "(" + returnsStr + ")"
	}
	
	// Generate method signature
	buf.WriteString(fmt.Sprintf("// %s implements the query method\n", methodInfo.Name))
	buf.WriteString(fmt.Sprintf("func (r *%sRepository) %s(ctx context.Context", entityName, methodInfo.Name))
	if paramsStr != "" {
		buf.WriteString(", " + paramsStr)
	}
	buf.WriteString(fmt.Sprintf(") %s {\n", returnsStr))
	buf.WriteString("\t// TODO: Implement query method\n")
	buf.WriteString("\t// This method should be generated using jetorm-gen with full type information\n")
	buf.WriteString("\tpanic(\"not implemented\")\n")
	buf.WriteString("}\n")
	
	return buf.String()
}

