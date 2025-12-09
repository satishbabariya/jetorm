package main

import (
	"flag"
	"fmt"

	"github.com/satishbabariya/jetorm/generator"
)

// parseConfig parses configuration from command line flags and config file
func parseConfig() (*generator.Config, error) {
	var (
		configFile   = flag.String("config", "", "Path to configuration file (JSON)")
		typeName     = flag.String("type", "", "Entity type name")
		output       = flag.String("output", "", "Output file path")
		packageName  = flag.String("package", "", "Package name for generated code")
		inputFile    = flag.String("input", "", "Input Go source file")
		interfaceName = flag.String("interface", "", "Repository interface name")
		generateComments = flag.Bool("comments", true, "Generate documentation comments")
		generateTests = flag.Bool("tests", false, "Generate test files")
	)
	flag.Parse()

	var cfg *generator.Config

	// Load from config file if provided
	if *configFile != "" {
		var err error
		cfg, err = generator.LoadConfig(*configFile)
		if err != nil {
			return nil, fmt.Errorf("failed to load config: %w", err)
		}
	} else {
		cfg = generator.DefaultConfig()
	}

	// Override with command line flags
	if *typeName != "" {
		cfg.EntityType = *typeName
	}
	if *output != "" {
		cfg.OutputFile = *output
	}
	if *packageName != "" {
		cfg.OutputPackage = *packageName
	}
	if *inputFile != "" {
		cfg.InputFile = *inputFile
	}
	if *interfaceName != "" {
		cfg.InterfaceName = *interfaceName
	}
	if flag.NFlag() > 0 {
		cfg.GenerateComments = *generateComments
		cfg.GenerateTests = *generateTests
	}

	// Validate configuration
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return cfg, nil
}

// initConfigFile creates an initial configuration file
func initConfigFile(configPath string) error {
	cfg := generator.DefaultConfig()
	cfg.EntityType = "User"
	cfg.InterfaceName = "UserRepository"
	cfg.InputFile = "user.go"
	cfg.OutputFile = "user_repository_gen.go"
	
	return cfg.SaveConfig(configPath)
}

