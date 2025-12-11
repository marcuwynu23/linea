package internal

import (
	"fmt"
	"io"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

// ParseYAML reads and parses a YAML file into a CommandConfig
func ParseYAML(filePath string) (*CommandConfig, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %s: %w", filePath, err)
	}

	var config CommandConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse YAML: %w", err)
	}

	if config.Command == "" {
		return nil, fmt.Errorf("command field is required")
	}

	return &config, nil
}

// ParseMultiYAML reads and parses a YAML file with multiple documents (separated by ---)
// Returns a slice of CommandConfig, one for each document
func ParseMultiYAML(filePath string) ([]*CommandConfig, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %s: %w", filePath, err)
	}

	// Always use decoder to handle both single and multiple documents
	decoder := yaml.NewDecoder(strings.NewReader(string(data)))
	var configs []*CommandConfig

	for {
		var config CommandConfig
		err := decoder.Decode(&config)
		if err != nil {
			// Check if it's EOF (end of documents)
			if err == io.EOF {
				break
			}
			return nil, fmt.Errorf("failed to parse YAML document: %w", err)
		}

		if config.Command == "" {
			// Skip empty documents
			continue
		}

		configs = append(configs, &config)
	}

	if len(configs) == 0 {
		return nil, fmt.Errorf("no valid commands found in YAML file")
	}

	return configs, nil
}

