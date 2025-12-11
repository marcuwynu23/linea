package tests

import (
	"path/filepath"
	"strings"
	"testing"

	"linea/internal"
)

func TestBuildCommand(t *testing.T) {
	config := &internal.CommandConfig{
		Command:    "docker",
		Subcommand: "ps",
		Args:       []string{"-a"},
	}

	cmd, err := internal.BuildCommand(config, nil)
	if err != nil {
		t.Fatalf("BuildCommand failed: %v", err)
	}
	expected := []string{"docker", "ps", "-a"}

	if len(cmd) != len(expected) {
		t.Fatalf("Expected %d args, got %d", len(expected), len(cmd))
	}

	for i, v := range expected {
		if cmd[i] != v {
			t.Errorf("Expected '%s' at index %d, got '%s'", v, i, cmd[i])
		}
	}
}

func TestBuildCommandWithoutSubcommand(t *testing.T) {
	config := &internal.CommandConfig{
		Command: "echo",
		Args:    []string{"Hello", "World"},
	}

	cmd, err := internal.BuildCommand(config, nil)
	if err != nil {
		t.Fatalf("BuildCommand failed: %v", err)
	}
	expected := []string{"echo", "Hello", "World"}

	if len(cmd) != len(expected) {
		t.Fatalf("Expected %d args, got %d", len(expected), len(cmd))
	}
}

func TestBuildCommandWithVariables(t *testing.T) {
	config := &internal.CommandConfig{
		Command: "echo",
		Args:    []string{"Hello {name}!"},
		Variables: map[string]string{
			"name": "Linea",
		},
	}

	cmd, err := internal.BuildCommand(config, nil)
	if err != nil {
		t.Fatalf("BuildCommand failed: %v", err)
	}
		if !strings.Contains(cmd[1], "Linea") {
		t.Errorf("Expected variable substitution, got %v", cmd)
	}
}

func TestFormatCommand(t *testing.T) {
	cmd := []string{"docker", "ps", "-a"}
	formatted := internal.FormatCommand(cmd)
	expected := "docker ps -a"

	if formatted != expected {
		t.Errorf("Expected '%s', got '%s'", expected, formatted)
	}
}

func TestExecuteCommandEcho(t *testing.T) {
	// Echo should work on all platforms:
	// - On Unix: echo is typically in PATH
	// - On Windows: echo is a shell built-in, handled via cmd.exe
	cmd := []string{"echo", "test"}
	err := internal.ExecuteCommand(cmd)
	if err != nil {
		t.Errorf("ExecuteCommand failed: %v", err)
	}
}

func TestDryRun(t *testing.T) {
	cmd := []string{"docker", "ps", "-a"}
	// This should not panic and should print something
	// We can't easily test stdout, but we can ensure it doesn't crash
	internal.DryRun(cmd)
}

func TestBuildCommandWithPathVariables(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.txt")

	config := &internal.CommandConfig{
		Command: "cat",
		Args:    []string{"{file_path}"},
		Variables: map[string]string{
			"file_path": testFile,
		},
	}

	cmd, err := internal.BuildCommand(config, nil)
	if err != nil {
		t.Fatalf("BuildCommand failed: %v", err)
	}
	// Path should be normalized
	if cmd[1] != internal.NormalizePath(testFile) {
		t.Errorf("Expected normalized path, got '%s'", cmd[1])
	}
}

