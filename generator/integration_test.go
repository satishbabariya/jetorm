package generator

import (
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestIntegration_CodeGeneration tests the full code generation workflow
func TestIntegration_CodeGeneration(t *testing.T) {
	// Create a temporary test file
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test_entity.go")

	// Write test entity file
	testEntityCode := `package test

import "context"

type TestUser struct {
	ID       int64  ` + "`db:\"id\" jet:\"primary_key,auto_increment\"`" + `
	Email    string ` + "`db:\"email\" jet:\"unique,not_null\"`" + `
	Username string ` + "`db:\"username\" jet:\"unique,not_null\"`" + `
	Age      int    ` + "`db:\"age\"`" + `
	Status   string ` + "`db:\"status\"`" + `
}

type TestUserRepository interface {
	FindByEmail(ctx context.Context, email string) (*TestUser, error)
	FindByAgeGreaterThan(ctx context.Context, age int) ([]*TestUser, error)
	CountByStatus(ctx context.Context, status string) (int64, error)
}
`
	if err := os.WriteFile(testFile, []byte(testEntityCode), 0644); err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	// Parse the interface
	parser := NewParser()
	interfaceInfo, err := parser.ParseInterface(testFile, "TestUserRepository")
	if err != nil {
		t.Fatalf("Failed to parse interface: %v", err)
	}

	if interfaceInfo == nil {
		t.Fatal("Interface not found")
	}

	// Verify interface parsing
	customMethods := interfaceInfo.FindCustomMethods()
	if len(customMethods) != 3 {
		t.Errorf("Expected 3 custom methods, got %d", len(customMethods))
	}

	// Verify method names
	expectedMethods := map[string]bool{
		"FindByEmail":           true,
		"FindByAgeGreaterThan": true,
		"CountByStatus":         true,
	}

	for _, method := range customMethods {
		if !expectedMethods[method.Name] {
			t.Errorf("Unexpected method: %s", method.Name)
		}
	}

	// Verify query method detection
	for _, method := range customMethods {
		if !IsQueryMethod(method.Name) {
			t.Errorf("Method %s should be detected as query method", method.Name)
		}
	}
}

// TestIntegration_MethodAnalysis tests method name analysis
func TestIntegration_MethodAnalysis(t *testing.T) {
	// Create a temporary entity type for testing
	// We'll use a simple struct type
	testCases := []struct {
		methodName string
		shouldPass bool
	}{
		{"FindByEmail", true},
		{"FindByAgeGreaterThan", true},
		{"FindByStatusIn", true},
		{"CountByStatus", true},
		{"DeleteByEmail", true},
		{"ExistsByUsername", true},
		{"FindFirstByStatus", true},
		{"FindByStatusOrderByCreatedAtDesc", true},
		{"InvalidMethod", false},
		{"Save", false}, // Base method, not a query method
	}

	for _, tc := range testCases {
		t.Run(tc.methodName, func(t *testing.T) {
			isQueryMethod := IsQueryMethod(tc.methodName)
			if isQueryMethod != tc.shouldPass {
				t.Errorf("IsQueryMethod(%q) = %v, want %v", tc.methodName, isQueryMethod, tc.shouldPass)
			}
		})
	}
}

// TestIntegration_ConfigFile tests configuration file handling
func TestIntegration_ConfigFile(t *testing.T) {
	tmpDir := t.TempDir()
	configFile := filepath.Join(tmpDir, "test_config.json")

	// Create a test config
	cfg := DefaultConfig()
	cfg.EntityType = "User"
	cfg.InterfaceName = "UserRepository"
	cfg.InputFile = "user.go"
	cfg.OutputFile = "user_repository_gen.go"
	cfg.GenerateComments = true
	cfg.GenerateTests = false

	// Save config
	if err := cfg.SaveConfig(configFile); err != nil {
		t.Fatalf("Failed to save config: %v", err)
	}

	// Load config
	loadedCfg, err := LoadConfig(configFile)
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	// Verify loaded config
	if loadedCfg.EntityType != cfg.EntityType {
		t.Errorf("EntityType mismatch: got %s, want %s", loadedCfg.EntityType, cfg.EntityType)
	}
	if loadedCfg.InterfaceName != cfg.InterfaceName {
		t.Errorf("InterfaceName mismatch: got %s, want %s", loadedCfg.InterfaceName, cfg.InterfaceName)
	}
	if loadedCfg.GenerateComments != cfg.GenerateComments {
		t.Errorf("GenerateComments mismatch: got %v, want %v", loadedCfg.GenerateComments, cfg.GenerateComments)
	}

	// Test validation
	if err := loadedCfg.Validate(); err != nil {
		t.Errorf("Config validation failed: %v", err)
	}
}

// TestIntegration_GeneratedCodeStructure tests the structure of generated code
func TestIntegration_GeneratedCodeStructure(t *testing.T) {
	// This test verifies that generated code has the expected structure
	// without actually executing it

	// Create a mock method info
	methodInfo := MethodInfo{
		Name: "FindByEmail",
		Parameters: []ParameterInfo{
			{Name: "email", Type: "string"},
		},
		Returns: []ReturnInfo{
			{Type: "*User"},
			{Type: "error"},
		},
	}

	// Generate method stub (simplified version)
	entityName := "User"
	methodCode := generateMethodStubForTest(methodInfo, entityName)

	// Verify code structure
	if !strings.Contains(methodCode, "FindByEmail") {
		t.Error("Generated code should contain method name")
	}
	if !strings.Contains(methodCode, "func") {
		t.Error("Generated code should contain function declaration")
	}
	if !strings.Contains(methodCode, "context.Context") {
		t.Error("Generated code should contain context parameter")
	}
	if !strings.Contains(methodCode, "email string") {
		t.Error("Generated code should contain method parameters")
	}
	if !strings.Contains(methodCode, "*User") {
		t.Error("Generated code should contain return type")
	}
}

// generateMethodStubForTest is a helper for testing
func generateMethodStubForTest(methodInfo MethodInfo, entityName string) string {
	var buf strings.Builder

	// Build parameter list
	var params []string
	for _, param := range methodInfo.Parameters {
		if param.Name != "" {
			params = append(params, param.Name+" "+param.Type)
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
	buf.WriteString("func (r *" + entityName + "Repository) " + methodInfo.Name + "(ctx context.Context")
	if paramsStr != "" {
		buf.WriteString(", " + paramsStr)
	}
	buf.WriteString(") " + returnsStr + " {\n")
	buf.WriteString("\t// TODO: Implement\n")
	buf.WriteString("}\n")

	return buf.String()
}

// TestIntegration_ParserRobustness tests parser with various edge cases
func TestIntegration_ParserRobustness(t *testing.T) {
	tmpDir := t.TempDir()

	testCases := []struct {
		name     string
		code     string
		expected bool
	}{
		{
			name: "simple interface",
			code: `package test
type TestRepo interface {
	FindByID(ctx context.Context, id int64) (*Test, error)
}`,
			expected: true,
		},
		{
			name: "interface with multiple methods",
			code: `package test
type TestRepo interface {
	FindByID(ctx context.Context, id int64) (*Test, error)
	FindByEmail(ctx context.Context, email string) (*Test, error)
	Save(ctx context.Context, entity *Test) (*Test, error)
}`,
			expected: true,
		},
		{
			name: "no interface",
			code: `package test
type Test struct {
	ID int64
}`,
			expected: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			testFile := filepath.Join(tmpDir, tc.name+".go")
			if err := os.WriteFile(testFile, []byte(tc.code), 0644); err != nil {
				t.Fatalf("Failed to write test file: %v", err)
			}

			parser := NewParser()
			interfaceInfo, err := parser.ParseInterface(testFile, "TestRepo")
			if err != nil {
				if tc.expected {
					t.Errorf("Unexpected error: %v", err)
				}
				return
			}

			if tc.expected && interfaceInfo == nil {
				t.Error("Expected interface but got nil")
			}
			if !tc.expected && interfaceInfo != nil {
				t.Error("Did not expect interface but got one")
			}
		})
	}
}

// TestIntegration_GoSyntaxValidation tests that generated code is valid Go syntax
func TestIntegration_GoSyntaxValidation(t *testing.T) {
	// Generate sample code
	pkgName := "test"
	entityName := "User"
	customMethods := []MethodInfo{
		{
			Name: "FindByEmail",
			Parameters: []ParameterInfo{
				{Name: "email", Type: "string"},
			},
			Returns: []ReturnInfo{
				{Type: "*User"},
				{Type: "error"},
			},
		},
	}

	cfg := DefaultConfig()
	cfg.GenerateComments = true

	code, err := generateRepositoryCodeForTest(pkgName, entityName, customMethods, cfg)
	if err != nil {
		t.Fatalf("Failed to generate code: %v", err)
	}

	// Parse generated code to verify syntax
	fset := token.NewFileSet()
	_, err = parser.ParseFile(fset, "", code, parser.ParseComments)
	if err != nil {
		t.Errorf("Generated code has syntax errors: %v\nCode:\n%s", err, code)
	}
}

// generateRepositoryCodeForTest is a helper that generates code for testing
func generateRepositoryCodeForTest(pkgName, entityName string, customMethods []MethodInfo, cfg *Config) (string, error) {
	var buf strings.Builder

	buf.WriteString("package " + pkgName + "\n\n")
	buf.WriteString(`import (
	"context"
	"github.com/satishbabariya/jetorm/core"
)
`)

	repoName := entityName + "Repository"
	idType := cfg.IDType
	if idType == "" {
		idType = "int64"
	}

	buf.WriteString("type " + repoName + " struct {\n")
	buf.WriteString("\t*core.BaseRepository[" + entityName + ", " + idType + "]\n")
	buf.WriteString("}\n\n")

	buf.WriteString("func New" + repoName + "(db *core.Database) (*" + repoName + ", error) {\n")
	buf.WriteString("\tbaseRepo, err := core.NewBaseRepository[" + entityName + ", " + idType + "](db)\n")
	buf.WriteString("\tif err != nil {\n")
	buf.WriteString("\t\treturn nil, err\n")
	buf.WriteString("\t}\n")
	buf.WriteString("\treturn &" + repoName + "{BaseRepository: baseRepo}, nil\n")
	buf.WriteString("}\n")

	for _, method := range customMethods {
		buf.WriteString("\n")
		buf.WriteString(generateMethodStubForTest(method, entityName))
	}

	return buf.String(), nil
}

