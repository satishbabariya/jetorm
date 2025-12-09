package generator

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// Config represents the configuration for code generation
type Config struct {
	// Entity configuration
	EntityType    string `json:"entity_type"`
	EntityPackage string `json:"entity_package,omitempty"`
	
	// Interface configuration
	InterfaceName string `json:"interface_name"`
	
	// Output configuration
	OutputFile    string `json:"output_file"`
	OutputPackage string `json:"output_package,omitempty"`
	
	// Input configuration
	InputFile string `json:"input_file"`
	
	// Generation options
	GenerateComments bool `json:"generate_comments,omitempty"`
	GenerateTests    bool `json:"generate_tests,omitempty"`
	
	// ID type (if not auto-detected)
	IDType string `json:"id_type,omitempty"`
}

// LoadConfig loads configuration from a file
func LoadConfig(configPath string) (*Config, error) {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}
	
	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}
	
	return &config, nil
}

// SaveConfig saves configuration to a file
func (c *Config) SaveConfig(configPath string) error {
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}
	
	if err := os.MkdirAll(filepath.Dir(configPath), 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}
	
	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}
	
	return nil
}

// Validate validates the configuration
func (c *Config) Validate() error {
	if c.EntityType == "" {
		return fmt.Errorf("entity_type is required")
	}
	if c.InterfaceName == "" {
		return fmt.Errorf("interface_name is required")
	}
	if c.OutputFile == "" {
		return fmt.Errorf("output_file is required")
	}
	if c.InputFile == "" {
		return fmt.Errorf("input_file is required")
	}
	return nil
}

// DefaultConfig returns a default configuration
func DefaultConfig() *Config {
	return &Config{
		GenerateComments: true,
		GenerateTests:    false,
		IDType:          "int64",
	}
}

