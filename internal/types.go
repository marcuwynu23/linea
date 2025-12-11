package internal

// CommandConfig represents the structure of a YAML command file
type CommandConfig struct {
	Command    string            `yaml:"command"`
	Subcommand string            `yaml:"subcommand,omitempty"`
	Args       []string          `yaml:"args,omitempty"`
	Variables  map[string]string `yaml:"variables,omitempty"`
}

