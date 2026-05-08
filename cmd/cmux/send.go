package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/hubermjonathan/cmux/internal/state"
	"github.com/hubermjonathan/cmux/internal/tmux"
)

func runSend(args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("usage: cmux send <pane> <text> [--session <name>]")
	}

	paneIdx, err := strconv.Atoi(args[0])
	if err != nil || paneIdx < 1 {
		return fmt.Errorf("pane must be a positive integer (1-indexed)")
	}

	session := ""
	var textParts []string
	for i := 1; i < len(args); i++ {
		if args[i] == "--session" && i+1 < len(args) {
			session = "cmux-" + args[i+1]
			i++
		} else {
			textParts = append(textParts, args[i])
		}
	}
	text := strings.Join(textParts, " ")

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
	target := fmt.Sprintf("%%%d", paneID)
	tmux.Run("send-keys", "-t", target, "-l", text)
	tmux.Run("send-keys", "-t", target, "Enter")
	return nil
}
