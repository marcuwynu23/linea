package internal

import (
	"fmt"
	"path/filepath"
	"runtime"
	"strings"
)

// DetectOS returns the current operating system
func DetectOS() string {
	return runtime.GOOS
}

// NormalizePath converts a path to use the correct path separators for the current OS
func NormalizePath(path string) string {
	if runtime.GOOS == "windows" {
		// Convert forward slashes to backslashes on Windows
		path = strings.ReplaceAll(path, "/", "\\")
	} else {
		// Convert backslashes to forward slashes on Unix-like systems
		path = strings.ReplaceAll(path, "\\", "/")
	}
	return filepath.Clean(path)
}

// SubstituteVariables replaces {variable} and $variable placeholders in strings with their values
func SubstituteVariables(s string, variables map[string]string) string {
	result := s
	for key, value := range variables {
		// Support {variable} syntax
		placeholder := "{" + key + "}"
		result = strings.ReplaceAll(result, placeholder, value)
		
		// Support $variable syntax
		dollarPlaceholder := "$" + key
		result = strings.ReplaceAll(result, dollarPlaceholder, value)
	}
	return result
}

// IsPathLike checks if a string looks like a file path rather than a flag or option
func IsPathLike(s string) bool {
	// Exclude common Windows flags that start with / or \
	// Examples: /?, /C, /D, \?, etc.
	if len(s) <= 3 && (strings.HasPrefix(s, "/") || strings.HasPrefix(s, "\\")) {
		// Single character flags like /?, /C, \? are not paths
		return false
	}
	
	// Check for Windows drive letters (C:, D:, etc.)
	if len(s) >= 2 && s[1] == ':' && ((s[0] >= 'A' && s[0] <= 'Z') || (s[0] >= 'a' && s[0] <= 'z')) {
		return true
	}
	
	// Check for relative paths starting with ./ or ../
	if strings.HasPrefix(s, "./") || strings.HasPrefix(s, "../") {
		return true
	}
	
	// Check for Unix absolute paths (but exclude short flags)
	if strings.HasPrefix(s, "/") && len(s) > 3 {
		// If it has multiple path segments, it's likely a path
		if strings.Count(s, "/") > 1 {
			return true
		}
		// Single segment but looks like a path (has extension or is longer)
		if strings.Contains(s, ".") || len(s) > 4 {
			return true
		}
	}
	
	// Check for Windows paths with backslashes
	if strings.Contains(s, "\\") {
		return true
	}
	
	// If it contains both forward and backslashes, it's likely a path
	if strings.Contains(s, "/") && strings.Contains(s, "\\") {
		return true
	}
	
	// Check for multiple path segments
	pathSeparators := strings.Count(s, "/") + strings.Count(s, "\\")
	if pathSeparators > 1 {
		return true
	}
	
	return false
}

// ExtractVariableReferences extracts all variable references from a string
// Returns a set of variable names (both {variable} and $variable syntax)
func ExtractVariableReferences(s string) map[string]bool {
	refs := make(map[string]bool)
	
	// Extract {variable} references
	start := -1
	for i, char := range s {
		if char == '{' {
			start = i
		} else if char == '}' && start != -1 {
			varName := s[start+1 : i]
			if varName != "" {
				refs[varName] = true
			}
			start = -1
		}
	}
	
	// Extract $variable references
	// Look for $ followed by alphanumeric characters or underscore
	for i := 0; i < len(s); i++ {
		if s[i] == '$' && i+1 < len(s) {
			// Check if next character is valid for variable name
			if (s[i+1] >= 'a' && s[i+1] <= 'z') || 
			   (s[i+1] >= 'A' && s[i+1] <= 'Z') || 
			   s[i+1] == '_' {
				// Extract variable name
				j := i + 1
				for j < len(s) && ((s[j] >= 'a' && s[j] <= 'z') || 
					(s[j] >= 'A' && s[j] <= 'Z') || 
					(s[j] >= '0' && s[j] <= '9') || 
					s[j] == '_') {
					j++
				}
				varName := s[i+1 : j]
				if varName != "" {
					refs[varName] = true
				}
				i = j - 1
			}
		}
	}
	
	return refs
}

