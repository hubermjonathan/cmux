package main

import (
	"fmt"
	"strconv"

	"github.com/jon-huber/cmux/internal/state"
	"github.com/jon-huber/cmux/internal/tmux"
)

func runFocus(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("usage: cmux focus <pane> [--session <name>]")
	}

	paneIdx, err := strconv.Atoi(args[0])
	if err != nil || paneIdx < 1 {
		return fmt.Errorf("pane must be a positive integer (1-indexed)")
	}

	session := ""
	for i := 1; i < len(args); i++ {
		if args[i] == "--session" && i+1 < len(args) {
			session = "cmux-" + args[i+1]
			i++
		}
	}

	if session == "" {
		session = mostRecentSession()
	}
	if session == "" {
		return fmt.Errorf("no cmux session found")
	}

	panes := state.GetPaneList(session)
	if paneIdx > len(panes) {
		return fmt.Errorf("pane %d does not exist (session has %d panes)", paneIdx, len(panes))
	}

	paneID := panes[paneIdx-1]
	return tmux.RunSilent("select-pane", "-t", fmt.Sprintf("%%%d", paneID))
}
