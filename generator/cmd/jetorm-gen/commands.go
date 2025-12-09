package main

import (
	"fmt"
)

// Command represents a CLI command
type Command struct {
	Name        string
	Description string
	Execute     func([]string) error
}

// Available commands
var commands = []Command{
	{
		Name:        "init",
		Description: "Create a new configuration file",
		Execute:     cmdInit,
	},
	{
		Name:        "generate",
		Description: "Generate repository code",
		Execute:     cmdGenerate,
	},
	{
		Name:        "validate",
		Description: "Validate configuration",
		Execute:     cmdValidate,
	},
}

// cmdInit creates a configuration file
func cmdInit(args []string) error {
	configPath := "jetorm-gen.json"
	if len(args) > 0 {
		configPath = args[0]
	}

	return initConfigFile(configPath)
}

// cmdGenerate generates code
func cmdGenerate(args []string) error {
	cfg, err := parseConfig()
	if err != nil {
		return fmt.Errorf("failed to parse config: %w", err)
	}

	// Generate code (implementation from main.go)
	// This is a simplified version
	fmt.Printf("Generating code for %s...\n", cfg.EntityType)
	return nil
}

// cmdValidate validates configuration
func cmdValidate(args []string) error {
	cfg, err := parseConfig()
	if err != nil {
		return fmt.Errorf("invalid configuration: %w", err)
	}

	if err := cfg.Validate(); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	fmt.Println("Configuration is valid")
	return nil
}

// printUsage prints command usage
func printUsage() {
	fmt.Println("Usage: jetorm-gen [command] [options]")
	fmt.Println("\nCommands:")
	for _, cmd := range commands {
		fmt.Printf("  %-15s %s\n", cmd.Name, cmd.Description)
	}
	fmt.Println("\nOptions:")
	fmt.Println("  -config string    Configuration file path")
	fmt.Println("  -type string      Entity type name")
	fmt.Println("  -interface string  Repository interface name")
	fmt.Println("  -input string      Input Go source file")
	fmt.Println("  -output string     Output file path")
	fmt.Println("  -comments          Generate documentation comments")
	fmt.Println("  -tests             Generate test files")
}

// executeCommand executes a command
func executeCommand(name string, args []string) error {
	for _, cmd := range commands {
		if cmd.Name == name {
			return cmd.Execute(args)
		}
	}
	return fmt.Errorf("unknown command: %s", name)
}

