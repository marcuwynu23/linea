package main

import (
	"fmt"
	"os"

	"linea/cmd"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	subcommand := os.Args[1]
	args := os.Args[2:]

	switch subcommand {
	case "run":
		cmd.RunCommandMain(args)
	case "test":
		cmd.TestCommandMain(args)
	case "help":
		cmd.HelpCommandMain(args)
	case "init":
		cmd.InitCommandMain(args)
	default:
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "  ❌ Error: unknown subcommand '%s'\n", subcommand)
		fmt.Fprintf(os.Stderr, "\n")
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Fprintf(os.Stderr, "\n")
	fmt.Fprintf(os.Stderr, "Linea - Commandline Workflow Tool\n")
	fmt.Fprintf(os.Stderr, "\n")
	fmt.Fprintf(os.Stderr, "  USAGE:\n")
	fmt.Fprintf(os.Stderr, "    linea <subcommand> [options] <yaml-file>\n")
	fmt.Fprintf(os.Stderr, "\n")
	fmt.Fprintf(os.Stderr, "  SUBCOMMANDS:\n")
	fmt.Fprintf(os.Stderr, "\n")
	fmt.Fprintf(os.Stderr, "    run    Execute the command defined in the YAML file\n")
	fmt.Fprintf(os.Stderr, "           \n")
	fmt.Fprintf(os.Stderr, "           Options:\n")
	fmt.Fprintf(os.Stderr, "             -v, --verbose              Show the command before executing\n")
	fmt.Fprintf(os.Stderr, "             --args <var>=<value>       Provide variable values (can be used multiple times)\n")
	fmt.Fprintf(os.Stderr, "\n")
	fmt.Fprintf(os.Stderr, "           Examples:\n")
	fmt.Fprintf(os.Stderr, "             linea run config.yml\n")
	fmt.Fprintf(os.Stderr, "             linea run -v config.yml\n")
	fmt.Fprintf(os.Stderr, "             linea run config.yml --args name=\"John\" --args age=30\n")
	fmt.Fprintf(os.Stderr, "\n")
	fmt.Fprintf(os.Stderr, "    test   Dry-run the command (print without executing)\n")
	fmt.Fprintf(os.Stderr, "           \n")
	fmt.Fprintf(os.Stderr, "           Options:\n")
	fmt.Fprintf(os.Stderr, "             --args <var>=<value>       Provide variable values for testing\n")
	fmt.Fprintf(os.Stderr, "\n")
	fmt.Fprintf(os.Stderr, "           Examples:\n")
	fmt.Fprintf(os.Stderr, "             linea test config.yml\n")
	fmt.Fprintf(os.Stderr, "             linea test config.yml --args variable=\"test\"\n")
	fmt.Fprintf(os.Stderr, "\n")
	fmt.Fprintf(os.Stderr, "    help   Display information about the command defined in YAML\n")
	fmt.Fprintf(os.Stderr, "           \n")
	fmt.Fprintf(os.Stderr, "           Examples:\n")
	fmt.Fprintf(os.Stderr, "             linea help config.yml\n")
	fmt.Fprintf(os.Stderr, "\n")
	fmt.Fprintf(os.Stderr, "    init   Initialize a new workflow YAML file with template\n")
	fmt.Fprintf(os.Stderr, "           \n")
	fmt.Fprintf(os.Stderr, "           Examples:\n")
	fmt.Fprintf(os.Stderr, "             linea init workflow.yml\n")
	fmt.Fprintf(os.Stderr, "             linea init my-commands.yml\n")
	fmt.Fprintf(os.Stderr, "\n")
	fmt.Fprintf(os.Stderr, "  For more information, visit: https://github.com/marcuwynu23/linea\n")
	fmt.Fprintf(os.Stderr, "\n")
}

