package main

import (
	_ "embed"
	"fmt"
	"strings"
)

//go:embed docs/cmux-docs.md
var docsContent string

func runDocs(args []string) error {
	if len(args) == 0 {
		fmt.Print(docsContent)
		return nil
	}

	// Filter to specific command section
	cmd := args[0]
	lines := strings.Split(docsContent, "\n")
	var output []string
	inSection := false
	for _, line := range lines {
		if strings.HasPrefix(line, "### cmux "+cmd) {
			inSection = true
		} else if inSection && strings.HasPrefix(line, "### cmux ") {
			break
		} else if inSection && strings.HasPrefix(line, "## ") {
			break
		}
		if inSection {
			output = append(output, line)
		}
	}

	if len(output) == 0 {
		return fmt.Errorf("no docs found for command %q", cmd)
	}
	fmt.Println(strings.Join(output, "\n"))
	return nil
}
