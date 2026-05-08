package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/hubermjonathan/cmux/internal/state"
)

func runStatus(args []string) error {
	isTmux := false
	session := ""

	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--tmux":
			isTmux = true
		case "--session":
			i++
			if i < len(args) {
				session = args[i]
			}
		default:
			if session == "" {
				session = args[i]
			}
		}
	}

	if session == "" {
		session = mostRecentSession()
	}
	if session == "" {
		return fmt.Errorf("no cmux session found")
	}
	if !strings.HasPrefix(session, sessionPrefix) {
		session = sessionPrefix + session
	}

	if isTmux {
		return printTmuxStatus(session)
	}
	return printTableStatus(session)
}

func printTmuxStatus(session string) error {
	counts := state.CountByStatus(session)
	var parts []string
	if n := counts[state.Working]; n > 0 {
		parts = append(parts, fmt.Sprintf("#[fg=green]●%d#[default]", n))
	}
	if n := counts[state.Waiting]; n > 0 {
		parts = append(parts, fmt.Sprintf("#[fg=yellow]◌%d#[default]", n))
	}
	if n := counts[state.Idle]; n > 0 {
		parts = append(parts, fmt.Sprintf("✓%d", n))
	}
	fmt.Printf("[cmux: %s]", strings.Join(parts, " "))
	return nil
}

func printTableStatus(session string) error {
	panes := state.GetAllPanes(session)
	if len(panes) == 0 {
		return fmt.Errorf("no panes found in session %s", session)
	}
	fmt.Printf("%-6s %-10s %-10s %s\n", "PANE", "STATUS", "DURATION", "LAST EVENT")
	for i, p := range panes {
		duration := "-"
		if p.LastEvent > 0 {
			duration = time.Since(time.Unix(p.LastEvent, 0)).Truncate(time.Second).String()
		}
		fmt.Printf("%-6d %-10s %-10s %s\n", i+1, p.Status, duration, eventName(p.Status))
	}
	return nil
}

func eventName(s state.Status) string {
	switch s {
	case state.Working:
		return "prompt-submit"
	case state.Waiting:
		return "stop"
	default:
		return "-"
	}
}
