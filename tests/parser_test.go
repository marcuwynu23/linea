package tests

import (
	"os"
	"path/filepath"
	"testing"

	"linea/internal"
)

func TestParseYAML(t *testing.T) {
	// Create a temporary YAML file
	tmpFile := filepath.Join(t.TempDir(), "test.yml")
	yamlContent := `command: echo
args:
  - "Hello, World"
variables:
  name: "Test"
`
	if err := os.WriteFile(tmpFile, []byte(yamlContent), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	config, err := internal.ParseYAML(tmpFile)
	if err != nil {
		t.Fatalf("Failed to parse YAML: %v", err)
	}

	if config.Command != "echo" {
		t.Errorf("Expected command 'echo', got '%s'", config.Command)
	}

	if len(config.Args) != 1 || config.Args[0] != "Hello, World" {
		t.Errorf("Expected args ['Hello, World'], got %v", config.Args)
	}

	if config.Variables["name"] != "Test" {
		t.Errorf("Expected variable 'name' to be 'Test', got '%s'", config.Variables["name"])
	}
}

func TestParseYAMLWithSubcommand(t *testing.T) {
	tmpFile := filepath.Join(t.TempDir(), "test.yml")
	yamlContent := `command: docker
subcommand: ps
args:
  - -a
`
	if err := os.WriteFile(tmpFile, []byte(yamlContent), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	config, err := internal.ParseYAML(tmpFile)
	if err != nil {
		t.Fatalf("Failed to parse YAML: %v", err)
	}

	if config.Command != "docker" {
		t.Errorf("Expected command 'docker', got '%s'", config.Command)
	}

	if config.Subcommand != "ps" {
		t.Errorf("Expected subcommand 'ps', got '%s'", config.Subcommand)
	}
}

func TestParseYAMLMissingCommand(t *testing.T) {
	tmpFile := filepath.Join(t.TempDir(), "test.yml")
	yamlContent := `args:
  - "test"
`
	if err := os.WriteFile(tmpFile, []byte(yamlContent), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	_, err := internal.ParseYAML(tmpFile)
	if err == nil {
		t.Error("Expected error for missing command field")
	}
}

