package cmd

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

func parseFlags(args []string) (paths []string, json bool, human bool) {
	for _, a := range args {
		switch a {
		case "--json":
			json = true
		case "--human":
			human = true
		default:
			if !strings.HasPrefix(a, "-") {
				paths = append(paths, a)
			}
		}
	}
	if len(paths) == 0 {
		paths = []string{"."}
	}
	return
}

func resolvePath(p string) string {
	if p == "." {
		cwd, err := os.Getwd()
		if err == nil {
			return cwd
		}
	}
	abs, err := os.Stat(p)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error accessing %s: %v\n", p, err)
		os.Exit(1)
	}
	if !abs.IsDir() {
		return p
	}
	return p
}

func comma(n int64) string {
	s := strconv.FormatInt(n, 10)
	if len(s) <= 3 {
		return s
	}
	var parts []string
	for i := len(s); i > 0; i -= 3 {
		start := i - 3
		if start < 0 {
			start = 0
		}
		parts = append([]string{s[start:i]}, parts...)
	}
	return strings.Join(parts, ",")
}
