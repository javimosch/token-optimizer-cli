package reporter

import (
	"fmt"

	"github.com/javimosch/token-optimizer-cli/scanner"
)

func PrintScan(r *scanner.ScanResult) {
	fmt.Println("Token Scan Results")
	fmt.Println("==================")
	fmt.Printf("Files scanned:  %d\n", r.TotalFiles)
	fmt.Printf("Total tokens:  %d (~$%.2f @ GPT-4o rates)\n", r.TotalTokens, float64(r.TotalTokens)*0.00001)
	fmt.Printf("Total size:    %d bytes\n", r.TotalSize)
	fmt.Printf("Bloat files:   %d (%d tokens)\n", r.BloatFiles, r.BloatTokens)
	fmt.Println()
	fmt.Println("Top files by token count:")
	fmt.Println("-------------------------")

	type fileEntry struct {
		path  string
		tokens int
		lang  string
	}
	var sorted []fileEntry
	for _, f := range r.Files {
		sorted = append(sorted, fileEntry{f.Path, f.Tokens.TokenEstimate, f.Language})
	}
	for i := 0; i < len(sorted); i++ {
		for j := i + 1; j < len(sorted); j++ {
			if sorted[j].tokens > sorted[i].tokens {
				sorted[i], sorted[j] = sorted[j], sorted[i]
			}
		}
	}

	limit := 20
	if len(sorted) < limit {
		limit = len(sorted)
	}
	for i := 0; i < limit; i++ {
		marker := ""
		if sorted[i].tokens > 1000 {
			marker = " ⚠️"
		}
		fmt.Printf("  %5d tokens  %-15s %s%s\n", sorted[i].tokens, sorted[i].lang, sorted[i].path, marker)
	}

	if r.BloatFiles > 0 {
		fmt.Println()
		fmt.Printf("Tip: %d files exceed 1000 tokens. Run 'audit' for optimization suggestions.\n", r.BloatFiles)
	}
}

func PrintAudit(r *scanner.ScanResult) {
	fmt.Println("Token Audit Report")
	fmt.Println("==================")
	fmt.Printf("Total files: %d\n", r.TotalFiles)
	fmt.Printf("Total tokens: %d\n", r.TotalTokens)
	fmt.Printf("Estimated cost: $%.4f (GPT-4o input)\n", float64(r.TotalTokens)*0.00001)
	fmt.Println()

	bloatCount := 0
	for _, f := range r.Files {
		if len(f.Issues) > 0 {
			bloatCount++
			fmt.Printf("📄 %s (%s)\n", f.Path, f.Language)
			for _, issue := range f.Issues {
				sev := "[i]"
				if issue.Severity == "warning" {
					sev = "[w]"
				}
				fmt.Printf("  %s %s\n", sev, issue.Message)
				if issue.Suggestion != "" {
					fmt.Printf("     → %s\n", issue.Suggestion)
				}
			}
			fmt.Println()
		}
	}

	if bloatCount == 0 {
		fmt.Println("No optimization issues found. Clean project!")
	} else {
		fmt.Printf("Summary: %d/%d files have optimization opportunities (%d total tokens).\n",
			bloatCount, r.TotalFiles, r.TotalTokens)
	}
}

func PrintCheck(fi *scanner.FileInfo) {
	fmt.Printf("File: %s\n", fi.Path)
	fmt.Printf("Size: %d bytes\n", fi.Size)
	fmt.Printf("Language: %s\n", fi.Language)
	fmt.Printf("Lines: %d\n", fi.Tokens.Lines)
	fmt.Printf("Chars: %d\n", fi.Tokens.Chars)
	fmt.Println()
	fmt.Println("Token estimates:")
	fmt.Printf("  GPT-4:     %d\n", fi.Tokens.TokensGPT4)
	fmt.Printf("  GPT-4o:    %d\n", fi.Tokens.TokensGPT4o)
	fmt.Printf("  Claude-3:  %d\n", fi.Tokens.TokensClaude)
	fmt.Printf("  Estimate:  %d\n", fi.Tokens.TokenEstimate)
	fmt.Println()

	if len(fi.Issues) > 0 {
		fmt.Println("Issues:")
		for _, issue := range fi.Issues {
			fmt.Printf("  [%s] %s\n", issue.Severity, issue.Message)
			if issue.Suggestion != "" {
				fmt.Printf("        → %s\n", issue.Suggestion)
			}
		}
	} else {
		fmt.Println("No issues detected.")
	}
}

func PrintSummary(r *scanner.ScanResult) {
	fmt.Println("Project Token Summary")
	fmt.Println("=====================")
	fmt.Printf("Directory:     .\n")
	fmt.Printf("Files:         %d\n", r.TotalFiles)
	fmt.Printf("Total tokens:  %d\n", r.TotalTokens)
	fmt.Printf("Total size:    %d bytes (%d KB)\n", r.TotalSize, r.TotalSize/1024)
	fmt.Println()
	fmt.Printf("Token distribution:\n")
	fmt.Printf("  Small files (<500 tok):   %d\n", countByThreshold(r, 0, 500))
	fmt.Printf("  Medium (500-1000 tok):    %d\n", countByThreshold(r, 500, 1000))
	fmt.Printf("  Large (1000-5000 tok):    %d\n", countByThreshold(r, 1000, 5000))
	fmt.Printf("  XL (>5000 tok):           %d\n", countByThreshold(r, 5000, -1))
	fmt.Println()
	fmt.Printf("Bloat: %d files (%d tokens, %.1f%% of total)\n",
		r.BloatFiles, r.BloatTokens, pct(r.BloatTokens, r.TotalTokens))
	fmt.Println()
	fmt.Printf("Rough cost at GPT-4o: $%.4f per full context load\n", float64(r.TotalTokens)*0.00001)
}

func countByThreshold(r *scanner.ScanResult, min, max int) int {
	n := 0
	for _, f := range r.Files {
		t := f.Tokens.TokenEstimate
		if t >= min && (max < 0 || t < max) {
			n++
		}
	}
	return n
}

func pct(a, b int) float64 {
	if b == 0 {
		return 0
	}
	return float64(a) / float64(b) * 100
}