// ValidateVariables checks if all referenced variables are defined
// Returns an error listing missing variables if any
func ValidateVariables(args []string, variables map[string]string) error {
	allRefs := make(map[string]bool)
	
	// Extract all variable references from all arguments
	for _, arg := range args {
		refs := ExtractVariableReferences(arg)
		for ref := range refs {
			allRefs[ref] = true
		}
	}
	
	// Check which variables are missing
	missing := []string{}
	for ref := range allRefs {
		if _, exists := variables[ref]; !exists {
			missing = append(missing, ref)
		}
	}
	
	if len(missing) > 0 {
		return fmt.Errorf("undefined variables: %s (use -s or --set to provide values)", strings.Join(missing, ", "))
	}
	
	return nil
}

// SubstituteVariablesInArgs applies variable substitution to all arguments
func SubstituteVariablesInArgs(args []string, variables map[string]string) []string {
	result := make([]string, len(args))
	for i, arg := range args {
		result[i] = SubstituteVariables(arg, variables)
		// Only normalize paths, not flags or options
		if IsPathLike(result[i]) {
			result[i] = NormalizePath(result[i])
		}
	}
	return result
}

// SubstituteVariablesInArgsWithSeparateMaps applies variable substitution with separate maps
// {name} uses yamlVars only, $name uses dollarVars
func SubstituteVariablesInArgsWithSeparateMaps(args []string, yamlVars map[string]string, dollarVars map[string]string) []string {
	result := make([]string, len(args))
	for i, arg := range args {
		result[i] = SubstituteVariablesWithSeparateMaps(arg, yamlVars, dollarVars)
		// Only normalize paths, not flags or options
		if IsPathLike(result[i]) {
			result[i] = NormalizePath(result[i])
		}
	}
	return result
}

// SubstituteVariablesWithSeparateMaps substitutes variables with separate maps
// {name} uses yamlVars only (not overridable), $name uses dollarVars (overridable)
func SubstituteVariablesWithSeparateMaps(s string, yamlVars map[string]string, dollarVars map[string]string) string {
	result := s
	
	// First substitute {variable} using ONLY YAML variables (not overridable)
	for key, value := range yamlVars {
		placeholder := "{" + key + "}"
		result = strings.ReplaceAll(result, placeholder, value)
	}
	
	// Then substitute $variable using dollarVars (overridable, includes -s/--set)
	for key, value := range dollarVars {
		// Replace ${VAR} first (more specific)
		placeholder2 := "${" + key + "}"
		result = strings.ReplaceAll(result, placeholder2, value)
		
		// Replace $VAR (but not if it's part of a longer variable name)
		placeholder1 := "$" + key
		for {
			idx := strings.Index(result, placeholder1)
			if idx == -1 {
				break
			}
			
			// Check if it's a valid variable reference (not part of a longer variable)
			afterIdx := idx + len(placeholder1)
			if afterIdx >= len(result) {
				// End of string, valid replacement
				result = result[:idx] + value + result[afterIdx:]
			} else {
				nextChar := result[afterIdx]
				// Valid if next char is not alphanumeric or underscore
				if !((nextChar >= 'a' && nextChar <= 'z') || 
					 (nextChar >= 'A' && nextChar <= 'Z') || 
					 (nextChar >= '0' && nextChar <= '9') || 
					 nextChar == '_') {
					result = result[:idx] + value + result[afterIdx:]
				} else {
					// Skip this occurrence, it's part of a longer variable
					result = result[:idx+1] + result[idx+1:]
				}
			}
		}
	}
	
	return result
}

// GetHelpFlag returns the appropriate help flag for the current OS
func GetHelpFlag() string {
	if runtime.GOOS == "windows" {
		return "/?"
	}
	return "--help"
}

