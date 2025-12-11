package tests

import (
	"runtime"
	"strings"
	"testing"

	"linea/internal"
)

func TestDetectOS(t *testing.T) {
	os := internal.DetectOS()
	if os == "" {
		t.Error("DetectOS should return a non-empty string")
	}
	// Should match runtime.GOOS
	if os != runtime.GOOS {
		t.Errorf("Expected %s, got %s", runtime.GOOS, os)
	}
}

func TestNormalizePath(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"C:/Users/Test/file.txt", "C:\\Users\\Test\\file.txt"},
		{"path/to/file", "path\\to\\file"},
	}

	if runtime.GOOS == "windows" {
		for _, tt := range tests {
			result := internal.NormalizePath(tt.input)
			if result != tt.expected {
				t.Errorf("NormalizePath(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		}
	} else {
		// On Unix-like systems, paths should use forward slashes
		test := struct {
			input    string
			expected string
		}{"C:\\Users\\Test\\file.txt", "C:/Users/Test/file.txt"}
		result := internal.NormalizePath(test.input)
		if result != test.expected {
			t.Errorf("NormalizePath(%q) = %q, want %q", test.input, result, test.expected)
		}
	}
}

func TestSubstituteVariables(t *testing.T) {
	variables := map[string]string{
		"name": "John",
		"path": "/home/user",
	}

	tests := []struct {
		input    string
		expected string
	}{
		{"Hello {name}!", "Hello John!"},
		{"Path: {path}/file.txt", "Path: /home/user/file.txt"},
		{"No variables here", "No variables here"},
		{"{name} and {name}", "John and John"},
	}

	for _, tt := range tests {
		result := internal.SubstituteVariables(tt.input, variables)
		if result != tt.expected {
			t.Errorf("SubstituteVariables(%q) = %q, want %q", tt.input, result, tt.expected)
		}
	}
}

func TestSubstituteVariablesInArgs(t *testing.T) {
	variables := map[string]string{
		"file": "test.txt",
		"dir":  "C:/Users/Test",
	}

	args := []string{"-f", "{file}", "-d", "{dir}"}
	result := internal.SubstituteVariablesInArgs(args, variables)

	if len(result) != 4 {
		t.Fatalf("Expected 4 args, got %d", len(result))
	}

	if result[1] != "test.txt" {
		t.Errorf("Expected 'test.txt', got '%s'", result[1])
	}

	if result[3] != internal.NormalizePath("C:/Users/Test") {
		t.Errorf("Expected normalized path, got '%s'", result[3])
	}
}

func TestGetHelpFlag(t *testing.T) {
	flag := internal.GetHelpFlag()
	if runtime.GOOS == "windows" {
		if flag != "/?" {
			t.Errorf("Expected '/?' on Windows, got '%s'", flag)
		}
	} else {
		if flag != "--help" {
			t.Errorf("Expected '--help' on Unix, got '%s'", flag)
		}
	}
}

func TestIsPathLike(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
		desc     string
	}{
		{"/?", false, "Windows help flag should not be normalized"},
		{"/C", false, "Short Windows flag should not be normalized"},
		{"C:/Users/Test/file.txt", true, "Windows absolute path should be normalized"},
		{"./file.txt", true, "Relative path should be normalized"},
		{"../file.txt", true, "Parent relative path should be normalized"},
		{"/home/user/file", true, "Unix absolute path should be normalized"},
		{"path/to/file", true, "Path with multiple segments should be normalized"},
		{"--help", false, "Unix help flag should not be normalized"},
		{"-v", false, "Short flag should not be normalized"},
	}

	for _, tt := range tests {
		result := internal.IsPathLike(tt.input)
		if result != tt.expected {
			t.Errorf("IsPathLike(%q) = %v, want %v (%s)", tt.input, result, tt.expected, tt.desc)
		}
	}
}

func TestSubstituteVariablesInArgsPreservesFlags(t *testing.T) {
	args := []string{"/?", "-v", "--help", "C:/Users/file.txt"}
	result := internal.SubstituteVariablesInArgs(args, nil)

	// Flags should be preserved
	if result[0] != "/?" {
		t.Errorf("Expected '/?' to be preserved, got '%s'", result[0])
	}
	if result[1] != "-v" {
		t.Errorf("Expected '-v' to be preserved, got '%s'", result[1])
	}
	if result[2] != "--help" {
		t.Errorf("Expected '--help' to be preserved, got '%s'", result[2])
	}

	// Paths should be normalized
	if runtime.GOOS == "windows" {
		if !strings.Contains(result[3], "\\") {
			t.Errorf("Expected Windows path to be normalized with backslashes, got '%s'", result[3])
		}
	}
}

func TestSubstituteVariablesDollarSyntax(t *testing.T) {
	variables := map[string]string{
		"name":     "Hello",
		"variable": "World",
	}

	// Test $variable syntax
	result := internal.SubstituteVariables("{name}, $variable", variables)
	expected := "Hello, World"
	if result != expected {
		t.Errorf("Expected '%s', got '%s'", expected, result)
	}

	// Test $variable alone
	result2 := internal.SubstituteVariables("$variable", variables)
	if result2 != "World" {
		t.Errorf("Expected 'World', got '%s'", result2)
	}
}

