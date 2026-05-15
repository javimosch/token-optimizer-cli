package main

import (
	"fmt"
	"os"

	"github.com/javimosch/token-optimizer-cli/cmd"
)

const version = "1.0.0"

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	subcommand := os.Args[1]
	args := os.Args[2:]

	switch subcommand {
	case "scan":
		cmd.Scan(args)
	case "audit":
		cmd.Audit(args)
	case "check":
		cmd.CheckFile(args)
	case "summary":
		cmd.Summary(args)
	case "version", "--version", "-v":
		fmt.Println(version)
	case "help", "--help", "-h":
		printUsage()
	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", subcommand)
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println(`token-optimizer-cli - Token optimization for AI coding agents

Usage:
  token-optimizer-cli scan [path]       Scan files for token usage
  token-optimizer-cli audit [path]      Deep audit with bloat detection
  token-optimizer-cli check <file>      Single file analysis
  token-optimizer-cli summary [path]    Project-level summary
  token-optimizer-cli version           Print version

Path defaults to current directory. Use --json for JSON output.`)
}
