package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/javimosch/token-optimizer-cli/reporter"
	"github.com/javimosch/token-optimizer-cli/scanner"
)

func Scan(args []string) {
	paths, useJSON, useHuman := parseFlags(args)
	if !useJSON && !useHuman {
		useHuman = true
	}

	root := resolvePath(paths[0])
	result, err := scanner.Scan(root, useJSON)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Scan error: %v\n", err)
		os.Exit(1)
	}

	if useJSON {
		b, _ := json.MarshalIndent(result, "", "  ")
		fmt.Println(string(b))
		return
	}

	reporter.PrintScan(result)
}

func Audit(args []string) {
	paths, useJSON, useHuman := parseFlags(args)
	if !useJSON && !useHuman {
		useHuman = true
	}

	root := resolvePath(paths[0])
	result, err := scanner.Scan(root, useJSON)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Audit error: %v\n", err)
		os.Exit(1)
	}

	if useJSON {
		b, _ := json.MarshalIndent(result, "", "  ")
		fmt.Println(string(b))
		return
	}

	reporter.PrintAudit(result)
}

func CheckFile(args []string) {
	paths, useJSON, _ := parseFlags(args)
	if len(paths) < 1 || paths[0] == "." {
		fmt.Fprintln(os.Stderr, "Usage: token-optimizer-cli check <file> [--json]")
		os.Exit(1)
	}

	fi, err := scanner.CheckSingleFile(paths[0])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	if useJSON {
		b, _ := json.MarshalIndent(fi, "", "  ")
		fmt.Println(string(b))
		return
	}

	reporter.PrintCheck(fi)
}

func Summary(args []string) {
	paths, useJSON, useHuman := parseFlags(args)
	if !useJSON && !useHuman {
		useHuman = true
	}

	root := resolvePath(paths[0])
	result, err := scanner.Scan(root, useJSON)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Summary error: %v\n", err)
		os.Exit(1)
	}

	if useJSON {
		type summaryData struct {
			TotalFiles  int   `json:"total_files"`
			TotalTokens int   `json:"total_tokens"`
			TotalSize   int64 `json:"total_size_bytes"`
			BloatFiles  int   `json:"bloat_files"`
			BloatTokens int   `json:"bloat_tokens"`
		}
		s := summaryData{
			TotalFiles:  result.TotalFiles,
			TotalTokens: result.TotalTokens,
			TotalSize:   result.TotalSize,
			BloatFiles:  result.BloatFiles,
			BloatTokens: result.BloatTokens,
		}
		b, _ := json.MarshalIndent(s, "", "  ")
		fmt.Println(string(b))
		return
	}

	reporter.PrintSummary(result)
}
